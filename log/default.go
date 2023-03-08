package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func DefaultEncoderConfig() zapcore.EncoderConfig {
	var encoderConfig = zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   //时间编码
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder //使用大写字母记录日志级别
	return encoderConfig
}

//默认日志切割配置
func DefaultLumberJackLoggger() *lumberjack.Logger {
	return &lumberjack.Logger{
		MaxSize:   5, //切割之前，日志最大大小（以MB为单位）
		LocalTime: true,
		Compress:  true, //是否压缩旧文件
		//Filename 日志文件的位置
		// MaxAge 保留旧文件的最大天数
		//MaxBackups 保留旧文件的最大个数
	}
}

func DefaultOptions() []zap.Option {
	var stackTraceLevel zap.LevelEnablerFunc = func(l zapcore.Level) bool {
		return l >= zap.DPanicLevel //当日志级别大于等于DPanic级别才输出堆栈信息
	}
	return []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(stackTraceLevel),
	}
}
