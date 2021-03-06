package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/TouchBistro/cannon/action"
	"github.com/stretchr/testify/assert"
)

func setup() {
	tmp := os.TempDir()
	os.RemoveAll(tmp + ".cannon")
	os.Setenv("HOME", strings.TrimSuffix(tmp, "/"))
}

func TestInitConfig(t *testing.T) {
	setup()

	reader := strings.NewReader(`title: PR title
branchName: new-example-branch
commitMessage: Committing cannon changes
description: |-
  Description line 1

  Description line 2
reviewers:
  users:
    - user1
  teams:
    - team1
repos:
  - name: TouchBistro/cannon
  - name: TouchBistro/example
    base: develop
actions:
  - type: replaceLine
    source: DB_USER=core
    target: DB_USER=SA
    path: .env.example
  - type: runCommand
    run: yarn
`)

	expectedConfig := CannonConfig{
		Title:         "PR title",
		BranchName:    "new-example-branch",
		CommitMessage: "Committing cannon changes",
		Description: `Description line 1

Description line 2`,
		Reviewers: map[string][]string{
			"users": {
				"user1",
			},
			"teams": {
				"team1",
			},
		},
		Repos: []Repo{
			Repo{
				Name: "TouchBistro/cannon",
			},
			Repo{
				Name: "TouchBistro/example",
				Base: "develop",
			},
		},
		Actions: []action.Action{
			action.Action{
				Type:   action.ActionReplaceLine,
				Source: "DB_USER=core",
				Target: "DB_USER=SA",
				Path:   ".env.example",
			},
			action.Action{
				Type: action.ActionRunCommand,
				Run:  "yarn",
			},
		},
	}
	expectedCannonDir := filepath.Join(os.TempDir(), ".cannon")

	err := Init(reader)

	assert := assert.New(t)

	assert.NoError(err)
	assert.DirExists(filepath.Join(os.TempDir(), ".cannon"))
	assert.Equal(expectedConfig, *Config())
	assert.Equal(expectedCannonDir, CannonDir())
	assert.Equal("master", Config().Repos[0].BaseBranch())
	assert.Equal("develop", Config().Repos[1].BaseBranch())
}
