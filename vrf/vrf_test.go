package vrf

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/r2ishiguro/vrf/go/vrf_ed25519"
)

const givenSeed string = "064e1cee5e9a1032691741b8d8f7a43c32ef81185a5d46ea0a1aa4accc9d24f4"

func TestGetSeed(t *testing.T) {
	givenSeedInBytes, err := hex.DecodeString(givenSeed)
	if err != nil {
		panic("failed to convert string to hex")
	}

	privateKey := ed25519.NewKeyFromSeed(givenSeedInBytes)
	seed := privateKey.Seed()
	if !bytes.Equal(seed, givenSeedInBytes) {
		panic("seed test failed: it is different from given seed")
	}
}

func TestKeyGenerationFromGivenSeed(t *testing.T) {
	wantPublicKey, wantPrivateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic("error during the key generation")
	}
	seed := wantPrivateKey.Seed()
	if err != nil {
		panic("error during the seed decoding into hex")
	}

	gotPrivateKey := ed25519.NewKeyFromSeed(seed)
	gotPublicKey := gotPrivateKey.Public().(ed25519.PublicKey)

	if !bytes.Equal(wantPrivateKey, gotPrivateKey) {
		fmt.Printf("wantPrivateKey: %v\n", wantPrivateKey)
		fmt.Printf("gotPrivateKey: %v\n", gotPrivateKey)
		panic("key generated from the seed is not identical from the original")
	}

	if !bytes.Equal(wantPublicKey, gotPublicKey) {
		fmt.Printf("wantPublicKey: %v\n", wantPublicKey)
		fmt.Printf("gotPublicKey: %v\n", gotPublicKey)
		panic("key generated from the seed is not identical from the original")
	}

}

func TestVrfProofToVerify(t *testing.T) {
	const message = "m"
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("error during the key generation | %s\n", err)
	}

	val := vrf_ed25519.ECVRF_hash_to_curve([]byte(message), publicKey)

	var hashToBytes [32]byte
	val.ToBytes(&hashToBytes)

	proof, err := vrf_ed25519.ECVRF_prove(publicKey, privateKey, []byte(message))
	if err != nil {
		log.Fatalf("error during the calculating proof | %s\n", err)
	}

	wantResult := true

	if verifyResult, err := vrf_ed25519.ECVRF_verify(publicKey, proof, []byte(message)); verifyResult != wantResult || err != nil {
		log.Fatalf("error during the verify | %s\n", err)
	}
}

func TestGetRatioForSinglePrivateKey(t *testing.T) {
	mesasge := "hello go test"
	privateKey := getKeyFromGivenSeed()
	publicKey := privateKey.Public().(ed25519.PublicKey)

	want := 0.3044079188621445 //for the input message: "hello go test"
	for i := 0; i < 10; i++ {
		val := vrf_ed25519.ECVRF_hash_to_curve([]byte(mesasge), publicKey)
		var vrfOutput [32]byte
		val.ToBytes(&vrfOutput)
		got := GetRatio(vrfOutput[:])
		if want != got {
			panic("inconsistent ratio for the same key found")
		}
	}
}

func TestDistributionOfGetRatio(t *testing.T) {
	privateKey := getKeyFromGivenSeed()
	publicKey := privateKey.Public().(ed25519.PublicKey)

	var ratios []float64
	for i := 0; i < 1000; i++ {
		message := sha256.Sum256([]byte(fmt.Sprintf("%d", rand.Int())))
		vrfOutput := vrf_ed25519.ECVRF_hash_to_curve([]byte(message[:]), publicKey)
		var vrfOutputInBytes [32]byte
		vrfOutput.ToBytes(&vrfOutputInBytes)
		got := GetRatio(vrfOutputInBytes[:])
		ratios = append(ratios, got)
	}
	fmt.Println("average is ", GetAverage(ratios))
	fmt.Println("variance is ", GetVariance(ratios))
}

func TestBasicStatistics(t *testing.T) {
	data := [5]float64{
		1.0, 3.0, 5.0, 7.0, 9.0,
	}
	want := 5.0
	got := GetAverage(data[:])

	if want != got {
		panic("incorrect calculation of average in given dataset")
	}

	want = 8.0
	got = GetVariance(data[:])

	if want != got {
		panic("incorrect calculation of variance in given dataset")
	}

}

func getKeyFromGivenSeed() ed25519.PrivateKey {
	givenSeedInBytes, err := hex.DecodeString(givenSeed)
	if err != nil {
		panic("failed to decode string to hex")
	}
	privateKey := ed25519.NewKeyFromSeed(givenSeedInBytes)
	return privateKey
}
