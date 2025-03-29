package supabase

import (
	"bytes"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/darwishdev/devkit-cli/pkg/config"
	supaapigo "github.com/darwishdev/supaapi-go"
	"github.com/kolesa-team/go-webp/encoder"
	webpencoder "github.com/kolesa-team/go-webp/webp"
	"github.com/rs/zerolog/log"
	"github.com/supabase-community/auth-go/types"
	storage_go "github.com/supabase-community/storage-go"
)

type SupabaseClientInterface interface {
	UserCreateUpdate(conf *config.ProjectConfig, req types.AdminUpdateUserRequest) error
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
	supporedTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".svg":  "image/svg",
		".webp": "image/webp",
		".pdf":  "application/pdf",
	}
	ext := filepath.Ext(filePath)
	mappedType, ok := supporedTypes[ext]
	if !ok {
		return &contentType
	}
	return &mappedType
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
			contentType := *s.GetContentType(filePath)
			upsert := true
			fileOpts := storage_go.FileOptions{
				ContentType: s.GetContentType(filePath),
				Upsert:      &upsert,
			}
			file, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}
			if !strings.HasPrefix(contentType, "image/") {
				fileReader := bytes.NewReader(file)
				_, err = supaapi.StorageClient.UploadFile(bucket, fileName, fileReader, fileOpts)
				if err != nil {
					return err
				}
				continue
			}
			if strings.Contains(contentType, "jpeg") {
				fileOpened, err := os.Open(filePath)
				if err != nil {
					return err
				}

				img, err := jpeg.Decode(fileOpened)
				if err != nil {

					log.Debug().Msg("jpeeeg")
					return err
				}
				options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
				if err != nil {
					log.Debug().Msg("NewLossyEncoderOptions new")
					return err
				}

				var bufferOutput bytes.Buffer
				if err := webpencoder.Encode(&bufferOutput, img, options); err != nil {
					log.Debug().Msg("ebpencoder new")
					return err
				}
				webpData := bufferOutput.Bytes()
				fileReader := bytes.NewReader(webpData)
				newFileName := strings.Replace(fileName, "jpeg", "webp", 1)
				newFileName = strings.Replace(newFileName, "jpg", "webp", 1)
				newContenType := "image/webp"
				fileOpts := storage_go.FileOptions{
					ContentType: &newContenType,
					Upsert:      &upsert,
				}
				_, err = supaapi.StorageClient.UploadFile(bucket, newFileName, fileReader, fileOpts)
				if err != nil {
					return err
				}
				continue

			}
			log.Debug().Interface("filename ", filePath).Msg("process")
			compressedFile, err := s.CompressImage(file, 70, contentType)
			if err != nil {
				return err
			}

			fileReader := bytes.NewReader(compressedFile)
			_, err = supaapi.StorageClient.UploadFile(bucket, fileName, fileReader, fileOpts)
			if err != nil {
				return err
			}

		}
	}
	return err
}

func (s *SupabaseClient) CompressImage(buffer []byte, quality float32, contentType string) ([]byte, error) {
	// Decode the input image

	if !strings.Contains(contentType, "jpg") {
		return buffer, nil
		// img, err := jpeg.Decode(bytes.NewReader(buffer))
		// if err != nil {
		//
		// 	log.Debug().Msg("jpeeeg")
		// 	return nil, err
		// }
		// options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
		// if err != nil {
		// 	log.Debug().Msg("NewLossyEncoderOptions new")
		// 	return nil, err
		// }
		// output, err := os.Create("tmp.webp")
		// if err != nil {
		// 	log.Debug().Msg("Create new")
		// 	return nil, err
		// }
		// if err := webpencoder.Encode(output, img, options); err != nil {
		// 	log.Debug().Msg("ebpencoder new")
		// 	return nil, err
		// }

		// log.Debug().Interface("output is ", output).Msg("file new")

	}
	img, format, err := image.Decode(bytes.NewReader(buffer))
	log.Printf("Image format: %s\n", format)
	if err != nil {
		log.Debug().Interface("error happens there", err.Error()).Msg("eror")
		// return nil, err
	}
	log.Printf("Image format: %s\n", format)
	// Compress the image to WebP format
	var compressed bytes.Buffer
	options := &webp.Options{Quality: quality} // Adjust quality (0-100)
	if err := webp.Encode(&compressed, img, options); err != nil {
		log.Debug().Interface("error happens here", err).Msg("eror")
		return nil, err
	}

	return compressed.Bytes(), nil
}
func (s *SupabaseClient) UserCreateUpdate(conf *config.ProjectConfig, req types.AdminUpdateUserRequest) error {
	api := s.OpenConnection(conf)
	_, err := api.UserCreateUpdate(req)
	if err != nil {
		return err
	}
	return nil
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
	if conf.State == "prod" || conf.State == "PROD" {
		env = supaapigo.PROD
	}
	supaapi := supaapigo.NewSupaapi(supaapigo.SupaapiConfig{
		ProjectRef:     conf.DBProjectREF,
		Env:            env,
		Port:           conf.SupabaseAPIPort,
		ServiceRoleKey: conf.SupabaseServiceRoleKey,
		ApiKey:         conf.SupabaseApiKey,
	})
	return supaapi
}
