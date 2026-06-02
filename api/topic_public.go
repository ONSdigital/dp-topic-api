package api

import (
	"context"
	"net/http"

	dprequest "github.com/ONSdigital/dp-net/v3/request"
	"github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

const topicRoot = "topic_root" // access specific document to retrieve list

// getRootTopicsPublicHandler is a handler that gets a public list of top level root topics by a specific id from MongoDB for Web
func (api *API) getRootTopicsPublicHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	logdata := log.Data{
		logKeyReqID:   ctx.Value(dprequest.RequestIdKey),
		logKeyTopicID: topicRoot,
		logKeyFunc:    "getTopicsListPublicHandler",
	}

	// The mongo document with id: `topic_root` contains the list of subtopics,
	// so we directly return that list
	api.getSubtopicsPublicByID(ctx, topicRoot, logdata, w)
}

// getTopicPublicHandler is a handler that gets a topic by its id from MongoDB for Web
func (api *API) getTopicPublicHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	logdata := log.Data{
		logKeyReqID:   ctx.Value(dprequest.RequestIdKey),
		logKeyTopicID: id,
		logKeyFunc:    "getTopicPublicHandler",
	}

	if id == topicRoot {
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
	log.Info(ctx, "request successful", logdata) // NOTE: name of function is in logdata
}

// getSubtopicsPublicHandler is a handler that gets a topic by its id from MongoDB for Web
func (api *API) getSubtopicsPublicHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	logdata := log.Data{
		logKeyReqID:   ctx.Value(dprequest.RequestIdKey),
		logKeyTopicID: id,
		logKeyFunc:    "getSubtopicsPublicHandler",
	}

	if id == topicRoot {
		handleError(ctx, w, apierrors.ErrTopicNotFound, logdata)
		return
	}

	api.getSubtopicsPublicByID(ctx, id, logdata, w)
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

	if topic.Current.SubtopicIds == nil || len(*topic.Current.SubtopicIds) == 0 {
		// no subtopics exist for the requested ID
		handleError(ctx, w, apierrors.ErrNotFound, logdata)
		return
	}

	for _, subTopicID := range *topic.Current.SubtopicIds {
		// get sub topic from mongoDB by subTopicID
		topic, err := api.dataStore.Backend.GetTopic(ctx, subTopicID)
		if err != nil {
			logdata["missing subtopic for id"] = subTopicID
			log.Error(ctx, "missing subtopic for id", err, logdata)
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
	log.Info(ctx, "request successful", logdata) // NOTE: name of function is in logdata
}
