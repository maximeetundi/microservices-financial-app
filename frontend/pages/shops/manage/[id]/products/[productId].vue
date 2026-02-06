<template>
  <div class="products-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">✏️ Modifier le produit</h1>
        <p class="page-subtitle">{{ product?.name || 'Chargement...' }}</p>
      </div>
      <NuxtLink :to="`/shops/manage/${shopSlug}/products`" class="btn-back">
        ← Retour aux produits
      </NuxtLink>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Chargement du produit...</p>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="error-state">
      <div class="error-icon">❌</div>
      <h3>Erreur</h3>
      <p>{{ error }}</p>
      <NuxtLink :to="`/shops/manage/${shopSlug}/products`" class="btn-primary">
        Retour aux produits
      </NuxtLink>
    </div>

    <!-- Product Form -->
    <div v-else class="form-container">
      <form @submit.prevent="saveProduct" class="product-form">
        <!-- Basic Info -->
        <div class="form-section">
          <h2 class="section-title">📋 Informations générales</h2>
          
          <div class="form-group">
            <label for="name">Nom du produit *</label>
            <input 
              id="name" 
              v-model="form.name" 
              type="text" 
              required 
              class="form-input"
              placeholder="Ex: iPhone 15 Pro"
            >
          </div>

          <div class="form-group">
            <label for="description">Description</label>
            <textarea 
              id="description" 
              v-model="form.description" 
              rows="4" 
              class="form-input"
              placeholder="Décrivez votre produit..."
            ></textarea>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="price">Prix *</label>
              <input 
                id="price" 
                v-model.number="form.price" 
                type="number" 
                min="0" 
                step="1" 
                required 
                class="form-input"
              >
            </div>
            
            <div class="form-group">
              <label for="compare_at_price">Prix barré</label>
              <input 
                id="compare_at_price" 
                v-model.number="form.compare_at_price" 
                type="number" 
                min="0" 
                step="1" 
                class="form-input"
                placeholder="Ancien prix"
              >
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="stock">Stock *</label>
              <input 
                id="stock" 
                v-model.number="form.stock" 
                type="number" 
                min="0" 
                required 
                class="form-input"
              >
            </div>
            
            <div class="form-group">
              <label for="sku">SKU / Référence</label>
              <input 
                id="sku" 
                v-model="form.sku" 
                type="text" 
                class="form-input"
                placeholder="REF-001"
              >
            </div>
          </div>

          <div class="form-group">
            <label for="status">Statut *</label>
            <select id="status" v-model="form.status" class="form-input">
              <option value="active">Actif</option>
              <option value="draft">Brouillon</option>
              <option value="archived">Archivé</option>
            </select>
          </div>

          <div class="form-group checkbox-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="form.is_featured">
              <span>⭐ Produit en vedette</span>
            </label>
          </div>
        </div>

        <!-- Images -->
        <div class="form-section">
          <h2 class="section-title">📷 Images</h2>
          
          <div class="images-grid">
            <div 
              v-for="(img, index) in form.images" 
              :key="index" 
              class="image-preview"
            >
              <img :src="img" alt="">
              <button type="button" @click="removeImage(index)" class="remove-image">×</button>
            </div>
            
            <div class="add-image">
              <input 
                type="url" 
                v-model="newImageUrl" 
                placeholder="URL de l'image" 
                class="form-input"
              >
              <button type="button" @click="addImage" class="btn-add-image">
                + Ajouter
              </button>
            </div>
          </div>
        </div>

        <!-- Submit -->
        <div class="form-actions">
          <button type="button" @click="$router.back()" class="btn-cancel">
            Annuler
          </button>
          <button type="submit" :disabled="saving" class="btn-submit">
            {{ saving ? 'Enregistrement...' : '💾 Enregistrer les modifications' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useShopApi, type Product } from '@/composables/useShopApi'

const route = useRoute()
const router = useRouter()
const shopApi = useShopApi()

const shopSlug = route.params.id as string
const productId = route.params.productId as string

const product = ref<Product | null>(null)
const loading = ref(true)
const saving = ref(false)
const error = ref<string | null>(null)
const newImageUrl = ref('')

const form = ref({
  name: '',
  description: '',
  price: 0,
  compare_at_price: 0,
  stock: 0,
  sku: '',
  status: 'active',
  is_featured: false,
  images: [] as string[]
})

const fetchProduct = async () => {
  try {
    loading.value = true
    const response = await shopApi.getProductById(productId)
    product.value = response
    
    // Fill form with product data
    form.value = {
      name: response.name,
      description: response.description || '',
      price: response.price,
      compare_at_price: response.compare_at_price || 0,
      stock: response.stock,
      sku: response.sku || '',
      status: response.status,
      is_featured: response.is_featured || false,
      images: [...(response.images || [])]
    }
  } catch (e: any) {
    console.error('Error fetching product:', e)
    error.value = e.message || 'Produit introuvable'
  } finally {
    loading.value = false
  }
}

const addImage = () => {
  if (newImageUrl.value && newImageUrl.value.startsWith('http')) {
    form.value.images.push(newImageUrl.value)
    newImageUrl.value = ''
  }
}

const removeImage = (index: number) => {
  form.value.images.splice(index, 1)
}

const saveProduct = async () => {
  try {
    saving.value = true
    
    await shopApi.updateProduct(productId, {
      name: form.value.name,
      description: form.value.description,
      price: form.value.price,
      compare_at_price: form.value.compare_at_price || undefined,
      stock: form.value.stock,
      sku: form.value.sku || undefined,
      status: form.value.status,
      is_featured: form.value.is_featured,
      images: form.value.images
    })
    
    router.push(`/shops/manage/${shopSlug}/products`)
  } catch (e: any) {
    alert('Erreur: ' + e.message)
  } finally {
    saving.value = false
  }
}

onMounted(fetchProduct)

definePageMeta({
  middleware: ['auth'],
  layout: 'shop-admin'
})
</script>

<style scoped>
.products-page {
  padding: 24px;
  width: 100%;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
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

.btn-back {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-muted, #6b7280);
  text-decoration: none;
  font-weight: 500;
  transition: color 0.2s;
}

.btn-back:hover {
  color: #6366f1;
}

/* Form */
.form-container {
  background: var(--surface, white);
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 16px;
  overflow: hidden;
}

.form-section {
  padding: 24px;
  border-bottom: 1px solid var(--border, #e5e7eb);
}

.form-section:last-of-type {
  border-bottom: none;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary, #1f2937);
  margin: 0 0 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary, #1f2937);
  margin-bottom: 8px;
}

.form-input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 10px;
  font-size: 15px;
  background: var(--surface, white);
  color: var(--text-primary, #1f2937);
  transition: border-color 0.2s, box-shadow 0.2s;
}

.form-input:focus {
  outline: none;
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.checkbox-group {
  margin-top: 8px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  font-weight: normal;
}

.checkbox-label input[type="checkbox"] {
  width: 20px;
  height: 20px;
  accent-color: #6366f1;
}

/* Images */
.images-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.image-preview {
  position: relative;
  width: 100px;
  height: 100px;
  border-radius: 12px;
  overflow: hidden;
  border: 2px solid var(--border, #e5e7eb);
}

.image-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.remove-image {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: rgba(239, 68, 68, 0.9);
  color: white;
  border: none;
  font-size: 18px;
  line-height: 1;
  cursor: pointer;
}

.add-image {
  display: flex;
  gap: 8px;
  flex: 1;
  min-width: 250px;
}

.btn-add-image {
  padding: 12px 20px;
  background: var(--surface-hover, #f3f4f6);
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 10px;
  font-weight: 500;
  cursor: pointer;
  white-space: nowrap;
}

/* Actions */
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px 24px;
  background: var(--surface-hover, #f9fafb);
}

.btn-cancel {
  padding: 12px 24px;
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 10px;
  background: white;
  color: var(--text-primary, #1f2937);
  font-weight: 500;
  cursor: pointer;
}

.btn-submit {
  padding: 12px 28px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-submit:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.4);
}

.btn-submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* States */
.loading-state, .error-state {
  text-align: center;
  padding: 60px 20px;
  background: var(--surface, white);
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 16px;
}

.error-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--border, #e5e7eb);
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

.btn-primary {
  display: inline-block;
  margin-top: 16px;
  padding: 12px 24px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  border-radius: 12px;
  font-weight: 600;
  text-decoration: none;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (max-width: 768px) {
  .form-row {
    grid-template-columns: 1fr;
  }
  
  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }
  
  .form-actions {
    flex-direction: column-reverse;
  }
  
  .btn-cancel, .btn-submit {
    width: 100%;
  }
}

/* Dark mode */
:root.dark .form-input,
:root.dark .btn-cancel,
:root.dark .btn-add-image {
  background: var(--surface, #1e293b);
  border-color: var(--border, #334155);
  color: var(--text-primary, #f1f5f9);
}
</style>
