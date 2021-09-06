package notification

import "encoding/json"

type CurOperation string

const (
	OpCreate = CurOperation("create")
	OpUpdate = CurOperation("update")
	OpDelete = CurOperation("delete")
)

func NewCurNotification(id uint64, op CurOperation) *CurNotification {
	return &CurNotification{
		ID:     id,
		Method: string(op),
	}
}

type CurNotification struct {
	ID     uint64 `json:"id"`
	Method string `json:"method"`

	encoded []byte
	err     error
}

func (cur *CurNotification) ensureEncoded() {
	if cur.encoded == nil && cur.err == nil {
		cur.encoded, cur.err = json.Marshal(cur)
	}
}

func (cur *CurNotification) Length() int {
	cur.ensureEncoded()
	return len(cur.encoded)
}

func (cur *CurNotification) Encode() ([]byte, error) {
	cur.ensureEncoded()
	return cur.encoded, cur.err
}
