package cmd

type Response struct {
	Ok    bool `json:"ok"`
	Data  any  `json:"data"`
	Error any  `json:"error"`
}
