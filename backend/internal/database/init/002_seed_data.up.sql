-- Insert test users (password_hash would be properly hashed in real implementation)
INSERT INTO users (email, password_hash, username) VALUES
('aaron@test.com', '$2a$10$dummyhashvalue', 'aaronhotard'),
('aiyana@test.com', '$2a$10$dummyhashvalue', 'aiyanathomas');

-- Insert corresponding profiles
INSERT INTO user_profiles (user_id, first_name, last_name, date_of_birth, gender, height_inches) VALUES
(1, 'Aaron', 'Hotard', '1999-11-19', 'Male', 74),
(2, 'Aiyana', 'Thomas', '2000-07-05', 'Female', 63);