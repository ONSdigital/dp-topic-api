package sdk

import (
	"context"

	healthcheck "github.com/ONSdigital/dp-api-clients-go/v2/health"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-topic-api/models"
)

//go:generate moq -out ./mocks/client.go -pkg mocks . Clienter

type Clienter interface {
	Checker(ctx context.Context, check *health.CheckState) error
	GetRootTopicsPrivate(ctx context.Context, reqHeaders Headers) (*models.PrivateSubtopics, error)
	GetRootTopicsPublic(ctx context.Context, reqHeaders Headers) (*models.PublicSubtopics, error)
	GetSubtopicsPrivate(ctx context.Context, reqHeaders Headers, id string) (*models.PrivateSubtopics, error)
	GetSubtopicsPublic(ctx context.Context, reqHeaders Headers, id string) (*models.PublicSubtopics, error)
	Health() *healthcheck.Client
	URL() string
}
