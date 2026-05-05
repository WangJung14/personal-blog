# 📝 Personal Blog — Tài liệu thiết kế

## Tổng quan

Blog cá nhân viết bằng **Go**, dùng để lưu trữ và đọc bài viết. Có admin panel để tạo/sửa/xóa bài với rich text editor. Không cần đăng ký tài khoản, không có comment, không có tags — đơn giản, tập trung vào nội dung.

---

## Stack kỹ thuật

| Thành phần | Lựa chọn | Ghi chú |
|---|---|---|
| Language | Go (latest) | |
| Router | `net/http` + `chi` | Nhẹ, idiomatic Go |
| Database | **PostgreSQL** (`pgx` hoặc `database/sql` + `lib/pq`) | Đã quen dùng, dễ scale |
| Template | `html/template` (built-in) | Không cần dependency thêm |
| Rich Text Editor | **Quill.js** (CDN) | Nhẹ, đẹp, lưu HTML |
| Auth | Session-based, 1 tài khoản admin | Username/password trong `.env` |
| CSS | Viết tay, không dùng framework | Kiểm soát hoàn toàn giao diện |

---

## Cấu trúc thư mục

```
blog/
├── main.go                  # Entry point, khởi động server
├── .env                     # ADMIN_USER, ADMIN_PASS, PORT, SECRET_KEY
├── blog.db                  
├── go.mod
├── go.sum
│
├── handlers/
│   ├── public.go            # GET /           — danh sách bài
│   │                        # GET /post/{slug} — đọc bài chi tiết
│   ├── admin.go             # GET/POST /admin/login
│   │                        # GET       /admin/dashboard
│   │                        # GET/POST  /admin/post/new
│   │                        # GET/POST  /admin/post/{id}/edit
│   │                        # POST      /admin/post/{id}/delete
│   └── middleware.go        # Auth middleware cho các route /admin/*
│
├── models/
│   └── post.go              # Struct Post, các hàm CRUD
│
├── storage/
│   └── db.go           
│
├── templates/
│   ├── layout.html          # Base layout (public)
│   ├── index.html           # Danh sách bài viết
│   ├── post.html            # Trang đọc bài
│   └── admin/
│       ├── layout.html      # Base layout (admin)
│       ├── login.html       # Form đăng nhập
│       ├── dashboard.html   # Danh sách bài + nút tạo/sửa/xóa
│       └── editor.html      # Quill.js editor — tạo & sửa bài
│
└── static/
    ├── css/
    │   ├── public.css       # Style trang đọc bài
    │   └── admin.css        # Style admin panel
    └── js/
        └── editor.js        # Init Quill, submit form
```

---

## Database schema

```sql
CREATE TABLE posts (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    title      TEXT    NOT NULL,
    slug       TEXT    NOT NULL UNIQUE,   -- URL-friendly, tự generate từ title
    content    TEXT    NOT NULL,          -- HTML từ Quill.js
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

---

## Các trang & route

### Public (không cần đăng nhập)

| Route | Mô tả |
|---|---|
| `GET /` | Danh sách tất cả bài, sắp xếp mới nhất trước |
| `GET /post/{slug}` | Đọc nội dung 1 bài |
| `GET /static/*` | Serve file tĩnh (CSS, JS) |

### Admin (cần đăng nhập)

| Route | Mô tả |
|---|---|
| `GET  /admin/login` | Trang đăng nhập |
| `POST /admin/login` | Xử lý đăng nhập, tạo session |
| `GET  /admin/logout` | Xóa session, redirect về login |
| `GET  /admin` | Dashboard — danh sách bài + quản lý |
| `GET  /admin/post/new` | Form tạo bài mới (Quill editor) |
| `POST /admin/post/new` | Lưu bài mới vào DB |
| `GET  /admin/post/{id}/edit` | Form sửa bài (load nội dung cũ vào Quill) |
| `POST /admin/post/{id}/edit` | Lưu thay đổi |
| `POST /admin/post/{id}/delete` | Xóa bài |

---

## Chi tiết tính năng

### Auth
- Đọc `ADMIN_USER` và `ADMIN_PASS` từ `.env` khi start server
- Đăng nhập thành công → tạo session cookie (ký bằng `SECRET_KEY`)
- Middleware kiểm tra cookie trên mọi route `/admin/*` (trừ `/admin/login`)
- Sai mật khẩu → hiển thị lỗi, không redirect

### Tạo & sửa bài (Editor)
- Dùng **Quill.js** load qua CDN
- Toolbar: Bold, Italic, Underline, Heading, Quote, Code block, Link, Image (base64 hoặc URL), Bullet list, Ordered list
- Submit form: JS copy nội dung HTML từ Quill vào `<input type="hidden">` rồi POST lên server
- Slug tự động generate từ title (lowercase, thay space bằng `-`, bỏ ký tự đặc biệt)
- Nếu slug đã tồn tại → append `-2`, `-3`,...

### Trang đọc bài
- Render HTML từ DB trực tiếp (Quill output)
- Hiển thị title, ngày tạo, nội dung
- Link về trang chủ

### Trang chủ
- Danh sách bài: title + ngày đăng + link đọc
- Sắp xếp `created_at DESC`

---

## Biến môi trường (`.env`)

```env
PORT=8080
ADMIN_USER=admin
ADMIN_PASS=your_secure_password
SECRET_KEY=random_string_for_signing_session
```

---

## Danh sách task

### Setup
- [ ] `go mod init blog`
- [ ] Cài dependencies: `chi`, `modernc.org`, `joho/godotenv`
- [ ] Tạo cấu trúc thư mục

### Storage & Model
- [ ] Viết `storage/db.go` — kết nối postgresql chạy migration tạo bảng `posts`
- [ ] Viết `models/post.go` — struct `Post` và các hàm: `GetAll`, `GetBySlug`, `GetByID`, `Create`, `Update`, `Delete`
- [ ] Viết hàm `generateSlug(title string) string`

### Auth & Middleware
- [ ] Viết `handlers/middleware.go` — đọc session cookie, redirect về `/admin/login` nếu chưa đăng nhập
- [ ] Viết login handler (GET + POST)
- [ ] Viết logout handler

### Handler — Public
- [ ] `GET /` — query tất cả bài, render `index.html`
- [ ] `GET /post/{slug}` — query 1 bài, render `post.html`

### Handler — Admin
- [ ] `GET /admin` — dashboard, danh sách bài
- [ ] `GET /admin/post/new` + `POST /admin/post/new`
- [ ] `GET /admin/post/{id}/edit` + `POST /admin/post/{id}/edit`
- [ ] `POST /admin/post/{id}/delete`

### Templates
- [ ] `templates/layout.html` — public base (header, footer)
- [ ] `templates/index.html` — danh sách bài
- [ ] `templates/post.html` — đọc bài
- [ ] `templates/admin/layout.html` — admin base
- [ ] `templates/admin/login.html` — form login
- [ ] `templates/admin/dashboard.html` — bảng quản lý bài
- [ ] `templates/admin/editor.html` — Quill editor + form

### Static
- [ ] `static/css/public.css` — style trang đọc
- [ ] `static/css/admin.css` — style admin
- [ ] `static/js/editor.js` — init Quill, xử lý submit

### Main
- [ ] `main.go` — load `.env`, init DB, đăng ký routes, start server

### Kiểm tra & hoàn thiện
- [ ] Test tạo bài, sửa bài, xóa bài
- [ ] Test auth: truy cập `/admin` khi chưa login
- [ ] Test slug trùng
- [ ] Đảm bảo HTML từ Quill render đúng trên trang đọc bài
- [ ] Viết `Dockerfile` (tuỳ chọn, để deploy dễ hơn)
- [ ] Viết `README.md` hướng dẫn chạy local và deploy

---

## Ghi chú deploy

Vì dùng SQLite, chỉ cần copy 1 binary + file `blog.db` + thư mục `static/` + `templates/` lên server là chạy được. Phù hợp với VPS nhỏ, Railway, Fly.io hoặc Render (dùng persistent disk cho file `.db`).