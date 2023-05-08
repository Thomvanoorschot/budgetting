package http

import (
	"budgetting/api/http/handler"
	"budgetting/api/http/routes"
	"budgetting/config"
	"budgetting/repositories/planetscale"
	"budgetting/services/banking"
	"budgetting/services/banking/nordigen"
	"budgetting/services/profile"
	"budgetting/services/profile/auth0"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func ListenAndServe(config *config.Config) error {
	e := gin.Default()
	h := createHandler(config)
	r := routes.NewRouter(config)
	r.SetupRoutes(e, h)
	return e.Run(fmt.Sprintf("%s:%s", config.ApiHost, config.ApiPort))
}

func createHandler(config *config.Config) *handler.Handler {
	repository, err := planetscale.NewRepository(config)

	nordigenClient := nordigen.NewClient(config)
	bankingService := banking.NewService(nordigenClient, repository)

	if err != nil {
		log.Fatal(err)
		return nil
	}
	auth0Client := auth0.NewClient(config)
	profileService := profile.NewService(repository, auth0Client)

	return handler.NewHandler(bankingService, profileService)
}
