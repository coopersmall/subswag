package env

import (
	"sync"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/cache"
	"github.com/coopersmall/subswag/clients"
	"github.com/coopersmall/subswag/db"
	"github.com/coopersmall/subswag/gateways"
	"github.com/coopersmall/subswag/repos"
	"github.com/coopersmall/subswag/services"
	"github.com/coopersmall/subswag/streams/publishers"
	"github.com/coopersmall/subswag/streams/subscribers"
	"github.com/coopersmall/subswag/utils"
)

type IEnv interface {
	GetLogger(name string) utils.ILogger
	GetTracer(service string) apm.ITracer
	GetRepos() (repos.IRepos, func())
	GetPublishers() (publishers.IPublishers, func())
	GetSubscribers() (subscribers.ISubscribers, func())
	GetClients() (clients.IClients, func())
	GetGateways() (gateways.IGateways, func())
	GetCaches() (cache.ICache, func())
	GetServices() (services.IServices, func())
	GetUsersService() (services.IUsersService, func())
	GetAuthenticationService() (services.IAuthenticationService, func())
	GetQuerier() db.IQuerier
	GetDB() db.IDBManager
	Shutdown() error
}

type Env struct {
	logger          utils.ILogger
	db              db.IDBManager
	apm             apm.ITraceProvider
	usersPool       sync.Pool
	authPool        sync.Pool
	servicesPool    sync.Pool
	cachePool       sync.Pool
	reposPool       sync.Pool
	clientsPool     sync.Pool
	gatewaysPool    sync.Pool
	publishersPool  sync.Pool
	subscribersPool sync.Pool
}

func (e *Env) GetLogger(name string) utils.ILogger {
	return utils.NewLogger(name)
}

func (e *Env) GetTracer(service string) apm.ITracer {
	return e.apm.NewTracer(service)
}

func (e *Env) GetUsersService() (services.IUsersService, func()) {
	usersService := e.usersPool.Get()
	return usersService.(services.IUsersService), func() {
		e.usersPool.Put(usersService)
	}
}

func (e *Env) GetAuthenticationService() (services.IAuthenticationService, func()) {
	auth := e.authPool.Get()
	return auth.(services.IAuthenticationService), func() {
		e.authPool.Put(auth)
	}
}

func (e *Env) GetServices() (services.IServices, func()) {
	s := e.servicesPool.Get()
	return s.(services.IServices), func() {
		e.servicesPool.Put(s)
	}
}

func (e *Env) GetCaches() (cache.ICache, func()) {
	c := e.cachePool.Get()
	return c.(cache.ICache), func() {
		e.cachePool.Put(c)
	}
}

func (e *Env) GetRepos() (repos.IRepos, func()) {
	r := e.reposPool.Get()
	return r.(repos.IRepos), func() {
		e.reposPool.Put(r)
	}
}

func (e *Env) GetPublishers() (publishers.IPublishers, func()) {
	p := e.publishersPool.Get()
	return p.(publishers.IPublishers), func() {
		e.publishersPool.Put(p)
	}
}

func (e *Env) GetSubscribers() (subscribers.ISubscribers, func()) {
	s := e.subscribersPool.Get()
	return s.(subscribers.ISubscribers), func() {
		e.subscribersPool.Put(s)
	}
}

func (e *Env) GetQuerier() db.IQuerier {
	return db.NewQuerier(e.db)
}

func (e *Env) GetDB() db.IDBManager {
	return e.db
}

func (e *Env) GetGateways() (gateways.IGateways, func()) {
	retrieved := e.gatewaysPool.Get()
	return retrieved.(gateways.IGateways), func() {
		e.gatewaysPool.Put(retrieved)
	}
}

func (e *Env) GetClients() (clients.IClients, func()) {
	retrieved := e.clientsPool.Get()
	return retrieved.(clients.IClients), func() {
		e.clientsPool.Put(retrieved)
	}
}

func (e *Env) Shutdown() error {
	if err := e.db.Shutdown(); err != nil {
		return err
	}
	return nil
}

func MustGetEnv(envVars IEnvVars) IEnv {
	env, err := GetEnv(envVars)
	if err != nil {
		panic(err)
	}
	return env
}

func GetEnv(envVars IEnvVars) (IEnv, error) {
	logger := utils.NewLogger("env")
	c := clients.GetClients(envVars)

	dbManager, err := db.NewDBManager(
		envVars,
		c,
	)
	if err != nil {
		return nil, err
	}

	traceProvider, err := apm.GetTraceProvider(c.APMClient())
	if err != nil {
		return nil, err
	}

	env := &Env{
		logger: logger,
		db:     dbManager,
		apm:    traceProvider,
	}

	env.reposPool = sync.Pool{
		New: func() interface{} {
			return repos.GetRepos(env)
		},
	}

	env.clientsPool = sync.Pool{
		New: func() interface{} {
			return clients.GetClients(envVars)
		},
	}

	env.gatewaysPool = sync.Pool{
		New: func() interface{} {
			clients, closeClients := env.GetClients()
			defer closeClients()
			return gateways.GetGateways(env, clients)
		},
	}

	env.cachePool = sync.Pool{
		New: func() interface{} {
			gateways, closeGateways := env.GetGateways()
			defer closeGateways()
			return cache.GetCaches(env, gateways)
		},
	}

	env.servicesPool = sync.Pool{
		New: func() interface{} {
			gateways, closeGateways := env.GetGateways()
			defer closeGateways()
			repos, closeRepos := env.GetRepos()
			defer closeRepos()
			cache, closeCache := env.GetCaches()
			defer closeCache()
			publish, closePublish := env.GetPublishers()
			defer closePublish()
			return services.GetServices(
				env,
				envVars,
				repos,
				gateways,
				cache,
				publish,
			)
		},
	}

	env.subscribersPool = sync.Pool{
		New: func() interface{} {
			gateways, closeGateways := env.GetGateways()
			defer closeGateways()
			services, closeServices := env.GetServices()
			defer closeServices()
			return subscribers.GetSubscribers(env, gateways, services)
		},
	}

	env.publishersPool = sync.Pool{
		New: func() interface{} {
			gateways, closeGateways := env.GetGateways()
			defer closeGateways()
			return publishers.GetPublishers(env, gateways)
		},
	}

	env.usersPool = sync.Pool{
		New: func() interface{} {
			services, closeServices := env.GetServices()
			defer closeServices()
			return services.UsersService()
		},
	}

	env.authPool = sync.Pool{
		New: func() interface{} {
			services, closeServices := env.GetServices()
			defer closeServices()
			return services.AuthenticationService()
		},
	}

	return env, nil
}
