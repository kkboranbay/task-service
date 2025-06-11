package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kkboranbay/task-service/internal/config"
	"github.com/kkboranbay/task-service/internal/model"
	"github.com/rs/zerolog"
	"net/http"
	"strings"
	"time"
)

type UserClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTMiddleware struct {
	config config.AuthConfig
	log    *zerolog.Logger
}

func NewJWTMiddleware(config config.AuthConfig, logger *zerolog.Logger) *JWTMiddleware {
	return &JWTMiddleware{config, logger}
}

func (m *JWTMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "токен аутентификации отсутствует",
			})
			c.Abort()
			return
		}

		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "неверный формат токена",
			})
			c.Abort()
			return
		}

		tokenString := splitToken[1]

		claims := &UserClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("неожиданный алгоритм подписи")
			}

			return []byte(m.config.JWTSecret), nil
		})

		if err != nil {
			m.log.Error().Err(err).Str("token", tokenString).Msg("ошибка парсинга JWT токена")
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "недействительный токен аутентификации",
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "недействительный токен аутентификации",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)

		c.Next()
	}
}

func (m *JWTMiddleware) GenerateToken(userID int64) (string, error) {
	now := time.Now()
	claims := UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.TokenExpireDelta)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(m.config.JWTSecret))
	if err != nil {
		m.log.Error().Err(err).Int64("user_id", userID).Msg("ошибка подписи JWT токена")
		return "", err
	}

	return signedToken, nil
}
