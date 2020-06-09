package server

import (
	"reflect"

	"github.com/micro/go-micro/v2/registry"
)

type rpcHandler struct {
	name      string
	handler   interface{}
	endpoints []*registry.Endpoint
	opts      HandlerOptions
}

func newRpcHandler(handler interface{}, opts ...HandlerOption) Handler {
	options := HandlerOptions{
		Metadata: make(map[string]map[string]string),
	}

	for _, o := range opts {
		o(&options)
	}

	typ := reflect.TypeOf(handler)
	hdlr := reflect.ValueOf(handler)
	//Indirect获取指向元素的值，如果为指针则返回指针指向的值，否则返回本值, Name()获取的名字不包含包名
	name := reflect.Indirect(hdlr).Type().Name()

	var endpoints []*registry.Endpoint //这里端点代表的是方法, exam 某个结构体包含的成员函数

	for m := 0; m < typ.NumMethod(); m++ {
		if e := extractEndpoint(typ.Method(m)); e != nil {
			e.Name = name + "." + e.Name

			for k, v := range options.Metadata[e.Name] {
				e.Metadata[k] = v
			}

			endpoints = append(endpoints, e)
		}
	}

	/*
		这里以handler作为struct分析，struct的名字作为handler名字，endpoints为struct所暴露的函数
		type Endpoint struct {
			Name     string            `json:"name"`
			Request  *Value            `json:"request"`
			Response *Value            `json:"response"`
			Metadata map[string]string `json:"metadata"`
		}
		Endpoint->name为函数名字 Request/Response 为参数的类型 如果参数是struct，则会存放struct的成员变量的名字/类型
	*/
	return &rpcHandler{
		name:      name,
		handler:   handler,
		endpoints: endpoints,
		opts:      options,
	}
}

func (r *rpcHandler) Name() string {
	return r.name
}

func (r *rpcHandler) Handler() interface{} {
	return r.handler
}

func (r *rpcHandler) Endpoints() []*registry.Endpoint {
	return r.endpoints
}

func (r *rpcHandler) Options() HandlerOptions {
	return r.opts
}
