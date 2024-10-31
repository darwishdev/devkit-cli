package supabase

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/darwishdev/devkit-cli/pkg/config"
	supaapigo "github.com/darwishdev/supaapi-go"
	"github.com/supabase-community/auth-go/types"
	storage_go "github.com/supabase-community/storage-go"
)

type SupabaseClientInterface interface {
	UsersCreateUpdate(conf *config.ProjectConfig, rows [][]string) error
	StorageSeed(conf *config.ProjectConfig, filesPath string) error
	OpenConnection(conf *config.ProjectConfig) supaapigo.Supaapi
}

type SupabaseClient struct {
}

func NewSupabaseClient() SupabaseClientInterface {
	return &SupabaseClient{}
}
func (s *SupabaseClient) GetContentType(filePath string) *string {
	contentType := "application/octet-stream" // Default content type
	if strings.HasSuffix(filePath, ".jpg") || strings.HasSuffix(filePath, ".jpeg") {
		contentType = "image/jpeg"
	} else if strings.HasSuffix(filePath, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(filePath, ".webp") {
		contentType = "image/webp"
	}
	return &contentType
}
func (s *SupabaseClient) StorageSeed(conf *config.ProjectConfig, filesPath string) error {
	supaapi := s.OpenConnection(conf)
	foldersMap := map[string][]string{}
	err := filepath.Walk(filesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != filesPath { // Exclude the root directory
			bucketName := filepath.Base(path) // Use the folder name as the bucket name
			_, err = supaapi.StorageClient.CreateBucket(bucketName, storage_go.BucketOptions{
				Public: true,
			})
			if err != nil && !strings.Contains(err.Error(), "already exists") {
				return err
			}
			files, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			for _, file := range files {
				if !file.IsDir() {
					filePath := filepath.Join(path, file.Name())
					foldersMap[bucketName] = append(foldersMap[bucketName], filePath)
				}

			}
		}
		return nil
	})

	if err != nil {
		return err
	}
	for bucket, value := range foldersMap {
		for _, filePath := range value {
			fileName := filepath.Base(filePath)
			upsert := true
			fileOpts := storage_go.FileOptions{
				ContentType: s.GetContentType(filePath),
				Upsert:      &upsert,
			}
			fileReader, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer fileReader.Close()
			_, err = supaapi.StorageClient.UploadFile(bucket, fileName, fileReader, fileOpts)
			if err != nil {
				return err
			}

		}
	}
	return err
}
func (s *SupabaseClient) UsersCreateUpdate(conf *config.ProjectConfig, rows [][]string) error {
	columns := rows[0]
	api := s.OpenConnection(conf)
	for _, row := range rows[1:] { // Start from the second row (index 1)
		supabasRequest := types.AdminUpdateUserRequest{}
		for colIndex, colCell := range row {
			if columns[colIndex] == "user_email" {
				supabasRequest.Email = colCell
			}
			if columns[colIndex] == "user_password#" {
				supabasRequest.Password = colCell
			}
		}

		_, err := api.UserCreateUpdate(supabasRequest)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *SupabaseClient) OpenConnection(conf *config.ProjectConfig) supaapigo.Supaapi {
	env := supaapigo.DEV
	if conf.Environmet == "prod" {
		env = supaapigo.PROD
	}
	supaapi := supaapigo.NewSupaapi(supaapigo.SupaapiConfig{
		ProjectRef:     conf.DBProjectREF,
		Env:            env,
		Port:           conf.DBPort,
		ServiceRoleKey: conf.SupabaseServiceRoleKey,
		ApiKey:         conf.SupabaseApiKey,
	})
	return supaapi
}
