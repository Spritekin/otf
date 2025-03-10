package cli

import (
	"bytes"
	"testing"

	"github.com/leg100/otf/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTeam_NewCommand(t *testing.T) {
	team := &auth.Team{Name: "owners", Organization: "acme-corp"}
	cmd := fakeApp(withTeam(team)).teamNewCommand()

	cmd.SetArgs([]string{"owners", "--organization", "acme-corp"})
	got := bytes.Buffer{}
	cmd.SetOut(&got)
	require.NoError(t, cmd.Execute())

	assert.Equal(t, "Successfully created team owners\n", got.String())
}

func TestTeam_DeleteCommand(t *testing.T) {
	team := &auth.Team{Name: "owners", Organization: "acme-corp"}
	cmd := fakeApp(withTeam(team)).teamDeleteCommand()

	cmd.SetArgs([]string{"owners", "--organization", "acme-corp"})
	got := bytes.Buffer{}
	cmd.SetOut(&got)
	require.NoError(t, cmd.Execute())

	assert.Equal(t, "Successfully deleted team owners\n", got.String())
}

func TestTeam_AddMembership(t *testing.T) {
	team := &auth.Team{Name: "owners", Organization: "acme-corp"}
	cmd := fakeApp(withTeam(team)).addTeamMembershipCommand()

	cmd.SetArgs([]string{"bobby", "--organization", "acme-corp", "--team", "owners"})
	got := bytes.Buffer{}
	cmd.SetOut(&got)
	require.NoError(t, cmd.Execute())

	assert.Equal(t, "Successfully added bobby to owners\n", got.String())
}

func TestTeam_RemoveMembership(t *testing.T) {
	team := &auth.Team{Name: "owners", Organization: "acme-corp"}
	cmd := fakeApp(withTeam(team)).deleteTeamMembershipCommand()

	cmd.SetArgs([]string{"bobby", "--organization", "acme-corp", "--team", "owners"})
	got := bytes.Buffer{}
	cmd.SetOut(&got)
	require.NoError(t, cmd.Execute())

	assert.Equal(t, "Successfully removed bobby from owners\n", got.String())
}
