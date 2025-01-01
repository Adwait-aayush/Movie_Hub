import React, { useState } from 'react';
import './RegistrationForm.css';
import { useNavigate, useOutletContext } from 'react-router-dom';
import swal from 'sweetalert';
const Registration = () => {
  const navigate = useNavigate();
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
 
  const handleSubmit = async (e) => {
    e.preventDefault();

    const data = {
      username: username,
      email: email,
      password: password,
    };

    const headers = new Headers();
    headers.append('Content-Type', 'application/json');

    const reqOptions = {
      method: 'POST',
      headers: headers,
      credentials:"include",
      body: JSON.stringify(data),
    };

    try {
      const response = await fetch('http://localhost:4000/Register', reqOptions);
      const result = await response.json();
      if(result.status===200){
        
        swal("Registered!", "You have been Signed in!", "success");
     navigate("/")
      }
      else{
        swal("Error!", `${result.message}`, "error");
      }
  
    } catch (error) {
      console.error('Error during registration:', error);
      alert('An error occurred during registration.');
    }
  };

  return (
    <div className="registration-container">
      <form onSubmit={handleSubmit} className="registration-form">
        <h2>Register</h2>
        <div className="form-group">
          <label htmlFor="username">Username</label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        </div>
        <div className="form-group">
          <label htmlFor="email">Email</label>
          <input
            type="email"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div className="form-group">
          <label htmlFor="password">Password</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <button type="submit" className="submit-btn">Register</button>
      </form>
    </div>
  );
};

export default Registration;
