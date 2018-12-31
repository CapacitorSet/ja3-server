ja3-server
============

A proof of concept of fingerprinting TLS clients with JA3 (part of an article [here](https://capacitorset.github.io/ja3/)).

Creates an HTTPS server which responds with the JA3 fingerprint of the client, and stores it into Redis for analytics purposes.

## Thanks/licenses

 * Original algorithm: [Salesforce](https://github.com/salesforce/ja3)
 * Golang implementation: [Remco Verhoef](https://github.com/honeytrap/honeytrap/commit/192795147948103a24d34dc06dba74eecdeb086b), copyright DutchSec, AGPL 3.
 * Golang stdlib (`crypto/tls`, `net/http`): copyright the Go authors, BSD.