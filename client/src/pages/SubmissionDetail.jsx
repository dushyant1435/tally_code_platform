import React, { useEffect, useState } from 'react';
import {
  Box,
  Button,
  CircularProgress,
  Container,
  Paper,
  Stack,
  Typography,
} from '@mui/material';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { Editor } from '@monaco-editor/react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import NavBar from '../components/NavBar';
import { StatusBadge } from '../components/Badges';
import { api } from '../api';

const SubmissionDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [sub, setSub] = useState(null);
  const [error, setError] = useState('');

  useEffect(() => {
    let cancelled = false;
    api
      .get(`/api/v1/submissions/${id}`)
      .then((data) => !cancelled && setSub(data))
      .catch((err) => !cancelled && setError(err.message));
  }, [id]);

  return (
    <>
      <NavBar />
      <Container maxWidth="md" sx={{ mt: 3 }}>
        <Button
          startIcon={<ArrowBackIcon />}
          onClick={() => navigate(-1)}
          sx={{ mb: 2 }}
        >
          Back
        </Button>

        {error ? (
          <Typography color="error">{error}</Typography>
        ) : !sub ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', mt: 6 }}>
            <CircularProgress />
          </Box>
        ) : (
          <Paper sx={{ p: 3 }} elevation={1}>
            <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 1 }}>
              <StatusBadge status={sub.status} />
              <Typography variant="h6">
                Submission #{sub.id}
              </Typography>
            </Stack>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              Problem{' '}
              <Link to={`/problem/${sub.problem_id}`}>
                {sub.problem_id}. {sub.problem_name}
              </Link>{' '}
              · by {sub.username} · {sub.language} ·{' '}
              {new Date(sub.created_at).toLocaleString()}
              {sub.runtime_seconds != null && ` · ${sub.runtime_seconds.toFixed(3)}s`}
              {sub.failed_test ? ` · failed test #${sub.failed_test}` : ''}
            </Typography>
            {sub.message && (
              <Paper variant="outlined" sx={{ p: 1.5, mb: 2, backgroundColor: '#fafafa' }}>
                <Typography variant="caption" color="text.secondary">
                  Details
                </Typography>
                <Box component="pre" sx={{ m: 0, fontSize: 13, whiteSpace: 'pre-wrap' }}>
                  {sub.message}
                </Box>
              </Paper>
            )}
            <Editor
              height="55vh"
              theme="vs-dark"
              language={sub.language || 'python'}
              value={sub.code || ''}
              options={{ readOnly: true, minimap: { enabled: false }, fontSize: 14 }}
            />
          </Paper>
        )}
      </Container>
    </>
  );
};

export default SubmissionDetail;
