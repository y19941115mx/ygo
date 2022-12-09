package main

import "ygo/framework"

func registerRoute(core *framework.Core) {
	core.Get("foo", FooController)
}
