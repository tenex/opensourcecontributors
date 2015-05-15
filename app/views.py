import time
import math
from app import app, mongo
from flask import render_template
from flask.ext.pymongo import ASCENDING, DESCENDING
from app.tools import jsonify

PAGE_SIZE = 50

@app.route('/')
def index():
    return app.send_static_file('index.html')

@app.route('/search/<query>')
@app.route('/search/<query>/<int:page>')
def search(query, page=1):
    collection = mongo.db.contributions
    criteria = {
        '_user_lower': query.lower(),
        # 'type': 'CommitCommentEvent'
    }
    total = collection.find(criteria).count()
    skip = (page-1) * PAGE_SIZE
    total_pages = math.ceil(float(total) / PAGE_SIZE)
    repos = collection.find(criteria).distinct('repo')

    events = collection.find(criteria)
    events = events.sort("created_at", DESCENDING)
    events = events.skip(skip).limit(PAGE_SIZE)
    events = list(events)
    events = {
        "repos": repos,
        "events": events,
        "user": query,
        "total": total,
        "pageCount": total_pages,
        "start": skip+1,
        "end": skip+len(events),
        "currentPage": page,
        "size": len(events)
    }
    return jsonify(**events)
