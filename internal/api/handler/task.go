package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kkboranbay/task-service/internal/model"
	"github.com/kkboranbay/task-service/internal/service"
	"github.com/rs/zerolog"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	taskService *service.TaskService
	log         *zerolog.Logger
}

func NewTaskHandler(taskService *service.TaskService, log *zerolog.Logger) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
		log:         log,
	}
}

func (h *TaskHandler) Register(router *gin.RouterGroup) {
	tasks := router.Group("/tasks")
	{
		tasks.POST("", h.Create)
		tasks.GET("", h.List)
		tasks.GET("/:id", h.GetByID)
		tasks.PUT("/:id", h.Update)
		tasks.DELETE("/:id", h.Delete)
	}
}

func (h *TaskHandler) getUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.log.Error().Msg("user_id не найден в контексте")
		return 0, false
	}
	id, ok := userID.(int64)
	if !ok {
		h.log.Error().Interface("user_id", userID).Msg("некорректный тип user_id в контексте")
		return 0, false
	}
	return id, true
}

func (h *TaskHandler) Create(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	var req model.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error().Err(err).Msg("ошибка разбора JSON")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "некорректные данные запроса",
		})
		return
	}

	task, err := h.taskService.CreateTask(c.Request.Context(), userID, req)
	if err != nil {
		h.log.Error().Err(err).Msg("ошибка создания задачи")
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "не удалось создать задачу",
		})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Error().Err(err).Str("id", idStr).Msg("ошибка парсинга ID")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "некорректный ID задачи",
		})
		return
	}

	task, err := h.taskService.GetTaskByID(c.Request.Context(), id, userID)
	if err != nil {
		h.log.Error().Err(err).Int64("id", id).Msg("ошибка получения задачи")
		c.JSON(http.StatusNotFound, model.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "задача не найдена",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) List(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	tasks, err := h.taskService.GetTaskList(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		h.log.Error().Err(err).Msg("ошибка получения списка задач")
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "не удалось получить список задач",
		})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) Update(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Error().Err(err).Str("id", idStr).Msg("ошибка парсинга ID")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "некорректный ID задачи",
		})
		return
	}

	var req model.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error().Err(err).Msg("ошибка разбора JSON")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "некорректные данные запроса",
		})
		return
	}

	task, err := h.taskService.UpdateTask(c.Request.Context(), id, userID, req)
	if err != nil {
		h.log.Error().Err(err).Int64("id", id).Msg("ошибка обновления задачи")
		c.JSON(http.StatusNotFound, model.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "задача не найдена или не удалось обновить",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Error().Err(err).Str("id", idStr).Msg("ошибка парсинга ID")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "некорректный ID задачи",
		})
		return
	}

	err = h.taskService.DeleteTask(c.Request.Context(), id, userID)
	if err != nil {
		h.log.Error().Err(err).Int64("id", id).Msg("ошибка удаления задачи")
		c.JSON(http.StatusNotFound, model.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "задача не найдена или не удалось удалить",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
