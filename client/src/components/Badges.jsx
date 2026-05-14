import React from 'react';
import { Chip } from '@mui/material';

const difficultyColors = {
  easy:   { bg: '#E7F6EC', fg: '#1E7F37' },
  medium: { bg: '#FFF3DD', fg: '#A0660A' },
  hard:   { bg: '#FCE8E8', fg: '#B0233A' },
};

export const DifficultyBadge = ({ difficulty }) => {
  const c = difficultyColors[difficulty] || { bg: '#eee', fg: '#444' };
  const label = difficulty
    ? difficulty[0].toUpperCase() + difficulty.slice(1)
    : 'Unknown';
  return (
    <Chip
      label={label}
      size="small"
      sx={{
        backgroundColor: c.bg,
        color: c.fg,
        fontWeight: 600,
        textTransform: 'capitalize',
      }}
    />
  );
};

const statusMeta = {
  accepted:            { label: 'Accepted',             bg: '#E7F6EC', fg: '#1E7F37' },
  wrong_answer:        { label: 'Wrong Answer',         bg: '#FCE8E8', fg: '#B0233A' },
  time_limit_exceeded: { label: 'Time Limit Exceeded',  bg: '#FFF3DD', fg: '#A0660A' },
  runtime_error:       { label: 'Runtime Error',        bg: '#FCE8E8', fg: '#B0233A' },
  compilation_error:   { label: 'Compilation Error',    bg: '#FFF3DD', fg: '#A0660A' },
  no_test_cases:       { label: 'No Test Cases',        bg: '#eee',    fg: '#555'    },
  server_error:        { label: 'Server Error',         bg: '#FCE8E8', fg: '#B0233A' },
};

export const StatusBadge = ({ status }) => {
  const m = statusMeta[status] || { label: status || 'Unknown', bg: '#eee', fg: '#444' };
  return (
    <Chip
      label={m.label}
      size="small"
      sx={{
        backgroundColor: m.bg,
        color: m.fg,
        fontWeight: 600,
      }}
    />
  );
};
