import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
 import { faUser } from "@fortawesome/free-regular-svg-icons"
 import { useNavigate } from "react-router-dom";
 import { useState,useEffect } from "react";
export default function Signin(){
    const Navigate = useNavigate();
     const[name,setname]=useState("")
          
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
      const navigate = useNavigate();
    return(
        <>
        {name ===""&&(
        <button className="signinbut" onClick={()=>Navigate("/Login")}>
                <FontAwesomeIcon icon={faUser} size="lg" /> <p>Signin</p>
              </button>
              )}
        </>
    )
}