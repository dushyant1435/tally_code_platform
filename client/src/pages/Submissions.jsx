import React, { useEffect, useState } from 'react';
import {
  Box,
  CircularProgress,
  Container,
  Paper,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from '@mui/material';
import { Link, useNavigate } from 'react-router-dom';
import NavBar from '../components/NavBar';
import { StatusBadge } from '../components/Badges';
import { api } from '../api';
import { useAuth } from '../auth/AuthContext';

const Submissions = () => {
  const { isAdmin } = useAuth();
  const navigate = useNavigate();
  const [rows, setRows] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    let cancelled = false;
    const path = isAdmin ? '/api/v1/admin/submissions' : '/api/v1/submissions';
    api
      .get(path)
      .then((data) => {
        if (!cancelled) setRows(Array.isArray(data) ? data : []);
      })
      .catch((err) => !cancelled && setError(err.message || 'Failed to fetch submissions'))
      .finally(() => !cancelled && setLoading(false));
    return () => {
      cancelled = true;
    };
  }, [isAdmin]);

  return (
    <>
      <NavBar />
      <Container maxWidth="lg" sx={{ mt: 4 }}>
        <Typography variant="h4" fontWeight={700} sx={{ mb: 3 }}>
          {isAdmin ? 'All submissions (admin)' : 'My submissions'}
        </Typography>

        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', mt: 6 }}>
            <CircularProgress />
          </Box>
        ) : error ? (
          <Typography color="error">{error}</Typography>
        ) : rows.length === 0 ? (
          <Typography color="text.secondary">
            No submissions yet. Solve a problem to see it here.
          </Typography>
        ) : (
          <TableContainer component={Paper}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Status</TableCell>
                  <TableCell>Problem</TableCell>
                  {isAdmin && <TableCell>User</TableCell>}
                  <TableCell>Language</TableCell>
                  <TableCell>Runtime</TableCell>
                  <TableCell>When</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {rows.map((s) => (
                  <TableRow
                    key={s.id}
                    hover
                    sx={{ cursor: 'pointer' }}
                    onClick={() => navigate(`/submissions/${s.id}`)}
                  >
                    <TableCell>
                      <StatusBadge status={s.status} />
                    </TableCell>
                    <TableCell>
                      <Link
                        to={`/problem/${s.problem_id}`}
                        onClick={(e) => e.stopPropagation()}
                        style={{ color: '#1976d2', textDecoration: 'none' }}
                      >
                        {s.problem_id}. {s.problem_name}
                      </Link>
                    </TableCell>
                    {isAdmin && <TableCell>{s.username}</TableCell>}
                    <TableCell>{s.language}</TableCell>
                    <TableCell>
                      {s.runtime_seconds != null ? `${s.runtime_seconds.toFixed(3)}s` : '—'}
                    </TableCell>
                    <TableCell sx={{ color: 'text.secondary', fontSize: 13 }}>
                      {new Date(s.created_at).toLocaleString()}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        )}
      </Container>
    </>
  );
};

export default Submissions;
