package main

import (
	"os"

	"github.com/mugnialby/arsip-backend/internal/api"
	"github.com/mugnialby/arsip-backend/internal/appcontext"
	"github.com/mugnialby/arsip-backend/internal/config"
	"github.com/mugnialby/arsip-backend/internal/repository"
	"github.com/mugnialby/arsip-backend/internal/service"
	"github.com/mugnialby/arsip-backend/internal/utils"
	"github.com/mugnialby/arsip-backend/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	/*------ LOGGER ------*/
	logger.Init()
	defer logger.Log.Sync()

	/*------ CONFIG ------*/
	cfg := config.Load()

	ctx, err := appcontext.NewAppContext(cfg)
	if err != nil {
		logger.Log.Error("main.failed",
			zap.Error(err),
		)
	}

	/*------ SERVICES ------*/
	archiveRepo := repository.NewArchiveRepository(ctx.DB)
	archiveService := service.NewArchiveService(archiveRepo)

	archiveAttachmentRepo := repository.NewArchiveAttachmentRepository(ctx.DB)
	archiveAttachmentService := service.NewArchiveAttachmentService(archiveAttachmentRepo)

	userRepo := repository.NewUserRepository(ctx.DB)
	userService := service.NewUserService(userRepo)

	roleRepo := repository.NewRoleRepository(ctx.DB)
	roleService := service.NewRoleService(roleRepo)

	departmentRepo := repository.NewDepartmentRepository(ctx.DB)
	departmentService := service.NewDepartmentService(departmentRepo)

	archiveTypeRepo := repository.NewArchiveTypeRepository(ctx.DB)
	archiveTypeService := service.NewArchiveTypeService(archiveTypeRepo)

	archiveCharacteristicRepo := repository.NewArchiveCharacteristicRepository(ctx.DB)
	archiveCharacteristicService := service.NewArchiveCharacteristicService(archiveCharacteristicRepo)

	archiveRoleAccessRepo := repository.NewArchiveRoleAccessRepository(ctx.DB)
	archiveRoleAccessService := service.NewArchiveRoleAccessService(archiveRoleAccessRepo)

	/*------ ROUTERS ------*/
	router := api.NewRouter(
		userService,
		roleService,
		departmentService,
		archiveService,
		archiveAttachmentService,
		archiveTypeService,
		archiveCharacteristicService,
		archiveRoleAccessService,
	)

	logger.Log.Info("main.success",
		zap.Any("message", "Server starting in"+cfg.AppEnv+" mode on port "+cfg.Port),
	)

	workingDir, err := os.Getwd()
	if err != nil {
		logger.Log.Error("config.env.get_storage_location.failed",
			zap.Error(err),
			zap.String("storage_location", workingDir),
			zap.String("message", "Failed to get storage path"),
		)
	}
	logger.Log.Info(workingDir)

	storageLocation, err := utils.GetStorageLocation()
	if err != nil {
		logger.Log.Error("config.env.get_storage_location.failed",
			zap.Error(err),
			zap.String("storage_location", storageLocation),
			zap.String("message", "Failed to get storage path"),
		)
	}
	logger.Log.Info(storageLocation)

	router.Run(":" + cfg.Port)
}
