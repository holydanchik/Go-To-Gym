CREATE TABLE IF NOT EXISTS workouts
(
    id              SERIAL PRIMARY KEY,
    created_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name            VARCHAR(255)                NOT NULL,
    description     TEXT,
    exercises       TEXT[]                      not null,
    calories_burned INTEGER,
    version         integer                     NOT NULL DEFAULT 1
);

INSERT INTO workouts (name, description, exercises, calories_burned)
VALUES ('Full Body Strength Training',
        'This workout targets all major muscle groups to build strength and endurance.',
        ARRAY['Squats', 'Push-ups', 'Rows', 'Lunges', 'Overhead press'],
        400);

INSERT INTO workouts (name, description, exercises, calories_burned)
VALUES ('Cardio HIIT',
        'High-Intensity Interval Training to improve cardiovascular health and burn calories.',
        ARRAY['Jumping jacks', 'Burpees', 'Mountain climbers', 'High knees', 'Jumping rope'],
        350);

INSERT INTO workouts (name, description, exercises, calories_burned)
VALUES ('Yoga for Flexibility',
        'A gentle yoga flow to improve flexibility, balance, and core strength.',
        ARRAY['Downward-Facing Dog', 'Warrior Pose', 'Triangle Pose', 'Cat-Cow', 'Childs Pose'],
        250);

INSERT INTO workouts (name, description, exercises, calories_burned)
VALUES ('Back and Bicep Burner',
        'Target your back and biceps with this intense workout.',
        ARRAY['Pull-ups', 'Rows', 'Bicep curls', 'Hammer curls'],
        300);

INSERT INTO workouts (name, description, exercises, calories_burned)
VALUES ('Legs and Core Challenge',
        'Strengthen your legs and core with this effective workout.',
        ARRAY['Squats', 'Lunges', 'Leg press', 'Plank', 'Crunches'],
        275);

INSERT INTO workouts (name, description, exercises, calories_burned)
VALUES ('Full Body Cardio Blast',
        'Get your heart rate up and burn calories with this full-body cardio workout.',
        ARRAY['Jumping jacks', 'Burpees', 'High knees', 'Mountain climbers', 'Running in place'],
        450);