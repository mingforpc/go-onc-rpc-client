package oncrpc

import (
	"bytes"
	"encoding/binary"
	"net"
)

// OpaqueAuth RPC struct
type OpaqueAuth struct {
	Flavor AuthFlavor
	Body   []byte
}

// AuthSysParms while Falovor is AuthSys
type AuthSysParms struct {
	Stamp       uint32
	MachineName string // max lenght 255
	UID         uint32
	GID         uint32
	GIDS        []int32 // max lenght 16
}

// MsgCall RPC msg call
type MsgCall struct {
	XID   uint32
	Type  MsgType  `xdr:"union"`
	CBody CallBody `xdr:"unioncase=0"`
	// RBody ReplyBody `xdr:"unioncase=1"`
}

// MsgReply RPC msg reply
type MsgReply struct {
	XID  uint32
	Type MsgType `xdr:"union"`
	// CBody CallBody `xdr:"unioncase=0"`
	RBody ReplyBody `xdr:"unioncase=1"`
}

// CallBody the call type of rpc msg
type CallBody struct {
	RPCVers uint32 // must be equal 2
	Prog    uint32
	Vers    uint32
	Proc    uint32
	Cred    OpaqueAuth // authentication credential
	Verf    OpaqueAuth // authentication verifier

	// procedure-specific parameters start here
	ProcParams interface{}
}

// ReplyBody the reply type of rpc msg
type ReplyBody struct {
	Stat   ReplyStat     `xdr:"union"`
	Areply AcceptedReply `xdr:"unioncase=0"`
	Rreply RejectedReply `xdr:"unioncase=1"`
}

// MisMatchInfo the info of MisMatch
type MisMatchInfo struct {
	Low  uint32
	High uint32
}

// AcceptedReply the reply type of rpc msg while accept
type AcceptedReply struct {
	Verf         OpaqueAuth
	Stat         AcceptStat   `xdr:"union"`
	MismatchInfo MisMatchInfo `xdr:"unioncase=2"`

	// procedure-specific results start here
}

// RejectedReply the reply type of rpc msg while reject
type RejectedReply struct {
	Stat         RejectStat   `xdr:"union"`
	MismatchInfo MisMatchInfo `xdr:"unioncase=0"`
	AuthStat     AuthStat     `xdr:"unioncase=1"`
}

// ReplyReader reply reader
type ReplyReader struct {
	conn         net.Conn
	fragmentSize uint32
	lastFragment bool
}

func (r *ReplyReader) Read(p []byte) (n int, err error) {
	lenght := len(p)

	if r.fragmentSize >= uint32(lenght) || r.lastFragment {

		n, err = r.conn.Read(p)
		r.fragmentSize = r.fragmentSize - uint32(n)

	} else {

		temp := make([]byte, r.fragmentSize)

		n, err = r.conn.Read(temp)
		buff := bytes.NewBuffer(temp)
		r.fragmentSize = r.fragmentSize - uint32(n)
		if err != nil && n <= 0 {
			return n, err
		}

		n, _ = buff.Read(p)

		if err == nil {
			var fragmentHeader uint32
			err2 := binary.Read(r.conn, binary.BigEndian, &fragmentHeader)
			if err2 == nil {
				r.fragmentSize, r.lastFragment = getFragment(fragmentHeader)
			}
		}
	}

	return n, err
}
