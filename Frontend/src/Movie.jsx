import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import "./Movie.css";
import Comment from './Comment';

export default function Movie() {
    const { id } = useParams();
    const [movie, setMovie] = useState(null);

    useEffect(() => {
        const headers = new Headers();
        headers.append('Content-Type', 'application/json');
        const reqOptions = {
            method: 'GET',
            headers: headers
        };
        fetch(`http://localhost:4000/movie/${id}`, reqOptions)
            .then(response => response.json())
            .then(data => setMovie(data));
    }, [id]);

    if (!movie) {
        return <div className="loading">Loading...</div>;
    }

    return (
        <div className="container">
            <div className="movie-container">
                <div className="movie-poster-container">
                    <img src={`https://image.tmdb.org/t/p/w500/${movie.poster_path}`} alt={movie.title} className="movie-poster" />
                </div>
                <div className="movie-info">
                    <h1 className="movie-title">{movie.title}</h1>
                    <p className="movie-overview">{movie.overview}</p>
                    <div className="movie-meta">
                        <div className="meta-item">
                            <span className="meta-label">Release Date</span>
                            <span className="meta-value">{movie.release_date}</span>
                        </div>
                        <div className="meta-item">
                            <span className="meta-label">Rating</span>
                            <span className="meta-value">{movie.vote_average.toFixed(1)}/10</span>
                        </div>
                        <div className="meta-item">
                            <span className="meta-label">Popularity</span>
                            <span className="meta-value">{movie.popularity.toFixed(0)}</span>
                        </div>
                        <div className="meta-item">
                            <span className="meta-label">Language</span>
                            <span className="meta-value">{movie.original_language.toUpperCase()}</span>
                        </div>
                    </div>
                </div>
            </div>
            <Comment />
        </div>
    );
}

