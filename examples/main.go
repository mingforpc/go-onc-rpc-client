package main

import (
	"bytes"
	"fmt"

	"github.com/mingforpc/go-onc-rpc-client/oncrpc"
	"github.com/rasky/go-xdr/xdr2"
)

type Cred struct {
	Stamp uint32
	MName []byte
	Uid   uint32
	Gid   uint32
	Gids  []uint32
}

type MounthService struct {
	Path []byte
}

func main() {

	// call MOUNT

	c := Cred{Stamp: 0x010716a5, MName: []byte("ming-MI"), Uid: 0, Gid: 0, Gids: []uint32{0}}
	writer := bytes.NewBuffer(nil)
	_, _ = xdr.Marshal(writer, c)

	ms := MounthService{Path: []byte("/volume1/test")}

	cred := &oncrpc.OpaqueAuth{Flavor: oncrpc.AuthFlavorAuthSys, Body: writer.Bytes()}
	verf := &oncrpc.OpaqueAuth{Flavor: oncrpc.AuthFlavorNone}
	client := oncrpc.GetClient("192.168.50.99:892", 100005, 3, cred, verf)

	err := client.Connect()
	fmt.Printf("Connect err:%+v \n", err)

	err = client.Call(2718113503, 1, ms)
	fmt.Printf("call err:%+v \n", err)

	msg, reader, err := client.Read()
	if err != nil {
		fmt.Printf("read err:%+v \n", err)
	}

	fmt.Printf("msg: %+v \n", msg)

	if reader != nil {
		fmt.Printf("reader:%+v \n", reader)
		a := make([]byte, 100)
		n, _ := reader.Read(a)
		fmt.Printf("% x \n", a[:n])
	}

}
