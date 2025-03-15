import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import { UserGatewayProvider } from './contexts/useUserGateway'
import { HttpClient } from '../clients/HttpClient'

// Define types for the page modules
interface RouteModule {
  default: React.ComponentType
  loader?: () => Promise<object> | object
  action?: () => Promise<object> | object
  ErrorBoundary?: React.ComponentType
}

// Type for the import.meta.glob result
type Pages = Record<string, RouteModule>

// Load all page components
const pages: Pages = import.meta.glob('./pages/**/*.tsx', { eager: true })

const routes = []
for (const path of Object.keys(pages)) {
  const fileName = path.match(/\.\/pages\/(.*)\.tsx$/)?.[1]
  if (!fileName) {
    continue
  }

  // Convert $param to :param for dynamic routes
  const normalizedPathName = fileName.includes('$')
    ? fileName.replace('$', ':')
    : fileName.replace(/\/index/, '')

  routes.push({
    path: fileName === 'index' ? '/' : `/${normalizedPathName.toLowerCase()}`,
    Element: pages[path].default,
    loader: pages[path]?.loader,
    action: pages[path]?.action,
    ErrorBoundary: pages[path]?.ErrorBoundary,
  })
}

const router = createBrowserRouter(
  routes.map(({ Element, ErrorBoundary, ...rest }) => ({
    ...rest,
    element: <Element />,
    ...(ErrorBoundary && { errorElement: <ErrorBoundary /> }),
  }))
)

const App = () => {
  const httpClient = new HttpClient({
    baseURL: 'http://localhost:9000',
    beforeRequest: (config) => {
      const authHeader = `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0aWQiOjEyMzQ1LCJ1aWQiOjEyMzQ1fQ.eL_py7_fUcApLIG8TRg1_CShz4-LASJQfKYYBCHO9Vw`
      config.headers.Authorization = authHeader
      config.headers['Content-Type'] = 'application/json'
    },
  })
  return (
    <UserGatewayProvider httpClient={httpClient}>
      <RouterProvider router={router} />
    </UserGatewayProvider>
  )
}

export default App
