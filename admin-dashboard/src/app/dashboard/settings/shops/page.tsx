'use client';

import { useState, useEffect } from 'react';
import {
    ShoppingBagIcon,
    ArrowPathIcon,
    PencilSquareIcon,
    XMarkIcon,
    CheckCircleIcon,
    CurrencyDollarIcon,
    TruckIcon,
    ReceiptPercentIcon,
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

// Define which config keys belong to shop service
const SERVICE_CONFIG_PREFIXES = [
    'shop_',
    'ecommerce_',
    'order_',
    'product_',
    'shipping_',
];

const FRIENDLY_NAMES: Record<string, string> = {
    'shop_commission_rate': 'Taux de Commission',
    'shop_min_payout': 'Paiement Minimum',
    'shop_payout_delay_days': 'Délai de Paiement (jours)',
    'shop_max_products': 'Produits Maximum par Boutique',
    'shop_max_categories': 'Catégories Maximum par Boutique',
    'ecommerce_platform_fee': 'Frais Plateforme E-commerce',
    'order_min_amount': 'Montant Minimum Commande',
    'order_max_amount': 'Montant Maximum Commande',
    'order_cancellation_fee': 'Frais d\'Annulation',
    'shipping_default_fee': 'Frais de Livraison par Défaut',
    'shipping_free_threshold': 'Seuil Livraison Gratuite',
    'product_max_images': 'Images Maximum par Produit',
    'product_max_variants': 'Variantes Maximum par Produit',
};

function getFriendlyName(key: string, defaultName: string): string {
    return FRIENDLY_NAMES[key] || defaultName;
}

function matchesService(key: string): boolean {
    return SERVICE_CONFIG_PREFIXES.some(prefix => key.startsWith(prefix));
}

function getCategory(key: string): string {
    if (key.startsWith('shop_')) return 'Boutique';
    if (key.startsWith('order_')) return 'Commandes';
    if (key.startsWith('shipping_')) return 'Livraison';
    if (key.startsWith('product_')) return 'Produits';
    return 'Général';
}

function getCategoryIcon(category: string) {
    switch (category) {
        case 'Boutique': return ShoppingBagIcon;
        case 'Commandes': return ReceiptPercentIcon;
        case 'Livraison': return TruckIcon;
        case 'Produits': return CurrencyDollarIcon;
        default: return ShoppingBagIcon;
    }
}

export default function ShopSettingsPage() {
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

            // Filter only shop-related configs
            const shopConfigs = allConfigs.filter((c: FeeConfig) => matchesService(c.key));
            setConfigs(shopConfigs);
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
                    <div className="p-3 bg-gradient-to-br from-pink-500 to-rose-600 rounded-xl shadow-lg">
                        <ShoppingBagIcon className="h-8 w-8 text-white" />
                    </div>
                    <div>
                        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                            Configuration Boutiques
                        </h1>
                        <p className="text-gray-500 dark:text-gray-400">
                            Commissions, frais et paramètres e-commerce
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

            {/* Shop Features Info */}
            <div className="mb-8 grid grid-cols-1 md:grid-cols-4 gap-4">
                <div className="bg-gradient-to-br from-pink-50 to-rose-50 dark:from-pink-900/20 dark:to-rose-900/20 border border-pink-100 dark:border-pink-800 rounded-xl p-4">
                    <div className="flex items-center gap-2 mb-1">
                        <ShoppingBagIcon className="h-5 w-5 text-pink-600 dark:text-pink-400" />
                        <span className="font-semibold text-pink-900 dark:text-pink-300">Boutiques</span>
                    </div>
                    <p className="text-xs text-pink-700 dark:text-pink-400">Configuration des marchands</p>
                </div>
                <div className="bg-gradient-to-br from-orange-50 to-amber-50 dark:from-orange-900/20 dark:to-amber-900/20 border border-orange-100 dark:border-orange-800 rounded-xl p-4">
                    <div className="flex items-center gap-2 mb-1">
                        <ReceiptPercentIcon className="h-5 w-5 text-orange-600 dark:text-orange-400" />
                        <span className="font-semibold text-orange-900 dark:text-orange-300">Commandes</span>
                    </div>
                    <p className="text-xs text-orange-700 dark:text-orange-400">Limites et frais</p>
                </div>
                <div className="bg-gradient-to-br from-cyan-50 to-teal-50 dark:from-cyan-900/20 dark:to-teal-900/20 border border-cyan-100 dark:border-cyan-800 rounded-xl p-4">
                    <div className="flex items-center gap-2 mb-1">
                        <TruckIcon className="h-5 w-5 text-cyan-600 dark:text-cyan-400" />
                        <span className="font-semibold text-cyan-900 dark:text-cyan-300">Livraison</span>
                    </div>
                    <p className="text-xs text-cyan-700 dark:text-cyan-400">Frais et seuils</p>
                </div>
                <div className="bg-gradient-to-br from-violet-50 to-purple-50 dark:from-violet-900/20 dark:to-purple-900/20 border border-violet-100 dark:border-violet-800 rounded-xl p-4">
                    <div className="flex items-center gap-2 mb-1">
                        <CurrencyDollarIcon className="h-5 w-5 text-violet-600 dark:text-violet-400" />
                        <span className="font-semibold text-violet-900 dark:text-violet-300">Produits</span>
                    </div>
                    <p className="text-xs text-violet-700 dark:text-violet-400">Limites par produit</p>
                </div>
            </div>

            {/* Configs by Category */}
            {loading ? (
                <div className="flex items-center justify-center h-64">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-pink-500"></div>
                </div>
            ) : configs.length === 0 ? (
                <div className="text-center py-12 bg-gray-50 dark:bg-gray-800 rounded-xl">
                    <ShoppingBagIcon className="h-12 w-12 mx-auto text-gray-400 mb-4" />
                    <p className="text-gray-500 dark:text-gray-400">
                        Aucune configuration boutique trouvée
                    </p>
                    <p className="text-sm text-gray-400 dark:text-gray-500 mt-2">
                        Ajoutez des configurations avec le préfixe "shop_", "order_", etc.
                    </p>
                </div>
            ) : (
                <div className="space-y-8">
                    {Object.entries(groupedConfigs).map(([category, categoryConfigs]) => {
                        const CategoryIcon = getCategoryIcon(category);
                        return (
                            <div key={category}>
                                <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
                                    <CategoryIcon className="h-5 w-5 text-pink-500" />
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
                                                {config.fixed_amount > 0 && (
                                                    <div className="flex justify-between text-sm">
                                                        <span className="text-gray-500 dark:text-gray-400">Montant</span>
                                                        <span className="font-medium text-gray-900 dark:text-white">
                                                            {config.fixed_amount.toLocaleString()} {config.currency}
                                                        </span>
                                                    </div>
                                                )}
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
                                                className="w-full flex items-center justify-center gap-2 px-4 py-2 bg-pink-50 dark:bg-pink-900/20 text-pink-600 dark:text-pink-400 rounded-lg hover:bg-pink-100 dark:hover:bg-pink-900/40 transition-colors"
                                            >
                                                <PencilSquareIcon className="h-4 w-4" />
                                                Modifier
                                            </button>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        );
                    })}
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
                                    Montant Fixe
                                </label>
                                <input
                                    type="number"
                                    step="0.01"
                                    value={formData.fixed_amount}
                                    onChange={(e) => setFormData({ ...formData, fixed_amount: parseFloat(e.target.value) || 0 })}
                                    className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-pink-500 focus:border-transparent"
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
                                    className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-pink-500 focus:border-transparent"
                                />
                            </div>

                            <div className="flex items-center gap-3">
                                <input
                                    type="checkbox"
                                    id="is_enabled"
                                    checked={formData.is_enabled}
                                    onChange={(e) => setFormData({ ...formData, is_enabled: e.target.checked })}
                                    className="w-4 h-4 text-pink-600 border-gray-300 rounded focus:ring-pink-500"
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
                                    className="flex-1 px-4 py-2 bg-pink-600 text-white rounded-lg hover:bg-pink-700 disabled:opacity-50 disabled:cursor-not-allowed"
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
