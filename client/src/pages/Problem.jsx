import React, { useEffect, useState } from 'react';
import { Button, Avatar, Select, MenuItem, InputLabel, FormControl, TextField, Container, Typography, Paper, Grid } from '@mui/material';
import CodeEditor from '../components/CodeEditor';
import NavBar from '../components/NavBar';
import { green, red } from '@mui/material/colors';
import { useParams } from 'react-router-dom';

import { CODE_SNIPPETS } from '../constants';

import { useNavigate } from 'react-router-dom';

const Problem = () => {

  const [value, setValue] = useState(CODE_SNIPPETS['python']);
    const { id } = useParams(); // Extract the problem ID from the URL
    let num = parseInt(id);
    const [problemId, setProblemId] = useState(id)
    const [problem, setProblem] = useState(null)
    const [testCases, setTestCases] = useState([]); 

    const getSample = async () => {
        const response = await fetch(`http://localhost:8080/api/v1/problem/${problemId}/sampleTestCases`, {
            method: "GET",
        });
        const data = await response.json();
        setTestCases(data);
    };

    const getProblem = async () => {
        const response = await fetch(`http://localhost:8080/api/v1/problem/${problemId}`, {
            method: "GET",
        });
        const data = await response.json();
        setProblem(data);
    };

    const [success,setSuccess] = useState()

    const getStatus = async () => {
      const response = await fetch(`http://localhost:8080/api/v1/runCode`, {
          method: "POST",
          headers: {
            'Content-Type': 'application/json', // Set the content type to JSON
          },
        body: JSON.stringify({
            "id": num,
            "code": value,
            "user_id": 123,
        }),
      });
      const data = await response.json();
      
      setSuccess(data.success);
      if(data.success==true){
        alert("YAYYY")
      }
      // console.log(data.success)
  };

    const handleSubmit = () => {
        // Replace this with your actual logic for what should happen when you click "Submit"
        console.log("Submit button clicked");
        getStatus()
    };


    const [output,setOutput] = useState([])
    const getOutPut = async () => {
      const response = await fetch(`http://localhost:8080/api/v1/runSampleCode`, {
          method: "POST",
          headers: {
            'Content-Type': 'application/json', // Set the content type to JSON
          },
        body: JSON.stringify({
            "id": num,
            "code": value,
            "user_id": 123,
        }),
      });
      const data = await response.json();
      
      setOutput(data.results);

      // console.log(data.results[0].output)
  };

    const handleRunSample = () => {
        // Replace this with your actual logic for what should happen when you click "Submit"
        console.log("Run Sample button clicked");
        getOutPut()
    };

    useEffect(() => {
        getProblem();
        getSample();
    }, [id]);

  //   useEffect(() => {
  //     console.log(output[0].output)
  // }, [output]);
  const navigate = useNavigate();
  const handleNav = (route) => {
    navigate(route);
  };
  return(
    <>
        <NavBar />
        <br />
      <Grid container spacing={0}>
      <Grid item xs={3}>
      <Button 
        variant="contained" 
        sx={{ backgroundColor: 'blue' , width:'290px'}} 
        onClick={() => handleNav(`/problem/${id}/testcase`)}
    >
        ADD TEST CASES
    </Button>
      </Grid>
      </Grid>

      <br />
        <Container maxWidth="lg">
            <div className="flex flex-col min-h-screen">
                <main className="flex-1 py-8 md:py-12 grid md:grid-cols-2 gap-8 md:gap-12">
                    <br />
                    <Paper elevation={3} sx={{ padding: 3, backgroundColor: '#fff', color: '#333' }}>
                        {problem ? (
                            <>
                                <Typography variant="h4" gutterBottom>
                                    {problem.name} {/* This will print the problem name */}
                                </Typography>
                                <Typography variant="h6">Problem Description</Typography>
                                <Typography paragraph>{problem.description}</Typography>
                                <Typography variant="h6">Constraints</Typography>
                                <Typography paragraph>{problem.constraints}</Typography>
                                <Typography variant="h6">Input</Typography>
                                <Typography paragraph>{problem.input_format}</Typography>
                                <Typography variant="h6">Output</Typography>
                                <Typography paragraph>{problem.output_format}</Typography>
                            </>
                        ) : (
                            <Typography variant="h6" gutterBottom>
                                Loading...
                            </Typography>
                        )}
                    </Paper>
                    <br />
                    <Paper elevation={3} sx={{ padding: 3, backgroundColor: '#fff', color: '#333' }}>
                        <form>
                            <CodeEditor value={value} setValue={setValue}/>
                            <Grid container spacing={2}>
                                <Grid item xs={6}>
                                    {/* <Button variant="contained" color="primary" type="SUBMIT" fullWidth onClick={handleRunSample}> */}
                                    <Button 
                                        variant="contained" 
                                        sx={{ backgroundColor: 'primary', '&:hover': { backgroundColor: 'darkgreen' } }} 
                                        fullWidth
                                        onClick={handleRunSample} // Add onClick handler here
                                    >
                                        Run Sample
                                    </Button>
                                </Grid>
                                <Grid item xs={6}>
                                    <Button 
                                        variant="contained" 
                                        sx={{ backgroundColor: 'green', '&:hover': { backgroundColor: 'darkgreen' } }} 
                                        fullWidth
                                        onClick={handleSubmit} // Add onClick handler here
                                    >
                                        Submit
                                    </Button>
                                </Grid>
                            </Grid>
                        </form>
                    </Paper>
                    <br />
                    <Paper elevation={3} sx={{ padding: 3, backgroundColor: '#fff', color: '#333' }}>
                        <Typography variant="h6">Test Cases</Typography>
                        {testCases !=null ? (
                            testCases.map((testCase, index) => (
                                <div key={index}>
                                    <Typography variant="subtitle1">Test Case {index + 1}</Typography>
                                    <Typography variant="body2">Input: {testCase.input}</Typography>
                                    <Typography variant="body2">Expected Output: {testCase.output}</Typography>
                                    {/* <Typography variant="body2">Your Output: {output[index].output}</Typography> */}
                                    <br />
                                </div>
                            ))
                        ) : (
                            <Typography variant="body2">No test cases available</Typography>
                        )}

                        <Typography variant="h6">YOUR OUTPUT</Typography>
                        {output!= null ? (
                            output.map((out, index) => (
                                <div key={index}>
                                  <Typography variant="subtitle1">Test Case {index + 1}</Typography>
                                  <Typography
                                  variant="body2"
                                  color={out.result == false ? 'red' : 'green'}
                                >
                                      Your Output: {out.output}
                                      </Typography>
                                      <Typography
                                  variant="body2"
                                  color='grey'
                                >
                                      {out.runtime}
                                      </Typography>
                                      <Typography
                                  variant="body2"
                                  color='grey'
                                >
                                      {out.memory_used}
                                      </Typography>
                                    <br />
                                </div>
                            ))
                        ) : (
                            <Typography variant="body2">PLEASE RUN THE CODE</Typography>
                        )}
                    </Paper>
                </main>
                <footer className="bg-gray-900 text-white px-4 md:px-6 py-3 flex items-center justify-between">
                    <Typography variant="body2">&copy; 2024 TALLY. All rights reserved.</Typography>
                </footer>
            </div>
        </Container>
    </>
  );
}

export default Problem;
