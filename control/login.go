package control

import (
	"fmt"
	"sifu-box/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Login(admin bool, secret string, logger *zap.Logger) (string, error) {
	uuid := uuid.New().String()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		model.Jwt{
			Admin: admin,
			UUID:  uuid,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(model.EXPIRE_TIME)),
			},
		}).SignedString([]byte(secret))
	if err != nil {
		logger.Error(fmt.Sprintf("创建JWT失败: [%s]", err.Error()))
		return "", fmt.Errorf("生成token失败")
	}
	return token, nil
}

func Verify(authorization, secret string, logger *zap.Logger) (string, bool, error) {
	token, err := jwt.ParseWithClaims(authorization, &model.Jwt{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		logger.Error(fmt.Sprintf(`解析"authorization"字段失败: [%s]`, err.Error()))
		return "", false, fmt.Errorf(`解析"authorization"字段失败`)
	}
	if !token.Valid || token == nil {
		return "", false, fmt.Errorf("token已经失效")
	}
	if claims, ok := token.Claims.(*model.Jwt); ok {
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
			model.Jwt{
				Admin: claims.Admin,
				UUID:  claims.UUID,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(model.EXPIRE_TIME)),
				},
			}).SignedString([]byte(secret))
		if err != nil {
			logger.Error(fmt.Sprintf("生成token失败: [%s]", err.Error()))
			return "", false, fmt.Errorf("生成token失败")
		}
		return token, claims.Admin, nil
	}
	logger.Error("未知错误")
	return "", false, fmt.Errorf("未知错误")
}
