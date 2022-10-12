package http

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"fotogramm/example-restful-api-server/auth"
	"fotogramm/example-restful-api-server/auth/usecase"
	"fotogramm/example-restful-api-server/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler_SignUp(t *testing.T) {
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		args    args
		code    int
		wantErr bool
	}{{
		name: "correct input",
		args: args{
			username: "testuser",
			password: "testpass",
		},
		code:    200,
		wantErr: false,
	},
		{
			name: "empty input",
			args: args{
				username: "",
				password: "",
			},
			code:    406,
			wantErr: false,
		}}

	uc := new(usecase.AuthUsecaseMock)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testInput := &userInput{
				Username: tt.args.username,
				Password: tt.args.username,
			}

			uc.On("SignUp", testInput.Username, testInput.Password).Return(nil)

			w := httptest.NewRecorder()

			body, err := json.Marshal(testInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUP error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			req, _ := http.NewRequest("POST", "/auth/sign-up", bytes.NewBuffer(body))

			r := gin.Default()
			RegisterRoutes(r, uc)
			r.ServeHTTP(w, req)

			if w.Code != tt.code {
				t.Errorf("SignUp = %v, want %v", w.Code, tt.code)
			}
		})

	}
}

func TestHandler_SignIn(t *testing.T) {
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		args    args
		hAuth   string
		code    int
		jwt     string
		wantErr bool
	}{{
		name: "correct input",
		args: args{
			username: "testuser",
			password: "testpass",
		},
		hAuth:   "Basic",
		code:    200,
		jwt:     "{\"token\": \"jwt\"}",
		wantErr: false,
	},
		{
			name: "empty input",
			args: args{
				username: "",
				password: "",
			},
			hAuth:   "Basic",
			code:    400,
			jwt:     "",
			wantErr: false,
		}}
	uc := new(usecase.AuthUsecaseMock)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testInput := &userInput{
				Username: tt.args.username,
				Password: tt.args.username,
			}

			j := ""
			if tt.jwt != "" {
				j = "jwt"
			}
			uc.On("SignIn", testInput.Username, testInput.Password).Return(&j, nil)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("POST", "/auth/sign-in", bytes.NewBuffer(nil))
			cred := fmt.Sprintf("%s:%s", testInput.Username, testInput.Password)
			cred = base64.StdEncoding.EncodeToString([]byte(cred))
			s := fmt.Sprintf("%s %s", tt.hAuth, cred)
			req.Header.Set("Authorization", s)

			r := gin.Default()
			RegisterRoutes(r, uc)
			r.ServeHTTP(w, req)

			if w.Code != tt.code || w.Body.String() != tt.jwt {
				t.Errorf("SignIn = %v, %v, want %v, %v", w.Code, w.Body.String(), tt.code, tt.jwt)
			}
		})

	}
}

func TestHandler_Delete(t *testing.T) {
	type args struct {
		delete string
	}
	tests := []struct {
		name    string
		args    args
		code    int
		wantErr bool
	}{{
		name: "correct input",
		args: args{
			delete: "delete",
		},
		code:    200,
		wantErr: false,
	}}

	testUser := &models.User{
		Username: "testuser",
		Password: "testpass",
	}

	r := gin.Default()
	group := r.Group("/delete", func(c *gin.Context) {
		c.Set(auth.CtxUserKey, testUser)
	})

	uc := new(usecase.AuthUsecaseMock)

	RegisterMidRoutes(group, uc)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testInput := &deleteInput{
				Delete: tt.args.delete,
			}

			body, err := json.Marshal(testInput)
			assert.NoError(t, err)

			uc.On("DeleteUser", testUser).Return(nil)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("DELETE", "/delete/user", bytes.NewBuffer(body))

			r.ServeHTTP(w, req)

			if w.Code != tt.code {
				t.Errorf("Delete = %v, want %v", w.Code, tt.code)
			}
		})

	}
}
