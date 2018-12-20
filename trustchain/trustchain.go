package trustchain

type TrustBlock struct {
	Comment   string
	Score     uint8
	Address   string
	Signature string
}

func NewBlock() TrustBlock {
	return TrustBlock{
		Comment:   "",
		Score:     0,
		Address:   "",
		Signature: "",
	}
}
