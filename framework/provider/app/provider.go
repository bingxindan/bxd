package app

import (
	"bxd/framework"
	"bxd/framework/contract"
)

type BxdAppProvider struct {
	app *BxdApp

	BasePath string
}

func (provider *BxdAppProvider) Register(c framework.Container) framework.NewInstance {
	return NewBxdApp
}

func (provider *BxdAppProvider) Boot(c framework.Container) {

}

func (provider *BxdAppProvider) IsDefer() bool {
	return false
}

func (provider *BxdAppProvider) Params() []interface{} {
	return []interface{}{provider.BasePath}
}

func (provider *BxdAppProvider) Name() string {
	return contract.AppKey
}