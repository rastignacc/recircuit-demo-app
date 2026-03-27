export type Role = 'buyer' | 'seller'
export type Condition = 'like_new' | 'excellent' | 'good' | 'fair'
export type OrderStatus = 'pending' | 'confirmed' | 'shipped' | 'delivered' | 'cancelled'

export interface User {
  id: number
  email: string
  name: string
  role: Role
  created_at: string
}

export interface AuthResponse {
  token: string
  user: User
}

export interface Category {
  id: number
  name: string
  slug: string
}

export interface Product {
  id: number
  seller_id: number
  category_id: number
  brand: string
  model: string
  condition: Condition
  price: number
  description: string
  image_url: string
  specs: Record<string, string>
  stock: number
  created_at: string
  updated_at: string
  category_name?: string
  seller_name?: string
}

export interface ProductListResponse {
  products: Product[]
  total: number
  page: number
  per_page: number
}

export interface ProductFilter {
  category_id?: number
  brand?: string
  condition?: Condition
  min_price?: number
  max_price?: number
  search?: string
  page?: number
  per_page?: number
  sort?: string
}

export interface CreateProductRequest {
  category_id: number
  brand: string
  model: string
  condition: Condition
  price: number
  description: string
  image_url: string
  specs: Record<string, string>
  stock: number
}

export interface OrderItem {
  id: number
  order_id: number
  product_id: number
  quantity: number
  unit_price: number
  product_name?: string
}

export interface Order {
  id: number
  buyer_id: number
  status: OrderStatus
  total: number
  created_at: string
  items?: OrderItem[]
}

export interface CreateOrderRequest {
  items: { product_id: number; quantity: number }[]
}

export interface SellerStats {
  total_listings: number
  total_sold: number
  total_revenue: number
}

export interface CartItem {
  product: Product
  quantity: number
}
