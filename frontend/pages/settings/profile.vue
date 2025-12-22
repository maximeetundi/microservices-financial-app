<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-4xl mx-auto animate-fade-in-up">
      <!-- Header -->
      <div class="mb-10">
        <h1 class="text-3xl font-bold text-base mb-2">Mon profil 👤</h1>
        <p class="text-muted">Vos informations personnelles</p>
      </div>

      <!-- Security Notice -->
      <div class="glass-card mb-8 p-6 bg-blue-500/10 border border-blue-500/30">
        <div class="flex items-start gap-4">
          <div class="w-12 h-12 rounded-xl bg-blue-500/20 flex items-center justify-center flex-shrink-0">
            <svg class="w-6 h-6 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
            </svg>
          </div>
          <div>
            <h3 class="font-bold text-blue-500 mb-1">🔒 Protection anti-fraude activée</h3>
            <p class="text-sm text-muted">
              Pour votre sécurité, la modification des informations personnelles est désactivée. 
              Pour effectuer un changement, veuillez contacter notre équipe support avec les justificatifs nécessaires.
            </p>
          </div>
        </div>
      </div>

      <!-- Profile Card -->
      <div class="glass-card mb-6">
        <div class="flex items-center gap-6 mb-8 p-6 border-b border-[var(--border-color)]">
          <!-- Avatar -->
          <div class="w-24 h-24 rounded-full bg-gradient-to-br from-primary to-purple-600 flex items-center justify-center text-3xl font-bold text-white">
            {{ userInitials }}
          </div>
          <div>
            <h2 class="text-2xl font-bold text-base">{{ profile.first_name }} {{ profile.last_name }}</h2>
            <p class="text-muted">Membre depuis {{ formatDate(profile.created_at) }}</p>
            <div class="flex items-center gap-2 mt-2">
              <span class="px-3 py-1 rounded-full bg-success/10 text-success text-sm font-medium flex items-center gap-1">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                </svg>
                Compte vérifié
              </span>
            </div>
          </div>
        </div>

        <!-- Personal Information -->
        <div class="p-6 space-y-6">
          <h3 class="text-sm font-bold text-muted mb-4">INFORMATIONS PERSONNELLES</h3>
          
          <div class="grid md:grid-cols-2 gap-6">
            <ReadOnlyField 
              label="Prénom" 
              :value="profile.first_name" 
              icon="person"
            />
            <ReadOnlyField 
              label="Nom" 
              :value="profile.last_name" 
              icon="person"
            />
            <ReadOnlyField 
              label="Email" 
              :value="profile.email" 
              icon="email"
              verified
            />
            <ReadOnlyField 
              label="Téléphone" 
              :value="profile.phone || 'Non renseigné'" 
              icon="phone"
              :verified="!!profile.phone"
            />
            <ReadOnlyField 
              label="Date de naissance" 
              :value="profile.birth_date || 'Non renseigné'" 
              icon="calendar"
            />
            <ReadOnlyField 
              label="Nationalité" 
              :value="profile.nationality || 'Non renseigné'" 
              icon="flag"
            />
          </div>
        </div>
      </div>

      <!-- Address -->
      <div class="glass-card mb-6 p-6">
        <h3 class="text-sm font-bold text-muted mb-4">ADRESSE</h3>
        <div class="grid md:grid-cols-2 gap-6">
          <ReadOnlyField 
            label="Adresse" 
            :value="profile.address || 'Non renseigné'" 
            icon="home"
          />
          <ReadOnlyField 
            label="Ville" 
            :value="profile.city || 'Non renseigné'" 
            icon="location"
          />
          <ReadOnlyField 
            label="Code postal" 
            :value="profile.postal_code || 'Non renseigné'" 
            icon="tag"
          />
          <ReadOnlyField 
            label="Pays" 
            :value="profile.country || 'Non renseigné'" 
            icon="globe"
          />
        </div>
      </div>

      <!-- Verification Status -->
      <div class="glass-card mb-6 p-6">
        <h3 class="text-sm font-bold text-muted mb-4">STATUT DE VÉRIFICATION</h3>
        <div class="space-y-4">
          <VerificationRow label="KYC (Identité)" :verified="profile.kyc_verified" />
          <VerificationRow label="Email" :verified="profile.email_verified" />
          <VerificationRow label="Téléphone" :verified="profile.phone_verified" />
          <VerificationRow label="Adresse" :verified="profile.address_verified" />
        </div>
      </div>

      <!-- Contact Support Button -->
      <div class="text-center">
        <NuxtLink 
          to="/support" 
          class="inline-flex items-center gap-2 px-8 py-4 rounded-xl bg-gradient-to-r from-primary to-purple-600 text-white font-bold hover:scale-105 transition-transform"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 5.636l-3.536 3.536m0 5.656l3.536 3.536M9.172 9.172L5.636 5.636m3.536 9.192l-3.536 3.536M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-5 0a4 4 0 11-8 0 4 4 0 018 0z"/>
          </svg>
          Contacter le support pour modifier mes informations
        </NuxtLink>
        <p class="text-sm text-muted mt-4">
          Une demande de modification sera traitée sous 24 à 48 heures après vérification des documents.
        </p>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { userAPI } from '@/composables/useApi'

const profile = ref({
  first_name: 'John',
  last_name: 'Doe',
  email: 'john.doe@email.com',
  phone: '+33 6 12 34 56 78',
  birth_date: '15/03/1990',
  nationality: 'France',
  address: '123 Rue de Paris',
  city: 'Paris',
  postal_code: '75001',
  country: 'France',
  created_at: new Date(),
  kyc_verified: true,
  email_verified: true,
  phone_verified: true,
  address_verified: true,
})

const userInitials = computed(() => {
  const first = profile.value.first_name?.[0] || ''
  const last = profile.value.last_name?.[0] || ''
  return (first + last).toUpperCase()
})

onMounted(async () => {
  try {
    const res = await userAPI.getProfile()
    if (res.data) {
      profile.value = { ...profile.value, ...res.data }
    }
  } catch (e) {
    console.error('Error loading profile:', e)
  }
})

function formatDate(date) {
  if (!date) return 'N/A'
  return new Date(date).toLocaleDateString('fr-FR', { month: 'long', year: 'numeric' })
}

definePageMeta({
  layout: false,
  middleware: 'auth'
})
</script>

<script>
// ReadOnlyField component
const ReadOnlyField = {
  props: ['label', 'value', 'icon', 'verified'],
  template: `
    <div class="p-4 rounded-xl bg-surface-hover">
      <div class="flex items-center justify-between mb-1">
        <span class="text-xs text-muted uppercase tracking-wider">{{ label }}</span>
        <div class="flex items-center gap-1">
          <svg v-if="verified" class="w-4 h-4 text-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          <svg class="w-4 h-4 text-muted/50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
          </svg>
        </div>
      </div>
      <p class="font-medium text-base">{{ value }}</p>
    </div>
  `
}

// VerificationRow component  
const VerificationRow = {
  props: ['label', 'verified'],
  template: `
    <div class="flex items-center justify-between p-4 rounded-xl bg-surface-hover">
      <span class="font-medium text-base">{{ label }}</span>
      <span 
        class="px-3 py-1 rounded-full text-sm font-medium flex items-center gap-1"
        :class="verified ? 'bg-success/10 text-success' : 'bg-warning/10 text-warning'"
      >
        <svg v-if="verified" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
        </svg>
        <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01"/>
        </svg>
        {{ verified ? 'Vérifié' : 'En attente' }}
      </span>
    </div>
  `
}

export default {
  components: { ReadOnlyField, VerificationRow }
}
</script>

<style scoped>
.animate-fade-in-up {
  animation: fadeInUp 0.5s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
