package tasks

import (
	"fmt"
	"titan/internal/core"
	"titan/internal/utils"
	"titan/pkg/types"
)

func StartTasks(errorChannel chan error, container *core.Container) {

	for _, task := range container.ConfigData.Profile.Tasks {
		go func() {
			// We only have application type tasks. If we ever add any other type we should add the relevant logic here
			app, err := getApp(container, task.Name)
			if err != nil {
				errorChannel <- err
				return
			}

			action, err := getAppAction(app, task.Action)
			if err != nil {
				errorChannel <- err
				return
			}
			container.Logger.Info("Action executed on project", "action", task.Action, "project", app.Name)

			options := utils.NewExecCommandOptions(container.SharedEnvironment, app.Path, action.Command, action.Args...)
			err = utils.ExecCommand(options)
			if err != nil {
				errorChannel <- err
				return
			}
		}()
	}
}

func getApp(container *core.Container, appName string) (*types.Application, error) {
	if app, found := container.ConfigData.Config.Server.Applications[appName]; found {
		return &app, nil
	}

	return nil, fmt.Errorf("Application [%v] not found in config", appName)
}

func getAppAction(app *types.Application, actionName string) (*types.ActionData, error) {
	if action, found := app.Actions[actionName]; found {
		return &action, nil
	}

	return nil, fmt.Errorf("Action [%v] not found in Application [%v] config", actionName, app.Name)
}
