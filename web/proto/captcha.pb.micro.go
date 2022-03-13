// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/captcha.proto

package captcha

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/asim/go-micro/v3/api"
	client "github.com/asim/go-micro/v3/client"
	server "github.com/asim/go-micro/v3/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Captcha service

func NewCaptchaEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Captcha service

type CaptchaService interface {
	Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
}

type captchaService struct {
	c    client.Client
	name string
}

func NewCaptchaService(name string, c client.Client) CaptchaService {
	return &captchaService{
		c:    c,
		name: name,
	}
}

func (c *captchaService) Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Captcha.Call", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Captcha service

type CaptchaHandler interface {
	Call(context.Context, *Request, *Response) error
}

func RegisterCaptchaHandler(s server.Server, hdlr CaptchaHandler, opts ...server.HandlerOption) error {
	type captcha interface {
		Call(ctx context.Context, in *Request, out *Response) error
	}
	type Captcha struct {
		captcha
	}
	h := &captchaHandler{hdlr}
	return s.Handle(s.NewHandler(&Captcha{h}, opts...))
}

type captchaHandler struct {
	CaptchaHandler
}

func (h *captchaHandler) Call(ctx context.Context, in *Request, out *Response) error {
	return h.CaptchaHandler.Call(ctx, in, out)
}