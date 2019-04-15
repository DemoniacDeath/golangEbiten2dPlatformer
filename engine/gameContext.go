package engine

type Context struct {
	Settings 	*Settings
	Quit     	bool
}

func NewContext(settings *Settings) *Context {
	return &Context{Settings: settings}
}