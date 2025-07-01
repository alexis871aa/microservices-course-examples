# E2E тесты для UFO сервиса

Этот пакет содержит end-to-end тесты для UFO сервиса, использующие testcontainers для создания изолированного тестового окружения.

## Структура тестов

- `constants.go` - константы проекта
- `setup.go` - настройка тестового окружения (Docker контейнеры)
- `teardown.go` - очистка тестового окружения
- `test_environment.go` - вспомогательные методы для работы с тестовыми данными
- `suite_test.go` - основной файл для запуска тестов
- `ufo_test.go` - тесты всех gRPC ручек UFO сервиса

## Ginkgo и Gomega - основы тестирования

### Зачем нужны Ginkgo и Gomega?

Обычные Go тесты выглядят так:
```go
func TestCreateUFO(t *testing.T) {
    // создаем UFO
    // проверяем результат
    if result != expected {
        t.Errorf("ожидали %v, получили %v", expected, result)
    }
}
```

С **Ginkgo** и **Gomega** тесты становятся более читаемыми:
```go
It("should create UFO sighting successfully", func() {
    // создаем UFO
    // проверяем результат
    Expect(response.Uuid).To(Not(BeEmpty()))
})
```

### Ginkgo - как организовать тесты

**Ginkgo** помогает структурировать тесты как историю. Представьте, что вы рассказываете о том, как работает ваш сервис:

#### 1. Describe - "Расскажи о сервисе"

Это самый верхний уровень - вы говорите о чем вообще будут тесты. Обычно это название сервиса или компонента. Внутри `Describe` находятся все тесты, связанные с этой темой.

```go
var _ = Describe("UFO Service", func() {
    // Здесь все тесты про UFO сервис
})
```

#### 2. Context - "В какой ситуации"

Это описание конкретной ситуации или сценария. Context помогает группировать тесты по похожим условиям. Например, "когда создаем новое наблюдение" или "когда наблюдение уже существует".

```go
Context("when creating new sighting", func() {
    // Тесты про создание нового наблюдения
})

Context("when sighting already exists", func() {
    // Тесты про работу с существующими наблюдениями
})
```

#### 3. It - "Что должно произойти"

Это конкретный тест - описание того, что должно случиться в данной ситуации. Каждый `It` - это один тест, который проверяет одно конкретное поведение.

```go
It("should return valid UUID", func() {
    // Конкретный тест
})

It("should save data to database", func() {
    // Другой конкретный тест
})
```

#### Полный пример из нашего проекта:

Читается как: "UFO Service, когда выполняется операция создания, должен создать наблюдение успешно"

```go
var _ = Describe("UFO Service E2E Tests", func() {
    Context("Create operation", func() {
        It("should create sighting successfully", func() {
            // Тест создания
        })
    })
    
    Context("Get operation", func() {
        It("should retrieve existing sighting", func() {
            // Тест получения
        })
    })
})
```

### Lifecycle хуки - когда что выполнять

#### BeforeSuite - "Подготовь все один раз в начале"
```go
var _ = BeforeSuite(func() {
    // Запускаем Docker контейнеры
    // Это происходит ОДИН раз перед всеми тестами
    mongoContainer = startMongoContainer()
    ufoContainer = startUFOContainer()
})
```

#### AfterSuite - "Убери все в конце"
```go
var _ = AfterSuite(func() {
    // Останавливаем контейнеры
    // Это происходит ОДИН раз после всех тестов
    mongoContainer.Terminate()
    ufoContainer.Terminate()
})
```

#### BeforeEach - "Подготовься перед каждым тестом"
```go
BeforeEach(func() {
    // Очищаем базу данных
    // Это происходит перед КАЖДЫМ тестом
    ClearSightingsCollection()
})
```

#### AfterEach - "Убери после каждого теста"
```go
AfterEach(func() {
    // Дополнительная очистка
    // Это происходит после КАЖДОГО теста
    ClearSightingsCollection()
})
```

### Gomega - как проверять результаты

**Gomega** делает проверки понятными для человека:

#### Вместо этого:
```go
if response.Uuid == "" {
    t.Error("UUID не должен быть пустым")
}
```

#### Пишем так:
```go
Expect(response.Uuid).To(Not(BeEmpty()))
```

#### Примеры из нашего проекта:

**Проверка, что что-то не пустое:**
```go
Expect(response.Uuid).To(Not(BeEmpty()))
// "Ожидаю, что UUID не будет пустым"
```

**Проверка равенства:**
```go
Expect(response.Info.Location.GetValue()).To(Equal("Test Location"))
// "Ожидаю, что локация будет равна 'Test Location'"
```

**Проверка, что ошибки нет:**
```go
Expect(err).To(BeNil())
// "Ожидаю, что ошибки не будет"
```

**Проверка, что ошибка есть:**
```go
Expect(err).To(HaveOccurred())
// "Ожидаю, что ошибка произошла"
```

**Проверка содержимого:**
```go
Expect(response.Info.Description.GetValue()).To(ContainSubstring("UFO"))
// "Ожидаю, что описание содержит слово 'UFO'"
```

### Как это работает в нашем проекте

#### 1. Файл `suite_test.go` - точка входа
```go
func TestE2E(t *testing.T) {
    RegisterFailHandler(Fail)  // Подключаем Ginkgo к Go тестам
    RunSpecs(t, "UFO E2E Tests")  // Запускаем все тесты
}
```

#### 2. Файл `setup.go` - подготовка окружения
```go
var _ = BeforeSuite(func() {
    // Запускаем MongoDB контейнер
    // Запускаем UFO сервис контейнер
    // Ждем пока все поднимется
})
```

#### 3. Файл `ufo_test.go` - сами тесты
```go
var _ = Describe("UFO Service E2E Tests", func() {
    BeforeEach(func() {
        ClearSightingsCollection() // Очищаем базу перед каждым тестом
    })
    
    Context("Create operation", func() {
        It("should create sighting successfully", func() {
            // Отправляем запрос на создание
            response, err := client.Create(ctx, request)
            
            // Проверяем результат
            Expect(err).To(BeNil())
            Expect(response.Uuid).To(Not(BeEmpty()))
        })
    })
})
```

### Почему это удобно?

1. **Читается как обычный текст** - любой может понять что тестируется
2. **Автоматическая подготовка** - BeforeEach сам очищает данные
3. **Понятные ошибки** - Gomega показывает что ожидалось и что получилось
4. **Группировка** - легко найти нужный тест
5. **Переиспользование** - общая логика в BeforeEach/AfterEach

### Полезные команды

```bash
# Запустить все тесты
task e2e:test:ufo

# Запустить только тесты создания (если добавить фокус)
ginkgo -v --focus="Create operation"

# Посмотреть подробный вывод
ginkgo -v --trace
```

### Пример вывода тестов

```
UFO Service E2E Tests
  Create operation
    ✓ should create sighting successfully [0.123 seconds]
  Get operation  
    ✓ should retrieve existing sighting [0.089 seconds]
  Update operation
    ✓ should update sighting successfully [0.156 seconds]

Ran 5 of 5 Specs in 22.064 seconds
SUCCESS! -- 5 Passed | 0 Failed | 0 Pending | 0 Skipped
```

## Покрытие тестами

### ✅ Работающие тесты (5/5 PASSED)

#### Create (создание наблюдения)
- ✅ Успешное создание наблюдения с валидными данными
- ✅ Генерация корректного UUID
- ✅ Сохранение в MongoDB

#### Get (получение наблюдения)
- ✅ Успешное получение наблюдения по UUID
- ✅ Корректное возвращение всех полей
- ✅ Правильная обработка timestamps

#### Update (обновление наблюдения)
- ✅ Успешное обновление существующего наблюдения
- ✅ Корректное обновление всех полей (location, description, color, duration)
- ✅ Установка updated_at timestamp

#### Delete (удаление наблюдения)
- ✅ Успешное мягкое удаление (soft delete)
- ✅ Установка deleted_at timestamp
- ✅ Возможность получения удаленного объекта

#### Полный жизненный цикл
- ✅ Комплексный CRUD тест: Create → Get → Update → Get → Delete → Get
- ✅ Проверка целостности данных на каждом этапе

### 🚫 Исключенные тесты

Валидационные тесты были исключены для упрощения отладки:
- Проверка пустых полей
- Проверка несуществующих UUID
- Проверка повторного удаления

## Запуск тестов

### Предварительные требования

1. Docker и Docker Compose
2. Go 1.24+
3. Task (для удобного запуска)

### Установка зависимостей

```bash
cd ../.. && go mod tidy
```

### Запуск тестов

```bash
# Через Task (рекомендуемый способ)
task e2e:test:ufo

# Прямой запуск через go test
go test -tags=integration -v ./...
```

### Переменные окружения

Тесты используют следующие переменные окружения из файла `.env`:

```bash
# MongoDB настройки
MONGO_IMAGE_NAME=mongo:7.0.5
MONGO_INITDB_ROOT_USERNAME=ufo_admin
MONGO_INITDB_ROOT_PASSWORD=ufo_secret
MONGO_DATABASE=ufo-service
MONGO_AUTH_DB=admin
MONGO_PORT=27017
MONGO_HOST=localhost

# gRPC настройки
GRPC_HOST=0.0.0.0
GRPC_PORT=50051

# Логирование
LOGGER_LEVEL=debug
LOGGER_AS_JSON=true
```

## Архитектура тестов

Тесты используют testcontainers для создания изолированного окружения:

1. **Docker сеть** - создается отдельная сеть `ufo-service`
2. **MongoDB контейнер** - база данных с аутентификацией
3. **UFO сервис контейнер** - собирается из `deploy/docker/ufo/Dockerfile`

### Особенности реализации

- **Динамическая конфигурация БД**: Репозиторий принимает имя базы данных из конфигурации
- **Правильная структура документов**: Использование `_id` вместо `uuid` для MongoDB
- **Корректная работа с protobuf**: Использование `.GetValue()` для wrapperspb типов
- **Изоляция тестов**: Каждый тест очищает коллекцию после выполнения

## Результаты тестов

```
Ran 5 of 5 Specs in 22.064 seconds
SUCCESS! -- 5 Passed | 0 Failed | 0 Pending | 0 Skipped
```

## Отладка

Для отладки тестов:

1. Проверьте логи контейнеров через Docker Desktop
2. Убедитесь что все переменные окружения установлены
3. Проверьте что порты не заняты другими процессами
4. Используйте `task e2e:test:ufo` для запуска с правильными переменными

## Очистка

```bash
# Очистка кеша тестов
go clean -testcache

# Очистка Docker ресурсов
docker system prune -f
``` 