<template>
  <div>
    <header>
      <div class="container">
        <div class="header-content">
          <router-link to="/" class="logo">📧 Newsletter</router-link>
          <router-link to="/login" class="login-btn">Login</router-link>
        </div>
      </div>
    </header>
    
    <main>
      <div class="container">
        <div class="hero">
          <h1>Stay Updated</h1>
          <p>Get the latest news, insights, and updates delivered straight to your inbox. Join thousands of subscribers who trust us to keep them informed.</p>
          
          <NewsletterForm />
        </div>
        
        <div class="features">
          <FeatureCard 
            icon="🚀"
            title="Latest Updates"
            description="Be the first to know about new features, product launches, and important announcements."
          />
          
          <FeatureCard 
            icon="💡"
            title="Expert Insights"
            description="Get valuable insights and tips from industry experts to help you stay ahead of the curve."
          />
          
          <FeatureCard 
            icon="🎯"
            title="Curated Content"
            description="Receive handpicked content tailored to your interests, delivered weekly to your inbox."
          />
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import NewsletterForm from './NewsletterForm.vue'
import FeatureCard from './FeatureCard.vue'
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

:deep(body) {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  line-height: 1.6;
  color: #333;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

header {
  padding: 20px 0;
  position: relative;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  font-size: 2rem;
  font-weight: bold;
  color: white;
  text-decoration: none;
}

.login-btn {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 2px solid rgba(255, 255, 255, 0.3);
  padding: 12px 24px;
  border-radius: 25px;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.login-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  border-color: rgba(255, 255, 255, 0.5);
  transform: translateY(-2px);
}

main {
  text-align: center;
  padding: 80px 0;
  color: white;
}

.hero {
  max-width: 800px;
  margin: 0 auto;
}

.hero h1 {
  font-size: 3.5rem;
  margin-bottom: 20px;
  font-weight: 700;
}

.hero p {
  font-size: 1.3rem;
  margin-bottom: 40px;
  opacity: 0.9;
}

.features {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 30px;
  margin-top: 80px;
  padding: 0 20px;
}

@media (max-width: 768px) {
  .hero h1 {
    font-size: 2.5rem;
  }
  
  .hero p {
    font-size: 1.1rem;
  }
  
  .header-content {
    flex-direction: column;
    gap: 20px;
  }
}
</style>
