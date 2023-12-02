from flask import Flask, jsonify, request, g
from service import *
import MercurySQLite

app = Flask(__name__)

@app.route('/')
def index():
    return 'Hello World!'

@app.route('/api/v1', methods=['GET'])
def ping_heartbeat():
    # return jsonify({'status': 'OK'})
    return jsonify({'status': 'OK', 'code': 200, 'data': None})

def get_db():
    if 'db' not in g:
        g.db = MercurySQLite.DataBase('test.db')
        init_db(g.db)
    return g.db

@app.teardown_appcontext
def close_db(error):
    if 'db' in g:
        del g.db
        # g.db.close()

@app.route('/api/v1/heartbeat', methods=['GET'])
def heartbeat():
    db = get_db()
    result = do_heartbeat(db)
    result = {'status': 'OK', 'code': 200, 'data': None}
    return jsonify(result)

# 没有匹配的路由
@app.errorhandler(404)
def page_not_found(error):
    return jsonify({'status': 'ERROR', 'code': 404, 'message': 'Page not found'}), 404

if __name__ == '__main__':
    app.run(debug=True)
