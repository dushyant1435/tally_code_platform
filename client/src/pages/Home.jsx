import React from 'react';
import {
  Box,
  Button,
  Card,
  CardActionArea,
  CardContent,
  Container,
  Grid,
  Stack,
  Typography,
} from '@mui/material';
import SportsTennisIcon from '@mui/icons-material/SportsTennis';
import ListAltOutlinedIcon from '@mui/icons-material/ListAltOutlined';
import SportsKabaddiIcon from '@mui/icons-material/SportsKabaddi';
import { Link, useNavigate } from 'react-router-dom';
import NavBar from '../components/NavBar';
import { useAuth } from '../auth/AuthContext';

const cards = [
  {
    to: '/playground',
    icon: <SportsTennisIcon sx={{ fontSize: 96, color: '#58A399' }} />,
    title: 'Playground',
    desc: 'Open editor with your own input. Run Python freely.',
  },
  {
    to: '/problem',
    icon: <ListAltOutlinedIcon sx={{ fontSize: 96, color: '#58A399' }} />,
    title: 'Coding arena',
    desc: 'Solve curated problems with real verdicts.',
  },
  {
    to: '/code-battle',
    icon: <SportsKabaddiIcon sx={{ fontSize: 96, color: '#888' }} />,
    title: 'Code battle',
    desc: 'Compete head-to-head — coming soon.',
    disabled: true,
  },
];

const Home = () => {
  const { user, isAuthenticated } = useAuth();
  const navigate = useNavigate();

  return (
    <>
      <NavBar />
      <Box
        sx={{
          background: 'linear-gradient(135deg, #1e1e2f 0%, #2c3e50 100%)',
          color: '#fff',
          py: 8,
        }}
      >
        <Container maxWidth="md">
          <Typography variant="h2" fontWeight={800} sx={{ letterSpacing: '.05em' }}>
            TALLY
          </Typography>
          <Typography variant="h6" sx={{ opacity: 0.8, mb: 3 }}>
            A small online judge for competitive coding practice.
          </Typography>
          {isAuthenticated ? (
            <Stack direction="row" spacing={2}>
              <Button variant="contained" size="large" component={Link} to="/problem" sx={{ backgroundColor: '#58A399' }}>
                Browse problems
              </Button>
              <Button variant="outlined" size="large" component={Link} to="/playground" sx={{ color: '#fff', borderColor: 'rgba(255,255,255,0.5)' }}>
                Open playground
              </Button>
            </Stack>
          ) : (
            <Stack direction="row" spacing={2}>
              <Button variant="contained" size="large" component={Link} to="/signup" sx={{ backgroundColor: '#58A399' }}>
                Get started
              </Button>
              <Button variant="outlined" size="large" component={Link} to="/login" sx={{ color: '#fff', borderColor: 'rgba(255,255,255,0.5)' }}>
                Sign in
              </Button>
            </Stack>
          )}
        </Container>
      </Box>

      <Container maxWidth="lg" sx={{ mt: 6, mb: 8 }}>
        {isAuthenticated && (
          <Typography variant="body1" sx={{ mb: 3 }}>
            Welcome back, <strong>{user.username}</strong>.
          </Typography>
        )}
        <Grid container spacing={3}>
          {cards.map((c) => (
            <Grid item xs={12} md={4} key={c.to}>
              <Card
                sx={{
                  height: '100%',
                  opacity: c.disabled ? 0.55 : 1,
                  transition: 'transform 0.2s ease',
                  '&:hover': { transform: c.disabled ? 'none' : 'translateY(-4px)' },
                }}
              >
                <CardActionArea
                  disabled={c.disabled}
                  onClick={() => !c.disabled && navigate(c.to)}
                  sx={{ p: 3, height: '100%' }}
                >
                  <Box sx={{ textAlign: 'center' }}>{c.icon}</Box>
                  <CardContent>
                    <Typography variant="h5" align="center" gutterBottom>
                      {c.title}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" align="center">
                      {c.desc}
                    </Typography>
                  </CardContent>
                </CardActionArea>
              </Card>
            </Grid>
          ))}
        </Grid>
      </Container>
    </>
  );
};

export default Home;
