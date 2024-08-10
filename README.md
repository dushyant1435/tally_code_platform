# tally_code_platform
create table Problems (
	id int primary key not null, 
	user_id int not null,
	name varchar(50) not null, 
	description text not null, 
	constraints text, 
	input_format text, 
	output_format text, 
);

create table testcases(
	id int not null,
	input text not null,
	output text not null,
	sample bool
);


CREATE TABLE submission (
    id int not null,
    user_id INT NOT NULL
);


-- Test cases for problem id 1 (Sum of Two Numbers)
INSERT INTO testcases (id, input, output, sample) VALUES
(1, '3 5', '8', true),
(1, '10 20', '30', true),
(1, '100 200', '300', true);

-- Test cases for problem id 2 (Palindrome Check)
INSERT INTO testcases (id, input, output, sample) VALUES
(2, 'madam', 'Yes', true),
(2, 'hello', 'No', true),
(2, 'racecar', 'Yes', true);

-- Test cases for problem id 3 (Maximum Element in Array)
INSERT INTO testcases (id, input, output, sample) VALUES
(3, '5 1 2 3 4 5', '5', true),
(3, '3 10 20 30', '30', true),
(3, '4 7 8 9 6', '9', true);

-- Test cases for problem id 4 (Factorial Calculation)
INSERT INTO testcases (id, input, output, sample) VALUES
(4, '5', '120', true),
(4, '0', '1', true),
(4, '7', '5040', true);

-- Test cases for problem id 5 (Prime Number Test)
INSERT INTO testcases (id, input, output, sample) VALUES
(5, '2', 'Prime', true),
(5, '4', 'Not Prime', true),
(5, '13', 'Prime', true);

-- Test cases for problem id 6 (Example Problem)
INSERT INTO testcases (id, input, output, sample) VALUES
(6, '1', '42', true),
(6, '500', '42', true),
(6, '1000', '42', true);



curl -X POST http://localhost:8080/api/v1/runCode \
-H "Content-Type: application/json" \
-d '{
  "id": 1,
  "code": "a, b = map(int, input().split()); print(a + b)",
  "user_id": 101
}'

