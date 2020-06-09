package main

import (
	"github.com/Electronic-Scavenger/ESGNIC/orm"
	"github.com/Electronic-Scavenger/ESGNIC/web"
)

func main() {
	orm.Init()
	var ch = make(chan struct{})
	web.Start()
	<-ch
}
