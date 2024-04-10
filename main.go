package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shweta20039/shortme/conf"
	"github.com/shweta20039/shortme/short"
	"github.com/shweta20039/shortme/web"
)

func main() {
	cfgFile := flag.String("c", "config.conf", "configuration file")
	version := flag.Bool("v", false, "Version")

	flag.Parse()

	if *version {
		fmt.Println(conf.Version)
		os.Exit(0)
	}

	// parse config
	conf.MustParseConfig(*cfgFile)

	// short service
	short.Start()

	// api
	web.Start()
}
