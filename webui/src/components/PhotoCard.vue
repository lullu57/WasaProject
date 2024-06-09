<template>
  <div class="photo-card">
    <img :src="'data:image/jpeg;base64,' + photoData.imageData" alt="Photo" class="photo-image"/>
    <div class="photo-info">
      <h4>{{ photoData.username }}</h4>
      <p>{{ formatDate(photoData.timestamp) }}</p>
      <div class="photo-actions">
        <button @click="toggleLike">{{ isLiked ? 'Unlike' : 'Like' }} ({{ photoData.likesCount }})</button>
        <button @click="toggleComments">Comments ({{ photoData.comments.length }})</button>
        <!-- Delete photo button, visible only to the photo owner -->
        <button v-if="photoData.userId === userId" @click="deletePhoto(photoData.photoId)" class="delete-photo">Delete</button>
      </div>
      <div v-if="showComments" class="comments-section">
        <div class="comment-form">
          <input v-model="newComment" placeholder="Write a comment..." class="comment-input"/>
          <button @click="postComment" class="post-comment">Post</button>
        </div>
        <div class="comment" v-for="comment in photoData.comments" :key="comment.commentId">
          <strong>{{ comment.username }}</strong>: {{ comment.content }}
          <button v-if="comment.userId === userId" @click="deleteComment(comment.commentId)" class="delete-comment">Delete</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import api from '@/services/axios';

export default {
  props: {
    photo: Object,
    userId: String // Add userId as a prop
  },
  data() {
    return {
      showComments: true,
      newComment: '',
      photoData: { ...this.photo }, // Initialize local state with prop data
      isLiked: false
    };
  },
  mounted() {
    this.checkIfLiked();
    this.$root.$on('usernameChanged', this.updateUsername);
  },
  beforeUnmount() {
    this.$root.$off('usernameChanged', this.updateUsername);
  },
  methods: {
    async checkIfLiked() {
      const config = {
        headers: {
          Authorization: this.userId
        }
      };
      try {
        const response = await api.get(`/photos/${this.photoData.photoId}/likes`, config);
        this.isLiked = response.data.liked;
      } catch (error) {
        console.error('Failed to check like status', error);
      }
    },
    async toggleLike() {
      const config = {
        headers: {
          Authorization: this.userId
        }
      };
      try {
        if (!this.isLiked) {
          const response = await api.post(`/photos/${this.photoData.photoId}/likes`, {}, config);         
          this.photoData.likesCount++;
        } else {
          const response = await api.delete(`/photos/${this.photoData.photoId}/likes`, config);
          this.photoData.likesCount--;
        }
        this.isLiked = !this.isLiked;
      } catch (error) {
        console.error('Failed to toggle like', error);
      }
    },
    async postComment() {
      if (this.newComment.trim() !== '') {
        const config = {
          headers: {
            Authorization: this.userId
          }
        };
        try {
          const response = await api.post(`/photos/${this.photoData.photoId}/comments`, { content: this.newComment }, config);
          let username = 'You'; // Ideally fetch from server or use global state
          this.photoData.comments.push({
            username,
            content: this.newComment,
            commentId: response.data.commentId,
            userId: this.userId // Use the computed property
          });
          this.newComment = '';
        } catch (error) {
          console.error('Failed to post comment', error);
        }
      }
    },
    async deleteComment(commentId) {
      try {
        await api.delete(`/comments/${commentId}`, {
          headers: {
            Authorization: this.userId
          }
        });
        this.photoData.comments = this.photoData.comments.filter(comment => comment.commentId !== commentId);
      } catch (error) {
        console.error('Failed to delete comment', error);
      }
    },
    async deletePhoto(photoId) {
      try {
        await api.delete(`/photos/${photoId}`, {
          headers: {
            Authorization: this.userId
          }
        });
        // Emit an event to the parent component to remove the photo from the list
        this.$emit('photoDeleted', photoId);
      } catch (error) {
        console.error('Failed to delete photo', error);
      }
    },
    updateUsername(newUsername) {
      this.photoData.username = newUsername;
      this.photoData.comments.forEach(comment => {
        if (comment.userId === this.userId) {
          comment.username = newUsername;
        }
      });
    },
    formatDate(value) {
      return new Date(value).toLocaleString();
    }
  }
}
</script>

<style scoped>
.photo-card {
  border: 2px solid #ccc; /* increased border thickness for better definition */
  border-radius: 4px;
  padding: 10px; /* increased padding for better spacing */
  margin-bottom: 10px; /* more vertical space between cards */
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 4px 6px rgba(0,0,0,0.1); /* slightly heavier shadow for depth */
  width: 100%; /* ensures it's responsive to container size */
  max-width: 300px; /* suitable max-width for content balancing */
  height: auto; /* allows for variable height based on content */
}

.photo-image {
  width: 100%; /* ensures the image is responsive */
  height: auto; /* maintains aspect ratio */
  max-height: 200px; /* increased maximum height for better display */
  object-fit: contain; /* ensures the image fits without distortion */
  border-radius: 2px;
}

.photo-info {
  width: 100%; /* ensures text alignment container is full-width */
  text-align: center;
  padding-top: 10px;
}

.photo-actions {
  display: flex;
  justify-content: space-around; /* changes from space-between to space-around for better distribution */
  padding: 5px 0; /* adds padding around buttons */
}

.photo-actions button {
  background-color: #ff5e8a; /* softer pink for like button */
  color: white;
  border: none;
  border-radius: 4px;
  padding: 5px 10px; /* larger buttons for easier interaction */
  cursor: pointer;
  transition: background-color 0.3s;
}

.photo-actions button:hover {
  background-color: #df0b35; /* deeper pink on hover for interaction feedback */
}

.comments-section {
  width: 100%; /* full width to align with the photo card */
  margin-top: 10px;
}

.comment-form {
  display: flex;
  justify-content: space-between;
  width: 100%; /* full width for better layout control */
  margin-top: 10px; /* added space above the comment form */
  margin-bottom: 5px;
}

.comment-input {
  border: 1px solid #ccc; /* refined border styling */
  border-radius: 4px;
  padding: 5px; /* increased padding for better text entry */
}

button.post-comment {
  background-color: #007bff; /* blue for post button */
  color: white; /* text color for readability */
  border-radius: 4px; /* rounded corners match other design elements */
  padding: 8px 10px; /* adjusted padding for size fitting */
  font-size: 0.9rem; /* increased font size for readability */
}

.comment {
  background-color: #f0f0f0; /* light gray background for comments */
  padding: 8px; /* increased padding for aesthetic spacing */
  border-radius: 4px; /* consistent rounded corners */
  margin-top: 4px; /* space between comments */
}
</style>
