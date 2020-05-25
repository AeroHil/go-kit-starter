package test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	abendpoint "aerobisoft.com/platform/pkg/endpoint"
	abservice "aerobisoft.com/platform/pkg/service"
	abhttp "aerobisoft.com/platform/pkg/transport/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

var logger log.Logger

func TestWiring(t *testing.T) {
	logger = log.NewLogfmtLogger(os.Stderr)

	// input empty middlewares
	svc := abservice.New([]abservice.Middleware{})
	eps := abendpoint.New(svc, map[string][]endpoint.Middleware{})
	// input empty server options
	options := map[string][]httptransport.ServerOption{}

	mux := abhttp.NewHTTPHandler(eps, options)
	srv := httptest.NewServer(mux)
	defer srv.Close()

	for _, testcase := range []struct {
		method, url, body, want string
	}{
		{"GET", srv.URL + "/greeting/John", "", `{"greeting":"GO-KIT Hello John"}`},
		{"GET", srv.URL + "/health", "", `{"healthy":true}`},
	} {
		req, _ := http.NewRequest(testcase.method, testcase.url, strings.NewReader(testcase.body))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Log(err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if want, have := testcase.want, strings.TrimSpace(string(body)); want != have {
			t.Errorf("%s %s %s: want %q, have %q", testcase.method, testcase.url, testcase.body, want, have)
		}
	}
}