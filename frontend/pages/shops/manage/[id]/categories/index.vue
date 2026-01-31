<template>
  <ShopLayout>
    <div class="categories-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">🏷️ Catégories</h1>
        <p class="page-subtitle">Organisez vos produits par rayons pour une meilleure navigation</p>
      </div>
      <button @click="openCreateModal" class="btn-create">
        <span class="icon">+</span>
        Nouvelle catégorie
      </button>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">📁</div>
        <div class="stat-info">
          <span class="stat-value">{{ categories.length }}</span>
          <span class="stat-label">Catégories</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">📦</div>
        <div class="stat-info">
          <span class="stat-value">{{ totalProducts }}</span>
          <span class="stat-label">Produits classés</span>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Chargement des catégories...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="categories.length === 0" class="empty-state">
      <div class="empty-icon">📂</div>
      <h3>Aucune catégorie</h3>
      <p>Créez des catégories pour mieux organiser vos produits</p>
      <button @click="openCreateModal" class="btn-primary">Créer une catégorie</button>
    </div>

    <!-- Categories Grid -->
    <div v-else class="categories-grid">
      <div 
        v-for="category in categories" 
        :key="category.id" 
        class="category-card"
      >
        <div class="category-icon" :style="{ background: getCategoryGradient(category) }">
          {{ category.icon || '📁' }}
        </div>
        <div class="category-content">
          <h3 class="category-name">{{ category.name }}</h3>
          <p class="category-description">{{ category.description || 'Aucune description' }}</p>
          <div class="category-stats">
            <span class="product-count">📦 {{ category.product_count || 0 }} produits</span>
          </div>
        </div>
        <div class="category-actions">
          <button @click="editCategory(category)" class="btn-edit" title="Modifier">
            ✏️
          </button>
          <button @click="deleteCategory(category)" class="btn-delete" title="Supprimer">
            🗑️
          </button>
        </div>
      </div>
    </div>

    <!-- Helper Card -->
    <div class="helper-card">
      <h3>💡 Conseils</h3>
      <ul>
        <li>Restez simple avec 5 à 8 catégories maximum</li>
        <li>Utilisez des noms clairs et descriptifs</li>
        <li>Ajoutez des émojis pour rendre les catégories plus visuelles</li>
      </ul>
    </div>

    <!-- Create/Edit Modal -->
    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal-content">
          <button class="close-btn" @click="closeModal">✕</button>
          <h2 class="modal-title">{{ isEditing ? 'Modifier la catégorie' : 'Nouvelle catégorie' }}</h2>
          
          <form @submit.prevent="saveCategory" class="modal-form">
            <!-- Emoji Selector -->
            <div class="form-group">
              <label>Icône</label>
              <div class="emoji-selector">
                <button 
                  v-for="emoji in availableEmojis" 
                  :key="emoji"
                  type="button"
                  class="emoji-btn"
                  :class="{ active: newCategory.icon === emoji }"
                  @click="newCategory.icon = emoji"
                >
                  {{ emoji }}
                </button>
              </div>
            </div>

            <div class="form-group">
              <label>Nom de la catégorie *</label>
              <input 
                v-model="newCategory.name" 
                type="text" 
                required
                placeholder="ex: Vêtements Femme"
                class="input"
              >
            </div>

            <div class="form-group">
              <label>Description</label>
              <textarea 
                v-model="newCategory.description" 
                rows="3"
                placeholder="Description optionnelle..."
                class="input"
              ></textarea>
            </div>

            <div class="modal-actions">
              <button type="button" @click="closeModal" class="btn-secondary">Annuler</button>
              <button type="submit" class="btn-primary" :disabled="submitting">
                <span v-if="submitting" class="spinner-small"></span>
                {{ isEditing ? 'Enregistrer' : 'Créer' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </div>
  </ShopLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useShopApi, type Category } from '@/composables/useShopApi'
import ShopLayout from '@/components/shops/ShopLayout.vue'

const route = useRoute()
const shopApi = useShopApi()
const slug = route.params.id as string

const categories = ref<Category[]>([])
const loading = ref(true)
const showModal = ref(false)
const submitting = ref(false)
const isEditing = ref(false)
const editingId = ref<string | null>(null)
const shopId = ref('')

const newCategory = ref({
  name: '',
  description: '',
  icon: '📁'
})

const availableEmojis = ['📁', '👕', '👗', '👠', '💍', '🎒', '📱', '💻', '🎮', '📚', '🍔', '🍕', '☕', '🏠', '🚗', '⚽', '🎨', '💄', '🧴', '🎁']

const totalProducts = computed(() => 
  categories.value.reduce((sum, cat) => sum + (cat.product_count || 0), 0)
)

const getCategoryGradient = (category: any) => {
  const gradients = [
    'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
    'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
    'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)'
  ]
  const index = category.name?.length % gradients.length || 0
  return gradients[index]
}

const fetchCategories = async () => {
  try {
    loading.value = true
    const res = await shopApi.listCategories(slug)
    categories.value = res.categories || []
    
    if (!shopId.value) {
      const shop = await shopApi.getShop(slug)
      shopId.value = shop.id
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  isEditing.value = false
  editingId.value = null
  newCategory.value = { name: '', description: '', icon: '📁' }
  showModal.value = true
}

const editCategory = (cat: Category) => {
  isEditing.value = true
  editingId.value = cat.id
  newCategory.value = {
    name: cat.name,
    description: cat.description || '',
    icon: cat.icon || '📁'
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const saveCategory = async () => {
  if (!shopId.value) return
  
  try {
    submitting.value = true
    
    if (isEditing.value && editingId.value) {
      await shopApi.updateCategory(editingId.value, newCategory.value)
    } else {
      await shopApi.createCategory({
        shop_id: shopId.value,
        ...newCategory.value
      })
    }
    
    closeModal()
    await fetchCategories()
  } catch (e: any) {
    alert('Erreur: ' + (e.message || 'Impossible de sauvegarder'))
  } finally {
    submitting.value = false
  }
}

const deleteCategory = async (cat: Category) => {
  if (!confirm(`Supprimer la catégorie "${cat.name}" ?`)) return
  try {
    await shopApi.deleteCategory(cat.id)
    await fetchCategories()
  } catch (e) {
    alert('Impossible de supprimer')
  }
}

onMounted(fetchCategories)

definePageMeta({
  middleware: ['auth'],
  layout: 'dashboard'
})
</script>

<style scoped>
.categories-page {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
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
  color: var(--text-primary);
  margin: 0;
}

.page-subtitle {
  color: var(--text-muted);
  margin-top: 4px;
}

.btn-create {
  display: flex;
  align-items: center;
  gap: 8px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  padding: 12px 24px;
  border-radius: 12px;
  border: none;
  cursor: pointer;
  font-weight: 600;
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-create:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.4);
}

/* Stats */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.stat-card {
  background: var(--surface, white);
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  font-size: 32px;
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(99, 102, 241, 0.1);
  border-radius: 12px;
}

.stat-value {
  display: block;
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary, #1f2937);
}

.stat-label {
  color: var(--text-muted, #6b7280);
  font-size: 14px;
}

/* Categories Grid */
.categories-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.category-card {
  background: var(--surface, white);
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 16px;
  padding: 20px;
  display: flex;
  gap: 16px;
  align-items: flex-start;
  transition: transform 0.2s, box-shadow 0.2s;
}

.category-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.1);
}

.category-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  flex-shrink: 0;
}

.category-content {
  flex: 1;
  min-width: 0;
}

.category-name {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary, #1f2937);
  margin: 0 0 4px 0;
}

.category-description {
  color: var(--text-muted, #6b7280);
  font-size: 14px;
  margin: 0 0 8px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.product-count {
  font-size: 13px;
  color: var(--text-muted, #6b7280);
}

.category-actions {
  display: flex;
  gap: 8px;
}

.btn-edit, .btn-delete {
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 16px;
  transition: background 0.2s;
  background: transparent;
}

.btn-edit:hover {
  background: rgba(99, 102, 241, 0.1);
}

.btn-delete:hover {
  background: rgba(239, 68, 68, 0.1);
}

/* Helper Card */
.helper-card {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.1) 0%, rgba(139, 92, 246, 0.1) 100%);
  border-radius: 16px;
  padding: 20px;
  margin-top: 24px;
}

.helper-card h3 {
  margin: 0 0 12px 0;
  color: var(--text-primary, #1f2937);
}

.helper-card ul {
  margin: 0;
  padding-left: 20px;
  color: var(--text-muted, #6b7280);
}

.helper-card li {
  margin-bottom: 8px;
}

/* Empty & Loading */
.empty-state, .loading-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-muted, #6b7280);
}

.empty-icon {
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

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 16px;
}

.modal-content {
  background: var(--surface, white);
  border-radius: 20px;
  width: 100%;
  max-width: 480px;
  padding: 32px;
  position: relative;
  animation: modalIn 0.3s ease-out;
}

@keyframes modalIn {
  from {
    opacity: 0;
    transform: scale(0.9) translateY(20px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: var(--text-muted, #6b7280);
}

.modal-title {
  font-size: 22px;
  font-weight: 700;
  margin: 0 0 24px 0;
  color: var(--text-primary, #1f2937);
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--text-primary, #1f2937);
}

.input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 12px;
  font-size: 15px;
  transition: border-color 0.2s, box-shadow 0.2s;
  background: var(--surface, white);
  color: var(--text-primary, #1f2937);
}

.input:focus {
  outline: none;
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.emoji-selector {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.emoji-btn {
  width: 44px;
  height: 44px;
  border: 2px solid var(--border, #e5e7eb);
  border-radius: 10px;
  background: var(--surface, white);
  cursor: pointer;
  font-size: 20px;
  transition: all 0.2s;
}

.emoji-btn:hover {
  background: rgba(99, 102, 241, 0.1);
}

.emoji-btn.active {
  border-color: #6366f1;
  background: rgba(99, 102, 241, 0.15);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

.btn-primary {
  padding: 12px 24px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-secondary {
  padding: 12px 24px;
  background: var(--surface, white);
  border: 1px solid var(--border, #e5e7eb);
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  color: var(--text-primary, #1f2937);
}

.spinner-small {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }
  
  .categories-grid {
    grid-template-columns: 1fr;
  }
}
</style>
