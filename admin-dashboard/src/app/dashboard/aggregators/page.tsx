'use client';

import { useState, useEffect } from 'react';
import {
    AdjustmentsHorizontalIcon,
    CheckCircleIcon,
    XCircleIcon,
    ExclamationTriangleIcon,
    MagnifyingGlassIcon,
    ArrowPathIcon,
    TrashIcon,
    Cog6ToothIcon,
    XMarkIcon,
    CurrencyDollarIcon,
    ArrowDownTrayIcon,
    ArrowUpTrayIcon
} from '@heroicons/react/24/outline';
import Link from 'next/link';

interface Aggregator {
    id: string;
    name: string;
    display_name: string;
    provider_type: string;
    logo_url: string;
    is_active: boolean;
    is_demo_mode: boolean;
    deposit_enabled: boolean;
    withdraw_enabled: boolean;
    fee_percentage: number;
    fee_fixed: number;
    min_transaction: number;
    max_transaction: number;
    daily_limit: number;
    priority: number;
    countries: { country_code: string; country_name: string; currency: string; is_active: boolean }[];
}

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';

export default function AggregatorsPage() {
    const [aggregators, setAggregators] = useState<Aggregator[]>([]);
    const [loading, setLoading] = useState(true);
    const [searchQuery, setSearchQuery] = useState('');
    const [statusFilter, setStatusFilter] = useState<'all' | 'enabled' | 'disabled'>('all');
    const [settingsModal, setSettingsModal] = useState<Aggregator | null>(null);
    const [saving, setSaving] = useState(false);

    useEffect(() => {
        loadAggregators();
    }, []);

    const getToken = () => localStorage.getItem('admin_token');

    const loadAggregators = async () => {
        setLoading(true);
        try {
            const response = await fetch(`${API_URL}/api/v1/admin/payment-providers`, {
                headers: {
                    'Authorization': `Bearer ${getToken()}`,
                    'Content-Type': 'application/json',
                },
            });

            if (response.ok) {
                const data = await response.json();
                setAggregators(data.providers || []);
            }
        } catch (error) {
            console.error('Error loading aggregators:', error);
        } finally {
            setLoading(false);
        }
    };

    const toggleStatus = async (id: string, isActive: boolean) => {
        try {
            await fetch(`${API_URL}/api/v1/admin/payment-providers/${id}/toggle-status`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${getToken()}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ is_active: !isActive }),
            });
            loadAggregators();
        } catch (error) {
            console.error('Error toggling status:', error);
        }
    };

    const toggleDeposit = async (id: string, enabled: boolean) => {
        try {
            await fetch(`${API_URL}/api/v1/admin/payment-providers/${id}/toggle-deposit`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${getToken()}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ enabled: !enabled }),
            });
            loadAggregators();
        } catch (error) {
            console.error('Error toggling deposit:', error);
        }
    };

    const toggleWithdraw = async (id: string, enabled: boolean) => {
        try {
            await fetch(`${API_URL}/api/v1/admin/payment-providers/${id}/toggle-withdraw`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${getToken()}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ enabled: !enabled }),
            });
            loadAggregators();
        } catch (error) {
            console.error('Error toggling withdraw:', error);
        }
    };

    const deleteAggregator = async (id: string) => {
        if (!confirm('Êtes-vous sûr de vouloir supprimer cet agrégateur ? Cette action est irréversible.')) return;

        try {
            const response = await fetch(`${API_URL}/api/v1/admin/payment-providers/${id}`, {
                method: 'DELETE',
                headers: { 'Authorization': `Bearer ${getToken()}` },
            });

            if (response.ok) {
                loadAggregators();
            } else {
                alert('Erreur lors de la suppression');
            }
        } catch (error) {
            console.error('Error deleting aggregator:', error);
        }
    };

    const saveSettings = async () => {
        if (!settingsModal) return;
        setSaving(true);

        try {
            // Update limits
            await fetch(`${API_URL}/api/v1/admin/payment-providers/${settingsModal.id}/limits`, {
                method: 'PUT',
                headers: {
                    'Authorization': `Bearer ${getToken()}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    min_transaction: settingsModal.min_transaction,
                    max_transaction: settingsModal.max_transaction,
                    daily_limit: settingsModal.daily_limit,
                }),
            });

            // Update fees
            await fetch(`${API_URL}/api/v1/admin/payment-providers/${settingsModal.id}/fees`, {
                method: 'PUT',
                headers: {
                    'Authorization': `Bearer ${getToken()}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    fee_percentage: settingsModal.fee_percentage,
                    fee_fixed: settingsModal.fee_fixed,
                }),
            });

            // Update priority
            await fetch(`${API_URL}/api/v1/admin/payment-providers/${settingsModal.id}/priority`, {
                method: 'PUT',
                headers: {
                    'Authorization': `Bearer ${getToken()}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ priority: settingsModal.priority }),
            });

            setSettingsModal(null);
            loadAggregators();
        } catch (error) {
            console.error('Error saving settings:', error);
            alert('Erreur lors de la sauvegarde');
        } finally {
            setSaving(false);
        }
    };

    const filteredAggregators = aggregators.filter(agg => {
        const matchesSearch = agg.display_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
            agg.name.toLowerCase().includes(searchQuery.toLowerCase());
        const matchesStatus = statusFilter === 'all' ||
            (statusFilter === 'enabled' && agg.is_active) ||
            (statusFilter === 'disabled' && !agg.is_active);
        return matchesSearch && matchesStatus;
    });

    const enabledCount = aggregators.filter(a => a.is_active).length;
    const disabledCount = aggregators.filter(a => !a.is_active).length;

    const formatAmount = (amount: number) => {
        if (amount >= 1000000) return `${(amount / 1000000).toFixed(1)}M`;
        if (amount >= 1000) return `${(amount / 1000).toFixed(0)}K`;
        return amount.toString();
    };

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
                            className={`card hover:shadow-xl transition-all duration-300 ${!agg.is_active ? 'opacity-60' : ''}`}
                        >
                            {/* Header */}
                            <div className="flex items-start justify-between mb-4">
                                <div className="flex items-center gap-3">
                                    <div className="w-12 h-12 rounded-xl bg-gray-100 flex items-center justify-center overflow-hidden">
                                        {agg.logo_url ? (
                                            <img src={agg.logo_url} alt={agg.display_name} className="w-8 h-8 object-contain" />
                                        ) : (
                                            <span className="text-xl font-bold text-gray-400">
                                                {agg.display_name[0]}
                                            </span>
                                        )}
                                    </div>
                                    <div>
                                        <h3 className="font-semibold text-gray-900">{agg.display_name}</h3>
                                        <p className="text-xs text-gray-500 font-mono">{agg.name}</p>
                                    </div>
                                </div>

                                <div className="flex items-center gap-1">
                                    {/* Settings */}
                                    <button
                                        onClick={() => setSettingsModal({ ...agg })}
                                        className="p-2 text-gray-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors"
                                        title="Configurer"
                                    >
                                        <Cog6ToothIcon className="w-5 h-5" />
                                    </button>

                                    {/* Instances */}
                                    <Link
                                        href={`/dashboard/aggregators/${agg.id}`}
                                        className="p-2 text-gray-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors"
                                        title="Gérer les instances"
                                    >
                                        <AdjustmentsHorizontalIcon className="w-5 h-5" />
                                    </Link>

                                    {/* Master Toggle */}
                                    <button
                                        onClick={() => toggleStatus(agg.id, agg.is_active)}
                                        className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${agg.is_active ? 'bg-emerald-500' : 'bg-gray-300'}`}
                                    >
                                        <span className={`inline-block h-4 w-4 transform rounded-full bg-white shadow-lg transition-transform ${agg.is_active ? 'translate-x-6' : 'translate-x-1'}`} />
                                    </button>

                                    {/* Delete (Demo only) */}
                                    {agg.name === 'demo' && (
                                        <button
                                            onClick={() => deleteAggregator(agg.id)}
                                            className="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                                            title="Supprimer"
                                        >
                                            <TrashIcon className="w-5 h-5" />
                                        </button>
                                    )}
                                </div>
                            </div>

                            {/* Demo Badge */}
                            {agg.is_demo_mode && (
                                <div className="flex items-center gap-2 p-2 mb-3 bg-amber-100 rounded-lg text-amber-700 text-sm">
                                    <ExclamationTriangleIcon className="w-4 h-4 flex-shrink-0" />
                                    <span>Mode Démo</span>
                                </div>
                            )}

                            {/* Deposit/Withdraw Toggles */}
                            <div className="space-y-3 mb-4">
                                <div className="flex items-center justify-between">
                                    <div className="flex items-center gap-2">
                                        <ArrowDownTrayIcon className="w-4 h-4 text-gray-500" />
                                        <span className="text-sm text-gray-600">Dépôts</span>
                                    </div>
                                    <button
                                        onClick={() => toggleDeposit(agg.id, agg.deposit_enabled)}
                                        disabled={!agg.is_active}
                                        className={`px-3 py-1 text-xs font-medium rounded-full transition-colors ${agg.deposit_enabled && agg.is_active
                                            ? 'bg-emerald-100 text-emerald-700 hover:bg-emerald-200'
                                            : 'bg-gray-100 text-gray-500 hover:bg-gray-200'
                                            }`}
                                    >
                                        {agg.deposit_enabled ? 'Activé' : 'Désactivé'}
                                    </button>
                                </div>

                                <div className="flex items-center justify-between">
                                    <div className="flex items-center gap-2">
                                        <ArrowUpTrayIcon className="w-4 h-4 text-gray-500" />
                                        <span className="text-sm text-gray-600">Retraits</span>
                                    </div>
                                    <button
                                        onClick={() => toggleWithdraw(agg.id, agg.withdraw_enabled)}
                                        disabled={!agg.is_active}
                                        className={`px-3 py-1 text-xs font-medium rounded-full transition-colors ${agg.withdraw_enabled && agg.is_active
                                            ? 'bg-emerald-100 text-emerald-700 hover:bg-emerald-200'
                                            : 'bg-gray-100 text-gray-500 hover:bg-gray-200'
                                            }`}
                                    >
                                        {agg.withdraw_enabled ? 'Activé' : 'Désactivé'}
                                    </button>
                                </div>
                            </div>

                            {/* Fees & Limits */}
                            <div className="pt-3 border-t border-gray-100 grid grid-cols-2 gap-3 text-xs">
                                <div>
                                    <span className="text-gray-500">Frais</span>
                                    <p className="font-semibold text-gray-900">
                                        {agg.fee_percentage > 0 ? `${agg.fee_percentage}%` : ''}
                                        {agg.fee_percentage > 0 && agg.fee_fixed > 0 ? ' + ' : ''}
                                        {agg.fee_fixed > 0 ? `${formatAmount(agg.fee_fixed)}` : ''}
                                        {agg.fee_percentage === 0 && agg.fee_fixed === 0 ? '0%' : ''}
                                    </p>
                                </div>
                                <div>
                                    <span className="text-gray-500">Min / Max</span>
                                    <p className="font-semibold text-gray-900">
                                        {formatAmount(agg.min_transaction)} - {agg.max_transaction > 0 ? formatAmount(agg.max_transaction) : '∞'}
                                    </p>
                                </div>
                            </div>

                            {/* Countries */}
                            <div className="mt-3 pt-3 border-t border-gray-100">
                                <span className="text-xs text-gray-500">Pays ({agg.countries?.length || 0})</span>
                                <div className="flex flex-wrap gap-1 mt-1">
                                    {(agg.countries || []).slice(0, 5).map((country) => (
                                        <span key={country.country_code} className={`px-2 py-0.5 text-xs rounded ${country.is_active ? 'bg-gray-100 text-gray-600' : 'bg-red-100 text-red-600'}`}>
                                            {country.country_code}
                                        </span>
                                    ))}
                                    {(agg.countries || []).length > 5 && (
                                        <span className="px-2 py-0.5 bg-indigo-100 text-indigo-600 text-xs rounded">
                                            +{agg.countries.length - 5}
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

            {/* Settings Modal */}
            {settingsModal && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
                    <div className="bg-white rounded-2xl shadow-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
                        <div className="p-6 border-b border-gray-100 flex items-center justify-between">
                            <div className="flex items-center gap-3">
                                <div className="w-10 h-10 rounded-xl bg-indigo-100 flex items-center justify-center">
                                    <Cog6ToothIcon className="w-5 h-5 text-indigo-600" />
                                </div>
                                <div>
                                    <h2 className="text-lg font-semibold text-gray-900">Configuration</h2>
                                    <p className="text-sm text-gray-500">{settingsModal.display_name}</p>
                                </div>
                            </div>
                            <button
                                onClick={() => setSettingsModal(null)}
                                className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
                            >
                                <XMarkIcon className="w-5 h-5 text-gray-500" />
                            </button>
                        </div>

                        <div className="p-6 space-y-6">
                            {/* Priority */}
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-2">
                                    Priorité (1-100)
                                </label>
                                <input
                                    type="number"
                                    min="1"
                                    max="100"
                                    value={settingsModal.priority}
                                    onChange={(e) => setSettingsModal({ ...settingsModal, priority: parseInt(e.target.value) || 50 })}
                                    className="input"
                                />
                                <p className="text-xs text-gray-500 mt-1">Plus la priorité est élevée, plus l'agrégateur sera privilégié</p>
                            </div>

                            {/* Fees */}
                            <div className="p-4 bg-gray-50 rounded-xl space-y-4">
                                <div className="flex items-center gap-2">
                                    <CurrencyDollarIcon className="w-5 h-5 text-gray-600" />
                                    <h3 className="font-medium text-gray-900">Frais</h3>
                                </div>
                                <div className="grid grid-cols-2 gap-4">
                                    <div>
                                        <label className="block text-sm text-gray-600 mb-1">Pourcentage (%)</label>
                                        <input
                                            type="number"
                                            step="0.01"
                                            min="0"
                                            max="100"
                                            value={settingsModal.fee_percentage}
                                            onChange={(e) => setSettingsModal({ ...settingsModal, fee_percentage: parseFloat(e.target.value) || 0 })}
                                            className="input"
                                        />
                                    </div>
                                    <div>
                                        <label className="block text-sm text-gray-600 mb-1">Fixe (XOF)</label>
                                        <input
                                            type="number"
                                            min="0"
                                            value={settingsModal.fee_fixed}
                                            onChange={(e) => setSettingsModal({ ...settingsModal, fee_fixed: parseFloat(e.target.value) || 0 })}
                                            className="input"
                                        />
                                    </div>
                                </div>
                            </div>

                            {/* Limits */}
                            <div className="p-4 bg-gray-50 rounded-xl space-y-4">
                                <h3 className="font-medium text-gray-900">Limites de transaction</h3>
                                <div className="grid grid-cols-2 gap-4">
                                    <div>
                                        <label className="block text-sm text-gray-600 mb-1">Minimum</label>
                                        <input
                                            type="number"
                                            min="0"
                                            value={settingsModal.min_transaction}
                                            onChange={(e) => setSettingsModal({ ...settingsModal, min_transaction: parseFloat(e.target.value) || 0 })}
                                            className="input"
                                        />
                                    </div>
                                    <div>
                                        <label className="block text-sm text-gray-600 mb-1">Maximum</label>
                                        <input
                                            type="number"
                                            min="0"
                                            value={settingsModal.max_transaction}
                                            onChange={(e) => setSettingsModal({ ...settingsModal, max_transaction: parseFloat(e.target.value) || 0 })}
                                            className="input"
                                        />
                                    </div>
                                </div>
                                <div>
                                    <label className="block text-sm text-gray-600 mb-1">Limite journalière par utilisateur</label>
                                    <input
                                        type="number"
                                        min="0"
                                        value={settingsModal.daily_limit}
                                        onChange={(e) => setSettingsModal({ ...settingsModal, daily_limit: parseFloat(e.target.value) || 0 })}
                                        className="input"
                                    />
                                </div>
                            </div>
                        </div>

                        <div className="p-6 border-t border-gray-100 flex justify-end gap-3">
                            <button
                                onClick={() => setSettingsModal(null)}
                                className="btn-secondary"
                            >
                                Annuler
                            </button>
                            <button
                                onClick={saveSettings}
                                disabled={saving}
                                className="btn-primary"
                            >
                                {saving ? 'Enregistrement...' : 'Enregistrer'}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
