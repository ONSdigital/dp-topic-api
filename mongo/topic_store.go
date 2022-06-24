package mongo

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-topic-api/api"
	errs "github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/models"

	mongohealth "github.com/ONSdigital/dp-mongodb/v3/health"
	mongodriver "github.com/ONSdigital/dp-mongodb/v3/mongodb"

	"go.mongodb.org/mongo-driver/bson"
)

type Mongo struct {
	mongodriver.MongoDriverConfig

	Connection   *mongodriver.MongoConnection
	healthClient *mongohealth.CheckMongoClient
}

// NewDBConnection creates a new Mongo object encapsulating a connection to the mongo server/cluster with the given configuration,
// and a health client to check the health of the mongo server/cluster
func NewDBConnection(_ context.Context, cfg config.MongoConfig) (m *Mongo, err error) {
	m = &Mongo{MongoDriverConfig: cfg}
	m.Connection, err = mongodriver.Open(&m.MongoDriverConfig)
	if err != nil {
		return nil, err
	}

	databaseCollectionBuilder := map[mongohealth.Database][]mongohealth.Collection{
		mongohealth.Database(m.Database): {
			mongohealth.Collection(m.ActualCollectionName(config.TopicsCollection)),
			mongohealth.Collection(m.ActualCollectionName(config.ContentCollection)),
		},
	}
	m.healthClient = mongohealth.NewClientWithCollections(m.Connection, databaseCollectionBuilder)

	return m, nil
}

// Close closes the mongo session and returns any error
// It is an error to call m.Close if m.Init() returned an error, and there is no open connection
func (m *Mongo) Close(ctx context.Context) error {
	return m.Connection.Close(ctx)
}

// Checker is called by the healthcheck library to check the health state of this mongoDB instance
func (m *Mongo) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return m.healthClient.Checker(ctx, state)
}

// GetTopic retrieves a topic document by its ID
func (m *Mongo) GetTopic(ctx context.Context, id string) (*models.TopicResponse, error) {
	var topic models.TopicResponse

	err := m.Connection.Collection(m.ActualCollectionName(config.TopicsCollection)).FindOne(ctx, bson.M{"id": id}, &topic)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocumentFound) {
			return nil, errs.ErrTopicNotFound
		}
		return nil, err
	}

	return &topic, nil
}

// CheckTopicExists checks that the topic exists
func (m *Mongo) CheckTopicExists(ctx context.Context, id string) error {

	count, err := m.Connection.Collection(m.ActualCollectionName(config.TopicsCollection)).Count(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}

	if count == 0 {
		return errs.ErrTopicNotFound
	}

	return nil
}

// GetContent retrieves a content document by its ID
func (m *Mongo) GetContent(ctx context.Context, id string, queryTypeFlags int) (*models.ContentResponse, error) {
	var content models.ContentResponse
	// init default, used to minimise the mongo response to minimise go HEAP usage
	contentSelect := bson.M{
		"ID":            1,
		"next.id":       1,
		"next.state":    1,
		"current.id":    1,
		"current.state": 1,
	}

	// Add spotlight first
	if (queryTypeFlags & api.QuerySpotlightFlag) != 0 {
		contentSelect["next.spotlight"] = 1
		contentSelect["current.spotlight"] = 1
	}

	// then Publications
	if (queryTypeFlags & api.QueryArticlesFlag) != 0 {
		contentSelect["next.articles"] = 1
		contentSelect["current.articles"] = 1
	}

	if (queryTypeFlags & api.QueryBulletinsFlag) != 0 {
		contentSelect["next.bulletins"] = 1
		contentSelect["current.bulletins"] = 1
	}

	if (queryTypeFlags & api.QueryMethodologiesFlag) != 0 {
		contentSelect["next.methodologies"] = 1
		contentSelect["current.methodologies"] = 1
	}

	if (queryTypeFlags & api.QueryMethodologyArticlesFlag) != 0 {
		contentSelect["next.methodology_articles"] = 1
		contentSelect["current.methodology_articles"] = 1
	}

	// then Datasets
	if (queryTypeFlags & api.QueryStaticDatasetsFlag) != 0 {
		contentSelect["next.static_datasets"] = 1
		contentSelect["current.static_datasets"] = 1
	}

	if (queryTypeFlags & api.QueryTimeseriesFlag) != 0 {
		contentSelect["next.timeseries"] = 1
		contentSelect["current.timeseries"] = 1
	}

	err := m.Connection.Collection(m.ActualCollectionName(config.ContentCollection)).FindOne(ctx, bson.M{"id": id}, &content, mongodriver.Projection(contentSelect))
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocumentFound) {
			return nil, errs.ErrContentNotFound
		}
		return nil, err
	}

	return &content, nil
}
