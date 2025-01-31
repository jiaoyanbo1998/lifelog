package loggerx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

// ZapLogger Logger接口的适配类 --- zap框架实现
type ZapLogger struct {
	// zap.Logger zap框架的核心日志记录器
	zapLogger *zap.Logger
}

// ZapNoLogger 不执行的zap
type ZapNoLogger struct {
	zapLogger *zap.Logger
}

// NewZapNoLogger 创建一个“不执行的zap”
func NewZapNoLogger() *ZapLogger {
	return &ZapLogger{
		zapLogger: zap.NewNop(),
	}
}

// NewZapLogger 使用构造函数初始化ZapLogger
func NewZapLogger() *ZapLogger {
	// 构建日志核心组件，支持同时输出到文件和控制台
	core := zapcore.NewTee(
		// 输出到控制台
		zapcore.NewCore(
			getConsoleEncoder(),        // 控制台日志编码器
			zapcore.AddSync(os.Stdout), // 输出到控制台
			zapcore.DebugLevel), // 日志级别为Debug，输出比debug高的所有级别日志
		// 输出到文件
		zapcore.NewCore(
			getJSONEncoder(), // json格式日志编码器
			getLogWriter(),   // 输出到文件
			zapcore.DebugLevel),
	)
	logger := zap.New(core, zap.AddCaller()) // 启用调用者信息
	return &ZapLogger{
		zapLogger: logger,
	}
}

// 将[]Field转为zap.Field
func (z *ZapLogger) toZapField(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}
	return zapFields
}

// Debug 记录Debug级别的日志
func (z *ZapLogger) Debug(msg string, args ...Field) {
	z.zapLogger.Debug(msg, z.toZapField(args)...)
}

// Info 记录Info级别的日志
func (z *ZapLogger) Info(msg string, args ...Field) {
	z.zapLogger.Info(msg, z.toZapField(args)...)
}

// Warn 记录Warn级别的日志
func (z *ZapLogger) Warn(msg string, args ...Field) {
	z.zapLogger.Warn(msg, z.toZapField(args)...)
}

// Error 记录Error级别的日志
func (z *ZapLogger) Error(msg string, args ...Field) {
	z.zapLogger.Error(msg, z.toZapField(args)...)
}

// getConsoleEncoder 控制台日志编码器
func getConsoleEncoder() zapcore.Encoder {
	// 开发模式
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	// 设置时间格式
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	// 小写带颜色的日志级别
	encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getJSONEncoder json日志编码器
func getJSONEncoder() zapcore.Encoder {
	// 开发模式
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	// 设置时间键为"time"
	encoderConfig.TimeKey = "time"
	// 设置时间格式
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	// 小写日志级别
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getLogWriter 同步日志记录器
func getLogWriter() zapcore.WriteSyncer {
	// 分片
	lumberLogger := &lumberjack.Logger{
		Filename:   "./log/lifeLog.log", // 日志文件路径
		MaxSize:    10,                  // 日志文件最大体积（MB）
		MaxAge:     30,                  // 日志文件最大保存天数
		MaxBackups: 5,                   // 最大保留的日志文件数量
	}
	return zapcore.AddSync(lumberLogger)
}

// getLogWriterSync 异步日志记录器
func getLogWriterSync() *zapcore.BufferedWriteSyncer {
	// 分片
	lumberLogger := &lumberjack.Logger{
		Filename:   "./log/lifeLog.log", // 日志文件路径
		MaxSize:    10,                  // 日志文件最大体积（MB）
		MaxAge:     30,                  // 日志文件最大保存天数
		MaxBackups: 5,                   // 最大保留的日志文件数量
	}
	return &zapcore.BufferedWriteSyncer{
		WS:   zapcore.AddSync(lumberLogger),
		Size: 4096,
	}
}
