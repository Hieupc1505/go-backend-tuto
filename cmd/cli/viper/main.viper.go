package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	Database []struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
	} `mapstructure:"database"`
}

func main() {
	viper := viper.New()
	viper.AddConfigPath("./configs")
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")

	//read configuaration
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Failed to read configuration %w \n", err))
	}
	//read server configuration
	fmt.Println("Server port::", viper.GetInt("server.port"))
	fmt.Println("Server Port::", viper.GetString("security.jwt.key"))

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("Failed to read configuration %w \n", err))
	}

	fmt.Println("Server port::", config.Server.Port)
	for _, db := range config.Database {
		fmt.Println("Database name:", db)

	}

}
