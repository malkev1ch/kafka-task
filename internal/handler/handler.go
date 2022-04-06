package handler

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/malkev1ch/kafka-task/internal/logger"
	"github.com/malkev1ch/kafka-task/internal/producer"
	"net/http"
	"strconv"
)

type Handler struct {
	log  logger.Logger
	prod *producer.Producer
}

// NewHandler Handler constructor
func NewHandler(prod *producer.Producer, log logger.Logger) *Handler {
	return &Handler{
		prod: prod,
		log:  log,
	}
}

// ProduceMessage receives request's params from query path and produces messages based on them
func (h *Handler) ProduceMessage(c echo.Context) error {
	sendTimeSecond, err := strconv.Atoi(c.QueryParam("seconds"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Please provide valid value of seconds")
	}
	messagesCountPerSecond, err := strconv.Atoi(c.QueryParam("amount"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Please provide valid amount value of operations")
	}
	go func() {
		h.log.Info("producer starts push message in topic")
		if err := h.prod.WriteMessages(sendTimeSecond, messagesCountPerSecond); err != nil {
			h.log.Fatalf("error while sending messages - %e", err)
		}
	}()
	return c.String(http.StatusOK, fmt.Sprintf("Successfully produced %d messages", messagesCountPerSecond*sendTimeSecond))
}
