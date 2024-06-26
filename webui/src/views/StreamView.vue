<template>
  <div class="stream-view">
    <div v-if="photos.length > 0" class="gallery">
      <PhotoCard 
        v-for="photo in photos" 
        :key="photo.photoId"
        :photo="photo"
        :user-id="userId" 
      />
    </div>
    <div v-else>
      <p>No photos to display. Start following people to see their photos here.</p>
    </div>
  </div>
</template>

<script>
import PhotoCard from '@/components/PhotoCard.vue';
import api from '@/services/axios';

export default {
  components: {
    PhotoCard
  },
  data() {
    return {
      photos: [],
      error: '',
      userId: localStorage.getItem('userId') // Store userId in a data property
    };
  },
  async mounted() {
    await this.fetchStreamPhotos();
  },
  methods: {
    async fetchStreamPhotos() {
      try {
        const response = await api.get('/stream', {
          headers: { Authorization: this.userId }
        });
        const photoIds = response.data; // Assuming this is an array of photo IDs
        if (photoIds && photoIds.length > 0) {
          await this.fetchPhotoDetails(photoIds);
        } else {
          this.photos = [];
        }
      } catch (error) {
        console.error('Failed to fetch stream photos:', error);
        this.error = "Failed to load photos. Please try again later.";
      }
    },
    async fetchPhotoDetails(photoIds) {
      this.photos = await Promise.all(photoIds.map(async (photoId) => {
        try {
          const res = await api.get(`/photos/${photoId}`, {
            headers: { Authorization: this.userId }
          });
          const photo = res.data;
          // Process comments
          photo.comments = await Promise.all(photo.comments.map(async (comment) => {
            const userResponse = await api.get(`/users/${comment.userId}/username`);
            comment.username = userResponse.data.username;
            comment.isOwner = comment.userId === this.userId;
            return comment;
          }));
          return photo;
        } catch (error) {
          console.error("Error fetching photo details:", error);
          return {}; // Return empty object or handle errors appropriately
        }
      }));
    }
  }
}
</script>

<style scoped>
.stream-view {
  padding: 20px;
}

p {
  color: #666;
  text-align: center;
}

.gallery {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); /* Adjust minmax for desired card width */
  gap: 20px; /* Adjust gap for spacing between cards */
  justify-content: center; /* Center cards in the gallery if they don't fill all columns */
  align-items: start; /* Align items at the start of the grid line */
}
</style>
