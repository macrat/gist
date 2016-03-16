package gist

type FileInfo struct {
	FileName string `json:"filename"`
	Language string `json:"language"`
	RawURL   string `json:"raw_url"`
	Size     int    `json:"size"`
	Type     string `json:"type"`
}

type User struct {
	AvatarURL         string `json:"avatar_url"`
	EventsURL         string `json:"events_url"`
	FollowersURL      string `json:"fllowers_url"`
	FollowingURL      string `json:"fllowing_url"`
	GistsURL          string `json:"gists_url"`
	GravatarID        string `json:"gravatar_id"`
	HTMLURL           string `json:"html_URL"`
	ID                int    `json:"id"`
	Login             string `json:"login"`
	OrganizationsURL  string `json:"organizations_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	ReposURL          string `json:"repos_url"`
	SiteAdmin         bool   `json:"site_admin"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	Type              string `json:"type"`
	URL               string `json:"url"`
}

type Overview struct {
	Comments    int                 `json:"comments"`
	CommentsURL string              `json:"comments_url"`
	CommitsURL  string              `json:"commits_url"`
	CreatedAt   string              `json:"created_at"`
	Description string              `json:"description"`
	Files       map[string]FileInfo `json:"files"`
	ForksURL    string              `json:"forks_url"`
	GitPullURL  string              `json:"git_pull_url"`
	GitPushURL  string              `json:"git_push_url"`
	HTMLURL     string              `json:"html_url"`
	ID          string              `json:"id"`
	Owner       User                `json:"owner"`
	Public      bool                `json:"public"`
	Truncated   bool                `json:"truncated"`
	URL         string              `json:"url"`
	UpdatedAt   string              `json:"updated_at"`
	User        User                `json:"user"`
}

type File struct {
	Content   string `json:"content"`
	Language  string `json:"language"`
	RawURL    string `json:"raw_url"`
	Size      int    `json:"size"`
	Truncated bool   `json:"truncated"`
	Type      string `json:"type"`
}

type Fork struct {
	CreatedAt string `json:"created_at"`
	ID        string `json:"id"`
	URL       string `json:"url"`
	UpdatedAt string `json:"updated_at"`
	User      User   `json:"user"`
}

type ChangeStatus struct {
	Additions int `json:"additions"`
	Deletions int `json:"deletions"`
	Total     int `json:"total"`
}

type History struct {
	ChangeStatus ChangeStatus `json:"change_status"`
	CommittedAt  string       `json:"committed_at"`
	URL          string       `json:"url"`
	User         User         `json:"user"`
	version      string       `json:"version"`
}

type Gist struct {
	Comments    int             `json:"comments"`
	CommentsURL string          `json:"comments_url"`
	CommitsURL  string          `json:"commits_url"`
	CreatedAt   string          `json:"created_at"`
	Description string          `json:"description"`
	Files       map[string]File `json:"files"`
	Forks       []Fork          `json:"forks"`
	ForksURL    string          `json:"forks_url"`
	GitPullURL  string          `json:"git_pull_url"`
	GitPushURL  string          `json:"git_push_url"`
	HTMLURL     string          `json:"html_url"`
	History     []History       `json:"history"`
	ID          string          `json:"id"`
	Owner       User            `json:"owner"`
	Public      bool            `json:"public"`
	Truncated   bool            `json:"truncated"`
	URL         string          `json:"url"`
	UpdatedAt   string          `json:"updated_at"`
	User        User            `json:"user"`
}

type NewGistFile struct {
	Content string `json:"content"`
}

type NewGist struct {
	Description string                 `json:"description"`
	Files       map[string]NewGistFile `json:"files"`
	Public      bool                   `json:"public"`
}
