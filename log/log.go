package log

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(core zapcore.Core, options ...zap.Option) *zap.Logger {
	return zap.New(core, append(DefaultOptions(), options...)...)
}

func NewCore(writer zapcore.WriteSyncer, enabler zapcore.LevelEnabler) zapcore.Core {
	return zapcore.NewCore(zapcore.NewJSONEncoder(DefaultEncoderConfig()), writer, enabler)
}

func NewStderrCore(enable zapcore.LevelEnabler) zapcore.Core {
	return NewCore(zapcore.Lock(zapcore.AddSync(os.Stderr)), enable)
}

func NewStdoutCore(enable zapcore.LevelEnabler) zapcore.Core {
	return NewCore(zapcore.Lock(zapcore.AddSync(os.Stdout)), enable)
}

func NewFileCore(filePath string, enabler zapcore.LevelEnabler) (zapcore.Core, io.Closer) {
	var writer = DefaultLumberJackLoggger()
	writer.Filename = filePath
	return NewCore(zapcore.AddSync(writer), enabler), writer
}
