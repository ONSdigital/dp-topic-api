package api

import (
	"context"
	"net/http"

	dprequest "github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// getRootTopicsPrivateHandler is a handler that gets a private list of top level root topics by a specific id from MongoDB for Publishing
func (api *API) getRootTopicsPrivateHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id := "topic_root" // access specific document to retrieve list
	logdata := log.Data{
		"request_id": ctx.Value(dprequest.RequestIdKey),
		"topic_id":   id,
		"function":   "getTopicsListPrivateHandler",
	}

	// The mongo document with id: `topic_root` contains the list of subtopics,
	// so we directly return that list
	api.getSubtopicsPrivateByID(ctx, id, logdata, w)
}

// getTopicPrivateHandler is a handler that gets a topic by its id from MongoDB for Publishing
func (api *API) getTopicPrivateHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	logdata := log.Data{
		"request_id": ctx.Value(dprequest.RequestIdKey),
		"topic_id":   id,
		"function":   "getTopicPrivateHandler",
	}

	// get topic from mongoDB by id
	topic, err := api.dataStore.Backend.GetTopic(ctx, id)
	if err != nil {
		handleError(ctx, w, err, logdata)
		return
	}

	// User has valid authentication to get raw topic document
	if err := WriteJSONBody(ctx, topic, w, logdata); err != nil {
		// WriteJSONBody has already logged the error
		return
	}
	log.Info(ctx, "request successful", logdata) // NOTE: name of function is in logdata
}

// getSubtopicsPrivateHandler is a handler that gets a topic by its id from MongoDB for Publishing
func (api *API) getSubtopicsPrivateHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	logdata := log.Data{
		"request_id": ctx.Value(dprequest.RequestIdKey),
		"topic_id":   id,
		"function":   "getSubtopicsPrivateHandler",
	}

	api.getSubtopicsPrivateByID(ctx, id, logdata, w)
}

func (api *API) getSubtopicsPrivateByID(ctx context.Context, id string, logdata log.Data, w http.ResponseWriter) {
	// get topic from mongoDB by id
	topic, err := api.dataStore.Backend.GetTopic(ctx, id)
	if err != nil {
		// no topic found to retrieve the subtopics from
		handleError(ctx, w, err, logdata)
		return
	}

	// User has valid authentication to get raw full topic document(s)
	var result models.PrivateSubtopics

	if topic.Next == nil {
		handleError(ctx, w, apierrors.ErrInternalServer, logdata)
		return
	}

	if len(topic.Next.SubtopicIds) == 0 {
		// no subtopics exist for the requested ID
		handleError(ctx, w, apierrors.ErrNotFound, logdata)
		return
	}

	for _, subTopicID := range topic.Next.SubtopicIds {
		// get topic from mongoDB by subTopicID
		topic, err := api.dataStore.Backend.GetTopic(ctx, subTopicID)
		if err != nil {
			logdata["missing subtopic for id"] = subTopicID
			log.Error(ctx, "missing subtopic for id", err, logdata)
			continue
		}

		if result.PrivateItems == nil {
			result.PrivateItems = &[]models.TopicResponse{*topic}
		} else {
			*result.PrivateItems = append(*result.PrivateItems, *topic)
		}

		result.TotalCount++
	}
	if result.TotalCount == 0 {
		handleError(ctx, w, apierrors.ErrInternalServer, logdata)
		return
	}

	if err := WriteJSONBody(ctx, result, w, logdata); err != nil {
		// WriteJSONBody has already logged the error
		return
	}
	log.Info(ctx, "request successful", logdata) // NOTE: name of function is in logdata
}
