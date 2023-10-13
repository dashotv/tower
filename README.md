# Tower

If tower fucks up seer, use this to clean up dates in `mongosh`:

> db.media.update({_id: ObjectId("5d02f69f6b696d1074020000")}, {$unset: {"paths.$[].created_at":1}})
