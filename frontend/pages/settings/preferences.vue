<template>
  <NuxtLayout name="dashboard">
    <div class="prefs-page">
      <!-- Header -->
      <div class="page-header">
        <NuxtLink to="/settings" class="back-link">← Paramètres</NuxtLink>
        <h1>🎨 Préférences</h1>
        <p>Personnalisez votre expérience</p>
      </div>

      <!-- Appearance -->
      <div class="section">
        <h2>Apparence</h2>
        
        <div class="pref-card">
          <div class="pref-item">
            <div class="pref-icon">🌙</div>
            <div class="pref-info">
              <h4>Thème</h4>
              <p>Apparence de l'application</p>
            </div>
            <div class="theme-toggle">
              <button 
                v-for="theme in themes" 
                :key="theme.id"
                @click="setTheme(theme.id)"
                :class="['theme-btn', { active: currentTheme === theme.id }]"
              >
                {{ theme.icon }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Language & Region -->
      <div class="section">
        <h2>Langue & Région</h2>
        
        <div class="pref-card">
          <div class="pref-item">
            <div class="pref-icon">🌍</div>
            <div class="pref-info">
              <h4>Langue</h4>
              <p>{{ getLanguageName(selectedLanguage) }}</p>
            </div>
            <select v-model="selectedLanguage" class="pref-select">
              <option value="fr">Français</option>
              <option value="en">English</option>
              <option value="es">Español</option>
              <option value="ar">العربية</option>
            </select>
          </div>

          <div class="pref-divider"></div>

          <div class="pref-item">
            <div class="pref-icon">💵</div>
            <div class="pref-info">
              <h4>Devise par défaut</h4>
              <p>Affichage des montants</p>
            </div>
            <select v-model="selectedCurrency" class="pref-select">
              <option value="USD">USD ($)</option>
              <option value="EUR">EUR (€)</option>
              <option value="XOF">XOF (FCFA)</option>
              <option value="XAF">XAF (FCFA)</option>
              <option value="GBP">GBP (£)</option>
            </select>
          </div>

          <div class="pref-divider"></div>

          <div class="pref-item">
            <div class="pref-icon">🕐</div>
            <div class="pref-info">
              <h4>Fuseau horaire</h4>
              <p>{{ selectedTimezone }}</p>
            </div>
            <select v-model="selectedTimezone" class="pref-select">
              <option value="Europe/Paris">Paris (UTC+1)</option>
              <option value="Africa/Dakar">Dakar (UTC+0)</option>
              <option value="Africa/Lagos">Lagos (UTC+1)</option>
              <option value="America/New_York">New York (UTC-5)</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Display -->
      <div class="section">
        <h2>Affichage</h2>
        
        <div class="pref-card">
          <div class="pref-item">
            <div class="pref-icon">🔢</div>
            <div class="pref-info">
              <h4>Format des nombres</h4>
              <p>1 000,00 ou 1,000.00</p>
            </div>
            <select v-model="numberFormat" class="pref-select">
              <option value="fr">1 000,00 (FR)</option>
              <option value="en">1,000.00 (EN)</option>
            </select>
          </div>

          <div class="pref-divider"></div>

          <div class="pref-item">
            <div class="pref-icon">📅</div>
            <div class="pref-info">
              <h4>Format de date</h4>
              <p>Ordre jour/mois/année</p>
            </div>
            <select v-model="dateFormat" class="pref-select">
              <option value="DD/MM/YYYY">DD/MM/YYYY</option>
              <option value="MM/DD/YYYY">MM/DD/YYYY</option>
              <option value="YYYY-MM-DD">YYYY-MM-DD</option>
            </select>
          </div>

          <div class="pref-divider"></div>

          <div class="pref-item">
            <div class="pref-icon">👁️</div>
            <div class="pref-info">
              <h4>Masquer les soldes</h4>
              <p>Pour plus de confidentialité</p>
            </div>
            <label class="toggle">
              <input type="checkbox" v-model="hideBalances">
              <span class="slider"></span>
            </label>
          </div>
        </div>
      </div>

      <!-- Privacy -->
      <div class="section">
        <h2>Confidentialité</h2>
        
        <div class="pref-card">
          <div class="pref-item">
            <div class="pref-icon">📊</div>
            <div class="pref-info">
              <h4>Analytiques</h4>
              <p>Aider à améliorer l'app</p>
            </div>
            <label class="toggle">
              <input type="checkbox" v-model="analyticsEnabled">
              <span class="slider"></span>
            </label>
          </div>

          <div class="pref-divider"></div>

          <div class="pref-item">
            <div class="pref-icon">🛡️</div>
            <div class="pref-info">
              <h4>Verrouillage auto</h4>
              <p>Après inactivité</p>
            </div>
            <select v-model="autoLockTime" class="pref-select">
              <option value="1">1 minute</option>
              <option value="5">5 minutes</option>
              <option value="15">15 minutes</option>
              <option value="30">30 minutes</option>
              <option value="0">Jamais</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Save Button -->
      <button @click="savePreferences" :disabled="saving" class="save-btn">
        {{ saving ? 'Enregistrement...' : '💾 Enregistrer les préférences' }}
      </button>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'

const saving = ref(false)

// Theme
const themes = [
  { id: 'system', icon: '💻' },
  { id: 'light', icon: '☀️' },
  { id: 'dark', icon: '🌙' },
]
const currentTheme = ref('dark')

// Language & Region
const selectedLanguage = ref('fr')
const selectedCurrency = ref('XOF')
const selectedTimezone = ref('Europe/Paris')

// Display
const numberFormat = ref('fr')
const dateFormat = ref('DD/MM/YYYY')
const hideBalances = ref(false)

// Privacy
const analyticsEnabled = ref(true)
const autoLockTime = ref('5')

const getLanguageName = (code) => {
  const names = { fr: 'Français', en: 'English', es: 'Español', ar: 'العربية' }
  return names[code] || code
}

const setTheme = (theme) => {
  currentTheme.value = theme
  
  if (theme === 'system') {
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
    document.documentElement.classList.toggle('dark', prefersDark)
  } else {
    document.documentElement.classList.toggle('dark', theme === 'dark')
  }
  
  localStorage.setItem('theme', theme)
}

const savePreferences = async () => {
  saving.value = true
  try {
    const prefs = {
      theme: currentTheme.value,
      language: selectedLanguage.value,
      currency: selectedCurrency.value,
      timezone: selectedTimezone.value,
      numberFormat: numberFormat.value,
      dateFormat: dateFormat.value,
      hideBalances: hideBalances.value,
      analyticsEnabled: analyticsEnabled.value,
      autoLockTime: autoLockTime.value,
    }
    
    localStorage.setItem('userPreferences', JSON.stringify(prefs))
    
    // Simulate API save
    await new Promise(r => setTimeout(r, 500))
    
    alert('Préférences enregistrées!')
  } catch (e) {
    console.error('Save error:', e)
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  // Load saved preferences
  const saved = localStorage.getItem('userPreferences')
  if (saved) {
    try {
      const prefs = JSON.parse(saved)
      currentTheme.value = prefs.theme || 'dark'
      selectedLanguage.value = prefs.language || 'fr'
      selectedCurrency.value = prefs.currency || 'XOF'
      selectedTimezone.value = prefs.timezone || 'Europe/Paris'
      numberFormat.value = prefs.numberFormat || 'fr'
      dateFormat.value = prefs.dateFormat || 'DD/MM/YYYY'
      hideBalances.value = prefs.hideBalances || false
      analyticsEnabled.value = prefs.analyticsEnabled !== false
      autoLockTime.value = prefs.autoLockTime || '5'
    } catch (e) {
      console.error('Error loading prefs:', e)
    }
  }
  
  // Load theme
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme) {
    currentTheme.value = savedTheme
  }
})

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>

<style scoped>
.prefs-page {
  width: 100%;
  max-width: 600px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 1.5rem;
}

.back-link {
  display: inline-block;
  color: #888;
  text-decoration: none;
  font-size: 0.875rem;
  margin-bottom: 0.5rem;
}

.page-header h1 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #fff;
  margin: 0 0 0.25rem 0;
}

.page-header p {
  font-size: 0.875rem;
  color: #888;
  margin: 0;
}

.section {
  margin-bottom: 1.5rem;
}

.section h2 {
  font-size: 0.75rem;
  font-weight: 600;
  color: #888;
  text-transform: uppercase;
  margin: 0 0 0.75rem 0;
}

.pref-card {
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 1rem;
  overflow: hidden;
}

.pref-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
}

.pref-divider {
  height: 1px;
  background: rgba(255,255,255,0.06);
  margin: 0 1rem;
}

.pref-icon {
  font-size: 1.5rem;
  width: 2.5rem;
  text-align: center;
}

.pref-info {
  flex: 1;
  min-width: 0;
}

.pref-info h4 {
  font-size: 0.875rem;
  font-weight: 500;
  color: #fff;
  margin: 0;
}

.pref-info p {
  font-size: 0.75rem;
  color: #666;
  margin: 0;
}

.pref-select {
  padding: 0.5rem 0.75rem;
  border-radius: 0.5rem;
  border: 1px solid rgba(255,255,255,0.1);
  background: rgba(255,255,255,0.05);
  color: #fff;
  font-size: 0.75rem;
  outline: none;
  max-width: 120px;
}

.theme-toggle {
  display: flex;
  gap: 0.25rem;
  background: rgba(255,255,255,0.05);
  border-radius: 0.5rem;
  padding: 0.25rem;
}

.theme-btn {
  width: 2rem;
  height: 2rem;
  border: none;
  background: transparent;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: all 0.2s;
}

.theme-btn.active {
  background: #6366f1;
}

.toggle {
  position: relative;
  width: 48px;
  height: 28px;
  cursor: pointer;
}

.toggle input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  inset: 0;
  background: rgba(255,255,255,0.1);
  border-radius: 28px;
  transition: 0.3s;
}

.slider::before {
  content: '';
  position: absolute;
  width: 22px;
  height: 22px;
  left: 3px;
  bottom: 3px;
  background: #fff;
  border-radius: 50%;
  transition: 0.3s;
}

.toggle input:checked + .slider {
  background: #6366f1;
}

.toggle input:checked + .slider::before {
  transform: translateX(20px);
}

.save-btn {
  width: 100%;
  padding: 1rem;
  border-radius: 0.875rem;
  border: none;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #fff;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  margin-top: 1rem;
}

.save-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.save-btn:not(:disabled):hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.3);
}
</style>
