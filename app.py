from flask import Flask, jsonify, request
from service import *
import MercurySQLite

db = MercurySQLite.DataBase('test.db')

app = Flask(__name__)

@app.route('/')
def index():
    return 'Hello World!'

@app.route('/api/v1.0/mercury', methods=['GET'])
def ping_heartbeat():
    return jsonify({'status': 'OK'})

@app.route('/api/v1.0/mercury/<string:>', methods=['GET'])
def ping_heartbeat():
    return jsonify({'status': 'OK'})

if __name__ == '__main__':
    app.run(debug=True)
