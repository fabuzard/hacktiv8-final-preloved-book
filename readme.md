Final project 


Golang 

Framework : Echo
Database: Postgresql
Folder structure: SOLID
Architecture: Microservices

3rd APIs:
-Midtrans
-verifyright

Tema:
Toko jual beli buku bekas
Category 4 Pendidikan berkualitas

Database:

Users:
ID
Fullname
Email
Password
Address
Role (buyer or seller)
Created_at
Balance

Books:
ID
Seller_id 
Name
Description
Author
Stock
Costs
Category
Created_at


Transactions:
Id
Name
Amount
Book_id
User_id
Transaction_date
Status
Invoice_url
Expiry_date 
Created_at
Deleted_at 



Endpoints:


1. User Management
-  POST User/register
-  POST User/Login
- GET/User
- PUT/User
* PUT/User/balance
* DELETE/User/:id

2. Book Management
POST /book (post buku untuk dijual)
GET /book (get all book)
GET /book
?category=Fiction (optional)
PUT/book/:id
DELETE/book/:id

3. Transaction Management
POST/transaction 
GET /transaction
PUT /transaction (dipake webhook 3rd party buat update data balance & status)



Additional Ideas:
- Cron job to soft delete failed transactions 
- 

Future implementations:
- gRPC
RabbitMQ
RajaOngkir (3rd API)
Connection Pool (tambah env)
Grafana K6s (load test)




https://checkout-staging.xendit.co/web/687a0f6de89114bd8ae96b58











-- =========================================
-- DDL: Table Definitions
-- =========================================

-- Drop existing tables if necessary
DROP TABLE IF EXISTS Transactions;
DROP TABLE IF EXISTS Books;
DROP TABLE IF EXISTS Users;

-- Users Table
CREATE TABLE Users (
    ID SERIAL PRIMARY KEY,
    Fullname VARCHAR(100) NOT NULL,
    Email VARCHAR(100) UNIQUE NOT NULL,
    Password VARCHAR(255) NOT NULL,
    Address TEXT,
    Role VARCHAR(10) NOT NULL,
    Balance DECIMAL(12, 2) DEFAULT 0.00,
    Created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Books Table
CREATE TABLE Books (
    ID SERIAL PRIMARY KEY,
    Seller_id INT REFERENCES Users(ID) ON DELETE CASCADE,
    Name VARCHAR(255) NOT NULL,
    Description TEXT,
    Author VARCHAR(100),
    Stock INT DEFAULT 0,
    Costs DECIMAL(10, 2) NOT NULL,
    Category VARCHAR(100),
    Created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Transactions Table
CREATE TABLE Transactions (
    ID SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    Amount INT NOT NULL,
    Book_id INT REFERENCES Books(ID) ON DELETE SET NULL,
    User_id INT REFERENCES Users(ID) ON DELETE SET NULL,
    Transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    Status VARCHAR(50) NOT NULL,
    Invoice_url TEXT,
    Expiry_date TIMESTAMP,
    Created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    Deleted_at TIMESTAMP
);

-- =========================================
-- DML: Sample Data Insertions
-- =========================================

-- Insert users
INSERT INTO Users (Fullname, Email, Password, Address, Role, Balance)
VALUES 
('Joki', 'joki@example.com', 'hashed_password_1', '123 Buyer St', 'buyer', 500.00),
('Bowo', 'bowo@example.com', 'hashed_password_2', '456 Seller Ave', 'seller', 0.00);

-- Insert books
INSERT INTO Books (Seller_id, Name, Description, Author, Stock, Costs, Category)
VALUES
(2, 'The Go Programming Language', 'An intro to Go', 'Alan A. A. Donovan', 20, 45.50, 'Programming'),
(2, 'Clean Code', 'A Handbook of Agile Software Craftsmanship', 'Robert C. Martin', 15, 39.99, 'Software Engineering');

-- Insert transactions
INSERT INTO Transactions (Name, Amount, Book_id, User_id, Status, Invoice_url, Expiry_date)
VALUES
('Purchase: The Go Programming Language', 1, 1, 1, 'pending', 'http://invoice.url/123', NOW() + INTERVAL '1 day'),
('Purchase: Clean Code', 2, 2, 1, 'paid', 'http://invoice.url/124', NOW() + INTERVAL '1 day');




