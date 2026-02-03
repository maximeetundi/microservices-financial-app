'use client';

import { useState, useEffect } from 'react';
import {
    BuildingOfficeIcon,
    ArrowPathIcon,
    PencilSquareIcon,
    XMarkIcon,
    CheckCircleIcon,
    UserGroupIcon,
    BanknotesIcon,
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

// Define which config keys belong to enterprise service
const SERVICE_CONFIG_PREFIXES = [
    'limit_daily_enterprise',
    'limit_monthly_enterprise',
    'limit_max_balance_enterprise',
    'enterprise_',
    'payroll_',
    'multi_sig_',
];

const FRIENDLY_NAMES: Record<string, string> = {
    'limit_daily_enterprise': 'Limite Quotidienne Entreprise',
    'limit_monthly_enterprise': 'Limite Mensuelle Entreprise',
    'limit_max_balance_enterprise': 'Solde Maximum Entreprise',
    'enterprise_wallet_fee': 'Frais Portefeuille Entreprise',
    'enterprise_transfer_fee': 'Frais Transfert Entreprise',
    'enterprise_payroll_fee': 'Frais Paie Entreprise',
    'payroll_min_amount': 'Montant Minimum Paie',
    'payroll_max_employees': 'Employés Maximum par Paie',
    'multi_sig_min_approvers': 'Approbateurs Minimum Multi-Sig',
    'multi_sig_max_approvers': 'Approbateurs Maximum Multi-Sig',
};

function getFriendlyName(key: string, defaultName: string): string {
    return FRIENDLY_NAMES[key] || defaultName;
}

function matchesService(key: string): boolean {
    return SERVICE_CONFIG_PREFIXES.some(prefix => key.startsWith(prefix));
}

function getCategory(key: string): string {
    if (key.startsWith('limit_')) return 'Limites';
    if (key.startsWith('payroll_')) return 'Paie';
    if (key.startsWith('multi_sig_')) return 'Multi-Signature';
    return 'Général';
}

export default function EnterpriseSettingsPage() {
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

            // Filter only enterprise-related configs
            const enterpriseConfigs = allConfigs.filter((c: FeeConfig) => matchesService(c.key));
            setConfigs(enterpriseConfigs);
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
                    <div className="p-3 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-xl shadow-lg">
                        <BuildingOfficeIcon className="h-8 w-8 text-white" />
                    </div>
                    <div>
                        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                            Configuration Entreprise
                        </h1>
                        <p className="text-gray-500 dark:text-gray-400">
                            Limites, frais et paramètres pour les comptes entreprise
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

            {/* Enterprise Features Info */}
            <div className="mb-8 grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="bg-gradient-to-br from-indigo-50 to-purple-50 dark:from-indigo-900/20 dark:to-purple-900/20 border border-indigo-100 dark:border-indigo-800 rounded-xl p-4">
                    <div className="flex items-center gap-3 mb-2">
                        <UserGroupIcon className="h-6 w-6 text-indigo-600 dark:text-indigo-400" />
                        <h3 className="font-semibold text-indigo-900 dark:text-indigo-300">Multi-Utilisateurs</h3>
                    </div>
                    <p className="text-sm text-indigo-700 dark:text-indigo-400">Gestion des équipes avec rôles et permissions</p>
                </div>
                <div className="bg-gradient-to-br from-purple-50 to-pink-50 dark:from-purple-900/20 dark:to-pink-900/20 border border-purple-100 dark:border-purple-800 rounded-xl p-4">
                    <div className="flex items-center gap-3 mb-2">
                        <BanknotesIcon className="h-6 w-6 text-purple-600 dark:text-purple-400" />
                        <h3 className="font-semibold text-purple-900 dark:text-purple-300">Paie Automatisée</h3>
                    </div>
                    <p className="text-sm text-purple-700 dark:text-purple-400">Paiements en masse avec approbation multi-signature</p>
                </div>
                <div className="bg-gradient-to-br from-pink-50 to-rose-50 dark:from-pink-900/20 dark:to-rose-900/20 border border-pink-100 dark:border-pink-800 rounded-xl p-4">
                    <div className="flex items-center gap-3 mb-2">
                        <BuildingOfficeIcon className="h-6 w-6 text-pink-600 dark:text-pink-400" />
                        <h3 className="font-semibold text-pink-900 dark:text-pink-300">Limites Élevées</h3>
                    </div>
                    <p className="text-sm text-pink-700 dark:text-pink-400">Plafonds adaptés aux besoins professionnels</p>
                </div>
            </div>

            {/* Configs by Category */}
            {loading ? (
                <div className="flex items-center justify-center h-64">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-500"></div>
                </div>
            ) : configs.length === 0 ? (
                <div className="text-center py-12 bg-gray-50 dark:bg-gray-800 rounded-xl">
                    <BuildingOfficeIcon className="h-12 w-12 mx-auto text-gray-400 mb-4" />
                    <p className="text-gray-500 dark:text-gray-400">
                        Aucune configuration entreprise trouvée
                    </p>
                    <p className="text-sm text-gray-400 dark:text-gray-500 mt-2">
                        Les configurations seront initialisées automatiquement
                    </p>
                </div>
            ) : (
                <div className="space-y-8">
                    {Object.entries(groupedConfigs).map(([category, categoryConfigs]) => (
                        <div key={category}>
                            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
                                <span className="w-2 h-2 rounded-full bg-indigo-500"></span>
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
                                                    {config.description || config.key}
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
                                                <span className="text-gray-500 dark:text-gray-400">Valeur</span>
                                                <span className="font-medium text-gray-900 dark:text-white">
                                                    {config.fixed_amount > 0
                                                        ? `${config.fixed_amount.toLocaleString()} ${config.currency || ''}`
                                                        : config.percentage_amount > 0
                                                            ? `${config.percentage_amount}%`
                                                            : '-'}
                                                </span>
                                            </div>
                                        </div>

                                        <button
                                            onClick={() => handleEdit(config)}
                                            className="w-full flex items-center justify-center gap-2 px-4 py-2 bg-indigo-50 dark:bg-indigo-900/20 text-indigo-600 dark:text-indigo-400 rounded-lg hover:bg-indigo-100 dark:hover:bg-indigo-900/40 transition-colors"
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
                                <p className="text-sm text-gray-500">{editingConfig.description}</p>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                    Valeur ({editingConfig.currency || 'EUR'})
                                </label>
                                <input
                                    type="number"
                                    step="0.01"
                                    value={formData.fixed_amount}
                                    onChange={(e) => setFormData({ ...formData, fixed_amount: parseFloat(e.target.value) || 0 })}
                                    className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
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
                                    className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                                />
                            </div>

                            <div className="flex items-center gap-3">
                                <input
                                    type="checkbox"
                                    id="is_enabled"
                                    checked={formData.is_enabled}
                                    onChange={(e) => setFormData({ ...formData, is_enabled: e.target.checked })}
                                    className="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
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
                                    className="flex-1 px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed"
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
