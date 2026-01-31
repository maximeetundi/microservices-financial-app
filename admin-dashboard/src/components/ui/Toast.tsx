'use client';

import { useEffect, useState, createContext, useContext, useCallback } from 'react';
import {
    CheckCircleIcon,
    ExclamationTriangleIcon,
    InformationCircleIcon,
    XCircleIcon,
    XMarkIcon,
} from '@heroicons/react/24/outline';

type ToastType = 'success' | 'error' | 'warning' | 'info';

interface Toast {
    id: string;
    message: string;
    type: ToastType;
    duration?: number;
}

interface ToastContextType {
    showToast: (message: string, type: ToastType, duration?: number) => void;
    success: (message: string) => void;
    error: (message: string) => void;
    warning: (message: string) => void;
    info: (message: string) => void;
}

const ToastContext = createContext<ToastContextType | null>(null);

export function useToast() {
    const context = useContext(ToastContext);
    if (!context) {
        throw new Error('useToast must be used within a ToastProvider');
    }
    return context;
}

export function ToastProvider({ children }: { children: React.ReactNode }) {
    const [toasts, setToasts] = useState<Toast[]>([]);

    const showToast = useCallback((message: string, type: ToastType, duration = 4000) => {
        const id = Math.random().toString(36).substr(2, 9);
        setToasts((prev) => [...prev, { id, message, type, duration }]);
    }, []);

    const removeToast = useCallback((id: string) => {
        setToasts((prev) => prev.filter((t) => t.id !== id));
    }, []);

    const success = useCallback((message: string) => showToast(message, 'success'), [showToast]);
    const error = useCallback((message: string) => showToast(message, 'error'), [showToast]);
    const warning = useCallback((message: string) => showToast(message, 'warning'), [showToast]);
    const info = useCallback((message: string) => showToast(message, 'info'), [showToast]);

    return (
        <ToastContext.Provider value={{ showToast, success, error, warning, info }}>
            {children}
            <ToastContainer toasts={toasts} onRemove={removeToast} />
        </ToastContext.Provider>
    );
}

function ToastContainer({ toasts, onRemove }: { toasts: Toast[]; onRemove: (id: string) => void }) {
    return (
        <div className="fixed bottom-6 right-6 z-[100] flex flex-col gap-3">
            {toasts.map((toast) => (
                <ToastItem key={toast.id} toast={toast} onRemove={() => onRemove(toast.id)} />
            ))}
        </div>
    );
}

function ToastItem({ toast, onRemove }: { toast: Toast; onRemove: () => void }) {
    useEffect(() => {
        const timer = setTimeout(onRemove, toast.duration || 4000);
        return () => clearTimeout(timer);
    }, [toast.duration, onRemove]);

    const config = {
        success: {
            bg: 'from-green-500 to-emerald-600',
            icon: CheckCircleIcon,
        },
        error: {
            bg: 'from-red-500 to-pink-600',
            icon: XCircleIcon,
        },
        warning: {
            bg: 'from-amber-500 to-orange-600',
            icon: ExclamationTriangleIcon,
        },
        info: {
            bg: 'from-blue-500 to-indigo-600',
            icon: InformationCircleIcon,
        },
    };

    const { bg, icon: Icon } = config[toast.type];

    return (
        <div
            className={`flex items-center gap-3 px-5 py-4 bg-gradient-to-r ${bg} text-white rounded-xl shadow-2xl animate-slide-up min-w-[300px] max-w-md`}
        >
            <Icon className="w-5 h-5 flex-shrink-0" />
            <span className="font-medium flex-1">{toast.message}</span>
            <button
                onClick={onRemove}
                className="p-1 hover:bg-white/20 rounded-lg transition-colors flex-shrink-0"
            >
                <XMarkIcon className="w-4 h-4" />
            </button>
        </div>
    );
}

// Simple standalone toast component (for pages without provider)
export function SimpleToast({
    message,
    type,
    onClose,
}: {
    message: string;
    type: ToastType;
    onClose: () => void;
}) {
    useEffect(() => {
        const timer = setTimeout(onClose, 4000);
        return () => clearTimeout(timer);
    }, [onClose]);

    const config = {
        success: { bg: 'from-green-500 to-emerald-600', icon: CheckCircleIcon },
        error: { bg: 'from-red-500 to-pink-600', icon: XCircleIcon },
        warning: { bg: 'from-amber-500 to-orange-600', icon: ExclamationTriangleIcon },
        info: { bg: 'from-blue-500 to-indigo-600', icon: InformationCircleIcon },
    };

    const { bg, icon: Icon } = config[type];

    return (
        <div className="fixed bottom-6 right-6 z-[100] animate-slide-up">
            <div className={`flex items-center gap-3 px-5 py-4 bg-gradient-to-r ${bg} text-white rounded-xl shadow-2xl`}>
                <Icon className="w-5 h-5" />
                <span className="font-medium">{message}</span>
                <button onClick={onClose} className="ml-2 hover:bg-white/20 p-1 rounded-lg transition-colors">
                    <XMarkIcon className="w-4 h-4" />
                </button>
            </div>
        </div>
    );
}
