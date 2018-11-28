# ArtifactDeckCodeGo
Artifact Card Game Deck Code Encoder & Decoder based on the [original api](https://github.com/ValveSoftware/ArtifactDeckCode) code written in PHP.

# Installation

`go get github.com/exts/ArtifactDeckCodeGo`

# Example Decode Usage

```go
// cardDeck returns a CardDeck struct reference
cardDeck, err := ArtifactDeckCodeGo.ParseDeck("ADCJQEAZX0ivAGABwA4XSXdAUEGQgEFAmIBRF0BDAkYAQUHAwEIBSQBMQFwASgBTw__")
if err != nil {
	// handle error
}

// cardDeck.name => string
// cardDeck.heroes => []Card
// cardDeck.cards => []Card
```

# Example Encode Usage

Encoding takes a CardDeck struct reference and spits out a deck code string. So provide the CardDeck struct with a full 40 card deck + 9 items and then 5 heroes and **excluding** signature cards (hero cards) to get a proper deck code back.

```go
// starter deck
codeStr := "ADCJQEAZX0ivAGABwA4XSXdAUEGQgEFAmIBRF0BDAkYAQUHAwEIBSQBMQFwASgBTw__"
cardDeck, err := ArtifactDeckCodeGo.ParseDeck(codeStr)
if err != nil {
	// handle error
}

deckCode, err := ArtifactDeckCodeGo.EncodeDeck(cardDeck)
if err != nil {
	// handle error
}

println(codeStr == deckCode) // true
```

# Notes
- Use card database to cross reference hero card id's to get the missing hero signature cards (each hero has a certain number of additional cards paired with them)