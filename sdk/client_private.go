package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-topic-api/models"
	apiError "github.com/ONSdigital/dp-topic-api/sdk/errors"
)

// GetRootTopicsPrivate gets the private list of top level root topics for Publishing which returns both Next and Current document(s) in the response
func (cli *Client) GetRootTopicsPrivate(ctx context.Context, reqHeaders Headers) (*models.PrivateSubtopics, apiError.Error) {
	path := fmt.Sprintf("%s/topics", cli.hcCli.URL)

	b, apiErr := cli.callTopicAPI(ctx, path, http.MethodGet, reqHeaders, nil)
	if apiErr.Error() != "nil" {
		return nil, apiErr
	}

	var rootTopics models.PrivateSubtopics

	if err := json.Unmarshal(b, &rootTopics); err != nil {
		return nil, apiError.StatusError{
			Err:  fmt.Errorf("failed to unmarshal rootTopics - error is: %v", err),
			Code: apiErr.Status(),
		}
	}

	return &rootTopics, apiError.StatusError{
		Err:  nil,
		Code: apiErr.Status(),
	}
}

// GetSubtopicsPrivate gets the private list of subtopics of a topic for Publishing which returns both Next and Current document(s) in the response
func (cli *Client) GetSubtopicsPrivate(ctx context.Context, reqHeaders Headers, id string) (*models.PrivateSubtopics, apiError.Error) {
	path := fmt.Sprintf("%s/topics/%s/subtopics", cli.hcCli.URL, id)

	b, apiErr := cli.callTopicAPI(ctx, path, http.MethodGet, reqHeaders, nil)
	if apiErr.Error() != "nil" {
		return nil, apiErr
	}

	var subtopics models.PrivateSubtopics

	if err := json.Unmarshal(b, &subtopics); err != nil {
		return nil, apiError.StatusError{
			Err:  fmt.Errorf("failed to unmarshal subtopics - error is: %v", err),
			Code: apiErr.Status(),
		}
	}

	return &subtopics, apiError.StatusError{
		Err:  nil,
		Code: apiErr.Status(),
	}
}
