package product

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

type Product interface {
	GetProductByID(context.Context, int) (GetProductByIDResponse, error)
	DeductQuantityByID(context.Context, int, int) (DeductQuantityByIDResponse, error)
}

type product struct {
	httpClient *gorequest.SuperAgent
}

func New(httpClient *gorequest.SuperAgent) Product {
	return &product{httpClient: httpClient}
}

type BaseResponse struct {
	ResponseCode string       `json:"response_code,omitempty"`
	ResponseDesc ResponseDesc `json:"response_desc"`
}

type ResponseDesc struct {
	ID string `json:"id"`
	EN string `json:"en"`
}

type ProductPayload struct {
	ID        int        `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name"`
	Price     int        `json:"price"`
	Quantity  int        `json:"quantity"`
}

type GetProductByIDResponse struct {
	BaseResponse
	Data ProductPayload `json:"data"`
}

type DeductQuantityByIDResponse struct {
	BaseResponse
	Data ProductPayload `json:"data"`
}

func (p *product) GetProductByID(ctx context.Context, id int) (getProductByIDResp GetProductByIDResponse, err error) {
	heartbeat := common.StartHeartbeat(ctx, 1)
	defer heartbeat.Stop()

	activity.GetLogger(ctx).Info("GetProductByID called.")

	resp, bodyBytes, errs := p.httpClient.
		Get(fmt.Sprintf("%s/%d", "http://localhost:4000/v1/products", id)).
		EndBytes()

	if len(errs) > 0 {
		err = errs[0]
		return
	}

	err = json.Unmarshal(bodyBytes, &getProductByIDResp)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(getProductByIDResp.ResponseDesc.EN)
		return
	}

	return
}

func (p *product) DeductQuantityByID(ctx context.Context, id int, quantity int) (deductQuantityByIDResp DeductQuantityByIDResponse, err error) {
	heartbeat := common.StartHeartbeat(ctx, 1)
	defer heartbeat.Stop()

	activity.GetLogger(ctx).Info("DeductQuantityByID called.")

	resp, bodyBytes, errs := p.httpClient.
		Patch(fmt.Sprintf("%s/%d/%s", "http://localhost:4000/v1/products", id, "deduct")).
		SendMap(map[string]interface{}{"quantity": quantity}).
		EndBytes()

	if len(errs) > 0 {
		err = errs[0]
		return
	}

	err = json.Unmarshal(bodyBytes, &deductQuantityByIDResp)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(deductQuantityByIDResp.ResponseDesc.EN)
		return
	}

	return
}
