package aa

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-chujang/demo-aa/model"
)

type (
	PackedUserOperation   model.PackedUserOperation
	UserOperationWithHash struct {
		UserOperation *PackedUserOperation `json:"userOperation"`
		Hash          hexutil.Bytes        `json:"hash"`
	}
)

func (o PackedUserOperation) Message() model.PackedUserOperation { return model.PackedUserOperation(o) }
func (o PackedUserOperation) toModel() model.PackedUserOperation { return model.PackedUserOperation(o) }

var (
	_ json.Marshaler   = (*PackedUserOperation)(nil)
	_ json.Unmarshaler = (*PackedUserOperation)(nil)
)

func (o PackedUserOperation) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.toModel())
}

func (o *PackedUserOperation) UnmarshalJSON(input []byte) error {
	var dec model.PackedUserOperation
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	*o = PackedUserOperation(dec)
	return nil
}
