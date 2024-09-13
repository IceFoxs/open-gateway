package main

import (
	con "github.com/IceFoxs/open-gateway/server/consul"
	na "github.com/IceFoxs/open-gateway/server/nacos"
	"sync"
)

var (
	wg      sync.WaitGroup
	localIP = "your ip"
)

func main() {
	wg.Add(2)
	go func() {
		defer wg.Done()
		con.CreateConsulServer()
	}()
	go func() {
		defer wg.Done()
		na.CreateNacosServer()
	}()
	wg.Wait()
}
