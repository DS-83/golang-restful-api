package server

import (
	"context"
	"database/sql"
	"example-restful-api-server/auth"
	httpAuth "example-restful-api-server/auth/delivery/http"
	gormUser "example-restful-api-server/auth/repo/gorm"
	mongoToken "example-restful-api-server/auth/repo/mongodb"
	"example-restful-api-server/auth/usecase"
	"example-restful-api-server/photogramm"
	httpPhoto "example-restful-api-server/photogramm/delivery/http"
	gormPhoto "example-restful-api-server/photogramm/repo/gorm"
	"example-restful-api-server/photogramm/repo/local"
	photoUsecase "example-restful-api-server/photogramm/usecase"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "example-restful-api-server/docs"

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
	sqlDB := initDB()
	db := initGormDB(sqlDB)
	mongoDB := initMongoDB()

	userRepo := gormUser.NewUserRepo(db, viper.GetString("default_album_name"))
	photoDBRepo := gormPhoto.NewPhotoRepo(db, viper.GetString("default_album_name"))
	photoLocalRepo := local.NewPhotoRepo(viper.GetString("local_storage"))
	albumRepo := gormPhoto.NewAlbumRepo(db)
	tokenRepo := mongoToken.NewTokenRepo(mongoDB)

	return &App{
		authUC: usecase.NewAuthUsecase(
			userRepo,
			tokenRepo,
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
	delRoute := r.Group("/user", authMiddleware)
	httpAuth.RegisterMidRoutes(delRoute, a.authUC)

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

func initGormDB(db *sql.DB) *gorm.DB {
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return gormDB
}

func initMongoDB() *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(viper.GetString("mongo.uri")))

	if err != nil {
		log.Fatal(err)
	}

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	return client.Database(viper.GetString("mongo.db"))
}
