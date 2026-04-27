-- Move Start of A/A test to now - you will have to manually clear the redis cache as atm it will only clock this change at midnight
UPDATE prism.experiments set aa_start_time = now(), aa_end_time = now() + interval '7 day' WHERE id = '1af09025-c50b-4d4b-8634-91fbd62e7bc4';

-- End A/A Test

UPDATE prism.experiments set aa_end_time = now() WHERE id = '1af09025-c50b-4d4b-8634-91fbd62e7bc4'
