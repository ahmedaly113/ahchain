package slashing

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSlashArgument{}, "ahchain/MsgSlashArgument", nil)
	cdc.RegisterConcrete(MsgAddAdmin{}, "slashing/MsgAddAdmin", nil)
	cdc.RegisterConcrete(MsgRemoveAdmin{}, "slashing/MsgRemoveAdmin", nil)
	cdc.RegisterConcrete(MsgUpdateParams{}, "slashing/MsgUpdateParams", nil)

	cdc.RegisterConcrete(Slash{}, "ahchain/Slash", nil)
}

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)
	ModuleCodec.Seal()
}
