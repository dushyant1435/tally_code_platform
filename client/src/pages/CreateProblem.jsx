import React, { useState } from 'react';
import { TextField, Button, Box, Typography, Container } from '@mui/material';
import NavBar from '../components/NavBar';
import { useNavigate } from 'react-router-dom';
import { API_BASE, CURRENT_USER_ID } from '../config';

const CreateProblem = () => {
  const navigate = useNavigate();
  const [problem, setProblem] = useState({
    name: '',
    description: '',
    constraints: '',
    inputFormat: '',
    outputFormat: '',
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setProblem((prev) => ({ ...prev, [name]: value }));
  };

  const postProblem = async () => {
    try {
      const response = await fetch(`${API_BASE}/api/v1/newproblem`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: problem.name,
          constraints: problem.constraints,
          description: problem.description,
          input_format: problem.inputFormat,
          output_format: problem.outputFormat,
          user_id: CURRENT_USER_ID,
        }),
      });
      if (!response.ok) {
        const err = await response.json().catch(() => ({}));
        alert(`Failed to create problem: ${err.error || response.status}`);
        return;
      }
      const data = await response.json();
      if (data?.id) {
        navigate(`/problem/${data.id}/testcase`);
      }
    } catch (err) {
      console.error('postProblem', err);
      alert('Failed to reach server.');
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    postProblem();
  };

  return (
    <>
      <NavBar />
      <Container maxWidth="sm">
        <Box
          component="form"
          onSubmit={handleSubmit}
          sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 5 }}
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
            multiline
            rows={2}
            fullWidth
          />
          <TextField
            label="Input Format"
            name="inputFormat"
            value={problem.inputFormat}
            onChange={handleChange}
            multiline
            rows={2}
            fullWidth
          />
          <TextField
            label="Output Format"
            name="outputFormat"
            value={problem.outputFormat}
            onChange={handleChange}
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
