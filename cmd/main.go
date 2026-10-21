package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/naufalfmm/gdc-task-management/config"
	"github.com/naufalfmm/gdc-task-management/internal/domains/idempotencykeys"
	"github.com/naufalfmm/gdc-task-management/internal/domains/tasks"
	teamenumdomain "github.com/naufalfmm/gdc-task-management/internal/domains/teamenums"
	teamdomain "github.com/naufalfmm/gdc-task-management/internal/domains/teams"
	userdomain "github.com/naufalfmm/gdc-task-management/internal/domains/users"
	authctrl "github.com/naufalfmm/gdc-task-management/internal/infrastructures/rest/controllers/auth"
	taskctrl "github.com/naufalfmm/gdc-task-management/internal/infrastructures/rest/controllers/tasks"
	"github.com/naufalfmm/gdc-task-management/internal/infrastructures/rest/routes"
	"github.com/naufalfmm/gdc-task-management/internal/middlewares"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/authtoken/jwt"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/logger/slog"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/notification/dummy"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/password/bcrypt"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/txmanager/pgx"
	v10 "github.com/naufalfmm/gdc-task-management/internal/pkg/validator/v10"
	idemkeyrepo "github.com/naufalfmm/gdc-task-management/internal/repositories/idempotencykeys"
	taskrepo "github.com/naufalfmm/gdc-task-management/internal/repositories/tasks"
	teamenumrepo "github.com/naufalfmm/gdc-task-management/internal/repositories/teamenums"
	teamrepo "github.com/naufalfmm/gdc-task-management/internal/repositories/teams"
	userrepo "github.com/naufalfmm/gdc-task-management/internal/repositories/users"
	authusc "github.com/naufalfmm/gdc-task-management/internal/usecases/auth"
	taskusc "github.com/naufalfmm/gdc-task-management/internal/usecases/tasks"
)

func main() {
	ge := gin.New()

	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	passHasher := bcrypt.NewBcryptHasher(conf.BcryptCost)
	jwtoken := jwt.NewJWT(conf.JwtPrivateKey, conf.JwtPublicKey, conf.JwtAlgorithm, conf.JwtAccessTokenExpiration, conf.JwtRefreshTokenExpiration)
	validator, err := v10.NewV10()
	if err != nil {
		panic(err)
	}

	logger := slog.New()
	notif := dummy.NewDummyNotification(logger)

	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, conf.GetDSN())
	if err != nil {
		panic(err)
	}
	defer dbpool.Close()

	if err := dbpool.Ping(ctx); err != nil {
		log.Fatalf("Database ping failed: %v\n", err)
	}

	teamQuerier := teamdomain.New(dbpool)
	teamRepositories := teamrepo.NewRepository(teamQuerier)

	teamEnumQuerier := teamenumdomain.New(dbpool)
	teamEnumRepositories := teamenumrepo.NewRepository(teamEnumQuerier)

	userQuerier := userdomain.New(dbpool)
	userRepositories := userrepo.NewRepository(userQuerier)

	idempotencyKeyQuerier := idempotencykeys.New(dbpool)
	idempotencyKeyRepositories := idemkeyrepo.NewRepository(idempotencyKeyQuerier)

	taskQuerier := tasks.New(dbpool)
	taskRepositories := taskrepo.NewRepository(taskQuerier)

	pgxtxmanager := pgx.NewPgxTxManager(dbpool)

	authUsecase := authusc.NewAuthUsecase(pgxtxmanager, teamRepositories, userRepositories, passHasher, jwtoken)
	taskUsecase := taskusc.NewTaskUsecase(pgxtxmanager, notif, taskRepositories, idempotencyKeyRepositories, userRepositories, teamRepositories, teamEnumRepositories, *conf)

	middleware := middlewares.NewMiddleware(jwtoken)
	authControllers := authctrl.NewAuthController(authUsecase)

	taskControllers := taskctrl.NewTaskController(taskUsecase)

	routes := routes.NewRoutes(authControllers, taskControllers, middleware)
	routes.Register(ge)

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", conf.ServicePort),
		Handler: ge,
	}

	binding.Validator = validator

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(
		sc,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	<-sc

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctxShutDown); err != nil {
		panic(err)
	}
}
