package redisclient

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	DB             int
	UserName       string
	Password       string
	UseSsl         bool
	SkipVerifyCert bool

	// Sentinel mode

	MasterName    string
	SentinelAddrs string

	// Standalone mode

	Addr string

	// Cluster mode
	ClusterAddrs string
}

func New(c *Config) (redis.UniversalClient, error) {
	var tlsConfig *tls.Config
	if c.UseSsl {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: c.SkipVerifyCert,
			MinVersion:         tls.VersionTLS12,
		}
	}

	var rCli redis.UniversalClient
	if c.MasterName == "" {
		// standalone mode
		if c.Addr != "" {
			rCli = redis.NewClient(
				&redis.Options{
					Addr:      c.Addr,
					Username:  c.UserName,
					Password:  c.Password,
					DB:        c.DB,
					TLSConfig: tlsConfig,
				},
			)
		} else if c.ClusterAddrs != "" {
			// cluster mode
			rCli = redis.NewClusterClient(
				&redis.ClusterOptions{
					Addrs:     strings.Split(c.ClusterAddrs, ","),
					Username:  c.UserName,
					Password:  c.Password,
					TLSConfig: tlsConfig,
				},
			)
		} else {
			return nil, fmt.Errorf("empty redis config")
		}
	} else {
		// sentinel mode
		if c.SentinelAddrs == "" {
			return nil, fmt.Errorf("empty redis config sentinel_addrs")
		}

		rCli = redis.NewFailoverClient(
			&redis.FailoverOptions{
				MasterName:    c.MasterName,
				SentinelAddrs: strings.Split(c.SentinelAddrs, ","),
				Password:      c.Password,
				DB:            c.DB,
				TLSConfig:     tlsConfig,
			},
		)
	}

	return rCli, nil
}
