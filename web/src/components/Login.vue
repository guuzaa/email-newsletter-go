<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-header">
        <router-link to="/" class="back-link">← Back to Home</router-link>
        <div class="logo">📧</div>
        <h1>Welcome Back</h1>
        <p>Sign in to access your newsletter dashboard</p>
      </div>
    
    <div v-if="errorMessage" class="error-message">
      {{ errorMessage }}
    </div>
    
    <div v-if="successMessage" class="success-message">
      {{ successMessage }}
    </div>
    
    <form @submit.prevent="login" method="POST" action="/login">
      <div class="form-group">
        <label for="username">Username</label>
        <input 
          type="text" 
          id="username" 
          name="username"
          v-model="form.username" 
          placeholder="Enter your username"
          autocomplete="username"
          required
        >
      </div>
      
      <div class="form-group">
        <label for="password">Password</label>
        <input 
          type="password" 
          id="password" 
          name="password"
          v-model="form.password" 
          placeholder="Enter your password"
          autocomplete="current-password"
          required
        >
      </div>
      
      <button type="submit" class="login-btn" :disabled="isSubmitting">
        {{ isSubmitting ? 'Signing in...' : 'Sign In' }}
      </button>
    </form>
    
    <div class="forgot-password">
      <a href="#" @click.prevent="showForgotPassword">Forgot your password?</a>
    </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface LoginForm {
  username: string
  password: string
}

const form = ref<LoginForm>({
  username: '',
  password: ''
})

const isSubmitting = ref<boolean>(false)
const errorMessage = ref<string>('')
const successMessage = ref<string>('')

onMounted(() => {
  // Check for URL parameters for messages
  const urlParams = new URLSearchParams(window.location.search)
  const error = urlParams.get('error')
  const success = urlParams.get('success')
  
  if (error) {
    errorMessage.value = 'Invalid username or password. Please try again.'
  }
  if (success) {
    successMessage.value = 'Login successful! Redirecting...'
  }
})

const login = async (): Promise<void> => {
  isSubmitting.value = true
  errorMessage.value = ''
  successMessage.value = ''
  
  // For now, we'll use the traditional form submission
  // since the Go backend expects form data
  const formElement = document.querySelector('form') as HTMLFormElement
  formElement.submit()
}

const showForgotPassword = (): void => {
  alert('Please contact support to reset your password.')
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-container {
  background: rgba(255, 255, 255, 0.95);
  padding: 40px;
  border-radius: 20px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
  backdrop-filter: blur(10px);
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.back-link {
  display: block;
  text-align: left;
  margin-bottom: 20px;
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
  font-size: 0.9rem;
}

.back-link:hover {
  text-decoration: underline;
}

.logo {
  font-size: 2rem;
  margin-bottom: 10px;
}

.login-header h1 {
  color: #333;
  margin-bottom: 10px;
  font-size: 1.8rem;
}

.login-header p {
  color: #666;
  font-size: 0.95rem;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #555;
  font-weight: 500;
}

.form-group input {
  width: 100%;
  padding: 15px;
  border: 2px solid #e1e5e9;
  border-radius: 10px;
  font-size: 1rem;
  transition: border-color 0.3s ease;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
}

.login-btn {
  width: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 15px;
  border-radius: 10px;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.3s ease;
  margin-bottom: 20px;
}

.login-btn:hover {
  transform: translateY(-2px);
}

.login-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
}

.back-link {
  text-align: center;
}

.back-link a {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
}

.back-link a:hover {
  text-decoration: underline;
}

.error-message {
  background: #f8d7da;
  color: #721c24;
  padding: 15px;
  border-radius: 10px;
  margin-bottom: 20px;
  border: 1px solid #f5c6cb;
}

.success-message {
  background: #d4edda;
  color: #155724;
  padding: 15px;
  border-radius: 10px;
  margin-bottom: 20px;
  border: 1px solid #c3e6cb;
}

.forgot-password {
  text-align: center;
  margin-top: 15px;
}

.forgot-password a {
  color: #666;
  text-decoration: none;
  font-size: 0.9rem;
}

.forgot-password a:hover {
  text-decoration: underline;
}

@media (max-width: 480px) {
  .login-container {
    padding: 30px 20px;
  }
  
  .login-header h1 {
    font-size: 1.5rem;
  }
}
</style>
