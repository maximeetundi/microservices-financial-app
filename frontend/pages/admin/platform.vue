<template>
  <div class="platform-accounts-page">
    <!-- Header -->
    <div class="page-header">
      <h1>💰 Comptes Plateforme</h1>
      <p class="subtitle">Gestion des comptes fiat et crypto de l'entreprise</p>
    </div>

    <!-- Tabs -->
    <div class="tabs">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        :class="['tab', { active: activeTab === tab.id }]"
        @click="activeTab = tab.id"
      >
        <span class="icon">{{ tab.icon }}</span>
        {{ tab.label }}
      </button>
    </div>

    <!-- Fiat Accounts Tab -->
    <div v-if="activeTab === 'fiat'" class="tab-content">
      <div class="section-header">
        <h2>Comptes Fiat</h2>
        <button class="btn btn-primary" @click="showCreateFiatModal = true">
          + Nouveau Compte
        </button>
      </div>

      <div class="accounts-grid">
        <div
          v-for="account in fiatAccounts"
          :key="account.id"
          class="account-card"
        >
          <div class="account-header">
            <span class="currency-badge">{{ account.currency }}</span>
            <span class="type-badge">{{ account.account_type }}</span>
            <span class="priority-badge" :class="getPriorityClass(account.priority)">
              P{{ account.priority || 50 }}
            </span>
          </div>
          <h3>{{ account.name }}</h3>
          <p class="balance">{{ formatAmount(account.balance, account.currency) }}</p>
          <div class="balance-limits">
            <span v-if="account.min_balance > 0">Min: {{ formatAmount(account.min_balance, account.currency) }}</span>
            <span v-if="account.max_balance > 0">Max: {{ formatAmount(account.max_balance, account.currency) }}</span>
          </div>
          <p class="description">{{ account.description }}</p>
          <div class="account-actions">
            <button class="btn btn-success btn-sm" @click="openCreditModal(account)">
              ↑ Créditer
            </button>
            <button class="btn btn-danger btn-sm" @click="openDebitModal(account)">
              ↓ Débiter
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Crypto Wallets Tab -->
    <div v-if="activeTab === 'crypto'" class="tab-content">
      <div class="section-header">
        <h2>Portefeuilles Crypto</h2>
        <button class="btn btn-primary" @click="showCreateCryptoModal = true">
          + Nouvelle Adresse
        </button>
      </div>

      <div class="wallets-table">
        <table>
          <thead>
            <tr>
              <th>Crypto</th>
              <th>Réseau</th>
              <th>Type</th>
              <th>Priorité</th>
              <th>Label</th>
              <th>Adresse</th>
              <th>Balance</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="wallet in cryptoWallets" :key="wallet.id">
              <td><span class="currency-badge crypto">{{ wallet.currency }}</span></td>
              <td>{{ getNetworkLabel(wallet.network) }}</td>
              <td>{{ wallet.wallet_type }}</td>
              <td>
                <span class="priority-badge" :class="getPriorityClass(wallet.priority)">
                  P{{ wallet.priority || 50 }}
                </span>
              </td>
              <td>{{ wallet.label || '-' }}</td>
              <td class="address">
                <span class="truncated" v-if="wallet.address">{{ truncateAddress(wallet.address) }}</span>
                <span v-else class="no-address">Non configurée</span>
                <button v-if="wallet.address" class="copy-btn" @click="copyAddress(wallet.address)">📋</button>
              </td>
              <td>{{ wallet.balance }} {{ wallet.currency }}</td>
              <td>
                <button class="btn btn-sm" @click="syncWalletBalance(wallet.id)">🔄 Sync</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Transactions Tab -->
    <div v-if="activeTab === 'transactions'" class="tab-content">
      <div class="section-header">
        <h2>Journal Comptable</h2>
        <div class="reconciliation-badge">
          <span>Réconciliation:</span>
          <span v-for="(balance, currency) in reconciliation" :key="currency" class="balance-item">
            {{ currency }}: {{ formatAmount(balance, currency) }}
          </span>
        </div>
      </div>

      <div class="transactions-table">
        <table>
          <thead>
            <tr>
              <th>Date</th>
              <th>Opération</th>
              <th>Débit</th>
              <th>Crédit</th>
              <th>Montant</th>
              <th>Description</th>
              <th>Par</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="tx in transactions" :key="tx.id">
              <td>{{ formatDate(tx.transaction_date) }}</td>
              <td><span class="op-badge" :class="tx.operation_type">{{ tx.operation_type }}</span></td>
              <td>{{ tx.debit_account_type }}</td>
              <td>{{ tx.credit_account_type }}</td>
              <td class="amount">{{ formatAmount(tx.amount, tx.currency) }}</td>
              <td>{{ tx.description || '-' }}</td>
              <td>{{ tx.performed_by || 'system' }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create Fiat Account Modal -->
    <div v-if="showCreateFiatModal" class="modal-overlay" @click.self="showCreateFiatModal = false">
      <div class="modal">
        <h3>Créer un Compte Fiat</h3>
        <form @submit.prevent="createFiatAccount">
          <div class="form-group">
            <label>Devise</label>
            <select v-model="newFiatAccount.currency" required>
              <option value="FCFA">FCFA</option>
              <option value="XOF">XOF</option>
              <option value="EUR">EUR</option>
              <option value="USD">USD</option>
            </select>
          </div>
          <div class="form-group">
            <label>Type</label>
            <select v-model="newFiatAccount.account_type" required>
              <option value="reserve">Réserve</option>
              <option value="fees">Frais</option>
              <option value="operations">Opérations</option>
              <option value="pending">En attente</option>
            </select>
          </div>
          <div class="form-group">
            <label>Nom</label>
            <input v-model="newFiatAccount.name" type="text" required placeholder="Réserve principale" />
          </div>
          <div class="form-group">
            <label>Description</label>
            <textarea v-model="newFiatAccount.description" placeholder="Description du compte"></textarea>
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>Priorité (1-100)</label>
              <input v-model.number="newFiatAccount.priority" type="number" min="1" max="100" placeholder="50" />
              <small>Plus élevé = sélectionné en premier</small>
            </div>
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>Balance minimum</label>
              <input v-model.number="newFiatAccount.min_balance" type="number" min="0" step="1000" placeholder="0" />
            </div>
            <div class="form-group">
              <label>Balance maximum (0 = illimité)</label>
              <input v-model.number="newFiatAccount.max_balance" type="number" min="0" step="10000" placeholder="0" />
            </div>
          </div>
          <div class="modal-actions">
            <button type="button" class="btn btn-secondary" @click="showCreateFiatModal = false">Annuler</button>
            <button type="submit" class="btn btn-primary">Créer</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Create Crypto Wallet Modal -->
    <div v-if="showCreateCryptoModal" class="modal-overlay" @click.self="showCreateCryptoModal = false">
      <div class="modal">
        <h3>Ajouter une Adresse Crypto</h3>
        <form @submit.prevent="createCryptoWallet">
          <div class="form-group">
            <label>Crypto</label>
            <select v-model="newCryptoWallet.currency" required>
              <option value="BTC">Bitcoin (BTC)</option>
              <option value="ETH">Ethereum (ETH)</option>
              <option value="USDT">Tether (USDT)</option>
              <option value="USDC">USD Coin (USDC)</option>
              <option value="BNB">Binance Coin (BNB)</option>
              <option value="SOL">Solana (SOL)</option>
            </select>
          </div>
          <div class="form-group">
            <label>Réseau</label>
            <select v-model="newCryptoWallet.network" required>
              <option value="bitcoin">Bitcoin</option>
              <option value="ethereum">Ethereum</option>
              <option value="bsc">Binance Smart Chain</option>
              <option value="tron">Tron</option>
              <option value="polygon">Polygon</option>
              <option value="solana">Solana</option>
            </select>
          </div>
          <div class="form-group">
            <label>Type</label>
            <select v-model="newCryptoWallet.wallet_type" required>
              <option value="hot">Hot Wallet</option>
              <option value="cold">Cold Wallet</option>
              <option value="reserve">Réserve</option>
            </select>
          </div>
          <div class="form-group">
            <label>Adresse</label>
            <input v-model="newCryptoWallet.address" type="text" required placeholder="0x..." />
          </div>
          <div class="form-group">
            <label>Label</label>
            <input v-model="newCryptoWallet.label" type="text" placeholder="ETH Hot Wallet 1" />
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>Priorité (1-100)</label>
              <input v-model.number="newCryptoWallet.priority" type="number" min="1" max="100" placeholder="50" />
              <small>Plus élevé = sélectionné en premier</small>
            </div>
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>Balance minimum</label>
              <input v-model.number="newCryptoWallet.min_balance" type="number" min="0" step="0.0001" placeholder="0" />
            </div>
            <div class="form-group">
              <label>Balance maximum (0 = illimité)</label>
              <input v-model.number="newCryptoWallet.max_balance" type="number" min="0" step="0.001" placeholder="0" />
            </div>
          </div>
          <div class="modal-actions">
            <button type="button" class="btn btn-secondary" @click="showCreateCryptoModal = false">Annuler</button>
            <button type="submit" class="btn btn-primary">Ajouter</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Credit/Debit Modal -->
    <div v-if="showCreditDebitModal" class="modal-overlay" @click.self="showCreditDebitModal = false">
      <div class="modal">
        <h3>{{ creditDebitMode === 'credit' ? 'Créditer' : 'Débiter' }} le compte</h3>
        <p class="modal-subtitle">{{ selectedAccount?.name }} ({{ selectedAccount?.currency }})</p>
        <form @submit.prevent="executeCreditDebit">
          <div class="form-group">
            <label>Montant</label>
            <input v-model.number="creditDebitAmount" type="number" step="0.01" min="0.01" required />
          </div>
          <div class="form-group">
            <label>Description</label>
            <input v-model="creditDebitDescription" type="text" required placeholder="Raison de l'opération" />
          </div>
          <div class="form-group">
            <label>Référence externe</label>
            <input v-model="creditDebitReference" type="text" placeholder="N° virement, etc." />
          </div>
          <div class="modal-actions">
            <button type="button" class="btn btn-secondary" @click="showCreditDebitModal = false">Annuler</button>
            <button type="submit" class="btn" :class="creditDebitMode === 'credit' ? 'btn-success' : 'btn-danger'">
              {{ creditDebitMode === 'credit' ? 'Créditer' : 'Débiter' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { adminPlatformAPI } from '~/composables/useApi'

definePageMeta({
  layout: 'admin'
})

const tabs = [
  { id: 'fiat', label: 'Comptes Fiat', icon: '💵' },
  { id: 'crypto', label: 'Portefeuilles Crypto', icon: '₿' },
  { id: 'transactions', label: 'Journal', icon: '📒' },
]

const activeTab = ref('fiat')

// Data
const fiatAccounts = ref<any[]>([])
const cryptoWallets = ref<any[]>([])
const transactions = ref<any[]>([])
const reconciliation = ref<Record<string, number>>({})

// Modals
const showCreateFiatModal = ref(false)
const showCreateCryptoModal = ref(false)
const showCreditDebitModal = ref(false)
const creditDebitMode = ref<'credit' | 'debit'>('credit')
const selectedAccount = ref<any>(null)
const creditDebitAmount = ref(0)
const creditDebitDescription = ref('')
const creditDebitReference = ref('')

// New account forms
const newFiatAccount = ref({
  currency: 'FCFA',
  account_type: 'reserve',
  name: '',
  description: '',
  priority: 50,
  min_balance: 0,
  max_balance: 0
})

const newCryptoWallet = ref({
  currency: 'ETH',
  network: 'ethereum',
  wallet_type: 'hot',
  address: '',
  label: '',
  priority: 50,
  min_balance: 0,
  max_balance: 0
})

// Fetch data
const fetchAccounts = async () => {
  try {
    const res = await adminPlatformAPI.getAccounts()
    fiatAccounts.value = res.data.accounts || []
  } catch (e) {
    console.error('Error fetching accounts:', e)
  }
}

const fetchCryptoWallets = async () => {
  try {
    const res = await adminPlatformAPI.getCryptoWallets()
    cryptoWallets.value = res.data.wallets || []
  } catch (e) {
    console.error('Error fetching crypto wallets:', e)
  }
}

const fetchTransactions = async () => {
  try {
    const res = await adminPlatformAPI.getTransactions()
    transactions.value = res.data.transactions || []
  } catch (e) {
    console.error('Error fetching transactions:', e)
  }
}

const fetchReconciliation = async () => {
  try {
    const res = await adminPlatformAPI.getReconciliation()
    reconciliation.value = res.data.balances || {}
  } catch (e) {
    console.error('Error fetching reconciliation:', e)
  }
}

// Actions
const createFiatAccount = async () => {
  try {
    await adminPlatformAPI.createAccount(newFiatAccount.value)
    showCreateFiatModal.value = false
    newFiatAccount.value = { currency: 'FCFA', account_type: 'reserve', name: '', description: '', priority: 50, min_balance: 0, max_balance: 0 }
    await fetchAccounts()
  } catch (e) {
    console.error('Error creating account:', e)
    alert('Erreur lors de la création du compte')
  }
}

const createCryptoWallet = async () => {
  try {
    await adminPlatformAPI.createCryptoWallet(newCryptoWallet.value)
    showCreateCryptoModal.value = false
    newCryptoWallet.value = { currency: 'ETH', network: 'ethereum', wallet_type: 'hot', address: '', label: '', priority: 50, min_balance: 0, max_balance: 0 }
    await fetchCryptoWallets()
  } catch (e) {
    console.error('Error creating crypto wallet:', e)
    alert('Erreur lors de l\'ajout du wallet')
  }
}

const openCreditModal = (account: any) => {
  selectedAccount.value = account
  creditDebitMode.value = 'credit'
  creditDebitAmount.value = 0
  creditDebitDescription.value = ''
  creditDebitReference.value = ''
  showCreditDebitModal.value = true
}

const openDebitModal = (account: any) => {
  selectedAccount.value = account
  creditDebitMode.value = 'debit'
  creditDebitAmount.value = 0
  creditDebitDescription.value = ''
  creditDebitReference.value = ''
  showCreditDebitModal.value = true
}

const executeCreditDebit = async () => {
  if (!selectedAccount.value) return
  try {
    const data = {
      amount: creditDebitAmount.value,
      description: creditDebitDescription.value,
      reference: creditDebitReference.value
    }
    if (creditDebitMode.value === 'credit') {
      await adminPlatformAPI.creditAccount(selectedAccount.value.id, data)
    } else {
      await adminPlatformAPI.debitAccount(selectedAccount.value.id, data)
    }
    showCreditDebitModal.value = false
    await fetchAccounts()
    await fetchTransactions()
    await fetchReconciliation()
  } catch (e: any) {
    console.error('Error executing credit/debit:', e)
    alert(e.response?.data?.error || 'Erreur lors de l\'opération')
  }
}

const syncWalletBalance = async (walletId: string) => {
  try {
    await adminPlatformAPI.syncCryptoWalletBalance(walletId)
    await fetchCryptoWallets()
  } catch (e) {
    console.error('Error syncing wallet:', e)
  }
}

// Helpers
const formatAmount = (amount: number, currency: string) => {
  return new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: currency === 'FCFA' || currency === 'XOF' ? 'XOF' : currency
  }).format(amount)
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleString('fr-FR')
}

const getNetworkLabel = (network: string) => {
  const labels: Record<string, string> = {
    ethereum: 'Ethereum',
    bsc: 'BSC',
    tron: 'Tron',
    bitcoin: 'Bitcoin',
    polygon: 'Polygon',
    solana: 'Solana'
  }
  return labels[network] || network
}

const truncateAddress = (address: string) => {
  if (!address) return ''
  return `${address.slice(0, 8)}...${address.slice(-6)}`
}

const copyAddress = async (address: string) => {
  await navigator.clipboard.writeText(address)
  alert('Adresse copiée!')
}

const getPriorityClass = (priority: number) => {
  if (priority >= 80) return 'high'
  if (priority >= 50) return 'medium'
  return 'low'
}

onMounted(() => {
  fetchAccounts()
  fetchCryptoWallets()
  fetchTransactions()
  fetchReconciliation()
})
</script>

<style scoped>
.platform-accounts-page {
  padding: 2rem;
  background: var(--bg-primary, #0a0a0f);
  min-height: 100vh;
  color: var(--text-primary, #fff);
}

.page-header {
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.subtitle {
  color: var(--text-secondary, #888);
}

.tabs {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  border-bottom: 1px solid var(--border-color, #333);
  padding-bottom: 1rem;
}

.tab {
  background: transparent;
  border: none;
  color: var(--text-secondary, #888);
  cursor: pointer;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.tab:hover, .tab.active {
  background: var(--bg-secondary, #1a1a2e);
  color: var(--text-primary, #fff);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.accounts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.account-card {
  background: var(--bg-secondary, #1a1a2e);
  border-radius: 12px;
  padding: 1.5rem;
  border: 1px solid var(--border-color, #333);
}

.account-header {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.currency-badge {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-size: 0.875rem;
  font-weight: 600;
}

.currency-badge.crypto {
  background: linear-gradient(135deg, #f59e0b, #d97706);
}

.type-badge {
  background: var(--bg-tertiary, #252540);
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-size: 0.75rem;
}

.account-card h3 {
  margin-bottom: 0.5rem;
}

.balance {
  font-size: 1.5rem;
  font-weight: 700;
  color: #10b981;
  margin-bottom: 0.5rem;
}

.description {
  color: var(--text-secondary, #888);
  font-size: 0.875rem;
  margin-bottom: 1rem;
}

.account-actions {
  display: flex;
  gap: 0.5rem;
}

/* Table styles */
.wallets-table, .transactions-table {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  background: var(--bg-secondary, #1a1a2e);
  border-radius: 12px;
  overflow: hidden;
}

th, td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid var(--border-color, #333);
}

th {
  background: var(--bg-tertiary, #252540);
  font-weight: 600;
}

.address {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.truncated {
  font-family: monospace;
  background: var(--bg-tertiary, #252540);
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
}

.copy-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1rem;
}

.op-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
}

.op-badge.admin_credit { background: #10b981; }
.op-badge.admin_debit { background: #ef4444; }
.op-badge.exchange { background: #6366f1; }
.op-badge.fee { background: #f59e0b; }

.reconciliation-badge {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.balance-item {
  background: var(--bg-tertiary, #252540);
  padding: 0.5rem 1rem;
  border-radius: 8px;
}

/* Buttons */
.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.2s;
}

.btn-sm {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
}

.btn-primary {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: white;
}

.btn-success {
  background: #10b981;
  color: white;
}

.btn-danger {
  background: #ef4444;
  color: white;
}

.btn-secondary {
  background: var(--bg-tertiary, #252540);
  color: var(--text-primary, #fff);
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: var(--bg-secondary, #1a1a2e);
  border-radius: 16px;
  padding: 2rem;
  width: 100%;
  max-width: 500px;
  border: 1px solid var(--border-color, #333);
}

.modal h3 {
  margin-bottom: 0.5rem;
}

.modal-subtitle {
  color: var(--text-secondary, #888);
  margin-bottom: 1.5rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--border-color, #333);
  border-radius: 8px;
  background: var(--bg-tertiary, #252540);
  color: var(--text-primary, #fff);
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 1.5rem;
}

/* Priority badge styles */
.priority-badge {
  padding: 0.2rem 0.5rem;
  border-radius: 12px;
  font-size: 0.7rem;
  font-weight: 600;
}

.priority-badge.high {
  background: #10b981;
  color: white;
}

.priority-badge.medium {
  background: #6366f1;
  color: white;
}

.priority-badge.low {
  background: #6b7280;
  color: white;
}

/* Balance limits */
.balance-limits {
  display: flex;
  gap: 1rem;
  font-size: 0.75rem;
  color: var(--text-secondary, #888);
  margin-bottom: 0.5rem;
}

.balance-limits span {
  background: var(--bg-tertiary, #252540);
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
}

/* Form row for side-by-side inputs */
.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-group small {
  display: block;
  color: var(--text-secondary, #888);
  font-size: 0.75rem;
  margin-top: 0.25rem;
}

.no-address {
  color: #ef4444;
  font-style: italic;
  font-size: 0.875rem;
}
</style>
