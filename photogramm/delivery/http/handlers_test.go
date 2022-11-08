// nolint
package http

import (
	"bytes"
	"encoding/json"
	"example-restful-api-server/auth"
	"example-restful-api-server/models"
	"example-restful-api-server/photogramm/usecase"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var user = &models.User{
	Username: "test",
	Password: "test",
}

// func TestHandler_Upload(t *testing.T) {
// 	type fields struct {
// 		useCase usecase.PhotogrammUsecaseMock
// 	}
// 	type args struct {
// 		fileName string
// 		albName  string
// 		responce string
// 		formName string
// 		content  string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		code   int
// 	}{{
// 		name: "Not an image",
// 		fields: fields{
// 			useCase: usecase.PhotogrammUsecaseMock{},
// 		},
// 		args: args{
// 			fileName: "upload_test/123.txt",
// 			albName:  "",
// 			responce: `{"response":"not image"}`,
// 			formName: "photo",
// 			content:  "application/octet-stream",
// 		},
// 		code: 406,
// 	},
// 		{
// 			name: "Valid request",
// 			fields: fields{
// 				useCase: usecase.PhotogrammUsecaseMock{},
// 			},
// 			args: args{
// 				fileName: "upload_test/image.png",
// 				albName:  "",
// 				responce: `{"response":"id"}`,
// 				formName: "photo",
// 				content:  "image/png",
// 			},
// 			code: 200,
// 		},
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			r := gin.Default()
// 			group := r.Group("/api", func(c *gin.Context) {
// 				c.Set(auth.CtxUserKey, user)
// 			})
// 			RegisterRoutes(group, &tt.fields.useCase)

// 			fileDir, _ := os.Getwd()
// 			fileName := tt.args.fileName
// 			filePath := path.Join(fileDir, fileName)

// 			file, _ := os.Open(filePath)
// 			defer file.Close()

// 			body := &bytes.Buffer{}
// 			writer := multipart.NewWriter(body)

// 			partHeader := make(textproto.MIMEHeader, 5)
// 			part, _ := writer.CreateFormFile("photo", filepath.Base(file.Name()))
// 			disp := fmt.Sprintf(`form-data; name="data"; filename="%s"`, fileName)
// 			partHeader.Add("Content-Disposition", disp)
// 			partHeader.Add("Content-Type", "image/jpeg")
// 			part, _ = writer.CreatePart(partHeader)
// 			io.Copy(part, file)
// 			writer.Close()

// 			w := httptest.NewRecorder()
// 			req, _ := http.NewRequest("POST", "/api/photogramm/upload", body)
// 			req.Header.Set("Content-Type", tt.args.content)
// 			fmt.Println(req.Header.Get("Content-Type"))
// 			multi, _, _ := req.FormFile(tt.args.albName)
// 			tt.fields.useCase.On("UploadPhoto", user, tt.args.albName, multi).Return("id", nil)
// 			req.Header.Set("Content-Type", writer.FormDataContentType())
// 			r.ServeHTTP(w, req)

// 			assert.Equal(t, tt.code, w.Code)
// 			assert.Equal(t, []byte(tt.args.responce), w.Body.Bytes())
// 		})
// 	}
// }

func TestHandler_GetPhoto(t *testing.T) {
	type fields struct {
		useCase usecase.PhotogrammUsecaseMock
	}
	type args struct {
		id       string
		response *models.Photo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		code   int
	}{{
		name: "Valid input",
		fields: fields{
			useCase: usecase.PhotogrammUsecaseMock{},
		},
		args: args{
			id: "123",
			response: &models.Photo{
				ID:        "123",
				Username:  user.Username,
				UserID:    1,
				AlbumName: "test",
			},
		},
		code: 200,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := gin.Default()
			group := r.Group("/api", func(c *gin.Context) {
				c.Set(auth.CtxUserKey, user)
			})
			RegisterRoutes(group, &tt.fields.useCase)

			tt.fields.useCase.On("GetPhoto", user, tt.args.id).Return(&models.Photo{
				ID:        tt.args.id,
				Username:  user.Username,
				UserID:    1,
				AlbumName: "test",
			}, nil)

			w := httptest.NewRecorder()
			url := fmt.Sprintf("/api/photogramm/getphoto/%s", tt.args.id)
			req, _ := http.NewRequest("GET", url, &bytes.Buffer{})
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			resp, _ := json.Marshal(tt.args.response)
			assert.Equal(t, resp, w.Body.Bytes())
		})
	}
}

func TestHandler_RemovePhoto(t *testing.T) {
	type fields struct {
		useCase usecase.PhotogrammUsecaseMock
	}
	type args struct {
		id       string
		response string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		code   int
	}{{
		name: "Valid input",
		fields: fields{
			useCase: usecase.PhotogrammUsecaseMock{},
		},
		args: args{
			id:       "123",
			response: `{"response":"delete success"}`,
		},
		code: 200,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := gin.Default()
			group := r.Group("/api", func(c *gin.Context) {
				c.Set(auth.CtxUserKey, user)
			})
			RegisterRoutes(group, &tt.fields.useCase)

			tt.fields.useCase.On("RemovePhoto", user, tt.args.id).Return(nil)

			w := httptest.NewRecorder()
			input := fmt.Sprintf(`{"id":"%s"}`, tt.args.id)
			req, _ := http.NewRequest("DELETE", "/api/photogramm/removephoto", bytes.NewBuffer([]byte(input)))
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			assert.Equal(t, []byte(tt.args.response), w.Body.Bytes())
		})
	}
}

func TestHandler_CreateAlbum(t *testing.T) {
	type fields struct {
		useCase usecase.PhotogrammUsecaseMock
	}
	type args struct {
		name     string
		response string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		code   int
	}{{
		name: "Valid input",
		fields: fields{
			useCase: usecase.PhotogrammUsecaseMock{},
		},
		args: args{
			name:     "test_album",
			response: `{"response":"success"}`,
		},
		code: 200,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := gin.Default()
			group := r.Group("/api", func(c *gin.Context) {
				c.Set(auth.CtxUserKey, user)
			})
			RegisterRoutes(group, &tt.fields.useCase)

			tt.fields.useCase.On("CreateAlbum", user, tt.args.name).Return(nil)

			w := httptest.NewRecorder()
			input := fmt.Sprintf(`{"name":"%s"}`, tt.args.name)
			req, _ := http.NewRequest("POST", "/api/photogramm/createalbum", bytes.NewBuffer([]byte(input)))
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			assert.Equal(t, []byte(tt.args.response), w.Body.Bytes())
		})
	}
}

func TestHandler_GetAlbum(t *testing.T) {
	type fields struct {
		useCase usecase.PhotogrammUsecaseMock
	}
	type args struct {
		name     string
		response *models.PhotoAlbum
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		code   int
	}{{
		name: "Valid input",
		fields: fields{
			useCase: usecase.PhotogrammUsecaseMock{},
		},
		args: args{
			name: "test_album",
			response: &models.PhotoAlbum{
				Name:     "test_album",
				UserID:   1,
				PhotosID: []string{},
				Total:    0,
			},
		},
		code: 200,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := gin.Default()
			group := r.Group("/api", func(c *gin.Context) {
				c.Set(auth.CtxUserKey, user)
			})
			RegisterRoutes(group, &tt.fields.useCase)

			tt.fields.useCase.On("GetAlbum", user, tt.args.name).Return(&models.PhotoAlbum{
				Name:     tt.args.name,
				UserID:   1,
				PhotosID: []string{},
				Total:    0,
			}, nil)

			w := httptest.NewRecorder()
			url := fmt.Sprintf("/api/photogramm/getalbum/%s", tt.args.name)
			req, _ := http.NewRequest("GET", url, &bytes.Buffer{})
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			resp, _ := json.Marshal(tt.args.response)
			assert.Equal(t, resp, w.Body.Bytes())
		})
	}
}

func TestHandler_RemoveAlbum(t *testing.T) {
	type fields struct {
		useCase usecase.PhotogrammUsecaseMock
	}
	type args struct {
		name     string
		response string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		code   int
	}{{
		name: "Valid input",
		fields: fields{
			useCase: usecase.PhotogrammUsecaseMock{},
		},
		args: args{
			name:     "test_album",
			response: `{"response":"success"}`,
		},
		code: 200,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := gin.Default()
			group := r.Group("/api", func(c *gin.Context) {
				c.Set(auth.CtxUserKey, user)
			})
			RegisterRoutes(group, &tt.fields.useCase)

			tt.fields.useCase.On("RemoveAlbum", user, tt.args.name).Return(nil)

			w := httptest.NewRecorder()
			input := fmt.Sprintf(`{"name":"%s"}`, tt.args.name)
			req, _ := http.NewRequest("DELETE", "/api/photogramm/removealbum", bytes.NewBuffer([]byte(input)))
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			assert.Equal(t, []byte(tt.args.response), w.Body.Bytes())
		})
	}
}

func TestHandler_GetInfo(t *testing.T) {
	type fields struct {
		useCase usecase.PhotogrammUsecaseMock
	}
	type args struct {
		name     string
		response *models.User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		code   int
	}{{
		name: "Valid input",
		fields: fields{
			useCase: usecase.PhotogrammUsecaseMock{},
		},
		args: args{
			name: "test_album",
			response: &models.User{
				ID:          1,
				Username:    "test",
				Password:    "test",
				PhotoAlbums: []models.PhotoAlbum{},
			},
		},
		code: 200,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := gin.Default()
			group := r.Group("/api", func(c *gin.Context) {
				c.Set(auth.CtxUserKey, user)
			})
			RegisterRoutes(group, &tt.fields.useCase)

			tt.fields.useCase.On("GetInfo", user).Return(&models.User{
				ID:          1,
				Username:    user.Username,
				Password:    user.Password,
				PhotoAlbums: []models.PhotoAlbum{},
			}, nil)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/photogramm/getinfo", &bytes.Buffer{})
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
			resp, _ := json.Marshal(tt.args.response)
			assert.Equal(t, resp, w.Body.Bytes())
		})
	}

}
