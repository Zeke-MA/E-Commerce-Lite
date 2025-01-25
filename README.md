# E-Commerce Lite

A **powerful yet simple backend** built for e-commerce using **Go**. This project provides user authentication, authorization, and an admin endpoint to manage products. Future updates will introduce additional features like caching with Redis to optimize performance.

## Table of Contents
- [Features](#features)
- [Technologies Used](#technologies-used)
- [Getting Started](#getting-started)
  - [Installation](#installation)

---

## Features

1. **User Authentication & Authorization**  
   - Secure registration and login system (JWT or session-based, depending on your implementation)  
   - Role-based access control (e.g., admin vs. regular user)

2. **Admin Management**  
   - Create, update, and delete products

3. **Checkout Process**  
   - Users can add products to cart and proceed to checkout  
   - Basic order creation flow (cart → order)

4. **RESTful API Endpoints**  
   - Essential endpoints for products, users, and orders using Go’s standard `net/http` library and Gorilla Mux

5. **Pluggable Architecture (Future)**  
   - Plans to incorporate Redis caching for improved performance  
   - Additional optimizations for larger user and product bases

---

## Technologies Used

- **Backend**: Go with `net/http` package
- **Router**: [Gorilla Mux](https://github.com/gorilla/mux)  
- **Database**: PostgreSQL  
- **Authentication**: *(JWT)* (bcrypt & golang-jwt/jwt/v5)
- **Cache (Planned)**: Redis  
- **Containerization (Database)**: Docker  

---

## Getting Started


### Installation

1. **Clone this repository**
   ```bash
   git clone https://github.com/Zeke-MA/E-Commerce-Lite.git
   cd E-Commerce-Lite
   ```

2. **Set up your environment variables**
    Create a .env file with the following values in the root of the project:
   ```bash
    PORT=PORT
    DATABASE_URL=DATABASE_URL
    POSTGRES_PASSWORD=POSTGRES_PASSWORD
    POSTGRES_USER=POSTGRES_USER
    POSTGRES_DB=POSTGRES_DB
    JWT_SECRET=JWT_SECRET
   ```  

3. **Build and Run**
    ```bash
    # Build the Go application
    go build -o ecom-lite

    # Run the binary
    ./ecom-lite
    ```