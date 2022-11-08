package server

import (
	"context"
	"database/sql"
	"example-restful-api-server/auth"
	httpAuth "example-restful-api-server/auth/delivery/http"
	mysqlUser "example-restful-api-server/auth/repo/mysqldb"
	"example-restful-api-server/auth/usecase"
	"example-restful-api-server/photogramm"
	httpPhoto "example-restful-api-server/photogramm/delivery/http"
	"example-restful-api-server/photogramm/repo/local"
	mysqlPhoto "example-restful-api-server/photogramm/repo/mysqldb"
	photoUsecase "example-restful-api-server/photogramm/usecase"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "example-restful-api-server/docs"

	pb "example-restful-api-server/authrpc"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	httpSrv *http.Server
	authUC  auth.UseCase
	photoUC photogramm.UseCase
}

func NewApp() *App {
	db := initDB()
	// db := initGormDB(sqlDB)
	rpcConn := initRpcClient()

	userRepo := mysqlUser.NewUserRepo(db, viper.GetString("default_album_name"))
	photoDBRepo := mysqlPhoto.NewPhotoRepo(db, viper.GetString("default_album_name"))
	photoLocalRepo := local.NewPhotoRepo(viper.GetString("local_storage"))
	albumRepo := mysqlPhoto.NewAlbumRepo(db)
	rpcClient := pb.NewAuthServiceClient(rpcConn)

	return &App{
		authUC: usecase.NewAuthUsecase(
			userRepo,
			rpcClient,
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

	// Set up http handlers
	// SignUp/SignIn endpoints
	httpAuth.RegisterRoutes(r, a.authUC)

	// User endpoint
	authMiddleware := httpAuth.NewAuthMiddleware(a.authUC)
	user := r.Group("/user", authMiddleware)
	httpAuth.RegisterMidRoutes(user, a.authUC)

	// API endpoints
	api := r.Group("/api", authMiddleware)
	httpPhoto.RegisterRoutes(api, a.photoUC)

	// Swagger
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// HTTP Server
	a.httpSrv = &http.Server{
		Addr:           ":" + port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host=localhost
	if err := a.httpSrv.ListenAndServe(); err != nil {
		// if err := a.httpSrv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
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

// func initGormDB(db *sql.DB) *gorm.DB {
// 	gormDB, err := gorm.Open(mysql.New(mysql.Config{
// 		Conn: db,
// 	}), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return gormDB
// }

func initRpcClient() *grpc.ClientConn {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(viper.GetString("grpc_auth_srv"), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	return conn
}
