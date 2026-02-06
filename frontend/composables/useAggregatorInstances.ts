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
    credentials: Record<string, string>
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
            // Use the correct backend endpoint
            const url = aggregatorId
                ? `/api/v1/admin/payment-providers/${aggregatorId}/instances`
                : '/api/v1/admin/provider-instances'

            const response = await $fetch<{ instances: AggregatorInstance[] }>(url)
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
            // POST /api/v1/admin/payment-providers/:id/instances
            const response = await $fetch<AggregatorInstance>(
                `/api/v1/admin/payment-providers/${data.aggregator_id}/instances`,
                {
                    method: 'POST',
                    body: {
                        name: data.instance_name,
                        priority: data.priority,
                        is_test_mode: data.is_test_mode,
                        // Should map other fields if needed, backend expects 'name', 'priority', 'credentials'
                        credentials: data.credentials
                    }
                }
            )
            instances.value.push(response)
            return response
        } catch (err: any) {
            error.value = err.message || 'Failed to create instance'
            throw err
        } finally {
            loading.value = false
        }
    }

    const updateInstance = async (instanceId: string, data: Partial<CreateInstanceRequest> & { aggregator_id?: string }) => {
        loading.value = true
        error.value = null

        try {
            // We need the provider ID (aggregator_id) to construct the URL
            // If not in data, we might need to find it in existing instances or require it
            const instance = instances.value.find(i => i.id === instanceId)
            const aggregatorId = data.aggregator_id || instance?.aggregator_id || instance?.provider_code // fallback?

            if (!aggregatorId) {
                throw new Error("Aggregator ID required for update")
            }

            const response = await $fetch<AggregatorInstance>(
                `/api/v1/admin/payment-providers/${aggregatorId}/instances/${instanceId}`,
                {
                    method: 'PUT',
                    body: {
                        name: data.instance_name,
                        priority: data.priority,
                        is_active: true, // assume active on update? or passed
                        credentials: data.credentials
                    }
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

    const deleteInstance = async (instanceId: string, aggregatorId: string) => {
        loading.value = true
        error.value = null

        try {
            await $fetch(`/api/v1/admin/payment-providers/${aggregatorId}/instances/${instanceId}`, {
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
