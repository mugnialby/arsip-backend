package handler

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mugnialby/arsip-backend/internal/model"
	archiveRequest "github.com/mugnialby/arsip-backend/internal/model/dto/request/archive"
	archiveRoleAccessRequest "github.com/mugnialby/arsip-backend/internal/model/dto/request/archiveRoleAccess"
	"github.com/mugnialby/arsip-backend/internal/service"
	"github.com/mugnialby/arsip-backend/internal/utils"
	"github.com/mugnialby/arsip-backend/pkg/logger"
	"github.com/mugnialby/arsip-backend/pkg/response"
	"go.uber.org/zap"
)

type ArchiveHandler struct {
	archiveService           *service.ArchiveService
	archiveAttachmentService *service.ArchiveAttachmentService
	archiveRoleAccessService *service.ArchiveRoleAccessService
}

var CacheTTL = 5 * time.Minute

type CacheItem struct {
	Data      []byte
	ExpiresAt time.Time
}

func NewArchiveHandler(
	archiveService *service.ArchiveService,
	archiveAttachmentService *service.ArchiveAttachmentService,
	archiveRoleAccessService *service.ArchiveRoleAccessService,
) *ArchiveHandler {
	return &ArchiveHandler{
		archiveService:           archiveService,
		archiveAttachmentService: archiveAttachmentService,
		archiveRoleAccessService: archiveRoleAccessService,
	}
}

func (h *ArchiveHandler) GetAllArchives(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	archives, err := h.archiveService.GetAllArchives()
	if err != nil {
		logger.Log.Error("archive.get_all.failed",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to get all data")
		return
	}

	logger.Log.Info("archive.get_all.success",
		zap.String("request_id", requestID.(string)),
		zap.Int("count", len(archives)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, archives)
}

func (h *ArchiveHandler) GetAllArchivesByData(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var getArchiveByDataRequest archiveRequest.GetArchiveByDataRequest
	if err := c.ShouldBindJSON(&getArchiveByDataRequest); err != nil {
		logger.Log.Warn("archive.get_by_data.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
		)

		response.Error(c, http.StatusBadRequest, "JSON request is not valid")
		return
	}

	archives, err := h.archiveService.GetAllArchivesByData(getArchiveByDataRequest)
	if err != nil {
		logger.Log.Error("archive.get_by_data.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("filters", getArchiveByDataRequest),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to get all data")
		return
	}

	logger.Log.Info("archive.get_by_data.success",
		zap.String("request_id", requestID.(string)),
		zap.Int("count", len(archives)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, archives)
}

func (h *ArchiveHandler) GetArchiveByID(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.Warn("archive.get_by_id.invalid_id",
			zap.String("request_id", requestID.(string)),
			zap.String("param", c.Param("id")),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	archive, err := h.archiveService.GetArchiveByID(uint(id))
	if err != nil {
		logger.Log.Info("archive.get_by_id.failed",
			zap.String("request_id", requestID.(string)),
			zap.Uint("archive_id", uint(id)),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusNotFound, "Failed to get data")
		return
	}

	logger.Log.Info("archive.get_by_id.archive_hdr.success",
		zap.String("request_id", requestID.(string)),
		zap.Uint("archive_id", uint(id)),
		zap.Int("attachments", len(archive.ArchiveAttachments)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	for _, archiveAttachment := range archive.ArchiveAttachments {
		if archiveAttachment.FileLocation == "" {
			continue
		}

		fileBytes, err := os.ReadFile(archiveAttachment.FileLocation)
		if err != nil {
			archiveAttachment.FileBase64 = ""
			continue
		}

		// Encode to Base64
		base64Data := base64.StdEncoding.EncodeToString(fileBytes)

		ext := strings.ToLower(filepath.Ext(archiveAttachment.FileName))

		switch ext {
		case ".jpg", ".jpeg", ".png", ".gif", ".webp":
			archiveAttachment.FileBase64 = "data:image/" + strings.TrimPrefix(ext, ".") + ";base64," + base64Data

		case ".pdf":
			archiveAttachment.FileBase64 = "data:application/pdf;base64," + base64Data

		default:
			archiveAttachment.FileBase64 = base64Data
		}
	}

	logger.Log.Info("archive.get_by_id.success",
		zap.String("request_id", requestID.(string)),
		zap.Uint("archive_id", uint(id)),
		zap.Int("attachments", len(archive.ArchiveAttachments)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, archive)
}

func (h *ArchiveHandler) CreateArchive(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var newArchiveRequest archiveRequest.NewArchiveRequest
	if err := c.ShouldBindJSON(&newArchiveRequest); err != nil {
		logger.Log.Warn("archive.create.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
		)

		response.Error(c, http.StatusBadRequest, "JSON request is not valid")
		return
	}

	newArchive := model.ArchiveHdr{
		ArchiveDate:             newArchiveRequest.ArchiveDate,
		ArchiveNumber:           newArchiveRequest.ArchiveNumber,
		ArchiveName:             newArchiveRequest.ArchiveName,
		ArchiveCharacteristicID: newArchiveRequest.ArchiveCharacteristicID,
		ArchiveTypeID:           newArchiveRequest.ArchiveTypeID,
		DepartmentID:            newArchiveRequest.DepartmentID,
		Status:                  "Y",
		CreatedBy:               newArchiveRequest.SubmittedBy,
	}

	if err := h.archiveService.CreateArchive(&newArchive); err != nil {
		logger.Log.Error("archive.create.create_archive_hdr.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("payload", newArchive),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to create archive hdr")
		return
	}

	logger.Log.Info("archive.create.archive_hdr.success",
		zap.String("request_id", requestID.(string)),
		zap.Uint("archive_id", newArchive.ID),
		zap.Int("attachments", len(newArchiveRequest.ListArchiveAttachments)),
		zap.Int("role_access", len(newArchiveRequest.RoleAccess)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	for _, roleAccess := range newArchiveRequest.RoleAccess {
		newArchiveRoleAccess := model.ArchiveRoleAccess{
			ArchiveHdrID: newArchive.ID,
			RoleID:       roleAccess.RoleID,
			DepartmentID: roleAccess.DepartmentID,
			Status:       "Y",
			CreatedBy:    newArchiveRequest.SubmittedBy,
		}

		if err := h.archiveRoleAccessService.CreateArchiveRoleAccess(&newArchiveRoleAccess); err != nil {
			logger.Log.Error("archive.create.create_archive_role_access.failed",
				zap.String("request_id", requestID.(string)),
				zap.Any("payload", roleAccess),
				zap.Error(err),
				zap.Duration("duration_ms", time.Since(start)),
			)

			response.Error(c, http.StatusInternalServerError, "Failed to create archive role access")
			return
		}
	}

	for _, archiveAttachment := range newArchiveRequest.ListArchiveAttachments {
		if archiveAttachment.IsNew {
			base64Data := archiveAttachment.FileBase64

			if strings.Contains(base64Data, ",") {
				parts := strings.SplitN(base64Data, ",", 2)
				base64Data = parts[1]
			}

			fileExt := DetectBase64Extension(archiveAttachment.FileBase64)
			if fileExt == "" {
				logger.Log.Error("archive.create.detect_base64_extension.failed",
					zap.String("request_id", requestID.(string)),
					zap.Duration("duration_ms", time.Since(start)),
				)

				response.Error(c, http.StatusBadRequest, "File extension not found")
				return
			}

			if !isAllowedFileType(fileExt) {
				logger.Log.Error("archive.create.file_extension_not_allowed.failed",
					zap.String("request_id", requestID.(string)),
					zap.Duration("duration_ms", time.Since(start)),
				)

				response.Error(c, http.StatusBadRequest, "File extension not allowed")
				return
			}

			decodedBytes, err := base64.StdEncoding.DecodeString(base64Data)
			if err != nil {
				logger.Log.Error("archive.create.base64_invalid.failed",
					zap.String("request_id", requestID.(string)),
					zap.Duration("duration_ms", time.Since(start)),
				)

				response.Error(c, http.StatusBadRequest, "Invalid base64 file data")
				return
			}

			storageLocation, err := utils.GetStorageLocation()
			if err != nil {
				logger.Log.Error("archive.create.get_storage_location.failed",
					zap.String("request_id", requestID.(string)),
					zap.Error(err),
					zap.Duration("duration_ms", time.Since(start)),
				)

				response.Error(c, http.StatusInternalServerError, "Failed to get storage location directory")
				return
			}

			uploadDir := filepath.Join(storageLocation, "uploads", "archives", strconv.Itoa(int(newArchive.ID)))
			if err := os.MkdirAll(uploadDir, 0755); err != nil {
				logger.Log.Error("archive.create.create_upload_directory.failed",
					zap.String("request_id", requestID.(string)),
					zap.Error(err),
					zap.Duration("duration_ms", time.Since(start)),
				)

				response.Error(c, http.StatusInternalServerError, "Failed to create upload directory")
				return
			}

			fileName := fmt.Sprintf("%d_%d.%s", newArchive.ID, time.Now().UnixNano(), fileExt)
			fileLocation := filepath.Join(uploadDir, fileName)
			if err := os.WriteFile(fileLocation, decodedBytes, 0644); err != nil {
				logger.Log.Error("archive.create.write_file.failed",
					zap.String("request_id", requestID.(string)),
					zap.Error(err),
					zap.Duration("duration_ms", time.Since(start)),
				)

				response.Error(c, http.StatusInternalServerError, "Failed to write file")
				return
			}

			newArchiveAttachment := model.ArchiveAttachment{
				ArchiveHdrID: newArchive.ID,
				FileName:     fileName,
				FileLocation: fileLocation,
				Status:       "Y",
				CreatedBy:    newArchiveRequest.SubmittedBy,
			}

			if err := h.archiveAttachmentService.CreateArchiveAttachment(&newArchiveAttachment); err != nil {
				logger.Log.Error("archive.create.create_archive_attachment_data.failed",
					zap.String("request_id", requestID.(string)),
					zap.Error(err),
					zap.Duration("duration_ms", time.Since(start)),
				)

				response.Error(c, http.StatusInternalServerError, "Failed to save archive attachment")
				return
			}
		}
	}

	logger.Log.Info("archive.create.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	c.Status(http.StatusCreated)
}

func (h *ArchiveHandler) UpdateArchiveById(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var updateArchiveRequest archiveRequest.UpdateArchiveRequest
	if err := c.ShouldBindJSON(&updateArchiveRequest); err != nil {
		logger.Log.Warn("archive.update.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
		)

		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	archive, err := h.archiveService.GetArchiveByID(updateArchiveRequest.ID)
	if err != nil {
		logger.Log.Error("archive.update.get_by_id.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("payload", updateArchiveRequest.ID),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusNotFound, "Failed to get data")
		return
	}

	timeNow := time.Now()

	archive.ArchiveDate = updateArchiveRequest.ArchiveDate
	archive.ArchiveNumber = updateArchiveRequest.ArchiveNumber
	archive.ArchiveName = updateArchiveRequest.ArchiveName
	archive.ArchiveCharacteristicID = updateArchiveRequest.ArchiveCharacteristicID
	archive.ArchiveTypeID = updateArchiveRequest.ArchiveTypeID
	archive.ModifiedBy = &updateArchiveRequest.SubmittedBy
	archive.ModifiedAt = &timeNow

	if err := h.archiveService.UpdateArchive(archive); err != nil {
		logger.Log.Error("archive.update.save.failed",
			zap.String("request_id", requestID.(string)),
			zap.Any("payload", updateArchiveRequest),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to update data")
		return
	}

	for _, roleAccess := range updateArchiveRequest.RoleAccess {
		if roleAccess.IsNew {
			newArchiveRoleAccess := model.ArchiveRoleAccess{
				ArchiveHdrID: updateArchiveRequest.ID,
				RoleID:       roleAccess.RoleID,
				DepartmentID: roleAccess.DepartmentID,
				Status:       "Y",
				CreatedBy:    updateArchiveRequest.SubmittedBy,
			}

			if err := h.archiveRoleAccessService.CreateArchiveRoleAccess(&newArchiveRoleAccess); err != nil {
				logger.Log.Error("archive.update.create_role_access.failed",
					zap.String("request_id", requestID.(string)),
					zap.Any("payload", newArchiveRoleAccess),
					zap.Error(err),
					zap.Duration("duration_ms", time.Since(start)),
				)

				response.Error(c, http.StatusInternalServerError, "Failed to create data")
				return
			}
		}

		if roleAccess.IsDelete {
			deleteArchiveRoleAccess := archiveRoleAccessRequest.DeleteArchiveRoleAccessRequest{
				ID:          roleAccess.ID,
				SubmittedBy: updateArchiveRequest.SubmittedBy,
			}

			if err := h.archiveRoleAccessService.DeleteArchiveRoleAccess(&deleteArchiveRoleAccess); err != nil {
				logger.Log.Error("archive.update.delete_role_access.failed",
					zap.String("request_id", requestID.(string)),
					zap.Any("payload", deleteArchiveRoleAccess),
					zap.Error(err),
					zap.Duration("duration_ms", time.Since(start)),
				)

				response.Error(c, http.StatusInternalServerError, "Failed to delete data")
				return
			}
		}
	}

	for _, archiveAttachment := range updateArchiveRequest.ListArchiveAttachments {
		if archiveAttachment.IsNew {
			base64Data := archiveAttachment.FileBase64

			if strings.Contains(base64Data, ",") {
				parts := strings.SplitN(base64Data, ",", 2)
				base64Data = parts[1]
			}

			fileExt := DetectBase64Extension(archiveAttachment.FileBase64)
			if fileExt == "" {
				logger.Log.Error("archive.update.detect_base64_ext.failed",
					zap.String("request_id", requestID.(string)),
					zap.Duration("duration_ms", time.Since(start)),
				)

				response.Error(c, http.StatusBadRequest, "File extension not found")
				return
			}

			if !isAllowedFileType(fileExt) {
				logger.Log.Error("archive.update.file_extension.failed",
					zap.String("request_id", requestID.(string)),
					zap.Duration("duration_ms", time.Since(start)),
				)

				response.Error(c, http.StatusBadRequest, "File extension not allowed")
				return
			}

			decodedBytes, err := base64.StdEncoding.DecodeString(base64Data)
			if err != nil {
				logger.Log.Error("archive.update.base64_conversion.failed",
					zap.String("request_id", requestID.(string)),
					zap.Error(err),
					zap.Duration("duration_ms", time.Since(start)),
				)

				response.Error(c, http.StatusBadRequest, "Invalid base64 file data")
				return
			}

			storageLocation, err := utils.GetStorageLocation()
			if err != nil {
				logger.Log.Error("archive.update.get_storage_location.failed",
					zap.String("request_id", requestID.(string)),
					zap.Error(err),
					zap.Duration("duration_ms", time.Since(start)),
				)

				response.Error(c, http.StatusInternalServerError, "Failed to get storage location directory")
				return
			}

			uploadDir := filepath.Join(storageLocation, "uploads", "archives", strconv.Itoa(int(updateArchiveRequest.ID)))
			if err := os.MkdirAll(uploadDir, 0755); err != nil {
				logger.Log.Error("archive.update.create_upload_directory.failed",
					zap.String("request_id", requestID.(string)),
					zap.Error(err),
					zap.Duration("duration_ms", time.Since(start)),
				)

				response.Error(c, http.StatusInternalServerError, "Failed to create upload directory")
				return
			}

			fileName := fmt.Sprintf("%d_%d.%s", updateArchiveRequest.ID, time.Now().UnixNano(), fileExt)
			fileLocation := filepath.Join(uploadDir, fileName)
			if err := os.WriteFile(fileLocation, decodedBytes, 0644); err != nil {
				logger.Log.Error("archive.update.write_file.failed",
					zap.String("request_id", requestID.(string)),
					zap.Error(err),
					zap.Duration("duration_ms", time.Since(start)),
				)

				response.Error(c, http.StatusInternalServerError, "Failed to write file")
				return
			}

			newArchiveAttachment := model.ArchiveAttachment{
				ArchiveHdrID: updateArchiveRequest.ID,
				FileName:     fileName,
				FileLocation: fileLocation,
				Status:       "Y",
				CreatedBy:    updateArchiveRequest.SubmittedBy,
			}

			if err := h.archiveAttachmentService.CreateArchiveAttachment(&newArchiveAttachment); err != nil {
				logger.Log.Error("archive.update.create_archive_attachment.failed",
					zap.String("request_id", requestID.(string)),
					zap.Error(err),
					zap.Any("payload", &newArchiveAttachment),
					zap.Duration("duration_ms", time.Since(start)),
				)

				response.Error(c, http.StatusInternalServerError, "Failed to create data")
				return
			}

			cacheFilePath := filepath.Join(storageLocation, "cache", "archives", strconv.Itoa(int(updateArchiveRequest.ID)), fmt.Sprintf("archive_%d.pdf", updateArchiveRequest.ID))
			if err := os.Remove(cacheFilePath); err == nil {
				logger.Log.Info("archive.update.create_process.delete_cached_data.success",
					zap.String("request_id", requestID.(string)),
					zap.Any("payload", &cacheFilePath),
					zap.Duration("duration_ms", time.Since(start)),
				)
			} else if os.IsNotExist(err) {
				logger.Log.Info("archive.update.create_process.delete_cached_data.not_found",
					zap.String("request_id", requestID.(string)),
					zap.Any("payload", &cacheFilePath),
					zap.Duration("duration_ms", time.Since(start)),
				)
			} else {
				logger.Log.Error("archive.update.create_process.delete_cached_data.failed",
					zap.String("request_id", requestID.(string)),
					zap.Error(err),
					zap.Any("payload", &cacheFilePath),
					zap.Duration("duration_ms", time.Since(start)),
				)
			}
		}

		if archiveAttachment.IsDelete {
			archiveAttachment, err := h.archiveAttachmentService.GetArchiveAttachmentByID(archiveAttachment.ID)
			if err != nil {
				logger.Log.Error("archive.update.get_archive_attachment_by_id.failed",
					zap.String("request_id", requestID.(string)),
					zap.Error(err),
					zap.Any("payload", &archiveAttachment),
					zap.Duration("duration_ms", time.Since(start)),
				)

				response.Error(c, http.StatusNotFound, "Failed to get data")
				return
			}

			archiveAttachment.Status = "N"
			if err := h.archiveAttachmentService.UpdateArchiveAttachment(archiveAttachment); err != nil {
				logger.Log.Error("archive.update.update_archive_attachment.failed",
					zap.String("request_id", requestID.(string)),
					zap.Error(err),
					zap.Any("payload", archiveAttachment),
					zap.Duration("duration_ms", time.Since(start)),
				)

				response.Error(c, http.StatusInternalServerError, "Failed to update data")
				return
			}

			storageLocation, err := utils.GetStorageLocation()
			if err != nil {
				logger.Log.Error("archive.update.delete_process.get_storage_location.failed",
					zap.String("request_id", requestID.(string)),
					zap.Error(err),
					zap.Duration("duration_ms", time.Since(start)),
				)

				response.Error(c, http.StatusInternalServerError, "Failed to get storage location directory")
				return
			}

			cacheFilePath := filepath.Join(storageLocation, "cache", "archives", strconv.Itoa(int(updateArchiveRequest.ID)), fmt.Sprintf("archive_%d.pdf", updateArchiveRequest.ID))
			if err := os.Remove(cacheFilePath); err == nil {
				logger.Log.Info("archive.update.delete_process.delete_cached_data.success",
					zap.String("request_id", requestID.(string)),
					zap.Any("payload", &cacheFilePath),
					zap.Duration("duration_ms", time.Since(start)),
				)
			} else if os.IsNotExist(err) {
				logger.Log.Info("archive.update.delete_process.delete_cached_data.not_found",
					zap.String("request_id", requestID.(string)),
					zap.Any("payload", &cacheFilePath),
					zap.Duration("duration_ms", time.Since(start)),
				)
			} else {
				logger.Log.Error("archive.update.delete_process.delete_cached_data.failed",
					zap.String("request_id", requestID.(string)),
					zap.Error(err),
					zap.Any("payload", &cacheFilePath),
					zap.Duration("duration_ms", time.Since(start)),
				)
			}
		}
	}

	logger.Log.Info("archive.update.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, archive)
}

func (h *ArchiveHandler) DeleteArchiveById(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var deleteArchiveRequest archiveRequest.DeleteArchiveRequest
	if err := c.ShouldBindJSON(&deleteArchiveRequest); err != nil {
		logger.Log.Warn("archive.delete.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
		)

		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	if err := h.archiveService.DeleteArchive(&deleteArchiveRequest); err != nil {
		logger.Log.Error("archive.delete.archive.failed",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Any("payload", deleteArchiveRequest),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to delete data")
		return
	}

	if err := h.archiveAttachmentService.DeleteArchiveAttachmentByArchiveID(deleteArchiveRequest.ID, deleteArchiveRequest.SubmittedBy); err != nil {
		logger.Log.Error("archive.delete.archive_attachment.failed",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Any("payload", deleteArchiveRequest),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to delete data")
		return
	}

	if err := h.archiveRoleAccessService.DeleteArchiveRoleAccessByArchiveID(deleteArchiveRequest.ID, deleteArchiveRequest.SubmittedBy); err != nil {
		logger.Log.Error("archive.delete.archive_role_access.failed",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Any("payload", deleteArchiveRequest),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to delete data")
		return
	}

	logger.Log.Info("archive.delete.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	c.Status(http.StatusOK)
}

func (h *ArchiveHandler) FindArchiveByQuery(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	query := c.Param("query")
	if query == "" {
		logger.Log.Warn("archive.find_by_query.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Any("message", "query is not valid"),
		)

		response.Error(c, http.StatusBadRequest, "JSON request is not valid")
		return
	}

	archives, err := h.archiveService.FindArchiveByQuery(query)
	if err != nil {
		logger.Log.Error("archive.find_by_query.failed",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Any("payload", query),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusNotFound, "Failed to get data")
		return
	}

	logger.Log.Info("archive.find_by_query.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, archives)
}

func (h *ArchiveHandler) FindArchiveByAdvanceQuery(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	var advancedSearchRequest archiveRequest.AdvancedSearchRequest
	if err := c.ShouldBindJSON(&advancedSearchRequest); err != nil {
		logger.Log.Warn("archive.find_by_advance_query.invalid_request",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
		)

		response.Error(c, http.StatusBadRequest, "JSON request is not valid")
		return
	}

	archives, err := h.archiveService.FindArchiveByAdvanceQuery(advancedSearchRequest)
	if err != nil {
		logger.Log.Error("archive.find_by_advance_query.failed",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Any("payload", advancedSearchRequest),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusNotFound, "Failed to get data")
		return
	}

	logger.Log.Info("archive.find_by_query.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	response.Success(c, archives)
}

func (h *ArchiveHandler) StreamMergedPDF(c *gin.Context) {
	start := time.Now()
	requestID, _ := c.Get("request_id")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.Warn("archive.stream.invalid_id",
			zap.String("request_id", requestID.(string)),
			zap.String("param", c.Param("id")),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	archive, err := h.archiveService.GetArchiveByID(uint(id))
	if err != nil {
		logger.Log.Error("archive.stream.get_by_id.failed",
			zap.String("request_id", requestID.(string)),
			zap.String("param", c.Param("id")),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusBadRequest, "Failed to get data")
		return
	}

	if len(archive.ArchiveAttachments) == 0 {
		logger.Log.Info("archive.stream.get_by_id.no_attachment",
			zap.String("request_id", requestID.(string)),
			zap.String("param", c.Param("id")),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusBadRequest, "Attachment not found")
		return
	}

	storageLocation, err := utils.GetStorageLocation()
	if err != nil {
		logger.Log.Error("archive.stream.get_storage_location.failed",
			zap.String("request_id", requestID.(string)),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to get storage location directory")
		return
	}

	cacheDir := filepath.Join(storageLocation, "cache", "archives", strconv.Itoa(int(id)))
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		logger.Log.Error("archive.stream.create_cache_directory.failed",
			zap.String("request_id", requestID.(string)),
			zap.String("param", c.Param("id")),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)
	}

	finalPDF := filepath.Join(cacheDir, fmt.Sprintf("archive_%d.pdf", id))
	if stat, err := os.Stat(finalPDF); err == nil {
		if time.Since(stat.ModTime()) < CacheTTL {
			streamFileChunked(c, finalPDF, requestID, start)
			return
		}
	}

	var imageFiles []string
	var pdfFiles []string

	for _, att := range archive.ArchiveAttachments {
		if _, err := os.Stat(att.FileLocation); err != nil {
			continue
		}

		ext := strings.ToLower(filepath.Ext(att.FileLocation))
		switch ext {
		case ".jpg", ".jpeg", ".png", ".webp":
			imageFiles = append(imageFiles, att.FileLocation)
		case ".pdf":
			pdfFiles = append(pdfFiles, att.FileLocation)
		}
	}

	var tempImagePDF string
	if len(imageFiles) > 0 {
		tempImagePDF = filepath.Join(cacheDir, fmt.Sprintf("images_%d.pdf", id))
		if err := convertImagesToPDF(imageFiles, tempImagePDF); err != nil {
			logger.Log.Error("archive.stream.convert_images_to_pdf.failed",
				zap.String("request_id", requestID.(string)),
				zap.String("param", c.Param("id")),
				zap.Error(err),
				zap.Duration("duration_ms", time.Since(start)),
			)

			response.Error(c, http.StatusInternalServerError, "Failed to convert images to pdf")
			return
		}

		pdfFiles = append([]string{tempImagePDF}, pdfFiles...)
	}

	if len(pdfFiles) == 0 {
		logger.Log.Warn("archive.stream.pdf_files.not_found",
			zap.String("request_id", requestID.(string)),
			zap.String("param", c.Param("id")),
			zap.Duration("duration_ms", time.Since(start)),
		)
	}

	if err := mergePDFs(pdfFiles, finalPDF); err != nil {
		logger.Log.Error("archive.stream.merge_pdfs.failed",
			zap.String("request_id", requestID.(string)),
			zap.String("param", c.Param("id")),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to merge pdfs")
		return
	}

	if tempImagePDF != "" {
		_ = os.Remove(tempImagePDF)
	}

	logger.Log.Info("archive.stream.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)

	streamFileChunked(c, finalPDF, requestID, start)
}

func DetectBase64Extension(base64Str string) string {
	header := ""

	if strings.Contains(base64Str, ",") {
		parts := strings.SplitN(base64Str, ",", 2)
		header = parts[0]
	} else {
		return ""
	}

	switch {
	case strings.Contains(header, "image/jpeg"):
		return "jpg"
	case strings.Contains(header, "image/png"):
		return "png"
	case strings.Contains(header, "application/pdf"):
		return "pdf"
	case strings.Contains(header, "image/webp"):
		return "webp"
	default:
		return ""
	}
}

var allowedExtensions = map[string]bool{
	"jpg":  true,
	"jpeg": true,
	"png":  true,
	"pdf":  true,
}

// Validate file extension
func isAllowedFileType(ext string) bool {
	ext = strings.ToLower(ext)
	return allowedExtensions[ext]
}

func streamFileChunked(c *gin.Context, filePath string, requestID any, start time.Time) {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Log.Error("archive.stream.chunked.failed",
			zap.String("request_id", requestID.(string)),
			zap.String("param", c.Param("id")),
			zap.Error(err),
			zap.Duration("duration_ms", time.Since(start)),
		)

		response.Error(c, http.StatusInternalServerError, "Failed to open file")
		return
	}
	defer file.Close()

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Header("Content-Type", "application/pdf")

	buf := make([]byte, 4096)
	for {
		n, err := file.Read(buf)
		if n > 0 {
			if _, writeErr := c.Writer.Write(buf[:n]); writeErr != nil {
				return
			}
			c.Writer.Flush()
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
	}

	logger.Log.Info("archive.stream.chunked.success",
		zap.String("request_id", requestID.(string)),
		zap.Duration("duration_ms", time.Since(start)),
	)
}

func convertImagesToPDF(imageFiles []string, outputPDF string) error {
	if len(imageFiles) == 0 {
		return nil
	}

	args := append(imageFiles, outputPDF)
	cmd := exec.Command("magick", args...)

	utils.ApplySysProcAttr(cmd)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("imagemagick error: %v | %s", err, string(out))
	}

	return nil
}

func mergePDFs(inputs []string, output string) error {
	if len(inputs) == 0 {
		return fmt.Errorf("no PDF files to merge")
	}

	args := append(inputs, output)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command(`C:\poppler\Library\bin\pdfunite.exe`, args...)
		utils.ApplySysProcAttr(cmd)
	} else {
		cmd = exec.Command("pdfunite", args...)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pdfunite error: %v | %s", err, string(out))
	}

	return nil
}
