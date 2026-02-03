'use client';

import { useState, useEffect, useRef } from 'react';
import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import {
    HomeIcon,
    UsersIcon,
    CreditCardIcon,
    ArrowsRightLeftIcon,
    DocumentTextIcon,
    Cog6ToothIcon,
    ShieldCheckIcon,
    ArrowRightOnRectangleIcon,
    WalletIcon,
    ChatBubbleLeftRightIcon,
    Bars3Icon,
    XMarkIcon,
    SparklesIcon,
    BellIcon,
    CheckCircleIcon,
    ExclamationTriangleIcon,
    UserIcon,
    BuildingOfficeIcon,
    AdjustmentsHorizontalIcon,
    ShoppingBagIcon,
    MegaphoneIcon,
} from '@heroicons/react/24/outline';
import { BellIcon as BellSolidIcon } from '@heroicons/react/24/solid';
import { logout } from '@/lib/api';
import api from '@/lib/api';
import { useAuthStore } from '@/stores/authStore';
import clsx from 'clsx';

interface Notification {
    id: string;
    type: string;
    title: string;
    message: string;
    is_read: boolean;
    created_at: string;
}

const notificationTypeConfig: Record<string, { color: string; bgColor: string; icon: React.ComponentType<{ className?: string }> }> = {
    admin: { color: 'text-purple-600', bgColor: 'bg-purple-100', icon: ShieldCheckIcon },
    account: { color: 'text-blue-600', bgColor: 'bg-blue-100', icon: UserIcon },
    security: { color: 'text-red-600', bgColor: 'bg-red-100', icon: ExclamationTriangleIcon },
    transaction: { color: 'text-green-600', bgColor: 'bg-green-100', icon: CheckCircleIcon },
    transfer: { color: 'text-emerald-600', bgColor: 'bg-emerald-100', icon: ArrowsRightLeftIcon },
    card: { color: 'text-amber-600', bgColor: 'bg-amber-100', icon: CreditCardIcon },
    kyc: { color: 'text-indigo-600', bgColor: 'bg-indigo-100', icon: ShieldCheckIcon },
};

interface SidebarProps {
    children: React.ReactNode;
}

interface NavigationItem {
    name: string;
    href: string;
    icon: React.ComponentType<{ className?: string }>;
    badge?: string;
}

interface NavigationSection {
    title: string;
    items: NavigationItem[];
}

const navigationSections: NavigationSection[] = [
    {
        title: 'Tableau de bord',
        items: [
            { name: 'Vue d\'ensemble', href: '/dashboard', icon: HomeIcon },
            { name: 'Comptes Plateforme', href: '/dashboard/platform', icon: BuildingOfficeIcon },
        ]
    },
    {
        title: 'Opérations',
        items: [
            { name: 'Utilisateurs', href: '/dashboard/users', icon: UsersIcon },
            { name: 'Entreprises', href: '/dashboard/enterprises', icon: BuildingOfficeIcon },
            { name: 'Transactions', href: '/dashboard/transactions', icon: ArrowsRightLeftIcon },
            { name: 'Portefeuilles', href: '/dashboard/wallets', icon: WalletIcon },
            { name: 'Cartes', href: '/dashboard/cards', icon: CreditCardIcon },
        ]
    },
    {
        title: 'Paiements',
        items: [
            { name: 'Agrégateurs', href: '/dashboard/aggregators', icon: AdjustmentsHorizontalIcon, badge: 'Nouveau' },
            { name: 'Gestion Crédits', href: '/dashboard/credits', icon: SparklesIcon },
        ]
    },
    {
        title: 'Compliance',
        items: [
            { name: 'KYC', href: '/dashboard/kyc', icon: ShieldCheckIcon },
            { name: 'Associations', href: '/dashboard/associations', icon: UsersIcon },
        ]
    },
    {
        title: 'Support & Comm.',
        items: [
            { name: 'Notifications', href: '/dashboard/notifications', icon: BellIcon },
            { name: 'Support', href: '/dashboard/support', icon: ChatBubbleLeftRightIcon },
            { name: 'Événements', href: '/dashboard/events', icon: SparklesIcon },
        ]
    },
    {
        title: 'Configuration',
        items: [
            { name: 'Système', href: '/dashboard/settings', icon: Cog6ToothIcon },
            { name: 'Portefeuilles', href: '/dashboard/settings/wallet', icon: WalletIcon },
            { name: 'Transferts', href: '/dashboard/settings/transfer', icon: ArrowsRightLeftIcon },
            { name: 'Échange', href: '/dashboard/settings/exchange', icon: SparklesIcon },
            { name: 'Cartes', href: '/dashboard/settings/card', icon: CreditCardIcon },
            { name: 'Entreprises', href: '/dashboard/settings/enterprise', icon: BuildingOfficeIcon },
            { name: 'Boutiques', href: '/dashboard/settings/shops', icon: ShoppingBagIcon },
            { name: 'Notifications', href: '/dashboard/settings/notifications', icon: BellIcon },
            { name: 'Événements', href: '/dashboard/settings/events', icon: MegaphoneIcon },
        ]
    },
    {
        title: 'Administration',
        items: [
            { name: 'Administrateurs', href: '/dashboard/admins', icon: Cog6ToothIcon },
            { name: 'Logs d\'audit', href: '/dashboard/logs', icon: DocumentTextIcon },
        ]
    },
];


// Helper function to get relative time
function getRelativeTime(dateString: string): string {
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
}

export default function DashboardLayout({ children }: SidebarProps) {
    const pathname = usePathname();
    const router = useRouter();
    const { admin } = useAuthStore();
    const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

    // Notification state
    const [notifications, setNotifications] = useState<Notification[]>([]);
    const [notificationDropdownOpen, setNotificationDropdownOpen] = useState(false);
    const notificationRef = useRef<HTMLDivElement>(null);

    const closeMobileMenu = () => setMobileMenuOpen(false);

    // Handle notification click - navigate to relevant page
    const handleNotificationClick = (notification: Notification) => {
        setNotificationDropdownOpen(false);

        const type = notification.type?.toLowerCase() || '';

        if (type === 'support' || type === 'conversation' || type === 'ticket') {
            router.push('/dashboard/support');
        } else if (type === 'transfer' || type === 'transaction') {
            router.push('/dashboard/transactions');
        } else if (type === 'card') {
            router.push('/dashboard/cards');
        } else if (type === 'kyc') {
            router.push('/dashboard/kyc');
        } else if (type === 'wallet') {
            router.push('/dashboard/wallets');
        } else if (type === 'security' || type === 'admin') {
            router.push('/dashboard/logs');
        } else {
            router.push('/dashboard/notifications');
        }
    };

    // Load notifications from admin-service (admin-specific notifications)
    const loadNotifications = async () => {
        try {
            // Use admin-service API via Kong gateway
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = typeof window !== 'undefined' ? localStorage.getItem('admin_token') : null;

            if (!token) return;

            const response = await fetch(`${API_URL}/api/v1/admin/notifications?limit=10`, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
            });

            if (response.ok) {
                const data = await response.json();
                const notifList = data?.notifications || data || [];
                setNotifications(Array.isArray(notifList) ? notifList : []);
            }
        } catch (error) {
            console.error('Error loading admin notifications:', error);
        }
    };

    useEffect(() => {
        loadNotifications();
        const interval = setInterval(loadNotifications, 30000);
        return () => clearInterval(interval);
    }, []);

    // Close dropdown when clicking outside
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (notificationRef.current && !notificationRef.current.contains(event.target as Node)) {
                setNotificationDropdownOpen(false);
            }
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const unreadCount = notifications.filter(n => {
        // If we're on the support page, exclude support-related notifications from count
        if (pathname?.startsWith('/dashboard/support')) {
            const type = n.type?.toLowerCase() || '';
            if (type === 'support' || type === 'conversation' || type === 'ticket' || type === 'message') {
                return false;
            }
        }
        return !n.is_read;
    }).length;

    const getNotificationIcon = (type: string) => {
        const config = notificationTypeConfig[type] || notificationTypeConfig.admin;
        const IconComponent = config.icon;
        return <IconComponent className={`w-4 h-4 ${config.color}`} />;
    };

    return (
        <div className="flex min-h-screen bg-gradient-to-br from-slate-50 via-white to-slate-100">
            {/* Mobile menu overlay */}
            {mobileMenuOpen && (
                <div
                    className="fixed inset-0 bg-black/60 backdrop-blur-sm z-40 lg:hidden transition-opacity"
                    onClick={closeMobileMenu}
                />
            )}

            {/* Mobile header */}
            <div className="fixed top-0 left-0 right-0 h-16 bg-gradient-to-r from-slate-900 via-slate-800 to-slate-900 border-b border-slate-700/50 flex items-center justify-between px-4 z-30 lg:hidden shadow-lg">
                <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center shadow-lg">
                        <img src="/logo.png" alt="Zekora" className="w-6 h-6 object-contain" />
                    </div>
                    <span className="text-white font-bold text-lg">Zekora Admin</span>
                </div>
                <div className="flex items-center gap-2">
                    {/* Mobile Notification Bell */}
                    <div className="relative" ref={notificationRef}>
                        <button
                            onClick={() => setNotificationDropdownOpen(!notificationDropdownOpen)}
                            className="p-2.5 text-white hover:bg-white/10 rounded-xl transition-colors relative"
                        >
                            {unreadCount > 0 ? (
                                <BellSolidIcon className="w-6 h-6 text-amber-400" />
                            ) : (
                                <BellIcon className="w-6 h-6" />
                            )}
                            {unreadCount > 0 && (
                                <span className="absolute -top-1 -right-1 w-5 h-5 bg-red-500 text-white text-xs font-bold rounded-full flex items-center justify-center animate-pulse">
                                    {unreadCount > 9 ? '9+' : unreadCount}
                                </span>
                            )}
                        </button>
                    </div>
                    <button
                        onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
                        className="p-2.5 text-white hover:bg-white/10 rounded-xl transition-colors"
                    >
                        {mobileMenuOpen ? (
                            <XMarkIcon className="w-6 h-6" />
                        ) : (
                            <Bars3Icon className="w-6 h-6" />
                        )}
                    </button>
                </div>
            </div>

            {/* Sidebar */}
            <aside className={clsx(
                'fixed lg:sticky lg:top-0 inset-y-0 left-0 z-50 w-72 flex flex-col transform transition-all duration-300 ease-out h-screen',
                'bg-gradient-to-b from-slate-900 via-slate-800 to-slate-900',
                mobileMenuOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
            )}>
                {/* Logo */}
                <div className="flex-shrink-0 p-6 border-b border-white/10">
                    <div className="flex items-center gap-3">
                        <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center shadow-lg shadow-indigo-500/30">
                            <img src="/logo.png" alt="Zekora" className="w-7 h-7 object-contain" />
                        </div>
                        <div>
                            <h1 className="text-xl font-bold text-white">Admin Panel</h1>
                            <p className="text-indigo-300 text-sm font-medium">Zekora</p>
                        </div>
                    </div>
                </div>

                {/* Navigation - With proper scrolling */}
                <nav className="flex-1 min-h-0 px-3 py-4 overflow-y-auto space-y-6 scrollbar-thin scrollbar-track-slate-800 scrollbar-thumb-slate-600 hover:scrollbar-thumb-slate-500">
                    {navigationSections.map((section, sectionIndex) => (
                        <div key={section.title}>
                            {/* Section Title */}
                            <p className="px-4 mb-2 text-xs font-semibold text-slate-400/80 uppercase tracking-wider">
                                {section.title}
                            </p>

                            {/* Section Items */}
                            <div className="space-y-1">
                                {section.items.map((item) => {
                                    const isActive = pathname === item.href || pathname?.startsWith(item.href + '/');
                                    return (
                                        <Link
                                            key={item.name}
                                            href={item.href}
                                            onClick={closeMobileMenu}
                                            className={clsx(
                                                'group flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition-all duration-200',
                                                isActive
                                                    ? 'bg-gradient-to-r from-indigo-500 to-purple-600 text-white shadow-lg shadow-indigo-500/25'
                                                    : 'text-slate-300 hover:bg-white/10 hover:text-white'
                                            )}
                                        >
                                            <item.icon className={clsx(
                                                'w-5 h-5 transition-all duration-200',
                                                isActive ? 'scale-110' : 'group-hover:scale-105'
                                            )} />
                                            <span className="flex-1">{item.name}</span>

                                            {/* Badge for new features */}
                                            {item.badge && (
                                                <span className={clsx(
                                                    'px-2 py-0.5 text-[10px] font-bold rounded-full uppercase tracking-wide',
                                                    item.badge === 'Nouveau'
                                                        ? 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/30'
                                                        : item.badge === 'Dev'
                                                            ? 'bg-amber-500/20 text-amber-300 border border-amber-500/30'
                                                            : 'bg-indigo-500/20 text-indigo-300 border border-indigo-500/30'
                                                )}>
                                                    {item.badge}
                                                </span>
                                            )}

                                            {isActive && (
                                                <div className="w-1.5 h-1.5 rounded-full bg-white animate-pulse" />
                                            )}
                                        </Link>
                                    );
                                })}
                            </div>

                            {/* Section Divider (except last) */}
                            {sectionIndex < navigationSections.length - 1 && (
                                <div className="mt-4 mx-4 border-t border-white/5" />
                            )}
                        </div>
                    ))}
                </nav>

                {/* User Profile - Fixed at bottom */}
                <div className="flex-shrink-0 p-4 border-t border-white/10 bg-black/20">
                    <div className="flex items-center gap-3 px-3 py-3 rounded-xl bg-white/5">
                        <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-bold shadow-lg">
                            {admin?.first_name?.[0] || 'A'}
                        </div>
                        <div className="flex-1 min-w-0">
                            <p className="text-sm font-semibold text-white truncate">
                                {admin?.first_name} {admin?.last_name}
                            </p>
                            <p className="text-xs text-indigo-300 truncate">
                                {admin?.role?.name || 'Administrateur'}
                            </p>
                        </div>
                    </div>
                    <button
                        onClick={() => {
                            closeMobileMenu();
                            logout();
                        }}
                        className="flex items-center gap-3 w-full px-4 py-2.5 mt-3 text-sm text-slate-400 hover:text-white hover:bg-red-500/20 rounded-xl transition-all duration-200"
                    >
                        <ArrowRightOnRectangleIcon className="w-5 h-5" />
                        Déconnexion
                    </button>
                </div>
            </aside>

            {/* Main content wrapper */}
            <div className="flex-1 flex flex-col lg:ml-0">
                {/* Desktop Header with Notification Bell */}
                <header className="hidden lg:flex h-16 bg-white/80 backdrop-blur-xl border-b border-gray-200/50 items-center justify-end px-8 sticky top-0 z-20">
                    {/* Notification Bell */}
                    <div className="relative" ref={notificationRef}>
                        <button
                            onClick={() => setNotificationDropdownOpen(!notificationDropdownOpen)}
                            className="relative p-3 rounded-xl hover:bg-gray-100 transition-all duration-200 group"
                        >
                            {unreadCount > 0 ? (
                                <BellSolidIcon className="w-6 h-6 text-indigo-600 group-hover:scale-110 transition-transform" />
                            ) : (
                                <BellIcon className="w-6 h-6 text-gray-500 group-hover:text-indigo-600 group-hover:scale-110 transition-all" />
                            )}
                            {unreadCount > 0 && (
                                <span className="absolute -top-0.5 -right-0.5 w-5 h-5 bg-gradient-to-r from-red-500 to-pink-500 text-white text-xs font-bold rounded-full flex items-center justify-center shadow-lg shadow-red-500/30 animate-bounce">
                                    {unreadCount > 9 ? '9+' : unreadCount}
                                </span>
                            )}
                        </button>

                        {/* Notification Dropdown */}
                        {notificationDropdownOpen && (
                            <div className="absolute right-0 mt-2 w-96 bg-white/95 backdrop-blur-xl rounded-2xl shadow-2xl border border-gray-100 overflow-hidden z-50 animate-slide-up">
                                {/* Header */}
                                <div className="px-5 py-4 bg-gradient-to-r from-indigo-500 to-purple-600 text-white">
                                    <div className="flex items-center justify-between">
                                        <h3 className="font-bold text-lg">Notifications</h3>
                                        {unreadCount > 0 && (
                                            <span className="px-2.5 py-1 bg-white/20 rounded-full text-sm font-medium">
                                                {unreadCount} nouvelle{unreadCount > 1 ? 's' : ''}
                                            </span>
                                        )}
                                    </div>
                                </div>

                                {/* Notifications List */}
                                <div className="max-h-80 overflow-y-auto">
                                    {notifications.length === 0 ? (
                                        <div className="p-8 text-center">
                                            <BellIcon className="w-12 h-12 mx-auto text-gray-300 mb-3" />
                                            <p className="text-gray-500">Aucune notification</p>
                                        </div>
                                    ) : (
                                        <div className="divide-y divide-gray-100">
                                            {notifications.slice(0, 5).map((notification) => {
                                                const config = notificationTypeConfig[notification.type] || notificationTypeConfig.admin;
                                                return (
                                                    <div
                                                        key={notification.id}
                                                        onClick={() => handleNotificationClick(notification)}
                                                        className={clsx(
                                                            'px-5 py-4 hover:bg-gray-50 transition-colors cursor-pointer',
                                                            !notification.is_read && 'bg-indigo-50/50'
                                                        )}
                                                    >
                                                        <div className="flex items-start gap-3">
                                                            <div className={`p-2 rounded-xl ${config.bgColor} flex-shrink-0`}>
                                                                {getNotificationIcon(notification.type)}
                                                            </div>
                                                            <div className="flex-1 min-w-0">
                                                                <p className={clsx(
                                                                    'text-sm font-medium truncate',
                                                                    !notification.is_read ? 'text-gray-900' : 'text-gray-700'
                                                                )}>
                                                                    {notification.title}
                                                                </p>
                                                                <p className="text-xs text-gray-500 mt-0.5 line-clamp-2">
                                                                    {notification.message}
                                                                </p>
                                                                <p className="text-xs text-gray-400 mt-1">
                                                                    {getRelativeTime(notification.created_at)}
                                                                </p>
                                                            </div>
                                                            {!notification.is_read && (
                                                                <span className="w-2 h-2 rounded-full bg-indigo-600 flex-shrink-0 mt-2" />
                                                            )}
                                                        </div>
                                                    </div>
                                                );
                                            })}
                                        </div>
                                    )}
                                </div>

                                {/* Footer */}
                                <div className="px-5 py-3 bg-gray-50 border-t border-gray-100">
                                    <Link
                                        href="/dashboard/notifications"
                                        onClick={() => setNotificationDropdownOpen(false)}
                                        className="block w-full text-center py-2.5 text-sm font-semibold text-indigo-600 hover:text-indigo-700 hover:bg-indigo-50 rounded-xl transition-colors"
                                    >
                                        Voir toutes les notifications →
                                    </Link>
                                </div>
                            </div>
                        )}
                    </div>

                    {/* Admin info */}
                    <div className="flex items-center gap-3 ml-4 pl-4 border-l border-gray-200">
                        <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-bold text-sm shadow-md">
                            {admin?.first_name?.[0] || 'A'}
                        </div>
                        <div className="hidden xl:block">
                            <p className="text-sm font-semibold text-gray-900">
                                {admin?.first_name} {admin?.last_name}
                            </p>
                            <p className="text-xs text-gray-500">
                                {admin?.role?.name || 'Administrateur'}
                            </p>
                        </div>
                    </div>
                </header>

                {/* Main content */}
                <main className="flex-1 mt-16 lg:mt-0 p-6 lg:p-8 overflow-auto">
                    <div className="max-w-7xl mx-auto">
                        {children}
                    </div>
                </main>
            </div>
        </div>
    );
}
