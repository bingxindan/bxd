package app

type BxdApp struct {
	basePath string
}

func NewBxdApp(params ...interface{}) (interface{}, error) {
	var basePath string
	if len(params) == 1 {
		basePath = params[0].(string)
	}
	return &BxdApp{basePath: basePath}, nil
}

func (app *BxdApp) Version() string {
	return "0.0.1"
}

func (app *BxdApp) BasePath() string {
	return app.basePath
}

func (app *BxdApp) ConfigPath() string {
	return app.BasePath() + "config/"
}

func (app *BxdApp) EnvironmentPath() string {
	return app.BasePath()
}

func (app *BxdApp) StoragePath() string {
	return app.BasePath() + "storage/"
}

func (app *BxdApp) LogPath() string {
	return app.BasePath() + "logs/"
}
