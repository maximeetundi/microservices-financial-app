'use client';

import { useEffect, useState } from 'react';
import { getUsers, approveKYC, rejectKYC, getUserKYCDocuments, getKYCDocumentURL, getKYCDownloadURL } from '@/lib/api';
import api from '@/lib/api';
import {
    XMarkIcon,
    DocumentIcon,
    CheckCircleIcon,
    XCircleIcon,
    ArrowPathIcon,
    EyeIcon,
    ShieldCheckIcon,
    ClockIcon,
    ExclamationTriangleIcon,
    UserIcon,
} from '@heroicons/react/24/outline';

interface User {
    id: string;
    email: string;
    first_name: string;
    last_name: string;
    phone: string;
    kyc_level: number;
    kyc_status: string;
    is_active: boolean;
    created_at: string;
}

interface KYCDocument {
    id: string;
    type: string;
    status: string;
    file_url?: string;
    file_path?: string;
    file_name?: string;
    file_size?: number;
    mime_type?: string;
    uploaded_at: string;
    reviewed_at?: string;
    rejection_reason?: string;
}

// Toast Component
function Toast({ message, type, onClose }: { message: string; type: 'success' | 'error'; onClose: () => void }) {
    useEffect(() => {
        const timer = setTimeout(onClose, 4000);
        return () => clearTimeout(timer);
    }, [onClose]);

    return (
        <div className="fixed bottom-6 right-6 z-50 animate-slide-up">
            <div className={`flex items-center gap-3 px-5 py-4 rounded-xl shadow-2xl text-white ${type === 'success' ? 'bg-gradient-to-r from-green-500 to-emerald-600' : 'bg-gradient-to-r from-red-500 to-pink-600'}`}>
                {type === 'success' ? <CheckCircleIcon className="w-5 h-5" /> : <XCircleIcon className="w-5 h-5" />}
                <span className="font-medium">{message}</span>
                <button onClick={onClose} className="ml-2 hover:bg-white/20 p-1 rounded-lg">
                    <XMarkIcon className="w-4 h-4" />
                </button>
            </div>
        </div>
    );
}

// Modal Component
function Modal({ isOpen, onClose, title, children, size = 'md' }: { isOpen: boolean; onClose: () => void; title: string; children: React.ReactNode; size?: string }) {
    if (!isOpen) return null;

    const sizeClass = size === 'lg' ? 'max-w-3xl' : size === 'xl' ? 'max-w-5xl' : 'max-w-md';

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div className="absolute inset-0 bg-black/50 backdrop-blur-sm" onClick={onClose} />
            <div className={`relative bg-white rounded-2xl shadow-2xl w-full ${sizeClass} animate-slide-up overflow-hidden max-h-[90vh] flex flex-col`}>
                <div className="flex items-center justify-between p-6 border-b border-gray-100">
                    <h3 className="text-xl font-bold text-gray-900">{title}</h3>
                    <button onClick={onClose} className="p-2 rounded-lg hover:bg-gray-100 transition-colors">
                        <XMarkIcon className="w-5 h-5 text-gray-500" />
                    </button>
                </div>
                <div className="p-6 overflow-y-auto flex-1">
                    {children}
                </div>
            </div>
        </div>
    );
}

// Secure Document Viewer - Loads document via presigned URL
function SecureDocumentViewer({ doc }: { doc: KYCDocument }) {
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [secureUrl, setSecureUrl] = useState<string | null>(null);

    const filePath = doc.file_url || doc.file_path;
    const isPdf = doc.mime_type?.includes('pdf') || filePath?.endsWith('.pdf');

    useEffect(() => {
        if (!filePath) return;

        const loadSecureUrl = async () => {
            setLoading(true);
            setError(null);
            try {
                const response = await getKYCDocumentURL(filePath);
                setSecureUrl(response.data.url);
            } catch (err) {
                console.error('Failed to get secure URL:', err);
                setError('Impossible de charger le document');
            } finally {
                setLoading(false);
            }
        };

        loadSecureUrl();
    }, [filePath]);

    const handleDownload = async () => {
        if (!filePath) return;
        try {
            const response = await getKYCDownloadURL(filePath, doc.file_name);
            window.open(response.data.url, '_blank');
        } catch (err) {
            console.error('Failed to get download URL:', err);
        }
    };

    if (loading) {
        return (
            <div className="p-8 flex justify-center">
                <div className="w-8 h-8 border-4 border-indigo-100 border-t-indigo-500 rounded-full animate-spin"></div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="p-4 text-center text-red-500">
                <ExclamationTriangleIcon className="w-8 h-8 mx-auto mb-2" />
                <p>{error}</p>
            </div>
        );
    }

    return (
        <div className="p-4 space-y-3">
            {secureUrl && !isPdf && (
                <img
                    src={secureUrl}
                    alt={`Document ${doc.type}`}
                    className="w-full max-h-96 object-contain rounded-lg bg-gray-100"
                />
            )}
            {secureUrl && isPdf && (
                <div className="text-center py-8 bg-gray-50 rounded-lg">
                    <DocumentIcon className="w-16 h-16 text-red-500 mx-auto mb-3" />
                    <p className="text-gray-600 font-medium">{doc.file_name || 'Document PDF'}</p>
                </div>
            )}
            <button
                onClick={handleDownload}
                className="w-full py-2 px-4 rounded-lg bg-indigo-500 hover:bg-indigo-600 text-white font-medium transition-colors flex items-center justify-center gap-2"
            >
                <DocumentIcon className="w-5 h-5" />
                Télécharger le document
            </button>
        </div>
    );
}

export default function KYCPage() {
    const [users, setUsers] = useState<User[]>([]);
    const [allUsers, setAllUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(true);
    const [actionLoading, setActionLoading] = useState<string | null>(null);
    const [filter, setFilter] = useState('pending');
    const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' } | null>(null);

    // Document viewer modal
    const [docModal, setDocModal] = useState<{ isOpen: boolean; user: User | null; documents: KYCDocument[] }>({
        isOpen: false,
        user: null,
        documents: [],
    });
    const [loadingDocs, setLoadingDocs] = useState(false);

    // Approve/Reject modal
    const [actionModal, setActionModal] = useState<{ isOpen: boolean; user: User | null; action: 'approve' | 'reject' }>({
        isOpen: false,
        user: null,
        action: 'approve',
    });
    const [rejectReason, setRejectReason] = useState('');

    const fetchUsers = async () => {
        try {
            setLoading(true);
            const response = await getUsers();
            const userList = response.data?.users || [];
            setAllUsers(userList);
        } catch (error) {
            console.error('Failed to fetch users:', error);
            setAllUsers([]);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchUsers();
    }, []);

    useEffect(() => {
        if (filter === 'all') {
            setUsers(allUsers);
        } else if (filter === 'none') {
            setUsers(allUsers.filter((u: User) => !u.kyc_status || u.kyc_status === 'none' || u.kyc_status === ''));
        } else {
            setUsers(allUsers.filter((u: User) => u.kyc_status === filter));
        }
    }, [filter, allUsers]);

    const fetchUserDocuments = async (user: User) => {
        setLoadingDocs(true);
        try {
            // Try to get KYC documents for this user
            const response = await getUserKYCDocuments(user.id);
            const docs = response.data?.documents || response.data || [];
            console.log('KYC documents loaded:', docs);
            setDocModal({ isOpen: true, user, documents: Array.isArray(docs) ? docs : [] });
        } catch (error) {
            console.error('KYC documents API error:', error);
            // Show empty documents on error instead of mock
            setDocModal({ isOpen: true, user, documents: [] });
        } finally {
            setLoadingDocs(false);
        }
    };

    const handleApprove = async () => {
        if (!actionModal.user) return;

        setActionLoading(actionModal.user.id);
        try {
            await approveKYC(actionModal.user.id, 'verified');
            await fetchUsers();
            setToast({ message: `KYC de ${actionModal.user.first_name} ${actionModal.user.last_name} approuvé`, type: 'success' });
            setActionModal({ isOpen: false, user: null, action: 'approve' });
        } catch (error) {
            console.error('Failed to approve KYC:', error);
            setToast({ message: "Erreur lors de l'approbation", type: 'error' });
        } finally {
            setActionLoading(null);
        }
    };

    const handleReject = async () => {
        if (!actionModal.user || !rejectReason.trim()) return;

        setActionLoading(actionModal.user.id);
        try {
            await rejectKYC(actionModal.user.id, rejectReason);
            await fetchUsers();
            setToast({ message: `KYC de ${actionModal.user.first_name} ${actionModal.user.last_name} rejeté`, type: 'success' });
            setActionModal({ isOpen: false, user: null, action: 'reject' });
            setRejectReason('');
        } catch (error) {
            console.error('Failed to reject KYC:', error);
            setToast({ message: 'Erreur lors du rejet', type: 'error' });
        } finally {
            setActionLoading(null);
        }
    };

    const getStatusBadge = (status: string) => {
        switch (status) {
            case 'pending':
                return 'bg-amber-100 text-amber-700';
            case 'verified':
                return 'bg-green-100 text-green-700';
            case 'rejected':
                return 'bg-red-100 text-red-700';
            case 'none':
            case '':
                return 'bg-gray-100 text-gray-600';
            default:
                return 'bg-gray-100 text-gray-600';
        }
    };

    const getStatusLabel = (status: string) => {
        switch (status) {
            case 'pending': return 'En attente';
            case 'verified': return 'Vérifié';
            case 'rejected': return 'Rejeté';
            case 'none':
            case '':
                return 'Non soumis';
            default: return 'Non soumis';
        }
    };

    const getDocTypeName = (type: string) => {
        switch (type) {
            case 'identity': return "🪪 Pièce d'identité";
            case 'selfie': return '🤳 Selfie';
            case 'address': return '🏠 Justificatif de domicile';
            default: return '📄 Document';
        }
    };

    const formatDate = (dateString: string) => {
        if (!dateString) return '-';
        try {
            return new Date(dateString).toLocaleDateString('fr-FR', {
                day: '2-digit',
                month: '2-digit',
                year: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            });
        } catch {
            return '-';
        }
    };

    // Stats
    const noneCount = allUsers.filter(u => !u.kyc_status || u.kyc_status === 'none' || u.kyc_status === '').length;
    const pendingCount = allUsers.filter(u => u.kyc_status === 'pending').length;
    const verifiedCount = allUsers.filter(u => u.kyc_status === 'verified').length;
    const rejectedCount = allUsers.filter(u => u.kyc_status === 'rejected').length;

    return (
        <div className="space-y-6">
            {/* Toast */}
            {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}

            {/* Header */}
            <div className="flex items-center justify-between flex-wrap gap-4">
                <div>
                    <h1 className="text-3xl font-bold bg-gradient-to-r from-slate-900 to-slate-600 bg-clip-text text-transparent">
                        Vérification KYC
                    </h1>
                    <p className="text-gray-500 mt-1">Gérez les demandes de vérification d'identité</p>
                </div>
                <button
                    onClick={fetchUsers}
                    className="flex items-center gap-2 px-4 py-2.5 bg-white border border-gray-200 rounded-xl hover:bg-gray-50 transition-colors shadow-sm"
                >
                    <ArrowPathIcon className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
                    Actualiser
                </button>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
                <div className="stat-card-primary flex items-center gap-3">
                    <div className="w-10 h-10 rounded-xl bg-white/20 flex items-center justify-center">
                        <UserIcon className="w-5 h-5" />
                    </div>
                    <div>
                        <p className="text-white/70 text-sm">Total</p>
                        <p className="text-2xl font-bold">{allUsers.length}</p>
                    </div>
                </div>
                <div className="bg-white rounded-xl p-4 border border-gray-200 shadow-sm flex items-center gap-3">
                    <div className="w-10 h-10 rounded-xl bg-gray-100 flex items-center justify-center">
                        <DocumentIcon className="w-5 h-5 text-gray-500" />
                    </div>
                    <div>
                        <p className="text-gray-500 text-sm">Non soumis</p>
                        <p className="text-2xl font-bold text-gray-600">{noneCount}</p>
                    </div>
                </div>
                <div className="stat-card-warning flex items-center gap-3">
                    <div className="w-10 h-10 rounded-xl bg-white/20 flex items-center justify-center">
                        <ClockIcon className="w-5 h-5" />
                    </div>
                    <div>
                        <p className="text-white/70 text-sm">En attente</p>
                        <p className="text-2xl font-bold">{pendingCount}</p>
                    </div>
                </div>
                <div className="stat-card-success flex items-center gap-3">
                    <div className="w-10 h-10 rounded-xl bg-white/20 flex items-center justify-center">
                        <ShieldCheckIcon className="w-5 h-5" />
                    </div>
                    <div>
                        <p className="text-white/70 text-sm">Vérifiés</p>
                        <p className="text-2xl font-bold">{verifiedCount}</p>
                    </div>
                </div>
                <div className="stat-card-danger flex items-center gap-3">
                    <div className="w-10 h-10 rounded-xl bg-white/20 flex items-center justify-center">
                        <XCircleIcon className="w-5 h-5" />
                    </div>
                    <div>
                        <p className="text-white/70 text-sm">Rejetés</p>
                        <p className="text-2xl font-bold">{rejectedCount}</p>
                    </div>
                </div>
            </div>

            {/* Filters */}
            <div className="flex gap-2 flex-wrap">
                {[
                    { key: 'pending', label: `🕐 En attente`, count: pendingCount },
                    { key: 'none', label: `📝 Non soumis`, count: noneCount },
                    { key: 'verified', label: `✅ Vérifiés`, count: verifiedCount },
                    { key: 'rejected', label: `❌ Rejetés`, count: rejectedCount },
                    { key: 'all', label: 'Tous', count: allUsers.length }
                ].map(f => (
                    <button
                        key={f.key}
                        onClick={() => setFilter(f.key)}
                        className={`px-4 py-2.5 rounded-xl text-sm font-medium transition-all ${filter === f.key
                            ? 'bg-gradient-to-r from-indigo-500 to-purple-600 text-white shadow-lg shadow-indigo-500/30'
                            : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50 shadow-sm'
                            }`}
                    >
                        {f.label}
                        <span className={`ml-2 px-2 py-0.5 rounded-full text-xs ${filter === f.key ? 'bg-white/20' : 'bg-gray-100'}`}>
                            {f.count}
                        </span>
                    </button>
                ))}
            </div>

            {/* Content */}
            {loading ? (
                <div className="flex items-center justify-center h-64">
                    <div className="text-center">
                        <div className="w-16 h-16 rounded-full border-4 border-indigo-100 border-t-indigo-500 animate-spin mx-auto"></div>
                        <p className="mt-4 text-gray-500 font-medium">Chargement...</p>
                    </div>
                </div>
            ) : users.length === 0 ? (
                <div className="bg-white rounded-2xl p-12 text-center border border-gray-100 shadow-sm">
                    <p className="text-4xl mb-4">📋</p>
                    <p className="text-gray-500 font-medium">Aucune demande {filter === 'pending' ? 'en attente' : ''}</p>
                </div>
            ) : (
                <div className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
                    <table className="w-full">
                        <thead className="bg-gray-50 border-b border-gray-100">
                            <tr>
                                <th className="text-left px-6 py-4 text-xs font-semibold text-gray-500 uppercase">Utilisateur</th>
                                <th className="text-left px-6 py-4 text-xs font-semibold text-gray-500 uppercase">Contact</th>
                                <th className="text-left px-6 py-4 text-xs font-semibold text-gray-500 uppercase">Niveau</th>
                                <th className="text-left px-6 py-4 text-xs font-semibold text-gray-500 uppercase">Statut</th>
                                <th className="text-left px-6 py-4 text-xs font-semibold text-gray-500 uppercase">Inscription</th>
                                <th className="text-right px-6 py-4 text-xs font-semibold text-gray-500 uppercase">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100">
                            {users.map((user) => (
                                <tr key={user.id} className="hover:bg-indigo-50/50 transition-colors group">
                                    <td className="px-6 py-4">
                                        <div className="flex items-center gap-3">
                                            <div className={`w-10 h-10 rounded-xl flex items-center justify-center font-semibold text-white bg-gradient-to-br ${user.kyc_status === 'pending' ? 'from-amber-500 to-orange-600' :
                                                user.kyc_status === 'verified' ? 'from-green-500 to-emerald-600' :
                                                    'from-gray-400 to-gray-500'
                                                }`}>
                                                {(user.first_name?.[0] || '')}{(user.last_name?.[0] || '')}
                                            </div>
                                            <div>
                                                <p className="font-semibold text-slate-900">
                                                    {user.first_name} {user.last_name}
                                                </p>
                                                <p className="text-xs text-gray-400 font-mono">{user.id?.slice(0, 8)}...</p>
                                            </div>
                                        </div>
                                    </td>
                                    <td className="px-6 py-4">
                                        <p className="text-sm text-slate-900">{user.email}</p>
                                        <p className="text-sm text-gray-400">{user.phone || '-'}</p>
                                    </td>
                                    <td className="px-6 py-4">
                                        <span className="inline-flex items-center px-2.5 py-1 rounded-lg text-xs font-medium bg-blue-100 text-blue-700">
                                            Niveau {user.kyc_level || 0}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4">
                                        <span className={`inline-flex items-center px-2.5 py-1 rounded-lg text-xs font-medium ${getStatusBadge(user.kyc_status)}`}>
                                            {getStatusLabel(user.kyc_status)}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4 text-sm text-gray-500">
                                        {formatDate(user.created_at)}
                                    </td>
                                    <td className="px-6 py-4 text-right">
                                        <div className="flex gap-2 justify-end">
                                            {user.kyc_status === 'pending' && (
                                                <>
                                                    <button
                                                        onClick={() => fetchUserDocuments(user)}
                                                        className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-indigo-700 bg-indigo-50 hover:bg-indigo-100 rounded-lg transition-colors"
                                                    >
                                                        <EyeIcon className="w-3.5 h-3.5" />
                                                        Voir docs
                                                    </button>
                                                    <button
                                                        onClick={() => setActionModal({ isOpen: true, user, action: 'approve' })}
                                                        className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-green-700 bg-green-50 hover:bg-green-100 rounded-lg transition-colors"
                                                    >
                                                        <CheckCircleIcon className="w-3.5 h-3.5" />
                                                        Approuver
                                                    </button>
                                                    <button
                                                        onClick={() => setActionModal({ isOpen: true, user, action: 'reject' })}
                                                        className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-red-700 bg-red-50 hover:bg-red-100 rounded-lg transition-colors"
                                                    >
                                                        <XCircleIcon className="w-3.5 h-3.5" />
                                                        Rejeter
                                                    </button>
                                                </>
                                            )}
                                            {user.kyc_status === 'verified' && (
                                                <span className="text-green-600 text-sm font-medium flex items-center gap-1">
                                                    <CheckCircleIcon className="w-4 h-4" /> Approuvé
                                                </span>
                                            )}
                                            {user.kyc_status === 'rejected' && (
                                                <button
                                                    onClick={() => setActionModal({ isOpen: true, user, action: 'approve' })}
                                                    className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-blue-700 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors"
                                                >
                                                    Réexaminer
                                                </button>
                                            )}
                                            {(!user.kyc_status || user.kyc_status === 'none' || user.kyc_status === '') && (
                                                <span className="text-gray-400 text-sm">Aucun document</span>
                                            )}
                                        </div>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}

            {/* Document Viewer Modal */}
            <Modal
                isOpen={docModal.isOpen}
                onClose={() => setDocModal({ isOpen: false, user: null, documents: [] })}
                title={`Documents KYC - ${docModal.user?.first_name} ${docModal.user?.last_name}`}
                size="lg"
            >
                {loadingDocs ? (
                    <div className="flex justify-center py-12">
                        <div className="w-10 h-10 border-4 border-indigo-100 border-t-indigo-500 rounded-full animate-spin"></div>
                    </div>
                ) : docModal.documents.length === 0 ? (
                    <div className="text-center py-12">
                        <DocumentIcon className="w-12 h-12 text-gray-300 mx-auto mb-4" />
                        <p className="text-gray-500">Aucun document soumis</p>
                    </div>
                ) : (
                    <div className="space-y-6">
                        {docModal.documents.map((doc) => (
                            <div key={doc.id} className="border border-gray-200 rounded-xl overflow-hidden">
                                <div className="flex items-center justify-between p-4 bg-gray-50 border-b border-gray-200">
                                    <div className="flex items-center gap-3">
                                        <span className="text-2xl">{doc.type === 'identity' ? '🪪' : doc.type === 'selfie' ? '🤳' : '🏠'}</span>
                                        <div>
                                            <p className="font-semibold text-gray-900">{getDocTypeName(doc.type)}</p>
                                            <p className="text-xs text-gray-500">Envoyé le {formatDate(doc.uploaded_at)}</p>
                                        </div>
                                    </div>
                                    <span className={`px-3 py-1 rounded-full text-xs font-medium ${doc.status === 'approved' ? 'bg-green-100 text-green-700' :
                                        doc.status === 'rejected' ? 'bg-red-100 text-red-700' :
                                            'bg-amber-100 text-amber-700'
                                        }`}>
                                        {doc.status === 'approved' ? 'Approuvé' : doc.status === 'rejected' ? 'Rejeté' : 'En attente'}
                                    </span>
                                </div>
                                {(doc.file_url || doc.file_path) && (
                                    <SecureDocumentViewer doc={doc} />
                                )}
                                {doc.rejection_reason && (
                                    <div className="p-4 bg-red-50 border-t border-red-100">
                                        <p className="text-sm text-red-700">
                                            <strong>Raison du rejet:</strong> {doc.rejection_reason}
                                        </p>
                                    </div>
                                )}
                            </div>
                        ))}

                        <div className="flex gap-3 pt-4 border-t border-gray-100">
                            <button
                                onClick={() => {
                                    setDocModal({ isOpen: false, user: null, documents: [] });
                                    if (docModal.user) {
                                        setActionModal({ isOpen: true, user: docModal.user, action: 'approve' });
                                    }
                                }}
                                className="flex-1 py-3 px-4 rounded-xl font-semibold text-white bg-gradient-to-r from-green-500 to-green-600 hover:from-green-600 hover:to-green-700 transition-all shadow-lg shadow-green-500/25"
                            >
                                ✓ Approuver le KYC
                            </button>
                            <button
                                onClick={() => {
                                    setDocModal({ isOpen: false, user: null, documents: [] });
                                    if (docModal.user) {
                                        setActionModal({ isOpen: true, user: docModal.user, action: 'reject' });
                                    }
                                }}
                                className="flex-1 py-3 px-4 rounded-xl font-semibold text-white bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 transition-all shadow-lg shadow-red-500/25"
                            >
                                ✗ Rejeter le KYC
                            </button>
                        </div>
                    </div>
                )}
            </Modal>

            {/* Action Modal */}
            <Modal
                isOpen={actionModal.isOpen}
                onClose={() => { setActionModal({ isOpen: false, user: null, action: 'approve' }); setRejectReason(''); }}
                title={actionModal.action === 'approve' ? 'Approuver le KYC' : 'Rejeter le KYC'}
            >
                <div className="space-y-4">
                    <div className="flex items-center gap-3 p-4 bg-gray-50 rounded-xl">
                        <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-semibold">
                            {actionModal.user?.first_name?.[0]}{actionModal.user?.last_name?.[0]}
                        </div>
                        <div>
                            <p className="font-semibold text-gray-900">{actionModal.user?.first_name} {actionModal.user?.last_name}</p>
                            <p className="text-sm text-gray-500">{actionModal.user?.email}</p>
                        </div>
                    </div>

                    {actionModal.action === 'approve' ? (
                        <div className="bg-green-50 border border-green-100 p-4 rounded-xl">
                            <div className="flex items-start gap-3">
                                <CheckCircleIcon className="w-5 h-5 text-green-600 flex-shrink-0 mt-0.5" />
                                <p className="text-sm text-green-800">
                                    L'utilisateur aura accès à toutes les fonctionnalités après approbation.
                                </p>
                            </div>
                        </div>
                    ) : (
                        <>
                            <div className="bg-red-50 border border-red-100 p-4 rounded-xl">
                                <div className="flex items-start gap-3">
                                    <ExclamationTriangleIcon className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" />
                                    <p className="text-sm text-red-800">
                                        L'utilisateur devra soumettre de nouveaux documents après le rejet.
                                    </p>
                                </div>
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-2">
                                    Raison du rejet <span className="text-red-500">*</span>
                                </label>
                                <textarea
                                    value={rejectReason}
                                    onChange={(e) => setRejectReason(e.target.value)}
                                    placeholder="Documents flous, informations illisibles, etc."
                                    rows={3}
                                    className="w-full px-4 py-3 rounded-xl border border-gray-200 focus:border-red-500 focus:ring-2 focus:ring-red-500/20 outline-none transition-all resize-none"
                                />
                            </div>
                        </>
                    )}

                    <div className="flex gap-3 pt-2">
                        <button
                            onClick={() => { setActionModal({ isOpen: false, user: null, action: 'approve' }); setRejectReason(''); }}
                            className="flex-1 py-3 px-4 rounded-xl font-semibold text-gray-700 bg-gray-100 hover:bg-gray-200 transition-colors"
                        >
                            Annuler
                        </button>
                        <button
                            onClick={actionModal.action === 'approve' ? handleApprove : handleReject}
                            disabled={actionLoading === actionModal.user?.id || (actionModal.action === 'reject' && !rejectReason.trim())}
                            className={`flex-1 py-3 px-4 rounded-xl font-semibold text-white transition-all shadow-lg disabled:opacity-50 disabled:cursor-not-allowed ${actionModal.action === 'approve'
                                ? 'bg-gradient-to-r from-green-500 to-green-600 hover:from-green-600 hover:to-green-700 shadow-green-500/25'
                                : 'bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 shadow-red-500/25'
                                }`}
                        >
                            {actionLoading === actionModal.user?.id ? (
                                <span className="flex items-center justify-center gap-2">
                                    <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                                    Traitement...
                                </span>
                            ) : actionModal.action === 'approve' ? 'Confirmer l\'approbation' : 'Confirmer le rejet'}
                        </button>
                    </div>
                </div>
            </Modal>
        </div>
    );
}
