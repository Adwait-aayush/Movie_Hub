import React, { useState, useEffect } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faHouse } from "@fortawesome/free-solid-svg-icons";
import { faSearch } from "@fortawesome/free-solid-svg-icons";
import { faEdit } from "@fortawesome/free-solid-svg-icons";
import { faUpload } from "@fortawesome/free-solid-svg-icons";
import { faArrowTrendUp } from "@fortawesome/free-solid-svg-icons";
import { faArrowRightFromBracket } from "@fortawesome/free-solid-svg-icons";
import './Navbuttons.css'
import { useNavigate } from "react-router-dom";

export default function Navbuttons() {
  const navigate = useNavigate();
  const [name, setName] = useState("");

  const Logout = () => {
    const headers = new Headers();
    headers.append('Content-Type', 'application/json');

    const reqOptions = {
      method: 'POST',
      headers: headers,
      credentials: "include"
    };

    fetch(`http://localhost:4000/Logout`, reqOptions)
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
        setName(""); // Clear the username after logout
        window.location.reload();
      });
  };

  useEffect(() => {
    const headers = new Headers();
    headers.append('Content-Type', 'application/json');
    const reqOptions = {
      method: 'GET',
      headers: headers,
      credentials: 'include'
    };

    fetch(`http://localhost:4000/Username`, reqOptions)
      .then(response => response.json())
      .then(data => setName(data.username));

    
    return () => setName(""); 
  }, []);

  return (
    <>
      <div className="navbuttons">
        <button className="home" onClick={() => navigate("/")}>
          <FontAwesomeIcon icon={faHouse} size="3x" /> <p>Home</p>
        </button>

        <button className="home" onClick={()=>navigate("/Search")}>
          <FontAwesomeIcon icon={faSearch} size="3x" /> <p>Search</p>
        </button>

        {name !== "" && (
          <button onClick={() => navigate("/addmov")} className="home">
            <FontAwesomeIcon icon={faEdit} size="3x" /> <p>Add Movie</p>
          </button>
        )}

        {name !== "" && (
          <button onClick={()=>navigate("/UserMovies")} className="home">
            <FontAwesomeIcon icon={faUpload} size="3x" />
            <p>Your Movies</p>
          </button>
        )}

        
      </div>
      <hr className="three-d-rule2" />

      {name !== "" && (
        <div className="logout">
          <button className="home" onClick={Logout}>
            <FontAwesomeIcon icon={faArrowRightFromBracket} size="3x" /> <p>Logout</p>
          </button>
        </div>
      )}
    </>
  );
}
