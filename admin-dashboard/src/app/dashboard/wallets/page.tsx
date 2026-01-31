'use client';

import { useEffect, useState, useMemo } from 'react';
import { getWallets, freezeWallet } from '@/lib/api';
import { format } from 'date-fns';
import {
    WalletIcon,
    MagnifyingGlassIcon,
    FunnelIcon,
    ArrowPathIcon,
    LockClosedIcon,
    LockOpenIcon,
    EyeIcon,
    CurrencyDollarIcon,
    ChevronLeftIcon,
    ChevronRightIcon,
    ExclamationCircleIcon,
} from '@heroicons/react/24/outline';
import { CheckCircleIcon, XCircleIcon } from '@heroicons/react/24/solid';
import PageHeader from '@/components/ui/PageHeader';
import StatsCard from '@/components/ui/StatsCard';
import Modal, { ModalFooter, ModalButton } from '@/components/ui/Modal';
import InputDialog from '@/components/ui/InputDialog';
import { SimpleToast } from '@/components/ui/Toast';

interface Wallet {
    id: string;
    user_id: string;
    user_email?: string;
    user_name?: string;
    currency: string;
    balance: number;
    status: string;
    address?: string;
    network?: string;
    created_at: string;
    updated_at?: string;
}

type FilterCurrency = 'all' | 'crypto' | 'fiat';
type FilterStatus = 'all' | 'active' | 'frozen' | 'blocked';

const ITEMS_PER_PAGE = 10;

const CRYPTO_CURRENCIES = ['BTC', 'ETH', 'SOL', 'XRP', 'USDT', 'USDC', 'TRX', 'BNB', 'LTC', 'TON'];
const FIAT_CURRENCIES = ['EUR', 'USD', 'XOF', 'GBP', 'CHF'];

export default function WalletsPage() {
    // Data state
    const [wallets, setWallets] = useState<Wallet[]>([]);
    const [loading, setLoading] = useState(true);
    const [refreshing, setRefreshing] = useState(false);

    // Filters
    const [search, setSearch] = useState('');
    const [currencyFilter, setCurrencyFilter] = useState<FilterCurrency>('all');
    const [statusFilter, setStatusFilter] = useState<FilterStatus>('all');
    const [specificCurrency, setSpecificCurrency] = useState<string>('');

    // Pagination
    const [currentPage, setCurrentPage] = useState(1);

    // Modal state
    const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(null);
    const [detailsOpen, setDetailsOpen] = useState(false);
    const [freezeDialogOpen, setFreezeDialogOpen] = useState(false);
    const [actionLoading, setActionLoading] = useState(false);

    // Toast state
    const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' } | null>(null);

    const fetchWallets = async (silent = false) => {
        if (!silent) setLoading(true);
        else setRefreshing(true);

        try {
            const response = await getWallets();
            setWallets(response.data?.wallets || []);
        } catch (error) {
            console.error('Failed to fetch wallets:', error);
            setToast({ message: 'Erreur lors du chargement des portefeuilles', type: 'error' });
        } finally {
            setLoading(false);
            setRefreshing(false);
        }
    };

    useEffect(() => {
        fetchWallets();
    }, []);

    // Computed stats
    const stats = useMemo(() => {
        const active = wallets.filter(w => w.status === 'active').length;
        const frozen = wallets.filter(w => w.status === 'frozen').length;
        const crypto = wallets.filter(w => CRYPTO_CURRENCIES.includes(w.currency)).length;
        const fiat = wallets.filter(w => FIAT_CURRENCIES.includes(w.currency)).length;

        // Calculate total balances per major currency
        const btcTotal = wallets
            .filter(w => w.currency === 'BTC')
            .reduce((sum, w) => sum + (w.balance || 0), 0);
        const ethTotal = wallets
            .filter(w => w.currency === 'ETH')
            .reduce((sum, w) => sum + (w.balance || 0), 0);

        return { total: wallets.length, active, frozen, crypto, fiat, btcTotal, ethTotal };
    }, [wallets]);

    // Filtered wallets
    const filteredWallets = useMemo(() => {
        return wallets.filter(wallet => {
            // Search filter
            if (search) {
                const s = search.toLowerCase();
                const matchesSearch =
                    wallet.id.toLowerCase().includes(s) ||
                    wallet.user_id.toLowerCase().includes(s) ||
                    wallet.user_email?.toLowerCase().includes(s) ||
                    wallet.currency.toLowerCase().includes(s) ||
                    wallet.address?.toLowerCase().includes(s);
                if (!matchesSearch) return false;
            }

            // Currency type filter
            if (currencyFilter === 'crypto' && !CRYPTO_CURRENCIES.includes(wallet.currency)) return false;
            if (currencyFilter === 'fiat' && !FIAT_CURRENCIES.includes(wallet.currency)) return false;

            // Specific currency filter
            if (specificCurrency && wallet.currency !== specificCurrency) return false;

            // Status filter
            if (statusFilter !== 'all' && wallet.status !== statusFilter) return false;

            return true;
        });
    }, [wallets, search, currencyFilter, statusFilter, specificCurrency]);

    // Pagination
    const totalPages = Math.ceil(filteredWallets.length / ITEMS_PER_PAGE);
    const paginatedWallets = filteredWallets.slice(
        (currentPage - 1) * ITEMS_PER_PAGE,
        currentPage * ITEMS_PER_PAGE
    );

    // Reset page when filters change
    useEffect(() => {
        setCurrentPage(1);
    }, [search, currencyFilter, statusFilter, specificCurrency]);

    // Handlers
    const handleViewDetails = (wallet: Wallet) => {
        setSelectedWallet(wallet);
        setDetailsOpen(true);
    };

    const handleFreezeClick = (wallet: Wallet) => {
        setSelectedWallet(wallet);
        setFreezeDialogOpen(true);
    };

    const handleFreezeSubmit = async (reason: string) => {
        if (!selectedWallet) return;

        setActionLoading(true);
        try {
            await freezeWallet(selectedWallet.id, reason);
            setToast({ message: 'Portefeuille gelé avec succès', type: 'success' });
            setFreezeDialogOpen(false);
            await fetchWallets(true);
        } catch (error) {
            setToast({ message: 'Erreur lors du gel du portefeuille', type: 'error' });
        } finally {
            setActionLoading(false);
        }
    };

    const getCurrencyIcon = (currency: string) => {
        const icons: Record<string, string> = {
            BTC: '₿', ETH: 'Ξ', SOL: '◎', XRP: '✕', USDT: '₮', USDC: '$',
            EUR: '€', USD: '$', XOF: 'CFA', GBP: '£', CHF: 'Fr',
        };
        return icons[currency] || currency;
    };

    const getStatusBadge = (status: string) => {
        const config: Record<string, { bg: string; text: string; icon: React.ReactNode }> = {
            active: {
                bg: 'bg-green-100',
                text: 'text-green-700',
                icon: <CheckCircleIcon className="w-4 h-4" />,
            },
            frozen: {
                bg: 'bg-amber-100',
                text: 'text-amber-700',
                icon: <LockClosedIcon className="w-4 h-4" />,
            },
            blocked: {
                bg: 'bg-red-100',
                text: 'text-red-700',
                icon: <XCircleIcon className="w-4 h-4" />,
            },
        };
        const { bg, text, icon } = config[status] || config.active;
        return (
            <span className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold ${bg} ${text}`}>
                {icon}
                {status === 'active' ? 'Actif' : status === 'frozen' ? 'Gelé' : 'Bloqué'}
            </span>
        );
    };

    if (loading) {
        return (
            <div className="flex items-center justify-center h-96">
                <div className="text-center">
                    <div className="w-16 h-16 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin mx-auto mb-4" />
                    <p className="text-gray-500">Chargement des portefeuilles...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="pb-10">
            <PageHeader
                title="Portefeuilles"
                subtitle={`${wallets.length} portefeuilles au total`}
                icon="💰"
                onRefresh={() => fetchWallets(true)}
                loading={refreshing}
            />

            {/* Stats Cards */}
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
                <StatsCard
                    title="Total"
                    value={stats.total}
                    icon={WalletIcon}
                    color="primary"
                />
                <StatsCard
                    title="Actifs"
                    value={stats.active}
                    icon={LockOpenIcon}
                    color="success"
                />
                <StatsCard
                    title="Gelés"
                    value={stats.frozen}
                    icon={LockClosedIcon}
                    color="warning"
                />
                <StatsCard
                    title="Crypto"
                    value={stats.crypto}
                    icon={CurrencyDollarIcon}
                    color="purple"
                />
            </div>

            {/* Filters */}
            <div className="bg-white rounded-2xl border border-gray-200 shadow-sm p-4 mb-6">
                <div className="flex flex-col lg:flex-row gap-4">
                    {/* Search */}
                    <div className="flex-1 relative">
                        <MagnifyingGlassIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                        <input
                            type="text"
                            value={search}
                            onChange={(e) => setSearch(e.target.value)}
                            placeholder="Rechercher par ID, email, adresse..."
                            className="w-full pl-10 pr-4 py-2.5 rounded-xl border border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 outline-none transition-all"
                        />
                    </div>

                    {/* Currency Type Filter */}
                    <div className="flex gap-2">
                        {[
                            { id: 'all', label: 'Tous' },
                            { id: 'crypto', label: '🔗 Crypto' },
                            { id: 'fiat', label: '💵 Fiat' },
                        ].map((option) => (
                            <button
                                key={option.id}
                                onClick={() => {
                                    setCurrencyFilter(option.id as FilterCurrency);
                                    setSpecificCurrency('');
                                }}
                                className={`px-4 py-2 rounded-xl text-sm font-medium transition-all ${currencyFilter === option.id
                                        ? 'bg-indigo-600 text-white shadow-md'
                                        : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                                    }`}
                            >
                                {option.label}
                            </button>
                        ))}
                    </div>

                    {/* Specific Currency */}
                    <select
                        value={specificCurrency}
                        onChange={(e) => setSpecificCurrency(e.target.value)}
                        className="px-4 py-2.5 rounded-xl border border-gray-200 bg-white text-gray-700 focus:border-indigo-500 outline-none"
                    >
                        <option value="">Toutes devises</option>
                        <optgroup label="Crypto-monnaies">
                            {CRYPTO_CURRENCIES.map(c => <option key={c} value={c}>{c}</option>)}
                        </optgroup>
                        <optgroup label="Devises Fiat">
                            {FIAT_CURRENCIES.map(c => <option key={c} value={c}>{c}</option>)}
                        </optgroup>
                    </select>

                    {/* Status Filter */}
                    <select
                        value={statusFilter}
                        onChange={(e) => setStatusFilter(e.target.value as FilterStatus)}
                        className="px-4 py-2.5 rounded-xl border border-gray-200 bg-white text-gray-700 focus:border-indigo-500 outline-none"
                    >
                        <option value="all">Tous statuts</option>
                        <option value="active">✅ Actif</option>
                        <option value="frozen">⏸️ Gelé</option>
                        <option value="blocked">🚫 Bloqué</option>
                    </select>
                </div>
            </div>

            {/* Table */}
            <div className="bg-white rounded-2xl border border-gray-200 shadow-sm overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="w-full">
                        <thead>
                            <tr className="bg-gradient-to-r from-gray-50 to-gray-100 border-b border-gray-200">
                                <th className="text-left px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Devise</th>
                                <th className="text-left px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Utilisateur</th>
                                <th className="text-right px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Solde</th>
                                <th className="text-left px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Statut</th>
                                <th className="text-left px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Créé le</th>
                                <th className="text-right px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100">
                            {paginatedWallets.map((wallet) => (
                                <tr key={wallet.id} className="hover:bg-gray-50/50 transition-colors">
                                    <td className="px-6 py-4">
                                        <div className="flex items-center gap-3">
                                            <div className={`w-10 h-10 rounded-xl flex items-center justify-center font-bold text-lg ${CRYPTO_CURRENCIES.includes(wallet.currency)
                                                    ? 'bg-gradient-to-br from-amber-400 to-orange-500 text-white'
                                                    : 'bg-gradient-to-br from-green-400 to-emerald-500 text-white'
                                                }`}>
                                                {getCurrencyIcon(wallet.currency)}
                                            </div>
                                            <div>
                                                <p className="font-bold text-gray-900">{wallet.currency}</p>
                                                {wallet.network && (
                                                    <p className="text-xs text-gray-400">{wallet.network}</p>
                                                )}
                                            </div>
                                        </div>
                                    </td>
                                    <td className="px-6 py-4">
                                        <p className="text-sm text-gray-700">{wallet.user_email || 'N/A'}</p>
                                        <p className="text-xs text-gray-400 font-mono">{wallet.user_id.slice(0, 8)}...</p>
                                    </td>
                                    <td className="px-6 py-4 text-right">
                                        <p className="font-bold text-gray-900">
                                            {(wallet.balance || 0).toLocaleString('fr-FR', {
                                                minimumFractionDigits: CRYPTO_CURRENCIES.includes(wallet.currency) ? 8 : 2,
                                                maximumFractionDigits: CRYPTO_CURRENCIES.includes(wallet.currency) ? 8 : 2,
                                            })}
                                        </p>
                                        <p className="text-xs text-gray-400">{wallet.currency}</p>
                                    </td>
                                    <td className="px-6 py-4">
                                        {getStatusBadge(wallet.status)}
                                    </td>
                                    <td className="px-6 py-4 text-sm text-gray-500">
                                        {wallet.created_at ? format(new Date(wallet.created_at), 'dd/MM/yyyy') : '-'}
                                    </td>
                                    <td className="px-6 py-4">
                                        <div className="flex items-center justify-end gap-2">
                                            <button
                                                onClick={() => handleViewDetails(wallet)}
                                                className="p-2 rounded-lg hover:bg-gray-100 text-gray-500 hover:text-indigo-600 transition-colors"
                                                title="Voir détails"
                                            >
                                                <EyeIcon className="w-5 h-5" />
                                            </button>
                                            {wallet.status !== 'blocked' && wallet.status !== 'frozen' && (
                                                <button
                                                    onClick={() => handleFreezeClick(wallet)}
                                                    className="p-2 rounded-lg hover:bg-amber-50 text-gray-500 hover:text-amber-600 transition-colors"
                                                    title="Geler"
                                                >
                                                    <LockClosedIcon className="w-5 h-5" />
                                                </button>
                                            )}
                                        </div>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>

                {/* Empty State */}
                {filteredWallets.length === 0 && (
                    <div className="py-16 text-center">
                        <WalletIcon className="w-16 h-16 text-gray-200 mx-auto mb-4" />
                        <h3 className="text-lg font-semibold text-gray-900 mb-2">Aucun portefeuille trouvé</h3>
                        <p className="text-gray-500 text-sm">Modifiez vos filtres pour voir plus de résultats</p>
                    </div>
                )}

                {/* Pagination */}
                {totalPages > 1 && (
                    <div className="flex items-center justify-between px-6 py-4 border-t border-gray-100">
                        <p className="text-sm text-gray-500">
                            Affichage de {(currentPage - 1) * ITEMS_PER_PAGE + 1} à{' '}
                            {Math.min(currentPage * ITEMS_PER_PAGE, filteredWallets.length)} sur{' '}
                            {filteredWallets.length} résultats
                        </p>
                        <div className="flex items-center gap-2">
                            <button
                                onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                                disabled={currentPage === 1}
                                className="p-2 rounded-lg border border-gray-200 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                                <ChevronLeftIcon className="w-5 h-5" />
                            </button>
                            <span className="px-4 py-2 text-sm font-medium text-gray-700">
                                Page {currentPage} / {totalPages}
                            </span>
                            <button
                                onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
                                disabled={currentPage === totalPages}
                                className="p-2 rounded-lg border border-gray-200 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                                <ChevronRightIcon className="w-5 h-5" />
                            </button>
                        </div>
                    </div>
                )}
            </div>

            {/* Details Modal */}
            <Modal
                isOpen={detailsOpen}
                onClose={() => setDetailsOpen(false)}
                title="Détails du Portefeuille"
                subtitle={selectedWallet?.currency}
                size="lg"
            >
                {selectedWallet && (
                    <div className="space-y-6">
                        {/* Header */}
                        <div className="flex items-center gap-4 p-4 bg-gradient-to-r from-gray-50 to-gray-100 rounded-xl">
                            <div className={`w-16 h-16 rounded-2xl flex items-center justify-center font-bold text-2xl ${CRYPTO_CURRENCIES.includes(selectedWallet.currency)
                                    ? 'bg-gradient-to-br from-amber-400 to-orange-500 text-white'
                                    : 'bg-gradient-to-br from-green-400 to-emerald-500 text-white'
                                }`}>
                                {getCurrencyIcon(selectedWallet.currency)}
                            </div>
                            <div className="flex-1">
                                <h3 className="text-xl font-bold text-gray-900">{selectedWallet.currency}</h3>
                                <p className="text-gray-500">{selectedWallet.network || 'Standard'}</p>
                            </div>
                            {getStatusBadge(selectedWallet.status)}
                        </div>

                        {/* Balance */}
                        <div className="bg-gradient-to-br from-indigo-500 to-purple-600 rounded-xl p-6 text-white">
                            <p className="text-white/70 text-sm mb-1">Solde</p>
                            <p className="text-3xl font-bold">
                                {(selectedWallet.balance || 0).toLocaleString('fr-FR', {
                                    minimumFractionDigits: 2,
                                    maximumFractionDigits: 8,
                                })}
                                <span className="text-lg ml-2">{selectedWallet.currency}</span>
                            </p>
                        </div>

                        {/* Details Grid */}
                        <div className="grid grid-cols-2 gap-4">
                            <div className="p-4 bg-gray-50 rounded-xl">
                                <p className="text-xs text-gray-400 uppercase font-bold mb-1">ID Wallet</p>
                                <p className="font-mono text-sm text-gray-900 break-all">{selectedWallet.id}</p>
                            </div>
                            <div className="p-4 bg-gray-50 rounded-xl">
                                <p className="text-xs text-gray-400 uppercase font-bold mb-1">ID Utilisateur</p>
                                <p className="font-mono text-sm text-gray-900 break-all">{selectedWallet.user_id}</p>
                            </div>
                            {selectedWallet.address && (
                                <div className="col-span-2 p-4 bg-gray-50 rounded-xl">
                                    <p className="text-xs text-gray-400 uppercase font-bold mb-1">Adresse Blockchain</p>
                                    <p className="font-mono text-sm text-gray-900 break-all">{selectedWallet.address}</p>
                                </div>
                            )}
                            <div className="p-4 bg-gray-50 rounded-xl">
                                <p className="text-xs text-gray-400 uppercase font-bold mb-1">Créé le</p>
                                <p className="text-sm text-gray-900">
                                    {selectedWallet.created_at
                                        ? format(new Date(selectedWallet.created_at), 'dd/MM/yyyy HH:mm')
                                        : '-'
                                    }
                                </p>
                            </div>
                            <div className="p-4 bg-gray-50 rounded-xl">
                                <p className="text-xs text-gray-400 uppercase font-bold mb-1">Dernière mise à jour</p>
                                <p className="text-sm text-gray-900">
                                    {selectedWallet.updated_at
                                        ? format(new Date(selectedWallet.updated_at), 'dd/MM/yyyy HH:mm')
                                        : '-'
                                    }
                                </p>
                            </div>
                        </div>

                        {/* Actions */}
                        <ModalFooter>
                            <ModalButton variant="secondary" onClick={() => setDetailsOpen(false)}>
                                Fermer
                            </ModalButton>
                            {selectedWallet.status === 'active' && (
                                <ModalButton
                                    variant="danger"
                                    onClick={() => {
                                        setDetailsOpen(false);
                                        handleFreezeClick(selectedWallet);
                                    }}
                                >
                                    <LockClosedIcon className="w-4 h-4" />
                                    Geler ce portefeuille
                                </ModalButton>
                            )}
                        </ModalFooter>
                    </div>
                )}
            </Modal>

            {/* Freeze Dialog */}
            <InputDialog
                isOpen={freezeDialogOpen}
                onClose={() => setFreezeDialogOpen(false)}
                onSubmit={handleFreezeSubmit}
                title="Geler le Portefeuille"
                subtitle={`${selectedWallet?.currency} - ${selectedWallet?.user_email || selectedWallet?.user_id}`}
                label="Raison du gel"
                placeholder="Ex: Activité suspecte détectée, enquête en cours..."
                multiline
                rows={3}
                submitText="Confirmer le gel"
                variant="warning"
                loading={actionLoading}
            />

            {/* Toast */}
            {toast && (
                <SimpleToast
                    message={toast.message}
                    type={toast.type}
                    onClose={() => setToast(null)}
                />
            )}
        </div>
    );
}
