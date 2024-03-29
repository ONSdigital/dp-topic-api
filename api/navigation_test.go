package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-topic-api/config"
	storetest "github.com/ONSdigital/dp-topic-api/store/mock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNavigationGetNavigationHandler(t *testing.T) {
	Convey("Given a topic API in publishing mode (private endpoints enabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		mockedDataStore := &storetest.StorerMock{}

		api := GetAPIWithMocks(cfg, mockedDataStore)
		w := httptest.NewRecorder()
		So(w.Header(), ShouldBeEmpty)

		var request *http.Request
		request, err = createRequestWithAuth("GET", "http://localhost:25300/navigation", nil)
		So(err, ShouldBeNil)

		api.getNavigationHandler(w, request)
		So(w.Code, ShouldEqual, http.StatusOK)
		So(w.Header(), ShouldNotBeEmpty)
		So(w.Header().Get("Cache-Control"), ShouldEqual, "public, max-age=1800")
	})
}
