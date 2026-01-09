package types

type IBlock interface {
	Height() uint64
	Hash() string
	ParentHash() string
	BlockTime() uint64
}
