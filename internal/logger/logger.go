package logger

import (
	"fmt"
	"go.uber.org/zap"
)

func NewZapLoggerForEnv(env string, callerSkip int) (*zap.SugaredLogger, error) {
	if env == "prod" {
		logger, err := zap.NewProduction(zap.AddCallerSkip(callerSkip), zap.AddStacktrace(zap.ErrorLevel))
		return logger.Sugar(), err
	}

	logger, err := zap.NewDevelopment(zap.AddCallerSkip(callerSkip))
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println(logger)
	//return logger.Sugar(), err
	return nil, nil

}
