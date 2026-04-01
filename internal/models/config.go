package models

type Config struct {
	Server     ServerConfig
	Logging    LoggingConfig
	Identities []Identity
	Actions    []Action
	ToolGroups []ToolGroup
}

type ServerConfig struct {
	Listen                string
	ReadTimeoutSeconds    int
	WriteTimeoutSeconds   int
	ShutdownTimeoutSeconds int
}

type LoggingConfig struct {
	Format string
	Level  string
}

type Identity struct {
	Name        string
	Username    string
	Email       string
	DisplayName  string
	Principal    string
	Domain       string
	Groups       []string
	Enabled      bool
	Secrets      map[string]string
}

type Action struct {
	Name          string
	Type          string
	Identity      string
	Description   string
	Enabled       bool
	TimeoutSeconds int
	InputSchema   map[string]any
}

type ToolGroup struct {
	Name           string
	Description    string
	Enabled        bool
	AccessTokenEnv string
	Actions        []string
}

func (c Config) HasIdentities() bool {
	return len(c.Identities) > 0
}

func (c Config) HasActions() bool {
	return len(c.Actions) > 0
}

func (c Config) HasToolGroups() bool {
	return len(c.ToolGroups) > 0
}
