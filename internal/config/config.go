package config

type Config struct {
	PostgresURL    string   `env:"POSTGRES_URL" envDefault:"postgres://postgres:qwerty@postgresdb:5432/postgres"`
	KafkaTopic     string   `env:"KAFKA_TOPIC" envDefault:"messages-topic"`
	KafkaGroupID   string   `env:"KAFKA_GROUP_ID" envDefault:"consumers-group"`
	LogLevel       string   `env:"LOG_LEVEL" envDefault:"info"`
	LoggerEncoding string   `env:"LOG_ENCODING" envDefault:"console"`
	NumConsumer    int      `env:"NUM_CONSUMERS" envDefault:"3"`
	Brokers        []string `env:"BROKERS" envSeparator:"," envDefault:"kafka1:19091,kafka2:19092,kafka3:19093"`
}
