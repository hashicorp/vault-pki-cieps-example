# vault-pki-cieps-example

This repository holds an example [Certificate Issuance External Policy
Service (CIEPS)](https://developer.hashicorp.com/vault/docs/v1.15.x/secrets/pki/cieps)
implementation. This protocol allows Vault PKI operators to customize
certificate validation and templating, adding subject attributes and
extensions not natively supported by Vault.

[API Docs](https://developer.hashicorp.com/vault/api-docs/v1.15.x/secret/pki#set-certificate-issuance-external-policy-service-cieps-configuration) / [Protocol Docs](https://developer.hashicorp.com/vault/docs/v1.15.x/secrets/pki/cieps)

This service implementation responds to requests under the `/evaluate`
endpoint, but any endpoint can be chosen.

**Note**: This service only works with Vault Enterprise v1.15.0+.

---

## Building

The `Makefile` contains a default `build` target for building the binary:

```bash
$ make
go build github.com/hashicorp/vault-pki-cieps-example/cli/cieps-server
go: downloading github.com/hashicorp/vault/sdk v0.10.0
$ ls cieps-server
cieps-server
```

This builds an example service binary, `cieps-server`, which can be run for
use with Vault. The `certs` build target can be used to generate temporary
certificates which can be used for listening against `localhost`:

```bash
$ make certs
openssl req -x509 -newkey rsa:2048 -keyout server.key -out server.crt -sha256 -days 3650 -nodes -subj "/CN=localhost" -addext "subjectAltName = DNS:localhost"
.+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*.......................+..+.+..+...+.+...+.....+.............+.....+...+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*............+.........+...+...+..........+............+..+...............+...+.+.........+...+...+...........+.........+.+........+...+....+........+...+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
.......+.....+.......+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*.+...+....+......+.....+...+....+.....+.+...............+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*..+..+.........+....+..+...................+............+..+.......+........+......+.+.........+..+....+........+...+...+.+..............+.......+..+...+.......+..+...+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
-----
```

## Running

To run the service, execute the binary:

```bash
$ ./cieps-server -help
Usage of ./cieps-server:
  -listen string
    	Path to the server key file corresponding to the given certificate file (default ":443")
  -server-cert string
    	Path to the server certificate file (default "server.crt")
  -server-key string
    	Path to the server key file corresponding to the given certificate file (default "server.key")
$ ./cieps-server -server-cert server.crt -server-key server.key -listen localhost:8443
```

Note that no messages are logged until a request comes in.

To enable this service in Vault, [write the CIEPS configuration](https://developer.hashicorp.com/vault/api-docs/v1.15.x/secret/pki#set-certificate-issuance-external-policy-service-cieps-configuration)
with a PKI operator token:

```bash
$ vault write pki/config/external-policy enabled=true external_service_url=https://localhost:8443/evaluate trusted_leaf_certificate_bundle=@/path/to/server.crt
Key                                 Value
---                                 -----
enabled                             true
entity_jmespath                     n/a
external_service_last_updated       2023-09-18T11:26:36-04:00
external_service_url                https://localhost:8443/evaluate
external_service_validated          false
group_jmespath                      n/a
last_successful_request             n/a
timeout                             15000000000
trusted_ca                          n/a
trusted_leaf_certificate_bundle     -----BEGIN CERTIFICATE-----
MIIDHzCCAgegAwIBAgIUHM49XOxUgTBSZeyjToeNqINn/F8wDQYJKoZIhvcNAQEL
BQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTIzMDkxODE1MjYxNVoXDTMzMDkx
NTE1MjYxNVowFDESMBAGA1UEAwwJbG9jYWxob3N0MIIBIjANBgkqhkiG9w0BAQEF
AAOCAQ8AMIIBCgKCAQEAmkceybr18iRT1ennPlAm2uDlxT0HV4Df48fSZ4w6hHXo
RFrA9+t2zyvitFjvUHaawfCLvqDWBo7TuAzvgEuSOazTakQgryyHuHCVryx0eF7P
bpboqmLL20IHeYFlOElUsFlvYVwZechQ0F2Kz1+mNBnVkR/DAhZjydTOX++BAaka
UlfsVP1MSVmF1eD2kxv7bPvpEiQr5ABVRfX5uhHKpXfW8h1/8vcMMd9XUUzbtPOJ
HuzOhZCbmMuOMA5HmBghIS6SBnFvX4KwVnpXEoPQvlubhFsAo1czhwFnJIxa3vPz
N/mhixaIiyqaaO0DaYuFXvxLOy7JwEW9Q+ySAbF8pQIDAQABo2kwZzAdBgNVHQ4E
FgQUyKRXRvoUtaVdQqg61P2FtVoIWAEwHwYDVR0jBBgwFoAUyKRXRvoUtaVdQqg6
1P2FtVoIWAEwDwYDVR0TAQH/BAUwAwEB/zAUBgNVHREEDTALgglsb2NhbGhvc3Qw
DQYJKoZIhvcNAQELBQADggEBAHG8BUO2QY/nKKM3bxX8ZPDSphI4b8X6+TV7kQCT
5HphFSh+rIqDqi1FVFwUR6ZHGSVB/cnHapQWY4V2y6I2IRaRPkjd8ZKAWk8n8vOp
twN0aPnAtwCTJnPKE+bTINjKbGiXQiGDqELNmmFSoI97TV2ER4jczEs2kF6qUJ2V
+IG7Dppk/9qhsspQhVn2HwNWkORs7Qsubaq/2w0Wa2KDKvCbM+eEhajwYYqf3uUP
0xErpm3lS2VjzDMIi2NmXtkE5F1A6Z8dwlnoXPtvCzn4cqs39AnqKHd84CdkeCax
xCrID0z6QfCShsy85My8l7Bj9i2bGw7Aj61nsFPok2OPm+Q=
-----END CERTIFICATE-----
vault_client_cert_bundle_no_keys    n/a
```

Then, as a user, issuing via the [external-policy](https://developer.hashicorp.com/vault/api-docs/v1.15.x/secret/pki#generate-certificate-and-key-with-external-policy)
endpoints will work:

```
$ vault write pki/external-policy/issue/my-policy-name common_name="localhost" key_type=ec key_bits=256
WARNING! The following warnings were returned from Vault:

  * result from demo server; no validation occurred

  * Endpoint ignored these unrecognized parameters: [common_name]

Key                 Value
---                 -----
certificate         -----BEGIN CERTIFICATE-----
MIIBtTCCAVugAwIBAgIUNd2iuUU/Ou15EWBAcrPFPrHOn8YwCgYIKoZIzj0EAwIw
EjEQMA4GA1UEAxMHUm9vdCBSMTAeFw0yMzA5MTgxNTI2MDZaFw0yMzA5MjgxNTI2
MzZaMAAwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAT7VK20VNfwAJJ6VPrzbj5F
ATkN1GQnkYShvUn4QKeVEChMBbnQ+yOPSs7PA6F0kUNcV9ICZs2bl+SVX85EPJ6F
o4GgMIGdMB8GA1UdIwQYMBaAFI4MmL0GhyKASc6JMICpRlbCK5R9MA4GA1UdDwEB
/wQEAwIDuDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwHQYDVR0OBBYE
FCOSXnQvev/TGigJnKp5kkiDk6UZMCwGCWCGSAGG+EIBDQQfEx1DSUVQUyBEZW1v
IFNlcnZlciBDZXJ0aWZpY2F0ZTAKBggqhkjOPQQDAgNIADBFAiEA3Y3vhWSl1ovQ
6KnqM/IQoZQAkcm0Dbh+QK9AsNyFitsCIEI/STMzVISe15NJVvM+wKdpFmYQQHhj
y67geRD3GJy1
-----END CERTIFICATE-----
expiration          1695914796
issuing_ca          -----BEGIN CERTIFICATE-----
MIIBiTCCAS+gAwIBAgIUBqXN6wIae8jGxXiPCRLENE4Q+aAwCgYIKoZIzj0EAwIw
EjEQMA4GA1UEAxMHUm9vdCBSMTAeFw0yMzA5MTgxNTI0NDJaFw0yMzEwMjAxNTI1
MTJaMBIxEDAOBgNVBAMTB1Jvb3QgUjEwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNC
AARRyfTMh/zuNSco+BYVrhnJqGZSkIHjG80xe8ryW6CUhSSdsRa4CJyQGKhj/G1z
J1o/Xf2Cpf2P/kAyfU+J0dQWo2MwYTAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/
BAUwAwEB/zAdBgNVHQ4EFgQUjgyYvQaHIoBJzokwgKlGVsIrlH0wHwYDVR0jBBgw
FoAUjgyYvQaHIoBJzokwgKlGVsIrlH0wCgYIKoZIzj0EAwIDSAAwRQIhANUZrtb3
dXDCwdGT1D268aLZwi70U5RCcxmxaqwp4WU6AiBx7GbTL+Qjguuc/b0kOgEDLsxG
pqh5dIBFTWg7FFwvoA==
-----END CERTIFICATE-----
private_key         -----BEGIN EC PRIVATE KEY-----
MHcCAQEEINVIwjkWXYMtOl8UJ5NLfMuUx6tzfSYt2iLo/GsaR67AoAoGCCqGSM49
AwEHoUQDQgAE+1SttFTX8ACSelT6824+RQE5DdRkJ5GEob1J+ECnlRAoTAW50Psj
j0rOzwOhdJFDXFfSAmbNm5fklV/ORDyehQ==
-----END EC PRIVATE KEY-----
private_key_type    ec
serial_number       35:dd:a2:b9:45:3f:3a:ed:79:11:60:40:72:b3:c5:3e:b1:ce:9f:c6
```

**Note**: this service is just for example purposes only. The service should
not be used in production as-is.
