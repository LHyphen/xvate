<h1 align="center" style="border-bottom: none;">xvate</h1>
<h3 align="center">A simple asymmetric encryption program.</h3>



## Initialization

Alice and Bob execute initialization commands respectively

```
xvate init
```

The two parties exchange public key

1. Alice sends the `public_alice` in the `/self` directory to Bob
2. Bob sends the `public_bob` in the `/self` directory to Alice
3. Alice stores `public_bob` in the `/other` directory
4. Bob stores `public_alice` in the `/other` directory
