# Configuration

Titan requires a configuration file in YAML format. By default it tries to use `titan.yaml` on the same place where the
binary is executed.

Via a flag, we can instruct titan another path and file name to get the config from. For example:

```bash
titan -c /path/to/my/config/file.yaml
# OR
titan <command> -c /path/to/my/config/file.yaml
```

A sample config file can be found [here](../sample.titan.yaml) which resembles the one used
whilst developing titan

## Sections

**root**

| Section      |Description                                                                         | Required |
| ------------ | ---------------------------------------------------------------------------------- | -------- |
| versions     | indicates the version to use for some tooling, like pnpm or node                   | ✅       |
| repo-actions | for each repository action titan can do, like fetch or build, it allows configure  | ➖       |
|              | the commands to execute, even using conditionals to add some commands or not when  |          |
|              | the condition/s is/are met. If any is missing, default values are used instead     |          |
| server       | it has all the data to run the proxy server as well as tasks                       | ✅       |


**versions**

| Section | Description                                                      | Required |
| ------- | ---------------------------------------------------------------- | -------- |
|node     | NodeJS version to use                                            | ✅       |
|pnpm     | PNPM version to use so it can be installed globally when setting | ✅       |
|         | the environment for the scripts                                  |          |


**repo-actions**

| Section        |Description                                                                   | Required |
| -------------- | ---------------------------------------------------------------------------- | -------- |
| repositories   | indicates the repositories that will be affected by the actions              | ✅       |
| scripts-output | indicates if the output of the scripts run for each action should be dumped  | ➖       |
|                | the console or to a file. Defaults to console                                |          |
| actions        | we can define specific configuration for each action: fetch, install, bild   | ➖       |
|                | and clean. See **actions** section for specific                              |          |

**actions**
| Section |Description                                                                   | Required |
| ------- | ---------------------------------------------------------------------------- | -------- |
| fetch   | we can indicate the commands to run for the action                           | ➖       |
| install | we can indicate the commands to run for the action                           | ➖       |
| build   | we can indicate the commands to run for the action                           | ➖       |
| clean   | we can indicate the commands to run for the action                           | ➖       |

**comands**
| Section   |Description                                                                   | Required |
| --------- | ---------------------------------------------------------------------------- | -------- |
| value     | what command will be executed. For instance `pnpm install` could be used for | ✅       |
|           | the install action                                                           |          |
| condition | we can provide simple conditions that would be evaluated to include or not   | ➖       |
|           | the associateed **value** into the script to execute                         |          |

**condition**
Currently, the only available condition token is to use `projectName`. If more tokens are required, those
would need to be taken into consideration code wise.
For now we keep it simple which means that adding a new condition token would mean a code change. If we see that
we need a lot of them, we may do some research to see if it can be done via mere configuration to avoid having to
change the code every time.
