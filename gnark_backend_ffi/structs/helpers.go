package structs

import (
	"encoding/hex"
	"log"

	fr_bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type Witness = uint32
type Witnesses = []Witness

func DeserializeFelt(encodedFelt string) fr_bn254.Element {
	// Decode the received felt.
	decodedFelt, err := hex.DecodeString(encodedFelt)
	if err != nil {
		log.Fatal(err)
	}

	// Deserialize the decoded felt.
	var deserializedFelt fr_bn254.Element
	deserializedFelt.SetBytes(decodedFelt)

	return deserializedFelt
}

func DeserializeFelts(encodedFelts string) fr_bn254.Vector {
	// Decode the received felts.
	decodedFelts, err := hex.DecodeString(encodedFelts)
	if err != nil {
		log.Fatal(err)
	}

	// Unpack and deserialize the decoded felts.
	var deserializedFelts fr_bn254.Vector
	deserializedFelts.UnmarshalBinary(decodedFelts)

	return deserializedFelts
}

// Samples a felt and returns the encoded felt and the non-encoded felt.
func SampleEncodedFelt() (string, fr_bn254.Element) {
	var felt fr_bn254.Element
	felt.SetRandom()

	return hex.EncodeToString(felt.Marshal()), felt
}

// Samples a felts vector and returns the encoded felts and the non-encoded felts vector.
func SampleEncodedFelts() (string, fr_bn254.Vector) {
	var felt1 fr_bn254.Element
	felt1.SetRandom()

	var felt2 fr_bn254.Element
	felt2.SetRandom()

	felts := fr_bn254.Vector{felt1, felt2}

	binaryFelts, _ := felts.MarshalBinary()

	return hex.EncodeToString(binaryFelts), felts
}
