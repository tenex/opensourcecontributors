use contributions
db.contributions.createIndex( { "_org_lower": 1 });
db.contributions.createIndex( { "_repo_lower": 1 });
db.contributions.createIndex( { "_user_lower": 1 });
db.contributions.createIndex( { "created_at": 1 });
db.contributions.createIndex( { "_event_id": 1 }, { unique: true }); // TODO use sparse; not unique with nulls
