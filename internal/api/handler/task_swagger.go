package handler

// Swagger аннотации для Task хендлеров

// CreateTask создает новую задачу
// @Summary Создать новую задачу
// @Description Создает новую задачу для авторизованного пользователя
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task body model.CreateTaskRequestSwagger true "Данные новой задачи"
// @Success 201 {object} model.TaskSwagger "Задача успешно создана"
// @Failure 400 {object} model.ErrorResponseSwagger "Некорректные данные запроса"
// @Failure 401 {object} model.ErrorResponseSwagger "Не авторизован"
// @Failure 500 {object} model.ErrorResponseSwagger "Внутренняя ошибка сервера"
// @Router /api/v1/tasks [post]
func (h *TaskHandler) CreateTaskDoc() {}

// GetTask получает задачу по ID
// @Summary Получить задачу по ID
// @Description Возвращает задачу по её уникальному идентификатору
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID задачи" minimum(1)
// @Success 200 {object} model.TaskSwagger "Данные задачи"
// @Failure 400 {object} model.ErrorResponseSwagger "Некорректный ID задачи"
// @Failure 401 {object} model.ErrorResponseSwagger "Не авторизован"
// @Failure 404 {object} model.ErrorResponseSwagger "Задача не найдена"
// @Failure 500 {object} model.ErrorResponseSwagger "Внутренняя ошибка сервера"
// @Router /api/v1/tasks/{id} [get]
func (h *TaskHandler) GetTaskDoc() {}

// ListTasks получает список задач с пагинацией
// @Summary Получить список задач
// @Description Возвращает список задач пользователя с поддержкой пагинации
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Номер страницы" minimum(1) default(1)
// @Param page_size query int false "Количество элементов на странице" minimum(1) maximum(100) default(10)
// @Success 200 {object} model.TaskListResponseSwagger "Список задач"
// @Failure 400 {object} model.ErrorResponseSwagger "Некорректные параметры запроса"
// @Failure 401 {object} model.ErrorResponseSwagger "Не авторизован"
// @Failure 500 {object} model.ErrorResponseSwagger "Внутренняя ошибка сервера"
// @Router /api/v1/tasks [get]
func (h *TaskHandler) ListTasksDoc() {}

// UpdateTask обновляет существующую задачу
// @Summary Обновить задачу
// @Description Обновляет данные существующей задачи (частичное обновление)
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID задачи" minimum(1)
// @Param task body model.UpdateTaskRequestSwagger true "Обновляемые данные задачи"
// @Success 200 {object} model.TaskSwagger "Обновленная задача"
// @Failure 400 {object} model.ErrorResponseSwagger "Некорректные данные запроса"
// @Failure 401 {object} model.ErrorResponseSwagger "Не авторизован"
// @Failure 404 {object} model.ErrorResponseSwagger "Задача не найдена"
// @Failure 500 {object} model.ErrorResponseSwagger "Внутренняя ошибка сервера"
// @Router /api/v1/tasks/{id} [put]
func (h *TaskHandler) UpdateTaskDoc() {}

// DeleteTask удаляет задачу
// @Summary Удалить задачу
// @Description Удаляет задачу по её уникальному идентификатору
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID задачи" minimum(1)
// @Success 204 "Задача успешно удалена"
// @Failure 400 {object} model.ErrorResponseSwagger "Некорректный ID задачи"
// @Failure 401 {object} model.ErrorResponseSwagger "Не авторизован"
// @Failure 404 {object} model.ErrorResponseSwagger "Задача не найдена"
// @Failure 500 {object} model.ErrorResponseSwagger "Внутренняя ошибка сервера"
// @Router /api/v1/tasks/{id} [delete]
func (h *TaskHandler) DeleteTaskDoc() {}
