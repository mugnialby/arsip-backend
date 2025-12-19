package api

import (
	"regexp"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/api/handler"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/config"
	"github.com/mugnialby/perpustakaan-kejari-kota-bogor-backend/internal/service"
)

func NewRouter(
	userService *service.UserService,
	roleService *service.RoleService,
	departmentService *service.DepartmentService,
	// authService *service.AuthService,
	archiveService *service.ArchiveService,
	archiveAttachmentService *service.ArchiveAttachmentService,
	archiveTypeService *service.ArchiveTypeService,
	archiveCharacteristicService *service.ArchiveCharacteristicService,
	archiveRoleAccessService *service.ArchiveRoleAccessService,
) *gin.Engine {
	r := gin.Default()

	// Apply CORS middleware FIRST
	r.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"http://localhost:3000", "http://192.168.50.52:3000"},
		AllowOriginFunc: func(origin string) bool {
			// Always allow localhost for development
			if origin == "http://localhost:3000" {
				return true
			}

			// Allow any local IP in 192.168.x.x or 10.x.x.x or 172.16–31.x.x range
			matched192, _ := regexp.MatchString(`^http://192\.168\.\d+\.\d+:3000$`, origin)
			matched10, _ := regexp.MatchString(`^http://10\.\d+\.\d+\.\d+:3000$`, origin)
			matched172, _ := regexp.MatchString(`^http://172\.(1[6-9]|2[0-9]|3[0-1])\.\d+\.\d+:3000$`, origin)

			return matched192 || matched10 || matched172
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Access-Control-Allow-Origin", "Origin", "Content-Type", "Accept", "Authorization", "Content-Disposition", "Cache-Control"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Optional manual handler for preflight
	r.OPTIONS("/*any", func(c *gin.Context) {
		c.Status(204)
	})

	config.Load()
	// jwtService := utils.NewJWTService(cfg.JWTSecret, cfg.AppName)

	userHandler := handler.NewUserHandler(userService)
	roleHandler := handler.NewRoleHandler(roleService)
	departmentHandler := handler.NewDepartmentHandler(departmentService)
	authHandler := handler.NewAuthHandler(userService)
	archiveHandler := handler.NewArchiveHandler(archiveService, archiveAttachmentService, archiveRoleAccessService)
	archiveTypeHandler := handler.NewArchiveTypeHandler(archiveTypeService)
	archiveCharacteristicHandler := handler.NewArchiveCharacteristicHandler(archiveCharacteristicService)

	api := r.Group("/api")
	{
		// Public routes
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/authenticate", authHandler.Login)

			// api.POST("/login", func(c *gin.Context) {
			// Here you’d normally validate credentials, but we mock it for example:
			// token, _ := jwtService.GenerateToken(1) // 1 = user ID
			// c.JSON(200, gin.H{"token": token})
			// })
		}

		master := api.Group("/master")
		// master.Use(middleware.JWTAuth(jwtService))
		{
			users := master.Group("/users")
			{
				users.GET("/", userHandler.GetAllUsers)
				users.GET("/:id", userHandler.GetUserByID)
				users.POST("/", userHandler.CreateUser)
				users.PUT("/", userHandler.UpdateUserById)
				users.PATCH("/", userHandler.DeleteUserById)
			}

			roles := master.Group("/roles")
			{
				roles.GET("/", roleHandler.GetAllRoles)
				roles.GET("/:id", roleHandler.GetRoleByID)
				roles.POST("/", roleHandler.CreateRole)
				roles.PUT("/", roleHandler.UpdateRoleById)
				roles.PATCH("/", roleHandler.DeleteRoleById)
				roles.GET("/findByQuery/department/:id", roleHandler.GetRoleByDepartmentID)
			}

			department := master.Group("/departments")
			{
				department.GET("/", departmentHandler.GetAllDepartments)
				department.GET("/:id", departmentHandler.GetDepartmentByID)
				department.POST("/", departmentHandler.CreateDepartment)
				department.PUT("/", departmentHandler.UpdateDepartmentById)
				department.PATCH("/", departmentHandler.DeleteDepartmentById)
			}

			archiveType := master.Group("/archiveTypes")
			{
				archiveType.GET("/", archiveTypeHandler.GetAllArchiveTypes)
				archiveType.GET("/:id", archiveTypeHandler.GetArchiveTypeByID)
				archiveType.POST("/", archiveTypeHandler.CreateArchiveType)
				archiveType.PUT("/", archiveTypeHandler.UpdateArchiveTypeById)
				archiveType.PATCH("/", archiveTypeHandler.DeleteArchiveTypeById)
			}

			archiveCharacteristic := master.Group("/archiveCharacteristics")
			{
				archiveCharacteristic.GET("/", archiveCharacteristicHandler.GetAllArchiveCharacteristics)
				archiveCharacteristic.GET("/:id", archiveCharacteristicHandler.GetArchiveCharacteristicByID)
				archiveCharacteristic.POST("/", archiveCharacteristicHandler.CreateArchiveCharacteristic)
				archiveCharacteristic.PUT("/", archiveCharacteristicHandler.UpdateArchiveCharacteristicById)
				archiveCharacteristic.PATCH("/", archiveCharacteristicHandler.DeleteArchiveCharacteristicById)
			}
		}

		archives := api.Group("/archives")
		{
			archives.GET("/", archiveHandler.GetAllArchives)
			archives.POST("/getByData", archiveHandler.GetAllArchivesByData)
			archives.GET("/:id", archiveHandler.GetArchiveByID)
			archives.POST("/", archiveHandler.CreateArchive)
			archives.PUT("/", archiveHandler.UpdateArchiveById)
			archives.PATCH("/", archiveHandler.DeleteArchiveById)
			archives.GET("/find/:query", archiveHandler.FindArchiveByQuery)
			archives.POST("/findByQuery/advanced", archiveHandler.FindArchiveByAdvanceQuery)
			archives.GET("/:id/pdf", archiveHandler.StreamMergedPDF)
		}
	}

	return r
}

// userID, exists := c.Get("user_id")
// if !exists {
//     c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
//     return
// }
