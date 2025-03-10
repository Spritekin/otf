#!/usr/bin/env bash

# Run latest upstream go-tfe tests against tfe. The main bash script defaults
# to running tests against our fork of the go-tfe repo, but we are half-way
# through transitioning from our fork to the latest upstream.

set -e

function join_by { local IFS="$1"; shift; echo "$*"; }

export GO_TFE_REPO=github.com/hashicorp/go-tfe@latest

tests=()
tests+=('TestStateVersionOutputsRead')
tests+=('TestOrganizationTags')
tests+=('TestWorkspaces_(Add|Remove)Tags')
tests+=('TestWorkspacesList/when_searching_using_a_tag')
tests+=('TestNotificationConfigurationCreate/with_a')
tests+=('TestNotificationConfigurationCreate/without_a')
tests+=('TestNotificationConfigurationDelete')
tests+=('TestNotificationConfigurationUpdate/with_options')
tests+=('TestNotificationConfigurationUpdate/without_options')
tests+=('TestNotificationConfigurationUpdate/^when')
all=$(join_by '|' "${tests[@]}")

./hack/go-tfe-tests.bash $all
