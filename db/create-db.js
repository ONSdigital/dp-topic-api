load('./db-config.js');

function createCollections() {
  topicDatabase.createCollection(topicCollectionName);
  topicCollection.createIndex({ id: 1 }, { name: 'topics_id' });

  topicDatabase.createCollection(contentCollectionName);
  contentCollection.createIndex({ id: 1 }, { name: 'topics_content_id' });
}

function createRootTopic() {
  topicCollection.insert({
    id: 'topic_root',
    current: {
      id: 'topic_root',
      state: 'published',
      subtopics_ids: [],
    },
    next: {
      id: 'topic_root',
      state: 'published',
      subtopics_ids: [],
    },
  });
}
