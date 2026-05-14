import React, { useRef, useState } from 'react';
import {
  Box,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';
import { Editor } from '@monaco-editor/react';
import { CODE_SNIPPETS } from '../constants';

const CodeEditor = ({
  value,
  setValue,
  language: controlledLanguage,
  setLanguage: setControlledLanguage,
}) => {
  const editorRef = useRef(null);
  const [uncontrolledLanguage, setUncontrolledLanguage] = useState('python');

  // Component can be used standalone (Playground) or with language lifted up
  // (Problem detail).
  const language = controlledLanguage ?? uncontrolledLanguage;
  const setLanguage = setControlledLanguage ?? setUncontrolledLanguage;

  const onSelect = (event) => {
    const next = event.target.value;
    setLanguage(next);
    const isUnchanged = Object.values(CODE_SNIPPETS).some((snip) => snip === value);
    if (!value || isUnchanged) {
      setValue(CODE_SNIPPETS[next] || '');
    }
  };

  const onMount = (editor) => {
    editorRef.current = editor;
  };

  return (
    <Box>
      <FormControl fullWidth variant="outlined" size="small" sx={{ mb: 1 }}>
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
        height="55vh"
        theme="vs-dark"
        language={language}
        value={value}
        onMount={onMount}
        onChange={(next) => setValue(next ?? '')}
        options={{ minimap: { enabled: false }, fontSize: 14 }}
      />
    </Box>
  );
};

export default CodeEditor;
