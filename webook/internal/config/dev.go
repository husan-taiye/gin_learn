//go:build !k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:l913687515@tcp(localhost:30001)/webook",
	},
	Redis: RedisConfig{
		Addr: "localhost:30003",
	},
}
