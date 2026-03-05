package v3

import (
	"testing"

	daml "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/daml/lf/archive/daml_lf_2"
	"github.com/stretchr/testify/require"
)

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
