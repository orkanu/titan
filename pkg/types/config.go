package types

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
