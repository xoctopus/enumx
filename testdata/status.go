package testdata

// Status
// +genx:enum
// @def attr.Gender=Gender
type Status int8

const (
	STATUS_UNKNOWN Status = iota
	// STATUS__ENABLED
	// @attr Gender=1
	STATUS__ENABLED // 关闭
	// STATUS__DISABLED
	// @attr Gender=2
	STATUS__DISABLED // 开启
	_                // placeholder will be ignored
)
