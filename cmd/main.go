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
	teamdomain "github.com/naufalfmm/gdc-task-management/internal/domains/teams"
	userdomain "github.com/naufalfmm/gdc-task-management/internal/domains/users"
	authctrl "github.com/naufalfmm/gdc-task-management/internal/infrastructures/rest/controllers/auth"
	"github.com/naufalfmm/gdc-task-management/internal/infrastructures/rest/routes"
	"github.com/naufalfmm/gdc-task-management/internal/middlewares"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/authtoken/jwt"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/password/bcrypt"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/txmanager/pgx"
	v10 "github.com/naufalfmm/gdc-task-management/internal/pkg/validator/v10"
	teamrepo "github.com/naufalfmm/gdc-task-management/internal/repositories/teams"
	userrepo "github.com/naufalfmm/gdc-task-management/internal/repositories/users"
	authusc "github.com/naufalfmm/gdc-task-management/internal/usecases/auth"
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

	userQuerier := userdomain.New(dbpool)
	userRepositories := userrepo.NewRepository(userQuerier)

	pgxtxmanager := pgx.NewPgxTxManager(dbpool)

	authUsecase := authusc.NewAuthUsecase(pgxtxmanager, teamRepositories, userRepositories, passHasher, jwtoken)

	middleware := middlewares.NewMiddleware()
	authControllers := authctrl.NewAuthController(authUsecase)

	routes := routes.NewRoutes(authControllers, middleware)
	routes.Register(ge)

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", conf.ServicePort),
		Handler: ge,
	}

	binding.Validator = validator

	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}

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
