# VRF
This repository contains basic code for verifiable random function(vrf.go), and a simple selection mechanism(sortition.go).

Note that vrf implementation is originally from Yahoo's work(2017, Apache 2.0), which retrieved from [here](https://github.com/r2ishiguro/vrf/tree/master/go/vrf_ed25519).

In this repository, I modified it because 1) it was far away from go convention(I'm not good at the convention though), 2) It was not good for utilizing vrf output, which can be used for cryptographic sortition or selection mechanism in blockchain technologies.

So, here's the change log. 1) all the function names changed to be carmelCase instead of snake_case. 2) all the functions became private except Prove(), Hash(), Verify(). 3) Prove() function now returns not only proof(pi) but also vrf output so that users can easily use them without calling Hash() function.

In addition, I made a simple selection mechanism(can be called a kind of cryptographic sortition). This may help you to understand how to use vrf output. For more details, [click here]().

Any kind of contribution will be welcomed. Thanks!

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

# A simple selection mechanism
> This also called as cryptographic sortition
## 1. Calculate Random number from hash(vrf output)
* Can calculate a random ratio range in [0, 1] from vrf output which is unique for a message, and verifiable for all the others who have issuer's public key and its proof
* The Ratio can be calculated as follows:
    * ratio = hash / (2^hashlen)
* And **its probability is uniformly distributed**
> To calculate the result by yourself, just run the main function. It's ready for you! e.g. $go run .
![Probability mass when N=150000](https://github.com/yoseplee/vrf-go/blob/master/visualize/probabilityMass(n=150000).png?raw=true)

## 2. Implement cryptographic sortition
![an overview of simple cryptographic sortition](https://github.com/yoseplee/vrf-go/blob/master/simple-sortition-overview.png?raw=true)
* Now we can implement a cryptographic sortition using VRF by setting a threshold or range which can represents selection by itself
* Example
    * let's say we have set range [0, 0.1] and any ratios which value is in it means the selected one
    * Peer 'A' calculated ratio and its value is 0.03
    * Then 'A' can claim that he/she is selected and can verify it by providing the proof
    * Peer 'B' calculated ratio and its value is 0.5
    * Then 'B' cannot claim that he/she is selected as its value is out of range [0, 0.1]

## 3. Result
* In the code written in sortition, the threshold set for 0.3, which means that only participants who got the value under 0.3 will be selected
* To execute experiment to see if its expected rate of selection, test code is ready for run sortition for 1000 times and count the ratio of success
```sh
# at the root directory of the project
cd sortition/
go test
```
* As the random variable from vrf output is from the uniform distribution, expected ratio of success will be almost the same as the threshold
* if you are very lucky, you would see fail as this is probability case.

# Other VRF Implementations
1. ed25519
    * [r2ishiguro: go](https://github.com/r2ishiguro/vrf/tree/master/go/vrf_ed25519)
    * [CONIKS: go](https://github.com/coniks-sys/coniks-go/tree/master/crypto/vrf)
    * [Algorand: go](https://github.com/algorand/go-algorand/tree/master/crypto)
    * [Witnet: rust](https://github.com/witnet/vrf-rs)
2. P-256
    * [NSEC5 Project: C/C++](https://github.com/fcelda/nsec5-crypto)
    * [Google Key Transparency(EC-VRF-P256-SHA256): go](https://github.com/google/keytransparency/blob/master/core/crypto/vrf/p256/p256.go)
        * uses SHA512 instead of SHA256 
