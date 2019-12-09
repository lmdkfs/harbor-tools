package controllers

type Tags struct {
	Name   string  `json:"name"`
	Labels []Lable `json:"labels"`
}

type Lable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Repositories struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ProjectID int    `json:"project_id"`
	TagsCount int    `json:"tags_count"`
}

type Projects struct {
	ProjectID int                    `json:"project_id"`
	Name      string                 `json:"name"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type ProjectInfo struct {
	ProjectName string                 `json:"project_name"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type ProjectMemberEntity struct {
	EntityID   int    `json:"entity_id"`
	RoleName   string `json:"role_name"`
	EntityName string `json:"entity_name"`
	EntityType string `json:"entity_type"`
	ProjectID  int    `json:"project_id"`
	ID         int    `json:"id"`
	RoleID     int    `json:"role_id"`
}
