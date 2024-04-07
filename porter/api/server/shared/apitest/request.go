package apitest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/porter-dev/porter/api/server/shared"
	"github.com/porter-dev/porter/api/server/shared/apierrors"
	"github.com/porter-dev/porter/api/server/shared/config"
)

func GetRequestAndRecorder(t *testing.T, method, route string, requestObj interface{}) (*http.Request, *httptest.ResponseRecorder) {
	var reader io.Reader = nil

	if requestObj != nil {
		data, err := json.Marshal(requestObj)
		if err != nil {
			t.Fatal(err)
		}

		reader = strings.NewReader(string(data))
	}

	// method and route don't actually matter since this is meant to test handlers
	req, err := http.NewRequest(method, route, reader)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	return req, rr
}

func WithURLParams(t *testing.T, req *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	routeParams := &chi.RouteParams{}

	for key, val := range params {
		routeParams.Add(key, val)
	}

	rctx.URLParams = *routeParams

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	return req
}

type failingDecoderValidator struct {
	config *config.Config
}

func (f *failingDecoderValidator) DecodeAndValidate(
	w http.ResponseWriter,
	r *http.Request,
	v interface{},
) (ok bool) {
	apierrors.HandleAPIError(f.config.Logger, f.config.Alerter, w, r, apierrors.NewErrInternal(fmt.Errorf("fake error")), true)
	return false
}

func (f *failingDecoderValidator) DecodeAndValidateNoWrite(
	r *http.Request,
	v interface{},
) error {
	return fmt.Errorf("fake error")
}

func NewFailingDecoderValidator(config *config.Config) shared.RequestDecoderValidator {
	return &failingDecoderValidator{config}
}
