import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import App from './App';
import Homepage from './Homepage';
import './index.css';
import MovieForm from './MovieFrom';
import Registration from './Registration';
import LoginForm from './LoginForm';
import Movie from './Movie';
import Usermovs from './Usermovs';
import Search from './search';

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      { index: true, element: <Homepage /> },
      { path: '/addmov', element: <MovieForm /> },
      {path :'/movie/:id',element:<Movie/>},
      {path :'/Search',element:<Search/>},
      {path :'/UserMovies',element:<Usermovs/>},
    ],
  },
  { path: '/Registration', element: <Registration /> }, {path:"/Login",element:<LoginForm/>},
]);

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>
);
