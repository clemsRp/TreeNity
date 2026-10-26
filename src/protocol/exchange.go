package protocol

import "encoding/json"

// Server Input structure
type ToServer struct {
	Command string          `json:"command"`
	Datas   json.RawMessage `json:"datas"`
	Path    string          `json:"path"`
}

// Server Output structure
type ToClient struct {
	Err      *Error          `json:"error"`
	Response string          `json:"response"`
	Datas    json.RawMessage `json:"datas"`
}
