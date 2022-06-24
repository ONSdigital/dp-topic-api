package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/ONSdigital/log.go/v2/log"
)

// GetRootTopicsPrivate gets the private list of top level root topics for Publishing which returns both Next and Current document(s) in the response
func (cli *Client) GetRootTopicsPrivate(ctx context.Context, reqHeaders Headers) (*models.PrivateSubtopics, error) {
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

	var rootTopics models.PrivateSubtopics

	if err = json.Unmarshal(b, &rootTopics); err != nil {
		logData := log.Data{
			"response_bytes": b,
		}
		log.Error(ctx, "failed to unmarshal bytes into root topics", err, logData)
		return nil, err
	}

	return &rootTopics, nil
}

// GetSubtopicsPrivate gets the private list of subtopics of a topic for Publishing which returns both Next and Current document(s) in the response
func (cli *Client) GetSubtopicsPrivate(ctx context.Context, reqHeaders Headers, id string) (*models.PrivateSubtopics, error) {
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

	var subtopics models.PrivateSubtopics

	if err = json.Unmarshal(b, &subtopics); err != nil {
		logData := log.Data{
			"response_bytes": b,
		}
		log.Error(ctx, "failed to unmarshal bytes into subtopics", err, logData)
		return nil, err
	}

	return &subtopics, nil
}
