-- Users and Authentication
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP WITH TIME ZONE
);

-- User Profiles
CREATE TABLE user_profiles (
    profile_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    date_of_birth DATE,
    gender VARCHAR(20),
    height_inches INTEGER,
    weight_pounds INTEGER,
    profile_picture_url VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id)
);

CREATE TABLE workouts (
    workout_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    workout_date DATE NOT NULL,
    title TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, workout_date) -- to prevent >1 workouts/day/user, since implicitly on any given calendar day in the app 1 workout can occur
);

CREATE TABLE workout_exercises (
    workout_exercise_id SERIAL PRIMARY KEY,
    workout_id INTEGER REFERENCES workouts(workout_id),
    exercise_name VARCHAR(100) NOT NULL,
    sets INTEGER NOT NULL,
    reps INTEGER NOT NULL,
    resistance_type VARCHAR(50) NOT NULL,  -- 'weight', 'band', 'bodyweight', etc.
    resistance_value INTEGER,              -- Optional; amt in lbs or band 'level'
    resistance_detail VARCHAR(100),        -- Optional; 'red band', 'blue band', 'assisted', etc.
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- This is mostly going to be a reference table of exercises.
-- Eventually we'll add a multi-valued user_id attribute to show which users have performed which exercises, and maybe we can use that for caching recently-performed exercises for users
CREATE TABLE exercises (
    exercise_id SERIAL PRIMARY KEY,
    exercise_name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE workout_sets (
    workout_id INTEGER REFERENCES workouts(workout_id),
    exercise_id INTEGER REFERENCES exercises(exercise_id),
    set_number SERIAL, -- auto-increment from 1
    reps INTEGER, -- user can leave reps null intentionally - would mean user hasn't yet performed exercise (or skipped it and wanted to make it clear in their logs).
    resistance_type VARCHAR(50) NOT NULL,
    resistance_value INTEGER, -- another user-nullable field
    resistance_detail VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (workout_id, exercise_id, set_number)
);

-- This will play a part in calculating weekly volume reports showing muscle groups over/underworked
CREATE TABLE muscles (
    muscle_id SERIAL PRIMARY KEY,
    muscle_name VARCHAR(50) NOT NULL UNIQUE,
    muscle_group VARCHAR(50) NOT NULL,  -- e.g., 'Back', 'Chest', 'Legs'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Junction table for exercise-muscle relationships
CREATE TABLE exercise_muscles (
    exercise_id INTEGER REFERENCES exercises(exercise_id),
    muscle_id INTEGER REFERENCES muscles(muscle_id),
    involvement_level VARCHAR(20) NOT NULL,  -- 'Primary', 'Secondary', or 'Stabilizer'
    PRIMARY KEY (exercise_id, muscle_id)
);

-- For frequently queried fields
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);