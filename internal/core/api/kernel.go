// Generated file. DO NOT EDIT.

package api

import (
	"context"
)

type Kernel interface {
	AuthEsv(context.Context, *AuthEsvInput) (*AuthEsvOutput, error)
	SaveEsvRefreshToken(context.Context, *SaveEsvRefreshTokenInput) (*SaveEsvRefreshTokenOutput, error)
	UnauthEsv(context.Context, *UnauthEsvInput) (*UnauthEsvOutput, error)
	GetEsvUser(context.Context, *GetEsvUserInput) (*GetEsvUserOutput, error)
	CreateProject(context.Context, *CreateProjectInput) (*CreateProjectOutput, error)
	DescribeTemplates(context.Context, *DescribeTemplatesInput) (*DescribeTemplatesOutput, error)
	CreateWorkspace(context.Context, *CreateWorkspaceInput) (*CreateWorkspaceOutput, error)
	DescribeWorkspaces(context.Context, *DescribeWorkspacesInput) (*DescribeWorkspacesOutput, error)
	ResolveWorkspace(context.Context, *ResolveWorkspaceInput) (*ResolveWorkspaceOutput, error)
	// Debug method to test what happens when the service panics.
	Panic(context.Context, *PanicInput) (*PanicOutput, error)
	// Retrieves the installed and current version of exo.
	GetVersion(context.Context, *GetVersionInput) (*GetVersionOutput, error)
	// Upgrades exo to the latest version.
	Upgrade(context.Context, *UpgradeInput) (*UpgradeOutput, error)
	// Checks whether server is up.
	Ping(context.Context, *PingInput) (*PingOutput, error)
	// Gracefully shutdown the exo daemon.
	Exit(context.Context, *ExitInput) (*ExitOutput, error)
	DescribeTasks(context.Context, *DescribeTasksInput) (*DescribeTasksOutput, error)
	GetUserHomeDir(context.Context, *GetUserHomeDirInput) (*GetUserHomeDirOutput, error)
	ReadDir(context.Context, *ReadDirInput) (*ReadDirOutput, error)
}

type AuthEsvInput struct {
}

type AuthEsvOutput struct {
	AuthURL  string `json:"authUrl"`
	AuthCode string `json:"authCode"`
}

type SaveEsvRefreshTokenInput struct {
	RefreshToken string `json:"refreshToken"`
}

type SaveEsvRefreshTokenOutput struct {
}

type UnauthEsvInput struct {
}

type UnauthEsvOutput struct {
}

type GetEsvUserInput struct {
	VaultURL string `json:"vaultUrl"`
}

type GetEsvUserOutput struct {
	Email string `json:"email"`
}

type CreateProjectInput struct {
	Root        string  `json:"root"`
	TemplateURL *string `json:"templateUrl"`
}

type CreateProjectOutput struct {
	WorkspaceID string `json:"workspaceId"`
}

type DescribeTemplatesInput struct {
}

type DescribeTemplatesOutput struct {
	Templates []TemplateDescription `json:"templates"`
}

type CreateWorkspaceInput struct {
	Root string `json:"root"`
}

type CreateWorkspaceOutput struct {
	ID string `json:"id"`
}

type DescribeWorkspacesInput struct {
}

type DescribeWorkspacesOutput struct {
	Workspaces []WorkspaceDescription `json:"workspaces"`
}

type ResolveWorkspaceInput struct {
	Ref string `json:"ref"`
}

type ResolveWorkspaceOutput struct {
	ID *string `json:"id"`
}

type PanicInput struct {
	Message string `json:"message"`
}

type PanicOutput struct {
}

type GetVersionInput struct {
}

type GetVersionOutput struct {
	Installed string  `json:"installed"`
	Latest    *string `json:"latest"`
	Current   bool    `json:"current"`
	Managed   bool    `json:"managed"`
}

type UpgradeInput struct {
}

type UpgradeOutput struct {
}

type PingInput struct {
}

type PingOutput struct {
}

type ExitInput struct {
}

type ExitOutput struct {
}

type DescribeTasksInput struct {

	// If supplied, filters tasks by job.
	JobIDs []string `json:"jobIds"`
}

type DescribeTasksOutput struct {
	Tasks []TaskDescription `json:"tasks"`
}

type GetUserHomeDirInput struct {
}

type GetUserHomeDirOutput struct {
	Path string `json:"path"`
}

type ReadDirInput struct {
	Path string `json:"path"`
}

type ReadDirOutput struct {
	Directory DirectoryEntry   `json:"directory"`
	Parent    *DirectoryEntry  `json:"parent"`
	Entries   []DirectoryEntry `json:"entries"`
}


type DirectoryEntry struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	IsDirectory bool   `json:"isDirectory"`
}

type TemplateDescription struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	IconGlyph   string `json:"iconGlyph"`
	URL         string `json:"url"`
}

type TaskDescription struct {
	ID string `json:"id"`
	// ID of root task in this tree.
	JobID    string  `json:"jobId"`
	ParentID *string `json:"parentId"`
	Name     string  `json:"name"`
	Status   string  `json:"status"`
	// Most recent log message. Single-line of text.
	Message  string        `json:"message"`
	Created  string        `json:"created"`
	Updated  string        `json:"updated"`
	Started  *string       `json:"started"`
	Finished *string       `json:"finished"`
	Progress *TaskProgress `json:"progress"`
}

type TaskProgress struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}
