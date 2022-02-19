package v1alpha1

import (
	"bytes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"context"
	"fmt"
	"github.com/linde12/gowol"
	"golang.org/x/crypto/ssh"
	"knative.dev/pkg/logging"
	"net"
	"time"
)

// TODO
type BackendType interface {
}

// +k8s:deepcopy-gen=false
type WOLProxy struct {
	ProxyClient *ssh.Client
	WakeOnLAN   WakeOnLAN
}

func (b *Backend) Init(ctx context.Context, client kubernetes.Interface) error {
	logger := logging.FromContext(ctx)
	if b.Spec.Type == BackendTypeBareMetal {
		if len(b.Spec.Machines) < 1 {
			return fmt.Errorf("no machines specified for the backend %v", b.ObjectMeta.Name)
		}

		var wolProxy WOLProxy
		// let's connect to the WOL proxy first
		if b.Spec.WakeOnLAN.Proxy != "" {
			logger.Infof("connecting to proxy for WOL: %v", b.Spec.WakeOnLAN.Proxy)
			wolProxy.WakeOnLAN = b.Spec.WakeOnLAN
			var err error
			proxyClient, err := getSSHClientFromSecret(ctx, client, b.Spec.WakeOnLAN.SSHSecret, b.ObjectMeta.Namespace, b.Spec.WakeOnLAN.Proxy)
			if err != nil {
				return err
			}
			wolProxy.ProxyClient = proxyClient
		}

		// TODO: make this concurrent
		for _, machine := range b.Spec.Machines {
			if err := machine.Init(ctx, client, wolProxy, b.ObjectMeta.Namespace); err != nil {
				b.Status.MarkBackendUnavailable()
				return err
			}
		}
	}
	b.Status.MarkBackendAvailable()
	return nil
}

func (m *Machine) Init(ctx context.Context, client kubernetes.Interface, wolProxy WOLProxy, namespace string) error {
	logger := logging.FromContext(ctx)
	logger.Infof("initializing machine: %v", m.Host)
	// check if Host is reachable on port 22
	conn, err := net.DialTimeout("tcp", m.Host+":22", 5*time.Second)
	if err != nil {
		logger.Errorf("unable to reach host: %v", m.Host)

		// see if WOL is enabled so we can wake the host up
		if len(m.WakeOnLan) > 0 {
			if err := wakeHost(ctx, wolProxy, m.Host, m.WakeOnLan); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	// conn might not be nil since we are not returning all the errors above due to WOL
	if conn != nil {
		defer conn.Close()
	}

	// now that the host is reachable, let's see if we can connect over SSH

	// let's start by getting SSH credentials from the secret
	if _, err := getSSHClientFromSecret(ctx, client, m.SSHSecret, namespace, m.Host); err != nil {
		return err
	}
	logger.Infof("successfully established SSH connection to host: %v", m.Host)
	return nil
}

func getSSHClientFromSecret(ctx context.Context, client kubernetes.Interface, secretName, namespace, host string) (*ssh.Client, error) {
	logger := logging.FromContext(ctx)
	secret, err := client.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	sshUser := secret.Data["user"]
	sshPassword := secret.Data["password"]

	sshConfig := ssh.ClientConfig{
		User: string(sshUser),
		Auth: []ssh.AuthMethod{
			ssh.Password(string(sshPassword)),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	logger.Infof("trying to establish SSH connection to %v...", host)
	sshClient, err := ssh.Dial("tcp", host+":22", &sshConfig)
	if err != nil {
		return nil, err
	}
	return sshClient, nil
}

func wakeHost(ctx context.Context, wolProxy WOLProxy, host string, mac string) error {
	logger := logging.FromContext(ctx)
	time.Sleep(5 * time.Second)

	logger.Infof("trying to wake host: %v via proxy %v", host, wolProxy.WakeOnLAN.Proxy)
	proxySession, err := wolProxy.ProxyClient.NewSession()
	if err != nil {
		return err
	}
	defer proxySession.Close()

	logger.Infof("running wake-on-lan command on proxy %v for host %v", wolProxy.WakeOnLAN.Proxy, host)
	var b bytes.Buffer
	proxySession.Stdout = &b
	if err := proxySession.Run(wolProxy.WakeOnLAN.Command + " " + mac); err != nil {
		return err
	}
	logger.Info(b.String())

	packet, err := gowol.NewMagicPacket(mac)
	if err != nil {
		return err
	}
	if err := packet.Send("255.255.255.255"); err != nil {
		return err
	}

	// check if host is reachable now
	conn, err := net.DialTimeout("tcp", host+":22", 10*time.Second)
	if err != nil {
		// so it's not reachable yet, let's try again
		logger.Infof("host %v is still not up, trying again ...", host)
		return wakeHost(ctx, wolProxy, host, mac)
	}
	defer conn.Close()
	return nil
}
