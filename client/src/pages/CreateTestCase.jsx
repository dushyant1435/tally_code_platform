import React, { useState } from 'react';
import {
  TextField,
  Button,
  Box,
  Typography,
  Container,
  Checkbox,
  FormControlLabel,
} from '@mui/material';
import NavBar from '../components/NavBar';
import { useParams } from 'react-router-dom';
import { API_BASE } from '../config';

const CreateTestCase = () => {
  const { id } = useParams();
  const num = parseInt(id, 10);

  const [testCase, setTestCase] = useState({
    id: num,
    input: '',
    output: '',
    sample: false,
  });
  const [statusMsg, setStatusMsg] = useState('');

  const handleChange = (e) => {
    const { name, value } = e.target;
    setTestCase((prev) => ({ ...prev, [name]: value }));
  };

  const handleCheckboxChange = (e) => {
    setTestCase((prev) => ({ ...prev, sample: e.target.checked }));
  };

  const postTestCase = async () => {
    try {
      const response = await fetch(`${API_BASE}/api/v1/createTestCase`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(testCase),
      });
      if (response.ok) {
        setStatusMsg('Test case created.');
        setTestCase({ id: num, input: '', output: '', sample: false });
      } else {
        const err = await response.json().catch(() => ({}));
        setStatusMsg(`Failed: ${err.error || response.status}`);
      }
    } catch (err) {
      console.error('postTestCase', err);
      setStatusMsg('Failed to reach server.');
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    postTestCase();
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
            Add New Test Case
          </Typography>

          <TextField
            label="Input"
            name="input"
            value={testCase.input}
            onChange={handleChange}
            required
            multiline
            rows={2}
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
          {statusMsg && (
            <Typography
              variant="body2"
              sx={{ color: statusMsg.startsWith('Test') ? 'green' : 'red' }}
            >
              {statusMsg}
            </Typography>
          )}
        </Box>
      </Container>
      <br />
    </>
  );
};

export default CreateTestCase;
