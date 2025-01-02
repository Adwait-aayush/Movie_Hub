import { useEffect, useState, useRef } from 'react';
import './Homepage.css';
import Logo from "./Logo.jpeg";
import Signin from './Signin';
import Register from './Register';
import { Link, Navigate, useNavigate } from 'react-router-dom';

export default function Homepage() {
    const [movies, setMovies] = useState([]);
    const movieListRef = useRef(null);
   const navigate=useNavigate()
    useEffect(() => {
        const headers = new Headers();
        headers.append('Content-Type', 'application/json');
        const reqOptions = {
            method: 'GET',
            headers: headers
        };
        fetch('http://localhost:4000/pop', reqOptions)
            .then(response => response.json())
            .then(data => setMovies(data));
    }, []);

    const slideLeft = () => {
        if (movieListRef.current) {
            movieListRef.current.scrollBy({
                left: -250,
                behavior: 'smooth'
            });
        }
    };

    const slideRight = () => {
        if (movieListRef.current) {
            movieListRef.current.scrollBy({
                left: 250,
                behavior: 'smooth'
            });
        }
    };

    const getStars = (rating) => {
        const fullStars = Math.floor(rating / 2);
        const halfStar = rating % 2 !== 0;
        const stars = [];

        for (let i = 0; i < fullStars; i++) {
            stars.push(<span key={i}>★</span>);
        }
        if (halfStar) {
            stars.push(<span key={fullStars}>☆</span>);
        }

        return stars;
    };

    return (
        <>

            <div className="Homepage">
                <div className='Header'>
                    <div className="logo-container">
                        <img className="Logo" src={Logo} alt="" />
                        <h1>Movie_Hub</h1>
                    </div>
                    <div className="button-container">
                        <Signin/>
                        <Register/>
                    </div>
                </div>
                <hr />
                <div className="LatestMovies">
                    <h1>Latest Movies</h1>
                    <div className="slider-container">
                        <div ref={movieListRef} className="movie-list">
                            {movies.map((movie, index) => (
                                <div key={index} className="movie-card" onClick={() => navigate(`/movie/${movie.id}`)}>
                                    <img src={`https://image.tmdb.org/t/p/w200/${movie.poster_path}`} alt={movie.title} />
                                    <h3>{movie.original_title}</h3>
                                    <div className="star-rating">
                                        {getStars(movie.vote_average)}
                                    </div>
                                    <p>Release Date: {movie.release_date}</p>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        </>
    );
}
