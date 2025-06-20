package handler

// Swagger аннотации для Auth хендлеров

// Login авторизация пользователя
// @Summary Авторизация пользователя
// @Description Выполняет вход пользователя в систему и возвращает JWT токен
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body model.LoginRequestSwagger true "Данные для входа"
// @Success 200 {object} model.LoginResponseSwagger "Успешная авторизация"
// @Failure 400 {object} model.ErrorResponseSwagger "Некорректные данные запроса"
// @Failure 401 {object} model.ErrorResponseSwagger "Неверные учетные данные"
// @Failure 500 {object} model.ErrorResponseSwagger "Внутренняя ошибка сервера"
// @Router /auth/login [post]
func (h *AuthHandler) LoginDoc() {}
