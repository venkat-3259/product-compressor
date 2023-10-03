package configs

type Config struct {
	Fiber    *FiberConfig    `json:"fiber" validate:"required"`
	Postgres *PostgresConfig `json:"postgres" validate:"required"`
}

type FiberConfig struct {
	Host              string `json:"host" validate:"required"`
	Port              uint   `json:"port" validate:"required"`
	ServerReadTimeout uint   `json:"serverReadTimeout" validate:"required"`
}

type PostgresConfig struct {
	Type                            string `json:"type" validate:"required"`
	Host                            string `json:"host" validate:"required"`
	Port                            uint   `json:"port" validate:"required"`
	DBName                          string `json:"dbName" validate:"required"`
	UserName                        string `json:"userName" validate:"required"`
	Password                        string `json:"password" validate:"required"`
	SSLMode                         string `json:"sslMode" validate:"required"`
	TimeZone                        string `json:"timeZone" validate:"required"`
	MaxConnections                  uint   `json:"maxConnections" validate:"required"`
	MaxIdleConnectionTimeoutMinutes uint   `json:"maxIdleConnectionTimeoutMinutes" validate:"required"`
	MaxConnectionLifeTimeMinutes    uint   `json:"maxConnectionLifeTimeMinutes" validate:"required"`
}
