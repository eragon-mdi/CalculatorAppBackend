# CalculatorApp

CalculatorApp — мой мини учебный проект. Подробности в п.4

## 1. Как запустить проект

### a. Клонируйте репозиторий:

```bash
git clone https://github.com/eragon-mdi/CalculatorAppBackend
cd CalculatorAppBackend
```

### b. Запуск через Docker Compose из корня проекта:
``` bash
docker compose -f docker/docker-compose.yml -p kalc up -d
```

Примечание: При желании можно изменить настройки в файле `.env` (например, порты или параметры подключения к базе данных). Этот файл находится в папке `./docker/`.

### Adminer
Если вам нужно управлять базой данных, включите Adminer, раскомментировав соответствующие строки в `docker/docker-compose.yml`.
Для доступа к Adminer перейдите по адресу: [http://localhost:9000](http://localhost:9000) (Порт можно изменить в файле `.env`).

## 2. Как остановить проект
Чтобы остановить и удалить все контейнеры, выполните:
``` bash
docker compose -p kalc down
```

## 3. Фронт 
Frontend для проекта находится в репозитории [CalculatorAppFrontendPantela](https://github.com/EvgenyMentor/CalculatorAppFrontendPantela)


## 4. **Дополнительная информация о проекте**

### Используемые технологии:

- __Go__ (в том числе Echo и GORM)
- __REST API__
- __PostgreSQL__
- __Docker__, __Docker Compose__ (для контейнеризации)

### Структура проекта:

- `cmd/` — основная логика приложения.
- `internal/` — слои приложения (handlers, service, repository).
- `docker/` — файлы конфигурации Docker, включая собственный образ на базе `golang:1.24-alpine3.21`.

