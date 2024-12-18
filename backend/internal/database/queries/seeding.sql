-- name: CheckInitialSeed :one
SELECT EXISTS (
  SELECT 1 FROM app_state WHERE key = 'initial_seed_completed'
);

-- name: MarkInitialSeed :exec
INSERT INTO app_state (key, value) VALUES ('initial_seed_completed', 'true');