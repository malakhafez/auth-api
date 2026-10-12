# Auth API

> Restaurant Branch Management System — Authentication & User Management API

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat&logo=postgresql)
![JWT](https://img.shields.io/badge/JWT-HS256-000000?style=flat&logo=jsonwebtokens)

---

## 📖 Overview

This service handles **Authentication** and **User Management** for the Restaurant Branch Management System.

It provides:
- **Login / Refresh / Logout** with JWT access tokens + refresh tokens.
- **User Management** by the CEO (create Branch Managers / Supply Managers).
- **Password & Email Change** for logged-in users.
- **Brute-force protection**: account lockout + per-IP rate limiting.

Designed to run as a **standalone API** behind an API Gateway.

---

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Port `9001` available

### 1. Run Migrations

Run these **in order**:

```bash
psql -U postgres -d restaurant_db -f migrations/auth_001.sql
psql -U postgres -d restaurant_db -f migrations/auth_002.sql
psql -U postgres -d restaurant_db -f migrations/auth_003_branch_fk.sql