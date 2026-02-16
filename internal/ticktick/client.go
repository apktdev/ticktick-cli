package ticktick

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/adrian/ticktick-cli/internal/config"
)

const (
	authBase = "https://ticktick.com/oauth"
	apiBase  = "https://api.ticktick.com/open/v1"
)

type Client struct {
	http *http.Client
	cfg  *config.Config
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
}

type Project struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Color      string `json:"color,omitempty"`
	SortOrder  int64  `json:"sortOrder,omitempty"`
	ViewMode   string `json:"viewMode,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Closed     bool   `json:"closed,omitempty"`
	GroupID    string `json:"groupId,omitempty"`
	Permission string `json:"permission,omitempty"`
}

type ChecklistItem struct {
	ID            string `json:"id,omitempty"`
	Title         string `json:"title,omitempty"`
	Status        int    `json:"status,omitempty"`
	CompletedTime string `json:"completedTime,omitempty"`
	IsAllDay      bool   `json:"isAllDay,omitempty"`
	SortOrder     int64  `json:"sortOrder,omitempty"`
	StartDate     string `json:"startDate,omitempty"`
	TimeZone      string `json:"timeZone,omitempty"`
}

type Task struct {
	ID            string          `json:"id,omitempty"`
	ProjectID     string          `json:"projectId,omitempty"`
	Title         string          `json:"title,omitempty"`
	Content       string          `json:"content,omitempty"`
	Desc          string          `json:"desc,omitempty"`
	IsAllDay      bool            `json:"isAllDay,omitempty"`
	StartDate     string          `json:"startDate,omitempty"`
	DueDate       string          `json:"dueDate,omitempty"`
	TimeZone      string          `json:"timeZone,omitempty"`
	Reminders     []string        `json:"reminders,omitempty"`
	RepeatFlag    string          `json:"repeatFlag,omitempty"`
	Priority      int             `json:"priority,omitempty"`
	Status        int             `json:"status,omitempty"`
	CompletedTime string          `json:"completedTime,omitempty"`
	SortOrder     int64           `json:"sortOrder,omitempty"`
	Items         []ChecklistItem `json:"items,omitempty"`
	Kind          string          `json:"kind,omitempty"`
}

type ProjectColumn struct {
	ID        string `json:"id,omitempty"`
	ProjectID string `json:"projectId,omitempty"`
	Name      string `json:"name,omitempty"`
	SortOrder int64  `json:"sortOrder,omitempty"`
}

type ProjectData struct {
	Project Project         `json:"project"`
	Tasks   []Task          `json:"tasks"`
	Columns []ProjectColumn `json:"columns,omitempty"`
}

func New(cfg *config.Config) *Client {
	return &Client{
		http: &http.Client{Timeout: 20 * time.Second},
		cfg:  cfg,
	}
}

func AuthURL(cfg *config.Config, scope, state string) (string, error) {
	if cfg.ClientID == "" {
		return "", fmt.Errorf("client_id is missing; run `tick auth set-client`")
	}
	if cfg.RedirectURI == "" {
		return "", fmt.Errorf("redirect_uri is missing; run `tick auth set-client`")
	}

	v := url.Values{}
	v.Set("client_id", cfg.ClientID)
	v.Set("scope", scope)
	v.Set("state", state)
	v.Set("redirect_uri", cfg.RedirectURI)
	v.Set("response_type", "code")

	return authBase + "/authorize?" + v.Encode(), nil
}

func (c *Client) ExchangeCode(ctx context.Context, code string) error {
	if c.cfg.ClientID == "" || c.cfg.ClientSecret == "" || c.cfg.RedirectURI == "" {
		return fmt.Errorf("client configuration incomplete; run `tick auth set-client`")
	}

	v := url.Values{}
	v.Set("code", code)
	v.Set("grant_type", "authorization_code")
	v.Set("scope", "tasks:read tasks:write")
	v.Set("redirect_uri", c.cfg.RedirectURI)

	tr, err := c.tokenRequest(ctx, v)
	if err != nil {
		return err
	}
	c.applyToken(tr)
	return nil
}

func (c *Client) Refresh(ctx context.Context) error {
	if c.cfg.RefreshToken == "" {
		return fmt.Errorf("no refresh token in config")
	}

	v := url.Values{}
	v.Set("grant_type", "refresh_token")
	v.Set("refresh_token", c.cfg.RefreshToken)

	tr, err := c.tokenRequest(ctx, v)
	if err != nil {
		return err
	}
	c.applyToken(tr)
	return nil
}

func (c *Client) tokenRequest(ctx context.Context, values url.Values) (*tokenResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, authBase+"/token", strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	auth := base64.StdEncoding.EncodeToString([]byte(c.cfg.ClientID + ":" + c.cfg.ClientSecret))
	req.Header.Set("Authorization", "Basic "+auth)

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, _ := io.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("token request failed (%d): %s", res.StatusCode, strings.TrimSpace(string(b)))
	}

	var tr tokenResponse
	if err := json.Unmarshal(b, &tr); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}
	return &tr, nil
}

func (c *Client) applyToken(tr *tokenResponse) {
	c.cfg.AccessToken = tr.AccessToken
	if tr.RefreshToken != "" {
		c.cfg.RefreshToken = tr.RefreshToken
	}
	c.cfg.TokenType = tr.TokenType
	c.cfg.Scope = tr.Scope
	if tr.ExpiresIn > 0 {
		c.cfg.Expiry = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	}
}

func (c *Client) ensureValidToken(ctx context.Context) error {
	if c.cfg.AccessToken == "" {
		return fmt.Errorf("no access token; run `tick auth login-url` and `tick auth exchange`")
	}
	if !c.cfg.Expiry.IsZero() && time.Now().After(c.cfg.Expiry.Add(-30*time.Second)) && c.cfg.RefreshToken != "" {
		if err := c.Refresh(ctx); err != nil {
			return fmt.Errorf("refresh token: %w", err)
		}
	}
	return nil
}

func (c *Client) doJSON(ctx context.Context, method, path string, in, out any) error {
	if err := c.ensureValidToken(ctx); err != nil {
		return err
	}

	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, apiBase+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.cfg.AccessToken)
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	b, _ := io.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("api request failed (%d): %s", res.StatusCode, strings.TrimSpace(string(b)))
	}
	if out == nil || len(b) == 0 {
		return nil
	}
	if err := json.Unmarshal(b, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

func (c *Client) ListProjects(ctx context.Context) ([]Project, error) {
	var projects []Project
	if err := c.doJSON(ctx, http.MethodGet, "/project", nil, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func (c *Client) GetProject(ctx context.Context, projectID string) (*Project, error) {
	var project Project
	if err := c.doJSON(ctx, http.MethodGet, "/project/"+url.PathEscape(projectID), nil, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

func (c *Client) ProjectData(ctx context.Context, projectID string) (*ProjectData, error) {
	var data ProjectData
	if err := c.doJSON(ctx, http.MethodGet, "/project/"+url.PathEscape(projectID)+"/data", nil, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) CreateProject(ctx context.Context, project Project) (*Project, error) {
	var created Project
	if err := c.doJSON(ctx, http.MethodPost, "/project", project, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

func (c *Client) UpdateProject(ctx context.Context, projectID string, project Project) (*Project, error) {
	var updated Project
	if err := c.doJSON(ctx, http.MethodPost, "/project/"+url.PathEscape(projectID), project, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (c *Client) DeleteProject(ctx context.Context, projectID string) error {
	return c.doJSON(ctx, http.MethodDelete, "/project/"+url.PathEscape(projectID), nil, nil)
}

func (c *Client) GetTask(ctx context.Context, projectID, taskID string) (*Task, error) {
	var task Task
	if err := c.doJSON(ctx, http.MethodGet, "/project/"+url.PathEscape(projectID)+"/task/"+url.PathEscape(taskID), nil, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (c *Client) CreateTask(ctx context.Context, task Task) (*Task, error) {
	var created Task
	if err := c.doJSON(ctx, http.MethodPost, "/task", task, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

func (c *Client) UpdateTask(ctx context.Context, taskID string, task Task) (*Task, error) {
	var updated Task
	if err := c.doJSON(ctx, http.MethodPost, "/task/"+url.PathEscape(taskID), task, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (c *Client) CompleteTask(ctx context.Context, projectID, taskID string) error {
	return c.doJSON(ctx, http.MethodPost, "/project/"+url.PathEscape(projectID)+"/task/"+url.PathEscape(taskID)+"/complete", nil, nil)
}

func (c *Client) DeleteTask(ctx context.Context, projectID, taskID string) error {
	return c.doJSON(ctx, http.MethodDelete, "/project/"+url.PathEscape(projectID)+"/task/"+url.PathEscape(taskID), nil, nil)
}
