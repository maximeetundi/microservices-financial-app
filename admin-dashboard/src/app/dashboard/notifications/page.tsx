'use client';

import { useState, useEffect } from 'react';
import {
    BellIcon,
    UserIcon,
    LockClosedIcon,
    ShieldCheckIcon,
    ExclamationTriangleIcon,
    CheckCircleIcon,
    XCircleIcon,
    ArrowPathIcon,
    FunnelIcon
} from '@heroicons/react/24/outline';
import api from '@/lib/api';
import { format } from 'date-fns';
import { fr } from 'date-fns/locale';

interface Notification {
    id: string;
    type: string;
    title: string;
    message: string;
    data?: string;
    is_read: boolean;
    created_at: string;
}

const typeConfig: Record<string, { label: string; color: string; icon: React.ComponentType<{ className?: string }> }> = {
    admin: { label: 'Admin', color: 'bg-purple-100 text-purple-800', icon: ShieldCheckIcon },
    account: { label: 'Compte', color: 'bg-blue-100 text-blue-800', icon: UserIcon },
    security: { label: 'Sécurité', color: 'bg-red-100 text-red-800', icon: ExclamationTriangleIcon },
    transaction: { label: 'Transaction', color: 'bg-green-100 text-green-800', icon: CheckCircleIcon },
    transfer: { label: 'Transfert', color: 'bg-emerald-100 text-emerald-800', icon: ArrowPathIcon },
    card: { label: 'Carte', color: 'bg-amber-100 text-amber-800', icon: LockClosedIcon },
    kyc: { label: 'KYC', color: 'bg-indigo-100 text-indigo-800', icon: ShieldCheckIcon },
};

export default function NotificationsPage() {
    const [notifications, setNotifications] = useState<Notification[]>([]);
    const [loading, setLoading] = useState(true);
    const [filter, setFilter] = useState<string>('all');
    const [searchTerm, setSearchTerm] = useState('');

    useEffect(() => {
        loadNotifications();
        // Poll for new notifications every 30 seconds
        const interval = setInterval(loadNotifications, 30000);
        return () => clearInterval(interval);
    }, []);

    const loadNotifications = async () => {
        try {
            setLoading(true);
            const response = await api.get('/notifications?limit=100');
            const data = response.data?.notifications || response.data || [];
            setNotifications(Array.isArray(data) ? data : []);
        } catch (error) {
            console.error('Error loading notifications:', error);
            // Use mock data for development
            setNotifications(getMockNotifications());
        } finally {
            setLoading(false);
        }
    };

    const getMockNotifications = (): Notification[] => {
        const now = new Date();
        return [
            {
                id: '1',
                type: 'admin',
                title: '👤 Nouveau compte créé',
                message: 'Un nouveau compte a été créé: john.doe@example.com',
                is_read: false,
                created_at: new Date(now.getTime() - 5 * 60000).toISOString(),
            },
            {
                id: '2',
                type: 'security',
                title: '🚫 Compte bloqué',
                message: 'Le compte user@example.com a été bloqué par Admin',
                is_read: false,
                created_at: new Date(now.getTime() - 15 * 60000).toISOString(),
            },
            {
                id: '3',
                type: 'admin',
                title: '🔓 PIN débloqué',
                message: 'Le PIN de jane@example.com a été débloqué',
                is_read: true,
                created_at: new Date(now.getTime() - 60 * 60000).toISOString(),
            },
            {
                id: '4',
                type: 'kyc',
                title: '✅ KYC Approuvé',
                message: 'KYC approuvé pour: client@example.com',
                is_read: true,
                created_at: new Date(now.getTime() - 2 * 60 * 60000).toISOString(),
            },
            {
                id: '5',
                type: 'transaction',
                title: '🛑 Transaction bloquée',
                message: 'Transaction suspecte de 5000 EUR bloquée',
                is_read: false,
                created_at: new Date(now.getTime() - 3 * 60 * 60000).toISOString(),
            },
        ];
    };

    const filteredNotifications = notifications.filter(n => {
        const matchesFilter = filter === 'all' || n.type === filter;
        const matchesSearch = searchTerm === '' ||
            n.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
            n.message.toLowerCase().includes(searchTerm.toLowerCase());
        return matchesFilter && matchesSearch;
    });

    const unreadCount = notifications.filter(n => !n.is_read).length;

    const getIcon = (type: string) => {
        const config = typeConfig[type] || typeConfig.admin;
        const IconComponent = config.icon;
        return <IconComponent className="w-5 h-5" />;
    };

    const formatTime = (dateString: string) => {
        try {
            return format(new Date(dateString), "dd MMM 'à' HH:mm", { locale: fr });
        } catch {
            return dateString;
        }
    };

    const getRelativeTime = (dateString: string) => {
        const now = new Date();
        const date = new Date(dateString);
        const diffMs = now.getTime() - date.getTime();
        const diffMins = Math.floor(diffMs / 60000);
        const diffHours = Math.floor(diffMs / 3600000);
        const diffDays = Math.floor(diffMs / 86400000);

        if (diffMins < 1) return "À l'instant";
        if (diffMins < 60) return `Il y a ${diffMins} min`;
        if (diffHours < 24) return `Il y a ${diffHours}h`;
        return `Il y a ${diffDays}j`;
    };

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <div>
                    <h1 className="text-2xl font-bold bg-gradient-to-r from-purple-600 to-indigo-600 bg-clip-text text-transparent">
                        Notifications
                    </h1>
                    <p className="text-gray-600 mt-1">
                        {unreadCount > 0
                            ? `${unreadCount} notification${unreadCount > 1 ? 's' : ''} non lue${unreadCount > 1 ? 's' : ''}`
                            : 'Toutes les notifications sont lues'}
                    </p>
                </div>
                <button
                    onClick={loadNotifications}
                    className="flex items-center gap-2 px-4 py-2 bg-white border border-gray-200 rounded-xl hover:bg-gray-50 transition-colors"
                >
                    <ArrowPathIcon className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
                    Actualiser
                </button>
            </div>

            {/* Filters */}
            <div className="flex flex-col sm:flex-row gap-4">
                <div className="relative flex-1">
                    <input
                        type="text"
                        placeholder="Rechercher dans les notifications..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="w-full px-4 py-3 pl-10 border border-gray-200 rounded-xl focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all"
                    />
                    <FunnelIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                </div>
                <div className="flex gap-2 overflow-x-auto pb-2 sm:pb-0">
                    <button
                        onClick={() => setFilter('all')}
                        className={`px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap transition-all ${filter === 'all'
                                ? 'bg-purple-600 text-white'
                                : 'bg-white text-gray-600 hover:bg-gray-50 border border-gray-200'
                            }`}
                    >
                        Tous
                    </button>
                    {Object.entries(typeConfig).map(([key, { label }]) => (
                        <button
                            key={key}
                            onClick={() => setFilter(key)}
                            className={`px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap transition-all ${filter === key
                                    ? 'bg-purple-600 text-white'
                                    : 'bg-white text-gray-600 hover:bg-gray-50 border border-gray-200'
                                }`}
                        >
                            {label}
                        </button>
                    ))}
                </div>
            </div>

            {/* Notifications List */}
            <div className="bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden">
                {loading ? (
                    <div className="p-12 text-center">
                        <ArrowPathIcon className="w-8 h-8 mx-auto text-gray-400 animate-spin" />
                        <p className="mt-2 text-gray-500">Chargement...</p>
                    </div>
                ) : filteredNotifications.length === 0 ? (
                    <div className="p-12 text-center">
                        <BellIcon className="w-12 h-12 mx-auto text-gray-300" />
                        <p className="mt-2 text-gray-500">Aucune notification</p>
                    </div>
                ) : (
                    <div className="divide-y divide-gray-100">
                        {filteredNotifications.map((notification) => {
                            const config = typeConfig[notification.type] || typeConfig.admin;
                            return (
                                <div
                                    key={notification.id}
                                    className={`p-4 hover:bg-gray-50 transition-colors cursor-pointer ${!notification.is_read ? 'bg-purple-50/50' : ''
                                        }`}
                                >
                                    <div className="flex items-start gap-4">
                                        {/* Icon */}
                                        <div className={`p-2 rounded-xl ${config.color}`}>
                                            {getIcon(notification.type)}
                                        </div>

                                        {/* Content */}
                                        <div className="flex-1 min-w-0">
                                            <div className="flex items-start justify-between gap-2">
                                                <h3 className={`font-medium ${!notification.is_read ? 'text-gray-900' : 'text-gray-700'}`}>
                                                    {notification.title}
                                                </h3>
                                                <span className="text-xs text-gray-400 whitespace-nowrap">
                                                    {getRelativeTime(notification.created_at)}
                                                </span>
                                            </div>
                                            <p className="text-sm text-gray-600 mt-1">
                                                {notification.message}
                                            </p>
                                            <div className="flex items-center gap-2 mt-2">
                                                <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${config.color}`}>
                                                    {config.label}
                                                </span>
                                                <span className="text-xs text-gray-400">
                                                    {formatTime(notification.created_at)}
                                                </span>
                                                {!notification.is_read && (
                                                    <span className="w-2 h-2 rounded-full bg-purple-600"></span>
                                                )}
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                )}
            </div>
        </div>
    );
}
