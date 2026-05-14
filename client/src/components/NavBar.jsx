import React, { useState } from 'react';
import {
  AppBar,
  Avatar,
  Box,
  Button,
  Chip,
  Container,
  IconButton,
  Menu,
  MenuItem,
  Toolbar,
  Tooltip,
  Typography,
  Divider,
} from '@mui/material';
import CodeOffIcon from '@mui/icons-material/CodeOff';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext';

const navLinks = [
  { label: 'Problems', to: '/problem' },
  { label: 'Submissions', to: '/submissions', authOnly: true },
  { label: 'Playground', to: '/playground' },
];

const NavBar = () => {
  const { user, isAuthenticated, isAdmin, logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const [anchor, setAnchor] = useState(null);
  const openMenu = (e) => setAnchor(e.currentTarget);
  const closeMenu = () => setAnchor(null);

  const onLogout = () => {
    closeMenu();
    logout();
    navigate('/');
  };

  const isActive = (to) =>
    to === '/' ? location.pathname === '/' : location.pathname.startsWith(to);

  return (
    <AppBar position="sticky" elevation={1} sx={{ backgroundColor: '#1e1e2f' }}>
      <Container maxWidth="xl">
        <Toolbar disableGutters sx={{ gap: 1 }}>
          <CodeOffIcon sx={{ mr: 1 }} />
          <Typography
            variant="h6"
            component={Link}
            to="/"
            sx={{
              fontFamily: 'monospace',
              fontWeight: 700,
              letterSpacing: '.2rem',
              color: 'inherit',
              textDecoration: 'none',
              mr: 3,
            }}
          >
            HardCode
          </Typography>

          <Box sx={{ flexGrow: 1, display: 'flex', gap: 0.5 }}>
            {navLinks
              .filter((l) => !l.authOnly || isAuthenticated)
              .map((l) => (
                <Button
                  key={l.to}
                  component={Link}
                  to={l.to}
                  sx={{
                    color: 'white',
                    opacity: isActive(l.to) ? 1 : 0.7,
                    borderBottom: isActive(l.to) ? '2px solid #58A399' : '2px solid transparent',
                    borderRadius: 0,
                    textTransform: 'none',
                    fontWeight: 500,
                  }}
                >
                  {l.label}
                </Button>
              ))}
            {isAdmin && (
              <Button
                component={Link}
                to="/admin"
                sx={{
                  color: 'white',
                  opacity: isActive('/admin') ? 1 : 0.7,
                  borderBottom: isActive('/admin') ? '2px solid #58A399' : '2px solid transparent',
                  borderRadius: 0,
                  textTransform: 'none',
                  fontWeight: 500,
                }}
              >
                Admin
              </Button>
            )}
          </Box>

          {isAuthenticated ? (
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
              {isAdmin && (
                <Chip
                  label="ADMIN"
                  size="small"
                  sx={{ backgroundColor: '#E5C100', color: '#000', fontWeight: 700 }}
                />
              )}
              <Tooltip title={user.username}>
                <IconButton onClick={openMenu} sx={{ p: 0 }}>
                  <Avatar sx={{ width: 36, height: 36, backgroundColor: '#58A399' }}>
                    {user.username[0]?.toUpperCase()}
                  </Avatar>
                </IconButton>
              </Tooltip>
              <Menu
                anchorEl={anchor}
                open={!!anchor}
                onClose={closeMenu}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
                transformOrigin={{ vertical: 'top', horizontal: 'right' }}
                sx={{ mt: 1 }}
              >
                <Box sx={{ px: 2, py: 1 }}>
                  <Typography variant="body2" fontWeight={600}>
                    {user.username}
                  </Typography>
                  <Typography variant="caption" color="text.secondary">
                    {user.email}
                  </Typography>
                </Box>
                <Divider />
                <MenuItem
                  onClick={() => {
                    closeMenu();
                    navigate('/submissions');
                  }}
                >
                  My submissions
                </MenuItem>
                <MenuItem onClick={onLogout}>Log out</MenuItem>
              </Menu>
            </Box>
          ) : (
            <Box sx={{ display: 'flex', gap: 1 }}>
              <Button
                variant="text"
                sx={{ color: 'white', textTransform: 'none' }}
                component={Link}
                to="/login"
              >
                Sign in
              </Button>
              <Button
                variant="contained"
                sx={{ backgroundColor: '#58A399', textTransform: 'none' }}
                component={Link}
                to="/signup"
              >
                Sign up
              </Button>
            </Box>
          )}
        </Toolbar>
      </Container>
    </AppBar>
  );
};

export default NavBar;
