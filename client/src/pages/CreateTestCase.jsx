import React, { useState } from 'react';
import { TextField, Button, Box, Typography, Container, Checkbox, FormControlLabel } from '@mui/material';
import NavBar from '../components/NavBar';
import { useParams } from 'react-router-dom';

const CreateTestCase = () => {
  const { id } = useParams();
  let num = parseInt(id);

  const [testCase, setTestCase] = useState({
    id: num,
    input: '',
    output: '',
    sample: false,
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setTestCase({
      ...testCase,
      [name]: value,
    });
  };

  const handleCheckboxChange = (e) => {
    setTestCase({
      ...testCase,
      sample: e.target.checked,
    });
  };

  const postTestCase = async () => {
    const response = await fetch(`http://localhost:8080/api/v1/createTestCase`, {
      method: "POST",
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(testCase),
    });
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    console.log(testCase)
    postTestCase();
  };

  return (
    <>
      <NavBar />
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
            Add New Test Case
          </Typography>

          <TextField
            label="Input"
            name="input"
            value={testCase.input}
            onChange={handleChange}
            required
            fullWidth
          />
          <TextField
            label="Output"
            name="output"
            value={testCase.output}
            onChange={handleChange}
            required
            multiline
            rows={4}
            fullWidth
          />
          <FormControlLabel
            control={
              <Checkbox
                checked={testCase.sample}
                onChange={handleCheckboxChange}
                color="primary"
              />
            }
            label="Sample"
          />
          <Button type="submit" variant="contained" color="primary" fullWidth>
            CREATE TEST CASE
          </Button>
        </Box>
      </Container>
      <br />
    </>
  );
};

export default CreateTestCase;
