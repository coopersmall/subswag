package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/coopersmall/subswag/utils"
)

type IDBManager interface {
	ReadOnly() *sql.DB
	ReadWrite() *sql.DB
	SetSchema(schema string) error
	WaitForConnection(ctx context.Context) error
	Shutdown() error
}

type DBManager struct {
	readonly  *sql.DB
	readwrite *sql.DB
}

func NewDBManager(
	vars iEnvVars,
	clients iClients,
) (*DBManager, error) {
	sharedPostgresURL, err := vars.GetSharedPostgresURL()
	if err != nil {
		return nil, err
	}
	standardPostgresURL, err := vars.GetStandardPostgresURL()
	if err != nil {
		return nil, err
	}
	readonly, err := clients.PostgresClient(
		sharedPostgresURL,
		&struct {
			MaxOpenConns int
			MaxIdleConns int
		}{
			MaxOpenConns: 25,
			MaxIdleConns: 25,
		},
	)
	if err != nil {
		return nil, err
	}

	readwrite, err := clients.PostgresClient(
		standardPostgresURL,
		&struct {
			MaxOpenConns int
			MaxIdleConns int
		}{
			MaxOpenConns: 25,
			MaxIdleConns: 25,
		},
	)
	if err != nil {
		return nil, err
	}

	return &DBManager{
		readonly:  readonly,
		readwrite: readwrite,
	}, nil
}

func (m *DBManager) ReadOnly() *sql.DB {
	return m.readonly
}

func (m *DBManager) ReadWrite() *sql.DB {
	return m.readwrite
}

func (m *DBManager) SetSchema(schema string) error {
	tx, err := m.readonly.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(schema)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (m *DBManager) WaitForConnection(ctx context.Context) error {
	for i := 0; i < 30; i++ {
		if err := m.readwrite.PingContext(ctx); err == nil {
			return nil
		}
		time.Sleep(time.Second)
	}
	return utils.NewInternalError("failed to connect to db")
}

func (m *DBManager) Shutdown() error {
	if m.readonly != nil {
		if err := m.readonly.Close(); err != nil {
			return err
		}
	}
	if m.readwrite != nil {
		if err := m.readwrite.Close(); err != nil {
			return err
		}
	}
	return nil
}

type iEnvVars interface {
	GetSharedPostgresURL() (string, error)
	GetStandardPostgresURL() (string, error)
}

type iClients interface {
	PostgresClient(url string, opts *struct {
		MaxOpenConns int
		MaxIdleConns int
	}) (*sql.DB, error)
}
