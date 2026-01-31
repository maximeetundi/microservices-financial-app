'use client';

import { useEffect, useState, useMemo } from 'react';
import { getCards, freezeCard, blockCard } from '@/lib/api';
import { format } from 'date-fns';
import {
    CreditCardIcon,
    MagnifyingGlassIcon,
    LockClosedIcon,
    XCircleIcon,
    EyeIcon,
    ChevronLeftIcon,
    ChevronRightIcon,
    DevicePhoneMobileIcon,
    BanknotesIcon,
    ShieldExclamationIcon,
} from '@heroicons/react/24/outline';
import { CheckCircleIcon, XCircleIcon as XCircleSolid } from '@heroicons/react/24/solid';
import PageHeader from '@/components/ui/PageHeader';
import StatsCard from '@/components/ui/StatsCard';
import Modal, { ModalFooter, ModalButton } from '@/components/ui/Modal';
import InputDialog from '@/components/ui/InputDialog';
import ConfirmDialog from '@/components/ui/ConfirmDialog';
import { SimpleToast } from '@/components/ui/Toast';

interface Card {
    id: string;
    user_id: string;
    user_email?: string;
    user_name?: string;
    card_type: 'virtual' | 'physical';
    status: string;
    currency: string;
    balance: number;
    last_four?: string;
    card_number_masked?: string;
    expiry_date?: string;
    created_at: string;
    frozen_at?: string;
    blocked_at?: string;
    freeze_reason?: string;
    block_reason?: string;
}

type FilterType = 'all' | 'virtual' | 'physical';
type FilterStatus = 'all' | 'active' | 'frozen' | 'blocked';

const ITEMS_PER_PAGE = 10;

export default function CardsPage() {
    // Data state
    const [cards, setCards] = useState<Card[]>([]);
    const [loading, setLoading] = useState(true);
    const [refreshing, setRefreshing] = useState(false);

    // Filters
    const [search, setSearch] = useState('');
    const [typeFilter, setTypeFilter] = useState<FilterType>('all');
    const [statusFilter, setStatusFilter] = useState<FilterStatus>('all');

    // Pagination
    const [currentPage, setCurrentPage] = useState(1);

    // Modal state
    const [selectedCard, setSelectedCard] = useState<Card | null>(null);
    const [detailsOpen, setDetailsOpen] = useState(false);
    const [freezeDialogOpen, setFreezeDialogOpen] = useState(false);
    const [blockConfirmOpen, setBlockConfirmOpen] = useState(false);
    const [blockReasonDialogOpen, setBlockReasonDialogOpen] = useState(false);
    const [blockReason, setBlockReason] = useState('');
    const [actionLoading, setActionLoading] = useState(false);

    // Toast state
    const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' } | null>(null);

    const fetchCards = async (silent = false) => {
        if (!silent) setLoading(true);
        else setRefreshing(true);

        try {
            const response = await getCards();
            setCards(response.data?.cards || []);
        } catch (error) {
            console.error('Failed to fetch cards:', error);
            setToast({ message: 'Erreur lors du chargement des cartes', type: 'error' });
        } finally {
            setLoading(false);
            setRefreshing(false);
        }
    };

    useEffect(() => {
        fetchCards();
    }, []);

    // Computed stats
    const stats = useMemo(() => {
        const total = cards.length;
        const virtual = cards.filter(c => c.card_type === 'virtual').length;
        const physical = cards.filter(c => c.card_type === 'physical').length;
        const active = cards.filter(c => c.status === 'active').length;
        const frozen = cards.filter(c => c.status === 'frozen').length;
        const blocked = cards.filter(c => c.status === 'blocked').length;
        const totalBalance = cards.reduce((sum, c) => sum + (c.balance || 0), 0);

        return { total, virtual, physical, active, frozen, blocked, totalBalance };
    }, [cards]);

    // Filtered cards
    const filteredCards = useMemo(() => {
        return cards.filter(card => {
            // Search filter
            if (search) {
                const s = search.toLowerCase();
                const matchesSearch =
                    card.id.toLowerCase().includes(s) ||
                    card.user_id.toLowerCase().includes(s) ||
                    card.user_email?.toLowerCase().includes(s) ||
                    card.last_four?.includes(s) ||
                    card.card_number_masked?.toLowerCase().includes(s);
                if (!matchesSearch) return false;
            }

            // Type filter
            if (typeFilter !== 'all' && card.card_type !== typeFilter) return false;

            // Status filter
            if (statusFilter !== 'all' && card.status !== statusFilter) return false;

            return true;
        });
    }, [cards, search, typeFilter, statusFilter]);

    // Pagination
    const totalPages = Math.ceil(filteredCards.length / ITEMS_PER_PAGE);
    const paginatedCards = filteredCards.slice(
        (currentPage - 1) * ITEMS_PER_PAGE,
        currentPage * ITEMS_PER_PAGE
    );

    // Reset page when filters change
    useEffect(() => {
        setCurrentPage(1);
    }, [search, typeFilter, statusFilter]);

    // Handlers
    const handleViewDetails = (card: Card) => {
        setSelectedCard(card);
        setDetailsOpen(true);
    };

    const handleFreezeClick = (card: Card) => {
        setSelectedCard(card);
        setFreezeDialogOpen(true);
    };

    const handleFreezeSubmit = async (reason: string) => {
        if (!selectedCard) return;

        setActionLoading(true);
        try {
            await freezeCard(selectedCard.id, reason);
            setToast({ message: 'Carte gelée avec succès', type: 'success' });
            setFreezeDialogOpen(false);
            await fetchCards(true);
        } catch (error) {
            setToast({ message: 'Erreur lors du gel de la carte', type: 'error' });
        } finally {
            setActionLoading(false);
        }
    };

    const handleBlockClick = (card: Card) => {
        setSelectedCard(card);
        setBlockReasonDialogOpen(true);
    };

    const handleBlockReasonSubmit = (reason: string) => {
        setBlockReason(reason);
        setBlockReasonDialogOpen(false);
        setBlockConfirmOpen(true);
    };

    const handleBlockConfirm = async () => {
        if (!selectedCard) return;

        setActionLoading(true);
        try {
            await blockCard(selectedCard.id, blockReason);
            setToast({ message: 'Carte bloquée définitivement', type: 'success' });
            setBlockConfirmOpen(false);
            setBlockReason('');
            await fetchCards(true);
        } catch (error) {
            setToast({ message: 'Erreur lors du blocage de la carte', type: 'error' });
        } finally {
            setActionLoading(false);
        }
    };

    const getStatusBadge = (status: string) => {
        const config: Record<string, { bg: string; text: string; icon: React.ReactNode; label: string }> = {
            active: {
                bg: 'bg-green-100',
                text: 'text-green-700',
                icon: <CheckCircleIcon className="w-4 h-4" />,
                label: 'Active',
            },
            frozen: {
                bg: 'bg-amber-100',
                text: 'text-amber-700',
                icon: <LockClosedIcon className="w-4 h-4" />,
                label: 'Gelée',
            },
            blocked: {
                bg: 'bg-red-100',
                text: 'text-red-700',
                icon: <XCircleSolid className="w-4 h-4" />,
                label: 'Bloquée',
            },
            pending: {
                bg: 'bg-blue-100',
                text: 'text-blue-700',
                icon: <CreditCardIcon className="w-4 h-4" />,
                label: 'En attente',
            },
        };
        const { bg, text, icon, label } = config[status] || config.pending;
        return (
            <span className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold ${bg} ${text}`}>
                {icon}
                {label}
            </span>
        );
    };

    const getTypeBadge = (type: string) => {
        if (type === 'virtual') {
            return (
                <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-indigo-100 text-indigo-700">
                    <DevicePhoneMobileIcon className="w-4 h-4" />
                    Virtuelle
                </span>
            );
        }
        return (
            <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-purple-100 text-purple-700">
                <CreditCardIcon className="w-4 h-4" />
                Physique
            </span>
        );
    };

    if (loading) {
        return (
            <div className="flex items-center justify-center h-96">
                <div className="text-center">
                    <div className="w-16 h-16 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin mx-auto mb-4" />
                    <p className="text-gray-500">Chargement des cartes...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="pb-10">
            <PageHeader
                title="Cartes Bancaires"
                subtitle={`${cards.length} cartes émises`}
                icon="💳"
                onRefresh={() => fetchCards(true)}
                loading={refreshing}
            />

            {/* Stats Cards */}
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
                <StatsCard
                    title="Total"
                    value={stats.total}
                    icon={CreditCardIcon}
                    color="primary"
                />
                <StatsCard
                    title="Virtuelles"
                    value={stats.virtual}
                    icon={DevicePhoneMobileIcon}
                    color="info"
                />
                <StatsCard
                    title="Physiques"
                    value={stats.physical}
                    icon={CreditCardIcon}
                    color="purple"
                />
                <StatsCard
                    title="Solde Total"
                    value={stats.totalBalance.toLocaleString('fr-FR', { style: 'currency', currency: 'EUR' })}
                    icon={BanknotesIcon}
                    color="success"
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
                            placeholder="Rechercher par ID, email, numéro de carte..."
                            className="w-full pl-10 pr-4 py-2.5 rounded-xl border border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 outline-none transition-all"
                        />
                    </div>

                    {/* Type Filter */}
                    <div className="flex gap-2">
                        {[
                            { id: 'all', label: 'Toutes' },
                            { id: 'virtual', label: '📱 Virtuelles' },
                            { id: 'physical', label: '💳 Physiques' },
                        ].map((option) => (
                            <button
                                key={option.id}
                                onClick={() => setTypeFilter(option.id as FilterType)}
                                className={`px-4 py-2 rounded-xl text-sm font-medium transition-all ${typeFilter === option.id
                                        ? 'bg-indigo-600 text-white shadow-md'
                                        : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                                    }`}
                            >
                                {option.label}
                            </button>
                        ))}
                    </div>

                    {/* Status Filter */}
                    <select
                        value={statusFilter}
                        onChange={(e) => setStatusFilter(e.target.value as FilterStatus)}
                        className="px-4 py-2.5 rounded-xl border border-gray-200 bg-white text-gray-700 focus:border-indigo-500 outline-none"
                    >
                        <option value="all">Tous statuts</option>
                        <option value="active">✅ Active</option>
                        <option value="frozen">⏸️ Gelée</option>
                        <option value="blocked">🚫 Bloquée</option>
                    </select>
                </div>
            </div>

            {/* Table */}
            <div className="bg-white rounded-2xl border border-gray-200 shadow-sm overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="w-full">
                        <thead>
                            <tr className="bg-gradient-to-r from-gray-50 to-gray-100 border-b border-gray-200">
                                <th className="text-left px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Carte</th>
                                <th className="text-left px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Type</th>
                                <th className="text-left px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Titulaire</th>
                                <th className="text-right px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Solde</th>
                                <th className="text-left px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Statut</th>
                                <th className="text-left px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Créée le</th>
                                <th className="text-right px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100">
                            {paginatedCards.map((card) => (
                                <tr key={card.id} className="hover:bg-gray-50/50 transition-colors">
                                    <td className="px-6 py-4">
                                        <div className="flex items-center gap-3">
                                            <div className={`w-12 h-8 rounded-lg flex items-center justify-center text-white font-bold text-xs ${card.card_type === 'virtual'
                                                    ? 'bg-gradient-to-br from-indigo-500 to-purple-600'
                                                    : 'bg-gradient-to-br from-gray-700 to-gray-900'
                                                }`}>
                                                {card.last_four ? `····${card.last_four}` : 'XXXX'}
                                            </div>
                                            <div>
                                                <p className="font-mono text-sm text-gray-900">
                                                    {card.card_number_masked || `**** **** **** ${card.last_four || '****'}`}
                                                </p>
                                                {card.expiry_date && (
                                                    <p className="text-xs text-gray-400">Exp: {card.expiry_date}</p>
                                                )}
                                            </div>
                                        </div>
                                    </td>
                                    <td className="px-6 py-4">
                                        {getTypeBadge(card.card_type)}
                                    </td>
                                    <td className="px-6 py-4">
                                        <p className="text-sm text-gray-700">{card.user_email || card.user_name || 'N/A'}</p>
                                        <p className="text-xs text-gray-400 font-mono">{card.user_id.slice(0, 8)}...</p>
                                    </td>
                                    <td className="px-6 py-4 text-right">
                                        <p className="font-bold text-gray-900">
                                            {(card.balance || 0).toLocaleString('fr-FR', {
                                                style: 'currency',
                                                currency: card.currency || 'EUR',
                                            })}
                                        </p>
                                    </td>
                                    <td className="px-6 py-4">
                                        {getStatusBadge(card.status)}
                                    </td>
                                    <td className="px-6 py-4 text-sm text-gray-500">
                                        {card.created_at ? format(new Date(card.created_at), 'dd/MM/yyyy') : '-'}
                                    </td>
                                    <td className="px-6 py-4">
                                        <div className="flex items-center justify-end gap-2">
                                            <button
                                                onClick={() => handleViewDetails(card)}
                                                className="p-2 rounded-lg hover:bg-gray-100 text-gray-500 hover:text-indigo-600 transition-colors"
                                                title="Voir détails"
                                            >
                                                <EyeIcon className="w-5 h-5" />
                                            </button>
                                            {card.status === 'active' && (
                                                <>
                                                    <button
                                                        onClick={() => handleFreezeClick(card)}
                                                        className="p-2 rounded-lg hover:bg-amber-50 text-gray-500 hover:text-amber-600 transition-colors"
                                                        title="Geler"
                                                    >
                                                        <LockClosedIcon className="w-5 h-5" />
                                                    </button>
                                                    <button
                                                        onClick={() => handleBlockClick(card)}
                                                        className="p-2 rounded-lg hover:bg-red-50 text-gray-500 hover:text-red-600 transition-colors"
                                                        title="Bloquer définitivement"
                                                    >
                                                        <XCircleIcon className="w-5 h-5" />
                                                    </button>
                                                </>
                                            )}
                                            {card.status === 'frozen' && (
                                                <button
                                                    onClick={() => handleBlockClick(card)}
                                                    className="p-2 rounded-lg hover:bg-red-50 text-gray-500 hover:text-red-600 transition-colors"
                                                    title="Bloquer définitivement"
                                                >
                                                    <XCircleIcon className="w-5 h-5" />
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
                {filteredCards.length === 0 && (
                    <div className="py-16 text-center">
                        <CreditCardIcon className="w-16 h-16 text-gray-200 mx-auto mb-4" />
                        <h3 className="text-lg font-semibold text-gray-900 mb-2">Aucune carte trouvée</h3>
                        <p className="text-gray-500 text-sm">Modifiez vos filtres pour voir plus de résultats</p>
                    </div>
                )}

                {/* Pagination */}
                {totalPages > 1 && (
                    <div className="flex items-center justify-between px-6 py-4 border-t border-gray-100">
                        <p className="text-sm text-gray-500">
                            Affichage de {(currentPage - 1) * ITEMS_PER_PAGE + 1} à{' '}
                            {Math.min(currentPage * ITEMS_PER_PAGE, filteredCards.length)} sur{' '}
                            {filteredCards.length} résultats
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
                title="Détails de la Carte"
                subtitle={selectedCard?.card_type === 'virtual' ? 'Carte Virtuelle' : 'Carte Physique'}
                size="lg"
            >
                {selectedCard && (
                    <div className="space-y-6">
                        {/* Card Preview */}
                        <div className={`relative p-6 rounded-2xl text-white overflow-hidden ${selectedCard.card_type === 'virtual'
                                ? 'bg-gradient-to-br from-indigo-600 via-purple-600 to-pink-500'
                                : 'bg-gradient-to-br from-gray-800 via-gray-700 to-gray-900'
                            }`}>
                            {/* Decorative elements */}
                            <div className="absolute top-0 right-0 w-48 h-48 bg-white/10 rounded-full -translate-y-1/2 translate-x-1/2" />
                            <div className="absolute bottom-0 left-0 w-32 h-32 bg-white/5 rounded-full translate-y-1/2 -translate-x-1/2" />

                            <div className="relative">
                                <div className="flex justify-between items-start mb-8">
                                    <div className="w-12 h-8 bg-yellow-400/80 rounded-md" />
                                    {getStatusBadge(selectedCard.status)}
                                </div>

                                <p className="font-mono text-xl tracking-wider mb-4">
                                    {selectedCard.card_number_masked || `**** **** **** ${selectedCard.last_four || '****'}`}
                                </p>

                                <div className="flex justify-between items-end">
                                    <div>
                                        <p className="text-white/60 text-xs uppercase mb-1">Titulaire</p>
                                        <p className="font-semibold">{selectedCard.user_name || selectedCard.user_email || 'N/A'}</p>
                                    </div>
                                    <div className="text-right">
                                        <p className="text-white/60 text-xs uppercase mb-1">Expiration</p>
                                        <p className="font-semibold">{selectedCard.expiry_date || 'MM/YY'}</p>
                                    </div>
                                </div>
                            </div>
                        </div>

                        {/* Balance */}
                        <div className="bg-gradient-to-br from-green-500 to-emerald-600 rounded-xl p-5 text-white">
                            <p className="text-white/70 text-sm mb-1">Solde Disponible</p>
                            <p className="text-2xl font-bold">
                                {(selectedCard.balance || 0).toLocaleString('fr-FR', {
                                    style: 'currency',
                                    currency: selectedCard.currency || 'EUR',
                                })}
                            </p>
                        </div>

                        {/* Details Grid */}
                        <div className="grid grid-cols-2 gap-4">
                            <div className="p-4 bg-gray-50 rounded-xl">
                                <p className="text-xs text-gray-400 uppercase font-bold mb-1">ID Carte</p>
                                <p className="font-mono text-sm text-gray-900 break-all">{selectedCard.id}</p>
                            </div>
                            <div className="p-4 bg-gray-50 rounded-xl">
                                <p className="text-xs text-gray-400 uppercase font-bold mb-1">ID Utilisateur</p>
                                <p className="font-mono text-sm text-gray-900 break-all">{selectedCard.user_id}</p>
                            </div>
                            <div className="p-4 bg-gray-50 rounded-xl">
                                <p className="text-xs text-gray-400 uppercase font-bold mb-1">Type</p>
                                <p className="text-sm text-gray-900">
                                    {selectedCard.card_type === 'virtual' ? 'Virtuelle' : 'Physique'}
                                </p>
                            </div>
                            <div className="p-4 bg-gray-50 rounded-xl">
                                <p className="text-xs text-gray-400 uppercase font-bold mb-1">Créée le</p>
                                <p className="text-sm text-gray-900">
                                    {selectedCard.created_at
                                        ? format(new Date(selectedCard.created_at), 'dd/MM/yyyy HH:mm')
                                        : '-'
                                    }
                                </p>
                            </div>
                        </div>

                        {/* Status Info */}
                        {selectedCard.status === 'frozen' && selectedCard.freeze_reason && (
                            <div className="p-4 bg-amber-50 border border-amber-200 rounded-xl">
                                <p className="text-xs text-amber-600 uppercase font-bold mb-1">Raison du gel</p>
                                <p className="text-sm text-amber-800">{selectedCard.freeze_reason}</p>
                                {selectedCard.frozen_at && (
                                    <p className="text-xs text-amber-500 mt-2">
                                        Gelée le {format(new Date(selectedCard.frozen_at), 'dd/MM/yyyy HH:mm')}
                                    </p>
                                )}
                            </div>
                        )}

                        {selectedCard.status === 'blocked' && selectedCard.block_reason && (
                            <div className="p-4 bg-red-50 border border-red-200 rounded-xl">
                                <p className="text-xs text-red-600 uppercase font-bold mb-1">Raison du blocage</p>
                                <p className="text-sm text-red-800">{selectedCard.block_reason}</p>
                                {selectedCard.blocked_at && (
                                    <p className="text-xs text-red-500 mt-2">
                                        Bloquée le {format(new Date(selectedCard.blocked_at), 'dd/MM/yyyy HH:mm')}
                                    </p>
                                )}
                            </div>
                        )}

                        {/* Actions */}
                        <ModalFooter>
                            <ModalButton variant="secondary" onClick={() => setDetailsOpen(false)}>
                                Fermer
                            </ModalButton>
                            {selectedCard.status === 'active' && (
                                <>
                                    <ModalButton
                                        variant="primary"
                                        onClick={() => {
                                            setDetailsOpen(false);
                                            handleFreezeClick(selectedCard);
                                        }}
                                    >
                                        <LockClosedIcon className="w-4 h-4" />
                                        Geler
                                    </ModalButton>
                                    <ModalButton
                                        variant="danger"
                                        onClick={() => {
                                            setDetailsOpen(false);
                                            handleBlockClick(selectedCard);
                                        }}
                                    >
                                        <XCircleIcon className="w-4 h-4" />
                                        Bloquer
                                    </ModalButton>
                                </>
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
                title="Geler la Carte"
                subtitle={`Carte ****${selectedCard?.last_four || '****'} - ${selectedCard?.user_email || selectedCard?.user_id}`}
                label="Raison du gel"
                placeholder="Ex: Activité suspecte détectée..."
                multiline
                rows={3}
                submitText="Confirmer le gel"
                variant="warning"
                loading={actionLoading}
            />

            {/* Block Reason Dialog */}
            <InputDialog
                isOpen={blockReasonDialogOpen}
                onClose={() => setBlockReasonDialogOpen(false)}
                onSubmit={handleBlockReasonSubmit}
                title="Bloquer la Carte"
                subtitle="Cette action est irréversible"
                label="Raison du blocage"
                placeholder="Ex: Fraude confirmée, utilisation non autorisée..."
                multiline
                rows={3}
                submitText="Continuer"
                variant="danger"
            />

            {/* Block Confirm Dialog */}
            <ConfirmDialog
                isOpen={blockConfirmOpen}
                onClose={() => {
                    setBlockConfirmOpen(false);
                    setBlockReason('');
                }}
                onConfirm={handleBlockConfirm}
                title="Confirmation du blocage"
                message={
                    <div className="space-y-2">
                        <p>Êtes-vous sûr de vouloir <strong className="text-red-600">bloquer définitivement</strong> cette carte ?</p>
                        <p className="text-sm text-gray-500">
                            Cette action est <strong>irréversible</strong>. La carte ne pourra plus jamais être réactivée.
                        </p>
                        <div className="mt-3 p-3 bg-gray-50 rounded-lg text-left">
                            <p className="text-xs text-gray-400 mb-1">Raison :</p>
                            <p className="text-sm text-gray-700">{blockReason}</p>
                        </div>
                    </div>
                }
                confirmText="Bloquer définitivement"
                cancelText="Annuler"
                variant="danger"
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
