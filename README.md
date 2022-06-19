# Golang notifier about reviews for Devman.org

Бот для оповещений о проверках работ на Devman.org

## Переменные окружения

- `DVMN_TOKEN` - API токен с сайта dvmn.org
- `TG_BOT_TOKEN` - Токен ТГ бота в который будут прилетать нотификации
- `TG_CHAT_ID` - Айди пользователя в ТГ которому будут прилетать нотификации

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
