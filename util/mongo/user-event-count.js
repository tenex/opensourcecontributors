db.contributions.aggregate(
    [
        {
            "$match": {
                "_event_id": { "$exists" : true },
            },
        },
        {
            "$group": {
                "_id": { "_event_id": "$_event_id" },
                "uniqueIds": { "$addToSet": "$_id" },
                "count": { "$sum": 1 },
            },
        },
        {
            "$match": {
                "count": { "$gt": 1 },
            }
        },
        {
            "$sort": { "count": -1 },
        },
        {
            "$limit": 3,
        }
    ],
    {
        "allowDiskUse": true,
    }).forEach(function(doc) {
        printjson(doc);
    });
