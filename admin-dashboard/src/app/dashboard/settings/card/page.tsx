'use client';

import { useState, useEffect } from 'react';
import {
    CreditCardIcon,
    ArrowPathIcon,
    PencilSquareIcon,
    XMarkIcon,
    CheckCircleIcon,
} from '@heroicons/react/24/outline';
import { getFeeConfigs, updateFeeConfig } from '@/lib/api';

interface FeeConfig {
    id: string;
    key: string;
    name: string;
    description: string;
    type: string;
    fixed_amount: number;
    percentage_amount: number;
    currency: string;
    is_enabled: boolean;
    updated_at: string;
}

// Define which config keys belong to card service
const SERVICE_CONFIG_PREFIXES = [
    'card_',
    'fee_card_',
];

const FRIENDLY_NAMES: Record<string, string> = {
    'card_creation_fee': 'Frais Création Carte',
    'card_monthly_fee': 'Frais Mensuels',
    'card_annual_fee': 'Frais Annuels',
    'card_load_fee': 'Frais Rechargement',
    'card_atm_withdrawal_fee': 'Frais Retrait DAB',
    'card_foreign_transaction_fee': 'Frais Transaction Étrangère',
    'card_replacement_fee': 'Frais Remplacement',
    'card_inactive_fee': 'Frais Inactivité',
    'card_limit_daily': 'Limite Quotidienne',
    'card_limit_monthly': 'Limite Mensuelle',
    'card_limit_atm': 'Limite Retrait DAB',
    'fee_card_virtual': 'Frais Carte Virtuelle',
    'fee_card_physical': 'Frais Carte Physique',
    'fee_card_premium': 'Frais Carte Premium',
};

function getFriendlyName(key: string, defaultName: string): string {
    return FRIENDLY_NAMES[key] || defaultName;
}

function matchesService(key: string): boolean {
    return SERVICE_CONFIG_PREFIXES.some(prefix => key.startsWith(prefix));
}

function getCategory(key: string): string {
    if (key.includes('limit')) return 'Limites';
    if (key.includes('fee') || key.includes('Fee')) return 'Frais';
    return 'Paramètres';
}

export default function CardSettingsPage() {
    const [configs, setConfigs] = useState<FeeConfig[]>([]);
    const [loading, setLoading] = useState(true);
    const [editingConfig, setEditingConfig] = useState<FeeConfig | null>(null);
    const [formData, setFormData] = useState({
        type: '',
        fixed_amount: 0,
        percentage_amount: 0,
        is_enabled: true,
    });
    const [saving, setSaving] = useState(false);
    const [successMessage, setSuccessMessage] = useState('');

    const loadConfigs = async () => {
        try {
            setLoading(true);
            const response = await getFeeConfigs();
            const allConfigs = response.configurations || [];

            // Filter only card-related configs
            const cardConfigs = allConfigs.filter((c: FeeConfig) => matchesService(c.key));
            setConfigs(cardConfigs);
        } catch (error) {
            console.error('Failed to load configs:', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        loadConfigs();
    }, []);

    const handleEdit = (config: FeeConfig) => {
        setEditingConfig(config);
        setFormData({
            type: config.type,
            fixed_amount: config.fixed_amount,
            percentage_amount: config.percentage_amount,
            is_enabled: config.is_enabled,
        });
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!editingConfig) return;

        try {
            setSaving(true);
            await updateFeeConfig(editingConfig.key, formData);
            setSuccessMessage('Configuration mise à jour avec succès');
            setEditingConfig(null);
            loadConfigs();
            setTimeout(() => setSuccessMessage(''), 3000);
        } catch (error) {
            console.error('Failed to update config:', error);
        } finally {
            setSaving(false);
        }
    };

    // Group configs by category
    const groupedConfigs = configs.reduce((acc, config) => {
        const category = getCategory(config.key);
        if (!acc[category]) acc[category] = [];
        acc[category].push(config);
        return acc;
    }, {} as Record<string, FeeConfig[]>);

    return (
        <div className="p-6 max-w-7xl mx-auto">
            {/* Header */}
            <div className="flex items-center justify-between mb-8">
                <div className="flex items-center gap-3">
                    <div className="p-3 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-xl shadow-lg">
                        <CreditCardIcon className="h-8 w-8 text-white" />
                    </div>
                    <div>
                        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                            Configuration Cartes
                        </h1>
                        <p className="text-gray-500 dark:text-gray-400">
                            Gérez les frais et limites des cartes
                        </p>
                    </div>
                </div>
                <button
                    onClick={loadConfigs}
                    className="flex items-center gap-2 px-4 py-2 bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
                >
                    <ArrowPathIcon className={`h-5 w-5 ${loading ? 'animate-spin' : ''}`} />
                    Actualiser
                </button>
            </div>

            {/* Success Message */}
            {successMessage && (
                <div className="mb-6 flex items-center gap-2 p-4 bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-400 rounded-lg">
                    <CheckCircleIcon className="h-5 w-5" />
                    {successMessage}
                </div>
            )}

            {/* Configs by Category */}
            {loading ? (
                <div className="flex items-center justify-center h-64">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
                </div>
            ) : configs.length === 0 ? (
                <div className="text-center py-12 bg-gray-50 dark:bg-gray-800 rounded-xl">
                    <CreditCardIcon className="h-12 w-12 mx-auto text-gray-400 mb-4" />
                    <p className="text-gray-500 dark:text-gray-400">
                        Aucune configuration de carte trouvée
                    </p>
                </div>
            ) : (
                <div className="space-y-8">
                    {Object.entries(groupedConfigs).map(([category, categoryConfigs]) => (
                        <div key={category}>
                            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
                                <span className="w-2 h-2 rounded-full bg-blue-500"></span>
                                {category}
                            </h2>
                            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                                {categoryConfigs.map((config) => (
                                    <div
                                        key={config.id}
                                        className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6 hover:shadow-lg transition-shadow"
                                    >
                                        <div className="flex items-start justify-between mb-4">
                                            <div>
                                                <h3 className="font-medium text-gray-900 dark:text-white">
                                                    {getFriendlyName(config.key, config.name)}
                                                </h3>
                                                <p className="text-sm text-gray-500 dark:text-gray-400 truncate max-w-[200px]">
                                                    {config.key}
                                                </p>
                                            </div>
                                            <span
                                                className={`px-2 py-1 rounded-full text-xs font-medium ${config.is_enabled
                                                        ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
                                                        : 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
                                                    }`}
                                            >
                                                {config.is_enabled ? 'Actif' : 'Inactif'}
                                            </span>
                                        </div>

                                        <div className="space-y-2 mb-4">
                                            <div className="flex justify-between text-sm">
                                                <span className="text-gray-500 dark:text-gray-400">Montant</span>
                                                <span className="font-medium text-gray-900 dark:text-white">
                                                    {config.fixed_amount.toLocaleString()} {config.currency}
                                                </span>
                                            </div>
                                            {config.percentage_amount > 0 && (
                                                <div className="flex justify-between text-sm">
                                                    <span className="text-gray-500 dark:text-gray-400">Pourcentage</span>
                                                    <span className="font-medium text-gray-900 dark:text-white">
                                                        {config.percentage_amount}%
                                                    </span>
                                                </div>
                                            )}
                                        </div>

                                        <button
                                            onClick={() => handleEdit(config)}
                                            className="w-full flex items-center justify-center gap-2 px-4 py-2 bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 rounded-lg hover:bg-blue-100 dark:hover:bg-blue-900/40 transition-colors"
                                        >
                                            <PencilSquareIcon className="h-4 w-4" />
                                            Modifier
                                        </button>
                                    </div>
                                ))}
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {/* Edit Modal */}
            {editingConfig && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
                    <div className="bg-white dark:bg-gray-800 rounded-2xl p-6 w-full max-w-md mx-4 shadow-2xl">
                        <div className="flex items-center justify-between mb-6">
                            <h2 className="text-xl font-bold text-gray-900 dark:text-white">
                                Modifier la configuration
                            </h2>
                            <button
                                onClick={() => setEditingConfig(null)}
                                className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg"
                            >
                                <XMarkIcon className="h-5 w-5 text-gray-500" />
                            </button>
                        </div>

                        <form onSubmit={handleSubmit} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                    Configuration
                                </label>
                                <p className="text-gray-900 dark:text-white font-medium">
                                    {getFriendlyName(editingConfig.key, editingConfig.name)}
                                </p>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                    Montant ({editingConfig.currency})
                                </label>
                                <input
                                    type="number"
                                    step="0.01"
                                    value={formData.fixed_amount}
                                    onChange={(e) => setFormData({ ...formData, fixed_amount: parseFloat(e.target.value) || 0 })}
                                    className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                    Pourcentage (%)
                                </label>
                                <input
                                    type="number"
                                    step="0.01"
                                    value={formData.percentage_amount}
                                    onChange={(e) => setFormData({ ...formData, percentage_amount: parseFloat(e.target.value) || 0 })}
                                    className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                />
                            </div>

                            <div className="flex items-center gap-3">
                                <input
                                    type="checkbox"
                                    id="is_enabled"
                                    checked={formData.is_enabled}
                                    onChange={(e) => setFormData({ ...formData, is_enabled: e.target.checked })}
                                    className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                                />
                                <label htmlFor="is_enabled" className="text-sm text-gray-700 dark:text-gray-300">
                                    Configuration active
                                </label>
                            </div>

                            <div className="flex gap-3 pt-4">
                                <button
                                    type="button"
                                    onClick={() => setEditingConfig(null)}
                                    className="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700"
                                >
                                    Annuler
                                </button>
                                <button
                                    type="submit"
                                    disabled={saving}
                                    className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
                                >
                                    {saving ? 'Enregistrement...' : 'Enregistrer'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}
