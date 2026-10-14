package middlewares

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/naufalfmm/gdc-task-management/internal/models/dtos"
)

type PanicError struct {
	Err error
}

func (pe PanicError) Error() string {
	return fmt.Errorf("[PANIC RECOVER] %v", pe.Err).Error()
}

func (m middlewares) PanicRecover() gin.HandlerFunc {
	return func(gc *gin.Context) {
		defer func() {
			r := recover()
			if r == nil {
				gc.Next()
				return
			}

			var isBrokenPipe bool
			if ne, ok := r.(*net.OpError); ok {
				if se, ok := ne.Err.(*os.SyscallError); ok {
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						isBrokenPipe = true
					}
				}
			}

			if isBrokenPipe {
				gc.Error(r.(error))
				gc.Abort()
				return
			}

			err, ok := r.(error)
			if !ok {
				err = &PanicError{
					Err: err,
				}
			}

			gc.AbortWithStatusJSON(http.StatusInternalServerError, dtos.Response[error]{
				Code:    "PANIC_RECOVER",
				Message: err.Error(),
				Data:    err,
			})
		}()
		gc.Next()
	}
}
