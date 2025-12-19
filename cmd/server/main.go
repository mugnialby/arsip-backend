package main

import (
	"log"

	"github.com/mugnialby/arsip-backend/internal/api"
	"github.com/mugnialby/arsip-backend/internal/appcontext"
	"github.com/mugnialby/arsip-backend/internal/config"
	"github.com/mugnialby/arsip-backend/internal/repository"
	"github.com/mugnialby/arsip-backend/internal/service"
	"github.com/mugnialby/arsip-backend/pkg/logger"
)

func main() {
	cfg := config.Load()
	logger.InitLogger()
	defer logger.Sync()

	ctx, err := appcontext.NewAppContext(cfg)
	if err != nil {
		log.Fatal(err)
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

	// --- API Router
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

	logger.Log.Infof("🚀 Server starting in %s mode on port %s", cfg.AppEnv, cfg.Port)
	router.Run(":" + cfg.Port)
}
