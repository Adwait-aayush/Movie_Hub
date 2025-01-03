import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"

export default function Search() {
    const [query, setquery] = useState("")
    const [movies, setmovies] = useState([])
    const [error, setError] = useState(null) 
    const navigate = useNavigate()

    useEffect(() => {
        if (query.length > 3) {
            const headers = new Headers()
            headers.append('Content-Type', 'application/json')
            const reqOptions = {
                method: 'GET',
                headers: headers
            }

            fetch(`http://localhost:4000/Search?name=${query}`, reqOptions)
                .then((response) => response.json())
                .then((data) => {
                    if (data && Array.isArray(data)) {
                        setmovies(data)
                        setError(null) 
                    } else {
                        setmovies([]) 
                        setError("No movies found.") 
                    }
                })
                .catch((err) => {
                    setmovies([]) 
                    setError("An error occurred while fetching movies.") 
                })
        } else {
            setmovies([]) 
        }
    }, [query])

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
            <h1>Search</h1>
            <hr />
            <form action="">
                <input
                    type="text"
                    value={query}
                    onChange={(e) => setquery(e.target.value)}
                    placeholder="Search for movies..."
                />
            </form>
            <hr />

            
            {error && <p>{error}</p>}

            
            {movies.length === 0 && !error && query.length > 3 && (
                <p>No Movies found. Add the movie You would Like to discuss about</p>
               
            
            )}

            <div className="movie-list">
                {movies.length > 0 &&
                    movies.map((movie, index) => (
                        <div
                            key={index}
                            className="movie-card"
                            onClick={() => navigate(`/movie/${movie.id}`)}
                        >
                            <img
                                src={`https://image.tmdb.org/t/p/w200/${movie.poster_path}`}
                                alt={movie.original_title}
                            />
                            <h3>{movie.original_title}</h3>
                            <div className="star-rating">
                                {getStars(movie.vote_average)}
                            </div>
                            <p>Release Date: {movie.release_date}</p>
                        </div>
                    ))}
            </div>
        </>
    );
}
