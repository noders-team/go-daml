package model

import (
	"fmt"
	"time"
)

type Commands struct {
	WorkflowID          string
	UserID              string
	CommandID           string
	Commands            []*Command
	DeduplicationPeriod DeduplicationPeriod
	MinLedgerTimeAbs    *time.Time
	MinLedgerTimeRel    *time.Duration
	ActAs               []string
	ReadAs              []string
	SubmissionID        string
}

type DeduplicationPeriod interface {
	isDeduplicationPeriod()
}

type DeduplicationDuration struct {
	Duration time.Duration
}

func (DeduplicationDuration) isDeduplicationPeriod() {}

type DeduplicationOffset struct {
	Offset int64
}

func (DeduplicationOffset) isDeduplicationPeriod() {}

type Command struct {
	Command CommandType
}

type CommandType interface {
	isCommandType()
}

type CreateCommand struct {
	TemplateID string
	Arguments  map[string]interface{}
}

func (CreateCommand) isCommandType() {}

type ExerciseCommand struct {
	ContractID string
	TemplateID string
	Choice     string
	Arguments  map[string]interface{}
}

func (ExerciseCommand) isCommandType() {}

type ExerciseByKeyCommand struct {
	TemplateID string
	Key        map[string]interface{}
	Choice     string
	Arguments  map[string]interface{}
}

func (ExerciseByKeyCommand) isCommandType() {}

type CompletionStreamRequest struct {
	UserID         string
	Parties        []string
	BeginExclusive int64
}

type CompletionStreamResponse struct {
	Response CompletionResponse
}

type CompletionResponse interface {
	isCompletionResponse()
}

type Completion struct {
	CommandID     string
	Status        Status
	UpdateID      string
	TransactionID string
	SubmissionID  string
	CompletedAt   *time.Time
	Offset        int64
}

func (Completion) isCompletionResponse() {}

type OffsetCheckpoint struct {
	Offset int64
}

func (OffsetCheckpoint) isCompletionResponse() {}

type Status interface {
	isStatus()
}

type StatusOK struct{}

func (StatusOK) isStatus() {}

type StatusError struct {
	Code    int32
	Message string
}

func (StatusError) isStatus() {}

type SubmitRequest struct {
	Commands *Commands
}

type SubmitResponse struct{}

type SubmitAndWaitRequest struct {
	Commands *Commands
}

type SubmitAndWaitResponse struct {
	UpdateID         string
	CompletionOffset int64
}

// Event Query Service types
type GetEventsByContractIDRequest struct {
	ContractID        string
	RequestingParties []string
}

type GetEventsByContractIDResponse struct {
	CreateEvent  *CreatedEvent
	ArchiveEvent *ArchivedEvent
}

type CreatedEvent struct {
	Offset           int64
	NodeID           int32
	ContractID       string
	TemplateID       string
	ContractKey      interface{}
	CreateArguments  interface{}
	CreatedEventBlob []byte
	InterfaceViews   []*InterfaceView
	WitnessParties   []string
	Signatories      []string
	Observers        []string
	CreatedAt        *time.Time
	PackageName      string
}

type InterfaceView struct {
	InterfaceID string
	ViewStatus  *ViewStatus
	ViewValue   interface{}
}

type ViewStatus struct {
	Code    int32
	Message string
}

type ArchivedEvent struct {
	Offset                int64
	NodeID                int32
	ContractID            string
	TemplateID            string
	WitnessParties        []string
	PackageName           string
	ImplementedInterfaces []string
}

type ExercisedEvent struct {
	Offset                int64
	NodeID                int32
	ContractID            string
	TemplateID            string
	InterfaceID           string
	Choice                string
	ChoiceArgument        interface{}
	ActingParties         []string
	Consuming             bool
	WitnessParties        []string
	LastDescendantNodeID  int32
	ExerciseResult        interface{}
	PackageName           string
	ImplementedInterfaces []string
}

// Package Service types
type ListPackagesRequest struct{}

type ListPackagesResponse struct {
	PackageIDs []string
}

type GetPackageRequest struct {
	PackageID string
}

type GetPackageResponse struct {
	ArchivePayload []byte
	HashFunction   HashFunction
	Hash           string
}

type HashFunction int32

const (
	HashFunctionSHA256 HashFunction = 0
)

type GetPackageStatusRequest struct {
	PackageID string
}

type GetPackageStatusResponse struct {
	PackageStatus PackageStatus
}

type PackageStatus int32

const (
	PackageStatusUnknown    PackageStatus = 0
	PackageStatusRegistered PackageStatus = 1
)

// State Service types
type GetActiveContractsRequest struct {
	Filter         *TransactionFilter
	Verbose        bool
	ActiveAtOffset int64
	EventFormat    *EventFormat
}

type GetActiveContractsResponse struct {
	Offset               int64
	WorkflowID           string
	ActiveContracts      []*CreatedEvent
	IncompleteUnassigned *IncompleteUnassigned
	IncompleteAssigned   *IncompleteAssigned
}

type IncompleteUnassigned struct {
	CreatedEvent    *CreatedEvent
	UnassignedEvent *UnassignedEvent
}

type IncompleteAssigned struct {
	AssignedEvent *AssignedEvent
}

type UnassignedEvent struct {
	UnassignID            string
	ContractID            string
	TemplateID            string
	Source                string
	Target                string
	Submitter             string
	ReassignmentCounter   uint64
	AssignmentExclusivity *time.Time
	WitnessParties        []string
	PackageName           string
	Offset                int64
}

type AssignedEvent struct {
	Source              string
	Target              string
	UnassignID          string
	Submitter           string
	ReassignmentCounter uint64
	CreatedEvent        *CreatedEvent
}

type EventFormat struct {
	FiltersByParty     map[string]*Filters
	FiltersForAnyParty *Filters
	Verbose            bool
}

type TransactionFilter struct {
	FiltersByParty map[string]*Filters
}

type Filters struct {
	Inclusive *InclusiveFilters
}

type InclusiveFilters struct {
	TemplateFilters []*TemplateFilter
}

type TemplateFilter struct {
	TemplateID              string
	InterfaceFilter         *InterfaceFilter
	IncludeCreatedEventBlob bool
}

type InterfaceFilter struct {
	InterfaceID          string
	IncludeInterfaceView bool
}

type GetConnectedSynchronizersRequest struct{}

type GetConnectedSynchronizersResponse struct {
	ConnectedSynchronizers []*ConnectedSynchronizer
}

type ConnectedSynchronizer struct {
	SynchronizerID        string
	ParticipantPermission ParticipantPermission
}

type ParticipantPermission int32

const (
	ParticipantPermissionSubmission   ParticipantPermission = 0
	ParticipantPermissionConfirmation ParticipantPermission = 1
	ParticipantPermissionObservation  ParticipantPermission = 2
)

type GetLedgerEndRequest struct{}

type GetLedgerEndResponse struct {
	Offset int64
}

type GetLatestPrunedOffsetsRequest struct{}

type GetLatestPrunedOffsetsResponse struct {
	ParticipantPrunedUpToInclusive          int64
	AllDivulgedContractsPrunedUpToInclusive int64
}

// Update Service types
type GetUpdatesRequest struct {
	BeginExclusive int64
	EndInclusive   *int64
	Filter         *TransactionFilter
	Verbose        bool
}

type GetUpdatesResponse struct {
	Update *Update
}

type Update struct {
	Transaction      *Transaction
	Reassignment     *Reassignment
	OffsetCheckpoint *OffsetCheckpoint
}

type Transaction struct {
	UpdateID    string
	CommandID   string
	WorkflowID  string
	EffectiveAt *time.Time
	Events      []*Event
	Offset      int64
}

type Event struct {
	Created   *CreatedEvent
	Archived  *ArchivedEvent
	Exercised *ExercisedEvent
}

type Reassignment struct {
	UpdateID    string
	Offset      int64
	UnassignID  string
	Source      string
	Target      string
	Counter     int64
	SubmittedAt *time.Time
	Unassigned  *time.Time
	Reassigned  *time.Time
}

type GetTransactionByIDRequest struct {
	UpdateID          string
	RequestingParties []string
}

type GetTransactionResponse struct {
	Transaction *Transaction
}

type GetTransactionByOffsetRequest struct {
	Offset            int64
	RequestingParties []string
}

// Version Service types
type GetLedgerAPIVersionRequest struct{}

type GetLedgerAPIVersionResponse struct {
	Version  string
	Features *FeaturesDescriptor
}

type FeaturesDescriptor struct {
	UserManagement   bool
	PartyManagement  bool
	OffsetCheckpoint bool
}

// Interactive Submission Service types
type PrepareSubmissionRequest struct {
	UserID                       string
	CommandID                    string
	Commands                     []*Command
	MinLedgerTime                *MinLedgerTime
	ActAs                        []string
	ReadAs                       []string
	DisclosedContracts           []*DisclosedContract
	SynchronizerID               string
	PackageIDSelectionPreference []string
	VerboseHashing               bool
	PrefetchContractKeys         []*PrefetchContractKey
}

type MinLedgerTime struct {
	Time MinLedgerTimeValue
}

type MinLedgerTimeValue interface {
	isMinLedgerTimeValue()
}

type MinLedgerTimeAbs struct {
	Time time.Time
}

func (MinLedgerTimeAbs) isMinLedgerTimeValue() {}

type MinLedgerTimeRel struct {
	Duration time.Duration
}

func (MinLedgerTimeRel) isMinLedgerTimeValue() {}

type DisclosedContract struct {
	TemplateID       string
	ContractID       string
	CreatedEventBlob []byte
	SynchronizerID   string
}

type PrefetchContractKey struct {
	TemplateID  string
	ContractKey map[string]interface{}
}

type PrepareSubmissionResponse struct {
	PreparedTransaction     []byte
	PreparedTransactionHash []byte
	HashingSchemeVersion    HashingSchemeVersion
	HashingDetails          string
}

type HashingSchemeVersion int32

const (
	HashingSchemeVersionUnspecified HashingSchemeVersion = 0
	HashingSchemeVersionV2          HashingSchemeVersion = 2
)

type ExecuteSubmissionRequest struct {
	PreparedTransaction  []byte
	PartySignatures      []*SinglePartySignatures
	DeduplicationPeriod  DeduplicationPeriod
	SubmissionID         string
	UserID               string
	HashingSchemeVersion HashingSchemeVersion
	MinLedgerTime        *MinLedgerTime
}

type SinglePartySignatures struct {
	Party      string
	Signatures []*Signature
}

type Signature struct {
	Format               SignatureFormat
	Signature            []byte
	SignedBy             string
	SigningAlgorithmSpec SigningAlgorithmSpec
}

type SignatureFormat int32

const (
	SignatureFormatUnspecified SignatureFormat = 0
	SignatureFormatRaw         SignatureFormat = 1
	SignatureFormatDER         SignatureFormat = 2
	SignatureFormatConcat      SignatureFormat = 3
	SignatureFormatSymbolic    SignatureFormat = 10000
)

type SigningAlgorithmSpec int32

const (
	SigningAlgorithmSpecUnspecified SigningAlgorithmSpec = 0
	SigningAlgorithmSpecED25519     SigningAlgorithmSpec = 1
	SigningAlgorithmSpecECDSASHA256 SigningAlgorithmSpec = 2
	SigningAlgorithmSpecECDSASHA384 SigningAlgorithmSpec = 3
)

type ExecuteSubmissionResponse struct{}

type GetPreferredPackageVersionRequest struct {
	Parties        []string
	PackageName    string
	SynchronizerID string
	VettingValidAt *time.Time
}

type GetPreferredPackageVersionResponse struct {
	PackageReference *PackageReference
	SynchronizerID   string
}

type PackageReference struct {
	PackageID      string
	PackageName    string
	PackageVersion string
}

type SDKVersionMismatchError struct {
	NodeVersion     string
	ContractVersion string
}

func (e *SDKVersionMismatchError) Error() string {
	return fmt.Sprintf("SDK version mismatch: node version %s, contract compiled with %s", e.NodeVersion, e.ContractVersion)
}
