package testdata

// Gender
// +genx:enum
type Gender int8

const (
	GENDER_UNKNOWN Gender = iota
	GENDER__MALE          // 男
	GENDER__FEMALE        // 女
)
