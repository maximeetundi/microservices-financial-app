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
    const [providerName, setProviderName] = useState('');

    // Form states
    const [newInstance, setNewInstance] = useState({
        name: '',
        vault_secret_path: '',
        is_active: true,
        is_primary: false,
        priority: 50,
    });

    const [selectedWalletId, setSelectedWalletId] = useState('');

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
                setProviderName(data.display_name || data.name);
            }
        } catch (e) {
            console.error(e);
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
            // Mocking endpoint structure based on proxy
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

    const handleCreateInstance = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');

            // Auto-generate vault path if empty based on name
            const path = newInstance.vault_secret_path || `secret/aggregators/${providerName.toLowerCase().replace(/\s+/g, '_')}/${newInstance.name.toLowerCase().replace(/\s+/g, '_')}`;

            const response = await fetch(`${API_URL}/api/v1/admin/payment-providers/${providerId}/instances`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ ...newInstance, vault_secret_path: path })
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
            } else {
                alert('Failed to create instance');
            }
        } catch (e) {
            console.error(e);
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
                        <h1 className="text-2xl font-bold text-gray-900">Instances: {providerName}</h1>
                        <p className="text-gray-500 text-sm">Gérez les clés API et les connexions wallets</p>
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

                            {inst.hot_wallet_id ? (
                                <div className="bg-gray-50 p-3 rounded-lg mb-4 flex items-center justify-between">
                                    <div className="flex items-center gap-2 text-sm text-gray-700">
                                        <WalletIcon className="w-4 h-4 text-emerald-600" />
                                        <span>Wallet Lié</span>
                                    </div>
                                    <button
                                        onClick={() => {
                                            setSelectedInstance(inst);
                                            setShowLinkWalletModal(true);
                                        }}
                                        className="text-xs text-indigo-600 hover:text-indigo-800"
                                    >
                                        Modifier
                                    </button>
                                </div>
                            ) : (
                                <button
                                    onClick={() => {
                                        setSelectedInstance(inst);
                                        setShowLinkWalletModal(true);
                                    }}
                                    className="w-full py-2 mb-4 border border-dashed border-gray-300 rounded-lg text-sm text-gray-500 hover:border-emerald-500 hover:text-emerald-600 transition-colors flex items-center justify-center gap-2"
                                >
                                    <LinkIcon className="w-4 h-4" />
                                    Lier Hot Wallet
                                </button>
                            )}

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
                        <h2 className="text-xl font-bold mb-4">Lier Hot Wallet</h2>
                        <p className="text-sm text-gray-500 mb-4">Sélectionnez le compte plateforme (hot wallet) qui financera cette instance.</p>

                        <form onSubmit={handleLinkWallet}>
                            <div className="space-y-3 max-h-60 overflow-y-auto mb-6">
                                {hotWallets.map(wallet => (
                                    <label key={wallet.id} className="flex items-center gap-3 p-3 border rounded-xl cursor-pointer hover:bg-gray-50">
                                        <input
                                            type="radio"
                                            name="wallet"
                                            value={wallet.id}
                                            checked={selectedWalletId === wallet.id}
                                            onChange={e => setSelectedWalletId(e.target.value)}
                                            className="radio text-emerald-600"
                                        />
                                        <div>
                                            <p className="font-semibold text-gray-900">{wallet.alias || 'Compte sans nom'}</p>
                                            <p className="text-xs text-gray-500">{wallet.balance} {wallet.currency}</p>
                                        </div>
                                    </label>
                                ))}
                            </div>

                            <div className="flex justify-end gap-3">
                                <button type="button" onClick={() => setShowLinkWalletModal(false)} className="btn-secondary">Annuler</button>
                                <button type="submit" className="btn-primary" disabled={!selectedWalletId}>Confirmer</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}
