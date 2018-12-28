package trustchain

type TrustBlock struct {
	BPM uint8
}

func NewBlock() TrustBlock {
	return TrustBlock{
		BPM: 0,
	}
}
