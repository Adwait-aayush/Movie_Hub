import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faRegistered } from "@fortawesome/free-regular-svg-icons";
import { useNavigate } from "react-router-dom";
import { useState, useEffect } from "react";

export default function Register() {
  const navigate = useNavigate();
  const [name, setName] = useState("");

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

  }, []);

  return (
    <>
      {name === "" && (
        <button className="registerbtn" onClick={() => navigate("/Registration")}>
          <FontAwesomeIcon icon={faRegistered} size="lg" /> <p>Register</p>
        </button>
      )}
    </>
  );
}
