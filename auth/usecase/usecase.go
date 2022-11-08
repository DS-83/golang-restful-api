package usecase

import (
	"context"
	"example-restful-api-server/auth"
	pb "example-restful-api-server/authrpc"
	"example-restful-api-server/models"

	"github.com/golang-jwt/jwt/v4"
)

type AuthUsecase struct {
	userRepo  auth.UserRepo
	rpcClient pb.AuthServiceClient
}

type AuthClaims struct {
	User *models.User
	jwt.RegisteredClaims
}

func NewAuthUsecase(a auth.UserRepo, c pb.AuthServiceClient) *AuthUsecase {
	return &AuthUsecase{
		userRepo:  a,
		rpcClient: c,
	}
}

func (c *AuthUsecase) SignUp(ctx context.Context, username, pass string) (err error) {
	cred := &pb.SignUpRequest{
		Username: username,
		Password: pass,
	}

	resp, err := c.rpcClient.SignUp(ctx, cred)
	if err != nil {
		return err
	}

	user, err := c.userRepo.CreateUser(ctx, toModelsUser(resp))
	if err != nil {
		return err
	}

	// Send request with updated User
	req := &pb.UpdRequest{
		Filtr:  resp,
		Upd:    toPbUser(user),
		SignUp: true,
	}

	if _, err = c.rpcClient.Update(ctx, req); err != nil {
		return err
	}

	return nil
}

// Sign in user and get JWT string
func (c *AuthUsecase) SignIn(ctx context.Context, username, pass string) (*string, error) {
	req := &pb.SignInRequest{
		Username: username,
		Password: pass,
	}

	resp, err := c.rpcClient.SignIn(ctx, req)
	if err != nil {
		return nil, err
	}

	return &resp.Token, nil
}

func (c *AuthUsecase) ParseTokenFromString(ctx context.Context, tokenString string) (*models.User, error) {
	req := &pb.ParseRequest{
		Token: tokenString,
	}

	resp, err := c.rpcClient.ParseToken(ctx, req)
	if err != nil {
		return nil, err
	}

	return toModelsUser(resp), nil
}

func (c *AuthUsecase) UpdateUser(ctx context.Context, filt, upd *models.User, t string) error {
	req := &pb.UpdRequest{
		Filtr:  toPbUser(filt),
		Upd:    toPbUser(upd),
		Token:  t,
		SignUp: false,
	}

	resp, err := c.rpcClient.Update(ctx, req)
	if err != nil {
		return err
	}

	if err := c.userRepo.UpdateUser(ctx, filt, toModelsUser(resp)); err != nil {
		return err
	}
	return nil
}

func (c *AuthUsecase) DeleteUser(ctx context.Context, u *models.User, token string) error {
	req := &pb.DelRequest{
		User:  toPbUser(u),
		Token: token,
	}

	_, err := c.rpcClient.Delete(ctx, req)
	if err != nil {
		return err
	}

	if err := c.userRepo.DeleteUser(ctx, u); err != nil {
		return err
	}

	return nil
}

func toModelsUser(u *pb.User) *models.User {
	return &models.User{
		ID:       int(u.MysqlId),
		MongoID:  u.Id,
		Username: u.Username,
		Password: u.Password,
	}
}

func toPbUser(u *models.User) *pb.User {
	return &pb.User{
		Id:       u.MongoID,
		MysqlId:  int64(u.ID),
		Username: u.Username,
		Password: u.Password,
	}
}
