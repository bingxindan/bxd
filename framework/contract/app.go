package contract

const AppKey = "app"

type App interface {
	Version() string
	BasePath() string
	ConfigPath() string
	EnvironmentPath() string
	StoragePath() string
	LogPath() string
}
