# Standing up dp-topic-api database

This folder contains a few scripts used for modifying the `topics` database in mongodb.

Needs to have mongodb 3.6+ installed and running - this can be done via the dp-compose repository.

## Starting with blank database

Run:

`mongo <mongo_url> init-db.js <options>`

This will:

- create the database and it's collections
- create a root topic for the topics to hang off
- insert Census topic and it's subtopics

The `<mongo_url>` part, if supplied, should look like:

- `<host>:<port>` for example: `localhost:27017` (this is the default)
- If authentication is needed, use the format `mongodb://<username>:<password>@<host>:<port>`

Example of the (optional) `<options>` part:

- `--eval 'cfg={verbose:true}'` (e.g. use for debugging)
- `cfg` defaults to: `{verbose:false, insert: true}` (see db-config.js)
- if you specify `cfg`, all missing options default to `false`

## Wiping your database

`mongo wipe-db.js`
