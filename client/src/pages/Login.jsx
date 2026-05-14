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
  Divider,
} from '@mui/material';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext';

const Login = () => {
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const redirectTo = location.state?.from || '/problem';

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);

  const submit = async (e) => {
    e.preventDefault();
    setError('');
    setBusy(true);
    try {
      await login(username.trim(), password);
      navigate(redirectTo, { replace: true });
    } catch (err) {
      setError(err.message || 'Login failed');
    } finally {
      setBusy(false);
    }
  };

  return (
    <Container maxWidth="xs" sx={{ mt: 8 }}>
      <Paper sx={{ p: 4 }} elevation={3}>
        <Typography variant="h4" align="center" gutterBottom>
          Sign in
        </Typography>
        <Typography variant="body2" align="center" color="text.secondary" sx={{ mb: 3 }}>
          to HardCode
        </Typography>

        <Box component="form" onSubmit={submit} sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          {error && <Alert severity="error">{error}</Alert>}
          <TextField
            label="Username or email"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
            autoFocus
            fullWidth
          />
          <TextField
            label="Password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            fullWidth
          />
          <Button type="submit" variant="contained" disabled={busy} size="large">
            {busy ? 'Signing in…' : 'Sign in'}
          </Button>
        </Box>

        <Divider sx={{ my: 3 }}>or</Divider>

        <Typography variant="body2" align="center">
          New here?{' '}
          <MuiLink component={Link} to="/signup">
            Create an account
          </MuiLink>
        </Typography>

        <Box
          sx={{
            mt: 3,
            p: 2,
            backgroundColor: '#f5f5f5',
            borderRadius: 1,
            fontSize: 13,
          }}
        >
          <strong>Try the demo accounts:</strong>
          <ul style={{ margin: '6px 0 0 18px' }}>
            <li>
              <code>admin</code> / <code>admin123</code> &nbsp;(admin — can create problems)
            </li>
            <li>
              <code>demo</code> / <code>demo123</code> &nbsp;(regular user)
            </li>
          </ul>
        </Box>
      </Paper>
    </Container>
  );
};

export default Login;
