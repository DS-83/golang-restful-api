package server

import (
	"context"
	"database/sql"
	"example-restful-api-server/auth"
	httpAuth "example-restful-api-server/auth/delivery/http"
	"example-restful-api-server/auth/repo/mysqldb"
	"example-restful-api-server/auth/usecase"
	"example-restful-api-server/photogramm"
	httpPhoto "example-restful-api-server/photogramm/delivery/http"
	"example-restful-api-server/photogramm/repo/local"
	photoMysqldb "example-restful-api-server/photogramm/repo/mysqldb"
	photoUsecase "example-restful-api-server/photogramm/usecase"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type App struct {
	httpSrv *http.Server
	authUC  auth.UseCase
	photoUC photogramm.UseCase
}

func NewApp() *App {
	db := initDB()

	userRepo := mysqldb.NewUserRepo(db, viper.GetString("mysql.default_album_name"))
	photoDBRepo := photoMysqldb.NewPhotoRepo(db, viper.GetString("default_album_name"))
	photoLocalRepo := local.NewPhotoRepo(viper.GetString("local_storage"))
	albumRepo := photoMysqldb.NewAlbumRepo(db)

	return &App{
		authUC: usecase.NewAuthUsecase(
			userRepo,
			[]byte(viper.GetString("jwtKey")),
		),
		photoUC: photoUsecase.NewPhotogrammUsecase(
			photoDBRepo,
			photoLocalRepo,
			albumRepo,
		),
	}
}

func (a *App) Run(port string) error {
	// Init gin handler
	r := gin.Default()
	r.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	// Set up http handlers
	// SignUp/SignIn endpoints
	httpAuth.RegisterRoutes(r, a.authUC)

	// User delete endpoint
	authMiddleware := httpAuth.NewAuthMiddleware(a.authUC)
	delRoute := r.Group("/delete", authMiddleware)
	httpAuth.RegisterMidRoutes(delRoute, a.authUC)

	api := r.Group("/api", authMiddleware)
	httpPhoto.RegisterRoutes(api, a.photoUC)

	// HTTP Server
	a.httpSrv = &http.Server{
		Addr:           ":" + port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host=localhost
	if err := a.httpSrv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
		log.Fatalf("Failed to listen and serve: %+v", err)
	}
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpSrv.Shutdown(ctx)
}

func initDB() *sql.DB {
	db, err := sql.Open("mysql", viper.GetString("mysql.uri"))
	if err != nil {
		log.Fatal(err)
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
