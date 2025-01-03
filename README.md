# Movie Hub

Movie Hub is a platform where movie enthusiasts can discover, discuss, and share opinions on movies. It's similar to Reddit, but for movies. Users can browse popular movies, participate in discussions, and even add movies to the platform that they want to talk about.

## Features

- **Popular Movies:** Displays a list of trending movies pulled from the TMDB API. Users can click on a movie to view details and join discussions.
  
- **Movie Discussions:** Allows users to comment on movies and engage with others who share similar interests. 

- **Add Movies:** If a movie is not in the database, users can add movies by providing the title and other details. These movies will then be available for discussion.

- **User Authentication:** Users can create sessions to log in and personalize their experience, keeping track of their posts and comments.

## Technologies Used

- **Frontend:** React for building the user interface.
- **Backend:** Go with Chi for routing and handling backend logic.
- **Database:** MongoDB to store user data, comments, and movie details.
- **Authentication:** Session-based user authentication for personalized experiences.
- **API:** TMDB API to fetch popular movie data.


