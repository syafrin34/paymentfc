package config

type Config struct {
	App      AppConfig      `yaml:"app" validate:"required"`
	Database DatabaseConfig `yaml:"database" validate:"required"`
	Redis    RedisConfig    `yaml:"redis" validate:"required"`
	JWT      JWTConfig      `yaml:"jwt" validate:"required"`
	Kafka    KafkaConfig    `yaml:"kafka" validate:"required"`
	Xendit   XenditConfig   `yaml:"xendit" validate:"required"`
	Toggle   ToggleConfig   `yaml:"toggle" validate:"required"`
	Vault    VaultConfig    `yaml:"vault" validate:"required"`
}
type JWTConfig struct {
	JWTSecret string `yaml:"jwt_secret" validate:"required"`
}

type AppConfig struct {
	Port      string `yaml:"port" validate:"required"`
	JWTSecret string `yaml:"jwtsecret" validate:"required"`
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
	Broker string            `yaml:"broker" validate:"required"`
	Topics map[string]string `yaml:"topics" validate:"required"`
}

type XenditConfig struct {
	XenditAPIKey       string `yaml:"xendit_api_key" validate:"required"`
	XenditWebhookToken string `yaml:"xendit_webhook_token" validate:"required"`
}
type ToggleConfig struct {
	DisableCreateInvoiceDirectly bool `yaml:"disable_create_invoice_directly" validate:"required"`
}
type VaultConfig struct {
	Host  string `yaml:"host" validate:"required"`
	Token string `yaml:"token" validate:"required"`
	Path  string `yaml:"path" validate:"required"`
}

// secret config vault
type SecretVaultConfig struct {
	DatabaseSecret DatabaseSecretConfig `json:"database"`
	RedisSecret    RedisSecretConfig    `json:"redis"`
	JWTSecret      string               `json:"jwt_secret"`
	XenditSecret   XenditSecretConfig   `json:"xendit"`
}

type DatabaseSecretConfig struct {
	Password string `json:"password"`
}

type RedisSecretConfig struct {
	Password string `json:"password"`
}

type XenditSecretConfig struct {
	SecretAPIKey string `json:"secret_api_key"`
	WebhookToken string `json:"webhook_token"`
}
