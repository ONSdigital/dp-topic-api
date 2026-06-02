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

// getRootTopicsPrivateHandler is a handler that gets a private list of top level root topics by a specific id from MongoDB for Publishing
func (api *API) getRootTopicsPrivateHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	const id = "topic_root" // access specific document to retrieve list
	logdata := log.Data{
		logKeyReqID:   ctx.Value(dprequest.RequestIdKey),
		logKeyTopicID: id,
		logKeyFunc:    "getTopicsListPrivateHandler",
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
		logKeyReqID:   ctx.Value(dprequest.RequestIdKey),
		logKeyTopicID: id,
		logKeyFunc:    "getTopicPrivateHandler",
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
		logKeyReqID:   ctx.Value(dprequest.RequestIdKey),
		logKeyTopicID: id,
		logKeyFunc:    "getSubtopicsPrivateHandler",
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

	if topic.Next.SubtopicIds == nil || len(*topic.Next.SubtopicIds) == 0 {
		// no subtopics exist for the requested ID
		handleError(ctx, w, apierrors.ErrNotFound, logdata)
		return
	}

	for _, subTopicID := range *topic.Next.SubtopicIds {
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

func (api *API) putTopicReleaseDatePrivateHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	logdata := log.Data{
		logKeyTopicID: id,
		logKeyFunc:    "putTopicReleaseDatePrivateHandler",
	}

	topicRelease, err := models.ReadReleaseDate(req.Body)
	if err != nil {
		handleError(ctx, w, err, logdata)
		return
	}

	releaseDate, err := topicRelease.Validate()
	if err != nil {
		handleError(ctx, w, err, logdata)
		return
	}

	// update topic next.release_date in mongo db
	if err := api.dataStore.Backend.UpdateReleaseDate(ctx, id, *releaseDate); err != nil {
		handleError(ctx, w, err, logdata)
		return
	}

	w.WriteHeader(http.StatusOK)

	log.Info(ctx, "request successful", logdata)
}

func (api *API) putTopicStatePrivateHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	state := vars["state"]
	logdata := log.Data{
		logKeyTopicID: id,
		"state":       state,
		logKeyFunc:    "putTopicStatePrivateHandler",
	}

	_, err := models.ParseState(state)
	if err != nil {
		handleError(ctx, w, err, logdata)
		return
	}

	if state == models.StatePublished.String() {
		log.Info(ctx, "attempting to publish topic", logdata)
		if err := api.publishTopic(ctx, id); err != nil {
			handleError(ctx, w, err, logdata)
		}
	} else {
		// update topic next.state in mongo db
		if err := api.dataStore.Backend.UpdateState(ctx, id, state); err != nil {
			handleError(ctx, w, err, logdata)
			return
		}
	}

	w.WriteHeader(http.StatusOK)

	log.Info(ctx, "request successful", logdata)
}

func (api *API) putTopicPrivateHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	logdata := log.Data{
		logKeyTopicID: id,
		logKeyFunc:    "putTopicPrivateHandler",
	}

	topicUpdate, err := models.ReadTopicUpdate(req.Body)
	if err != nil {
		handleError(ctx, w, err, logdata)
		return
	}

	if err := topicUpdate.ValidateUpdate(); err != nil {
		handleError(ctx, w, err, logdata)
		return
	}

	if err := api.dataStore.Backend.CheckTopicExists(ctx, id); err != nil {
		handleError(ctx, w, err, logdata)
		return
	}

	// update topic in mongo db
	if err := api.dataStore.Backend.UpdateTopic(ctx, api.topicAPIURL, id, topicUpdate); err != nil {
		handleError(ctx, w, err, logdata)
		return
	}

	if topicUpdate.State == models.StatePublished.String() {
		log.Info(ctx, "attempting to publish topic", logdata)
		if err := api.publishTopic(ctx, id); err != nil {
			handleError(ctx, w, err, logdata)
		}
	}

	w.WriteHeader(http.StatusOK)

	log.Info(ctx, "request successful", logdata)
}

func (api *API) publishTopic(ctx context.Context, id string) error {
	// TODO - should lock resource, put this in a mongo db transaction or use eTags to
	// check if the resource has changed since initial request - as it is not a public
	// endpoint and is not currently used by the publishing system we can ignore this for
	// now

	// get topic
	topic, err := api.dataStore.Backend.GetTopic(ctx, id)
	if err != nil {
		return err
	}

	// set next state
	topic.Next.State = models.StatePublished.String()

	// update local copy of topic
	newTopic := syncNextAndCurrentTopic(topic)

	// update topic in mongo db
	err = api.dataStore.Backend.UpsertTopic(ctx, id, newTopic)
	if err != nil {
		return err
	}

	return nil
}

func syncNextAndCurrentTopic(topic *models.TopicResponse) *models.TopicResponse {
	return &models.TopicResponse{
		ID:      topic.ID,
		Next:    topic.Next,
		Current: topic.Next,
	}
}
