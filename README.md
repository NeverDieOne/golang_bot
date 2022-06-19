# Golang notifier about reviews for Devman.org

Бот для оповещений о проверках работ на Devman.org

## Переменные окружения

- `DVMN_TOKEN` - API токен с сайта dvmn.org

## Как запустить

### Запуск во время отладки

```console
go run main.go
```

### Запуск на прод-сервере

- Компиляция:

```console
go build
```

- Запуск

```console
./golang_bot
```
