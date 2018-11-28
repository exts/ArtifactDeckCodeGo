package ArtifactDeckCodeGo

// encoding & decoding
const CurrentVersion = 2
const EncodedPrefix = "ADC"

// encoding
const HeaderSize = 3
const MaxBytesForVarUint32 = 5

// card type we store card data in
type Card struct {
	Id    int
	Turn  int
	Count int
}

// card deck to keep track of all cards
type CardDeck struct {
	Name   string
	Heroes []Card
	Cards  []Card
}
