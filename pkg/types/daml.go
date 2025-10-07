package types

import (
	"math/big"
	"time"
)

type (
	PARTY       string
	TEXT        string
	INT64       int64
	BOOL        bool
	DECIMAL     *big.Int
	NUMERIC     *big.Int
	DATE        time.Time
	TIMESTAMP   time.Time
	UNIT        struct{}
	LIST        []string
	MAP         map[string]interface{}
	OPTIONAL    *interface{}
	GENMAP      map[string]interface{}
	CONTRACT_ID string
)

// VARIANT represents a DAML variant/union type
type VARIANT interface {
	GetVariantTag() string
	GetVariantValue() interface{}
}

// ENUM represents a DAML enum type
type ENUM interface {
	GetEnumConstructor() string
	GetEnumTypeID() string
}
