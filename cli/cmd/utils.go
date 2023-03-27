package cmd

type Response struct {
	Ok    bool           `json:"ok"`
	Data  map[string]any `json:"data"`
	Error any            `json:"error"`
}
