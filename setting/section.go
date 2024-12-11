package setting

import "time"

type Config struct {
	Server ServerSetting   `mapstructure:"server"`
	MySql  MySqlSetting    `mapstructure:"mysql"`
	Logger LoggerSetting   `mapstructure:"logger"`
	Redis  RedisSetting    `mapstructure:"redis"`
	PgDb   PostgresSetting `mapstructure:"postgres"`
	Token  TokenSetting    `mapstructure:"token"`
}

type ServerSetting struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type MySqlSetting struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Dbname          string `mapstructure:"dbname"`
	MaxIdleConns    int    `mapstructure:"maxIdleConns"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns"`
	ConnMaxLifetime int    `mapstructure:"connMaxLifetime"`
}

type LoggerSetting struct {
	Log_Level     string `mapstructure:"log_level"`
	File_log_name string `mapstructure:"file_log_name"`
	Max_backups   int    `mapstructure:"max_backups"`
	Max_age       int    `mapstructure:"max_age"`
	Max_size      int    `mapstructure:"max_size"`
	Compress      bool   `mapstructure:"compress"`
}

type RedisSetting struct {
	Host_name string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Password  string `mapstructure:"password"`
}

type PostgresSetting struct {
	DbDriver string `mapstructure:"dbDriver"`
	DbSource string `mapstructure:"dbSource"`
}

type TokenSetting struct {
	SecretKey            string        `mapstructure:"secret_key"`
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
}
