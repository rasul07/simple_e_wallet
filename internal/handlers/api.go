package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rasul07/alif-task/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/rasul07/alif-task/api/docs"
)

type API struct {
	router        *gin.Engine
	walletService service.WalletService
}

func NewAPI(walletService service.WalletService) *API {
	api := &API{
		router:        gin.New(),
		walletService: walletService,
	}

	api.setupRoutes()
	return api
}

func (api *API) setupRoutes() {
	cfg := cors.DefaultConfig()
	cfg.AllowHeaders = append(cfg.AllowHeaders, "*")
	cfg.AllowAllOrigins = true
	cfg.AllowCredentials = true
	api.router.Use(cors.New(cfg))

	handler := NewHandler(api.walletService)

	v1 := api.router.Group("/v1")
	v1.Use(AuthMiddleware())
	{
		v1.POST("/wallet/check", handler.CheckWalletExists)
		v1.POST("/wallet/topup", handler.TopUpWallet)
		v1.POST("/wallet/transactions", handler.GetTransactions)
		v1.POST("/wallet/balance", handler.GetBalance)
	}

	url := ginSwagger.URL("swagger/doc.json")
	api.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

func (api *API) Run(addr string) error {
	return api.router.Run(addr)
}
