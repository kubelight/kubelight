from flask import Flask, redirect, url_for, request
from .api.clusters.api import clusters

app = Flask(__name__)
app.register_blueprint(clusters)

mappings = {
    'startCluster.startCluster-custom': '/clusters/create',
}


@app.route('/', methods=['GET', 'POST'])
def redirection():
    action = request.json['queryResult']['action']
    route = mappings.get(action)
    if route is not None:
        return redirect(route, code=302)
    app.logger.info("Invalid action in query: {}".format(action))
    return "Invalid action in query: {}".format(action)


if __name__ == '__main__':
    app.run()
