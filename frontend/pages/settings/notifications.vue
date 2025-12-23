<template>
  <NuxtLayout name="dashboard">
    <div class="notif-prefs-page">
      <!-- Header -->
      <div class="page-header">
        <NuxtLink to="/settings" class="back-link">← Paramètres</NuxtLink>
        <h1>🔔 Notifications</h1>
        <p>Gérez vos alertes et communications</p>
      </div>

      <!-- Push Notifications -->
      <div class="section">
        <h2>Notifications Push</h2>
        
        <div class="pref-card">
          <div class="pref-item">
            <div class="pref-info">
              <h4>Activer les notifications push</h4>
              <p>Recevez des alertes en temps réel</p>
            </div>
            <label class="toggle">
              <input type="checkbox" v-model="prefs.push.enabled">
              <span class="slider"></span>
            </label>
          </div>
        </div>
      </div>

      <!-- Transaction Alerts -->
      <div class="section">
        <h2>Alertes de transaction</h2>
        
        <div class="pref-card">
          <div class="pref-item">
            <div class="pref-icon">💸</div>
            <div class="pref-info">
              <h4>Transferts reçus</h4>
              <p>Quand vous recevez de l'argent</p>
            </div>
            <label class="toggle">
              <input type="checkbox" v-model="prefs.transactions.received">
              <span class="slider"></span>
            </label>
          </div>

          <div class="pref-divider"></div>

          <div class="pref-item">
            <div class="pref-icon">📤</div>
            <div class="pref-info">
              <h4>Transferts envoyés</h4>
              <p>Confirmation d'envoi</p>
            </div>
            <label class="toggle">
              <input type="checkbox" v-model="prefs.transactions.sent">
              <span class="slider"></span>
            </label>
          </div>

          <div class="pref-divider"></div>

          <div class="pref-item">
            <div class="pref-icon">💳</div>
            <div class="pref-info">
              <h4>Paiements carte</h4>
              <p>Utilisation de vos cartes</p>
            </div>
            <label class="toggle">
              <input type="checkbox" v-model="prefs.transactions.card">
              <span class="slider"></span>
            </label>
          </div>

          <div class="pref-divider"></div>

          <div class="pref-item">
            <div class="pref-icon">⚠️</div>
            <div class="pref-info">
              <h4>Solde faible</h4>
              <p>Alerte quand le solde est bas</p>
            </div>
            <label class="toggle">
              <input type="checkbox" v-model="prefs.transactions.lowBalance">
              <span class="slider"></span>
            </label>
          </div>
        </div>
      </div>

      <!-- Security Alerts -->
      <div class="section">
        <h2>Alertes de sécurité</h2>
        
        <div class="pref-card">
          <div class="pref-item">
            <div class="pref-icon">🔐</div>
            <div class="pref-info">
              <h4>Nouvelle connexion</h4>
              <p>Connexion depuis un nouvel appareil</p>
            </div>
            <label class="toggle">
              <input type="checkbox" v-model="prefs.security.newLogin">
              <span class="slider"></span>
            </label>
          </div>

          <div class="pref-divider"></div>

          <div class="pref-item">
            <div class="pref-icon">🔑</div>
            <div class="pref-info">
              <h4>Changement de mot de passe</h4>
              <p>Modification des identifiants</p>
            </div>
            <label class="toggle">
              <input type="checkbox" v-model="prefs.security.passwordChange">
              <span class="slider"></span>
            </label>
          </div>

          <div class="pref-divider"></div>

          <div class="pref-item">
            <div class="pref-icon">📱</div>
            <div class="pref-info">
              <h4>Code OTP</h4>
              <p>Recevoir les codes par SMS</p>
            </div>
            <label class="toggle">
              <input type="checkbox" v-model="prefs.security.otpSms">
              <span class="slider"></span>
            </label>
          </div>
        </div>
      </div>

      <!-- Email Preferences -->
      <div class="section">
        <h2>Emails</h2>
        
        <div class="pref-card">
          <div class="pref-item">
            <div class="pref-icon">📧</div>
            <div class="pref-info">
              <h4>Résumé hebdomadaire</h4>
              <p>Bilan de vos transactions</p>
            </div>
            <label class="toggle">
              <input type="checkbox" v-model="prefs.email.weeklyReport">
              <span class="slider"></span>
            </label>
          </div>

          <div class="pref-divider"></div>

          <div class="pref-item">
            <div class="pref-icon">📰</div>
            <div class="pref-info">
              <h4>Newsletter</h4>
              <p>Nouveautés et offres</p>
            </div>
            <label class="toggle">
              <input type="checkbox" v-model="prefs.email.newsletter">
              <span class="slider"></span>
            </label>
          </div>

          <div class="pref-divider"></div>

          <div class="pref-item">
            <div class="pref-icon">🎁</div>
            <div class="pref-info">
              <h4>Promotions</h4>
              <p>Offres spéciales et réductions</p>
            </div>
            <label class="toggle">
              <input type="checkbox" v-model="prefs.email.promotions">
              <span class="slider"></span>
            </label>
          </div>
        </div>
      </div>

      <!-- Save Button -->
      <button @click="savePrefs" :disabled="saving" class="save-btn">
        {{ saving ? 'Enregistrement...' : '💾 Enregistrer' }}
      </button>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'

const saving = ref(false)

const prefs = reactive({
  push: {
    enabled: true,
  },
  transactions: {
    received: true,
    sent: true,
    card: true,
    lowBalance: true,
  },
  security: {
    newLogin: true,
    passwordChange: true,
    otpSms: true,
  },
  email: {
    weeklyReport: false,
    newsletter: false,
    promotions: false,
  }
})

const savePrefs = async () => {
  saving.value = true
  try {
    localStorage.setItem('notificationPrefs', JSON.stringify(prefs))
    await new Promise(r => setTimeout(r, 500))
    alert('Préférences de notification enregistrées!')
  } catch (e) {
    console.error('Save error:', e)
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  const saved = localStorage.getItem('notificationPrefs')
  if (saved) {
    try {
      const data = JSON.parse(saved)
      Object.assign(prefs, data)
    } catch (e) {
      console.error('Load error:', e)
    }
  }
})

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>

<style scoped>
.notif-prefs-page {
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
  font-size: 1.25rem;
  width: 2rem;
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

.toggle {
  position: relative;
  width: 48px;
  height: 28px;
  cursor: pointer;
  flex-shrink: 0;
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
  margin-top: 0.5rem;
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
