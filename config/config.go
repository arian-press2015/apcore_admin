package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    BackendURL string
    TokenFile  string
}

func NewConfig() *Config {
    viper.SetDefault("backend_url", "http://localhost:8080")
    viper.SetDefault("token_file", ".auth_token")
    
    viper.AutomaticEnv()
    
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")

    if err := viper.ReadInConfig(); err == nil {
        // Config file was found and successfully parsed
    } else {
        // Handle error if needed (optional)
    }
    
    return &Config{
        BackendURL: viper.GetString("backend_url"),
        TokenFile:  viper.GetString("token_file"),
    }
}
