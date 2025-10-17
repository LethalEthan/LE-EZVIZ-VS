package logging

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func init() {
	var err error
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{
		// "stdout",
		"./lez.log",
	}
	if Log, err = cfg.Build(); err != nil {
		panic(err)
	}
}
