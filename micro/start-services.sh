#!/bin/bash

# Japanese Teacher Booking System - Startup Script
cd "$(dirname "$0")"
ROOT_DIR="$(pwd)"
echo "🇯🇵 Starting JapanLearn - Japanese Teacher Booking System"
echo "========================================================="

# Function to check if a port is available
check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null ; then
        echo "⚠️  Port $port is already in use"
        return 1
    else
        echo "✅ Port $port is available"
        return 0
    fi
}

# Function to start a service
start_service() {
    local service_name=$1
    local service_dir=$2
    local port=$3
    
    echo ""
    echo "🚀 Starting $service_name on port $port..."
    
    if check_port $port; then
        cd "$ROOT_DIR/$service_dir"
        
        # Create .env file if it doesn't exist
        if [ ! -f .env ]; then
            echo "📝 Creating .env file for $service_name..."
            create_env_file $service_name $port
        fi
        
        # Build and run the service
        echo "🔨 Building $service_name..."
        go build -o bin/app cmd/app/main.go
        
        if [ $? -eq 0 ]; then
            echo "✅ $service_name built successfully"
            echo "🏃 Running $service_name..."
            ./bin/app &
            SERVICE_PID=$!
            echo "📋 $service_name PID: $SERVICE_PID"
            
            # Store PID for cleanup
            echo $SERVICE_PID >> ../pids.txt
        else
            echo "❌ Failed to build $service_name"
            return 1
        fi
        
        cd ..
    else
        echo "❌ Cannot start $service_name - port $port is in use"
        return 1
    fi
}

# Function to create .env files
create_env_file() {
    local service_name=$1
    local port=$2
    
    case $service_name in
        "User Service")
            cat > .env << EOF
# User Service Configuration
APP_PORT=$port
GIN_MODE=debug
MYSQL_DSN=root:123@tcp(127.0.0.1:9080)/japanlearn_users?charset=utf8mb4&parseTime=True&loc=Local
DB_MAX_CONECTION=20
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=30

# JWT Configuration
JWT_SECRET_KEY=your-secret-key-here
JWT_TOKEN_DURATION=24

# SMTP Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_EMAIL=noreply@japanlearn.com
SMTP_TEMPLATE_PATH=./templates
SMTP_TEMPLATE_LOGO_URL=https://japanlearn.com/logo.png
SMTP_TIMEOUT_DURATION=30
SMTP_INSECURE_SKIP_VERIFY=false
SMTP_USE_TLS=true
EOF
            ;;
        "Teacher Service")
            cat > .env << EOF
# Teacher Service Configuration
APP_PORT=$port
GIN_MODE=debug
MYSQL_DSN=root:123@tcp(127.0.0.1:9080)/japanlearn_teachers?charset=utf8mb4&parseTime=True&loc=Local
DB_MAX_CONECTION=20
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=30

# JWT Configuration
JWT_SECRET_KEY=your-secret-key-here
JWT_TOKEN_DURATION=24

# Supabase Configuration (for file uploads)
CLIENT_ENDPOINT=https://your-project.supabase.co
CLIENT_ACCESS_KEY=your-access-key
CLIENT_SECRET_KEY=your-secret-key
CLIENT_REGION=us-east-1
CLIENT_BUCKET_NAME=teacher-images
EOF
            ;;
        "Booking Service")
            cat > .env << EOF
# Booking Service Configuration
APP_PORT=$port
GIN_MODE=debug
MYSQL_DSN=root:123@tcp(127.0.0.1:9080)/japanlearn_bookings?charset=utf8mb4&parseTime=True&loc=Local
DB_MAX_CONECTION=20
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=30

# JWT Configuration
JWT_SECRET_KEY=your-secret-key-here
JWT_TOKEN_DURATION=24

# Teacher Service Configuration
SERVICE_SCHEDULE_HOST=localhost
SERVICE_SCHEDULE_PORT=8082

# Supabase Configuration
CLIENT_ENDPOINT=https://your-project.supabase.co
CLIENT_ACCESS_KEY=your-access-key
CLIENT_SECRET_KEY=your-secret-key
CLIENT_REGION=us-east-1
CLIENT_BUCKET_NAME=booking-images
EOF
            ;;
        "Payment Service")
            cat > .env << EOF
# Payment Service Configuration
APP_PORT=$port
GIN_MODE=debug
MYSQL_DSN=root:123@tcp(127.0.0.1:9080)/japanlearn_payments?charset=utf8mb4&parseTime=True&loc=Local
DB_MAX_CONECTION=20
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=30

# JWT Configuration
JWT_SECRET_KEY=your-secret-key-here
JWT_TOKEN_DURATION=24

# Midtrans Configuration
MIDTRANS_SERVER_KEY=your-midtrans-server-key
MIDTRANS_CLIENT_KEY=your-midtrans-client-key
MIDTRANS_ENVIRONMENT=sandbox

# Booking Service Configuration
SERVICE_BOOKING_HOST=localhost
SERVICE_BOOKING_PORT=8083
EOF
            ;;
    esac
}

# Function to start frontend
start_frontend() {
    echo ""
    echo "🌐 Starting Frontend Server with Vite..."
    
    # Check if port 8080 is available
    if check_port 8080; then
        cd "$ROOT_DIR/frontend-vue"
        
        # Check if node_modules exists, if not install dependencies
        if [ ! -d "node_modules" ]; then
            echo "📦 Installing frontend dependencies..."
            npm install
            if [ $? -ne 0 ]; then
                echo "❌ Failed to install frontend dependencies"
                cd ..
                return 1
            fi
        fi
        
        # Start the frontend development server
        echo "🚀 Starting Vite development server..."
        npm run dev &
        FRONTEND_PID=$!
        
        echo "📋 Frontend PID: $FRONTEND_PID"
        echo $FRONTEND_PID >> ../pids.txt
        cd ..
        
        # Wait a moment for the server to start
        sleep 3
        
        # Check if the frontend is actually running
        if kill -0 $FRONTEND_PID 2>/dev/null; then
            echo "✅ Frontend server started successfully"
        else
            echo "❌ Frontend server failed to start"
            return 1
        fi
    else
        echo "❌ Cannot start frontend - port 8080 is in use"
        return 1
    fi
}

# Function to cleanup processes
cleanup() {
    echo ""
    echo "🧹 Cleaning up processes..."
    
    if [ -f pids.txt ]; then
        while read pid; do
            if kill -0 $pid 2>/dev/null; then
                echo "🔪 Killing process $pid"
                kill $pid
            fi
        done < pids.txt
        rm pids.txt
    fi
    
    echo "✅ Cleanup completed"
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM

# Clear any existing PID file
rm -f pids.txt

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

echo "✅ Go is installed: $(go version)"

# Start services
echo ""
echo "🚀 Starting microservices..."

# Start User Service (Port 8081)
start_service "User Service" "user" 8081

# Wait a bit between services
sleep 2

# Start Teacher Service (Port 8082)
start_service "Teacher Service" "teacher" 8082

# Wait a bit between services
sleep 2

# Start Booking Service (Port 8083)
start_service "Booking Service" "booking" 8083

# Wait a bit between services
sleep 2

# Start Payment Service (Port 8084)
start_service "Payment Service" "payment" 8084

# Wait a bit before starting frontend
sleep 3

# Start Frontend (Port 8080)
start_frontend

echo ""
echo "🎉 All services started successfully!"
echo "========================================="
echo "📱 Frontend:        http://localhost:8080"
echo "👤 User Service:    http://localhost:8081"
echo "👨‍🏫 Teacher Service: http://localhost:8082"
echo "📅 Booking Service: http://localhost:8083"
echo "💳 Payment Service: http://localhost:8084"
echo ""
echo "📋 Service Status:"
echo "==================="

# Check service status
for port in 8080 8081 8082 8083 8084; do
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null ; then
        echo "✅ Port $port: Running"
    else
        echo "❌ Port $port: Not running"
    fi
done

echo ""
echo "💡 Tips:"
echo "- Press Ctrl+C to stop all services"
echo "- Check logs in terminal for any errors"
echo "- Make sure MySQL is running on localhost:3306"
echo "- Update .env files with your actual database credentials"
echo ""
echo "🇯🇵 JapanLearn is ready! Visit http://localhost:8080 to start learning!"

# Keep the script running
wait
