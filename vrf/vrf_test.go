package vrf

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
	"testing"

	"github.com/r2ishiguro/vrf/go/vrf_ed25519"
)

const givenSeed string = "064e1cee5e9a1032691741b8d8f7a43c32ef81185a5d46ea0a1aa4accc9d24f4"

func TestGetSeed(t *testing.T) {
	_, secretKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic("failed to generate a new key pair")
	}
	seed := secretKey.Seed()
	fmt.Println(seed)
	fmt.Printf("%x\n", seed)
}

func TestKeyGenerationFromGivenSeed(t *testing.T) {
	wantPublicKey, wantPrivateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic("error during the key generation")
	}
	// seed, err := hex.DecodeString(string(wantPrivateKey.Seed()))
	seed := wantPrivateKey.Seed()
	if err != nil {
		panic("error during the seed decoding into hex")
	}
	gotPrivateKey := ed25519.NewKeyFromSeed(seed)
	gotPublicKey := []byte(fmt.Sprintf("%v", gotPrivateKey.Public()))

	if bytes.Compare(wantPublicKey, gotPublicKey) == 1 || bytes.Compare(wantPrivateKey, gotPrivateKey) == 1 {
		fmt.Printf("wantPublicKey: %v\n", wantPublicKey)
		fmt.Printf("gotPublicKey: %v\n", gotPublicKey)
		fmt.Printf("wantPrivateKey: %v\n", wantPrivateKey)
		fmt.Printf("gotPrivateKey: %v\n", gotPrivateKey)
		panic("key generated from the seed is not identical from the original")
	}

}

func TestKeyGenerationBySeed(t *testing.T) {
	givenSeedInByteArray, _ := hex.DecodeString(givenSeed)
	privateKey := ed25519.NewKeyFromSeed(givenSeedInByteArray)
	publicKey := []byte(fmt.Sprintf("%v", privateKey.Public()))
	val := vrf_ed25519.ECVRF_hash_to_curve([]byte("message"), publicKey)
	fmt.Println(val)
}

func TestVrfProofToVerify(t *testing.T) {
	const message = "m"
	publicKey, secretKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("error during the key generation | %s\n", err)
	}
	fmt.Printf("skey %x\n", secretKey)

	val := vrf_ed25519.ECVRF_hash_to_curve([]byte(message), publicKey)
	fmt.Printf("hash: %x\n", val)
	var hashToBytes [32]byte
	val.ToBytes(&hashToBytes)
	fmt.Printf("hash: %x\n", hashToBytes)

	proof, err := vrf_ed25519.ECVRF_prove(publicKey, secretKey, []byte(message))
	if err != nil {
		log.Fatalf("error during the calculating proof | %s\n", err)
	}
	fmt.Println("proof: ", proof)

	verifyResult, err := vrf_ed25519.ECVRF_verify(publicKey, proof, []byte(message))
	if err != nil {
		log.Fatalf("error during the verify | %s\n", err)
	}
	fmt.Println("Result: ", verifyResult)
}

/*
func TestGetRatio(t *testing.T) {
	mesasge := "hello go test"
	var privateKey ed25519.PrivateKey
	var publicKey ed25519.PublicKey
	privateKey = []byte(givenPrivateKey)
	publicKey = []byte(privateKey.Public())
	val := vrf_ed25519.ECVRF_hash_to_curve([]byte(message), publicKey)
}
*/
