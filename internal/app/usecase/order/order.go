package order

import (
	"time"

	"github.com/harunnryd/tempokerja/internal/app/repo"
	"github.com/harunnryd/tempokerja/internal/app/repo/product"
	"github.com/harunnryd/tempokerja/internal/app/repo/transaction"
	"github.com/harunnryd/tempokerja/internal/pkg/jagatempo/saga"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type Order interface {
	CreateOrder(workflow.Context, int, int, int, int) (CreateOrderResponse, error)
}

type order struct {
	repo repo.Repo
}

func New(repo repo.Repo) Order {
	return &order{repo: repo}
}

type CreateOrderResponse struct {
	ID        int `json:"id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

func (o *order) CreateOrder(ctx workflow.Context, originID int, destinationID int, productID int, quantity int) (createOrderResp CreateOrderResponse, err error) {
	var (
		wg             workflow.WaitGroup
		product        product.GetProductByIDResponse
		transaction    transaction.CreateOrderResponse
		processedCount int
	)

	ctx = workflow.WithActivityOptions(ctx, o.withActivityOpts())
	ctx, cancelFn := workflow.WithCancel(ctx)
	wg = workflow.NewWaitGroup(ctx)

	saga := saga.New()

	selector := workflow.NewNamedSelector(ctx, "create-order-selector")

	wg.Add(1)
	future, settable := workflow.NewFuture(ctx)

	workflow.Go(ctx, func(ctx workflow.Context) {
		defer wg.Done()
		err = workflow.ExecuteActivity(ctx, o.repo.Product().GetProductByID, productID).Get(ctx, &product)

		settable.Set(product, err)
	})

	selector.AddFuture(future, func(f workflow.Future) {
		if err = f.Get(ctx, &product); err != nil {
			cancelFn()
			return
		}
	})

	processedCount += 1

	wg.Wait(ctx)

	wg.Add(1)
	future, settable = workflow.NewFuture(ctx)

	workflow.Go(ctx, func(ctx workflow.Context) {
		defer wg.Done()
		err = workflow.ExecuteActivity(ctx, o.repo.Transaction().CreateOrder, originID, destinationID, productID, quantity, (product.Data.Price*quantity)).Get(ctx, &transaction)

		settable.Set(transaction, err)
	})

	selector.AddFuture(future, func(f workflow.Future) {
		if err = f.Get(ctx, &transaction); err != nil {
			cancelFn()
			return
		}
	})

	saga.AddCompensation(ctx, func(ctx workflow.Context) (err error) {
		wg.Add(1)
		future, settable := workflow.NewFuture(ctx)

		workflow.Go(ctx, func(ctx workflow.Context) {
			defer wg.Done()
			err = workflow.ExecuteActivity(ctx, o.repo.Transaction().CancelOrderByID, transaction.Data.ID).Get(ctx, &transaction)

			settable.Set(transaction, err)
		})

		selector.AddFuture(future, func(f workflow.Future) {
			if err = f.Get(ctx, &transaction); err != nil {
				cancelFn()
				return
			}
		})

		processedCount += 1

		wg.Wait(ctx)

		return
	})

	processedCount += 1

	wg.Wait(ctx)

	wg.Add(1)
	future, settable = workflow.NewFuture(ctx)

	workflow.Go(ctx, func(ctx workflow.Context) {
		defer wg.Done()
		err = workflow.ExecuteActivity(ctx, o.repo.Product().DeductQuantityByID, productID, quantity).Get(ctx, &product)

		settable.Set(product, err)
	})

	selector.AddFuture(future, func(f workflow.Future) {
		if err = f.Get(ctx, &product); err != nil {
			cancelFn()
			return
		}
	})

	processedCount += 1

	wg.Wait(ctx)

	for i := 0; i < processedCount; i++ {
		selector.Select(ctx)
		if err != nil {
			saga.Compensate(ctx)
			cancelFn()
			return
		}
	}

	// map the output.
	createOrderResp = CreateOrderResponse{ID: transaction.Data.ID, ProductID: transaction.Data.ProductID, Quantity: transaction.Data.Quantity}

	return
}

func (o *order) withActivityOpts() (ao workflow.ActivityOptions) {
	ao = workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		HeartbeatTimeout:    10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    5 * time.Minute,
			MaximumAttempts:    5,
		},
	}

	return
}
