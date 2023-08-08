package logging

import (
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"google.golang.org/grpc/grpclog"
)

func CreateNewLogger(logger *zap.Logger, serviceName string) *otelzap.Logger {
	log := otelzap.New(logger.Named(serviceName), otelzap.WithTraceIDField(true))
	return log
}

type grpcLogger struct {
	l         *zap.SugaredLogger
	verbosity int
}

func (g *grpcLogger) Info(args ...interface{}) {
	g.l.Info(args...)
}

func (g *grpcLogger) Infoln(args ...interface{}) {
	g.l.Info(args...)
}

func (g *grpcLogger) Infof(format string, args ...interface{}) {
	g.l.Infof(format, args...)
}

func (g *grpcLogger) Warning(args ...interface{}) {
	g.l.Warn(args...)
}

func (g *grpcLogger) Warningln(args ...interface{}) {
	g.l.Warn(args...)
}

func (g *grpcLogger) Warningf(format string, args ...interface{}) {
	g.l.Warnf(format, args...)
}

func (g *grpcLogger) Error(args ...interface{}) {
	g.l.Error(args...)
}

func (g *grpcLogger) Errorln(args ...interface{}) {
	g.l.Error(args...)
}

func (g *grpcLogger) Errorf(format string, args ...interface{}) {
	g.l.Errorf(format, args...)
}

func (g *grpcLogger) Fatal(args ...interface{}) {
	g.l.Fatal(args...)
}

func (g *grpcLogger) Fatalln(args ...interface{}) {
	g.l.Fatal(args...)
}

func (g *grpcLogger) Fatalf(format string, args ...interface{}) {
	g.l.Fatalf(format, args...)
}

func (g *grpcLogger) V(l int) bool {
	return l <= g.verbosity
}

func SetInternalGRPCLogger(log *zap.Logger) {
	logger := &grpcLogger{
		l:         log.Sugar(),
		verbosity: 2,
	}
	grpclog.SetLoggerV2(logger)
}