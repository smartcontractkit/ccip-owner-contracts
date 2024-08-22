package operation

type Operation interface {
	GetOpType() string
	GetOpChainID() uint64
	GetOpData() interface{}
}
