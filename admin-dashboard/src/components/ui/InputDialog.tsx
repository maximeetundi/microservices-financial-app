'use client';

import { useState, useEffect } from 'react';
import Modal, { ModalFooter, ModalButton } from './Modal';

interface InputDialogProps {
    isOpen: boolean;
    onClose: () => void;
    onSubmit: (value: string) => void;
    title: string;
    subtitle?: string;
    label: string;
    placeholder?: string;
    defaultValue?: string;
    required?: boolean;
    multiline?: boolean;
    rows?: number;
    submitText?: string;
    cancelText?: string;
    variant?: 'primary' | 'danger' | 'warning';
    loading?: boolean;
}

export default function InputDialog({
    isOpen,
    onClose,
    onSubmit,
    title,
    subtitle,
    label,
    placeholder = '',
    defaultValue = '',
    required = true,
    multiline = false,
    rows = 3,
    submitText = 'Confirmer',
    cancelText = 'Annuler',
    variant = 'primary',
    loading = false,
}: InputDialogProps) {
    const [value, setValue] = useState(defaultValue);

    // Reset value when dialog opens
    useEffect(() => {
        if (isOpen) {
            setValue(defaultValue);
        }
    }, [isOpen, defaultValue]);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (required && !value.trim()) return;
        onSubmit(value);
    };

    const handleClose = () => {
        setValue('');
        onClose();
    };

    const variantMap = {
        primary: 'primary',
        danger: 'danger',
        warning: 'primary',
    } as const;

    return (
        <Modal isOpen={isOpen} onClose={handleClose} title={title} subtitle={subtitle} size="md">
            <form onSubmit={handleSubmit}>
                <div className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            {label}
                            {required && <span className="text-red-500 ml-1">*</span>}
                        </label>
                        {multiline ? (
                            <textarea
                                value={value}
                                onChange={(e) => setValue(e.target.value)}
                                placeholder={placeholder}
                                rows={rows}
                                className="w-full px-4 py-3 rounded-xl border-2 border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 outline-none transition-all resize-none text-gray-900"
                                autoFocus
                            />
                        ) : (
                            <input
                                type="text"
                                value={value}
                                onChange={(e) => setValue(e.target.value)}
                                placeholder={placeholder}
                                className="w-full px-4 py-3 rounded-xl border-2 border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 outline-none transition-all text-gray-900"
                                autoFocus
                            />
                        )}
                    </div>
                </div>

                <ModalFooter>
                    <ModalButton type="button" variant="secondary" onClick={handleClose} disabled={loading}>
                        {cancelText}
                    </ModalButton>
                    <ModalButton
                        type="submit"
                        variant={variantMap[variant]}
                        disabled={required && !value.trim()}
                        loading={loading}
                    >
                        {submitText}
                    </ModalButton>
                </ModalFooter>
            </form>
        </Modal>
    );
}
