# crypto/envelope

Primitivas de cifrado para el modelo zero-knowledge de EduGo (ADR 0029). Dos capas independientes:
cifrado simétrico de datos (AES-256-GCM) y sellado asimétrico hacia un destinatario (X25519 vía NaCl box).

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/crypto/envelope
```

Solo depende de la stdlib y `golang.org/x/crypto`. Sin algoritmos caseros.

## Capa simétrica — AES-256-GCM

Cifra blobs en reposo con una DEK (Data Encryption Key) de 32 bytes. Nonce aleatorio de 12 bytes
por valor, prefijado al ciphertext; formato `nonce(12B) || ciphertext || tag(16B)`, overhead fijo de
28 bytes. Si la DEK es incorrecta o el blob fue manipulado, `Open` falla en la verificación del tag.

```go
env, err := envelope.NewEnvelope(dek)      // dek de 32 bytes
sealed, err := env.Seal(plaintext)         // nonce||ciphertext||tag
plaintext, err := env.Open(sealed)         // error si DEK mala o manipulado
n := env.Overhead()                        // 28
```

Constantes: `DEKSize` (32), `Overhead` (28). Errores: `ErrKeySize`, `ErrBlobTooShort`.

## Capa asimétrica — sellado X25519 (NaCl box anónimo)

Sella un blob (típicamente una DEK) hacia un destinatario por su clave pública; solo su privada lo abre.
El emisor no necesita identidad propia: cada sello usa un par efímero interno, así que dos sellados del
mismo dato difieren.

```go
pub, priv, err := envelope.GenerateKeyPair()   // X25519, 32B cada una
sealed, err := envelope.SealFor(pub, dek)       // cifra hacia el destinatario
dek, err := envelope.OpenWith(priv, sealed)     // solo la privada correcta abre
```

Constantes: `PublicKeySize` (32), `PrivateKeySize` (32). Errores: `ErrPublicKeySize`,
`ErrPrivateKeySize`, `ErrOpenFailed`.

### Flujo zero-knowledge (ADR 0029)

- **Subida:** el dispositivo genera una DEK, cifra sus datos con la capa simétrica y sella la DEK con
  `Ks_pub` (pública del servidor). El servidor abre la DEK con `Ks_priv` y descifra.
- **Pairing:** el servidor sella la DEK con `Kd_pub` (pública del dispositivo) para entregársela.

## Decisión de diseño

El sellado asimétrico usa `golang.org/x/crypto/nacl/box` (`SealAnonymous`/`OpenAnonymous`) en vez de
ensamblar `crypto/ecdh` + HKDF + AEAD a mano: NaCl box es una construcción X25519 auditada y estándar
que cubre exactamente el caso (sellado anónimo hacia una pública), con menos superficie de error. Ver
godoc del paquete para el detalle.

## Navegación

- [Changelog](CHANGELOG.md)

## Comandos disponibles

```bash
make build     # Compilar
make test      # Tests
make check     # Lint y validación
```
