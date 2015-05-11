from flask import Flask
from flask.ext.pymongo import PyMongo

app = Flask('githubcontributionarchive')
app.config.from_object('config')
mongo = PyMongo(app)
from app import views
