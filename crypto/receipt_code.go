package crypto

// ReceiptCode indicates status of receipt after tx application
type ReceiptCode byte

// ReceiptCode values
const (
	ReceiptCodeOK               ReceiptCode = 0x0
	ReceiptCodeOutOfGas         ReceiptCode = 0x1
	ReceiptCodeIgniteError      ReceiptCode = 0x2
	ReceiptCodeContractNotFound ReceiptCode = 0x3
	ReceiptCodeMethodNotFound   ReceiptCode = 0x4
)
