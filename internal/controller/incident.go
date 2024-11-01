package controller

import "github.com/gin-gonic/gin"

func (controller *Controller) queryIncidents() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, map[string]interface{}{"a": 1})
	}
}
