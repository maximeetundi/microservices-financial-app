'use client';

import { useEffect, useState } from 'react';
import { format } from 'date-fns';
import {
    BuildingStorefrontIcon,
    MagnifyingGlassIcon,
    XMarkIcon,
    CheckCircleIcon,
    ArrowPathIcon,
    EyeIcon,
    NoSymbolIcon,
    ExclamationTriangleIcon,
    CubeIcon,
    ShoppingCartIcon,
    CurrencyDollarIcon,
    UserGroupIcon,
} from '@heroicons/react/24/outline';

interface Shop {
    id: string;
    owner_id: string;
    owner_type: string;
    name: string;
    slug: string;
    description: string;
    logo_url: string;
    banner_url: string;
    is_public: boolean;
    wallet_id: string;
    currency: string;
    tags: string[];
    status: string;
    stats: {
        total_products: number;
        total_orders: number;
        total_revenue: number;
    };
    created_at: string;
}

// Toast Component
function Toast({ message, type, onClose }: { message: string; type: 'success' | 'error' | 'info'; onClose: () => void }) {
    useEffect(() => {
        const timer = setTimeout(onClose, 4000);
        return () => clearTimeout(timer);
    }, [onClose]);

    const colors = {
        success: 'from-green-500 to-emerald-600',
        error: 'from-red-500 to-pink-600',
        info: 'from-indigo-500 to-purple-600',
    };

    return (
        <div className="fixed bottom-6 right-6 z-50 animate-slide-up">
            <div className={`flex items-center gap-3 px-5 py-4 bg-gradient-to-r ${colors[type]} text-white rounded-xl shadow-2xl`}>
                {type === 'success' && <CheckCircleIcon className="w-5 h-5" />}
                {type === 'error' && <ExclamationTriangleIcon className="w-5 h-5" />}
                <span className="font-medium">{message}</span>
                <button onClick={onClose} className="ml-2 hover:bg-white/20 p-1 rounded-lg transition-colors">
                    <XMarkIcon className="w-4 h-4" />
                </button>
            </div>
        </div>
    );
}

// Modal Component
function Modal({ isOpen, onClose, title, children, size = 'md' }: { isOpen: boolean; onClose: () => void; title: string; children: React.ReactNode; size?: 'sm' | 'md' | 'lg' | 'xl' }) {
    if (!isOpen) return null;
    const sizeClasses = { sm: 'max-w-sm', md: 'max-w-md', lg: 'max-w-2xl', xl: 'max-w-4xl' };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div className="absolute inset-0 bg-black/50 backdrop-blur-sm" onClick={onClose} />
            <div className={`relative bg-white rounded-2xl shadow-2xl w-full ${sizeClasses[size]} max-h-[90vh] flex flex-col overflow-hidden`}>
                <div className="flex items-center justify-between p-5 border-b border-gray-100">
                    <h3 className="text-xl font-bold text-gray-900">{title}</h3>
                    <button onClick={onClose} className="p-2 rounded-lg hover:bg-gray-100 transition-colors">
                        <XMarkIcon className="w-5 h-5 text-gray-500" />
                    </button>
                </div>
                <div className="p-5 overflow-y-auto">{children}</div>
            </div>
        </div>
    );
}

type FilterTab = 'all' | 'active' | 'suspended' | 'public' | 'private';

export default function ShopsPage() {
    const [shops, setShops] = useState<Shop[]>([]);
    const [loading, setLoading] = useState(true);
    const [actionLoading, setActionLoading] = useState<string | null>(null);
    const [searchQuery, setSearchQuery] = useState('');
    const [activeTab, setActiveTab] = useState<FilterTab>('all');
    const [detailsModal, setDetailsModal] = useState<{ isOpen: boolean; shop: Shop | null }>({ isOpen: false, shop: null });
    const [suspendModal, setSuspendModal] = useState<{ isOpen: boolean; shop: Shop | null }>({ isOpen: false, shop: null });
    const [suspendReason, setSuspendReason] = useState('');
    const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' | 'info' } | null>(null);

    const showToast = (message: string, type: 'success' | 'error' | 'info') => setToast({ message, type });

    const getAuthToken = () => {
        if (typeof window !== 'undefined') {
            return localStorage.getItem('adminToken');
        }
        return null;
    };

    const fetchShops = async () => {
        try {
            const token = getAuthToken();
            const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/shop-service/api/v1/shops?page=1&page_size=100`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            const data = await response.json();
            setShops(data.shops || []);
        } catch (error) {
            console.error('Failed to fetch shops:', error);
            // Demo data for testing
            setShops([]);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => { fetchShops(); }, []);

    const handleSuspend = async () => {
        if (!suspendModal.shop || !suspendReason.trim()) return;
        setActionLoading(suspendModal.shop.id);
        try {
            const token = getAuthToken();
            await fetch(`${process.env.NEXT_PUBLIC_API_URL}/shop-service/api/v1/shops/${suspendModal.shop.id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({ status: 'suspended' })
            });
            await fetchShops();
            setSuspendModal({ isOpen: false, shop: null });
            setSuspendReason('');
            showToast('Boutique suspendue', 'success');
        } catch (error) {
            showToast('Erreur lors de la suspension', 'error');
        } finally {
            setActionLoading(null);
        }
    };

    const handleReactivate = async (shop: Shop) => {
        setActionLoading(shop.id);
        try {
            const token = getAuthToken();
            await fetch(`${process.env.NEXT_PUBLIC_API_URL}/shop-service/api/v1/shops/${shop.id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({ status: 'active' })
            });
            await fetchShops();
            showToast('Boutique réactivée', 'success');
        } catch (error) {
            showToast('Erreur lors de la réactivation', 'error');
        } finally {
            setActionLoading(null);
        }
    };

    // Filtering
    const filteredShops = shops.filter(shop => {
        const matchesSearch =
            shop.name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
            shop.slug?.toLowerCase().includes(searchQuery.toLowerCase()) ||
            shop.tags?.some(t => t.toLowerCase().includes(searchQuery.toLowerCase()));

        switch (activeTab) {
            case 'active': return matchesSearch && shop.status === 'active';
            case 'suspended': return matchesSearch && shop.status === 'suspended';
            case 'public': return matchesSearch && shop.is_public;
            case 'private': return matchesSearch && !shop.is_public;
            default: return matchesSearch;
        }
    });

    // Stats
    const activeShops = shops.filter(s => s.status === 'active').length;
    const suspendedShops = shops.filter(s => s.status === 'suspended').length;
    const publicShops = shops.filter(s => s.is_public).length;
    const totalProducts = shops.reduce((sum, s) => sum + (s.stats?.total_products || 0), 0);
    const totalOrders = shops.reduce((sum, s) => sum + (s.stats?.total_orders || 0), 0);

    const tabs = [
        { key: 'all' as FilterTab, label: 'Toutes', count: shops.length },
        { key: 'active' as FilterTab, label: 'Actives', count: activeShops },
        { key: 'suspended' as FilterTab, label: 'Suspendues', count: suspendedShops },
        { key: 'public' as FilterTab, label: 'Publiques', count: publicShops },
        { key: 'private' as FilterTab, label: 'Privées', count: shops.length - publicShops },
    ];

    const formatPrice = (amount: number, currency: string) => {
        return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: currency || 'XOF' }).format(amount);
    };

    if (loading) {
        return (
            <div className="flex items-center justify-center h-64">
                <div className="w-12 h-12 rounded-full border-4 border-indigo-100 border-t-indigo-500 animate-spin"></div>
            </div>
        );
    }

    return (
        <div className="space-y-6">
            {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}

            {/* Header */}
            <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900">🛍️ Boutiques</h1>
                    <p className="text-gray-500 text-sm">{shops.length} boutiques enregistrées</p>
                </div>
                <button onClick={() => fetchShops()} className="flex items-center gap-2 px-4 py-2 bg-white border border-gray-200 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors text-sm">
                    <ArrowPathIcon className="w-4 h-4" />
                    Actualiser
                </button>
            </div>

            {/* Stats Mini */}
            <div className="flex gap-4 overflow-x-auto pb-2">
                <div className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-indigo-500 to-purple-600 text-white rounded-xl">
                    <BuildingStorefrontIcon className="w-5 h-5" />
                    <span className="font-bold">{shops.length}</span>
                    <span className="text-white/70 text-sm">Boutiques</span>
                </div>
                <div className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-green-500 to-emerald-600 text-white rounded-xl">
                    <CheckCircleIcon className="w-5 h-5" />
                    <span className="font-bold">{activeShops}</span>
                    <span className="text-white/70 text-sm">Actives</span>
                </div>
                <div className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-cyan-500 to-blue-600 text-white rounded-xl">
                    <CubeIcon className="w-5 h-5" />
                    <span className="font-bold">{totalProducts}</span>
                    <span className="text-white/70 text-sm">Produits</span>
                </div>
                <div className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-amber-500 to-orange-600 text-white rounded-xl">
                    <ShoppingCartIcon className="w-5 h-5" />
                    <span className="font-bold">{totalOrders}</span>
                    <span className="text-white/70 text-sm">Commandes</span>
                </div>
            </div>

            {/* Tabs & Search */}
            <div className="flex flex-col lg:flex-row gap-4 items-start lg:items-center justify-between">
                <div className="flex gap-2 overflow-x-auto">
                    {tabs.map(tab => (
                        <button
                            key={tab.key}
                            onClick={() => setActiveTab(tab.key)}
                            className={`px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-all ${activeTab === tab.key
                                ? 'bg-indigo-600 text-white'
                                : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                                }`}
                        >
                            {tab.label} ({tab.count})
                        </button>
                    ))}
                </div>
                <div className="relative w-full lg:w-72">
                    <MagnifyingGlassIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                    <input
                        type="text"
                        placeholder="Rechercher..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="w-full pl-10 pr-4 py-2.5 bg-white border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none text-sm"
                    />
                </div>
            </div>

            {/* Shops List */}
            <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
                <div className="grid grid-cols-12 gap-4 px-4 py-3 bg-gray-50 border-b border-gray-200 text-xs font-semibold text-gray-500 uppercase">
                    <div className="col-span-4">Boutique</div>
                    <div className="col-span-2 text-center">Statut</div>
                    <div className="col-span-2 text-center">Produits</div>
                    <div className="col-span-2 text-center">Revenus</div>
                    <div className="col-span-2 text-right">Actions</div>
                </div>
                <div className="divide-y divide-gray-100">
                    {filteredShops.map((shop) => (
                        <div key={shop.id} className="grid grid-cols-12 gap-4 px-4 py-3 items-center hover:bg-gray-50 transition-colors">
                            {/* Shop Info */}
                            <div className="col-span-4 flex items-center gap-3">
                                <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white text-lg font-bold flex-shrink-0 overflow-hidden">
                                    {shop.logo_url ? (
                                        <img src={shop.logo_url} alt="" className="w-full h-full object-cover" />
                                    ) : (
                                        shop.name?.[0]
                                    )}
                                </div>
                                <div className="min-w-0">
                                    <p className="font-medium text-gray-900 truncate">{shop.name}</p>
                                    <p className="text-xs text-gray-500 truncate">/{shop.slug}</p>
                                </div>
                            </div>

                            {/* Status */}
                            <div className="col-span-2 flex justify-center">
                                <div className="flex items-center gap-1.5">
                                    <span className={`w-2 h-2 rounded-full ${shop.status === 'active' ? 'bg-green-500' : shop.status === 'suspended' ? 'bg-red-500' : 'bg-gray-400'}`}></span>
                                    <span className={`text-xs font-medium ${shop.status === 'active' ? 'text-green-600' : shop.status === 'suspended' ? 'text-red-600' : 'text-gray-600'}`}>
                                        {shop.status === 'active' ? 'Active' : shop.status === 'suspended' ? 'Suspendue' : 'Fermée'}
                                    </span>
                                </div>
                            </div>

                            {/* Products */}
                            <div className="col-span-2 text-center">
                                <span className="font-semibold text-gray-900">{shop.stats?.total_products || 0}</span>
                            </div>

                            {/* Revenue */}
                            <div className="col-span-2 text-center">
                                <span className="font-semibold text-green-600">{formatPrice(shop.stats?.total_revenue || 0, shop.currency)}</span>
                            </div>

                            {/* Actions */}
                            <div className="col-span-2 flex justify-end gap-1">
                                <button
                                    onClick={() => setDetailsModal({ isOpen: true, shop })}
                                    className="p-2 text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors"
                                    title="Voir les détails"
                                >
                                    <EyeIcon className="w-4 h-4" />
                                </button>
                                {shop.status === 'active' ? (
                                    <button
                                        onClick={() => setSuspendModal({ isOpen: true, shop })}
                                        className="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                                        title="Suspendre"
                                    >
                                        <NoSymbolIcon className="w-4 h-4" />
                                    </button>
                                ) : (
                                    <button
                                        onClick={() => handleReactivate(shop)}
                                        disabled={actionLoading === shop.id}
                                        className="p-2 text-green-600 hover:bg-green-50 rounded-lg transition-colors disabled:opacity-50"
                                        title="Réactiver"
                                    >
                                        <CheckCircleIcon className="w-4 h-4" />
                                    </button>
                                )}
                            </div>
                        </div>
                    ))}
                </div>
                {filteredShops.length === 0 && (
                    <div className="text-center py-12">
                        <BuildingStorefrontIcon className="w-10 h-10 text-gray-300 mx-auto mb-3" />
                        <p className="text-gray-500">Aucune boutique trouvée</p>
                    </div>
                )}
            </div>

            {/* Shop Details Modal */}
            <Modal
                isOpen={detailsModal.isOpen}
                onClose={() => setDetailsModal({ isOpen: false, shop: null })}
                title="Détails de la boutique"
                size="lg"
            >
                {detailsModal.shop && (
                    <div className="space-y-6">
                        {/* Shop Header */}
                        <div className="flex items-center gap-4 p-4 bg-gradient-to-r from-indigo-50 to-purple-50 rounded-xl">
                            <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white text-2xl font-bold overflow-hidden">
                                {detailsModal.shop.logo_url ? (
                                    <img src={detailsModal.shop.logo_url} alt="" className="w-full h-full object-cover" />
                                ) : (
                                    detailsModal.shop.name?.[0]
                                )}
                            </div>
                            <div>
                                <h3 className="text-xl font-bold text-gray-900">{detailsModal.shop.name}</h3>
                                <div className="flex items-center gap-3 mt-1">
                                    <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${detailsModal.shop.status === 'active' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                                        {detailsModal.shop.status === 'active' ? '✓ Active' : '✗ Suspendue'}
                                    </span>
                                    <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${detailsModal.shop.is_public ? 'bg-indigo-100 text-indigo-700' : 'bg-gray-100 text-gray-600'}`}>
                                        {detailsModal.shop.is_public ? '🌍 Publique' : '🔒 Privée'}
                                    </span>
                                </div>
                            </div>
                        </div>

                        {/* Description */}
                        {detailsModal.shop.description && (
                            <div className="p-4 bg-gray-50 rounded-xl">
                                <p className="text-gray-700">{detailsModal.shop.description}</p>
                            </div>
                        )}

                        {/* Stats */}
                        <div className="grid grid-cols-3 gap-4">
                            <div className="p-4 bg-gradient-to-br from-cyan-50 to-blue-50 rounded-xl text-center">
                                <CubeIcon className="w-8 h-8 text-cyan-600 mx-auto mb-2" />
                                <p className="text-2xl font-bold text-gray-900">{detailsModal.shop.stats?.total_products || 0}</p>
                                <p className="text-sm text-gray-500">Produits</p>
                            </div>
                            <div className="p-4 bg-gradient-to-br from-amber-50 to-orange-50 rounded-xl text-center">
                                <ShoppingCartIcon className="w-8 h-8 text-amber-600 mx-auto mb-2" />
                                <p className="text-2xl font-bold text-gray-900">{detailsModal.shop.stats?.total_orders || 0}</p>
                                <p className="text-sm text-gray-500">Commandes</p>
                            </div>
                            <div className="p-4 bg-gradient-to-br from-green-50 to-emerald-50 rounded-xl text-center">
                                <CurrencyDollarIcon className="w-8 h-8 text-green-600 mx-auto mb-2" />
                                <p className="text-2xl font-bold text-gray-900">{formatPrice(detailsModal.shop.stats?.total_revenue || 0, detailsModal.shop.currency)}</p>
                                <p className="text-sm text-gray-500">Revenus</p>
                            </div>
                        </div>

                        {/* Info */}
                        <div className="grid grid-cols-2 gap-4">
                            <div className="p-3 bg-gray-50 rounded-lg">
                                <p className="text-xs text-gray-500">Slug</p>
                                <p className="font-medium text-gray-900">/{detailsModal.shop.slug}</p>
                            </div>
                            <div className="p-3 bg-gray-50 rounded-lg">
                                <p className="text-xs text-gray-500">Devise</p>
                                <p className="font-medium text-gray-900">{detailsModal.shop.currency}</p>
                            </div>
                            <div className="p-3 bg-gray-50 rounded-lg">
                                <p className="text-xs text-gray-500">Type propriétaire</p>
                                <p className="font-medium text-gray-900">{detailsModal.shop.owner_type === 'user' ? 'Utilisateur' : 'Entreprise'}</p>
                            </div>
                            <div className="p-3 bg-gray-50 rounded-lg">
                                <p className="text-xs text-gray-500">Créée le</p>
                                <p className="font-medium text-gray-900">
                                    {detailsModal.shop.created_at ? format(new Date(detailsModal.shop.created_at), 'dd/MM/yyyy') : '—'}
                                </p>
                            </div>
                        </div>

                        {/* Tags */}
                        {detailsModal.shop.tags?.length > 0 && (
                            <div>
                                <p className="text-sm font-medium text-gray-700 mb-2">Tags</p>
                                <div className="flex flex-wrap gap-2">
                                    {detailsModal.shop.tags.map(tag => (
                                        <span key={tag} className="px-3 py-1 bg-indigo-50 text-indigo-700 rounded-full text-sm">{tag}</span>
                                    ))}
                                </div>
                            </div>
                        )}

                        {/* Actions */}
                        <div className="flex gap-3 pt-4 border-t border-gray-100">
                            {detailsModal.shop.status === 'active' ? (
                                <button
                                    onClick={() => { setDetailsModal({ isOpen: false, shop: null }); setSuspendModal({ isOpen: true, shop: detailsModal.shop }); }}
                                    className="flex-1 flex items-center justify-center gap-2 py-2.5 bg-red-50 text-red-700 rounded-lg font-medium hover:bg-red-100 transition-colors"
                                >
                                    <NoSymbolIcon className="w-4 h-4" />
                                    Suspendre
                                </button>
                            ) : (
                                <button
                                    onClick={() => { handleReactivate(detailsModal.shop!); }}
                                    className="flex-1 flex items-center justify-center gap-2 py-2.5 bg-green-50 text-green-700 rounded-lg font-medium hover:bg-green-100 transition-colors"
                                >
                                    <CheckCircleIcon className="w-4 h-4" />
                                    Réactiver
                                </button>
                            )}
                        </div>
                    </div>
                )}
            </Modal>

            {/* Suspend Shop Modal */}
            <Modal
                isOpen={suspendModal.isOpen}
                onClose={() => { setSuspendModal({ isOpen: false, shop: null }); setSuspendReason(''); }}
                title="Suspendre la boutique"
            >
                <div className="space-y-4">
                    <div className="flex items-center gap-3 p-4 bg-red-50 rounded-xl">
                        <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-semibold overflow-hidden">
                            {suspendModal.shop?.logo_url ? (
                                <img src={suspendModal.shop.logo_url} alt="" className="w-full h-full object-cover" />
                            ) : (
                                suspendModal.shop?.name?.[0]
                            )}
                        </div>
                        <div>
                            <p className="font-semibold text-gray-900">{suspendModal.shop?.name}</p>
                            <p className="text-sm text-gray-500">/{suspendModal.shop?.slug}</p>
                        </div>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Raison de la suspension <span className="text-red-500">*</span>
                        </label>
                        <textarea
                            value={suspendReason}
                            onChange={(e) => setSuspendReason(e.target.value)}
                            placeholder="Décrivez la raison..."
                            rows={3}
                            className="w-full px-4 py-3 rounded-xl border border-gray-200 focus:border-red-500 focus:ring-2 focus:ring-red-500/20 outline-none transition-all resize-none"
                        />
                    </div>
                    <div className="flex gap-3">
                        <button
                            onClick={() => { setSuspendModal({ isOpen: false, shop: null }); setSuspendReason(''); }}
                            className="flex-1 py-2.5 rounded-xl font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 transition-colors"
                        >
                            Annuler
                        </button>
                        <button
                            onClick={handleSuspend}
                            disabled={!suspendReason.trim() || actionLoading === suspendModal.shop?.id}
                            className="flex-1 py-2.5 rounded-xl font-medium text-white bg-red-600 hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
                        >
                            {actionLoading === suspendModal.shop?.id ? 'Suspension...' : 'Confirmer'}
                        </button>
                    </div>
                </div>
            </Modal>
        </div>
    );
}
