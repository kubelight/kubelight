# Kubelight - ephemeral Kubernetes cluster pools

---

Kubelight maintains a pool of ready-to-use ephemeral Kubernetes or OpenShift clusters.

### Pool

A pool defines the desired count and configuration of the clusters.

### Backend

A backend defines the resources that are available to a pool. A pool can have multiple backends.

##### OpenShift cluster pool being run on CRC on a bare metal.

```yaml
apiVersion: kubelight.dev/v1alpha1
kind: Pool
metadata:
  name: openshift-dev-pool
spec:
  backends:
    - name: m93p-machines
      type: crc
  count: 3
  flavor:
    name: openshift
    version: 4.9.x
  config:
    nodes: 1
    resources:
      nodes: 1
      cpu: 2
      memory: 16G
      disk: 50G
  strategy: aggressive
```

```yaml
apiVersion: kubelight.dev/v1alpha1
kind: Backend
metadata:
  name: m93p-machines
spec:
  type: bare-metal
  machines:
    - host: server1.concaf
      credentials: server1-secret
      wake-on-lan: true
    - host: server2.concaf
      credentials: server2-secret
      wake-on-lan: true
    - host: server3.concaf
      credentials: server3-secret
      wake-on-lan: true
```

##### Kubernetes pool being run on GKE and AWS

```yaml
apiVersion: kubelight.dev/v1alpha1
kind: Pool
metadata:
  name: kubernetes-dev-pool
spec:
  backends:
    - name: gke-backend
      type: gke
    - name: aws-backend
      type: aws
  count: 3
  flavor:
    name: openshift
    version: 4.9.x
  config:
    nodes: 1
    resources:
      nodes: 1
      cpu: 2
      memory: 16G
      disk: 50G
  strategy: aggressive
```

```yaml
apiVersion: kubelight.dev/v1alpha1
kind: Backend
metadata:
  name: gke-backend
spec:
  type: gke
  credentials: <secret>
  zone: us-central-1c
```

```yaml
apiVersion: kubelight.dev/v1alpha1
kind: Backend
metadata:
  name: aws-backend
spec:
  type: aws
  credentials: <secret>
```

### CLI reference

```console
kubelight claim <pool> <claim>
kubelight release <claim>

kubelight claims list
kubelight pools list
kubelight backends list
```
