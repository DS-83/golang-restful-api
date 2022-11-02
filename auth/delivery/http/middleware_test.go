package http

import (
	"example-restful-api-server/auth/usecase"
	e "example-restful-api-server/err"
	"example-restful-api-server/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_Handle(t *testing.T) {
	type fields struct {
		uc          usecase.AuthUsecaseMock
		returnArg   interface{}
		returnError error
		methodName  string
		token       string
		methodArgs  interface{}
	}
	type args struct {
		auth     string
		authType string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		code   int
	}{
		// TODO: Add test cases.
		{
			name: "No auth header",
			fields: fields{
				uc: usecase.AuthUsecaseMock{},
			},
			code: 401,
		},
		{
			name: "Empty auth header",
			fields: fields{
				uc: usecase.AuthUsecaseMock{},
			},
			args: args{
				auth:     "Authorization",
				authType: "",
			},
			code: 401,
		},
		{
			name: "Bearer Auth Header with no token",
			fields: fields{
				uc:          usecase.AuthUsecaseMock{},
				returnArg:   &models.User{},
				returnError: e.ErrInvalidAccessToken,
				methodName:  "ParseTokenFromString",
				token:       "",
			},
			args: args{
				auth:     "Authorization",
				authType: "Bearer ",
			},
			code: 401,
		},
		{
			name: "Valid request",
			fields: fields{
				uc:          usecase.AuthUsecaseMock{},
				returnArg:   &models.User{},
				returnError: nil,
				methodName:  "ParseTokenFromString",
				token:       "token",
			},
			args: args{
				auth:     "Authorization",
				authType: "Bearer token",
			},
			code: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()

			m := NewAuthMiddleware(&tt.fields.uc)

			r.POST("/api/test", m, func(ctx *gin.Context) {
				ctx.Status(http.StatusOK)
			})
			w := httptest.NewRecorder()
			if tt.fields.methodName != "" {
				tt.fields.methodArgs = tt.fields.token
				tt.fields.uc.On(tt.fields.methodName, tt.fields.methodArgs).Return(tt.fields.returnArg, tt.fields.returnError)
			}
			req, _ := http.NewRequest("POST", "/api/test", nil)
			req.Header.Set(tt.args.auth, tt.args.authType)
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.code, w.Code)
		})
	}
}
