package internal

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

// Zap 配置
type Zap struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"`                            // 级别
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                         // 日志前缀
	Format        string `mapstructure:"format" json:"format" yaml:"format"`                         // 输出
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`                  // 日志文件夹
	EncodeLevel   string `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"`       // 编码级
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"` // 栈名
	MaxAge        int    `mapstructure:"max-age" json:"max-age" yaml:"max-age"`                      // 日志留存时间
	ShowLine      bool   `mapstructure:"show-line" json:"show-line" yaml:"show-line"`                // 显示行
	LogInConsole  bool   `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console"` // 输出控制台
}

// zapEncodeLevel 根据配置文件中的encode-level解析出对应的LevelEncoder
func (z *Zap) zapEncodeLevel() zapcore.LevelEncoder {
	switch {
	case z.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	case z.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case z.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case z.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

// transportLevel 根据配置文件中的level解析出对应的Level
func (z *Zap) transportLevel() zapcore.Level {
	z.Level = strings.ToLower(z.Level)
	switch z.Level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.WarnLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

// getLevelPriority 通过zapcore.Level 获取 zap.LevelEnablerFunc
func (z *Zap) getLevelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	switch level {
	case zapcore.DebugLevel:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	case zapcore.InfoLevel:
		return func(level zapcore.Level) bool { // 日志级别
			return level == zap.InfoLevel
		}
	case zapcore.WarnLevel:
		return func(level zapcore.Level) bool { // 警告级别
			return level == zap.WarnLevel
		}
	case zapcore.ErrorLevel:
		return func(level zapcore.Level) bool { // 错误级别
			return level == zap.ErrorLevel
		}
	case zapcore.DPanicLevel:
		return func(level zapcore.Level) bool { // dpanic级别
			return level == zap.DPanicLevel
		}
	case zapcore.PanicLevel:
		return func(level zapcore.Level) bool { // panic级别
			return level == zap.PanicLevel
		}
	case zapcore.FatalLevel:
		return func(level zapcore.Level) bool { // 终止级别
			return level == zap.FatalLevel
		}
	default:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	}
}

// getEncoderConfig 获取日志编码格式
func (z *Zap) getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: z.StacktraceKey,
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   z.zapEncodeLevel(),
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(z.Prefix + t.Format("2006/01/02 - 15:04:05.000"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

// getEncoder 获取编码器
func (z *Zap) getEncoder() zapcore.Encoder {
	if z.Format == "json" {
		return zapcore.NewJSONEncoder(z.getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(z.getEncoderConfig())
}

// getEncoderCore 获取编码器核心
func (z *Zap) getEncoderCore(l zapcore.Level, level zap.LevelEnablerFunc) zapcore.Core {
	writer, err := FileRotateLogs.GetWriteSyncer(z, l.String()) // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("获取日志切割器失败:%v", err.Error())
		return nil
	}
	return zapcore.NewCore(z.getEncoder(), writer, level)
}

// getZapCores 获取日志编码器核心切片
func (z *Zap) getZapCores() []zapcore.Core {
	cores := make([]zapcore.Core, 0, 7)
	for level := z.transportLevel(); level <= zapcore.FatalLevel; level++ {
		cores = append(cores, z.getEncoderCore(level, z.getLevelPriority(level)))
	}
	return cores
}

// NewZap NewZap
func (z *Zap) NewZap() *zap.Logger {
	cores := z.getZapCores()                    // 获取 zap 核心切片
	logger := zap.New(zapcore.NewTee(cores...)) // 初始化 zap.Logger
	if z.ShowLine {                             // 判断是否显示行
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}
