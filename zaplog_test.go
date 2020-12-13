package zaplog

import (
	"github.com/SedovSG/zaplog/models"
	"go.uber.org/zap"
	"testing"
)

type DataModel struct {
	Name string
	Age  int
}

func TestThrowEmpty(test *testing.T) {
	var result interface{} = Throw()

	if val, ok := result.(*zap.Logger); !ok && val == nil {
		test.Error("Неверный тип возвращаемого значения функции")
	}
}

func TestThrow(test *testing.T) {
	var result interface{} = Throw("HEAD", "/api/v1/guides?", 500, map[string]interface{}{"data": DataModel{"Alex", 28}})

	if val, ok := result.(*zap.Logger); !ok && val == nil {
		test.Error("Неверный тип возвращаемого значения функции")
	}
}

func Example() {
	Throw().Info("Text message")

	Throw("field").Debug("Text message")

	Throw("GET", "/api/v1/index", 301).Warn("Text message")

	Throw("HEAD", "/api/v1/guides?", 500, map[string]interface{}{"data": DataModel{"Alex", 28}}).Error("Text message")

	Throw(models.Request{}, models.Response{}, map[string]interface{}{"other": DataModel{"Ann", 17}}).Fatal("Text message")
}

func Benchmark(bench *testing.B) {
	for i := 0; i < bench.N; i++ {
		Throw("HEAD", "/api/v1/guides?", 500, map[string]interface{}{"data": DataModel{"Alex", 28}}).Error("Text message")
	}
}
