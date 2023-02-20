package db

type Message struct {
	From string
	To   string
	On   string
	Tag  string
}

var msgs []*Message

func MessagesOn(on string) ([]Message, error) {
	m := make([]Message, 0, len(msgs))
	for _, msg := range msgs {
		if msg.On == on {
			m = append(m, *msg)
		}
	}
	return m, nil
}

func MessagesTo(to string) ([]Message, error) {
	m := make([]Message, 0, len(msgs))
	for _, msg := range msgs {
		if msg.To == to {
			m = append(m, *msg)
		}
	}
	return m, nil
}
