package dbx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgreSQLContainerConfig struct {
	ImageTag   string
	User       string
	Password   string
	Host       string
	MappedPort string
	Database   string
	Migrations fs.FS
}

func (c *PostgreSQLContainerConfig) Clone() *PostgreSQLContainerConfig {
	return &PostgreSQLContainerConfig{
		ImageTag:   c.ImageTag,
		User:       c.User,
		Password:   c.Password,
		Host:       c.Host,
		MappedPort: c.MappedPort,
		Database:   c.Database,
		Migrations: c.Migrations,
	}
}

func (c *PostgreSQLContainerConfig) SetDatabase(database string) *PostgreSQLContainerConfig {
	c.Database = database
	return c
}

func (c *PostgreSQLContainerConfig) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Host,
		c.MappedPort,
		c.Database,
	)
}

type PostgreSQLContainerOption func(c *PostgreSQLContainerConfig)

type PostgreSQLContainer struct {
	testcontainers.Container
	config *PostgreSQLContainerConfig
}

func NewPostgreSQLContainer(ctx context.Context, opts ...PostgreSQLContainerOption) (*PostgreSQLContainer, error) {
	const (
		psqlImage = "postgres"
		psqlPort  = "5432"
	)

	config := &PostgreSQLContainerConfig{
		ImageTag:   "14.4",
		User:       "postgres",
		Password:   "postgres",
		MappedPort: "5432",
		Database:   "test",
	}

	for _, opt := range opts {
		opt(config)
	}

	containerPort := psqlPort + "/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Env: map[string]string{
				"POSTGRES_USER":     config.User,
				"POSTGRES_PASSWORD": config.Password,
				"POSTGRES_DB":       config.Database,
			},
			ExposedPorts: []string{
				containerPort,
			},
			Image:      fmt.Sprintf("%s:%s", psqlImage, config.ImageTag),
			WaitingFor: wait.ForListeningPort(nat.Port(containerPort)),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("getting request provider: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting host for: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(containerPort))
	if err != nil {
		return nil, fmt.Errorf("getting mapped port for (%s): %w", containerPort, err)
	}
	config.MappedPort = mappedPort.Port()
	config.Host = host

	return &PostgreSQLContainer{
		Container: container,
		config:    config,
	}, nil
}

func (c *PostgreSQLContainer) RunMigrations() error {
	if c.config.Migrations == nil {
		return nil
	}

	db, dbErr := sql.Open("postgres", c.GetDSN())
	if dbErr != nil {
		return fmt.Errorf("could not run migrations: %w", dbErr)
	}
	defer db.Close()

	if mErr := Migrate(db, c.config.Migrations); errors.Is(mErr, ErrNilVersion) {
		return fmt.Errorf("failed not enforce changes to database: %w", ErrNilVersion)
	} else if errors.Is(mErr, ErrNoChange) {
		return fmt.Errorf("exceptional result from migrations: %w", ErrNoChange)
	} else if mErr != nil {
		return fmt.Errorf("migration failure: %w", mErr)
	}

	if version, dirty, vErr := Version(db); vErr != nil {
		return fmt.Errorf("could not resolve migration version: %w", vErr)
	} else if dirty {
		return fmt.Errorf("migration dirty state enabled: version %d", version)
	}

	return nil
}

func (c *PostgreSQLContainer) GetDSN() string {
	return c.config.GetDSN()
}

func (c *PostgreSQLContainer) CloneConfig() *PostgreSQLContainerConfig {
	return c.config.Clone()
}

func (c *PostgreSQLContainer) Database() string {
	return c.config.Database
}
