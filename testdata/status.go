package testdata

// Status
// +genx:enum
// @def attr.gender=Gender
type Status int8

const (
	STATUS_UNKNOWN Status = iota
	// STATUS__ENABLED
	// @attr gender=GENDER__MALE
	STATUS__ENABLED // 关闭
	// STATUS__DISABLED
	// @attr gender=GENDER__FEMALE
	STATUS__DISABLED // 开启
	_                // placeholder will be ignored
)
