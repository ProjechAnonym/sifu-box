package control

import (
	"fmt"
	"sifu-box/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Login(admin bool, secret string, logger *zap.Logger) (string, error) {
	uuid := uuid.New().String()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, 
		models.Jwt{
			Admin: admin,
			UUID:  uuid,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(models.EXPIRETIME)),
			},
		}).SignedString([]byte(secret))
	if err != nil {
		logger.Error(fmt.Sprintf("创建JWT失败: [%s]", err.Error()))
		return "", fmt.Errorf("生成token失败")
	}
	return token, nil
}

func Verify(authorization, secret string, logger *zap.Logger) (string, error){
	token, err := jwt.ParseWithClaims(authorization, &models.Jwt{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		logger.Error(fmt.Sprintf("解析'authorization'字段失败: [%s]", err.Error()))
		return "", fmt.Errorf("解析'authorization'字段失败")
	}
	if !token.Valid || token == nil {
		return "nil", fmt.Errorf("token已经失效")
	}
	if claims, ok := token.Claims.(*models.Jwt); ok {
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, 
			models.Jwt{
				Admin: claims.Admin,
				UUID:  claims.UUID,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(models.EXPIRETIME)),
				},
			}).SignedString([]byte(secret))
		if err != nil {
			logger.Error(fmt.Sprintf("生成token失败: [%s]", err.Error()))
			return "", fmt.Errorf("生成token失败")
		}
		return token, nil
	}
	logger.Error("未知错误")
	return "", fmt.Errorf("未知错误")
}