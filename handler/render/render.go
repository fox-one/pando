package render

import (
	"encoding/json"
	"net/http"

	"github.com/twitchtv/twirp"
)

type H map[string]interface{}

func JSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	_ = enc.Encode(v)
}

func Text(w http.ResponseWriter, t string) {
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(t))
}

func Error(w http.ResponseWriter, err error) {
	_ = twirp.WriteError(w, err)
}

func BadRequest(w http.ResponseWriter, err error) {
	if _, ok := err.(twirp.Error); !ok {
		err = twirp.NewError(twirp.Malformed, err.Error())
	}

	Error(w, err)
}
