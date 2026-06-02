package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ONSdigital/dp-authorisation/auth"
	dphandlers "github.com/ONSdigital/dp-net/v3/handlers"
	"github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/store"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

const (
	logKeyReqID   = "request_id"
	logKeyFunc    = "function"
	logKeyTopicID = "topic_id"
)

// nolint:unused // just for this block
var (
	createPermission = auth.Permissions{Create: true}
	readPermission   = auth.Permissions{Read: true}
	updatePermission = auth.Permissions{Update: true}
	deletePermission = auth.Permissions{Delete: true}
)

// AuthHandler provides authorisation checks on requests
type AuthHandler interface {
	Require(required auth.Permissions, handler http.HandlerFunc) http.HandlerFunc
}

// API provides a struct to wrap the api around
type API struct {
	Router                 *mux.Router
	dataStore              store.DataStore
	enablePrivateEndpoints bool
	navigationCacheMaxAge  string
	permissions            AuthHandler
	topicAPIURL            string
}

// Setup function sets up the api and returns an api
func Setup(ctx context.Context, cfg *config.Config, router *mux.Router, dataStore store.DataStore, permissions AuthHandler, topicAPIURL string) *API {
	api := &API{
		Router:                 router,
		dataStore:              dataStore,
		enablePrivateEndpoints: cfg.EnablePrivateEndpoints,
		navigationCacheMaxAge:  fmt.Sprintf("%.0f", cfg.NavigationCacheMaxAge.Seconds()),
		permissions:            permissions,
		topicAPIURL:            topicAPIURL,
	}

	if cfg.EnablePrivateEndpoints {
		// create publishing related endpoints ...
		log.Info(ctx, "enabling private endpoints for topic api")

		api.enablePrivateTopicEndpoints()
	} else {
		// create web related endpoints ...

		log.Info(ctx, "enabling only public endpoints for topic api")
		api.enablePublicEndpoints()
	}

	return api
}

// enablePublicEndpoints register only the public GET endpoints.
func (api *API) enablePublicEndpoints() {
	api.get("/navigation", api.getNavigationHandler)
	api.get("/topics", api.getRootTopicsPublicHandler)
	api.get("/topics/{id}", api.getTopicPublicHandler)
	api.get("/topics/{id}/content", api.getContentPublicHandler)
	api.get("/topics/{id}/subtopics", api.getSubtopicsPublicHandler)
}

// enablePrivateTopicEndpoints register the topics endpoints with the appropriate authentication and authorisation
// checks required when running the topic API in publishing (private) mode.
func (api *API) enablePrivateTopicEndpoints() {
	api.get(
		"/topics/{id}",
		api.isAuthenticated(
			api.isAuthorised(readPermission, api.getTopicPrivateHandler)),
	)

	api.get(
		"/topics/{id}/subtopics",
		api.isAuthenticated(
			api.isAuthorised(readPermission, api.getSubtopicsPrivateHandler)),
	)

	api.get(
		"/topics/{id}/content",
		api.isAuthenticated(
			api.isAuthorised(readPermission, api.getContentPrivateHandler)),
	)

	api.get(
		"/topics",
		api.isAuthenticated(
			api.isAuthorised(readPermission, api.getRootTopicsPrivateHandler)),
	)

	api.get(
		"/navigation",
		api.isAuthenticated(
			api.isAuthorised(readPermission, api.getNavigationHandler)),
	)

	api.put(
		"/topics/{id}/release-date",
		api.isAuthenticated(
			api.isAuthorised(updatePermission, api.putTopicReleaseDatePrivateHandler)),
	)

	api.put(
		"/topics/{id}/state/{state}",
		api.isAuthenticated(
			api.isAuthorised(updatePermission, api.putTopicStatePrivateHandler)),
	)

	api.put(
		"/topics/{id}",
		api.isAuthenticated(
			api.isAuthorised(updatePermission, api.putTopicPrivateHandler)),
	)
}

// isAuthenticated wraps a http handler func in another http handler func that checks the caller is authenticated to
// perform the requested action. handler is the http.HandlerFunc to wrap in an
// authentication check. The wrapped handler is only called if the caller is authenticated
func (api *API) isAuthenticated(handler http.HandlerFunc) http.HandlerFunc {
	return dphandlers.CheckIdentity(handler)
}

// isAuthorised wraps a http.HandlerFunc another http.HandlerFunc that checks the caller is authorised to perform the
// requested action. required is the permissions required to perform the action, handler is the http.HandlerFunc to
// apply the check to. The wrapped handler is only called if the caller has the required permissions.
func (api *API) isAuthorised(required auth.Permissions, handler http.HandlerFunc) http.HandlerFunc {
	return api.permissions.Require(required, handler)
}

// get register a GET http.HandlerFunc.
func (api *API) get(path string, handler http.HandlerFunc) {
	api.Router.HandleFunc(path, handler).Methods("GET")
}

// get register a PUT http.HandlerFunc.
func (api *API) put(path string, handler http.HandlerFunc) {
	api.Router.HandleFunc(path, handler).Methods("PUT")
}

// get register a POST http.HandlerFunc.
//
//nolint:unused // Lint:Method is not currently used but will keep as it's a good reference/completeness for future use.
func (api *API) post(path string, handler http.HandlerFunc) {
	api.Router.HandleFunc(path, handler).Methods("POST")
}

// get register a DELETE http.HandlerFunc.
//
//nolint:unused // Lint:Method is not currently used but will keep as it's a good reference/completeness for future use.
func (api *API) delete(path string, handler http.HandlerFunc) {
	api.Router.HandleFunc(path, handler).Methods("DELETE")
}

// WriteJSONBody marshals the provided interface into json, and writes it to the response body.
func WriteJSONBody(ctx context.Context, v interface{}, w http.ResponseWriter, data log.Data) error {
	// Set headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Marshal provided model
	payload, err := json.Marshal(v)
	if err != nil {
		handleError(ctx, w, apierrors.ErrInternalServer, data)
		return err
	}

	// Write payload to body
	if _, err := w.Write(payload); err != nil {
		// a stack trace is added for Non User errors
		data["response_status"] = http.StatusInternalServerError
		log.Error(ctx, "request unsuccessful", err, data)
		return err
	}

	return nil
}

// ReadJSONBody reads the bytes from the provided body, and marshals it to the provided model interface.
func ReadJSONBody(ctx context.Context, body io.ReadCloser, v interface{}, w http.ResponseWriter, data log.Data) error {
	defer body.Close()

	// Set headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Get Body bytes
	payload, err := io.ReadAll(body)
	if err != nil {
		handleError(ctx, w, apierrors.ErrUnableToReadMessage, data)
		return err
	}

	// Unmarshal body bytes to model
	if err := json.Unmarshal(payload, v); err != nil {
		handleError(ctx, w, apierrors.ErrUnableToParseJSON, data)
		return err
	}

	return nil
}

// handleError is a utility function that maps api errors to an http status code and sets the provided responseWriter accordingly
func handleError(ctx context.Context, w http.ResponseWriter, err error, data log.Data) {
	var status int
	if err != nil {
		switch err {
		case apierrors.ErrTopicNotFound,
			apierrors.ErrContentNotFound,
			apierrors.ErrNotFound:
			status = http.StatusNotFound
		case apierrors.ErrUnableToReadMessage,
			apierrors.ErrUnableToParseJSON:
			status = http.StatusInternalServerError
		case apierrors.ErrContentUnrecognisedParameter,
			apierrors.ErrEmptyRequestBody,
			apierrors.ErrInvalidReleaseDate,
			apierrors.ErrTopicInvalidState,
			apierrors.ErrTopicMissingFields:
			status = http.StatusBadRequest
		case apierrors.ErrTopicStateTransitionNotAllowed:
			status = http.StatusForbidden
		default:
			status = http.StatusInternalServerError
		}
	}

	if data == nil {
		data = log.Data{}
	}

	switch status {
	case http.StatusNotFound, http.StatusForbidden, http.StatusBadRequest:
		data["response_status"] = status
		data["user_error"] = err.Error()
		log.Error(ctx, "request unsuccessful", errors.New("request unsuccessful"), data)
		http.Error(w, err.Error(), status)
	default:
		// a stack trace is added for Non User errors
		data["response_status"] = status
		log.Error(ctx, "request unsuccessful", err, data)
		http.Error(w, apierrors.ErrInternalServer.Error(), status)
	}
}
