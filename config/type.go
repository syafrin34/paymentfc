package config

type Config struct {
	App      AppConfig      `yaml:"app" validate:"required"`
	Database DatabaseConfig `yaml:"database" validate:"required"`
	Redis    RedisConfig    `yaml:"redis" validate:"required"`
	Secret   SecretConfig   `yaml:"secret" validate:"required"`
	Kafka    KafkaConfig    `yaml:"kafka" validate:"required"`
	Xendit   XenditConfig   `yaml:"kafka" validate:"required"`
	Toggle   ToggleConfig   `yaml:"toggle" validate:"required"`
}
type SecretConfig struct {
	JWTSecret string `yaml:"jwt_secret" validate:"required"`
}

type AppConfig struct {
	Port string `yaml:"port" validate:"required"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Name     string `yaml:"name" validate:"required"`
	Port     string `yaml:"port" validate:"required"`
}
type RedisConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Port     string `yaml:"port" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}
type KafkaConfig struct {
	Broker      string            `yaml:"broker" validate:"required"`
	KafkaTopics map[string]string `yaml:"topics" validate:"required"`
}

type XenditConfig struct {
	XenditAPIKey  string `yaml:"xendit_api_key" validate:"required"`
	XenditWebhook string `yaml:"xendit_webhook_token" validate:"required"`
}
type ToggleConfig struct {
	DisableCreateInvoiceDirectly bool `yaml:"disable_create_invoice_directly" validate:"required"`
}
