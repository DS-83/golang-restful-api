package usecase

import (
	"context"
	"fotogramm/example-restful-api-server/auth/repo/mysqldb"
	"fotogramm/example-restful-api-server/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthFlow(t *testing.T) {
	type fields struct {
		repo *mysqldb.UserRepoMock
	}
	type args struct {
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		code   int
	}{
		{
			name: "Sign Up",
			fields: fields{
				repo: &mysqldb.UserRepoMock{},
			},
			args: args{},
			code: 0,
		},
	}

	var (
		ctx      = context.Background()
		username = "user"
		password = "pass"

		user = &models.User{
			Username: username,
			Password: "$2a$08$CejQA7X35DwLG.dTHu.QO.asf.YAedzbqG9GVJbM7xSPomLY54lmy", // sha1 of pass+salt
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewAuthUsecase(tt.fields.repo, []byte("secret"))
			tt.fields.repo.On("CreateUser", user).Return(nil)
			err := uc.SignUp(ctx, username, password)
			assert.NoError(t, err)

		})
	}
}
