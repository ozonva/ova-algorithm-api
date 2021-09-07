package notification

import "encoding/json"

// CurOperation is a notification type sent to kafka
type CurOperation string

const (
	OpCreate = CurOperation("create")
	OpUpdate = CurOperation("update")
	OpDelete = CurOperation("delete")
)

// NewCurNotification creates new notification for sarama kafka client
func NewCurNotification(id uint64, op CurOperation) *CurNotification {
	return &CurNotification{
		id:     id,
		method: string(op),
	}
}

// CurNotification a notification for sarama kafka client.
// Should be created with NewCurNotification
type CurNotification struct {
	id     uint64 `json:"id"`
	method string `json:"method"`

	encoded []byte
	err     error
}

func (cur *CurNotification) ensureEncoded() {
	if cur.encoded == nil && cur.err == nil {
		cur.encoded, cur.err = json.Marshal(cur)
	}
}

// Length of encoded data. It's required for sarama kafka client
func (cur *CurNotification) Length() int {
	cur.ensureEncoded()
	return len(cur.encoded)
}

// Encode returns encoded data. Its required for sarama kafka client
func (cur *CurNotification) Encode() ([]byte, error) {
	cur.ensureEncoded()
	return cur.encoded, cur.err
}
