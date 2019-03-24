package oncrpc

import (
	"encoding/binary"
	"io"
	"net"
)

// Client onc rpc client
type Client struct {
	conn    net.Conn
	Address string
	Prog    uint32      // remote program
	Vers    uint32      // remote program version number
	Cred    *OpaqueAuth // authentication credential
	Verf    *OpaqueAuth // authentication verifier
}

// GetClient Get client
func GetClient(address string, prog uint32, vers uint32, cred *OpaqueAuth, verf *OpaqueAuth) *Client {

	client := &Client{Address: address, Prog: prog, Vers: vers, Cred: cred, Verf: verf}

	return client
}

// Connect create connection
func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.Address)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) writeFull(data []byte) error {

	dataSize := len(data)

	lastFragment := false

	startIndex := 0

	for !lastFragment {

		fragmentSize := MaxRecordFragmentSize
		if dataSize <= MaxRecordFragmentSize {
			lastFragment = true
			fragmentSize = dataSize
		}

		header := createFragmentHeader(uint32(fragmentSize), lastFragment)

		headerBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(headerBytes, header)

		written, err := c.conn.Write(append(headerBytes, data[startIndex:fragmentSize]...))
		if err != nil {
			return err
		}

		startIndex += written
	}

	return nil
}

// Call remote procedure
func (c *Client) Call(xid uint32, proc uint32, procParams interface{}) error {

	cbody := CallBody{
		RPCVers:    RPCVERS,
		Prog:       c.Prog,
		Vers:       c.Vers,
		Proc:       proc,
		Cred:       *c.Cred,
		Verf:       *c.Verf,
		ProcParams: procParams,
	}

	msg := &MsgCall{XID: xid, Type: MsgTypeCall, CBody: cbody}
	data, err := msg.Encode()
	if err != nil {
		return err
	}

	err = c.writeFull(data)
	if err != nil {
		return err
	}

	return nil
}

// ReadObj read data to an object
func (c *Client) ReadObj(obj interface{}) (*MsgReply, error) {
	// TODO:
	return nil, nil
}

// ReadHeader read ONC RPC header from conn
func (c *Client) Read() (msg *MsgReply, reader io.Reader, err error) {

	msg = &MsgReply{}

	var fragmentHeader uint32
	err = binary.Read(c.conn, binary.BigEndian, &fragmentHeader)
	if err != nil {
		return nil, nil, err
	}

	fragmentSize, lastFragment := getFragment(fragmentHeader)
	// XID
	err = binary.Read(c.conn, binary.BigEndian, &msg.XID)
	if err != nil {
		return nil, nil, err
	}
	// Type
	err = binary.Read(c.conn, binary.BigEndian, &msg.Type)
	if err != nil {
		return nil, nil, err
	}
	// reply stat
	err = binary.Read(c.conn, binary.BigEndian, &msg.RBody.Stat)
	if err != nil {
		return nil, nil, err
	}

	if msg.RBody.Stat == ReplyStatAccepted {
		// ReplyStatAccepted
		accepted := AcceptedReply{}

		// opaque_auth.flavorv
		err = binary.Read(c.conn, binary.BigEndian, &accepted.Verf.Flavor)
		if err != nil {
			return nil, nil, err
		}
		// opaque_auth.body
		var verFlavorLen uint32
		err = binary.Read(c.conn, binary.BigEndian, &verFlavorLen)
		if err != nil {
			return nil, nil, err
		}
		if verFlavorLen > 0 {
			accepted.Verf.Body = make([]byte, verFlavorLen)
			err = binary.Read(c.conn, binary.BigEndian, &accepted.Verf.Flavor)
			if err != nil {
				return nil, nil, err
			}
		}
		// stat
		err = binary.Read(c.conn, binary.BigEndian, &accepted.Stat)
		if err != nil {
			return nil, nil, err
		}

		msg.RBody.Areply = accepted

		switch accepted.Stat {
		case AcceptSuccess:
			// procedure-specific results
			arHeaderSize := acceptSuccessReplyFixedSize + verFlavorLen
			if fragmentSize > arHeaderSize {
				reader = &ReplyReader{conn: c.conn, fragmentSize: uint32(fragmentSize - arHeaderSize), lastFragment: lastFragment}
			}
		case AcceptProgMismatch:
			err = binary.Read(c.conn, binary.BigEndian, &accepted.MismatchInfo)
		default:
			// void
		}
		if err != nil {
			return nil, nil, err
		}
	} else {
		// ReplyStatDenied
		reject := RejectedReply{}

		// reject stat
		err = binary.Read(c.conn, binary.BigEndian, &reject.Stat)
		if err != nil {
			return nil, nil, err
		}

		switch reject.Stat {
		case RejectRPCMismatch:
			err = binary.Read(c.conn, binary.BigEndian, &reject.MismatchInfo)
		case RejectAuthError:
			err = binary.Read(c.conn, binary.BigEndian, &reject.AuthStat)
		}
		if err != nil {
			return nil, nil, err
		}

		msg.RBody.Rreply = reject
	}

	return msg, reader, nil
}
