package ledger

import (
	"encoding/json"
	"math/big"
	"reflect"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	v2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/types"
	"github.com/rs/zerolog/log"
)

func parseTemplateID(templateID string) (packageID, moduleName, entityName string) {
	parts := strings.Split(templateID, ":")
	if len(parts) == 3 {
		return parts[0], parts[1], parts[2]
	} else if len(parts) == 2 {
		return "", parts[0], parts[1]
	}
	return "", "", templateID
}

func commandsToProto(cmd *model.Commands) *v2.Commands {
	pbCmd := &v2.Commands{
		WorkflowId:   cmd.WorkflowID,
		UserId:       cmd.UserID,
		CommandId:    cmd.CommandID,
		Commands:     commandsArrayToProto(cmd.Commands),
		ActAs:        cmd.ActAs,
		ReadAs:       cmd.ReadAs,
		SubmissionId: cmd.SubmissionID,
	}

	if cmd.MinLedgerTimeAbs != nil {
		pbCmd.MinLedgerTimeAbs = timestamppb.New(*cmd.MinLedgerTimeAbs)
	}

	if cmd.MinLedgerTimeRel != nil {
		pbCmd.MinLedgerTimeRel = durationpb.New(*cmd.MinLedgerTimeRel)
	}

	switch dp := cmd.DeduplicationPeriod.(type) {
	case model.DeduplicationDuration:
		pbCmd.DeduplicationPeriod = &v2.Commands_DeduplicationDuration{
			DeduplicationDuration: durationpb.New(dp.Duration),
		}
	case model.DeduplicationOffset:
		pbCmd.DeduplicationPeriod = &v2.Commands_DeduplicationOffset{
			DeduplicationOffset: dp.Offset,
		}
	}

	return pbCmd
}

func commandsArrayToProto(cmds []*model.Command) []*v2.Command {
	result := make([]*v2.Command, len(cmds))
	for i, cmd := range cmds {
		result[i] = commandToProto(cmd)
	}
	return result
}

func commandToProto(cmd *model.Command) *v2.Command {
	pbCmd := &v2.Command{}

	switch c := cmd.Command.(type) {
	case *model.CreateCommand:
		packageID, moduleName, entityName := parseTemplateID(c.TemplateID)
		pbCmd.Command = &v2.Command_Create{
			Create: &v2.CreateCommand{
				TemplateId: &v2.Identifier{
					PackageId:  packageID,
					ModuleName: moduleName,
					EntityName: entityName,
				},
				CreateArguments: convertToRecord(c.Arguments),
			},
		}
	case *model.ExerciseCommand:
		packageID, moduleName, entityName := parseTemplateID(c.TemplateID)
		pbCmd.Command = &v2.Command_Exercise{
			Exercise: &v2.ExerciseCommand{
				ContractId: c.ContractID,
				TemplateId: &v2.Identifier{
					PackageId:  packageID,
					ModuleName: moduleName,
					EntityName: entityName,
				},
				Choice:         c.Choice,
				ChoiceArgument: mapToValue(c.Arguments),
			},
		}
	case *model.ExerciseByKeyCommand:
		packageID, moduleName, entityName := parseTemplateID(c.TemplateID)
		pbCmd.Command = &v2.Command_ExerciseByKey{
			ExerciseByKey: &v2.ExerciseByKeyCommand{
				TemplateId: &v2.Identifier{
					PackageId:  packageID,
					ModuleName: moduleName,
					EntityName: entityName,
				},
				ContractKey:    mapToValue(c.Key),
				Choice:         c.Choice,
				ChoiceArgument: mapToValue(c.Arguments),
			},
		}
	}

	return pbCmd
}

func transactionFilterToProto(filter *model.TransactionFilter) *v2.TransactionFilter {
	if filter == nil {
		return nil
	}

	pbFilter := &v2.TransactionFilter{
		FiltersByParty: make(map[string]*v2.Filters),
	}

	for party, filters := range filter.FiltersByParty {
		pbFilter.FiltersByParty[party] = filtersToProto(filters)
	}

	return pbFilter
}

func filtersToProto(filters *model.Filters) *v2.Filters {
	if filters == nil {
		return nil
	}

	pbFilters := &v2.Filters{}

	// Convert template filters to cumulative filters
	if filters.Inclusive != nil {
		for _, tf := range filters.Inclusive.TemplateFilters {
			pbFilters.Cumulative = append(pbFilters.Cumulative, &v2.CumulativeFilter{
				IdentifierFilter: &v2.CumulativeFilter_TemplateFilter{
					TemplateFilter: templateFilterToProto(tf),
				},
			})
		}
	}

	return pbFilters
}

func templateFilterToProto(tf *model.TemplateFilter) *v2.TemplateFilter {
	if tf == nil {
		return nil
	}

	packageID, moduleName, entityName := parseTemplateID(tf.TemplateID)
	return &v2.TemplateFilter{
		TemplateId: &v2.Identifier{
			PackageId:  packageID,
			ModuleName: moduleName,
			EntityName: entityName,
		},
		IncludeCreatedEventBlob: tf.IncludeCreatedEventBlob,
	}
}

func eventFormatToProto(format *model.EventFormat) *v2.EventFormat {
	if format == nil {
		return nil
	}
	filtersByParty := make(map[string]*v2.Filters)
	if format.FiltersByParty != nil {
		for key, val := range format.FiltersByParty {
			filtersByParty[key] = filtersToProto(val)
		}
	}
	return &v2.EventFormat{
		FiltersForAnyParty: filtersToProto(format.FiltersForAnyParty),
		FiltersByParty:     filtersByParty,
		Verbose:            format.Verbose,
	}
}

func createdEventFromProto(pb *v2.CreatedEvent) *model.CreatedEvent {
	if pb == nil {
		return nil
	}

	event := &model.CreatedEvent{
		ContractID:  pb.ContractId,
		TemplateID:  pb.TemplateId.String(),
		Signatories: pb.Signatories,
		Observers:   pb.Observers,
	}

	// Convert proto Values to map[string]interface{}
	if pb.CreateArguments != nil {
		event.CreateArguments = valueFromRecord(pb.CreateArguments)
	}
	if pb.ContractKey != nil {
		if key := valueFromProto(pb.ContractKey); key != nil {
			if m, ok := key.(map[string]interface{}); ok {
				event.ContractKey = m
			}
		}
	}

	return event
}

func archivedEventFromProto(pb *v2.ArchivedEvent) *model.ArchivedEvent {
	if pb == nil {
		return nil
	}

	return &model.ArchivedEvent{
		ContractID: pb.ContractId,
		TemplateID: pb.TemplateId.String(),
	}
}

func convertBigIntToNumeric(i *big.Int, scale int) *big.Rat {
	den := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(scale)), nil)
	return new(big.Rat).SetFrac(i, den)
}

func valueFromProto(pb *v2.Value) interface{} {
	if pb == nil {
		return nil
	}

	switch v := pb.Sum.(type) {
	case *v2.Value_Unit:
		return map[string]interface{}{"_type": "unit"}
	case *v2.Value_Bool:
		return v.Bool
	case *v2.Value_Int64:
		return v.Int64
	case *v2.Value_Text:
		return v.Text
	case *v2.Value_Numeric:
		return v.Numeric
	case *v2.Value_Party:
		return v.Party
	case *v2.Value_ContractId:
		return v.ContractId
	case *v2.Value_Date:
		return v.Date
	case *v2.Value_Timestamp:
		return v.Timestamp
	case *v2.Value_Optional:
		if v.Optional.Value != nil {
			return valueFromProto(v.Optional.Value)
		}
		return nil
	case *v2.Value_List:
		result := make([]interface{}, len(v.List.Elements))
		for i, elem := range v.List.Elements {
			result[i] = valueFromProto(elem)
		}
		return result
	case *v2.Value_Record:
		if v.Record == nil {
			return nil
		}
		record := make(map[string]interface{})
		for _, field := range v.Record.Fields {
			record[field.Label] = valueFromProto(field.Value)
		}
		return record
	case *v2.Value_TextMap:
		if v.TextMap == nil {
			return nil
		}
		result := make(map[string]interface{})
		for _, entry := range v.TextMap.Entries {
			result[entry.Key] = valueFromProto(entry.Value)
		}
		return result
	default:
		return nil
	}
}

func mapToValue(data interface{}) *v2.Value {
	if data == nil {
		return nil
	}

	// handle custom pointer types first before dereferencing
	switch v := data.(type) {
	case types.NUMERIC:
		return &v2.Value{Sum: &v2.Value_Numeric{Numeric: convertBigIntToNumeric((*big.Int)(v), 10).FloatString(10)}}
	case types.DECIMAL:
		return &v2.Value{Sum: &v2.Value_Numeric{Numeric: convertBigIntToNumeric((*big.Int)(v), 10).FloatString(10)}}
	case *big.Int:
		return &v2.Value{Sum: &v2.Value_Numeric{Numeric: convertBigIntToNumeric(v, 10).FloatString(10)}}
	case []types.INT64, []types.TEXT, []types.BOOL, []int64, []string:
		rv := reflect.ValueOf(v)
		elements := make([]*v2.Value, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			elements[i] = mapToValue(rv.Index(i).Interface())
		}
		return &v2.Value{
			Sum: &v2.Value_List{
				List: &v2.List{Elements: elements},
			},
		}
	case types.LIST:
		return &v2.Value{Sum: &v2.Value_List{List: &v2.List{Elements: mapValues(v)}}}
	case types.VARIANT:
		return &v2.Value{
			Sum: &v2.Value_Variant{
				Variant: &v2.Variant{
					Constructor: v.GetVariantTag(),
					Value:       mapToValue(v.GetVariantValue()),
				},
			},
		}
	case types.ENUM:
		return &v2.Value{
			Sum: &v2.Value_Enum{
				Enum: &v2.Enum{
					Constructor: v.GetEnumConstructor(),
				},
			},
		}
	}

	// handle pointers by dereferencing them
	if reflect.TypeOf(data).Kind() == reflect.Ptr {
		val := reflect.ValueOf(data)
		if val.IsNil() {
			return nil
		}
		return mapToValue(val.Elem().Interface())
	}

	// Handle custom types before other type checking
	switch v := data.(type) {
	case types.INT64:
		return &v2.Value{Sum: &v2.Value_Int64{Int64: int64(v)}}
	case types.TEXT:
		return &v2.Value{Sum: &v2.Value_Text{Text: string(v)}}
	case types.BOOL:
		return &v2.Value{Sum: &v2.Value_Bool{Bool: bool(v)}}
	case types.PARTY:
		return &v2.Value{Sum: &v2.Value_Party{Party: string(v)}}
	case types.CONTRACT_ID:
		return &v2.Value{Sum: &v2.Value_ContractId{ContractId: string(v)}}
	case types.DATE:
		return &v2.Value{Sum: &v2.Value_Date{Date: int32((time.Time)(v).Unix() / 86400)}}
	case types.TIMESTAMP:
		return &v2.Value{Sum: &v2.Value_Timestamp{Timestamp: int64((time.Time)(v).Unix())}}
	case bool:
		return &v2.Value{Sum: &v2.Value_Bool{Bool: v}}
	case int64:
		return &v2.Value{Sum: &v2.Value_Int64{Int64: v}}
	case int:
		return &v2.Value{Sum: &v2.Value_Int64{Int64: int64(v)}}
	case string:
		return &v2.Value{Sum: &v2.Value_Text{Text: v}}
	case []interface{}:
		elements := make([]*v2.Value, len(v))
		for i, elem := range v {
			elements[i] = mapToValue(elem)
		}
		return &v2.Value{
			Sum: &v2.Value_List{
				List: &v2.List{Elements: elements},
			},
		}
	case map[string]interface{}:
		if typeVal, hasType := v["_type"]; hasType && typeVal == "optional" {
			if val, hasValue := v["value"]; hasValue {
				return &v2.Value{
					Sum: &v2.Value_Optional{
						Optional: &v2.Optional{
							Value: mapToValue(val),
						},
					},
				}
			} else {
				return &v2.Value{
					Sum: &v2.Value_Optional{
						Optional: &v2.Optional{
							Value: nil,
						},
					},
				}
			}
		}

		if typeStr, ok := v["_type"].(string); ok && typeStr == "unit" {
			return &v2.Value{Sum: &v2.Value_Unit{Unit: &emptypb.Empty{}}}
		}

		if typeStr, ok := v["_type"].(string); ok && typeStr == "party" {
			if partyValue, ok := v["value"].(string); ok {
				return &v2.Value{Sum: &v2.Value_Party{Party: partyValue}}
			}
		}

		if typeStr, ok := v["_type"].(string); ok && typeStr == "genmap" {
			if mapValue, ok := v["value"].(map[string]interface{}); ok {
				entries := make([]*v2.TextMap_Entry, 0, len(mapValue))
				for key, val := range mapValue {
					entries = append(entries, &v2.TextMap_Entry{
						Key:   key,
						Value: mapToValue(val),
					})
				}
				return &v2.Value{
					Sum: &v2.Value_TextMap{
						TextMap: &v2.TextMap{
							Entries: entries,
						},
					},
				}
			} else if genMapValue, ok := v["value"].(types.GENMAP); ok {
				entries := make([]*v2.TextMap_Entry, 0, len(genMapValue))
				for key, val := range genMapValue {
					entries = append(entries, &v2.TextMap_Entry{
						Key:   key,
						Value: mapToValue(val),
					})
				}
				return &v2.Value{
					Sum: &v2.Value_TextMap{
						TextMap: &v2.TextMap{
							Entries: entries,
						},
					},
				}
			}
		}
		fields := make([]*v2.RecordField, 0, len(v))
		for key, val := range v {
			if key != "_type" && key != "value" {
				fields = append(fields, &v2.RecordField{
					Label: key,
					Value: mapToValue(val),
				})
			}
		}
		return &v2.Value{
			Sum: &v2.Value_Record{
				Record: &v2.Record{Fields: fields},
			},
		}
	case time.Time:
		return &v2.Value{Sum: &v2.Value_Date{Date: int32(v.Unix() / 86400)}}
	case interface{}:
		// Check if it implements VARIANT interface
		if variant, ok := v.(types.VARIANT); ok {
			return &v2.Value{
				Sum: &v2.Value_Variant{
					Variant: &v2.Variant{
						Constructor: variant.GetVariantTag(),
						Value:       mapToValue(variant.GetVariantValue()),
					},
				},
			}
		}
		return mapToValue(structToMap(v))
	default:
		return nil
	}
}

func structToMap(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		b, _ := json.Marshal(v)
		json.Unmarshal(b, &result)
		return result
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		if !fieldValue.CanInterface() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		tagName := jsonTag
		hasOmitEmpty := false
		if idx := strings.Index(jsonTag, ","); idx != -1 {
			tagName = jsonTag[:idx]
			options := jsonTag[idx+1:]
			hasOmitEmpty = strings.Contains(options, "omitempty")
		}

		if tagName == "" {
			tagName = strings.ToLower(field.Name)
		}

		actualValue := fieldValue.Interface()
		if hasOmitEmpty && actualValue == nil {
			continue
		}

		if hasOmitEmpty && fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}

		if hasOmitEmpty && fieldValue.IsZero() && fieldValue.Kind() != reflect.Ptr {
			continue
		}

		result[tagName] = actualValue
	}

	return result
}

func mapValues(values []string) []*v2.Value {
	result := make([]*v2.Value, len(values))
	for i, v := range values {
		result[i] = mapToValue(v)
	}
	return result
}

func convertToRecord(data map[string]interface{}) *v2.Record {
	if data == nil {
		return nil
	}

	fields := make([]*v2.RecordField, 0, len(data))
	for key, val := range data {
		val := mapToValue(val)
		if val == nil {
			log.Warn().Msgf("unsupported type %s for field %s, ignoring", reflect.TypeOf(val), key)
			continue
		}
		fields = append(fields, &v2.RecordField{
			Label: key,
			Value: val,
		})
	}

	return &v2.Record{Fields: fields}
}

func valueFromRecord(record *v2.Record) map[string]interface{} {
	if record == nil {
		return nil
	}

	result := make(map[string]interface{})
	for _, field := range record.Fields {
		result[field.Label] = valueFromProto(field.Value)
	}
	return result
}

func prepareSubmissionRequestToProto(req *model.PrepareSubmissionRequest) *interactive.PrepareSubmissionRequest {
	if req == nil {
		return nil
	}

	pbReq := &interactive.PrepareSubmissionRequest{
		UserId:                       req.UserID,
		CommandId:                    req.CommandID,
		Commands:                     commandsArrayToProto(req.Commands),
		ActAs:                        req.ActAs,
		ReadAs:                       req.ReadAs,
		SynchronizerId:               req.SynchronizerID,
		PackageIdSelectionPreference: req.PackageIDSelectionPreference,
		VerboseHashing:               req.VerboseHashing,
	}

	if req.MinLedgerTime != nil {
		pbReq.MinLedgerTime = minLedgerTimeToProto(req.MinLedgerTime)
	}

	if req.DisclosedContracts != nil {
		pbReq.DisclosedContracts = disclosedContractsToProto(req.DisclosedContracts)
	}

	if req.PrefetchContractKeys != nil {
		pbReq.PrefetchContractKeys = prefetchContractKeysToProto(req.PrefetchContractKeys)
	}

	return pbReq
}

func minLedgerTimeToProto(mlt *model.MinLedgerTime) *interactive.MinLedgerTime {
	if mlt == nil || mlt.Time == nil {
		return nil
	}

	pbMLT := &interactive.MinLedgerTime{}

	switch t := mlt.Time.(type) {
	case model.MinLedgerTimeAbs:
		pbMLT.Time = &interactive.MinLedgerTime_MinLedgerTimeAbs{
			MinLedgerTimeAbs: timestamppb.New(t.Time),
		}
	case model.MinLedgerTimeRel:
		pbMLT.Time = &interactive.MinLedgerTime_MinLedgerTimeRel{
			MinLedgerTimeRel: durationpb.New(t.Duration),
		}
	}

	return pbMLT
}

func disclosedContractsToProto(contracts []*model.DisclosedContract) []*v2.DisclosedContract {
	if contracts == nil {
		return nil
	}

	result := make([]*v2.DisclosedContract, len(contracts))
	for i, contract := range contracts {
		result[i] = disclosedContractToProto(contract)
	}
	return result
}

func disclosedContractToProto(contract *model.DisclosedContract) *v2.DisclosedContract {
	if contract == nil {
		return nil
	}

	packageID, moduleName, entityName := parseTemplateID(contract.TemplateID)
	pbContract := &v2.DisclosedContract{
		TemplateId: &v2.Identifier{
			PackageId:  packageID,
			ModuleName: moduleName,
			EntityName: entityName,
		},
		ContractId:       contract.ContractID,
		SynchronizerId:   contract.SynchronizerID,
		CreatedEventBlob: contract.CreatedEventBlob,
	}

	return pbContract
}

func prefetchContractKeysToProto(keys []*model.PrefetchContractKey) []*v2.PrefetchContractKey {
	if keys == nil {
		return nil
	}

	result := make([]*v2.PrefetchContractKey, len(keys))
	for i, key := range keys {
		packageID, moduleName, entityName := parseTemplateID(key.TemplateID)
		result[i] = &v2.PrefetchContractKey{
			TemplateId: &v2.Identifier{
				PackageId:  packageID,
				ModuleName: moduleName,
				EntityName: entityName,
			},
			ContractKey: mapToValue(key.ContractKey),
		}
	}
	return result
}

func prepareSubmissionResponseFromProto(pb *interactive.PrepareSubmissionResponse) *model.PrepareSubmissionResponse {
	if pb == nil {
		return nil
	}

	resp := &model.PrepareSubmissionResponse{
		PreparedTransactionHash: pb.PreparedTransactionHash,
		HashingSchemeVersion:    model.HashingSchemeVersion(pb.HashingSchemeVersion),
	}

	if pb.PreparedTransaction != nil {
		data, _ := proto.Marshal(pb.PreparedTransaction)
		resp.PreparedTransaction = data
	}

	if pb.HashingDetails != nil {
		resp.HashingDetails = *pb.HashingDetails
	}

	return resp
}

func executeSubmissionRequestToProto(req *model.ExecuteSubmissionRequest) *interactive.ExecuteSubmissionRequest {
	if req == nil {
		return nil
	}

	pbReq := &interactive.ExecuteSubmissionRequest{
		SubmissionId:         req.SubmissionID,
		UserId:               req.UserID,
		HashingSchemeVersion: interactive.HashingSchemeVersion(req.HashingSchemeVersion),
	}

	if req.PreparedTransaction != nil {
		pt := &interactive.PreparedTransaction{}
		proto.Unmarshal(req.PreparedTransaction, pt)
		pbReq.PreparedTransaction = pt
	}

	if len(req.PartySignatures) > 0 {
		pbReq.PartySignatures = &interactive.PartySignatures{
			Signatures: singlePartySignaturesToProto(req.PartySignatures),
		}
	}

	if req.MinLedgerTime != nil {
		pbReq.MinLedgerTime = minLedgerTimeToProto(req.MinLedgerTime)
	}

	switch dp := req.DeduplicationPeriod.(type) {
	case model.DeduplicationDuration:
		pbReq.DeduplicationPeriod = &interactive.ExecuteSubmissionRequest_DeduplicationDuration{
			DeduplicationDuration: durationpb.New(dp.Duration),
		}
	case model.DeduplicationOffset:
		pbReq.DeduplicationPeriod = &interactive.ExecuteSubmissionRequest_DeduplicationOffset{
			DeduplicationOffset: dp.Offset,
		}
	}

	return pbReq
}

func singlePartySignaturesToProto(sigs []*model.SinglePartySignatures) []*interactive.SinglePartySignatures {
	if sigs == nil {
		return nil
	}

	result := make([]*interactive.SinglePartySignatures, len(sigs))
	for i, sig := range sigs {
		result[i] = &interactive.SinglePartySignatures{
			Party:      sig.Party,
			Signatures: signaturesToProto(sig.Signatures),
		}
	}
	return result
}

func signaturesToProto(sigs []*model.Signature) []*interactive.Signature {
	if sigs == nil {
		return nil
	}

	result := make([]*interactive.Signature, len(sigs))
	for i, sig := range sigs {
		result[i] = &interactive.Signature{
			Format:               interactive.SignatureFormat(sig.Format),
			Signature:            sig.Signature,
			SignedBy:             sig.SignedBy,
			SigningAlgorithmSpec: interactive.SigningAlgorithmSpec(sig.SigningAlgorithmSpec),
		}
	}
	return result
}

func identifierToString(id *v2.Identifier) string {
	if id == nil {
		return ""
	}
	if id.PackageId != "" {
		return id.PackageId + ":" + id.ModuleName + ":" + id.EntityName
	}
	return id.ModuleName + ":" + id.EntityName
}

func unassignedEventFromProto(pb *v2.UnassignedEvent) *model.UnassignedEvent {
	if pb == nil {
		return nil
	}

	var assignmentExclusivity *time.Time
	if pb.AssignmentExclusivity != nil {
		t := pb.AssignmentExclusivity.AsTime()
		assignmentExclusivity = &t
	}

	var templateID string
	if pb.TemplateId != nil {
		templateID = identifierToString(pb.TemplateId)
	}

	return &model.UnassignedEvent{
		UnassignID:            pb.UnassignId,
		ContractID:            pb.ContractId,
		TemplateID:            templateID,
		Source:                pb.Source,
		Target:                pb.Target,
		Submitter:             pb.Submitter,
		ReassignmentCounter:   pb.ReassignmentCounter,
		AssignmentExclusivity: assignmentExclusivity,
		WitnessParties:        pb.WitnessParties,
		PackageName:           pb.PackageName,
		Offset:                pb.Offset,
	}
}

func assignedEventFromProto(pb *v2.AssignedEvent) *model.AssignedEvent {
	if pb == nil {
		return nil
	}

	return &model.AssignedEvent{
		Source:              pb.Source,
		Target:              pb.Target,
		UnassignID:          pb.UnassignId,
		Submitter:           pb.Submitter,
		ReassignmentCounter: pb.ReassignmentCounter,
		CreatedEvent:        createdEventFromProto(pb.CreatedEvent),
	}
}

func exercisedEventFromProto(pb *v2.ExercisedEvent) *model.ExercisedEvent {
	if pb == nil {
		return nil
	}

	event := &model.ExercisedEvent{
		Offset:               pb.Offset,
		NodeID:               pb.NodeId,
		ContractID:           pb.ContractId,
		Choice:               pb.Choice,
		ActingParties:        pb.ActingParties,
		Consuming:            pb.Consuming,
		WitnessParties:       pb.WitnessParties,
		LastDescendantNodeID: pb.LastDescendantNodeId,
		PackageName:          pb.PackageName,
	}

	if pb.TemplateId != nil {
		event.TemplateID = identifierToString(pb.TemplateId)
	}

	if pb.InterfaceId != nil {
		event.InterfaceID = identifierToString(pb.InterfaceId)
	}

	if pb.ChoiceArgument != nil {
		if arg := valueFromProto(pb.ChoiceArgument); arg != nil {
			if m, ok := arg.(map[string]interface{}); ok {
				event.ChoiceArgument = m
			}
		}
	}

	if pb.ExerciseResult != nil {
		event.ExerciseResult = valueFromProto(pb.ExerciseResult)
	}

	for _, iface := range pb.ImplementedInterfaces {
		event.ImplementedInterfaces = append(event.ImplementedInterfaces, identifierToString(iface))
	}

	return event
}
