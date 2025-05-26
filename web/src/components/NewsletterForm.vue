<template>
  <div class="newsletter-form">
    <h2>Subscribe to Our Newsletter</h2>
    
    <div v-if="successMessage" class="success-message">
      {{ successMessage }}
    </div>
    
    <div v-if="errorMessage" class="error-message">
      {{ errorMessage }}
    </div>
    
    <form @submit.prevent="subscribe">
      <div class="form-group">
        <label for="name">Full Name</label>
        <input 
          type="text" 
          id="name" 
          v-model="form.name" 
          placeholder="Enter your full name"
          required
        >
      </div>
      
      <div class="form-group">
        <label for="email">Email Address</label>
        <input 
          type="email" 
          id="email" 
          v-model="form.email" 
          placeholder="Enter your email address"
          required
        >
      </div>
      
      <button type="submit" class="submit-btn" :disabled="isSubmitting">
        {{ isSubmitting ? 'Subscribing...' : 'Subscribe Now' }}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface NewsletterForm {
  name: string
  email: string
}

const form = ref<NewsletterForm>({
  name: '',
  email: ''
})

const isSubmitting = ref<boolean>(false)
const successMessage = ref<string>('')
const errorMessage = ref<string>('')

const subscribe = async (): Promise<void> => {
  isSubmitting.value = true
  successMessage.value = ''
  errorMessage.value = ''
  
  try {
    const response = await fetch('/subscriptions', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name: form.value.name,
        email: form.value.email
      })
    })
    
    if (response.ok) {
      successMessage.value = 'Thank you for subscribing! Please check your email to confirm your subscription.'
      form.value.name = ''
      form.value.email = ''
    } else {
      const errorData = await response.json()
      errorMessage.value = errorData.message || 'Something went wrong. Please try again.'
    }
  } catch (error) {
    errorMessage.value = 'Network error. Please check your connection and try again.'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
.newsletter-form {
  background: rgba(255, 255, 255, 0.95);
  padding: 40px;
  border-radius: 20px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  max-width: 500px;
  margin: 0 auto;
  backdrop-filter: blur(10px);
}

.newsletter-form h2 {
  color: #333;
  margin-bottom: 20px;
  font-size: 1.8rem;
}

.form-group {
  margin-bottom: 20px;
  text-align: left;
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

.submit-btn {
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
}

.submit-btn:hover {
  transform: translateY(-2px);
}

.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
}

.success-message {
  background: #d4edda;
  color: #155724;
  padding: 15px;
  border-radius: 10px;
  margin-bottom: 20px;
  border: 1px solid #c3e6cb;
}

.error-message {
  background: #f8d7da;
  color: #721c24;
  padding: 15px;
  border-radius: 10px;
  margin-bottom: 20px;
  border: 1px solid #f5c6cb;
}

@media (max-width: 768px) {
  .newsletter-form {
    padding: 30px 20px;
    margin: 0 20px;
  }
}
</style>
