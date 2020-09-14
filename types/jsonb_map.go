package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jackc/pgx/pgtype"

	"github.com/pkg/errors"
)

// TODO: consider using pgx intSlice
// TODO[ak|06/2018]: Scan() and Value() should work as is
// TODO[ak|06/2018]: I'm not sure about the text/binary related functions they are similar to pgtype.JSON and pgtype.JSONB
// Sourced from pgtype.JSON and pgype.JSONB
type JSONBMap map[string]interface{}

func (dst *JSONBMap) DecodeText(ci *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		*dst = map[string]interface{}{}
		return nil
	}
	return json.Unmarshal(src, dst)
}

// Scan implements the database/sql Scanner interface.
func (dst *JSONBMap) Scan(src interface{}) error {
	if src == nil {
		*dst = map[string]interface{}{}
		return nil
	}

	switch src := src.(type) {
	case string:
		return dst.DecodeText(nil, []byte(src))
	case []byte:
		srcCopy := make([]byte, len(src))
		copy(srcCopy, src)
		return dst.DecodeText(nil, srcCopy)
	}

	return errors.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src JSONBMap) Value() (driver.Value, error) {
	if src == nil {
		return "{}", nil
	}
	b, err := json.Marshal(src)
	return b, err
}
