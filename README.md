### 1. Clone Repository
```bash
  git clone git@github.com:aditya3232/ms-totp.git
  cd ms-totp
```

### 2. Install Dependency
```bash
  go mod tidy
```

### 3. Konfigurasi Env
```bash
  .env
```

### 4. Jalankan Semua Migrasi
```bash
  go run main.go migrate up
```

### 5. Rollback Migrasi
```bash
  go run main.go migrate down
```

### 6. Cek Migrasi Status
```bash
  go run main.go migrate status
```

### 7. Jalankan Seed (admin & role)
```bash
  go run main.go seed
```

### 8. Menjalankan Service
```bash
  go run main.go start
```

### 9. Menjalankan Salah Satu Test
```bash
  go test ./test -run TestSignIn_Success -v
```

### 10. Cek Coverage
```bash
  go test -coverpkg=./... ./test -coverprofile=coverage.out && go tool cover -func=coverage.out
```

### 11. Get Detail Coverage
```bash
go tool cover -func=coverage.out \
 | grep -E "user_handler.go|user_service.go" \
 | grep -E "SignIn"

```
