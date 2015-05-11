from app import app, mongo
from flask import render_template, jsonify

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/search/<query>')
def search(query):
    events = mongo.db.contributions.find({
        '_user_lower' : query.lower()
    })
    return jsonify(events)
