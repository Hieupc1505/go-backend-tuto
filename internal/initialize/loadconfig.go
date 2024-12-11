package initialize

import (
	"fmt"

	"github.com/spf13/viper"
	"hieupc05.github/backend-server/global"
)

func LoadConfig(path string) {
	viper := viper.New()
	viper.AddConfigPath(path)
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")

	//read configuaration
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Failed to read configuration %w \n", err))
	}
	//read server configuration
	fmt.Println("Server port::", viper.GetInt("server.port"))

	if err := viper.Unmarshal(&global.Config); err != nil {
		panic(fmt.Errorf("Failed to read configuration %w \n", err))
	}

}
