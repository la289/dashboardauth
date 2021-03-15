import logo from './logo.svg';
import './App.css';
import LogInForm from './components/LogInForm.js';
import Dashboard from './components/Dashboard';

function App() {
  return (
    <div className="App">
      <Dashboard />
    </div>
  );

  // if state.loggedIn()
  // return Dashboard
  // else
  // return LogInForm
}


export default App;
