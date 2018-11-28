package ArtifactDeckCodeGo

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
)

func EncodeDeck(deck *CardDeck) (string, error) {
	if deck == nil {
		return "", errors.New("Provided an invalid card deck")
	}

	bytes, err := EncodeBytes(deck)
	if err != nil {
		return "", err
	}

	deckCode, err := EncodeBytesToString(bytes)
	if err != nil {
		return "", err
	}

	return *deckCode, nil
}

func EncodeBytes(deckContents *CardDeck) ([]byte, error) {
	if deckContents == nil || len(deckContents.Heroes) <= 0 || len(deckContents.Cards) <= 0 {
		return nil, errors.New("Provided an invalid card deck")
	}

	SortCardsById(deckContents.Cards)
	SortCardsById(deckContents.Heroes)

	countHeroes := len(deckContents.Heroes)
	allCards := append(deckContents.Heroes, deckContents.Cards...)

	bytes := make([]byte, 0)

	// our version and hero count
	version := CurrentVersion<<4 | ExtraNBitsWithCarry(countHeroes, 3)
	if !AddByte(&bytes, version) {
		return nil, errors.New("Couldn't add version to byte array")
	}

	//checksum white will be updated at the end
	checksumByte := len(bytes)
	if !AddByte(&bytes, 0) {
		return nil, errors.New("Couldn't add checksum to byte array")
	}

	//write the name size
	namelen := 0
	name := deckContents.Name

	if len(deckContents.Name) > 0 {
		actualLength := len(deckContents.Name)
		for actualLength > 63 {
			trimAmount := int(math.Floor(float64(actualLength - 63/4)))
			if trimAmount <= 1 {
				trimAmount = 1
			}

			name = name[:(len(name) - int(trimAmount))]
			actualLength = len(name)
		}

		namelen = len(name)
	}

	if !AddByte(&bytes, namelen) {
		return nil, errors.New("Couldn't add name length to byte array")
	}

	if !AddRemainingNumberToBuffer(countHeroes, 3, &bytes) {
		return nil, errors.New("Couldn't add the remaning hero count to byte array")
	}

	prevCardId := 0
	for currHero := 0; currHero < countHeroes; currHero++ {
		card := allCards[currHero]
		if card.Turn <= 0 {
			return nil, errors.New("Invalid hero card, first 5 cards should be heroes with id/turn data")
		}

		if !AddCardToBuffer(card.Turn, card.Id-prevCardId, &bytes) {
			return nil, errors.New("Issue adding hero card to byte array")
		}

		prevCardId = card.Id
	}

	//reset card offset
	prevCardId = 0

	// now add all of the cards
	for currCard := countHeroes; currCard < len(allCards); currCard++ {
		card := allCards[currCard]
		if card.Count <= 0 {
			return nil, errors.New("Invalid card, card count should be at least 1 or more")
		}

		if card.Id <= 0 {
			return nil, errors.New("Invalid card, card id missing")
		}

		// record this set of cards and advance
		if !AddCardToBuffer(card.Count, card.Id-prevCardId, &bytes) {
			return nil, errors.New("Issue adding card to byte array")
		}

		prevCardId = card.Id
	}

	//save off the pre string bytes for the checksum
	stringByteCount := len(bytes)

	// write the string
	nameBytes := []byte(name)
	for nameByte := range nameBytes {
		if !AddByte(&bytes, int(nameBytes[nameByte])) {
			return nil, errors.New("Issue adding name bytes to byte array")
		}
	}

	fullChecksum := ComputeCheckSum(&bytes, stringByteCount-HeaderSize)
	smallChecksum := fullChecksum & 0x0FF

	bytes[checksumByte] = byte(smallChecksum)

	return bytes, nil
}

func EncodeBytesToString(bytes []byte) (*string, error) {
	byteCount := len(bytes)

	//if we have an empty buffer exit
	if byteCount == 0 {
		return nil, errors.New("Empty buffer provided")
	}

	encoded := base64.StdEncoding.EncodeToString(bytes)
	deckString := fmt.Sprintf("%s%s", EncodedPrefix, encoded)

	// search & replace
	strFindReplace := map[string]string{
		"/": "-",
		"=": "_",
	}

	var deckCode string
	for strKey, strVal := range strFindReplace {
		deckCode = strings.Replace(deckString, strKey, strVal, -1)
	}

	return &deckCode, nil
}

func SortCardsById(cards []Card) {
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Id < cards[j].Id
	})
}

func ExtraNBitsWithCarry(value int, bits uint) int {
	unLimitBit := 1 << bits
	unResult := value & (unLimitBit - 1)
	if value >= unLimitBit {
		unResult |= unLimitBit
	}
	return unResult
}

func AddByte(bytes *[]byte, added int) bool {
	if added > 255 {
		return false
	}
	*bytes = append(*bytes, byte(added))
	return true
}

//utility to write the rest of a number into a buffer. This will first strip the specified N bits
//off, and then write a series of bytes of the structure of 1 overflow bit and 7 data bits
func AddRemainingNumberToBuffer(value int, alreadyWrittenBits uint, bytes *[]byte) bool {
	value >>= alreadyWrittenBits
	numBytes := 0
	for value > 0 {
		nextByte := ExtraNBitsWithCarry(value, 7)
		value >>= 7
		if !AddByte(bytes, nextByte) {
			return false
		}
		numBytes += 1
	}
	return true
}

func AddCardToBuffer(count int, value int, bytes *[]byte) bool {
	if count == 0 {
		return false
	}

	countBytesStart := len(*bytes)

	//determine our count. We can only store 2 bits, and we know the value is at least one,
	// so we can encode values 1-5. However, we set both bits to indicate an extended count encoding
	firstByteMaxCount := 0x03
	extendedCount := count-1 >= firstByteMaxCount

	//determine our first byte, which contains our count, a continue flag, and the first few bits of our value
	firstByteCount := count - 1
	if extendedCount {
		firstByteCount = firstByteMaxCount
	}

	firstByte := firstByteCount << 6
	firstByte |= ExtraNBitsWithCarry(value, 5)

	if !AddByte(bytes, firstByte) {
		return false
	}

	if !AddRemainingNumberToBuffer(value, 5, bytes) {
		return false
	}

	countBytesEnd := len(*bytes)

	// check if something went horribly wrong
	if countBytesEnd-countBytesStart > 11 {
		return false
	}

	return true
}

func ComputeCheckSum(bytes *[]byte, numBytes int) int {
	checksum := 0
	for addCheck := HeaderSize; addCheck < (numBytes + HeaderSize); addCheck++ {
		cByte := byte((*bytes)[addCheck])
		checksum += int(cByte)
	}
	return checksum
}
