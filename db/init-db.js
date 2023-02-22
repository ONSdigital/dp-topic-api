load('./create-db.js');
load('./insert_topics.js');
load('./insert-census-topics.js');

function init() {
  createCollections();
  createRootTopic();
  createCensusTopics();
}

init();
