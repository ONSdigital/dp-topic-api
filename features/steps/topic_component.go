package steps

import (
	"context"
	"net/http"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-component-test/utils"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/mongo"
	"github.com/ONSdigital/dp-topic-api/service"
	serviceMock "github.com/ONSdigital/dp-topic-api/service/mock"
	"github.com/ONSdigital/dp-topic-api/store"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/cucumber/godog"
	"github.com/gofrs/uuid"

	mongodriver "github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

type TopicComponent struct {
	ErrorFeature   componenttest.ErrorFeature
	svc            *service.Service
	errorChan      chan error
	MongoClient    *mongo.Mongo
	Config         *config.Config
	HTTPServer     *http.Server
	ServiceRunning bool
}

func NewTopicComponent(mongoURL, zebedeeURL string) (*TopicComponent, error) {

	f := &TopicComponent{
		HTTPServer:     &http.Server{},
		errorChan:      make(chan error),
		ServiceRunning: false,
	}

	var err error
	f.Config, err = config.Get()
	if err != nil {
		return nil, err
	}

	f.Config.ClusterEndpoint = mongoURL
	f.Config.ZebedeeURL = zebedeeURL
	f.Config.Database = utils.RandomDatabase()
	f.Config.EnablePrivateEndpoints = false
	// The following is to reset the Username and Password that have been set is Config from the previous
	// config.Get()
	f.Config.Username, f.Config.Password = "", ""
	f.Config.MongoConfig.Username, f.Config.MongoConfig.Password = createCredsInDB(&f.Config.MongoConfig)

	f.MongoClient, err = mongo.NewDBConnection(context.TODO(), f.Config.MongoConfig)
	if err != nil {
		return nil, err
	}

	initMock := &serviceMock.InitialiserMock{
		DoGetMongoDBFunc:     f.DoGetMongoDB,
		DoGetHealthCheckFunc: f.DoGetHealthcheckOk,
		DoGetHTTPServerFunc:  f.DoGetHTTPServer,
	}

	f.svc = service.New(f.Config, service.NewServiceList(initMock))

	return f, nil
}

func createCredsInDB(mongoConfig *mongodriver.MongoDriverConfig) (string, string) {
	mongoConnection, err := mongodriver.Open(mongoConfig)
	if err != nil {
		panic("expected db connection to be opened")
	}

	username := "admin"
	password, _ := uuid.NewV4()
	createCollectionResponse := mongoConnection.RunCommand(context.TODO(), bson.D{{"create", "test"}})

	if createCollectionResponse != nil {
		panic("expected collection creation to go through")
	}
	userCreationResponse := mongoConnection.RunCommand(context.TODO(), bson.D{
		{Key: "createUser", Value: username},
		{Key: "pwd", Value: password.String()},
		{Key: "roles", Value: []bson.M{
			{"role": "root", "db": "admin"},
		}},
	})
	if userCreationResponse != nil {
		panic("expected user creation to go through")
	}
	return username, password.String()
}

func (f *TopicComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^private endpoints are enabled$`, f.privateEndpointsAreEnabled)
	ctx.Step(`^I have these topics:$`, f.iHaveTheseTopics)
	ctx.Step(`^I have these contents:$`, f.iHaveTheseContents)
}

func (f *TopicComponent) Close() error {
	ctx := context.Background()
	err := f.MongoClient.Connection.DropDatabase(ctx)
	if err != nil {
		log.Warn(ctx, "error dropping database on Close()", log.Data{"err": err.Error()})
	}
	if f.svc != nil && f.ServiceRunning {
		err = f.svc.Close(ctx)
		if err != nil {
			log.Warn(ctx, "error closing service on Close()", log.Data{"err": err.Error()})
		}
		f.ServiceRunning = false
	}

	return nil
}

func (f *TopicComponent) InitialiseService() (http.Handler, error) {
	if err := f.svc.Run(context.Background(), "1", "", "", f.errorChan); err != nil {
		return nil, err
	}
	f.ServiceRunning = true
	return f.HTTPServer.Handler, nil
}

func (f *TopicComponent) DoGetHealthcheckOk(_ *config.Config, _ string, _ string, _ string) (service.HealthChecker, error) {
	return &serviceMock.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
	}, nil
}

func (f *TopicComponent) DoGetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	f.HTTPServer.Addr = bindAddr
	f.HTTPServer.Handler = router
	return f.HTTPServer
}

// DoGetMongoDB returns a MongoDB
func (f *TopicComponent) DoGetMongoDB(_ context.Context, _ config.MongoConfig) (store.MongoDB, error) {
	return f.MongoClient, nil
}
