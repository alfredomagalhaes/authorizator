package types

// Application defines a application that will have
// access control.
// All applications should have users, groups and roles
type Application struct {
	Base
	Name       string `json:"name"`
	ExternalID string `json:"external_id" gorm:"unique"`
}
