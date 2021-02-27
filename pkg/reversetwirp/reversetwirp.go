package reversetwirp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/go-chi/chi"
	"github.com/oxtoacart/bpool"
)

type ReverseTwirp struct {
	Target http.Handler
	Path   string
}

type ParamTransfer func(key, value string) interface{}

func passThrough(_, value string) interface{} {
	return value
}

func (t *ReverseTwirp) Handle(method string, tr ParamTransfer) http.HandlerFunc {
	if tr == nil {
		tr = passThrough
	}

	bufferPool := bpool.NewBufferPool(64)

	return func(w http.ResponseWriter, r *http.Request) {
		body := make(map[string]interface{})
		if ctx := chi.RouteContext(r.Context()); ctx != nil {
			params := ctx.URLParams
			for idx, key := range params.Keys {
				body[key] = tr(key, params.Values[idx])
			}

			ctx.Reset()
		}

		for key, items := range r.URL.Query() {
			value := strings.Join(items, ",")
			body[key] = tr(key, value)
		}

		if len(body) > 0 {
			_ = json.NewDecoder(r.Body).Decode(&body)
			_ = r.Body.Close()

			b := bufferPool.Get()
			defer bufferPool.Put(b)

			_ = json.NewEncoder(b).Encode(body)
			r.Body = ioutil.NopCloser(b)
			r.ContentLength = int64(b.Len())
		}

		r.Method = http.MethodPost
		r.URL.RawQuery = ""
		r.URL.Path = path.Join(t.Path, method)
		r.Header.Set("Content-Type", "application/json")

		t.Target.ServeHTTP(w, r)
	}
}

type TwirpServer interface {
	http.Handler
	// PathPrefix returns the HTTP URL path prefix for all methods handled by this
	// service. This can be used with an HTTP mux to route twirp requests
	// alongside non-twirp requests on one HTTP listener.
	PathPrefix() string
}

func NewSingleTwirpServerProxy(svr TwirpServer) *ReverseTwirp {
	return &ReverseTwirp{
		Target: svr,
		Path:   svr.PathPrefix(),
	}
}
