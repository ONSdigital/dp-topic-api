insert-census-topics
====================

This utility inserts the Census topic and all its subtopics.

### How to run the utility

Run
```
mongo <mongo_url> <options> insert-census-topics.js
```

The `<mongo_url>` part should look like:
- `<host>:<port>/topics` for example: `localhost:27017/topics`
- If authentication is needed, use the format `mongodb://<username>:<password>@<host>:<port>/topics`
- in the above, `/topics` indicates the database to be modified

Example of the (optional) `<options>` part:

- `--eval 'cfg={verbose:true}'` (e.g. use for debugging)
  - `cfg` defaults to: `{verbose:false, insert: true}`
  - if you specify `cfg`, all missing options default to `false`

It is recommended to perform a dry run and check the result looks as expected:

```
mongo localhost:27017/topics --eval 'cfg={verbose:true, insert:false}' insert-census-topics.js
```

Note: when connecting to a TLS-enabled DocumentDB cluster (sandbox or prod), you'll need to add the following options:
- `--tls`
- `--tlsCAFile=<pem>` where `<pem>` is the path to the Certificate Authority .pem file

For example:
```
mongo mongodb://$MONGO_USER:$MONGO_PASS@$MONGO_HOST/topics --tls --tlsCAFile=./cert.pem insert-census-topics.js
```
