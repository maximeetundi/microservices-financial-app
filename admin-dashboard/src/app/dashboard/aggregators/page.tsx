'use client';

import { useState, useEffect } from 'react';
import {
    AdjustmentsHorizontalIcon,
    CheckCircleIcon,
    XCircleIcon,
    ExclamationTriangleIcon,
    MagnifyingGlassIcon,
    ArrowPathIcon,
} from '@heroicons/react/24/outline';

interface Aggregator {
    id: string;
    provider_code: string;
    provider_name: string;
    logo_url: string;
    is_enabled: boolean;
    deposit_enabled: boolean;
    withdraw_enabled: boolean;
    supported_regions: string[];
    priority: number;
    min_amount: number;
    max_amount: number;
    fee_percent: number;
    fee_fixed: number;
    fee_currency: string;
    maintenance_mode: boolean;
    maintenance_msg: string;
}

export default function AggregatorsPage() {
    const [aggregators, setAggregators] = useState<Aggregator[]>([]);
    const [loading, setLoading] = useState(true);
    const [searchQuery, setSearchQuery] = useState('');
    const [statusFilter, setStatusFilter] = useState<'all' | 'enabled' | 'disabled'>('all');

    useEffect(() => {
        loadAggregators();
    }, []);

    const loadAggregators = async () => {
        setLoading(true);
        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');

            const response = await fetch(`${API_URL}/api/v1/admin/aggregators`, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
            });

            if (response.ok) {
                const data = await response.json();
                setAggregators(data.aggregators || []);
            }
        } catch (error) {
            console.error('Error loading aggregators:', error);
        } finally {
            setLoading(false);
        }
    };

    const toggleAggregator = async (code: string, enabled: boolean) => {
        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');
            const endpoint = enabled ? 'enable' : 'disable';

            await fetch(`${API_URL}/api/v1/admin/aggregators/${code}/${endpoint}`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
            });

            loadAggregators();
        } catch (error) {
            console.error('Error toggling aggregator:', error);
        }
    };

    const toggleDeposit = async (code: string, enabled: boolean) => {
        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');

            await fetch(`${API_URL}/api/v1/admin/aggregators/${code}/toggle-deposit`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ enabled }),
            });

            loadAggregators();
        } catch (error) {
            console.error('Error toggling deposit:', error);
        }
    };

    const toggleWithdraw = async (code: string, enabled: boolean) => {
        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');

            await fetch(`${API_URL}/api/v1/admin/aggregators/${code}/toggle-withdraw`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ enabled }),
            });

            loadAggregators();
        } catch (error) {
            console.error('Error toggling withdraw:', error);
        }
    };

    const filteredAggregators = aggregators.filter(agg => {
        const matchesSearch = agg.provider_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
            agg.provider_code.toLowerCase().includes(searchQuery.toLowerCase());
        const matchesStatus = statusFilter === 'all' ||
            (statusFilter === 'enabled' && agg.is_enabled) ||
            (statusFilter === 'disabled' && !agg.is_enabled);
        return matchesSearch && matchesStatus;
    });

    const enabledCount = aggregators.filter(a => a.is_enabled).length;
    const disabledCount = aggregators.filter(a => !a.is_enabled).length;

    return (
        <div className="space-y-6 animate-fadeIn">
            {/* Header */}
            <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900">Agrégateurs de Paiement</h1>
                    <p className="text-gray-500 mt-1">Gérez les fournisseurs de paiement disponibles sur la plateforme</p>
                </div>
                <button
                    onClick={loadAggregators}
                    className="btn-secondary flex items-center gap-2"
                >
                    <ArrowPathIcon className={`w-5 h-5 ${loading ? 'animate-spin' : ''}`} />
                    Actualiser
                </button>
            </div>

            {/* Stats Cards */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="stat-card flex items-center gap-4">
                    <div className="p-3 rounded-xl bg-indigo-100">
                        <AdjustmentsHorizontalIcon className="w-6 h-6 text-indigo-600" />
                    </div>
                    <div>
                        <p className="text-2xl font-bold text-gray-900">{aggregators.length}</p>
                        <p className="text-sm text-gray-500">Total Agrégateurs</p>
                    </div>
                </div>

                <div className="stat-card flex items-center gap-4">
                    <div className="p-3 rounded-xl bg-emerald-100">
                        <CheckCircleIcon className="w-6 h-6 text-emerald-600" />
                    </div>
                    <div>
                        <p className="text-2xl font-bold text-emerald-600">{enabledCount}</p>
                        <p className="text-sm text-gray-500">Actifs</p>
                    </div>
                </div>

                <div className="stat-card flex items-center gap-4">
                    <div className="p-3 rounded-xl bg-red-100">
                        <XCircleIcon className="w-6 h-6 text-red-600" />
                    </div>
                    <div>
                        <p className="text-2xl font-bold text-red-600">{disabledCount}</p>
                        <p className="text-sm text-gray-500">Désactivés</p>
                    </div>
                </div>
            </div>

            {/* Filters */}
            <div className="card">
                <div className="flex flex-col md:flex-row gap-4">
                    {/* Search */}
                    <div className="relative flex-1">
                        <MagnifyingGlassIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                        <input
                            type="text"
                            placeholder="Rechercher un agrégateur..."
                            className="input pl-10"
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                        />
                    </div>

                    {/* Status Filter */}
                    <select
                        className="select w-full md:w-48"
                        value={statusFilter}
                        onChange={(e) => setStatusFilter(e.target.value as any)}
                    >
                        <option value="all">Tous les statuts</option>
                        <option value="enabled">Actifs seulement</option>
                        <option value="disabled">Désactivés</option>
                    </select>
                </div>
            </div>

            {/* Aggregators Grid */}
            {loading ? (
                <div className="flex justify-center py-12">
                    <div className="spinner w-8 h-8" />
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {filteredAggregators.map((agg) => (
                        <div
                            key={agg.id}
                            className={`card hover:shadow-xl transition-all duration-300 ${!agg.is_enabled ? 'opacity-60' : ''
                                } ${agg.maintenance_mode ? 'border-amber-300 bg-amber-50/50' : ''}`}
                        >
                            {/* Header */}
                            <div className="flex items-start justify-between mb-4">
                                <div className="flex items-center gap-3">
                                    <div className="w-12 h-12 rounded-xl bg-gray-100 flex items-center justify-center overflow-hidden">
                                        {agg.logo_url ? (
                                            <img src={agg.logo_url} alt={agg.provider_name} className="w-8 h-8 object-contain" />
                                        ) : (
                                            <span className="text-xl font-bold text-gray-400">
                                                {agg.provider_name[0]}
                                            </span>
                                        )}
                                    </div>
                                    <div>
                                        <h3 className="font-semibold text-gray-900">{agg.provider_name}</h3>
                                        <p className="text-xs text-gray-500 font-mono">{agg.provider_code}</p>
                                    </div>
                                </div>

                                {/* Master Toggle */}
                                <button
                                    onClick={() => toggleAggregator(agg.provider_code, !agg.is_enabled)}
                                    className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${agg.is_enabled ? 'bg-emerald-500' : 'bg-gray-300'
                                        }`}
                                >
                                    <span
                                        className={`inline-block h-4 w-4 transform rounded-full bg-white shadow-lg transition-transform ${agg.is_enabled ? 'translate-x-6' : 'translate-x-1'
                                            }`}
                                    />
                                </button>
                            </div>

                            {/* Maintenance Warning */}
                            {agg.maintenance_mode && (
                                <div className="flex items-center gap-2 p-2 mb-3 bg-amber-100 rounded-lg text-amber-700 text-sm">
                                    <ExclamationTriangleIcon className="w-4 h-4 flex-shrink-0" />
                                    <span>En maintenance</span>
                                </div>
                            )}

                            {/* Feature Toggles */}
                            <div className="space-y-3 mb-4">
                                <div className="flex items-center justify-between">
                                    <span className="text-sm text-gray-600">Dépôts</span>
                                    <button
                                        onClick={() => toggleDeposit(agg.provider_code, !agg.deposit_enabled)}
                                        disabled={!agg.is_enabled}
                                        className={`px-3 py-1 text-xs font-medium rounded-full transition-colors ${agg.deposit_enabled && agg.is_enabled
                                                ? 'bg-emerald-100 text-emerald-700'
                                                : 'bg-gray-100 text-gray-500'
                                            }`}
                                    >
                                        {agg.deposit_enabled ? 'Activé' : 'Désactivé'}
                                    </button>
                                </div>

                                <div className="flex items-center justify-between">
                                    <span className="text-sm text-gray-600">Retraits</span>
                                    <button
                                        onClick={() => toggleWithdraw(agg.provider_code, !agg.withdraw_enabled)}
                                        disabled={!agg.is_enabled}
                                        className={`px-3 py-1 text-xs font-medium rounded-full transition-colors ${agg.withdraw_enabled && agg.is_enabled
                                                ? 'bg-emerald-100 text-emerald-700'
                                                : 'bg-gray-100 text-gray-500'
                                            }`}
                                    >
                                        {agg.withdraw_enabled ? 'Activé' : 'Désactivé'}
                                    </button>
                                </div>
                            </div>

                            {/* Stats */}
                            <div className="pt-3 border-t border-gray-100 grid grid-cols-2 gap-3 text-xs">
                                <div>
                                    <span className="text-gray-500">Frais</span>
                                    <p className="font-semibold text-gray-900">
                                        {agg.fee_percent}% {agg.fee_fixed > 0 ? `+ ${agg.fee_fixed} ${agg.fee_currency}` : ''}
                                    </p>
                                </div>
                                <div>
                                    <span className="text-gray-500">Min / Max</span>
                                    <p className="font-semibold text-gray-900">
                                        {agg.min_amount} - {agg.max_amount > 0 ? agg.max_amount : '∞'}
                                    </p>
                                </div>
                            </div>

                            {/* Regions */}
                            <div className="mt-3 pt-3 border-t border-gray-100">
                                <span className="text-xs text-gray-500">Régions</span>
                                <div className="flex flex-wrap gap-1 mt-1">
                                    {(agg.supported_regions || []).slice(0, 5).map((region) => (
                                        <span key={region} className="px-2 py-0.5 bg-gray-100 text-gray-600 text-xs rounded">
                                            {region}
                                        </span>
                                    ))}
                                    {(agg.supported_regions || []).length > 5 && (
                                        <span className="px-2 py-0.5 bg-indigo-100 text-indigo-600 text-xs rounded">
                                            +{agg.supported_regions.length - 5}
                                        </span>
                                    )}
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {/* Empty State */}
            {!loading && filteredAggregators.length === 0 && (
                <div className="text-center py-12">
                    <AdjustmentsHorizontalIcon className="w-12 h-12 mx-auto text-gray-300 mb-3" />
                    <p className="text-gray-500">Aucun agrégateur trouvé</p>
                </div>
            )}
        </div>
    );
}
