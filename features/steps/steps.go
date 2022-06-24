package steps

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ONSdigital/dp-topic-api/config"

	dpMongoDriver "github.com/ONSdigital/dp-mongodb/v3/mongodb"

	componentModels "github.com/ONSdigital/dp-topic-api/features/models"
	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/cucumber/godog"
	"go.mongodb.org/mongo-driver/bson"
)

func (f *TopicComponent) iHaveTheseTopics(topicsWriteJson *godog.DocString) error {
	ctx := context.Background()
	topics := []componentModels.TopicWrite{}
	m := f.MongoClient

	err := json.Unmarshal([]byte(topicsWriteJson.Content), &topics)
	if err != nil {
		return err
	}

	for _, topicsDoc := range topics {
		if err := f.putTopicInDatabase(ctx, m.Connection.Collection(m.ActualCollectionName(config.TopicsCollection)), topicsDoc); err != nil {
			return err
		}
	}

	return nil
}

func (f *TopicComponent) putTopicInDatabase(ctx context.Context, mongoCollection *dpMongoDriver.Collection, topicDoc componentModels.TopicWrite) error {
	update := bson.M{
		"$set": topicDoc,
		"$setOnInsert": bson.M{
			"last_updated": time.Now(),
		},
	}
	_, err := mongoCollection.UpsertById(ctx, topicDoc.ID, update)
	if err != nil {
		return err
	}
	return nil
}

func (f *TopicComponent) iHaveTheseContents(contentJson *godog.DocString) error {
	ctx := context.Background()
	collection := []models.ContentResponse{}
	m := f.MongoClient

	err := json.Unmarshal([]byte(contentJson.Content), &collection)
	if err != nil {
		return err
	}

	for _, topicsDoc := range collection {
		if err := f.putContentInDatabase(ctx, m.Connection.Collection(m.ActualCollectionName(config.ContentCollection)), topicsDoc); err != nil {
			return err
		}
	}

	return nil
}

func (f *TopicComponent) putContentInDatabase(ctx context.Context, mongoCollection *dpMongoDriver.Collection, contentDoc models.ContentResponse) error {
	update := bson.M{
		"$set": contentDoc,
		"$setOnInsert": bson.M{
			"last_updated": time.Now(),
		},
	}
	_, err := mongoCollection.UpsertById(ctx, contentDoc.ID, update)
	if err != nil {
		return err
	}
	return nil
}

func (f *TopicComponent) privateEndpointsAreEnabled() error {
	f.Config.EnablePrivateEndpoints = true
	return nil
}
