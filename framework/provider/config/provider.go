package config

import (
	"bxd/framework"
	"bxd/framework/contract"
)

type BxdConfigProvider struct {
	Folder string
	envMaps map[string]string
}

func (provider *BxdConfigProvider) Register(c framework.Container) framework.NewInstance {
	return NewBxdConfig
}

func (provider *BxdConfigProvider) Boot(c framework.Container) {
	if provider.Folder == "" && c.IsBind(contract.AppKey) {
		provider.Folder = c.MustMake(contract.AppKey).(contract.App).ConfigPath()
	}
	if c.IsBind(contract.EnvKey) {
		provider.envMaps = c.MustMake(contract.EnvKey).(contract.Env).All()
	}
}

func (provider *BxdConfigProvider) IsDefer(c framework.Container) bool {
	return false
}

func (provider *BxdConfigProvider) Params() []interface{} {
	return []interface{}{provider.Folder, provider.envMaps}
}

func (provider *BxdConfigProvider) Name() string {
	return contract.ConfigKey
}