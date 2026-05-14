import React, { useState } from 'react';
import CodeEditor from '../components/CodeEditor';
import NavBar from '../components/NavBar';
import TextBox from '../components/Textbox';
import { Button, Grid } from '@mui/material';
import { CODE_SNIPPETS } from '../constants';
import { API_BASE } from '../config';

const Playground = () => {
  const [value, setValue] = useState(CODE_SNIPPETS['python']);
  const [inputValue, setInputValue] = useState('');
  const [outputValue, setOutputValue] = useState('');

  const handleSubmit = async () => {
    try {
      const response = await fetch(`${API_BASE}/api/v1/runCustomCode`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          code: value,
          input: inputValue ?? '',
        }),
      });
      const data = await response.json();
      setOutputValue(data.output ?? '');
    } catch (err) {
      console.error('Playground run', err);
      setOutputValue('Error: failed to reach server');
    }
  };

  return (
    <>
      <NavBar />
      <div>
        <CodeEditor value={value} setValue={setValue} />
      </div>

      <Grid container spacing={2} sx={{ px: 3 }}>
        <Grid item xs={12}>
          <Button
            variant="contained"
            sx={{ backgroundColor: 'green', width: '100%' }}
            onClick={handleSubmit}
          >
            RUN CODE
          </Button>
        </Grid>
      </Grid>

      <Grid container spacing={2} sx={{ px: 3, mt: 2 }}>
        <Grid item xs={12} md={6}>
          <Button variant="contained" sx={{ backgroundColor: 'blue', width: '100%' }}>
            INPUT
          </Button>
          <TextBox value={inputValue} setValue={setInputValue} />
        </Grid>
        <Grid item xs={12} md={6}>
          <Button variant="contained" sx={{ backgroundColor: 'blue', width: '100%' }}>
            OUTPUT
          </Button>
          <TextBox value={outputValue} setValue={setOutputValue} />
        </Grid>
      </Grid>
    </>
  );
};

export default Playground;
