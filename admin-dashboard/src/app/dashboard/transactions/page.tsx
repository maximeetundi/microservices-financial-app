'use client';

import { useEffect, useState } from 'react';
import { getTransactions, blockTransaction, refundTransaction } from '@/lib/api';
import { format } from 'date-fns';

interface Transaction {
    id: string;
    wallet_id: string;
    type: string;
    amount: number;
    currency: string;
    status: string;
    created_at: string;
}

export default function TransactionsPage() {
    const [transactions, setTransactions] = useState<Transaction[]>([]);
    const [loading, setLoading] = useState(true);
    const [actionLoading, setActionLoading] = useState<string | null>(null);

    const fetchTransactions = async () => {
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

    const handleBlock = async (txId: string) => {
        const reason = prompt('Raison du blocage:');
        if (!reason) return;

        setActionLoading(txId);
        try {
            await blockTransaction(txId, reason);
            await fetchTransactions();
            alert('Transaction bloquée');
        } catch (error) {
            alert('Erreur lors du blocage');
        } finally {
            setActionLoading(null);
        }
    };

    const handleRefund = async (txId: string) => {
        const reason = prompt('Raison du remboursement:');
        if (!reason) return;

        setActionLoading(txId);
        try {
            await refundTransaction(txId, reason);
            await fetchTransactions();
            alert('Remboursement initié');
        } catch (error) {
            alert('Erreur lors du remboursement');
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
                <h1 className="text-2xl font-bold text-slate-900">Transactions</h1>
                <p className="text-slate-500 mt-1">{transactions.length} transactions</p>
            </div>

            <div className="table-container">
                <table className="w-full">
                    <thead className="bg-gray-50 border-b">
                        <tr>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">ID</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Type</th>
                            <th className="text-right px-6 py-4 text-sm font-medium text-slate-600">Montant</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Status</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Date</th>
                            <th className="text-right px-6 py-4 text-sm font-medium text-slate-600">Actions</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-100">
                        {transactions.map((tx) => (
                            <tr key={tx.id} className="hover:bg-gray-50">
                                <td className="px-6 py-4">
                                    <p className="font-mono text-sm text-slate-700 truncate max-w-[150px]">{tx.id}</p>
                                </td>
                                <td className="px-6 py-4">
                                    <span className={`badge ${getTypeBadge(tx.type)}`}>
                                        {tx.type}
                                    </span>
                                </td>
                                <td className="px-6 py-4 text-right font-medium">
                                    {tx.amount?.toLocaleString()} {tx.currency}
                                </td>
                                <td className="px-6 py-4">
                                    <span className={`badge ${getStatusBadge(tx.status)}`}>
                                        {tx.status}
                                    </span>
                                </td>
                                <td className="px-6 py-4 text-slate-500 text-sm">
                                    {tx.created_at ? format(new Date(tx.created_at), 'dd/MM/yyyy HH:mm') : '-'}
                                </td>
                                <td className="px-6 py-4 text-right space-x-2">
                                    <button
                                        onClick={() => handleBlock(tx.id)}
                                        disabled={actionLoading === tx.id}
                                        className="btn-danger text-sm px-3 py-1"
                                    >
                                        Bloquer
                                    </button>
                                    <button
                                        onClick={() => handleRefund(tx.id)}
                                        disabled={actionLoading === tx.id}
                                        className="btn-secondary text-sm px-3 py-1"
                                    >
                                        Rembourser
                                    </button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
                {transactions.length === 0 && (
                    <div className="text-center py-12 text-slate-500">
                        Aucune transaction trouvée
                    </div>
                )}
            </div>
        </div>
    );
}

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
        default:
            return 'badge-warning';
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
