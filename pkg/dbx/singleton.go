package dbx

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var once sync.Once

var pgTestContainer *PostgreSQLContainer

func bootstrap(t *testing.T, opts ...PostgreSQLContainerOption) {
	ctx := context.Background()

	container, containerErr := NewPostgreSQLContainer(ctx, opts...)
	require.NoError(t, containerErr)
	require.NotEmpty(t, container)

	mErr := container.RunMigrations()
	require.NoError(t, mErr, "Database migrations failed")

	pgTestContainer = container
}

func ContainerTestInstance(t *testing.T, opts ...PostgreSQLContainerOption) (dsn string) {
	once.Do(func() {
		bootstrap(t, opts...)
	})

	return pgTestContainer.GetDSN()
}
