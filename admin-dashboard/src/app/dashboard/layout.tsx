'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import {
    HomeIcon,
    UsersIcon,
    CreditCardIcon,
    BanknotesIcon,
    ArrowsRightLeftIcon,
    DocumentTextIcon,
    Cog6ToothIcon,
    ShieldCheckIcon,
    ArrowRightOnRectangleIcon,
    WalletIcon,
    ChatBubbleLeftRightIcon,
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
    { name: 'KYC', href: '/dashboard/kyc', icon: ShieldCheckIcon },
    { name: 'Support', href: '/dashboard/support', icon: ChatBubbleLeftRightIcon },
    { name: 'Administrateurs', href: '/dashboard/admins', icon: Cog6ToothIcon },
    { name: 'Logs d\'audit', href: '/dashboard/logs', icon: DocumentTextIcon },
];


export default function DashboardLayout({ children }: SidebarProps) {
    const pathname = usePathname();
    const { admin } = useAuthStore();

    return (
        <div className="flex min-h-screen bg-gray-50">
            {/* Sidebar */}
            <aside className="sidebar flex flex-col">
                <div className="p-6 border-b border-slate-700">
                    <h1 className="text-xl font-bold text-white">🏦 Admin Panel</h1>
                    <p className="text-slate-400 text-sm mt-1">Crypto Bank</p>
                </div>

                <nav className="flex-1 px-4 py-6 space-y-1">
                    {navigation.map((item) => {
                        const isActive = pathname === item.href || pathname?.startsWith(item.href + '/');
                        return (
                            <Link
                                key={item.name}
                                href={item.href}
                                className={clsx(
                                    'flex items-center gap-3 px-4 py-3 rounded-lg text-sm font-medium transition-colors',
                                    isActive
                                        ? 'bg-primary-600 text-white'
                                        : 'text-slate-300 hover:bg-slate-800 hover:text-white'
                                )}
                            >
                                <item.icon className="w-5 h-5" />
                                {item.name}
                            </Link>
                        );
                    })}
                </nav>

                <div className="p-4 border-t border-slate-700">
                    <div className="flex items-center gap-3 px-4 py-2">
                        <div className="w-8 h-8 rounded-full bg-primary-500 flex items-center justify-center text-white font-medium">
                            {admin?.first_name?.[0] || 'A'}
                        </div>
                        <div className="flex-1 min-w-0">
                            <p className="text-sm font-medium text-white truncate">
                                {admin?.first_name} {admin?.last_name}
                            </p>
                            <p className="text-xs text-slate-400 truncate">
                                {admin?.role?.name || 'Admin'}
                            </p>
                        </div>
                    </div>
                    <button
                        onClick={logout}
                        className="flex items-center gap-3 w-full px-4 py-2 mt-2 text-sm text-slate-300 hover:bg-slate-800 rounded-lg"
                    >
                        <ArrowRightOnRectangleIcon className="w-5 h-5" />
                        Déconnexion
                    </button>
                </div>
            </aside>

            {/* Main content */}
            <main className="main-content">
                {children}
            </main>
        </div>
    );
}
