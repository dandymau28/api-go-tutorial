package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/api-go/config"
	"gitlab.com/api-go/controller"
	"gitlab.com/api-go/middleware"
	"gitlab.com/api-go/repository"
	"gitlab.com/api-go/service"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB                  = config.DatabaseConnection()
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	bookRepository repository.BookRepository = repository.NewBookRepository(db)
	jwtService     service.JWTService        = service.NewJWTService()
	userService    service.UserService       = service.NewUserService(userRepository)
	authService    service.AuthService       = service.NewAuthService(userRepository)
	bookService    service.BookService       = service.NewBookService(bookRepository)
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
	userController controller.UserController = controller.NewUserController(userService, jwtService)
	bookController controller.BookController = controller.NewBookController(bookService, jwtService)
)

func main() {
	defer config.CloseConnection(db)
	r := gin.Default()

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user")
	{
		userRoutes.PUT("/update", userController.Update)
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.GET("/", userController.All)
	}

	bookRoutes := r.Group("api/books", middleware.AuthorizeJWT(jwtService))
	{
		bookRoutes.POST("/", bookController.Insert)
		bookRoutes.PUT("/:id", bookController.Update)
		bookRoutes.DELETE("/:id", bookController.Delete)
		bookRoutes.GET("/", bookController.All)
		bookRoutes.GET("/:id", bookController.FindById)
	}

	r.Run()

}
