package handler

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mugnialby/arsip-backend/internal/model"
	archiveRequest "github.com/mugnialby/arsip-backend/internal/model/dto/request/archive"
	archiveRoleAccessRequest "github.com/mugnialby/arsip-backend/internal/model/dto/request/archiveRoleAccess"
	"github.com/mugnialby/arsip-backend/internal/service"
	"github.com/mugnialby/arsip-backend/pkg/response"
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
	archives, err := h.archiveService.GetAllArchives()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, archives)
}

func (h *ArchiveHandler) GetAllArchivesByData(c *gin.Context) {
	var getArchiveByDataRequest archiveRequest.GetArchiveByDataRequest
	if err := c.ShouldBindJSON(&getArchiveByDataRequest); err != nil {
		response.Error(c, http.StatusBadRequest, "Form data is not valid")
		return
	}

	archives, err := h.archiveService.GetAllArchivesByData(getArchiveByDataRequest)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, archives)
}

func (h *ArchiveHandler) GetArchiveByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	archive, err := h.archiveService.GetArchiveByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Archive not found"})
		return
	}

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

	c.JSON(http.StatusOK, archive)
}

func (h *ArchiveHandler) CreateArchive(c *gin.Context) {
	var newArchiveRequest archiveRequest.NewArchiveRequest
	if err := c.ShouldBindJSON(&newArchiveRequest); err != nil {
		response.Error(c, http.StatusBadRequest, "Form data is not valid")
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
		response.Error(c, http.StatusInternalServerError, "Failed to create archive hdr")
		return
	}

	for _, roleAccess := range newArchiveRequest.RoleAccess {
		newArchiveRoleAccess := model.ArchiveRoleAccess{
			ArchiveHdrID: newArchive.ID,
			RoleID:       roleAccess.RoleID,
			DepartmentID: roleAccess.DepartmentID,
			Status:       "Y",
			CreatedBy:    newArchiveRequest.SubmittedBy,
		}

		if err := h.archiveRoleAccessService.CreateArchiveRoleAccess(&newArchiveRoleAccess); err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to save archive RoleAccess")
			return
		}
	}

	for _, archiveAttachment := range newArchiveRequest.ListArchiveAttachments {
		if archiveAttachment.IsNew {
			base64Data := archiveAttachment.FileBase64

			// remove base64 header
			if strings.Contains(base64Data, ",") {
				parts := strings.SplitN(base64Data, ",", 2)
				base64Data = parts[1]
			}

			fileExt := DetectBase64Extension(archiveAttachment.FileBase64)
			if fileExt == "" {
				response.Error(c, http.StatusBadRequest, "Unsupported file type")
				return
			}

			// validate allowed file type
			if !isAllowedFileType(fileExt) {
				response.Error(c, http.StatusBadRequest, "File type not allowed")
				return
			}

			decodedBytes, err := base64.StdEncoding.DecodeString(base64Data)
			if err != nil {
				response.Error(c, http.StatusBadRequest, "Invalid base64 file data")
				return
			}

			uploadDir := createUploadDirectory(newArchive.ID)
			fileName := fmt.Sprintf("%d_%d.%s", newArchive.ID, time.Now().UnixNano(), fileExt)
			fileLocation := filepath.Join(uploadDir, fileName)

			if err := os.WriteFile(fileLocation, decodedBytes, 0644); err != nil {
				response.Error(c, http.StatusInternalServerError, "Failed to save file")
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
				response.Error(c, http.StatusInternalServerError, "Failed to save archive attachment")
				return
			}
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Archive created successfully"})
}

func (h *ArchiveHandler) UpdateArchiveById(c *gin.Context) {
	var updateArchiveRequest archiveRequest.UpdateArchiveRequest
	if err := c.ShouldBind(&updateArchiveRequest); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	archive, err := h.archiveService.GetArchiveByID(updateArchiveRequest.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "data not found"})
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
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
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
				response.Error(c, http.StatusInternalServerError, "Failed to save archive RoleAccess")
				return
			}
		}

		if roleAccess.IsDelete {
			deleteArchiveRoleAccess := archiveRoleAccessRequest.DeleteArchiveRoleAccessRequest{
				ID:          roleAccess.ID,
				SubmittedBy: updateArchiveRequest.SubmittedBy,
			}

			if err := h.archiveRoleAccessService.DeleteArchiveRoleAccess(&deleteArchiveRoleAccess); err != nil {
				// tambah logger di sini
				response.Error(c, http.StatusInternalServerError, "API Fail")
				return
			}

			c.Status(http.StatusOK)
		}
	}

	for _, archiveAttachment := range updateArchiveRequest.ListArchiveAttachments {
		if archiveAttachment.IsNew {
			base64Data := archiveAttachment.FileBase64

			// remove base64 header
			if strings.Contains(base64Data, ",") {
				parts := strings.SplitN(base64Data, ",", 2)
				base64Data = parts[1]
			}

			fileExt := DetectBase64Extension(archiveAttachment.FileBase64)
			if fileExt == "" {
				response.Error(c, http.StatusBadRequest, "Unsupported file type")
				return
			}

			// validate allowed file type
			if !isAllowedFileType(fileExt) {
				response.Error(c, http.StatusBadRequest, "File type not allowed")
				return
			}

			decodedBytes, err := base64.StdEncoding.DecodeString(base64Data)
			if err != nil {
				response.Error(c, http.StatusBadRequest, "Invalid base64 file data")
				return
			}

			uploadDir := createUploadDirectory(updateArchiveRequest.ID)
			fileName := fmt.Sprintf("%d_%d.%s", updateArchiveRequest.ID, time.Now().UnixNano(), fileExt)
			fileLocation := filepath.Join(uploadDir, fileName)

			if err := os.WriteFile(fileLocation, decodedBytes, 0644); err != nil {
				response.Error(c, http.StatusInternalServerError, "Failed to save file")
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
				response.Error(c, http.StatusInternalServerError, "Failed to save archive attachment")
				return
			}

			projectRoot := getProjectRoot()
			cacheFolder := "storage/cache"
			cacheFilePath := filepath.Join(projectRoot, cacheFolder, fmt.Sprintf("archive_%d.pdf", updateArchiveRequest.ID))

			if err := os.Remove(cacheFilePath); err == nil {
				log.Println("Cache deleted:", cacheFilePath)
			} else if os.IsNotExist(err) {
				log.Println("No cache found:", cacheFilePath)
			} else {
				log.Println("Failed to delete cache:", err)
			}
		}

		// ❌ Handle deleted attachments
		if archiveAttachment.IsDelete {
			archiveAttachment, err := h.archiveAttachmentService.GetArchiveAttachmentByID(archiveAttachment.ID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Archive attachment not found"})
				return
			}

			archiveAttachment.Status = "N"
			if err := h.archiveAttachmentService.UpdateArchiveAttachment(archiveAttachment); err != nil {
				// tambah logger di sini
				response.Error(c, http.StatusInternalServerError, "API Fail")
				return
			}

			c.JSON(http.StatusOK, archiveAttachment)

			projectRoot := getProjectRoot()
			cacheFolder := "storage/cache"
			cacheFilePath := filepath.Join(projectRoot, cacheFolder, fmt.Sprintf("archive_%d.pdf", updateArchiveRequest.ID))

			if err := os.Remove(cacheFilePath); err == nil {
				log.Println("Cache deleted:", cacheFilePath)
			} else if os.IsNotExist(err) {
				log.Println("No cache found:", cacheFilePath)
			} else {
				log.Println("Failed to delete cache:", err)
			}
		}
	}

	c.JSON(http.StatusOK, archive)
}

func (h *ArchiveHandler) DeleteArchiveById(c *gin.Context) {
	var deleteArchiveRequest archiveRequest.DeleteArchiveRequest
	if err := c.ShouldBindJSON(&deleteArchiveRequest); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusBadRequest, "JSON Request is not valid")
		return
	}

	if err := h.archiveService.DeleteArchive(&deleteArchiveRequest); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	if err := h.archiveAttachmentService.DeleteArchiveAttachmentByArchiveID(deleteArchiveRequest.ID, deleteArchiveRequest.SubmittedBy); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	if err := h.archiveRoleAccessService.DeleteArchiveRoleAccessByArchiveID(deleteArchiveRequest.ID, deleteArchiveRequest.SubmittedBy); err != nil {
		// tambah logger di sini
		response.Error(c, http.StatusInternalServerError, "API Fail")
		return
	}

	c.Status(http.StatusOK)
}

func (h *ArchiveHandler) FindArchiveByQuery(c *gin.Context) {
	query := c.Param("query")
	archives, err := h.archiveService.FindArchiveByQuery(query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Archives not found"})
		return
	}

	c.JSON(http.StatusOK, archives)
}

func (h *ArchiveHandler) FindArchiveByAdvanceQuery(c *gin.Context) {
	var advancedSearchRequest archiveRequest.AdvancedSearchRequest
	if err := c.ShouldBindJSON(&advancedSearchRequest); err != nil {
		response.Error(c, http.StatusBadRequest, "JSON is not valid")
		return
	}

	archives, err := h.archiveService.FindArchiveByAdvanceQuery(advancedSearchRequest)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Archives not found"})
		return
	}

	c.JSON(http.StatusOK, archives)
}

func (h *ArchiveHandler) StreamMergedPDF(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	archive, err := h.archiveService.GetArchiveByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Archive not found"})
		return
	}

	if len(archive.ArchiveAttachments) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No attachments"})
		return
	}

	projectRoot := getProjectRoot()
	cacheDir := filepath.Join(projectRoot, "storage", "cache")
	os.MkdirAll(cacheDir, 0755)

	finalPDF := filepath.Join(cacheDir, fmt.Sprintf("archive_%d.pdf", id))

	// ✅ Serve cache
	if stat, err := os.Stat(finalPDF); err == nil {
		if time.Since(stat.ModTime()) < CacheTTL {
			streamFileChunked(c, finalPDF)
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		pdfFiles = append([]string{tempImagePDF}, pdfFiles...)
	}

	if len(pdfFiles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid files"})
		return
	}

	if err := mergePDFs(pdfFiles, finalPDF); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Cleanup temp
	if tempImagePDF != "" {
		_ = os.Remove(tempImagePDF)
	}

	streamFileChunked(c, finalPDF)
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

func streamFileChunked(c *gin.Context, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
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
}

func getProjectRoot() string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal("Cannot determine executable path:", err)
	}

	dir := filepath.Dir(exe)

	// Walk up until we find the project folder "arsip-backend"
	for {
		if strings.HasSuffix(dir, "arsip-backend") {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			log.Fatal("Unable to find project root (arsip-backend)")
		}
		dir = parent
	}
}

func createUploadDirectory(archiveId uint) string {
	projectRoot := getProjectRoot()

	convertedArchiveId := strconv.Itoa(int(archiveId))
	uploadDir := filepath.Join(projectRoot, "storage", "uploads", "archives", convertedArchiveId)

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			log.Fatal("Failed to create upload dir:", err)
		}
		log.Println("Created upload directory:", uploadDir)
	} else {
		log.Println("Upload directory exists:", uploadDir)
	}

	return uploadDir
}

func getImageMagickCmd() string {
	switch runtime.GOOS {
	case "windows":
		return "magick" // ImageMagick 7+
	default:
		return "magick"
	}
}

func convertImagesToPDF(imageFiles []string, outputPDF string) error {
	if len(imageFiles) == 0 {
		return nil
	}

	args := append(imageFiles, outputPDF)
	cmd := exec.Command(getImageMagickCmd(), args...)

	// Windows safety
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

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
		cmd = exec.Command(getPDFUnitePath(), args...)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	} else {
		cmd = exec.Command("pdfunite", args...)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pdfunite error: %v | %s", err, string(out))
	}

	return nil
}

func getPDFUnitePath() string {
	// OPTION A: if added to PATH
	// return "pdfunite"

	// OPTION B (recommended): absolute path
	return `C:\poppler\Library\bin\pdfunite.exe`
}
