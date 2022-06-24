package models

// Navigation is used to get high level list of topics and subtopics with links and description for site navigation.
type Navigation struct {
	Description string                 `json:"description"`
	Links       *TopicLinks            `json:"links,omitempty"`
	Items       *[]TopicNonReferential `json:"items,omitempty"`
}

// TopicNonReferential is used to create a single comprehensive list of topics and subtopics.
type TopicNonReferential struct {
	Description   string                 `json:"description,omitempty"`
	Label         string                 `json:"label"`
	Links         *TopicLinks            `json:"links,omitempty"`
	Name          string                 `json:"name"`
	SubtopicItems *[]TopicNonReferential `json:"subtopics,omitempty"`
	Title         string                 `json:"title"`
	Uri           string                 `json:"uri"`
}
