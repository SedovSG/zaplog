# Пакет предназначен для формирования логов в соответсвии со стандартом GELF

[![Build Status](https://travis-ci.org/SedovSG/zaplog.svg?branch=main)](https://travis-ci.org/SedovSG/zaplog)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/SedovSG/zaplog)
![GitHub](https://img.shields.io/github/license/SedovSG/zaplog)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/SedovSG/zaplog)
[![Go Report Card](https://goreportcard.com/badge/github.com/SedovSG/zaplog)](https://goreportcard.com/report/github.com/SedovSG/zaplog)

[GELF](https://docs.graylog.org/en/3.1/pages/gelf.html) - сокращение от **Graylog Extended Log Format** - специального формата логов, разработанного для [Graylog](https://github.com/Graylog2/graylog2-server). Согласно стандарту GELF, в каждом лог-сообщении существует следующий набор полей:
- версия;
- хост (откуда пришло сообщение);
- время (timestamp);
- сокращенный и полный вариант сообщения;
- другие поля, настроенные самостоятельно.

Пример минимально возможного json-объекта лога:
```js
{
  "host": "127.0.0.1",
  "timestamp": "2010-10-10 11:55:36 -0700",
  "method": "GET",
  "url": "/apache_pb.gif",
  "protocol": "HTTP/1.0",
  "code": 200,
  "level": "ERROR",
  "message": "Text error",
  "tag": "my-tag",
  "type": "gelf"
}
```

#### Схема работы:
+ Выставление лога
+ Отрпавка его на порт udp:12201
+ Обработка лога Logstash или Graylog
+ И дальнейшая запись в Elasticsearch (в случае ELK)

## Структура
* Зависимости: Go 1.14.0
* Go-пакеты: [go.uber.org/zap](https://github.com/uber-go/zap)

## Влючение поддержки GELF драйвера

В docker-compose.yml добавить:
```bash
---
version: '3.7'

services:
  service:
    ...
    logging:
      driver: gelf
      options:
        gelf-address: "udp://195.201.108.163:12201"
        tag: "service-name"
    ...
```

Название тега нужно выбирать, по возможности, уникальным, т.к. из него гененируются индексы в Elasticsearch.

## Установка
```bash
$: go get -u github.com/SedovSG/zaplog
```

## Использование
```go
zaplog.Throw().Info("Text message")
zaplog.Throw("field").Debug("Text message")
zaplog.Throw("GET", "/api/v1/index", 301).Warn("Text message")
zaplog.Throw("HEAD", "/api/v1/guides?", 500, map[string]interface{}{"data": DataModel}).Error("Text message")
zaplog.Throw(Request, Response, map[string]interface{}{"other": DataModel}).Fatal("Text message")
```

Вывод:
```json
{
  "level":"INFO",
  "caller":"main/index.go:21",
  "function": "main.index",
  "message":"Text message",
  "method":"HEAD",
  "url":"/api/v1/guides?",
  "code":500,
  "stacktrace": "",
  "extend":[{
      "data": {
        "name": "Alex",
        "money": 1250
      }
  }]
}
```

## Полезные ссылки
+ Официальная документация по активации драйвера **GELF**: [https://docs.docker.com/config/containers/logging/gelf/](https://docs.docker.com/config/containers/logging/gelf/)
+ Документация по пакету **go.uber.org/zap**: [https://pkg.go.dev/go.uber.org/zap](https://pkg.go.dev/go.uber.org/zap)
