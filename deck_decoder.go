package ArtifactDeckCode

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// returns a map[string][]map[string]int or error equivalent to the php version
// data["heroes"] 	=> []map[string]int
// data["cards"] 	=> []map[string]int
// current need help with getting deck name
func ParseDeck(strDeckCode string) (CData, error) {
	deckBytes, err := DecodeDeckString(strDeckCode)
	if err != nil {
		return nil, err
	}
	return ParseDeckInternal(strDeckCode, deckBytes)
}

func DecodeDeckString(strDeckCode string) ([]byte, error) {

	// check for prefix
	if !strings.HasPrefix(strDeckCode, EncodedPrefix) {
		return nil, fmt.Errorf("Deck Code doesn't have correct prefix")
	}

	// strip prefix
	var deckCode string = strDeckCode[3:]

	// search & replace
	strFindReplace := map[string]string{
		"-": "/",
		"_": "=",
	}

	for strKey, strVal := range strFindReplace {
		deckCode = strings.Replace(deckCode, strKey, strVal, -1)
	}

	// decode string
	decode, err := base64.StdEncoding.DecodeString(deckCode)
	if err != nil {
		return nil, err
	}

	return decode, nil
}

func ParseDeckInternal(strDeckCode string, deckBytes []byte) (CData, error) {

	byteIndex := 0
	totalBytes := len(deckBytes)

	// check version num
	versionAndHeroes := deckBytes[byteIndex]
	byteIndex += 1

	version := versionAndHeroes >> 4

	if version != CurrentVersion && version != 1 {
		return nil, fmt.Errorf("Invalid Version from byte data")
	}

	// do checksum check
	checksum := deckBytes[byteIndex]
	byteIndex += 1

	strLength := 0
	if version > 1 {
		strLength = int(deckBytes[byteIndex])
		byteIndex += 1
	}

	totalCardBytes := totalBytes - strLength

	// grab string size
	{
		computedChecksum := 0
		for i := byteIndex; i < totalCardBytes; i++ {
			computedChecksum += int(deckBytes[i])
		}

		masked := byte(computedChecksum & 0xFF)
		if checksum != masked {
			return nil, errors.New("String size mismatch")
		}
	}

	//read in our hero count (part of the bits are in the version, but we can overflow bits here
	numHeroes := 0
	if !ReadVarEncodedUint32(int(versionAndHeroes), 3, deckBytes, &byteIndex, totalCardBytes, &numHeroes) {
		return nil, errors.New("Couldn't get hero count")
	}

	// now read in the heroes
	heroes := make([]map[string]int, 0)

	{
		prevCardBase := 0
		for currHero := 0; currHero < numHeroes; currHero++ {
			heroTurn := 0
			heroCardId := 0
			if !ReadSerializedCard(deckBytes, &byteIndex, totalCardBytes, &prevCardBase, &heroTurn, &heroCardId) {
				return nil, errors.New("Couldn't get hero card data")
			}

			heroes = append(heroes, map[string]int{
				"id":   heroCardId,
				"turn": heroTurn,
			})
		}
	}

	cards := make([]map[string]int, 0)

	prevCardBase := 0
	for byteIndex < totalCardBytes {
		cardId := 0
		cardCount := 0
		if !ReadSerializedCard(deckBytes, &byteIndex, totalCardBytes, &prevCardBase, &cardCount, &cardId) {
			return nil, errors.New("Couldn't get card data")
		}

		cards = append(cards, map[string]int{
			"id":    cardId,
			"count": cardCount,
		})
		println(byteIndex)
	}

	cardData := CData{
		"heroes": heroes,
		"cards":  cards,
	}

	return cardData, nil
}

//TODO remove debug code here
//handles decoding a card that was serialized
func ReadSerializedCard(data []byte, indexStart *int, indexEnd int,
	prevCardBase *int, outCount *int, outCardId *int) bool {

	//end of the memory block?
	if *indexStart > indexEnd {
		return false
	}

	//header contains the count (2 bits), a continue flag, and 5 bits of offset data. If we have 11 for the count bits we have the count
	//encoded after the offset
	header := int(data[*indexStart])
	*indexStart += 1
	hasExtendedCount := (header >> 6) == 0x03

	//read in the delta, which has 5 bits in the header, then additional bytes while the value is set
	cardDelta := 0
	if !ReadVarEncodedUint32(header, 5, data, indexStart, indexEnd, &cardDelta) {
		println("nope!")
		return false
	}

	*outCardId = *prevCardBase + cardDelta

	if hasExtendedCount {
		println("aye", indexStart)
		if !ReadVarEncodedUint32(0, 0, data, indexStart, indexEnd, outCount) {
			return false
		}
	} else {
		*outCount = (header >> 6) + 1
	}

	println("made it!")
	*prevCardBase = *outCardId

	return true
}

func ReadBitsChunk(nChunk int, nNumBits uint, nCurrShift uint, nOutBits *int) bool {
	continueBit := 1 << nNumBits
	newBits := nChunk & (continueBit - 1)
	*nOutBits |= newBits << nCurrShift

	return (nChunk & continueBit) != 0
}

func ReadVarEncodedUint32(nBaseValue int, nBaseBits uint, data []byte,
	indexStart *int, indexEnd int, outValue *int) bool {

	*outValue = 0
	var deltaShift uint = 0

	if (nBaseBits == 0) || ReadBitsChunk(nBaseValue, nBaseBits, deltaShift, outValue) {
		deltaShift += nBaseBits

		for {
			if *indexStart > indexEnd {
				return false
			}

			nextByte := int(data[*indexStart])
			*indexStart += 1

			if !ReadBitsChunk(nextByte, 7, deltaShift, outValue) {
				break
			}

			deltaShift += 7
		}
	}

	return true
}
