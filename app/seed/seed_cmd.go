package seed

import (
	"github.com/darwishdev/devkit-cli/pkg/config"
	"github.com/darwishdev/devkit-cli/pkg/db"
	"github.com/darwishdev/devkit-cli/pkg/fileutils"
	"github.com/darwishdev/devkit-cli/pkg/supabase"
	"github.com/darwishdev/sqlseeder"
	"github.com/spf13/pflag"
)

type SeedCmdInterface interface {
	NewSeed(args []string, flags *pflag.FlagSet)
	StorageSeed(flags *pflag.FlagSet)
}
type SeedCmd struct {
	config     config.ConfigInterface
	sqlSeeder  sqlseeder.SeederInterface
	dbUtils    db.DbInterface
	supaClient supabase.SupabaseClientInterface
	fileUtils  fileutils.FileUtilsInterface
}

func NewSeedCmd(config config.ConfigInterface, fileUtils fileutils.FileUtilsInterface, sqlseeder sqlseeder.SeederInterface, dbUtils db.DbInterface, supaClient supabase.SupabaseClientInterface) SeedCmdInterface {
	return &SeedCmd{
		config:     config,
		dbUtils:    dbUtils,
		sqlSeeder:  sqlseeder,
		supaClient: supaClient,
		fileUtils:  fileUtils,
	}
}
