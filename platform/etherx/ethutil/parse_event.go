package ethutil

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrParseLogNilTopics    = errors.New("nil topics")
	ErrParseLogUnknownEvent = errors.New("unknown event")
	ErrParseErrInvalidData  = errors.New("invalid data")
	ErrParseErrUnknownError = errors.New("unknown error")
)

func ParseLog(log types.Log, eventMap map[string]abi.Event) (string, map[string]interface{}, error) {
	if log.Topics == nil {
		return "", nil, ErrParseLogNilTopics
	}
	event, ok := eventMap[log.Topics[0].Hex()]
	if !ok {
		return "", nil, ErrParseLogUnknownEvent
	}

	var (
		parsed     = make(map[string]interface{})
		indexed    abi.Arguments
		nonIndexed = event.Inputs.NonIndexed()
	)
	for _, arg := range event.Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopicsIntoMap(parsed, indexed, log.Topics[1:]); err != nil {
		return "", nil, err
	}
	if err := nonIndexed.UnpackIntoMap(parsed, log.Data); err != nil {
		return "", nil, err
	}
	return event.Name, parsed, nil
}

func ParseError(data []byte, errorMap map[string]abi.Error) (string, map[string]interface{}, error) {
	if data == nil || len(data) < 4 {
		return "", nil, ErrParseErrInvalidData
	}

	var (
		selector   = data[:4]
		errorID    = hexutil.Encode(selector)
		errDef, ok = errorMap[errorID]
	)
	if !ok {
		return "", nil, ErrParseErrUnknownError
	}
	args, err := errDef.Inputs.Unpack(data[4:])
	if err != nil {
		return "", nil, err
	}

	parsed := make(map[string]interface{})
	for i, arg := range errDef.Inputs {
		parsed[arg.Name] = args[i]
	}
	return errDef.Name, parsed, nil
}
