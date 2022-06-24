package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/models"
	storeMock "github.com/ONSdigital/dp-topic-api/store/datastoretest"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetTopicPrivateHandler(t *testing.T) {
	Convey("Given a topic API in publishing mode (private endpoints enabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		cfg.EnablePrivateEndpoints = true
		Convey("And a topic API with mongoDB returning 'created' and 'full' topics", func() {

			mongoDBMock := &storeMock.MongoDBMock{
				GetTopicFunc: func(ctx context.Context, id string) (*models.TopicResponse, error) {
					switch id {
					case testTopicID1:
						return dbTopic(models.StateCreated), nil
					default:
						return nil, apierrors.ErrTopicNotFound
					}
				},
			}
			topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

			Convey("When an existing 'created' topic is requested with the valid Topic-Id context value", func() {
				request, err := createRequestWithAuth(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s", testTopicID1), nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected topic is returned with status code 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retTopic := models.TopicResponse{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic, ShouldResemble, *createdTopicAll())
				})
			})

			Convey("Requesting an nonexistent topic ID results in a NotFound response", func() {
				request, err := createRequestWithAuth(http.MethodGet, "http://localhost:25300/topics/inexistent", nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}

func TestGetSubtopicsPrivateHandler(t *testing.T) {
	Convey("Given a topic API in web mode (private endpoints enabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		cfg.EnablePrivateEndpoints = true
		Convey("And a topic API with mongoDB returning 'next' and 'current' topics", func() {

			mongoDBMock := &storeMock.MongoDBMock{
				GetTopicFunc: func(ctx context.Context, id string) (*models.TopicResponse, error) {
					switch id {
					case "1":
						return dbTopic1(models.StatePublished), nil
					case "2":
						return dbTopic2(models.StatePublished), nil
					case "3":
						return dbTopic3(models.StatePublished), nil
					case "4":
						return dbTopic4(models.StatePublished), nil
					case "topic_root":
						return dbTopic1(models.StatePublished), nil
					default:
						return nil, apierrors.ErrTopicNotFound
					}
				},
			}

			topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

			// 1 has subtopics & points to 2 & 3
			Convey("When an existing 'published' subtopic is requested with the valid Topic-Id value 1", func() {
				request, err := createRequestWithAuth(http.MethodGet, "http://localhost:25300/topics/1/subtopics", nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected sub-documents is returned with status code 200, and documents with ID's 2 & 3 returned", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retTopic := models.PrivateSubtopics{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic.TotalCount, ShouldEqual, 2)
					So((*retTopic.PrivateItems)[0].Current.ID, ShouldEqual, "2")
					So((*retTopic.PrivateItems)[1].Current.ID, ShouldEqual, "3")
				})
			})

			// 2 has subtopics & points to 4, 6 (but ID 6 does not exist)
			Convey("When an existing 'published' subtopic is requested with the valid Topic-Id value 2", func() {
				request, err := createRequestWithAuth(http.MethodGet, "http://localhost:25300/topics/2/subtopics", nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected sub-document is returned with status code 200, and document with ID 4 is returned", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retTopic := models.PrivateSubtopics{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic.TotalCount, ShouldEqual, 1)
					So((*retTopic.PrivateItems)[0].Current.ID, ShouldEqual, "4")
				})
			})

			// 3 has subtopics, but the ID 5 in the list does not exist
			Convey("When an existing 'published' subtopic is requested with the valid Topic-Id value 3", func() {
				request, err := createRequestWithAuth(http.MethodGet, "http://localhost:25300/topics/3/subtopics", nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then no sub-documents are returned and we get status code 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					So(payload, ShouldResemble, []byte("internal error\n"))
				})
			})

			// 4 has NO subtopics, so is an end node that has a content link
			Convey("When an existing 'published' subtopic is requested with the valid Topic-Id value 4", func() {
				request, err := createRequestWithAuth(http.MethodGet, "http://localhost:25300/topics/4/subtopics", nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then no sub-documents are returned and we get status code 404", func() {
					So(w.Code, ShouldEqual, http.StatusNotFound)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					So(payload, ShouldResemble, []byte("not found\n"))
				})
			})

			Convey("Requesting an nonexistent topic ID results in a NotFound response", func() {
				request, err := createRequestWithAuth(http.MethodGet, "http://localhost:25300/topics/inexistent/subtopics", nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})

			// topic_root for test uses dbTopic1 which has subtopics & points to 2 & 3
			Convey("When an existing 'published' /topics/topic_root/subtopics document is requested", func() {
				request, err := createRequestWithAuth(http.MethodGet, "http://localhost:25300/topics/topic_root/subtopics", nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected sub-documents is returned with status code 200, and documents with ID's 2 & 3 returned", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)

					So(err, ShouldBeNil)
					retTopic := models.PrivateSubtopics{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic.TotalCount, ShouldEqual, 2)
					So((*retTopic.PrivateItems)[0].Current.ID, ShouldEqual, "2")
					So((*retTopic.PrivateItems)[1].Current.ID, ShouldEqual, "3")
				})
			})

		})
	})
}

func TestGetTopicsListPrivateHandler(t *testing.T) {
	Convey("Given a topic API in web mode (private endpoints enabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		cfg.EnablePrivateEndpoints = true
		Convey("And a topic API with mongoDB returning 'next' and 'current' topics", func() {

			mongoDBMock := &storeMock.MongoDBMock{
				GetTopicFunc: func(ctx context.Context, id string) (*models.TopicResponse, error) {
					switch id {
					case "2":
						return dbTopic2(models.StatePublished), nil
					case "3":
						return dbTopic3(models.StatePublished), nil
					case "topic_root":
						return dbTopic1(models.StatePublished), nil
					default:
						return nil, apierrors.ErrTopicNotFound
					}
				},
			}

			topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

			// topic_root for test uses 1 which has subtopics & points to 2 & 3
			Convey("When an existing 'published' /topics list is requested", func() {
				request, err := createRequestWithAuth(http.MethodGet, "http://localhost:25300/topics", nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected sub-documents is returned with status code 200, and documents with ID's 2 & 3 returned", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retTopic := models.PrivateSubtopics{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic.TotalCount, ShouldEqual, 2)
					So((*retTopic.PrivateItems)[0].Current.ID, ShouldEqual, "2")
					So((*retTopic.PrivateItems)[1].Current.ID, ShouldEqual, "3")
				})
			})

			// topic_root for test uses dbTopic1 which has subtopics & points to 2 & 3
			Convey("When an existing 'published' /topics/topic_root document is requested", func() {
				request, err := createRequestWithAuth(http.MethodGet, "http://localhost:25300/topics/topic_root", nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected sub-documents is returned with status code 200, and documents with ID's 2 & 3 returned", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)

					retTopic := models.TopicResponse{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic.ID, ShouldEqual, "1")
					So(retTopic.Next.ID, ShouldEqual, "1")
				})
			})
		})
	})
}
