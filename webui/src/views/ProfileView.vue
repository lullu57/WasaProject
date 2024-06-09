<template>
  <div class="profile-view">
    <div v-if="userProfile && !isBanned && !isBannedByProfileOwner" class="info-container">
      <p>Username: {{ userProfile.username }}</p>
      <input v-if="isOwnProfile" v-model="newUsername" placeholder="Change username" />
      <button v-if="isOwnProfile" @click="changeUsername">Change Username</button>
      <p>Followers: {{ userProfile.followers?.length || '0' }}</p>
      <p>Following: {{ userProfile.following?.length || '0' }}</p>
      <p>Posts: {{ detailedPhotos.length || '0' }}</p>
      <!-- Follow/Unfollow button -->
      <button v-if="!isOwnProfile" @click="userProfile.isFollowing ? unfollowUser() : followUser()">
        {{ userProfile.isFollowing ? 'Unfollow' : 'Follow' }}
      </button>
      <!-- Ban/Unban button -->
      <button v-if="!isOwnProfile" @click="userProfile.isBanned ? unbanUser() : banUser()">
        {{ userProfile.isBanned ? 'Unban' : 'Ban' }}
      </button>
    </div>
    <div v-else-if="isBanned || isBannedByProfileOwner">
      <p>This profile is not accessible.</p>
    </div>
    <div v-else>
      <p>Loading profile...</p>
    </div>
    <div class="gallery" v-if="userProfile && !isBanned && !isBannedByProfileOwner">
      <PhotoCard
        v-for="photo in detailedPhotos"
        :key="photo.photoId"
        :photo="photo"
        :user-id="localStorageUserId"
        @photoDeleted="handlePhotoDeleted"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import api from "@/services/axios";
import PhotoCard from '@/components/PhotoCard.vue';

const route = useRoute();
const router = useRouter();
const userId = ref(route.params.profileId);
const userProfile = ref(null);
const newUsername = ref('');
const detailedPhotos = ref([]);
const localStorageUserId = localStorage.getItem('userId');
const isOwnProfile = computed(() => userId.value === localStorageUserId);
const isBanned = ref(false);
const isBannedByProfileOwner = ref(false);

const fetchUserProfile = async () => {
  try {
    const response = await api.get(`/users/${userId.value}`);
    userProfile.value = response.data;
    if (!isOwnProfile.value) { // Check if the profile is not the user's own
      await checkIfUserIsFollowed(); // Check if the user is following the profile user
      await checkIfUserIsBanned(); // Check if the user has banned the profile user
      await checkIfUserIsBannedByProfileOwner(); // Check if the user is banned by the profile owner
      if (isBannedByProfileOwner.value || isBanned.value) {
        return; // If the user is banned, stop further processing
      }
    }
    if (userProfile.value && userProfile.value.photos) {
      await fetchPhotoDetails(userProfile.value.photos);
    }
  } catch (error) {
    console.error("Error fetching user profile:", error);
    console.log(`Request failed with status code ${error.response?.status}: ${error.response?.data}`);
  }
};

const fetchPhotoDetails = async (photoIds) => {
  try {
    detailedPhotos.value = await Promise.all(photoIds.map(async (id) => {
      const res = await api.get(`/photos/${id}`, {
            headers: { Authorization: this.userId }
          });
      const photo = res.data;
      photo.comments = await Promise.all(photo.comments.map(async (comment) => {
        const userResponse = await api.get(`/users/${comment.userId}/username`);
        comment.username = userResponse.data.username;
        return comment;
      }));
      return photo;
    }));
  } catch (error) {
    console.error("Error fetching photo details:", error);
  }
};

const checkIfUserIsFollowed = async () => {
  try {
    const response = await api.get(`/follows/${userId.value}`, {
      headers: {
        Authorization: `${localStorage.getItem('userId')}`
      }
    });
    userProfile.value.isFollowing = response.data.isFollowed; // Ensure this matches the key returned by your API
  } catch (error) {
    console.error("Error checking if user is followed:", error);
  }
};

const checkIfUserIsBanned = async () => {
  try {
    const response = await api.get(`/bans/${userId.value}`, {
      headers: {
        Authorization: `${localStorage.getItem('userId')}`
      }
    });
    isBanned.value = response.data.banned; // Ensure this matches the key returned by your API
  } catch (error) {
    console.error("Error checking if user is banned:", error);
  }
};

const checkIfUserIsBannedByProfileOwner = async () => {
  try {
    const response = await api.get(`/bans/${localStorageUserId}`, {
      headers: {
        Authorization: `${userId.value}`
      }
    });
    isBannedByProfileOwner.value = response.data.banned; // Ensure this matches the key returned by your API
  } catch (error) {
    console.error("Error checking if user is banned by profile owner:", error);
  }
};

const followUser = async () => {
  await api.post(`/users/${userId.value}/followers`, {}, {
    headers: { Authorization: localStorageUserId }
  });
  userProfile.value.isFollowing = true;
};

const unfollowUser = async () => {
  await api.delete(`/users/${userId.value}/followers`, {
    headers: { Authorization: localStorageUserId }
  });
  userProfile.value.isFollowing = false;
};

const banUser = async () => {
  await api.post(`/users/${userId.value}/bans`, {}, {
    headers: { Authorization: localStorageUserId }
  });
  userProfile.value.isBanned = true;
  isBanned.value = true;
};

const unbanUser = async () => {
  await api.delete(`/users/${userId.value}/bans`, {
    headers: { Authorization: localStorageUserId }
  });
  userProfile.value.isBanned = false;
  isBanned.value = false;
};

const changeUsername = async () => {
  try {
    await api.patch(`/users/username`, {
      newUsername: newUsername.value
    }, {
      headers: { Authorization: localStorageUserId }
    });
    userProfile.value.username = newUsername.value; // Update the username in the view
    newUsername.value = ''; // Clear the input field
    alert('Username changed successfully!');
    router.push('/profile/' + userProfile.value.userId); // navigate to the new profile url
    window.location.reload(); // refresh the page

  } catch (error) {
    console.error("Error changing username:", error);
    if (error.response && error.response.status === 409) {
      alert('Username already taken. Please choose another one.');
    } else {
      alert('Failed to change username.');
    }
  }
};

const handlePhotoDeleted = (photoId) => {
  detailedPhotos.value = detailedPhotos.value.filter(photo => photo.photoId !== photoId);
};

onMounted(fetchUserProfile);

watch(route, async (newRoute) => {
  userId.value = newRoute.params.profileId;
  await fetchUserProfile();
});
</script>

<style scoped>
.profile-view {
  padding: 20px;
}

.info-container {
  background-color: #f4f4f4;
  padding: 20px;
  margin-bottom: 20px;
}

.gallery {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); /* Adjust minmax for desired card width */
  gap: 20px; /* Adjust gap for spacing between cards */
  justify-content: center; /* Center cards in the gallery if they don't fill all columns */
  align-items: start; /* Align items at the start of the grid line */
}

input[type="text"] {
  display: block;
  margin-top: 10px;
  padding: 8px;
  width: 100%;
}

button {
  margin-top: 10px;
}
</style>
