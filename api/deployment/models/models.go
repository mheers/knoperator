package models

type JobCreateRequest struct {
	Name        string
	Image       string
	Command     []string
	Args        []string
	Env         map[string]string
	MountPoints map[string]string
	WorkingDir  string
}

type DeploymentCreateRequest struct {
	Name    string
	Image   string
	Command []string
	Args    []string
	Env     map[string]string
}
