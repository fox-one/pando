package call

import (
	"context"
	"encoding/json"

	"github.com/fox-one/pando/cmd/pando-cli/internal/cfg"
	"github.com/go-resty/resty/v2"
)

var client = resty.New()

func R(ctx context.Context) *resty.Request {
	r := client.SetHostURL(cfg.GetApiHost()).NewRequest()
	r = r.SetContext(ctx)
	r = r.SetAuthToken(cfg.GetAuthToken())

	return r
}

func DecodeResponse(r *resty.Response) ([]byte, error) {
	var body struct {
		Error `json:"error,omitempty"`
		Data  json.RawMessage `json:"data,omitempty"`
	}

	if err := json.Unmarshal(r.Body(), &body); err != nil {
		return nil, err
	}

	if body.Error.Code > 0 {
		return nil, body.Error
	}

	return body.Data, nil
}

func UnmarshalResponse(r *resty.Response, v interface{}) error {
	data, err := DecodeResponse(r)
	if err != nil {
		return err
	}

	if v != nil {
		return json.Unmarshal(data, v)
	}

	return nil
}
