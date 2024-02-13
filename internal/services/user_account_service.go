package services

import (
	"github.com/andreevym/gofermart/internal/repository/mem"
)

type UserAccountService struct {
}

func NewUserAccountService(repository *mem.MemUserAccountRepository) *UserAccountService {
	return &UserAccountService{}
}
