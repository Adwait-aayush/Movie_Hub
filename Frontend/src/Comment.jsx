import React, { useState, useEffect } from 'react';
import './Comment.css';
import { useParams } from 'react-router-dom';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

const Comment = () => {
  const { id } = useParams();
  const [name, setName] = useState("Guest");
  const [comments, setComments] = useState([]);
  const [newComment, setNewComment] = useState('');
  const [replyingTo, setReplyingTo] = useState(null);
  const [reply, setReply] = useState('');
  const [shownReplies, setShownReplies] = useState({});

  useEffect(() => {
    const headers = new Headers();
    headers.append('Content-Type', 'application/json');
    const reqOptions = {
      method: 'GET',
      headers: headers,
      credentials: 'include',
    };

    fetch('http://localhost:4000/Username', reqOptions)
      .then(response => response.json())
      .then(data => setName(data.username))
      .catch((err) => console.error("Error fetching username:", err));
  }, []);

  useEffect(() => {
    const headers = new Headers();
    headers.append('Content-Type', 'application/json');
    const reqOptions = {
      method: 'GET',
      headers: headers,
    };

    fetch(`http://localhost:4000/comments/${id}`, reqOptions)
      .then(response => response.json())
      .then(data => setComments(data || []))
      .catch((err) => {
        console.error("Error fetching comments:", err);
        setComments([]);
      });
  }, [id]);

  const handleAddComment = () => {
    if (newComment.trim()) {
      const body = {
        movie_id: id,
        comment_id: "",
        message: newComment,
        likes: 0,
        dislikes: 0,
        replies: [],
        author: name,
      };

      const headers = new Headers();
      headers.append('Content-Type', 'application/json');
      const reqOptions = {
        method: 'POST',
        headers: headers,
        body: JSON.stringify(body),
      };

      fetch('http://localhost:4000/comments', reqOptions)
        .then(response => response.json())
        .then(() => {
          setNewComment('');
          toast.success("Comment added successfully!");
          
          fetch(`http://localhost:4000/comments/${id}`)
            .then(response => response.json())
            .then(data => setComments(data || []));
        })
        .catch((err) => console.error("Error posting comment:", err));
    }
  };

  const handleReply = (commentId) => {
    setReplyingTo(commentId === replyingTo ? null : commentId);
  };

  const handleAddReply = async (commentId) => {
    try {
      const response = await fetch(`http://localhost:4000/replycomments/${commentId}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        }
      });

      if (!response.ok) {
        throw new Error('Failed to fetch comment data');
      }

      const commentData = await response.json();
      console.log('Fetched comment data:', commentData);

      const updatedCommentData = {
        ...commentData,
        replies: [
          ...commentData.replies,
          {
            movie_id: id,
            comment_id: "",
            message: reply,
            likes: 0,
            dislikes: 0,
            replies: [],
            author: name,
          }
        ]
      };

      const updateResponse = await fetch('http://localhost:4000/addreply', {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(updatedCommentData)
      });

      if (!updateResponse.ok) {
        throw new Error('Failed to add reply');
      }

      const result = await updateResponse.json();
      console.log('Reply added successfully:', result);

      setReply('');
      setReplyingTo(null);

      // Refresh comments
      const refreshResponse = await fetch(`http://localhost:4000/comments/${id}`);
      const refreshedComments = await refreshResponse.json();
      setComments(refreshedComments || []);

    } catch (error) {
      console.error('Error adding reply:', error);
    }
  };

  const handleLike = (commentId) => {
    const updatedComments = comments.map(comment => {
      if (comment.comment_id === commentId) {
        return { ...comment, likes: comment.likes + 1 };
      }
      return comment;
    });
    setComments(updatedComments);
  };

  const handleDislike = (commentId) => {
    const updatedComments = comments.map(comment => {
      if (comment.comment_id === commentId) {
        return { ...comment, dislikes: comment.dislikes + 1 };
      }
      return comment;
    });
    setComments(updatedComments);
  };

  const handleDelete = (commentId) => {
    const updatedComments = comments.filter(comment => comment.comment_id !== commentId);
    setComments(updatedComments);
    toast.success("Comment deleted successfully!");
  };

  const toggleReplies = (commentId) => {
    setShownReplies(prev => ({
      ...prev,
      [commentId]: !prev[commentId]
    }));
  };

  const renderComment = (comment, isReply = false) => (
    <div key={comment.comment_id} className={`comment ${isReply ? 'reply' : ''}`}>
      <h4>{comment.author}</h4>
      <p>{comment.message}</p>
      <div className="comment-actions">
        {!isReply && <button onClick={() => handleReply(comment.comment_id)}>Reply</button>}
        <button onClick={() => handleLike(comment.comment_id)}>Like ({comment.likes})</button>
        <button onClick={() => handleDislike(comment.comment_id)}>Dislike ({comment.dislikes})</button>
        {comment.author === name && (
          <button onClick={() => handleDelete(comment.comment_id)}>Delete</button>
        )}
      </div>
      {replyingTo === comment.comment_id && (
        <div className="add-comment">
          <textarea
            value={reply}
            onChange={(e) => setReply(e.target.value)}
            placeholder="Write a reply..."
          />
          <button onClick={() => handleAddReply(comment.comment_id)}>Post Reply</button>
        </div>
      )}
      {!isReply && comment.replies && comment.replies.length > 0 && (
        <>
          <button 
            className="show-replies-button" 
            onClick={() => toggleReplies(comment.comment_id)}
          >
            {shownReplies[comment.comment_id] ? 'Hide' : 'Show'} {comment.replies.length} {comment.replies.length === 1 ? 'reply' : 'replies'}
          </button>
          {shownReplies[comment.comment_id] && (
            <div className="replies-container">
              {comment.replies.map(reply => renderComment(reply, true))}
            </div>
          )}
        </>
      )}
    </div>
  );

  return (
    <div className="comment-section">
      <ToastContainer />
      <h2>Comments</h2>
      {comments && comments.length > 0 ? (
        comments.map(comment => (
          <div key={comment.comment_id}>
            {renderComment(comment)}
          </div>
        ))
      ) : (
        <p>No comments yet. Be the first to comment!</p>
      )}
      <div className="add-comment">
        <textarea
          value={newComment}
          onChange={(e) => setNewComment(e.target.value)}
          placeholder="Add a comment..."
        />
        <button onClick={handleAddComment}>Post Comment</button>
      </div>
    </div>
  );
};

export default Comment;

