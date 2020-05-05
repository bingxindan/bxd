package framework

import (
	"github.com/pkg/errors"
	"sync"
)

type Container interface {
	Bind(provider ServiceProvider, isSingleton bool) error

	Singleton(provider ServiceProvider) error

	IsBind(key string) bool

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
	return bxd.Bind(provider, true)
}

func (bxd *BxdContainer) IsBind(key string) bool {
	return bxd.findServiceProvider(key) != nil
}

func (bxd *BxdContainer) findServiceProvider(key string) ServiceProvider {
	for _, sp := range bxd.providers {
		if sp.Name() == key {
			return sp
		}
	}
	return nil
}

func (bxd *BxdContainer) Make(key string) (interface{}, error) {
	return bxd.make(key, nil)
}

func (bxd *BxdContainer) make(key string, params []interface{}) (interface{}, error) {
	if bxd.findServiceProvider(key) == nil {
		return nil, errors.New("contract " + key + " have not register")
	}

	if ins, ok := bxd.instances[key]; ok {
		return ins, nil
	}

	method := bxd.methods[key]
	prov := bxd.findServiceProvider(key)
	isSingle := bxd.isSingletons[key]
	prov.Boot(bxd)

	if params == nil {
		params = prov.Params()
	}

	ins, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if isSingle {
		bxd.instances[key] = ins
		return ins, nil
	}

	return ins, nil
}

func (bxd *BxdContainer) MustMake(key string) interface{} {
	serv, err := bxd.make(key, nil)
	if err != nil {
		panic(err)
	}
	return serv
}

func (bxd *BxdContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return bxd.make(key, params)
}
