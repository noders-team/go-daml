package gen_clients

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/noders-team/go-daml/pkg/codec"
	"github.com/noders-team/go-daml/pkg/model"
	. "github.com/noders-team/go-daml/pkg/types"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = fmt.Sprintf
	_ = codec.NewJsonCodec
	_ = model.CreateCommand{}
	_ = NewNumericFromDecimal
)

const packageNameSpliceApiFeaturedAppV1100 = "splice-api-featured-app-v1"

// IFeaturedAppActivityMarker is a DAML interface
type IFeaturedAppActivityMarker interface {
	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand
}

// IFeaturedAppRight is a DAML interface
type IFeaturedAppRight interface {
	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// FeaturedAppRightCreateActivityMarker executes the FeaturedAppRight_CreateActivityMarker choice
	FeaturedAppRightCreateActivityMarker(contractID string, args FeaturedAppRightCreateActivityMarker) *model.ExerciseCommand
}

// AppRewardBeneficiary is a Record type
type AppRewardBeneficiary struct {
	Beneficiary PARTY   `json:"beneficiary"`
	Weight      NUMERIC `json:"weight"`
}

// ToMap converts AppRewardBeneficiary to a map for DAML arguments
func (t AppRewardBeneficiary) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["beneficiary"] = t.Beneficiary.ToMap()

	m["weight"] = (*big.Int)(t.Weight)

	return m
}

// MarshalJSON implements custom JSON marshaling for AppRewardBeneficiary using JsonCodec
func (t AppRewardBeneficiary) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AppRewardBeneficiary using JsonCodec
func (t *AppRewardBeneficiary) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FeaturedAppActivityMarkerView is a Record type
type FeaturedAppActivityMarkerView struct {
	Dso         PARTY   `json:"dso"`
	Provider    PARTY   `json:"provider"`
	Beneficiary PARTY   `json:"beneficiary"`
	Weight      NUMERIC `json:"weight"`
}

// ToMap converts FeaturedAppActivityMarkerView to a map for DAML arguments
func (t FeaturedAppActivityMarkerView) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["dso"] = t.Dso.ToMap()

	m["provider"] = t.Provider.ToMap()

	m["beneficiary"] = t.Beneficiary.ToMap()

	m["weight"] = (*big.Int)(t.Weight)

	return m
}

// MarshalJSON implements custom JSON marshaling for FeaturedAppActivityMarkerView using JsonCodec
func (t FeaturedAppActivityMarkerView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeaturedAppActivityMarkerView using JsonCodec
func (t *FeaturedAppActivityMarkerView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FeaturedAppRightView is a Record type
type FeaturedAppRightView struct {
	Dso      PARTY `json:"dso"`
	Provider PARTY `json:"provider"`
}

// ToMap converts FeaturedAppRightView to a map for DAML arguments
func (t FeaturedAppRightView) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["dso"] = t.Dso.ToMap()

	m["provider"] = t.Provider.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for FeaturedAppRightView using JsonCodec
func (t FeaturedAppRightView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeaturedAppRightView using JsonCodec
func (t *FeaturedAppRightView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FeaturedAppRightCreateActivityMarker is a Record type
type FeaturedAppRightCreateActivityMarker struct {
	Beneficiaries []AppRewardBeneficiary `json:"beneficiaries"`
}

// ToMap converts FeaturedAppRightCreateActivityMarker to a map for DAML arguments
func (t FeaturedAppRightCreateActivityMarker) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["beneficiaries"] = t.Beneficiaries

	return m
}

// MarshalJSON implements custom JSON marshaling for FeaturedAppRightCreateActivityMarker using JsonCodec
func (t FeaturedAppRightCreateActivityMarker) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeaturedAppRightCreateActivityMarker using JsonCodec
func (t *FeaturedAppRightCreateActivityMarker) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FeaturedAppRightCreateActivityMarkerResult is a Record type
type FeaturedAppRightCreateActivityMarkerResult struct {
	ActivityMarkerCids []CONTRACT_ID `json:"activityMarkerCids"`
}

// ToMap converts FeaturedAppRightCreateActivityMarkerResult to a map for DAML arguments
func (t FeaturedAppRightCreateActivityMarkerResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["activityMarkerCids"] = t.ActivityMarkerCids

	return m
}

// MarshalJSON implements custom JSON marshaling for FeaturedAppRightCreateActivityMarkerResult using JsonCodec
func (t FeaturedAppRightCreateActivityMarkerResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeaturedAppRightCreateActivityMarkerResult using JsonCodec
func (t *FeaturedAppRightCreateActivityMarkerResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// IFeaturedAppActivityMarkerInterfaceID returns the interface ID for the IFeaturedAppActivityMarker interface
func IFeaturedAppActivityMarkerInterfaceID(packageID *string) string {
	pkgName := packageNameSpliceApiFeaturedAppV1100
	if packageID != nil {
		pkgName = *packageID
	}
	return fmt.Sprintf("#%s:%s:%s", pkgName, "Splice.Api.FeaturedAppRightV1", "FeaturedAppActivityMarker")
}

// IFeaturedAppRightInterfaceID returns the interface ID for the IFeaturedAppRight interface
func IFeaturedAppRightInterfaceID(packageID *string) string {
	pkgName := packageNameSpliceApiFeaturedAppV1100
	if packageID != nil {
		pkgName = *packageID
	}
	return fmt.Sprintf("#%s:%s:%s", pkgName, "Splice.Api.FeaturedAppRightV1", "FeaturedAppRight")
}
