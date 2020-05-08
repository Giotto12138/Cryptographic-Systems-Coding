# The implementation is “from scratch”,

# problem1

## This practice is to implement a basic Diﬃe-Hellman cryptosystem  
the key is 1024-bit long. There are three separate programs. dh-alice1 models Alice’s initial message to Bob, and outputs a secret key to be stored for later. dh-bob models Bob’s receipt of the message from Alice, and outputs a response message back to Alice. dh-alice2 models Alice’s receipt of Bob’s response. 

```
dh-alice1 <filename for message to Bob> <filename to store secret key>. Outputs decimal-formatted ( p,g,ga ) to Bob, writes (p,g,a) to a second ﬁle.   
dh-bob <filename of message from Alice> <filename of message back to Alice>. Reads in Alice’s message, outputs ( gb ) to Alice, prints the shared secret gab.   
dh-alice2 <filename of message from Bob> <filename to read secret key>. Reads in Bob’s message and Alice’s stored secret, prints the shared secret gab.
```


# problem2

## This practice is to implement a variant of the Elgamal public-key encryption scheme with AES-GCM mode and SHA256  
The variant uses a hash function (SHA256) to compute a key for the AES-GCM encryption scheme  

 ```
elg-keygen <filename to store public key> <filename to store secret key>. This program should be identical to the dh-alice1 program.   
elg-encrypt <message text as a string with quotes> <filename of public key> <filename of ciphertext>.   
elg-decrypt <filename of ciphertext> <filename to read secret key>. 
```

# problem3

## This practice is to use brute force method and baby-step-giant-step algorithm to break Diﬃe-Hellman. 
The two programs are to ﬁnd an integer x such that gx ≡ h mod p  

```
dl-brute <filename for inputs>. On input a ﬁle containing decimal-formatted ( p,g,h ), prints x to standard output.   
dl-efficient <filename for inputs>. On input a ﬁle containing decimal-formatted ( p,g,h ), prints x to standard output.
```