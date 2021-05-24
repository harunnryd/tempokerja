package order

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/harunnryd/tempokerja/internal/app/usecase"
	"go.temporal.io/sdk/client"
)

type Order interface {
	CreateOrder(http.ResponseWriter, *http.Request) (interface{}, error)
}

type order struct {
	temporalClient client.Client
	usecase        usecase.Usecase
}

func New(temporalClient client.Client, usecase usecase.Usecase) Order {
	return &order{temporalClient: temporalClient, usecase: usecase}
}

type CreateOrderRequest struct {
	OriginID      int `json:"origin_id"`
	DestinationID int `json:"destination_id"`
	ProductID     int `json:"product_id"`
	Quantity      int `json:"quantity"`
}

func (o *order) CreateOrder(w http.ResponseWriter, r *http.Request) (resp interface{}, err error) {
	opts := client.StartWorkflowOptions{
		ID:        "create-order-workflow",
		TaskQueue: "CREATE_ORDER",
	}

	var createUserParam CreateOrderRequest

	if err = json.NewDecoder(r.Body).Decode(&createUserParam); err != nil {
		return
	}

	workflowRun, err := o.temporalClient.ExecuteWorkflow(
		r.Context(),
		opts,
		o.usecase.Order().CreateOrder,
		createUserParam.OriginID,
		createUserParam.DestinationID,
		createUserParam.ProductID,
		createUserParam.Quantity,
	)

	if err != nil {
		fmt.Println("error1", err)
		return
	}

	err = workflowRun.Get(r.Context(), &resp)
	if err != nil {
		return
	}
	return
}
