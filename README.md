# PreviewDock

> **Every pull request deserves its own environment.**

PreviewDock, GitHub pull request'leri için izole ve geçici preview environment'lar oluşturan açık kaynak bir self-hosted platformdur. Repository bağlayın, PR açın; PreviewDock güvenlik policy'lerini değerlendirir, doğru commit SHA'yı işler ve başarılı doğrulama sonrasında preview URL'ini yayınlar.

Bu proje açık kaynak olarak geliştirilmektedir ve **[muhammedkoca.com.tr](https://muhammedkoca.com.tr)** tarafından hazırlanmıştır.

## Neden PreviewDock?

- Pull request başına ayrı preview environment
- Exact commit SHA ile deterministic deployment planı
- GitHub webhook HMAC doğrulaması ve fork PR trust gate
- Public API'ye Docker socket açmayan restricted agent boundary
- PostgreSQL / Redis / Docker Compose tabanlı self-hosted foundation

## Mevcut durum

PreviewDock aktif geliştirme aşamasındaki bir open-source project'tir. API health endpoint'leri, webhook signature verification, deployment state transition validation, `previewdock.yml` validation, GitHub App JWT client, PostgreSQL başlangıç migration'ı, Next.js dashboard ve Compose stack hazırdır.

Gerçek GitHub App onboarding, durable queue/persistence adapter'ları, Docker build/runtime controller, reverse proxy routing, RBAC ve full end-to-end deployment workflow üzerinde çalışmalar sürmektedir. Bu nedenle mevcut sürüm production deployment için henüz önerilmez.

## Quick start

### Gereksinimler

- Docker Engine ve Docker Compose
- GitHub App (webhook testleri ve gerçek integration için)
- Node.js 22+ (web development için)
- Go current stable (API/agent development için)

```bash
cp .env.example .env
# .env içindeki değerleri güçlü, benzersiz secret'larla doldurun.
docker compose up --build
```

- Web dashboard: `http://localhost:3000`
- API health: `http://localhost:8080/healthz`

`GITHUB_WEBHOOK_SECRET` boşsa webhook endpoint'i fail-closed çalışır ve hiçbir isteği kabul etmez.

## GitHub App setup

Her self-hosted instance kendi GitHub App bilgilerini kullanır. Private key, webhook secret veya production credential repository içinde saklanmaz.

1. GitHub'da GitHub App oluşturun.
2. Yalnızca gerekli permission'ları verin: repository metadata, pull request read ve deployment/check write.
3. Webhook URL: `https://your-control-plane.example/api/v1/webhooks/github`
4. App ID, private key ve webhook secret değerlerini secret manager veya local `.env` üzerinden sağlayın.
5. `*.pem`, `.env` ve token'ları asla Git'e commit etmeyin.

## `previewdock.yml`

Başlangıç için [`previewdock.yml.example`](previewdock.yml.example) dosyasını kullanın.

```yaml
version: 1
build:
  context: .
  dockerfile: Dockerfile
service:
  container_port: 3000
health:
  path: /health
  expected_status: 200
resources:
  cpu: 1
  memory: 1024MB
preview:
  ttl: 48h
  authentication: required
```

Configuration strict validation ile kontrol edilir. Absolute path, parent traversal, geçersiz port, unsafe health path ve hatalı resource limit reddedilir.

## Security & trust model

Pull request code **untrusted** kabul edilir. Fork PR'ları otomatik privileged build veya secret erişimi alamaz; untrusted workload için ayrı hardened host ya da VM/microVM runner kullanılmalıdır. Docker container'ları tek başına hostile code için complete security boundary değildir.

## Architecture

```text
Web UI → API → PostgreSQL / Job Queue → Deployment Controller
       → Authenticated Deployment Agent → Container Runtime → Preview Router
```

## Documentation

- [Architecture](ARCHITECTURE.md)
- [Configuration](CONFIGURATION.md)
- [Deployment](DEPLOYMENT.md)
- [Security](SECURITY.md)
- [Threat model](THREAT_MODEL.md)
- [Türkçe README](README_TR.md)
- [Contributing](CONTRIBUTING.md)

## License

PreviewDock [MIT License](LICENSE) ile lisanslanmıştır.

---

Built with care by [muhammedkoca.com.tr](https://muhammedkoca.com.tr) · Open source for developers who want better pull request workflows.
