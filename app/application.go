package app

type GlobalFlags struct {
	Force       bool
	Tag         string
	HttpProxy   string
	GithubProxy string
}

// Define an interface for managing applications
type Manager interface {
	GetName() string
	Install(flags *GlobalFlags) error
	Update(flags *GlobalFlags) error
	Delete(flags *GlobalFlags) error
	IsInstalled() bool
}
