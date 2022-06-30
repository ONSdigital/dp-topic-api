package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	healthcheck "github.com/ONSdigital/dp-api-clients-go/v2/health"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-topic-api/models"
	apiError "github.com/ONSdigital/dp-topic-api/sdk/errors"
	"github.com/ONSdigital/log.go/v2/log"
)

const (
	service = "dp-topic-api"
)

type Client struct {
	hcCli *healthcheck.Client
}

// New creates a new instance of Client with a given topic api url
func New(topicAPIURL string) *Client {
	return &Client{
		hcCli: healthcheck.NewClient(service, topicAPIURL),
	}
}

// NewWithHealthClient creates a new instance of topic API Client,
// reusing the URL and Clienter from the provided healthcheck client
func NewWithHealthClient(hcCli *healthcheck.Client) *Client {
	return &Client{
		hcCli: healthcheck.NewClientWithClienter(service, hcCli.URL, hcCli.Client),
	}
}

// URL returns the URL used by this client
func (cli *Client) URL() string {
	return cli.hcCli.URL
}

// Health returns the underlying Healthcheck Client for this topic API client
func (cli *Client) Health() *healthcheck.Client {
	return cli.hcCli
}

// Checker calls topic api health endpoint and returns a check object to the caller
func (cli *Client) Checker(ctx context.Context, check *health.CheckState) error {
	return cli.hcCli.Checker(ctx, check)
}

// GetRootTopicsPublic gets the public list of top level root topics for Web which returns the Current document(s) in the response
func (cli *Client) GetRootTopicsPublic(ctx context.Context, reqHeaders Headers) (*models.PublicSubtopics, error) {
	path := fmt.Sprintf("%s/topics", cli.hcCli.URL)

	b, err := cli.callTopicAPI(ctx, path, http.MethodGet, reqHeaders, nil)
	if err != nil {
		logData := log.Data{
			"path":        path,
			"method":      http.MethodGet,
			"req_headers": reqHeaders,
			"body":        nil,
		}
		log.Error(ctx, "failed to call topic api", err, logData)
		return nil, err
	}

	var rootTopics models.PublicSubtopics

	if err = json.Unmarshal(b, &rootTopics); err != nil {
		logData := log.Data{
			"response_bytes": b,
		}
		log.Error(ctx, "failed to unmarshal bytes into root topics", err, logData)
		return nil, err
	}

	return &rootTopics, nil
}

// GetSubtopicsPublic gets the public list of subtopics of a topic for Web which returns the Current document(s) in the response
func (cli *Client) GetSubtopicsPublic(ctx context.Context, reqHeaders Headers, id string) (*models.PublicSubtopics, error) {
	path := fmt.Sprintf("%s/topics/%s/subtopics", cli.hcCli.URL, id)

	b, err := cli.callTopicAPI(ctx, path, http.MethodGet, reqHeaders, nil)
	if err != nil {
		logData := log.Data{
			"path":        path,
			"method":      http.MethodGet,
			"req_headers": reqHeaders,
			"body":        nil,
		}
		log.Error(ctx, "failed to call topic api", err, logData)
		return nil, err
	}

	var subtopics models.PublicSubtopics

	if err = json.Unmarshal(b, &subtopics); err != nil {
		logData := log.Data{
			"response_bytes": b,
		}
		log.Error(ctx, "failed to unmarshal bytes into subtopics", err, logData)
		return nil, err
	}

	return &subtopics, nil
}

// GetNavigationPublic gets the public list of navigation items
func (cli *Client) GetNavigationPublic(ctx context.Context, reqHeaders Headers, lang string) (*models.Navigation, error) {
	path := fmt.Sprintf("%s/navigation?lang=%s", cli.hcCli.URL, lang)

	b, err := cli.callTopicAPI(ctx, path, http.MethodGet, reqHeaders, nil)
	if err != nil {
		logData := log.Data{
			"path":        path,
			"method":      http.MethodGet,
			"req_headers": reqHeaders,
			"body":        nil,
		}
		log.Error(ctx, "failed to call topic api", err, logData)
		return nil, err
	}

	var navigation models.Navigation

	if err = json.Unmarshal(b, &navigation); err != nil {
		logData := log.Data{
			"response_bytes": b,
		}
		log.Error(ctx, "failed to unmarshal bytes into navigation", err, logData)
		return nil, err
	}

	return &navigation, nil
}

// callTopicAPI calls the Topic API endpoint given by path for the provided REST method, request headers, and body payload.
// It returns the response body and any error that occurred.
func (cli *Client) callTopicAPI(ctx context.Context, path, method string, headers Headers, payload []byte) ([]byte, error) {
	URL, err := url.Parse(path)
	if err != nil {
		return nil, apiError.StatusError{
			Err:  fmt.Errorf("failed to parse path: \"%v\" error is: %v", path, err),
			Code: http.StatusInternalServerError,
		}
	}

	path = URL.String()

	var req *http.Request

	if payload != nil {
		req, err = http.NewRequest(method, path, bytes.NewReader(payload))
	} else {
		req, err = http.NewRequest(method, path, http.NoBody)
	}

	// check req, above, didn't error
	if err != nil {
		return nil, apiError.StatusError{
			Err:  fmt.Errorf("failed to create request for call to topic api, error is: %v", err),
			Code: http.StatusInternalServerError,
		}
	}

	if payload != nil {
		req.Header.Add("Content-type", "application/json")
	}

	err = headers.Add(req)
	if err != nil {
		return nil, apiError.StatusError{
			Err:  fmt.Errorf("failed to add headers to request, error is: %v", err),
			Code: http.StatusInternalServerError,
		}
	}

	resp, err := cli.hcCli.Client.Do(ctx, req)
	if err != nil {
		return nil, apiError.StatusError{
			Err:  fmt.Errorf("failed to call topic api, error is: %v", err),
			Code: http.StatusInternalServerError,
		}
	}
	defer func() {
		err = closeResponseBody(resp)
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 400 {
		return nil, apiError.StatusError{
			Err:  fmt.Errorf("failed as unexpected code from topic api: %v", resp.StatusCode),
			Code: resp.StatusCode,
		}
	}

	if resp.Body == nil {
		return nil, nil
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, apiError.StatusError{
			Err:  fmt.Errorf("failed to read response body from call to topic api, error is: %v", err),
			Code: http.StatusInternalServerError,
		}
	}

	return b, nil
}

// closeResponseBody closes the response body and logs an error if unsuccessful
func closeResponseBody(resp *http.Response) error {
	if resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			return apiError.StatusError{
				Err:  fmt.Errorf("error closing http response body from call to topic api, error is: %v", err),
				Code: http.StatusInternalServerError,
			}
		}
	}

	return nil
}
