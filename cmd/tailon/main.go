package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fulldump/goconfig"

	"github.com/fulldump/tailon/api"
)

var VERSION = "dev"

type Config struct {
	HttpAddr string `usage:"Service address"`
	Statics  string `usage:"statics directory or http address"`
	Version  bool   `usage:"Show version and exit"`
}

func main() {

	c := &Config{
		HttpAddr: ":8080",
	}
	goconfig.Read(c)

	if c.Version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	b := api.Build(VERSION, c.Statics)

	b.WithInterceptors(
		api.AccessLog(log.Default()),
		api.RecoverFromPanic,
		api.PrettyErrorInterceptor,
	)

	s := &http.Server{
		Addr:    c.HttpAddr,
		Handler: b,
	}

	fmt.Println("Server listening on", s.Addr)
	s.ListenAndServe()
}
