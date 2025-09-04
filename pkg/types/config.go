package types

type ActionData struct {
	Command string   `yaml:"command"`
	Args    []string `yaml:"args"`
}

type Application struct {
	Name    string                `yaml:"name"`
	Path    string                `yaml:"path"`
	Actions map[string]ActionData `yaml:"actions"`
}

type Profile struct {
	Parameters map[string]string `yaml:"parameters"`
	Tasks      []struct {
		Type   string `yaml:"type"`
		Name   string `yaml:"name"`
		Action string `yaml:"action"`
	} `yaml:"tasks"`
	Routes []string `yaml:"routes"`
}

type Server struct {
	// Host value
	Host string `yaml:"host"`
	// Port value
	Port int `yaml:"port"`
	// SSL configuration
	SSL struct {
		Port int    `yaml:"port"`
		Cert string `yaml:"cert"`
		Key  string `yaml:"key"`
	} `yaml:"ssl"`

	// Routes to proxy to
	Routes map[string]struct {
		Source string `yaml:"source"`
		Target string `yaml:"target"`
	} `yaml:"routes"`

	Applications map[string]Application `yaml:"applications"`

	Profiles map[string]Profile `yaml:"profiles"`
}

type Versions struct {
	// Node version to install/use in system via NVM
	Node string `yaml:"node"`
	// PNOM version to install
	PNPM string `yaml:"pnpm"`
}

type RepoCommands struct {
	Value     string `yaml:"value"`
	Condition string `yaml:"condition,omitempty"`
}

type RepoAction struct {
	Commands []RepoCommands `yaml:"commands"`
}

type RepoActions struct {
	// map[string]RepoAction
	// List of respositories
	ScriptsOutput string                 `yaml:"scripts-output,omitempty"`
	Repositories  map[string]string      `yaml:"repositories"`
	Actions       map[string]*RepoAction `yaml:"actions"`
}

// Config struct for titan
type Config struct {
	Versions Versions `yaml:"versions"`

	// Repository actions
	RepoActions RepoActions `yaml:"repo-actions,omitempty"`

	// Proxy server configuration
	Server Server `yaml:"server"`
}
