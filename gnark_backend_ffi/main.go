package main

import "C"
import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/recoilme/btreeset"

	fr_bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	cs_bn254 "github.com/consensys/gnark/constraint/bn254"
)

// TODO: Deserialize rawR1CS.

type RawR1CS struct {
	Gates          []RawGate
	PublicInputs   btreeset.BTreeSet
	Values         fr_bn254.Vector
	NumVariables   uint
	NumConstraints uint
}

type RawGate struct {
	MulTerms     []MulTerm
	AddTerms     []AddTerm
	ConstantTerm fr_bn254.Element
}

type MulTerm struct {
	Coefficient  fr_bn254.Element
	Multiplicand uint32
	Multiplier   uint32
}

type AddTerm struct {
	Coefficient fr_bn254.Element
	Sum         uint32
}

//export ProveWithMeta
func ProveWithMeta(rawR1CS string) *C.char {
	// Deserialize rawR1CS.
	var r RawR1CS
	err := json.Unmarshal([]byte(rawR1CS), &r)
	if err != nil {
		log.Fatal(err)
	}

	// Create R1CS.
	r1cs := cs_bn254.NewR1CS(0)

	// Add variables.
	witness, err := witness.New(r1cs.CurveID().ScalarField())
	if err != nil {
		log.Fatal(err)
	}
	witness.Fill(0, 0, nil)

	// Add constraints.

	// Setup.
	pk, _, err := groth16.Setup(r1cs)
	if err != nil {
		log.Fatal(err)
	}

	// Prove.
	proof, err := groth16.Prove(r1cs, pk, witness)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Proof: ", proof)

	// Serialize proof
	var serialized_proof bytes.Buffer
	proof.WriteTo(&serialized_proof)
	proof_string := serialized_proof.String()

	return C.CString(proof_string)
}

//export ProveWithPK
func ProveWithPK(rawR1CS string, provingKey string) *C.char {
	// Create R1CS.
	r1cs := cs_bn254.NewR1CS(1)

	// Add variables.
	witness, err := witness.New(r1cs.CurveID().ScalarField())
	if err != nil {
		log.Fatal(err)
	}
	witness.Fill(0, 0, nil)

	// Add constraints.

	// Deserialize proving key.
	pk := groth16.NewProvingKey(r1cs.CurveID())
	_, err = pk.ReadFrom(bytes.NewReader([]byte(provingKey)))
	if err != nil {
		log.Fatal(err)
	}

	// Prove.
	proof, err := groth16.Prove(r1cs, pk, witness)
	if err != nil {
		log.Fatal(err)
	}

	// Serialize proof
	var serialized_proof bytes.Buffer
	proof.WriteTo(&serialized_proof)
	proof_string := serialized_proof.String()

	return C.CString(proof_string)
}

//export VerifyWithMeta
func VerifyWithMeta(rawr1cs string, proof string) bool {
	// Create R1CS.
	r1cs := cs_bn254.NewR1CS(1)

	// Add variables.
	witness, err := witness.New(r1cs.CurveID().ScalarField())
	if err != nil {
		log.Fatal(err)
	}
	witness.Fill(0, 0, nil)

	// Add constraints.

	// Deserialize proof.
	p := groth16.NewProof(r1cs.CurveID())
	_, err = p.ReadFrom(bytes.NewReader([]byte(proof)))
	if err != nil {
		log.Fatal(err)
	}

	// Setup.
	_, vk, err := groth16.Setup(r1cs)
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve public inputs.
	publicInputs, err := witness.Public()
	if err != nil {
		log.Fatal(err)
	}

	// Verify.
	if groth16.Verify(p, vk, publicInputs) != nil {
		return false
	}

	return true
}

//export VerifyWithVK
func VerifyWithVK(rawr1cs string, proof string, verifyingKey string) bool {
	// Create R1CS.
	r1cs := cs_bn254.NewR1CS(1)

	// Add variables.
	witness, err := witness.New(r1cs.CurveID().ScalarField())
	if err != nil {
		log.Fatal(err)
	}
	witness.Fill(0, 0, nil)

	// Add constraints.

	// Deserialize proof.
	p := groth16.NewProof(r1cs.CurveID())
	_, err = p.ReadFrom(bytes.NewReader([]byte(proof)))
	if err != nil {
		log.Fatal(err)
	}

	// Deserialize verifying key.
	vk := groth16.NewVerifyingKey(r1cs.CurveID())
	_, err = vk.ReadFrom(bytes.NewReader([]byte(verifyingKey)))
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve public inputs.
	publicInputs, err := witness.Public()
	if err != nil {
		log.Fatal(err)
	}

	// Verify.
	if groth16.Verify(p, vk, publicInputs) != nil {
		return false
	}

	return true
}

//export Preprocess
func Preprocess(rawR1CS string) (*C.char, *C.char) {
	// Create R1CS.
	r1cs := cs_bn254.NewR1CS(1)

	// Add variables.
	witness, err := witness.New(r1cs.CurveID().ScalarField())
	if err != nil {
		log.Fatal(err)
	}
	witness.Fill(0, 0, nil)

	// Add constraints.

	// Setup.
	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		log.Fatal(err)
	}

	// Serialize proving key.
	var serialized_pk bytes.Buffer
	pk.WriteTo(&serialized_pk)
	pk_string := serialized_pk.String()

	// Serialize verifying key.
	var serialized_vk bytes.Buffer
	vk.WriteTo(&serialized_vk)
	vk_string := serialized_vk.String()

	return C.CString(pk_string), C.CString(vk_string)
}

//export TestFeltSerialization
func TestFeltSerialization(encodedFelt string) *C.char {
	// Decode the received felt.
	decodedFelt, err := hex.DecodeString(encodedFelt)
	if err != nil {
		log.Fatal(err)
	}

	// Deserialize the decoded felt.
	var deserializedFelt fr_bn254.Element
	deserializedFelt.SetBytes(decodedFelt)
	fmt.Printf("| GO |\n%v\n", deserializedFelt)

	// Serialize the felt.
	serializedFelt := deserializedFelt.Bytes()

	// Encode the serialized felt.
	serializedFeltString := hex.EncodeToString(serializedFelt[:])

	return C.CString(serializedFeltString)
}

//export TestFeltsSerialization
func TestFeltsSerialization(encodedFelts string) *C.char {
	// Decode the received felts.
	decodedFelts, err := hex.DecodeString(encodedFelts)
	if err != nil {
		log.Fatal(err)
	}

	// Unpack and deserialize the decoded felts.
	var deserializedFelts fr_bn254.Vector
	deserializedFelts.UnmarshalBinary(decodedFelts)

	// Serialize the felt.
	serializedFelts, err := deserializedFelts.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}

	// Encode the serialized felt.
	serializedFeltsString := hex.EncodeToString(serializedFelts[:])

	return C.CString(serializedFeltsString)
}

//export TestU64Serialization
func TestU64Serialization(number uint64) uint64 {
	fmt.Println(number)
	return number
}

//export TestMulTermSerialization
func TestMulTermSerialization(encodedMulTerm string) *C.char {
	return C.CString("unimplemented")
}

//export TestMulTermsSerialization
func TestMulTermsSerialization(encodedMulTerms string) *C.char {
	return C.CString("unimplemented")
}

//export TestAddTermSerialization
func TestAddTermSerialization(encodedAddTerm string) *C.char {
	return C.CString("unimplemented")
}

//export TestAddTermsSerialization
func TestAddTermsSerialization(encodedAddTerms string) *C.char {
	return C.CString("unimplemented")
}

//export TestRawGateSerialization
func TestRawGateSerialization(encodedRawGate string) *C.char {
	return C.CString("unimplemented")
}

//export TestRawGatesSerialization
func TestRawGatesSerialization(encodedRawGates string) *C.char {
	return C.CString("unimplemented")
}

//export TestRawR1CSSerialization
func TestRawR1CSSerialization(encodedR1CS string) *C.char {
	return C.CString("unimplemented")
}

func main() {}
