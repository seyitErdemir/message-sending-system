# Message Sending System

Modern bir mesajlaşma sistemi için Go tabanlı API projesi.

## Gereksinimler 🛠

- Go 1.16+
- Docker & Docker Compose
- Git

## Hızlı Başlangıç 🚀

1. Projeyi klonlayın:
```bash
git clone <repository-url>
cd message-sending-system
```

2. Geliştirme ortamını başlatın:
```bash
docker-compose -f docker-compose.dev.yml up --build
```

Uygulama http://localhost:3000 adresinde çalışmaya başlayacaktır.

## Geliştirme Ortamı 💻

### Docker ile Çalıştırma

```bash
# Geliştirme modunda başlatma
docker-compose -f docker-compose.dev.yml up --build

# Servisleri durdurma
docker-compose down

# Logları görüntüleme
docker-compose logs -f
```

### Yerel Ortamda Çalıştırma

```bash
# Bağımlılıkları yükleme
go mod tidy

# Hot-reload ile çalıştırma (Air kullanarak)
air
```

## Proje Yapısı 📁

```
.
├── cmd/
│   └── api/          # Ana uygulama
├── pkg/
│   ├── handlers/     # HTTP handlers
│   ├── models/       # Veri modelleri
│   └── database/     # Veritabanı işlemleri
├── docker/          
└── configs/          # Konfigürasyon dosyaları
```

## Ortam Değişkenleri 🔧

Farklı ortamlar için konfigürasyon dosyaları:
- `.env.dev` - Geliştirme ortamı
- `.env.test` - Test ortamı
- `.env` - Lokal ortam

## Test 🧪

```bash
# Test ortamını ayarlama
cp .env.test .env

# Testleri çalıştırma
go test -v ./... -env=test
``` 