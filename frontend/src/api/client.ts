import axios from 'axios'
import type {
  AuthResponse,
  Category,
  CreateOrderRequest,
  CreateProductRequest,
  Order,
  Product,
  ProductFilter,
  ProductListResponse,
  SellerStats,
} from './types'

const api = axios.create({
  baseURL: '/api/v1',
  headers: { 'Content-Type': 'application/json' },
  withCredentials: true,
})

api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('user')
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    return Promise.reject(err)
  }
)

export const authApi = {
  register(email: string, password: string, name: string, role: string) {
    return api.post<AuthResponse>('/register', { email, password, name, role })
  },
  login(email: string, password: string) {
    return api.post<AuthResponse>('/login', { email, password })
  },
  logout() {
    return api.post('/logout')
  },
}

export const productApi = {
  list(filter: ProductFilter = {}) {
    const params = new URLSearchParams()
    if (filter.category_id) params.set('category_id', String(filter.category_id))
    if (filter.brand) params.set('brand', filter.brand)
    if (filter.condition) params.set('condition', filter.condition)
    if (filter.min_price) params.set('min_price', String(filter.min_price))
    if (filter.max_price) params.set('max_price', String(filter.max_price))
    if (filter.search) params.set('search', filter.search)
    if (filter.page) params.set('page', String(filter.page))
    if (filter.per_page) params.set('per_page', String(filter.per_page))
    if (filter.sort) params.set('sort', filter.sort)
    return api.get<ProductListResponse>(`/products?${params}`)
  },
  getById(id: number) {
    return api.get<Product>(`/products/${id}`)
  },
  create(data: CreateProductRequest) {
    return api.post<Product>('/products', data)
  },
  update(id: number, data: Partial<CreateProductRequest>) {
    return api.put<Product>(`/products/${id}`, data)
  },
  delete(id: number) {
    return api.delete(`/products/${id}`)
  },
}

export const categoryApi = {
  list() {
    return api.get<Category[]>('/categories')
  },
}

export const orderApi = {
  create(data: CreateOrderRequest) {
    return api.post<Order>('/orders', data)
  },
  list() {
    return api.get<Order[]>('/orders')
  },
  getById(id: number) {
    return api.get<Order>(`/orders/${id}`)
  },
}

export const sellerApi = {
  products(page = 1) {
    return api.get<ProductListResponse>(`/seller/products?page=${page}`)
  },
  stats() {
    return api.get<SellerStats>('/seller/stats')
  },
}

export default api
