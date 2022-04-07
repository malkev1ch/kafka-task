package handler

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/malkev1ch/kafka-task/internal/config"
	"github.com/malkev1ch/kafka-task/internal/consumer"
	"github.com/malkev1ch/kafka-task/internal/logger"
	"github.com/malkev1ch/kafka-task/internal/producer"
	"github.com/malkev1ch/kafka-task/internal/repository"
	"net/http"
	"strconv"
)

type Handler struct {
	repo      *repository.PostgresRepository
	cfg       *config.Config
	log       logger.Logger
	producer  *producer.Producer
	consGroup *consumer.ConsumerGroup
}

// NewHandler Handler constructor
func NewHandler(
	prod *producer.Producer,
	consGroup *consumer.ConsumerGroup,
	cfg *config.Config,
	repo *repository.PostgresRepository,
	log logger.Logger,
) *Handler {
	return &Handler{
		producer:  prod,
		log:       log,
		consGroup: consGroup,
		cfg:       cfg,
		repo:      repo,
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
		if err := h.producer.WriteMessages(sendTimeSecond, messagesCountPerSecond); err != nil {
			h.log.Fatalf("error while sending messages - %e", err)
		}
	}()
	return c.String(http.StatusOK, fmt.Sprintf("Successfully produced %d messages", messagesCountPerSecond*sendTimeSecond))
}

// StartConsumer receives signal to start consumers group
func (h *Handler) StartConsumer(c echo.Context) error {

	h.consGroup.RunConsumers(h.cfg.KafkaTopic, h.cfg.NumConsumer, h.repo)

	return c.String(http.StatusOK, "consumer launch process started")
}

// StartProducer receives signal to start consumers group
func (h *Handler) StartProducer(c echo.Context) error {
	h.producer.Writer = h.producer.GetNewKafkaWriter(h.cfg.KafkaTopic)

	return c.String(http.StatusOK, "producer launch process started")
}
