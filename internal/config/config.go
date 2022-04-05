package config

type Config struct {
	PostgresURL    string   `env:"POSTGRES_URL" envDefault:"postgres://postgres:qwerty@localhost:5433/postgres"`
	KafkaPort      string   `env:"KAFKA_PORT" envDefault:"9092"`
	KafkaHost      string   `env:"KAFKA_HOST" envDefault:"localhost"`
	KafkaTopic     string   `env:"KAFKA_TOPIC" envDefault:"kafka-task"`
	KafkaGroupID   string   `env:"KAFKA_GROUP_ID" envDefault:"consumers"`
	LogLevel       string   `env:"LOG_LEVEL" envDefault:"info"`
	LoggerEncoding string   `env:"LOG_ENCODING" envDefault:"console"`
	Brokers        []string `env:"BROKERS" envSeparator:"," envDefault:"localhost:9092"`
	Development    bool     `env:"DEV_BOOL" envDefault:"true"`
}
