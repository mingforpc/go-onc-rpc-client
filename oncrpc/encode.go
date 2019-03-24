package oncrpc

import (
	"bytes"

	xdr "github.com/davecgh/go-xdr/xdr2"
)

// Encode encode to xdr bytes
func (m *MsgCall) Encode() ([]byte, error) {

	writer := bytes.NewBuffer(nil)
	bytesWritten, err := xdr.Marshal(writer, m)
	if err != nil {
		return nil, err
	}

	return writer.Bytes()[:bytesWritten], nil
}
