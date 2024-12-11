package auth

import (
	"context"
	"doc-server/internal/models"
	"doc-server/internal/utils"
	"fmt"
)

func (s *serv) Register(ctx context.Context, token string, login string, pswd string) error {
	if token != s.configApp.AdminToken() {
		return fmt.Errorf("error: admin token is not equals")
	}

	if !utils.IsValidatePassword(pswd) {
		return fmt.Errorf("error: the password does not meet the conditions")
	}

	passHash, err := utils.HashPasswordT(pswd)
	if err != nil {
		return err
	}

	err = s.userRepo.Create(ctx, models.User{
		Login:    login,
		Password: passHash,
	})
	if err != nil {
		return err
	}

	return nil
}
