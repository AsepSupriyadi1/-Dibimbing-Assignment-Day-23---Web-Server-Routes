# 🚀 Assignment Day 23 – RESTful API Project

Selamat datang di project RESTful API ini!  
Project ini dikembangkan sebagai bagian dari pembelajaran backend, menggunakan Go, GORM, dan Gin framework.

## 🧾 Fitur Utama
- Autentikasi (Login & Register)
- Proteksi route dengan JWT
- CRUD Inventaris dan Produk
- Manajemen Pesanan
- Validasi input dan transaksi database

---

## ⚙️ Cara Menjalankan Project

### 1. 📦 Clone Repo
```bash
git clone https://github.com/username/assignment-day-23.git
cd assignment-day-23
```

### 2. 🧪 Buat Database
Buat database baru di lokal dengan nama:
```
assignment-day-23
```

> ✅ Nama ini **harus sesuai** dengan `DB_NAME` di file `.env`.

### 3. 🔐 Konfigurasi `.env`
> ⚠️ File `.env` **sudah tersedia** di repo untuk kemudahan development.

Pastikan environment variable seperti berikut:
```env
DB_USER=your_db_user
DB_PASS=your_db_password
DB_NAME=assignment-day-23
DB_PORT=5432
JWT_SECRET=your_jwt_secret
```

### 4. ▶️ Jalankan Server
```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

---

## 🔍 Pengetesan API

### 1. Import Collection ke Postman
- Buka **Postman**
- Import file bernama `postman_collection.json` dari repo ini

### 2. Alur Pengetesan
1. **Login / Register**
2. Gunakan token JWT dari hasil login
3. Akses berbagai **protected route** lainnya seperti Inventori, Produk, dan Pesanan

---

## 📎 Catatan
- Semua endpoint telah disusun dengan clean architecture.
- Menggunakan GORM ORM untuk query database.
- Transactional operation untuk memastikan konsistensi data.

---

## 🙏 Terima Kasih
Terima kasih telah mencoba project ini.  
Jangan ragu untuk memberikan masukan atau kontribusi!
