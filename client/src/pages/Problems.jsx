import React, { useEffect, useMemo, useState } from 'react';
import {
  Box,
  Button,
  Chip,
  CircularProgress,
  Container,
  FormControl,
  IconButton,
  InputAdornment,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Typography,
  tableCellClasses,
} from '@mui/material';
import { styled } from '@mui/material/styles';
import SearchIcon from '@mui/icons-material/Search';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import RadioButtonUncheckedIcon from '@mui/icons-material/RadioButtonUnchecked';
import AddIcon from '@mui/icons-material/Add';
import { Link, useNavigate } from 'react-router-dom';
import NavBar from '../components/NavBar';
import { DifficultyBadge } from '../components/Badges';
import { api } from '../api';
import { useAuth } from '../auth/AuthContext';

const StyledHeadCell = styled(TableCell)(() => ({
  [`&.${tableCellClasses.head}`]: {
    backgroundColor: '#1e1e2f',
    color: '#fff',
    fontWeight: 600,
    fontSize: 13,
    letterSpacing: '.05em',
    textTransform: 'uppercase',
  },
}));

const StyledRow = styled(TableRow)(() => ({
  '&:nth-of-type(odd)': { backgroundColor: '#fafafa' },
  '&:hover': {
    backgroundColor: '#e8f3f0',
    cursor: 'pointer',
  },
}));

const Problems = () => {
  const navigate = useNavigate();
  const { isAdmin } = useAuth();

  const [problems, setProblems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [search, setSearch] = useState('');
  const [difficulty, setDifficulty] = useState('all');
  const [statusFilter, setStatusFilter] = useState('all');

  useEffect(() => {
    let cancelled = false;
    api
      .get('/api/v1/problems')
      .then((data) => {
        if (!cancelled) setProblems(Array.isArray(data) ? data : []);
      })
      .catch((err) => {
        if (!cancelled) setError(err.message || 'Failed to fetch problems');
      })
      .finally(() => !cancelled && setLoading(false));
    return () => {
      cancelled = true;
    };
  }, []);

  const filtered = useMemo(() => {
    const needle = search.trim().toLowerCase();
    return problems.filter((p) => {
      if (difficulty !== 'all' && p.difficulty !== difficulty) return false;
      if (statusFilter === 'solved' && !p.status) return false;
      if (statusFilter === 'unsolved' && p.status) return false;
      if (needle) {
        const hay = [
          p.name,
          ...(p.tags || []),
          p.author,
          String(p.id),
        ]
          .filter(Boolean)
          .join(' ')
          .toLowerCase();
        if (!hay.includes(needle)) return false;
      }
      return true;
    });
  }, [problems, search, difficulty, statusFilter]);

  const stats = useMemo(() => {
    const total = problems.length;
    const solved = problems.filter((p) => p.status).length;
    return { total, solved };
  }, [problems]);

  return (
    <>
      <NavBar />
      <Container maxWidth="lg" sx={{ mt: 4 }}>
        <Stack
          direction={{ xs: 'column', md: 'row' }}
          justifyContent="space-between"
          alignItems={{ xs: 'flex-start', md: 'center' }}
          spacing={2}
          sx={{ mb: 3 }}
        >
          <Box>
            <Typography variant="h4" fontWeight={700}>
              Problems
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Solved {stats.solved} / {stats.total}
            </Typography>
          </Box>
          {isAdmin && (
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              component={Link}
              to="/admin/problem/create"
              sx={{ backgroundColor: '#58A399' }}
            >
              New problem
            </Button>
          )}
        </Stack>

        <Paper sx={{ p: 2, mb: 2 }} elevation={1}>
          <Stack direction={{ xs: 'column', md: 'row' }} spacing={2}>
            <TextField
              placeholder="Search by title, tag, author…"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              size="small"
              fullWidth
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon fontSize="small" />
                  </InputAdornment>
                ),
              }}
            />
            <FormControl size="small" sx={{ minWidth: 160 }}>
              <InputLabel>Difficulty</InputLabel>
              <Select
                label="Difficulty"
                value={difficulty}
                onChange={(e) => setDifficulty(e.target.value)}
              >
                <MenuItem value="all">All</MenuItem>
                <MenuItem value="easy">Easy</MenuItem>
                <MenuItem value="medium">Medium</MenuItem>
                <MenuItem value="hard">Hard</MenuItem>
              </Select>
            </FormControl>
            <FormControl size="small" sx={{ minWidth: 160 }}>
              <InputLabel>Status</InputLabel>
              <Select
                label="Status"
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
              >
                <MenuItem value="all">All</MenuItem>
                <MenuItem value="solved">Solved</MenuItem>
                <MenuItem value="unsolved">Unsolved</MenuItem>
              </Select>
            </FormControl>
          </Stack>
        </Paper>

        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', mt: 6 }}>
            <CircularProgress />
          </Box>
        ) : error ? (
          <Typography color="error">{error}</Typography>
        ) : (
          <TableContainer component={Paper} elevation={1}>
            <Table>
              <TableHead>
                <TableRow>
                  <StyledHeadCell sx={{ width: 60 }}>Status</StyledHeadCell>
                  <StyledHeadCell sx={{ width: 60 }}>#</StyledHeadCell>
                  <StyledHeadCell>Title</StyledHeadCell>
                  <StyledHeadCell>Tags</StyledHeadCell>
                  <StyledHeadCell sx={{ width: 110 }}>Difficulty</StyledHeadCell>
                  <StyledHeadCell sx={{ width: 140 }}>Author</StyledHeadCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {filtered.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} align="center" sx={{ py: 4, color: 'text.secondary' }}>
                      No problems match your filters.
                    </TableCell>
                  </TableRow>
                ) : (
                  filtered.map((p) => (
                    <StyledRow key={p.id} onClick={() => navigate(`/problem/${p.id}`)} hover>
                      <TableCell>
                        {p.status ? (
                          <CheckCircleIcon sx={{ color: '#1E7F37' }} />
                        ) : (
                          <RadioButtonUncheckedIcon sx={{ color: '#bbb' }} />
                        )}
                      </TableCell>
                      <TableCell>{p.id}</TableCell>
                      <TableCell sx={{ fontWeight: 500 }}>{p.name}</TableCell>
                      <TableCell>
                        <Stack direction="row" spacing={0.5} flexWrap="wrap" useFlexGap>
                          {(p.tags || []).map((t) => (
                            <Chip key={t} label={t} size="small" variant="outlined" />
                          ))}
                        </Stack>
                      </TableCell>
                      <TableCell>
                        <DifficultyBadge difficulty={p.difficulty} />
                      </TableCell>
                      <TableCell sx={{ color: 'text.secondary' }}>
                        {p.author || '—'}
                      </TableCell>
                    </StyledRow>
                  ))
                )}
              </TableBody>
            </Table>
          </TableContainer>
        )}
      </Container>
    </>
  );
};

export default Problems;
