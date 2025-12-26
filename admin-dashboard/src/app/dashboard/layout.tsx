'use client';

import { useState } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
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
} from '@heroicons/react/24/outline';
import { logout } from '@/lib/api';
import { useAuthStore } from '@/stores/authStore';
import clsx from 'clsx';

interface SidebarProps {
    children: React.ReactNode;
}

const navigation = [
    { name: 'Dashboard', href: '/dashboard', icon: HomeIcon },
    { name: 'Utilisateurs', href: '/dashboard/users', icon: UsersIcon },
    { name: 'Transactions', href: '/dashboard/transactions', icon: ArrowsRightLeftIcon },
    { name: 'Cartes', href: '/dashboard/cards', icon: CreditCardIcon },
    { name: 'Portefeuilles', href: '/dashboard/wallets', icon: WalletIcon },
    { name: 'Notifications', href: '/dashboard/notifications', icon: BellIcon },
    { name: 'KYC', href: '/dashboard/kyc', icon: ShieldCheckIcon },
    { name: 'Support', href: '/dashboard/support', icon: ChatBubbleLeftRightIcon },
    { name: 'Administrateurs', href: '/dashboard/admins', icon: Cog6ToothIcon },
    { name: 'Logs d\'audit', href: '/dashboard/logs', icon: DocumentTextIcon },
];


export default function DashboardLayout({ children }: SidebarProps) {
    const pathname = usePathname();
    const { admin } = useAuthStore();
    const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

    const closeMobileMenu = () => setMobileMenuOpen(false);

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
                        <SparklesIcon className="w-6 h-6 text-white" />
                    </div>
                    <span className="text-white font-bold text-lg">Zekora Admin</span>
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

            {/* Sidebar */}
            <aside className={clsx(
                'fixed lg:static inset-y-0 left-0 z-50 w-72 flex flex-col transform transition-all duration-300 ease-out',
                'bg-gradient-to-b from-slate-900 via-slate-800 to-slate-900',
                mobileMenuOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
            )}>
                {/* Logo */}
                <div className="p-6 border-b border-white/10">
                    <div className="flex items-center gap-3">
                        <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center shadow-lg shadow-indigo-500/30">
                            <SparklesIcon className="w-7 h-7 text-white" />
                        </div>
                        <div>
                            <h1 className="text-xl font-bold text-white">Admin Panel</h1>
                            <p className="text-indigo-300 text-sm font-medium">Zekora</p>
                        </div>
                    </div>
                </div>

                {/* Navigation */}
                <nav className="flex-1 px-4 py-6 space-y-1.5 overflow-y-auto">
                    {navigation.map((item) => {
                        const isActive = pathname === item.href || pathname?.startsWith(item.href + '/');
                        return (
                            <Link
                                key={item.name}
                                href={item.href}
                                onClick={closeMobileMenu}
                                className={clsx(
                                    'flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium transition-all duration-200',
                                    isActive
                                        ? 'bg-gradient-to-r from-indigo-500 to-purple-600 text-white shadow-lg shadow-indigo-500/30'
                                        : 'text-slate-300 hover:bg-white/10 hover:text-white'
                                )}
                            >
                                <item.icon className={clsx(
                                    'w-5 h-5 transition-transform',
                                    isActive && 'scale-110'
                                )} />
                                {item.name}
                                {isActive && (
                                    <div className="ml-auto w-1.5 h-1.5 rounded-full bg-white animate-pulse" />
                                )}
                            </Link>
                        );
                    })}
                </nav>

                {/* User Profile */}
                <div className="p-4 border-t border-white/10 bg-black/20">
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

            {/* Main content */}
            <main className="flex-1 lg:ml-0 mt-16 lg:mt-0 p-6 lg:p-8 overflow-auto">
                <div className="max-w-7xl mx-auto">
                    {children}
                </div>
            </main>
        </div>
    );
}

