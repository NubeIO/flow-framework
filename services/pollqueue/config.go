package pollqueue

type Config struct {
	EnablePolling bool   `yaml:"enable_polling"`
	LogLevel      string `yaml:"log_level"`
}
