package logger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	_ "time/tzdata"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

type Logger struct {
	ctxTags   grpc.UnaryServerInterceptor
	zap       grpc.UnaryServerInterceptor
	ZapLogger *zap.Logger
}

type AppType string

type ApplicationResponse struct {
	Type       AppType
	StatusCode int
	Level      zapcore.Level
	AppCode    string
	Message    string
	Detail     interface{}
}

const (
	// DEBUG
	DEBUG = zapcore.DebugLevel
	// INFO
	INFO = zapcore.InfoLevel
	// WARN
	WARN = zapcore.WarnLevel
	// ERROR
	ERROR = zapcore.ErrorLevel
	// FATAL
	FATAL = zapcore.FatalLevel
)

/*
	Zap Logger
*/
func NewLogger() error {
	loc, err := time.LoadLocation(os.Getenv("TZ"))
	if err != nil {
		return err
	}
	time.Local = loc
	l := &Logger{}
	if err := l.getLogger(); err != nil {
		return err
	}

	Info("00", "Logger initialized successfully.")
	return nil
}

func (l *Logger) getLogger() error {
	env := "production"
	if env == "production" {
		z, err := zap.NewProduction()
		if err != nil {
			return NewApplicationError(err).Init("21", "Failed to initialize logger.")
		}
		z.Sync()
		l.ZapLogger = z
		zap.ReplaceGlobals(l.ZapLogger)
	} else {
		z, err := zap.NewDevelopment()
		if err != nil {
			return NewApplicationError(err).Init("21", "Failed to initialize logger.")
		}
		l.ZapLogger = z
		zap.ReplaceGlobals(l.ZapLogger)
	}
	return nil
}

type ResponseJson struct {
	Response []*ResponseRequest `json:"response"`
}

type Header struct {
	Status string
}

type ResponseRequest struct {
	Code    int    `json:"Code"`
	Message string `json:"Message"`
}

func ResponseJSON(w http.ResponseWriter, object interface{}, appResp *ApplicationResponse, tags ...*Tag) {
	resp, err := json.Marshal(object)
	if err != nil {
		Error(
			GetApplicationError(err).AddMessage("An error has occured JsonMarshal "),
			tags...,
		) // 暫定
		// ErrorResponseのJsonつくらないといけない
	}

	print(appResp.AppCode, Level(appResp.Level), appResp.Message, tags...)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(resp)))
	w.WriteHeader(appResp.StatusCode)
	w.Write(resp)
}

func print(appCode string, level Level, message string, tags ...*Tag) {
	// 簡易
	// zap.L().Info(message, zap.String("key", "value"), zap.Time("now", time.Now()))
	serviceField := zap.String("service", "WebApplication")
	var fields []zapcore.Field
	if tags != nil {
		fields = make([]zapcore.Field, len(tags)+2)
		fields[0], fields[1] = serviceField, zap.String("code", appCode)
		for i, tag := range tags {
			fields[i+2] = tag.Log()
		}
	} else {
		fields := make([]zapcore.Field, 2)
		fields[0], fields[1] = serviceField, zap.String("code", appCode)
	}
	switch zapcore.Level(level) {
	case zapcore.DebugLevel:
		zap.L().Debug(message, fields...)
	case zapcore.InfoLevel:
		zap.L().Info(message, fields...)
	case zapcore.ErrorLevel:
		zap.L().Error(message, fields...)
	case zapcore.FatalLevel:
		zap.L().Fatal(message, fields...)
	}
}

func Info(code string, message string, tags ...*Tag) {
	appCode := "WebAppEnv " + code
	print(appCode, Level(INFO), message, tags...)
}

func Fatal(appErr *ApplicationError, tags ...*Tag) {
	var message string
	if appErr.Detail != "" {
		message = fmt.Sprintf("%s, ( DETAILS: %s)", appErr.Message, appErr.Detail)
	} else {
		message = appErr.Message
	}
	print(appErr.AppCode, Level(FATAL), message, tags...)
}

func Error(appErr *ApplicationError, tags ...*Tag) {
	var message string

	if appErr.Detail != "" {
		message = fmt.Sprintf("%s, ( DETAILS: %s )", appErr.Message, appErr.Detail)
	} else {
		message = appErr.Message
	}
	print(appErr.AppCode, Level(ERROR), message, tags...)
}

func Debug(code string, message string, tags ...*Tag) {
	appCode := "WebAppEnv " + code
	print(appCode, Level(DEBUG), message, tags...)
}
