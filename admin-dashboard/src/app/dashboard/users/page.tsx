'use client';

import { useEffect, useState } from 'react';
import { getUsers, blockUser, unblockUser, unlockUserPin, getUserKYCDocuments, getKYCDocumentURL, getKYCDownloadURL } from '@/lib/api';
import { format } from 'date-fns';
import {
    UsersIcon,
    ShieldCheckIcon,
    NoSymbolIcon,
    KeyIcon,
    MagnifyingGlassIcon,
    XMarkIcon,
    ExclamationTriangleIcon,
    LockClosedIcon,
    CheckCircleIcon,
    ArrowPathIcon,
    LockOpenIcon,
    DocumentIcon,
    EyeIcon,
    ArrowDownTrayIcon,
    IdentificationIcon,
    UserCircleIcon,
    EnvelopeIcon,
    PhoneIcon,
    CalendarIcon,
    MapPinIcon,
} from '@heroicons/react/24/outline';

interface User {
    id: string;
    email: string;
    first_name: string;
    last_name: string;
    phone: string;
    is_active: boolean;
    kyc_level: number | string;
    kyc_status?: string;
    created_at: string;
    pin_locked?: boolean;
    country?: string;
    address?: string;
    date_of_birth?: string;
}

interface KYCDocument {
    id: string;
    type: string;
    status: string;
    file_url?: string;
    file_path?: string;
    file_name?: string;
    uploaded_at: string;
    identity_sub_type?: string;
    document_number?: string;
    expiry_date?: string;
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

// Secure Document Viewer Component
function SecureDocumentViewer({ doc }: { doc: KYCDocument }) {
    const [loading, setLoading] = useState(false);
    const [imageUrl, setImageUrl] = useState<string | null>(null);

    useEffect(() => {
        const loadUrl = async () => {
            setLoading(true);
            try {
                const response = await getKYCDocumentURL(doc.id);
                setImageUrl(response.data.url);
            } catch (err) {
                console.error('Failed to load document:', err);
            } finally {
                setLoading(false);
            }
        };
        loadUrl();
    }, [doc.id]);

    const handleDownload = async () => {
        try {
            const response = await getKYCDownloadURL(doc.id);
            window.open(response.data.url, '_blank');
        } catch (err) {
            console.error('Failed to download:', err);
        }
    };

    const isPDF = doc.file_name?.toLowerCase().endsWith('.pdf');

    return (
        <div className="space-y-3">
            {loading && (
                <div className="flex justify-center py-6">
                    <div className="w-8 h-8 border-4 border-indigo-100 border-t-indigo-500 rounded-full animate-spin"></div>
                </div>
            )}
            {imageUrl && !isPDF && (
                <img src={imageUrl} alt={doc.file_name} className="w-full rounded-lg max-h-48 object-contain bg-gray-100" />
            )}
            {imageUrl && isPDF && (
                <div className="flex items-center justify-center gap-3 p-4 bg-gray-50 rounded-lg">
                    <DocumentIcon className="w-10 h-10 text-red-500" />
                    <span className="text-gray-600">{doc.file_name || 'Document PDF'}</span>
                </div>
            )}
            <button onClick={handleDownload} className="w-full flex items-center justify-center gap-2 px-3 py-2 bg-indigo-50 hover:bg-indigo-100 text-indigo-700 rounded-lg text-sm font-medium transition-colors">
                <ArrowDownTrayIcon className="w-4 h-4" />
                Télécharger
            </button>
        </div>
    );
}

type FilterTab = 'all' | 'active' | 'blocked' | 'kyc_pending';

export default function UsersPage() {
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(true);
    const [actionLoading, setActionLoading] = useState<string | null>(null);
    const [searchQuery, setSearchQuery] = useState('');
    const [activeTab, setActiveTab] = useState<FilterTab>('all');

    // Modal states
    const [detailsModal, setDetailsModal] = useState<{ isOpen: boolean; user: User | null; documents: KYCDocument[] }>({ isOpen: false, user: null, documents: [] });
    const [blockModal, setBlockModal] = useState<{ isOpen: boolean; user: User | null }>({ isOpen: false, user: null });
    const [blockReason, setBlockReason] = useState('');
    const [loadingDocs, setLoadingDocs] = useState(false);
    const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' | 'info' } | null>(null);

    const showToast = (message: string, type: 'success' | 'error' | 'info') => setToast({ message, type });

    const fetchUsers = async () => {
        try {
            const response = await getUsers();
            setUsers(response.data.users || []);
        } catch (error) {
            console.error('Failed to fetch users:', error);
        } finally {
            setLoading(false);
        }
    };

    const openUserDetails = async (user: User) => {
        setDetailsModal({ isOpen: true, user, documents: [] });
        setLoadingDocs(true);
        try {
            const response = await getUserKYCDocuments(user.id);
            setDetailsModal(prev => ({ ...prev, documents: response.data.documents || [] }));
        } catch (error) {
            console.error('Failed to fetch documents:', error);
        } finally {
            setLoadingDocs(false);
        }
    };

    useEffect(() => { fetchUsers(); }, []);

    const handleBlock = async () => {
        if (!blockModal.user || !blockReason.trim()) return;
        setActionLoading(blockModal.user.id);
        try {
            await blockUser(blockModal.user.id, blockReason);
            await fetchUsers();
            setBlockModal({ isOpen: false, user: null });
            setBlockReason('');
            showToast(`Compte bloqué avec succès`, 'success');
        } catch (error) {
            showToast('Erreur lors du blocage', 'error');
        } finally {
            setActionLoading(null);
        }
    };

    const handleUnblock = async (user: User) => {
        setActionLoading(user.id);
        try {
            await unblockUser(user.id);
            await fetchUsers();
            showToast(`Compte débloqué`, 'success');
        } catch (error) {
            showToast('Erreur lors du déblocage', 'error');
        } finally {
            setActionLoading(null);
        }
    };

    const handleUnlockPin = async (user: User) => {
        setActionLoading(`pin-${user.id}`);
        try {
            await unlockUserPin(user.id);
            await fetchUsers();
            showToast(`PIN débloqué`, 'success');
        } catch (error) {
            showToast('Erreur lors du déblocage du PIN', 'error');
        } finally {
            setActionLoading(null);
        }
    };

    // Filtering
    const filteredUsers = users.filter(user => {
        const matchesSearch =
            user.email?.toLowerCase().includes(searchQuery.toLowerCase()) ||
            user.first_name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
            user.last_name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
            user.phone?.includes(searchQuery);

        switch (activeTab) {
            case 'active': return matchesSearch && user.is_active;
            case 'blocked': return matchesSearch && !user.is_active;
            case 'kyc_pending': return matchesSearch && user.kyc_status === 'pending';
            default: return matchesSearch;
        }
    });

    // Stats
    const activeUsers = users.filter(u => u.is_active).length;
    const blockedUsers = users.filter(u => !u.is_active).length;
    const kycPending = users.filter(u => u.kyc_status === 'pending').length;

    const tabs = [
        { key: 'all' as FilterTab, label: 'Tous', count: users.length },
        { key: 'active' as FilterTab, label: 'Actifs', count: activeUsers },
        { key: 'blocked' as FilterTab, label: 'Bloqués', count: blockedUsers },
        { key: 'kyc_pending' as FilterTab, label: 'KYC En attente', count: kycPending },
    ];

    const getKYCBadge = (user: User) => {
        const status = user.kyc_status || (String(user.kyc_level) === '3' ? 'verified' : 'none');
        switch (status) {
            case 'verified': return <span className="w-2 h-2 rounded-full bg-green-500" title="KYC Vérifié"></span>;
            case 'pending': return <span className="w-2 h-2 rounded-full bg-amber-500" title="KYC En attente"></span>;
            case 'rejected': return <span className="w-2 h-2 rounded-full bg-red-500" title="KYC Refusé"></span>;
            default: return <span className="w-2 h-2 rounded-full bg-gray-300" title="Non vérifié"></span>;
        }
    };

    const getIdentitySubTypeName = (subType?: string) => {
        switch (subType) {
            case 'cni': return "Carte Nationale d'Identité";
            case 'passport': return 'Passeport';
            case 'permis': return 'Permis de Conduire';
            default: return null;
        }
    };

    const getDocTypeName = (type: string) => {
        switch (type) {
            case 'identity': return "Pièce d'identité";
            case 'selfie': return 'Selfie';
            case 'address': return 'Justificatif de domicile';
            default: return 'Document';
        }
    };

    const formatExpiryDate = (dateString?: string) => {
        if (!dateString) return null;
        try {
            return new Date(dateString).toLocaleDateString('fr-FR');
        } catch {
            return dateString;
        }
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
                    <h1 className="text-2xl font-bold text-gray-900">👥 Utilisateurs</h1>
                    <p className="text-gray-500 text-sm">{users.length} utilisateurs enregistrés</p>
                </div>
                <button onClick={() => fetchUsers()} className="flex items-center gap-2 px-4 py-2 bg-white border border-gray-200 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors text-sm">
                    <ArrowPathIcon className="w-4 h-4" />
                    Actualiser
                </button>
            </div>

            {/* Stats Mini */}
            <div className="flex gap-4 overflow-x-auto pb-2">
                <div className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-indigo-500 to-purple-600 text-white rounded-xl">
                    <UsersIcon className="w-5 h-5" />
                    <span className="font-bold">{users.length}</span>
                    <span className="text-white/70 text-sm">Total</span>
                </div>
                <div className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-green-500 to-emerald-600 text-white rounded-xl">
                    <CheckCircleIcon className="w-5 h-5" />
                    <span className="font-bold">{activeUsers}</span>
                    <span className="text-white/70 text-sm">Actifs</span>
                </div>
                <div className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-amber-500 to-orange-600 text-white rounded-xl">
                    <IdentificationIcon className="w-5 h-5" />
                    <span className="font-bold">{kycPending}</span>
                    <span className="text-white/70 text-sm">KYC</span>
                </div>
                <div className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-red-500 to-pink-600 text-white rounded-xl">
                    <NoSymbolIcon className="w-5 h-5" />
                    <span className="font-bold">{blockedUsers}</span>
                    <span className="text-white/70 text-sm">Bloqués</span>
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

            {/* Users List - Simplified */}
            <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
                <div className="grid grid-cols-12 gap-4 px-4 py-3 bg-gray-50 border-b border-gray-200 text-xs font-semibold text-gray-500 uppercase">
                    <div className="col-span-5">Utilisateur</div>
                    <div className="col-span-2 text-center">Statut</div>
                    <div className="col-span-2 text-center">KYC</div>
                    <div className="col-span-3 text-right">Actions</div>
                </div>
                <div className="divide-y divide-gray-100">
                    {filteredUsers.map((user) => (
                        <div key={user.id} className="grid grid-cols-12 gap-4 px-4 py-3 items-center hover:bg-gray-50 transition-colors">
                            {/* User Info - Minimal */}
                            <div className="col-span-5 flex items-center gap-3">
                                <div className="w-10 h-10 rounded-full bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white text-sm font-bold flex-shrink-0">
                                    {user.first_name?.[0]}{user.last_name?.[0]}
                                </div>
                                <div className="min-w-0">
                                    <p className="font-medium text-gray-900 truncate">{user.first_name} {user.last_name}</p>
                                    <p className="text-xs text-gray-500 truncate">{user.email}</p>
                                </div>
                            </div>

                            {/* Status */}
                            <div className="col-span-2 flex justify-center">
                                <div className="flex items-center gap-1.5">
                                    <span className={`w-2 h-2 rounded-full ${user.is_active ? 'bg-green-500' : 'bg-red-500'}`}></span>
                                    <span className={`text-xs font-medium ${user.is_active ? 'text-green-600' : 'text-red-600'}`}>
                                        {user.is_active ? 'Actif' : 'Bloqué'}
                                    </span>
                                </div>
                            </div>

                            {/* KYC */}
                            <div className="col-span-2 flex justify-center">
                                {getKYCBadge(user)}
                            </div>

                            {/* Actions */}
                            <div className="col-span-3 flex justify-end gap-1">
                                <button
                                    onClick={() => openUserDetails(user)}
                                    className="p-2 text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors"
                                    title="Voir les détails"
                                >
                                    <EyeIcon className="w-4 h-4" />
                                </button>
                                {user.pin_locked && (
                                    <button
                                        onClick={() => handleUnlockPin(user)}
                                        disabled={actionLoading === `pin-${user.id}`}
                                        className="p-2 text-amber-600 hover:bg-amber-50 rounded-lg transition-colors disabled:opacity-50"
                                        title="Débloquer PIN"
                                    >
                                        <LockOpenIcon className="w-4 h-4" />
                                    </button>
                                )}
                                {user.is_active ? (
                                    <button
                                        onClick={() => setBlockModal({ isOpen: true, user })}
                                        className="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                                        title="Bloquer"
                                    >
                                        <NoSymbolIcon className="w-4 h-4" />
                                    </button>
                                ) : (
                                    <button
                                        onClick={() => handleUnblock(user)}
                                        disabled={actionLoading === user.id}
                                        className="p-2 text-green-600 hover:bg-green-50 rounded-lg transition-colors disabled:opacity-50"
                                        title="Débloquer"
                                    >
                                        <CheckCircleIcon className="w-4 h-4" />
                                    </button>
                                )}
                            </div>
                        </div>
                    ))}
                </div>
                {filteredUsers.length === 0 && (
                    <div className="text-center py-12">
                        <UsersIcon className="w-10 h-10 text-gray-300 mx-auto mb-3" />
                        <p className="text-gray-500">Aucun utilisateur trouvé</p>
                    </div>
                )}
            </div>

            {/* User Details Modal */}
            <Modal
                isOpen={detailsModal.isOpen}
                onClose={() => setDetailsModal({ isOpen: false, user: null, documents: [] })}
                title="Détails de l'utilisateur"
                size="lg"
            >
                {detailsModal.user && (
                    <div className="space-y-6">
                        {/* User Info */}
                        <div className="flex items-center gap-4 p-4 bg-gradient-to-r from-indigo-50 to-purple-50 rounded-xl">
                            <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white text-2xl font-bold">
                                {detailsModal.user.first_name?.[0]}{detailsModal.user.last_name?.[0]}
                            </div>
                            <div>
                                <h3 className="text-xl font-bold text-gray-900">{detailsModal.user.first_name} {detailsModal.user.last_name}</h3>
                                <div className="flex items-center gap-3 mt-1">
                                    <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${detailsModal.user.is_active ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                                        {detailsModal.user.is_active ? '✓ Actif' : '✗ Bloqué'}
                                    </span>
                                    <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${detailsModal.user.kyc_status === 'verified' ? 'bg-green-100 text-green-700' :
                                            detailsModal.user.kyc_status === 'pending' ? 'bg-amber-100 text-amber-700' :
                                                'bg-gray-100 text-gray-600'
                                        }`}>
                                        KYC: {detailsModal.user.kyc_status === 'verified' ? 'Vérifié' : detailsModal.user.kyc_status === 'pending' ? 'En attente' : 'Non vérifié'}
                                    </span>
                                </div>
                            </div>
                        </div>

                        {/* Contact Info */}
                        <div className="grid grid-cols-2 gap-4">
                            <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
                                <EnvelopeIcon className="w-5 h-5 text-gray-400" />
                                <div>
                                    <p className="text-xs text-gray-500">Email</p>
                                    <p className="text-sm font-medium text-gray-900">{detailsModal.user.email}</p>
                                </div>
                            </div>
                            <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
                                <PhoneIcon className="w-5 h-5 text-gray-400" />
                                <div>
                                    <p className="text-xs text-gray-500">Téléphone</p>
                                    <p className="text-sm font-medium text-gray-900">{detailsModal.user.phone || '—'}</p>
                                </div>
                            </div>
                            <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
                                <CalendarIcon className="w-5 h-5 text-gray-400" />
                                <div>
                                    <p className="text-xs text-gray-500">Inscrit le</p>
                                    <p className="text-sm font-medium text-gray-900">
                                        {detailsModal.user.created_at ? format(new Date(detailsModal.user.created_at), 'dd/MM/yyyy') : '—'}
                                    </p>
                                </div>
                            </div>
                            <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
                                <MapPinIcon className="w-5 h-5 text-gray-400" />
                                <div>
                                    <p className="text-xs text-gray-500">Pays</p>
                                    <p className="text-sm font-medium text-gray-900">{detailsModal.user.country || '—'}</p>
                                </div>
                            </div>
                        </div>

                        {/* KYC Documents */}
                        <div>
                            <h4 className="font-semibold text-gray-900 mb-3 flex items-center gap-2">
                                <IdentificationIcon className="w-5 h-5 text-indigo-600" />
                                Documents KYC
                            </h4>

                            {loadingDocs ? (
                                <div className="flex justify-center py-8">
                                    <div className="w-8 h-8 border-4 border-indigo-100 border-t-indigo-500 rounded-full animate-spin"></div>
                                </div>
                            ) : detailsModal.documents.length === 0 ? (
                                <div className="text-center py-8 bg-gray-50 rounded-xl">
                                    <DocumentIcon className="w-10 h-10 text-gray-300 mx-auto mb-2" />
                                    <p className="text-gray-500 text-sm">Aucun document soumis</p>
                                </div>
                            ) : (
                                <div className="space-y-4">
                                    {detailsModal.documents.map((doc) => (
                                        <div key={doc.id} className="border border-gray-200 rounded-xl overflow-hidden">
                                            <div className="flex items-center justify-between p-3 bg-gray-50 border-b border-gray-200">
                                                <div className="flex items-center gap-2">
                                                    <span className="text-xl">{doc.type === 'identity' ? '🪪' : doc.type === 'selfie' ? '🤳' : '🏠'}</span>
                                                    <div>
                                                        <p className="font-medium text-gray-900 text-sm">{getDocTypeName(doc.type)}</p>
                                                        {doc.type === 'identity' && doc.identity_sub_type && (
                                                            <p className="text-xs text-indigo-600">{getIdentitySubTypeName(doc.identity_sub_type)}</p>
                                                        )}
                                                    </div>
                                                </div>
                                                <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${doc.status === 'approved' ? 'bg-green-100 text-green-700' :
                                                        doc.status === 'rejected' ? 'bg-red-100 text-red-700' :
                                                            'bg-amber-100 text-amber-700'
                                                    }`}>
                                                    {doc.status === 'approved' ? 'Approuvé' : doc.status === 'rejected' ? 'Rejeté' : 'En attente'}
                                                </span>
                                            </div>

                                            {/* Metadata */}
                                            {doc.type === 'identity' && (doc.document_number || doc.expiry_date) && (
                                                <div className="px-3 py-2 bg-indigo-50/50 border-b border-indigo-100 flex gap-4 text-xs">
                                                    {doc.document_number && (
                                                        <span><strong>N°:</strong> {doc.document_number}</span>
                                                    )}
                                                    {doc.expiry_date && (
                                                        <span><strong>Expire:</strong> {formatExpiryDate(doc.expiry_date)}</span>
                                                    )}
                                                </div>
                                            )}

                                            <div className="p-3">
                                                <SecureDocumentViewer doc={doc} />
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>

                        {/* Quick Actions */}
                        <div className="flex gap-3 pt-4 border-t border-gray-100">
                            {detailsModal.user.pin_locked && (
                                <button
                                    onClick={() => { handleUnlockPin(detailsModal.user!); }}
                                    className="flex-1 flex items-center justify-center gap-2 py-2.5 bg-amber-50 text-amber-700 rounded-lg font-medium hover:bg-amber-100 transition-colors"
                                >
                                    <LockOpenIcon className="w-4 h-4" />
                                    Débloquer PIN
                                </button>
                            )}
                            {detailsModal.user.is_active ? (
                                <button
                                    onClick={() => { setDetailsModal({ isOpen: false, user: null, documents: [] }); setBlockModal({ isOpen: true, user: detailsModal.user }); }}
                                    className="flex-1 flex items-center justify-center gap-2 py-2.5 bg-red-50 text-red-700 rounded-lg font-medium hover:bg-red-100 transition-colors"
                                >
                                    <NoSymbolIcon className="w-4 h-4" />
                                    Bloquer
                                </button>
                            ) : (
                                <button
                                    onClick={() => { handleUnblock(detailsModal.user!); }}
                                    className="flex-1 flex items-center justify-center gap-2 py-2.5 bg-green-50 text-green-700 rounded-lg font-medium hover:bg-green-100 transition-colors"
                                >
                                    <CheckCircleIcon className="w-4 h-4" />
                                    Débloquer
                                </button>
                            )}
                        </div>
                    </div>
                )}
            </Modal>

            {/* Block User Modal */}
            <Modal
                isOpen={blockModal.isOpen}
                onClose={() => { setBlockModal({ isOpen: false, user: null }); setBlockReason(''); }}
                title="Bloquer l'utilisateur"
            >
                <div className="space-y-4">
                    <div className="flex items-center gap-3 p-4 bg-red-50 rounded-xl">
                        <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-semibold">
                            {blockModal.user?.first_name?.[0]}{blockModal.user?.last_name?.[0]}
                        </div>
                        <div>
                            <p className="font-semibold text-gray-900">{blockModal.user?.first_name} {blockModal.user?.last_name}</p>
                            <p className="text-sm text-gray-500">{blockModal.user?.email}</p>
                        </div>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Raison du blocage <span className="text-red-500">*</span>
                        </label>
                        <textarea
                            value={blockReason}
                            onChange={(e) => setBlockReason(e.target.value)}
                            placeholder="Décrivez la raison..."
                            rows={3}
                            className="w-full px-4 py-3 rounded-xl border border-gray-200 focus:border-red-500 focus:ring-2 focus:ring-red-500/20 outline-none transition-all resize-none"
                        />
                    </div>
                    <div className="flex gap-3">
                        <button
                            onClick={() => { setBlockModal({ isOpen: false, user: null }); setBlockReason(''); }}
                            className="flex-1 py-2.5 rounded-xl font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 transition-colors"
                        >
                            Annuler
                        </button>
                        <button
                            onClick={handleBlock}
                            disabled={!blockReason.trim() || actionLoading === blockModal.user?.id}
                            className="flex-1 py-2.5 rounded-xl font-medium text-white bg-red-600 hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
                        >
                            {actionLoading === blockModal.user?.id ? 'Blocage...' : 'Confirmer'}
                        </button>
                    </div>
                </div>
            </Modal>
        </div>
    );
}
