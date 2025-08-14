#!/bin/bash

# Database Setup Script for JapanLearn
echo "🗄️  Setting up JapanLearn Databases"
echo "===================================="

# Default MySQL credentials (change these as needed)
MYSQL_HOST="127.0.0.1"
MYSQL_PORT="3306"
MYSQL_USER="root"
MYSQL_PASSWORD=""

# Function to check if MySQL is running
check_mysql() {
    if ! command -v mysql &> /dev/null; then
        echo "❌ MySQL is not installed or not in PATH"
        echo "Please install MySQL and try again"
        exit 1
    fi

    if ! mysqladmin ping -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" --protocol=TCP --silent 2>/dev/null; then
        echo "❌ Cannot connect to MySQL server"
        echo "Please check your MySQL server is running and credentials are correct"
        echo "Host: $MYSQL_HOST:$MYSQL_PORT"
        echo "User: $MYSQL_USER"
        exit 1
    fi

    echo "✅ MySQL connection successful"
}

# Function to create databases
create_databases() {
    echo ""
    echo "📝 Creating databases..."

    databases=("japanlearn_users" "japanlearn_teachers" "japanlearn_bookings" "japanlearn_payments")

    for db in "${databases[@]}"; do
        echo "Creating database: $db"
        mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" --protocol=TCP -e "CREATE DATABASE IF NOT EXISTS $db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null
        
        if [ $? -eq 0 ]; then
            echo "✅ Database $db created successfully"
        else
            echo "❌ Failed to create database $db"
            exit 1
        fi
    done
}

# Function to create tables (basic structure)
create_tables() {
    echo ""
    echo "🏗️  Creating tables..."

    # Users table
    echo "Creating users table..."
    mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" --protocol=TCP japanlearn_users << 'EOF'
CREATE TABLE IF NOT EXISTS users (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role ENUM('user', 'admin', 'teacher') NOT NULL DEFAULT 'user',
    profile_image VARCHAR(255),
    reset_token VARCHAR(255),
    reset_expiration TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS activity_logs (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNSIGNED NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



CREATE TABLE IF NOT EXISTS favorite_teachers (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNSIGNED NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    teacher_id INT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (user_id, teacher_id)
);

CREATE TABLE `hero_images` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `key_image` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `image_url` text COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_key_image` (`key_image`)
);

EOF

    # Teachers table
    echo "Creating teachers table..."
    mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" --protocol=TCP japanlearn_teachers << 'EOF'
CREATE TABLE IF NOT EXISTS teachers (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNSIGNED,
    name VARCHAR(100) NOT NULL,
    bio TEXT,
    language_level VARCHAR(20) DEFAULT 'beginner',
    price_per_hour DECIMAL(10,2) NOT NULL,
    available_start_time VARCHAR(8) DEFAULT '09:00',
    available_end_time VARCHAR(8) DEFAULT '17:00',
    profile_image VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS schedules (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    teacher_id INT UNSIGNED NOT NULL,
    date DATE NOT NULL,
    start_time VARCHAR(8) NOT NULL,
    end_time VARCHAR(8) NOT NULL,
    status ENUM('available', 'booked', 'cancelled') DEFAULT 'available',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE CASCADE,
    INDEX idx_teacher_date (teacher_id, date),
    INDEX idx_status (status)
);
EOF

    # Bookings table
    echo "Creating bookings table..."
    mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" --protocol=TCP japanlearn_bookings << 'EOF'
CREATE TABLE `bookings` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned DEFAULT NULL,
  `schedule_id` bigint unsigned DEFAULT NULL,
  `status` enum('pending','paid','cancelled','rescheduled','completed') COLLATE utf8mb4_unicode_ci DEFAULT 'pending',
  `payment_id` bigint unsigned DEFAULT NULL,
  `reschedule_from` bigint unsigned DEFAULT NULL,
  `note` longtext COLLATE utf8mb4_unicode_ci,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `total_price` double DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_schedule_id` (`schedule_id`),
  KEY `idx_status` (`status`),
  KEY `idx_bookings_user_id` (`user_id`),
  KEY `idx_bookings_schedule_id` (`schedule_id`)
);
EOF

    # Payments table
    echo "Creating payments table..."
    mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" --protocol=TCP japanlearn_payments << 'EOF'
CREATE TABLE IF NOT EXISTS payments (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    midtrans_transaction_id VARCHAR(100),
    amount DECIMAL(10,2) NOT NULL,
    status ENUM('pending', 'settlement', 'failed', 'cancel') DEFAULT 'pending',
    payment_method VARCHAR(50),
    booking_id INT UNSIGNED NOT NULL,
    paid_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_booking_id (booking_id),
    INDEX idx_status (status)
);

CREATE TABLE IF NOT EXISTS payment_methods (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
EOF

    echo "✅ Tables created successfully"
}

# Function to load demo data
load_demo_data() {
    echo ""
    echo "📊 Loading demo data..."
    
    if [ -f "demo-data.sql" ]; then
        mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" --protocol=TCP < demo-data.sql
        
        if [ $? -eq 0 ]; then
            echo "✅ Demo data loaded successfully"
        else
            echo "❌ Failed to load demo data"
            exit 1
        fi
    else
        echo "⚠️  demo-data.sql not found, skipping demo data"
    fi
}

# Function to show database info
show_info() {
    echo ""
    echo "📋 Database Setup Complete!"
    echo "=========================="
    echo "Host: $MYSQL_HOST:$MYSQL_PORT"
    echo "User: $MYSQL_USER"
    echo ""
    echo "Databases created:"
    echo "- japanlearn_users"
    echo "- japanlearn_teachers" 
    echo "- japanlearn_bookings"
    echo "- japanlearn_payments"
    echo ""
    echo "🔧 Next steps:"
    echo "1. Update .env files in each service with your database credentials"
    echo "2. Run ./start-services.sh to start the application"
    echo "3. Visit http://localhost:PORT to access the frontend"
    echo ""
    echo "📝 Demo accounts (password: password123):"
    echo "- john.smith@example.com"
    echo "- sarah.johnson@example.com"
    echo "- admin@japanlearn.com (admin)"
}

# Main execution
echo "Starting database setup..."

# Prompt for MySQL credentials if needed
if [ -n "$1" ]; then
    MYSQL_PASSWORD="$1"
fi

if [ -z "$MYSQL_PASSWORD" ]; then
    echo ""
    echo "Please enter your MySQL root password (input will be hidden):"
    read -s MYSQL_PASSWORD
fi

# Run setup steps
check_mysql
create_databases
create_tables
load_demo_data
show_info

echo ""
echo "🎉 Database setup completed successfully!"
