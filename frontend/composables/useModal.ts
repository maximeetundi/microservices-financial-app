// Modal management composable
// Replaces alert(), confirm(), prompt() with elegant modals

import { ref, reactive, computed } from 'vue'

export type ModalType = 'success' | 'error' | 'warning' | 'info' | 'confirm'

export interface ModalOptions {
    type: ModalType
    title: string
    message: string
    confirmText?: string
    cancelText?: string
    showCancel?: boolean
    icon?: string
    persistent?: boolean // Can't close by clicking outside
}

export interface ModalState extends ModalOptions {
    isOpen: boolean
    resolve?: (value: boolean) => void
}

const modalState = reactive<ModalState>({
    isOpen: false,
    type: 'info',
    title: '',
    message: '',
    confirmText: 'OK',
    cancelText: 'Annuler',
    showCancel: false,
    persistent: false,
})

export function useModal() {
    const show = (options: ModalOptions): Promise<boolean> => {
        return new Promise((resolve) => {
            Object.assign(modalState, {
                ...options,
                isOpen: true,
                confirmText: options.confirmText || 'OK',
                cancelText: options.cancelText || 'Annuler',
                showCancel: options.showCancel ?? false,
                persistent: options.persistent ?? false,
                resolve,
            })
        })
    }

    const close = (confirmed: boolean = false) => {
        if (modalState.resolve) {
            modalState.resolve(confirmed)
        }
        modalState.isOpen = false
    }

    // Convenience methods
    const success = (title: string, message: string) => {
        return show({ type: 'success', title, message })
    }

    const error = (title: string, message: string) => {
        return show({ type: 'error', title, message })
    }

    const warning = (title: string, message: string) => {
        return show({ type: 'warning', title, message })
    }

    const info = (title: string, message: string) => {
        return show({ type: 'info', title, message })
    }

    const confirm = (title: string, message: string, confirmText = 'Confirmer', cancelText = 'Annuler') => {
        return show({
            type: 'confirm',
            title,
            message,
            showCancel: true,
            confirmText,
            cancelText,
        })
    }

    return {
        state: modalState,
        show,
        close,
        success,
        error,
        warning,
        info,
        confirm,
    }
}
