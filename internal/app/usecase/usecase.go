package usecase

import (
	"github.com/harunnryd/tempokerja/internal/app/repo"
	"github.com/harunnryd/tempokerja/internal/app/usecase/order"
)

type Usecase interface {
	Order() order.Order
}

type usecase struct {
	order order.Order
}

func New(repo repo.Repo) Usecase {
	u := new(usecase)
	u.order = order.New(repo)
	return u
}

func (u *usecase) Order() order.Order {
	return u.order
}
