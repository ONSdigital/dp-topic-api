package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/pkg/errors"

	healthcheck "github.com/ONSdigital/dp-api-clients-go/v2/health"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/dp-topic-api/models"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	testHost = "http://localhost:25700"
)

var (
	initialTestState = healthcheck.CreateCheckState(service)

	testPublicTopics = models.PublicSubtopics{
		Count:       2,
		Offset:      0,
		Limit:       100,
		TotalCount:  2,
		PublicItems: &[]models.Topic{testPublicTopic1, testPublicTopic2},
	}

	testPublicTopic1 = models.Topic{
		ID:          "1234",
		Description: "Root Topic 1",
		Title:       "Root Topic 1",
		Keywords:    []string{"test"},
		State:       "published",
	}

	testPublicTopic2 = models.Topic{
		ID:          "5678",
		Description: "Root Topic 2",
		Title:       "Root Topic 2",
		Keywords:    []string{"test"},
		State:       "published",
	}
)

func newMockHTTPClient(r *http.Response, err error) *dphttp.ClienterMock {
	return &dphttp.ClienterMock{
		SetPathsWithNoRetriesFunc: func(paths []string) {
		},
		DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return r, err
		},
		GetPathsWithNoRetriesFunc: func() []string {
			return []string{"/healthcheck"}
		},
	}
}

func newTopicAPIClient(t *testing.T, httpClient *dphttp.ClienterMock) *Client {
	healthClient := healthcheck.NewClientWithClienter(service, testHost, httpClient)
	return NewWithHealthClient(healthClient)
}

func TestHealthCheckerClient(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	timePriorHealthCheck := time.Now().UTC()
	path := "/health"

	Convey("Given clienter.Do returns an error", t, func() {
		clientError := errors.New("unexpected error")
		httpClient := newMockHTTPClient(&http.Response{}, clientError)
		topicAPIClient := newTopicAPIClient(t, httpClient)
		check := initialTestState

		Convey("When topic API client Checker is called", func() {
			err := topicAPIClient.Checker(ctx, &check)
			So(err, ShouldBeNil)

			Convey("Then the expected check is returned", func() {
				So(check.Name(), ShouldEqual, service)
				So(check.Status(), ShouldEqual, health.StatusCritical)
				So(check.StatusCode(), ShouldEqual, 0)
				So(check.Message(), ShouldEqual, clientError.Error())
				So(*check.LastChecked(), ShouldHappenAfter, timePriorHealthCheck)
				So(check.LastSuccess(), ShouldBeNil)
				So(*check.LastFailure(), ShouldHappenAfter, timePriorHealthCheck)
			})

			Convey("And client.Do should be called once with the expected parameters", func() {
				doCalls := httpClient.DoCalls()
				So(doCalls, ShouldHaveLength, 1)
				So(doCalls[0].Req.URL.Path, ShouldEqual, path)
			})
		})
	})

	Convey("Given a 500 response for health check", t, func() {
		httpClient := newMockHTTPClient(&http.Response{StatusCode: http.StatusInternalServerError}, nil)
		topicAPIClient := newTopicAPIClient(t, httpClient)
		check := initialTestState

		Convey("When topic API client Checker is called", func() {
			err := topicAPIClient.Checker(ctx, &check)
			So(err, ShouldBeNil)

			Convey("Then the expected check is returned", func() {
				So(check.Name(), ShouldEqual, service)
				So(check.Status(), ShouldEqual, health.StatusCritical)
				So(check.StatusCode(), ShouldEqual, 500)
				So(check.Message(), ShouldEqual, service+healthcheck.StatusMessage[health.StatusCritical])
				So(*check.LastChecked(), ShouldHappenAfter, timePriorHealthCheck)
				So(check.LastSuccess(), ShouldBeNil)
				So(*check.LastFailure(), ShouldHappenAfter, timePriorHealthCheck)
			})

			Convey("And client.Do should be called once with the expected parameters", func() {
				doCalls := httpClient.DoCalls()
				So(doCalls, ShouldHaveLength, 1)
				So(doCalls[0].Req.URL.Path, ShouldEqual, path)
			})
		})
	})
}

func TestGetRootTopicsPublic(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	Convey("Given public root topics is returned successfully", t, func() {
		body, err := json.Marshal(testPublicTopics)
		if err != nil {
			t.Errorf("failed to setup test data, error: %v", err)
		}

		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			},
			nil)

		topicAPIClient := newTopicAPIClient(t, httpClient)

		Convey("When GetRootTopicsPublic is called", func() {
			respRootTopics, err := topicAPIClient.GetRootTopicsPublic(ctx, Headers{})

			Convey("Then the expected public root topics is returned", func() {
				So(*respRootTopics, ShouldResemble, testPublicTopics)

				Convey("And no error is returned", func() {
					So(err, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, "/topics")
					})
				})
			})
		})
	})

	Convey("Given a 500 response from topic api", t, func() {
		httpClient := newMockHTTPClient(&http.Response{StatusCode: http.StatusInternalServerError}, nil)
		topicAPIClient := newTopicAPIClient(t, httpClient)

		Convey("When GetRootTopicsPublic is called", func() {
			respRootTopics, err := topicAPIClient.GetRootTopicsPublic(ctx, Headers{})

			Convey("Then an error should be returned ", func() {
				So(err, ShouldNotBeNil)

				Convey("And the expected public root topics should be nil", func() {
					So(respRootTopics, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, "/topics")
					})
				})
			})
		})
	})

	Convey("Given the client returns an unexpected error", t, func() {
		clientError := errors.New("unexpected error")
		httpClient := newMockHTTPClient(&http.Response{}, clientError)

		topicAPIClient := newTopicAPIClient(t, httpClient)

		Convey("When GetRootTopicsPublic is called", func() {
			respRootTopics, err := topicAPIClient.GetRootTopicsPublic(ctx, Headers{})

			Convey("Then an error should be returned ", func() {
				So(err, ShouldNotBeNil)

				Convey("And the expected public root topics should be nil", func() {
					So(respRootTopics, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, "/topics")
					})
				})
			})
		})
	})
}

func TestGetSubtopicsPublic(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	Convey("Given public subtopics is returned successfully", t, func() {
		body, err := json.Marshal(testPublicTopics)
		if err != nil {
			t.Errorf("failed to setup test data, error: %v", err)
		}

		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			},
			nil)

		topicAPIClient := newTopicAPIClient(t, httpClient)

		Convey("When GetSubtopicsPublic is called", func() {
			respSubtopics, err := topicAPIClient.GetSubtopicsPublic(ctx, Headers{}, "1357")

			Convey("Then the expected public subtopics is returned", func() {
				So(*respSubtopics, ShouldResemble, testPublicTopics)

				Convey("And no error is returned", func() {
					So(err, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, "/topics/1357/subtopics")
					})
				})
			})
		})
	})

	Convey("Given a 500 response from topic api", t, func() {
		httpClient := newMockHTTPClient(&http.Response{StatusCode: http.StatusInternalServerError}, nil)
		topicAPIClient := newTopicAPIClient(t, httpClient)

		Convey("When GetSubtopicsPublic is called", func() {
			respSubtopics, err := topicAPIClient.GetSubtopicsPublic(ctx, Headers{}, "1357")

			Convey("Then an error should be returned ", func() {
				So(err, ShouldNotBeNil)

				Convey("And the expected public subtopics should be nil", func() {
					So(respSubtopics, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, "/topics/1357/subtopics")
					})
				})
			})
		})
	})

	Convey("Given the client returns an unexpected error", t, func() {
		clientError := errors.New("unexpected error")
		httpClient := newMockHTTPClient(&http.Response{}, clientError)

		topicAPIClient := newTopicAPIClient(t, httpClient)

		Convey("When GetSubtopicsPublic is called", func() {
			respSubtopics, err := topicAPIClient.GetSubtopicsPublic(ctx, Headers{}, "1357")

			Convey("Then an error should be returned ", func() {
				So(err, ShouldNotBeNil)

				Convey("And the expected public subtopics should be nil", func() {
					So(respSubtopics, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, "/topics/1357/subtopics")
					})
				})
			})
		})
	})
}
