-- Demo Data for JapanLearn - Japanese Teacher Booking System
-- Run this script after setting up the databases

-- Use the teachers database
USE japanlearn_teachers;

-- Insert sample teachers
INSERT INTO teachers (name, bio, language_level, price_per_hour, available_start_time, available_end_time, profile_image, created_at, updated_at) VALUES
('Tanaka Hiroshi', 'Native Japanese speaker with 5 years of teaching experience. Specializes in conversational Japanese and business communication. I love helping students build confidence in speaking Japanese naturally.', 'intermediate', 25.00, '09:00', '17:00', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop&crop=face', NOW(), NOW()),

('Sato Yuki', 'Certified Japanese language teacher with JLPT expertise. I have been teaching Japanese for over 8 years and specialize in grammar, reading comprehension, and test preparation.', 'advanced', 35.00, '10:00', '18:00', 'https://images.unsplash.com/photo-1494790108755-2616b612b786?w=150&h=150&fit=crop&crop=face', NOW(), NOW()),

('Yamamoto Kenji', 'Friendly and patient teacher perfect for beginners. I use interactive methods and real-life situations to make learning Japanese fun and practical. Let\'s start your Japanese journey together!', 'beginner', 20.00, '08:00', '16:00', 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face', NOW(), NOW()),

('Nakamura Akiko', 'Professional Japanese instructor with university teaching background. I specialize in academic Japanese, formal writing, and advanced grammar structures for serious learners.', 'advanced', 40.00, '11:00', '19:00', 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=150&h=150&fit=crop&crop=face', NOW(), NOW()),

('Suzuki Takeshi', 'Experienced conversation partner focusing on natural Japanese communication. Former business professional who can teach both casual and business Japanese effectively.', 'intermediate', 30.00, '09:30', '17:30', 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=150&h=150&fit=crop&crop=face', NOW(), NOW()),

('Watanabe Mei', 'Young and energetic teacher who makes learning Japanese enjoyable. I use modern teaching methods including games, multimedia, and cultural activities to engage students.', 'beginner', 22.00, '10:00', '18:00', 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=150&h=150&fit=crop&crop=face', NOW(), NOW()),

('Ito Masaki', 'JLPT specialist with proven track record of student success. I have helped over 200 students pass their JLPT exams from N5 to N1 level. Structured and goal-oriented approach.', 'advanced', 45.00, '12:00', '20:00', 'https://images.unsplash.com/photo-1507591064344-4c6ce005b128?w=150&h=150&fit=crop&crop=face', NOW(), NOW()),

('Kobayashi Rina', 'Cultural enthusiast who teaches Japanese through cultural immersion. Learn not just the language but also Japanese customs, traditions, and modern culture.', 'intermediate', 28.00, '09:00', '17:00', 'https://images.unsplash.com/photo-1487412720507-e7ab37603c6f?w=150&h=150&fit=crop&crop=face', NOW(), NOW());

-- Insert sample schedules for the next 7 days
-- Teacher 1 (Tanaka Hiroshi) schedules
INSERT INTO schedules (teacher_id, date, start_time, end_time, status, created_at, updated_at) VALUES
(1, CURDATE(), '09:00', '10:00', 'available', NOW(), NOW()),
(1, CURDATE(), '10:00', '11:00', 'available', NOW(), NOW()),
(1, CURDATE(), '14:00', '15:00', 'available', NOW(), NOW()),
(1, CURDATE(), '15:00', '16:00', 'available', NOW(), NOW()),
(1, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '09:00', '10:00', 'available', NOW(), NOW()),
(1, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '11:00', '12:00', 'available', NOW(), NOW()),
(1, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '14:00', '15:00', 'available', NOW(), NOW()),
(1, DATE_ADD(CURDATE(), INTERVAL 2 DAY), '10:00', '11:00', 'available', NOW(), NOW()),
(1, DATE_ADD(CURDATE(), INTERVAL 2 DAY), '15:00', '16:00', 'available', NOW(), NOW());

-- Teacher 2 (Sato Yuki) schedules
INSERT INTO schedules (teacher_id, date, start_time, end_time, status, created_at, updated_at) VALUES
(2, CURDATE(), '10:00', '11:00', 'available', NOW(), NOW()),
(2, CURDATE(), '11:00', '12:00', 'available', NOW(), NOW()),
(2, CURDATE(), '15:00', '16:00', 'available', NOW(), NOW()),
(2, CURDATE(), '16:00', '17:00', 'available', NOW(), NOW()),
(2, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '10:00', '11:00', 'available', NOW(), NOW()),
(2, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '13:00', '14:00', 'available', NOW(), NOW()),
(2, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '15:00', '16:00', 'available', NOW(), NOW()),
(2, DATE_ADD(CURDATE(), INTERVAL 3 DAY), '11:00', '12:00', 'available', NOW(), NOW()),
(2, DATE_ADD(CURDATE(), INTERVAL 3 DAY), '16:00', '17:00', 'available', NOW(), NOW());

-- Teacher 3 (Yamamoto Kenji) schedules
INSERT INTO schedules (teacher_id, date, start_time, end_time, status, created_at, updated_at) VALUES
(3, CURDATE(), '08:00', '09:00', 'available', NOW(), NOW()),
(3, CURDATE(), '09:00', '10:00', 'available', NOW(), NOW()),
(3, CURDATE(), '13:00', '14:00', 'available', NOW(), NOW()),
(3, CURDATE(), '14:00', '15:00', 'available', NOW(), NOW()),
(3, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '08:00', '09:00', 'available', NOW(), NOW()),
(3, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '10:00', '11:00', 'available', NOW(), NOW()),
(3, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '13:00', '14:00', 'available', NOW(), NOW()),
(3, DATE_ADD(CURDATE(), INTERVAL 2 DAY), '09:00', '10:00', 'available', NOW(), NOW()),
(3, DATE_ADD(CURDATE(), INTERVAL 2 DAY), '14:00', '15:00', 'available', NOW(), NOW());

-- Teacher 4 (Nakamura Akiko) schedules
INSERT INTO schedules (teacher_id, date, start_time, end_time, status, created_at, updated_at) VALUES
(4, CURDATE(), '11:00', '12:00', 'available', NOW(), NOW()),
(4, CURDATE(), '13:00', '14:00', 'available', NOW(), NOW()),
(4, CURDATE(), '17:00', '18:00', 'available', NOW(), NOW()),
(4, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '11:00', '12:00', 'available', NOW(), NOW()),
(4, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '14:00', '15:00', 'available', NOW(), NOW()),
(4, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '17:00', '18:00', 'available', NOW(), NOW()),
(4, DATE_ADD(CURDATE(), INTERVAL 3 DAY), '12:00', '13:00', 'available', NOW(), NOW()),
(4, DATE_ADD(CURDATE(), INTERVAL 3 DAY), '16:00', '17:00', 'available', NOW(), NOW());

-- Teacher 5 (Suzuki Takeshi) schedules
INSERT INTO schedules (teacher_id, date, start_time, end_time, status, created_at, updated_at) VALUES
(5, CURDATE(), '09:30', '10:30', 'available', NOW(), NOW()),
(5, CURDATE(), '11:30', '12:30', 'available', NOW(), NOW()),
(5, CURDATE(), '14:30', '15:30', 'available', NOW(), NOW()),
(5, CURDATE(), '16:30', '17:30', 'available', NOW(), NOW()),
(5, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '09:30', '10:30', 'available', NOW(), NOW()),
(5, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '13:30', '14:30', 'available', NOW(), NOW()),
(5, DATE_ADD(CURDATE(), INTERVAL 2 DAY), '10:30', '11:30', 'available', NOW(), NOW()),
(5, DATE_ADD(CURDATE(), INTERVAL 2 DAY), '15:30', '16:30', 'available', NOW(), NOW());

-- Teacher 6 (Watanabe Mei) schedules
INSERT INTO schedules (teacher_id, date, start_time, end_time, status, created_at, updated_at) VALUES
(6, CURDATE(), '10:00', '11:00', 'available', NOW(), NOW()),
(6, CURDATE(), '12:00', '13:00', 'available', NOW(), NOW()),
(6, CURDATE(), '15:00', '16:00', 'available', NOW(), NOW()),
(6, CURDATE(), '17:00', '18:00', 'available', NOW(), NOW()),
(6, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '10:00', '11:00', 'available', NOW(), NOW()),
(6, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '14:00', '15:00', 'available', NOW(), NOW()),
(6, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '16:00', '17:00', 'available', NOW(), NOW()),
(6, DATE_ADD(CURDATE(), INTERVAL 3 DAY), '11:00', '12:00', 'available', NOW(), NOW());

-- Teacher 7 (Ito Masaki) schedules
INSERT INTO schedules (teacher_id, date, start_time, end_time, status, created_at, updated_at) VALUES
(7, CURDATE(), '12:00', '13:00', 'available', NOW(), NOW()),
(7, CURDATE(), '14:00', '15:00', 'available', NOW(), NOW()),
(7, CURDATE(), '18:00', '19:00', 'available', NOW(), NOW()),
(7, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '12:00', '13:00', 'available', NOW(), NOW()),
(7, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '15:00', '16:00', 'available', NOW(), NOW()),
(7, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '18:00', '19:00', 'available', NOW(), NOW()),
(7, DATE_ADD(CURDATE(), INTERVAL 2 DAY), '13:00', '14:00', 'available', NOW(), NOW()),
(7, DATE_ADD(CURDATE(), INTERVAL 2 DAY), '17:00', '18:00', 'available', NOW(), NOW());

-- Teacher 8 (Kobayashi Rina) schedules
INSERT INTO schedules (teacher_id, date, start_time, end_time, status, created_at, updated_at) VALUES
(8, CURDATE(), '09:00', '10:00', 'available', NOW(), NOW()),
(8, CURDATE(), '11:00', '12:00', 'available', NOW(), NOW()),
(8, CURDATE(), '14:00', '15:00', 'available', NOW(), NOW()),
(8, CURDATE(), '16:00', '17:00', 'available', NOW(), NOW()),
(8, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '09:00', '10:00', 'available', NOW(), NOW()),
(8, DATE_ADD(CURDATE(), INTERVAL 1 DAY), '12:00', '13:00', 'available', NOW(), NOW()),
(8, DATE_ADD(CURDATE(), INTERVAL 2 DAY), '10:00', '11:00', 'available', NOW(), NOW()),
(8, DATE_ADD(CURDATE(), INTERVAL 2 DAY), '15:00', '16:00', 'available', NOW(), NOW());

-- Use the users database
USE japanlearn_users;

-- Insert sample users (passwords are hashed for 'password123')
INSERT INTO users (name, email, password_hash, role, created_at, updated_at) VALUES
('John Smith', 'john.smith@example.com', '$2a$10$KN9OeCthb7Cp9sUssP0yfOfgxBHU7kMe8TcjC8otZ3Lw2slRs3sDi', 'user', NOW(), NOW()),
('Sarah Johnson', 'sarah.johnson@example.com', '$2a$10$KN9OeCthb7Cp9sUssP0yfOfgxBHU7kMe8TcjC8otZ3Lw2slRs3sDi', 'user', NOW(), NOW()),
('Mike Davis', 'mike.davis@example.com', '$2a$10$KN9OeCthb7Cp9sUssP0yfOfgxBHU7kMe8TcjC8otZ3Lw2slRs3sDi', 'user', NOW(), NOW()),
('Emily Wilson', 'emily.wilson@example.com', '$2a$10$KN9OeCthb7Cp9sUssP0yfOfgxBHU7kMe8TcjC8otZ3Lw2slRs3sDi', 'user', NOW(), NOW()),
('Admin User', 'admin@japanlearn.com', '$2a$10$KN9OeCthb7Cp9sUssP0yfOfgxBHU7kMe8TcjC8otZ3Lw2slRs3sDi', 'admin', NOW(), NOW());

-- Use the bookings database
USE japanlearn_bookings;

-- Insert sample bookings
INSERT INTO bookings (user_id, schedule_id, status, note, created_at, updated_at) VALUES
(1, 1, 'paid', 'Looking forward to learning basic conversation skills', NOW(), NOW()),
(2, 10, 'pending', 'Need help with JLPT N3 preparation', NOW(), NOW()),
(3, 20, 'paid', 'Want to focus on business Japanese', NOW(), NOW()),
(1, 25, 'completed', 'Great lesson on hiragana and katakana', DATE_SUB(NOW(), INTERVAL 7 DAY), DATE_SUB(NOW(), INTERVAL 7 DAY));

-- Use the payments database
USE japanlearn_payments;

-- Insert sample payments
INSERT INTO payments (midtrans_transaction_id, amount, status, payment_method, booking_id, paid_at, created_at, updated_at) VALUES
('TXN001', 25.00, 'settlement', 'credit_card', 1, NOW(), NOW(), NOW()),
('TXN002', 30.00, 'pending', 'bank_transfer', 2, NULL, NOW(), NOW()),
('TXN003', 28.00, 'settlement', 'e_wallet', 3, NOW(), NOW(), NOW()),
('TXN004', 20.00, 'settlement', 'credit_card', 4, DATE_SUB(NOW(), INTERVAL 7 DAY), DATE_SUB(NOW(), INTERVAL 7 DAY), DATE_SUB(NOW(), INTERVAL 7 DAY));

-- Insert sample payment methods
INSERT INTO payment_methods (code,name, is_active, created_at, updated_at) VALUES
('credit_card','Credit Card', 1, NOW(), NOW()),
('bank_transfer','Bank Transfer', 1, NOW(), NOW()),
('e_wallet','E-Wallet', 1, NOW(), NOW()),
('paypal','PayPal', 0, NOW(), NOW());

COMMIT;

-- Display summary
SELECT 'Demo data inserted successfully!' as message;
SELECT 'Teachers created:' as info, COUNT(*) as count FROM japanlearn_teachers.teachers;
SELECT 'Schedules created:' as info, COUNT(*) as count FROM japanlearn_teachers.schedules;
SELECT 'Users created:' as info, COUNT(*) as count FROM japanlearn_users.users;
SELECT 'Bookings created:' as info, COUNT(*) as count FROM japanlearn_bookings.bookings;
SELECT 'Payments created:' as info, COUNT(*) as count FROM japanlearn_payments.payments;
