# envsync

> Securely sync `.env` files across your team using a shared secret store backend.

---

## Installation

```bash
go install github.com/yourorg/envsync@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourorg/envsync/releases).

---

## Usage

**Push your local `.env` to the shared store:**

```bash
envsync push --env .env --store s3://my-bucket/project
```

**Pull the latest `.env` from the shared store:**

```bash
envsync pull --store s3://my-bucket/project --out .env
```

**Initialize a new project config:**

```bash
envsync init
```

This creates an `envsync.yaml` config file in your project root, where you can define your backend, encryption key source, and environment targets.

---

## Configuration

```yaml
# envsync.yaml
store: s3://my-bucket/myproject
key_source: aws-kms
environments:
  - development
  - staging
  - production
```

---

## Supported Backends

- AWS S3
- GCS (Google Cloud Storage)
- HashiCorp Vault
- Local filesystem *(for testing)*

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)