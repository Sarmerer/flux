# ⚡ FlowDB MVP Roadmap

### ✅ Phase 1 – Core Setup (Week 1)

* **Backend (Go or Rust)**

  * User auth: signup/login/logout with JWT.
  * Basic API: CRUD for `tables` and `columns` metadata.
  * Postgres setup with migration system.
  * Seeded demo account for portfolio/demo.

* **Frontend (React or Vue)**

  * Auth pages (login/signup).
  * Dashboard layout (sidebar + topbar).
  * Page: “My Tables”.

---

### ✅ Phase 2 – Visual Database Builder (Week 2)

* **Backend**

  * API endpoints: create/delete table, add/remove column.
  * Map metadata changes → actual Postgres migrations (ALTER TABLE).
  * Enforce column properties (required, unique, etc).

* **Frontend**

  * Drag-and-drop column creator (React Flow or simple list reorder).
  * Column property editor (type, required, etc).
  * Save + sync with backend.

---

### ✅ Phase 3 – Data Management (Week 3)

* **Backend**

  * CRUD API for rows inside tables.
  * Validation (check column constraints).
  * Pagination/filtering.

* **Frontend**

  * Spreadsheet-like data grid (TanStack Table or AG Grid).
  * Add/edit/delete rows.
  * Inline validation (e.g., can’t save empty required field).

---

### ✅ Phase 4 – Workflow Builder Lite (Week 4)

* **Backend**

  * Workflow schema (JSON with nodes/edges).
  * Trigger type: *on row created*.
  * One action type: *send webhook* (simplest automation).

* **Frontend**

  * Workflow canvas (React Flow).
  * Add trigger node + action node, connect them.
  * Save workflow to backend.

---

### ✅ Phase 5 – Polish & Demo (Week 5–6)

* **Frontend polish**

  * Activity log page: “Row X created, Workflow Y ran”.
  * UI cleanup: Tailwind styling, nice dashboards.
* **Demo mode**

  * Public “Try It” account (auto-reset daily).
* **Hosting**

  * Frontend on Vercel/Netlify.
  * Backend on Render/Fly.io.
  * Postgres on Supabase/Neon.

---

# 🎨 What to Show in Portfolio

* **Video demo**: 2–3 minutes showing:

  * Create a table → add columns.
  * Enter some data in grid.
  * Create workflow: “When new row added → send webhook”.
  * Add row → workflow fires.
* **Live demo link** with test account.
* **GitHub repo**: clean commits, nice README, architecture diagram.

---

# 🚀 Stretch Goals (after MVP)

* More workflow actions: send email, update row, Slack notification.
* Roles/permissions.
* Real-time updates with WebSockets.
* Payments (Stripe sandbox).
