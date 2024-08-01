package options

import (
	"github.com/ahang7/go-IAM/pkg/db"
	"github.com/spf13/pflag"
	"gorm.io/gorm"
	"time"
)

type MySQLOptions struct {
	Host                  string        `json:"host" mapstructure:"host"`
	Username              string        `json:"username" mapstructure:"username"`
	Password              string        `json:"password" mapstructure:"password"`
	Database              string        `json:"database" mapstructure:"database"`
	MaxIdleConnections    int           `json:"max-idle-connections,omitempty" mapstructure:"max_idle_connections"`
	MaxOpenConnections    int           `json:"max-open-connections,omitempty" mapstructure:"max_open_connections"`
	MaxConnectionLifeTime time.Duration `json:"max-connection-life-time,omitempty" mapstructure:"max_connection_life_time"`
	LogLevel              int           `json:"log-level,omitempty" mapstructure:"log_level"`
}

// NewMySQLOptionsNil create a "" MySQLOptions
func NewMySQLOptionsNil() *MySQLOptions {
	return &MySQLOptions{
		Host:                  "127.0.0.1:3306",
		Username:              "",
		Password:              "",
		Database:              "",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: time.Duration(10) * time.Second,
		LogLevel:              1, // Silent
	}
}

func (o *MySQLOptions) Validate() []error {
	var errs []error
	return errs
}

// AddFlags adds flags to the given pflag.flagSet.
func (o *MySQLOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Host, "mysql.host", o.Host, ""+
		"MySQL service host address. If left blank, the following related mysql options will be ignored.")

	fs.StringVar(&o.Username, "mysql.username", o.Username, ""+
		"Username for access to mysql service.")

	fs.StringVar(&o.Password, "mysql.password", o.Password, ""+
		"Password for access to mysql, should be used pair with password.")

	fs.StringVar(&o.Database, "mysql.database", o.Database, ""+
		"Database name for the logicServer to use.")

	fs.IntVar(&o.MaxIdleConnections, "mysql.max-idle-connections", o.MaxOpenConnections, ""+
		"Maximum idle connections allowed to connect to mysql.")

	fs.IntVar(&o.MaxOpenConnections, "mysql.max-open-connections", o.MaxOpenConnections, ""+
		"Maximum open connections allowed to connect to mysql.")

	fs.DurationVar(&o.MaxConnectionLifeTime, "mysql.max-connection-life-time", o.MaxConnectionLifeTime, ""+
		"Maximum connection life time allowed to connect to mysql.")

	fs.IntVar(&o.LogLevel, "mysql.log-mode", o.LogLevel, ""+
		"Specify gorm log level.")
}

// NewClient new mysql client with options.
func (o *MySQLOptions) NewClient() (*gorm.DB, error) {
	opts := db.MySQLOptions{
		Host:                  o.Host,
		UserName:              o.Username,
		Password:              o.Password,
		Database:              o.Database,
		MaxIdleConnections:    o.MaxIdleConnections,
		MaxOpenConnections:    o.MaxOpenConnections,
		MaxConnectionLifeTime: o.MaxConnectionLifeTime,
		MaxConnectionIdleTime: time.Minute * 24,
		LogLevel:              o.LogLevel,
		//Logger:                nil,
	}
	return db.NewMySQLClient(&opts)
}
