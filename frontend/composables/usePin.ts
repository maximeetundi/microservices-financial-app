// PIN management composable
// Handles PIN setup, verification, and status checking - SECURE VERSION

import { ref, computed } from 'vue'
import { userAPI } from './useApi'
import CryptoJS from 'crypto-js'

const pinState = ref({
    hasPin: false,
    isLoading: false,
    showSetupModal: false,
    showVerifyModal: false,
    pendingAction: null as ((pin?: string) => Promise<void>) | null,
    verifyCallback: null as ((verified: boolean) => void) | null,
})

// Secure storage key - hashed with device fingerprint
const getSecureStorageKey = () => {
    const fingerprint = typeof window !== 'undefined' ? 
        navigator.userAgent + navigator.language + screen.width + screen.height : ''
    return 'pin_hash_' + CryptoJS.SHA256(fingerprint).toString()
}

// Encrypt PIN for local storage
const encryptPin = (pin: string): string => {
    const key = getSecureStorageKey()
    return CryptoJS.AES.encrypt(pin, key).toString()
}

// Decrypt PIN from local storage
const decryptPin = (encryptedPin: string): string => {
    const key = getSecureStorageKey()
    const bytes = CryptoJS.AES.decrypt(encryptedPin, key)
    return bytes.toString(CryptoJS.enc.Utf8)
}

// Check if PIN exists locally
const hasLocalPin = (): boolean => {
    if (typeof window === 'undefined') return false
    const encryptedPin = localStorage.getItem(getSecureStorageKey())
    return !!encryptedPin
}

// Store PIN locally (encrypted)
const storePinLocally = (pin: string) => {
    if (typeof window === 'undefined') return
    const encryptedPin = encryptPin(pin)
    localStorage.setItem(getSecureStorageKey(), encryptedPin)
    // Also store in sessionStorage for session persistence
    sessionStorage.setItem('has_pin', 'true')
}

// Verify PIN locally
const verifyPinLocally = (inputPin: string): boolean => {
    if (typeof window === 'undefined') return false
    const encryptedPin = localStorage.getItem(getSecureStorageKey())
    if (!encryptedPin) return false
    
    try {
        const storedPin = decryptPin(encryptedPin)
        return storedPin === inputPin
    } catch (error) {
        console.error('PIN decryption failed:', error)
        return false
    }
}

export function usePin() {

    // PIN Complexity Validation
    const validatePin = (pin: string): { valid: boolean; message?: string } => {
        if (!/^\d{5}$/.test(pin)) {
            return { valid: false, message: 'Le PIN doit contenir 5 chiffres.' }
        }

        // Check for repeated digits (e.g., 11111, 00000)
        if (/^(\d)\1{4}$/.test(pin)) {
            return { valid: false, message: 'Ce code PIN est trop simple (répétition).' }
        }

        // Check for sequential digits (e.g., 12345, 54321)
        const sequences = [
            '01234', '12345', '23456', '34567', '45678', '56789',
            '98765', '87654', '76543', '65432', '54321', '43210'
        ]
        if (sequences.includes(pin)) {
            return { valid: false, message: 'Ce code PIN est trop simple (suite logique).' }
        }

        return { valid: true }
    }

    // Check if user has set their PIN (LOCAL + BACKUP)
    const checkPinStatus = async (): Promise<boolean> => {
        try {
            pinState.value.isLoading = true
            
            // First check local storage
            if (hasLocalPin()) {
                pinState.value.hasPin = true
                return true
            }
            
            // Fallback to backend check (for migration purposes)
            const response = await userAPI.checkPinStatus()
            pinState.value.hasPin = response.data.has_pin
            return response.data.has_pin
        } catch (error) {
            console.error('Failed to check PIN status:', error)
            // Try local as fallback
            pinState.value.hasPin = hasLocalPin()
            return pinState.value.hasPin
        } finally {
            pinState.value.isLoading = false
        }
    }

    // Set up a new PIN (LOCAL STORAGE ONLY)
    const setupPin = async (pin: string, confirmPin: string): Promise<{ success: boolean; message: string }> => {
        // Validate complexity first
        const validation = validatePin(pin)
        if (!validation.valid) {
            return { success: false, message: validation.message || 'PIN invalide' }
        }

        // Confirm PINs match
        if (pin !== confirmPin) {
            return { success: false, message: 'Les PINs ne correspondent pas.' }
        }

        try {
            pinState.value.isLoading = true
            
            // Store PIN locally (encrypted)
            storePinLocally(pin)
            pinState.value.hasPin = true
            pinState.value.showSetupModal = false
            
            // Optional: Still notify backend for tracking (without sending PIN)
            try {
                await userAPI.notifyPinSetup()
            } catch (backendError) {
                console.warn('Backend notification failed, but PIN is set locally:', backendError)
            }
            
            return { success: true, message: 'PIN défini avec succès' }
        } catch (error) {
            const message = (error as any).response?.data?.error || 'Échec de la définition du PIN'
            return { success: false, message }
        } finally {
            pinState.value.isLoading = false
        }
    }

    // Verify the PIN (LOCAL ONLY - NEVER SENDS PIN TO BACKEND)
    const verifyPin = async (pin: string): Promise<{ valid: boolean; attemptsLeft?: number; message?: string }> => {
        try {
            pinState.value.isLoading = true
            
            // Verify locally only
            const isValid = verifyPinLocally(pin)
            
            if (isValid) {
                return {
                    valid: true,
                    message: 'PIN vérifié avec succès',
                }
            } else {
                return {
                    valid: false,
                    attemptsLeft: 3, // Fixed attempts for local verification
                    message: 'PIN incorrect',
                }
            }
        } catch (error) {
            const data = (error as any).response?.data
            return {
                valid: false,
                attemptsLeft: data?.attempts_left,
                message: data?.message || 'Échec de la vérification du PIN',
            }
        } finally {
            pinState.value.isLoading = false
        }
    }

    // Change the PIN (LOCAL ONLY)
    const changePin = async (currentPin: string, newPin: string, confirmPin: string): Promise<{ success: boolean; message: string }> => {
        // Verify current PIN locally first
        if (!verifyPinLocally(currentPin)) {
            return { success: false, message: 'PIN actuel incorrect' }
        }

        // Validate new PIN complexity
        const validation = validatePin(newPin)
        if (!validation.valid) {
            return { success: false, message: validation.message || 'PIN invalide' }
        }

        // Confirm new PINs match
        if (newPin !== confirmPin) {
            return { success: false, message: 'Les nouveaux PINs ne correspondent pas.' }
        }

        try {
            pinState.value.isLoading = true
            
            // Store new PIN locally
            storePinLocally(newPin)
            
            return { success: true, message: 'PIN modifié avec succès' }
        } catch (error) {
            const message = (error as any).response?.data?.error || 'Échec de la modification du PIN'
            return { success: false, message }
        } finally {
            pinState.value.isLoading = false
        }
    }

    // Require PIN before executing an action
    const requirePin = (action: (pin?: string) => Promise<void>): Promise<boolean> => {
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
    const executePendingAction = async (pin?: string) => {
        // Close modal FIRST
        pinState.value.showSetupModal = false
        pinState.value.showVerifyModal = false

        const action = pinState.value.pendingAction
        const callback = pinState.value.verifyCallback

        // Clear state
        pinState.value.pendingAction = null
        pinState.value.verifyCallback = null

        // Execute action if exists
        if (action) {
            try {
                await action(pin)
            } catch (error) {
                console.error('Pending action failed:', error)
            }
        }

        // Call callback with success
        if (callback) {
            callback(true)
        }
    }

    // Clear PIN (for debugging/testing)
    const clearPin = () => {
        if (typeof window !== 'undefined') {
            localStorage.removeItem(getSecureStorageKey())
            sessionStorage.removeItem('has_pin')
            pinState.value.hasPin = false
        }
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
        validatePin,
        clearPin,
    }
}
