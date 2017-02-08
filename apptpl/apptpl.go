package apptpl

// Config implements the Template Action
type Template struct {
	App           string `yaml:"app"`
	Core          string `yaml:"core"`
	ConfigPath    string `yaml:"configPath"`
	VHost         string `yaml:"vHost"`
	DBPort        int    `yaml:"dbPort"`
	DBPrefix      string `yaml:"dbPrefix"`
	PublicFiles   string `yaml:"publicFiles"`
	PrivateFiles  string `yaml:"privateFiles"`
	ConfigSyncDir string `yaml:"configSyncDir"`
}

// App is the interface that eash plugin must implement
type App interface {
	New()
	Write() error
}

// AppType is used so we can choose which Action implementation to use
type AppType struct {
	App string `yaml:"app"` // which action is being called
}
