DELETE FROM archive.documents WHERE data @> '[{"title": "Документ 1"}]'::jsonb
   OR data @> '[{"title": "Документ 2"}]'::jsonb
   OR data @> '[{"title": "Документ 3"}]'::jsonb
   OR data @> '[{"title": "Документ 4"}]'::jsonb
   OR data @> '[{"title": "Документ 5"}]'::jsonb
   OR data @> '[{"title": "Документ 6"}]'::jsonb
   OR data @> '[{"title": "Документ 7"}]'::jsonb
   OR data @> '[{"title": "Документ 8"}]'::jsonb
   OR data @> '[{"title": "Документ 9"}]'::jsonb
   OR data @> '[{"title": "Документ 10"}]'::jsonb;

