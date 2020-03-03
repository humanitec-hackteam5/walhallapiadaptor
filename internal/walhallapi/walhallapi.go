package walhallapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (a *APIState) GetCurrentUser() string {
	return a.claims.Username
}

// ListOrgs returns a map from org name to UUIDs for the user's orgs - excluding the self org.
func (a *APIState) ListOrgs() (map[string]string, error) {
	cacheKey := "ListOrgs()"
	cachedResult, ok := a.cache[cacheKey]
	if ok {
		result := cachedResult.(map[string]string)
		return result, nil
	}
	url := "/api/walhalluser/" + a.claims.UserUUID
	resp, err := a.makeRequest(http.MethodGet, url, nil)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("list orgs: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, NewHTTPError(resp.StatusCode, a.claims.Username)
	}
	decoder := json.NewDecoder(resp.Body)
	var userDetail struct {
		Orgs []struct {
			UUID string `json:"organization_uuid"`
			Name string `json:"name"`
		} `json:"organizations"`
	}
	err = decoder.Decode(&userDetail)
	if err != nil {
		return nil, fmt.Errorf("Response from %s: %v ", url, err)
	}

	orgs := make(map[string]string)
	for _, org := range userDetail.Orgs {
		if org.Name != a.claims.Username {
			orgs[org.Name] = org.UUID
		}
	}
	a.cache[cacheKey] = orgs
	return orgs, nil
}

func (a *APIState) ListApps(orgName string) (map[string]string, error) {
	cacheKey := fmt.Sprintf(`ListApps("%s")`, orgName)
	cachedResult, ok := a.cache[cacheKey]
	if ok {
		result := cachedResult.(map[string]string)
		return result, nil
	}
	orgs, err := a.ListOrgs()
	if err != nil {
		return nil, fmt.Errorf("list apps: %v", err)
	}
	orgUUID, ok := orgs[orgName]
	if !ok {
		return nil, ErrNotFound
	}

	url := "/api/application?limit=100&organization_uuid=" + orgUUID
	resp, err := a.makeRequest(http.MethodGet, url, nil)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("list apps: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response from %s: Status %d ", url, resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var appDetail struct {
		Results []struct {
			UUID string `json:"app_uuid"`
			Name string `json:"name"`
		} `json:"results"`
	}
	err = decoder.Decode(&appDetail)
	if err != nil {
		return nil, fmt.Errorf("Response from %s: %v ", url, err)
	}

	apps := make(map[string]string)
	for _, app := range appDetail.Results {

		apps[app.Name] = app.UUID
	}
	a.cache[cacheKey] = apps
	return apps, nil
}

func (a *APIState) ListModules(orgName string) ([]Module, error) {
	orgs, err := a.ListOrgs()
	if err != nil {
		return nil, fmt.Errorf("list modules: %v", err)
	}
	orgUUID, ok := orgs[orgName]
	if !ok {
		return nil, ErrNotFound
	}

	url := fmt.Sprintf("/api/logicmodule?organization=%s&limit=50&status=internal", orgUUID)
	resp, err := a.makeRequest(http.MethodGet, url, nil)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("list modules: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response from %s: Status %d ", url, resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var moduleResponse struct {
		Results []Module `json:"results"`
	}
	err = decoder.Decode(&moduleResponse)
	if err != nil {
		return nil, fmt.Errorf("Response from %s: %v ", url, err)
	}
	return moduleResponse.Results, nil
}

func (a *APIState) RefreshModules(orgName string) (string, error) {
	orgs, err := a.ListOrgs()
	if err != nil {
		return "", fmt.Errorf("list modules: %v", err)
	}
	orgUUID, ok := orgs[orgName]
	if !ok {
		return "", ErrNotFound
	}

	url := fmt.Sprintf("/api/repositories/github/sync?organization_uuid=%s", orgUUID)
	resp, err := a.makeRequest(http.MethodPost, url, nil)
	defer resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("list modules: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Response from %s: Status %d ", url, resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var syncStatus struct {
		Status string `json:"status"`
	}
	err = decoder.Decode(&syncStatus)
	if err != nil {
		return "", fmt.Errorf("Response from %s: %v ", url, err)
	}
	return syncStatus.Status, nil
}
func (a *APIState) GetRefreshModulesStatus(orgName string) (string, error) {
	orgs, err := a.ListOrgs()
	if err != nil {
		return "", fmt.Errorf("refresh modules: %v", err)
	}
	orgUUID, ok := orgs[orgName]
	if !ok {
		return "", ErrNotFound
	}

	url := fmt.Sprintf("/api/repositories/github/status?organization_uuid=%s", orgUUID)
	resp, err := a.makeRequest(http.MethodGet, url, nil)
	defer resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("get refresh modules status: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Response from %s: Status %d ", url, resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var syncStatus struct {
		Status string `json:"status"`
	}
	err = decoder.Decode(&syncStatus)
	if err != nil {
		return "", fmt.Errorf("Response from %s: %v ", url, err)
	}
	return syncStatus.Status, nil
}

func (a *APIState) ListEnvs(orgName, appName string) ([]Environment, error) {
	cacheKey := fmt.Sprintf(`ListEnvs("%s","%s")`, orgName, appName)
	cachedResult, ok := a.cache[cacheKey]
	if ok {
		result := cachedResult.([]Environment)
		return result, nil
	}
	apps, err := a.ListApps(orgName)
	if err != nil {
		return nil, fmt.Errorf("get environment as deployment set: %v", err)
	}
	appUUID, ok := apps[appName]
	if !ok {
		return nil, ErrNotFound
	}

	url := "/api/environments?application=" + appUUID
	resp, err := a.makeRequest(http.MethodGet, url, nil)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("list apps: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response from %s: Status %d ", url, resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var envDetails struct {
		Results []Environment `json:"results"`
	}
	err = decoder.Decode(&envDetails)
	if err != nil {
		return nil, fmt.Errorf("Response from %s: %v ", url, err)
	}
	a.cache[cacheKey] = envDetails.Results
	return envDetails.Results, nil
}

func (a *APIState) GetEnv(orgName, appName, envName string) (Environment, error) {
	envDetails, err := a.ListEnvs(orgName, appName)
	if err != nil {
		return Environment{}, err
	}
	for _, envDetail := range envDetails {
		if envDetail.Name == envName {
			return envDetail, nil
		}
	}
	return Environment{}, ErrNotFound
}

func (a *APIState) PatchEnv(env Environment, moduleVersions []int) (Environment, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	var moduleVersionsWrapper struct {
		LogicModuleVersionIds []int `json:"logic_module_version_ids"`
	}
	moduleVersionsWrapper.LogicModuleVersionIds = moduleVersions
	encoder.Encode(moduleVersionsWrapper)

	resp, err := a.makeRequest(http.MethodPatch, "/api/environments/"+env.UUID, &buffer)
	defer resp.Body.Close()
	if err != nil {
		return Environment{}, fmt.Errorf("patch environment: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return Environment{}, fmt.Errorf("patch environment: expected 200, got %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var envDetail Environment
	decoder.Decode(&envDetail)
	return envDetail, nil
}

func (a *APIState) DeleteModuleVersionFromEnv(env Environment, mv ModuleVersion) (Environment, error) {
	deleteMVURL := fmt.Sprintf("/api/environments/%s/remove/%s", env.UUID, mv.UUID)
	resp, err := a.makeRequest(http.MethodDelete, deleteMVURL, nil)
	defer resp.Body.Close()
	if err != nil {
		return Environment{}, fmt.Errorf("delete module version from env: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return Environment{}, fmt.Errorf("delete module version from env: expected 200, got %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var envDetail Environment
	decoder.Decode(&envDetail)
	return envDetail, nil
}

func (a *APIState) GetConfigsForModuleVersionInEnv(env Environment, mv ModuleVersion) ([]Config, error) {
	getConfigsURL := fmt.Sprintf("/api/configuration?logic_module_version=%d&environment=%s", mv.ID, env.UUID)
	resp, err := a.makeRequest(http.MethodGet, getConfigsURL, nil)
	defer resp.Body.Close()
	if err != nil {
		return []Config{}, fmt.Errorf("get config for module in env: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return []Config{}, fmt.Errorf("get config for module in env: expected 200, got %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var configResults struct {
		Results []Config `json:"results"`
	}
	err = decoder.Decode(&configResults)
	if err != nil {
		return []Config{}, fmt.Errorf("get config for module in env: %v", err)
	}
	return configResults.Results, nil
}

// UpdateConfiguration updates the configuration to match the supplied config
// PUT /api/configuration/<config.ID>
// Returns the supplied config back
func (a *APIState) UpdateConfiguration(config Config) (Config, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.Encode(config)
	resp, err := a.makeRequest(http.MethodPut, "/api/configuration/"+strconv.Itoa(config.ID), &buffer)
	defer resp.Body.Close()
	if err != nil {
		return Config{}, fmt.Errorf("put config: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return Config{}, fmt.Errorf("put config: expected 200, got %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var configResult Config
	err = decoder.Decode(&configResult)
	if err != nil {
		return Config{}, fmt.Errorf("put config: %v", err)
	}
	return configResult, nil
}

// CreateConfiguration creates a new a configuration of a given type in an environment
// POST /api/configuration
// Body:
// {
//   type: <configType>,
//   logic_module_version: <mv.ID>,
//   environment: <env.UUID>
// }
// Returns the created configuration (useful for a subsequent call to UpdateConfiguration
func (a *APIState) CreateConfiguration(env Environment, mv ModuleVersion, configType string) (Config, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	configDef := struct {
		Type            string `json:"type"`
		ModuleVersionID int    `json:"logic_module_version"`
		EnvUUID         string `json:"environment"`
	}{Type: configType, ModuleVersionID: mv.ID, EnvUUID: env.UUID}
	encoder.Encode(configDef)
	resp, err := a.makeRequest(http.MethodPost, "/api/configuration", &buffer)
	defer resp.Body.Close()
	if err != nil {
		return Config{}, fmt.Errorf("post config for module in env: %v", err)
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return Config{}, fmt.Errorf("post config for module in env: expected 201, got %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var configResult Config
	err = decoder.Decode(&configResult)
	if err != nil {
		return Config{}, fmt.Errorf("put config for module in env: %v", err)
	}
	return configResult, nil
}

// DeleteConfiguration deletes a configuration given the supplied Configuration ID
// DELETE /api/configuration/<ConfigID>
// Nil error indicates success
func (a *APIState) DeleteConfiguration(configID int) error {
	deleteConfigURL := fmt.Sprintf("/api/configuration/%d", configID)
	resp, err := a.makeRequest(http.MethodDelete, deleteConfigURL, nil)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("post config for module in env: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("post config for module in env: expected 200, got %d", resp.StatusCode)
	}
	return nil
}

// DeployToEnvironment deployes the current state of the configurations to the
// cluster defined in the environment
func (a *APIState) DeployToEnvironment(env Environment) error {
	putDeployURL := fmt.Sprintf("/api/environments/%s/deploy", env.UUID)
	resp, err := a.makeRequest(http.MethodPut, putDeployURL, bytes.NewBuffer([]byte("{}")))
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("deploy to env %s: %w", env.UUID, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("post config for module in env: expected 200, got %d", resp.StatusCode)
	}
	return nil
}
