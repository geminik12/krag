package core

import (
	"net/http"

	"github.com/geminik12/krag/internal/pkg/errorsx"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

func WriteResponse(c *gin.Context, data any, err error) {
	if err != nil {
		errx := errorsx.FromError(err)
		c.JSON(errx.Code, ErrorResponse{
			Reason:  errx.Reason,
			Message: errx.Message,
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
