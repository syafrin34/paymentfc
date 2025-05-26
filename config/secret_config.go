package config

import (
	"encoding/json"
	"fmt"
	"log"

	vaultAPI "github.com/hashicorp/vault/api"
)

func LoadSecretConfig(cfg Config) Config {
	secretConfig := vaultAPI.DefaultConfig()
	secretConfig.Address = cfg.Vault.Host

	tempDebug15, _ := json.Marshal(cfg.Vault.Host)
	fmt.Printf("\n===== DEBUG secret_config.go - Line: 15 ==== \n\n%s\n\n=============\n\n\n", string(tempDebug15))

	vaultClient, err := vaultAPI.NewClient(secretConfig)
	if err != nil {
		log.Fatalf("failed init vault client: %v", err)
	}

	vaultClient.SetToken(cfg.Vault.Token)

	// red secret
	secret, err := vaultClient.Logical().Read(cfg.Vault.Path)
	if err != nil {
		log.Fatalf("failed read secret data: %v", err)
	}
	if secret == nil || secret.Data == nil {
		log.Fatalf("secret data not found")
	}

	// vault data --> map[string]interface{}
	vaultData, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		log.Fatalf("invalid secret format !")
	}
	vaultDataJSON, err := json.Marshal(vaultData)
	if err != nil {
		log.Fatalf("failed marshal vault data")
	}
	var secretVaultConfig SecretVaultConfig
	err = json.Unmarshal(vaultDataJSON, &secretVaultConfig)
	if err != nil {
		log.Fatalf("failed unmarshal vaul data to secret config: %v", err)
	}

	//removee me - debug purpose only
	tempDebug36, _ := json.Marshal(secretVaultConfig)
	fmt.Printf("\n===== DEBUG secret_config.go - Line: 36 ==== \n\n%s\n\n=============\n\n\n", string(tempDebug36))

	// construct config with seccret config
	cfg.Database.Password = secretVaultConfig.DatabaseSecret.Password
	cfg.Redis.Password = secretVaultConfig.RedisSecret.Password
	cfg.JWT.JWTSecret = secretVaultConfig.JWTSecret
	cfg.Xendit.XenditAPIKey = secretVaultConfig.XenditSecret.SecretAPIKey
	cfg.Xendit.XenditWebhookToken = secretVaultConfig.XenditSecret.WebhookToken
	return cfg
}
