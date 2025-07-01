# Redis Caching Configuration

## Переменные окружения для Redis

Для работы кеширования необходимо настроить следующие переменные окружения:

```bash
# Redis server connection
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_CONNECTION_TIMEOUT=5s

# Connection pool settings
REDIS_MAX_IDLE=10
REDIS_IDLE_TIMEOUT=240s

# Cache TTL (default: 1 hour)
REDIS_CACHE_TTL=1h
```

## Архитектура кеширования

Реализована чистая архитектура с разделением ответственности:

```
Service Layer (логика кеширования)
    ↓           ↓
UFO Repo + UFO Cache Repo
```

### Репозитории:
- **UFO Repository** (`internal/repository/ufo/`) - работа только с MongoDB
- **UFO Cache Repository** (`internal/repository/ufo_cache/`) - работа только с Redis
- **UFO Service** - оркестрация и логика кеширования

### Логика кеширования в сервисном слое:

- **Get**: Проверка UFO Cache Repo → fallback на UFO Repo → сохранение в UFO Cache Repo
- **Create**: Сохранение в UFO Repo → кеширование через UFO Cache Repo
- **Update**: Обновление в UFO Repo → инвалидация через UFO Cache Repo
- **Delete**: Удаление в UFO Repo → инвалидация через UFO Cache Repo

## Ключи кеширования

Используется префикс `ufo:sighting:{uuid}` для ключей кеша.

## Graceful degradation

При недоступности Redis сервис продолжает работать с MongoDB без кеширования.
Ошибки кеширования игнорируются и не влияют на основную бизнес-логику.

## Преимущества архитектуры

- **Принцип единственной ответственности** - каждый репозиторий отвечает только за свою БД
- **Явная логика** - вся логика кеширования видна в сервисном слое
- **Тестируемость** - легко мокать каждый репозиторий отдельно
- **Гибкость** - можно менять стратегии кеширования без изменения репозиториев
- **Предметность** - названия отражают бизнес-логику (UFO наблюдения) 