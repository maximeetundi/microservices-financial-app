'use client';

import React, { useState, useEffect } from 'react';
import {
    BanknotesIcon,
    QrCodeIcon,
    ArrowsRightLeftIcon,
    PlusIcon,
    ArrowPathIcon,
    WalletIcon,
    ShieldCheckIcon,
    BoltIcon, // For Hot Wallet
    LockClosedIcon, // For Cold Wallet
    ServerStackIcon
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
    syncPlatformCryptoWallet,
    consolidateUserFunds
} from '@/lib/api';

const safeFormatCurrency = (amount: number, currency: string) => {
    try {
        const code = currency === 'FCFA' ? 'XOF' : currency;
        return new Intl.NumberFormat('fr-FR', {
            style: 'currency',
            currency: code
        }).format(amount);
    } catch (e) {
        return new Intl.NumberFormat('fr-FR', {
            style: 'decimal',
            minimumFractionDigits: 2,
            maximumFractionDigits: 6
        }).format(amount) + ' ' + currency;
    }
};

export default function PlatformAccountsPage() {
    const [activeTab, setActiveTab] = useState<'fiat' | 'crypto' | 'transactions'>('crypto'); // Default to crypto for verification
    const [loading, setLoading] = useState(true);

    // Data
    const [fiatAccounts, setFiatAccounts] = useState<any[]>([]);
    const [cryptoWallets, setCryptoWallets] = useState<any[]>([]);
    const [transactions, setTransactions] = useState<any[]>([]);
    const [reconciliation, setReconciliation] = useState<any>({});

    const [showCreateFiat, setShowCreateFiat] = useState(false);
    const [showCreateCrypto, setShowCreateCrypto] = useState(false);
    const [showCreditDebit, setShowCreditDebit] = useState(false);
    const [showConsolidate, setShowConsolidate] = useState(false);
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

    const [consolidateForm, setConsolidateForm] = useState({
        target_type: 'hot',
        amount: 0,
        currency: 'ETH'
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
            alert('Erreur lors de la création du compte');
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
            alert('Erreur lors de la création du portefeuille');
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
            alert('Opération échouée');
        }
    };

    const handleSync = async (id: string) => {
        try {
            await syncPlatformCryptoWallet(id);
            alert('Synchronisation lancée');
            loadData(); // Refresh to see update if immediate
        } catch (error) {
            alert('Échec de la synchronisation');
        }
    };

    const handleConsolidate = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await consolidateUserFunds(consolidateForm);
            alert('Consolidation lancée avec succès');
            setShowConsolidate(false);
            loadData();
        } catch (error) {
            console.error(error);
            alert('Échec de la consolidation des fonds');
        }
    };

    // Helper Component for Compact Display
    const CompactBalance = ({ amount, currency }: { amount: number, currency: string }) => {
        const [showFull, setShowFull] = useState(false);

        const fullValue = safeFormatCurrency(amount, currency);

        const getCompact = () => {
            if (amount >= 1_000_000_000) return (amount / 1_000_000_000).toFixed(2).replace(/\.00$/, '') + ' Md';
            if (amount >= 1_000_000) return (amount / 1_000_000).toFixed(2).replace(/\.00$/, '') + ' M';
            return fullValue;
        };

        return (
            <div
                className="cursor-pointer relative group inline-flex items-baseline gap-1"
                onClick={() => setShowFull(!showFull)}
                title="Cliquez pour basculer"
            >
                <span className={`font-mono font-bold tracking-tight ${amount < 0 ? 'text-red-500' : ''}`}>
                    {showFull ? fullValue : getCompact()}
                </span>
                {!showFull && amount >= 1_000_000 && (
                    <span className="text-[10px] text-slate-400 font-normal uppercase">{currency}</span>
                )}
                {/* Hover Tooltip */}
                {!showFull && amount >= 1_000_000 && (
                    <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-1 px-2 py-1 bg-slate-800 text-white text-xs rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap pointer-events-none z-10">
                        {fullValue}
                    </div>
                )}
            </div>
        );
    };

    const formatCurrency = (amount: number, currency: string) => {
        return safeFormatCurrency(amount, currency);
    };

    // Helper for Crypto Icons
    const getCryptoIcon = (currency: string) => {
        const symbol = currency.toLowerCase();
        // Fallback or use a reliable CDN
        return `https://raw.githubusercontent.com/spothq/cryptocurrency-icons/master/128/color/${symbol}.png`;
    };

    const getNetworkBadge = (network: string) => {
        const n = network.toLowerCase();
        let color = 'bg-slate-100 text-slate-600';
        if (n.includes('bitcoin') || n === 'btc') color = 'bg-orange-100 text-orange-700 border-orange-200';
        else if (n.includes('ethereum') || n === 'eth') color = 'bg-blue-100 text-blue-700 border-blue-200';
        else if (n.includes('bsc') || n === 'bnb') color = 'bg-yellow-100 text-yellow-700 border-yellow-200';
        else if (n.includes('tron') || n === 'trc') color = 'bg-red-100 text-red-700 border-red-200';
        else if (n.includes('solana') || n === 'sol') color = 'bg-purple-100 text-purple-700 border-purple-200';
        else if (n.includes('ton')) color = 'bg-sky-100 text-sky-700 border-sky-200'; // TON Color

        return (
            <span className={`px-2.5 py-1 rounded-md text-[10px] font-bold uppercase border ${color}`}>
                {network}
            </span>
        );
    };

    return (
        <div className="space-y-6 pb-20">
            {/* Header */}
            <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
                <div>
                    <h1 className="text-2xl font-bold text-slate-900 flex items-center gap-3">
                        <ServerStackIcon className="w-8 h-8 text-indigo-600" />
                        Comptes Plateforme & Liquidité
                    </h1>
                    <p className="text-slate-500 mt-1">Superviser les réserves Fiat et les portefeuilles Crypto "Hot/Cold"</p>
                </div>
                <div className="flex gap-3">
                    <button
                        onClick={loadData}
                        className="p-2.5 rounded-xl bg-slate-50 border border-slate-200 text-slate-500 hover:text-indigo-600 hover:bg-indigo-50 transition-all"
                        title="Actualiser"
                    >
                        <ArrowPathIcon className={`w-6 h-6 ${loading ? 'animate-spin' : ''}`} />
                    </button>
                    {activeTab === 'fiat' && (
                        <button onClick={() => setShowCreateFiat(true)} className="flex items-center gap-2 px-5 py-2.5 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700 font-bold shadow-lg shadow-indigo-200 hover:shadow-indigo-300 transition-all">
                            <PlusIcon className="w-5 h-5" /> Nouveau Compte Fiat
                        </button>
                    )}
                    {activeTab === 'crypto' && (
                        <>
                            <button onClick={() => setShowConsolidate(true)} className="flex items-center gap-2 px-5 py-2.5 bg-white border border-slate-200 text-slate-700 rounded-xl hover:bg-slate-50 font-bold shadow-sm transition-all">
                                <BanknotesIcon className="w-5 h-5" /> Consolider Fonds
                            </button>
                            <button onClick={() => setShowCreateCrypto(true)} className="flex items-center gap-2 px-5 py-2.5 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700 font-bold shadow-lg shadow-indigo-200 hover:shadow-indigo-300 transition-all">
                                <PlusIcon className="w-5 h-5" /> Ajouter Wallet
                            </button>
                        </>
                    )}
                </div>
            </div>

            {/* Tabs */}
            <div className="flex p-1 bg-slate-100 rounded-xl w-fit">
                <button
                    onClick={() => setActiveTab('crypto')}
                    className={`px-6 py-2.5 rounded-lg text-sm font-bold transition-all ${activeTab === 'crypto' ? 'bg-white text-indigo-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}
                >
                    <div className="flex items-center gap-2">
                        <QrCodeIcon className="w-5 h-5" /> Portefeuilles Crypto
                    </div>
                </button>
                <button
                    onClick={() => setActiveTab('fiat')}
                    className={`px-6 py-2.5 rounded-lg text-sm font-bold transition-all ${activeTab === 'fiat' ? 'bg-white text-indigo-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}
                >
                    <div className="flex items-center gap-2">
                        <BanknotesIcon className="w-5 h-5" /> Comptes Fiat
                    </div>
                </button>
                <button
                    onClick={() => setActiveTab('transactions')}
                    className={`px-6 py-2.5 rounded-lg text-sm font-bold transition-all ${activeTab === 'transactions' ? 'bg-white text-indigo-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}
                >
                    <div className="flex items-center gap-2">
                        <ArrowsRightLeftIcon className="w-5 h-5" /> Transactions & Réconciliation
                    </div>
                </button>
            </div>

            {/* Content: Crypto Wallets */}
            {activeTab === 'crypto' && (
                <div className="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
                    <div className="overflow-x-auto">
                        <table className="w-full text-left">
                            <thead className="bg-slate-50 border-b border-slate-100">
                                <tr>
                                    <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">Monaie</th>
                                    <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">Réseau / Adresse</th>
                                    <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">Type (Stockage/Ops)</th>
                                    <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider text-right">Solde</th>
                                    <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider text-right">Actions</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-slate-100">
                                {cryptoWallets.map(wallet => (
                                    <tr key={wallet.id} className="hover:bg-slate-50/80 transition-colors">
                                        <td className="px-6 py-4">
                                            <div className="flex items-center gap-3">
                                                <div className="w-10 h-10 rounded-full bg-slate-50 flex items-center justify-center border border-slate-100 p-2 overflow-hidden">
                                                    <img
                                                        src={getCryptoIcon(wallet.currency)}
                                                        alt={wallet.currency}
                                                        className="w-full h-full object-contain"
                                                        onError={(e: any) => {
                                                            (e.target as HTMLImageElement).style.display = 'none';
                                                            (e.target as HTMLImageElement).parentElement!.innerText = wallet.currency.substring(0, 2);
                                                        }}
                                                    />
                                                </div>
                                                <div>
                                                    <div className="font-bold text-slate-900">{wallet.label || wallet.currency}</div>
                                                    <div className="text-xs text-slate-500 font-mono">ID: {wallet.currency}</div>
                                                </div>
                                            </div>
                                        </td>
                                        <td className="px-6 py-4">
                                            <div className="flex flex-col gap-2">
                                                <div>{getNetworkBadge(wallet.network)}</div>
                                                {wallet.address ? (
                                                    <div className="flex items-center gap-2 group cursor-pointer" title={wallet.address}>
                                                        <code className="text-xs text-slate-600 bg-slate-100 px-2 py-1 rounded border border-slate-200 font-mono max-w-[200px] truncate block">
                                                            {wallet.address}
                                                        </code>
                                                    </div>
                                                ) : (
                                                    <span className="text-xs text-slate-400 italic">Adresse non générée</span>
                                                )}
                                            </div>
                                        </td>
                                        <td className="px-6 py-4">
                                            {wallet.wallet_type === 'cold' ? (
                                                <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-indigo-50 text-indigo-700 border border-indigo-100 text-xs font-bold">
                                                    <LockClosedIcon className="w-3.5 h-3.5" /> Cold Storage
                                                </div>
                                            ) : (
                                                <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-orange-50 text-orange-700 border border-orange-100 text-xs font-bold">
                                                    <BoltIcon className="w-3.5 h-3.5" /> Hot Wallet
                                                </div>
                                            )}
                                        </td>
                                        <td className="px-6 py-4 text-right">
                                            <div className="font-mono font-bold text-slate-900">
                                                <CompactBalance amount={Number(wallet.balance) || 0} currency={wallet.currency} />
                                            </div>
                                        </td>
                                        <td className="px-6 py-4 text-right">
                                            <div className="flex items-center justify-end gap-2">
                                                <button
                                                    onClick={() => handleSync(wallet.id)}
                                                    className="px-3 py-1.5 rounded-lg bg-white border border-slate-200 text-slate-600 text-xs font-bold hover:bg-slate-50 transition-colors"
                                                >
                                                    Sync
                                                </button>
                                                <button
                                                    onClick={() => { setSelectedAccount(wallet); setOpMode('credit'); setShowCreditDebit(true); }} // Reusing CreditDebit for now as wrapper
                                                    className="px-3 py-1.5 rounded-lg bg-indigo-50 border border-indigo-100 text-indigo-600 text-xs font-bold hover:bg-indigo-100 transition-colors"
                                                >
                                                    Détails
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                                {cryptoWallets.length === 0 && (
                                    <tr>
                                        <td colSpan={5} className="px-6 py-8 text-center text-slate-500 italic">
                                            Aucun portefeuille crypto trouvé. Initialisez le système.
                                        </td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>
            )}

            {/* Content: Fiat Accounts */}
            {activeTab === 'fiat' && (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {fiatAccounts.map(account => (
                        <div key={account.id} className="group bg-white p-6 rounded-2xl border border-slate-200 shadow-sm hover:shadow-xl hover:border-indigo-200 transition-all duration-300">
                            <div className="flex justify-between items-start mb-4">
                                <div className="p-3 bg-indigo-50 rounded-xl text-indigo-600 group-hover:scale-110 transition-transform">
                                    <BanknotesIcon className="w-6 h-6" />
                                </div>
                                <span className={`px-2.5 py-1 rounded-full text-[10px] font-bold uppercase tracking-wide border ${account.account_type === 'reserve' ? 'bg-purple-50 text-purple-700 border-purple-100' :
                                    account.account_type === 'operations' ? 'bg-emerald-50 text-emerald-700 border-emerald-100' :
                                        'bg-slate-50 text-slate-600 border-slate-200'
                                    }`}>
                                    {account.account_type === 'reserve' ? 'Réserve (Stockage)' :
                                        account.account_type === 'operations' ? 'Opérationnel' :
                                            account.account_type}
                                </span>
                            </div>

                            <h3 className="font-bold text-slate-900 text-lg mb-1 line-clamp-1">{account.name}</h3>
                            <div className="text-3xl font-extrabold text-slate-800 mb-6 tracking-tight">
                                <CompactBalance amount={account.balance} currency={account.currency} />
                            </div>

                            <div className="space-y-3 mb-6 bg-slate-50 p-4 rounded-xl border border-slate-100">
                                <div className="flex justify-between items-center text-sm">
                                    <span className="text-slate-500 font-medium">Priorité</span>
                                    <span className="font-mono font-bold text-slate-700 bg-white px-2 py-0.5 rounded shadow-sm border border-slate-100">{account.priority}</span>
                                </div>
                                <div className="flex justify-between items-center text-sm">
                                    <span className="text-slate-500 font-medium">Devise Base</span>
                                    <span className="font-bold text-indigo-600">{account.currency}</span>
                                </div>
                            </div>

                            <div className="grid grid-cols-2 gap-3">
                                <button
                                    onClick={() => { setSelectedAccount(account); setOpMode('credit'); setShowCreditDebit(true); }}
                                    className="px-4 py-2 bg-emerald-50 text-emerald-700 border border-emerald-100 hover:bg-emerald-100 rounded-xl text-sm font-bold transition-all"
                                >
                                    Créditer
                                </button>
                                <button
                                    onClick={() => { setSelectedAccount(account); setOpMode('debit'); setShowCreditDebit(true); }}
                                    className="px-4 py-2 bg-red-50 text-red-700 border border-red-100 hover:bg-red-100 rounded-xl text-sm font-bold transition-all"
                                >
                                    Débiter
                                </button>
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {/* Content: Transactions */}
            {activeTab === 'transactions' && (
                <div className="space-y-8">
                    {/* Reconciliation Summary */}
                    <div className="bg-slate-900 text-white p-8 rounded-3xl shadow-2xl relative overflow-hidden">
                        <div className="absolute top-0 right-0 w-64 h-64 bg-indigo-500 rounded-full blur-[100px] opacity-20 -mr-20 -mt-20"></div>

                        <h3 className="text-xl font-bold mb-8 flex items-center gap-3 relative z-10">
                            <ShieldCheckIcon className="w-6 h-6 text-emerald-400" />
                            Réconciliation Globale des Fonds
                        </h3>

                        <div className="grid grid-cols-2 md:grid-cols-4 gap-8 relative z-10">
                            {Object.entries(reconciliation).map(([currency, amount]) => (
                                <div key={currency} className="bg-white/5 p-4 rounded-2xl backdrop-blur-sm border border-white/10">
                                    <div className="text-slate-400 text-xs font-bold uppercase tracking-wider mb-2">{currency} Total Passif</div>
                                    <div className="text-2xl font-mono font-bold text-emerald-400">
                                        {formatCurrency(Number(amount), currency)}
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>

                    <div className="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
                        <div className="px-6 py-4 border-b border-slate-100 bg-slate-50/50">
                            <h3 className="font-bold text-slate-800">Journal des Opérations</h3>
                        </div>
                        <table className="w-full text-left">
                            <thead className="bg-slate-50 text-slate-500 text-xs uppercase font-bold tracking-wider">
                                <tr>
                                    <th className="px-6 py-4">Date</th>
                                    <th className="px-6 py-4">Opération</th>
                                    <th className="px-6 py-4">Source (Débit)</th>
                                    <th className="px-6 py-4">Destination (Crédit)</th>
                                    <th className="px-6 py-4 text-right">Montant</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-slate-100">
                                {transactions.map(tx => (
                                    <tr key={tx.id} className="hover:bg-slate-50 transition-colors">
                                        <td className="px-6 py-4 text-sm text-slate-500 font-medium">
                                            {new Date(tx.created_at || tx.transaction_date).toLocaleString()}
                                        </td>
                                        <td className="px-6 py-4">
                                            <span className="px-2.5 py-1 rounded-md text-xs font-bold bg-slate-100 text-slate-700 uppercase tracking-wide border border-slate-200">
                                                {tx.operation_type}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 text-sm font-mono text-slate-600">
                                            {tx.debit_account_type}
                                        </td>
                                        <td className="px-6 py-4 text-sm font-mono text-slate-600">
                                            {tx.credit_account_type}
                                        </td>
                                        <td className="px-6 py-4 text-right font-mono font-bold text-slate-900">
                                            {formatCurrency(tx.amount, tx.currency)}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            )}

            {/* Modals - Simplified for display, functionally fully wired */}
            {showCreateFiat && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4 transition-all">
                    <div className="bg-white rounded-2xl p-8 w-full max-w-lg shadow-2xl animate-slide-up">
                        <h2 className="text-2xl font-bold mb-6 text-slate-900">Créer Compte Fiat</h2>
                        <form onSubmit={handleCreateFiat} className="space-y-6">
                            <div>
                                <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Nom du Compte</label>
                                <input className="w-full px-4 py-3 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none transition-all font-semibold"
                                    value={fiatForm.name} onChange={(e: any) => setFiatForm({ ...fiatForm, name: e.target.value })} required placeholder="Ex: Réserve Principale EUR" />
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Devise</label>
                                    <select className="w-full px-4 py-3 border border-slate-200 rounded-xl bg-slate-50 font-medium outline-none"
                                        value={fiatForm.currency} onChange={(e: any) => setFiatForm({ ...fiatForm, currency: e.target.value })}>
                                        <option value="EUR">EUR</option>
                                        <option value="USD">USD</option>
                                        <option value="FCFA">FCFA</option>
                                        <option value="XOF">XOF</option>
                                    </select>
                                </div>
                                <div>
                                    <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Type</label>
                                    <select className="w-full px-4 py-3 border border-slate-200 rounded-xl bg-slate-50 font-medium outline-none"
                                        value={fiatForm.account_type} onChange={(e: any) => setFiatForm({ ...fiatForm, account_type: e.target.value })}>
                                        <option value="reserve">Réserve (Stockage)</option>
                                        <option value="operations">Opérationnel</option>
                                        <option value="fees">Frais (Accumulateur)</option>
                                    </select>
                                </div>
                            </div>
                            <div className="flex justify-end gap-3 mt-8">
                                <button type="button" onClick={() => setShowCreateFiat(false)} className="px-5 py-2.5 text-slate-600 hover:bg-slate-50 rounded-xl font-bold transition-colors">Annuler</button>
                                <button type="submit" className="px-6 py-2.5 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700 font-bold shadow-lg shadow-indigo-500/30 transition-all">Créer le Compte</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            {/* Create Crypto Modal */}
            {showCreateCrypto && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4 transition-all">
                    <div className="bg-white rounded-2xl p-8 w-full max-w-lg shadow-2xl animate-slide-up">
                        <h2 className="text-2xl font-bold mb-6 text-slate-900">Ajouter Portefeuille Crypto</h2>
                        <form onSubmit={handleCreateCrypto} className="space-y-6">
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Devise / Token</label>
                                    <input className="w-full px-4 py-3 border border-slate-200 rounded-xl font-bold uppercase transition-all"
                                        value={cryptoForm.currency} onChange={(e: any) => setCryptoForm({ ...cryptoForm, currency: e.target.value.toUpperCase() })} required placeholder="Ex: ETH" />
                                </div>
                                <div>
                                    <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Réseau</label>
                                    <input className="w-full px-4 py-3 border border-slate-200 rounded-xl font-bold uppercase transition-all"
                                        value={cryptoForm.network} onChange={(e: any) => setCryptoForm({ ...cryptoForm, network: e.target.value.toLowerCase() })} required placeholder="Ex: ethereum" />
                                </div>
                            </div>

                            <div>
                                <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Type de Portefeuille</label>
                                <div className="grid grid-cols-2 gap-4">
                                    <button type="button"
                                        onClick={() => setCryptoForm({ ...cryptoForm, wallet_type: 'hot' })}
                                        className={`px-4 py-3 rounded-xl border-2 font-bold transition-all ${cryptoForm.wallet_type === 'hot' ? 'border-orange-500 bg-orange-50 text-orange-700' : 'border-slate-200 text-slate-500'}`}>
                                        Hot Wallet (Ops)
                                    </button>
                                    <button type="button"
                                        onClick={() => setCryptoForm({ ...cryptoForm, wallet_type: 'cold' })}
                                        className={`px-4 py-3 rounded-xl border-2 font-bold transition-all ${cryptoForm.wallet_type === 'cold' ? 'border-indigo-500 bg-indigo-50 text-indigo-700' : 'border-slate-200 text-slate-500'}`}>
                                        Cold Storage
                                    </button>
                                </div>
                            </div>

                            <div>
                                <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Label (Optionnel)</label>
                                <input className="w-full px-4 py-3 border border-slate-200 rounded-xl font-medium transition-all"
                                    value={cryptoForm.label} onChange={(e: any) => setCryptoForm({ ...cryptoForm, label: e.target.value })} placeholder="Ex: ETH Hot Wallet 2" />
                            </div>

                            <div className="flex justify-end gap-3 mt-8">
                                <button type="button" onClick={() => setShowCreateCrypto(false)} className="px-5 py-2.5 text-slate-600 hover:bg-slate-50 rounded-xl font-bold transition-colors">Annuler</button>
                                <button type="submit" className="px-6 py-2.5 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700 font-bold shadow-lg shadow-indigo-500/30 transition-all">Générer Clés & Créer</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
            {/* Consolidate Funds Modal */}
            {showConsolidate && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4 transition-all">
                    <div className="bg-white rounded-2xl p-8 w-full max-w-lg shadow-2xl animate-slide-up">
                        <h2 className="text-2xl font-bold mb-6 text-slate-900 flex items-center gap-2">
                            <BanknotesIcon className="w-6 h-6 text-indigo-600" />
                            Consolider Fonds Utilisateurs
                        </h2>
                        <p className="text-sm text-slate-500 mb-6 bg-slate-50 p-4 rounded-xl border border-slate-100">
                            Cette opération va balayer les fonds disponibles sur les portefeuilles virtuels des utilisateurs vers le portefeuille de plateforme sélectionné.
                            Seuls les montants excédant les seuils minimums seront transférés (interne base de données).
                        </p>
                        <form onSubmit={handleConsolidate} className="space-y-6">
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Devise</label>
                                    <input className="w-full px-4 py-3 border border-slate-200 rounded-xl font-bold uppercase transition-all"
                                        value={consolidateForm.currency} onChange={(e: any) => setConsolidateForm({ ...consolidateForm, currency: e.target.value.toUpperCase() })} required placeholder="Ex: ETH" />
                                </div>
                                <div>
                                    <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Cible</label>
                                    <select className="w-full px-4 py-3 border border-slate-200 rounded-xl bg-slate-50 font-medium outline-none"
                                        value={consolidateForm.target_type} onChange={(e: any) => setConsolidateForm({ ...consolidateForm, target_type: e.target.value })}>
                                        <option value="hot">Hot Wallet</option>
                                        <option value="cold">Cold Wallet</option>
                                    </select>
                                </div>
                            </div>

                            <div>
                                <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Montant Minimum à Consolider</label>
                                <div className="relative">
                                    <input type="number" step="any" min="0" className="w-full pl-4 pr-16 py-3 border border-slate-200 rounded-xl text-xl font-bold focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none transition-all"
                                        value={consolidateForm.amount} onChange={(e: any) => setConsolidateForm({ ...consolidateForm, amount: parseFloat(e.target.value) })} required placeholder="0.00" />
                                    <div className="absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 font-bold">
                                        {consolidateForm.currency}
                                    </div>
                                </div>
                                <p className="text-xs text-slate-400 mt-1">Laissez 0 pour tout consolider.</p>
                            </div>

                            <div className="flex justify-end gap-3 mt-8">
                                <button type="button" onClick={() => setShowConsolidate(false)} className="px-5 py-2.5 text-slate-600 hover:bg-slate-50 rounded-xl font-bold transition-colors">Annuler</button>
                                <button type="submit" className="px-6 py-2.5 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700 font-bold shadow-lg shadow-indigo-500/30 transition-all">
                                    Lancer Consolidation
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            {/* Credit/Debit Operation Modal */}
            {showCreditDebit && selectedAccount && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4 transition-all">
                    <div className="bg-white rounded-2xl p-8 w-full max-w-md shadow-2xl animate-slide-up">
                        <h2 className={`text-2xl font-bold mb-6 ${opMode === 'credit' ? 'text-emerald-600' : 'text-red-600'}`}>
                            {opMode === 'credit' ? 'Créditer / Recevoir' : 'Débiter / Envoyer'}
                        </h2>
                        <div className="mb-6 p-4 bg-slate-50 rounded-xl border border-slate-100 flex items-center gap-4">
                            {selectedAccount.network ? (
                                <div className="w-12 h-12 p-2 bg-white rounded-full border border-slate-100 mb-0">
                                    <img src={getCryptoIcon(selectedAccount.currency)} className="w-full h-full object-contain" onError={(e) => { (e.target as HTMLImageElement).style.display = 'none'; }} />
                                </div>
                            ) : null}
                            <div>
                                <div className="text-sm text-slate-500 font-bold uppercase mb-1">Compte Cible</div>
                                <div className="font-bold text-slate-900 text-lg">{selectedAccount.name || selectedAccount.label || selectedAccount.currency}</div>
                                <div className="text-sm text-slate-600 font-mono mt-1">{selectedAccount.currency} ({selectedAccount.account_type || selectedAccount.wallet_type})</div>
                            </div>
                        </div>

                        <form onSubmit={handleOperation} className="space-y-6">
                            <div>
                                <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Montant</label>
                                <div className="relative">
                                    <input type="number" step="any" min="0" className="w-full pl-4 pr-16 py-3 border border-slate-200 rounded-xl text-xl font-bold focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none transition-all"
                                        value={opForm.amount} onChange={(e: any) => setOpForm({ ...opForm, amount: parseFloat(e.target.value) })} required placeholder="0.00" />
                                    <div className="absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 font-bold">
                                        {selectedAccount.currency}
                                    </div>
                                </div>
                            </div>
                            <div>
                                <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Description / Motif</label>
                                <input className="w-full px-4 py-3 border border-slate-200 rounded-xl font-medium focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none transition-all"
                                    value={opForm.description} onChange={(e: any) => setOpForm({ ...opForm, description: e.target.value })} required placeholder="Ex: Injection de liquidité initiale" />
                            </div>

                            {/* Address Display for Credit (Receive) on Crypto */}
                            {opMode === 'credit' && selectedAccount.network && (
                                <div className="p-3 bg-blue-50 rounded-xl border border-blue-100 text-sm">
                                    <div className="font-bold text-blue-700 mb-1">Adresse de réception :</div>
                                    <code className="block break-all font-mono text-blue-600 bg-white p-2 rounded border border-blue-100">
                                        {selectedAccount.address}
                                    </code>
                                </div>
                            )}

                            <div className="flex justify-end gap-3 mt-8">
                                <button type="button" onClick={() => setShowCreditDebit(false)} className="px-5 py-2.5 text-slate-600 hover:bg-slate-50 rounded-xl font-bold transition-colors">Annuler</button>
                                <button type="submit"
                                    className={`px-6 py-2.5 text-white rounded-xl font-bold shadow-lg transition-all ${opMode === 'credit' ? 'bg-emerald-600 hover:bg-emerald-700 shadow-emerald-500/30' : 'bg-red-600 hover:bg-red-700 shadow-red-500/30'}`}>
                                    Confirmer {opMode === 'credit' ? 'Crédit' : 'Débit'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}
