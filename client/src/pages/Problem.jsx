import React, { useCallback, useEffect, useMemo, useState } from 'react';
import {
  Alert,
  Box,
  Button,
  Chip,
  CircularProgress,
  Container,
  Divider,
  Grid,
  Paper,
  Stack,
  Tab,
  Tabs,
  Typography,
} from '@mui/material';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import SendIcon from '@mui/icons-material/Send';
import { Link as RouterLink, useNavigate, useParams } from 'react-router-dom';
import NavBar from '../components/NavBar';
import CodeEditor from '../components/CodeEditor';
import { DifficultyBadge, StatusBadge } from '../components/Badges';
import { CODE_SNIPPETS } from '../constants';
import { api } from '../api';
import { useAuth } from '../auth/AuthContext';

const TabPanel = ({ value, index, children }) =>
  value === index ? <Box sx={{ pt: 2 }}>{children}</Box> : null;

const Problem = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { isAuthenticated, isAdmin } = useAuth();
  const num = useMemo(() => parseInt(id, 10), [id]);

  const [tab, setTab] = useState(0);
  const [code, setCode] = useState(CODE_SNIPPETS.python);
  const [language, setLanguage] = useState('python');
  const [problem, setProblem] = useState(null);
  const [samples, setSamples] = useState([]);
  const [sampleResults, setSampleResults] = useState([]);
  const [submissions, setSubmissions] = useState([]);
  const [runBusy, setRunBusy] = useState(false);
  const [submitBusy, setSubmitBusy] = useState(false);
  const [verdict, setVerdict] = useState(null);
  const [error, setError] = useState('');

  const loadProblem = useCallback(async () => {
    setError('');
    try {
      const [p, s] = await Promise.all([
        api.get(`/api/v1/problem/${id}`),
        api.get(`/api/v1/problem/${id}/sampleTestCases`),
      ]);
      setProblem(p);
      setSamples(Array.isArray(s) ? s : []);
    } catch (err) {
      setError(err.message || 'Failed to load problem');
    }
  }, [id]);

  const loadSubmissions = useCallback(async () => {
    if (!isAuthenticated) return;
    try {
      const data = await api.get(`/api/v1/problem/${id}/submissions`);
      setSubmissions(Array.isArray(data) ? data : []);
    } catch (err) {
      // 404 / 401 are OK here; just leave empty.
      console.warn('loadSubmissions', err.message);
    }
  }, [id, isAuthenticated]);

  useEffect(() => {
    loadProblem();
    loadSubmissions();
  }, [loadProblem, loadSubmissions]);

  const onRunSample = async () => {
    setRunBusy(true);
    try {
      const data = await api.post('/api/v1/runSampleCode', {
        id: num,
        code,
        language,
      });
      setSampleResults(Array.isArray(data.results) ? data.results : []);
    } catch (err) {
      setSampleResults([
        { input: '', expected: '', output: err.message, result: false, runtime: '-' },
      ]);
    } finally {
      setRunBusy(false);
    }
  };

  const onSubmit = async () => {
    if (!isAuthenticated) {
      navigate('/login', { state: { from: `/problem/${id}` } });
      return;
    }
    setSubmitBusy(true);
    setVerdict(null);
    try {
      const data = await api.post('/api/v1/runCode', {
        id: num,
        code,
        language,
      });
      setVerdict(data);
      loadSubmissions();
    } catch (err) {
      setVerdict({ status: 'server_error', message: err.message });
    } finally {
      setSubmitBusy(false);
    }
  };

  if (error) {
    return (
      <>
        <NavBar />
        <Container maxWidth="md" sx={{ mt: 4 }}>
          <Alert severity="error">{error}</Alert>
        </Container>
      </>
    );
  }

  if (!problem) {
    return (
      <>
        <NavBar />
        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 8 }}>
          <CircularProgress />
        </Box>
      </>
    );
  }

  return (
    <>
      <NavBar />
      <Container maxWidth="xl" sx={{ mt: 3, mb: 6 }}>
        <Grid container spacing={2}>
          {/* LEFT: description / submissions tabs */}
          <Grid item xs={12} md={6}>
            <Paper sx={{ p: 3 }} elevation={1}>
              <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 1 }}>
                <Typography variant="h5" fontWeight={700}>
                  {problem.id}. {problem.name}
                </Typography>
              </Stack>
              <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap sx={{ mb: 2 }}>
                <DifficultyBadge difficulty={problem.difficulty} />
                {(problem.tags || []).map((t) => (
                  <Chip key={t} label={t} size="small" variant="outlined" />
                ))}
              </Stack>

              <Tabs value={tab} onChange={(_, v) => setTab(v)}>
                <Tab label="Description" />
                <Tab label={`Submissions${submissions.length ? ` (${submissions.length})` : ''}`} />
              </Tabs>
              <Divider />

              <TabPanel value={tab} index={0}>
                <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap', mb: 2 }}>
                  {problem.description}
                </Typography>
                {problem.constraints && (
                  <>
                    <Typography variant="subtitle2">Constraints</Typography>
                    <Typography variant="body2" sx={{ whiteSpace: 'pre-wrap', mb: 2 }}>
                      {problem.constraints}
                    </Typography>
                  </>
                )}
                {problem.input_format && (
                  <>
                    <Typography variant="subtitle2">Input format</Typography>
                    <Typography variant="body2" sx={{ whiteSpace: 'pre-wrap', mb: 2 }}>
                      {problem.input_format}
                    </Typography>
                  </>
                )}
                {problem.output_format && (
                  <>
                    <Typography variant="subtitle2">Output format</Typography>
                    <Typography variant="body2" sx={{ whiteSpace: 'pre-wrap', mb: 2 }}>
                      {problem.output_format}
                    </Typography>
                  </>
                )}

                <Typography variant="subtitle2" sx={{ mt: 2 }}>
                  Examples
                </Typography>
                {samples.length === 0 ? (
                  <Typography variant="body2" color="text.secondary">
                    No sample test cases for this problem yet.
                  </Typography>
                ) : (
                  samples.map((tc, i) => (
                    <Paper
                      key={i}
                      variant="outlined"
                      sx={{ p: 2, mt: 1, backgroundColor: '#fafafa' }}
                    >
                      <Typography variant="caption" color="text.secondary">
                        Example {i + 1}
                      </Typography>
                      <Box component="pre" sx={{ m: 0, fontSize: 13 }}>
                        Input:  {tc.input}
                        {'\n'}Output: {tc.output}
                      </Box>
                    </Paper>
                  ))
                )}

                {isAdmin && (
                  <Box sx={{ mt: 3 }}>
                    <Button
                      variant="outlined"
                      component={RouterLink}
                      to={`/admin/problem/${id}/testcase`}
                    >
                      Add test cases
                    </Button>
                  </Box>
                )}
              </TabPanel>

              <TabPanel value={tab} index={1}>
                {!isAuthenticated ? (
                  <Typography variant="body2" color="text.secondary">
                    Sign in to see your submissions.
                  </Typography>
                ) : submissions.length === 0 ? (
                  <Typography variant="body2" color="text.secondary">
                    No submissions yet for this problem.
                  </Typography>
                ) : (
                  <Stack spacing={1}>
                    {submissions.map((s) => (
                      <Paper
                        key={s.id}
                        variant="outlined"
                        sx={{
                          p: 1.5,
                          display: 'flex',
                          alignItems: 'center',
                          gap: 2,
                        }}
                      >
                        <StatusBadge status={s.status} />
                        <Box sx={{ flex: 1 }}>
                          <Typography variant="body2">
                            {s.language} ·{' '}
                            {s.runtime_seconds != null
                              ? `${s.runtime_seconds.toFixed(3)}s`
                              : '—'}
                            {s.failed_test ? ` · failed test #${s.failed_test}` : ''}
                          </Typography>
                          <Typography variant="caption" color="text.secondary">
                            {new Date(s.created_at).toLocaleString()}
                          </Typography>
                        </Box>
                      </Paper>
                    ))}
                  </Stack>
                )}
              </TabPanel>
            </Paper>
          </Grid>

          {/* RIGHT: editor + run/submit */}
          <Grid item xs={12} md={6}>
            <Paper sx={{ p: 3 }} elevation={1}>
              <CodeEditor
                value={code}
                setValue={setCode}
                language={language}
                setLanguage={setLanguage}
              />

              <Stack direction="row" spacing={1} sx={{ mt: 1 }}>
                <Button
                  variant="outlined"
                  startIcon={runBusy ? <CircularProgress size={16} /> : <PlayArrowIcon />}
                  onClick={onRunSample}
                  disabled={runBusy}
                  fullWidth
                >
                  Run
                </Button>
                <Button
                  variant="contained"
                  color="success"
                  startIcon={submitBusy ? <CircularProgress size={16} color="inherit" /> : <SendIcon />}
                  onClick={onSubmit}
                  disabled={submitBusy}
                  fullWidth
                >
                  Submit
                </Button>
              </Stack>

              {verdict && (
                <Box sx={{ mt: 2 }}>
                  <Alert
                    severity={verdict.success ? 'success' : 'error'}
                    icon={false}
                  >
                    <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 0.5 }}>
                      <StatusBadge status={verdict.status} />
                      {verdict.totalRuntime != null && (
                        <Typography variant="caption" color="text.secondary">
                          {verdict.totalRuntime.toFixed(3)}s
                        </Typography>
                      )}
                      {verdict.failedAt > 0 && (
                        <Typography variant="caption" color="text.secondary">
                          failed test #{verdict.failedAt}
                        </Typography>
                      )}
                    </Stack>
                    {verdict.message && (
                      <Typography variant="caption" sx={{ whiteSpace: 'pre-wrap' }}>
                        {verdict.message}
                      </Typography>
                    )}
                  </Alert>
                </Box>
              )}

              {sampleResults.length > 0 && (
                <Box sx={{ mt: 2 }}>
                  <Typography variant="subtitle2" sx={{ mb: 1 }}>
                    Sample test results
                  </Typography>
                  <Stack spacing={1}>
                    {sampleResults.map((r, i) => (
                      <Paper
                        key={i}
                        variant="outlined"
                        sx={{
                          p: 1.5,
                          backgroundColor: r.result ? '#f1faf2' : '#fdecec',
                        }}
                      >
                        <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 0.5 }}>
                          <Typography variant="body2" fontWeight={600}>
                            Test {i + 1}
                          </Typography>
                          <Chip
                            label={r.result ? 'Pass' : 'Fail'}
                            size="small"
                            color={r.result ? 'success' : 'error'}
                          />
                          <Typography variant="caption" color="text.secondary">
                            {r.runtime}
                          </Typography>
                        </Stack>
                        <Box component="pre" sx={{ m: 0, fontSize: 12, whiteSpace: 'pre-wrap' }}>
                          input:    {r.input}
                          {'\n'}expected: {r.expected}
                          {'\n'}got:      {r.output}
                        </Box>
                      </Paper>
                    ))}
                  </Stack>
                </Box>
              )}
            </Paper>
          </Grid>
        </Grid>
      </Container>
    </>
  );
};

export default Problem;
