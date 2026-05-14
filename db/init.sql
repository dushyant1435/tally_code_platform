-- Schema for tally_code_platform.
-- Loaded by Postgres on first container start via
-- /docker-entrypoint-initdb.d (see docker-compose.yml).

-- ---------------------------------------------------------------------------
-- Users + roles
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(50)  UNIQUE NOT NULL,
    email         VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT         NOT NULL,
    role          VARCHAR(20)  NOT NULL DEFAULT 'user'
                  CHECK (role IN ('user', 'admin')),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Seeded accounts (bcrypt cost 10, generated with cmd/hashgen):
--   admin / admin123   (role=admin)
--   demo  / demo123    (role=user)
INSERT INTO users (id, username, email, password_hash, role) VALUES
(1, 'admin', 'admin@tally.local',
    '$2a$10$ral5/rOEQfWvPDqMYIFmGOvUn1JxN0Jw2FvJ7gr.JzHDtKHHYFzRa', 'admin'),
(2, 'demo',  'demo@tally.local',
    '$2a$10$gf0AtcqSg1Mz1vYFj0w8ee1FmJN21lAc7fDiCH./cDUkN1yLm7rPe', 'user')
ON CONFLICT (id) DO NOTHING;
SELECT setval('users_id_seq', GREATEST((SELECT MAX(id) FROM users), 1));

-- ---------------------------------------------------------------------------
-- Problems + difficulty + tags
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS problems (
    id            SERIAL PRIMARY KEY,
    user_id       INT          NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    name          VARCHAR(100) NOT NULL,
    description   TEXT         NOT NULL,
    constraints   TEXT,
    input_format  TEXT,
    output_format TEXT,
    difficulty    VARCHAR(10)  NOT NULL DEFAULT 'easy'
                  CHECK (difficulty IN ('easy', 'medium', 'hard')),
    tags          TEXT[]       NOT NULL DEFAULT '{}',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

INSERT INTO problems (id, user_id, name, description, constraints,
                      input_format, output_format, difficulty, tags) VALUES
(1, 1, 'Sum of Two Numbers',
    'Read two integers a and b on one line separated by a space and print their sum.',
    '-10^9 <= a, b <= 10^9',
    'A single line with two integers a and b separated by a space.',
    'A single integer: a + b.',
    'easy', ARRAY['math', 'warmup']),
(2, 1, 'Palindrome Check',
    'Given a lowercase string s, print "Yes" if it is a palindrome and "No" otherwise.',
    '1 <= |s| <= 100',
    'A single lowercase string.',
    '"Yes" or "No".',
    'easy', ARRAY['string']),
(3, 1, 'Maximum Element in Array',
    'Given an integer n followed by n integers on the same line, print the maximum.',
    '1 <= n <= 1000',
    'First integer n, followed by n integers, all on one line.',
    'The maximum integer.',
    'easy', ARRAY['array']),
(4, 1, 'Factorial Calculation',
    'Read a non-negative integer n and print n!.',
    '0 <= n <= 12',
    'A single non-negative integer n.',
    'n! as a single integer.',
    'medium', ARRAY['math', 'recursion']),
(5, 1, 'Prime Number Test',
    'Read an integer n and print "Prime" if n is prime, otherwise "Not Prime".',
    '2 <= n <= 10^6',
    'A single integer n.',
    '"Prime" or "Not Prime".',
    'medium', ARRAY['math']),
(6, 1, 'Example Problem',
    'Read any integer and print 42. Used to demonstrate constant output.',
    '1 <= n <= 1000',
    'Any single integer.',
    'The integer 42.',
    'hard', ARRAY['fun'])
ON CONFLICT (id) DO NOTHING;
SELECT setval('problems_id_seq', GREATEST((SELECT MAX(id) FROM problems), 1));

-- ---------------------------------------------------------------------------
-- Test cases
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS testcases (
    id      INT     NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    input   TEXT    NOT NULL,
    output  TEXT    NOT NULL,
    sample  BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_testcases_id ON testcases(id);

INSERT INTO testcases (id, input, output, sample) VALUES
(1, '3 5',        '8',         true),
(1, '10 20',      '30',        true),
(1, '100 200',    '300',       true),

(2, 'madam',      'Yes',       true),
(2, 'hello',      'No',        true),
(2, 'racecar',    'Yes',       true),

(3, '5 1 2 3 4 5','5',         true),
(3, '3 10 20 30', '30',        true),
(3, '4 7 8 9 6',  '9',         true),

(4, '5',          '120',       true),
(4, '0',          '1',         true),
(4, '7',          '5040',      true),

(5, '2',          'Prime',     true),
(5, '4',          'Not Prime', true),
(5, '13',         'Prime',     true),

(6, '1',          '42',        true),
(6, '500',        '42',        true),
(6, '1000',       '42',        true);

-- ---------------------------------------------------------------------------
-- Submissions (full history, every run)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS submissions (
    id               SERIAL PRIMARY KEY,
    problem_id       INT          NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    user_id          INT          NOT NULL REFERENCES users(id)    ON DELETE CASCADE,
    language         VARCHAR(20)  NOT NULL DEFAULT 'python',
    code             TEXT         NOT NULL,
    status           VARCHAR(30)  NOT NULL,
        -- one of: accepted, wrong_answer, time_limit_exceeded,
        -- runtime_error, compilation_error, no_test_cases, server_error
    runtime_seconds  DOUBLE PRECISION,
    failed_test      INT,
    message          TEXT,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_submissions_user    ON submissions(user_id);
CREATE INDEX IF NOT EXISTS idx_submissions_problem ON submissions(problem_id);
