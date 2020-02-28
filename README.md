# Concept of VRF(Verifiable Random Function)
![the concept of vrf](https://github.com/yoseplee/vrf-go/blob/master/vrf-concept.png?raw=true)
* A pseudorandom number can be verified by anyone who has sender's public key
* A sender can generate a pseudorandom number with his/her private key and message
* its result(a random number) and the proof is returned and throw them to a receiver
* A receiver can verify the number that sender generated that pseudorandom number with (sender's public key, proof, pseudorandom number, message)

## Functions in VRF
> Generally, VRF implementation has 3 function below
1. Keygen(VRF_GEN): generates key pair(secret key, public key)
2. Evaluate(VRF_EVAL): generates pseudorandom number and its proof
3. Verify(VRF_VER): verify the random number with proof

## 3 Properties of VRF
> [Gorka Irazoqui Apecechea's article posted to Medium - see how it works would be great for you](https://medium.com/witnet/cryptographic-sortition-in-blockchains-the-importance-of-vrfs-ad5c20a4e018)
1. Collision resistance: it is hard to find two inputs that map to the same output
2. Pseudorandomness: the output is indistinguishable from random by anyone not knowing the secret key
3. Trusted Uniqueness: That requires that, given a public key, a VRF input m corresponding to a unique output for the same input value, result should be unique

# Cryptographic sortition using VRF
## 1. Calculate Random number from hash(vrf output)
* Can calculate a random ratio range in [0, 1] from vrf output which is unique for a message, and verifiable for all the others who have issuer's public key and its proof
* The Ratio can be calculated as follows:
    * ratio = hash / (2^hashlen)
* And **its probability is uniformly distributed**
> To calculate the result by yourself, just run the main function. It's ready for you! e.g. $go run .
* N=100
![Probability mass when N=100](https://github.com/yoseplee/vrf-go/blob/master/visualize/probabilityMass(n=100).png?raw=true)
* N=500
![Probability mass when N=500](https://github.com/yoseplee/vrf-go/blob/master/visualize/probabilityMass(n=500).png?raw=true)
* N=1000
![Probability mass when N=1000](https://github.com/yoseplee/vrf-go/blob/master/visualize/probabilityMass(n=1000).png?raw=true)
* N=5000
![Probability mass when N=1500](https://github.com/yoseplee/vrf-go/blob/master/visualize/probabilityMass(n=1500).png?raw=true)
* N=10000
![Probability mass when N=10000](https://github.com/yoseplee/vrf-go/blob/master/visualize/probabilityMass(n=10000).png?raw=true)
* N=15000
![Probability mass when N=150000](https://github.com/yoseplee/vrf-go/blob/master/visualize/probabilityMass(n=150000).png?raw=true)

## 2. Implement cryptographic sortition
* Now we can implement a cryptographic sortition using VRF by setting a threshold or range which can represents selection by itself
* Example
    * let's say we have set range [0, 0.1] and any ratios which value is in it means the selected one
    * Peer 'A' calaulated ratio and its value is 0.03
    * Then 'A' can claim that he/she is selected and can verify it by providing the proof
    * Peer 'B' calcuated ratio and its value is 0.5
    * Then 'B' can not claim that he/she is selected as its value is out of range [0, 0.1]

## 3. Result
TBU

# Possible Implementations
1. ed25519
    * [r2ishiguro: go](https://github.com/r2ishiguro/vrf/tree/master/go/vrf_ed25519)
    * [CONIKS: go](https://github.com/coniks-sys/coniks-go/tree/master/crypto/vrf)
    * [Algorand: go](https://github.com/algorand/go-algorand/tree/master/crypto)
    * [Witnet: rust](https://github.com/witnet/vrf-rs)
2. P-256
    * [NSEC5 Project: C/C++](https://github.com/fcelda/nsec5-crypto)
    * [Google Key Transparency(EC-VRF-P256-SHA256): go](https://github.com/google/keytransparency/blob/master/core/crypto/vrf/p256/p256.go)
        * uses SHA512 instead of SHA256 


# Goal
1. Utilize ed25519 based VRF
2. Generate Secret Key(private key) and Public key
3. Generate Random hash value
4. Verify random hash value generated at 2