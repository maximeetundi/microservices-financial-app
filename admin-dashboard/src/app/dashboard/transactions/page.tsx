'use client';

import { useEffect, useState } from 'react';
import { getTransactions, blockTransaction, refundTransaction } from '@/lib/api';
import { format, isToday, isYesterday, parseISO } from 'date-fns';
import { fr } from 'date-fns/locale';
import {
    ArrowsRightLeftIcon,
    MagnifyingGlassIcon,
    XMarkIcon,
    ExclamationTriangleIcon,
    CheckCircleIcon,
    ArrowPathIcon,
    NoSymbolIcon,
    ArrowUturnLeftIcon,
    BanknotesIcon,
    ClockIcon,
    CurrencyDollarIcon,
    FunnelIcon,
} from '@heroicons/react/24/outline';

interface Transaction {
    id: string;
    wallet_id: string;
    type: string;
    amount: number;
    currency: string;
    status: string;
    created_at: string;
    description?: string;
    user_email?: string;
}

// Toast Component
interface ToastProps {
    message: string;
    type: 'success' | 'error' | 'info';
    onClose: () => void;
}

function Toast({ message, type, onClose }: ToastProps) {
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
interface ModalProps {
    isOpen: boolean;
    onClose: () => void;
    title: string;
    icon?: React.ReactNode;
    iconColor?: string;
    children: React.ReactNode;
}

function Modal({ isOpen, onClose, title, icon, iconColor = 'bg-indigo-100 text-indigo-600', children }: ModalProps) {
    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div
                className="absolute inset-0 bg-black/50 backdrop-blur-sm animate-fade-in"
                onClick={onClose}
            />
            <div className="relative bg-white rounded-2xl shadow-2xl w-full max-w-md animate-slide-up overflow-hidden">
                <div className="flex items-center gap-4 p-6 border-b border-gray-100">
                    {icon && (
                        <div className={`w-12 h-12 rounded-xl flex items-center justify-center ${iconColor}`}>
                            {icon}
                        </div>
                    )}
                    <h3 className="text-xl font-bold text-gray-900">{title}</h3>
                    <button
                        onClick={onClose}
                        className="ml-auto p-2 rounded-lg hover:bg-gray-100 transition-colors"
                    >
                        <XMarkIcon className="w-5 h-5 text-gray-500" />
                    </button>
                </div>
                <div className="p-6">
                    {children}
                </div>
            </div>
        </div>
    );
}

type FilterTab = 'all' | 'completed' | 'pending' | 'blocked' | 'failed';
type TypeFilter = 'all' | 'deposit' | 'withdrawal' | 'transfer' | 'exchange';

export default function TransactionsPage() {
    const [transactions, setTransactions] = useState<Transaction[]>([]);
    const [loading, setLoading] = useState(true);
    const [actionLoading, setActionLoading] = useState<string | null>(null);
    const [searchQuery, setSearchQuery] = useState('');
    const [statusFilter, setStatusFilter] = useState<FilterTab>('all');
    const [typeFilter, setTypeFilter] = useState<TypeFilter>('all');

    // Modal states
    const [blockModal, setBlockModal] = useState<{ isOpen: boolean; tx: Transaction | null }>({ isOpen: false, tx: null });
    const [refundModal, setRefundModal] = useState<{ isOpen: boolean; tx: Transaction | null }>({ isOpen: false, tx: null });
    const [reason, setReason] = useState('');

    // Toast state
    const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' | 'info' } | null>(null);

    const showToast = (message: string, type: 'success' | 'error' | 'info') => {
        setToast({ message, type });
    };

    const fetchTransactions = async () => {
        setLoading(true);
        try {
            const response = await getTransactions();
            setTransactions(response.data.transactions || []);
        } catch (error) {
            console.error('Failed to fetch transactions:', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchTransactions();
    }, []);

    const handleBlock = async () => {
        if (!blockModal.tx || !reason.trim()) return;

        setActionLoading(blockModal.tx.id);
        try {
            await blockTransaction(blockModal.tx.id, reason);
            await fetchTransactions();
            setBlockModal({ isOpen: false, tx: null });
            setReason('');
            showToast('Transaction bloquée avec succès', 'success');
        } catch (error) {
            showToast('Erreur lors du blocage de la transaction', 'error');
        } finally {
            setActionLoading(null);
        }
    };

    const handleRefund = async () => {
        if (!refundModal.tx || !reason.trim()) return;

        setActionLoading(refundModal.tx.id);
        try {
            await refundTransaction(refundModal.tx.id, reason);
            await fetchTransactions();
            setRefundModal({ isOpen: false, tx: null });
            setReason('');
            showToast('Remboursement initié avec succès', 'success');
        } catch (error) {
            showToast('Erreur lors du remboursement', 'error');
        } finally {
            setActionLoading(null);
        }
    };

    // Filter transactions
    const filteredTransactions = transactions.filter(tx => {
        const matchesSearch =
            tx.id?.toLowerCase().includes(searchQuery.toLowerCase()) ||
            tx.wallet_id?.toLowerCase().includes(searchQuery.toLowerCase()) ||
            tx.user_email?.toLowerCase().includes(searchQuery.toLowerCase());

        const matchesStatus = statusFilter === 'all' ||
            tx.status?.toLowerCase() === statusFilter ||
            (statusFilter === 'completed' && tx.status?.toLowerCase() === 'confirmed');

        const matchesType = typeFilter === 'all' ||
            tx.type?.toLowerCase() === typeFilter ||
            (typeFilter === 'transfer' && (tx.type?.toLowerCase() === 'send' || tx.type?.toLowerCase() === 'receive'));

        return matchesSearch && matchesStatus && matchesType;
    });

    // Compute stats
    const totalAmount = transactions.reduce((sum, tx) => sum + (tx.amount || 0), 0);
    const completedCount = transactions.filter(tx => ['completed', 'confirmed'].includes(tx.status?.toLowerCase())).length;
    const pendingCount = transactions.filter(tx => tx.status?.toLowerCase() === 'pending').length;
    const blockedCount = transactions.filter(tx => tx.status?.toLowerCase() === 'blocked').length;

    const statusTabs: { key: FilterTab; label: string; count: number }[] = [
        { key: 'all', label: 'Toutes', count: transactions.length },
        { key: 'completed', label: 'Complétées', count: completedCount },
        { key: 'pending', label: 'En attente', count: pendingCount },
        { key: 'blocked', label: 'Bloquées', count: blockedCount },
    ];

    const typeTabs: { key: TypeFilter; label: string }[] = [
        { key: 'all', label: 'Tous types' },
        { key: 'deposit', label: 'Dépôts' },
        { key: 'withdrawal', label: 'Retraits' },
        { key: 'transfer', label: 'Transferts' },
        { key: 'exchange', label: 'Échanges' },
    ];

    if (loading) {
        return (
            <div className="flex items-center justify-center h-64">
                <div className="text-center">
                    <div className="w-16 h-16 rounded-full border-4 border-indigo-100 border-t-indigo-500 animate-spin mx-auto"></div>
                    <p className="mt-4 text-gray-500 font-medium">Chargement des transactions...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="space-y-6">
            {/* Toast */}
            {toast && (
                <Toast
                    message={toast.message}
                    type={toast.type}
                    onClose={() => setToast(null)}
                />
            )}

            {/* Header */}
            <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
                <div>
                    <h1 className="text-3xl font-bold bg-gradient-to-r from-slate-900 to-slate-600 bg-clip-text text-transparent">
                        Gestion des Transactions
                    </h1>
                    <p className="text-slate-500 mt-1">
                        Surveillez et gérez toutes les transactions du système
                    </p>
                </div>
                <button
                    onClick={fetchTransactions}
                    className="btn-secondary flex items-center gap-2"
                >
                    <ArrowPathIcon className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
                    Actualiser
                </button>
            </div>

            {/* Stats Cards */}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                <div className="stat-card-primary flex items-center gap-4">
                    <div className="w-12 h-12 rounded-xl bg-white/20 flex items-center justify-center">
                        <ArrowsRightLeftIcon className="w-6 h-6" />
                    </div>
                    <div>
                        <p className="text-white/70 text-sm">Total</p>
                        <p className="text-2xl font-bold">{transactions.length}</p>
                    </div>
                </div>
                <div className="stat-card-success flex items-center gap-4">
                    <div className="w-12 h-12 rounded-xl bg-white/20 flex items-center justify-center">
                        <CheckCircleIcon className="w-6 h-6" />
                    </div>
                    <div>
                        <p className="text-white/70 text-sm">Complétées</p>
                        <p className="text-2xl font-bold">{completedCount}</p>
                    </div>
                </div>
                <div className="stat-card-warning flex items-center gap-4">
                    <div className="w-12 h-12 rounded-xl bg-white/20 flex items-center justify-center">
                        <ClockIcon className="w-6 h-6" />
                    </div>
                    <div>
                        <p className="text-white/70 text-sm">En attente</p>
                        <p className="text-2xl font-bold">{pendingCount}</p>
                    </div>
                </div>
                <div className="stat-card-info flex items-center gap-4">
                    <div className="w-12 h-12 rounded-xl bg-white/20 flex items-center justify-center">
                        <CurrencyDollarIcon className="w-6 h-6" />
                    </div>
                    <div>
                        <p className="text-white/70 text-sm">Volume total</p>
                        <p className="text-2xl font-bold">{totalAmount.toLocaleString()} €</p>
                    </div>
                </div>
            </div>

            {/* Filters */}
            <div className="flex flex-col gap-4">
                {/* Status Tabs */}
                <div className="flex gap-2 overflow-x-auto pb-2">
                    {statusTabs.map(tab => (
                        <button
                            key={tab.key}
                            onClick={() => setStatusFilter(tab.key)}
                            className={`px-4 py-2.5 rounded-xl text-sm font-medium whitespace-nowrap transition-all ${statusFilter === tab.key
                                    ? 'bg-gradient-to-r from-indigo-500 to-purple-600 text-white shadow-lg shadow-indigo-500/30'
                                    : 'bg-white text-gray-600 hover:bg-gray-50 border border-gray-200 shadow-sm'
                                }`}
                        >
                            {tab.label}
                            <span className={`ml-2 px-2 py-0.5 rounded-full text-xs ${statusFilter === tab.key ? 'bg-white/20' : 'bg-gray-100'
                                }`}>
                                {tab.count}
                            </span>
                        </button>
                    ))}
                </div>

                {/* Type Filter & Search */}
                <div className="flex flex-col lg:flex-row gap-4 items-start lg:items-center justify-between">
                    {/* Type Filter */}
                    <div className="flex items-center gap-2">
                        <FunnelIcon className="w-4 h-4 text-gray-400" />
                        <select
                            value={typeFilter}
                            onChange={(e) => setTypeFilter(e.target.value as TypeFilter)}
                            className="px-4 py-2 bg-white border border-gray-200 rounded-xl text-sm font-medium text-gray-600 focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                        >
                            {typeTabs.map(tab => (
                                <option key={tab.key} value={tab.key}>{tab.label}</option>
                            ))}
                        </select>
                    </div>

                    {/* Search */}
                    <div className="relative w-full lg:w-80">
                        <MagnifyingGlassIcon className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
                        <input
                            type="text"
                            placeholder="Rechercher par ID, wallet..."
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            className="input pl-12"
                        />
                    </div>
                </div>
            </div>

            {/* Transactions Table */}
            <div className="table-container">
                <table className="w-full min-w-[800px]">
                    <thead className="bg-gray-50 border-b">
                        <tr>
                            <th className="text-left px-6 py-4 text-sm font-semibold text-slate-600">ID</th>
                            <th className="text-left px-6 py-4 text-sm font-semibold text-slate-600">Type</th>
                            <th className="text-right px-6 py-4 text-sm font-semibold text-slate-600">Montant</th>
                            <th className="text-left px-6 py-4 text-sm font-semibold text-slate-600">Status</th>
                            <th className="text-left px-6 py-4 text-sm font-semibold text-slate-600">Date</th>
                            <th className="text-right px-6 py-4 text-sm font-semibold text-slate-600">Actions</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-100">
                        {filteredTransactions.map((tx) => (
                            <tr key={tx.id} className="hover:bg-indigo-50/50 transition-colors group">
                                <td className="px-6 py-4">
                                    <div className="flex items-center gap-3">
                                        <div className={`w-10 h-10 rounded-xl flex items-center justify-center ${getTypeBgColor(tx.type)}`}>
                                            {getTypeIcon(tx.type)}
                                        </div>
                                        <p className="font-mono text-sm text-slate-700 truncate max-w-[120px]">{tx.id}</p>
                                    </div>
                                </td>
                                <td className="px-6 py-4">
                                    <span className={`badge ${getTypeBadge(tx.type)}`}>
                                        {getTypeLabel(tx.type)}
                                    </span>
                                </td>
                                <td className="px-6 py-4 text-right">
                                    <span className={`font-semibold ${getAmountColor(tx.type)}`}>
                                        {getAmountPrefix(tx.type)}{tx.amount?.toLocaleString()} {tx.currency}
                                    </span>
                                </td>
                                <td className="px-6 py-4">
                                    <span className={`badge ${getStatusBadge(tx.status)}`}>
                                        {getStatusLabel(tx.status)}
                                    </span>
                                </td>
                                <td className="px-6 py-4 text-slate-500 text-sm">
                                    {tx.created_at ? format(new Date(tx.created_at), 'dd/MM/yyyy HH:mm') : '-'}
                                </td>
                                <td className="px-6 py-4">
                                    <div className="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                                        <button
                                            onClick={() => { setBlockModal({ isOpen: true, tx }); setReason(''); }}
                                            className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-red-700 bg-red-50 hover:bg-red-100 rounded-lg transition-colors"
                                            title="Bloquer"
                                        >
                                            <NoSymbolIcon className="w-3.5 h-3.5" />
                                            Bloquer
                                        </button>
                                        <button
                                            onClick={() => { setRefundModal({ isOpen: true, tx }); setReason(''); }}
                                            className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-amber-700 bg-amber-50 hover:bg-amber-100 rounded-lg transition-colors"
                                            title="Rembourser"
                                        >
                                            <ArrowUturnLeftIcon className="w-3.5 h-3.5" />
                                            Rembourser
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
                {filteredTransactions.length === 0 && (
                    <div className="text-center py-16">
                        <ArrowsRightLeftIcon className="w-12 h-12 text-slate-300 mx-auto mb-4" />
                        <p className="text-slate-500 font-medium">Aucune transaction trouvée</p>
                        <p className="text-slate-400 text-sm mt-1">
                            {searchQuery ? 'Essayez une autre recherche' : 'Les transactions apparaîtront ici'}
                        </p>
                    </div>
                )}
            </div>

            {/* Block Transaction Modal */}
            <Modal
                isOpen={blockModal.isOpen}
                onClose={() => { setBlockModal({ isOpen: false, tx: null }); setReason(''); }}
                title="Bloquer la transaction"
                icon={<NoSymbolIcon className="w-6 h-6" />}
                iconColor="bg-red-100 text-red-600"
            >
                <div className="space-y-4">
                    {blockModal.tx && (
                        <div className="p-4 bg-gray-50 rounded-xl border border-gray-100">
                            <div className="flex items-center justify-between mb-2">
                                <span className="text-sm text-gray-500">ID Transaction</span>
                                <span className="font-mono text-sm">{blockModal.tx.id}</span>
                            </div>
                            <div className="flex items-center justify-between mb-2">
                                <span className="text-sm text-gray-500">Type</span>
                                <span className={`badge ${getTypeBadge(blockModal.tx.type)}`}>{getTypeLabel(blockModal.tx.type)}</span>
                            </div>
                            <div className="flex items-center justify-between">
                                <span className="text-sm text-gray-500">Montant</span>
                                <span className="font-semibold">{blockModal.tx.amount?.toLocaleString()} {blockModal.tx.currency}</span>
                            </div>
                        </div>
                    )}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Raison du blocage <span className="text-red-500">*</span>
                        </label>
                        <textarea
                            value={reason}
                            onChange={(e) => setReason(e.target.value)}
                            placeholder="Décrivez la raison du blocage..."
                            rows={3}
                            className="w-full px-4 py-3 rounded-xl border border-gray-200 focus:border-red-500 focus:ring-2 focus:ring-red-500/20 outline-none transition-all resize-none"
                        />
                    </div>
                    <div className="flex gap-3 pt-2">
                        <button
                            onClick={() => { setBlockModal({ isOpen: false, tx: null }); setReason(''); }}
                            className="flex-1 py-3 px-4 rounded-xl font-semibold text-gray-700 bg-gray-100 hover:bg-gray-200 transition-colors"
                        >
                            Annuler
                        </button>
                        <button
                            onClick={handleBlock}
                            disabled={!reason.trim() || actionLoading === blockModal.tx?.id}
                            className="flex-1 py-3 px-4 rounded-xl font-semibold text-white bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-lg shadow-red-500/25"
                        >
                            {actionLoading === blockModal.tx?.id ? (
                                <span className="flex items-center justify-center gap-2">
                                    <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                                    Blocage...
                                </span>
                            ) : 'Bloquer la transaction'}
                        </button>
                    </div>
                </div>
            </Modal>

            {/* Refund Transaction Modal */}
            <Modal
                isOpen={refundModal.isOpen}
                onClose={() => { setRefundModal({ isOpen: false, tx: null }); setReason(''); }}
                title="Rembourser la transaction"
                icon={<ArrowUturnLeftIcon className="w-6 h-6" />}
                iconColor="bg-amber-100 text-amber-600"
            >
                <div className="space-y-4">
                    {refundModal.tx && (
                        <div className="p-4 bg-gray-50 rounded-xl border border-gray-100">
                            <div className="flex items-center justify-between mb-2">
                                <span className="text-sm text-gray-500">ID Transaction</span>
                                <span className="font-mono text-sm">{refundModal.tx.id}</span>
                            </div>
                            <div className="flex items-center justify-between mb-2">
                                <span className="text-sm text-gray-500">Type</span>
                                <span className={`badge ${getTypeBadge(refundModal.tx.type)}`}>{getTypeLabel(refundModal.tx.type)}</span>
                            </div>
                            <div className="flex items-center justify-between">
                                <span className="text-sm text-gray-500">Montant à rembourser</span>
                                <span className="font-semibold text-green-600">+{refundModal.tx.amount?.toLocaleString()} {refundModal.tx.currency}</span>
                            </div>
                        </div>
                    )}
                    <div className="bg-amber-50 border border-amber-100 p-4 rounded-xl">
                        <div className="flex items-start gap-3">
                            <ExclamationTriangleIcon className="w-5 h-5 text-amber-600 flex-shrink-0 mt-0.5" />
                            <p className="text-sm text-amber-800">
                                Cette action créditera le montant sur le portefeuille de l'utilisateur. Cette opération est irréversible.
                            </p>
                        </div>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Raison du remboursement <span className="text-red-500">*</span>
                        </label>
                        <textarea
                            value={reason}
                            onChange={(e) => setReason(e.target.value)}
                            placeholder="Décrivez la raison du remboursement..."
                            rows={3}
                            className="w-full px-4 py-3 rounded-xl border border-gray-200 focus:border-amber-500 focus:ring-2 focus:ring-amber-500/20 outline-none transition-all resize-none"
                        />
                    </div>
                    <div className="flex gap-3 pt-2">
                        <button
                            onClick={() => { setRefundModal({ isOpen: false, tx: null }); setReason(''); }}
                            className="flex-1 py-3 px-4 rounded-xl font-semibold text-gray-700 bg-gray-100 hover:bg-gray-200 transition-colors"
                        >
                            Annuler
                        </button>
                        <button
                            onClick={handleRefund}
                            disabled={!reason.trim() || actionLoading === refundModal.tx?.id}
                            className="flex-1 py-3 px-4 rounded-xl font-semibold text-white bg-gradient-to-r from-amber-500 to-orange-500 hover:from-amber-600 hover:to-orange-600 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-lg shadow-amber-500/25"
                        >
                            {actionLoading === refundModal.tx?.id ? (
                                <span className="flex items-center justify-center gap-2">
                                    <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                                    Remboursement...
                                </span>
                            ) : 'Confirmer le remboursement'}
                        </button>
                    </div>
                </div>
            </Modal>
        </div>
    );
}

// Helper functions
function getTypeBadge(type: string): string {
    switch (type?.toLowerCase()) {
        case 'deposit':
        case 'receive':
            return 'badge-success';
        case 'withdrawal':
        case 'send':
            return 'badge-danger';
        case 'exchange':
            return 'badge-info';
        case 'transfer':
            return 'badge-purple';
        default:
            return 'badge-warning';
    }
}

function getTypeLabel(type: string): string {
    switch (type?.toLowerCase()) {
        case 'deposit':
            return 'Dépôt';
        case 'receive':
            return 'Réception';
        case 'withdrawal':
            return 'Retrait';
        case 'send':
            return 'Envoi';
        case 'exchange':
            return 'Échange';
        case 'transfer':
            return 'Transfert';
        default:
            return type || 'Inconnu';
    }
}

function getTypeBgColor(type: string): string {
    switch (type?.toLowerCase()) {
        case 'deposit':
        case 'receive':
            return 'bg-green-100';
        case 'withdrawal':
        case 'send':
            return 'bg-red-100';
        case 'exchange':
            return 'bg-blue-100';
        default:
            return 'bg-gray-100';
    }
}

function getTypeIcon(type: string) {
    const iconClass = "w-5 h-5";
    switch (type?.toLowerCase()) {
        case 'deposit':
        case 'receive':
            return <BanknotesIcon className={`${iconClass} text-green-600`} />;
        case 'withdrawal':
        case 'send':
            return <ArrowsRightLeftIcon className={`${iconClass} text-red-600`} />;
        case 'exchange':
            return <ArrowsRightLeftIcon className={`${iconClass} text-blue-600`} />;
        default:
            return <ArrowsRightLeftIcon className={`${iconClass} text-gray-600`} />;
    }
}

function getAmountColor(type: string): string {
    switch (type?.toLowerCase()) {
        case 'deposit':
        case 'receive':
            return 'text-green-600';
        case 'withdrawal':
        case 'send':
            return 'text-red-600';
        default:
            return 'text-gray-900';
    }
}

function getAmountPrefix(type: string): string {
    switch (type?.toLowerCase()) {
        case 'deposit':
        case 'receive':
            return '+';
        case 'withdrawal':
        case 'send':
            return '-';
        default:
            return '';
    }
}

function getStatusBadge(status: string): string {
    switch (status?.toLowerCase()) {
        case 'completed':
        case 'confirmed':
            return 'badge-success';
        case 'pending':
            return 'badge-warning';
        case 'failed':
        case 'blocked':
            return 'badge-danger';
        default:
            return 'badge-info';
    }
}

function getStatusLabel(status: string): string {
    switch (status?.toLowerCase()) {
        case 'completed':
        case 'confirmed':
            return '✓ Complétée';
        case 'pending':
            return '⏳ En attente';
        case 'failed':
            return '✗ Échouée';
        case 'blocked':
            return '🚫 Bloquée';
        default:
            return status || 'Inconnu';
    }
}
