import React, { useState } from 'react';
import {
  Box,
  Button,
  CircularProgress,
  Container,
  Grid,
  Paper,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import NavBar from '../components/NavBar';
import CodeEditor from '../components/CodeEditor';
import { CODE_SNIPPETS } from '../constants';
import { api } from '../api';

const Playground = () => {
  const [code, setCode] = useState(CODE_SNIPPETS.python);
  const [language, setLanguage] = useState('python');
  const [input, setInput] = useState('');
  const [output, setOutput] = useState('');
  const [busy, setBusy] = useState(false);

  const run = async () => {
    setBusy(true);
    try {
      const data = await api.post('/api/v1/runCustomCode', { code, input, language });
      setOutput(data.output ?? '');
    } catch (err) {
      setOutput(`Error: ${err.message}`);
    } finally {
      setBusy(false);
    }
  };

  return (
    <>
      <NavBar />
      <Container maxWidth="xl" sx={{ mt: 3, mb: 6 }}>
        <Typography variant="h4" fontWeight={700} sx={{ mb: 2 }}>
          Playground
        </Typography>
        <Grid container spacing={2}>
          <Grid item xs={12} md={7}>
            <Paper sx={{ p: 2 }} elevation={1}>
              <CodeEditor
                value={code}
                setValue={setCode}
                language={language}
                setLanguage={setLanguage}
              />
              <Button
                fullWidth
                variant="contained"
                color="success"
                onClick={run}
                disabled={busy}
                startIcon={busy ? <CircularProgress size={18} color="inherit" /> : <PlayArrowIcon />}
                sx={{ mt: 1 }}
              >
                Run
              </Button>
            </Paper>
          </Grid>
          <Grid item xs={12} md={5}>
            <Stack spacing={2}>
              <Paper sx={{ p: 2 }} elevation={1}>
                <Typography variant="subtitle2" sx={{ mb: 1 }}>
                  stdin
                </Typography>
                <TextField
                  multiline
                  minRows={6}
                  fullWidth
                  value={input}
                  onChange={(e) => setInput(e.target.value)}
                  placeholder="Type any input your program reads from stdin…"
                />
              </Paper>
              <Paper sx={{ p: 2, backgroundColor: '#1e1e1e', color: '#e0e0e0' }} elevation={1}>
                <Typography variant="subtitle2" sx={{ mb: 1, color: '#e0e0e0' }}>
                  output
                </Typography>
                <Box
                  component="pre"
                  sx={{
                    m: 0,
                    fontFamily: 'monospace',
                    fontSize: 13,
                    minHeight: 140,
                    whiteSpace: 'pre-wrap',
                  }}
                >
                  {output || '\u00a0'}
                </Box>
              </Paper>
            </Stack>
          </Grid>
        </Grid>
      </Container>
    </>
  );
};

export default Playground;
