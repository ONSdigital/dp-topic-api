package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	dprequest "github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/mocks"
	"github.com/ONSdigital/dp-topic-api/store"
	storetest "github.com/ONSdigital/dp-topic-api/store/datastoretest"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	mu          sync.Mutex
	testContext = context.Background()
)

func TestPublishedSubnetEndpointsAreDisabled(t *testing.T) {

	type testEndpoint struct {
		Method string
		URL    string
	}

	publishSubnetEndpoints := map[testEndpoint]int{
		// Topics endpoints
		{Method: "GET", URL: "http://localhost:25300/navigation"}: http.StatusNotFound,
	}

	Convey("When the API is started with private endpoints disabled and permission auth disabled", t, func() {

		for endpoint, expectedStatusCode := range publishSubnetEndpoints {
			Convey("The following endpoint "+endpoint.URL+"(Method:"+endpoint.Method+") should return 404", func() {
				request, err := createRequestWithAuth(endpoint.Method, endpoint.URL, nil)
				So(err, ShouldBeNil)

				cfg, err := config.Get()
				So(err, ShouldBeNil)
				cfg.EnablePrivateEndpoints = false
				cfg.EnablePermissionsAuth = false

				w := httptest.NewRecorder()
				mockedDataStore := &storetest.StorerMock{}
				api := GetWebAPIWithMocks(testContext, cfg, mockedDataStore, nil)

				api.Router.ServeHTTP(w, request)

				So(w.Code, ShouldEqual, expectedStatusCode)
			})
		}
	})
}

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		permissions := mocks.NewAuthHandlerMock()

		Convey("When created in Publishing mode", func() {
			cfg, err := config.Get()
			So(err, ShouldBeNil)
			cfg.EnablePrivateEndpoints = true
			mockedDataStore := &storetest.StorerMock{}

			api := GetWebAPIWithMocks(testContext, cfg, mockedDataStore, permissions)

			Convey("When created the following routes should have been added", func() {
				So(api, ShouldNotBeNil)
				So(hasRoute(api.Router, "/topics", "GET"), ShouldBeTrue)
				So(hasRoute(api.Router, "/topics/{id}", "GET"), ShouldBeTrue)
				So(hasRoute(api.Router, "/topics/{id}/subtopics", "GET"), ShouldBeTrue)
				So(hasRoute(api.Router, "/topics/{id}/content", "GET"), ShouldBeTrue)
				So(hasRoute(api.Router, "/navigation", "GET"), ShouldBeTrue)
			})
		})

		Convey("When created in Web mode", func() {
			cfg, err := config.Get()
			So(err, ShouldBeNil)
			cfg.EnablePrivateEndpoints = false

			mockedDataStore := &storetest.StorerMock{}

			Convey("And when the app auth token is not set", func() {
				cfg.AppAuthToken = ""
				cfg.EnablePermissionsAuth = true
				api := GetWebAPIWithMocks(testContext, cfg, mockedDataStore, permissions)

				Convey("Then the following routes should have been added", func() {
					So(api, ShouldNotBeNil)
					So(hasRoute(api.Router, "/topics", "GET"), ShouldBeTrue)
					So(hasRoute(api.Router, "/topics/{id}", "GET"), ShouldBeTrue)
					So(hasRoute(api.Router, "/topics/{id}/subtopics", "GET"), ShouldBeTrue)
					So(hasRoute(api.Router, "/topics/{id}/content", "GET"), ShouldBeTrue)
					So(hasRoute(api.Router, "/navigation", "GET"), ShouldBeFalse)
				})
			})

			Convey("and when app auth token is set", func() {
				cfg.AppAuthToken = "4946B8C8-77F7-47E2-BABA-B71E0B264A93"
				api := GetWebAPIWithMocks(testContext, cfg, mockedDataStore, permissions)

				Convey("Then the following routes should have been added", func() {
					So(api, ShouldNotBeNil)
					So(hasRoute(api.Router, "/topics", "GET"), ShouldBeTrue)
					So(hasRoute(api.Router, "/topics/{id}", "GET"), ShouldBeTrue)
					So(hasRoute(api.Router, "/topics/{id}/subtopics", "GET"), ShouldBeTrue)
					So(hasRoute(api.Router, "/topics/{id}/content", "GET"), ShouldBeTrue)
					So(hasRoute(api.Router, "/navigation", "GET"), ShouldBeTrue)
				})
			})
		})
	})
}

// GetWebAPIWithMocks also used in other tests, so exported
func GetWebAPIWithMocks(ctx context.Context, cfg *config.Config, mockedDataStore store.Storer, permissions AuthHandler) *API {
	mu.Lock()
	defer mu.Unlock()

	return Setup(ctx, cfg, mux.NewRouter(), store.DataStore{Backend: mockedDataStore}, permissions)
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, nil)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}

func createRequestWithAuth(method, url string, body io.Reader) (*http.Request, error) {
	request := httptest.NewRequest(method, url, body)

	ctx := request.Context()
	ctx = dprequest.SetCaller(ctx, "someone@ons.gov.uk")
	request = request.WithContext(ctx)
	return request, nil
}
