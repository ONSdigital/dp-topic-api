package sdk

import (
	"net/http"

	dprequest "github.com/ONSdigital/dp-net/v2/request"
	"github.com/ONSdigital/log.go/v2/log"
)

type Headers struct {
	ServiceAuthToken string `json:"-"`
	UserAuthToken    string `json:"-"`
}

func (h *Headers) Add(req *http.Request) error {
	ctx := req.Context()

	if h == nil {
		log.Info(ctx, "the Headers struct is nil so there are no headers to add to the request")
		return nil
	}

	if h.ServiceAuthToken != "" {
		dprequest.AddServiceTokenHeader(req, h.ServiceAuthToken)
	}

	if h.UserAuthToken != "" {
		dprequest.AddFlorenceHeader(req, h.UserAuthToken)
	}

	return nil
}
