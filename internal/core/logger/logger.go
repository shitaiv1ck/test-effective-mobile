package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func NewLogger(config Config) (*Logger, error) {
	zaplvl := zap.NewAtomicLevel()
	if err := zaplvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, err
	}

	zapcfg := zap.NewDevelopmentEncoderConfig()
	zapcfg.EncodeTime = zapcore.ISO8601TimeEncoder

	zapEncoder := zapcore.NewConsoleEncoder(zapcfg)

	core := zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zaplvl)

	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{
		Logger: zapLogger,
	}, nil
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
	}
}

func FromContext(ctx context.Context) *Logger {
	return ctx.Value("logger").(*Logger)
}
