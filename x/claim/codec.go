package claim

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(c *codec.Codec) {
	c.RegisterConcrete(MsgCreateClaim{}, "ahchain/MsgCreateClaim", nil)
	c.RegisterConcrete(MsgEditClaim{}, "ahchain/MsgEditClaim", nil)
	c.RegisterConcrete(MsgDeleteClaim{}, "ahchain/MsgDeleteClaim", nil)
	c.RegisterConcrete(MsgAddAdmin{}, "claim/MsgAddAdmin", nil)
	c.RegisterConcrete(MsgRemoveAdmin{}, "claim/MsgRemoveAdmin", nil)
	c.RegisterConcrete(MsgUpdateParams{}, "claim/MsgUpdateParams", nil)

	c.RegisterConcrete(Claim{}, "ahchain/Claim", nil)
}

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)
	ModuleCodec.Seal()
}
