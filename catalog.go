package osbapi

// TODO: Metadata and Schemas
type Catalog struct {
	Services []*Service `json:"services"`
}

type Service struct {
	Name                 string             `json:"name"`
	Description          string             `json:"description"`
	ID                   string             `json:"id"`
	Tags                 []string           `json:"tags"`
	Bindable             bool               `json:"bindable"`
	PlanUpdateable       bool               `json:"plan_updateable"`
	BindingsRetrievable  bool               `json:"bindings_retrievable"`
	InstancesRetrievable bool               `json:"instances_retrievable"`
	Plans                []*Plan            `json:"plans"`
	Requires             []string           `json:"requires"`
	DashboardClient      []*DashboardClient `json:"dashboard_client"`
}

type DashboardClient struct {
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	RedirectURI string `json:"redirect_uri"`
}

type Plan struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ID          string `json:"id"`
	Free        bool   `json:"free"`
	Bindable    bool   `json:"bindable"`
}
