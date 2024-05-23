<template>
    <div class="discover-users">
      <h2>Discover Users</h2>
      <ul class="user-list">
        <li v-for="user in users" :key="user.userId">
          <router-link :to="{ name: 'Profile', params: { profileId: user.userId } }">
            {{ user.username }}
          </router-link>
          <button @click="toggleFollow(user)" :disabled="user.processing">
            {{ user.isFollowing ? 'Unfollow' : 'Follow' }}
          </button>
          <!-- Toggle Ban/Unban -->
          <button @click="toggleBan(user)" :disabled="user.processing">
            {{ user.isBanned ? 'Unban' : 'Ban' }}
          </button>
        </li>
      </ul>
    </div>
  </template>
  
  <script>
  import api from '@/services/axios';
  
  export default {
    data() {
      return {
        users: []
      };
    },
    async mounted() {
      await this.fetchUsers();
    },
    methods: {
      async fetchUsers() {
        try {
          const response = await api.get('/users', {
            headers: { Authorization: localStorage.getItem('userId') }
          });
          this.users = response.data.map(user => ({
            ...user,
            isFollowing: false,
            isBanned: false,
            processing: false
          }));
          await this.checkFollowAndBanStatus();
        } catch (error) {
          console.error('Failed to fetch users:', error);
        }
      },
      async checkFollowAndBanStatus() {
        try {
          await Promise.all(this.users.map(async (user) => {
            const followRes = await api.get(`/follows/${user.userId}`, {
              headers: { Authorization: localStorage.getItem('userId') }
            });
            user.isFollowing = followRes.data.isFollowed;
  
            const banRes = await api.get(`/bans/${user.userId}`, {
              headers: { Authorization: localStorage.getItem('userId') }
            });
            user.isBanned = banRes.data.banned;
          }));
        } catch (error) {
          console.error('Failed to fetch follow/ban status:', error);
        }
      },
      async toggleFollow(user) {
        user.processing = true;
        try {
          if (user.isFollowing) {
            await api.delete(`/users/${user.userId}/follows`, {
              headers: { Authorization: localStorage.getItem('userId') }
            });
            user.isFollowing = false;
          } else {
            await api.post(`/users/${user.userId}/follows`, {}, {
              headers: { Authorization: localStorage.getItem('userId') }
            });
            user.isFollowing = true;
          }
        } catch (error) {
          console.error('Failed to toggle follow:', error);
        } finally {
          user.processing = false;
        }
      },
      async toggleBan(user) {
        user.processing = true;
        try {
          if (user.isBanned) {
            await api.delete(`/users/${user.userId}/bans`, {
              headers: { Authorization: localStorage.getItem('userId') }
            });
            user.isBanned = false;
          } else {
            await api.post(`/users/${user.userId}/bans`, {}, {
              headers: { Authorization: localStorage.getItem('userId') }
            });
            user.isBanned = true;
          }
        } catch (error) {
          console.error('Failed to toggle ban:', error);
        } finally {
          user.processing = false;
        }
      }
    }
  }
  </script>
  
  <style scoped>
  .discover-users {
    padding: 20px;
  }
  
  .user-list {
    list-style-type: none;
    padding: 0;
  }
  
  .user-list li {
    margin-bottom: 10px;
    display: flex;
    align-items: center;
  }
  
  button {
    margin-left: 10px;
  }
  </style>
  