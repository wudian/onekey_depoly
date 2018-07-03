package handlers

import (
	at "gitlab.zhonganinfo.com/tech_bighealth/angine/types"
	"gitlab.zhonganinfo.com/tech_bighealth/ann-module/lib/go-wire"
)

var _ = wire.RegisterInterface(
	struct{ at.RPCResult }{},
	wire.ConcreteType{&at.ResultGenesis{}, at.ResultTypeGenesis},
	wire.ConcreteType{&at.ResultBlockchainInfo{}, at.ResultTypeBlockchainInfo},
	wire.ConcreteType{&at.ResultBlock{}, at.ResultTypeBlock},
	wire.ConcreteType{&at.ResultStatus{}, at.ResultTypeStatus},
	wire.ConcreteType{&at.ResultNetInfo{}, at.ResultTypeNetInfo},
	wire.ConcreteType{&at.ResultDialSeeds{}, at.ResultTypeDialSeeds},
	wire.ConcreteType{&at.ResultValidators{}, at.ResultTypeValidators},
	wire.ConcreteType{&at.ResultDumpConsensusState{}, at.ResultTypeDumpConsensusState},
	wire.ConcreteType{&at.ResultBroadcastTx{}, at.ResultTypeBroadcastTx},
	wire.ConcreteType{&at.ResultBroadcastTxCommit{}, at.ResultTypeBroadcastTxCommit},
	wire.ConcreteType{&at.ResultRequestSpecialOP{}, at.ResultTypeRequestSpecialOP},
	wire.ConcreteType{&at.ResultUnconfirmedTxs{}, at.ResultTypeUnconfirmedTxs},
	wire.ConcreteType{&at.ResultSubscribe{}, at.ResultTypeSubscribe},
	wire.ConcreteType{&at.ResultUnsubscribe{}, at.ResultTypeUnsubscribe},
	wire.ConcreteType{&at.ResultEvent{}, at.ResultTypeEvent},
	wire.ConcreteType{&at.ResultUnsafeSetConfig{}, at.ResultTypeUnsafeSetConfig},
	wire.ConcreteType{&at.ResultUnsafeProfile{}, at.ResultTypeUnsafeStartCPUProfiler},
	wire.ConcreteType{&at.ResultUnsafeProfile{}, at.ResultTypeUnsafeStopCPUProfiler},
	wire.ConcreteType{&at.ResultUnsafeProfile{}, at.ResultTypeUnsafeWriteHeapProfile},
	wire.ConcreteType{&at.ResultUnsafeFlushMempool{}, at.ResultTypeUnsafeFlushMempool},
	wire.ConcreteType{&at.ResultQuery{}, at.ResultTypeQuery},
	wire.ConcreteType{&at.ResultInfo{}, at.ResultTypeInfo},
	wire.ConcreteType{&at.ResultSurveillance{}, at.ResultTypeSurveillance},
	wire.ConcreteType{&at.ResultCoreVersion{}, at.ResultTypeCoreVersion},
)
