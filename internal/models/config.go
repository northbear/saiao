package models

type Config struct {
	Server     ServerConfig  `yaml:"server"`
	Logging    LoggingConfig `yaml:"logging"`
	Identities []Identity    `yaml:"identities"`
	Actions    []Action      `yaml:"actions"`
	ToolGroups []ToolGroup   `yaml:"tool_groups"`
}

type ServerConfig struct {
	Listen                 string `yaml:"listen"`
	ReadTimeoutSeconds     int    `yaml:"read_timeout_seconds"`
	WriteTimeoutSeconds    int    `yaml:"write_timeout_seconds"`
	ShutdownTimeoutSeconds int    `yaml:"shutdown_timeout_seconds"`
}

type LoggingConfig struct {
	Format string `yaml:"format"`
	Level  string `yaml:"level"`
}

type Identity struct {
	Name        string            `yaml:"name"`
	Username    string            `yaml:"username"`
	Email       string            `yaml:"email"`
	DisplayName string            `yaml:"display_name"`
	Principal   string            `yaml:"principal"`
	Domain      string            `yaml:"domain"`
	Groups      []string          `yaml:"groups"`
	Enabled     bool              `yaml:"enabled"`
	Secrets     map[string]string `yaml:"secrets"`
}

type Action struct {
	Name             string            `yaml:"name"`
	Type             string            `yaml:"type"`
	Identity         string            `yaml:"identity"`
	Description      string            `yaml:"description"`
	Enabled          bool              `yaml:"enabled"`
	TimeoutSeconds   int               `yaml:"timeout_seconds"`
	InputSchema      map[string]any    `yaml:"input_schema"`
	Method           string            `yaml:"method"`
	URL              string            `yaml:"url"`
	Headers          map[string]string `yaml:"headers"`
	BodyTemplate     string            `yaml:"body_template"`
	Host             string            `yaml:"host"`
	Port             int               `yaml:"port"`
	CommandTemplate  string            `yaml:"command_template"`
	WorkingDirectory string            `yaml:"working_directory"`
	Environment      map[string]string `yaml:"environment"`
	To               []string          `yaml:"to"`
	Cc               []string          `yaml:"cc"`
	Bcc              []string          `yaml:"bcc"`
	SubjectTemplate  string            `yaml:"subject_template"`
}

type ToolGroup struct {
	Name           string   `yaml:"name"`
	Description    string   `yaml:"description"`
	Enabled        bool     `yaml:"enabled"`
	AccessTokenEnv string   `yaml:"access_token_env"`
	Actions        []string `yaml:"actions"`
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
