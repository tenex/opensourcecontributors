# Fabrice Aneche (@akhenakh)
# https://github.com/akhenakh
try:
    import simplejson as json
except ImportError:
    try:
        import json
    except ImportError:
        raise ImportError
import datetime
from bson.objectid import ObjectId
from werkzeug import Response

class MongoJsonEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, (datetime.datetime, datetime.date)):
            return obj.isoformat()
        elif isinstance(obj, ObjectId):
            return str(obj)
        return json.JSONEncoder.default(self, obj)

def jsonify(*args, **kwargs):
    """ jsonify with support for MongoDB ObjectId
    """
    obj = json.dumps(dict(*args, **kwargs),
                     cls=MongoJsonEncoder,
                     sort_keys=True,
                     indent=4)
    return Response(obj, mimetype='application/json')
