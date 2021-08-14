# VRF
This repository contains basic code for verifiable random function (vrf.go), and a simple selection mechanism (sortition.go).

Note that this vrf implementation is originally from Yahoo's work (2017, Apache 2.0), which can be retrieved [here](https://github.com/r2ishiguro/vrf/tree/master/go/vrf_ed25519).

In this repository, I modified it because 1) it was far away from the go convention (though I'm not good at the convention), 2) it was not good for utilizing vrf output for cryptographic sortition or selection mechanisms in blockchain technologies.

This is the change log: 1) All the function names were changed to be carmelCase instead of snake_case. 2) All the functions became private except Prove(), Hash(), Verify(). 3) Prove() function now returns not only proof(pi) but also vrf output so that users can easily use them without calling Hash() function.

In addition, I made a simple selection mechanism (this can be called a kind of cryptographic sortition). This may help you to understand how to use vrf output. For more details, [click here](https://github.com/yoseplee/vrf#3-a-simple-selection-mechanism).

Any kind of contribution will be welcome. Thanks! 

# Appendix
## 1. Available VRF Implementations
1. ed25519
    * [r2ishiguro: go](https://github.com/r2ishiguro/vrf/tree/master/go/vrf_ed25519)
    * [CONIKS: go](https://github.com/coniks-sys/coniks-go/tree/master/crypto/vrf)
    * [Algorand: go](https://github.com/algorand/go-algorand/tree/master/crypto)
    * [Witnet: rust](https://github.com/witnet/vrf-rs)
2. P-256
    * [NSEC5 Project: C/C++](https://github.com/fcelda/nsec5-crypto)
    * [Google Key Transparency(EC-VRF-P256-SHA256): go](https://github.com/google/keytransparency/blob/master/core/crypto/vrf/p256/p256.go)
        * uses SHA512 instead of SHA256
## 2. Concept of VRF(Verifiable Random Function)
![the concept of vrf](https://github.com/yoseplee/vrf/blob/master/resources/vrf-concept.png?raw=true)
* A pseudorandom number can be verified by anyone who has a sender's public key
* A sender can generate a pseudorandom number with their private key and message
* the result (a random number) and the proof is returned and both are sent to a receiver
* The receiver can verify the number the sender generated with the sender's public key, proof, pseudorandom number, and message

### 2.1. Functions in VRF
> Generally, VRF implementation has the 3 functions below
1. Keygen (VRF_GEN): generates a key pair (secret key, public key)
2. Evaluate (VRF_EVAL): generates a pseudorandom number and its proof
3. Verify (VRF_VER): verifies the random number with proof

### 2.2. The Three Properties of VRF
> [Gorka Irazoqui Apecechea's article posted to Medium - see how it works would be great for you](https://medium.com/witnet/cryptographic-sortition-in-blockchains-the-importance-of-vrfs-ad5c20a4e018)
1. Collision resistance: it is hard to find two inputs that map to the same output
2. Pseudorandomness: the output is unidentifiable as a random number for anyone not knowing the secret key
3. Trusted Uniqueness: This requires that, given a public key, for a VRF input m corresponding to a unique output for the same input value, the result should be unique

## 3. A simple selection mechanism
> This is also called cryptographic sortition
### 3.1. Calculate Random number from hash(vrf output)
![Probability mass](https://github.com/yoseplee/vrf/blob/master/resources/sortitionProbMass.gif?raw=true)
* Can calculate a random ratio range [0, 1] from the vrf output which is unique to a message and verifiable for everyone who has the issuer's public key and its proof
* The Ratio can be calculated as follows:
    * ratio = hash / (2^hashlen)
* And **its probability is uniformly distributed**
> To calculate the result by yourself, just run the main function. It's ready for you! e.g. $go run .

### 3.2. Implement the selection mechanism
![an overview of simple cryptographic sortition](https://github.com/yoseplee/vrf/blob/master/resources/simple-sortition-overview.png?raw=true)
* Now we can implement a cryptographic sortition using VRF by setting a threshold or range which can represent a selection by itself.
* Example
    * Let's say we have set a range [0, 0.1] and any ratio whose value falls in that range is a selected value
    * Peer 'A' calculated a ratio and its value is 0.03
    * Then 'A' can claim that they have selected a value and can verify it by providing the proof
    * Peer 'B' calculated a ratio and its value is 0.5
    * Then 'B' cannot claim that they have selected a value as its value is outside of the range [0, 0.1]

### 3.3. Result
* In the code written in the sortition, the threshold was set to 0.3, which means that only participants who got a value under 0.3 will be selected
* To see if it falls in the expected selection ratio, I ran the test code, which runs sortition 1000 times and counts the ratio of success
```sh
# at the root directory of the project
cd sortition/
go test
```
* As the random variable from vrf output is from a uniform distribution, the expected ratio of success will be almost the same as the threshold
* As this is probability-based test, there is a small chance that you get lucky and see it fail.
