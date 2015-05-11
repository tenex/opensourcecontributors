db.contributions.find().snapshot().forEach(
    function (e) {
        if (typeof e.created_at === 'string') {
            e.created_at = ISODate(e.created_at);
            db.contributions.save(e);
        }
    }
)
