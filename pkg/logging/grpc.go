package logging

import (
	"context"
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func NewGRPCOption( /* *GRPCLogCfg */ ) ([]grpc_logging.Option, error) {

	var logOpts []grpc_logging.Option

	logOpts = []grpc_logging.Option{
		grpc_logging.WithLevels(DefaultCodeToLevelGRPC),
	}

	globalLevel, globalStart, globalEnd, err := fillGlobalOptionConfig( /* *GRPCLogCfg */ )

	if err != nil {
		return logOpts, err
	}
	if !globalStart && !globalEnd {
		//Nothing to do
		return nil, nil
	}

	if err := validateLevel(globalLevel); err != nil {
		return logOpts, err
	}

	logOpts = []grpc_logging.Option{
		grpc_logging.WithLevels(DefaultCodeToLevelGRPC),
	}

	reqLogDecision, err := getGRPCLoggingOption(globalStart, globalEnd)
	if err != nil {
		return logOpts, err
	}

	logOpts = append(logOpts, reqLogDecision)

	return logOpts, nil
}

func fillGlobalOptionConfig( /* *GRPCLogCfg */ ) (string, bool, bool, error) {
	globalLevel := "ERROR"

	globalStart := true
	globalEnd := true

	return globalLevel, globalStart, globalEnd, nil
}

func validateLevel(level string) error {
	if level == "" {
		return fmt.Errorf("level field in YAML file is empty")
	}
	if level == "INFO" || level == "DEBUG" || level == "ERROR" || level == "WARNING" {
		return nil
	}
	return fmt.Errorf("the format of level is invalid. Expected INFO/DEBUG/ERROR/WARNING, got this %v", level)
}

func getGRPCLoggingOption(logStart, logEnd bool) (grpc_logging.Option, error) {
	if !logStart && !logEnd {
		return grpc_logging.WithLogOnEvents(), nil
	}
	if !logStart && logEnd {
		return grpc_logging.WithLogOnEvents(grpc_logging.FinishCall), nil
	}
	if logStart && logEnd {
		return grpc_logging.WithLogOnEvents(grpc_logging.StartCall, grpc_logging.FinishCall), nil
	}
	return nil, fmt.Errorf("log decision combination is not supported")
}

func InterceptorLogger(l log.Logger) grpc_logging.Logger {
	return grpc_logging.LoggerFunc(func(_ context.Context, lvl grpc_logging.Level, msg string, fields ...any) {
		largs := append([]any{"msg", msg}, fields...)
		switch lvl {
		case grpc_logging.LevelDebug:
			_ = level.Debug(l).Log(largs...)
		case grpc_logging.LevelInfo:
			_ = level.Info(l).Log(largs...)
		case grpc_logging.LevelWarn:
			_ = level.Warn(l).Log(largs...)
		case grpc_logging.LevelError:
			_ = level.Error(l).Log(largs...)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
