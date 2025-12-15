'use client';

import { useEffect, useState } from 'react';
import { getWallets, freezeWallet } from '@/lib/api';
import { format } from 'date-fns';

interface Wallet {
    id: string;
    user_id: string;
    currency: string;
    balance: number;
    status: string;
    created_at: string;
}

export default function WalletsPage() {
    const [wallets, setWallets] = useState<Wallet[]>([]);
    const [loading, setLoading] = useState(true);
    const [actionLoading, setActionLoading] = useState<string | null>(null);

    const fetchWallets = async () => {
        try {
            const response = await getWallets();
            setWallets(response.data.wallets || []);
        } catch (error) {
            console.error('Failed to fetch wallets:', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchWallets();
    }, []);

    const handleFreeze = async (walletId: string) => {
        const reason = prompt('Raison du gel:');
        if (!reason) return;

        setActionLoading(walletId);
        try {
            await freezeWallet(walletId, reason);
            await fetchWallets();
            alert('Portefeuille gelé');
        } catch (error) {
            alert('Erreur lors du gel');
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
                <h1 className="text-2xl font-bold text-slate-900">Portefeuilles</h1>
                <p className="text-slate-500 mt-1">{wallets.length} portefeuilles</p>
            </div>

            <div className="table-container">
                <table className="w-full">
                    <thead className="bg-gray-50 border-b">
                        <tr>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">ID</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Devise</th>
                            <th className="text-right px-6 py-4 text-sm font-medium text-slate-600">Solde</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Status</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Créé le</th>
                            <th className="text-right px-6 py-4 text-sm font-medium text-slate-600">Actions</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-100">
                        {wallets.map((wallet) => (
                            <tr key={wallet.id} className="hover:bg-gray-50">
                                <td className="px-6 py-4">
                                    <p className="font-mono text-sm text-slate-700 truncate max-w-[150px]">{wallet.id}</p>
                                </td>
                                <td className="px-6 py-4 font-medium text-slate-700">{wallet.currency}</td>
                                <td className="px-6 py-4 text-right font-medium">
                                    {wallet.balance?.toLocaleString()} {wallet.currency}
                                </td>
                                <td className="px-6 py-4">
                                    <span className={`badge ${wallet.status === 'active' ? 'badge-success' : 'badge-warning'}`}>
                                        {wallet.status}
                                    </span>
                                </td>
                                <td className="px-6 py-4 text-slate-500 text-sm">
                                    {wallet.created_at ? format(new Date(wallet.created_at), 'dd/MM/yyyy') : '-'}
                                </td>
                                <td className="px-6 py-4 text-right">
                                    <button
                                        onClick={() => handleFreeze(wallet.id)}
                                        disabled={actionLoading === wallet.id || wallet.status === 'frozen'}
                                        className="btn-secondary text-sm px-3 py-1"
                                    >
                                        {actionLoading === wallet.id ? '...' : 'Geler'}
                                    </button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
                {wallets.length === 0 && (
                    <div className="text-center py-12 text-slate-500">
                        Aucun portefeuille trouvé
                    </div>
                )}
            </div>
        </div>
    );
}
