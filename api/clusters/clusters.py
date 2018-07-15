import config
from googleapiclient import discovery
from google.oauth2 import service_account


class Clusters:
    def __init__(self):
        self.credentials = service_account.Credentials.from_service_account_file(config.SERVICE_ACCOUNT_FILE)
        self.gke = discovery.build('container', 'v1', credentials=self.credentials)
        self.gke_clusters = self.gke.projects().locations().clusters()
        self.project_id = self.credentials.project_id

    def create(self, name, nodes, kubernetes_version, machine_type, location):
        return self.gke_clusters.create(
            parent='projects/{}/locations/{}'.format(self.project_id, location),
            body={
                'cluster': {
                    'name': name,
                    'initialNodeCount': nodes,
                    'initialClusterVersion': kubernetes_version,
                    # 'config': {
                    #     'machineType': machine_type,
                    # }
                },
            },
        ).execute()

    def list(self):
        return self.gke_clusters.list(
            parent='projects/{}/locations/-'.format(self.project_id)
        ).execute()
