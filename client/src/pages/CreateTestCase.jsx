import React, { useState } from 'react';
import {
  Alert,
  Box,
  Button,
  Checkbox,
  Container,
  FormControlLabel,
  Paper,
  TextField,
  Typography,
} from '@mui/material';
import { useParams } from 'react-router-dom';
import NavBar from '../components/NavBar';
import { api } from '../api';

const CreateTestCase = () => {
  const { id } = useParams();
  const num = parseInt(id, 10);

  const [testCase, setTestCase] = useState({ id: num, input: '', output: '', sample: false });
  const [statusMsg, setStatusMsg] = useState('');
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);

  const onChange = (e) => {
    const { name, value } = e.target;
    setTestCase((prev) => ({ ...prev, [name]: value }));
  };

  const submit = async (e) => {
    e.preventDefault();
    setError('');
    setStatusMsg('');
    setBusy(true);
    try {
      await api.post('/api/v1/createTestCase', testCase);
      setStatusMsg('Test case created.');
      setTestCase({ id: num, input: '', output: '', sample: false });
    } catch (err) {
      setError(err.message || 'Failed to create test case');
    } finally {
      setBusy(false);
    }
  };

  return (
    <>
      <NavBar />
      <Container maxWidth="sm" sx={{ mt: 4 }}>
        <Paper sx={{ p: 4 }} elevation={2}>
          <Typography variant="h4" align="center" gutterBottom>
            Add test case
          </Typography>
          <Typography variant="body2" align="center" color="text.secondary" sx={{ mb: 2 }}>
            Problem #{num}
          </Typography>

          <Box component="form" onSubmit={submit} sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
            {error && <Alert severity="error">{error}</Alert>}
            {statusMsg && <Alert severity="success">{statusMsg}</Alert>}

            <TextField label="Input" name="input" value={testCase.input} onChange={onChange} required multiline rows={3} fullWidth />
            <TextField label="Expected output" name="output" value={testCase.output} onChange={onChange} required multiline rows={3} fullWidth />
            <FormControlLabel
              control={
                <Checkbox
                  checked={testCase.sample}
                  onChange={(e) => setTestCase((p) => ({ ...p, sample: e.target.checked }))}
                />
              }
              label="Sample (visible to users on the problem page)"
            />
            <Button type="submit" variant="contained" disabled={busy} size="large">
              {busy ? 'Saving…' : 'Add test case'}
            </Button>
          </Box>
        </Paper>
      </Container>
    </>
  );
};

export default CreateTestCase;
