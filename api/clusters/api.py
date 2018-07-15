from flask import Blueprint, current_app, request, jsonify
from .clusters import Clusters


clusters = Blueprint('clusters', __name__)
cluster = Clusters()


@clusters.route('/clusters/create', methods=['POST'])
def create():
    fulfillment_text = request.json['queryResult']['fulfillmentText']
    current_app.logger.info("Creating a cluster for fulfillment text: '{}'".format(fulfillment_text))

    parameters = request.json['queryResult']['parameters']
    return jsonify(cluster.create(
        name=parameters.get('name'),
        machine_type=parameters.get('machineType'),
        location=parameters.get('location'),
        kubernetes_version=parameters.get('kubernetesVersion'),
        nodes=parameters.get('nodes')
    ))

@clusters.route('/clusters/list', methods=['GET'])
def list():
    return jsonify(cluster.list())
