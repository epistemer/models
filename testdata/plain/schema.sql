-- A schema with no enums and no prefix, the second fixture for models'
-- generators (roadmap M1.2): the enums generator writes no file for it.
CREATE TABLE note (
  id bigint PRIMARY KEY,
  body text NOT NULL
);
