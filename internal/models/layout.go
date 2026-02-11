package models

type ComponentType string

type Component struct {
	ID        int           `json:"id"`
	Component ComponentType `json:"component"`
	Model     any           `json:"model"`
}

type Layout []Component
