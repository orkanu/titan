package actions

import (
	"fmt"
	"slices"
	"titan/internal/utils"
)

type BuildAction struct {
	name     string
	commands []utils.Action
}

func NewBuildAction() BuildAction {
	return BuildAction{
		name:     "build",
		commands: []utils.Action{utils.ALL, utils.BUILD},
	}
}

func (ba BuildAction) Name() string {
	return ba.name
}

func (ba BuildAction) ShouldExecute(command utils.Action) bool {
	return slices.Contains(ba.commands, command)
}

func (ba BuildAction) Execute(repoPath string, projectName string, env []string) error {
	options := utils.NewExecCommandOptions(env, repoPath, "pnpm", "run", "build:local")
	fmt.Printf("Action [build] on project [%v]\n", projectName)
	if err := utils.ExecCommand(options); err != nil {
		return fmt.Errorf("Error executing build action script: %v", err)
	}
	return nil
}
