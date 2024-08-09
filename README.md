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
	status bool
);