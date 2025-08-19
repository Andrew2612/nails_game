package dtos

type Config struct {
	GameSettings struct {
		LineSize int `json:"lineSize"`
	}
	Database struct {
		Host     string `json:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
		Port     string `json:"port"`
	}
}
