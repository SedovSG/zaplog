package logger

import (
	"encoding/json"
	"github.com/SedovSG/zaplog/models"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// InitZapLogger настаривает конфигурацию логера
func InitZapLogger() (*zap.Logger, error) {
	var mode string
	mode, ok := os.LookupEnv("MODE")
	if !ok {
		mode = "DEV"
	}

	var config = zap.Config{}

	if mode == "PROD" {
		config = zap.NewProductionConfig()
		config.DisableStacktrace = true
	} else {
		config = zap.NewDevelopmentConfig()
		config.Development = true
	}

	config.Encoding = "json"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.NameKey = "name"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.FunctionKey = "function"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.EncoderConfig.StacktraceKey = "stacktrace"
	config.EncoderConfig.TimeKey = ""

	logger, error := config.Build()
	if error != nil {
		log.Fatal(error)
	}

	return logger, error
}

// AddFields добавляет поля в лог
func AddFields(fields []interface{}, logger *zap.Logger) *zap.Logger {
	var (
		method   string
		url      string
		code     int
		request  []byte
		response []byte
	)

	extend := make([]interface{}, 0, len(fields))

	methods := []string{"GET", "PUT", "POST", "HEAD", "DELETE", "OPTIONS"}

	for _, val := range fields {
		if _, ok := val.(int); ok {
			code = val.(int)
		}

		if _, ok := val.(string); ok {
			if contains(methods, val.(string)) {
				method = val.(string)
			} else if strings.Contains(val.(string), "/") {
				url = val.(string)
			}
		}

		switch val.(type) {
		case *http.Request:
			result := val.(*http.Request)

			var body []byte

			if result.Body != nil {
				body, _ = ioutil.ReadAll(result.Body)
				defer result.Body.Close()
			}

			if method == "" {
				method = result.Method
			}

			if url == "" {
				url = result.URL.String()
			}

			request, _ = json.Marshal(models.Request{
				Method:     result.Method,
				Proto:      result.Proto,
				URL:        result.URL.String(),
				RequestURI: result.RequestURI,
				Body:       string(body),
				Header:     result.Header,
				Cookies:    result.Cookies(),
				RemoteAddr: result.RemoteAddr,
				UserAgent:  result.UserAgent(),
				Referer:    result.Referer(),
			})
		case *http.Response:
			result := val.(*http.Response)

			var body []byte
			var location string

			if result.Body != nil {
				body, _ = ioutil.ReadAll(result.Body)
				defer result.Body.Close()
			}

			if loc, _ := result.Location(); loc != nil {
				location = loc.String()
			}

			if code == 0 {
				code = result.StatusCode
			}

			response, _ = json.Marshal(models.Response{
				Status:     result.Status,
				StatusCode: result.StatusCode,
				Proto:      result.Proto,
				Location:   location,
				Body:       string(body),
				Header:     result.Header,
				Cookies:    result.Cookies(),
			})
		case map[string]interface{}:
			// Поле для добавления в лог дополнительных данных контекста
			extend = append(extend, val)
		}
	}

	methodField := zap.Field{Key: "method", Type: zapcore.StringType, String: method}
	urlField := zap.Field{Key: "url", Type: zapcore.StringType, String: url}
	codeField := zap.Field{Key: "code", Type: zapcore.Int64Type, Integer: int64(code)}
	requestField := zap.Field{Key: "request", Type: zapcore.ByteStringType, Interface: request}
	responseField := zap.Field{Key: "response", Type: zapcore.ByteStringType, Interface: response}
	extendField := zap.Field{Key: "extend", Type: zapcore.ReflectType, Interface: extend}

	result := logger.With(
		methodField, urlField, codeField, requestField, responseField, extendField,
	)

	return result
}

// contains ищет подстроку в строке
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
