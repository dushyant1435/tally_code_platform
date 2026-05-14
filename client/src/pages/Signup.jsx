import React, { useState } from 'react';
import {
  Box,
  Button,
  Container,
  Paper,
  TextField,
  Typography,
  Alert,
  Link as MuiLink,
} from '@mui/material';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext';

const Signup = () => {
  const { signup } = useAuth();
  const navigate = useNavigate();

  const [form, setForm] = useState({ username: '', email: '', password: '' });
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);

  const onChange = (e) =>
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));

  const submit = async (e) => {
    e.preventDefault();
    setError('');
    setBusy(true);
    try {
      await signup(form.username.trim(), form.email.trim(), form.password);
      navigate('/problem', { replace: true });
    } catch (err) {
      setError(err.message || 'Signup failed');
    } finally {
      setBusy(false);
    }
  };

  return (
    <Container maxWidth="xs" sx={{ mt: 8 }}>
      <Paper sx={{ p: 4 }} elevation={3}>
        <Typography variant="h4" align="center" gutterBottom>
          Create account
        </Typography>
        <Typography variant="body2" align="center" color="text.secondary" sx={{ mb: 3 }}>
          on TALLY Code Platform
        </Typography>

        <Box component="form" onSubmit={submit} sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          {error && <Alert severity="error">{error}</Alert>}
          <TextField
            label="Username"
            name="username"
            value={form.username}
            onChange={onChange}
            required
            autoFocus
            fullWidth
            helperText="3-30 chars, letters / digits / underscore"
          />
          <TextField
            label="Email"
            name="email"
            type="email"
            value={form.email}
            onChange={onChange}
            required
            fullWidth
          />
          <TextField
            label="Password"
            name="password"
            type="password"
            value={form.password}
            onChange={onChange}
            required
            fullWidth
            helperText="Minimum 6 characters"
          />
          <Button type="submit" variant="contained" disabled={busy} size="large">
            {busy ? 'Creating account…' : 'Create account'}
          </Button>
        </Box>

        <Typography variant="body2" align="center" sx={{ mt: 3 }}>
          Already have an account?{' '}
          <MuiLink component={Link} to="/login">
            Sign in
          </MuiLink>
        </Typography>
      </Paper>
    </Container>
  );
};

export default Signup;
