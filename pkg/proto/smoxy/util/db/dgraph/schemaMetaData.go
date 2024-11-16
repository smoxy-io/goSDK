package dgraph

import (
	"encoding/json"
	"strconv"
)

func (m MigrationStatus) Name() string {
	return MigrationStatus_name[int32(m)]
}

func (m MigrationStatus) Value() int32 {
	return int32(m)
}

func (m *MigrationStatus) UnmarshalJSON(d []byte) error {
	// Migrations status can be a numeric base 10 string or an integer depending on if it's being sourced from
	// a protobuf data transmission (integer) or the database (numeric base 10 string)
	var si int32

	if err := json.Unmarshal(d, &si); err == nil {
		// value is an integer
		*m = MigrationStatus(si)
		return nil
	}

	var ss string

	if err := json.Unmarshal(d, &ss); err != nil {
		// value is not a string either
		return err
	}

	s, err := strconv.Atoi(ss)

	if err != nil {
		// not a numeric base 10 string
		return err
	}

	*m = MigrationStatus(s)

	return nil
}

func MigrationStatusFromString(status string) MigrationStatus {
	return MigrationStatus(MigrationStatus_value[status])
}
