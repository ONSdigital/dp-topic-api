// insert-census-topics.js
//
// Create census topic and subtopics

const topicsCollection = 'topics'
const contentCollection = 'content'
const idSize = 4
const rootId = "topic_root"
const idAlphabet = '123456789'
const apiUrl = "http://localhost:25300/topics/"

if (typeof(cfg) == "undefined") {
    // default, but can be changed on command-line, see README
    cfg = {
        verbose:  false,    // display the new documents
        insert:   true      // set to false to avoid inserts
    }
}

//////////////////////////

var subtopics = [ 
                { title: "Ageing", description: "Facts and figures to help understand the needs of an ageing population in England and Wales including changes over time and the characteristics of certain age groups." },
                { title: "Demography", description: "Facts and figures to help understand the population of England and Wales, including age, sex, household composition and legal partnerships." },
                { title: "Education", description: "Facts and figures to help understand the education level of people in England and Wales, and how it varies by age, sex, area, ethnicity, disability and other characteristics." },
                { title: "Equalities", description: "Facts and figures to help understand equality for people in England and Wales, including what life is like for people with certain characteristics." },
                { title: "Ethnic group, national identity, language and religion", description: "Facts and figures to help understand people's ethnicity, national identity, language and religion in England and Wales, including changes over time and Welsh language." },
                { title: "Health, disability and unpaid care", description: "Facts and figures to help understand health, disability and unpaid care for people in England and Wales, and how it varies by areas and other characteristics." },
                { title: "Historic census", description: "Facts and figures to help understand how life has changed over time, using data from the 2011 Census and earlier." },
                { title: "Housing", description: "Facts and figures to help understand the types of housing people live in, in England and Wales, including how it has changed over time, by area and household and property characteristics." },
                { title: "International migration", description: "Facts and figures to help understand people who have moved in and out of the UK within England and Wales, including changes over time, people with more than one passport and second-generation migration." },
                { title: "Labour market", description: "Facts and figures to help understand the labour market for people in England and Wales, including how it varies by area and other characteristics." },
                { title: "Sexual orientation and gender identity", description: "Facts and figures to help understand the sexual orientation and gender identity of people in England and Wales, including how it varies by area and characteristics such as demography, housing, employment and education." },
                { title: "Travel to work", description: "Facts and figures to help understand how people in England and Wales travel to work, including mode of transport, distance and how it changes across rural and urban areas." },
                { title: "Veterans", description: "Facts and figures to help understand the veteran population of England and Wales, including topics such as housing, education, employment and skills." }
            ]

function isUsedId(id) {
    return db.getCollection(topicsCollection).find({id:id}).hasNext()
}

function makeId() {
    var result = ''
    for (var i = 0; i < idSize; i++) {
      result += idAlphabet.charAt(Math.floor(Math.random() * idAlphabet.length))
    }
    return result
}

function createTopic(title, description) {
    do {
        var id = makeId()
    } while (isUsedId(id))
    var topic = {
        id: id,
        next : {
            id : id,
            description : description,
            title : title,
            state : "published",
            links : {
                self : {
                    href : apiUrl + id,
                    id : id
                },
                content : {
                    href : apiUrl + id + "/content"
                }
            }
        },
        current : {
            id : id,
            description : description,
            title : title,
            state : "published",
            links : {
                self : {
                    href : apiUrl + id,
                    id : id
                },
                content : {
                    href : apiUrl + id + "/content"
                }
            }
        }
    }

    return topic;
}

function createContent(id) {
    return {
        id : id,
        next : {
            state : "published"
        },
        current : {
            state : "published"
        }
    }
}


//////////////////////////

// Create Census topic
var censusTopic = createTopic("Census", "Census")
censusTopic.next.links.subtopics = {
    href : apiUrl + censusTopic.id + "/subtopics",
}
censusTopic.current.links.subtopics = {
    href : apiUrl + censusTopic.id + "/subtopics",
}
censusTopic.next.subtopics_ids = []
censusTopic.current.subtopics_ids = []

// Add census topic to subtopics of root
var rootTopicCursor = db.getCollection(topicsCollection).find({id:rootId})
if (!rootTopicCursor.hasNext()) {
    print("Error: Couldn't find the root topic")
    quit(0)
}
var rootTopic = rootTopicCursor.next()
rootTopic.next.subtopics_ids.push(censusTopic.id)
rootTopic.current.subtopics_ids.push(censusTopic.id)

// Create census subtopics
for (var idx in subtopics) {
    var subtopicDefiniton = subtopics[idx];
    var topic = createTopic(subtopicDefiniton.title, subtopicDefiniton.description)
    var content = createContent(topic.id)

    if (cfg.verbose) {
        print("New subtopic document")
        print(JSON.stringify(topic))        
        print("New content document")
        print(JSON.stringify(content))
    }

    if (cfg.insert) {
        db.getCollection(topicsCollection).insertOne(topic)
        db.getCollection(contentCollection).insertOne(content)
    }

    // Add subtopic to Census topic
    censusTopic.next.subtopics_ids.push(topic.id)
    censusTopic.current.subtopics_ids.push(topic.id)
}

var censusContent = createContent(censusTopic.id)
if (cfg.verbose) {
    print("New Census topic")
    print(JSON.stringify(censusTopic))
    print("New Census content")
    print(JSON.stringify(censusContent))
    print("New topic root")
    print(JSON.stringify(rootTopic))
}

if (cfg.insert) {
    db.getCollection(topicsCollection).insertOne(censusTopic)
    db.getCollection(contentCollection).insertOne(censusContent)
    db.getCollection(topicsCollection).updateOne({id:rootId}, {$set : rootTopic} )
}
