package integration

import (
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/require"
)

// TestAutoApply tests auto-apply functionality, using the UI to enable
// auto-apply on a workspace first before invoking 'terraform apply'.
func TestAutoApply(t *testing.T) {
	t.Parallel()

	svc := setup(t, nil)
	user, ctx := svc.createUserCtx(t, ctx)
	org := svc.createOrganization(t, ctx)

	// create workspace and enable auto-apply
	browser := createBrowserCtx(t)
	err := chromedp.Run(browser, chromedp.Tasks{
		newSession(t, ctx, svc.Hostname(), user.Username, svc.Secret),
		createWorkspace(t, svc.Hostname(), org.Name, t.Name()),
		chromedp.Tasks{
			// go to workspace
			chromedp.Navigate(workspaceURL(svc.Hostname(), org.Name, t.Name())),
			screenshot(t),
			// go to workspace settings
			chromedp.Click(`//a[text()='settings']`, chromedp.NodeVisible),
			screenshot(t),
			// enable auto-apply
			chromedp.Click("input#auto_apply", chromedp.NodeVisible, chromedp.ByQuery),
			screenshot(t),
			// submit form
			chromedp.Click(`//button[text()='Save changes']`, chromedp.NodeVisible),
			screenshot(t),
			// confirm workspace updated
			matchText(t, ".flash-success", "updated workspace"),
		},
	})
	require.NoError(t, err)

	// create terraform config
	configPath := newRootModule(t, svc.Hostname(), org.Name, t.Name())
	svc.tfcli(t, ctx, "init", configPath)
	// terraform apply - note we are not passing the -auto-approve flag yet we
	// expect it to auto-apply because the workspace is set to auto-apply.
	out := svc.tfcli(t, ctx, "apply", configPath)
	require.Contains(t, string(out), "Apply complete! Resources: 1 added, 0 changed, 0 destroyed.")
}
