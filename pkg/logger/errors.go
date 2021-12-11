package logger

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap/zapcore"
)

/*
	StatusCodeをhttp/netのStatus利用
	Go HTTPのステータスにマッピングすればいい
	stripe_converter.goが参考になる
*/

type ApplicationError struct {
	Type       ErrorType
	StatusCode int
	Level      zapcore.Level
	AppCode    string
	Message    string
	Detail     interface{}
}

/*
	ここをHTTP使用のエラー定義になる
*/

// Level type represents zapcore.Level.
type Level zapcore.Level

type ErrorType string

const (
	//
	ERR_MYSQL = ErrorType("MysqlError")

	// ERR_INTERNAL is errors that was occured internally.
	ERR_INTERNAL = ErrorType("ApplicationError")
)

// ApplicationError with Default Params.
func NewApplicationError(v interface{}) *ApplicationError {
	return &ApplicationError{
		Type:  ERR_INTERNAL,
		Level: ERROR,
		// StatusCode: 500,
		Detail: fmt.Sprint(v),
	}
}

// func "Error" returns message field in ApplicationError.
func (e *ApplicationError) Error() string {
	return e.Message
}

// Init
func (e *ApplicationError) Init(code string, message interface{}) *ApplicationError {
	// e.AppCode = e.Type.String() + code
	e.AppCode = code
	e.Message = fmt.Sprint(message)
	return e
}

// Error Code Handling
func GetApplicationError(err error) (appErr *ApplicationError) {
	// ex) MySQL PostgreSQL Cast
	switch err.(type) {
	case *mysql.MySQLError:

		mysqlErr := err.(*mysql.MySQLError)
		appErr = NewApplicationError(mysqlErr)
		appErr.AppCode = fmt.Sprintf("%v", mysqlErrorStatus(int(mysqlErr.Number)))
		appErr.Type = ERR_MYSQL
		tags := []*Tag{
			NewTag("Type", appErr.Type),
			NewTag("MySQLErrorCode", appErr.AppCode),
		}
		print(appErr.AppCode, Level(ERROR), mysqlErr.Error(), tags...)
		return
	case error:
		appErr = NewApplicationError(err.Error())
		// エラー種別関数　追加 mysqlErrorStatusのような
		appErr.StatusCode = http.StatusInternalServerError
		return
	default:
		return
	}
}

func mysqlErrorStatus(errNumber int) (statusCode int) {
	switch errNumber {
	case 1452:
		// 外部キーエラー
		return 1452
	case 1054:
		// 指定したカラムがあっていない
		return 1054
	case 1062:
		// PKが重複している
		return 1062
	case 1064:
		// SQL記述にエラー
		return 1064
	case 1146:
		// 指定したテーブル名が解決できない
		return 1146
	default:
		return http.StatusInternalServerError
	}
}

// Add original message.
func (e *ApplicationError) AddMessage(m string) *ApplicationError {
	e.Message = fmt.Sprintf("%v %v", m, e.Message)
	return e
}

type errorJson struct {
	ErrorMessage *messages `json:"errors"`
}

type messages struct {
	StatusCode int    `json:"stats_code"`
	Message    string `json:"message"`
}

func (e *ApplicationError) ErrorResponseJSON(w http.ResponseWriter, r *http.Request) *ApplicationError {
	msg := fmt.Sprintf(`%s url="%s" msg="%s"`, r.Method, r.URL, e.Detail)
	errJson, err := json.Marshal(&errorJson{
		ErrorMessage: &messages{
			StatusCode: e.StatusCode,
			Message:    msg,
		},
	})

	if err != nil {
		Error(
			GetApplicationError(err).AddMessage("An error has occured JsonMarshal "),
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(errJson)))
	w.WriteHeader(e.StatusCode)
	w.Write(errJson)
	return e
}
