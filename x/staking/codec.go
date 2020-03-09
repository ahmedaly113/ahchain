package staking

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers all the necessary types and interfaces for the module
func RegisterCodec(c *codec.Codec) {
	c.RegisterConcrete(MsgSubmitArgument{}, "ahchain/MsgSubmitArgument", nil)
	c.RegisterConcrete(MsgSubmitUpvote{}, "ahchain/MsgUpvoteArgument", nil)
	c.RegisterConcrete(MsgEditArgument{}, "ahchain/MsgEditArgument", nil)
	c.RegisterConcrete(MsgAddAdmin{}, "staking/MsgAddAdmin", nil)
	c.RegisterConcrete(MsgRemoveAdmin{}, "staking/MsgRemoveAdmin", nil)
	c.RegisterConcrete(MsgUpdateParams{}, "staking/MsgUpdateParams", nil)

	c.RegisterConcrete(Stake{}, "ahchain/Stake", nil)
	c.RegisterConcrete(Argument{}, "ahchain/Argument", nil)

}

// ModuleCodec encodes module codec
var ModuleCodec *codec.Codec

func init() {
	ModuleCodec = codec.New()
	RegisterCodec(ModuleCodec)
	codec.RegisterCrypto(ModuleCodec)
	ModuleCodec.Seal()
}
