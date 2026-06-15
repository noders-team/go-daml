package model

import "time"

type AuthorizeRequest struct {
	Proposal              *TopologyTransactionProposal
	TransactionHash       string
	MustFullyAuthorize    bool
	ForceChanges          []ForceFlag
	SignedBy              []string
	Store                 *StoreID
	WaitToBecomeEffective *time.Duration
}

type AuthorizeResponse struct {
	Transaction *SignedTopologyTransaction
}

type AddTransactionsRequest struct {
	Transactions          []*SignedTopologyTransaction
	ForceChanges          []ForceFlag
	Store                 *StoreID
	WaitToBecomeEffective *time.Duration
}

type AddTransactionsResponse struct{}

type ListNamespaceDelegationRequest struct {
	BaseQuery                  *BaseQuery
	FilterNamespace            string
	FilterTargetKeyFingerprint string
}

type ListNamespaceDelegationResponse struct {
	Results []*NamespaceDelegationResult
}

// Deprecated: party-to-key mappings are deprecated in Canton; see PartyToKeyMapping.
type ListPartyToKeyMappingRequest struct {
	BaseQuery   *BaseQuery
	FilterParty string
}

// Deprecated: party-to-key mappings are deprecated in Canton; see PartyToKeyMapping.
type ListPartyToKeyMappingResponse struct {
	Results []*PartyToKeyMappingResult
}

type ListPartyToParticipantRequest struct {
	BaseQuery         *BaseQuery
	FilterParty       string
	FilterParticipant string
}

type ListPartyToParticipantResponse struct {
	Results []*PartyToParticipantResult
}

type BaseQuery struct {
	Store           *StoreID
	Proposals       bool
	TimeQuery       *TimeQuery
	Operation       Operation
	FilterSignedKey string
	ProtocolVersion *int32
}

type StoreID struct {
	Value string
}

type TimeQuery struct {
	Serial *int64
	Range  *TimeRange
}

type TimeRange struct {
	From  *time.Time
	Until *time.Time
}

type Operation int32

const (
	OperationUnspecified Operation = 0
	OperationAddReplace  Operation = 1
	OperationRemove      Operation = 2
)

type ForceFlag int32

const (
	ForceFlagUnspecified                           ForceFlag = 0
	ForceFlagAlienMember                           ForceFlag = 1
	ForceFlagLedgerTimeRecordTimeToleranceIncrease ForceFlag = 2
)

type SignedTopologyTransaction struct {
	Transaction                []byte
	Signatures                 []TopologyTransactionSignature
	MultiTransactionSignatures []*MultiTransactionSignatures
	Proposal                   bool
}

type MultiTransactionSignatures struct {
	TransactionHashes [][]byte
	Signatures        []TopologyTransactionSignature
}

type TopologyTransactionSignature struct {
	SignedBy        string
	Signature       []byte
	SignatureFormat int32
}

type TopologyTransactionProposal struct {
	Operation Operation
	Mapping   TopologyMapping
	Serial    uint32
}

type TopologyMapping interface {
	isTopologyMapping()
}

type NamespaceDelegationMapping struct {
	Namespace string
	TargetKey PublicKey
	// Deprecated: is_root_delegation is deprecated in Canton. A root delegation is
	// now expressed through the namespace delegation restriction (the ability to
	// sign all mappings); see the NamespaceDelegation restriction in protocol v30.
	IsRootDelegation bool
}

func (*NamespaceDelegationMapping) isTopologyMapping() {}

// Deprecated: PartyToKeyMapping is deprecated in Canton. Protocol signing keys for
// externally signed parties now live in PartyToParticipantMapping.SigningKeys.
type PartyToKeyMapping struct {
	Party       string
	Threshold   uint32
	SigningKeys []PublicKey
}

func (*PartyToKeyMapping) isTopologyMapping() {}

// Deprecated: key scheme is deprecated in Canton crypto; use SigningKeySpec
// (the key's algorithm is now carried by KeySpec, not Scheme).
type SigningKeyScheme int32

const (
	SigningKeySchemeUnspecified SigningKeyScheme = 0
	SigningKeySchemeED25519     SigningKeyScheme = 1
	SigningKeySchemeECDSAP256   SigningKeyScheme = 2
	SigningKeySchemeECDSAP384   SigningKeyScheme = 3
)

type SigningKeySpec int32

const (
	SigningKeySpecUnspecified SigningKeySpec = 0
	SigningKeySpecCurve25519  SigningKeySpec = 1
	SigningKeySpecP256        SigningKeySpec = 2
	SigningKeySpecP384        SigningKeySpec = 3
	SigningKeySpecSecp256k1   SigningKeySpec = 4
)

type SigningKeyUsage int32

const (
	SigningKeyUsageUnspecified             SigningKeyUsage = 0
	SigningKeyUsageNamespace               SigningKeyUsage = 1
	SigningKeyUsageIdentityDelegation      SigningKeyUsage = 2
	SigningKeyUsageSequencerAuthentication SigningKeyUsage = 3
	SigningKeyUsageProtocol                SigningKeyUsage = 4
	SigningKeyUsageProofOfOwnership        SigningKeyUsage = 5
)

type PublicKey struct {
	Format int32
	Key    []byte
	ID     string
	// Deprecated: key scheme is deprecated in Canton crypto; use KeySpec instead
	// (the algorithm is now carried by the key spec, not the scheme).
	Scheme  int32
	KeySpec int32
	Usage   []int32
}

type TopologyTransactionResult struct {
	Transaction  *SignedTopologyTransaction
	Status       ResultStatus
	ErrorMessage string
}

type ResultStatus int32

const (
	ResultStatusSuccess   ResultStatus = 0
	ResultStatusFailure   ResultStatus = 1
	ResultStatusDuplicate ResultStatus = 2
)

type NamespaceDelegationResult struct {
	Context *BaseResult
	Item    *NamespaceDelegationMapping
}

// Deprecated: party-to-key mappings are deprecated in Canton; see PartyToKeyMapping.
type PartyToKeyMappingResult struct {
	Context *BaseResult
	Item    *PartyToKeyMapping
}

type PartyToParticipantResult struct {
	Context *BaseResult
	Item    *PartyToParticipantMapping
}

type PartyToParticipantMapping struct {
	Party        string
	Threshold    uint32
	Participants []HostingParticipant
	// SigningKeys, when set, marks the party as externally signed: the protocol
	// signing keys live in the PartyToParticipant mapping (the modern replacement
	// for the deprecated PartyToKeyMapping). SigningKeysThreshold defaults to 1.
	SigningKeys          []PublicKey
	SigningKeysThreshold uint32
}

func (*PartyToParticipantMapping) isTopologyMapping() {}

type HostingParticipant struct {
	ParticipantUID string
	Permission     ParticipantPermission
}

type BaseResult struct {
	Store                *StoreID
	Sequenced            *time.Time
	ValidFrom            *time.Time
	ValidUntil           *time.Time
	Operation            Operation
	TransactionHash      []byte
	Serial               int32
	SignedByFingerprints []string
}

type FilterTargetKeyOrFingerprint struct {
	Fingerprint string
}

type SignTransactionsRequest struct {
	Transactions []*SignedTopologyTransaction
	SignedBy     []string
	Store        *StoreID
	ForceFlags   []ForceFlag
}

type SignTransactionsResponse struct {
	Transactions []*SignedTopologyTransaction
}

type GenerateTransactionsRequest struct {
	Proposals []*GenerateTransactionProposal
}

type GenerateTransactionProposal struct {
	Operation Operation
	Serial    uint32
	Mapping   TopologyMapping
	Store     *StoreID
}

type GenerateTransactionsResponse struct {
	GeneratedTransactions []*GeneratedTransaction
}

type GeneratedTransaction struct {
	SerializedTransaction []byte
	TransactionHash       []byte
}

type CreateTemporaryTopologyStoreRequest struct {
	Name            string
	ProtocolVersion uint32
}

type CreateTemporaryTopologyStoreResponse struct {
	StoreID *StoreID
}

type DropTemporaryTopologyStoreRequest struct {
	StoreID *StoreID
}

type DropTemporaryTopologyStoreResponse struct{}

type ImportTopologySnapshotRequest struct {
	TopologySnapshot      []byte
	Store                 *StoreID
	WaitToBecomeEffective *time.Duration
}

type ImportTopologySnapshotResponse struct{}
