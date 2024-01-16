package transaction

import (
	"gorm.io/gorm"
)

type ITransactionService interface {
	Begin() error
	Commit() error
	Rollback() error
}

type transactionService struct {
	db *gorm.DB
}

func NewTransactionService(db *gorm.DB) ITransactionService {
	return &transactionService{db: db}
}

func (t *transactionService) Begin() error {
	return t.db.Begin().Error
}

func (t *transactionService) Commit() error {
	return t.db.Commit().Error
}

func (t *transactionService) Rollback() error {
	return t.db.Rollback().Error
}
