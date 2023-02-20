package db

type Message struct {
	From    string
	To      string
	On      string
	Tag     string
	Data    []byte
	Keyword []byte
}

var msgs []*Message

func AddMessage(m *Message) (ok bool, err error) {
	msgs = append(msgs, m)
	return true, nil
}

func MessagesOn(tag string) ([]Message, error) {
	m := make([]Message, 0, len(msgs))
	for _, msg := range msgs {
		if msg.Tag == tag {
			m = append(m, *msg)
		}
	}
	return m, nil
}
