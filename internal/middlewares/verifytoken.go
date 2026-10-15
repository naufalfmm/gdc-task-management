package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/naufalfmm/gdc-task-management/internal/consts"
	"github.com/naufalfmm/gdc-task-management/internal/models/dtos"
)

func (m middlewares) VerifyToken() gin.HandlerFunc {
	return func(gc *gin.Context) {
		tokenAuth := gc.Request.Header.Get("Authorization")
		if tokenAuth == "" {
			gc.AbortWithStatusJSON(http.StatusUnauthorized, dtos.Response[dtos.ErrorResponse]{
				Status:    false,
				Code:      consts.ErrorCodeInvalidToken,
				Message:   consts.ErrInvalidToken.Error(),
				Timestamp: time.Now(),
				Data: dtos.ErrorResponse{
					Error: consts.ErrInvalidToken,
				},
			})
			return
		}

		token := ""
		tokenSplit := strings.Split(tokenAuth, " ")
		if len(tokenSplit) > 1 {
			token = tokenSplit[1]
		}

		tokenData, err := m.jwt.Decode(token, &dtos.LoginClaims{})
		if err != nil {
			gc.AbortWithStatusJSON(http.StatusUnauthorized, dtos.Response[dtos.ErrorResponse]{
				Status:    false,
				Code:      consts.ErrorCodeInvalidToken,
				Message:   consts.ErrInvalidToken.Error(),
				Timestamp: time.Now(),
				Data: dtos.ErrorResponse{
					Error: consts.ErrInvalidToken,
				},
			})
			return
		}

		userTokenData := tokenData.(*dtos.LoginClaims)

		gc.Set(consts.XUserHeader, userTokenData.UserAuth)
		gc.Next()
	}
}
