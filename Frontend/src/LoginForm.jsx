import React, { useState } from 'react';
import './LoginForm.css';
import swal from 'sweetalert';
import { useNavigate, useOutletContext } from 'react-router-dom';
const LoginForm = () => {

 
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
const navigate=useNavigate()
  const handleSubmit = async(e) => {
    e.preventDefault();
    const data ={
      email: email,
      password: password
    }
    const headers=new Headers();
    headers.append('Content-Type', 'application/json');
    const reqOptions={
      method: 'POST',
      headers: headers,
      credentials:"include",
      body: JSON.stringify(data)
    }
    const response=await fetch('http://localhost:4000/Login', reqOptions);
    const result=await response.json();
    if(result.status===200){
      
      swal("Signed in!", "You have been Signed in!", "success");
   navigate("/")
    }
    else{
      swal("Error!", `${result.message}`, "error");
    }

  };

  return (
    <div className="login-container">
      <form onSubmit={handleSubmit} className="login-form">
        <h2>Login</h2>
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
        <button type="submit" className="submit-btn">Login</button>
      </form>
    </div>
  );
};

export default LoginForm;

