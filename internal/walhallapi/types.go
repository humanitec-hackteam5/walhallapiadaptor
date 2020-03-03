package walhallapi

// Config represents a configuration in Walhall Core
type Config struct {
	ID              int                    `json:"id"`
	UUID            string                 `json:"configuration_uuid"`
	Name            string                 `json:"name"`
	Spec            map[string]interface{} `json:"specification"`
	Type            string                 `json:"type"`
	Status          string                 `json:"status"`
	EnvUUID         string                 `json:"environment"`
	ModuleVersionID int                    `json:"logic_module_version"`
	CreatedAt       string                 `json:"create_date"`
	EditedAt        string                 `json:"edit_date"`
}

// Environment represents an environment in Walhall Core
type Environment struct {
	ModuleVersions []struct {
		ModuleVersion
		Module Module `json:"logic_module"`
	} `json:"logic_module_versions"`
	UUID string `json:"env_uuid"`
	Name string `json:"name"`
}

// Module represents a logic model in Walhall Core
type Module struct {
	Name     string          `json:"name"`
	UUID     string          `json:"module_uuid"`
	Repo     string          `json:"github_repo"`
	Image    string          `json:"image"`
	Versions []ModuleVersion `json:"versions"`
}

// ModuleVersion represents a logic module version in Walhall Core
type ModuleVersion struct {
	ID      int    `json:"id"`
	UUID    string `json:"version_uuid"`
	Version string `json:"version"`
}

// Interface to allow easy mocking of walhallapi
type WalhallAPIer interface {
	GetCurrentUser() string
	ListOrgs() (map[string]string, error)
	ListApps(orgName string) (map[string]string, error)
	ListModules(orgName string) ([]Module, error)
	RefreshModules(orgName string) (string, error)
	GetRefreshModulesStatus(orgName string) (string, error)
	ListEnvs(orgName, appName string) ([]Environment, error)
	GetEnv(orgName, appName, envName string) (Environment, error)
	PatchEnv(env Environment, moduleVersions []int) (Environment, error)
	DeleteModuleVersionFromEnv(env Environment, mv ModuleVersion) (Environment, error)
	GetConfigsForModuleVersionInEnv(env Environment, mv ModuleVersion) ([]Config, error)
	UpdateConfiguration(config Config) (Config, error)
	CreateConfiguration(env Environment, mv ModuleVersion, configType string) (Config, error)
	DeleteConfiguration(configID int) error
	DeployToEnvironment(env Environment) error
}
