import React from 'react';
import {
  Box,
  Card,
  CardActionArea,
  CardContent,
  Container,
  Grid,
  Typography,
} from '@mui/material';
import AddBoxIcon from '@mui/icons-material/AddBox';
import HistoryIcon from '@mui/icons-material/History';
import { Link } from 'react-router-dom';
import NavBar from '../components/NavBar';

const cards = [
  {
    to: '/admin/problem/create',
    icon: <AddBoxIcon sx={{ fontSize: 48, color: '#58A399' }} />,
    title: 'Create a problem',
    desc: 'Add a new challenge for users to solve.',
  },
  {
    to: '/submissions',
    icon: <HistoryIcon sx={{ fontSize: 48, color: '#58A399' }} />,
    title: 'All submissions',
    desc: "View every user's submission history.",
  },
];

const AdminHome = () => (
  <>
    <NavBar />
    <Container maxWidth="md" sx={{ mt: 4 }}>
      <Typography variant="h4" fontWeight={700} gutterBottom>
        Admin
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
        Manage problems and review activity across the platform.
      </Typography>
      <Grid container spacing={2}>
        {cards.map((c) => (
          <Grid item xs={12} sm={6} key={c.to}>
            <Card>
              <CardActionArea component={Link} to={c.to} sx={{ p: 2 }}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                  {c.icon}
                  <CardContent sx={{ p: '8px !important' }}>
                    <Typography variant="h6">{c.title}</Typography>
                    <Typography variant="body2" color="text.secondary">
                      {c.desc}
                    </Typography>
                  </CardContent>
                </Box>
              </CardActionArea>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Container>
  </>
);

export default AdminHome;
