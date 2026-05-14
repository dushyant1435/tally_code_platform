import React, { useState } from 'react';
import {
  Alert,
  Box,
  Button,
  Chip,
  Container,
  FormControl,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import NavBar from '../components/NavBar';
import { api } from '../api';

const CreateProblem = () => {
  const navigate = useNavigate();
  const [problem, setProblem] = useState({
    name: '',
    description: '',
    constraints: '',
    inputFormat: '',
    outputFormat: '',
    difficulty: 'easy',
  });
  const [tags, setTags] = useState([]);
  const [tagInput, setTagInput] = useState('');
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);

  const onChange = (e) => {
    const { name, value } = e.target;
    setProblem((prev) => ({ ...prev, [name]: value }));
  };

  const addTag = () => {
    const t = tagInput.trim().toLowerCase();
    if (t && !tags.includes(t)) setTags([...tags, t]);
    setTagInput('');
  };

  const removeTag = (t) => setTags(tags.filter((x) => x !== t));

  const submit = async (e) => {
    e.preventDefault();
    setError('');
    setBusy(true);
    try {
      const data = await api.post('/api/v1/newproblem', {
        name: problem.name,
        description: problem.description,
        constraints: problem.constraints || null,
        input_format: problem.inputFormat || null,
        output_format: problem.outputFormat || null,
        difficulty: problem.difficulty,
        tags,
      });
      if (data?.id) {
        navigate(`/admin/problem/${data.id}/testcase`);
      }
    } catch (err) {
      setError(err.message || 'Failed to create problem');
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
            New problem
          </Typography>
          <Box component="form" onSubmit={submit} sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}>
            {error && <Alert severity="error">{error}</Alert>}
            <TextField label="Problem name" name="name" value={problem.name} onChange={onChange} required fullWidth />
            <TextField
              label="Description"
              name="description"
              value={problem.description}
              onChange={onChange}
              required
              multiline
              rows={4}
              fullWidth
            />
            <FormControl fullWidth>
              <InputLabel>Difficulty</InputLabel>
              <Select
                label="Difficulty"
                name="difficulty"
                value={problem.difficulty}
                onChange={onChange}
              >
                <MenuItem value="easy">Easy</MenuItem>
                <MenuItem value="medium">Medium</MenuItem>
                <MenuItem value="hard">Hard</MenuItem>
              </Select>
            </FormControl>

            <Box>
              <Stack direction="row" spacing={1}>
                <TextField
                  label="Add tag"
                  size="small"
                  value={tagInput}
                  onChange={(e) => setTagInput(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                      e.preventDefault();
                      addTag();
                    }
                  }}
                  fullWidth
                />
                <Button onClick={addTag} variant="outlined">
                  Add
                </Button>
              </Stack>
              <Stack direction="row" spacing={0.5} flexWrap="wrap" useFlexGap sx={{ mt: 1 }}>
                {tags.map((t) => (
                  <Chip key={t} label={t} size="small" onDelete={() => removeTag(t)} />
                ))}
              </Stack>
            </Box>

            <TextField label="Constraints" name="constraints" value={problem.constraints} onChange={onChange} multiline rows={2} fullWidth />
            <TextField label="Input format" name="inputFormat" value={problem.inputFormat} onChange={onChange} multiline rows={2} fullWidth />
            <TextField label="Output format" name="outputFormat" value={problem.outputFormat} onChange={onChange} multiline rows={2} fullWidth />
            <Button type="submit" variant="contained" disabled={busy} size="large">
              {busy ? 'Creating…' : 'Create problem'}
            </Button>
          </Box>
        </Paper>
      </Container>
    </>
  );
};

export default CreateProblem;
