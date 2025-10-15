# FlowDB Frontend

A modern Vue 3 + TypeScript frontend for FlowDB, built with shadcn-vue components and Tailwind CSS.

## Features

- 🔐 **Authentication** - Login/Register with JWT tokens
- 📊 **Dashboard** - Overview of projects, tables, and workflows
- 🗂️ **Project Management** - Create and manage database projects
- 🗃️ **Table Builder** - Visual table schema designer with drag-and-drop
- 📋 **Data Grid** - Spreadsheet-like interface for data management
- 🔄 **Workflow Builder** - Visual workflow automation designer
- 📈 **Activity Log** - Real-time activity monitoring
- ⚙️ **Settings** - User preferences and configuration
- 🔌 **Real-time Updates** - WebSocket integration for live updates

## Tech Stack

- **Vue 3** - Progressive JavaScript framework
- **TypeScript** - Type-safe JavaScript
- **Vite** - Fast build tool and dev server
- **Tailwind CSS** - Utility-first CSS framework
- **shadcn-vue** - Beautiful UI components
- **Pinia** - State management
- **Vue Router** - Client-side routing
- **Lucide Vue** - Icon library

## Getting Started

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Installation

1. Install dependencies:
```bash
npm install
```

2. Create environment file:
```bash
cp .env.example .env
```

3. Update environment variables in `.env`:
```env
VITE_API_BASE_URL=http://localhost:8080/api
VITE_WS_URL=ws://localhost:8080/ws
```

4. Start development server:
```bash
npm run dev
```

5. Open http://localhost:5173 in your browser

### Building for Production

```bash
npm run build
```

The built files will be in the `dist` directory.

## Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── ui/             # shadcn-vue components
│   └── AppSidebar.vue  # Main navigation sidebar
├── lib/                # Utility libraries
│   ├── api.ts          # API client
│   └── websocket.ts    # WebSocket service
├── pages/              # Page components
│   ├── auth/           # Authentication pages
│   ├── Dashboard.vue   # Main dashboard
│   ├── Projects.vue    # Project management
│   ├── Tables.vue      # Table management
│   ├── TableBuilder.vue # Visual table designer
│   ├── DataGrid.vue    # Data management interface
│   ├── Workflows.vue   # Workflow management
│   ├── WorkflowBuilder.vue # Visual workflow designer
│   ├── Activity.vue    # Activity log
│   └── Settings.vue    # User settings
├── stores/             # Pinia stores
│   ├── auth.ts         # Authentication state
│   ├── projects.ts     # Project state
│   └── ui/             # UI state
├── types/              # TypeScript type definitions
│   ├── api.ts          # API types
│   └── sidebar.ts      # Sidebar types
├── router/             # Vue Router configuration
└── main.ts             # Application entry point
```

## Key Features

### Authentication
- JWT-based authentication
- Login/Register forms with validation
- Protected routes
- User profile management

### Project Management
- Create, edit, and delete projects
- Project overview with statistics
- Project-specific navigation

### Table Builder
- Visual table schema designer
- Drag-and-drop column management
- Real-time SQL preview
- Column type validation
- Primary key and foreign key support

### Data Grid
- Spreadsheet-like data interface
- Sortable columns
- Search and filtering
- Pagination
- Inline editing
- Row operations (create, update, delete)

### Workflow Builder
- Visual workflow designer
- Multiple trigger types (row events, scheduled, webhook)
- Multiple action types (webhook, email, database operations)
- Real-time workflow preview
- Test execution

### Activity Log
- Real-time activity monitoring
- Filterable by type and status
- Detailed activity information
- Pagination and search

### Settings
- User profile management
- Application preferences
- Security settings (password, 2FA, API keys)
- Database configuration

## API Integration

The frontend integrates with the FlowDB backend API through:

- **REST API** - For CRUD operations
- **WebSocket** - For real-time updates
- **JWT Authentication** - For secure API access

## Development

### Code Style

- Use TypeScript for type safety
- Follow Vue 3 Composition API patterns
- Use shadcn-vue components for consistency
- Apply Tailwind CSS for styling
- Write descriptive component and function names

### State Management

- Use Pinia stores for global state
- Keep component state local when possible
- Use reactive refs for reactive data
- Use computed properties for derived state

### Component Structure

```vue
<script setup lang="ts">
// Imports
// Reactive data
// Computed properties
// Methods
// Lifecycle hooks
</script>

<template>
  <!-- Template with shadcn-vue components -->
</template>
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see LICENSE file for details