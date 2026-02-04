import { ref } from 'vue'

export interface AggregatorInstance {
    id: string
    aggregator_id: string
    instance_name: string
    provider_code: string
    provider_name: string
    enabled: boolean
    priority: number
    is_test_mode: boolean
    restricted_countries: string[]
    daily_limit?: number
    monthly_limit?: number
    daily_usage: number
    monthly_usage: number
    total_transactions: number
    total_volume: number
    last_used_at?: string
    wallets: InstanceWallet[]
    notes?: string
    created_at: string
    updated_at: string
}

export interface InstanceWallet {
    id: string
    instance_id: string
    hot_wallet_id: string
    is_primary: boolean
    priority: number
    min_balance?: number
    max_balance?: number
    auto_recharge_enabled: boolean
    recharge_threshold?: number
    recharge_target?: number
    recharge_source_wallet_id?: string
    wallet_currency: string
    wallet_balance: number
    total_deposits: number
    total_withdrawals: number
    transaction_count: number
    enabled: boolean
}

export interface CreateInstanceRequest {
    aggregator_id: string
    instance_name: string
    api_credentials: Record<string, string>
    priority?: number
    daily_limit?: number
    monthly_limit?: number
    restricted_countries?: string[]
    is_test_mode?: boolean
    notes?: string
}

export interface LinkWalletRequest {
    hot_wallet_id: string
    is_primary?: boolean
    priority?: number
    min_balance?: number
    max_balance?: number
    auto_recharge_enabled?: boolean
    recharge_threshold?: number
    recharge_target?: number
    recharge_source_wallet_id?: string
}

export const useAggregatorInstances = () => {
    const instances = ref<AggregatorInstance[]>([])
    const loading = ref(false)
    const error = ref<string | null>(null)

    const fetchInstances = async (aggregatorId?: string) => {
        loading.value = true
        error.value = null

        try {
            const params = aggregatorId ? `?aggregator_id=${aggregatorId}` : ''
            const response = await $fetch<{ instances: AggregatorInstance[] }>(
                `/api/admin/instances${params}`
            )
            instances.value = response.instances
        } catch (err: any) {
            error.value = err.message || 'Failed to fetch instances'
            console.error('Error fetching instances:', err)
        } finally {
            loading.value = false
        }
    }

    const createInstance = async (data: CreateInstanceRequest) => {
        loading.value = true
        error.value = null

        try {
            const response = await $fetch<AggregatorInstance>('/api/admin/instances', {
                method: 'POST',
                body: data
            })
            instances.value.push(response)
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to create instance'
            throw err
        } finally {
            loading.value = false
        }
    }

    const updateInstance = async (instanceId: string, data: Partial<CreateInstanceRequest>) => {
        loading.value = true
        error.value = null

        try {
            const response = await $fetch<AggregatorInstance>(
                `/api/admin/instances/${instanceId}`,
                {
                    method: 'PUT',
                    body: data
                }
            )

            const index = instances.value.findIndex(i => i.id === instanceId)
            if (index !== -1) {
                instances.value[index] = response
            }

            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to update instance'
            throw err
        } finally {
            loading.value = false
        }
    }

    const deleteInstance = async (instanceId: string) => {
        loading.value = true
        error.value = null

        try {
            await $fetch(`/api/admin/instances/${instanceId}`, {
                method: 'DELETE'
            })

            instances.value = instances.value.filter(i => i.id !== instanceId)
        } catch (err: any) {
            error.value = err.message || 'Failed to delete instance'
            throw err
        } finally {
            loading.value = false
        }
    }

    const linkWallet = async (instanceId: string, data: LinkWalletRequest) => {
        loading.value = true
        error.value = null

        try {
            const response = await $fetch<InstanceWallet>(
                `/api/admin/instances/${instanceId}/wallets`,
                {
                    method: 'POST',
                    body: data
                }
            )

            // Update instance wallets
            const instance = instances.value.find(i => i.id === instanceId)
            if (instance) {
                if (!instance.wallets) instance.wallets = []
                instance.wallets.push(response)
            }

            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to link wallet'
            throw err
        } finally {
            loading.value = false
        }
    }

    const unlinkWallet = async (instanceId: string, walletId: string) => {
        loading.value = true
        error.value = null

        try {
            await $fetch(`/api/admin/instances/${instanceId}/wallets/${walletId}`, {
                method: 'DELETE'
            })

            // Update instance wallets
            const instance = instances.value.find(i => i.id === instanceId)
            if (instance && instance.wallets) {
                instance.wallets = instance.wallets.filter(w => w.id !== walletId)
            }
        } catch (err: any) {
            error.value = err.message || 'Failed to unlink wallet'
            throw err
        } finally {
            loading.value = false
        }
    }

    const updateWallet = async (
        instanceId: string,
        walletId: string,
        data: Partial<LinkWalletRequest>
    ) => {
        loading.value = true
        error.value = null

        try {
            const response = await $fetch<InstanceWallet>(
                `/api/admin/instances/${instanceId}/wallets/${walletId}`,
                {
                    method: 'PUT',
                    body: data
                }
            )

            // Update instance wallet
            const instance = instances.value.find(i => i.id === instanceId)
            if (instance && instance.wallets) {
                const index = instance.wallets.findIndex(w => w.id === walletId)
                if (index !== -1) {
                    instance.wallets[index] = response
                }
            }

            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to update wallet'
            throw err
        } finally {
            loading.value = false
        }
    }

    return {
        instances,
        loading,
        error,
        fetchInstances,
        createInstance,
        updateInstance,
        deleteInstance,
        linkWallet,
        unlinkWallet,
        updateWallet
    }
}
