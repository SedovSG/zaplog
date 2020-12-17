// Copyright (c) 2020 Sedov Stanislav <SedovSG@yandex.ru>.
//
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file or at
// https://opensource.org/licenses/BSD-3-Clause.

// Package zaplog предоставляет удобный способ логирования согласно спецификации
// GELF, используя в своей основе пакет uber/zap
package zaplog

import (
	"github.com/SedovSG/zaplog/logger"
	"go.uber.org/zap"
	"log"
)

// Throw добавляет новые поля в лог
// Для добавления данных контекста, отличных от Request и Response,
// необходимо передавать данные в виде map[string]interface{}, которые
// добавляются в поле "extend". После чего можно непосредственно вызвать метод
// пакета zap.Logger, для формирования лога.
//
// Пример использования:
// Throw("HEAD", "/api/v1/guides?", 500, map[string]interface{}{"data": DataModel{"Alex", 28}}).Info("Text message")
func Throw(fields ...interface{}) *zap.Logger {
	zapLogger, err := logger.InitZapLogger()
	if err != nil {
		log.Printf("Ошибка инициализации лога: %s", err.Error())
		return nil
	}

	return logger.AddFields(fields, zapLogger)
}
