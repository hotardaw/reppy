CREATE TABLE IF NOT EXISTS app_state (
    key TEXT PRIMARY KEY,
    value TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO workouts (user_id, workout_date, title) VALUES
    (1, CURRENT_DATE, 'Upper Body Strength Day');

INSERT INTO exercises (exercise_name, description) VALUES
    ('Bench Press', 'A compound upper body exercise targeting the chest, shoulders, and triceps'),
    ('Squat', 'A compound lower body exercise targeting the quadriceps, hamstrings, and glutes'),
    ('Deadlift', 'A compound full body exercise targeting the back, glutes, and hamstrings'),
    ('Pull-up', 'An upper body exercise targeting the back and biceps'),
    ('Push-up', 'A bodyweight exercise targeting the chest, shoulders, and triceps');

-- Get the workout_id we just created - we'll need it for the sets
DO $$ 
DECLARE 
    new_workout_id INTEGER;
BEGIN
    SELECT workout_id INTO new_workout_id FROM workouts 
    WHERE user_id = 1 AND workout_date = CURRENT_DATE;

    -- Add sets for bench press
    INSERT INTO workout_sets (workout_id, exercise_id, reps, resistance_type, resistance_value) VALUES
    (new_workout_id, 1, 8, 'weight', 135),
    (new_workout_id, 1, 8, 'weight', 135),
    (new_workout_id, 1, 6, 'weight', 155);

    -- Add sets for pull-ups
    INSERT INTO workout_sets (workout_id, exercise_id, reps, resistance_type, resistance_detail) VALUES
    (new_workout_id, 4, 10, 'bodyweight', NULL),
    (new_workout_id, 4, 8, 'bodyweight', NULL),
    (new_workout_id, 4, 6, 'bodyweight', NULL);

    -- Add sets for push-ups
    INSERT INTO workout_sets (workout_id, exercise_id, reps, resistance_type, resistance_detail) VALUES
    (new_workout_id, 5, 15, 'bodyweight', NULL),
    (new_workout_id, 5, 12, 'bodyweight', NULL),
    (new_workout_id, 5, 10, 'bodyweight', NULL);
END $$;