package envoyerapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ProjectsResource struct {
	Client
}

type Project struct {
	ID                       int              `json:"id"`
	UserId                   int              `json:"user_id"`
	Version                  int              `json:"version"`
	Name                     string           `json:"name"`
	Provider                 string           `json:"provider"`
	PlainRepository          string           `json:"plain_repository"`
	Repository               string           `json:"repository"`
	Type                     string           `json:"type"`
	Branch                   string           `json:"branch"`
	PushToDeploy             bool             `json:"push_to_deploy"`
	WebhookId                *int             `json:"webhook_id"`
	Status                   *string          `json:"status"`
	ShouldDeployAgain        int              `json:"should_deploy_again"`
	DeployAgainTarget        *string          `json:"deploy_again_target"`
	DeployAgainTargetType    *string          `json:"deploy_again_target_type"`
	DeploymentStartedAt      *string          `json:"deployment_started_at"`
	DeploymentFinishedAt     *string          `json:"deployment_finished_at"`
	LastDeploymentStatus     *string          `json:"last_deployment_status"`
	DailyDeploys             int              `json:"daily_deploys"`
	WeeklyDeploys            int              `json:"weekly_deploys"`
	LastDeploymentTook       int              `json:"last_deployment_took"`
	RetainDeployments        int              `json:"retain_deployments"`
	HighVolume               bool             `json:"high_volume"`
	EnvironmentServers       []int            `json:"environment_servers"`
	Folders                  []string         `json:"folders"`
	Monitor                  *string          `json:"monitor"`
	NewYorkStatus            *string          `json:"new_york_status"`
	LondonStatus             *string          `json:"london_status"`
	SingaporeStatus          *string          `json:"singapore_status"`
	Token                    string           `json:"token"`
	CreatedAt                string           `json:"created_at"`
	UpdatedAt                string           `json:"updated_at"`
	InstallDevDependencies   bool             `json:"install_dev_dependencies"`
	InstallDependencies      bool             `json:"install_dependencies"`
	QuietComposer            bool             `json:"quiet_composer"`
	Servers                  []map[string]any `json:"servers"`
	Collaborators            []map[string]any `json:"collaborators"`
	User                     map[string]any   `json:"user"`
	HasEnvironment           bool             `json:"has_environment"`
	HasMonitoringError       bool             `json:"has_monitoring_error"`
	HasMissingHeartbeats     bool             `json:"has_missing_heartbeats"`
	LastDeployedBranch       *string          `json:"last_deployed_branch"`
	LastDeploymentId         *int             `json:"last_deployment_id"`
	LastDeploymentAuthor     *string          `json:"last_deployment_author"`
	LastDeploymentAvatar     *string          `json:"last_deployment_avatar"`
	LastDeploymentHash       *string          `json:"last_deployment_hash"`
	LastDeploymentBranch     *string          `json:"last_deployment_branch"`
	LastDeploymentTimestamp  *string          `json:"last_deployment_timestamp"`
	CanBeDeployed            bool             `json:"can_be_deployed"`
	DisplayableRepositoryUrl string           `json:"displayable_repository_url"`
}

type ProjectResponse struct {
	Project Project `json:"project"`
}

func (p *ProjectsResource) List() string {
	return "List of projects"
}

func (p *ProjectsResource) Get(projectId string) (*Project, error) {
	req, err := p.Client.NewRequest(http.MethodGet, "projects/"+projectId, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		p.Client.logger.Debug("Failed to send request", slog.String("method", req.Method), slog.String("url", req.URL.String()), slog.Any("req_headers", req.Header), slog.Any("body", req.Body), slog.Any("error", err))

		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		httpError := NewHTTPError(req, resp, &err)

		p.Client.logger.Debug("error response", slog.Int("status", httpError.Status), slog.String("body", httpError.BodySnippet), slog.Any("resp_headers", httpError.Headers))

		return nil, httpError
	}

	var projectResponse ProjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&projectResponse); err != nil {
		return nil, err
	}

	return &projectResponse.Project, nil
}

func (p *ProjectsResource) LinkedFolders(projectId int) *LinkedFolders {
	return &LinkedFolders{
		Client: p.Client,
	}
}
