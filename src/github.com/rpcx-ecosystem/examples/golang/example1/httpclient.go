package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/smallnest/rpcx/codec"
	"github.com/rpcx-ecosystem/rpcx-gateway"
)

type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}

func main() {
	//编码
	cc := &codec.MsgpackCodec{}

	args := &Args{
		A: 10,
		B: 20,
	}
	//编码
	data, _ := cc.Encode(args)

	//请求
	req, err := http.NewRequest("POST",
		"http://localhost:8972/",
		bytes.NewReader(data))
	if err != nil {
		log.Fatal("failed to create request: ", err)
		return
	}

	h := req.Header
	h.Set(gateway.XMessageID, "10000")
	h.Set(gateway.XMessageType, "0")
	h.Set(gateway.XSerializeType, "3")
	h.Set(gateway.XServicePath, "Arith")
	h.Set(gateway.XServiceMethod, "Mul")
	//发送请求
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("failed to call: ", err)
	}
	defer res.Body.Close()

	// handle http response
	// 处理响应
	replyData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("failed to read response: ", err)
	}

	reply := &Reply{}
	//解码
	err = cc.Decode(replyData, reply)
	if err != nil {
		log.Fatal("failed to decode reply: ", err)
	}

	log.Printf("%d * %d = %d", args.A, args.B, reply.C)
}
