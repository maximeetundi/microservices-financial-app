<template>
  <NuxtLayout name="dashboard">
    <div class="max-w-7xl mx-auto animate-fade-in-up">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
        <div>
          <h1 class="text-3xl font-extrabold text-gray-900 dark:text-white mb-2">My Wallets 👛</h1>
          <p class="text-gray-500 dark:text-gray-400">Manage your fiat and crypto currencies</p>
        </div>
        <div class="flex gap-3">
           <NuxtLink to="/cards" class="btn-secondary flex items-center gap-2">
            <span class="text-xl">💳</span>
            <span>Commander une carte</span>
          </NuxtLink>
          <button @click="showCreateWallet = true" class="btn-primary flex items-center gap-2 shadow-lg shadow-indigo-500/30">
            <span class="text-xl">+</span>
            <span>Nouveau Portefeuille</span>
          </button>
        </div>
      </div>

      <!-- Total Balance Card -->
      <div class="glass-card mb-8 p-8 dark:bg-slate-900/80 border border-gray-200 dark:border-white/10 relative overflow-hidden group">
        <!-- Background Effects -->
        <div class="absolute top-0 right-0 w-64 h-64 bg-indigo-500/10 rounded-full blur-3xl group-hover:bg-indigo-500/20 transition-all duration-500"></div>
        
        <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-6 relative z-10">
          <div>
            <p class="text-gray-500 dark:text-gray-400 font-medium mb-1 uppercase tracking-wider text-sm">Valeur Totale</p>
            <div class="flex items-baseline gap-3">
              <div v-if="loading" class="h-14 w-64 bg-gray-200 dark:bg-slate-800 rounded-xl animate-pulse my-1"></div>
              <h2 v-else class="text-5xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-gray-900 to-gray-700 dark:from-white dark:to-gray-300">
                {{ formatMoney(totalBalance) }}
              </h2>
            </div>
            <p class="text-sm mt-3 flex items-center gap-2 text-emerald-600 dark:text-emerald-400 font-medium bg-emerald-50 dark:bg-emerald-500/10 px-3 py-1 rounded-full w-fit">
              <span class="relative flex h-2 w-2">
                <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
              </span>
              Calculé en temps réel
            </p>
          </div>
          <div class="flex flex-wrap gap-3">
             <!-- Recharger / Deposit -->
            <button @click="openTopUpModal" class="flex items-center gap-2 px-6 py-3 rounded-xl bg-indigo-600 text-white font-bold hover:bg-indigo-700 shadow-lg shadow-indigo-500/30 transition-all transform hover:-translate-y-0.5">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
              Recharger
            </button>
            <!-- Envoyer / Send -->
            <NuxtLink to="/transfer" class="flex items-center gap-2 px-6 py-3 rounded-xl bg-white/50 dark:bg-slate-800 text-gray-700 dark:text-white border border-gray-200 dark:border-gray-700 font-bold hover:bg-gray-50 dark:hover:bg-slate-700 transition-all transform hover:-translate-y-0.5">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"/></svg>
              Envoyer
            </NuxtLink>
          </div>
        </div>
      </div>

      <!-- Wallets Grid -->
      <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
        <div v-for="i in 3" :key="i" class="glass-card p-6 h-64 animate-pulse bg-white/50 dark:bg-slate-900/50 border border-gray-200 dark:border-white/5 relative overflow-hidden">
           <div class="flex items-center justify-between mb-6">
              <div class="flex items-center gap-4">
                 <div class="w-12 h-12 rounded-xl bg-gray-200 dark:bg-slate-800"></div>
                 <div class="space-y-2">
                    <div class="h-4 w-20 bg-gray-200 dark:bg-slate-800 rounded"></div>
                    <div class="h-3 w-32 bg-gray-200 dark:bg-slate-800 rounded"></div>
                 </div>
              </div>
           </div>
           <div class="h-8 w-40 bg-gray-200 dark:bg-slate-800 rounded mb-2"></div>
           <div class="h-4 w-24 bg-gray-200 dark:bg-slate-800 rounded"></div>
        </div>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
        <div v-for="wallet in wallets" :key="wallet.id" 
            class="glass-card p-6 cursor-pointer hover:border-indigo-500 dark:hover:border-indigo-500 transition-all duration-300 relative overflow-hidden group dark:bg-slate-900 border border-gray-200 dark:border-white/10"
            :class="{'ring-2 ring-indigo-500 dark:ring-indigo-400': selectedWallet?.id === wallet.id}"
            @click="selectWallet(wallet)">
            
           <!-- Background Glow -->
           <div class="absolute -right-10 -top-10 w-32 h-32 rounded-full blur-3xl opacity-0 group-hover:opacity-10 transition-opacity duration-500"
                :class="getCurrencyBg(wallet.currency)"></div>

          <div class="flex items-center justify-between mb-6 relative z-10">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-xl flex items-center justify-center text-2xl shadow-md text-white font-bold" 
                   :class="getCurrencyBg(wallet.currency)">
                {{ getCurrencyIcon(wallet.currency) }}
              </div>
              <div>
                <p class="font-bold text-gray-900 dark:text-white text-lg">{{ wallet.currency }}</p>
                <p class="text-xs text-gray-500 dark:text-gray-400 font-medium uppercase tracking-wider">{{ wallet.name }}</p>
              </div>
            </div>
            <div class="flex flex-col items-end gap-1">
               <span class="px-2.5 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wide border" 
                  :class="wallet.status === 'active' ? 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-200 dark:border-emerald-500/20' : 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-200 dark:border-amber-500/20'">
              {{ wallet.status }}
            </span>
             <span class="text-[10px] font-mono text-gray-400">{{ wallet.type }}</span>
            </div>
          </div>

          <div class="mb-4 relative z-10">
            <p class="text-2xl font-extrabold text-gray-900 dark:text-white tracking-tight">
              {{ formatCrypto(wallet.balance, wallet.currency) }}
            </p>
            <p class="text-sm text-gray-500 dark:text-gray-400 font-medium">
              ≈ {{ formatMoney(wallet.balanceUSD) }}
            </p>
          </div>

          <div class="flex justify-between items-center text-sm relative z-10 pt-4 border-t border-gray-100 dark:border-gray-800">
            <span class="text-gray-400 font-medium">Adresse</span>
            <div class="flex items-center gap-2">
                 <span class="text-gray-600 dark:text-gray-300 font-mono text-xs bg-gray-100 dark:bg-slate-800 px-2 py-1 rounded">{{ truncateAddress(wallet.wallet_address) }}</span>
            </div>
          </div>
          
          <!-- Action Buttons per Wallet -->
          <div class="flex gap-2 mt-4 pt-4 border-t border-gray-100 dark:border-gray-800 relative z-10">
            <button @click.stop="openTopUpForWallet(wallet)" 
                    class="flex-1 flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl bg-indigo-600 text-white text-sm font-semibold hover:bg-indigo-700 transition-all">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
              Recharger
            </button>
            <NuxtLink :to="`/transfer?wallet=${wallet.id}${wallet.wallet_type === 'crypto' ? '&type=crypto' : ''}`" @click.stop
                    class="flex-1 flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl bg-white dark:bg-slate-800 text-gray-700 dark:text-white text-sm font-semibold border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-slate-700 transition-all">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"/></svg>
              Envoyer
            </NuxtLink>

          </div>
        </div>
      </div>

      <!-- Selected Wallet Actions (If any layout specific needed, handled by modal) -->

      <!-- Create Wallet Modal -->
      <div v-if="showCreateWallet" class="fixed inset-0 bg-black/60 backdrop-blur-md z-50 flex items-center justify-center p-4 animate-fade-in-up">
        <div class="bg-white dark:bg-slate-900 rounded-2xl p-8 max-w-md w-full shadow-2xl border border-gray-100 dark:border-gray-800">
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-xl font-bold text-gray-900 dark:text-white">Nouveau Portefeuille</h3>
            <button @click="showCreateWallet = false" class="text-gray-400 hover:text-gray-600 dark:hover:text-white transition-colors">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>

          <form @submit.prevent="createWallet" class="space-y-5">
            <div>
              <label class="block text-sm font-bold text-gray-700 dark:text-gray-300 mb-2">Type de devise</label>
              <select v-model="newWallet.type" class="input-premium w-full p-3 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500 outline-none transition-all">
                <option value="fiat">Devise Fiat (USD, EUR...)</option>
                <option value="crypto">Crypto-monnaie (BTC, ETH...)</option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-bold text-gray-700 dark:text-gray-300 mb-2">Devise</label>
              <select v-model="newWallet.currency" class="input-premium w-full p-3 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500 outline-none transition-all">
                <template v-if="newWallet.type === 'fiat'">
                  <option value="USD">🇺🇸 Dollar US (USD)</option>
                  <option value="EUR">🇪🇺 Euro (EUR)</option>
                  <option value="GBP">🇬🇧 Livre Sterling (GBP)</option>
                  <option value="XOF">🇨🇮 Franc CFA (XOF)</option>
                  <option value="XAF">🇨🇲 Franc CFA (XAF)</option>
                </template>
                <template v-else>
                  <option value="BTC">₿ Bitcoin (BTC)</option>
                  <option value="ETH">Ξ Ethereum (ETH)</option>
                  <option value="USDT">💵 Tether (USDT)</option>
                  <option value="SOL">◎ Solana (SOL)</option>
                </template>
              </select>
            </div>

            <div>
              <label class="block text-sm font-bold text-gray-700 dark:text-gray-300 mb-2">Nom du portefeuille</label>
              <input 
                v-model="newWallet.name" 
                type="text" 
                placeholder="ex: Mon portefeuille principal"
                class="input-premium w-full p-3 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800 focus:ring-2 focus:ring-indigo-500 outline-none transition-all"
              />
            </div>

            <div class="flex gap-3 mt-8">
              <button 
                type="button" 
                @click="showCreateWallet = false"
                class="flex-1 py-3 px-4 rounded-xl font-bold text-gray-600 dark:text-gray-300 bg-gray-100 dark:bg-slate-800 hover:bg-gray-200 dark:hover:bg-slate-700 transition-colors"
              >
                Annuler
              </button>
              <button 
                type="submit" 
                :disabled="creatingWallet"
                class="flex-1 py-3 px-4 rounded-xl font-bold text-white bg-indigo-600 hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-500/30"
              >
                {{ creatingWallet ? 'Création...' : 'Créer le portefeuille' }}
              </button>
            </div>
          </form>
        </div>
      </div>

       <!-- Recharge / Top Up Modal -->
      <div v-if="showTopUpModal" class="fixed inset-0 bg-black/60 backdrop-blur-md z-50 flex items-center justify-center p-4 animate-fade-in-up">
        <div class="bg-white dark:bg-slate-900 rounded-2xl p-0 max-w-lg w-full shadow-2xl border border-gray-100 dark:border-gray-800 overflow-hidden flex flex-col max-h-[90vh]">
           <div class="p-6 border-b border-gray-100 dark:border-gray-800 flex justify-between items-center">
             <h3 class="text-xl font-bold text-gray-900 dark:text-white">
               {{ selectedWallet?.wallet_type === 'crypto' || selectedWallet?.type === 'crypto' ? 'Recevoir Crypto' : 'Recharger Compte' }}
             </h3>
             <button @click="closeTopUpModal" class="text-gray-400 hover:text-gray-900 dark:hover:text-white">✕</button>
           </div>
           
           <div class="p-6 overflow-y-auto custom-scrollbar">
             <!-- Wallet Selection Display -->
             <div class="mb-6 p-4 bg-gray-50 dark:bg-slate-800 rounded-xl flex items-center gap-4">
               <div class="w-10 h-10 rounded-full flex items-center justify-center text-xl font-bold text-white shadow-sm" :class="getCurrencyBg(selectedWallet?.currency)">
                 {{ getCurrencyIcon(selectedWallet?.currency) }}
               </div>
               <div>
                  <p class="text-sm text-gray-500 dark:text-gray-400">Portefeuille cible</p>
                  <p class="font-bold text-gray-900 dark:text-white">{{ selectedWallet?.name }} ({{ selectedWallet?.currency }})</p>
               </div>
             </div>

             <!-- CRYPTO DEPOSIT VIEW -->
             <div v-if="selectedWallet?.wallet_type === 'crypto' || selectedWallet?.type === 'crypto'" class="text-center space-y-6">
                <div class="p-6 bg-white rounded-xl shadow-inner border border-gray-100 inline-block">
                  <img 
                    :src="`https://api.qrserver.com/v1/create-qr-code/?size=160x160&data=${selectedWallet?.wallet_address}`" 
                    alt="Wallet QR Code" 
                    class="w-40 h-40 mix-blend-multiply" 
                  />
                </div>
                
                <div class="space-y-2">
                   <p class="text-sm text-gray-500 dark:text-gray-400">Votre adresse {{ selectedWallet?.currency }}</p>
                   <div class="relative group">
                     <div class="p-4 bg-gray-100 dark:bg-slate-800 rounded-xl break-all font-mono text-sm text-gray-700 dark:text-gray-300 border border-transparent group-hover:border-indigo-500 transition-colors">
                        {{ selectedWallet?.wallet_address }}
                     </div>
                     <button @click="copyAddress" class="absolute right-2 top-2 p-2 bg-indigo-500 text-white rounded-lg shadow-lg hover:scale-105 transition-transform">
                       <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3"/></svg>
                     </button>
                   </div>
                   <p class="text-xs text-yellow-600 dark:text-yellow-500 mt-4 bg-yellow-50 dark:bg-yellow-900/20 p-3 rounded-lg border border-yellow-200 dark:border-yellow-700/30">
                     ⚠️ Envoyez uniquement du <strong>{{ selectedWallet?.currency }}</strong> sur cette adresse. Tout autre jeton sera perdu définitivement.
                   </p>
                </div>
             </div>

             <!-- FIAT DEPOSIT VIEW -->
             <div v-else>
               <!-- Amount Input -->
               <div class="mb-6">
                 <label class="block text-sm font-bold text-gray-700 dark:text-gray-300 mb-2">Montant à déposer</label>
                 <div class="relative">
                   <input 
                     v-model.number="depositAmount" 
                     type="number" 
                     placeholder="0.00"
                     min="1"
                     class="input-premium w-full p-4 text-2xl font-bold rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-slate-800 dark:text-white focus:ring-2 focus:ring-indigo-500 outline-none"
                   />
                   <span class="absolute right-4 top-1/2 -translate-y-1/2 font-bold text-gray-500">{{ selectedWallet?.currency }}</span>
                 </div>
                 <div class="flex gap-2 mt-3">
                   <button v-for="amt in [1000, 5000, 10000, 50000]" :key="amt" 
                     @click="depositAmount = amt"
                     class="flex-1 py-2 text-sm font-medium rounded-lg bg-gray-100 dark:bg-slate-800 hover:bg-indigo-100 dark:hover:bg-indigo-500/20 transition-colors text-gray-700 dark:text-gray-300">
                     {{ amt.toLocaleString() }}
                   </button>
                 </div>
               </div>

               <!-- Payment Method Selection -->
               <div class="mb-6">
                 <label class="block text-sm font-bold text-gray-700 dark:text-gray-300 mb-3">Méthode de paiement</label>
                 <div class="grid grid-cols-1 gap-3">
                   <!-- Mobile Money Section -->
                   <p class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mt-2 mb-1">📱 Mobile Money</p>
                   <button @click="depositMethod = 'orange'" 
                     :class="depositMethod === 'orange' ? 'border-orange-500 bg-orange-50 dark:bg-orange-500/10' : 'border-gray-200 dark:border-gray-700'"
                     class="flex items-center gap-4 p-4 rounded-xl border hover:border-orange-500 transition-all group">
                     <span class="text-2xl">🟠</span>
                     <div class="text-left flex-1">
                       <p class="font-bold text-gray-900 dark:text-white">Orange Money</p>
                       <p class="text-xs text-gray-500 dark:text-gray-400">Paiement instantané</p>
                     </div>
                     <div v-if="depositMethod === 'orange'" class="w-5 h-5 rounded-full bg-orange-500 flex items-center justify-center">
                       <svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                         <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"/>
                       </svg>
                     </div>
                   </button>
                   <button @click="depositMethod = 'mtn'" 
                     :class="depositMethod === 'mtn' ? 'border-yellow-500 bg-yellow-50 dark:bg-yellow-500/10' : 'border-gray-200 dark:border-gray-700'"
                     class="flex items-center gap-4 p-4 rounded-xl border hover:border-yellow-500 transition-all group">
                     <span class="text-2xl">🟡</span>
                     <div class="text-left flex-1">
                       <p class="font-bold text-gray-900 dark:text-white">MTN Mobile Money</p>
                       <p class="text-xs text-gray-500 dark:text-gray-400">Paiement instantané</p>
                     </div>
                     <div v-if="depositMethod === 'mtn'" class="w-5 h-5 rounded-full bg-yellow-500 flex items-center justify-center">
                       <svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                         <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"/>
                       </svg>
                     </div>
                   </button>
                   <button @click="depositMethod = 'wave'" 
                     :class="depositMethod === 'wave' ? 'border-blue-500 bg-blue-50 dark:bg-blue-500/10' : 'border-gray-200 dark:border-gray-700'"
                     class="flex items-center gap-4 p-4 rounded-xl border hover:border-blue-500 transition-all group">
                     <span class="text-2xl">🌊</span>
                     <div class="text-left flex-1">
                       <p class="font-bold text-gray-900 dark:text-white">Wave</p>
                       <p class="text-xs text-gray-500 dark:text-gray-400">Paiement instantané</p>
                     </div>
                     <div v-if="depositMethod === 'wave'" class="w-5 h-5 rounded-full bg-blue-500 flex items-center justify-center">
                       <svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                         <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"/>
                       </svg>
                     </div>
                   </button>

                   <!-- Bank Transfer Section -->
                   <p class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mt-4 mb-1">🏦 Virement</p>
                   <button @click="depositMethod = 'bank'" 
                     :class="depositMethod === 'bank' ? 'border-emerald-500 bg-emerald-50 dark:bg-emerald-500/10' : 'border-gray-200 dark:border-gray-700'"
                     class="flex items-center gap-4 p-4 rounded-xl border hover:border-emerald-500 transition-all group">
                     <span class="text-2xl">🏦</span>
                     <div class="text-left flex-1">
                       <p class="font-bold text-gray-900 dark:text-white">Virement Bancaire</p>
                       <p class="text-xs text-gray-500 dark:text-gray-400">IBAN / RIB • 1-3 jours</p>
                     </div>
                     <div v-if="depositMethod === 'bank'" class="w-5 h-5 rounded-full bg-emerald-500 flex items-center justify-center">
                       <svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                         <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"/>
                       </svg>
                     </div>
                   </button>

                   <!-- Card Section -->
                   <p class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mt-4 mb-1">💳 Carte</p>
                   <button @click="depositMethod = 'card'" 
                     :class="depositMethod === 'card' ? 'border-purple-500 bg-purple-50 dark:bg-purple-500/10' : 'border-gray-200 dark:border-gray-700'"
                     class="flex items-center gap-4 p-4 rounded-xl border hover:border-purple-500 transition-all group">
                     <span class="text-2xl">💳</span>
                     <div class="text-left flex-1">
                       <p class="font-bold text-gray-900 dark:text-white">Carte Bancaire</p>
                       <p class="text-xs text-gray-500 dark:text-gray-400">Visa, Mastercard • Instantané</p>
                     </div>
                     <div v-if="depositMethod === 'card'" class="w-5 h-5 rounded-full bg-purple-500 flex items-center justify-center">
                       <svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                         <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"/>
                       </svg>
                     </div>
                   </button>
                 </div>
               </div>

               <!-- Error/Success Message -->
               <div v-if="depositError" class="mb-4 p-3 rounded-xl bg-red-50 dark:bg-red-500/10 border border-red-200 dark:border-red-500/20 text-red-600 dark:text-red-400 text-sm">
                 {{ depositError }}
               </div>
               <div v-if="depositSuccess" class="mb-4 p-3 rounded-xl bg-green-50 dark:bg-green-500/10 border border-green-200 dark:border-green-500/20 text-green-600 dark:text-green-400 text-sm">
                 {{ depositSuccess }}
               </div>

               <!-- Submit Button -->
               <button 
                 @click="submitDeposit"
                 :disabled="!depositAmount || depositAmount <= 0 || !depositMethod || depositLoading"
                 class="w-full py-4 rounded-xl font-bold text-white bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-lg shadow-indigo-500/30"
               >
                 <span v-if="depositLoading" class="flex items-center justify-center gap-2">
                   <svg class="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
                     <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                     <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
                   </svg>
                   Traitement...
                 </span>
                 <span v-else>Déposer {{ depositAmount ? depositAmount.toLocaleString() : '' }} {{ selectedWallet?.currency }}</span>
               </button>
             </div>
           </div>
        </div>
      </div>

    </div>
    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="fixed inset-0 bg-black/60 backdrop-blur-md z-50 flex items-center justify-center p-4 animate-fade-in-up">
      <div class="bg-white dark:bg-slate-900 rounded-2xl p-6 max-w-md w-full shadow-2xl border border-gray-100 dark:border-gray-800">
         <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Supprimer le portefeuille ?</h3>
         
         <div v-if="walletToDelete" class="mb-6">
           <p class="text-gray-600 dark:text-gray-400 mb-4">
             Vous êtes sur le point de supprimer <strong>{{ walletToDelete.name }}</strong>.
           </p>

           <div v-if="walletToDelete.balance > 0" class="bg-yellow-50 dark:bg-yellow-900/20 p-4 rounded-xl border border-yellow-200 dark:border-yellow-700/50 mb-4">
             <div class="flex items-start gap-3">
               <span class="text-2xl">⚠️</span>
               <div>
                 <p class="font-bold text-yellow-800 dark:text-yellow-200 text-sm">Fonds restants : {{ walletToDelete.balance }} {{ walletToDelete.currency }}</p>
                 <p class="text-xs text-yellow-700 dark:text-yellow-300 mt-1">
                   Les fonds seront automatiquement convertis et envoyés vers votre portefeuille principal ({{ mainWallet?.name }}).
                 </p>
               </div>
             </div>
           </div>

           <p class="text-sm text-gray-500">Cette action est irréversible.</p>
         </div>

         <div v-if="deleteError" class="mb-4 p-3 rounded-xl bg-red-50 dark:bg-red-500/10 border border-red-200 dark:border-red-500/20 text-red-600 dark:text-red-400 text-sm">
           {{ deleteError }}
         </div>

         <div class="flex gap-3 mt-6">
           <button @click="showDeleteModal = false" class="flex-1 py-3 px-4 rounded-xl font-bold text-gray-700 dark:text-gray-200 bg-gray-100 dark:bg-slate-800 hover:bg-gray-200 dark:hover:bg-slate-700 transition-colors">
             Annuler
           </button>
           <button 
             @click="confirmDelete" 
             :disabled="deleteLoading"
             class="flex-1 py-3 px-4 rounded-xl font-bold text-white bg-red-600 hover:bg-red-700 transition-colors shadow-lg shadow-red-500/30 flex justify-center items-center"
           >
             <svg v-if="deleteLoading" class="animate-spin h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24">
               <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
               <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
             </svg>
             {{ deleteLoading ? 'Suppression...' : 'Confirmer' }}
           </button>
         </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { walletAPI, transferAPI } from '~/composables/useApi'
import { useRouter } from 'vue-router'
import { usePin } from '~/composables/usePin'

import { storeToRefs } from 'pinia'
import { useWalletStore } from '~/stores/wallet'

const router = useRouter()
const { requirePin } = usePin()
const walletStore = useWalletStore()

// Delete Logic
const showDeleteModal = ref(false)
const walletToDelete = ref(null)
const deleteLoading = ref(false)
const deleteError = ref('')

// Use store refs
const { wallets, loading } = storeToRefs(walletStore)
// Use store totalBalance directly or via computed if specific formatting needed?
// The store has totalBalance as number. The template uses formatMoney(totalBalance).
// So relying on storeToRefs(walletStore).totalBalance is correct.
const { totalBalance } = storeToRefs(walletStore)

const mainWallet = computed(() => {
  if (!wallets.value.length) return null
  // Oldest wallet is considered Main
  return [...wallets.value].sort((a, b) => new Date(a.created_at) - new Date(b.created_at))[0]
})

const requestDelete = (wallet) => {
  if (mainWallet.value && wallet.id === mainWallet.value.id) {
    alert("Impossible de supprimer votre portefeuille principal.")
    return
  }
  walletToDelete.value = wallet
  deleteError.value = ''
  showDeleteModal.value = true
}

const confirmDelete = async () => {
  if (!walletToDelete.value) return
  
  await requirePin(async (pin) => {
    try {
      deleteLoading.value = true
      deleteError.value = ''

      // 1. Check for funds and transfer if needed
      if (walletToDelete.value.balance > 0 && mainWallet.value) {
        if (walletToDelete.value.currency === mainWallet.value.currency) {
           // Same currency transfer
           await transferAPI.create({
              type: 'internal', 
              amount: walletToDelete.value.balance,
              currency: walletToDelete.value.currency, 
              recipient: mainWallet.value.id,
              description: `Clôture du portefeuille ${walletToDelete.value.name}`
           })
        } else {
           // Different currency
           await transferAPI.create({
              type: 'internal', 
              amount: walletToDelete.value.balance, 
              currency: walletToDelete.value.currency, 
              recipient: mainWallet.value.id,
              description: `Clôture du portefeuille ${walletToDelete.value.name} (Conversion)`
           })
        }
      }

      // 2. Delete Wallet with PIN
      // We assume walletAPI.delete will be updated to accept PIN/options
      await walletAPI.delete(walletToDelete.value.id, pin)
      
      // 3. Refresh
      await fetchWallets()
      showDeleteModal.value = false
      
    } catch (e) {
      console.error(e)
      deleteError.value = e.response?.data?.error || e.message || "Erreur lors de la suppression"
    } finally {
      deleteLoading.value = false
    }
  })
}


// Selected Wallet state
const selectedWallet = ref(null)
const showCreateWallet = ref(false)
const showTopUpModal = ref(false)
const creatingWallet = ref(false)
// loading is now from store

// Deposit form
const depositAmount = ref(5000)
const depositMethod = ref('orange')
const depositLoading = ref(false)
const depositError = ref('')
const depositSuccess = ref('')

// New wallet form
const newWallet = ref({
  type: 'fiat',
  currency: 'USD',
  name: ''
})

// totalBalance is now from store

const formatMoney = (amount) => {
  const val = Number(amount)
  if (amount === undefined || amount === null || isNaN(val)) return '$0.00'
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(val)
}

const formatCrypto = (amount, currency) => {
  const val = Number(amount)
  if (amount === undefined || amount === null || isNaN(val)) return `0.00 ${currency}`
  if (['BTC', 'ETH', 'SOL'].includes(currency)) {
    return `${val.toFixed(6)} ${currency}`
  }
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency }).format(val)
}

const getCurrencyIcon = (currency) => {
  const icons = { BTC: '₿', ETH: 'Ξ', USD: '$', EUR: '€', GBP: '£', SOL: '◎' }
  return icons[currency] || '💰'
}

const getCurrencyBg = (currency) => {
  const bgs = { 
    BTC: 'bg-amber-500', 
    ETH: 'bg-indigo-500', 
    USD: 'bg-emerald-500', 
    EUR: 'bg-blue-600', 
    SOL: 'bg-purple-500',
    GBP: 'bg-rose-500',
    XOF: 'bg-cyan-600',
    XAF: 'bg-teal-600'
  }
  return bgs[currency] || 'bg-slate-500'
}

const truncateAddress = (address) => {
  if (!address || address === 'N/A') return 'N/A'
  return `${address.slice(0, 10)}...${address.slice(-4)}`
}

const selectWallet = (wallet) => {
  selectedWallet.value = wallet
}

const openTopUpModal = () => {
    if (!selectedWallet.value && wallets.value.length > 0) {
        selectedWallet.value = wallets.value[0]
    }
    showTopUpModal.value = true
}

const openTopUpForWallet = (wallet) => {
    selectedWallet.value = wallet
    showTopUpModal.value = true
}

const copyAddress = () => {
  if (selectedWallet.value?.wallet_address) {
    navigator.clipboard.writeText(selectedWallet.value.wallet_address)
    alert('Adresse copiée !')
  }
}

const closeTopUpModal = () => {
  showTopUpModal.value = false
  depositError.value = ''
  depositSuccess.value = ''
}

const submitDeposit = async () => {
  if (!selectedWallet.value || !depositAmount.value || depositAmount.value <= 0) return
  
  depositLoading.value = true
  depositError.value = ''
  depositSuccess.value = ''
  
  try {
    const response = await walletAPI.deposit(selectedWallet.value.id, depositAmount.value, depositMethod.value)
    
    if (response.data) {
      depositSuccess.value = `Dépôt de ${depositAmount.value.toLocaleString()} ${selectedWallet.value.currency} réussi via ${depositMethod.value.toUpperCase()}!`
      // Refresh wallets to show new balance
      await fetchWallets()
      // Reset form after 2 seconds
      setTimeout(() => {
        closeTopUpModal()
        depositAmount.value = 5000
      }, 2000)
    }
  } catch (e) {
    console.error('Deposit error:', e)
    depositError.value = e.response?.data?.error || 'Erreur lors du dépôt. Veuillez réessayer.'
  } finally {
    depositLoading.value = false
  }
}

const fetchWallets = async () => {
  await walletStore.fetchWallets()
}

const createWallet = async () => {
  creatingWallet.value = true
  try {
    const response = await walletAPI.create({
      currency: newWallet.value.currency,
      name: newWallet.value.name || `Mon Portefeuille ${newWallet.value.currency}`,
      wallet_type: newWallet.value.type
    })
    
    // Handle response structure variations
    const rawWallet = response.data?.wallet || response.data || response
    
    if (rawWallet) {
       const wallet = {
          ...rawWallet,
          type: rawWallet.wallet_type || rawWallet.type,
          wallet_type: rawWallet.wallet_type || rawWallet.type
      }
      wallets.value.push(wallet)
      showCreateWallet.value = false
      newWallet.value = { type: 'fiat', currency: 'USD', name: '' }
      fetchWallets() 
    }
  } catch (e) {
    console.error('Error creating wallet:', e)
  } finally {
    creatingWallet.value = false
  }
}

onMounted(async () => {
  // Initialize wallet store if needed (load from cache)
  walletStore.initialize()
  
  await fetchWallets()
  
  if (wallets.value.length > 0 && !selectedWallet.value) {
    selectedWallet.value = wallets.value[0]
  }
})

definePageMeta({
  layout: false, // Explicitly set layout to false since we use NuxtLayout inside template
  middleware: 'auth'
})
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

/* Custom Scrollbar for Modal */
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: rgba(0,0,0,0.05);
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(100,116,139,0.5);
  border-radius: 3px;
}
</style>