import { Outlet } from 'react-router-dom'
import Navigationbar from './Navigationbar'
import './App.css' 
import { useState,useEffect } from 'react'

function App() {


  
  return (
    <div className="app-container">
      <Navigationbar />
      <main className="content">
        <Outlet        
       />
      </main>
    </div>
  )
}

export default App

