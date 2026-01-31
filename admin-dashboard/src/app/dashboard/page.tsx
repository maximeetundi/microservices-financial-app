'use client';

import { useEffect, useState, useMemo } from 'react';
import { getDashboard } from '@/lib/api';
import Link from 'next/link';
import { format, subDays } from 'date-fns';
import { fr } from 'date-fns/locale';
import {
    AreaChart,
    Area,
    BarChart,
    Bar,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
    PieChart,
    Pie,
    Cell,
} from 'recharts';
import {
    UsersIcon,
    CreditCardIcon,
    ArrowsRightLeftIcon,
    WalletIcon,
    ExclamationTriangleIcon,
    ArrowTrendingUpIcon,
    ArrowTrendingDownIcon,
    ChevronRightIcon,
    ClockIcon,
    CheckCircleIcon,
    XCircleIcon,
    BanknotesIcon,
    ShieldCheckIcon,
    DocumentCheckIcon,
    UserPlusIcon,
    ArrowPathIcon,
    EyeIcon,
    Cog6ToothIcon,
    ChatBubbleLeftRightIcon,
    DocumentTextIcon,
} from '@heroicons/react/24/outline';
import { CheckCircleIcon as CheckSolid, XCircleIcon as XSolid } from '@heroicons/react/24/solid';
import PageHeader from '@/components/ui/PageHeader';
import StatsCard from '@/components/ui/StatsCard';

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
    // Trends
    users_change?: number;
    volume_change?: number;
}

interface Activity {
    id: string;
    type: 'user' | 'transaction' | 'kyc' | 'card' | 'alert';
    action: string;
    details: string;
    timestamp: string;
    status?: 'success' | 'warning' | 'error' | 'info';
}

// Mock data for charts (would come from API in real app)
const generateVolumeData = () => {
    return Array.from({ length: 7 }, (_, i) => ({
        date: format(subDays(new Date(), 6 - i), 'EEE', { locale: fr }),
        volume: Math.floor(Math.random() * 50000) + 10000,
        transactions: Math.floor(Math.random() * 100) + 20,
    }));
};

const generateUsersData = () => {
    return Array.from({ length: 7 }, (_, i) => ({
        date: format(subDays(new Date(), 6 - i), 'EEE', { locale: fr }),
        newUsers: Math.floor(Math.random() * 30) + 5,
        active: Math.floor(Math.random() * 200) + 50,
    }));
};

export default function DashboardPage() {
    const [stats, setStats] = useState<Stats | null>(null);
    const [loading, setLoading] = useState(true);
    const [refreshing, setRefreshing] = useState(false);
    const [recentActivities, setRecentActivities] = useState<Activity[]>([]);

    // Chart data
    const [volumeData] = useState(generateVolumeData);
    const [usersData] = useState(generateUsersData);

    const fetchData = async (silent = false) => {
        if (!silent) setLoading(true);
        else setRefreshing(true);

        try {
            const response = await getDashboard();
            setStats(response.data?.stats || null);

            // Generate mock activities (would come from API)
            setRecentActivities([
                { id: '1', type: 'user', action: 'Nouvel utilisateur', details: 'john.doe@email.com inscrit', timestamp: new Date().toISOString(), status: 'success' },
                { id: '2', type: 'kyc', action: 'KYC soumis', details: 'Document en attente de validation', timestamp: new Date(Date.now() - 3600000).toISOString(), status: 'warning' },
                { id: '3', type: 'transaction', action: 'Transfert effectué', details: '500€ envoyé à user@mail.com', timestamp: new Date(Date.now() - 7200000).toISOString(), status: 'success' },
                { id: '4', type: 'alert', action: 'Tentative suspecte', details: '3 échecs de connexion', timestamp: new Date(Date.now() - 10800000).toISOString(), status: 'error' },
                { id: '5', type: 'card', action: 'Nouvelle carte', details: 'Carte virtuelle créée', timestamp: new Date(Date.now() - 14400000).toISOString(), status: 'info' },
            ]);
        } catch (error) {
            console.error('Failed to fetch stats:', error);
        } finally {
            setLoading(false);
            setRefreshing(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, []);

    // KYC distribution for pie chart
    const kycData = useMemo(() => [
        { name: 'Vérifiés', value: stats?.verified_kyc || 0, color: '#22c55e' },
        { name: 'En attente', value: stats?.pending_kyc || 0, color: '#f59e0b' },
        { name: 'Non soumis', value: Math.max(0, (stats?.total_users || 0) - (stats?.verified_kyc || 0) - (stats?.pending_kyc || 0)), color: '#94a3b8' },
    ], [stats]);

    const formatCurrency = (amount: number) => {
        return new Intl.NumberFormat('fr-FR', {
            style: 'currency',
            currency: 'EUR',
            minimumFractionDigits: 0,
            maximumFractionDigits: 0,
        }).format(amount);
    };

    const getActivityIcon = (type: string, status?: string) => {
        const iconClass = "w-5 h-5";
        const icons: Record<string, React.ReactNode> = {
            user: <UsersIcon className={iconClass} />,
            transaction: <ArrowsRightLeftIcon className={iconClass} />,
            kyc: <DocumentCheckIcon className={iconClass} />,
            card: <CreditCardIcon className={iconClass} />,
            alert: <ExclamationTriangleIcon className={iconClass} />,
        };
        return icons[type] || <ClockIcon className={iconClass} />;
    };

    const getStatusColor = (status?: string) => {
        const colors: Record<string, string> = {
            success: 'bg-green-100 text-green-600',
            warning: 'bg-amber-100 text-amber-600',
            error: 'bg-red-100 text-red-600',
            info: 'bg-blue-100 text-blue-600',
        };
        return colors[status || 'info'] || colors.info;
    };

    if (loading) {
        return (
            <div className="flex items-center justify-center h-96">
                <div className="text-center">
                    <div className="w-16 h-16 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin mx-auto mb-4" />
                    <p className="text-gray-500">Chargement du tableau de bord...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="pb-10">
            <PageHeader
                title="Tableau de Bord"
                subtitle={`Dernière mise à jour: ${format(new Date(), 'dd MMMM yyyy, HH:mm', { locale: fr })}`}
                icon="📊"
                onRefresh={() => fetchData(true)}
                loading={refreshing}
            />

            {/* Alert Banner if pending KYC */}
            {stats?.pending_kyc && stats.pending_kyc > 0 && (
                <div className="mb-6 p-4 bg-gradient-to-r from-amber-50 to-orange-50 border border-amber-200 rounded-2xl flex items-center justify-between">
                    <div className="flex items-center gap-3">
                        <div className="p-2 bg-amber-100 rounded-xl">
                            <ExclamationTriangleIcon className="w-6 h-6 text-amber-600" />
                        </div>
                        <div>
                            <p className="font-semibold text-amber-900">
                                {stats.pending_kyc} demande{stats.pending_kyc > 1 ? 's' : ''} KYC en attente
                            </p>
                            <p className="text-sm text-amber-700">Vérifiez les documents soumis pour débloquer les comptes utilisateurs</p>
                        </div>
                    </div>
                    <Link
                        href="/dashboard/kyc"
                        className="flex items-center gap-2 px-4 py-2 bg-amber-600 text-white rounded-xl font-semibold hover:bg-amber-700 transition-colors"
                    >
                        Vérifier
                        <ChevronRightIcon className="w-4 h-4" />
                    </Link>
                </div>
            )}

            {/* Stats Grid */}
            <div className="grid grid-cols-2 lg:grid-cols-4 xl:grid-cols-5 gap-4 mb-8">
                <StatsCard
                    title="Utilisateurs"
                    value={stats?.total_users?.toLocaleString() || '0'}
                    icon={UsersIcon}
                    color="primary"
                    trend={stats?.users_change ? { value: stats.users_change, isPositive: stats.users_change > 0, label: 'cette semaine' } : undefined}
                />
                <StatsCard
                    title="Volume Aujourd'hui"
                    value={formatCurrency(stats?.volume_today || 0)}
                    icon={BanknotesIcon}
                    color="success"
                    trend={stats?.volume_change ? { value: stats.volume_change, isPositive: stats.volume_change > 0, label: 'vs hier' } : undefined}
                />
                <StatsCard
                    title="Transactions"
                    value={stats?.transactions_today?.toLocaleString() || '0'}
                    icon={ArrowsRightLeftIcon}
                    color="info"
                />
                <StatsCard
                    title="Cartes Actives"
                    value={stats?.active_cards?.toLocaleString() || '0'}
                    icon={CreditCardIcon}
                    color="purple"
                />
                <StatsCard
                    title="Nouveaux Inscrits"
                    value={stats?.new_users_today?.toLocaleString() || '0'}
                    icon={UserPlusIcon}
                    color="warning"
                />
            </div>

            {/* Charts Row */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
                {/* Volume Chart */}
                <div className="bg-white rounded-2xl border border-gray-200 shadow-sm p-6">
                    <div className="flex items-center justify-between mb-6">
                        <div>
                            <h3 className="text-lg font-bold text-gray-900">Volume des Transactions</h3>
                            <p className="text-sm text-gray-500">7 derniers jours</p>
                        </div>
                        <div className="flex items-center gap-2 text-sm">
                            <span className="flex items-center gap-1 text-green-600">
                                <ArrowTrendingUpIcon className="w-4 h-4" />
                                +12.5%
                            </span>
                        </div>
                    </div>
                    <div className="h-64">
                        <ResponsiveContainer width="100%" height="100%">
                            <AreaChart data={volumeData}>
                                <defs>
                                    <linearGradient id="volumeGradient" x1="0" y1="0" x2="0" y2="1">
                                        <stop offset="0%" stopColor="#6366f1" stopOpacity={0.3} />
                                        <stop offset="100%" stopColor="#6366f1" stopOpacity={0} />
                                    </linearGradient>
                                </defs>
                                <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
                                <XAxis dataKey="date" stroke="#9ca3af" fontSize={12} />
                                <YAxis stroke="#9ca3af" fontSize={12} tickFormatter={(v) => `${v / 1000}k`} />
                                <Tooltip
                                    contentStyle={{ borderRadius: '12px', border: '1px solid #e5e7eb' }}
                                    formatter={(value: number) => [formatCurrency(value), 'Volume']}
                                />
                                <Area
                                    type="monotone"
                                    dataKey="volume"
                                    stroke="#6366f1"
                                    strokeWidth={2}
                                    fill="url(#volumeGradient)"
                                />
                            </AreaChart>
                        </ResponsiveContainer>
                    </div>
                </div>

                {/* Users Chart */}
                <div className="bg-white rounded-2xl border border-gray-200 shadow-sm p-6">
                    <div className="flex items-center justify-between mb-6">
                        <div>
                            <h3 className="text-lg font-bold text-gray-900">Nouveaux Utilisateurs</h3>
                            <p className="text-sm text-gray-500">7 derniers jours</p>
                        </div>
                    </div>
                    <div className="h-64">
                        <ResponsiveContainer width="100%" height="100%">
                            <BarChart data={usersData}>
                                <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
                                <XAxis dataKey="date" stroke="#9ca3af" fontSize={12} />
                                <YAxis stroke="#9ca3af" fontSize={12} />
                                <Tooltip
                                    contentStyle={{ borderRadius: '12px', border: '1px solid #e5e7eb' }}
                                />
                                <Bar dataKey="newUsers" fill="#22c55e" radius={[4, 4, 0, 0]} name="Nouveaux" />
                            </BarChart>
                        </ResponsiveContainer>
                    </div>
                </div>
            </div>

            {/* Bottom Row */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                {/* Quick Actions */}
                <div className="bg-white rounded-2xl border border-gray-200 shadow-sm p-6">
                    <h3 className="text-lg font-bold text-gray-900 mb-4">Actions Rapides</h3>
                    <div className="grid grid-cols-2 gap-3">
                        {[
                            { href: '/dashboard/users', icon: UsersIcon, label: 'Utilisateurs', color: 'from-blue-500 to-indigo-600' },
                            { href: '/dashboard/kyc', icon: DocumentCheckIcon, label: 'KYC', color: 'from-amber-500 to-orange-600' },
                            { href: '/dashboard/transactions', icon: ArrowsRightLeftIcon, label: 'Transactions', color: 'from-green-500 to-emerald-600' },
                            { href: '/dashboard/cards', icon: CreditCardIcon, label: 'Cartes', color: 'from-purple-500 to-pink-600' },
                            { href: '/dashboard/support', icon: ChatBubbleLeftRightIcon, label: 'Support', color: 'from-cyan-500 to-blue-600' },
                            { href: '/dashboard/settings', icon: Cog6ToothIcon, label: 'Paramètres', color: 'from-gray-500 to-gray-700' },
                        ].map((action) => (
                            <Link
                                key={action.href}
                                href={action.href}
                                className="flex items-center gap-3 p-3 rounded-xl border border-gray-100 hover:border-gray-200 hover:shadow-md transition-all group"
                            >
                                <div className={`p-2 rounded-lg bg-gradient-to-br ${action.color} text-white group-hover:scale-110 transition-transform`}>
                                    <action.icon className="w-5 h-5" />
                                </div>
                                <span className="font-medium text-gray-700 group-hover:text-gray-900">{action.label}</span>
                            </Link>
                        ))}
                    </div>
                </div>

                {/* KYC Distribution */}
                <div className="bg-white rounded-2xl border border-gray-200 shadow-sm p-6">
                    <h3 className="text-lg font-bold text-gray-900 mb-4">Statut KYC</h3>
                    <div className="flex items-center">
                        <div className="h-48 flex-1">
                            <ResponsiveContainer width="100%" height="100%">
                                <PieChart>
                                    <Pie
                                        data={kycData}
                                        cx="50%"
                                        cy="50%"
                                        innerRadius={50}
                                        outerRadius={70}
                                        paddingAngle={3}
                                        dataKey="value"
                                    >
                                        {kycData.map((entry, index) => (
                                            <Cell key={`cell-${index}`} fill={entry.color} />
                                        ))}
                                    </Pie>
                                    <Tooltip
                                        contentStyle={{ borderRadius: '12px', border: '1px solid #e5e7eb' }}
                                    />
                                </PieChart>
                            </ResponsiveContainer>
                        </div>
                        <div className="space-y-3">
                            {kycData.map((item) => (
                                <div key={item.name} className="flex items-center gap-2">
                                    <div className="w-3 h-3 rounded-full" style={{ backgroundColor: item.color }} />
                                    <span className="text-sm text-gray-600">{item.name}</span>
                                    <span className="font-bold text-gray-900">{item.value}</span>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>

                {/* Recent Activity */}
                <div className="bg-white rounded-2xl border border-gray-200 shadow-sm p-6">
                    <div className="flex items-center justify-between mb-4">
                        <h3 className="text-lg font-bold text-gray-900">Activité Récente</h3>
                        <Link href="/dashboard/logs" className="text-sm text-indigo-600 hover:underline">
                            Voir tout
                        </Link>
                    </div>
                    <div className="space-y-3">
                        {recentActivities.map((activity) => (
                            <div key={activity.id} className="flex items-start gap-3 p-2 rounded-lg hover:bg-gray-50 transition-colors">
                                <div className={`p-2 rounded-lg ${getStatusColor(activity.status)}`}>
                                    {getActivityIcon(activity.type, activity.status)}
                                </div>
                                <div className="flex-1 min-w-0">
                                    <p className="font-medium text-gray-900 text-sm">{activity.action}</p>
                                    <p className="text-xs text-gray-500 truncate">{activity.details}</p>
                                </div>
                                <span className="text-xs text-gray-400 whitespace-nowrap">
                                    {format(new Date(activity.timestamp), 'HH:mm')}
                                </span>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
}
