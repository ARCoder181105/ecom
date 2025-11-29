# E-Commerce Backend API

A robust RESTful API built with Go (Golang) for an e-commerce platform, featuring user authentication, product management, and order processing with role-based access control.

## 🚀 Features

- **User Management**
  - User registration and authentication
  - JWT-based authorization
  - Role-based access control (Customer, Seller, Admin)
  - Secure password hashing with bcrypt

- **Product Management**
  - CRUD operations for products
  - Image upload with Cloudinary integration
  - Product search functionality
  - Pagination support
  - Stock quantity tracking

- **Order Management**
  - Place orders with multiple items
  - Order history tracking
  - Transaction-based order processing
  - Automatic stock management
  - Order status updates (Admin only)

## 🛠️ Tech Stack

- **Language**: Go 1.25
- **Web Framework**: Chi Router
- **Database**: PostgreSQL
- **Migration Tool**: Goose
- **Query Builder**: SQLC
- **Authentication**: JWT (golang-jwt/jwt)
- **Image Storage**: Cloudinary
- **Password Hashing**: bcrypt

## 📋 Prerequisites

- Go 1.25 or higher
- PostgreSQL 12+
- Goose (for database migrations)
- Cloudinary account (for image uploads)

## 🔧 Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/ARCoder181105/ecom.git
   cd ecom
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   
   Create a `.env` file in the root directory:
   ```env
   PORT=8080
   DATABASE_URL=postgresql://postgres:password@localhost:5432/ecomdb?sslmode=disable
   JWT_SECRET=your-secret-key-here
   CLOUDINARY_URL=cloudinary://api_key:api_secret@cloud_name
   FRONTEND_URL=http://localhost:5173
   ```

4. **Run database migrations**
   ```bash
   goose -dir db/migrate/migrations postgres "your-database-url" up
   ```

5. **Generate SQLC code** (if needed)
   ```bash
   sqlc generate
   ```

6. **Run the server**
   ```bash
   go run cmd/main.go
   ```

   Or use Air for hot-reloading:
   ```bash
   air
   ```

## 📁 Project Structure

```
.
├── cmd/
│   ├── main.go              # Application entry point
│   └── api/
│       └── api.go           # API server setup
├── db/
│   ├── db.go                # Database connection
│   └── migrate/
│       ├── migrations/      # SQL migration files
│       ├── queries/         # SQLC query files
│       └── sqlc/            # Generated SQLC code
├── services/
│   ├── user/                # User service handlers
│   ├── products/            # Product service handlers
│   └── orders/              # Order service handlers
├── types/
│   └── types.go             # Shared type definitions
├── utils/
│   ├── jwt.go               # JWT utilities
│   ├── middleware.go        # Authentication middleware
│   ├── cloudinary.go        # Image upload utilities
│   └── utils.go             # Common utilities
└── .env                     # Environment variables
```

## 🔐 API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/user/register` | Register new user | No |
| POST | `/api/v1/user/login` | Login user | No |
| GET | `/api/v1/user/profile` | Get user profile | Yes |

### Products

| Method | Endpoint | Description | Auth Required | Role |
|--------|----------|-------------|---------------|------|
| GET | `/api/v1/product/getAllProducts` | List all products (with pagination & search) | No | - |
| GET | `/api/v1/product/getProduct/{productId}` | Get single product | No | - |
| POST | `/api/v1/product/create` | Create new product | Yes | Seller/Admin |
| POST | `/api/v1/product/upload` | Upload product image | Yes | Any |
| PUT | `/api/v1/product/{productID}` | Update product | Yes | Owner/Admin |
| DELETE | `/api/v1/product/{productID}` | Delete product | Yes | Owner/Admin |

### Orders

| Method | Endpoint | Description | Auth Required | Role |
|--------|----------|-------------|---------------|------|
| GET | `/api/v1/orders/orders` | Get user's orders | Yes | Any |
| GET | `/api/v1/orders/orders/{orderID}` | Get order details | Yes | Owner |
| POST | `/api/v1/orders/placeOrder` | Place new order | Yes | Any |
| POST | `/api/v1/orders/updateOrderStatus` | Update order status | Yes | Admin |

## 📝 Request Examples

### Register User
```json
POST /api/v1/user/register
{
  "first_name": "John",
  "last_name": "Doe",
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepassword"
}
```

### Create Product
```json
POST /api/v1/product/create
{
  "name": "Product Name",
  "description": "Product description",
  "image": "https://cloudinary.com/image.jpg",
  "price": "29.99",
  "stock_quantity": 100
}
```

### Place Order
```json
POST /api/v1/orders/placeOrder
[
  {
    "product_id": "uuid-here",
    "quantity": 2
  },
  {
    "product_id": "uuid-here",
    "quantity": 1
  }
]
```

## 🔒 Authentication

The API uses JWT tokens stored in HTTP-only cookies. After successful login/registration, the token is automatically set in the cookie and included in subsequent requests.

**Token Expiration**: 24 hours

## 🗄️ Database Schema

### Users Table
- id (UUID, Primary Key)
- first_name, last_name, username
- email (Unique)
- password (Hashed)
- role (customer, seller, admin)
- created_at

### Products Table
- id (UUID, Primary Key)
- name, description, image
- price (Decimal)
- stock_quantity (Integer)
- user_id (Foreign Key to Users)
- created_at

### Orders Table
- id (UUID, Primary Key)
- user_id (Foreign Key to Users)
- total_price (Decimal)
- status (pending, shipped, completed, cancelled)
- created_at

### Order Items Table
- id (UUID, Primary Key)
- order_id (Foreign Key to Orders)
- product_id (Foreign Key to Products)
- quantity, price
- created_at

## 🧪 Development

### Running with Air (Hot Reload)
```bash
air
```

### Database Migrations

Create new migration:
```bash
goose -dir db/migrate/migrations create migration_name sql
```

Run migrations:
```bash
goose -dir db/migrate/migrations postgres "connection-string" up
```

Rollback:
```bash
goose -dir db/migrate/migrations postgres "connection-string" down
```

## 🚧 Error Handling

All errors return JSON responses:
```json
{
  "error": "Error message here"
}
```

Common HTTP Status Codes:
- `200`: Success
- `201`: Created
- `400`: Bad Request
- `401`: Unauthorized
- `403`: Forbidden
- `404`: Not Found
- `500`: Internal Server Error

## 🔐 Security Features

- Password hashing with bcrypt
- JWT-based authentication
- HTTP-only cookies
- CORS configuration
- Role-based access control
- SQL injection prevention (via SQLC)
- Transaction-based order processing

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is open source and available under the MIT License.

## 👤 Author

**ARCoder181105**

- GitHub: [@ARCoder181105](https://github.com/ARCoder181105)

## 🙏 Acknowledgments

- Chi Router for the excellent HTTP router
- SQLC for type-safe SQL queries
- Goose for database migrations
- Cloudinary for image hosting