import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom';
import './SignUpPage.css'

function SignUpPage() {

    const [token, setToken] = useState('');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const navigate = useNavigate();
    
  const handleSignUp = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/v1/signup', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, password })
      });

      if (response.ok) {
        const data = await response.json();
        setToken(data.token);
        console.log("Token->", token)
        // navigate('/dashboard');
      } else {
        // Handle error response
        console.log("Error", response)
      }
    } catch (error) {
      // Handle fetch error
      console.log("Error occured", error)
    }
  };

  return (
    <div className="signup-container">
      <h2>Sign Up</h2>
      <div className="form-container">
        <input type="text" placeholder="Username" value={username} onChange={(e) => setUsername(e.target.value)} />
        <input type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} />
        <button onClick={handleSignUp}>Sign Up</button>
      </div>
    </div>
  );
}

export default SignUpPage;