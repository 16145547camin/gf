// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// 分组路由测试
package ghttp_test

import (
    "github.com/gogf/gf/g"
    "github.com/gogf/gf/g/frame/gmvc"
    "github.com/gogf/gf/g/net/ghttp"
    "github.com/gogf/gf/g/os/gtime"
    "github.com/gogf/gf/g/test/gtest"
    "testing"
    "time"
)

// 执行对象
type Object struct {}

func (o *Object) Index(r *ghttp.Request) {
    r.Response.Write("Object Index")
}

func (o *Object) Show(r *ghttp.Request) {
    r.Response.Write("Object Show")
}

func (o *Object) Delete(r *ghttp.Request) {
    r.Response.Write("Object REST Delete")
}

// 控制器
type Controller struct {
    gmvc.Controller
}

func (c *Controller) Index() {
    c.Response.Write("Controller Index")
}

func (c *Controller) Show() {
    c.Response.Write("Controller Show")
}

func (c *Controller) Post() {
    c.Response.Write("Controller REST Post")
}

func Handler(r *ghttp.Request) {
    r.Response.Write("Handler")
}

func Test_Router_Group1(t *testing.T) {
    s   := g.Server(gtime.Nanosecond())
    obj := new(Object)
    ctl := new(Controller)
    // 分组路由方法注册
    g := s.Group("/api")
    g.ALL ("/handler",     Handler)
    g.ALL ("/ctl",         ctl)
    g.GET ("/ctl/my-show", ctl, "Show")
    g.REST("/ctl/rest",    ctl)
    g.ALL ("/obj",         obj)
    g.GET ("/obj/my-show", obj, "Show")
    g.REST("/obj/rest",    obj)
    s.SetPort(8200)
    s.SetDumpRouteMap(false)
    go s.Run()
    defer func() {
        s.Shutdown()
        time.Sleep(time.Second)
    }()
    time.Sleep(time.Second)
    gtest.Case(t, func() {
        client := ghttp.NewClient()
        client.SetPrefix("http://127.0.0.1:8200")

        gtest.Assert(client.GetContent ("/api/handler"),     "Handler")

        gtest.Assert(client.GetContent ("/api/ctl"),         "Controller Index")
        gtest.Assert(client.GetContent ("/api/ctl/"),        "Controller Index")
        gtest.Assert(client.GetContent ("/api/ctl/index"),   "Controller Index")
        gtest.Assert(client.GetContent ("/api/ctl/my-show"), "Controller Show")
        gtest.Assert(client.GetContent ("/api/ctl/post"),    "Controller REST Post")
        gtest.Assert(client.GetContent ("/api/ctl/show"),    "Controller Show")
        gtest.Assert(client.PostContent("/api/ctl/rest"),    "Controller REST Post")

        gtest.Assert(client.GetContent ("/api/obj"),         "Object Index")
        gtest.Assert(client.GetContent ("/api/obj/"),        "Object Index")
        gtest.Assert(client.GetContent ("/api/obj/index"),   "Object Index")
        gtest.Assert(client.GetContent ("/api/obj/delete"),  "Object REST Delete")
        gtest.Assert(client.GetContent ("/api/obj/my-show"), "Object Show")
        gtest.Assert(client.GetContent ("/api/obj/show"),    "Object Show")
        gtest.Assert(client.DeleteContent("/api/obj/rest"),  "Object REST Delete")

        // 测试404
        resp, err := client.Get("/ThisDoesNotExist")
        defer resp.Close()
        gtest.Assert(err,             nil)
        gtest.Assert(resp.StatusCode, 404)
    })
}

func Test_Router_Group2(t *testing.T) {
    s   := g.Server(gtime.Nanosecond())
    obj := new(Object)
    ctl := new(Controller)
    // 分组路由批量注册
    s.Group("/api").Bind("/api", []ghttp.GroupItem{
        {"ALL",  "/handler",     Handler},
        {"ALL",  "/ctl",         ctl},
        {"GET",  "/ctl/my-show", ctl, "Show"},
        {"REST", "/ctl/rest",    ctl},
        {"ALL",  "/obj",         obj},
        {"GET",  "/obj/my-show", obj, "Show"},
        {"REST", "/obj/rest",    obj},
    })
    s.SetPort(8300)
    s.SetDumpRouteMap(false)
    go s.Run()
    defer func() {
        s.Shutdown()
        time.Sleep(time.Second)
    }()
    time.Sleep(time.Second)
    gtest.Case(t, func() {
        client := ghttp.NewClient()
        client.SetPrefix("http://127.0.0.1:8300")

        gtest.Assert(client.GetContent ("/api/handler"),     "Handler")

        gtest.Assert(client.GetContent ("/api/ctl/my-show"), "Controller Show")
        gtest.Assert(client.GetContent ("/api/ctl/post"),    "Controller REST Post")
        gtest.Assert(client.GetContent ("/api/ctl/show"),    "Controller Show")
        gtest.Assert(client.PostContent("/api/ctl/rest"),    "Controller REST Post")

        gtest.Assert(client.GetContent ("/api/obj/delete"),  "Object REST Delete")
        gtest.Assert(client.GetContent ("/api/obj/my-show"), "Object Show")
        gtest.Assert(client.GetContent ("/api/obj/show"),    "Object Show")
        gtest.Assert(client.DeleteContent("/api/obj/rest"),  "Object REST Delete")

        // 测试404
        resp, err := client.Get("/ThisDoesNotExist")
        defer resp.Close()
        gtest.Assert(err,             nil)
        gtest.Assert(resp.StatusCode, 404)
    })
}
