package tests

import (
	"github.com/exts/ArtifactDeckCodeGo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodedDeckMatchesDeckString(t *testing.T) {

	deckString := "ADCJRwSJX2Dc7wBEAN4XUFBcN0BQmQBQWABRCgBCgN0AWUBbQFDbwEISEJsdWUvQmxhY2sgQ29udHJvbA__"

	// get an actual deck to test
	deck, _ := ArtifactDeckCodeGo.ParseDeck(deckString)

	// pass decoded deck to be parsed back into string form
	encodedDeckString, _ := ArtifactDeckCodeGo.EncodeDeck(deck)

	assert.Equal(t, deckString, encodedDeckString)
}

func TestEncodeBytesReturnsValidBytes(t *testing.T) {

	// get an actual deck to test
	deck, _ := ArtifactDeckCodeGo.ParseDeck("ADCJRwSJX2Dc7wBEAN4XUFBcN0BQmQBQWABRCgBCgN0AWUBbQFDbwEISEJsdWUvQmxhY2sgQ29udHJvbA__")

	_, encodedDeck := ArtifactDeckCodeGo.EncodeBytes(deck)
	assert.Nil(t, encodedDeck)
}

func TestEncodedBytesInvalidDeck(t *testing.T) {
	deck := &ArtifactDeckCodeGo.CardDeck{}
	deck2 := &ArtifactDeckCodeGo.CardDeck{
		Heroes: []ArtifactDeckCodeGo.Card{
			ArtifactDeckCodeGo.Card{
				Id: 10,
			},
		},
	}

	_, deckBytes := ArtifactDeckCodeGo.EncodeBytes(nil)
	_, deckBytes2 := ArtifactDeckCodeGo.EncodeBytes(deck)
	_, deckBytes3 := ArtifactDeckCodeGo.EncodeBytes(deck2)

	assert.NotNil(t, deckBytes)
	assert.NotNil(t, deckBytes2)
	assert.NotNil(t, deckBytes3)
}
