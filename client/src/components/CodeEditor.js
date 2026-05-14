import React, { useRef, useState } from 'react';
import {
  Box,
  Grid,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';
import { Editor } from '@monaco-editor/react';
import { CODE_SNIPPETS } from '../constants';

const CodeEditor = ({ value, setValue }) => {
  const [language, setLanguage] = useState('python');
  const editorRef = useRef(null);

  const onSelect = (event) => {
    const next = event.target.value;
    setLanguage(next);

    // Only seed the editor with a snippet if the user hasn't typed anything
    // custom yet. Otherwise switching languages would silently wipe their code.
    const isUnchanged = Object.values(CODE_SNIPPETS).some(
      (snip) => snip === value,
    );
    if (!value || isUnchanged) {
      setValue(CODE_SNIPPETS[next] || '');
    }
  };

  const onMount = (editor) => {
    editorRef.current = editor;
  };

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
            onChange={(next) => setValue(next ?? '')}
            options={{
              minimap: { enabled: false },
            }}
          />
        </Grid>
      </Grid>
    </Box>
  );
};

export default CodeEditor;
