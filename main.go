package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/maxik12233/blog/controllers"
	"github.com/maxik12233/blog/services"
	"github.com/maxik12233/blog/types"
	"github.com/maxik12233/blog/user"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	router    *gin.Engine
	muxrouter *mux.Router
	db        *gorm.DB
	logger    *zap.Logger
	err       error

	resservice    services.ResourceService
	rescontroller controllers.ResourceController
)

func InitializeLogger() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	logFile, _ := os.OpenFile("log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	logger = zap.New(core, zap.AddCaller())
}

func initialMigration() {

	db, err = gorm.Open(postgres.Open(os.Getenv("DB_URL")), &gorm.Config{})
	if err != nil {
		logger.Fatal("Unable to connect to database with gorm")
		os.Exit(1)
	}

	err = db.AutoMigrate(&user.User{}, &user.ContactInfo{}, &user.Location{}, &user.PersonalInfo{}, &types.Role{}, &types.Article{}, &types.Comment{}, &types.Like{})
	if err != nil {
		logger.Fatal("Failed automigration")
		os.Exit(1)
	}
}

func fillRoleData() {
	db.Create([]types.RoleData{
		{RoleName: "common"},
		{RoleName: "moderator"},
		{RoleName: "admin"},
	})
}

func main() {

	InitializeLogger()

	if err := godotenv.Load(); err != nil {
		logger.Fatal(err.Error())
	}

	initialMigration()
	defer func() {
		dbInstance, _ := db.DB()
		_ = dbInstance.Close()
	}()
	fillRoleData()

	router = gin.Default()
	muxrouter = mux.NewRouter()

	basepathGin := router.Group("/v1")
	basepathMux := muxrouter.PathPrefix("/v1").Subrouter()

	// res TODO rewrite to microservice
	resservice = services.NewResourceService(db, logger.With(zap.String("service", "resource_service")))
	rescontroller = controllers.NewResourceController(resservice)
	rescontroller.RegisterResourceRoutes(basepathGin)

	// user microservice
	repo := user.NewUserRepo(db, logger.With(zap.String("service", "user_repository")))
	svc := user.NewUserService(repo, logger.With(zap.String("service", "user_service")))
	userEndpoints := user.MakeUserEndpoints(svc)
	user.CreateNewServer(basepathMux, userEndpoints)

	var httpAddr = flag.String("http", os.Getenv("PORT"), "http lister address")
	// Start the servers
	go func() {
		fmt.Println("listening on port: ", *httpAddr)
		http.ListenAndServe(*httpAddr, muxrouter)
	}()
	router.Run(":6060")
}
