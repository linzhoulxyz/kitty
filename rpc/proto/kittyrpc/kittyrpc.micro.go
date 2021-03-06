// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/kittyrpc/kittyrpc.proto

package go_micro_srv_kittyrpc

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Kittyrpc service

type KittyrpcService interface {
	Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
}

type kittyrpcService struct {
	c    client.Client
	name string
}

func NewKittyrpcService(name string, c client.Client) KittyrpcService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "go.micro.srv.kittyrpc"
	}
	return &kittyrpcService{
		c:    c,
		name: name,
	}
}

func (c *kittyrpcService) Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Kittyrpc.Call", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Kittyrpc service

type KittyrpcHandler interface {
	Call(context.Context, *Request, *Response) error
}

func RegisterKittyrpcHandler(s server.Server, hdlr KittyrpcHandler, opts ...server.HandlerOption) error {
	type kittyrpc interface {
		Call(ctx context.Context, in *Request, out *Response) error
	}
	type Kittyrpc struct {
		kittyrpc
	}
	h := &kittyrpcHandler{hdlr}
	return s.Handle(s.NewHandler(&Kittyrpc{h}, opts...))
}

type kittyrpcHandler struct {
	KittyrpcHandler
}

func (h *kittyrpcHandler) Call(ctx context.Context, in *Request, out *Response) error {
	return h.KittyrpcHandler.Call(ctx, in, out)
}
