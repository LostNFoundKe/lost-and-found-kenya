package config

import "github.com/spf13/viper"

type Config struct {
	Port               int    `mapstructure:"PORT"`
	Environment        string `mapstructure:"ENVIRONMENT"`
	LogLevel           string `mapstructure:"LOG_LEVEL"`
	DatabaseURL        string `mapstructure:"DB_URL"`
	JWTSecret          string `mapstructure:"JWT_SECRET"`
	JWTExpiration      int    `mapstructure:"JWT_EXPIRATION"`
	GCSBucketName      string `mapstructure:"GCS_BUCKETNAME"`
	GCSProjectID       string `mapstructure:"GCS_PROGECT_ID"`
	GCSCredentialsFile string `mapstructure:"GCS_CREDENTIALS_FILE"`
	RedisURL           string `mapstructure:"REDIS_URL"`
}

func Load(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}

// func getEnv(key, fallback string) string {
// 	if value, exists := os.LookupEnv(key); exists {
// 		return value
// 	}
// 	return fallback
// }
