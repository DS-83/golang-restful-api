package http

import (
	"bytes"
	"example-restful-api-server/auth"
	"example-restful-api-server/models"
	"example-restful-api-server/photogramm/usecase"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Upload(t *testing.T) {
	type fields struct {
		useCase usecase.PhotogrammUsecaseMock
	}
	type args struct {
		fileHeader *multipart.FileHeader
		albName    string
		responce   string
		formName   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{{
		name: "Not an image",
		fields: fields{
			useCase: usecase.PhotogrammUsecaseMock{},
		},
		args: args{
			fileHeader: &multipart.FileHeader{
				Filename: "123",
				Header:   textproto.MIMEHeader{"Content-Type": []string{"doc"}},
				Size:     0,
			},
			albName:  "",
			responce: `{"response":"not image"}`,
			formName: "photo",
		},
	},
		{
			name: "Valid request",
			fields: fields{
				useCase: usecase.PhotogrammUsecaseMock{},
			},
			args: args{
				fileHeader: &multipart.FileHeader{
					Filename: "123",
					Header:   textproto.MIMEHeader{"Content-Type": []string{"image", "png"}},
					Size:     0,
				},
				albName:  "",
				responce: `{"response":"id"}`,
				formName: "photo",
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &models.User{
				Username: "test",
				Password: "test",
			}

			r := gin.Default()
			group := r.Group("/api", func(c *gin.Context) {
				c.Set(auth.CtxUserKey, user)
			})
			RegisterRoutes(group, &tt.fields.useCase)

			tt.fields.useCase.On("UploadPhoto", user, tt.args.albName).Return("id", nil)

			// File
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			part, _ := writer.CreatePart(tt.args.fileHeader.Header)
			file.Write([]byte(&body))
			file, _ = writer.CreateFormFile(tt.args.formName, tt.args.fileHeader.Filename)
			file.Write([]byte(`sample`))
			writer.Close()
			//

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/photogramm/upload", body)
			fmt.Println(writer.FormDataContentType())
			req.Header.Set("Content-Type", writer.FormDataContentType())
			r.ServeHTTP(w, req)

			assert.Equal(t, 400, w.Code)
			assert.Equal(t, []byte(tt.args.responce), w.Body.Bytes())
		})
	}
}
