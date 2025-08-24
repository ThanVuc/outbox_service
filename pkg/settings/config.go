package settings

type Config struct {
	PostgresAuth PostgresAuth `mapstructure:"postgres-auth" json:"postgres-auth" yaml:"postgres-auth"`
	Log          Log          `mapstructure:"log" json:"log" yaml:"log"`
	RabbitMQ     RabbitMQ     `mapstructure:"rabbitmq" json:"rabbitmq" yaml:"rabbitmq"`
}

type PostgresAuth struct {
	Host            string `mapstructure:"host" json:"host" yaml:"host"`
	Port            int    `mapstructure:"port" json:"port" yaml:"port"`
	User            string `mapstructure:"user" json:"user" yaml:"user"`
	Password        string `mapstructure:"password" json:"password" yaml:"password"`
	Database        string `mapstructure:"database" json:"database" yaml:"database"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
	MaxLifetime     int    `mapstructure:"conn_max_lifetime" json:"conn_max_lifetime" yaml:"conn_max_lifetime"`    // in seconds
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time" json:"conn_max_idle_time" yaml:"conn_max_idle_time"` // in seconds
}

type Log struct {
	Level string `mapstructure:"level" json:"level" yaml:"level"`
}

type RabbitMQ struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	User     string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}
