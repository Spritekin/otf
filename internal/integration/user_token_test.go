package integration

import (
	"context"
	"testing"

	"github.com/leg100/otf/internal"
	"github.com/leg100/otf/internal/auth"
	"github.com/leg100/otf/internal/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserToken(t *testing.T) {
	t.Parallel()

	// perform all actions as superuser
	ctx := internal.AddSubjectToContext(context.Background(), &auth.SiteAdmin)

	t.Run("create", func(t *testing.T) {
		svc := setup(t, nil)
		// create user and then add them to context so that it is their token
		// that is created.
		ctx := internal.AddSubjectToContext(ctx, svc.createUser(t, ctx))
		_, _, err := svc.CreateUserToken(ctx, tokens.CreateUserTokenOptions{
			Description: "lorem ipsum...",
		})
		require.NoError(t, err)
	})

	t.Run("list", func(t *testing.T) {
		svc := setup(t, nil)
		user := svc.createUser(t, ctx)
		// create user and then add them to context so that it is their token
		// that is created.
		ctx := internal.AddSubjectToContext(ctx, user)

		svc.createToken(t, ctx, user)
		svc.createToken(t, ctx, user)
		svc.createToken(t, ctx, user)

		got, err := svc.ListUserTokens(ctx)
		require.NoError(t, err)

		assert.Equal(t, 3, len(got))
	})

	t.Run("delete", func(t *testing.T) {
		svc := setup(t, nil)
		user := svc.createUser(t, ctx)
		// create user and then add them to context so that it is their token
		// that is created.
		ctx := internal.AddSubjectToContext(ctx, user)
		token, _ := svc.createToken(t, ctx, user)

		err := svc.DeleteUserToken(ctx, token.ID)
		require.NoError(t, err)
	})
}
