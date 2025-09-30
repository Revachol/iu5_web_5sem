package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetHistoricalRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	logrus.Info("Received ID:", idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	historical_object, err := h.Repository.GetHistoricalRequest(id)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	ctx.HTML(http.StatusOK, "historical_object.html", gin.H{
		"historical_object": historical_object,
	})
}
