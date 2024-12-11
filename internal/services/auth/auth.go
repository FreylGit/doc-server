package auth

import (
	"context"
	"doc-server/internal/utils"
	"fmt"
)

func (s *serv) Auth(ctx context.Context, login string, pswd string) (string, error) {
	user, err := s.userRepo.Get(ctx, login)
	if err != nil {
		return "", err
	}
	//passwordHash, err := utils.HashPasswordT(pswd)
	err = utils.CheckPassword(user.Password, pswd)
	if err != nil {
		return "", fmt.Errorf("password is not equal")
	}
	//TODO: засейвить мб?
	token := utils.GenerateAuthToken(user.Id, user.Login, s.configApp.SecretKey())
	return token, nil
}
