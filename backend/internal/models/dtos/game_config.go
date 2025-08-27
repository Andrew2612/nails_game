package dtos

type Config struct {
	GameSettings
	DatabaseConfig
}

type GameSettings struct {
	LineSize int `json:"lineSize"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	Port     string `json:"port"`
}
