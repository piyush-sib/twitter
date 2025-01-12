package main

import (
	"flag"
	"os"

	"go.uber.org/dig"
)

type flags struct {
	dig.Out

	AppName           string `name:"appname"`
	Environment       string `name:"environment"`
	StructuredLogFile string `name:"structured-log-file"`
}

func newFlags() *flags {
	return &flags{
		Environment: "development", // Keeping Default environment as development
	}
}

func getFlags() *flags {
	flg := newFlags()
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.StringVar(&flg.AppName, "appname", flg.AppName, "Application name")
	fs.StringVar(&flg.Environment, "environment", flg.Environment, "Environment")
	fs.StringVar(&flg.StructuredLogFile, "structured-log-file", flg.StructuredLogFile, "Structured log file")
	_ = fs.Parse(os.Args[1:]) // Ignore error, because it exits on error
	return flg
}
