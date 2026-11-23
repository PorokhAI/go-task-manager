ALTER TABLE tasks
ADD CONSTRAINT tasks_priority_check
CHECK (priority BETWEEN 0 AND 2);