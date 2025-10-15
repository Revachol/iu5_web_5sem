package handler

import (
	_ "github.com/Revachol/iu5_web_5sem/docs"
	"github.com/Revachol/iu5_web_5sem/internal/app/config"
	"github.com/Revachol/iu5_web_5sem/internal/app/redis"
	"github.com/Revachol/iu5_web_5sem/internal/app/repository"
	"github.com/Revachol/iu5_web_5sem/internal/app/role"

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
	// ---------------------------
	// Публичные маршруты
	// ---------------------------
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", h.GetHistoricalObjects)
	router.GET("/historical_object/:id", h.GetHistoricalObject)
	router.GET("/api/historical_object/:id", h.GetHistoricalObjectAPI) //GET одна запись
	router.GET("/api/historical_objects", h.GetHistoricalObjectsAPI)   //GET список с фильтрацией
	router.POST("/api/users/register", h.RegisterUserAPI)              // 	POST регистрация
	router.POST("/api/users/login", h.LoginUserAPI)                    // POST аутентификация
	router.POST("/api/users/logout", h.LogoutUserAPI)                  // POST деавторизация
	// router.POST("/add_to_cart/:id", h.AddServiceToRequest)
	// router.POST("/delete_request/:id", h.DeleteRequest)
	//router.GET("/historical_estimate/:id", h.GetHistoricalEstimate)

	// ---------------------------
	// Защищённые маршруты
	// ---------------------------

	// Все авторизованные пользователи (User и Admin)
	auth := router.Group("/api")
	auth.Use(h.AuthMiddleware(role.User, role.Admin))
	{

		auth.POST("/add_to_estimate/:id", h.AddHistoricalObjecsToRequestAPI)                                 //POST добавления в заявку-черновик.
		auth.GET("/historical_estimate_count", h.GetDraftRequestAPI)                                         //GET иконки корзины (без входных параметров, ид заявки вычисляется).
		auth.GET("/historical_estimate", h.GetAllHistoricalEstimateAPI)                                      //GET список (кроме удаленных и черновика, поля модератора и создателя через логины) с фильтрацией по диапазону даты формирования и статусу
		auth.GET("/historical_estimate/:id", h.GetHistoricalEstimateAPI)                                     //GET одна запись (поля заявки + ее услуги)
		auth.PUT("/historical_estimate/:id", h.UpdateHistoricalEstimateAPI)                                  //PUT изменения полей заявки по теме
		auth.PUT("/historical_estimate/:id/form", h.FormHistoricalEstimateAPI)                               //PUT сформировать создателем (дата формирования).
		auth.DELETE("/historical_estimate/:id", h.DeleteHistoricalEstimeteAPI)                               //DELETE удаление (дата формирования)
		auth.DELETE("/estimate/:estimate_id/historical_objects/:object_id", h.DeleteHObjectFromHEstimateAPI) //DELETE удаление услуги из заявки
		auth.PUT("/estimate/:estimate_id/historical_objects/:object_id/quantity", h.UpdateQuantityAPI)       //PUT изменение количества/порядка/значения в м-м (без PK м-м)
	}

	moder := router.Group("/api")
	moder.Use(h.AuthMiddleware(role.Admin))
	{
		moder.POST("/create_object", h.CreateHistoricalObjectAPI)                     //POST добавление (без изображения)
		moder.POST("/historical_object/:id/image", h.UploadHistoricalObjectImage)     //POST добавление изображения
		moder.PUT("/historical_object/:id", h.UpdateHistoricalObjectAPI)              //PUT изменение
		moder.DELETE("/historical_object/:id", h.DeleteHistoricalObjectAPI)           //DELETE удаление. Удаление изображения встроено в метод удаления услуги
		moder.PUT("/historical_estimate/:id/complete", h.FinishHistoricalEstimateAPI) //PUT завершить/отклонить модератором.
		moder.GET("/users/:id", h.GetUserAPI)                                         // GET полей пользователя после аутентификации (для личного кабинета)
		moder.PUT("/users/:id", h.UpdateUserAPI)                                      // PUT пользователя (личный кабинет)

	}
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
