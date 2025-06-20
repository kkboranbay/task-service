package model

import "time"

// Swagger документация для моделей

// Task представляет задачу
// @Description Модель задачи в системе
type TaskSwagger struct {
	// Уникальный идентификатор задачи
	// @example 1
	ID int64 `json:"id" example:"1"`

	// Идентификатор пользователя-владельца задачи
	// @example 123
	UserID int64 `json:"user_id" example:"123"`

	// Заголовок задачи
	// @example "Изучить Go"
	Title string `json:"title" example:"Изучить Go"`

	// Подробное описание задачи
	// @example "Пройти курс по Go и создать REST API"
	Description string `json:"description" example:"Пройти курс по Go и создать REST API"`

	// Статус выполнения задачи
	// @example "pending"
	// @Enum pending in_progress completed cancelled
	Status TaskStatus `json:"status" example:"pending" enums:"pending,in_progress,completed,cancelled"`

	// Время создания задачи
	// @example "2024-01-15T10:30:00Z"
	CreatedAt time.Time `json:"created_at" example:"2024-01-15T10:30:00Z"`

	// Время последнего обновления задачи
	// @example "2024-01-15T10:30:00Z"
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-15T10:30:00Z"`
}

// CreateTaskRequest запрос на создание задачи
// @Description Данные для создания новой задачи
type CreateTaskRequestSwagger struct {
	// Заголовок задачи (обязательное поле)
	// @example "Изучить Go"
	Title string `json:"title" binding:"required,min=1,max=255" example:"Изучить Go"`

	// Подробное описание задачи (необязательное поле)
	// @example "Пройти курс по Go и создать REST API"
	Description string `json:"description" binding:"max=1000" example:"Пройти курс по Go и создать REST API"`
}

// UpdateTaskRequest запрос на обновление задачи
// @Description Данные для обновления существующей задачи
type UpdateTaskRequestSwagger struct {
	// Заголовок задачи (необязательное поле)
	// @example "Изучить Go (обновлено)"
	Title *string `json:"title,omitempty" binding:"omitempty,min=1,max=255" example:"Изучить Go (обновлено)"`

	// Подробное описание задачи (необязательное поле)
	// @example "Пройти курс по Go, создать REST API и добавить тесты"
	Description *string `json:"description,omitempty" binding:"omitempty,max=1000" example:"Пройти курс по Go, создать REST API и добавить тесты"`

	// Статус выполнения задачи (необязательное поле)
	// @example "in_progress"
	// @Enum pending in_progress completed cancelled
	Status *TaskStatus `json:"status,omitempty" example:"in_progress" enums:"pending,in_progress,completed,cancelled"`
}

// TaskListResponse ответ со списком задач
// @Description Список задач с пагинацией
type TaskListResponseSwagger struct {
	// Массив задач
	Tasks []TaskSwagger `json:"tasks"`

	// Общее количество задач
	// @example 25
	Total int64 `json:"total" example:"25"`

	// Лимит записей на странице
	// @example 10
	Limit int `json:"limit" example:"10"`

	// Смещение (количество пропущенных записей)
	// @example 0
	Offset int `json:"offset" example:"0"`
}

// LoginRequest запрос на авторизацию
// @Description Данные для входа в систему
type LoginRequestSwagger struct {
	// Имя пользователя
	// @example "admin"
	Username string `json:"username" binding:"required" example:"admin"`

	// Пароль пользователя
	// @example "admin"
	Password string `json:"password" binding:"required" example:"admin"`
}

// LoginResponse ответ при успешной авторизации
// @Description Ответ с JWT токеном
type LoginResponseSwagger struct {
	// JWT токен для авторизации
	// @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`

	// Время истечения токена (в секундах)
	// @example 86400
	ExpiresIn int64 `json:"expires_in" example:"86400"`
}

// ErrorResponse модель ошибки
// @Description Стандартный формат ошибки API
type ErrorResponseSwagger struct {
	// HTTP код ошибки
	// @example 400
	Code int `json:"code" example:"400"`

	// Сообщение об ошибке
	// @example "Некорректные данные запроса"
	Message string `json:"message" example:"Некорректные данные запроса"`

	// Подробности ошибки (необязательное поле)
	// @example "Поле 'title' обязательно для заполнения"
	Details string `json:"details,omitempty" example:"Поле 'title' обязательно для заполнения"`
}

// HealthResponse ответ health check
// @Description Состояние сервиса
type HealthResponseSwagger struct {
	// Статус сервиса
	// @example "ok"
	Status string `json:"status" example:"ok"`

	// Временная метка
	// @example "2024-01-15T10:30:00Z"
	Timestamp time.Time `json:"timestamp" example:"2024-01-15T10:30:00Z"`

	// Версия API
	// @example "1.0.0"
	Version string `json:"version" example:"1.0.0"`
}
