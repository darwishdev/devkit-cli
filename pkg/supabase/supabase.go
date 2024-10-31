package supabase

import (
	"github.com/darwishdev/devkit-cli/pkg/config"
	supaapigo "github.com/darwishdev/supaapi-go"
	"github.com/supabase-community/auth-go/types"
)

type SupabaseClientInterface interface {
	UsersCreateUpdate(conf *config.ProjectConfig, rows [][]string) error
}

type SupabaseClient struct {
}

func NewSupabaseClient() SupabaseClientInterface {
	return &SupabaseClient{}
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
