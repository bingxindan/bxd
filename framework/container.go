package framework

import (
	"github.com/pkg/errors"
	"sync"
)

type Container interface {
	Bind(provider ServiceProvider, isSingleton bool) error

	Singleton(provider ServiceProvider) error

	IsBind(key string) error

	Make(key string) (interface{}, error)

	MustMake(key string) interface{}

	MakeNew(key string, params []interface{}) (interface{}, error)
}

type BxdContainer struct {
	Container
	providers    []ServiceProvider
	instances    map[string]interface{}
	methods      map[string]NewInstance
	isSingletons map[string]bool

	lock sync.RWMutex
}

func NewBxdContaier() *BxdContainer {
	return &BxdContainer{
		providers:    []ServiceProvider{},
		instances:    map[string]interface{}{},
		methods:      map[string]NewInstance{},
		isSingletons: map[string]bool{},
		lock:         sync.RWMutex{},
	}
}

func (bxd *BxdContainer) Bind(provider ServiceProvider, isSingleton bool) error {
	bxd.lock.RLock()
	defer bxd.lock.RUnlock()
	key := provider.Name()

	bxd.providers = append(bxd.providers, provider)
	bxd.isSingletons[key] = isSingleton
	bxd.methods[key] = provider.Register(bxd)

	if provider.IsDefer() == false {
		provider.Boot(bxd)
		params := provider.Params()
		method := bxd.methods[key]
		instance, err := method(params...)
		if err != nil {
			return errors.New(err.Error())
		}
		if isSingleton == true {
			bxd.instances[key] = instance
		}
	}
	return nil
}

func (bxd *BxdContainer) Singleton(provider ServiceProvider) error {
	return nil
}

func (bxd *BxdContainer) IsBind(key string) error {
	return nil
}

func (bxd *BxdContainer) findServiceProvider(key string) ServiceProvider {
	return nil
}

func (bxd *BxdContainer) Make(key string) (interface{}, error) {
	return key, nil
}

func (bxd *BxdContainer) MustMake(key string) interface{} {
	return nil
}

func (bxd *BxdContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return key, nil
}