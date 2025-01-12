package main

import (
	"errors"
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

func getFlags() (*flags, error) {
	flg := &flags{}
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.StringVar(&flg.AppName, "appname", flg.AppName, "Application name")
	fs.StringVar(&flg.Environment, "environment", flg.Environment, "Environment")
	fs.StringVar(&flg.StructuredLogFile, "structured-log-file", flg.StructuredLogFile, "Structured log file")
	_ = fs.Parse(os.Args[1:]) // Ignore error, because it exits on error
	err := checkFlags(flg)
	if err != nil {
		return nil, err
	}
	return flg, nil
}

func checkFlags(flg *flags) error {
	if flg.Environment == "" {
		return errors.New("environment flag is required")
	}
	if flg.StructuredLogFile == "" {
		return errors.New("log file path is required")
	}
	return nil
}
