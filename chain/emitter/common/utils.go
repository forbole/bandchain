package common

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Atoi converts the given string into an int64. Panics on errors.
func Atoi(val string) int64 {
	res, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		panic(err)
	}
	return res
}

func ParseBytes(b []byte) []byte {
	if len(b) == 0 {
		return []byte{}
	}
	return b
}

func ParseEvents(events sdk.StringEvents) EvMap {
	evMap := make(EvMap)
	for _, event := range events {
		for _, kv := range event.Attributes {
			key := event.Type + "." + kv.Key
			evMap[key] = append(evMap[key], kv.Value)
		}
	}
	return evMap
}