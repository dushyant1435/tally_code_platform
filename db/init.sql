-- Schema for tally_code_platform.
-- Loaded by Postgres on first container start via
-- /docker-entrypoint-initdb.d (see docker-compose.yml).

CREATE TABLE IF NOT EXISTS problems (
    id            SERIAL PRIMARY KEY,
    user_id       INT          NOT NULL,
    name          VARCHAR(100) NOT NULL,
    description   TEXT         NOT NULL,
    constraints   TEXT,
    input_format  TEXT,
    output_format TEXT
);

CREATE TABLE IF NOT EXISTS testcases (
    id      INT     NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    input   TEXT    NOT NULL,
    output  TEXT    NOT NULL,
    sample  BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_testcases_id ON testcases(id);

CREATE TABLE IF NOT EXISTS submission (
    id      INT NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    user_id INT NOT NULL,
    PRIMARY KEY (id, user_id)
);

-- ---------------------------------------------------------------------------
-- Seed data: six problems with three sample test cases each.
-- We insert with explicit ids and then realign the sequence so future inserts
-- via /api/v1/newproblem continue from 7.
-- ---------------------------------------------------------------------------

INSERT INTO problems (id, user_id, name, description, constraints, input_format, output_format) VALUES
(1, 123, 'Sum of Two Numbers',
    'Read two integers a and b on one line separated by a space and print their sum.',
    '-10^9 <= a, b <= 10^9',
    'A single line with two integers a and b separated by a space.',
    'A single integer: a + b.'),
(2, 123, 'Palindrome Check',
    'Given a lowercase string s, print "Yes" if it is a palindrome and "No" otherwise.',
    '1 <= |s| <= 100',
    'A single lowercase string.',
    '"Yes" or "No".'),
(3, 123, 'Maximum Element in Array',
    'Given an integer n followed by n integers on the same line, print the maximum.',
    '1 <= n <= 1000',
    'First integer n, followed by n integers, all on one line.',
    'The maximum integer.'),
(4, 123, 'Factorial Calculation',
    'Read a non-negative integer n and print n!.',
    '0 <= n <= 12',
    'A single non-negative integer n.',
    'n! as a single integer.'),
(5, 123, 'Prime Number Test',
    'Read an integer n and print "Prime" if n is prime, otherwise "Not Prime".',
    '2 <= n <= 10^6',
    'A single integer n.',
    '"Prime" or "Not Prime".'),
(6, 123, 'Example Problem',
    'Read any integer and print 42. Used to demonstrate constant output.',
    '1 <= n <= 1000',
    'Any single integer.',
    'The integer 42.')
ON CONFLICT (id) DO NOTHING;

SELECT setval('problems_id_seq', GREATEST((SELECT MAX(id) FROM problems), 1));

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
