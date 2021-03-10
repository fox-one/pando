package views

type Error struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

type JsonResult struct {
	Data interface{} `json:"data,omitempty"`
}
