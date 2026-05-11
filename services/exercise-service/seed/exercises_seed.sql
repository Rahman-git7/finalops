-- Seed muscle groups
INSERT INTO exercise.muscle_groups (name) VALUES
  ('chest'), ('back'), ('legs'), ('shoulders'), ('biceps'),
  ('triceps'), ('core'), ('glutes'), ('hamstrings'), ('calves')
ON CONFLICT (name) DO NOTHING;

-- Seed categories
INSERT INTO exercise.categories (name) VALUES
  ('compound'), ('isolation'), ('cardio'), ('bodyweight')
ON CONFLICT (name) DO NOTHING;

-- Seed exercises
INSERT INTO exercise.exercises (name, category_id, primary_muscle, secondary_muscle, equipment, description) VALUES
  -- Chest
  ('Développé couché', (SELECT id FROM exercise.categories WHERE name='compound'),
   (SELECT id FROM exercise.muscle_groups WHERE name='chest'),
   (SELECT id FROM exercise.muscle_groups WHERE name='triceps'),
   'barbell', 'Exercice de base pour la poitrine avec une barre'),
  ('Développé couché haltères', (SELECT id FROM exercise.categories WHERE name='compound'),
   (SELECT id FROM exercise.muscle_groups WHERE name='chest'),
   (SELECT id FROM exercise.muscle_groups WHERE name='triceps'),
   'dumbbell', 'Développé couché avec haltères pour plus d''amplitude'),
  ('Dips', (SELECT id FROM exercise.categories WHERE name='bodyweight'),
   (SELECT id FROM exercise.muscle_groups WHERE name='chest'),
   (SELECT id FROM exercise.muscle_groups WHERE name='triceps'),
   'bodyweight', 'Tractions en appui sur les barres parallèles'),
  ('Écarté couché', (SELECT id FROM exercise.categories WHERE name='isolation'),
   (SELECT id FROM exercise.muscle_groups WHERE name='chest'),
   NULL, 'dumbbell', 'Isolation de la poitrine avec haltères'),

  -- Back
  ('Tractions', (SELECT id FROM exercise.categories WHERE name='bodyweight'),
   (SELECT id FROM exercise.muscle_groups WHERE name='back'),
   (SELECT id FROM exercise.muscle_groups WHERE name='biceps'),
   'bodyweight', 'Tractions à la barre fixe, prise pronation ou supination'),
  ('Rowing barre', (SELECT id FROM exercise.categories WHERE name='compound'),
   (SELECT id FROM exercise.muscle_groups WHERE name='back'),
   (SELECT id FROM exercise.muscle_groups WHERE name='biceps'),
   'barbell', 'Tirage horizontal avec barre pour le dos'),
  ('Tirage vertical poulie', (SELECT id FROM exercise.categories WHERE name='isolation'),
   (SELECT id FROM exercise.muscle_groups WHERE name='back'),
   (SELECT id FROM exercise.muscle_groups WHERE name='biceps'),
   'cable', 'Tirage vertical à la poulie haute'),
  ('Soulevé de terre', (SELECT id FROM exercise.categories WHERE name='compound'),
   (SELECT id FROM exercise.muscle_groups WHERE name='back'),
   (SELECT id FROM exercise.muscle_groups WHERE name='legs'),
   'barbell', 'Exercice roi du dos et des jambes'),

  -- Legs
  ('Squat', (SELECT id FROM exercise.categories WHERE name='compound'),
   (SELECT id FROM exercise.muscle_groups WHERE name='legs'),
   (SELECT id FROM exercise.muscle_groups WHERE name='glutes'),
   'barbell', 'Le roi des exercices de jambes'),
  ('Leg press', (SELECT id FROM exercise.categories WHERE name='compound'),
   (SELECT id FROM exercise.muscle_groups WHERE name='legs'),
   (SELECT id FROM exercise.muscle_groups WHERE name='glutes'),
   'machine', 'Presse à cuisses'),
  ('Romanian deadlift', (SELECT id FROM exercise.categories WHERE name='compound'),
   (SELECT id FROM exercise.muscle_groups WHERE name='hamstrings'),
   (SELECT id FROM exercise.muscle_groups WHERE name='glutes'),
   'barbell', 'Soulevé de terre roumain pour ischios et fessiers'),
  ('Fentes', (SELECT id FROM exercise.categories WHERE name='compound'),
   (SELECT id FROM exercise.muscle_groups WHERE name='legs'),
   (SELECT id FROM exercise.muscle_groups WHERE name='glutes'),
   'bodyweight', 'Fentes avant ou arrière'),

  -- Shoulders
  ('Développé militaire', (SELECT id FROM exercise.categories WHERE name='compound'),
   (SELECT id FROM exercise.muscle_groups WHERE name='shoulders'),
   (SELECT id FROM exercise.muscle_groups WHERE name='triceps'),
   'barbell', 'Développé épaules debout ou assis'),
  ('Élévations latérales', (SELECT id FROM exercise.categories WHERE name='isolation'),
   (SELECT id FROM exercise.muscle_groups WHERE name='shoulders'),
   NULL, 'dumbbell', 'Isolation du deltoïde latéral'),

  -- Arms
  ('Curl biceps barre', (SELECT id FROM exercise.categories WHERE name='isolation'),
   (SELECT id FROM exercise.muscle_groups WHERE name='biceps'),
   NULL, 'barbell', 'Curl biceps à la barre droite ou EZ'),
  ('Curl haltères', (SELECT id FROM exercise.categories WHERE name='isolation'),
   (SELECT id FROM exercise.muscle_groups WHERE name='biceps'),
   NULL, 'dumbbell', 'Curl biceps avec haltères'),
  ('Extension triceps poulie', (SELECT id FROM exercise.categories WHERE name='isolation'),
   (SELECT id FROM exercise.muscle_groups WHERE name='triceps'),
   NULL, 'cable', 'Triceps pushdown à la poulie'),
  ('Skull crushers', (SELECT id FROM exercise.categories WHERE name='isolation'),
   (SELECT id FROM exercise.muscle_groups WHERE name='triceps'),
   NULL, 'barbell', 'Extension triceps allongé avec barre EZ'),

  -- Core
  ('Gainage', (SELECT id FROM exercise.categories WHERE name='bodyweight'),
   (SELECT id FROM exercise.muscle_groups WHERE name='core'),
   NULL, 'bodyweight', 'Planche frontale isométrique'),
  ('Crunch', (SELECT id FROM exercise.categories WHERE name='isolation'),
   (SELECT id FROM exercise.muscle_groups WHERE name='core'),
   NULL, 'bodyweight', 'Abdominaux classiques')
ON CONFLICT (name) DO NOTHING;
