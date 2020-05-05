package env

import (
	"bxd/framework"
	"bxd/framework/contract"
)

type BxdEnvProvider struct {
	Folder string
}

func (provider *BxdEnvProvider) Register(c framework.Container) framework.NewInstance {
	return NewBxdEnv
}

func (provider *BxdEnvProvider) Boot(c framework.Container) {
	if provider.Folder == "" {
		if c.IsBind(contract.AppKey) {
			app := c.MustMake(contract.AppKey).(contract.App)
			provider.Folder = app.EnvironmentPath()
		}
	}
}

func (provider *BxdEnvProvider) IsDefer() bool {
	return false
}

func (provider *BxdEnvProvider) Params() []interface{} {
	return []interface{}{provider.Folder}
}

func (provider *BxdEnvProvider) Name() string {
	return contract.EnvKey
}