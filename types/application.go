package types

// Application defines a application that will have
// access control.
// All applications should have users, groups and roles
type Application struct {
	Base       `bson:"inline"`
	Name       string `json:"name" bson:"name"`
	ExternalID string `json:"external_id" bson:"external_id"`
}
