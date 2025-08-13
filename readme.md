# 📚 Preloved Bookstore – Final Project

A backend system for a secondhand book marketplace built with **Golang** using the **Echo framework**. This project adopts **microservices architecture** and follows **SOLID principles** to ensure maintainability and scalability.

> 🎯 Aligned with SDG Goal: **Category 4 – Quality Education**

## 📝 Changelog

### v1.1 – Redis Caching for Books (Aug 2025)
- ✅ Implemented Redis caching for `Get book by id` and `Create book`
- ✅ Cache new and updated books automatically  
- ✅ Improved read performance for high-traffic endpoints  
- ✅ Fallback to PostgreSQL if cache miss occurs  

### v1.0 – Initial Release (July 2025)
- 🚀 Microservices architecture with Go + Echo  
- 📚 Book CRUD endpoints  
- 👤 User management with JWT authentication  
- 💳 Transaction management with Midtrans integration  
- 📧 Email notifications for transactions
---

## 👥 Team Members

- **Fahreza Abuzard Alghifary**
- **Muhammad Ariq Aziz**
- **Isa Fakhruddin**

---

## ⚙️ Tech Stack

- **Language**: Go (Golang)
- **Framework**: Echo
- **Architecture**: Microservices
- **Database**: PostgreSQL
- **Cache**: Redis (for book retrieval optimization)
- **Structure**: Based on SOLID Principles

### 🔗 3rd-Party Integrations

- **Midtrans** – Payment Gateway
- **SMTP** – Email service
---

## 🗂️ Features Overview

### 👤 User Management
- `POST /user/register` – Register new user  
- `POST /user/login` – Authenticate user  
- `GET /user` – Retrieve profile  
- `PUT /user` – Update profile  
- `PUT /user/balance` – Top-up balance  
- `DELETE /user/:id` – Delete account  

### 📚 Book Management
- `POST /book` – Add book for sale  
- `GET /book` – List all books  
- `GET /book?category=...` – Filter books by category  
- `PUT /book/:id` – Update book details  
- `DELETE /book/:id` – Remove book  

### 💳 Transaction Management
- `POST /transaction` – Create transaction  
- `GET /transaction` – List all user transactions  
- `PUT /transaction` – Update transaction (webhook from 3rd-party)

---

ADDITIONALS:

⏱️ Cron job for soft-deleting expired or failed transactions

📦 Inventory check before transactions

💰 Automatic balance deduction post-payment

✅ Role-based access control (RBAC)


🧪 Testing
Unit tests for core logic

Integration tests for services

Mocking external API calls


🚀 Future Implementations
gRPC – Microservice communication

RabbitMQ – Event-driven architecture

RajaOngkir API – Shipping cost integration

Connection Pooling – With .env config

Grafana & K6 – Monitoring and load testing
