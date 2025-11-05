# Flow

> A modern data platform for storing, managing, and automating your data workflows

Flow is a pet project that provides a visual interface for creating and managing PostgreSQL databases. Think of it as your personal data hub where you can design schemas, manage data through a spreadsheet-like interface, and automate workflows—all without writing SQL.

## What is Flow?

Flow lets you:

- **Create isolated databases** for different projects
- **Design table schemas visually** with drag-and-drop
- **Manage data** through an intuitive spreadsheet interface
- **Automate workflows** with triggers and actions
- **Collaborate** with role-based access control

Perfect for prototyping, personal projects, or learning database concepts.

## Features

### Visual Table Builder
Design your database schema without writing SQL. Drag and drop columns, set types, configure relationships—it all syncs automatically to your PostgreSQL database.

### Spreadsheet-like Data Grid
View and edit data with familiar spreadsheet controls. Sort, filter, paginate, and perform CRUD operations with a clean interface.

### Workflow Automation
Create automated workflows with visual node-based editor:
- **Triggers**: Row events, webhooks, schedules
- **Actions**: Send webhooks, emails, or mutate data
- **Real-time**: Test and monitor executions live

### Project Isolation
Each project gets its own PostgreSQL database, ensuring complete data separation and security.

### Role-Based Access
Invite team members with granular permissions: Owner, Admin, Editor, or Viewer.

## Tech Stack

**Frontend**
- Vue 3 + TypeScript
- Tailwind CSS + shadcn-vue
- TanStack Table
- Vue Flow

**Backend**
- Go 1.24+
- Chi Router
- PostgreSQL 14+
- JWT Authentication

## Getting Started

### Prerequisites

- PostgreSQL 14+
- Go 1.24+
- Node.js 18+

### Installation

1. Clone the repository
```bash
git clone https://github.com/yourusername/flow.git
cd flow
```

2. Set up the database
```bash
psql -U postgres -c "CREATE DATABASE flow_dev;"
psql -U postgres -d flow_dev -f backend/migrations/schema.sql
```

3. Configure environment
```bash
# Backend
cd backend
cp .env.example .env
# Edit .env with your database credentials

# Frontend
cd ../frontend
cp .env.example .env
```

4. Run the application
```bash
# Backend (from backend/)
go run cmd/api/main.go

# Frontend (from frontend/)
npm install
npm run dev
```

5. Open [http://localhost:5173](http://localhost:5173)

## Architecture

Flow follows a clean, layered architecture:

```
┌─────────────┐
│   Handlers  │  HTTP/WebSocket endpoints
└──────┬──────┘
       │
┌──────▼──────┐
│  Services   │  Business logic & validation
└──────┬──────┘
       │
┌──────▼──────┐
│Infrastructure│ Database operations
└─────────────┘
```

Each layer has a single responsibility, making the codebase maintainable and testable.

## Project Status

This is a personal pet project built for learning and experimentation. Core features are functional, but it's not production-ready. Contributions and feedback are welcome!

## License

MIT

---

Built with curiosity and caffeine ☕
