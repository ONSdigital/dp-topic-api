load('./db-config.js');

function wipe() {
  topicCollection.remove({});
  contentCollection.remove({});
  topicDatabase.dropDatabase();
}

wipe();
