<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
    <!-- Backdrop -->
    <div class="fixed inset-0 bg-black/60 backdrop-blur-sm transition-opacity" @click="closeModal"></div>

    <div class="flex min-h-full items-center justify-center p-4">
      <div class="relative w-full max-w-4xl bg-white dark:bg-slate-900 rounded-2xl shadow-2xl overflow-hidden">
        <!-- Header -->
        <div class="bg-gradient-to-r from-indigo-600 to-purple-600 px-8 py-6">
          <div class="flex items-center justify-between">
            <div>
              <h2 class="text-2xl font-bold text-white">Facture #{{ invoice?.invoice_number }}</h2>
              <p class="text-indigo-100 mt-1">{{ formatDate(invoice?.created_at) }}</p>
            </div>
            <div class="flex items-center gap-3">
              <button @click="downloadPDF" class="p-2 bg-white/20 hover:bg-white/30 rounded-lg transition-colors" title="Télécharger PDF">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                </svg>
              </button>
              <button @click="printInvoice" class="p-2 bg-white/20 hover:bg-white/30 rounded-lg transition-colors" title="Imprimer">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z"/>
                </svg>
              </button>
              <button @click="closeModal" class="p-2 bg-white/20 hover:bg-white/30 rounded-lg transition-colors">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                </svg>
              </button>
            </div>
          </div>
        </div>

        <!-- Invoice Content -->
        <div id="invoice-content" class="p-8">
          <!-- Company & Client Info -->
          <div class="grid grid-cols-2 gap-8 mb-8">
            <div>
              <h3 class="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-3">De</h3>
              <div class="space-y-1">
                <p class="font-bold text-lg text-gray-900 dark:text-white">Zekora Bank</p>
                <p class="text-gray-600 dark:text-gray-400">123 Avenue de la Finance</p>
                <p class="text-gray-600 dark:text-gray-400">75001 Paris, France</p>
                <p class="text-gray-600 dark:text-gray-400">contact@zekora.com</p>
              </div>
            </div>
            <div class="text-right">
              <h3 class="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-3">Facturé à</h3>
              <div class="space-y-1">
                <p class="font-bold text-lg text-gray-900 dark:text-white">{{ userName }}</p>
                <p class="text-gray-600 dark:text-gray-400">{{ userEmail }}</p>
                <p class="text-gray-600 dark:text-gray-400">{{ userAddress }}</p>
              </div>
            </div>
          </div>

          <!-- Order Details -->
          <div class="bg-gray-50 dark:bg-slate-800/50 rounded-xl p-6 mb-8">
            <div class="grid grid-cols-4 gap-4 text-center">
              <div>
                <p class="text-sm text-gray-500 mb-1">Ordre ID</p>
                <p class="font-semibold text-gray-900 dark:text-white">{{ order?.id }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-500 mb-1">Date</p>
                <p class="font-semibold text-gray-900 dark:text-white">{{ formatDate(order?.created_at) }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-500 mb-1">Statut</p>
                <span :class="getStatusClass(order?.status)">{{ order?.status }}</span>
              </div>
              <div>
                <p class="text-sm text-gray-500 mb-1">Type</p>
                <p class="font-semibold text-gray-900 dark:text-white capitalize">{{ order?.order_type }}</p>
              </div>
            </div>
          </div>

          <!-- Items Table -->
          <table class="w-full mb-8">
            <thead>
              <tr class="border-b border-gray-200 dark:border-gray-700">
                <th class="text-left py-3 px-4 text-sm font-semibold text-gray-600 dark:text-gray-400">Description</th>
                <th class="text-center py-3 px-4 text-sm font-semibold text-gray-600 dark:text-gray-400">Type</th>
                <th class="text-center py-3 px-4 text-sm font-semibold text-gray-600 dark:text-gray-400">Quantité</th>
                <th class="text-center py-3 px-4 text-sm font-semibold text-gray-600 dark:text-gray-400">Prix</th>
                <th class="text-right py-3 px-4 text-sm font-semibold text-gray-600 dark:text-gray-400">Total</th>
              </tr>
            </thead>
            <tbody>
              <tr class="border-b border-gray-100 dark:border-gray-800">
                <td class="py-4 px-4">
                  <div class="font-medium text-gray-900 dark:text-white">{{ order?.pair }}</div>
                  <div class="text-sm text-gray-500">{{ order?.side?.toUpperCase() }} - {{ order?.order_type }}</div>
                </td>
                <td class="py-4 px-4 text-center">
                  <span :class="order?.side === 'buy' ? 'text-green-600 bg-green-100' : 'text-red-600 bg-red-100'" class="px-2 py-1 text-xs rounded-full">
                    {{ order?.side?.toUpperCase() }}
                  </span>
                </td>
                <td class="py-4 px-4 text-center text-gray-900 dark:text-white">{{ order?.amount }}</td>
                <td class="py-4 px-4 text-center text-gray-900 dark:text-white">{{ formatPrice(order?.price) }}</td>
                <td class="py-4 px-4 text-right font-semibold text-gray-900 dark:text-white">{{ calculateTotal() }}</td>
              </tr>
            </tbody>
          </table>

          <!-- Totals -->
          <div class="flex justify-end mb-8">
            <div class="w-72 space-y-3">
              <div class="flex justify-between text-sm">
                <span class="text-gray-600 dark:text-gray-400">Sous-total</span>
                <span class="text-gray-900 dark:text-white">{{ calculateTotal() }}</span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-gray-600 dark:text-gray-400">Frais (0.1%)</span>
                <span class="text-gray-900 dark:text-white">{{ calculateFees() }}</span>
              </div>
              <div class="flex justify-between text-lg font-bold border-t border-gray-200 dark:border-gray-700 pt-3">
                <span class="text-gray-900 dark:text-white">Total</span>
                <span class="text-indigo-600 dark:text-indigo-400">{{ calculateGrandTotal() }}</span>
              </div>
            </div>
          </div>

          <!-- Payment Info -->
          <div class="bg-indigo-50 dark:bg-indigo-900/20 rounded-xl p-6">
            <div class="flex items-center justify-between">
              <div>
                <h4 class="font-semibold text-indigo-900 dark:text-indigo-300 mb-1">Informations de paiement</h4>
                <p class="text-sm text-indigo-700 dark:text-indigo-400">Payé via Zekora Bank - {{ invoice?.payment_method || 'Carte bancaire' }}</p>
                <p class="text-sm text-indigo-700 dark:text-indigo-400">Transaction ID: {{ invoice?.transaction_id }}</p>
              </div>
              <div class="text-right">
                <div class="inline-flex items-center gap-2 px-4 py-2 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 rounded-full">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                  </svg>
                  <span class="font-semibold">Payée</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Footer Note -->
          <div class="mt-8 text-center">
            <p class="text-sm text-gray-500">Merci pour votre confiance ! Cette facture a été générée automatiquement par Zekora Bank.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from '~/stores/auth'

const props = defineProps({
  isOpen: Boolean,
  order: Object,
  invoice: Object
})

const emit = defineEmits(['close'])

const authStore = useAuthStore()

const userName = computed(() => {
  if (authStore.user) {
    return `${authStore.user.first_name || ''} ${authStore.user.last_name || ''}`
  }
  return 'Client'
})

const userEmail = computed(() => authStore.user?.email || '')
const userAddress = computed(() => authStore.user?.address || 'Non spécifié')

const closeModal = () => {
  emit('close')
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('fr-FR', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatPrice = (price) => {
  if (!price) return 'Market'
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'USD' }).format(price)
}

const getStatusClass = (status) => {
  const classes = {
    'pending': 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400',
    'filled': 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
    'cancelled': 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
    'partially_filled': 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400'
  }
  return `px-2 py-1 text-xs rounded-full ${classes[status] || classes.pending}`
}

const calculateTotal = () => {
  if (!props.order) return '$0.00'
  const price = props.order.price || 0
  const total = price * props.order.amount
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'USD' }).format(total)
}

const calculateFees = () => {
  if (!props.order || !props.order.price) return '$0.00'
  const total = props.order.price * props.order.amount
  const fees = total * 0.001
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'USD' }).format(fees)
}

const calculateGrandTotal = () => {
  if (!props.order) return '$0.00'
  const price = props.order.price || 0
  const total = price * props.order.amount
  const fees = total * 0.001
  const grandTotal = total + fees
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'USD' }).format(grandTotal)
}

const downloadPDF = () => {
  // Implement PDF download logic
  alert('Téléchargement PDF en cours...')
}

const printInvoice = () => {
  const printContent = document.getElementById('invoice-content')
  const originalContents = document.body.innerHTML
  document.body.innerHTML = printContent.innerHTML
  window.print()
  document.body.innerHTML = originalContents
  window.location.reload()
}
</script>

<style scoped>
@media print {
  #invoice-content {
    padding: 0;
  }
}
</style>
