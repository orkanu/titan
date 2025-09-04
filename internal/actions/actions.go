package actions

import (
	"fmt"
	"log/slog"
	"strings"
	"titan/internal/utils"
	"titan/pkg/parser"
	"titan/pkg/types"
)

type ExecOptions struct {
	logger        *slog.Logger
	repoAction    *types.RepoAction
	repoPath      string
	projectName   string
	env           []string
	scriptsOutput string
}

func NewExecOptions(
	logger *slog.Logger,
	env []string,
	repoAction *types.RepoAction,
	repoPath string,
	projectName string,
	scriptsOutput string,
) *ExecOptions {
	return &ExecOptions{
		logger:        logger,
		env:           env,
		repoAction:    repoAction,
		repoPath:      repoPath,
		projectName:   projectName,
		scriptsOutput: scriptsOutput,
	}
}

// Action defines what an action is
type Action interface {
	// Name gives the name of the action
	Name() string
	// ShouldExecute evaluates is the action has to be executed based on the command requested
	ShouldExecute(command types.Action) bool
	// Execute executes the action
	Execute(options *ExecOptions) error
}

func getScriptFromConfig(actionName string, repoAction *types.RepoAction, parserCtx map[string]any, defaultScript string, logger *slog.Logger) string {
	var sb strings.Builder
	if repoAction != nil {
		logger.Debug("using configured repository command actions", "command", actionName)
		for _, cmd := range repoAction.Commands {
			if cmd.Condition == "" {
				sb.WriteString(cmd.Value)
				continue
			}
			p := parser.NewParser(cmd.Condition, parserCtx)
			if p.ParseExpression() {
				sb.WriteString(cmd.Value)
			} else {
				logger.Debug("skipping action due unmet condition", "command", actionName, "condition", cmd.Condition)
			}
		}
	} else {
		logger.Debug("using default repository command actions", "command", actionName)
		sb.WriteString(defaultScript)
	}

	return sb.String()
}

func executeScript(actionName string, scriptFromConfig string, logger *slog.Logger, repoPath string, projectName string, env []string) error {
	var sb strings.Builder
	sb.WriteString(`
		#!/bin/bash
		set -e
	`)
	sb.WriteString(scriptFromConfig)
	script := sb.String()

	logger.Info("executing action", "action", actionName, "project", projectName)
	if err := utils.ExecScript(script, env, repoPath); err != nil {
		return fmt.Errorf("failed executing [%v] action script: %v", actionName, err)
	}
	return nil
}
