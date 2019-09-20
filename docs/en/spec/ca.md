# CA

QOS CA center[kepler](https://github.com/QOSGroup/kepler.git).

## Concepts

To create QSC and initialize QCP, you need to apply for CA and obtain certificate files such as `qsc.crt` and `qcp.crt`.

commands such as `qoscli tx create-qsc`、`qoscli tx init-qcp` need `qsc.crt`、`qcp.crt` files, the common struct of these files:

```go
type CertificateSigningRequest struct {
	Subj      Subject       `json:"subj"`
	IsCa      bool          `json:"is_ca"`
	NotBefore time.Time     `json:"not_before"`
	NotAfter  time.Time     `json:"not_after"`
	PublicKey crypto.PubKey `json:"public_key"`
}
```
for different ca, `Subject` are different.

### QSC
```go
type QSCSubject struct {
	ChainId string        `json:"chain_id"` // QOS mainnet chain-id
	Name    string        `json:"name"`     // QSC tokens' name
	Banker  crypto.PubKey `json:"banker"`   // banker public key, ed25519
}
```

### QCP
```go
type QCPSubject struct {
	ChainId  string `json:"chain_id"`      // QOS mainnet chain-id
	QCPChain string `json:"qcp_chain"`     // QCP chain-id
}
```

## Apply

CA Center provides rpc for developers to apply for and query.
After the application is successful, wait for the CA management personnel to review. Then the certificate will be sent by email.

### QSC

apply:
```shell script
curl -X POST "{baseUrl}/qsc/apply?qscName={QSC name}&qosChainId={QOS mainnet chain-id}&qscPub={QSC public key}&bankerPub={QSC banker public key}&phone={mobile}&email={email}&info={apply info}" -H "accept: application/json"
```
required parameters:

- qscName, QSC token name, length lte 8, capital
- qosChainId, QOS mainnet chain-id
- qscPub, QSC public key, eg.{"type":"tendermint/PubKeyEd25519","value":"VxlG2mFfbAaWSqIOuntAB3wojZwQz+sK070LC4ZBAeg="}
- bankerPub, QSC banker public key, eg.{"type":"tendermint/PubKeyEd25519","value":"VxlG2mFfbAaWSqIOuntAB3wojZwQz+sK070LC4ZBAeg="}
- phone
- email
- info, apply information

query:
```shell script
curl -X GET "{baseUrl}/qsc/apply?phone={mobile}&email={email}" -H "accept: application/json"
```

required parameters:

- phone
- email

### QCP

apply:
```shell script
curl -X POST "{baseUrl}/qcp/apply?qcpChainId={QCP chain-id}&qosChainId={QOS mainnet chain-id}&qcpPub={QCP public key}&phone={mobile}&email={email}&info={apply info}" -H "accept: application/json"
```

required parameters:

- qcpChainId, QCP chain-id
- qosChainId, QOS mainnet chain-id
- qcpPub, QCP public key, eg.{"type":"tendermint/PubKeyEd25519","value":"VxlG2mFfbAaWSqIOuntAB3wojZwQz+sK070LC4ZBAeg="}
- phone
- email
- info, apply information

query:
```shell script
curl -X GET "{baseUrl}/qcp/apply?phone={mobile}&email={email}" -H "accept: application/json"
```

required parameters:

- phone
- email