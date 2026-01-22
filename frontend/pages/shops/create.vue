<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-3xl mx-auto py-8 px-4">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-8">🏪 Créer ma boutique</h1>

      <form @submit.prevent="createShop" class="space-y-6">
        <div class="bg-white dark:bg-slate-900 rounded-2xl p-6 border border-gray-100 dark:border-gray-800">
          <h2 class="text-lg font-bold text-gray-900 dark:text-white mb-4">Informations générales</h2>
          
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Nom de la boutique *</label>
              <input v-model="form.name" type="text" required class="w-full input-premium" placeholder="Ma Super Boutique">
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Description</label>
              <textarea v-model="form.description" rows="3" class="w-full input-premium" placeholder="Décrivez votre boutique..."></textarea>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Tags (séparés par des virgules)</label>
              <input v-model="tagsInput" type="text" class="w-full input-premium" placeholder="mode, vêtements, accessoires">
            </div>

            <div class="flex items-center gap-3">
              <input type="checkbox" v-model="form.is_public" id="isPublic" class="w-5 h-5 rounded">
              <label for="isPublic" class="text-gray-700 dark:text-gray-300">Boutique publique (visible dans le marketplace)</label>
            </div>
          </div>
        </div>

        <div class="bg-white dark:bg-slate-900 rounded-2xl p-6 border border-gray-100 dark:border-gray-800">
          <h2 class="text-lg font-bold text-gray-900 dark:text-white mb-4">💰 Configuration des paiements</h2>
          
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Portefeuille pour recevoir les paiements *</label>
              <select v-model="form.wallet_id" required class="w-full input-premium">
                <option value="">Sélectionner un portefeuille</option>
                <option v-for="w in wallets" :key="w.id" :value="w.id">
                  {{ w.currency }} - Solde: {{ formatPrice(w.balance, w.currency) }}
                </option>
              </select>
              <p class="text-sm text-gray-500 mt-1">Les paiements des clients seront crédités sur ce portefeuille</p>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Devise de la boutique *</label>
              <select v-model="form.currency" required class="w-full input-premium">
                <option value="">Sélectionner une devise</option>
                <option value="XOF">XOF - Franc CFA</option>
                <option value="EUR">EUR - Euro</option>
                <option value="USD">USD - Dollar US</option>
                <option value="GBP">GBP - Livre Sterling</option>
              </select>
            </div>
          </div>
        </div>

        <div class="bg-white dark:bg-slate-900 rounded-2xl p-6 border border-gray-100 dark:border-gray-800">
          <h2 class="text-lg font-bold text-gray-900 dark:text-white mb-4">🖼️ Images (optionnel)</h2>
          
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Logo</label>
              <div class="aspect-square border-2 border-dashed border-gray-200 dark:border-gray-700 rounded-xl flex items-center justify-center cursor-pointer hover:border-indigo-500 transition-colors">
                <input type="file" @change="uploadLogo" accept="image/*" class="hidden" ref="logoInput">
                <div v-if="form.logo_url" class="w-full h-full">
                  <img :src="form.logo_url" class="w-full h-full object-cover rounded-xl">
                </div>
                <div v-else @click="($refs.logoInput as HTMLInputElement).click()" class="text-center p-4">
                  <div class="text-4xl mb-2">📷</div>
                  <p class="text-sm text-gray-500">Cliquer pour ajouter</p>
                </div>
              </div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Bannière</label>
              <div class="aspect-video border-2 border-dashed border-gray-200 dark:border-gray-700 rounded-xl flex items-center justify-center cursor-pointer hover:border-indigo-500 transition-colors">
                <input type="file" @change="uploadBanner" accept="image/*" class="hidden" ref="bannerInput">
                <div v-if="form.banner_url" class="w-full h-full">
                  <img :src="form.banner_url" class="w-full h-full object-cover rounded-xl">
                </div>
                <div v-else @click="($refs.bannerInput as HTMLInputElement).click()" class="text-center p-4">
                  <div class="text-4xl mb-2">🖼️</div>
                  <p class="text-sm text-gray-500">Cliquer pour ajouter</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="flex gap-4">
          <NuxtLink to="/shops" class="flex-1 py-4 text-center bg-gray-100 dark:bg-slate-800 text-gray-700 dark:text-gray-300 rounded-xl font-bold hover:bg-gray-200 dark:hover:bg-slate-700 transition-colors">
            Annuler
          </NuxtLink>
          <button type="submit" :disabled="loading" class="flex-1 btn-premium py-4 disabled:opacity-50">
            {{ loading ? 'Création...' : 'Créer ma boutique' }}
          </button>
        </div>
      </form>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
import { useShopApi } from '~/composables/useShopApi'
import { useApi } from '~/composables/useApi'

definePageMeta({ middleware: 'auth' })

const shopApi = useShopApi()
const { walletApi } = useApi()

const loading = ref(false)
const wallets = ref<any[]>([])
const tagsInput = ref('')

const form = reactive({
  name: '',
  description: '',
  is_public: true,
  wallet_id: '',
  currency: 'XOF',
  logo_url: '',
  banner_url: '',
})

const formatPrice = (amount: number, currency: string) => {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency }).format(amount)
}

const uploadLogo = async (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return
  try {
    const result = await shopApi.uploadMedia(file)
    form.logo_url = result.url
  } catch (e) {
    console.error('Failed to upload logo', e)
  }
}

const uploadBanner = async (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return
  try {
    const result = await shopApi.uploadMedia(file)
    form.banner_url = result.url
  } catch (e) {
    console.error('Failed to upload banner', e)
  }
}

const createShop = async () => {
  loading.value = true
  try {
    const tags = tagsInput.value.split(',').map(t => t.trim()).filter(Boolean)
    const shop = await shopApi.createShop({
      ...form,
      tags,
    })
    navigateTo(`/shops/${shop.slug}`)
  } catch (e: any) {
    alert(e.message || 'Échec de la création')
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    const result = await walletApi.getWallets()
    wallets.value = result.wallets || []
  } catch (e) {
    console.error('Failed to load wallets', e)
  }
})
</script>
