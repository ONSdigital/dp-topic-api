package api

import (
	"net/http"

	dprequest "github.com/ONSdigital/dp-net/v3/request"
	"github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// getContentPublicHandler is a handler that gets content by its id from MongoDB for Web
func (api *API) getContentPublicHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	logdata := log.Data{
		logKeyReqID:  ctx.Value(dprequest.RequestIdKey),
		"content_id": id,
		logKeyFunc:   "getContentPublicHandler",
	}

	// get type from query parameters, or default value
	queryTypeFlags := getContentTypeParameter(req.URL.Query())
	if queryTypeFlags == 0 {
		handleError(ctx, w, apierrors.ErrContentUnrecognisedParameter, logdata)
		return
	}

	// check topic from mongoDB by id
	err := api.dataStore.Backend.CheckTopicExists(ctx, id)
	if err != nil {
		handleError(ctx, w, err, logdata)
		return
	}

	// get content from mongoDB by id
	content, err := api.dataStore.Backend.GetContent(ctx, id, queryTypeFlags)
	if err != nil {
		// no content found
		handleError(ctx, w, err, logdata)
		return
	}

	// User is not authenticated and hence has only access to current sub document(s)

	if content.Current == nil {
		handleError(ctx, w, apierrors.ErrContentNotFound, logdata)
		return
	}

	currentResult := getRequiredItems(queryTypeFlags, content.Current, content.ID)

	if currentResult.TotalCount == 0 {
		handleError(ctx, w, apierrors.ErrContentNotFound, logdata)
		return
	}

	if err := WriteJSONBody(ctx, currentResult, w, logdata); err != nil {
		// WriteJSONBody has already logged the error
		return
	}
	log.Info(ctx, "request successful", logdata) // NOTE: name of function is in logdata
}

// getContentPrivateHandler is a handler that gets content by its id from MongoDB for Publishing
func (api *API) getContentPrivateHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	logdata := log.Data{
		logKeyReqID:  ctx.Value(dprequest.RequestIdKey),
		"content_id": id,
		logKeyFunc:   "getContentPrivateHandler",
	}

	// get type from query parameters, or default value
	queryTypeFlags := getContentTypeParameter(req.URL.Query())
	if queryTypeFlags == 0 {
		handleError(ctx, w, apierrors.ErrContentUnrecognisedParameter, logdata)
		return
	}

	// check topic from mongoDB by id
	err := api.dataStore.Backend.CheckTopicExists(ctx, id)
	if err != nil {
		handleError(ctx, w, err, logdata)
		return
	}

	// get content from mongoDB by id
	content, err := api.dataStore.Backend.GetContent(ctx, id, queryTypeFlags)
	if err != nil {
		// no content found
		handleError(ctx, w, err, logdata)
		return
	}

	// User has valid authentication to get raw full content document(s)

	if content.Current == nil {
		/*
			TODO
			In the future: when the API becomes more than read-only
			When a document is first created, it will only have 'next' until it is published, when it gets 'current' populated.
			So current == nil is not an error.

			For now we return an error because we dont have publishing steps.
		*/
		handleError(ctx, w, apierrors.ErrInternalServer, logdata)
		return
	}
	if content.Next == nil {
		handleError(ctx, w, apierrors.ErrInternalServer, logdata)
		return
	}

	currentResult := getRequiredItems(queryTypeFlags, content.Current, content.ID)

	// The 'Next' type items may have a different length to the current, so we do the above again, but for Next
	nextResult := getRequiredItems(queryTypeFlags, content.Next, content.ID)

	if currentResult.TotalCount == 0 && nextResult.TotalCount == 0 {
		handleError(ctx, w, apierrors.ErrContentNotFound, logdata)
		return
	}

	var result models.PrivateContentResponseAPI
	result.Next = nextResult
	result.Current = currentResult

	if err := WriteJSONBody(ctx, result, w, logdata); err != nil {
		// WriteJSONBody has already logged the error
		return
	}
	log.Info(ctx, "request successful", logdata) // NOTE: name of function is in logdata
}
