from flask import Flask
from flask.ext.pymongo import PyMongo
app = Flask(__name__)
app.config['MONGO_DBNAME'] = 'contributions'
mongo = PyMongo(app)
