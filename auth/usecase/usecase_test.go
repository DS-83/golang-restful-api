package usecase

// import (
// 	"context"
// 	"example-restful-api-server/auth/repo/mock"
// 	"example-restful-api-server/models"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// type fields struct {
// 	userRepo  *mock.UserStorageMock
// 	tokenRepo *mock.TokenStorageMock
// }

// func TestAuthFlow(t *testing.T) {

// 	type args struct {
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		code   int
// 	}{
// 		{
// 			name: "Sign Up",
// 			fields: fields{
// 				userRepo:  &mock.UserStorageMock{},
// 				tokenRepo: &mock.TokenStorageMock{},
// 			},
// 			args: args{},
// 			code: 0,
// 		},
// 	}

// 	var (
// 		ctx      = context.Background()
// 		username = "user"
// 		password = "pass"

// 		user = &models.User{
// 			Username: username,
// 			Password: "$2a$08$CejQA7X35DwLG.dTHu.QO.asf.YAedzbqG9GVJbM7xSPomLY54lmy", // sha1 of pass+salt
// 		}
// 	)

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			uc := NewAuthUsecase(tt.fields.userRepo, tt.fields.tokenRepo, []byte("secret"))
// 			tt.fields.userRepo.On("CreateUser", user).Return(nil)
// 			err := uc.SignUp(ctx, username, password)
// 			assert.NoError(t, err)

// 		})
// 	}
// }
