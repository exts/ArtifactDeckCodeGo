package tests

import (
	"ArtifactDeckCode"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodedDeckMatchesDeckString(t *testing.T) {

	deckString := "ADCJRwSJX2Dc7wBEAN4XUFBcN0BQmQBQWABRCgBCgN0AWUBbQFDbwEISEJsdWUvQmxhY2sgQ29udHJvbA__"

	// get an actual deck to test
	deck, _ := ArtifactDeckCode.ParseDeck(deckString)

	// pass decoded deck to be parsed back into string form
	encodedDeckString, _ := ArtifactDeckCode.EncodeDeck(deck)

	assert.Equal(t, deckString, encodedDeckString)
}

func TestEncodeBytesReturnsValidBytes(t *testing.T) {

	// get an actual deck to test
	deck, _ := ArtifactDeckCode.ParseDeck("ADCJRwSJX2Dc7wBEAN4XUFBcN0BQmQBQWABRCgBCgN0AWUBbQFDbwEISEJsdWUvQmxhY2sgQ29udHJvbA__")

	_, encodedDeck := ArtifactDeckCode.EncodeBytes(deck)
	assert.Nil(t, encodedDeck)
}

func TestEncodedBytesInvalidDeck(t *testing.T) {
	deck := &ArtifactDeckCode.CardDeck{}
	deck2 := &ArtifactDeckCode.CardDeck{
		Heroes: []ArtifactDeckCode.Card{
			ArtifactDeckCode.Card{
				Id: 10,
			},
		},
	}

	_, deckBytes := ArtifactDeckCode.EncodeBytes(nil)
	_, deckBytes2 := ArtifactDeckCode.EncodeBytes(deck)
	_, deckBytes3 := ArtifactDeckCode.EncodeBytes(deck2)

	assert.NotNil(t, deckBytes)
	assert.NotNil(t, deckBytes2)
	assert.NotNil(t, deckBytes3)
}