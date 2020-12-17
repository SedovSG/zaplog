package logger

import (
	"encoding/json"
	"github.com/SedovSG/zaplog/models"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
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

	logger, err := config.Build()
	if err != nil {
		return logger, err
	}

	defer func() {
		err = logger.Sync()
	}()

	return logger, err
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

		switch val := val.(type) {
		case *http.Request:
			var body []byte

			if val.Body != nil {
				body, _ = ioutil.ReadAll(val.Body)
				defer val.Body.Close()
			}

			if method == "" {
				method = val.Method
			}

			if url == "" {
				url = val.URL.String()
			}

			request, _ = json.Marshal(models.Request{
				Method:     val.Method,
				Proto:      val.Proto,
				URL:        val.URL.String(),
				RequestURI: val.RequestURI,
				Body:       string(body),
				Header:     val.Header,
				Cookies:    val.Cookies(),
				RemoteAddr: val.RemoteAddr,
				UserAgent:  val.UserAgent(),
				Referer:    val.Referer(),
			})
		case *http.Response:
			var body []byte
			var location string

			if val.Body != nil {
				body, _ = ioutil.ReadAll(val.Body)
				defer val.Body.Close()
			}

			if loc, _ := val.Location(); loc != nil {
				location = loc.String()
			}

			if code == 0 {
				code = val.StatusCode
			}

			response, _ = json.Marshal(models.Response{
				Status:     val.Status,
				StatusCode: val.StatusCode,
				Proto:      val.Proto,
				Location:   location,
				Body:       string(body),
				Header:     val.Header,
				Cookies:    val.Cookies(),
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
