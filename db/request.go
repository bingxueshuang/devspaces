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

func RequestsTo(to string) ([]Request, error) {
	r := make([]Request, 0, len(requests))
	for _, v := range requests {
		if v.To == to {
			r = append(r, *v)
		}
	}
	return r, nil
}
