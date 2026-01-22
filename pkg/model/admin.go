package model

import "time"

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

type PartyDetails struct {
	Party              string
	IsLocal            bool
	LocalMetadata      map[string]string
	IdentityProviderID string
}

type ListKnownPartiesResponse struct {
	PartyDetails  []*PartyDetails
	NextPageToken string
}

type PruneRequest struct {
	PruneUpTo                 int64
	SubmissionID              string
	PruneAllDivulgedContracts bool
}

type PackageDetails struct {
	PackageID   string
	PackageSize uint64
	KnownSince  *time.Time
	Name        string
	Version     string
}

type CommandState int

const (
	CommandStateUnspecified CommandState = iota
	CommandStatePending
	CommandStateSucceeded
	CommandStateFailed
)

type CommandStatus struct {
	Started   *time.Time
	Completed *time.Time
	State     CommandState
}

type IdentityProviderConfig struct {
	IdentityProviderID string
	IsDeactivated      bool
	Issuer             string
	JwksURL            string
	Audience           string
}

type UpdateMask struct {
	Paths []string
}

type UpdateVettedPackagesRequest struct {
	Changes                        []*VettedPackagesChange
	DryRun                         bool
	SynchronizerID                 string
	ExpectedTopologySerial         *PriorTopologySerial
	UpdateVettedPackagesForceFlags []UpdateVettedPackagesForceFlag
}

type UpdateVettedPackagesResponse struct {
	PastVettedPackages *VettedPackages
	NewVettedPackages  *VettedPackages
}

type VettedPackagesChange struct {
	Vet   *VettedPackagesVet
	Unvet *VettedPackagesUnvet
}

type VettedPackagesVet struct {
	Packages               []*VettedPackagesRef
	NewValidFromInclusive  *time.Time
	NewValidUntilExclusive *time.Time
}

type VettedPackagesUnvet struct {
	Packages []*VettedPackagesRef
}

type VettedPackagesRef struct {
	PackageID      string
	PackageName    string
	PackageVersion string
}

type PriorTopologySerial struct {
	TopologySerial uint32
}

type UpdateVettedPackagesForceFlag int32

const (
	UpdateVettedPackagesForceFlagUnspecified                  UpdateVettedPackagesForceFlag = 0
	UpdateVettedPackagesForceFlagAllowVetIncompatibleUpgrades UpdateVettedPackagesForceFlag = 2
	UpdateVettedPackagesForceFlagAllowUnvettedDependencies    UpdateVettedPackagesForceFlag = 3
)
