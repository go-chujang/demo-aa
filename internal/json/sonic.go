package json

import (
	"github.com/bytedance/sonic"
)

/*
usage
- import this package
*/
var (
	sonicCfg      = sonic.ConfigStd
	Marshal       = sonicCfg.Marshal
	Unmarshal     = sonicCfg.Unmarshal
	MarshalIndent = sonicCfg.MarshalIndent
	NewEncoder    = sonicCfg.NewEncoder
	NewDecoder    = sonicCfg.NewDecoder
)
