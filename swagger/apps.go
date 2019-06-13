package swagger

import (
	cce "github.com/smartedgemec/controller-ce"
)

// AppSummary is a summary representation of the app.
type AppSummary struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Vendor      string `json:"vendor"`
	Description string `json:"description"`
}

// AppDetail is a detailed representation of the app.
type AppDetail struct {
	AppSummary
	Cores  int         `json:"cores"`
	Memory int         `json:"memory"`
	Ports  []PortProto `json:"ports"`
	Source string      `json:"source"`
}

// PortProto is a port and protocol combination.
type PortProto struct {
	cce.PortProto
}

// AppList is a list representation of apps.
type AppList struct {
	Apps []AppSummary `json:"apps"`
}
