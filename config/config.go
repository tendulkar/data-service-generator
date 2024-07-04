// config/config.go
package config

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Model    ModelConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ModelConfig struct {
	BaseDir              string
	TypesPath            string
	ValidationsPath      string
	PostgresTypeMapsPath string
	AttributesPath       string
}
