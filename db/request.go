package db

type Request struct {
	From   string
	On     string
	To     string
	Secret []byte
}

var requests []*Request

func AddRequest(r *Request) (ok bool, err error) {
	requests = append(requests, r)
	return true, nil
}
