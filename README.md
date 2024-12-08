# Backend Assignment

This project is a Go-based web application that uses Fiber for the web framework, GORM for ORM, Firebase Storage for image storage, and PostgreSQL database for data storage.

## Prerequisites

- Go 1.16 or later
- Firebase project with Storage enabled
- Firebase service account key JSON file
- PostgreSQL database

## Setup

### 1. Clone the Repository

```sh
git clone https://github.com/srahul099/Product-Management-System
cd Product-Management-System
```

### 2. Install Dependencies

```sh
go mod tidy
```

### 3. Set Up Environment Variables

Create a .env file in the root directory of the project and add the following environment variables:

```sh
DB_HOST=your_db_host
DB_PORT=your_db_port
DB_USER=your_db_user
DB_PASS=your_db_password
DB_NAME=your_db_name
DB_SSLMODE=disable
```

### 4. Set Up Firebase

1. Go to the Firebase Console.
2. Create a new project or select an existing project.
3. Enable Firebase Storage.
4. Generate a service account key JSON file and download it.
5. Place the serviceAccountKey.json file in the root directory of the project.

### 6. Run the Application

```sh
go run main.go
```

## Endpoint

### 1. Create User

- URL: `/api/create_user`
- Method: `POST`
- Body:
  ```javascript
  {
  "username": "example_user"
  }
  ```

### 2. Create product

- URL: `/api/create_products`
- Method: `POST`
- Body:
  ```javascript
  {
  "user_id": 1,
  "product_name": "example_product",
  "product_description": "This is an example product.",
  "product_images": ["https://example.com/image1.jpg", "https://example.com/image2.jpg"],
  "product_price": 100.0
  }
  ```

### 3. Get All Products

- URL: `/api/products`
- Method: `GET`

### 4. Get Product by ID

- URL: `/api/products/:id`
- Method: `GET`
