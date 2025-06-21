# golang_console_backup_ms_sql_server

A console-based Go application to back up tables from a Microsoft SQL Server database.  
Supports selective table backup and optional deletion after the backup completes.

---

## 🧩 Features

- 🔄 Backup individual or all tables from a specified SQL Server database
- 🗃 Save data to another MS SQL Server
- 🧹 Optionally delete tables after a successful backup
- ✅ Supports concurrent backup processing
- 🖥 Simple CLI interface for automation and scripting

---

## 🚀 Requirements

- Go 1.18 or higher
- SQL Server 2012 or later
- Access credentials with permission to read and drop tables

---

## ⚙️ Configuration

Edit the configuration in `config/config.go` or via environment variables if supported.

Key values include:
- Source server & database
- Backup target server & database
- Table selection mode
- Backup location
- Flags: `isDeleteAfterBackup`, etc.

---

## 📦 Usage

```bash
# Clone the repo
git clone https://github.com/AldythNahak/golang_console_backup_ms_sql_server.git
cd golang_console_backup_ms_sql_server

# Build or run the app
go run main.go
```
---

## 📦 Dependencies

The following Go packages are used:
```bash
go get github.com/denisenkom/go-mssqldb
go get github.com/jmoiron/sqlx
go get golang.org/x/term
```

## 📖 Goal

This repository serves as a learning log, a reference for interview prep, and a contribution to the open-source learning community.

---

## 🧑‍💻 Author

**Aldyth Nahak**  
[LinkedIn](https://linkedin.com/in/aldythnahak) | [GitHub](https://github.com/AldythNahak)

---

## ⭐️ Contribute or Follow

Feel free to fork, clone, or star this repo if you find it helpful!
