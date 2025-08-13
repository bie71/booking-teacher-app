# 🇯🇵 JapanLearn - Japanese Teacher Booking System

A modern, full-stack web application for booking Japanese language lessons with native teachers. Built with Go microservices backend and vanilla JavaScript frontend.

## 🌟 Features

### For Students
- **Teacher Discovery**: Browse and filter Japanese teachers by level, price, and availability
- **Easy Booking**: Book lessons with your preferred teachers
- **Schedule Management**: View, reschedule, or cancel your bookings
- **Payment Integration**: Secure payment processing with multiple payment methods
- **Progress Tracking**: Monitor your learning journey and completed lessons
- **User Dashboard**: Comprehensive dashboard to manage your account and bookings

### For Teachers
- **Profile Management**: Create and manage your teaching profile
- **Schedule Control**: Set your availability and manage time slots
- **Student Management**: View and manage your student bookings
- **Image Upload**: Upload profile pictures and teaching materials

### System Features
- **Microservices Architecture**: Scalable and maintainable service-oriented design
- **RESTful APIs**: Well-documented REST endpoints for all services
- **CORS Support**: Cross-origin resource sharing for frontend-backend communication
- **Responsive Design**: Mobile-friendly UI that works on all devices
- **Real-time Updates**: Dynamic content updates without page refresh

## 🏗️ Architecture

The system consists of 5 microservices:

1. **User Service** (Port 8081) - User authentication and profile management
2. **Teacher Service** (Port 8082) - Teacher profiles and schedule management
3. **Booking Service** (Port 8083) - Lesson booking and management
4. **Payment Service** (Port 8084) - Payment processing and transactions
5. **Frontend** (Port 8080) - Web interface for users

## 🛠️ Technology Stack

### Backend
- **Language**: Go 1.24+
- **Framework**: Gin (HTTP web framework)
- **Database**: MySQL with GORM ORM
- **Authentication**: JWT tokens
- **Payment**: Midtrans integration
- **File Storage**: Supabase integration
- **Architecture**: Microservices with REST APIs

### Frontend
- **Languages**: HTML5, CSS3, JavaScript (ES6+)
- **Styling**: Custom CSS with CSS Grid and Flexbox
- **Icons**: Font Awesome
- **Fonts**: Inter & Noto Sans JP
- **Architecture**: Single Page Application (SPA)

## 🚀 Quick Start

### Prerequisites
- Go 1.24 or higher
- MySQL 8.0 or higher
- Python 3 or Node.js (for frontend server)
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd micro
   ```

2. **Start all services**
   ```bash
   ./start-services.sh
   ```

   This script will:
   - Check system requirements
   - Create necessary .env files
   - Build and start all Go services
   - Start the frontend server
   - Display service status

3. **Access the application**
   - Frontend: http://localhost:8080
   - User Service: http://localhost:8081
   - Teacher Service: http://localhost:8082
   - Booking Service: http://localhost:8083
   - Payment Service: http://localhost:8084

### Manual Setup

If you prefer to set up services manually:

1. **Set up databases**
   ```sql
   CREATE DATABASE japanlearn_users;
   CREATE DATABASE japanlearn_teachers;
   CREATE DATABASE japanlearn_bookings;
   CREATE DATABASE japanlearn_payments;
   ```

2. **Configure environment variables**
   
   Create `.env` files in each service directory with appropriate configurations:

   **user/.env**
   ```env
   APP_PORT=8081
   MYSQL_DSN=root:password@tcp(localhost:3306)/japanlearn_users?charset=utf8mb4&parseTime=True&loc=Local
   JWT_SECRET_KEY=your-secret-key-here
   JWT_TOKEN_DURATION=24
   # ... other configurations
   ```

3. **Start services individually**
   ```bash
   # User Service
   cd user && go run cmd/app/main.go

   # Teacher Service
   cd teacher && go run cmd/app/main.go

   # Booking Service
   cd booking && go run cmd/app/main.go

   # Payment Service
   cd payment && go run cmd/app/main.go

   # Frontend (in a new terminal)
   cd frontend && python3 -m http.server 8080
   ```

## 📚 API Documentation

### User Service (Port 8081)
- `POST /api/v1/register` - Register new user
- `POST /api/v1/login` - User login
- `GET /api/v1/me` - Get current user profile
- `POST /api/v1/reset-password` - Reset password

### Teacher Service (Port 8082)
- `GET /api/v1/teachers` - Get all teachers
- `GET /api/v1/teachers/:id` - Get teacher by ID
- `POST /api/v1/teachers` - Create teacher profile
- `PUT /api/v1/teachers/:id` - Update teacher profile
- `GET /api/v1/schedule/teacher/:teacher_id` - Get teacher schedules
- `POST /api/v1/schedule` - Create schedule
- `PUT /api/v1/schedule/:id` - Update schedule

### Booking Service (Port 8083)
- `POST /api/v1/bookings` - Create booking
- `GET /api/v1/bookings` - Get all bookings
- `GET /api/v1/booking/:id` - Get booking by ID
- `GET /api/v1/bookings/user/:user_id` - Get user bookings
- `POST /api/v1/bookings/:id/reschedule` - Reschedule booking
- `POST /api/v1/bookings/:id/cancel` - Cancel booking

### Payment Service (Port 8084)
- `POST /api/v1/payments` - Create payment
- `GET /api/v1/payments` - Get all payments
- `GET /api/v1/payment/:id` - Get payment by ID
- `POST /api/v1/payments/callback` - Payment webhook

## 🎨 Frontend Structure

```
frontend/
├── index.html          # Main HTML file
├── css/
│   └── style.css       # Main stylesheet
└── js/
    ├── config.js       # Configuration and constants
    ├── auth.js         # Authentication management
    ├── api.js          # API client and services
    ├── teachers.js     # Teacher management
    ├── booking.js      # Booking management
    ├── dashboard.js    # User dashboard
    └── main.js         # Main application logic
```

## 🔧 Configuration

### Environment Variables

Each service requires specific environment variables. The startup script creates default `.env` files, but you should update them with your actual values:

**Database Configuration**
- `MYSQL_DSN`: MySQL connection string
- `DB_MAX_CONECTION`: Maximum database connections
- `DB_MAX_IDLE_CONNS`: Maximum idle connections

**JWT Configuration**
- `JWT_SECRET_KEY`: Secret key for JWT tokens
- `JWT_TOKEN_DURATION`: Token expiration time in hours

**External Services**
- `MIDTRANS_SERVER_KEY`: Midtrans payment gateway server key
- `MIDTRANS_CLIENT_KEY`: Midtrans payment gateway client key
- `CLIENT_ENDPOINT`: Supabase endpoint for file storage

## 🧪 Testing

### API Testing
You can test the APIs using tools like Postman or curl:

```bash
# Test user registration
curl -X POST http://localhost:8081/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"password123"}'

# Test teacher listing
curl http://localhost:8082/api/v1/teachers
```

### Frontend Testing
1. Open http://localhost:8080 in your browser
2. Register a new account
3. Browse teachers
4. Book a lesson
5. Check your dashboard

## 🚀 Deployment

### Docker Deployment (Recommended)

1. **Create Dockerfiles for each service**
2. **Use docker-compose for orchestration command: `docker compose up --build -d`**
3. **Set up environment-specific configurations**

### Traditional Deployment

1. **Build binaries for each service**
   ```bash
   cd user && go build -o bin/user cmd/app/main.go
   cd teacher && go build -o bin/teacher cmd/app/main.go
   cd booking && go build -o bin/booking cmd/app/main.go
   cd payment && go build -o bin/payment cmd/app/main.go
   ```

2. **Deploy to your server**
3. **Set up reverse proxy (Nginx)**
4. **Configure SSL certificates**
5. **Set up monitoring and logging**

## 🔒 Security Considerations

- **JWT Authentication**: Secure token-based authentication
- **CORS Configuration**: Properly configured cross-origin requests
- **Input Validation**: Server-side validation for all inputs
- **SQL Injection Prevention**: Using GORM ORM with parameterized queries
- **Password Hashing**: Secure password storage
- **HTTPS**: Use SSL/TLS in production

## 🐛 Troubleshooting

### Common Issues

1. **Port Already in Use**
   ```bash
   # Find and kill process using port
   lsof -ti:8080 | xargs kill -9
   ```

2. **Database Connection Failed**
   - Check MySQL is running
   - Verify database credentials in .env files
   - Ensure databases exist

3. **CORS Errors**
   - Check if all services have CORS middleware enabled
   - Verify frontend is accessing correct API URLs

4. **Build Failures**
   - Ensure Go version is 1.24+
   - Run `go mod tidy` in each service directory
   - Check for missing dependencies

### Logs and Debugging

- Service logs are displayed in the terminal
- Enable debug mode by setting `GIN_MODE=debug` in .env files
- Use browser developer tools for frontend debugging

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Go Community** for excellent frameworks and libraries
- **Gin Framework** for fast HTTP routing
- **GORM** for elegant ORM
- **Font Awesome** for beautiful icons
- **Google Fonts** for typography

## 📞 Support

If you encounter any issues or have questions:

1. Check the troubleshooting section
2. Search existing issues
3. Create a new issue with detailed information
4. Contact the development team

---

**Happy Learning Japanese! 🇯🇵 頑張って！**
