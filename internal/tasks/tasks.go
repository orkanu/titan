package tasks

import (
	"fmt"
	"sync"
	"titan/internal/container"
	"titan/internal/utils"
	"titan/pkg/types"
)

func StartTasks(container *container.Container, wg *sync.WaitGroup) {

	for _, task := range container.ConfigData.Profile.Tasks {
		go func() {
			wg.Add(1)
			defer wg.Done()

			// We only have application type tasks. If we ever add any other type we should add the relevant logic here
			app, err := getApp(container, task.Name)
			if err != nil {
				// TODO handle error
				return
			}

			action, err := getAppAction(app, task.Action)
			if err != nil {
				// TODO handle error
				return
			}

			fmt.Printf("Action [%v] on project [%v]\n", task.Action, app.Name)

			options := utils.NewExecCommandOptions(container.SharedEnvironment, app.Path, action.Command, action.Args...)
			err = utils.ExecCommand(options)
			if err != nil {
				// TODO handle error
				return
			}
		}()
	}
}

func getApp(container *container.Container, appName string) (*types.Application, error) {
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
