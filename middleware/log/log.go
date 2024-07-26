package log

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"gin.example.com/entity/config"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugarLogger *zap.SugaredLogger
var apiLogger *zap.SugaredLogger

func InitLogger(config *config.LogConfig) {
	core := zapcore.NewTee(
		zapcore.NewCore(getEncoder("console"), zapcore.AddSync(os.Stdout), getDebugLevel(config.Level)), // console 打印
		zapcore.NewCore(getEncoder(""), getLogWriter(config), getDebugLevel(config.Level)),              // 写入文件
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)) // 输出文件名和行号
	sugarLogger = logger.Sugar()
	apiLogger = zap.New(core).Sugar()
}

func getDebugLevel(level string) zapcore.LevelEnabler {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	}
	return zapcore.InfoLevel
}

func getEncoder(typ string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000") // 时间格式设置
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder                           // 日志级别大写
	if typ == "console" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	encoderConfig.ConsoleSeparator = "|" // 分隔符
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(config *config.LogConfig) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    int(config.MaxSize),
		MaxBackups: int(config.MaxBackups),
		MaxAge:     int(config.MaxAge),
		Compress:   false, // 是否压缩/归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger)
}

func Debug(args ...interface{}) {
	sugarLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	sugarLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	sugarLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	sugarLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	sugarLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	sugarLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	sugarLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	sugarLogger.Errorf(template, args...)
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

		param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path
		apiLogger.Infof(ginDefaultformatter(param))
	}
}

func ginDefaultformatter(param gin.LogFormatterParams) string {
	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("API|%s|%s|%3d|%v|%s|%s",
		param.Path,
		param.Method,
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.ErrorMessage,
	)
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					apiLogger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					apiLogger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					apiLogger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
