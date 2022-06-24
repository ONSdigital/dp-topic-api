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

func TestGetTopicPublicHandler(t *testing.T) {

	Convey("Given a topic API in web mode (private endpoints disabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		cfg.EnablePrivateEndpoints = false
		Convey("And a topic API with mongoDB returning 'next' and 'current' topics", func() {

			mongoDBMock := &storeMock.MongoDBMock{
				GetTopicFunc: func(ctx context.Context, id string) (*models.TopicResponse, error) {
					switch id {
					case testTopicID1:
						return dbTopic(models.StatePublished), nil
					default:
						return nil, apierrors.ErrTopicNotFound
					}
				},
			}

			topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

			Convey("When an existing 'published' topic is requested with the valid Topic-Id context value", func() {
				request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s", testTopicID1), nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected sub-document topic is returned with status code 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retTopic := models.Topic{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic, ShouldResemble, *createdTopicCurrent())
				})
			})

			Convey("Requesting an nonexistent topic ID results in a NotFound response", func() {
				request := httptest.NewRequest(http.MethodGet, "http://localhost:25300/topics/inexistent", nil)
				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}

func TestGetSubtopicsPublicHandler(t *testing.T) {

	Convey("Given a topic API in web mode (private endpoints disabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		cfg.EnablePrivateEndpoints = false
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
					default:
						return nil, apierrors.ErrTopicNotFound
					}
				},
			}

			topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

			// 1 has subtopics & points to 2 & 3
			Convey("When an existing 'published' subtopic is requested with the valid Topic-Id value 1", func() {
				request := httptest.NewRequest(http.MethodGet, "http://localhost:25300/topics/1/subtopics", nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected sub-documents is returned with status code 200, and documents with ID's 2 & 3 returned", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retTopic := models.PublicSubtopics{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic.TotalCount, ShouldEqual, 2)
					So((*retTopic.PublicItems)[0].ID, ShouldEqual, "2")
					So((*retTopic.PublicItems)[1].ID, ShouldEqual, "3")
				})
			})

			// 2 has subtopics & points to 4, 6 (but ID 6 does not exist)
			Convey("When an existing 'published' subtopic is requested with the valid Topic-Id value 2", func() {
				request := httptest.NewRequest(http.MethodGet, "http://localhost:25300/topics/2/subtopics", nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected sub-document is returned with status code 200, and document with ID 4 is returned", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retTopic := models.PublicSubtopics{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic.TotalCount, ShouldEqual, 1)
					So((*retTopic.PublicItems)[0].ID, ShouldEqual, "4")
				})
			})

			// 3 has subtopics, but the ID 5 in the list does not exist
			Convey("When an existing 'published' subtopic is requested with the valid Topic-Id value 3", func() {
				request := httptest.NewRequest(http.MethodGet, "http://localhost:25300/topics/3/subtopics", nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then no sub-documents are returned and we get status code 500", func() {
					So(w.Code, ShouldEqual, http.StatusNotFound)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					So(payload, ShouldResemble, []byte("content not found\n"))
				})
			})

			// 4 has NO subtopics, so is an end node that has a content link
			Convey("When an existing 'published' subtopic is requested with the valid Topic-Id value 4", func() {
				request := httptest.NewRequest(http.MethodGet, "http://localhost:25300/topics/4/subtopics", nil)

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
				request := httptest.NewRequest(http.MethodGet, "http://localhost:25300/topics/inexistent/subtopics", nil)
				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})

			// topic_root for test uses dbTopic1 which has subtopics & points to 2 & 3
			Convey("When an existing 'published' /topics/topic_root/subtopics document is requested", func() {
				request := httptest.NewRequest(http.MethodGet, "http://localhost:25300/topics/topic_root/subtopics", nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected status code 404 is returned, because this is not avaible for public web mode", func() {
					So(w.Code, ShouldEqual, http.StatusNotFound)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					So(payload, ShouldResemble, []byte("topic not found\n"))
				})
			})

		})
	})
}

func TestGetTopicsListPublicHandler(t *testing.T) {
	Convey("Given a topic API in web mode (private endpoints disabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		cfg.EnablePrivateEndpoints = false
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

			// topic_root for test uses dbTopic1 which has subtopics & points to 2 & 3
			Convey("When an existing 'published' /topics list is requested", func() {
				request := httptest.NewRequest(http.MethodGet, "http://localhost:25300/topics", nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected sub-documents is returned with status code 200, and documents with ID's 2 & 3 returned", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retTopic := models.PublicSubtopics{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic.TotalCount, ShouldEqual, 2)
					So((*retTopic.PublicItems)[0].ID, ShouldEqual, "2")
					So((*retTopic.PublicItems)[1].ID, ShouldEqual, "3")
				})
			})

			// topic_root for test uses dbTopic1 which has subtopics & points to 2 & 3
			Convey("When an existing 'published' /topics/topic_root document is requested", func() {
				request := httptest.NewRequest(http.MethodGet, "http://localhost:25300/topics/topic_root", nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected status code 404 is returned, because this is not avaible for public web mode", func() {
					So(w.Code, ShouldEqual, http.StatusNotFound)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					So(payload, ShouldResemble, []byte("topic not found\n"))
				})
			})

			// No more tests needed because getting the 'topic_root' makes use of
			// function getSubtopicsPublicByID() which is fully tested in
			// TestGetSubtopicsPublicHandler() above, preventing duplication of tests.
		})
	})
}
