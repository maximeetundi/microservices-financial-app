'use client';

import { ArrowPathIcon } from '@heroicons/react/24/outline';

interface PageHeaderProps {
    title: string;
    subtitle?: string;
    icon?: string; // Emoji
    onRefresh?: () => void;
    loading?: boolean;
    actions?: React.ReactNode;
}

export default function PageHeader({
    title,
    subtitle,
    icon,
    onRefresh,
    loading = false,
    actions,
}: PageHeaderProps) {
    return (
        <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4 mb-6">
            <div>
                <h1 className="text-2xl lg:text-3xl font-bold text-gray-900 flex items-center gap-3">
                    {icon && <span className="text-2xl">{icon}</span>}
                    {title}
                </h1>
                {subtitle && (
                    <p className="text-gray-500 text-sm mt-1">{subtitle}</p>
                )}
            </div>

            <div className="flex items-center gap-3">
                {actions}

                {onRefresh && (
                    <button
                        onClick={onRefresh}
                        disabled={loading}
                        className="flex items-center gap-2 px-4 py-2.5 bg-white border border-gray-200 rounded-xl text-gray-700 hover:bg-gray-50 transition-all shadow-sm disabled:opacity-50"
                    >
                        <ArrowPathIcon className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
                        <span className="hidden sm:inline">Actualiser</span>
                    </button>
                )}
            </div>
        </div>
    );
}
