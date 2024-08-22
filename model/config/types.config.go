package config

type Config struct {
	Server     ServerConfig   `json:"server"`
	Database   DatabaseConfig `json:"database"`
	Secrets    SecretConfig   `json:"secrets"`
	Stage      string         `json:"stage"`
	SMTP       SMTPConfig     `json:"smtp"`
	GRPCConfig GRPCConfig     `json:"grpc"`
}

type PostgreConfig struct {
	Host   string `json:"host"`
	DBName string `json:"db_name"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	URL      string `json:"url"`
}

type ServerConfig struct {
	GRPC PortConfig `json:"grpc"`
	Rest PortConfig `json:"rest"`
}

type PortConfig struct {
	Port string `json:"port"`
	Name string `json:"name"`
}

type DatabaseConfig struct {
	PostgreDB PostgreConfig `json:"postgre"`
	RedisDB   RedisConfig   `json:"redis"`
}

type SecretConfig struct {
	WaitingRoomSecret string `json:"waiting_room_secret"`
	MainRoomSecret    string `json:"main_room_secret"`
}

type SMTPConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type GRPCConfig struct {
	DashboardQueue string `json:"dashboard_queue"`
}
