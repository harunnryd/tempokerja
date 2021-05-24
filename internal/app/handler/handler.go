package handler

import (
	"github.com/harunnryd/tempokerja/config"
	"github.com/harunnryd/tempokerja/internal/app/handler/order"
	"github.com/harunnryd/tempokerja/internal/app/usecase"
	"go.temporal.io/sdk/client"
)

type Handler interface {
	Order() order.Order
}

type handler struct {
	order order.Order
}

func New(cfg config.Config, temporalClient client.Client, usecase usecase.Usecase) Handler {
	h := new(handler)

	h.order = order.New(temporalClient, usecase)

	return h
}

func (h *handler) Order() order.Order {
	return h.order
}
