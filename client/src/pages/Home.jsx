import * as React from 'react';
import { useNavigate } from 'react-router-dom'; // Import useNavigate
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Typography from '@mui/material/Typography';
import NavBar from '../components/NavBar';
import SportsTennisIcon from '@mui/icons-material/SportsTennis';
import SportsKabaddiIcon from '@mui/icons-material/SportsKabaddi';
import ListAltOutlinedIcon from '@mui/icons-material/ListAltOutlined';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';

const Home = () => {
  const navigate = useNavigate(); // Initialize useNavigate

  const handleNavigation = (path) => {
    navigate(path);
  };

  return (
    <>
      <NavBar />
      <Box sx={{ flexGrow: 1, p: 3 }}>
        <Grid container spacing={2} justifyContent="center">
          {/* First row with two cards */}
          <Grid item xs={12} md={6} sx={{ display: 'flex', justifyContent: 'flex-end', pr: 0 }}>
            <Card 
              sx={{ 
                maxWidth: 400, 
                transition: 'transform 0.3s ease-in-out', 
                '&:hover': { transform: 'scale(1.05)' } ,
                // backgroundColor: '#FFD3B6'
                
              }} 
              onClick={() => handleNavigation('/playground')} // Navigate to '/playground'
            >
              <SportsTennisIcon 
                sx={{ 
                  height: 150,  
                  width: 1,   
                  ml: 0
                }} 
              />
              <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                  PLAYGROUND
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Have fun trying different things on our Code Editor.
                </Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} md={6} sx={{ display: 'flex', justifyContent: 'flex-start', pl: 0 }}> 
            <Card 
              sx={{ 
                maxWidth: 400, 
                transition: 'transform 0.3s ease-in-out', 
                '&:hover': { transform: 'scale(1.05)' } 
              }} 
              onClick={() => handleNavigation('/problem')} // Navigate to '/coding-arena'
            >
              <ListAltOutlinedIcon 
                sx={{ 
                  height: 150,  
                  width: 1,   
                  ml: 0
                }} 
              />
              <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                  CODING ARENA
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Solve coding problems and improve your skills.
                </Typography>
              </CardContent>
            </Card>
          </Grid>

          {/* Second row with one card centered */}
          <Grid item xs={12}>
            <Box display="flex" justifyContent="center">
              <Card 
                sx={{ 
                  maxWidth: 400, 
                  transition: 'transform 0.3s ease-in-out', 
                  '&:hover': { transform: 'scale(1.05)' } 
                }} 
                onClick={() => handleNavigation('/code-battle')} // Navigate to '/code-battle'
              >
                <SportsKabaddiIcon 
                  sx={{ 
                    height: 150,  
                    width: 1,   
                    ml: 0
                  }} 
                />
                <CardContent>
                  <Typography gutterBottom variant="h5" component="div">
                    CODE BATTLE
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Compete in contests and become the ultimate master.
                  </Typography>
                </CardContent>
              </Card>
            </Box>
          </Grid>
        </Grid>
      </Box>
    </>
  );
}

export default Home;
