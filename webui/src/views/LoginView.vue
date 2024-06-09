<template>
  <div class="login-container">
    <h1>Login</h1>
    <input v-model="username" placeholder="Enter username" class="input-field"/>
    <button @click="login" class="login-button">Login</button>
    <p v-if="error" class="error-message">{{ error }}</p>
  </div>
</template>

<script>
import axios from 'axios';
import api from "@/services/axios"; 

export default {
  data() {
    return {
      username: '',
      error: ''
    };
  },
  methods: {
    async login() {
      try {
        const response = await api.post('/session', { name: this.username });
        localStorage.setItem("userId", response.data.token);
        axios.defaults.headers.common['Authorization'] = response.data;
        this.$router.push('/stream').then(() => {
          window.location.reload(); 
        });
      } catch (err) {
        this.error = 'Failed to login. Please try again.';
        console.error(err);
      }
    }
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding-top: 50px; /* Added padding at the top */
  width: 100%;
}

.input-field, .login-button {
  width: 80%; /* Full width of the container */
  padding: 10px;
  margin: 10px 0; /* Margin for spacing */
}

.login-button {
  background-color: #007bff; /* Blue color for button */
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.login-button:hover {
  background-color: #0056b3; /* Darker blue on hover */
}

.error-message {
  color: red;
}
</style>
