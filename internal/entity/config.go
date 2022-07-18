package entity

import (
	"os"
)

type EnvVar string

type Config struct {
	PostgresUser		EnvVar
	PostgresPassword	EnvVar
	PostgresDBName		EnvVar
	PostgresHost		EnvVar
	PostgresPort		EnvVar
	RouterHost			EnvVar
	RouterPort			EnvVar
}

func (ev *EnvVar) GetEnv(key string) {
	if val, ok := os.LookupEnv(key); ok {
		ev = EnvVar(val)
	} else {
		panic("Environment variable `" + key + "` is not set. Check spelling or injection")
	}
}

func (ev *EnvVar) String() string {
	return string(*ev)
}
