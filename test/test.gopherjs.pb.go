package test

import "github.com/gopherjs/gopherjs/js"

type MyMessage struct {
    *js.Object
    Msg string `js:"msg"`
    Num uint32 `js:"num"`
}
