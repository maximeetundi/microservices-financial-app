'use client';

import { useState, useEffect } from 'react';
import { useRouter, useParams } from 'next/navigation';
import {
    ArrowLeftIcon,
    PlusIcon,
    ServerIcon,
    CheckCircleIcon,
    XCircleIcon,
    TrashIcon,
    PencilIcon,
    LinkIcon,
    WalletIcon,
} from '@heroicons/react/24/outline';

interface ProviderInstance {
    id: string;
    provider_id: string;
    name: string;
    vault_secret_path: string;
    hot_wallet_id?: string;
    is_active: boolean;
    is_primary: boolean;
    priority: number;
    request_count: number;
    health_status: string;
    last_used_at?: string;
    last_error?: string;
    created_at: string;
}

interface PlatformAccount {
    id: string;
    currency: string;
    balance: number;
    alias: string;
    name?: string;
    account_type?: string;
}

interface InstanceWallet {
    id: string;
    instance_id: string;
    currency: string;
    hot_wallet_id: string;
    is_active: boolean;
    priority: number;
}

export default function AggregatorInstancesPage() {
    const params = useParams();
    const router = useRouter();
    const providerId = params.id as string;

    const [instances, setInstances] = useState<ProviderInstance[]>([]);
    const [loading, setLoading] = useState(true);
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [showLinkWalletModal, setShowLinkWalletModal] = useState(false);
    const [selectedInstance, setSelectedInstance] = useState<ProviderInstance | null>(null);
    const [hotWallets, setHotWallets] = useState<PlatformAccount[]>([]);
    const [countries, setCountries] = useState<any[]>([]);
    const [providerName, setProviderName] = useState('');
    const [selectedWalletIds, setSelectedWalletIds] = useState<string[]>([]);
    const [selectedWalletId, setSelectedWalletId] = useState('');
    const [selectedCurrency, setSelectedCurrency] = useState('');
    const [showAddCountryModal, setShowAddCountryModal] = useState(false);
    const [instanceWallets, setInstanceWallets] = useState<Record<string, InstanceWallet[]>>({});
    const [expandedInstance, setExpandedInstance] = useState<string | null>(null);
    const [newCountry, setNewCountry] = useState({
        country_code: '',
        country_name: '',
        currency: 'XOF',
        priority: 50,
        fee_percentage: 0
    });

    // Region presets for quick country addition
    const REGION_PRESETS: Record<string, { name: string; countries: { code: string; name: string; currency: string }[] }> = {
        UEMOA: {
            name: 'UEMOA (Zone XOF)',
            countries: [
                { code: 'BJ', name: 'Bénin', currency: 'XOF' },
                { code: 'BF', name: 'Burkina Faso', currency: 'XOF' },
                { code: 'CI', name: "Côte d'Ivoire", currency: 'XOF' },
                { code: 'GW', name: 'Guinée-Bissau', currency: 'XOF' },
                { code: 'ML', name: 'Mali', currency: 'XOF' },
                { code: 'NE', name: 'Niger', currency: 'XOF' },
                { code: 'SN', name: 'Sénégal', currency: 'XOF' },
                { code: 'TG', name: 'Togo', currency: 'XOF' },
            ]
        },
        CEMAC: {
            name: 'CEMAC (Zone XAF)',
            countries: [
                { code: 'CM', name: 'Cameroun', currency: 'XAF' },
                { code: 'CF', name: 'Centrafrique', currency: 'XAF' },
                { code: 'TD', name: 'Tchad', currency: 'XAF' },
                { code: 'CG', name: 'Congo', currency: 'XAF' },
                { code: 'GQ', name: 'Guinée Équatoriale', currency: 'XAF' },
                { code: 'GA', name: 'Gabon', currency: 'XAF' },
            ]
        },
        WEST_AFRICA: {
            name: 'Afrique de l\'Ouest (Autres)',
            countries: [
                { code: 'GH', name: 'Ghana', currency: 'GHS' },
                { code: 'NG', name: 'Nigeria', currency: 'NGN' },
                { code: 'GN', name: 'Guinée', currency: 'GNF' },
                { code: 'SL', name: 'Sierra Leone', currency: 'SLL' },
                { code: 'LR', name: 'Libéria', currency: 'LRD' },
                { code: 'GM', name: 'Gambie', currency: 'GMD' },
            ]
        },
        EAST_AFRICA: {
            name: 'Afrique de l\'Est',
            countries: [
                { code: 'KE', name: 'Kenya', currency: 'KES' },
                { code: 'TZ', name: 'Tanzanie', currency: 'TZS' },
                { code: 'UG', name: 'Ouganda', currency: 'UGX' },
                { code: 'RW', name: 'Rwanda', currency: 'RWF' },
            ]
        }
    };

    const addRegionCountries = async (regionKey: string) => {
        const region = REGION_PRESETS[regionKey];
        if (!region) return;

        const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
        const token = localStorage.getItem('admin_token');

        let addedCount = 0;
        for (const country of region.countries) {
            try {
                const response = await fetch(`${API_URL}/api/v1/admin/payment-providers/${providerId}/countries`, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        country_code: country.code,
                        country_name: country.name,
                        currency: country.currency,
                        priority: 50,
                        fee_percentage: 0
                    })
                });
                if (response.ok) addedCount++;
            } catch (e) {
                console.error(`Failed to add ${country.code}`, e);
            }
        }

        loadProviderDetails();
        alert(`${addedCount} pays ajoutés sur ${region.countries.length}`);
    };

    const [newInstance, setNewInstance] = useState({
        name: '',
        vault_secret_path: '',
        is_active: true,
        is_primary: false,
        priority: 50,
    });

    useEffect(() => {
        if (providerId) {
            loadInstances();
            loadHotWallets();
            loadProviderDetails();
        }
    }, [providerId]);

    const loadProviderDetails = async () => {
        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');
            const response = await fetch(`${API_URL}/api/v1/admin/payment-providers/${providerId}`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            if (response.ok) {
                const data = await response.json();
                setProviderName(data.provider.display_name || data.provider.name);
                setCountries(data.provider.countries || []);
            }
        } catch (e) {
            console.error(e);
        }
    };

    const toggleCountry = async (countryCode: string, isActive: boolean) => {
        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');

            // Assuming endpoint structure based on standard conventions in this project
            const response = await fetch(`${API_URL}/api/v1/admin/payment-providers/${providerId}/countries/${countryCode}/toggle`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ is_active: isActive })
            });

            if (response.ok) {
                loadProviderDetails(); // Reload to get updated state
            } else {
                alert('Erreur lors de la modification du pays');
            }
        } catch (e) {
            console.error('Error toggling country:', e);
        }
    };

    const loadInstances = async () => {
        setLoading(true);
        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');
            const response = await fetch(`${API_URL}/api/v1/admin/payment-providers/${providerId}/instances`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            if (response.ok) {
                const data = await response.json();
                setInstances(data.instances || []);
            }
        } finally {
            setLoading(false);
        }
    };

    const loadHotWallets = async () => {
        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');
            const response = await fetch(`${API_URL}/api/v1/admin/platform/accounts?type=operations`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            if (response.ok) {
                const data = await response.json();
                setHotWallets(data.accounts || []);
            }
        } catch (e) {
            console.error('Failed to load wallets', e);
        }
    };

    const loadInstanceWallets = async (instanceId: string) => {
        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');
            const response = await fetch(`${API_URL}/api/v1/admin/payment-providers/${providerId}/instances/${instanceId}/wallets`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            if (response.ok) {
                const data = await response.json();
                setInstanceWallets(prev => ({ ...prev, [instanceId]: data.wallets || [] }));
            }
        } catch (e) {
            console.error('Failed to load instance wallets', e);
        }
    };

    const addWalletToInstance = async (instanceId: string, currency: string, hotWalletId: string) => {
        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');
            const response = await fetch(`${API_URL}/api/v1/admin/payment-providers/${providerId}/instances/${instanceId}/wallets`, {
                method: 'POST',
                headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
                body: JSON.stringify({ currency, hot_wallet_id: hotWalletId, priority: 50 })
            });
            if (response.ok) {
                loadInstanceWallets(instanceId);
                setShowLinkWalletModal(false);
            }
        } catch (e) {
            console.error('Failed to add wallet', e);
        }
    };

    const removeWalletFromInstance = async (instanceId: string, walletLinkId: string) => {
        if (!confirm('Retirer ce wallet de l\'instance ?')) return;
        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');
            await fetch(`${API_URL}/api/v1/admin/payment-providers/${providerId}/instances/${instanceId}/wallets/${walletLinkId}`, {
                method: 'DELETE',
                headers: { 'Authorization': `Bearer ${token}` }
            });
            loadInstanceWallets(instanceId);
        } catch (e) {
            console.error('Failed to remove wallet', e);
        }
    };

    const toggleWalletSelection = (walletId: string) => {
        setSelectedWalletIds(prev =>
            prev.includes(walletId)
                ? prev.filter(id => id !== walletId)
                : [...prev, walletId]
        );
    };

    const toggleAllWallets = () => {
        if (selectedWalletIds.length === hotWallets.length) {
            setSelectedWalletIds([]);
        } else {
            setSelectedWalletIds(hotWallets.map(w => w.id));
        }
    };

    const handleCreateInstance = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');
            const path = newInstance.vault_secret_path || `secret/aggregators/${providerName.toLowerCase().replace(/\s+/g, '_')}/${newInstance.name.toLowerCase().replace(/\s+/g, '_')}`;

            // Build wallets array
            const walletsPayload = selectedWalletIds.map(id => {
                const w = hotWallets.find(hw => hw.id === id);
                return {
                    hot_wallet_id: id,
                    currency: w?.currency || 'XOF'
                };
            });

            const response = await fetch(`${API_URL}/api/v1/admin/payment-providers/${providerId}/instances`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    ...newInstance,
                    vault_secret_path: path,
                    wallets: walletsPayload
                })
            });

            if (response.ok) {
                setShowCreateModal(false);
                loadInstances();
                setNewInstance({
                    name: '',
                    vault_secret_path: '',
                    is_active: true,
                    is_primary: false,
                    priority: 50,
                });
                setSelectedWalletIds([]);
            } else {
                const errData = await response.json();
                alert(`Failed to create instance: ${errData.error || 'Unknown error'}`);
            }
        } catch (e) {
            console.error(e);
            alert('Failed to create instance: Network error');
        }
    };

    const handleLinkWallet = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!selectedInstance) return;

        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');

            const response = await fetch(`${API_URL}/api/v1/admin/payment-providers/${providerId}/instances/${selectedInstance.id}/link-wallet`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ hot_wallet_id: selectedWalletId })
            });

            if (response.ok) {
                setShowLinkWalletModal(false);
                loadInstances();
            } else {
                alert('Failed to link wallet');
            }
        } catch (e) {
            console.error(e);
        }
    };

    const deleteInstance = async (id: string) => {
        if (!confirm('Êtes-vous sûr de vouloir supprimer cette instance ?')) return;
        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');
            await fetch(`${API_URL}/api/v1/admin/payment-providers/${providerId}/instances/${id}`, {
                method: 'DELETE',
                headers: { 'Authorization': `Bearer ${token}` }
            });
            loadInstances();
        } catch (e) {
            console.error(e);
        }
    };

    const handleAddCountry = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');

            const response = await fetch(`${API_URL}/api/v1/admin/payment-providers/${providerId}/countries`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(newCountry)
            });

            if (response.ok) {
                setShowAddCountryModal(false);
                loadProviderDetails();
                setNewCountry({
                    country_code: '',
                    country_name: '',
                    currency: 'XOF',
                    priority: 50,
                    fee_percentage: 0
                });
            } else {
                alert('Erreur lors de l\'ajout du pays');
            }
        } catch (e) {
            console.error('Error adding country:', e);
        }
    };

    const getFlagEmoji = (countryCode: string) => {
        if (!countryCode) return '🏳️';
        const codePoints = countryCode
            .toUpperCase()
            .split('')
            .map(char => 127397 + char.charCodeAt(0));
        return String.fromCodePoint(...codePoints);
    };

    return (
        <div className="space-y-6 animate-fadeIn">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div className="flex items-center gap-4">
                    <button
                        onClick={() => router.back()}
                        className="p-2 rounded-lg hover:bg-gray-100 text-gray-500"
                    >
                        <ArrowLeftIcon className="w-5 h-5" />
                    </button>
                    <div>
                        <h1 className="text-2xl font-bold text-gray-900">{providerName}</h1>
                        <p className="text-gray-500 text-sm">Gestion des instances et des pays</p>
                    </div>
                </div>
                <button
                    onClick={() => setShowCreateModal(true)}
                    className="btn-primary flex items-center gap-2"
                >
                    <PlusIcon className="w-5 h-5" />
                    Nouvelle Instance
                </button>
            </div>

            {/* Countries Management Section */}
            <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
                <div className="flex items-center justify-between mb-4">
                    <h2 className="text-lg font-bold text-gray-900">Pays Supportés</h2>
                    <button
                        onClick={() => setShowAddCountryModal(true)}
                        className="btn-secondary flex items-center gap-2 text-sm"
                    >
                        <PlusIcon className="w-4 h-4" />
                        Ajouter Pays
                    </button>
                </div>
                {countries.length === 0 ? (
                    <p className="text-gray-500 italic">Aucun pays configuré pour cet agrégateur.</p>
                ) : (
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                        {countries.map((country: any) => (
                            <div key={country.country_code} className="flex items-center justify-between p-4 border rounded-xl bg-gray-50">
                                <div className="flex items-center gap-3">
                                    <span className="text-2xl">{getFlagEmoji(country.country_code)}</span>
                                    <div>
                                        <p className="font-semibold text-gray-900">{country.country_name}</p>
                                        <p className="text-xs text-gray-500">{country.currency}</p>
                                    </div>
                                </div>
                                <label className="relative inline-flex items-center cursor-pointer">
                                    <input
                                        type="checkbox"
                                        className="sr-only peer"
                                        checked={country.is_active}
                                        onChange={(e) => toggleCountry(country.country_code, e.target.checked)}
                                    />
                                    <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-indigo-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
                                </label>
                            </div>
                        ))}
                    </div>
                )}
            </div>

            <h2 className="text-xl font-bold text-gray-900 mt-8 mb-4">Instances Configurées</h2>

            {/* List */}
            {loading ? (

                <div className="flex justify-center py-12">
                    <div className="spinner w-8 h-8" />
                </div>
            ) : instances.length === 0 ? (
                <div className="text-center py-16 bg-white rounded-xl border border-gray-100">
                    <ServerIcon className="w-16 h-16 mx-auto text-gray-300 mb-4" />
                    <h3 className="text-lg font-medium text-gray-900">Aucune instance configurée</h3>
                    <p className="text-gray-500 mb-6">Ajoutez une première instance pour activer ce provider</p>
                    <button onClick={() => setShowCreateModal(true)} className="btn-secondary">
                        Ajouter une instance
                    </button>
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {instances.map(inst => (
                        <div key={inst.id} className="bg-white rounded-xl shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow">
                            <div className="flex justify-between items-start mb-4">
                                <div>
                                    <div className="flex items-center gap-2">
                                        <h3 className="font-bold text-gray-900">{inst.name}</h3>
                                        {inst.is_primary && (
                                            <span className="bg-indigo-100 text-indigo-700 text-xs px-2 py-0.5 rounded-full font-medium">Principal</span>
                                        )}
                                    </div>
                                    <p className="text-xs text-gray-500 font-mono mt-1 break-all">{inst.vault_secret_path}</p>
                                </div>
                                <div className={`w-3 h-3 rounded-full ${inst.is_active ? 'bg-emerald-500' : 'bg-gray-300'}`} />
                            </div>

                            <div className="space-y-3 mb-6">
                                <div className="flex justify-between text-sm">
                                    <span className="text-gray-500">Santé</span>
                                    <span className={`font-medium ${inst.health_status === 'healthy' ? 'text-emerald-600' :
                                        inst.health_status === 'error' ? 'text-red-600' : 'text-gray-500'
                                        }`}>
                                        {inst.health_status}
                                    </span>
                                </div>
                                <div className="flex justify-between text-sm">
                                    <span className="text-gray-500">Requêtes</span>
                                    <span className="font-medium text-gray-900">{inst.request_count}</span>
                                </div>
                                <div className="flex justify-between text-sm">
                                    <span className="text-gray-500">Priorité</span>
                                    <span className="font-medium text-gray-900">{inst.priority}</span>
                                </div>
                            </div>

                            {/* Multi-Wallet Section */}
                            <div className="mb-4">
                                <button
                                    onClick={() => {
                                        if (expandedInstance === inst.id) {
                                            setExpandedInstance(null);
                                        } else {
                                            setExpandedInstance(inst.id);
                                            loadInstanceWallets(inst.id);
                                        }
                                    }}
                                    className="w-full flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
                                >
                                    <div className="flex items-center gap-2 text-sm text-gray-700">
                                        <WalletIcon className="w-4 h-4 text-emerald-600" />
                                        <span>Hot Wallets</span>
                                        {instanceWallets[inst.id]?.length > 0 && (
                                            <span className="bg-emerald-100 text-emerald-700 text-xs px-2 py-0.5 rounded-full">
                                                {instanceWallets[inst.id].length} liés
                                            </span>
                                        )}
                                    </div>
                                    <span className="text-gray-400 text-xs">
                                        {expandedInstance === inst.id ? '▲' : '▼'}
                                    </span>
                                </button>

                                {expandedInstance === inst.id && (
                                    <div className="mt-2 space-y-2 p-3 border border-gray-200 rounded-lg">
                                        {/* Currency badges for linked wallets */}
                                        {instanceWallets[inst.id]?.length > 0 ? (
                                            <div className="flex flex-wrap gap-2 mb-3">
                                                {instanceWallets[inst.id].map(w => (
                                                    <div key={w.id} className="flex items-center gap-1 bg-emerald-50 border border-emerald-200 rounded-full px-2 py-1">
                                                        <span className="text-xs font-medium text-emerald-700">
                                                            {w.currency} <span className="opacity-75 text-[10px] hidden sm:inline">- {w.name || 'Wallet'}</span>
                                                        </span>
                                                        <button
                                                            onClick={() => removeWalletFromInstance(inst.id, w.id)}
                                                            className="text-red-400 hover:text-red-600 ml-1"
                                                            title="Retirer"
                                                        >
                                                            ×
                                                        </button>
                                                    </div>
                                                ))}
                                            </div>
                                        ) : (
                                            <p className="text-xs text-gray-500 italic mb-2">Aucun wallet lié</p>
                                        )}

                                        {/* Add wallet button */}
                                        <button
                                            onClick={() => {
                                                setSelectedInstance(inst);
                                                setShowLinkWalletModal(true);
                                            }}
                                            className="w-full py-2 border border-dashed border-gray-300 rounded-lg text-xs text-gray-500 hover:border-emerald-500 hover:text-emerald-600 transition-colors flex items-center justify-center gap-1"
                                        >
                                            <PlusIcon className="w-3 h-3" />
                                            Ajouter Wallet
                                        </button>
                                    </div>
                                )}
                            </div>

                            <div className="flex gap-2 pt-4 border-t border-gray-100">
                                <button className="flex-1 btn-secondary text-xs">Test</button>
                                <button
                                    onClick={() => deleteInstance(inst.id)}
                                    className="p-2 text-red-500 hover:bg-red-50 rounded-lg transition-colors"
                                >
                                    <TrashIcon className="w-4 h-4" />
                                </button>
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {/* Create Modal */}
            {showCreateModal && (
                <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
                    <div className="bg-white rounded-2xl w-full max-w-md p-6 animate-scaleIn">
                        <h2 className="text-xl font-bold mb-4">Nouvelle Instance</h2>
                        <form onSubmit={handleCreateInstance} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Nom (ex: Instance 1)</label>
                                <input
                                    type="text"
                                    required
                                    className="input w-full"
                                    value={newInstance.name}
                                    onChange={e => setNewInstance({ ...newInstance, name: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Chemin Vault (Optionnel)</label>
                                <input
                                    type="text"
                                    className="input w-full font-mono text-sm"
                                    placeholder="Généré automatiquement si vide"
                                    value={newInstance.vault_secret_path}
                                    onChange={e => setNewInstance({ ...newInstance, vault_secret_path: e.target.value })}
                                />
                            </div>

                            {/* Hot Wallet Selection (Multi) */}
                            <div>
                                <div className="flex justify-between items-center mb-2">
                                    <label className="block text-sm font-medium text-gray-700">Hot Wallets Associés</label>
                                    <button
                                        type="button"
                                        onClick={toggleAllWallets}
                                        className="text-xs text-blue-600 hover:text-blue-800"
                                    >
                                        {selectedWalletIds.length === hotWallets.length ? 'Tout désélectionner' : 'Tout sélectionner'}
                                    </button>
                                </div>
                                <div className="border rounded-md p-3 max-h-48 overflow-y-auto space-y-2 bg-gray-50">
                                    {hotWallets.length === 0 ? (
                                        <p className="text-gray-500 text-sm">Aucun hot wallet disponible.</p>
                                    ) : (
                                        hotWallets.map(wallet => (
                                            <label key={wallet.id} className="flex items-center space-x-2 cursor-pointer p-1 hover:bg-gray-100 rounded">
                                                <input
                                                    type="checkbox"
                                                    className="checkbox checkbox-xs"
                                                    checked={selectedWalletIds.includes(wallet.id)}
                                                    onChange={() => toggleWalletSelection(wallet.id)}
                                                />
                                                <span className="text-sm">
                                                    <span className="font-semibold">{wallet.currency}</span> - {wallet.name}
                                                    <span className="text-gray-500 text-xs ml-1">(Solde: {wallet.balance})</span>
                                                </span>
                                            </label>
                                        ))
                                    )}
                                </div>
                                <p className="text-xs text-gray-500 mt-1">
                                    Associer des wallets pour permettre les dépôts/retraits immédiats.
                                </p>
                            </div>

                            <div className="flex gap-4">
                                <div className="flex-1">
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Priorité</label>
                                    <input
                                        type="number"
                                        className="input w-full"
                                        value={newInstance.priority}
                                        onChange={e => setNewInstance({ ...newInstance, priority: parseInt(e.target.value) })}
                                    />
                                </div>
                                <div className="flex items-end mb-3">
                                    <label className="flex items-center gap-2 cursor-pointer">
                                        <input
                                            type="checkbox"
                                            className="checkbox"
                                            checked={newInstance.is_primary}
                                            onChange={e => setNewInstance({ ...newInstance, is_primary: e.target.checked })}
                                        />
                                        <span className="text-sm">Principal</span>
                                    </label>
                                </div>
                            </div>

                            <div className="flex justify-end gap-3 mt-6">
                                <button type="button" onClick={() => setShowCreateModal(false)} className="btn-secondary">Annuler</button>
                                <button type="submit" className="btn-primary">Créer</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            {/* Link Wallet Modal */}
            {showLinkWalletModal && (
                <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
                    <div className="bg-white rounded-2xl w-full max-w-md p-6 animate-scaleIn">
                        <h2 className="text-xl font-bold mb-4">Lier Hot Wallets</h2>
                        <p className="text-sm text-gray-500 mb-4">Sélectionnez les wallets à lier à cette instance.</p>

                        <div className="flex justify-between items-center mb-2">
                            <label className="text-sm font-medium text-gray-700">Wallets Disponibles</label>
                            <button
                                type="button"
                                onClick={() => {
                                    const linkedIds = instanceWallets[selectedInstance?.id || '']?.map(w => w.id) || [];
                                    const available = hotWallets.filter(w => !linkedIds.includes(w.id));

                                    // If all available are selected, deselect all. Otherwise select all available.
                                    const availableIds = available.map(w => w.id);
                                    const allSelected = availableIds.every(id => selectedWalletIds.includes(id));

                                    if (allSelected) {
                                        setSelectedWalletIds([]);
                                    } else {
                                        setSelectedWalletIds(availableIds);
                                    }
                                }}
                                className="text-xs text-blue-600 hover:text-blue-800"
                            >
                                Tout sélectionner / désélectionner
                            </button>
                        </div>

                        <div className="space-y-2 max-h-72 overflow-y-auto mb-6 border rounded-lg p-2 bg-gray-50">
                            {hotWallets.filter(w => !instanceWallets[selectedInstance?.id || '']?.some(linked => linked.id === w.id)).length === 0 && (
                                <p className="text-center text-gray-500 text-sm py-8 italic">Tous les wallets disponibles sont déjà liés.</p>
                            )}

                            {hotWallets.map(wallet => {
                                const isLinked = instanceWallets[selectedInstance?.id || '']?.some(w => w.id === wallet.id);
                                if (isLinked) return null;

                                return (
                                    <label
                                        key={wallet.id}
                                        className={`flex items-center p-3 rounded-lg border cursor-pointer transition-all ${selectedWalletIds.includes(wallet.id)
                                            ? 'border-indigo-500 bg-indigo-50 ring-1 ring-indigo-500'
                                            : 'border-gray-200 hover:border-gray-300 bg-white'
                                            }`}
                                    >
                                        <input
                                            type="checkbox"
                                            className="checkbox checkbox-sm mr-3"
                                            checked={selectedWalletIds.includes(wallet.id)}
                                            onChange={() => {
                                                if (selectedWalletIds.includes(wallet.id)) {
                                                    setSelectedWalletIds(prev => prev.filter(id => id !== wallet.id));
                                                } else {
                                                    setSelectedWalletIds(prev => [...prev, wallet.id]);
                                                }
                                            }}
                                        />
                                        <div className="flex-1">
                                            <div className="flex justify-between items-center">
                                                <span className="font-medium text-gray-900">{wallet.name || wallet.alias || 'Wallet'}</span>
                                                <div className="flex items-center gap-1 bg-gray-100 px-2 py-0.5 rounded text-xs font-mono">
                                                    {wallet.currency}
                                                </div>
                                            </div>
                                            <p className="text-xs text-gray-500 mt-0.5">
                                                Solde: {wallet.balance?.toLocaleString()} {wallet.currency}
                                            </p>
                                        </div>
                                    </label>
                                );
                            })}
                        </div>

                        <div className="flex justify-end gap-3">
                            <button
                                type="button"
                                onClick={() => {
                                    setShowLinkWalletModal(false);
                                    setSelectedWalletIds([]);
                                }}
                                className="btn-secondary"
                            >
                                Annuler
                            </button>
                            <button
                                onClick={async () => {
                                    if (selectedInstance && selectedWalletIds.length > 0) {
                                        // Loop through to add multiple
                                        for (const wid of selectedWalletIds) {
                                            const w = hotWallets.find(hw => hw.id === wid);
                                            if (w) {
                                                await addWalletToInstance(selectedInstance.id, w.currency, w.id);
                                            }
                                        }
                                        // Close and refresh happens in addWalletToInstance but we might need to manual refresh if loop
                                        // Actually addWalletToInstance calls loadInstanceWallets. 
                                        // To avoid race conditions visually, we just close.
                                        setShowLinkWalletModal(false);
                                        setSelectedWalletIds([]);
                                        // Force reload one last time to be sure
                                        setTimeout(() => loadInstanceWallets(selectedInstance.id), 500);
                                    }
                                }}
                                className="btn-primary"
                                disabled={selectedWalletIds.length === 0}
                            >
                                Lier ({selectedWalletIds.length})
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {/* Add Country Modal */}
            {showAddCountryModal && (
                <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
                    <div className="bg-white rounded-2xl w-full max-w-lg p-6 animate-scaleIn">
                        <h2 className="text-xl font-bold mb-4">Ajouter un Pays</h2>

                        {/* Quick Add Regions */}
                        <div className="mb-6 p-4 bg-gray-50 rounded-xl">
                            <p className="text-sm font-medium text-gray-700 mb-3">Ajout rapide par région</p>
                            <div className="flex flex-wrap gap-2">
                                {Object.entries(REGION_PRESETS).map(([key, region]) => (
                                    <button
                                        key={key}
                                        type="button"
                                        onClick={() => {
                                            setShowAddCountryModal(false);
                                            addRegionCountries(key);
                                        }}
                                        className="px-3 py-1.5 text-xs font-medium bg-indigo-100 text-indigo-700 rounded-full hover:bg-indigo-200 transition-colors"
                                    >
                                        + {region.name}
                                    </button>
                                ))}
                            </div>
                        </div>

                        <div className="border-t border-gray-200 pt-4">
                            <p className="text-sm font-medium text-gray-700 mb-3">Ou ajouter un pays individuel</p>
                            <form onSubmit={handleAddCountry} className="space-y-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Code Pays (ex: CI, SN)</label>
                                    <input
                                        type="text"
                                        required
                                        maxLength={3}
                                        className="input w-full uppercase"
                                        placeholder="CI"
                                        value={newCountry.country_code}
                                        onChange={e => setNewCountry({ ...newCountry, country_code: e.target.value.toUpperCase() })}
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Nom du Pays</label>
                                    <input
                                        type="text"
                                        required
                                        className="input w-full"
                                        placeholder="Côte d'Ivoire"
                                        value={newCountry.country_name}
                                        onChange={e => setNewCountry({ ...newCountry, country_name: e.target.value })}
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Devise</label>
                                    <select
                                        className="input w-full"
                                        value={newCountry.currency}
                                        onChange={e => setNewCountry({ ...newCountry, currency: e.target.value })}
                                    >
                                        <option value="XOF">XOF (Franc CFA)</option>
                                        <option value="XAF">XAF (Franc CFA CEMAC)</option>
                                        <option value="GNF">GNF (Franc Guinéen)</option>
                                        <option value="NGN">NGN (Naira)</option>
                                        <option value="GHS">GHS (Cedi)</option>
                                        <option value="KES">KES (Shilling Kenyan)</option>
                                        <option value="USD">USD (Dollar US)</option>
                                        <option value="EUR">EUR (Euro)</option>
                                    </select>
                                </div>
                                <div className="flex justify-end gap-3 mt-6">
                                    <button type="button" onClick={() => setShowAddCountryModal(false)} className="btn-secondary">Annuler</button>
                                    <button type="submit" className="btn-primary">Ajouter</button>
                                </div>
                            </form>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
