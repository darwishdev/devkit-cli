package seed

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/supabase-community/auth-go/types"
)

// This command creates a new domain within your Go backend application
// by generating the necessary files and directory structure.
func (c *SeedCmd) SeedSuperUser(flags *pflag.FlagSet) {
	email, _ := flags.GetString("email")
	password, _ := flags.GetString("password")
	phone, _ := flags.GetString("phone")
	name, _ := flags.GetString("name")
	conf, err := c.config.GetProjectConfig()
	if err != nil {
		log.Err(err).Msg("can't read the project config")
		os.Exit(1)
	}

	log.Debug().Interface("config is", conf.SupabaseAPIPort).Msg("debug conf")
	db, err := c.dbUtils.Open(conf.DBSource)
	defer db.Close()

	if err != nil {
		log.Err(err).Str("source", conf.DBSource).Msg("can't connect to the database")
		os.Exit(1)
	}
	var created_user_id int

	rollbackUserQuery := "DELETE FROM accounts_schema.user where user_id = $1"
	rollbackAuthUserQuery := "DELETE FROM auth.users where email = $1"
	userInsertQuery := "INSERT INTO accounts_schema.user (user_name, user_phone , user_email , user_password , user_type_id) VALUES ($1,$2,$3,$4,1) RETURNING user_id"
	userRoleInsertQuery := "INSERT INTO accounts_schema.user_role (user_id, role_id) VALUES ($1,1)"
	supabaseRequest := types.AdminUpdateUserRequest{
		Email:    email,
		Password: password,
	}
	err = c.supaClient.UserCreateUpdate(conf, supabaseRequest)
	if err != nil {
		log.Err(err).Msg("can't insert on supabase auth.users table")
		os.Exit(1)
	}
	log.Info().Msg("super admin user created on supabase auth.user table")
	err = db.QueryRow(userInsertQuery, name, phone, email, password).Scan(&created_user_id)
	if err != nil {
		log.Err(err).Msg("error inserting the user on accounts_schema.user")
		db.Exec(rollbackAuthUserQuery, email)
		os.Exit(1)
	}
	log.Info().Msg("super admin user created on accounts_schema.user")
	_, err = db.Exec(userRoleInsertQuery, created_user_id)
	if err != nil {
		log.Err(err).Msg("error executing the insert statement")
		db.Exec(rollbackAuthUserQuery, email)
		db.Exec(rollbackUserQuery, created_user_id)
		os.Exit(1)
	}
	log.Info().Msg("super admin user attached to super admin role on table accounts_schmema.user_role")
}
