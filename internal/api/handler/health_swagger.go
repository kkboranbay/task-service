package handler

// Health проверка состояния сервиса
// @Summary Проверка состояния сервиса
// @Description Возвращает статус работоспособности API сервиса
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} model.HealthResponseSwagger "Сервис работает нормально"
// @Failure 500 {object} model.ErrorResponseSwagger "Сервис недоступен"
// @Router /health [get]
func HealthCheckDoc() {}
