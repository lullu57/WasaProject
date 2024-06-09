<template>
  <div class="discover-users">
    <h2>Discover Users</h2>
    <input v-model="searchQuery" @input="searchUsers" placeholder="Search users..." class="search-box" />
    <ul class="user-list">
      <li v-for="user in filteredUsers" :key="user.userId">
        <router-link :to="{ name: 'Profile', params: { profileId: user.userId } }">
          {{ user.username }}
        </router-link>
        <!-- Hide follow and ban buttons for the user's own profile -->
        <button v-if="user.userId !== localStorageUserId" @click="toggleFollow(user)" :disabled="user.processing">
          {{ user.isFollowing ? 'Unfollow' : 'Follow' }}
        </button>
        <button v-if="user.userId !== localStorageUserId" @click="toggleBan(user)" :disabled="user.processing">
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
      users: [],
      searchQuery: '',
      localStorageUserId: localStorage.getItem('userId'), // Get the user's ID from localStorage
    };
  },
  computed: {
    filteredUsers() {
      const query = this.searchQuery.toLowerCase();
      return this.users.filter(user => user.username.toLowerCase().includes(query));
    }
  },
  async mounted() {
    await this.fetchUsers();
  },
  methods: {
    async fetchUsers() {
      try {
        const response = await api.get(`/users`, {
          headers: { Authorization: this.localStorageUserId }
        });
        const users = response.data.map(user => ({
          ...user,
          isFollowing: false,
          isBanned: false,
          processing: false
        }));

        // Separate the user's own profile and move it to the top
        const ownProfile = users.find(user => user.userId === this.localStorageUserId);
        this.users = ownProfile ? [ownProfile, ...users.filter(user => user.userId !== this.localStorageUserId)] : users;

        await this.checkFollowAndBanStatus();
      } catch (error) {
        console.error('Failed to fetch users:', error);
      }
    },
    async checkFollowAndBanStatus() {
      try {
        await Promise.all(this.users.map(async (user) => {
          const followRes = await api.get(`/follows/${user.userId}`, {
            headers: { Authorization: this.localStorageUserId }
          });
          user.isFollowing = followRes.data.isFollowed;

          const banRes = await api.get(`/bans/${user.userId}`, {
            headers: { Authorization: this.localStorageUserId }
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
          await api.delete(`/users/${user.userId}/followers`, {
            headers: { Authorization: this.localStorageUserId }
          });
          user.isFollowing = false;
        } else {
          await api.post(`/users/${user.userId}/followers`, {}, {
            headers: { Authorization: this.localStorageUserId }
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
            headers: { Authorization: this.localStorageUserId }
          });
          user.isBanned = false;
        } else {
          await api.post(`/users/${user.userId}/bans`, {}, {
            headers: { Authorization: this.localStorageUserId }
          });
          user.isBanned = true;
        }
      } catch (error) {
        console.error('Failed to toggle ban:', error);
      } finally {
        user.processing = false;
      }
    },
    searchUsers() {
      this.$forceUpdate();
    }
  }
}
</script>

<style scoped>
.discover-users {
  padding: 20px;
}

.search-box {
  margin-bottom: 20px;
  padding: 5px;
  width: 100%;
  box-sizing: border-box;
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
