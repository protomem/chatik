# Chatik

![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/protomem/chatik) ![GitHub commit activity (branch)](https://img.shields.io/github/commit-activity/t/protomem/chatik?color=green)

## Описание

**Chatik** - fullstack приложение, реализующие функционал мессенджера. (Пока доступны только публичные каналы и текстовые сообщения)

## Используемые технологии

- **Backend**: Go, Mongo
- **Frontend**: TypeScript, React, Redux, Mui
- **Поддерживаемые протоколы**: WebSocket, SSE

## Настройка и запуск

Есть 2 два варинта запуска: локально или через docker

### Локально

В основном используется для разработки и тестирования

```sh
make run-infra-local # команда для запуска окружения
make run-local

make run-web-local # команда для запуска frontend-а

make stop-infra-local # команда для остановки окружения
```

### Docker, docker-compose

Используется как основной способ запуска

```sh
make run-stage
# или
make

make run-stage API_URL="" # API_URL это адресс сервера на котором будет запущенно приложение (по-умолчанию: localhost:8080)

make stop-stage # команда для остановки приложения
```

### Настройка

Файлы настроек находятся в папке `/configs`. В ней есть подпапки `/local`, `/stage`. В каждой находятся настройки соответствующего режима.

Так же есть папка `/deploy`. С ней ситуация таже, что и с `/configs`. В ней находятся docker, docker-compose файлы.

## Примечания

- **Иногда сбоит WebSocket**. Нужно переключиться на SSE, совершить какую-нибудь активность(получить или отправить сообщение) и переключиться обратно.
