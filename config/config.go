package config

type Config struct {
	MongoDB *MongoConfig
	GraphDB *GraphConfig
	SECRET string
	StripeKey string
	Redis * RedisConfig
}

type GraphConfig struct {
	Dialect  string
	Username string
	Password string
	Charset  string
}

type RedisConfig struct {
	Host string
	Port string
	Prefix string
}

type MongoConfig struct {
	Database  string
	Server string
}

func GetConfig() *Config {
	return &Config{
		MongoDB: &MongoConfig{
			Database:  "invest",
			Server: "mongodb://gcic:dreamer6@ds163330.mlab.com:63330/invest",
		},
		Redis: &RedisConfig{
			Host: "http://localhost",
			Port: "6379",
			Prefix: "socket.io",
		},
		SECRET: "!n^e$t%4dc$2id*32+",
		StripeKey: "sk_test_BQokikJOvBiI2HlWgH4olfQ2",
	}
}