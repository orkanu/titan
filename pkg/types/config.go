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

// Config struct for titan
type Config struct {
	Versions Versions `yaml:"versions"`
	// Base path where the repositories are located
	BasePath string `yaml:"base_path"`

	// List of respositories
	Repositories map[string]string `yaml:"repositories"`

	// Proxy server configuration
	Server Server `yaml:"server"`
}
