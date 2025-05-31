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
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"smtp"`
	Kafka struct {
		Topic string `yaml:"topic"`
	} `yaml:"kafka"`
	Ports struct {
		HTTP  int      `yaml:"http"`
		GRPC  int      `yaml:"grpc"`
		WS    int      `yaml:"ws"`
		Kafka []string `yaml:"kafka"`
	} `yaml:"ports"`
}
