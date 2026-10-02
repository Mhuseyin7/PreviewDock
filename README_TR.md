# PreviewDock

**Her pull request kendi ortamını hak eder.** PreviewDock, güvenilen pull request’ler için kısa ömürlü önizleme ortamları oluşturan, kendi sunucunuzda çalıştırılabilen bir kontrol düzlemidir.

## Başlangıç

1. `.env.example` dosyasını `.env` olarak kopyalayın ve güçlü değerler girin.
2. `docker compose up --build` komutunu çalıştırın.
3. Arayüz `http://localhost:3000`, sağlık kontrolü `http://localhost:8080/healthz` adresindedir.

GitHub webhook endpoint’i, `GITHUB_WEBHOOK_SECRET` tanımlı değilse hiçbir isteği kabul etmez. Fork’tan gelen PR’lar otomatik deploy edilmez; onay ve ayrı, güçlendirilmiş bir runner gerektirir.

## Güvenlik notu

Docker container’ı tek başına kötü niyetli PR kodu için güvenlik sınırı değildir. Fork PR’ları için ayrı bir sunucu veya VM/microVM tabanlı runner kullanılmalıdır. Ayrıntılar için `SECURITY.md` ve `THREAT_MODEL.md` dosyalarına bakın.
