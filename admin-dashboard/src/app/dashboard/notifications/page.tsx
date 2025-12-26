'use client';

import { useState, useEffect } from 'react';
import {
    BellIcon,
    UserIcon,
    LockClosedIcon,
    ShieldCheckIcon,
    ExclamationTriangleIcon,
    CheckCircleIcon,
    ArrowPathIcon,
    FunnelIcon,
    EnvelopeOpenIcon,
    TrashIcon,
    CreditCardIcon,
    ArrowsRightLeftIcon,
} from '@heroicons/react/24/outline';
import { BellIcon as BellSolidIcon } from '@heroicons/react/24/solid';
import api from '@/lib/api';
import { format, isToday, isYesterday, parseISO } from 'date-fns';
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

const typeConfig: Record<string, { label: string; color: string; bgColor: string; icon: React.ComponentType<{ className?: string }> }> = {
    admin: { label: 'Admin', color: 'text-purple-600', bgColor: 'bg-purple-100', icon: ShieldCheckIcon },
    account: { label: 'Compte', color: 'text-blue-600', bgColor: 'bg-blue-100', icon: UserIcon },
    security: { label: 'Sécurité', color: 'text-red-600', bgColor: 'bg-red-100', icon: ExclamationTriangleIcon },
    transaction: { label: 'Transaction', color: 'text-green-600', bgColor: 'bg-green-100', icon: CheckCircleIcon },
    transfer: { label: 'Transfert', color: 'text-emerald-600', bgColor: 'bg-emerald-100', icon: ArrowsRightLeftIcon },
    card: { label: 'Carte', color: 'text-amber-600', bgColor: 'bg-amber-100', icon: CreditCardIcon },
    kyc: { label: 'KYC', color: 'text-indigo-600', bgColor: 'bg-indigo-100', icon: ShieldCheckIcon },
    pin: { label: 'PIN', color: 'text-orange-600', bgColor: 'bg-orange-100', icon: LockClosedIcon },
};

export default function NotificationsPage() {
    const [notifications, setNotifications] = useState<Notification[]>([]);
    const [loading, setLoading] = useState(true);
    const [filter, setFilter] = useState<string>('all');
    const [searchTerm, setSearchTerm] = useState('');
    const [markingAsRead, setMarkingAsRead] = useState(false);

    useEffect(() => {
        loadNotifications();
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
                type: 'pin',
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
                title: '🛑 Transaction suspecte',
                message: 'Transaction suspecte de 5000 EUR détectée',
                is_read: false,
                created_at: new Date(now.getTime() - 3 * 60 * 60000).toISOString(),
            },
            {
                id: '6',
                type: 'transfer',
                title: '💸 Transfert effectué',
                message: 'Transfert de 250 EUR vers +33612345678',
                is_read: true,
                created_at: new Date(now.getTime() - 24 * 60 * 60000).toISOString(),
            },
            {
                id: '7',
                type: 'card',
                title: '💳 Nouvelle carte demandée',
                message: 'Demande de carte virtuelle par user@example.com',
                is_read: true,
                created_at: new Date(now.getTime() - 48 * 60 * 60000).toISOString(),
            },
        ];
    };

    const handleMarkAllAsRead = async () => {
        setMarkingAsRead(true);
        try {
            await api.post('/notifications/mark-all-read');
            setNotifications(prev => prev.map(n => ({ ...n, is_read: true })));
        } catch (error) {
            // Fallback to local state update
            setNotifications(prev => prev.map(n => ({ ...n, is_read: true })));
        } finally {
            setMarkingAsRead(false);
        }
    };

    const filteredNotifications = notifications.filter(n => {
        const matchesFilter = filter === 'all' || n.type === filter;
        const matchesSearch = searchTerm === '' ||
            n.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
            n.message.toLowerCase().includes(searchTerm.toLowerCase());
        return matchesFilter && matchesSearch;
    });

    // Group notifications by date
    const groupedNotifications = filteredNotifications.reduce((groups, notification) => {
        const date = parseISO(notification.created_at);
        let groupKey: string;

        if (isToday(date)) {
            groupKey = "Aujourd'hui";
        } else if (isYesterday(date)) {
            groupKey = 'Hier';
        } else {
            groupKey = format(date, 'EEEE d MMMM', { locale: fr });
        }

        if (!groups[groupKey]) {
            groups[groupKey] = [];
        }
        groups[groupKey].push(notification);
        return groups;
    }, {} as Record<string, Notification[]>);

    const unreadCount = notifications.filter(n => !n.is_read).length;
    const todayCount = notifications.filter(n => isToday(parseISO(n.created_at))).length;

    const getIcon = (type: string) => {
        const config = typeConfig[type] || typeConfig.admin;
        const IconComponent = config.icon;
        return <IconComponent className={`w-5 h-5 ${config.color}`} />;
    };

    const formatTime = (dateString: string) => {
        try {
            return format(new Date(dateString), "HH:mm", { locale: fr });
        } catch {
            return '';
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
            {/* Header with Stats */}
            <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
                <div>
                    <h1 className="text-3xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
                        Centre de Notifications
                    </h1>
                    <p className="text-gray-500 mt-1">
                        Suivez toutes les activités importantes en temps réel
                    </p>
                </div>

                {/* Stats Cards */}
                <div className="flex gap-4">
                    <div className="flex items-center gap-3 px-5 py-3 bg-white rounded-2xl shadow-lg border border-gray-100">
                        <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center">
                            <BellSolidIcon className="w-5 h-5 text-white" />
                        </div>
                        <div>
                            <p className="text-2xl font-bold text-gray-900">{notifications.length}</p>
                            <p className="text-xs text-gray-500">Total</p>
                        </div>
                    </div>
                    <div className="flex items-center gap-3 px-5 py-3 bg-white rounded-2xl shadow-lg border border-gray-100">
                        <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-red-500 to-pink-600 flex items-center justify-center">
                            <EnvelopeOpenIcon className="w-5 h-5 text-white" />
                        </div>
                        <div>
                            <p className="text-2xl font-bold text-gray-900">{unreadCount}</p>
                            <p className="text-xs text-gray-500">Non lues</p>
                        </div>
                    </div>
                    <div className="flex items-center gap-3 px-5 py-3 bg-white rounded-2xl shadow-lg border border-gray-100">
                        <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-green-500 to-emerald-600 flex items-center justify-center">
                            <CheckCircleIcon className="w-5 h-5 text-white" />
                        </div>
                        <div>
                            <p className="text-2xl font-bold text-gray-900">{todayCount}</p>
                            <p className="text-xs text-gray-500">Aujourd'hui</p>
                        </div>
                    </div>
                </div>
            </div>

            {/* Actions Bar */}
            <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
                {/* Search */}
                <div className="relative flex-1 max-w-md">
                    <input
                        type="text"
                        placeholder="Rechercher dans les notifications..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="w-full px-4 py-3 pl-11 bg-white border border-gray-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all shadow-sm"
                    />
                    <FunnelIcon className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                </div>

                {/* Actions */}
                <div className="flex gap-2">
                    {unreadCount > 0 && (
                        <button
                            onClick={handleMarkAllAsRead}
                            disabled={markingAsRead}
                            className="flex items-center gap-2 px-4 py-2.5 bg-gradient-to-r from-indigo-500 to-purple-600 text-white rounded-xl font-medium hover:shadow-lg hover:shadow-indigo-500/30 transition-all disabled:opacity-50"
                        >
                            <CheckCircleIcon className={`w-4 h-4 ${markingAsRead ? 'animate-spin' : ''}`} />
                            Tout marquer comme lu
                        </button>
                    )}
                    <button
                        onClick={loadNotifications}
                        className="flex items-center gap-2 px-4 py-2.5 bg-white border border-gray-200 rounded-xl hover:bg-gray-50 transition-colors shadow-sm"
                    >
                        <ArrowPathIcon className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
                        Actualiser
                    </button>
                </div>
            </div>

            {/* Filter Tabs */}
            <div className="flex gap-2 overflow-x-auto pb-2 scrollbar-hide">
                <button
                    onClick={() => setFilter('all')}
                    className={`px-4 py-2.5 rounded-xl text-sm font-medium whitespace-nowrap transition-all ${filter === 'all'
                        ? 'bg-gradient-to-r from-indigo-500 to-purple-600 text-white shadow-lg shadow-indigo-500/30'
                        : 'bg-white text-gray-600 hover:bg-gray-50 border border-gray-200 shadow-sm'
                        }`}
                >
                    Tous ({notifications.length})
                </button>
                {Object.entries(typeConfig).map(([key, { label, bgColor, color }]) => {
                    const count = notifications.filter(n => n.type === key).length;
                    if (count === 0) return null;
                    return (
                        <button
                            key={key}
                            onClick={() => setFilter(key)}
                            className={`px-4 py-2.5 rounded-xl text-sm font-medium whitespace-nowrap transition-all ${filter === key
                                ? 'bg-gradient-to-r from-indigo-500 to-purple-600 text-white shadow-lg shadow-indigo-500/30'
                                : `bg-white ${color} hover:bg-gray-50 border border-gray-200 shadow-sm`
                                }`}
                        >
                            {label} ({count})
                        </button>
                    );
                })}
            </div>

            {/* Notifications List - Grouped by Date */}
            <div className="space-y-6">
                {loading ? (
                    <div className="bg-white rounded-2xl shadow-lg border border-gray-100 p-12 text-center">
                        <ArrowPathIcon className="w-10 h-10 mx-auto text-indigo-500 animate-spin" />
                        <p className="mt-3 text-gray-500 font-medium">Chargement des notifications...</p>
                    </div>
                ) : Object.keys(groupedNotifications).length === 0 ? (
                    <div className="bg-white rounded-2xl shadow-lg border border-gray-100 p-12 text-center">
                        <div className="w-16 h-16 mx-auto rounded-2xl bg-gray-100 flex items-center justify-center mb-4">
                            <BellIcon className="w-8 h-8 text-gray-400" />
                        </div>
                        <p className="text-gray-500 font-medium">Aucune notification</p>
                        <p className="text-sm text-gray-400 mt-1">
                            {searchTerm ? 'Essayez une autre recherche' : 'Les notifications apparaîtront ici'}
                        </p>
                    </div>
                ) : (
                    Object.entries(groupedNotifications).map(([dateGroup, groupNotifications]) => (
                        <div key={dateGroup}>
                            {/* Date Header */}
                            <div className="flex items-center gap-3 mb-3">
                                <h3 className="text-sm font-semibold text-gray-500 uppercase tracking-wider">
                                    {dateGroup}
                                </h3>
                                <div className="flex-1 h-px bg-gray-200" />
                                <span className="text-xs text-gray-400 bg-gray-100 px-2 py-1 rounded-full">
                                    {groupNotifications.length} notification{groupNotifications.length > 1 ? 's' : ''}
                                </span>
                            </div>

                            {/* Notifications */}
                            <div className="bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden divide-y divide-gray-100">
                                {groupNotifications.map((notification) => {
                                    const config = typeConfig[notification.type] || typeConfig.admin;
                                    return (
                                        <div
                                            key={notification.id}
                                            className={`p-5 hover:bg-gray-50/50 transition-all cursor-pointer group ${!notification.is_read ? 'bg-indigo-50/30' : ''
                                                }`}
                                        >
                                            <div className="flex items-start gap-4">
                                                {/* Icon */}
                                                <div className={`p-3 rounded-xl ${config.bgColor} transition-transform group-hover:scale-110`}>
                                                    {getIcon(notification.type)}
                                                </div>

                                                {/* Content */}
                                                <div className="flex-1 min-w-0">
                                                    <div className="flex items-start justify-between gap-3">
                                                        <h4 className={`font-semibold ${!notification.is_read ? 'text-gray-900' : 'text-gray-700'}`}>
                                                            {notification.title}
                                                        </h4>
                                                        <div className="flex items-center gap-2 flex-shrink-0">
                                                            <span className="text-xs text-gray-400">
                                                                {formatTime(notification.created_at)}
                                                            </span>
                                                            {!notification.is_read && (
                                                                <span className="w-2.5 h-2.5 rounded-full bg-indigo-600 animate-pulse" />
                                                            )}
                                                        </div>
                                                    </div>
                                                    <p className="text-sm text-gray-600 mt-1">
                                                        {notification.message}
                                                    </p>
                                                    <div className="flex items-center gap-3 mt-3">
                                                        <span className={`px-2.5 py-1 rounded-lg text-xs font-medium ${config.bgColor} ${config.color}`}>
                                                            {config.label}
                                                        </span>
                                                        <span className="text-xs text-gray-400">
                                                            {getRelativeTime(notification.created_at)}
                                                        </span>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    );
                                })}
                            </div>
                        </div>
                    ))
                )}
            </div>
        </div>
    );
}
