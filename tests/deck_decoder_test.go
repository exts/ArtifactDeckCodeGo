package tests

import (
	"github.com/exts/ArtifactDeckCodeGo"
	"github.com/stretchr/testify/assert"
	"testing"
)

// test if decode string works
func TestValidDecodeDeckString(t *testing.T) {
	code := "ADCJQEAZX0ivAGABwA4XSXdAUEGQgEFAmIBRF0BDAkYAQUHAwEIBSQBMQFwASgBTw__"
	codeRaw := "ADCJUQYL7kCQYsNCbhdgXfdAUqrAVqfpQNIRI9oARFUDVtBXSBTdHJpZmVDcm8gLSBVRyBDb21ibw"
	_, isNil := ArtifactDeckCodeGo.DecodeDeckString(code)
	_, isNil2 := ArtifactDeckCodeGo.DecodeDeckStringRaw(codeRaw)
	assert.Nil(t, isNil)
	assert.Nil(t, isNil2)
}

// test if code errors out if we don't have proper prefix
func TestDeckDecodeStringInvalidPrefix(t *testing.T) {
	code := "JQEAZX0ivAGABwA4XSXdAUEGQgEFAmIBRF0BDAkYAQUHAwEIBSQBMQFwASgBTw__"
	codeRaw := "JUQYL7kCQYsNCbhdgXfdAUqrAVqfpQNIRI9oARFUDVtBXSBTdHJpZmVDcm8gLSBVRyBDb21ibw"
	_, isNil := ArtifactDeckCodeGo.DecodeDeckString(code)
	_, isNil2 := ArtifactDeckCodeGo.DecodeDeckStringRaw(codeRaw)
	assert.NotNil(t, isNil)
	assert.NotNil(t, isNil2)
}

// test ParseDeckInternal
func TestParseDeckInternal(t *testing.T) {
	code := "ADCJQEAZX0ivAGABwA4XSXdAUEGQgEFAmIBRF0BDAkYAQUHAwEIBSQBMQFwASgBTw__"
	codeRaw := "ADCJUQYL7kCQYsNCbhdgXfdAUqrAVqfpQNIRI9oARFUDVtBXSBTdHJpZmVDcm8gLSBVRyBDb21ibw"
	deckBytes, _ := ArtifactDeckCodeGo.DecodeDeckString(code)
	deckBytesRaw, _ := ArtifactDeckCodeGo.DecodeDeckStringRaw(codeRaw)
	_, parsed := ArtifactDeckCodeGo.ParseDeckInternal(code, deckBytes)
	_, parsed2 := ArtifactDeckCodeGo.ParseDeckInternal(code, deckBytesRaw)

	assert.Nil(t, parsed)
	assert.Nil(t, parsed2)
}

func TestHeroCount(t *testing.T) {
	codes := []string{
		"ADCJQEAZX0ivAGABwA4XSXdAUEGQgEFAmIBRF0BDAkYAQUHAwEIBSQBMQFwASgBTw__",
		"ADCJcURIH0De7sBKAGQeF1BQWbdAVhHRwFIMQIECG0CTgIfRlBCdQFSZWQvR3JlZW4gQnJhd2xlcg__",
		"ADCJRwSJX2Dc7wBEAN4XUFBcN0BQmQBQWABRCgBCgN0AWUBbQFDbwEISEJsdWUvQmxhY2sgQ29udHJvbA__",
		"ADCJTgAI329uwEcCVE4XUEBJt0BR08BSEcBSEUCFxhLTQkEBTQChAQXCWAB",
	}

	for _, code := range codes {
		deck, _ := ArtifactDeckCodeGo.ParseDeck(code)
		assert.Equal(t, len(deck.Heroes), 5)
	}

	codesRaw := []string{
		"ADCJQEAZX0ivAGABwA4XSXdAUEGQgEFAmIBRF0BDAkYAQUHAwEIBSQBMQFwASgBTw",
		"ADCJcURIH0De7sBKAGQeF1BQWbdAVhHRwFIMQIECG0CTgIfRlBCdQFSZWQvR3JlZW4gQnJhd2xlcg",
		"ADCJRwSJX2Dc7wBEAN4XUFBcN0BQmQBQWABRCgBCgN0AWUBbQFDbwEISEJsdWUvQmxhY2sgQ29udHJvbA",
		"ADCJUQYL7kCQYsNCbhdgXfdAUqrAVqfpQNIRI9oARFUDVtBXSBTdHJpZmVDcm8gLSBVRyBDb21ibw",
	}

	for _, code := range codesRaw {
		deck, _ := ArtifactDeckCodeGo.ParseDeckRaw(code)
		assert.Equal(t, len(deck.Heroes), 5)
	}
}
