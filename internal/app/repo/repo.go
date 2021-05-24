package repo

import (
	"log"

	"github.com/harunnryd/tempokerja/config"
	"github.com/harunnryd/tempokerja/internal/app/driver/db"
	"github.com/harunnryd/tempokerja/internal/app/repo/product"
	"github.com/harunnryd/tempokerja/internal/app/repo/transaction"
	"github.com/parnurzeal/gorequest"
)

type Repo interface {
	Product() product.Product
	Transaction() transaction.Transaction
}

type repo struct {
	product     product.Product
	transaction transaction.Transaction
}

func New(cfg config.Config) Repo {
	dbase := db.New(db.WithConfig(cfg))
	_, err := dbase.Manager(db.PgsqlDialectParam)

	if err != nil {
		log.Fatalln("error1", err)
	}

	repo := new(repo)
	httpClient := gorequest.New()
	repo.product = product.New(httpClient)
	repo.transaction = transaction.New(httpClient)

	return repo
}

func (r *repo) Product() product.Product {
	return r.product
}

func (r *repo) Transaction() transaction.Transaction {
	return r.transaction
}
