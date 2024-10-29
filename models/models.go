package models

type ConfigMap map[string]interface{}
type ConfigModel struct {
	BasePath               string `yaml:"base_path"`
	ProtoServiceFilePath   string `yaml:"proto_service_file_path"`
	Environmet             string `yaml:"environment"`
	DBProjectREF           string `yaml:"db_project_ref"`
	SupabaseServiceRoleKey string `yaml:"supabase_service_role_key"`
	SupabaseApiKey         string `yaml:"supabase_api_key"`
	DBSource               string `yaml:"db_source"`
	BaseBuf                string `yaml:"base_buf"`
	AppName                string `yaml:"app_name"`
	GithubToken            string `yaml:"github_token"`
	ApiFilePath            string `yaml:"api_file_path"`
	ExePath                string
	ApiServiceName         string `yaml:"api_service_name"`
	QueryPath              string `yaml:"query_path"`
	StorePath              string `yaml:"store_path"`
	ProtoGenPath           string `yaml:"proto_gen_path"`
	ProtoPath              string `yaml:"proto_path"`
	ApiVersion             string `yaml:"api_version"`
}
