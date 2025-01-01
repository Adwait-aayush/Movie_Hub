import { useState,useEffect } from "react"
import './USnameandpp.css'
import Guest from './Guest.jpg'
import { useOutletContext } from "react-router-dom"
export default function USnameandpp(){
    const[name,setname]=useState("Guest")
    
useEffect(()=>{
    const headers=new Headers()
    headers.append('Content-Type','application/json')
    const reqoptions={
      method:'GET',
      headers:headers,
      credentials:'include'
  
    }
    fetch(`http://localhost:4000/Username`,reqoptions)
    .then(response=>response.json())
    .then(data=>setname(data.username))
  },[])
   
   const[pp,setpp]=useState(Guest)
    return(
        <>
        <div className="profile">
        <div className="profilepic">
            <img src={pp} alt="profilepic" />
        </div>
        <div className="username">
            <h2>{name}</h2>
        </div>
        </div>
        
        </>
    )
}