# Rustun Dashboard Frontend

Modern web dashboard for Rustun VPN management, built with Vue 3, Vite, Tailwind CSS, and Element Plus.

## Features

- ✅ **Modern UI** - Clean and responsive design
- ✅ **Vue 3** - Composition API with `<script setup>`
- ✅ **Vite** - Lightning-fast development
- ✅ **Element Plus** - Rich component library
- ✅ **Tailwind CSS** - Utility-first CSS
- ✅ **Pinia** - State management
- ✅ **Vue Router** - Navigation
- ✅ **Axios** - HTTP client with interceptors
- ✅ **Basic Auth** - Secure authentication

## Pages

- **Login** - Authentication page
- **Dashboard** - Overview with statistics
- **Clusters** - Manage VPN clusters
- **Cluster Detail** - View and manage cluster clients
- **Clients** - List and manage all clients

## Development

### Prerequisites

- Node.js 16+ and npm

### Install Dependencies

```bash
cd frontend
npm install
```

### Start Development Server

```bash
npm run dev
```

The app will be available at `http://localhost:3000`

### Build for Production

```bash
npm run build
```

Build output will be in `dist/` directory.

### Preview Production Build

```bash
npm run preview
```

## Project Structure

```
frontend/
├── public/              # Static assets
├── src/
│   ├── api/            # API services
│   │   └── index.js    # Axios instance and API methods
│   ├── assets/         # Styles and assets
│   │   └── style.css   # Global styles with Tailwind
│   ├── components/     # Reusable components
│   │   ├── Layout.vue     # Main layout with sidebar
│   │   └── ClientDialog.vue  # Add/Edit client dialog
│   ├── router/         # Vue Router
│   │   └── index.js    # Routes and navigation guards
│   ├── store/          # Pinia stores
│   │   └── index.js    # App state management
│   ├── views/          # Page components
│   │   ├── Login.vue
│   │   ├── Dashboard.vue
│   │   ├── Clusters.vue
│   │   ├── ClusterDetail.vue
│   │   └── Clients.vue
│   ├── App.vue         # Root component
│   └── main.js         # Application entry
├── index.html
├── vite.config.js
├── tailwind.config.js
├── postcss.config.js
└── package.json
```

## Configuration

### API Proxy

Development proxy is configured in `vite.config.js`:

```js
server: {
  port: 3000,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true
    }
  }
}
```

### Authentication

Basic Authentication credentials are stored in localStorage:
- `username` - Username
- `password` - Password

Default: `admin` / `admin123`

## API Integration

All API calls are handled through `src/api/index.js`:

```js
import { clusterAPI, clientAPI } from '@/api'

// Get all clusters
const clusters = await clusterAPI.getAll()

// Create a client
await clientAPI.create({
  cluster: 'production',
  identity: 'server-01',
  private_ip: '10.0.1.10',
  mask: '255.255.255.0',
  gateway: '10.0.1.254',
  ciders: []
})
```

## State Management

Using Pinia for state management:

```js
import { useAppStore } from '@/store'

const store = useAppStore()

// Fetch data
await store.fetchClusters()
await store.fetchClients()

// Access state
const clusters = store.clusters
const loading = store.loading
```

## Styling

### Tailwind CSS

Utility classes are available throughout:

```vue
<div class="flex items-center gap-4 p-6">
  <h1 class="text-2xl font-bold">Title</h1>
</div>
```

### Element Plus

Full component library:

```vue
<el-button type="primary" @click="handleClick">
  Click Me
</el-button>

<el-table :data="tableData">
  <el-table-column prop="name" label="Name" />
</el-table>
```

## Best Practices

1. **Composition API** - Use `<script setup>` for all components
2. **Type Safety** - Validate form inputs with Element Plus rules
3. **Error Handling** - All API calls have try-catch blocks
4. **Loading States** - Show loading indicators during async operations
5. **Responsive Design** - Mobile-friendly with Tailwind breakpoints
6. **Code Splitting** - Lazy load routes for better performance

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## License

MIT

