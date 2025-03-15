package testing

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/coopersmall/subswag/env"
	"github.com/coopersmall/subswag/utils/actions"
	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestEnv struct {
	env.IEnv
	postgres testcontainers.Container
	redis    testcontainers.Container
	envVars  *testEnvVars
	schema   string
	ctx      context.Context
}

type TestEnvOption func(*testEnvVars)

func WithExternalPostgres(url string) TestEnvOption {
	return func(vars *testEnvVars) {
		vars.SetStandardPostgresURL(url)
		vars.SetSharedPostgresURL(url)
	}
}

func WithExternalRedis(url string, password string) TestEnvOption {
	return func(vars *testEnvVars) {
		vars.SetRedisURL(url)
		vars.SetRedisPassword(password)
	}
}

func NewTestingEnvWithOptions(opts ...TestEnvOption) (*TestEnv, error) {
	ctx := context.Background()

	envVars, err := NewTestEnvVars()
	if err != nil {
		return nil, fmt.Errorf("failed to create test env vars: %w", err)
	}

	schema, err := actions.GetSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to get schema: %w", err)
	}

	for _, opt := range opts {
		opt(envVars)
	}

	return &TestEnv{
		envVars: envVars,
		schema:  schema,
		ctx:     ctx,
	}, nil
}

func (t *TestEnv) start() error {
	wg := sync.WaitGroup{}
	wg.Add(2)
	var err error
	go func() {
		defer wg.Done()
		err = t.startPostgres()
	}()
	go func() {
		defer wg.Done()
		err = t.startRedis()
	}()
	wg.Wait()

	if err != nil {
		return err
	}

	env, err := env.GetEnv(t.envVars)
	if err != nil {
		return fmt.Errorf("failed to get env: %w", err)
	}

	t.IEnv = env

	if err = t.resetDB(); err != nil {
		return err
	}

	if err = t.waitForPostgres(); err != nil {
		return err
	}

	return nil
}

func (t *TestEnv) Reset() error {
	wg := sync.WaitGroup{}

	wg.Add(2)
	var err error
	go func() {
		defer wg.Done()
		err = t.resetDB()
	}()
	go func() {
		defer wg.Done()
		err = t.resetCache()
	}()
	wg.Wait()

	return err
}

func (t *TestEnv) stop() error {
	wg := sync.WaitGroup{}
	wg.Add(2)
	var err error
	go func() {
		defer wg.Done()
		err = t.stopDB()
	}()
	go func() {
		defer wg.Done()
		if t.redis != nil {
			err = t.redis.Terminate(t.ctx)
		}
	}()
	wg.Wait()
	return err
}

func (t *TestEnv) startDB() error {
	return t.startPostgres()
}

func (t *TestEnv) resetDB() error {
	dbs := t.GetDB()
	return dbs.SetSchema(t.schema)
}

func (t *TestEnv) stopDB() error {
	return t.GetDB().Shutdown()
}

func (t *TestEnv) startCache() error {
	return t.startRedis()
}

func (t *TestEnv) resetCache() error {
	gateways, close := t.GetGateways()
	defer close()
	return gateways.RedisCacheGateway().DeleteAll(t.ctx)
}

func (t *TestEnv) stopCache() error {
	if t.redis != nil {
		return t.redis.Terminate(t.ctx)
	}
	return nil
}

func (t *TestEnv) startPostgres() error {
	postgres, err := testcontainers.GenericContainer(t.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:13",
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor: wait.ForAll(
				wait.ForListeningPort("5432/tcp"),
				wait.ForSQL("5432/tcp", "postgres", func(host string, port nat.Port) string {
					return fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", host, port.Port())
				}).WithStartupTimeout(45*time.Second),
			),
			Env: map[string]string{
				"POSTGRES_DB":       "testdb",
				"POSTGRES_USER":     "test",
				"POSTGRES_PASSWORD": "test",
			},
		},
		Started: true,
	})
	if err != nil {
		return fmt.Errorf("failed to start postgres container: %w", err)
	}

	t.postgres = postgres

	pgHost, err := postgres.Host(t.ctx)
	if err != nil {
		return fmt.Errorf("failed to get postgres host: %w", err)
	}
	pgPort, err := postgres.MappedPort(t.ctx, "5432")
	if err != nil {
		return fmt.Errorf("failed to get postgres port: %w", err)
	}

	pgURL := fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable",
		pgHost, pgPort.Port())

	t.envVars.SetStandardPostgresURL(pgURL)
	t.envVars.SetSharedPostgresURL(pgURL)
	return nil
}

func (t *TestEnv) waitForPostgres() error {
	db := t.GetDB()
	ctx, cancel := context.WithTimeout(t.ctx, 30*time.Second)
	defer cancel()
	return db.WaitForConnection(ctx)
}

func (t *TestEnv) startRedis() error {
	redis, err := testcontainers.GenericContainer(t.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:6-alpine",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForListeningPort("6379/tcp"),
		},
		Started: true,
	})
	if err != nil {
		return fmt.Errorf("failed to start redis container: %w", err)
	}

	t.redis = redis

	redisHost, err := redis.Host(t.ctx)
	if err != nil {
		return fmt.Errorf("failed to get redis host: %w", err)
	}
	redisPort, err := redis.MappedPort(t.ctx, "6379")
	if err != nil {
		return fmt.Errorf("failed to get redis port: %w", err)
	}

	t.envVars.SetRedisURL(fmt.Sprintf("%s:%s", redisHost, redisPort.Port()))
	return nil
}
