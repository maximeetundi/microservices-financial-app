<template>
    <div class="h-full flex flex-col p-6 space-y-6 overflow-hidden">
        <header class="flex justify-between items-center shrink-0">
            <div>
                <h1 class="text-2xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-indigo-400">
                    System Settings
                </h1>
                <p class="text-zinc-400 text-sm mt-1">Configure global platform settings and fees</p>
            </div>
            
            <div class="flex items-center gap-3">
                <button 
                  @click="refreshFees"
                  :disabled="loading"
                  class="p-2 rounded-lg bg-zinc-800/50 hover:bg-zinc-700/50 border border-zinc-700/50 text-zinc-400 transition-colors"
                >
                    <Icon name="lucide:refresh-cw" class="w-5 h-5" :class="{'animate-spin': loading}" />
                </button>
            </div>
        </header>

        <!-- Main Content -->
        <main class="flex-1 overflow-y-auto pr-2 custom-scrollbar space-y-8">
            
            <!-- Fee Configuration Section -->
            <section>
                <div class="flex items-center gap-3 mb-4">
                    <div class="p-2 rounded-lg bg-emerald-500/10 border border-emerald-500/20">
                        <Icon name="lucide:coins" class="w-5 h-5 text-emerald-400" />
                    </div>
                    <h2 class="text-lg font-semibold text-white">Fee Configuration</h2>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
                    <div v-for="fee in fees" :key="fee.id" 
                        class="group relative bg-[#1A1D24] p-5 rounded-xl border border-zinc-800 hover:border-zinc-700 transition-all duration-300"
                    >
                        <!-- Header -->
                        <div class="flex justify-between items-start mb-4">
                            <div class="flex-1">
                                <h3 class="font-medium text-white group-hover:text-blue-400 transition-colors">{{ fee.name }}</h3>
                                <p class="text-xs text-zinc-500 font-mono mt-1">{{ fee.key }}</p>
                            </div>
                            <div :class="[
                                'px-2 py-0.5 rounded text-[10px] font-medium border uppercase tracking-wider',
                                fee.is_enabled ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' : 'bg-red-500/10 text-red-400 border-red-500/20'
                            ]">
                                {{ fee.is_enabled ? 'Active' : 'Disabled' }}
                            </div>
                        </div>

                        <!-- Description -->
                        <p class="text-sm text-zinc-400 mb-6 min-h-[40px]">{{ fee.description }}</p>

                        <!-- Values Display -->
                        <div class="flex items-center gap-4 mb-6">
                            <div class="flex-1 p-3 rounded-lg bg-zinc-900/50 border border-zinc-800/50">
                                <div class="text-[10px] uppercase text-zinc-600 font-bold mb-1">Type</div>
                                <div class="text-sm text-zinc-300 capitalize">{{ fee.type }}</div>
                            </div>
                            
                            <div class="flex-1 p-3 rounded-lg bg-zinc-900/50 border border-zinc-800/50">
                                <template v-if="fee.type === 'percentage'">
                                    <div class="text-[10px] uppercase text-zinc-600 font-bold mb-1">Rate</div>
                                    <div class="text-sm text-white font-mono">{{ fee.percentage_amount }}%</div>
                                </template>
                                <template v-else-if="fee.type === 'flat'">
                                    <div class="text-[10px] uppercase text-zinc-600 font-bold mb-1">Fixed</div>
                                    <div class="text-sm text-white font-mono">{{ fee.fixed_amount }} {{ fee.currency }}</div>
                                </template>
                                <template v-else>
                                    <div class="text-[10px] uppercase text-zinc-600 font-bold mb-1">Combined</div>
                                    <div class="text-sm text-white font-mono flex flex-col leading-tight">
                                        <span>{{ fee.percentage_amount }}%</span>
                                        <span class="text-xs text-zinc-500">+ {{ fee.fixed_amount }} {{ fee.currency }}</span>
                                    </div>
                                </template>
                            </div>
                        </div>

                        <!-- Actions -->
                        <div class="flex items-center justify-between pt-4 border-t border-zinc-800">
                           <div class="text-xs text-zinc-600">
                                Updated: {{ formatDate(fee.updated_at) }}
                           </div>
                           <button 
                             @click="openEditModal(fee)"
                             class="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-zinc-800 hover:bg-blue-600 hover:text-white text-zinc-400 text-xs font-medium transition-all"
                           >
                                <Icon name="lucide:settings-2" class="w-3.5 h-3.5" />
                                Configure
                           </button>
                        </div>
                    </div>

                    <!-- Add New Fee Button (Optional/Future) -->
                    <button 
                        @click="openCreateModal"
                        class="flex flex-col items-center justify-center p-6 rounded-xl border border-zinc-800 border-dashed hover:border-blue-500/50 hover:bg-blue-500/5 transition-all group min-h-[250px]"
                    >
                        <div class="w-12 h-12 rounded-full bg-zinc-900 border border-zinc-800 group-hover:border-blue-500/30 flex items-center justify-center mb-3 transition-colors">
                            <Icon name="lucide:plus" class="w-6 h-6 text-zinc-600 group-hover:text-blue-500" />
                        </div>
                        <span class="text-sm font-medium text-zinc-500 group-hover:text-blue-400">Add New Fee Config</span>
                    </button>
                </div>
            </section>

        </main>

        <!-- Edit Modal -->
        <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
            <div class="bg-[#1A1D24] border border-zinc-800 rounded-xl w-full max-w-lg shadow-2xl transform transition-all">
                <div class="p-6 border-b border-zinc-800 flex justify-between items-center">
                    <div>
                        <h3 class="text-lg font-bold text-white">{{ isEditing ? 'Edit Fee Configuration' : 'New Fee Configuration' }}</h3>
                        <p class="text-sm text-zinc-500 mt-1">{{ form.name || 'Global Setting' }}</p>
                    </div>
                    <button @click="closeModal" class="text-zinc-500 hover:text-white">
                        <Icon name="lucide:x" class="w-5 h-5" />
                    </button>
                </div>
                
                <div class="p-6 space-y-5">
                    <!-- Key (Read-only if editing) -->
                    <div>
                        <label class="block text-xs font-medium text-zinc-400 mb-1.5 uppercase">Configuration Key</label>
                        <input 
                            v-model="form.key" 
                            type="text" 
                            :disabled="isEditing"
                            class="w-full bg-zinc-900/50 border border-zinc-800 rounded-lg px-4 py-2.5 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none disabled:opacity-50 disabled:cursor-not-allowed font-mono"
                            placeholder="e.g. transfer_fee_internal"
                        />
                    </div>

                    <!-- Name & Description -->
                    <div class="grid grid-cols-1 gap-4">
                        <div>
                            <label class="block text-xs font-medium text-zinc-400 mb-1.5 uppercase">Display Name</label>
                            <input 
                                v-model="form.name"
                                type="text"
                                class="w-full bg-zinc-900/50 border border-zinc-800 rounded-lg px-4 py-2.5 text-white text-sm focus:border-blue-500 outline-none"
                            />
                        </div>
                         <div>
                            <label class="block text-xs font-medium text-zinc-400 mb-1.5 uppercase">Description</label>
                            <textarea 
                                v-model="form.description"
                                rows="2"
                                class="w-full bg-zinc-900/50 border border-zinc-800 rounded-lg px-4 py-2.5 text-white text-sm focus:border-blue-500 outline-none resize-none"
                            ></textarea>
                        </div>
                    </div>

                    <!-- Configuration Type -->
                    <div>
                        <label class="block text-xs font-medium text-zinc-400 mb-1.5 uppercase">Fee Structure Type</label>
                        <div class="grid grid-cols-3 gap-2">
                            <button 
                                v-for="type in ['flat', 'percentage', 'hybrid']" 
                                :key="type"
                                @click="form.type = type"
                                :class="[
                                    'px-3 py-2 rounded-lg text-sm font-medium border transition-all capitalize',
                                    form.type === type 
                                        ? 'bg-blue-500/10 border-blue-500/50 text-blue-400' 
                                        : 'bg-zinc-900 border-zinc-800 text-zinc-500 hover:border-zinc-700'
                                ]"
                            >
                                {{ type }}
                            </button>
                        </div>
                    </div>

                    <!-- Values Input -->
                    <div class="grid grid-cols-2 gap-4">
                        <div v-if="form.type !== 'flat'">
                            <label class="block text-xs font-medium text-zinc-400 mb-1.5 uppercase">Percentage (%)</label>
                            <div class="relative">
                                <input 
                                    v-model.number="form.percentage_amount"
                                    type="number" step="0.01" min="0" max="100"
                                    class="w-full bg-zinc-900/50 border border-zinc-800 rounded-lg pl-3 pr-8 py-2.5 text-white text-sm focus:border-blue-500 outline-none"
                                />
                                <span class="absolute right-3 top-2.5 text-zinc-500 text-xs">%</span>
                            </div>
                        </div>
                        
                        <div v-if="form.type !== 'percentage'">
                            <label class="block text-xs font-medium text-zinc-400 mb-1.5 uppercase">Fixed Amount</label>
                            <div class="flex gap-2">
                                <input 
                                    v-model.number="form.fixed_amount"
                                    type="number" step="0.01" min="0"
                                    class="w-full bg-zinc-900/50 border border-zinc-800 rounded-lg px-3 py-2.5 text-white text-sm focus:border-blue-500 outline-none"
                                />
                                <input 
                                    v-model="form.currency"
                                    type="text"
                                    class="w-20 bg-zinc-900/50 border border-zinc-800 rounded-lg px-2 py-2.5 text-white text-sm text-center uppercase focus:border-blue-500 outline-none"
                                    placeholder="EUR"
                                />
                            </div>
                        </div>
                    </div>

                    <!-- Status Toggle -->
                     <div class="flex items-center justify-between p-4 bg-zinc-900/50 rounded-lg border border-zinc-800">
                        <div>
                            <div class="text-sm font-medium text-white">Status</div>
                            <div class="text-xs text-zinc-500">Enable or disable this fee rule</div>
                        </div>
                        <button 
                            @click="form.is_enabled = !form.is_enabled"
                            :class="[
                                'relative inline-flex h-6 w-11 items-center rounded-full transition-colors',
                                form.is_enabled ? 'bg-emerald-500' : 'bg-zinc-700'
                            ]"
                        >
                            <span 
                                :class="[
                                    'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
                                    form.is_enabled ? 'translate-x-[22px]' : 'translate-x-1'
                                ]" 
                            />
                        </button>
                    </div>

                </div>

                <div class="p-6 border-t border-zinc-800 flex justify-end gap-3 bg-zinc-900/30">
                    <button 
                        @click="closeModal"
                        class="px-4 py-2 rounded-lg text-sm font-medium text-zinc-400 hover:text-white transition-colors"
                    >
                        Cancel
                    </button>
                    <button 
                        @click="saveFee"
                        :disabled="saving"
                        class="px-4 py-2 rounded-lg bg-blue-600 hover:bg-blue-500 text-white text-sm font-medium transition-colors disabled:opacity-50 flex items-center gap-2"
                    >
                        <Icon v-if="saving" name="lucide:loader-2" class="w-4 h-4 animate-spin" />
                        {{ isEditing ? 'Save Changes' : 'Create Configuration' }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useApi } from '@/composables/useApi'

definePageMeta({
  layout: 'admin'
})

const { adminGlobalFeeApi } = useApi()

// Data
const fees = ref<any[]>([])
const loading = ref(true)
const saving = ref(false)
const showModal = ref(false)
const isEditing = ref(false)

// Form
const defaultForm = {
    key: '',
    name: '',
    description: '',
    type: 'percentage',
    fixed_amount: 0,
    percentage_amount: 0,
    currency: 'EUR',
    is_enabled: true
}
const form = ref({ ...defaultForm })

// Methods
const formatDate = (dateString: string) => {
    if (!dateString) return '-'
    return new Date(dateString).toLocaleDateString('en-US', { 
        year: 'numeric', month: 'short', day: 'numeric',
        hour: '2-digit', minute: '2-digit'
    })
}

const refreshFees = async () => {
    loading.value = true
    try {
        const res = await adminGlobalFeeApi.getAll()
        if (res.data?.fees) {
            fees.value = res.data.fees
        } else {
             // Fallback mocked data if API fails locally during dev
             fees.value = []
        }
    } catch (err) {
        console.error('Failed to load fees:', err)
        // Mock data for display preview
        /*
        fees.value = [
            { id: '1', key: 'transfer_internal', name: 'Internal Transfer Fee', description: 'Fee for transfers between users', type: 'percentage', percentage_amount: 0.5, is_enabled: true, updated_at: new Date().toISOString() },
            { id: '2', key: 'transfer_international', name: 'International Transfer Fee', description: 'Fee for international transfers', type: 'hybrid', fixed_amount: 5, currency: 'EUR', percentage_amount: 1.5, is_enabled: true, updated_at: new Date().toISOString() },
        ]
        */
    } finally {
        loading.value = false
    }
}

const openEditModal = (fee: any) => {
    isEditing.value = true
    form.value = { ...fee }
    showModal.value = true
}

const openCreateModal = () => {
    isEditing.value = false
    form.value = { ...defaultForm }
    showModal.value = true
}

const closeModal = () => {
    showModal.value = false
    form.value = { ...defaultForm }
}

const saveFee = async () => {
    saving.value = true
    try {
        if (isEditing.value) {
            await adminGlobalFeeApi.update(form.value.key, {
                type: form.value.type,
                fixed_amount: form.value.fixed_amount,
                percentage_amount: form.value.percentage_amount,
                is_enabled: form.value.is_enabled,
                currency: form.value.currency
            })
        } else {
            await adminGlobalFeeApi.create(form.value)
        }
        await refreshFees()
        closeModal()
    } catch (err) {
        console.error('Failed to save fee:', err)
        alert('Failed to save configuration provided.')
    } finally {
        saving.value = false
    }
}

onMounted(() => {
    refreshFees()
})
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
    width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
    background: rgba(255, 255, 255, 0.02);
}
.custom-scrollbar::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.1);
    border-radius: 2px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: rgba(255, 255, 255, 0.2);
}
</style>
