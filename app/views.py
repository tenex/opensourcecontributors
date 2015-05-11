from app import app, mongo
from flask import render_template
from flask.ext.pymongo import ASCENDING, DESCENDING
from app.tools import jsonify

@app.route('/')
def index():
    return app.send_static_file('index.html')

@app.route('/search/<query>')
def search(query):
    collection = mongo.db.contributions
    criteria = {'_user_lower': query.lower()}
    events = collection.find(criteria).sort('created_at', DESCENDING)
    events = { 'events': list(events) }
    return jsonify(**events)
