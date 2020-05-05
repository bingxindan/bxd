package contract

const (
	EnvProduction = "production"
	EnvTesting = "testing"
	EnvDevelopment = "development"
	EnvKey = "env"
)

type Env interface {
	AppEnv() string
	AppDebug() bool
	AppURL() string
	IsExist(string) bool
	Get(string) string
	All() map[string]string
}