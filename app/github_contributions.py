from flask import Flask, render_template
from flask.ext.pymongo import PyMongo, ASCENDING, DESCENDING
from tools import jsonify
import time
import math

app = Flask(__name__)
app.config.from_object('config')
mongo = PyMongo(app)

PAGE_SIZE = 50

@app.route('/')
def index():
    return app.send_static_file('index.html')

@app.route('/stats')
def stats():
    summary = {
    }
    return jsonify(**summary)

@app.route('/user/<username>')
def user(username):
    collection = mongo.db.contributions
    criteria = {
        '_user_lower': username.lower(),
    }
    repos = collection.find(criteria)
    repos = repos.distinct('repo')
    repos.sort(key=str.lower)

    event_count = collection.find(criteria).count()

    summary = {
        "username": username,
        "eventCount": event_count,
        "repos": repos,
    }
    return jsonify(**summary)

@app.route('/user/<username>/events')
@app.route('/user/<username>/events/<int:page>')
def events(username, page=1):
    collection = mongo.db.contributions
    criteria = {
        '_user_lower': username.lower(),
    }

    skip = (page-1) * PAGE_SIZE
    #total_pages = math.ceil(float(total) / PAGE_SIZE)

    events = collection.find(criteria)
    events = events.sort("created_at", DESCENDING)
    events = events.skip(skip).limit(PAGE_SIZE)
    events = list(events)
    events = {
        "events": events,
        "start": skip+1,
        "end": skip+len(events),
        "currentPage": page,
        "size": len(events)
    }
    return jsonify(**events)
