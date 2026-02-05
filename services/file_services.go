package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"upload/models"
	"upload/repositories"

	"github.com/google/uuid"
)

const (
	MaxPhotoSize       = 1 * 1024 * 1024
	MaxCertificateSize = 2 * 1024 * 1024
	MaxGeneralSize     = 10 * 1024 * 1024
)

var (
	AllowedPhotoTypes = map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
	}

	AllowedCertificateTypes = map[string]bool{
		"application/pdf": true,
	}

	AllowedGeneralTypes = map[string]bool{
		"image/jpeg":      true,
		"image/jpg":       true,
		"image/png":       true,
		"application/pdf": true,
	}
)

type FileService struct {
	repo       repositories.FileRepository
	uploadPath string
	baseURL    string
}

func NewFileService(repo repositories.FileRepository, uploadPath, baseURL string) *FileService {
	return &FileService{
		repo:       repo,
		uploadPath: uploadPath,
		baseURL:    baseURL,
	}
}

func (s *FileService) UploadPhoto(file *multipart.FileHeader, userID string) (*models.FileResponse, error) {
	return s.uploadWithValidation(file, userID, AllowedPhotoTypes, MaxPhotoSize, "photo")
}

func (s *FileService) UploadCertificate(file *multipart.FileHeader, userID string) (*models.FileResponse, error) {
	return s.uploadWithValidation(file, userID, AllowedCertificateTypes, MaxCertificateSize, "certificate")
}

func (s *FileService) Upload(file *multipart.FileHeader, userID string) (*models.FileResponse, error) {
	return s.uploadWithValidation(file, userID, AllowedGeneralTypes, MaxGeneralSize, "file")
}

func (s *FileService) uploadWithValidation(
	file *multipart.FileHeader,
	userID string,
	allowedTypes map[string]bool,
	maxSize int64,
	fileType string,
) (*models.FileResponse, error) {
	if file.Size > maxSize {
		return nil, fmt.Errorf("%s size exceeds maximum allowed size of %d bytes", fileType, maxSize)
	}

	mimeType := file.Header.Get("Content-Type")
	if !allowedTypes[mimeType] {
		allowedList := make([]string, 0, len(allowedTypes))
		for t := range allowedTypes {
			allowedList = append(allowedList, t)
		}
		return nil, fmt.Errorf("invalid %s type: %s. Allowed types: %s", fileType, mimeType, strings.Join(allowedList, ", "))
	}

	ext := filepath.Ext(file.Filename)
	storedName := uuid.New().String() + ext

	if err := os.MkdirAll(s.uploadPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	filePath := filepath.Join(s.uploadPath, storedName)
	if err := s.saveFile(file, filePath); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	fileModel := &models.File{
		UserID:       userID,
		OriginalName: file.Filename,
		StoredName:   storedName,
		MimeType:     mimeType,
		Size:         file.Size,
		Path:         filePath,
	}

	if err := s.repo.Create(fileModel); err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}

	response := fileModel.ToResponse(s.baseURL)
	return &response, nil
}

func (s *FileService) saveFile(file *multipart.FileHeader, destPath string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func (s *FileService) GetAll() ([]models.FileResponse, error) {
	files, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]models.FileResponse, len(files))
	for i, file := range files {
		responses[i] = file.ToResponse(s.baseURL)
	}

	return responses, nil
}

func (s *FileService) GetByID(id string) (*models.FileResponse, error) {
	file, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("file not found")
	}

	response := file.ToResponse(s.baseURL)
	return &response, nil
}

func (s *FileService) Delete(id string) error {
	file, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("file not found")
	}

	if err := os.Remove(file.Path); err != nil {
		fmt.Printf("Warning: failed to delete file from disk: %v\n", err)
	}

	return s.repo.Delete(id)
}
