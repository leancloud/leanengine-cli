package api

import (
	"github.com/leancloud/lean-cli/lean/apps"
)

// GetAppListResult is GetAppList function's result type
type GetAppListResult struct {
	AppID     string `json:"app_id"`
	AppKey    string `json:"app_key"`
	AppName   string `json:"app_name"`
	MasterKey string `json:"master_key"`
	AppDomain string `json:"app_domain"`
}

// GetAppList returns the current user's all LeanCloud application
func GetAppList() ([]*GetAppListResult, error) {
	client := NewClient()

	resp, err := client.get("/1/clients/self/apps", nil)
	if err != nil {
		return nil, err
	}

	var result []*GetAppListResult
	err = resp.JSON(&result)
	return result, err
}

// DeployImage will deploy the engine group with specify image tag
func DeployImage(appID string, groupName string, imageTag string) (string, error) {
	client := NewClient()
	opts, err := client.options()
	if err != nil {
		return "", err
	}
	opts.Headers["X-LC-Id"] = appID

	resp, err := client.put("/1.1/functions/_ops/groups/"+groupName+"/deploy", map[string]interface{}{
		"imageTag": imageTag,
		"async":    true,
	}, opts)

	if err != nil {
		return "", err
	}
	result := new(struct {
		EventToken string `json:"eventToken"`
	})

	err = resp.JSON(result)
	return result.EventToken, err
}

// DeployAppFromGit will deploy applications with user's git repo
// returns the event token for polling deploy log
func DeployAppFromGit(projectPath string, groupName string) (string, error) {
	client := NewClient()

	appID, err := apps.GetCurrentAppID(projectPath)
	if err != nil {
		return "", err
	}

	opts, err := client.options()
	if err != nil {
		return "", err
	}
	opts.Headers["X-LC-Id"] = appID

	resp, err := client.post("/1.1/functions/_ops/groups/"+groupName+"/buildAndDeploy", map[string]interface{}{
		"comment":             "",
		"noDependenciesCache": false,
		"async":               true,
	}, opts)

	if err != nil {
		return "", err
	}

	result := new(struct {
		EventToken string `json:"eventToken"`
	})
	err = resp.JSON(result)
	return result.EventToken, err
}

// DeployAppFromFile will deploy applications with specific file
// returns the event token for polling deploy log
func DeployAppFromFile(projectPath string, groupName string, fileURL string) (string, error) {
	client := NewClient()

	appID, err := apps.GetCurrentAppID(projectPath)
	if err != nil {
		return "", err
	}

	opts, err := client.options()
	if err != nil {
		return "", err
	}
	opts.Headers["X-LC-Id"] = appID

	resp, err := client.post("/1.1/functions/_ops/groups/"+groupName+"/buildAndDeploy", map[string]interface{}{
		"zipUrl":              fileURL,
		"comment":             "",
		"noDependenciesCache": false,
		"async":               true,
	}, opts)

	if err != nil {
		return "", err
	}

	result := new(struct {
		EventToken string `json:"eventToken"`
	})
	err = resp.JSON(result)
	return result.EventToken, err
}

// GetAppInfoResult is GetAppInfo function's result type
type GetAppInfoResult struct {
	AppID          string `json:"app_id"`
	AppKey         string `json:"app_key"`
	AppName        string `json:"app_name"`
	MasterKey      string `json:"master_key"`
	AppDomain      string `json:"app_domain"`
	LeanEngineMode string `json:"leanengine_mode"`
}

// GetAppInfo returns the application's detail info
func GetAppInfo(appID string) (*GetAppInfoResult, error) {
	client := NewClient()
	resp, err := client.get("/1.1/clients/self/apps/"+appID, nil)
	if err != nil {
		return nil, err
	}
	result := new(GetAppInfoResult)
	err = resp.JSON(result)
	return result, err
}

// GetGroupsResult is GetGroups's result struct
type GetGroupsResult struct {
	GroupName string `json:"groupName"`
	Prod      int    `json:"prod"`
	Instances []struct {
		Name  string `json:"name"`
		Quota int    `json:"quota"`
	} `json:"instances"`
	CurrentImage struct {
		Runtime  string `json:"runtime"`
		ImageTag string `json:"imageTag"`
	} `json:"currentImage"`
}

// GetGroups returns the application's engine groups
func GetGroups(appID string) ([]*GetGroupsResult, error) {
	client := NewClient()
	opts, err := client.options()
	if err != nil {
		return nil, err
	}
	opts.Headers["X-LC-Id"] = appID

	resp, err := client.get("/1.1/functions/_ops/groups", opts)
	if err != nil {
		return nil, err
	}

	var result []*GetGroupsResult
	err = resp.JSON(&result)

	return result, err
}
