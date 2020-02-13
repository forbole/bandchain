package zoracle

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/keeper"
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey

	EventTypeRequest         = types.EventTypeRequest
	EventTypeReport          = types.EventTypeReport
	EventTypeRequestExecuted = types.EventTypeRequestExecuted
	AttributeKeyRequestID    = types.AttributeKeyRequestID
	AttributeKeyCodeHash     = types.AttributeKeyCodeHash
	AttributeKeyPrepare      = types.AttributeKeyPrepare
	AttributeKeyResult       = types.AttributeKeyResult
	AttributeKeyValidator    = types.AttributeKeyValidator
	AttributeKeyCodeName     = types.AttributeKeyCodeName
)

var (
	NewKeeper            = keeper.NewKeeper
	NewQuerier           = keeper.NewQuerier
	ModuleCdc            = types.ModuleCdc
	RegisterCodec        = types.RegisterCodec
	NewMsgRequestData    = types.NewMsgRequestData
	NewMsgReportData     = types.NewMsgReportData
	RequestStoreKey      = types.RequestStoreKey
	ResultStoreKey       = types.ResultStoreKey
	DataSourceStoreKey   = types.DataSourceStoreKey
	OracleScriptStoreKey = types.OracleScriptStoreKey

	ParamKeyTable = types.ParamKeyTable
	NewParams     = types.NewParams
	DefaultParams = types.DefaultParams
)

type (
	Keeper         = keeper.Keeper
	MsgRequestData = types.MsgRequestData
	MsgReportData  = types.MsgReportData
	// RequestInfo    = types.RequestInfo
)
