<template>
  <ShopLayout>
    <div class="create-product-page">
      <!-- Page Header -->
      <div class="page-header">
        <NuxtLink :to="`/shops/manage/${slug}/products`" class="back-link">
          ← Retour aux produits
        </NuxtLink>
        <h1 class="page-title">📦 Nouveau Produit</h1>
        <p class="page-subtitle">Ajoutez un produit à votre catalogue</p>
      </div>

      <!-- Form -->
      <form @submit.prevent="submitProduct" class="product-form">
        
        <!-- Images Section -->
        <div class="form-section">
          <h3 class="section-title">📷 Images</h3>
          <div class="images-grid">
            <div 
              v-for="(image, idx) in images" 
              :key="idx"
              class="image-preview"
            >
              <img :src="image" alt="">
              <button type="button" @click="removeImage(idx)" class="remove-btn">✕</button>
            </div>
            <label v-if="images.length < 5" class="image-upload">
              <input type="file" accept="image/*" @change="handleImageUpload" hidden>
              <span class="upload-icon">+</span>
              <span class="upload-text">Ajouter</span>
            </label>
          </div>
          <p class="helper-text">Ajoutez jusqu'à 5 images (la première sera l'image principale)</p>
        </div>

        <!-- Basic Info -->
        <div class="form-section">
          <h3 class="section-title">📝 Informations</h3>
          
          <div class="form-group">
            <label>Nom du produit *</label>
            <input 
              v-model="form.name" 
              type="text" 
              required 
              placeholder="ex: T-shirt Vintage Premium"
              class="input"
            >
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>Prix *</label>
              <div class="price-input">
                <input 
                  v-model.number="form.price" 
                  type="number" 
                  min="0" 
                  step="1" 
                  required
                  class="input"
                  placeholder="0"
                >
                <span class="currency-badge">{{ form.currency }}</span>
              </div>
            </div>
            <div class="form-group">
              <label>Prix barré (optionnel)</label>
              <input 
                v-model.number="form.compare_at_price" 
                type="number" 
                min="0" 
                class="input"
                placeholder="0"
              >
            </div>
          </div>

          <div class="form-group">
            <label>Description</label>
            <textarea 
              v-model="form.description" 
              rows="4"
              placeholder="Décrivez votre produit de manière détaillée..."
              class="input"
            ></textarea>
          </div>

          <div class="form-group">
            <label>Catégorie</label>
            <select v-model="form.category_id" class="input">
              <option value="">Sélectionner une catégorie</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                {{ cat.icon }} {{ cat.name }}
              </option>
            </select>
          </div>
        </div>

        <!-- Inventory -->
        <div class="form-section">
          <h3 class="section-title">📊 Inventaire</h3>
          
          <div class="form-row">
            <div class="form-group">
              <label>Stock disponible</label>
              <input 
                v-model.number="form.stock" 
                type="number" 
                min="0"
                class="input"
              >
            </div>
            <div class="form-group">
              <label>SKU (référence)</label>
              <input 
                v-model="form.sku" 
                type="text"
                placeholder="ex: TSHIRT-001"
                class="input"
              >
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>Poids (kg)</label>
              <input 
                v-model.number="form.weight" 
                type="number" 
                min="0" 
                step="0.1"
                class="input"
                placeholder="0.0"
              >
            </div>
            <div class="form-group">
              <label>Statut</label>
              <select v-model="form.status" class="input">
                <option value="active">✅ Actif (visible)</option>
                <option value="draft">📝 Brouillon</option>
                <option value="out_of_stock">⚠️ Rupture</option>
              </select>
            </div>
          </div>
        </div>

        <!-- Options -->
        <div class="form-section">
          <h3 class="section-title">⚙️ Options</h3>
          
          <label class="checkbox-item">
            <input type="checkbox" v-model="form.is_featured">
            <span class="checkbox-label">⭐ Mettre en avant ce produit</span>
          </label>

          <label class="checkbox-item">
            <input type="checkbox" v-model="form.is_digital">
            <span class="checkbox-label">💻 Produit numérique (pas de livraison)</span>
          </label>
        </div>

        <!-- Actions -->
        <div class="form-actions">
          <NuxtLink :to="`/shops/manage/${slug}/products`" class="btn-secondary">
            Annuler
          </NuxtLink>
          <button type="submit" class="btn-primary" :disabled="submitting">
            <span v-if="submitting" class="spinner-small"></span>
            {{ submitting ? 'Création...' : '✓ Créer le produit' }}
          </button>
        </div>

      </form>
    </div>
  </ShopLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useShopApi, type Category } from '@/composables/useShopApi'
import ShopLayout from '@/components/shops/ShopLayout.vue'

const route = useRoute()
const router = useRouter()
const shopApi = useShopApi()
const slug = route.params.id as string
const shopId = ref('')

const images = ref<string[]>([])
const categories = ref<Category[]>([])

const form = ref({
  name: '',
  description: '',
  price: 0,
  compare_at_price: 0,
  currency: 'XOF',
  stock: 10,
  sku: '',
  weight: 0,
  is_featured: false,
  is_digital: false,
  status: 'active',
  category_id: ''
})

const submitting = ref(false)

const handleImageUpload = async (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return
  
  try {
    // Upload to server
    const result = await shopApi.uploadMedia(file)
    if (result.url) {
      images.value.push(result.url)
    }
  } catch (e) {
    // Fallback to local preview
    const reader = new FileReader()
    reader.onload = (e) => {
      if (e.target?.result) {
        images.value.push(e.target.result as string)
      }
    }
    reader.readAsDataURL(file)
  }
}

const removeImage = (idx: number) => {
  images.value.splice(idx, 1)
}

const submitProduct = async () => {
  if (!shopId.value) return
  
  try {
    submitting.value = true
    await shopApi.createProduct({
      shop_id: shopId.value,
      ...form.value,
      images: images.value,
      image_url: images.value[0] || ''
    })
    router.push(`/shops/manage/${slug}/products`)
  } catch (e: any) {
    alert('Erreur: ' + (e.message || 'Impossible de créer le produit'))
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  try {
    const shop = await shopApi.getShop(slug)
    shopId.value = shop.id
    form.value.currency = shop.currency || 'XOF'
    
    // Load categories
    const catResult = await shopApi.listCategories(slug)
    categories.value = catResult.categories || []
  } catch (e) {
    console.error('Failed to load shop info', e)
  }
})

definePageMeta({
  middleware: ['auth'],
  layout: 'dashboard'
})
</script>

<style scoped>
.create-product-page {
  padding: 24px;
  max-width: 800px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 32px;
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  color: var(--text-muted, #6b7280);
  font-size: 14px;
  margin-bottom: 8px;
  transition: color 0.2s;
}

.back-link:hover {
  color: #6366f1;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary, #1f2937);
  margin: 0;
}

.page-subtitle {
  color: var(--text-muted, #6b7280);
  margin-top: 4px;
}

/* Form */
.product-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-section {
  background: var(--surface, white);
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 16px;
  padding: 24px;
}

.section-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary, #1f2937);
  margin: 0 0 20px 0;
}

.form-group {
  margin-bottom: 16px;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary, #1f2937);
  margin-bottom: 6px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 12px;
  font-size: 15px;
  background: var(--surface, white);
  color: var(--text-primary, #1f2937);
  transition: border-color 0.2s, box-shadow 0.2s;
}

.input:focus {
  outline: none;
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.price-input {
  position: relative;
}

.price-input .input {
  padding-right: 60px;
}

.currency-badge {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;
}

/* Images */
.images-grid {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.image-preview {
  width: 100px;
  height: 100px;
  border-radius: 12px;
  overflow: hidden;
  position: relative;
}

.image-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.remove-btn {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 24px;
  height: 24px;
  background: rgba(0, 0, 0, 0.6);
  border: none;
  border-radius: 50%;
  color: white;
  cursor: pointer;
  font-size: 12px;
}

.image-upload {
  width: 100px;
  height: 100px;
  border: 2px dashed var(--border, #e5e7eb);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.image-upload:hover {
  border-color: #6366f1;
  background: rgba(99, 102, 241, 0.05);
}

.upload-icon {
  font-size: 24px;
  color: var(--text-muted, #9ca3af);
}

.upload-text {
  font-size: 12px;
  color: var(--text-muted, #9ca3af);
}

.helper-text {
  font-size: 13px;
  color: var(--text-muted, #6b7280);
  margin-top: 12px;
}

/* Checkboxes */
.checkbox-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--surface-hover, #f9fafb);
  border-radius: 12px;
  cursor: pointer;
  margin-bottom: 8px;
  transition: background 0.2s;
}

.checkbox-item:hover {
  background: rgba(99, 102, 241, 0.1);
}

.checkbox-item input[type="checkbox"] {
  width: 20px;
  height: 20px;
  accent-color: #6366f1;
}

.checkbox-label {
  font-size: 15px;
  color: var(--text-primary, #1f2937);
}

/* Actions */
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 24px;
}

.btn-primary {
  padding: 14px 28px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-weight: 700;
  font-size: 15px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.4);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.btn-secondary {
  padding: 14px 24px;
  background: var(--surface, white);
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 12px;
  font-weight: 600;
  color: var(--text-primary, #1f2937);
  cursor: pointer;
  transition: background 0.2s;
}

.btn-secondary:hover {
  background: var(--surface-hover, #f9fafb);
}

.spinner-small {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Dark Mode */
:global(.dark) .form-section {
  background: #1e293b;
  border-color: #334155;
}

:global(.dark) .checkbox-item {
  background: #334155;
}

:global(.dark) .image-upload {
  border-color: #475569;
}

@media (max-width: 640px) {
  .form-row {
    grid-template-columns: 1fr;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .btn-primary, .btn-secondary {
    width: 100%;
    justify-content: center;
  }
}
</style>
