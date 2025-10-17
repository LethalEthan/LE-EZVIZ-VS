package logging

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func init() {
	var err error
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{
		"stdout",
		// "./lez.log",
	}
	if Log, err = cfg.Build(); err != nil {
		panic(err)
	}
}

func CreateLogger(logFile, stdout bool) {
	var err error
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = make([]string, 0, 2)
	if stdout {
		cfg.OutputPaths = append(cfg.OutputPaths, "stdout")
	}
	if logFile {
		cfg.OutputPaths = append(cfg.OutputPaths, "./lez.log")
	}
	if Log, err = cfg.Build(); err != nil {
		panic(err)
	}
}
