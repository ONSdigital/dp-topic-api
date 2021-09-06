package api

import (
	"context"
	"net/http"

	dprequest "github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

// getTopicPublicHandler is a handler that gets a topic by its id from MongoDB for Web
func (api *API) getTopicPublicHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	logdata := log.Data{
		"request_id": ctx.Value(dprequest.RequestIdKey),
		"topic_id":   id,
		"function":   "getTopicPublicHandler",
	}

	if id == "topic_root" {
		handleError(ctx, w, apierrors.ErrTopicNotFound, logdata)
		return
	}

	// get topic from mongoDB by id
	topic, err := api.dataStore.Backend.GetTopic(ctx, id)
	if err != nil {
		handleError(ctx, w, err, logdata)
		return
	}

	// User is not authenticated and hence has only access to current sub document
	if err := WriteJSONBody(ctx, topic.Current, w, logdata); err != nil {
		// WriteJSONBody has already logged the error
		return
	}
	log.Event(ctx, "request successful", log.INFO, logdata) // NOTE: name of function is in logdata
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
	log.Event(ctx, "request successful", log.INFO, logdata) // NOTE: name of function is in logdata
}

func (api *API) getSubtopicsPublicByID(ctx context.Context, id string, logdata log.Data, w http.ResponseWriter) {
	// get topic from mongoDB by id
	topic, err := api.dataStore.Backend.GetTopic(ctx, id)
	if err != nil {

		// no topic found to retrieve the subtopics from


		
		handleError(ctx, w, err, logdata)
		return
	}

	// User is not authenticated and hence has only access to current sub document(s)
	var result models.PublicSubtopics

	if topic.Current == nil {
		handleError(ctx, w, apierrors.ErrContentNotFound, logdata)
		return
	}

	if len(topic.Current.SubtopicIds) == 0 {
		// no subtopics exist for the requested ID
		handleError(ctx, w, apierrors.ErrNotFound, logdata)
		return
	}

	for _, subTopicID := range topic.Current.SubtopicIds {
		// get sub topic from mongoDB by subTopicID
		topic, err := api.dataStore.Backend.GetTopic(ctx, subTopicID)
		if err != nil {
			logdata["missing subtopic for id"] = subTopicID
			log.Event(ctx, err.Error(), log.ERROR, logdata)
			continue
		}

		if result.PublicItems == nil {
			result.PublicItems = &[]models.Topic{*topic.Current}
		} else {
			*result.PublicItems = append(*result.PublicItems, *topic.Current)
		}

		result.TotalCount++
	}
	if result.TotalCount == 0 {
		handleError(ctx, w, apierrors.ErrContentNotFound, logdata)
		return
	}

	if err := WriteJSONBody(ctx, result, w, logdata); err != nil {
		// WriteJSONBody has already logged the error
		return
	}
	log.Event(ctx, "request successful", log.INFO, logdata) // NOTE: name of function is in logdata
}

// getSubtopicsPublicHandler is a handler that gets a topic by its id from MongoDB for Web
func (api *API) getSubtopicsPublicHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	logdata := log.Data{
		"request_id": ctx.Value(dprequest.RequestIdKey),
		"topic_id":   id,
		"function":   "getSubtopicsPublicHandler",
	}

	if id == "topic_root" {
		handleError(ctx, w, apierrors.ErrTopicNotFound, logdata)
		return
	}

	api.getSubtopicsPublicByID(ctx, id, logdata, w)
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
			log.Event(ctx, err.Error(), log.ERROR, logdata)
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
	log.Event(ctx, "request successful", log.INFO, logdata) // NOTE: name of function is in logdata
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

// getTopicsListPublicHandler is a handler that gets a public list of top level topics by a specific id from MongoDB for Web
func (api *API) getTopicsListPublicHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id := "topic_root" // access specific document to retrieve list
	logdata := log.Data{
		"request_id": ctx.Value(dprequest.RequestIdKey),
		"topic_id":   id,
		"function":   "getTopicsListPublicHandler",
	}

	// The mongo document with id: `topic_root` contains the list of sobtopics,
	// so we directly return that list
	api.getSubtopicsPublicByID(ctx, id, logdata, w)
}

// getTopicsListPrivateHandler is a handler that gets a private list of top level topics by a specific id from MongoDB for Web
func (api *API) getTopicsListPrivateHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id := "topic_root" // access specific document to retrieve list
	logdata := log.Data{
		"request_id": ctx.Value(dprequest.RequestIdKey),
		"topic_id":   id,
		"function":   "getTopicsListPrivateHandler",
	}

	// The mongo document with id: `topic_root` contains the list of sobtopics,
	// so we directly return that list
	api.getSubtopicsPrivateByID(ctx, id, logdata, w)
}
