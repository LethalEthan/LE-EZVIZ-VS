package logging

import "go.uber.org/zap"

var Log *zap.Logger

func init() {
	var err error
	if Log, err = zap.NewDevelopment(); err != nil {
		panic(err)
	}
}
