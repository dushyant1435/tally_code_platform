import React, { useEffect, useState } from 'react';
import {
  Button,
  Container,
  Typography,
  Paper,
  Grid,
} from '@mui/material';
import CodeEditor from '../components/CodeEditor';
import NavBar from '../components/NavBar';
import { useNavigate, useParams } from 'react-router-dom';
import { CODE_SNIPPETS } from '../constants';
import { API_BASE, CURRENT_USER_ID } from '../config';

const Problem = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const num = parseInt(id, 10);

  const [value, setValue] = useState(CODE_SNIPPETS['python']);
  const [problem, setProblem] = useState(null);
  const [testCases, setTestCases] = useState([]);
  const [output, setOutput] = useState([]);
  const [submissionMsg, setSubmissionMsg] = useState('');

  const getSample = async () => {
    try {
      const response = await fetch(
        `${API_BASE}/api/v1/problem/${id}/sampleTestCases`,
      );
      if (!response.ok) return;
      const data = await response.json();
      setTestCases(Array.isArray(data) ? data : []);
    } catch (err) {
      console.error('getSample', err);
    }
  };

  const getProblem = async () => {
    try {
      const response = await fetch(`${API_BASE}/api/v1/problem/${id}`);
      if (!response.ok) {
        setProblem(null);
        return;
      }
      const data = await response.json();
      setProblem(data);
    } catch (err) {
      console.error('getProblem', err);
    }
  };

  const handleSubmit = async () => {
    try {
      const response = await fetch(`${API_BASE}/api/v1/runCode`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          id: num,
          code: value,
          user_id: CURRENT_USER_ID,
        }),
      });
      const data = await response.json();
      if (data.success) {
        setSubmissionMsg('All test cases passed! Problem solved.');
      } else {
        setSubmissionMsg(
          data.message
            ? `Failed${data.failedAt ? ` at test ${data.failedAt}` : ''}: ${data.message}`
            : 'Submission failed.',
        );
      }
    } catch (err) {
      console.error('handleSubmit', err);
      setSubmissionMsg('Error submitting code.');
    }
  };

  const handleRunSample = async () => {
    try {
      const response = await fetch(`${API_BASE}/api/v1/runSampleCode`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          id: num,
          code: value,
          user_id: CURRENT_USER_ID,
        }),
      });
      const data = await response.json();
      setOutput(Array.isArray(data.results) ? data.results : []);
    } catch (err) {
      console.error('handleRunSample', err);
    }
  };

  useEffect(() => {
    getProblem();
    getSample();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);

  return (
    <>
      <NavBar />
      <br />
      <Grid container spacing={0}>
        <Grid item xs={3}>
          <Button
            variant="contained"
            sx={{ backgroundColor: 'blue', width: '290px' }}
            onClick={() => navigate(`/problem/${id}/testcase`)}
          >
            ADD TEST CASES
          </Button>
        </Grid>
      </Grid>

      <br />
      <Container maxWidth="lg">
        <div className="flex flex-col min-h-screen">
          <main className="flex-1 py-8 md:py-12 grid md:grid-cols-2 gap-8 md:gap-12">
            <Paper elevation={3} sx={{ padding: 3, backgroundColor: '#fff', color: '#333' }}>
              {problem ? (
                <>
                  <Typography variant="h4" gutterBottom>
                    {problem.name}
                  </Typography>
                  <Typography variant="h6">Problem Description</Typography>
                  <Typography paragraph>{problem.description}</Typography>
                  <Typography variant="h6">Constraints</Typography>
                  <Typography paragraph>{problem.constraints}</Typography>
                  <Typography variant="h6">Input</Typography>
                  <Typography paragraph>{problem.input_format}</Typography>
                  <Typography variant="h6">Output</Typography>
                  <Typography paragraph>{problem.output_format}</Typography>
                </>
              ) : (
                <Typography variant="h6" gutterBottom>
                  Loading...
                </Typography>
              )}
            </Paper>

            <Paper elevation={3} sx={{ padding: 3, backgroundColor: '#fff', color: '#333' }}>
              <CodeEditor value={value} setValue={setValue} />
              <Grid container spacing={2}>
                <Grid item xs={6}>
                  <Button
                    variant="contained"
                    sx={{ backgroundColor: 'primary', '&:hover': { backgroundColor: 'darkgreen' } }}
                    fullWidth
                    onClick={handleRunSample}
                  >
                    Run Sample
                  </Button>
                </Grid>
                <Grid item xs={6}>
                  <Button
                    variant="contained"
                    sx={{ backgroundColor: 'green', '&:hover': { backgroundColor: 'darkgreen' } }}
                    fullWidth
                    onClick={handleSubmit}
                  >
                    Submit
                  </Button>
                </Grid>
              </Grid>
              {submissionMsg && (
                <Typography
                  variant="body1"
                  sx={{ mt: 2, color: submissionMsg.startsWith('All') ? 'green' : 'red' }}
                >
                  {submissionMsg}
                </Typography>
              )}
            </Paper>

            <Paper elevation={3} sx={{ padding: 3, backgroundColor: '#fff', color: '#333' }}>
              <Typography variant="h6">Test Cases</Typography>
              {testCases.length > 0 ? (
                testCases.map((testCase, index) => (
                  <div key={index}>
                    <Typography variant="subtitle1">Test Case {index + 1}</Typography>
                    <Typography variant="body2">Input: {testCase.input}</Typography>
                    <Typography variant="body2">
                      Expected Output: {testCase.output}
                    </Typography>
                    <br />
                  </div>
                ))
              ) : (
                <Typography variant="body2">No test cases available</Typography>
              )}

              <Typography variant="h6">YOUR OUTPUT</Typography>
              {output.length > 0 ? (
                output.map((out, index) => (
                  <div key={index}>
                    <Typography variant="subtitle1">Test Case {index + 1}</Typography>
                    <Typography
                      variant="body2"
                      color={out.result ? 'green' : 'red'}
                    >
                      Your Output: {out.output}
                    </Typography>
                    <Typography variant="body2" color="grey">
                      runtime: {out.runtime}
                    </Typography>
                    <Typography variant="body2" color="grey">
                      memory: {out.memory_used}
                    </Typography>
                    <br />
                  </div>
                ))
              ) : (
                <Typography variant="body2">PLEASE RUN THE CODE</Typography>
              )}
            </Paper>
          </main>
          <footer className="bg-gray-900 text-white px-4 md:px-6 py-3 flex items-center justify-between">
            <Typography variant="body2">&copy; 2024 TALLY. All rights reserved.</Typography>
          </footer>
        </div>
      </Container>
    </>
  );
};

export default Problem;
