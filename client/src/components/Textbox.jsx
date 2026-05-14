import React, { useRef } from 'react';
import { Box, Grid } from '@mui/material';
import { Editor } from '@monaco-editor/react';

const TextBox = ({ value, setValue, language = 'plaintext' }) => {
  const editorRef = useRef(null);

  const onMount = (editor) => {
    editorRef.current = editor;
  };

  return (
    <Box p={3}>
      <Grid container spacing={4}>
        <Grid item xs={12} md={12}>
          <Editor
            height="20vh"
            theme="vs-dark"
            language={language}
            value={value ?? ''}
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

export default TextBox;
