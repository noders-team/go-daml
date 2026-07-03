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

const packageNameSpliceApiFeaturedAppV2100 = "splice-api-featured-app-v2"

// IFeaturedAppActivityMarker2 is a DAML interface
type IFeaturedAppActivityMarker2 interface {
	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand
}

// IFeaturedAppRight2 is a DAML interface
type IFeaturedAppRight2 interface {
	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// FeaturedAppRightCreateActivityMarker executes the FeaturedAppRight_CreateActivityMarker choice
	FeaturedAppRightCreateActivityMarker(contractID string, args FeaturedAppRightCreateActivityMarker) *model.ExerciseCommand
}

// AppRewardBeneficiary2 is a Record type
type AppRewardBeneficiary2 struct {
	Beneficiary PARTY   `json:"beneficiary"`
	Weight      NUMERIC `json:"weight"`
}

// ToMap converts AppRewardBeneficiary2 to a map for DAML arguments
func (t AppRewardBeneficiary2) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["beneficiary"] = t.Beneficiary.ToMap()

	m["weight"] = (*big.Int)(t.Weight)

	return m
}

// MarshalJSON implements custom JSON marshaling for AppRewardBeneficiary2 using JsonCodec
func (t AppRewardBeneficiary2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AppRewardBeneficiary2 using JsonCodec
func (t *AppRewardBeneficiary2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FeaturedAppActivityMarkerView2 is a Record type
type FeaturedAppActivityMarkerView2 struct {
	Dso         PARTY   `json:"dso"`
	Provider    PARTY   `json:"provider"`
	Beneficiary PARTY   `json:"beneficiary"`
	Weight      NUMERIC `json:"weight"`
}

// ToMap converts FeaturedAppActivityMarkerView2 to a map for DAML arguments
func (t FeaturedAppActivityMarkerView2) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["dso"] = t.Dso.ToMap()

	m["provider"] = t.Provider.ToMap()

	m["beneficiary"] = t.Beneficiary.ToMap()

	m["weight"] = (*big.Int)(t.Weight)

	return m
}

// MarshalJSON implements custom JSON marshaling for FeaturedAppActivityMarkerView2 using JsonCodec
func (t FeaturedAppActivityMarkerView2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeaturedAppActivityMarkerView2 using JsonCodec
func (t *FeaturedAppActivityMarkerView2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FeaturedAppRightView2 is a Record type
type FeaturedAppRightView2 struct {
	Dso      PARTY `json:"dso"`
	Provider PARTY `json:"provider"`
}

// ToMap converts FeaturedAppRightView2 to a map for DAML arguments
func (t FeaturedAppRightView2) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["dso"] = t.Dso.ToMap()

	m["provider"] = t.Provider.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for FeaturedAppRightView2 using JsonCodec
func (t FeaturedAppRightView2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeaturedAppRightView2 using JsonCodec
func (t *FeaturedAppRightView2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FeaturedAppRightCreateActivityMarker2 is a Record type
type FeaturedAppRightCreateActivityMarker2 struct {
	Beneficiaries []AppRewardBeneficiary2 `json:"beneficiaries"`
	Weight        *NUMERIC                `json:"weight"`
}

// ToMap converts FeaturedAppRightCreateActivityMarker2 to a map for DAML arguments
func (t FeaturedAppRightCreateActivityMarker2) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["beneficiaries"] = t.Beneficiaries

	if t.Weight != nil {
		m["weight"] = map[string]interface{}{
			"_type": "optional",
			"value": (*big.Int)(*t.Weight),
		}
	} else {
		m["weight"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for FeaturedAppRightCreateActivityMarker2 using JsonCodec
func (t FeaturedAppRightCreateActivityMarker2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeaturedAppRightCreateActivityMarker2 using JsonCodec
func (t *FeaturedAppRightCreateActivityMarker2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FeaturedAppRightCreateActivityMarkerResult2 is a Record type
type FeaturedAppRightCreateActivityMarkerResult2 struct {
	ActivityMarkerCids []CONTRACT_ID `json:"activityMarkerCids"`
}

// ToMap converts FeaturedAppRightCreateActivityMarkerResult2 to a map for DAML arguments
func (t FeaturedAppRightCreateActivityMarkerResult2) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["activityMarkerCids"] = t.ActivityMarkerCids

	return m
}

// MarshalJSON implements custom JSON marshaling for FeaturedAppRightCreateActivityMarkerResult2 using JsonCodec
func (t FeaturedAppRightCreateActivityMarkerResult2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeaturedAppRightCreateActivityMarkerResult2 using JsonCodec
func (t *FeaturedAppRightCreateActivityMarkerResult2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// IFeaturedAppActivityMarker2InterfaceID returns the interface ID for the IFeaturedAppActivityMarker2 interface
func IFeaturedAppActivityMarker2InterfaceID(packageID *string) string {
	pkgName := packageNameSpliceApiFeaturedAppV2100
	if packageID != nil {
		pkgName = *packageID
	}
	return fmt.Sprintf("#%s:%s:%s", pkgName, "Splice.Api.FeaturedAppRightV2", "FeaturedAppActivityMarker")
}

// IFeaturedAppRight2InterfaceID returns the interface ID for the IFeaturedAppRight2 interface
func IFeaturedAppRight2InterfaceID(packageID *string) string {
	pkgName := packageNameSpliceApiFeaturedAppV2100
	if packageID != nil {
		pkgName = *packageID
	}
	return fmt.Sprintf("#%s:%s:%s", pkgName, "Splice.Api.FeaturedAppRightV2", "FeaturedAppRight")
}
