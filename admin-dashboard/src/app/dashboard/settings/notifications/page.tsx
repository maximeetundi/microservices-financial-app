'use client';

import { useState, useEffect } from 'react';
import {
    BellIcon,
    ArrowPathIcon,
    PencilSquareIcon,
    XMarkIcon,
    CheckCircleIcon,
    ChatBubbleLeftRightIcon,
    EnvelopeIcon,
    DevicePhoneMobileIcon,
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

// Define which config keys belong to notification/messaging service
const SERVICE_CONFIG_PREFIXES = [
    'notification_',
    'sms_',
    'email_',
    'push_',
    'messaging_',
];

const FRIENDLY_NAMES: Record<string, string> = {
    'notification_enabled': 'Notifications Activées',
    'notification_transaction_alert': 'Alertes Transactions',
    'notification_marketing': 'Communications Marketing',
    'notification_security_alert': 'Alertes Sécurité',
    'sms_enabled': 'SMS Activés',
    'sms_cost_per_message': 'Coût par SMS',
    'sms_daily_limit': 'Limite SMS Quotidienne',
    'sms_verification_enabled': 'Vérification SMS',
    'email_enabled': 'Emails Activés',
    'email_daily_limit': 'Limite Emails Quotidienne',
    'email_marketing_enabled': 'Emails Marketing',
    'push_enabled': 'Push Notifications',
    'push_daily_limit': 'Limite Push Quotidienne',
    'messaging_enabled': 'Messagerie Interne',
    'messaging_max_attachments': 'Pièces Jointes Max',
};

function getFriendlyName(key: string, defaultName: string): string {
    return FRIENDLY_NAMES[key] || defaultName;
}

function matchesService(key: string): boolean {
    return SERVICE_CONFIG_PREFIXES.some(prefix => key.startsWith(prefix));
}

function getCategory(key: string): string {
    if (key.startsWith('sms_')) return 'SMS';
    if (key.startsWith('email_')) return 'Email';
    if (key.startsWith('push_')) return 'Push Notifications';
    if (key.startsWith('messaging_')) return 'Messagerie Interne';
    return 'Notifications Générales';
}

export default function NotificationSettingsPage() {
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
            const allConfigs = response.data.configurations || [];

            // Filter only notification-related configs
            const notificationConfigs = allConfigs.filter((c: FeeConfig) => matchesService(c.key));
            setConfigs(notificationConfigs);
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
                    <div className="p-3 bg-gradient-to-br from-sky-500 to-blue-600 rounded-xl shadow-lg">
                        <BellIcon className="h-8 w-8 text-white" />
                    </div>
                    <div>
                        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                            Configuration Notifications
                        </h1>
                        <p className="text-gray-500 dark:text-gray-400">
                            SMS, Emails, Push et messagerie interne
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

            {/* Channel Overview */}
            <div className="mb-8 grid grid-cols-1 md:grid-cols-4 gap-4">
                <div className="bg-gradient-to-br from-green-50 to-emerald-50 dark:from-green-900/20 dark:to-emerald-900/20 border border-green-100 dark:border-green-800 rounded-xl p-4">
                    <div className="flex items-center gap-2 mb-1">
                        <DevicePhoneMobileIcon className="h-5 w-5 text-green-600 dark:text-green-400" />
                        <span className="font-semibold text-green-900 dark:text-green-300">SMS</span>
                    </div>
                    <p className="text-xs text-green-700 dark:text-green-400">Vérifications et alertes</p>
                </div>
                <div className="bg-gradient-to-br from-blue-50 to-indigo-50 dark:from-blue-900/20 dark:to-indigo-900/20 border border-blue-100 dark:border-blue-800 rounded-xl p-4">
                    <div className="flex items-center gap-2 mb-1">
                        <EnvelopeIcon className="h-5 w-5 text-blue-600 dark:text-blue-400" />
                        <span className="font-semibold text-blue-900 dark:text-blue-300">Email</span>
                    </div>
                    <p className="text-xs text-blue-700 dark:text-blue-400">Confirmations et rapports</p>
                </div>
                <div className="bg-gradient-to-br from-purple-50 to-violet-50 dark:from-purple-900/20 dark:to-violet-900/20 border border-purple-100 dark:border-purple-800 rounded-xl p-4">
                    <div className="flex items-center gap-2 mb-1">
                        <BellIcon className="h-5 w-5 text-purple-600 dark:text-purple-400" />
                        <span className="font-semibold text-purple-900 dark:text-purple-300">Push</span>
                    </div>
                    <p className="text-xs text-purple-700 dark:text-purple-400">Notifications mobiles</p>
                </div>
                <div className="bg-gradient-to-br from-orange-50 to-amber-50 dark:from-orange-900/20 dark:to-amber-900/20 border border-orange-100 dark:border-orange-800 rounded-xl p-4">
                    <div className="flex items-center gap-2 mb-1">
                        <ChatBubbleLeftRightIcon className="h-5 w-5 text-orange-600 dark:text-orange-400" />
                        <span className="font-semibold text-orange-900 dark:text-orange-300">Chat</span>
                    </div>
                    <p className="text-xs text-orange-700 dark:text-orange-400">Messagerie intégrée</p>
                </div>
            </div>

            {/* Configs by Category */}
            {loading ? (
                <div className="flex items-center justify-center h-64">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-sky-500"></div>
                </div>
            ) : configs.length === 0 ? (
                <div className="text-center py-12 bg-gray-50 dark:bg-gray-800 rounded-xl">
                    <BellIcon className="h-12 w-12 mx-auto text-gray-400 mb-4" />
                    <p className="text-gray-500 dark:text-gray-400">
                        Aucune configuration de notification trouvée
                    </p>
                    <p className="text-sm text-gray-400 dark:text-gray-500 mt-2">
                        Ajoutez des configurations avec les préfixes "notification_", "sms_", "email_", etc.
                    </p>
                </div>
            ) : (
                <div className="space-y-8">
                    {Object.entries(groupedConfigs).map(([category, categoryConfigs]) => (
                        <div key={category}>
                            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
                                <span className="w-2 h-2 rounded-full bg-sky-500"></span>
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
                                                    <span className="text-gray-500 dark:text-gray-400">Valeur</span>
                                                    <span className="font-medium text-gray-900 dark:text-white">
                                                        {config.fixed_amount.toLocaleString()} {config.currency}
                                                    </span>
                                                </div>
                                            )}
                                        </div>

                                        <button
                                            onClick={() => handleEdit(config)}
                                            className="w-full flex items-center justify-center gap-2 px-4 py-2 bg-sky-50 dark:bg-sky-900/20 text-sky-600 dark:text-sky-400 rounded-lg hover:bg-sky-100 dark:hover:bg-sky-900/40 transition-colors"
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
                                    Valeur
                                </label>
                                <input
                                    type="number"
                                    step="0.01"
                                    value={formData.fixed_amount}
                                    onChange={(e) => setFormData({ ...formData, fixed_amount: parseFloat(e.target.value) || 0 })}
                                    className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-sky-500 focus:border-transparent"
                                />
                            </div>

                            <div className="flex items-center gap-3">
                                <input
                                    type="checkbox"
                                    id="is_enabled"
                                    checked={formData.is_enabled}
                                    onChange={(e) => setFormData({ ...formData, is_enabled: e.target.checked })}
                                    className="w-4 h-4 text-sky-600 border-gray-300 rounded focus:ring-sky-500"
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
                                    className="flex-1 px-4 py-2 bg-sky-600 text-white rounded-lg hover:bg-sky-700 disabled:opacity-50 disabled:cursor-not-allowed"
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
