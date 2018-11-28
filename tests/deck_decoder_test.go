package tests

import (
	"github.com/exts/ArtifactDeckCodeGo"
	"github.com/stretchr/testify/assert"
	"testing"
)

// test if decode string works
func TestValidDecodeDeckString(t *testing.T) {
	code := "ADCJQEAZX0ivAGABwA4XSXdAUEGQgEFAmIBRF0BDAkYAQUHAwEIBSQBMQFwASgBTw__"
	_, isNil := ArtifactDeckCodeGo.DecodeDeckString(code)
	assert.Nil(t, isNil)
}

// test if code errors out if we don't have proper prefix
func TestDeckDecodeStringInvalidPrefix(t *testing.T) {
	code := "JQEAZX0ivAGABwA4XSXdAUEGQgEFAmIBRF0BDAkYAQUHAwEIBSQBMQFwASgBTw__"
	_, isNil := ArtifactDeckCodeGo.DecodeDeckString(code)
	assert.NotNil(t, isNil)
}

// test ParseDeckInternal
func TestParseDeckInternal(t *testing.T) {
	code := "ADCJQEAZX0ivAGABwA4XSXdAUEGQgEFAmIBRF0BDAkYAQUHAwEIBSQBMQFwASgBTw__"
	deckBytes, _ := ArtifactDeckCodeGo.DecodeDeckString(code)
	_, parsed := ArtifactDeckCodeGo.ParseDeckInternal(code, deckBytes)

	assert.Nil(t, parsed)
}

func TestHeroCount(t *testing.T) {
	codes := []string{
		"ADCJQEAZX0ivAGABwA4XSXdAUEGQgEFAmIBRF0BDAkYAQUHAwEIBSQBMQFwASgBTw__",
		"ADCJcURIH0De7sBKAGQeF1BQWbdAVhHRwFIMQIECG0CTgIfRlBCdQFSZWQvR3JlZW4gQnJhd2xlcg__",
		"ADCJRwSJX2Dc7wBEAN4XUFBcN0BQmQBQWABRCgBCgN0AWUBbQFDbwEISEJsdWUvQmxhY2sgQ29udHJvbA__",
	}

	for _, code := range codes {
		deck, _ := ArtifactDeckCodeGo.ParseDeck(code)
		assert.Equal(t, len(deck.Heroes), 5)
	}
}
