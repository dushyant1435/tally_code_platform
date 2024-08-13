import React, { useEffect, useRef, useState } from 'react';
import {
  Box,
  Grid,
  Typography,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';
import { Editor } from '@monaco-editor/react';
import { CODE_SNIPPETS } from '../constants';
// import Output from './Output';

const CodeEditor = ({value, setValue}) => {
  
  const [language, setLanguage] = useState('python');
  const editorRef = useRef(null);

  const onSelect = (event) => {
    const selectedLanguage = event.target.value;
    setLanguage(selectedLanguage);
    setValue(CODE_SNIPPETS[selectedLanguage]);
  };

  const onMount = (editor) => {
    editorRef.current = editor;
  };

  useEffect(() => {
    console.log(value);
  }, [value]);

  return (
    <Box p={3}>
      <Grid container spacing={4}>
        <Grid item xs={12} md={12}>
          <FormControl fullWidth variant="outlined" sx={{ mb: 2 }}>
            <InputLabel id="language-selector-label">Language</InputLabel>
            <Select
              labelId="language-selector-label"
              value={language}
              onChange={onSelect}
              label="Language"
            >
              <MenuItem value="python">Python</MenuItem>
              <MenuItem value="javascript">JavaScript</MenuItem>
              <MenuItem value="cpp">C++</MenuItem>
              <MenuItem value="java">Java</MenuItem>
            </Select>
          </FormControl>
          <Editor
            height="50vh"
            theme="vs-dark"
            language={language}
            value={value}
            onMount={onMount}
            onChange={(value) => setValue(value)}
            options={{
              minimap: { enabled: false },
            }}
          />
        </Grid>
        {/* <Grid item xs={12} md={6}>
          <Output editorRef={editorRef} language={language} />
        </Grid> */}
      </Grid>
    </Box>
  );
};

export default CodeEditor;
