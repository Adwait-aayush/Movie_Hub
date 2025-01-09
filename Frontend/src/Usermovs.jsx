import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import Swal from 'sweetalert2';


export default function UserMovies() {
    const [movies, setMovies] = useState([]);
    const [name, setName] = useState("");
    const [error, setError] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const navigate = useNavigate();

    useEffect(() => {
        const fetchUsername = async () => {
            try {
                const response = await fetch("http://localhost:4000/Username", {
                    method: "GET",
                    headers: { "Content-Type": "application/json" },
                    credentials: "include",
                });
                if (!response.ok) throw new Error("Failed to fetch username");
                const data = await response.json();
                setName(data.username || "");
            } catch (err) {
                setError("Failed to fetch username");
                console.error("Error fetching username:", err);
            }
        };

        fetchUsername();
    }, []);

    useEffect(() => {
        const fetchMovies = async () => {
            if (name) {
                try {
                    const response = await fetch(`http://localhost:4000/user?name=${name}`, {
                        method: "GET",
                        headers: { "Content-Type": "application/json" },
                    });
                    if (!response.ok) throw new Error("You Have not uploaded any movies yet");
                    const data = await response.json();
                    setMovies(data.filter(movie => movie !== null));
                } catch (err) {
                    setError("You Have not uploaded any movies yet");
                    console.error("Error fetching movies:", err);
                } finally {
                    setIsLoading(false);
                }
            } else {
                setIsLoading(false);
            }
        };

        fetchMovies();
    }, [name]);

    const deleteMovie = async (movieId) => {
        try {
            const response = await fetch(`http://localhost:4000/Delete/${movieId}`, {
                method: 'DELETE',
                headers: { 'Content-Type': 'application/json' },
            });

            if (!response.ok) throw new Error("Failed to delete movie");

            const data = await response.json();
            console.log(data);

           
            await Swal.fire({
                title: 'Success!',
                text: 'Movie deleted successfully',
                icon: 'success',
                confirmButtonText: 'OK'
            });

           
            setMovies(movies.filter(movie => movie.id !== movieId));
        } catch (err) {
            console.error("Error deleting movie:", err);
            Swal.fire({
                title: 'Error!',
                text: 'Failed to delete movie',
                icon: 'error',
                confirmButtonText: 'OK'
            });
        }
    }

    if (isLoading) return <div className="loading">Loading...</div>;
    if (error) return <div className="error">{error}</div>;

    return (
        <div className="user-movies-container">
            <h2 className="user-movies-title">Your Movies</h2>
            {movies.length === 0 ? (
                <p className="no-movies-message">No movies found. Add the movies you would like to discuss.</p>
            ) : (
                <div className="movie-list">
                    {movies.map((movie, index) => (
                        <div key={index} className="movie-card">
                            <img
                                src={movie.poster_path ? `https://image.tmdb.org/t/p/w200/${movie.poster_path}` : '/placeholder.svg?height=300&width=200'}
                                alt={movie.original_title || 'Movie poster'}
                                className="movie-poster"
                            />
                            <h3 className="movie-title">{movie.original_title || 'Unknown Title'}</h3>
                            <p className="release-date">Release Date: {movie.release_date || 'Unknown'}</p>
                            <button
                                className="delete-movie-button"
                                onClick={() => deleteMovie(movie.id)}
                            >
                                Delete
                            </button>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}

