import React, { useState } from 'react';
import { TextField, Button, Box, Typography, Container } from '@mui/material';
import NavBar from '../components/NavBar';

const CreateProblem = () => {
  const [problem, setProblem] = useState({
    name: '',
    description: '',
    constraints: '',
    inputFormat: '',
    outputFormat: '',
  });

  const [pid,setPid] = useState()

  const handleChange = (e) => {
    const { name, value } = e.target;
    setProblem({
      ...problem,
      [name]: value,
    });
  };

  const postProblem = async () => {
    const response = await fetch(`http://localhost:8080/api/v1/newproblem`, {
          method: "POST",
          headers: {
            'Content-Type': 'application/json', // Set the content type to JSON
          },
        body: JSON.stringify({
            "name": problem.name,
            "constraints": problem.constraints,
            "description": problem.description,
            "input_format": problem.inputFormat,
            "output_format": problem.outputFormat,
            "user_id": 123,
        }),
    });
    const data = await response.json();
    setPid(data);
    console.log(data.id)
  };
  const handleSubmit = (e) => {
    e.preventDefault();
    // Handle form submission logic here
    postProblem()
    console.log(problem);
  };

  return (
    <>
    <NavBar></NavBar>
    <Container maxWidth="sm">
      <Box
        component="form"
        onSubmit={handleSubmit}
        sx={{
          display: 'flex',
          flexDirection: 'column',
          gap: 2,
          mt: 5,
        }}
      >
        <Typography variant="h4" component="h1" align="center" gutterBottom>
          Add New Problem
        </Typography>

        <TextField
          label="Problem Name"
          name="name"
          value={problem.name}
          onChange={handleChange}
          required
          fullWidth
        />
        <TextField
          label="Problem Description"
          name="description"
          value={problem.description}
          onChange={handleChange}
          required
          multiline
          rows={4}
          fullWidth
        />
        <TextField
          label="Constraints"
          name="constraints"
          value={problem.constraints}
          onChange={handleChange}
          required
          multiline
          rows={2}
          fullWidth
        />
        <TextField
          label="Input Format"
          name="inputFormat"
          value={problem.inputFormat}
          onChange={handleChange}
          required
          multiline
          rows={2}
          fullWidth
        />
        <TextField
          label="Output Format"
          name="outputFormat"
          value={problem.outputFormat}
          onChange={handleChange}
          required
          multiline
          rows={2}
          fullWidth
        />
        <Button type="submit" variant="contained" color="primary" fullWidth>
          CREATE CHALLENGE
        </Button>
      </Box>
    </Container>
    <br />
    </>
  );
};

export default CreateProblem;
