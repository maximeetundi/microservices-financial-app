// PIN management composable
// Handles PIN setup, verification, and status checking

import { ref, computed } from 'vue'
import { userAPI } from './useApi'

const pinState = ref({
    hasPin: false,
    isLoading: false,
    showSetupModal: false,
    showVerifyModal: false,
    pendingAction: null as (() => Promise<void>) | null,
    verifyCallback: null as ((verified: boolean) => void) | null,
})

export function usePin() {

    // Check if user has set their PIN
    const checkPinStatus = async (): Promise<boolean> => {
        try {
            pinState.value.isLoading = true
            const response = await userAPI.checkPinStatus()
            pinState.value.hasPin = response.data.has_pin
            return response.data.has_pin
        } catch (error) {
            console.error('Failed to check PIN status:', error)
            return false
        } finally {
            pinState.value.isLoading = false
        }
    }

    // Set up a new PIN
    const setupPin = async (pin: string, confirmPin: string): Promise<{ success: boolean; message: string }> => {
        try {
            pinState.value.isLoading = true
            await userAPI.setupPin({ pin, confirm_pin: confirmPin })
            pinState.value.hasPin = true
            pinState.value.showSetupModal = false
            return { success: true, message: 'PIN défini avec succès' }
        } catch (error: any) {
            const message = error.response?.data?.error || 'Échec de la définition du PIN'
            return { success: false, message }
        } finally {
            pinState.value.isLoading = false
        }
    }

    // Verify the PIN
    const verifyPin = async (pin: string): Promise<{ valid: boolean; attemptsLeft?: number; message?: string }> => {
        try {
            pinState.value.isLoading = true
            const response = await userAPI.verifyPin({ pin })
            return {
                valid: response.data.valid,
                attemptsLeft: response.data.attempts_left,
                message: response.data.message,
            }
        } catch (error: any) {
            const data = error.response?.data
            return {
                valid: false,
                attemptsLeft: data?.attempts_left,
                message: data?.message || 'Échec de la vérification du PIN',
            }
        } finally {
            pinState.value.isLoading = false
        }
    }

    // Change the PIN
    const changePin = async (currentPin: string, newPin: string, confirmPin: string): Promise<{ success: boolean; message: string }> => {
        try {
            pinState.value.isLoading = true
            await userAPI.changePin({
                current_pin: currentPin,
                new_pin: newPin,
                confirm_pin: confirmPin
            })
            return { success: true, message: 'PIN modifié avec succès' }
        } catch (error: any) {
            const message = error.response?.data?.error || 'Échec de la modification du PIN'
            return { success: false, message }
        } finally {
            pinState.value.isLoading = false
        }
    }

    // Require PIN before executing an action
    const requirePin = (action: () => Promise<void>): Promise<boolean> => {
        return new Promise((resolve) => {
            if (!pinState.value.hasPin) {
                // No PIN set, require setup first
                pinState.value.showSetupModal = true
                pinState.value.pendingAction = action
                pinState.value.verifyCallback = resolve
            } else {
                // Show verification modal
                pinState.value.showVerifyModal = true
                pinState.value.pendingAction = action
                pinState.value.verifyCallback = resolve
            }
        })
    }

    // Force show PIN setup modal
    const showPinSetup = () => {
        pinState.value.showSetupModal = true
    }

    // Close modals
    const closeModals = () => {
        pinState.value.showSetupModal = false
        pinState.value.showVerifyModal = false
        pinState.value.pendingAction = null
        if (pinState.value.verifyCallback) {
            pinState.value.verifyCallback(false)
            pinState.value.verifyCallback = null
        }
    }

    // Execute pending action after successful verification
    const executePendingAction = async () => {
        if (pinState.value.pendingAction) {
            try {
                await pinState.value.pendingAction()
            } catch (error) {
                console.error('Pending action failed:', error)
            }
            pinState.value.pendingAction = null
        }
        if (pinState.value.verifyCallback) {
            pinState.value.verifyCallback(true)
            pinState.value.verifyCallback = null
        }
        closeModals()
    }

    return {
        state: pinState,
        hasPin: computed(() => pinState.value.hasPin),
        isLoading: computed(() => pinState.value.isLoading),
        showSetupModal: computed(() => pinState.value.showSetupModal),
        showVerifyModal: computed(() => pinState.value.showVerifyModal),
        checkPinStatus,
        setupPin,
        verifyPin,
        changePin,
        requirePin,
        showPinSetup,
        closeModals,
        executePendingAction,
    }
}
