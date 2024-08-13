import React, { useEffect, useState } from 'react';
import CodeEditor from '../components/CodeEditor';
import NavBar from '../components/NavBar';
import TextBox from '../components/Textbox';
import { Button, Avatar, Select, MenuItem, InputLabel, FormControl, TextField, Container, Typography, Paper, Grid } from '@mui/material';
import { CODE_SNIPPETS } from '../constants';


const Playground = (() => {
    const [value, setValue] = useState(CODE_SNIPPETS['python']);
    const [inputValue, setInputValue] = useState();
    const [outputValue, setOutputValue] = useState();


    const getStatus = async () => {
        const response = await fetch(`http://localhost:8080/api/v1/runCustomCode`, {
            method: "POST",
            headers: {
              'Content-Type': 'application/json', // Set the content type to JSON
            },
          body: JSON.stringify({
              "code": value,
              "input": inputValue,
          }),
        });
        const data = await response.json();
        
        // setSuccess(data.success);
        // console.log(data)
        // console.log(data.success)
        setOutputValue(data.output)
    };
  
      const handleSubmit = () => {
          // Replace this with your actual logic for what should happen when you click "Submit"
          console.log("Submit button clicked");
          getStatus()
      };


  return (
   <>
    <NavBar></NavBar>
    <div>
    <CodeEditor value={value} setValue={setValue}/>
    </div>

    <Button 
        variant="contained" 
        sx={{ backgroundColor: 'green' , width:'1565px'}} 
        onClick={handleSubmit}
    >
        RUN CODE
    </Button>
    <br />
    <br />

    <Grid container spacing={2}>
    <Grid item xs={6}>

    <Button 
        variant="contained" 
        sx={{ backgroundColor: 'blue' , width:'750px'}} 
    >
        INPUT
    </Button>

    <TextBox value={inputValue} setValue={setInputValue}/>

    </Grid>
    <Grid item xs={6}>

    <Button 
        variant="contained" 
        sx={{ backgroundColor: 'blue' , width:'750px'}} 
    >
        OUTPUT
    </Button>
    <TextBox value={outputValue} setValue={setOutputValue}/>
    </Grid>
    </Grid>
    
   
   </> 
  )
})

export default Playground