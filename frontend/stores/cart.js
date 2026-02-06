import { defineStore } from 'pinia'

export const useCartStore = defineStore('cart', {
  state: () => ({
    items: [],
    loading: false
  }),

  getters: {
    totalItems: (state) => {
      return state.items.reduce((total, item) => total + item.quantity, 0)
    },
    
    // Alias pour totalItems (pour compatibilité)
    itemCount: (state) => {
      return state.items.reduce((total, item) => total + item.quantity, 0)
    },
    
    total: (state) => {
      return state.items.reduce((total, item) => total + (item.price * item.quantity), 0)
    },
    
    // Obtenir les articles par shop
    itemsByShop: (state) => {
      const grouped = {}
      state.items.forEach(item => {
        if (!grouped[item.shop_id]) {
          grouped[item.shop_id] = {
            shop_id: item.shop_id,
            shop_name: item.shop_name,
            items: [],
            subtotal: 0
          }
        }
        grouped[item.shop_id].items.push(item)
        grouped[item.shop_id].subtotal += item.price * item.quantity
      })
      return Object.values(grouped)
    }
  },

  actions: {
    // Charger le panier depuis le localStorage
    loadCart() {
      try {
        const savedCart = localStorage.getItem('cart')
        if (savedCart) {
          this.items = JSON.parse(savedCart)
        }
      } catch (error) {
        console.error('Erreur lors du chargement du panier:', error)
        this.items = []
      }
    },

    // Alias pour loadCart (pour compatibilité)
    loadFromStorage() {
      this.loadCart()
    },

    // Sauvegarder le panier dans le localStorage
    saveCart() {
      try {
        localStorage.setItem('cart', JSON.stringify(this.items))
      } catch (error) {
        console.error('Erreur lors de la sauvegarde du panier:', error)
      }
    },

    // Ajouter un article au panier
    addToCart(product, quantity = 1, shopId, shopName) {
      const existingItem = this.items.find(item => 
        item.id === product.id && item.shop_id === shopId
      )

      if (existingItem) {
        existingItem.quantity += quantity
      } else {
        this.items.push({
          id: product.id,
          name: product.name,
          price: product.price,
          image: product.image,
          quantity: quantity,
          shop_id: shopId,
          shop_name: shopName,
          added_at: new Date().toISOString()
        })
      }

      this.saveCart()
      
      // Notification
      this.showNotification(`${product.name} ajouté au panier`)
    },

    // Alias pour addToCart (pour compatibilité avec ancien format)
    addItem(itemData) {
      const existingItem = this.items.find(item => 
        item.id === itemData.id && item.shop_id === itemData.shopId
      )

      if (existingItem) {
        existingItem.quantity += itemData.quantity
      } else {
        this.items.push({
          id: itemData.id,
          name: itemData.name,
          price: itemData.price,
          image: itemData.image,
          quantity: itemData.quantity,
          shop_id: itemData.shopId || 'unknown',
          shop_name: itemData.shopName || 'Boutique',
          added_at: new Date().toISOString()
        })
      }

      this.saveCart()
      this.showNotification(`${itemData.name} ajouté au panier`)
    },

    // Mettre à jour la quantité d'un article
    updateQuantity(itemId, newQuantity) {
      const item = this.items.find(item => item.id === itemId)
      if (item) {
        if (newQuantity <= 0) {
          this.removeItem(itemId)
        } else {
          item.quantity = newQuantity
          this.saveCart()
        }
      }
    },

    // Supprimer un article du panier
    removeItem(itemId) {
      const index = this.items.findIndex(item => item.id === itemId)
      if (index > -1) {
        const item = this.items[index]
        this.items.splice(index, 1)
        this.saveCart()
        this.showNotification(`${item.name} retiré du panier`)
      }
    },

    // Vider le panier
    clearCart() {
      this.items = []
      this.saveCart()
      this.showNotification('Panier vidé')
    },

    // Afficher une notification
    showNotification(message) {
      // Créer une notification temporaire
      const notification = document.createElement('div')
      notification.className = 'fixed top-4 right-4 bg-green-500 text-white px-4 py-2 rounded-lg shadow-lg z-50 animate-fade-in-up'
      notification.textContent = message
      document.body.appendChild(notification)
      
      setTimeout(() => {
        notification.remove()
      }, 3000)
    },

    // Synchroniser avec le backend (quand l'utilisateur est connecté)
    async syncWithBackend() {
      if (!this.isAuthenticated) return
      
      this.loading = true
      try {
        // Appel API pour synchroniser le panier
        // const response = await $fetch('/api/cart/sync', {
        //   method: 'POST',
        //   body: { items: this.items }
        // })
        // this.items = response.items
      } catch (error) {
        console.error('Erreur de synchronisation:', error)
      } finally {
        this.loading = false
      }
    }
  },

  persist: {
    enabled: true,
    strategies: [
      {
        key: 'cart',
        storage: localStorage
      }
    ]
  }
})
