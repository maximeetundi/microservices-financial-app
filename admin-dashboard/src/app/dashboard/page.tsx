'use client';

import { useEffect, useState } from 'react';
import { getDashboard } from '@/lib/api';
import {
    UsersIcon,
    CreditCardIcon,
    ArrowsRightLeftIcon,
    WalletIcon,
    ExclamationTriangleIcon,
    ChartBarIcon,
} from '@heroicons/react/24/outline';

interface Stats {
    total_users: number;
    active_today: number;
    transactions_today: number;
    volume_today: number;
    total_cards: number;
    active_cards: number;
    total_wallets: number;
    pending_kyc: number;
    verified_kyc: number;
    new_users_today: number;
    transfers_today: number;
}


export default function DashboardPage() {
    const [stats, setStats] = useState<Stats | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchStats = async () => {
            try {
                const response = await getDashboard();
                setStats(response.data.stats);
            } catch (error) {
                console.error('Failed to fetch stats:', error);
            } finally {
                setLoading(false);
            }
        };
        fetchStats();
    }, []);

    if (loading) {
        return (
            <div className="flex items-center justify-center h-64">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
            </div>
        );
    }

    const statCards = [
        {
            name: 'Utilisateurs',
            value: stats?.total_users || 0,
            icon: UsersIcon,
            color: 'bg-blue-500',
            subtext: `${stats?.active_today || 0} actifs aujourd'hui`,
        },
        {
            name: 'Transactions',
            value: stats?.transactions_today || 0,
            icon: ArrowsRightLeftIcon,
            color: 'bg-green-500',
            subtext: `${formatCurrency(stats?.volume_today || 0)} de volume`,
        },
        {
            name: 'Cartes',
            value: stats?.total_cards || 0,
            icon: CreditCardIcon,
            color: 'bg-purple-500',
            subtext: 'Cartes émises',
        },
        {
            name: 'Portefeuilles',
            value: stats?.total_wallets || 0,
            icon: WalletIcon,
            color: 'bg-orange-500',
            subtext: 'Portefeuilles actifs',
        },
        {
            name: 'KYC en attente',
            value: stats?.pending_kyc || 0,
            icon: ExclamationTriangleIcon,
            color: 'bg-yellow-500',
            subtext: 'Vérifications requises',
        },
    ];

    return (
        <div>
            <div className="mb-8">
                <h1 className="text-2xl font-bold text-slate-900">Tableau de bord</h1>
                <p className="text-slate-500 mt-1">Vue d'ensemble du système</p>
            </div>

            {/* Stats Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-6 mb-8">
                {statCards.map((stat) => (
                    <div key={stat.name} className="stat-card">
                        <div className="flex items-center gap-4">
                            <div className={`p-3 rounded-lg ${stat.color}`}>
                                <stat.icon className="w-6 h-6 text-white" />
                            </div>
                            <div>
                                <p className="text-2xl font-bold text-slate-900">{stat.value.toLocaleString()}</p>
                                <p className="text-sm text-slate-500">{stat.name}</p>
                            </div>
                        </div>
                        <p className="text-xs text-slate-400 mt-3">{stat.subtext}</p>
                    </div>
                ))}
            </div>

            {/* Quick Actions */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <div className="card">
                    <h2 className="text-lg font-semibold text-slate-900 mb-4">Actions rapides</h2>
                    <div className="grid grid-cols-2 gap-4">
                        <a href="/dashboard/users" className="btn-secondary text-center">
                            👥 Utilisateurs
                        </a>
                        <a href="/dashboard/kyc" className="btn-secondary text-center">
                            📋 Vérifier KYC
                        </a>
                        <a href="/dashboard/transactions" className="btn-secondary text-center">
                            💸 Transactions
                        </a>
                        <a href="/dashboard/cards" className="btn-secondary text-center">
                            💳 Cartes
                        </a>
                        <a href="/dashboard/support" className="btn-secondary text-center">
                            💬 Support
                        </a>
                        <a href="/dashboard/logs" className="btn-secondary text-center">
                            📝 Logs d&apos;audit
                        </a>
                    </div>
                </div>

                <div className="card">
                    <h2 className="text-lg font-semibold text-slate-900 mb-4">Alertes récentes</h2>
                    <div className="space-y-3">
                        {stats?.pending_kyc && stats.pending_kyc > 0 ? (
                            <div className="flex items-center gap-3 p-3 bg-yellow-50 rounded-lg">
                                <ExclamationTriangleIcon className="w-5 h-5 text-yellow-600" />
                                <span className="text-sm text-yellow-800">
                                    {stats.pending_kyc} demandes KYC en attente de validation
                                </span>
                            </div>
                        ) : (
                            <p className="text-sm text-slate-500">Aucune alerte</p>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}

function formatCurrency(amount: number): string {
    return new Intl.NumberFormat('fr-FR', {
        style: 'currency',
        currency: 'USD',
    }).format(amount);
}
