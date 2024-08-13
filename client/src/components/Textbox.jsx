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

const TextBox = ({value, setValue}) => {
  
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
          <Editor
            height="20vh"
            // width="vh"
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

export default TextBox;
