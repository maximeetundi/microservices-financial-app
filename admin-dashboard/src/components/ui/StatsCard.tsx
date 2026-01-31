'use client';

import { ArrowTrendingUpIcon, ArrowTrendingDownIcon } from '@heroicons/react/24/solid';

interface StatsCardProps {
    title: string;
    value: string | number;
    icon: React.ComponentType<{ className?: string }>;
    trend?: {
        value: number;
        isPositive: boolean;
        label?: string;
    };
    color?: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'purple';
    className?: string;
}

const colorClasses = {
    primary: 'from-indigo-500 to-purple-600',
    success: 'from-green-500 to-emerald-600',
    warning: 'from-amber-500 to-orange-600',
    danger: 'from-red-500 to-pink-600',
    info: 'from-blue-500 to-cyan-600',
    purple: 'from-purple-500 to-pink-600',
};

export default function StatsCard({
    title,
    value,
    icon: Icon,
    trend,
    color = 'primary',
    className = '',
}: StatsCardProps) {
    return (
        <div
            className={`relative overflow-hidden bg-gradient-to-br ${colorClasses[color]} text-white rounded-2xl p-5 shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-1 ${className}`}
        >
            {/* Background decoration */}
            <div className="absolute top-0 right-0 w-32 h-32 bg-white/10 rounded-full -translate-y-1/2 translate-x-1/2" />
            <div className="absolute bottom-0 left-0 w-24 h-24 bg-white/5 rounded-full translate-y-1/2 -translate-x-1/2" />

            <div className="relative flex items-start justify-between">
                <div>
                    <p className="text-white/70 text-sm font-medium">{title}</p>
                    <p className="text-3xl font-bold mt-1">{value}</p>

                    {trend && (
                        <div className="flex items-center gap-1 mt-2">
                            {trend.isPositive ? (
                                <ArrowTrendingUpIcon className="w-4 h-4 text-green-300" />
                            ) : (
                                <ArrowTrendingDownIcon className="w-4 h-4 text-red-300" />
                            )}
                            <span className={`text-sm font-medium ${trend.isPositive ? 'text-green-300' : 'text-red-300'}`}>
                                {trend.isPositive ? '+' : ''}{trend.value}%
                            </span>
                            {trend.label && (
                                <span className="text-white/50 text-xs">{trend.label}</span>
                            )}
                        </div>
                    )}
                </div>

                <div className="p-3 bg-white/20 rounded-xl backdrop-blur-sm">
                    <Icon className="w-6 h-6" />
                </div>
            </div>
        </div>
    );
}

// Mini variant for inline stats
export function MiniStatsCard({
    title,
    value,
    icon: Icon,
    color = 'primary',
}: {
    title: string;
    value: string | number;
    icon: React.ComponentType<{ className?: string }>;
    color?: 'primary' | 'success' | 'warning' | 'danger' | 'info';
}) {
    return (
        <div className={`flex items-center gap-3 px-4 py-2.5 bg-gradient-to-r ${colorClasses[color]} text-white rounded-xl shadow-md`}>
            <Icon className="w-5 h-5" />
            <span className="font-bold">{value}</span>
            <span className="text-white/70 text-sm">{title}</span>
        </div>
    );
}
