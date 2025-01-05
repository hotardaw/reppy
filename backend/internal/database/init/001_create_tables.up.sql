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

-- This is intended as a reference table of exercises.
-- Eventually we'll add a multi-valued user_id attribute to show which users have performed which exercises, and maybe we can use that for caching recently-performed exercises for users
CREATE TABLE exercises (
    exercise_id SERIAL PRIMARY KEY,
    exercise_name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE resistance_type_enum AS ENUM ('weight', 'band', 'bodyweight');
CREATE TABLE workout_sets (
    workout_id INTEGER REFERENCES workouts(workout_id) NOT NULL,
    exercise_id INTEGER REFERENCES exercises(exercise_id) NOT NULL,
    set_number INTEGER NOT NULL,          -- We'll handle this in the application layer
    reps INTEGER,                         -- Optional - filled in when performed
    resistance_value INTEGER,             -- Optional - weight in lbs/kg
    resistance_type resistance_type_enum, -- Optional - 'weight', 'band', 'bodyweight' only
    resistance_detail VARCHAR(100),       -- Optional - band color, cable attachment, etc.
    rpe DECIMAL(3,1),                     -- Optional
    percent_1rm DECIMAL(4,1),             -- Optional, allows values like 77.5, 82.5, etc.
    notes TEXT,                           -- Optional
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (workout_id, exercise_id, set_number)
);

-- This will contribute to generating (weekly?) volume reports for users showing muscle groups they might've over/underworked
CREATE TABLE muscles (
    muscle_id SERIAL PRIMARY KEY,
    muscle_name VARCHAR(50) NOT NULL UNIQUE,
    muscle_group VARCHAR(50) NOT NULL,  -- e.g., 'Back', 'Chest', 'Legs'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Junction table for exercise-muscle relationships
CREATE TYPE involvement_level_enum AS ENUM ('primary', 'secondary');
CREATE TABLE exercise_muscles (
    exercise_id INTEGER REFERENCES exercises(exercise_id),
    muscle_id INTEGER REFERENCES muscles(muscle_id),
    involvement_level involvement_level_enum NOT NULL,
    PRIMARY KEY (exercise_id, muscle_id)
);




CREATE VIEW exercise_one_rm AS
SELECT 
    w.user_id,
    e.exercise_name,
    MAX(ws.resistance_value / (1.0278 - 0.0278 * ws.reps)) as estimated_1rm
FROM workouts w
JOIN workout_sets ws ON w.workout_id = ws.workout_id
JOIN exercises e ON ws.exercise_id = e.exercise_id
WHERE ws.resistance_type = 'weight'
  AND ws.resistance_value IS NOT NULL 
  AND ws.reps IS NOT NULL
GROUP BY w.user_id, e.exercise_name;


-- For frequently queried fields
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);