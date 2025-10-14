package handler

import (
	_ "github.com/Revachol/iu5_web_5sem/docs"
	"github.com/Revachol/iu5_web_5sem/internal/app/config"
	"github.com/Revachol/iu5_web_5sem/internal/app/redis"
	"github.com/Revachol/iu5_web_5sem/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	Repository *repository.Repository
	Config     *config.Config
	Redis      *redis.Client
}

func NewHandler(r *repository.Repository, cfg *config.Config, redis *redis.Client) *Handler {
	return &Handler{
		Repository: r,
		Config:     cfg,
		Redis:      redis,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", h.GetHistoricalObjects)
	router.GET("/historical_object/:id", h.GetHistoricalObject)
	router.GET("/api/historical_object/:id", h.GetHistoricalObjectAPI) //GET одна запись
	router.GET("/api/historical_objects", h.GetHistoricalObjectsAPI)   //GET список с фильтрацией

	router.POST("/add_to_cart/:id", h.AddServiceToRequest)
	router.POST("/delete_request/:id", h.DeleteRequest)
	router.POST("/api/create_object", h.CreateHistoricalObjectAPI)                 //POST добавление (без изображения)
	router.POST("/api/add_to_estimate/:id", h.AddHistoricalObjecsToRequestAPI)     //POST добавления в заявку-черновик.
	router.POST("/api/historical_object/:id/image", h.UploadHistoricalObjectImage) //POST добавление изображения

	router.PUT("/api/historical_object/:id", h.UpdateHistoricalObjectAPI)    //PUT изменение
	router.DELETE("/api/historical_object/:id", h.DeleteHistoricalObjectAPI) //DELETE удаление. Удаление изображения встроено в метод удаления услуги

	// Заявка
	router.GET("/historical_estimate/:id", h.GetHistoricalEstimate)
	router.GET("/api/historical_estimate_count", h.GetDraftRequestAPI)     //GET иконки корзины (без входных параметров, ид заявки вычисляется).
	router.GET("/api/historical_estimate", h.GetAllHistoricalEstimateAPI)  //GET список (кроме удаленных и черновика, поля модератора и создателя через логины) с фильтрацией по диапазону даты формирования и статусу
	router.GET("/api/historical_estimate/:id", h.GetHistoricalEstimateAPI) //GET одна запись (поля заявки + ее услуги)

	router.PUT("/api/historical_estimate/:id", h.UpdateHistoricalEstimateAPI)          //PUT изменения полей заявки по теме
	router.PUT("/api/historical_estimate/:id/form", h.FormHistoricalEstimateAPI)       //PUT сформировать создателем (дата формирования).
	router.PUT("/api/historical_estimate/:id/complete", h.FinishHistoricalEstimateAPI) //PUT завершить/отклонить модератором.
	router.DELETE("/api/historical_estimate/:id", h.DeleteHistoricalEstimeteAPI)       //DELETE удаление (дата формирования)

	//М-М
	router.DELETE("/api/estimate/:estimate_id/historical_objects/:object_id", h.DeleteHObjectFromHEstimateAPI) //DELETE удаление услуги из заявки
	router.PUT("/api/estimate/historical_objects/:estimate_id/:object_id/quantity", h.UpdateQuantityAPI)       //PUT изменение количества/порядка/значения в м-м (без PK м-м)

	//Users
	router.POST("/api/users/register", h.RegisterUserAPI) // 	POST регистрация
	router.GET("/api/users/:id", h.GetUserAPI)            // GET полей пользователя после аутентификации (для личного кабинета)
	router.PUT("/api/users/:id", h.UpdateUserAPI)         // PUT пользователя (личный кабинет)
	router.POST("/api/users/login", h.LoginUserAPI)       // POST аутентификация
	router.POST("/api/users/logout", h.LogoutUserAPI)     // POST деавторизация

}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
