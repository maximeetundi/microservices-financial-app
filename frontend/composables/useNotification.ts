import { ref } from 'vue'

interface Notification {
    visible: boolean
    type: 'success' | 'error' | 'warning' | 'info' | 'confirm'
    title?: string
    message: string
    confirmText?: string
    cancelText?: string
    showCancel?: boolean
    onConfirm?: () => void
    onCancel?: () => void
}

const notification = ref<Notification>({
    visible: false,
    type: 'info',
    message: '',
})

export const useNotification = () => {
    const showSuccess = (message: string, title = 'Succès') => {
        notification.value = {
            visible: true,
            type: 'success',
            title,
            message,
            confirmText: 'OK',
            showCancel: false,
        }
    }

    const showError = (message: string, title = 'Erreur') => {
        notification.value = {
            visible: true,
            type: 'error',
            title,
            message,
            confirmText: 'OK',
            showCancel: false,
        }
    }

    const showWarning = (message: string, title = 'Attention') => {
        notification.value = {
            visible: true,
            type: 'warning',
            title,
            message,
            confirmText: 'OK',
            showCancel: false,
        }
    }

    const showInfo = (message: string, title = 'Information') => {
        notification.value = {
            visible: true,
            type: 'info',
            title,
            message,
            confirmText: 'OK',
            showCancel: false,
        }
    }

    const confirm = (message: string, title = 'Confirmation', onConfirm?: () => void) => {
        notification.value = {
            visible: true,
            type: 'confirm',
            title,
            message,
            confirmText: 'Confirmer',
            cancelText: 'Annuler',
            showCancel: true,
            onConfirm,
        }
    }

    const close = () => {
        notification.value.visible = false
    }

    const handleConfirm = () => {
        if (notification.value.onConfirm) {
            notification.value.onConfirm()
        }
        close()
    }

    const handleCancel = () => {
        if (notification.value.onCancel) {
            notification.value.onCancel()
        }
        close()
    }

    return {
        notification,
        showSuccess,
        showError,
        showWarning,
        showInfo,
        confirm,
        close,
        handleConfirm,
        handleCancel,
    }
}
