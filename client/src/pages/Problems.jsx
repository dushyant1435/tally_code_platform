import React, { useEffect, useState } from 'react';
import NavBar from '../components/NavBar';
import { Button, Paper, Grid } from '@mui/material';
import { styled } from '@mui/material/styles';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell, { tableCellClasses } from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import { useNavigate } from 'react-router-dom';
import { API_BASE, CURRENT_USER_ID } from '../config';

const StyledTableCell = styled(TableCell)(({ theme }) => ({
  [`&.${tableCellClasses.head}`]: {
    backgroundColor: '#58A399',
    color: theme.palette.common.white,
  },
  [`&.${tableCellClasses.body}`]: {
    fontSize: 14,
  },
}));

const StyledTableRow = styled(TableRow)(() => ({
  '&:nth-of-type(odd)': {
    backgroundColor: '#E5DDC5',
    border: 0,
  },
  '&:last-child td, &:last-child th': {
    border: 0,
  },
  '&:nth-of-type(even)': {
    backgroundColor: '#F1EEDC',
    border: 0,
  },
  '&:hover': {
    transform: 'scale(1.01)',
    transition: 'transform 0.3s ease-in-out',
    cursor: 'pointer',
    zIndex: 1,
    position: 'relative',
  },
}));

const Problems = () => {
  const [problems, setProblems] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    const getProblems = async () => {
      try {
        const response = await fetch(
          `${API_BASE}/api/v1/problems?user_id=${CURRENT_USER_ID}`,
        );
        if (!response.ok) {
          console.error('Failed to fetch problems', response.status);
          return;
        }
        const data = await response.json();
        setProblems(Array.isArray(data) ? data : []);
      } catch (err) {
        console.error('Failed to fetch problems', err);
      }
    };
    getProblems();
  }, []);

  const handleRowClick = (route) => {
    navigate(route);
  };

  return (
    <>
      <NavBar />
      <br />
      <Grid container spacing={0}>
        <Grid item xs={3}>
          <Button
            variant="contained"
            sx={{ backgroundColor: 'blue', width: '290px' }}
            onClick={() => handleRowClick(`/problem/create`)}
          >
            CREATE CHALLENGE
          </Button>
        </Grid>
      </Grid>

      <br />
      <TableContainer
        component={Paper}
        sx={{
          width: '95%',
          margin: 'auto',
          borderRadius: 0,
          border: '0px solid #ddd',
        }}
      >
        <Table sx={{ minWidth: 700 }} aria-label="customized table">
          <TableHead>
            <TableRow>
              <StyledTableCell>PId</StyledTableCell>
              <StyledTableCell align="center">Title</StyledTableCell>
              <StyledTableCell align="center">Status</StyledTableCell>
              <StyledTableCell align="center">Difficulty</StyledTableCell>
              <StyledTableCell align="center">Created By</StyledTableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {problems.map((row) => (
              <StyledTableRow
                key={row.id}
                onClick={() => handleRowClick(`/problem/${row.id}`)}
              >
                <StyledTableCell component="th" scope="row">
                  {row.id}
                </StyledTableCell>
                <StyledTableCell align="center">{row.name}</StyledTableCell>
                <StyledTableCell align="center">
                  {row.status ? 'Solved' : 'Unsolved'}
                </StyledTableCell>
                <StyledTableCell align="center">hard</StyledTableCell>
                <StyledTableCell align="center">{row.user_id}</StyledTableCell>
              </StyledTableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </>
  );
};

export default Problems;
