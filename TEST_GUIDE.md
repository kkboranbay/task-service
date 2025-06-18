# Testing Guide

Это руководство описывает стратегию тестирования для Task Service, следуя лучшим практикам крупных технологических компаний.

## Архитектура тестирования

Мы используем **пирамиду тестирования** с тремя основными уровнями:

```
    /\
   /  \     E2E Tests (немного)
  /____\
 /      \   Integration Tests (средне)
/________\  Unit Tests (много)
```

### 1. Unit Tests (70%)
- **Цель**: Тестируют отдельные компоненты в изоляции
- **Файлы**: `*_test.go` рядом с тестируемым кодом
- **Инструменты**: testify/mock, testify/assert, testify/suite
- **Особенности**: Быстрые, надежные, используют моки

### 2. Integration Tests (20%)
- **Цель**: Тестируют взаимодействие с внешними системами (БД)
- **Файлы**: `internal/repository/postgres/*_test.go`
- **Инструменты**: testify/suite, реальная PostgreSQL
- **Особенности**: Изолированная тестовая БД для каждого теста

### 3. End-to-End Tests (10%)
- **Цель**: Тестируют полные пользовательские сценарии
- **Файлы**: `tests/e2e/*_test.go`
- **Инструменты**: HTTP клиент, реальный сервер
- **Особенности**: Полная интеграция всех компонентов

## Структура тестов

```
task-service/
├── internal/
│   ├── mocks/              # Моки для unit тестов
│   ├── testutils/          # Утилиты для тестирования
│   ├── service/*_test.go   # Unit тесты сервисов
│   ├── api/handler/*_test.go # Unit тесты хендлеров
│   └── repository/postgres/*_test.go # Integration тесты
├── tests/
│   └── e2e/               # End-to-end тесты
├── Makefile               # Команды для запуска тестов
```

## Принципы тестирования

### 1. F.I.R.S.T.
- **Fast**: Тесты выполняются быстро
- **Independent**: Тесты не зависят друг от друга
- **Repeatable**: Результат одинаковый в любой среде
- **Self-Validating**: Четкий pass/fail результат
- **Timely**: Пишутся вместе с кодом

### 2. AAA Pattern
```go
func TestCreateTask(t *testing.T) {
    // Arrange - подготовка данных
    service := setupService()
    request := createTaskRequest()
    
    // Act - выполнение действия
    result, err := service.CreateTask(ctx, userID, request)
    
    // Assert - проверка результатов
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### 3. Given-When-Then
```go
func TestTaskCreation(t *testing.T) {
    // Given - начальное состояние
    suite.testDB.Truncate(t)
    
    // When - действие
    task, err := suite.repo.Create(ctx, userID, request)
    
    // Then - ожидаемый результат
    require.NoError(t, err)
    assert.Equal(t, expectedTitle, task.Title)
}
```

## Запуск тестов

### Локальная разработка

```bash
# Установка зависимостей
make deps

# Настройка среды разработки
make dev-setup

# Запуск всех тестов
make test

# Запуск только unit тестов
make test-unit

# Запуск только integration тестов
make test-integration

# Запуск только E2E тестов
make test-e2e

# Запуск с покрытием кода
make test-coverage

# Запуск бенчмарков
make test-benchmark
```

### Docker окружение

```bash
# Запуск тестов в Docker
docker compose -f docker-compose.test.yml up --build --abort-on-container-exit

# Запуск отдельных видов тестов
docker compose -f docker-compose.test.yml run --rm unit-tests
docker compose -f docker-compose.test.yml run --rm integration-tests
docker compose -f docker-compose.test.yml run --rm e2e-tests
```

## Инструменты и библиотеки

### Основные
- **testify/suite**: Организация тестов в наборы
- **testify/assert**: Удобные assertions
- **testify/mock**: Создание моков
- **testify/require**: Критические проверки

### Специализированные
- **httptest**: Тестирование HTTP хендлеров
- **golangci-lint**: Статический анализ кода
- **gosec**: Проверка безопасности
- **pprof**: Профилирование производительности

## Best Practices

### 1. Использование table-driven тестов

```go
func TestCreateTask(t *testing.T) {
    tests := []struct {
        name        string
        request     CreateTaskRequest
        setupMock   func(*MockRepo)
        expectError bool
        errorMsg    string
    }{
        {
            name: "successful_creation",
            request: CreateTaskRequest{
                Title: "Test Task",
                Description: "Description",
            },
            setupMock: func(m *MockRepo) {
                m.On("Create", mock.Anything).Return(task, nil)
            },
            expectError: false,
        },
        // ... другие тест-кейсы
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // тест логика
        })
    }
}
```

### 2. Использование фикстур

```go
// testutils/fixtures.go
func TaskFixture(overrides ...func(*Task)) *Task {
    task := &Task{
        ID: 1,
        Title: "Default Task",
        Status: StatusPending,
    }
    
    for _, override := range overrides {
        override(task)
    }
    
    return task
}

// В тесте
task := testutils.TaskFixture(func(t *Task) {
    t.Title = "Custom Title"
    t.Status = StatusCompleted
})
```

### 3. Cleanup и изоляция

```go
func (suite *TestSuite) SetupTest() {
    // Очистка состояния перед каждым тестом
    suite.testDB.Truncate(suite.T())
}

func (suite *TestSuite) TearDownTest() {
    // Очистка после теста
    suite.mockRepo.AssertExpectations(suite.T())
}
```

### 4. Тестирование ошибок

```go
func TestCreateTask_DatabaseError_ReturnsError(t *testing.T) {
    // Arrange
    mockRepo.On("Create", mock.Anything).Return(nil, errors.New("db error"))
    
    // Act
    task, err := service.CreateTask(ctx, userID, request)
    
    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "не удалось создать задачу")
    assert.Nil(t, task)
}
```

## Метрики качества

### Покрытие кода
- **Минимум**: 80% общего покрытия
- **Цель**: 90%+ для критических компонентов
- **Исключения**: Автогенерированный код, моки

### Performance benchmarks
```go
func BenchmarkCreateTask(b *testing.B) {
    // setup
    for i := 0; i < b.N; i++ {
        service.CreateTask(ctx, userID, request)
    }
}
```

### Время выполнения
- **Unit тесты**: < 1 секунды всего набора
- **Integration тесты**: < 30 секунд
- **E2E тесты**: < 2 минут

## Отладка тестов

### Verbose режим
```bash
go test -v ./...
```

### Отладка конкретного теста
```bash
go test -run TestCreateTask_WithValidData -v ./internal/service
```

### Профилирование тестов
```bash
go test -cpuprofile cpu.prof -memprofile mem.prof -bench .
```

## Интеграция с IDE

### VS Code
- Установите Go extension
- Настройте `settings.json`:
```json
{
    "go.testFlags": ["-v", "-count=1"],
    "go.testTimeout": "300s"
}
```

### GoLand
- Включите автозапуск тестов при сохранении
- Настройте покрытие кода
- Используйте встроенный отладчик для тестов

## Мониторинг и отчетность

### Покрытие кода
```bash
# Генерация отчета
make test-coverage

# Просмотр в браузере
open coverage.html
```

## Заключение

Хорошие тесты - это инвестиция в:
- **Качество**: Раннее обнаружение багов
- **Скорость**: Быстрые итерации и рефакторинг
- **Уверенность**: Безопасные релизы
- **Документацию**: Живые примеры использования кода

Следуйте этим принципам, и ваши тесты будут надежной основой для качественного продукта!