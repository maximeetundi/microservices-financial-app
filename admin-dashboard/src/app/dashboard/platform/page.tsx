'use client';

import { useState, useEffect } from 'react';
import {
    BanknotesIcon,
    QrCodeIcon,
    ArrowsRightLeftIcon,
    PlusIcon,
    ArrowPathIcon,
    WalletIcon
} from '@heroicons/react/24/outline';
import {
    getPlatformAccounts,
    getPlatformCryptoWallets,
    getPlatformTransactions,
    getPlatformReconciliation,
    createPlatformAccount,
    createPlatformCryptoWallet,
    creditPlatformAccount,
    debitPlatformAccount,
    syncPlatformCryptoWallet
} from '@/lib/api';

export default function PlatformAccountsPage() {
    const [activeTab, setActiveTab] = useState<'fiat' | 'crypto' | 'transactions'>('fiat');
    const [loading, setLoading] = useState(true);

    // Data
    const [fiatAccounts, setFiatAccounts] = useState<any[]>([]);
    const [cryptoWallets, setCryptoWallets] = useState<any[]>([]);
    const [transactions, setTransactions] = useState<any[]>([]);
    const [reconciliation, setReconciliation] = useState<any>({});

    // Modals
    const [showCreateFiat, setShowCreateFiat] = useState(false);
    const [showCreateCrypto, setShowCreateCrypto] = useState(false);
    const [showCreditDebit, setShowCreditDebit] = useState(false);
    const [selectedAccount, setSelectedAccount] = useState<any>(null);
    const [opMode, setOpMode] = useState<'credit' | 'debit'>('credit');

    // Forms
    const [fiatForm, setFiatForm] = useState({
        currency: 'EUR',
        account_type: 'reserve',
        name: '',
        description: '',
        priority: 50,
        min_balance: 0,
        max_balance: 0
    });

    const [cryptoForm, setCryptoForm] = useState({
        currency: 'ETH',
        network: 'ethereum',
        wallet_type: 'hot',
        address: '',
        label: '',
        priority: 50,
        min_balance: 0,
        max_balance: 0
    });

    const [opForm, setOpForm] = useState({
        amount: 0,
        description: '',
        reference: ''
    });

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        setLoading(true);
        try {
            const [accRes, walletRes, txRes, recRes] = await Promise.all([
                getPlatformAccounts(),
                getPlatformCryptoWallets(),
                getPlatformTransactions(),
                getPlatformReconciliation()
            ]);

            setFiatAccounts(accRes.data.accounts || []);
            setCryptoWallets(walletRes.data.wallets || []);
            setTransactions(txRes.data.transactions || []);
            setReconciliation(recRes.data.balances || {});
        } catch (error) {
            console.error('Failed to load platform data', error);
        } finally {
            setLoading(false);
        }
    };

    const handleCreateFiat = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await createPlatformAccount(fiatForm);
            setShowCreateFiat(false);
            loadData();
        } catch (error) {
            console.error(error);
            alert('Failed to create account');
        }
    };

    const handleCreateCrypto = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await createPlatformCryptoWallet(cryptoForm);
            setShowCreateCrypto(false);
            loadData();
        } catch (error) {
            console.error(error);
            alert('Failed to create wallet');
        }
    };

    const handleOperation = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!selectedAccount) return;
        try {
            if (opMode === 'credit') {
                await creditPlatformAccount(selectedAccount.id, opForm);
            } else {
                await debitPlatformAccount(selectedAccount.id, opForm);
            }
            setShowCreditDebit(false);
            loadData();
        } catch (error) {
            console.error(error);
            alert('Operation failed');
        }
    };

    const handleSync = async (id: string) => {
        try {
            await syncPlatformCryptoWallet(id);
            alert('Sync started');
            loadData(); // Refresh to see update if immediate
        } catch (error) {
            alert('Sync failed');
        }
    };

    const formatCurrency = (amount: number, currency: string) => {
        return new Intl.NumberFormat('fr-FR', {
            style: 'currency',
            currency: currency === 'FCFA' ? 'XOF' : currency
        }).format(amount);
    };

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-2xl font-bold text-slate-900">Platform Accounts</h1>
                    <p className="text-slate-500">Manage company liquidity and verify funds</p>
                </div>
                <div className="flex gap-2">
                    <button onClick={loadData} className="p-2 rounded-lg border hover:bg-slate-50">
                        <ArrowPathIcon className={`w-5 h-5 ${loading ? 'animate-spin' : ''}`} />
                    </button>
                    {activeTab === 'fiat' && (
                        <button onClick={() => setShowCreateFiat(true)} className="flex items-center gap-2 px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700">
                            <PlusIcon className="w-5 h-5" /> New Fiat Account
                        </button>
                    )}
                    {activeTab === 'crypto' && (
                        <button onClick={() => setShowCreateCrypto(true)} className="flex items-center gap-2 px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700">
                            <PlusIcon className="w-5 h-5" /> New Wallet
                        </button>
                    )}
                </div>
            </div>

            {/* Tabs */}
            <div className="flex border-b border-slate-200">
                <button
                    onClick={() => setActiveTab('fiat')}
                    className={`px-6 py-3 text-sm font-medium border-b-2 transition-colors ${activeTab === 'fiat' ? 'border-indigo-600 text-indigo-600' : 'border-transparent text-slate-500 hover:text-slate-700'}`}
                >
                    <div className="flex items-center gap-2">
                        <BanknotesIcon className="w-4 h-4" /> Fiat Accounts
                    </div>
                </button>
                <button
                    onClick={() => setActiveTab('crypto')}
                    className={`px-6 py-3 text-sm font-medium border-b-2 transition-colors ${activeTab === 'crypto' ? 'border-indigo-600 text-indigo-600' : 'border-transparent text-slate-500 hover:text-slate-700'}`}
                >
                    <div className="flex items-center gap-2">
                        <QrCodeIcon className="w-4 h-4" /> Crypto Wallets
                    </div>
                </button>
                <button
                    onClick={() => setActiveTab('transactions')}
                    className={`px-6 py-3 text-sm font-medium border-b-2 transition-colors ${activeTab === 'transactions' ? 'border-indigo-600 text-indigo-600' : 'border-transparent text-slate-500 hover:text-slate-700'}`}
                >
                    <div className="flex items-center gap-2">
                        <ArrowsRightLeftIcon className="w-4 h-4" /> Ledger & Reconciliation
                    </div>
                </button>
            </div>

            {/* Content */}
            {activeTab === 'fiat' && (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {fiatAccounts.map(account => (
                        <div key={account.id} className="bg-white p-6 rounded-xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow">
                            <div className="flex justify-between items-start mb-4">
                                <div className="p-2 bg-indigo-50 rounded-lg text-indigo-600">
                                    <BanknotesIcon className="w-6 h-6" />
                                </div>
                                <span className={`px-2 py-1 rounded-full text-xs font-semibold ${account.is_active ? 'bg-emerald-50 text-emerald-600' : 'bg-slate-100 text-slate-500'}`}>
                                    {account.account_type.toUpperCase()}
                                </span>
                            </div>
                            <h3 className="font-bold text-slate-900 mb-1">{account.name}</h3>
                            <div className="text-2xl font-bold text-slate-800 mb-4">{formatCurrency(account.balance, account.currency)}</div>
                            <div className="text-sm text-slate-500 space-y-1 mb-6">
                                <div className="flex justify-between">
                                    <span>Currency</span>
                                    <span className="font-mono">{account.currency}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span>Priority</span>
                                    <span className="font-mono">P{account.priority}</span>
                                </div>
                            </div>
                            <div className="grid grid-cols-2 gap-3">
                                <button
                                    onClick={() => { setSelectedAccount(account); setOpMode('credit'); setShowCreditDebit(true); }}
                                    className="px-3 py-2 bg-emerald-50 text-emerald-600 hover:bg-emerald-100 rounded-lg text-sm font-medium transition-colors"
                                >
                                    Credit
                                </button>
                                <button
                                    onClick={() => { setSelectedAccount(account); setOpMode('debit'); setShowCreditDebit(true); }}
                                    className="px-3 py-2 bg-red-50 text-red-600 hover:bg-red-100 rounded-lg text-sm font-medium transition-colors"
                                >
                                    Debit
                                </button>
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {activeTab === 'crypto' && (
                <div className="bg-white rounded-xl border border-slate-200 overflow-hidden">
                    <table className="w-full text-left">
                        <thead className="bg-slate-50 text-slate-500 text-xs uppercase font-medium">
                            <tr>
                                <th className="px-6 py-4">Wallet</th>
                                <th className="px-6 py-4">Network/Address</th>
                                <th className="px-6 py-4">Type</th>
                                <th className="px-6 py-4 text-right">Balance</th>
                                <th className="px-6 py-4 text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-slate-100">
                            {cryptoWallets.map(wallet => (
                                <tr key={wallet.id} className="hover:bg-slate-50">
                                    <td className="px-6 py-4">
                                        <div className="font-medium text-slate-900">{wallet.label || wallet.currency}</div>
                                        <div className="text-xs text-slate-500">{wallet.currency}</div>
                                    </td>
                                    <td className="px-6 py-4">
                                        <div className="text-sm text-slate-900">{wallet.network}</div>
                                        <code className="text-xs text-slate-500 bg-slate-100 px-1 py-0.5 rounded">{wallet.address}</code>
                                    </td>
                                    <td className="px-6 py-4">
                                        <span className="px-2 py-1 rounded-full text-xs bg-slate-100 text-slate-600">{wallet.wallet_type}</span>
                                    </td>
                                    <td className="px-6 py-4 text-right font-mono font-medium">
                                        {wallet.balance} {wallet.currency}
                                    </td>
                                    <td className="px-6 py-4 text-right">
                                        <button
                                            onClick={() => handleSync(wallet.id)}
                                            className="text-indigo-600 hover:text-indigo-800 text-sm font-medium"
                                        >
                                            Sync
                                        </button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}

            {activeTab === 'transactions' && (
                <div className="space-y-6">
                    {/* Reconciliation Summary */}
                    <div className="bg-slate-900 text-white p-6 rounded-xl">
                        <h3 className="text-lg font-bold mb-4 flex items-center gap-2">
                            <WalletIcon className="w-5 h-5 text-emerald-400" />
                            Global Reconciliation
                        </h3>
                        <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
                            {Object.entries(reconciliation).map(([currency, amount]) => (
                                <div key={currency}>
                                    <div className="text-slate-400 text-xs uppercase mb-1">{currency} Total Liability</div>
                                    <div className="text-xl font-mono font-bold text-emerald-400">
                                        {formatCurrency(Number(amount), currency)}
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>

                    <div className="bg-white rounded-xl border border-slate-200 overflow-hidden">
                        <table className="w-full text-left">
                            <thead className="bg-slate-50 text-slate-500 text-xs uppercase font-medium">
                                <tr>
                                    <th className="px-6 py-4">Date</th>
                                    <th className="px-6 py-4">Operation</th>
                                    <th className="px-6 py-4">From (Debit)</th>
                                    <th className="px-6 py-4">To (Credit)</th>
                                    <th className="px-6 py-4 text-right">Amount</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-slate-100">
                                {transactions.map(tx => (
                                    <tr key={tx.id} className="hover:bg-slate-50">
                                        <td className="px-6 py-4 text-sm text-slate-500">
                                            {new Date(tx.created_at || tx.transaction_date).toLocaleString()}
                                        </td>
                                        <td className="px-6 py-4">
                                            <span className="px-2 py-1 rounded text-xs font-medium bg-slate-100 text-slate-700">
                                                {tx.operation_type}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 text-sm">
                                            {tx.debit_account_type}
                                        </td>
                                        <td className="px-6 py-4 text-sm">
                                            {tx.credit_account_type}
                                        </td>
                                        <td className="px-6 py-4 text-right font-mono font-medium text-slate-900">
                                            {formatCurrency(tx.amount, tx.currency)}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            )}

            {/* Modals would go here (simplified for brevity) - You can implement full forms if needed */}
            {showCreateFiat && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
                    <div className="bg-white rounded-xl p-6 w-full max-w-md">
                        <h2 className="text-xl font-bold mb-4">Create Fiat Account</h2>
                        <form onSubmit={handleCreateFiat} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium mb-1">Name</label>
                                <input className="w-full p-2 border rounded" value={fiatForm.name} onChange={e => setFiatForm({ ...fiatForm, name: e.target.value })} required />
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-sm font-medium mb-1">Currency</label>
                                    <select className="w-full p-2 border rounded" value={fiatForm.currency} onChange={e => setFiatForm({ ...fiatForm, currency: e.target.value })}>
                                        <option value="EUR">EUR</option>
                                        <option value="USD">USD</option>
                                        <option value="FCFA">FCFA</option>
                                        <option value="XOF">XOF</option>
                                    </select>
                                </div>
                                <div>
                                    <label className="block text-sm font-medium mb-1">Type</label>
                                    <select className="w-full p-2 border rounded" value={fiatForm.account_type} onChange={e => setFiatForm({ ...fiatForm, account_type: e.target.value })}>
                                        <option value="reserve">Reserve</option>
                                        <option value="fees">Fees</option>
                                        <option value="operations">Operations</option>
                                    </select>
                                </div>
                            </div>
                            <div className="flex justify-end gap-2 mt-6">
                                <button type="button" onClick={() => setShowCreateFiat(false)} className="px-4 py-2 text-slate-600 hover:bg-slate-100 rounded">Cancel</button>
                                <button type="submit" className="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700">Create</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            {showCreditDebit && selectedAccount && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
                    <div className="bg-white rounded-xl p-6 w-full max-w-md">
                        <h2 className="text-xl font-bold mb-1">{opMode === 'credit' ? 'Credit' : 'Debit'} Account</h2>
                        <p className="text-slate-500 mb-4">{selectedAccount.name} ({selectedAccount.currency})</p>
                        <form onSubmit={handleOperation} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium mb-1">Amount</label>
                                <input type="number" step="0.01" className="w-full p-2 border rounded" value={opForm.amount} onChange={e => setOpForm({ ...opForm, amount: parseFloat(e.target.value) })} required />
                            </div>
                            <div>
                                <label className="block text-sm font-medium mb-1">Description</label>
                                <input className="w-full p-2 border rounded" value={opForm.description} onChange={e => setOpForm({ ...opForm, description: e.target.value })} required />
                            </div>
                            <div className="flex justify-end gap-2 mt-6">
                                <button type="button" onClick={() => setShowCreditDebit(false)} className="px-4 py-2 text-slate-600 hover:bg-slate-100 rounded">Cancel</button>
                                <button type="submit" className={`px-4 py-2 text-white rounded ${opMode === 'credit' ? 'bg-emerald-600 hover:bg-emerald-700' : 'bg-red-600 hover:bg-red-700'}`}>Confirm</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}
