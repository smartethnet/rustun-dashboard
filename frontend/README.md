# Rustun Dashboard Frontend

Modern web dashboard for Rustun VPN management, built with Vue 3, Vite, Tailwind CSS, and Element Plus.

## Features

- ✅ **Modern UI** - Clean and responsive design
- ✅ **Network Topology** - Interactive visualization of VPN network structure
- ✅ **Multi-language** - Support for English and Chinese
- ✅ **Client Management** - Full CRUD operations for VPN clients
- ✅ **Real-time Stats** - Live statistics and monitoring
- ✅ **Responsive Design** - Works on desktop, tablet, and mobile

## Technology Stack

- **Vue 3** - Composition API with `<script setup>`
- **Vite** - Lightning-fast development
- **Element Plus** - Rich component library
- **Tailwind CSS** - Utility-first CSS
- **Pinia** - State management
- **Vue Router** - Navigation
- **Vue I18n** - Internationalization
- **vis-network** - Network topology visualization
- **Axios** - HTTP client

## Pages

- **Login** - Authentication page with language switcher
- **Dashboard** - Overview with statistics and quick access
- **Network Topology** - Interactive network visualization
  - Hierarchical layout showing clusters, clients, gateways, and routes
  - Zoom, pan, and fit controls
  - Real-time statistics
  - Click nodes for details
- **Clients** - List and manage all clients
  - Filter by cluster
  - Search by identity or IP
  - CRUD operations

## Quick Start

### Install Dependencies

```bash
cd frontend
npm install
```

### Development Server

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
│   ├── locales/        # i18n translations
│   │   ├── index.js    # i18n setup
│   │   ├── en.js       # English translations
│   │   └── zh.js       # Chinese translations
│   ├── router/         # Vue Router
│   │   └── index.js    # Routes and navigation guards
│   ├── store/          # Pinia stores
│   │   └── index.js    # App state management
│   ├── views/          # Page components
│   │   ├── Login.vue
│   │   ├── Dashboard.vue
│   │   ├── Topology.vue
│   │   └── Clients.vue
│   ├── App.vue         # Root component
│   └── main.js         # Application entry
├── docs/
│   └── I18N.md         # i18n documentation
├── index.html
├── vite.config.js
├── tailwind.config.js
├── postcss.config.js
└── package.json
```

## Network Topology

The topology visualization shows:

- **Clusters** - Blue boxes representing cluster groups
- **Gateways** - Orange diamonds showing network gateways
- **Clients** - Green circles for individual clients
- **Routes** - Dashed lines showing CIDR routes

### Features

- **Hierarchical Layout** - Organized from clusters down to routes
- **Interactive** - Zoom, pan, and click nodes for details
- **Controls** - Zoom in/out, fit view, toggle labels
- **Statistics** - Real-time node and connection counts
- **Legend** - Visual guide to node types

### Usage

```vue
<template>
  <div ref="networkContainer" class="network-container"></div>
</template>

<script setup>
import { Network } from 'vis-network/standalone'

const initNetwork = () => {
  const data = {
    nodes: [...],
    edges: [...]
  }
  
  const options = {
    layout: { hierarchical: { enabled: true } }
  }
  
  network.value = new Network(container, data, options)
}
</script>
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

## Multi-language Support

### Available Languages

- 🇺🇸 English (en)
- 🇨🇳 中文 (zh)

### Switch Language

- **Login page** - Click language buttons at the top
- **Main interface** - Click globe icon in header

### Add New Language

1. Create new language file in `src/locales/`
2. Import and register in `src/locales/index.js`
3. Update language switcher components

See [I18N.md](docs/I18N.md) for detailed documentation.

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
7. **i18n** - Use translation keys for all user-facing text

## Development Tips

### Hot Module Replacement

Vite provides instant HMR. Your changes will reflect immediately.

### Vue Devtools

Install Vue Devtools browser extension for debugging:
- Inspect component tree
- View Pinia state
- Track router navigation

### Console Logging

Use `console.log` for debugging, but remove before committing:

```js
// Development
console.log('Data:', data)

// Production
// Remove or use proper logging service
```

## Deployment

### Build

```bash
npm run build
```

### Serve Static Files

```bash
# Using nginx
server {
  listen 80;
  root /path/to/dist;
  
  location / {
    try_files $uri $uri/ /index.html;
  }
  
  location /api {
    proxy_pass http://localhost:8080;
  }
}
```

### Environment Variables

Create `.env.production`:

```env
VITE_API_BASE_URL=https://api.your-domain.com
```

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## Troubleshooting

### Port Already in Use

```bash
# Kill process on port 3000
lsof -ti:3000 | xargs kill -9
```

### Node Version

Requires Node.js 16+ and npm 7+

```bash
node --version  # Should be >= 16.0.0
npm --version   # Should be >= 7.0.0
```

### Clear Cache

```bash
rm -rf node_modules package-lock.json
npm install
```

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## License

MIT
