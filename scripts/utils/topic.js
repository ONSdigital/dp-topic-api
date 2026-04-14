load("./scripts/utils/db.js");
load("./scripts/utils/content.js");

/**
 * Adds a subtopic to the subtopic_ids array of a topic.
 * @param {string} topicId - The id of the topic to update. This is not the mongo _id
 * @param {string} subtopicId - The id of the subtopic to add. This is not the mongo _id
 */
function addSubtopic(topicId, subtopicId) {
  const parentTopic = findTopicByID(topicId);

  if (!parentTopic) {
    console.error(
      `couldn't find topic id ${topicId} for subtopic id ${subtopicId}`
    );
    return;
  }

  console.log(`Updating parent topic id ${topicId} with subtopics`);

  parentTopic.next.subtopics_ids.push(subtopicId);
  parentTopic.current.subtopics_ids.push(subtopicId);

  updateTopic(topicId, parentTopic);
}

/**
 * Create a topic from a set of data
 * @param {string} parentID - The id of the topic to update. This is not the mongo _id
 * @param {object} topicData - The topic data object - attributes: title, description, subtopics (which are then topic data objects)
 * @param {string} topicID - The id of the topic to add. This is not the mongo _id
 */
function createTopic(parentID, topicData, topicID) {
  console.log(
    `creating topic with topicID ${topicID} and title ${topicData.title}`
  );
  const hasSubtopics = topicData.subtopics && topicData.subtopics.length > 0;
  const topic = generateTopic(
    topicData.title,
    topicData.description,
    topicID,
    hasSubtopics
  );

  // Create subtopics
  if (hasSubtopics) {
    console.log(`creating subtopics for topicID ${topicID}`);

    topicData.subtopics.forEach((subtopic, index) => {
      const subtopicID = `${topicID}-sub-${index}`;
      topic.current.subtopics_ids.push(subtopicID);
      topic.next.subtopics_ids.push(subtopicID);
      createTopic(topic.id, subtopic, subtopicID);
    });
  }

  const inserted = insertTopic(topic);

  const content = generateContent(topic.id);
  insertContent(content);

  if (inserted) {
    addSubtopic(parentID, topic.id);
  }
}

/**
 * Find a topic by its slug.
 * @param {string} slug - The slug of the topic.
 * @returns {object} - The found topic object or null if not found.
 */
function findTopicBySlug(slug) {
  return getTopicCollection().findOne({ "current.slug": slug });
}

/**
 * Find a topic by its id.
 * @param {string} id - The id of the topic. This is not the mongo _id
 * @returns {object} - The found topic object or null if not found.
 */
function findTopicByID(id) {
  return getTopicCollection().findOne({ id: id });
}

/**
 * Generates a full topic document from the following params:
 * @param {string} title - The title of the topic.
 * @param {string} description - The description of the topic.
 * @param {string} id - The id of the topic. This is not the mongo _id
 * @param {boolean} hasSubtopics - Does it have subtopics
 * @returns {object} - The found topic object or null if not found.
 */
function generateTopic(title, description, id, hasSubtopics) {
  const topicInstance = {
    id,
    description,
    title,
    state: "published",
    slug: generateSlug(title),
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

/**
 * Inserts a topic document into mongodb
 * @param {object} topic - The topic document to insert
 */
function insertTopic(topic) {
  console.log(`inserting topic with id ${topic.id}`);
  const topicWithIDExists = findTopicByID(topic.id);
  const topicWithSlugExists = findTopicBySlug(topic.next.slug);

  if (!topicWithSlugExists && !topicWithIDExists) {
    if (cfg.insert) {
      try {
        getTopicCollection().insertOne(topic);
        return true;
      } catch (err) {
        console.log(err);
        return false;
      }
    }
    return false;
  } else {
    console.warn(
      `topic with slug ${topic.next.slug} or id ${topic.id} already exists`
    );
    return false;
  }
}

/**
 * Updates a topic document in mongodb
 * @param {string} topicID - The topic ID to update. This is not the mongo _id
 * @param {object} newTopic - The new topic document.
 *
 */
function updateTopic(topicID, newTopic) {
  console.log(`updating topic with id ${topicID}`);
  if (cfg.insert) {
    getTopicCollection().updateOne({ id: topicID }, { $set: newTopic });
  }
}
