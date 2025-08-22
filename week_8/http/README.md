# Week 8: HTTP + Envoy Gateway

## 🎯 Что реализовано

Полноценная архитектура с **Envoy Proxy** в качестве API Gateway для **Weather API**. Все настроено согласно production best practices с подробными комментариями для обучения.

## 🏗 Архитектура

```
┌─────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client    │────│  Envoy Gateway  │────│  Weather API    │
│             │    │   (Port 8080)   │    │   (Port 8080)   │
│ localhost   │    │                 │    │                 │
│             │    │ • Routing       │    │ • Chi router    │
│             │    │ • Load Balance  │    │ • Health checks │
│             │    │ • Health checks │    │ • CORS          │
│             │    │ • Circuit Break │    │                 │
└─────────────┘    └─────────────────┘    └─────────────────┘
                            │
                   ┌─────────────────┐
                   │ Envoy Admin     │
                   │  (Port 8082)    │
                   │ • /ready        │
                   │ • /stats        │
                   │ • /config_dump  │
                   └─────────────────┘
```

## 📁 Структура проекта

```
week_8/http/
├── cmd/
│   ├── http_server/     # Weather API сервер
│   └── http_client/     # Тестовый клиент
├── pkg/models/          # Модели данных
├── Dockerfile           # Multi-stage сборка для API
├── docker-compose.yml   # Orchestration Envoy + API
├── envoy.yaml          # Конфигурация Envoy (с комментариями!)
└── Taskfile.yaml       # Команды управления
```

## 🚀 Быстрый старт

### 1. Запуск инфраструктуры
```bash
# Собираем и запускаем все сервисы
task up

# Проверяем здоровье
task health
```

### 2. Тестирование API
```bash
# Через тестового клиента
task test-client

# Или вручную через curl
curl http://localhost:8080/api/v1/weather/healthcheck
curl -X PUT http://localhost:8080/api/v1/weather/Moscow \
  -H "Content-Type: application/json" \
  -d '{"temperature": 15.5}'
curl http://localhost:8080/api/v1/weather/Moscow
```

## 🔧 Управление

| Команда | Описание |
|---------|----------|
| `task up` | Запуск инфраструктуры |
| `task down` | Остановка |
| `task restart` | Перезапуск |
| `task logs` | Логи всех сервисов |
| `task logs-envoy` | Логи только Envoy |
| `task logs-api` | Логи только API |
| `task health` | Проверка здоровья |
| `task test-client` | Запуск тестового клиента |
| `task clean` | Очистка Docker ресурсов |

## 🌐 Endpoints

### Через Envoy Gateway (основной трафик)
- **Gateway**: http://localhost:8080
- **Health**: http://localhost:8080/api/v1/weather/healthcheck
- **Get Weather**: `GET http://localhost:8080/api/v1/weather/{city}`
- **Update Weather**: `PUT http://localhost:8080/api/v1/weather/{city}`

### Мониторинг и отладка
- **Envoy Admin**: http://localhost:8081
  - `/ready` - статус готовности
  - `/stats` - метрики
  - `/config_dump` - текущая конфигурация
- **API Direct**: http://localhost:8081 (для отладки)

## 📋 Особенности конфигурации Envoy

### 🎯 Listeners
- **HTTP Listener** на порту 8080 для клиентского трафика
- **Admin Listener** на порту 8081 для мониторинга

### 🛣 Routing
- Префикс `/api/v1/weather` → weather-api кластер
- Все остальные запросы → 404 с полезной информацией
- Настроены таймауты: 15s timeout

### 🏥 Health Checks
- HTTP проверки каждые 10 секунд
- Путь: `/api/v1/weather/healthcheck`
- 3 неудачи = unhealthy, 2 успеха = healthy


### 🔄 Load Balancing
- **Round Robin** между healthy endpoints
- **STRICT_DNS** для резолвинга имен контейнеров

## 🐳 Docker конфигурация

### Multi-stage Dockerfile
- **Stage 1**: Сборка Go приложения на Alpine
- **Stage 2**: Runtime образ с минимальными зависимостями
- Non-root пользователь для безопасности
- Health checks встроены в образ

### Docker Compose
- Isolated network для сервисов
- Health check dependencies
- Volume mounts для конфигурации
- Restart policies

## 🔍 Мониторинг и отладка

### Envoy Admin Interface
```bash
# Статус готовности
curl http://localhost:8081/ready

# Общая статистика
curl http://localhost:8081/stats

# Текущая конфигурация
curl http://localhost:8081/config_dump

# Кластеры и их статусы
curl http://localhost:8081/clusters
```

### Логи
```bash
# Все логи
docker compose logs -f

# Только Envoy (маршрутизация, ошибки)
docker compose logs -f envoy

# Только API (бизнес-логика)
docker compose logs -f weather-api
```

## 🎓 Архитектура Envoy: детальный разбор

### 🧩 Как устроен Envoy изнутри

Envoy построен по принципу **pipeline** - каждый запрос проходит через цепочку обработчиков:

```
Client Request → Listener → Filter Chain → Router → Cluster → Backend
     ↓              ↓           ↓           ↓         ↓         ↓
  HTTP/1.1      Accept      Parse HTTP   Choose    Load     Weather
    TCP         Socket      + Route      Target   Balance     API
```

### 📁 Структура конфигурации Envoy

Envoy имеет два типа конфигурации с разными назначениями:

#### `static_resources` - основные рабочие ресурсы
```yaml
static_resources:
  listeners: [...]    # Обслуживание клиентского трафика
  clusters: [...]     # Backend сервисы
  secrets: [...]      # TLS сертификаты и ключи
```

**Назначение:** Ресурсы для обработки **пользовательского трафика**
- Listeners обслуживают клиентов
- Clusters направляют запросы к backend'ам
- Можно перезагружать через admin API без рестарта

#### `admin` - административная конфигурация
```yaml
admin:
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 8081
```

**Назначение:** Управление и мониторинг **самого Envoy**
- Статистика (`/stats`)
- Проверка готовности (`/ready`)
- Дамп конфигурации (`/config_dump`)
- Изменение log level'ов
- **НЕ обслуживает пользовательский трафик!**

#### Почему разделение важно:
1. **Безопасность:** admin интерфейс обычно закрыт от внешней сети
2. **Мониторинг:** разные метрики для пользовательского vs административного трафика
3. **Жизненный цикл:** admin работает даже если основные listeners недоступны

### 🎯 Listeners: точки входа трафика

**Что это:** Слушатели определяют, на каких адресах и портах Envoy принимает входящие соединения.

```yaml
static_resources:  # ← Служебное слово Envoy
  listeners:       # ← Список всех listeners для клиентского трафика
  - name: weather_api_listener    # Имя для логов и метрик
    address:
      socket_address:
        address: 0.0.0.0          # Слушаем на всех интерфейсах
        port_value: 8080          # TCP порт
```

**Зачем несколько listeners:**
- Основной трафик (8080) vs административный (8081)
- HTTP vs gRPC vs TCP на разных портах
- Разные политики безопасности для разных типов клиентов

### 🔗 Filter Chains: обработка соединений

**Что это:** Каждое входящее соединение проходит через цепочку фильтров.

```yaml
filter_chains:
- filters:
  - name: envoy.filters.network.http_connection_manager
    # HTTP Connection Manager - главный фильтр для HTTP
```

**Как работает:**
1. **L4 фильтры** (Network Filters) - работают с TCP/UDP
2. **L7 фильтры** (HTTP Filters) - работают с HTTP внутри HTTP Connection Manager
3. Фильтры выполняются **по порядку** - порядок критичен!

#### 📚 Откуда брать фильтры

**Встроенные фильтры Envoy:** (префикс `envoy.filters.*`)
```yaml
# Network фильтры (L4):
- envoy.filters.network.http_connection_manager  # HTTP обработка
- envoy.filters.network.tcp_proxy               # TCP proxy
- envoy.filters.network.redis_proxy             # Redis proxy
- envoy.filters.network.mongo_proxy             # MongoDB proxy
- envoy.filters.network.mysql_proxy             # MySQL proxy

# HTTP фильтры (L7, внутри HTTP Connection Manager):
- envoy.filters.http.router                     # Маршрутизация (обязательный!)
- envoy.filters.http.rate_limit                 # Ограничение скорости
- envoy.filters.http.jwt_authn                  # JWT аутентификация
- envoy.filters.http.cors                       # CORS заголовки
- envoy.filters.http.gzip                       # Сжатие ответов
- envoy.filters.http.fault                      # Fault injection для тестов
- envoy.filters.http.ext_authz                  # Внешняя авторизация
```

**Внешние фильтры:** (WebAssembly, Lua, или собственные)
```yaml
- name: envoy.filters.http.wasm                 # WebAssembly фильтры
- name: envoy.filters.http.lua                  # Lua скрипты
- name: my_company.custom_filter                # Собственные фильтры
```

#### 🔧 HTTP Connection Manager детально

**Наш фильтр:**
```yaml
- name: envoy.filters.network.http_connection_manager
  typed_config:
    "@type": type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
    stat_prefix: weather_api    # Префикс для метрик
    route_config: [...]         # Конфигурация маршрутизации
    http_filters: [...]         # L7 HTTP фильтры
```

**Что делает HTTP Connection Manager:**
1. **HTTP Parsing** - парсит HTTP/1.1, HTTP/2, HTTP/3 запросы
2. **Connection Management** - управляет keep-alive, timeouts
3. **Statistics** - собирает метрики по запросам
4. **Access Logging** - логирует все HTTP запросы
5. **Header Manipulation** - добавляет/удаляет заголовки
6. **Compression** - управляет сжатием ответов
7. **Routing** - применяет правила маршрутизации

**Важные параметры:**
```yaml
http_connection_manager:
  stat_prefix: my_service           # Префикс метрик
  request_timeout: 30s              # Общий таймаут запроса
  stream_idle_timeout: 300s         # Таймаут неактивного stream
  drain_timeout: 60s                # Время для graceful shutdown
  generate_request_id: true         # Генерировать X-Request-ID
  preserve_external_request_id: true # Сохранять клиентский request ID
  use_remote_address: true          # Использовать реальный IP клиента
  access_log: [...]                 # Настройки логирования
```

#### 🔄 Пример цепочки для production:

**Network фильтры (порядок важен!):**
```yaml
filter_chains:
- filters:
  # 1. TLS терминация
  - name: envoy.filters.network.client_ssl_auth
  # 2. HTTP обработка  
  - name: envoy.filters.network.http_connection_manager
    typed_config:
      http_filters:
        # 3. Rate limiting (до авторизации!)
        - name: envoy.filters.http.rate_limit
        # 4. JWT проверка
        - name: envoy.filters.http.jwt_authn
        # 5. CORS заголовки
        - name: envoy.filters.http.cors
        # 6. Сжатие ответа
        - name: envoy.filters.http.gzip
        # 7. Маршрутизация (всегда последний!)
        - name: envoy.filters.http.router
```

**Логика выполнения:**
- **Request path:** 1 → 2 → 3 → 4 → 5 → 6 → 7 → Backend
- **Response path:** Backend → 7 → 6 → 5 → 4 → 3 → 2 → 1 → Client

### 🛣 Virtual Hosts: логическое разделение

**Что это:** Виртуальные хосты группируют маршруты по доменам или логическим сервисам.

#### 🔗 Связь Virtual Hosts с Routes

**Иерархия конфигурации:**
```yaml
route_config:
  name: main_routes
  virtual_hosts:          # ← Контейнер для маршрутов
  - name: weather_service # ← Virtual Host
    domains: ["*"]        # ← Домены для этого virtual host
    routes:               # ← Routes принадлежат Virtual Host
    - match: {...}        # ← Конкретный маршрут #1
      route: {...}
    - match: {...}        # ← Конкретный маршрут #2
      direct_response: {...}
```

**Логика выбора маршрута:**
1. **Первый шаг:** Envoy смотрит на `Host` заголовок запроса
2. **Выбор Virtual Host:** Находит virtual_host с подходящим доменом
3. **Выбор Route:** Внутри выбранного virtual_host ищет подходящий route
4. **Порядок важен:** Routes проверяются **по порядку** до первого совпадения

#### 📋 Реальный пример с несколькими сервисами:

```yaml
route_config:
  name: company_routes
  virtual_hosts:
  
  # API v1 - старая версия
  - name: api_v1
    domains: ["api.example.com", "api-v1.example.com"]
    routes:
    - match:
        prefix: "/v1/users"
      route:
        cluster: user_service_v1
        timeout: 10s
    - match:
        prefix: "/v1/orders" 
      route:
        cluster: order_service_v1
        timeout: 15s
    - match:
        prefix: "/v1/health"
      direct_response:
        status: 200
        body:
          inline_string: "API v1 OK"
  
  # API v2 - новая версия  
  - name: api_v2
    domains: ["api-v2.example.com"]
    routes:
    - match:
        prefix: "/v2/users"
      route:
        cluster: user_service_v2
        timeout: 5s    # ← Быстрее чем v1
    - match:
        prefix: "/v2/orders"
      route:
        cluster: order_service_v2
        timeout: 10s
        
  # Веб-сайт компании
  - name: website
    domains: ["example.com", "www.example.com"]
    routes:
    - match:
        prefix: "/api/"    # ← API запросы на сайте
      route:
        cluster: api_gateway
    - match:
        prefix: "/"        # ← Статические файлы
      route:
        cluster: web_frontend
        
  # Внутренние сервисы (только для мониторинга)
  - name: internal
    domains: ["internal.example.com"]
    routes:
    - match:
        prefix: "/metrics"
      route:
        cluster: prometheus_cluster
    - match:
        prefix: "/admin"
      route:
        cluster: admin_panel
        # Дополнительная авторизация через заголовки
        request_headers_to_add:
        - header:
            key: "X-Internal-Request"
            value: "true"
```

#### 🎯 Как работает выбор маршрута:

**Пример запроса:** `GET https://api-v2.example.com/v2/users/123`

1. **Host matching:** `api-v2.example.com` → подходит virtual_host `api_v2`
2. **Route matching:** `/v2/users/123` начинается с `/v2/users` → выбран первый route
3. **Forwarding:** Запрос идет на кластер `user_service_v2`

**Пример запроса:** `GET https://example.com/api/orders`

1. **Host matching:** `example.com` → подходит virtual_host `website` 
2. **Route matching:** `/api/orders` начинается с `/api/` → выбран первый route
3. **Forwarding:** Запрос идет на кластер `api_gateway`

#### ⚠️ Важные нюансы:

**Wildcard домены:**
```yaml
domains: ["*"]                    # Любой домен (catch-all)
domains: ["*.example.com"]        # Поддомены example.com
domains: ["api.*.example.com"]    # api.{что-то}.example.com
```

**Порядок virtual_hosts:**
- Более специфичные домены должны идти **первыми**
- Wildcard `["*"]` должен быть **последним**

**Fallback стратегия:**
```yaml
virtual_hosts:
- name: specific_service
  domains: ["api.example.com"]    # ← Специфичный домен первым
  routes: [...]
  
- name: catch_all
  domains: ["*"]                  # ← Wildcard последним
  routes:
  - match:
      prefix: "/"
    direct_response:
      status: 404
      body:
        inline_string: "Service not found"
```

**Зачем это нужно:**
- **Мультитенантность:** один Envoy для разных клиентов/проектов
- **A/B тестирование:** разные версии API на разных доменах
- **Безопасность:** внутренние сервисы на отдельных доменах
- **Миграция:** постепенный переход между версиями API

### 🎯 Routes: правила направления трафика

**Принцип работы:** Envoy сравнивает входящий запрос с правилами **в порядке объявления**.

```yaml
routes:
# Первое правило - самое специфичное
- match:
    prefix: "/api/v1/weather/healthcheck"
  direct_response:
    status: 200
    
# Второе правило - менее специфичное  
- match:
    prefix: "/api/v1/weather"
  route:
    cluster: weather_api_cluster
    
# Последнее правило - catch-all
- match:
    prefix: "/"
  direct_response:
    status: 404
```

**Типы совпадений:**
- `prefix: "/api"` - начинается с префикса
- `path: "/exact"` - точное совпадение
- `safe_regex: "^/api/v[0-9]+/"` - регулярные выражения

### 🏗 Clusters: группы backend серверов

**Что это:** Кластер описывает группу серверов, которые предоставляют один сервис.

#### 📍 Простой кластер (наш пример):
```yaml
clusters:
- name: weather_api_cluster
  type: STRICT_DNS           # Как находим серверы
  lb_policy: ROUND_ROBIN     # Как балансируем нагрузку
  load_assignment:           # Где находятся серверы
    cluster_name: weather_api_cluster
    endpoints:
    - lb_endpoints:
      - endpoint:
          address:
            socket_address:
              address: weather-api    # Docker DNS резолвинг
              port_value: 8080
```

#### 🌐 Кластер с множественными endpoints:

**Статический список серверов:**
```yaml
clusters:
- name: user_service_cluster
  type: STATIC              # Статический список IP
  lb_policy: ROUND_ROBIN
  load_assignment:
    cluster_name: user_service_cluster
    endpoints:
    # Первая группа серверов (например, дата-центр 1)
    - locality:
        region: "us-east-1"
        zone: "us-east-1a"
      lb_endpoints:
      - endpoint:
          address:
            socket_address:
              address: 10.0.1.100    # Сервер 1
              port_value: 8080
        load_balancing_weight: 100   # Равный вес
      - endpoint:
          address:
            socket_address:
              address: 10.0.1.101    # Сервер 2  
              port_value: 8080
        load_balancing_weight: 100
      - endpoint:
          address:
            socket_address:
              address: 10.0.1.102    # Сервер 3
              port_value: 8080
        load_balancing_weight: 50    # Меньший вес (слабый сервер)
        
    # Вторая группа серверов (например, дата-центр 2)  
    - locality:
        region: "us-west-2" 
        zone: "us-west-2a"
      lb_endpoints:
      - endpoint:
          address:
            socket_address:
              address: 10.0.2.100    # Сервер 4
              port_value: 8080
      - endpoint:
          address:
            socket_address:
              address: 10.0.2.101    # Сервер 5
              port_value: 8080
```

#### 🔄 Динамические адреса через DNS:

**STRICT_DNS - множественные A-записи:**
```yaml
clusters:
- name: api_cluster
  type: STRICT_DNS
  lb_policy: LEAST_REQUEST    # Лучше для разных нагрузок
  dns_refresh_rate: 30s       # Обновлять DNS каждые 30 сек
  load_assignment:
    cluster_name: api_cluster
    endpoints:
    - lb_endpoints:
      - endpoint:
          address:
            socket_address:
              # DNS должен возвращать несколько A-записей:
              # api.service.consul → 10.0.1.100, 10.0.1.101, 10.0.1.102
              address: api.service.consul
              port_value: 8080
```

**DNS настройка для примера выше:**
```bash
# В /etc/hosts или DNS сервере:
10.0.1.100  api.service.consul
10.0.1.101  api.service.consul  
10.0.1.102  api.service.consul
```

#### 🚀 EDS - External Discovery Service:

**Для Kubernetes или Consul:**
```yaml
clusters:
- name: dynamic_service
  type: EDS                    # External Discovery Service
  eds_cluster_config:
    eds_config:
      api_config_source:
        api_type: GRPC
        grpc_services:
        - envoy_grpc:
            cluster_name: xds_cluster   # Кластер для подключения к EDS
            
# Отдельный кластер для EDS сервера
- name: xds_cluster
  type: STATIC
  lb_policy: ROUND_ROBIN
  http2_protocol_options: {}    # EDS использует gRPC (HTTP/2)
  load_assignment:
    cluster_name: xds_cluster
    endpoints:
    - lb_endpoints:
      - endpoint:
          address:
            socket_address:
              # Consul Connect, Istio Pilot, или собственный EDS
              address: consul-connect.service.consul
              port_value: 8502
```

#### 🐳 Docker Swarm / Kubernetes примеры:

**Docker Swarm Services:**
```yaml
clusters:
- name: web_service
  type: STRICT_DNS
  lb_policy: ROUND_ROBIN
  load_assignment:
    cluster_name: web_service
    endpoints:
    - lb_endpoints:
      - endpoint:
          address:
            socket_address:
              # Docker Swarm автоматически балансирует между репликами
              address: web_service     # имя сервиса в docker-compose
              port_value: 80
```

**Kubernetes Service Discovery:**
```yaml
clusters:  
- name: k8s_service
  type: STRICT_DNS
  lb_policy: ROUND_ROBIN
  load_assignment:
    cluster_name: k8s_service
    endpoints:
    - lb_endpoints:
      - endpoint:
          address:
            socket_address:
              # Kubernetes DNS: <service>.<namespace>.svc.cluster.local
              address: user-api.production.svc.cluster.local
              port_value: 8080
```

#### 🏥 Health Checks для множественных endpoints:

```yaml
clusters:
- name: resilient_cluster
  type: STRICT_DNS
  lb_policy: ROUND_ROBIN
  
  # Детальные health checks
  health_checks:
  - timeout: 3s
    interval: 5s              # Более частые проверки
    unhealthy_threshold: 2    # Быстрее помечать как unhealthy
    healthy_threshold: 2
    http_health_check:
      path: "/health"
      expected_statuses:      # Принимать 200, 202, 204
      - start: 200
        end: 299
      request_headers_to_add:
      - header:
          key: "X-Health-Check"
          value: "envoy"
          
  # Circuit breaker для защиты backend'ов
  circuit_breakers:
    thresholds:
    - priority: DEFAULT
      max_connections: 100      # Макс соединений к кластеру
      max_pending_requests: 30  # Макс запросов в очереди  
      max_requests: 50          # Макс активных запросов
      max_retries: 3           # Макс повторных попыток
      
  # Настройки отключения плохих серверов
  outlier_detection:
    consecutive_5xx: 3         # 3 ошибки 5xx подряд → временно исключить
    interval: 30s             # Проверять каждые 30 сек
    base_ejection_time: 30s   # Исключить минимум на 30 сек
    max_ejection_percent: 50  # Не исключать больше 50% серверов
```

#### 🎯 Стратегии для динамических сред:

**1. Consul Service Discovery:**
```yaml
# Consul возвращает список healthy instances через DNS
clusters:
- name: consul_service
  type: STRICT_DNS
  dns_refresh_rate: 10s      # Часто обновлять для быстрого обнаружения
  load_assignment:
    endpoints:
    - lb_endpoints:
      - endpoint:
          address:
            socket_address:
              address: api.service.consul
              port_value: 8080
```

**2. AWS ELB/ALB через STRICT_DNS:**
```yaml
clusters:
- name: aws_service
  type: STRICT_DNS
  load_assignment:
    endpoints:
    - lb_endpoints:
      - endpoint:
          address:
            socket_address:
              # AWS Load Balancer DNS name
              address: api-lb-123456789.us-east-1.elb.amazonaws.com
              port_value: 80
```

**3. Service Mesh (Istio/Linkerd):**
```yaml
# В service mesh Envoy получает конфиг через xDS API
clusters:
- name: mesh_service
  type: EDS                  # Pilot/Controller даёт список endpoints
  eds_cluster_config:
    eds_config:
      ads: {}               # Используем Aggregated Discovery Service
```

#### 💡 Когда использовать какой тип:

- **STATIC** → Фиксированные серверы, редко меняются
- **STRICT_DNS** → Docker, простые микросервисы, Consul
- **LOGICAL_DNS** → Внешние API, CDN (резолвинг на каждый запрос)
- **EDS** → Kubernetes, Service Mesh, сложная оркестрация

### 📡 Discovery Types: как Envoy находит серверы

**STRICT_DNS** (используется в нашем конфиге):
```yaml
type: STRICT_DNS
load_assignment:
  endpoints:
  - lb_endpoints:
    - endpoint:
        address:
          socket_address:
            address: weather-api    # Docker resolves this
            port_value: 8080
```

**Как работает:**
1. Envoy делает DNS запрос для `weather-api`
2. Docker возвращает IP адрес контейнера
3. Envoy подключается к этому IP
4. При изменении IP (перезапуск контейнера) DNS обновляется автоматически

**Альтернативы:**
- **STATIC** - статический список IP:PORT
- **EDS** - External Discovery Service (для Kubernetes)
- **LOGICAL_DNS** - резолвинг DNS при каждом запросе

### 🏥 Health Checks: контроль доступности

**Механизм работы:**

```yaml
health_checks:
- timeout: 3s              # Ждем ответ максимум 3 сек
  interval: 10s            # Проверяем каждые 10 сек
  unhealthy_threshold: 3   # 3 сбоя = сервер недоступен
  healthy_threshold: 2     # 2 успеха = сервер доступен
  http_health_check:
    path: "/healthcheck"   # GET запрос на этот path
```

**Жизненный цикл endpoint'а:**
1. **HEALTHY** → Получает трафик
2. После 3 сбоев → **UNHEALTHY** → Трафик не получает
3. После 2 успехов → **HEALTHY** → Снова получает трафик

**Важные детали:**
- Health check трафик не считается в load balancing
- Unhealthy endpoints остаются в кластере для recovery
- Можно настроить разные health checks для разных протоколов

### ⚖️ Load Balancing: распределение нагрузки

**ROUND_ROBIN** (используется в конфиге):
```
Request 1 → Server A
Request 2 → Server B  
Request 3 → Server C
Request 4 → Server A (cycle repeats)
```

**Другие алгоритмы:**
- **LEAST_REQUEST** - на сервер с наименьшим количеством активных запросов
- **RANDOM** - случайный выбор
- **RING_HASH** - консистентное хеширование (для sticky sessions)

### 🔄 Поток обработки запроса

```
1. Client → HTTP Request → Envoy:8080
   │
2. weather_api_listener accepts connection
   │
3. HTTP Connection Manager parses HTTP headers
   │
4. Route matching:
   ├─ "/api/v1/weather/moscow" matches prefix "/api/v1/weather"
   └─ Route to weather_api_cluster
   │
5. Cluster selection:
   ├─ Check health status of endpoints
   ├─ Apply load balancing (Round Robin)
   └─ Select weather-api:8080
   │
6. Forward request → weather-api:8080
   │
7. Backend response ← weather-api:8080
   │
8. Response to client ← Envoy:8080
```

### 🔧 Envoy vs другие решения

**Почему STRICT_DNS в Docker:**
- Docker контейнеры имеют динамические IP
- Docker DNS автоматически обновляется при рестарте
- Проще настройки по сравнению с static IP
- Работает out-of-the-box с docker-compose

**Envoy vs Nginx:**
- **Nginx:** static config, reload для изменений
- **Envoy:** dynamic config через API, hot reload
- **Nginx:** L7 load balancer + web server  
- **Envoy:** specialized L7 proxy + observability

**Envoy vs HAProxy:**
- **HAProxy:** проще настройка, меньше ресурсов
- **Envoy:** больше возможностей, лучше observability
- **HAProxy:** TCP/HTTP load balancer
- **Envoy:** универсальный data plane для service mesh

#### 🕸️ Service Mesh: следующий уровень

**Что такое Service Mesh:**
Service Mesh — это выделенная инфраструктурная прослойка для обеспечения безопасной, быстрой и надежной связи между микросервисами.

**Архитектура Service Mesh:**
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Service A │    │   Service B │    │   Service C │
│             │    │             │    │             │
├─────────────┤    ├─────────────┤    ├─────────────┤
│ Envoy Proxy │◄──►│ Envoy Proxy │◄──►│ Envoy Proxy │  ← Data Plane
│  (Sidecar)  │    │  (Sidecar)  │    │  (Sidecar)  │
└─────────────┘    └─────────────┘    └─────────────┘
       ▲                   ▲                   ▲
       │                   │                   │
       └───────────────────┼───────────────────┘
                           │
                   ┌───────▼────────┐
                   │ Control Plane  │  ← Управление конфигом
                   │ (Istio/Linkerd)│
                   └────────────────┘
```

**Роль Envoy в Service Mesh:**
- **Data Plane:** Envoy обрабатывает весь межсервисный трафик
- **Sidecar Pattern:** Каждый сервис получает свой Envoy proxy
- **Автоматическая конфигурация:** Control Plane настраивает Envoy через xDS API

**Преимущества Envoy для Service Mesh:**

1. **Наблюдаемость (Observability):**
   ```yaml
   # Envoy автоматически генерирует метрики для каждого запроса:
   # - Latency, throughput, error rates
   # - По источнику и назначению
   # - Распределенная трассировка
   ```

2. **Безопасность:**
   ```yaml
   # Автоматический mTLS между всеми сервисами
   # JWT validation, RBAC policies
   # Шифрование в transit без изменений в коде приложения
   ```

3. **Управление трафиком:**
   ```yaml
   # Canary deployments, blue-green
   # Circuit breaking, retries, timeouts
   # Rate limiting per-service
   ```

4. **Dynamic Configuration:**
   ```yaml
   # Изменения конфигурации без рестарта
   # A/B тестирование через изменение маршрутов
   # Автоматическое обнаружение новых сервисов
   ```

**Популярные Service Mesh решения с Envoy:**

**1. Istio (Google/IBM/Lyft):**
- Самый функциональный, но сложный
- Pilot (control plane) + Envoy (data plane)
- Богатые возможности для enterprise

**2. Linkerd (Buoyant):**
- Простой в настройке и эксплуатации
- Rust-based data plane + Go control plane
- Фокус на производительности

**3. Consul Connect (HashiCorp):**
- Интегрируется с экосистемой HashiCorp
- Service discovery + service mesh
- Хорошо для hybrid cloud

**Пример конфига Envoy в Istio:**
```yaml
# Автоматически генерируется Pilot'ом
clusters:
- name: outbound|8080||productcatalog.default.svc.cluster.local
  type: EDS
  eds_cluster_config:
    eds_config:
      ads: {}               # Получаем endpoints от Pilot
  common_lb_config:
    locality_weighted_lb_config: {}
  transport_socket:
    name: envoy.transport_sockets.tls
    typed_config:
      "@type": type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext
      common_tls_context:
        tls_certificate_sds_secret_configs:
        - name: default      # Автоматический mTLS
```

**Когда использовать Service Mesh:**
- **✅ Больше 10+ микросервисов**
- **✅ Требования к безопасности (PCI DSS, SOX)**
- **✅ Сложная маршрутизация и A/B тестирование**
- **✅ Необходимость наблюдаемости без изменений кода**

**Когда НЕ нужен Service Mesh:**
- **❌ Простая архитектура (2-5 сервисов)**
- **❌ Команда не готова к операционной сложности**
- **❌ Критичны latency и throughput**
- **❌ Legacy приложения без containerization**

### 🎓 Ключевые концепции для понимания

1. **Data Plane vs Control Plane:**
   - Envoy = Data Plane (обрабатывает трафик)
   - Istio/Consul Connect = Control Plane (управляет конфигом)

2. **Статическая vs Динамическая конфигурация:**
   - Наш пример использует статический YAML
   - В production часто используется xDS API для динамики

3. **Observability из коробки:**
   - Все запросы логируются
   - Метрики доступны через admin interface
   - Distributed tracing поддерживается нативно

### В envoy.yaml найдешь:
- **200+ строк комментариев** объясняющих каждый параметр
- Примеры альтернативных настроек
- Объяснение потока обработки запроса
- Преимущества gateway архитектуры

## 🛡 Production Best Practices

✅ **Реализовано:**
- Non-root пользователь в контейнерах
- Health checks на всех уровнях
- Graceful shutdown
- Structured logging
- Resource limits через Circuit Breaker
- Isolated Docker networks
- Multi-stage builds для оптимизации размера образа

✅ **Готово к продакшену:**
- Подробное логирование и метрики
- Автоматическое исключение неработающих backend'ов
- Защита от каскадных сбоев
- Централизованная точка входа для всех API
- Легкое масштабирование (добавление новых сервисов)

## 🔧 Расширение функциональности

### Добавление нового сервиса:
1. Добавь новый кластер в `envoy.yaml`
2. Добавь маршрут в routes секцию
3. Добавь сервис в `docker-compose.yml`

### Примеры дополнительных возможностей:
- **Rate Limiting** - ограничение количества запросов
- **JWT Authentication** - проверка токенов
- **Request/Response трансформация** - изменение данных на лету
- **Retry policies** - умные повторные попытки
- **Distributed tracing** - трассировка запросов
- **gRPC поддержка** - для микросервисной архитектуры

---

**💡 Совет**: Изучи `envoy.yaml` построчно - там максимум полезной информации для понимания как работают современные API Gateway! 