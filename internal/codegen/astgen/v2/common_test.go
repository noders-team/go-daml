package v2

import (
	"testing"

	daml "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/daml/lf/archive/daml_lf_1"
	"github.com/stretchr/testify/require"
)

func TestParseKeyExpression(t *testing.T) {
	codeGen := &codeGenAst{}

	pkg := &daml.Package{
		InternedStrings: []string{
			"", // index 0 is usually empty
			"owner",
			"amount",
			"orderId",
			"customer",
		},
	}

	t.Run("Simple projection key", func(t *testing.T) {
		key := &daml.DefTemplate_DefKey{
			KeyExpr: &daml.DefTemplate_DefKey_Key{
				Key: &daml.KeyExpr{
					Sum: &daml.KeyExpr_Projections_{
						Projections: &daml.KeyExpr_Projections{
							Projections: []*daml.KeyExpr_Projection{
								{
									Field: &daml.KeyExpr_Projection_FieldInternedStr{
										FieldInternedStr: 1, // "owner"
									},
								},
							},
						},
					},
				},
			},
		}

		fieldNames := codeGen.parseKeyExpression(pkg, key)
		require.Len(t, fieldNames, 1)
		require.Equal(t, "owner", fieldNames[0])
	})

	t.Run("Record key with multiple fields", func(t *testing.T) {
		key := &daml.DefTemplate_DefKey{
			KeyExpr: &daml.DefTemplate_DefKey_Key{
				Key: &daml.KeyExpr{
					Sum: &daml.KeyExpr_Record_{
						Record: &daml.KeyExpr_Record{
							Fields: []*daml.KeyExpr_RecordField{
								{
									Field: &daml.KeyExpr_RecordField_FieldInternedStr{
										FieldInternedStr: 1, // "owner"
									},
								},
								{
									Field: &daml.KeyExpr_RecordField_FieldInternedStr{
										FieldInternedStr: 3, // "orderId"
									},
								},
							},
						},
					},
				},
			},
		}

		fieldNames := codeGen.parseKeyExpression(pkg, key)
		require.Len(t, fieldNames, 2)
		require.Equal(t, "owner", fieldNames[0])
		require.Equal(t, "orderId", fieldNames[1])
	})

	t.Run("Complex key expression", func(t *testing.T) {
		key := &daml.DefTemplate_DefKey{
			KeyExpr: &daml.DefTemplate_DefKey_ComplexKey{
				ComplexKey: &daml.Expr{
					// Complex expression - should return empty and log warning
				},
			},
		}

		fieldNames := codeGen.parseKeyExpression(pkg, key)
		require.Len(t, fieldNames, 0)
	})

	t.Run("String field names", func(t *testing.T) {
		key := &daml.DefTemplate_DefKey{
			KeyExpr: &daml.DefTemplate_DefKey_Key{
				Key: &daml.KeyExpr{
					Sum: &daml.KeyExpr_Projections_{
						Projections: &daml.KeyExpr_Projections{
							Projections: []*daml.KeyExpr_Projection{
								{
									Field: &daml.KeyExpr_Projection_FieldStr{
										FieldStr: "directFieldName",
									},
								},
							},
						},
					},
				},
			},
		}

		fieldNames := codeGen.parseKeyExpression(pkg, key)
		require.Len(t, fieldNames, 1)
		require.Equal(t, "directFieldName", fieldNames[0])
	})
}
