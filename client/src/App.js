import './App.css';
import { BrowserRouter, Navigate, Routes, Route } from 'react-router-dom';
import Problems from './pages/Problems';
import Problem from './pages/Problem';
import Home from './pages/Home';
import Playground from './pages/Playground';
import CreateProblem from './pages/CreateProblem';
import CreateTestCase from './pages/CreateTestCase.jsx';
function App() {
  return (
    <div className="App">
      <BrowserRouter>
        {/* <ThemeProvider theme={theme}> */}
        {/* <CssBaseline /> */}
        <Routes>
        <Route path="/" element={<Home />} />
          <Route path="/problem" element={<Problems />} />
          <Route path="/problem/:id" element={<Problem />} />
          <Route path="/playground" element={<Playground />} />
          <Route path="/problem/create" element={<CreateProblem />} />
          <Route path="/problem/:id/testcase" element={<CreateTestCase />} />
          {/* <Route
              path="/home"
              element={isAuth ? <HomePage /> : <Navigate to="/" />}
            /> */}
          {/* <Route
              path="/profile/:userId"
              element={isAuth ? <ProfilePage /> : <Navigate to="/" />}
            /> */}
        </Routes>
        {/* </ThemeProvider> */}
      </BrowserRouter>
    </div>
  );
}

export default App;
