CREATE TYPE mission_status AS ENUM ('in_progress', 'completed');

CREATE TYPE target_status AS ENUM ('in_progress', 'completed');

-- Spy Cats table
CREATE TABLE spy_cats (
      id SERIAL PRIMARY KEY,
      name VARCHAR(100) NOT NULL,
      years_of_experience INTEGER NOT NULL,
      breed VARCHAR(100) NOT NULL,
      salary DECIMAL(10, 2) NOT NULL,
      created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Missions table
CREATE TABLE missions (
      id SERIAL PRIMARY KEY,
      cat_id INTEGER REFERENCES spy_cats(id),
      status mission_status DEFAULT 'in_progress',
      created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Targets table
CREATE TABLE targets (
     id SERIAL PRIMARY KEY,
     mission_id INTEGER REFERENCES missions(id) ON DELETE CASCADE,
     name VARCHAR(100) NOT NULL,
     country VARCHAR(100) NOT NULL,
     notes TEXT,
     status target_status DEFAULT 'in_progress',
     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);