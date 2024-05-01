<template>
    <div class="upload-container">
      <input type="file" @change="handleFileChange" ref="fileInput" />
      <button @click="uploadImage" :disabled="!selectedFile">Upload</button>
    </div>
  </template>
  
  <script>
  import api from "@/services/axios";
  
  export default {
    data() {
      return {
        selectedFile: null,
      };
    },
    methods: {
      handleFileChange(event) {
        this.selectedFile = event.target.files[0];
      },
      async uploadImage() {
        if (!this.selectedFile) {
          alert("Please select a file to upload.");
          return;
        }
        const formData = new FormData();
        formData.append('image', this.selectedFile);
  
        try {
          const response = await api.post('/photos', formData, {
            headers: {
              'Content-Type': 'multipart/form-data',
               Authorization: `${localStorage.getItem('userId')}` // Assuming you use Bearer token
            }
          });
          alert('Upload successful!');
          console.log(response.data);
        } catch (error) {
          console.error('Upload failed:', error);
          alert('Upload failed!');
        }
      }
    }
  }
  </script>
  
  <style scoped>
  .upload-container {
    margin-top: 20px;
    justify-content: center;
    align-items: center;
    padding: 3px;
  }

  button {
    margin-top: 10px;
    background-color: #86c457;
    justify-content: center;
    align-items: center;
    margin-left: 60px;
  }
  
  </style>
  