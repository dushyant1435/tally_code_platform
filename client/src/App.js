import React from 'react';
import './App.css';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { ThemeProvider, createTheme, CssBaseline } from '@mui/material';

import { AuthProvider } from './auth/AuthContext';
import ProtectedRoute from './auth/ProtectedRoute';

import Home from './pages/Home';
import Login from './pages/Login';
import Signup from './pages/Signup';
import Problems from './pages/Problems';
import Problem from './pages/Problem';
import Playground from './pages/Playground';
import Submissions from './pages/Submissions';
import SubmissionDetail from './pages/SubmissionDetail';
import AdminHome from './pages/AdminHome';
import CreateProblem from './pages/CreateProblem';
import CreateTestCase from './pages/CreateTestCase';

const theme = createTheme({
  palette: {
    primary:  { main: '#1e1e2f' },
    secondary:{ main: '#58A399' },
    background: { default: '#f4f5f7' },
  },
  shape: { borderRadius: 8 },
  typography: { fontFamily: 'Inter, -apple-system, BlinkMacSystemFont, sans-serif' },
});

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <BrowserRouter>
        <AuthProvider>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/login" element={<Login />} />
            <Route path="/signup" element={<Signup />} />

            {/* Open routes */}
            <Route path="/problem" element={<Problems />} />
            <Route path="/problem/:id" element={<Problem />} />
            <Route path="/playground" element={<Playground />} />

            {/* User routes (auth required) */}
            <Route
              path="/submissions"
              element={
                <ProtectedRoute>
                  <Submissions />
                </ProtectedRoute>
              }
            />
            <Route
              path="/submissions/:id"
              element={
                <ProtectedRoute>
                  <SubmissionDetail />
                </ProtectedRoute>
              }
            />

            {/* Admin routes */}
            <Route
              path="/admin"
              element={
                <ProtectedRoute adminOnly>
                  <AdminHome />
                </ProtectedRoute>
              }
            />
            <Route
              path="/admin/problem/create"
              element={
                <ProtectedRoute adminOnly>
                  <CreateProblem />
                </ProtectedRoute>
              }
            />
            <Route
              path="/admin/problem/:id/testcase"
              element={
                <ProtectedRoute adminOnly>
                  <CreateTestCase />
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </BrowserRouter>
    </ThemeProvider>
  );
}

export default App;
