package sql

import (
	"context"
	"net/url"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetDefaultMaxConnections(t *testing.T) {
	tests := []struct {
		name    string
		connstr string
		want    string
	}{
		{"empty dsn", "", "pool_max_conns=20"},
		{"non-empty dsn", "user=louis host=localhost", "user=louis host=localhost pool_max_conns=20"},
		{"postgres url", "postgres:///otf", "postgres:///otf?pool_max_conns=20"},
		{"postgresql url", "postgresql:///otf", "postgresql:///otf?pool_max_conns=20"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := setDefaultMaxConnections(tt.connstr, 20)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestWaitAndLock tests acquiring a connection from a pool, obtaining a session
// lock and then releasing lock and the connection, and it does this several
// times, to demonstrate that it is returning resources and not running into
// limits.
func TestWaitAndLock(t *testing.T) {
	ctx := context.Background()
	db, err := New(ctx, Options{
		Logger:     logr.Discard(),
		ConnString: NewTestDB(t),
	})
	require.NoError(t, err)
	t.Cleanup(db.Close)

	for i := 0; i < 100; i++ {
		func() {
			err := db.WaitAndLock(ctx, 123, func() error { return nil })
			require.NoError(t, err)
		}()
	}
}

func TestConnStringParse(t *testing.T) {
	got := "postgres:///"
	u, err := url.Parse(got)
	u.Path = "/mydb"
	require.NoError(t, err)
	assert.Equal(t, "postgres:///mydb", u.String())
}
