package ArtifactDeckCode

const CurrentVersion = 2
const EncodedPrefix = "ADC"

type Card struct {
	Id    int
	Turn  int
	Count int
}

type CardDeck struct {
	Name   string
	Heroes []Card
	Cards  []Card
}
