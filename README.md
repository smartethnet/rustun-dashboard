# Rustun Dashboard

Dashboard for Rustun project with Vue frontend and Golang backend.

## Project Structure

```
rustun-dashboard/
├── frontend/              # Vue + Vite + Tailwind + ElementPlus
│   ├── src/
│   │   ├── components/   # Vue components
│   │   ├── views/        # Page views
│   │   ├── router/       # Vue Router
│   │   ├── store/        # Pinia state management
│   │   ├── assets/       # Static assets
│   │   ├── api/          # API services
│   │   ├── main.js       # Entry point
│   │   └── App.vue       # Root component
│   ├── public/           # Public assets
│   ├── package.json
│   ├── vite.config.js
│   ├── tailwind.config.js
│   └── index.html
│
└── backend/              # Golang backend
    ├── cmd/
    │   └── dashboard/    # Main application
    ├── internal/
    │   ├── handler/      # HTTP handlers
    │   ├── service/      # Business logic
    │   ├── model/        # Data models
    │   └── middleware/   # HTTP middleware
    ├── pkg/              # Public packages
    ├── api/              # API definitions
    ├── go.mod
    └── config.yaml
```

## Development

### Frontend

```bash
cd frontend
npm install
npm run dev
```

### Backend

```bash
cd backend
go mod download
go run cmd/dashboard/main.go
```

