# MyTheresa Backend Application

This is a backend application for managing products, categories, and discounts using SQLite and GORM, and Swagger for documentation.

## Features
- Product Management:
  - Create 
  - Get product
  - List products with discounts applied
- Category Management:
  - Create category
- Discount Rules:
  - Create discount types
  - Create new discounts
  - Get all discounts

## Prerequisites
- [Docker](https://docs.docker.com/get-docker/) installed on your system.

## Setup Instructions
 
1. Build and run the Docker container:
    ```bash
   docker-compose up --build
   
2. The application will be accessible at:
    ```
   http://localhost:8080
   
## Testing
1. To test the application, run the following command in the root of the project:
    ```bash
    go test ./...
   
## API Documentation

1. Access the complete documentation at:
    ```
   http://localhost:8080/swagger/index.html
