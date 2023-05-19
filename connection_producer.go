package ysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/mitchellh/mapstructure"

	_ "github.com/yugabyte/pgx/v4/stdlib"
)

// YugabyteConnectionProducer implements ConnectionProducer and provides a generic producer for most yugabyte databases
type YugabyteConnectionProducer struct {
	ConnectionURL            string      `json:"connection_url" mapstructure:"connection_url" structs:"connection_url"`
	MaxOpenConnections       int         `json:"max_open_connections" mapstructure:"max_open_connections" structs:"max_open_connections"`
	MaxIdleConnections       int         `json:"max_idle_connections" mapstructure:"max_idle_connections" structs:"max_idle_connections"`
	MaxConnectionLifetimeRaw interface{} `json:"max_connection_lifetime" mapstructure:"max_connection_lifetime" structs:"max_connection_lifetime"`
	Host                     string      `json:"host" mapstructure:"host" structs:"host"`
	Username                 string      `json:"username" mapstructure:"username" structs:"username"`
	Password                 string      `json:"password" mapstructure:"password" structs:"password"`
	Port                     int         `json:"port" mapstructure:"port" structs:"port"`
	DbName                   string      `json:"db" mapstructure:"db" structs:"db"`
	LoadBalance              bool        `json:"load_balance" mapstructure:"load_balance" structs:"load_balance"`
	YbServersRefreshInterval int         `json:"yb_servers_refresh_interval" mapstructure:"yb_servers_refresh_interval" structs:"yb_servers_refresh_interval"`
	TopologyKeys             string      `json:"topology_keys" mapstructure:"topology_keys" structs:"topology_keys"`

	Type                  string
	RawConfig             map[string]interface{}
	maxConnectionLifetime time.Duration
	Initialized           bool
	db                    *sql.DB
	sync.Mutex
}

var ErrNotInitialized = errors.New("connection has not been initialized")

func (c *YugabyteConnectionProducer) Initialize(ctx context.Context, conf map[string]interface{}, verifyConnection bool) error {
	_, err := c.Init(ctx, conf, verifyConnection)
	return err
}

func (c *YugabyteConnectionProducer) Init(ctx context.Context, conf map[string]interface{}, verifyConnection bool) (map[string]interface{}, error) {
	c.Lock()
	defer c.Unlock()

	c.RawConfig = conf

	decoderConfig := &mapstructure.DecoderConfig{
		Result:           c,
		WeaklyTypedInput: true,
		TagName:          "json",
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(conf)
	if err != nil {
		return nil, err
	}

	switch {
	case len(c.ConnectionURL) != 0:
		break //As the connection will be produced through it
	case len(c.Host) == 0:
		return nil, fmt.Errorf("host cannot be empty")
	case len(c.Username) == 0:
		return nil, fmt.Errorf("username cannot be empty")
	case len(c.Password) == 0:
		return nil, fmt.Errorf("password cannot be empty")
	}

	// Don't escape special characters for YugabyteDB password
	// Also don't escape special characters for the username and password if
	// the disable_escaping parameter is set to true
	username := c.Username
	password := c.Password

	// QueryHelper doesn't do any SQL escaping, but if it starts to do so
	// then maybe we won't be able to use it to do URL substitution any more.
	c.ConnectionURL = dbutil.QueryHelper(c.ConnectionURL, map[string]string{
		"username": username,
		"password": password,
	})

	// Set initialized to true at this point since all fields are set,
	// and the connection can be established at a later time.
	c.Initialized = true

	if verifyConnection {
		if _, err := c.Connection(ctx); err != nil {
			return nil, fmt.Errorf("error verifying connection: %s", err)
		}

		if err := c.db.PingContext(ctx); err != nil {
			return nil, fmt.Errorf("error verifying connection: %s", err)
		}
	}

	return c.RawConfig, nil
}

func (c *YugabyteConnectionProducer) Connection(ctx context.Context) (interface{}, error) {
	if !c.Initialized {
		return nil, ErrNotInitialized
	}

	// If we already have a DB, test it and return
	if c.db != nil {
		if err := c.db.PingContext(ctx); err == nil {
			return c.db, nil
		}
		// If the ping was unsuccessful, close it and ignore errors as we'll be
		// reestablishing anyways
		c.db.Close()
	}

	var conn string
	if c.TopologyKeys != "" {
		conn = fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable load_balance=%v yb_servers_refresh_interval=%d topology_keys=%s ", c.Host, c.Port, c.Username, c.Password, c.DbName, c.LoadBalance, c.YbServersRefreshInterval, c.TopologyKeys)
	} else {
		conn = fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable load_balance=%v yb_servers_refresh_interval=%d ", c.Host, c.Port, c.Username, c.Password, c.DbName, c.LoadBalance, c.YbServersRefreshInterval)
	}

	if len(c.ConnectionURL) != 0 {
		conn = c.ConnectionURL
	}

	//attempt to make connection
	var err error
	c.db, err = sql.Open("pgx", conn)
	if err != nil {
		return nil, err
	}

	return c.db, nil
}

// Close attempts to close the connection
func (c *YugabyteConnectionProducer) Close() error {
	// Grab the write lock
	c.Lock()
	defer c.Unlock()

	if c.db != nil {
		c.db.Close()
	}

	c.db = nil

	return nil
}
