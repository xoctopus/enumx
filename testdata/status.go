package testdata

// Status
// +genx:enum
type Status int8

const (
	STATUS_UNKNOWN   Status = iota
	STATUS__ENABLED         // 关闭
	STATUS__DISABLED        // 开启
	_                       // placeholder
)
