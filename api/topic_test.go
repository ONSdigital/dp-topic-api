package api

import (
	"fmt"

	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/mocks"
	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/ONSdigital/dp-topic-api/store"
	"github.com/gorilla/mux"
)

// Constants for testing
const (
	testTopicID1 = "topicTopicID1"
)

func dbTopicWithID(state models.State, id string) *models.TopicResponse {
	return &models.TopicResponse{
		ID: id,
		Next: &models.Topic{
			ID:          id,
			Description: "next test description - 1",
			Title:       "test title - 1",
			Keywords:    []string{"keyword 1", "keyword 2", "keyword 3"},
			State:       state.String(),
			Links: &models.TopicLinks{
				Self: &models.LinkObject{
					HRef: fmt.Sprintf("http://example.com/topics/%s", id),
					ID:   id,
				},
				Subtopics: &models.LinkObject{
					HRef: fmt.Sprintf("http://example.com/topics/%s/subtopics", id),
				},
				Content: &models.LinkObject{
					HRef: fmt.Sprintf("http://example.com/topics/%s/content", id),
				},
			},
		},
		Current: &models.Topic{
			ID:          id,
			Description: "current test description - 1",
			Title:       "test title - 1",
			Keywords:    []string{"keyword 1", "keyword 2", "keyword 3"},
			State:       state.String(),
			Links: &models.TopicLinks{
				Self: &models.LinkObject{
					HRef: fmt.Sprintf("http://example.com/topics/%s", id),
					ID:   id,
				},
				Subtopics: &models.LinkObject{
					HRef: fmt.Sprintf("http://example.com/topics/%s/subtopics", id),
				},
				Content: &models.LinkObject{
					HRef: fmt.Sprintf("http://example.com/topics/%s/content", id),
				},
			},
		},
	}
}

// DB model corresponding to a topic in the provided state, without any download variant
func dbTopic(state models.State) *models.TopicResponse {
	return dbTopicWithID(state, testTopicID1)
}

// API model corresponding to TopicResponse
func createdTopicAll() *models.TopicResponse {
	return dbTopic(models.StateCreated)
}

func dbTopicCurrentWithID(state models.State, id string) *models.Topic {
	return &models.Topic{
		ID:          id,
		Description: "current test description - 1",
		Title:       "test title - 1",
		Keywords:    []string{"keyword 1", "keyword 2", "keyword 3"},
		State:       state.String(),
		Links: &models.TopicLinks{
			Self: &models.LinkObject{
				HRef: fmt.Sprintf("http://example.com/topics/%s", id),
				ID:   id,
			},
			Subtopics: &models.LinkObject{
				HRef: fmt.Sprintf("http://example.com/topics/%s/subtopics", id),
			},
			Content: &models.LinkObject{
				HRef: fmt.Sprintf("http://example.com/topics/%s/content", id),
			},
		},
	}
}

// create just the 'current' sub-document
func dbTopicCurrent(state models.State) *models.Topic {
	return dbTopicCurrentWithID(state, testTopicID1)
}

func createdTopicCurrent() *models.Topic {
	return dbTopicCurrent(models.StatePublished)
}

// NOTE: The data within the following four sets of data returning functions
//       are interlinked with one another by the SubtopicsIds

// ================= - 1 has subtopics & points to 2 & 3
// DB model corresponding to a topic in the provided state, without any download variant
func dbTopic1(state models.State) *models.TopicResponse {
	return &models.TopicResponse{
		ID: "1",
		Next: &models.Topic{
			ID:    "1",
			State: state.String(),
			Links: &models.TopicLinks{
				Subtopics: &models.LinkObject{
					HRef: "http://example.com/topics/1/subtopics",
				},
			},
			SubtopicIds: []string{"2", "3"},
		},
		Current: &models.Topic{
			ID:    "1",
			State: state.String(),
			Links: &models.TopicLinks{
				Subtopics: &models.LinkObject{
					HRef: "http://example.com/topics/1/subtopics",
				},
			},
			SubtopicIds: []string{"2", "3"},
		},
	}
}

// ================= - 2 has subtopics & points to 4, 6 (but ID 6 does not exist)
// DB model corresponding to a topic in the provided state, without any download variant
func dbTopic2(state models.State) *models.TopicResponse {
	return &models.TopicResponse{
		ID: "2",
		Next: &models.Topic{
			ID:    "2",
			State: state.String(),
			Links: &models.TopicLinks{
				Subtopics: &models.LinkObject{
					HRef: "http://example.com/topics/2/subtopics",
				},
			},
			SubtopicIds: []string{"4", "6"},
		},
		Current: &models.Topic{
			ID:    "2",
			State: state.String(),
			Links: &models.TopicLinks{
				Subtopics: &models.LinkObject{
					HRef: "http://example.com/topics/2/subtopics",
				},
			},
			SubtopicIds: []string{"4", "6"},
		},
	}
}

// ================= - 3 has subtopics, but the ID 5 in the list does not exist
// DB model corresponding to a topic in the provided state, without any download variant
func dbTopic3(state models.State) *models.TopicResponse {
	return &models.TopicResponse{
		ID: "3",
		Next: &models.Topic{
			ID:    "3",
			State: state.String(),
			Links: &models.TopicLinks{
				Subtopics: &models.LinkObject{
					HRef: "http://example.com/topics/3/subtopics",
				},
			},
			SubtopicIds: []string{"5"},
		},
		Current: &models.Topic{
			ID:    "3",
			State: state.String(),
			Links: &models.TopicLinks{
				Subtopics: &models.LinkObject{
					HRef: "http://example.com/topics/3/subtopics",
				},
			},
			SubtopicIds: []string{"5"},
		},
	}
}

// ================= - 4 has NO subtopics, so is an end node that has a content link
// DB model corresponding to a topic in the provided state, without any download variant
func dbTopic4(state models.State) *models.TopicResponse {
	return &models.TopicResponse{
		ID: "4",
		Next: &models.Topic{
			ID:    "4",
			State: state.String(),
			Links: &models.TopicLinks{
				Content: &models.LinkObject{
					HRef: "http://example.com/topics/4/content",
				},
			},
		},
		Current: &models.Topic{
			ID:    "4",
			State: state.String(),
			Links: &models.TopicLinks{
				Content: &models.LinkObject{
					HRef: "http://example.com/topics/4/content",
				},
			},
		},
	}
}

// GetAPIWithMocks also used in other tests, so exported
func GetAPIWithMocks(cfg *config.Config, mockedDataStore store.Storer) *API {
	mu.Lock()
	defer mu.Unlock()
	//	urlBuilder := url.NewBuilder("http://example.com")

	permissions := mocks.NewAuthHandlerMock()

	return Setup(testContext, cfg, mux.NewRouter(), store.DataStore{Backend: mockedDataStore}, permissions)
}
