package model

type User struct {
	ID                 string
	PrimaryParty       string
	IsDeactivated      bool
	Metadata           map[string]string
	IdentityProviderID string
}

type Right struct {
	Type RightType
}

type RightType interface {
	isRightType()
}

type CanActAs struct {
	Party string
}

func (CanActAs) isRightType() {}

type CanReadAs struct {
	Party string
}

func (CanReadAs) isRightType() {}

type ParticipantAdmin struct{}

func (ParticipantAdmin) isRightType() {}

type IdentityProviderAdmin struct{}

func (IdentityProviderAdmin) isRightType() {}
