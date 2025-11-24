# Neonex CLI Tools

CLI สำหรับจัดการ Neonex Core projects และ modules

## การติดตั้ง

```bash
# Build CLI
go build -o neonex.exe ./cmd/neonex

# หรือติดตั้งแบบ global
go install ./cmd/neonex
```

## คำสั่งที่มี

### 1. สร้าง Project ใหม่

```bash
neonex new <project-name>
```

สร้างโครงสร้าง project พร้อม:
- โครงสร้างโฟลเดอร์ครบถ้วน
- ไฟล์ main.go
- go.mod
- .gitignore
- README.md
- .env.example
- .air.toml (สำหรับ hot reload)
- run.bat และ run.ps1

**ตัวอย่าง:**
```bash
neonex new my-api
cd my-api
go mod download
```

### 2. รัน Development Server

```bash
# รันแบบธรรมดา
neonex serve

# รันด้วย hot reload
neonex serve --hot

# กำหนด port
neonex serve --port 3000
```

**Options:**
- `-w, --hot` - เปิดใช้ hot reload ด้วย Air
- `-p, --port` - กำหนด port (default: 8080)

### 3. จัดการ Modules

#### สร้าง Module ใหม่

```bash
neonex module create <module-name>
```

สร้างไฟล์:
- `module.go` - Module definition
- `model.go` - GORM model
- `repository.go` - Data access layer
- `service.go` - Business logic
- `controller.go` - HTTP handlers
- `routes.go` - Route definitions
- `di.go` - Dependency injection
- `seeder.go` - Seed data
- `module.json` - Module metadata

**ตัวอย่าง:**
```bash
neonex module create product
```

จากนั้นต้อง register ใน `main.go`:
```go
import "yourproject/modules/product"

core.ModuleMap["product"] = func() core.Module { return product.New() }
app.RegisterModels(&product.Product{})
```

#### แสดง Modules ทั้งหมด

```bash
neonex module list
```

แสดงรายการ modules พร้อมสถานะ

## ตัวอย่างการใช้งาน

### สร้าง Project และ Module

```bash
# 1. สร้าง project
neonex new shop-api
cd shop-api

# 2. ติดตั้ง dependencies
go mod download

# 3. สร้าง modules
neonex module create product
neonex module create category
neonex module create order

# 4. แก้ไข main.go ให้ register modules

# 5. รัน server
neonex serve --hot
```

### Module ที่สร้างจะมีโครงสร้าง

```
modules/product/
├── product.go      # Module entry point
├── model.go        # Product model
├── repository.go   # Data access
├── service.go      # Business logic
├── controller.go   # HTTP handlers
├── routes.go       # API routes
├── di.go          # Dependency injection
├── seeder.go      # Sample data
└── module.json    # Metadata
```

### API Routes ที่ได้

```
GET    /product/         - Get all
GET    /product/:id      - Get by ID
POST   /product/         - Create
PUT    /product/:id      - Update
DELETE /product/:id      - Delete
GET    /product/search   - Search
```

## Hot Reload

ต้องติดตั้ง Air ก่อนใช้งาน:

```bash
go install github.com/cosmtrek/air@latest
```

จากนั้นใช้:
```bash
neonex serve --hot
```

หรือรัน Air โดยตรง:
```bash
air
```

## Version

```bash
neonex --version
```

## Help

```bash
neonex --help
neonex <command> --help
```
