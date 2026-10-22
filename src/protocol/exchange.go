package protocol

// Server Input structure
type ToServer struct {
	Command string `json:"command"`
	Datas   any    `json:"datas"`
}

// Server Output structure
type ToClient struct {
	Err      error  `json:"error"`
	Response string `json:"response"`
	Datas    any    `json:"datas"`
}
