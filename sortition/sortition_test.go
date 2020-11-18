package sortition

import (
	"crypto/ed25519"
	"encoding/hex"
	"log"
	"testing"

	"github.com/yoseplee/vrf-go"
)

const message string = "this is a message"
const givenSeed1 string = "dd16afe4ff9ecb20567c4a638cef8ee276e938d8617479296936497b9f80fd70"
const givenSeed2 string = "c20c468db4a665d3d53db9f0e9f08155a8052cabdddc8326e2c3bd2d90e42fea"

func TestGetRatioFromHash(t *testing.T) {
	givenSeedInBytes, err := hex.DecodeString(givenSeed1)
	if err != nil {
		log.Printf("failed to decode seed from given seed")
	}

	privateKey := ed25519.NewKeyFromSeed(givenSeedInBytes)
	publicKey := privateKey.Public().(ed25519.PublicKey)
	_, vrfOutput, err := vrf.Prove(publicKey, privateKey, []byte(message))
	want := 0.9398921304750247
	if got := HashRatio(vrfOutput); want != got {
		log.Fatalf("incorrect calculation of ratio. want = %v , but got = %v\n", want, got)
	}
}

func TestSortition(t *testing.T) {
	totalIteration := 1000
	var privateKeySlice []ed25519.PrivateKey
	var ratioSlice []float64
	for i := 0; i < totalIteration; i++ {
		publicKey, privateKey, _ := ed25519.GenerateKey(nil)
		privateKeySlice = append(privateKeySlice, privateKey)
		_, vrfOutput, err := vrf.Prove(publicKey, privateKey, []byte(message))
		if err != nil {
			log.Fatalf("incorrect calculation of ratio: %v", err)
		}
		ratio := HashRatio(vrfOutput)
		ratioSlice = append(ratioSlice, ratio)
	}

	success := 0
	for _, ratio := range ratioSlice {
		if Sortition(ratio) {
			success++
		}
	}
	rateOfSuccess := float64(success) / float64(totalIteration)
	if !(rateOfSuccess > sortitionThreshold-0.05 && rateOfSuccess < sortitionThreshold+0.05) {
		log.Fatalf("out of bound: success rate fails\nrate of success: %f\nleft bound: %f\nright bound: %f\n", rateOfSuccess, sortitionThreshold-0.05, sortitionThreshold+0.05)
	}
}
