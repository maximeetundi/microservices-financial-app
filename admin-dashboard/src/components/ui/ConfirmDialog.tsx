'use client';

import { ExclamationTriangleIcon, ShieldExclamationIcon, InformationCircleIcon } from '@heroicons/react/24/outline';
import Modal, { ModalFooter, ModalButton } from './Modal';

interface ConfirmDialogProps {
    isOpen: boolean;
    onClose: () => void;
    onConfirm: () => void;
    title: string;
    message: string | React.ReactNode;
    confirmText?: string;
    cancelText?: string;
    variant?: 'danger' | 'warning' | 'info';
    loading?: boolean;
}

export default function ConfirmDialog({
    isOpen,
    onClose,
    onConfirm,
    title,
    message,
    confirmText = 'Confirmer',
    cancelText = 'Annuler',
    variant = 'danger',
    loading = false,
}: ConfirmDialogProps) {
    const config = {
        danger: {
            icon: ShieldExclamationIcon,
            iconBg: 'bg-red-100',
            iconColor: 'text-red-600',
            buttonVariant: 'danger' as const,
        },
        warning: {
            icon: ExclamationTriangleIcon,
            iconBg: 'bg-amber-100',
            iconColor: 'text-amber-600',
            buttonVariant: 'primary' as const,
        },
        info: {
            icon: InformationCircleIcon,
            iconBg: 'bg-blue-100',
            iconColor: 'text-blue-600',
            buttonVariant: 'primary' as const,
        },
    };

    const { icon: Icon, iconBg, iconColor, buttonVariant } = config[variant];

    return (
        <Modal isOpen={isOpen} onClose={onClose} title={title} size="sm">
            <div className="text-center">
                <div className={`mx-auto w-16 h-16 ${iconBg} rounded-2xl flex items-center justify-center mb-4`}>
                    <Icon className={`w-8 h-8 ${iconColor}`} />
                </div>
                <div className="text-gray-600 text-sm leading-relaxed">
                    {typeof message === 'string' ? <p>{message}</p> : message}
                </div>
            </div>

            <ModalFooter>
                <ModalButton variant="secondary" onClick={onClose} disabled={loading}>
                    {cancelText}
                </ModalButton>
                <ModalButton variant={buttonVariant} onClick={onConfirm} loading={loading}>
                    {confirmText}
                </ModalButton>
            </ModalFooter>
        </Modal>
    );
}
