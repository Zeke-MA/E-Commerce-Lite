# E-Commerce Lite

A **powerful yet simple backend** built for e-commerce using **Go**. This project provides user authentication, authorization, and an admin portal to manage products and users. Future updates will introduce additional features like caching with Redis to optimize performance.

## Table of Contents
- [Features](#features)
- [Technologies Used](#technologies-used)
- [Getting Started](#getting-started)
  - [Installation](#installation)
- [Usage](#usage)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

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

- **Backend**: Go (Golang) with the built-in `net/http` standard library  
- **Router**: [Gorilla Mux](https://github.com/gorilla/mux)  
- **Database**: PostgreSQL  
- **Authentication**: *(JWT)*
- **Cache (Planned)**: Redis  
- **Containerization (Optional)**: Docker  

---

## Getting Started


### Installation

1. **Clone this repository**
   ```bash
   git clone https://github.com/Zeke-MA/E-Commerce-Lite.git
   cd E-Commerce-Lite
   ```

2. **Set up your environment variables**
    Create a .env file or export environment variables in your shell (e.g., .bashrc, .zshrc, or Windows Environment Variables). Below is an example using a .env file in the project root:
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
    go build -o ecomm-lite

    # Run the binary
    ./ecomm-lite
    ```