package main

import (
	root "custom-in-memory-db/cmd/client/cmd"
	conf2 "custom-in-memory-db/cmd/client/cmd/conf"
)

func main() {
	var conf conf2.Config
	conf2.InitConf()

	root.Execute(&conf)
}
