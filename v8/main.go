package main

import (
	"rogchap.com/v8go"
)

var iso *v8go.Isolate

func init() {
	iso, _ = v8go.NewIsolate()
}
