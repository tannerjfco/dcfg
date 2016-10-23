package drudconfig

type Action interface {
	GetPayload() string
	Pretty()
}

var TypeMap = map[string]Action{
	"command": &Command{},
	"write":   &Write{},
}

type Task struct {
	Name    string `yaml:"name"`    // name of the task
	Dest    string `yaml:"dest"`    // what this action will be performed on
	Workdir string `yaml:"workdir"` // where this action will be called from
	Wait    string `yaml:"wait"`    // how long to wait before this action is called
	Repeat  int    `yaml:"repeat"`  // how many times to run this action
	Ignore  bool   `yaml:"ignore"`  // ignore failures or not
}

type TaskType struct {
	Action string `yaml:"action"` // which action is being called
}
