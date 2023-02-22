const topicDatabaseName = 'topics';
const topicDatabase = db.getSiblingDB(topicDatabaseName);

const topicCollectionName = 'topics';
const topicCollection = topicDatabase.getCollection(topicCollectionName);
const contentCollectionName = 'content';
const contentCollection = topicDatabase.getCollection(topicCollectionName);

const apiUrl = 'http://localhost:25300/topics/';

const cfgDefault = {
  verbose: true, // display the new documents
  insert: true, // set to false to avoid inserts
};
