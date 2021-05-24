package transaction

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/harunnryd/tempokerja/internal/app/common"
	"github.com/parnurzeal/gorequest"
	"go.temporal.io/sdk/activity"
)

type Transaction interface {
	CreateOrder(context.Context, int, int, int, int, float64) (CreateOrderResponse, error)
	CancelOrderByID(context.Context, int) (CancelOrderByIDResponse, error)
}

type transaction struct {
	httpClient *gorequest.SuperAgent
}

func New(httpClient *gorequest.SuperAgent) Transaction {
	return &transaction{httpClient: httpClient}
}

type BaseResponse struct {
	ResponseCode string       `json:"response_code,omitempty"`
	ResponseDesc ResponseDesc `json:"response_desc"`
}

type ResponseDesc struct {
	ID string `json:"id"`
	EN string `json:"en"`
}

type CreateOrderResponse struct {
	BaseResponse
	Data OrderPayload `json:"data"`
}

type OrderPayload struct {
	ID        int        `json:"id"`
	ProductID int        `json:"product_id"`
	Quantity  int        `json:"quantity"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CancelOrderByIDResponse struct {
	BaseResponse
	Data interface{} `json:"data"`
}

func (t *transaction) CreateOrder(ctx context.Context, originID int, destinationID int, productID int, quantity int, amount float64) (createOrderResp CreateOrderResponse, err error) {
	heartbeat := common.StartHeartbeat(ctx, 1)
	defer heartbeat.Stop()

	activity.GetLogger(ctx).Info("CreateOrder called.")

	resp, bodyBytes, errs := t.httpClient.
		Post("http://localhost:4001/v1/transactions").
		SendMap(map[string]interface{}{
			"origin_id":      originID,
			"destination_id": destinationID,
			"product_id":     productID,
			"quantity":       quantity,
			"amount":         amount,
		}).
		EndBytes()

	if len(errs) > 0 {
		err = errs[0]
		return
	}

	err = json.Unmarshal(bodyBytes, &createOrderResp)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(createOrderResp.ResponseDesc.EN)
		return
	}

	return
}

func (t *transaction) CancelOrderByID(ctx context.Context, id int) (cancelOrderByIDResp CancelOrderByIDResponse, err error) {
	heartbeat := common.StartHeartbeat(ctx, 1)
	defer heartbeat.Stop()

	activity.GetLogger(ctx).Info("CancelOrderByID called.")

	resp, bodyBytes, errs := t.httpClient.
		Delete(fmt.Sprintf("%s/%d", "http://localhost:4001/v1/transactions", id)).
		EndBytes()

	if len(errs) > 0 {
		err = errs[0]
		return
	}

	err = json.Unmarshal(bodyBytes, &cancelOrderByIDResp)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(cancelOrderByIDResp.ResponseDesc.EN)
		return
	}

	return
}
