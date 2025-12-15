'use client';

import { useEffect, useState } from 'react';
import { getCards, freezeCard, blockCard } from '@/lib/api';
import { format } from 'date-fns';

interface Card {
    id: string;
    user_id: string;
    card_type: string;
    status: string;
    currency: string;
    balance: number;
    created_at: string;
}

export default function CardsPage() {
    const [cards, setCards] = useState<Card[]>([]);
    const [loading, setLoading] = useState(true);
    const [actionLoading, setActionLoading] = useState<string | null>(null);

    const fetchCards = async () => {
        try {
            const response = await getCards();
            setCards(response.data.cards || []);
        } catch (error) {
            console.error('Failed to fetch cards:', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchCards();
    }, []);

    const handleFreeze = async (cardId: string) => {
        const reason = prompt('Raison du gel:');
        if (!reason) return;

        setActionLoading(cardId);
        try {
            await freezeCard(cardId, reason);
            await fetchCards();
            alert('Carte gelée');
        } catch (error) {
            alert('Erreur lors du gel');
        } finally {
            setActionLoading(null);
        }
    };

    const handleBlock = async (cardId: string) => {
        const reason = prompt('Raison du blocage (irréversible):');
        if (!reason) return;

        if (!confirm('Êtes-vous sûr ? Cette action est irréversible.')) return;

        setActionLoading(cardId);
        try {
            await blockCard(cardId, reason);
            await fetchCards();
            alert('Carte bloquée définitivement');
        } catch (error) {
            alert('Erreur lors du blocage');
        } finally {
            setActionLoading(null);
        }
    };

    if (loading) {
        return (
            <div className="flex items-center justify-center h-64">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
            </div>
        );
    }

    return (
        <div>
            <div className="mb-8">
                <h1 className="text-2xl font-bold text-slate-900">Cartes</h1>
                <p className="text-slate-500 mt-1">{cards.length} cartes émises</p>
            </div>

            <div className="table-container">
                <table className="w-full">
                    <thead className="bg-gray-50 border-b">
                        <tr>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">ID Carte</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Type</th>
                            <th className="text-right px-6 py-4 text-sm font-medium text-slate-600">Solde</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Status</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Créée le</th>
                            <th className="text-right px-6 py-4 text-sm font-medium text-slate-600">Actions</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-100">
                        {cards.map((card) => (
                            <tr key={card.id} className="hover:bg-gray-50">
                                <td className="px-6 py-4">
                                    <p className="font-mono text-sm text-slate-700 truncate max-w-[150px]">{card.id}</p>
                                </td>
                                <td className="px-6 py-4">
                                    <span className={`badge ${card.card_type === 'virtual' ? 'badge-info' : 'badge-success'}`}>
                                        {card.card_type === 'virtual' ? 'Virtuelle' : 'Physique'}
                                    </span>
                                </td>
                                <td className="px-6 py-4 text-right font-medium">
                                    {card.balance?.toLocaleString()} {card.currency}
                                </td>
                                <td className="px-6 py-4">
                                    <span className={`badge ${getStatusBadge(card.status)}`}>
                                        {card.status}
                                    </span>
                                </td>
                                <td className="px-6 py-4 text-slate-500 text-sm">
                                    {card.created_at ? format(new Date(card.created_at), 'dd/MM/yyyy') : '-'}
                                </td>
                                <td className="px-6 py-4 text-right space-x-2">
                                    <button
                                        onClick={() => handleFreeze(card.id)}
                                        disabled={actionLoading === card.id || card.status === 'blocked'}
                                        className="btn-secondary text-sm px-3 py-1"
                                    >
                                        Geler
                                    </button>
                                    <button
                                        onClick={() => handleBlock(card.id)}
                                        disabled={actionLoading === card.id || card.status === 'blocked'}
                                        className="btn-danger text-sm px-3 py-1"
                                    >
                                        Bloquer
                                    </button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
                {cards.length === 0 && (
                    <div className="text-center py-12 text-slate-500">
                        Aucune carte trouvée
                    </div>
                )}
            </div>
        </div>
    );
}

function getStatusBadge(status: string): string {
    switch (status?.toLowerCase()) {
        case 'active':
            return 'badge-success';
        case 'frozen':
            return 'badge-warning';
        case 'blocked':
            return 'badge-danger';
        default:
            return 'badge-info';
    }
}
