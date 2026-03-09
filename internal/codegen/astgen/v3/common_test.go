package v3

import (
	"testing"

	daml "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/daml/lf/archive/daml_lf_2"
	"github.com/stretchr/testify/require"
)

func TestExtractFieldWithInternedTypeApplicationsV3(t *testing.T) {
	codeGen := &codeGenAst{}

	pkg := &daml.Package{
		InternedStrings: []string{
			"",
			"optionalTsField",
			"listTextField",
			"contractIdField",
			"Tenant",
		},
		InternedTypes: []*daml.Type{
			{
				Sum: &daml.Type_Builtin_{
					Builtin: &daml.Type_Builtin{Builtin: daml.BuiltinType_OPTIONAL},
				},
			},
			{
				Sum: &daml.Type_Builtin_{
					Builtin: &daml.Type_Builtin{Builtin: daml.BuiltinType_TIMESTAMP},
				},
			},
			{
				Sum: &daml.Type_Tapp{
					Tapp: &daml.Type_TApp{
						Lhs: &daml.Type{Sum: &daml.Type_InternedType{InternedType: 0}},
						Rhs: &daml.Type{Sum: &daml.Type_InternedType{InternedType: 1}},
					},
				},
			},
			{
				Sum: &daml.Type_Builtin_{
					Builtin: &daml.Type_Builtin{Builtin: daml.BuiltinType_LIST},
				},
			},
			{
				Sum: &daml.Type_Builtin_{
					Builtin: &daml.Type_Builtin{Builtin: daml.BuiltinType_TEXT},
				},
			},
			{
				Sum: &daml.Type_Tapp{
					Tapp: &daml.Type_TApp{
						Lhs: &daml.Type{Sum: &daml.Type_InternedType{InternedType: 3}},
						Rhs: &daml.Type{Sum: &daml.Type_InternedType{InternedType: 4}},
					},
				},
			},
			{
				Sum: &daml.Type_Builtin_{
					Builtin: &daml.Type_Builtin{Builtin: daml.BuiltinType_CONTRACT_ID},
				},
			},
			{
				Sum: &daml.Type_Con_{
					Con: &daml.Type_Con{
						Tycon: &daml.TypeConId{
							NameInternedDname: 1,
						},
					},
				},
			},
			{
				Sum: &daml.Type_Tapp{
					Tapp: &daml.Type_TApp{
						Lhs: &daml.Type{Sum: &daml.Type_InternedType{InternedType: 6}},
						Rhs: &daml.Type{Sum: &daml.Type_InternedType{InternedType: 7}},
					},
				},
			},
		},
		InternedDottedNames: []*daml.InternedDottedName{
			nil,
			{SegmentsInternedStr: []int32{4}},
		},
	}

	testCases := []struct {
		name         string
		fieldNameIdx int32
		typeIdx      int32
		expectedName string
		expectedType string
	}{
		{
			name:         "optional timestamp from interned type app",
			fieldNameIdx: 1,
			typeIdx:      2,
			expectedName: "optionalTsField",
			expectedType: "*TIMESTAMP",
		},
		{
			name:         "list text from interned type app",
			fieldNameIdx: 2,
			typeIdx:      5,
			expectedName: "listTextField",
			expectedType: "[]TEXT",
		},
		{
			name:         "contract id from interned type app",
			fieldNameIdx: 3,
			typeIdx:      8,
			expectedName: "contractIdField",
			expectedType: "CONTRACT_ID",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			field := &daml.FieldWithType{
				FieldInternedStr: tc.fieldNameIdx,
				Type:             &daml.Type{Sum: &daml.Type_InternedType{InternedType: tc.typeIdx}},
			}

			fieldName, fieldType, err := codeGen.extractField(pkg, field)
			require.NoError(t, err)
			require.Equal(t, tc.expectedName, fieldName)
			require.Equal(t, tc.expectedType, fieldType)
		})
	}
}

func TestParseKeyExpressionV3(t *testing.T) {
	codeGen := &codeGenAst{}

	// Create a mock package with interned strings
	pkg := &daml.Package{
		InternedStrings: []string{
			"", // index 0 is usually empty
			"owner",
			"amount",
			"orderId",
			"customer",
		},
	}

	t.Run("Record projection key", func(t *testing.T) {
		// Create a key expression with a record projection (e.g., this.owner)
		key := &daml.DefTemplate_DefKey{
			KeyExpr: &daml.Expr{
				Sum: &daml.Expr_RecProj_{
					RecProj: &daml.Expr_RecProj{
						FieldInternedStr: 1, // "owner"
						Record: &daml.Expr{
							Sum: &daml.Expr_VarInternedStr{
								VarInternedStr: 1, // template variable
							},
						},
					},
				},
			},
		}

		fieldNames := codeGen.parseKeyExpression(pkg, key)
		require.Len(t, fieldNames, 2) // Both the field and the variable
		require.Contains(t, fieldNames, "owner")
	})

	t.Run("Record construction key", func(t *testing.T) {
		// Create a key expression with a record construction (composite key)
		key := &daml.DefTemplate_DefKey{
			KeyExpr: &daml.Expr{
				Sum: &daml.Expr_RecCon_{
					RecCon: &daml.Expr_RecCon{
						Fields: []*daml.FieldWithExpr{
							{
								FieldInternedStr: 1, // "owner"
								Expr: &daml.Expr{
									Sum: &daml.Expr_VarInternedStr{
										VarInternedStr: 1,
									},
								},
							},
							{
								FieldInternedStr: 3, // "orderId"
								Expr: &daml.Expr{
									Sum: &daml.Expr_VarInternedStr{
										VarInternedStr: 3,
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
		require.Contains(t, fieldNames, "owner")
		require.Contains(t, fieldNames, "orderId")
	})

	t.Run("Variable reference key", func(t *testing.T) {
		// Create a key expression with a simple variable reference
		key := &daml.DefTemplate_DefKey{
			KeyExpr: &daml.Expr{
				Sum: &daml.Expr_VarInternedStr{
					VarInternedStr: 4, // "customer"
				},
			},
		}

		fieldNames := codeGen.parseKeyExpression(pkg, key)
		require.Len(t, fieldNames, 1)
		require.Equal(t, "customer", fieldNames[0])
	})

	t.Run("Empty key expression", func(t *testing.T) {
		// Test with nil key
		key := &daml.DefTemplate_DefKey{
			KeyExpr: nil,
		}

		fieldNames := codeGen.parseKeyExpression(pkg, key)
		require.Len(t, fieldNames, 0)
	})
}
