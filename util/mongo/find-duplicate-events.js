// This shouldn't ever need to be run unless the existing data needs re-importing
// Takes 40 minutes to clear 10 records :(
db.contributions.aggregate(
  [
    {
      "$match": {
        "_event_id": { "$exists" : true }
      },
    },
    {
      "$group": {
        "_id": { "_event_id": "$_event_id" },
        "uniqueIds": { "$push": "$_id" },
        "count": { "$sum": 1 }
      },
    },
    {
      "$match": {
        "count": { "$gt": 1 }
      }
    },
    {
      $out : "duplicates"
    }
  ],
  {
    "allowDiskUse": true
  }
);
// .forEach(function(doc) {
//     doc.uniqueIds.shift();
//     printjson(doc);
//     var wRes = db.contributions.remove(
//         {
//             "_id": {
//                 "$in": doc.uniqueIds,
//             },
//         },
//         {
//             "justOne": true,
//         });
//     print(wRes);
// });
