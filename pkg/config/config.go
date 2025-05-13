package config

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		Port     int    `yaml:"port"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
	Smtp struct {
		Sender   string `yaml:"sender_email"`
		Password string `yaml:"sender_password"`
	} `yaml:"smtp"`
}
