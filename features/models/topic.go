package models

import "github.com/ONSdigital/dp-topic-api/models"

// TopicWrite is used for component testing
type TopicWrite struct {
	ID      string  `bson:"id,omitempty"       json:"id,omitempty"`
	Next    *TopicW `bson:"next,omitempty"     json:"next,omitempty"`
	Current *TopicW `bson:"current,omitempty"  json:"current,omitempty"`
}

// TopicW is used for component testing
type TopicW struct {
	ID          string             `bson:"id,omitempty"             json:"id,omitempty"`
	Description string             `bson:"description,omitempty"    json:"description,omitempty"`
	Title       string             `bson:"title,omitempty"          json:"title,omitempty"`
	Keywords    []string           `bson:"keywords,omitempty"       json:"keywords,omitempty"`
	State       string             `bson:"state,omitempty"          json:"state,omitempty"`
	Links       *models.TopicLinks `bson:"links,omitempty"          json:"links,omitempty"`
	SubtopicIds []string           `bson:"subtopics_ids,omitempty"  json:"subtopics_ids,omitempty"`
}
