package config

type Config struct {
	Mysql *MysqlConf `yaml:"mysql"`
	Jwt   *JwtConf   `yaml:"jwt"`
	Log   *LogConf   `yaml:"log"`
	Otel  *OtelConf  `yaml:"otel"`
	Redis *RedisConf `yaml:"redis"`
	Smtp  *SmtpConf  `yaml:"smtp"`
	COS   *COSConf   `yaml:"cos"`
	Kafka *KafkaConf `yaml:"kafka"`
	Grpc  *GrpcConf  `yaml:"grpc"`
}

type MysqlConf struct {
	Dsn     string `yaml:"dsn"`
	Logfile string `yaml:"logfile"`
}

type JwtConf struct {
	SecretKey string `yaml:"secretKey"`
	EncKey    string `yaml:"encKey"`
	Timeout   int    `yaml:"timeout"`
}

type LogConf struct {
	File       string `yaml:"file"`
	MaxSize    int    `yaml:"maxSize"`
	MaxBackups int    `yaml:"maxBackups"`
	MaxAge     int    `yaml:"maxAge"`
	Compress   bool   `yaml:"compress"`
}

type RedisConf struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type SmtpConf struct {
	Secret string `yaml:"secret"`
	Server string `yaml:"server"`
	Addr   string `yaml:"addr"`
}

type COSConf struct {
	SecretID    string `yaml:"secretID"`
	SecretKey   string `yaml:"secretKey"`
	URL         string `yaml:"url"`
	InternalUrl string `yaml:"internalUrl"`
}

type KafkaConf struct {
	Addrs         []string `yaml:"addrs"`
	ConsumerGroup string   `yaml:"consumerGroup"`
}

type GrpcConf struct {
	Addr string `yaml:"addr"`
}

type OtelConf struct {
	Enabled         bool   `yaml:"enabled"`
	ServiceName     string `yaml:"serviceName"`
	ServiceVersion  string `yaml:"serviceVersion"`
	TraceExporter   string `yaml:"traceExporter"`
	MetricsExporter string `yaml:"metricsExporter"`
	Endpoint        string `yaml:"endpoint"`
	Insecure        bool   `yaml:"insecure"`
	MetricsInterval int    `yaml:"metricsInterval"`
}
