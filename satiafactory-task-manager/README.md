# Satisfactory Task Manager

Веб-приложение для задач в контексте Satisfactory: микросервисы на Go, PostgreSQL, Docker, балансировка task-service через NGINX.

## Быстрый старт (Docker)

### 1. Требования

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (Windows)
- Порты **8080**, **8081**, **5432**, **8090** свободны

### 2. Режимы запуска

| Файл | Назначение |
|------|------------|
| `docker-compose.balance.yml` | **Рекомендуется для сдачи**: 3 реплики task-service + NGINX |
| `docker-compose.minimal.yml` | Отладка: один task-service, без NGINX |
| `docker-compose.yml` | Полный стек (+ satisfactory-data-service) |

**Важно:** не смешивайте minimal и balance одновременно — остановите старые контейнеры перед сменой режима.

### 3. Запуск (балансировка, как в ТЗ)

```powershell
cd c:\satiafactory-task-manager

# Остановить всё и убрать «лишний» одиночный task-service, если был
docker compose -f docker-compose.minimal.yml down 2>$null
docker stop task-service 2>$null
docker rm task-service 2>$null

# Собрать и запустить
docker compose -f docker-compose.balance.yml up -d --build
```

Приложение: **http://localhost:8080**

### 4. Если «всё сломалось» после экспериментов

Частая причина — старый том PostgreSQL без таблиц `users` / `tasks`:

```powershell
docker compose -f docker-compose.balance.yml down -v
docker compose -f docker-compose.balance.yml up -d --build
```

Флаг `-v` удаляет volume `postgres_data`. Данные в БД сбросятся, схема создастся заново (init-скрипты + миграции в сервисах).

### 5. Проверка балансировки

После входа откройте DevTools → Network и несколько раз обновите список задач (`GET /tasks`). В ответах через gateway к nginx заголовок **X-Instance-ID** должен чередоваться: `tasks-1`, `tasks-2`, `tasks-3`.

Или из PowerShell (нужен JWT в cookie после логина в браузере — проще смотреть в DevTools).

Логи NGINX при падении реплики:

```powershell
docker logs nginx-lb --tail 20
```

### 6. Полезные команды

```powershell
docker compose -f docker-compose.balance.yml ps
docker compose -f docker-compose.balance.yml logs -f gateway user-service nginx
docker compose -f docker-compose.balance.yml down
```

## Архитектура (balance)

```
Браузер → gateway:8080 → user-service:8081 (auth)
                      → nginx:8090 → task-service-1|2|3:8082
                      → postgres:5432
```

## Рецепты Satisfactory

Сервис `satisfactory-data-service` при первом запуске импортирует `services/satisfactory-data-service/data/Docs.json` в PostgreSQL.

На странице «Мои задачи»:
1. Введите название рецепта в поле поиска (минимум 2 символа).
2. Выберите рецепт — появится превью ингредиентов и иконки (из [SatisfactoryTools](https://github.com/greeny/SatisfactoryTools)).
3. Укажите количество партий и создайте задачу — в карточке задачи отобразится, что нужно для крафта.

Ручной импорт (если нужно обновить данные):

```powershell
docker compose -f docker-compose.balance.yml run --rm satisfactory-data-service ./data-service -import
```

## Дальше для «нагруженной» работы (по требованию преподавателя)

Уже есть: горизонтальное масштабирование, NGINX round-robin, `X-Instance-ID`, healthcheck Postgres.

План расширения:

1. **Кэш** — Redis перед `GET /tasks`
2. **Очереди** — RabbitMQ/NATS для долгих операций
3. **Метрики** — Prometheus + Grafana
4. **Нагрузочное тестирование** — `k6` или `hey` на `http://localhost:8080/tasks`

## Локальный запуск без Docker

Запустите PostgreSQL, затем в отдельных терминалах: `user-service`, `task-service` (с `INSTANCE_ID`), `gateway`. URL задаётся через переменные `DATABASE_URL`, `JWT_SECRET`, `USER_SERVICE_URL`, `TASK_SERVICE_URL`.
