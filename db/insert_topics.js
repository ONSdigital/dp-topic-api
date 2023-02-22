const idSize = 4;
const rootId = 'topic_root';
const idAlphabet = '123456789';

if (typeof (cfg) == 'undefined') {
  // default, but can be changed on command-line, see README
  cfg = cfgDefault;
}

function deepCopy(obj) {
  return JSON.parse(JSON.stringify(obj));
}

function isUsedId(id) {
  return topicCollection.find({ id }).hasNext();
}

/**
 * Generates random id string based on idAlphabet and
 * idSize
 * @returns id string
 */
function makeId() {
  let result = '';
  for (let i = 0; i < idSize; i++) {
    result += idAlphabet.charAt(Math.floor(Math.random() * idAlphabet.length));
  }
  return result;
}

function generateUnusedId() {
  let id = '';
  do {
    id = makeId();
  } while (isUsedId(id));

  return id;
}

function createTopic(title, description, hasSubtopics) {
  const id = generateUnusedId();

  const topicInstance = {
    id,
    description,
    title,
    state: 'published',
    links: {
      self: {
        href: `${apiUrl}${id}`,
        id,
      },
      content: {
        href: `${apiUrl}${id}/content`,
      },
    },
  };

  if (hasSubtopics) {
    topicInstance.links.subtopics = {
      href: `${apiUrl}${id}/subtopics`,
    };
    topicInstance.subtopics_ids = [];
  }

  const topic = {
    id,
    next: deepCopy(topicInstance),
    current: deepCopy(topicInstance),
  };

  return topic;
}

function createContent(id) {
  return {
    id,
    next: {
      state: 'published',
    },
    current: {
      state: 'published',
    },
  };
}

function addTopicToParent(parentTopicId, topicId) {
  // Add topic id to subtopics of parent
  const parentTopicCursor = topicCollection.find({ id: parentTopicId });

  print(parentTopicCursor.hasNext());

  if (!parentTopicCursor.hasNext()) {
    print('Error: Couldn\'t find the parent topic');
    quit(0);
  }

  const parentTopic = parentTopicCursor.next();

  parentTopic.next.subtopics_ids.push(topicId);
  parentTopic.current.subtopics_ids.push(topicId);

  return parentTopic;
}

function insertTopic(topic) {
  if (cfg.verbose) {
    print('New topic document');
    print(JSON.stringify(topic));
  }

  if (cfg.insert) {
    topicCollection.insertOne(topic);
  }
}

function insertContent(content) {
  if (cfg.verbose) {
    print('New content document');
    print(JSON.stringify(content));
  }

  if (cfg.insert) {
    contentCollection.insertOne(content);
  }
}

function createTopLevelTopic(topicData) {
  const hasSubtopics = topicData.subtopics.length > 0;
  const topLevelTopic = createTopic(topicData.title, topicData.description, hasSubtopics);

  const rootTopic = addTopicToParent(rootId, topLevelTopic.id);
  print('New root topic');
  print(JSON.stringify(rootTopic));

  // Create census subtopics
  topicData.subtopics.forEach((subtopic) => {
    const topic = createTopic(subtopic.title, subtopic.description);
    const content = createContent(topic.id);

    insertTopic(topic);
    insertContent(content);

    // Add subtopic to Census topic
    topLevelTopic.next.subtopics_ids.push(topic.id);
    topLevelTopic.current.subtopics_ids.push(topic.id);
  });

  const topLevelContent = createContent(topLevelTopic.id);
  insertTopic(topLevelTopic);
  insertContent(topLevelContent);
  topicCollection.updateOne({ id: rootId }, { $set: rootTopic });
}

function createTopLevelTopics(topicsData) {
  print(JSON.stringify(topicsData));

  topicsData.forEach((topicData) => createTopLevelTopic(topicData));
}
