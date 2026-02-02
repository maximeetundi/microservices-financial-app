'use client';

import { useState, useEffect } from 'react';
import {
    SparklesIcon,
    BoltIcon,
    CreditCardIcon,
    UsersIcon,
    ClockIcon,
    UserIcon,
    CheckCircleIcon,
    XCircleIcon,
    ArrowPathIcon,
    BanknotesIcon,
    FunnelIcon,
    MagnifyingGlassIcon,
    CalendarIcon,
    GlobeAltIcon,
    ChartBarIcon,
    GiftIcon,
    TrophyIcon,
    ExclamationTriangleIcon,
} from '@heroicons/react/24/outline';
import clsx from 'clsx';

// Types
interface ReasonType {
    id: string;
    label: string;
    icon: string;
}

interface HotWallet {
    id: string;
    name: string;
    currency: string;
    balance: number;
    is_active: boolean;
}

interface Campaign {
    id: string;
    name: string;
    type: string;
    status: string;
    reason: string;
    reason_type: string;
    currency: string;
    total_amount: number;
    user_count: number;
    success_count: number;
    failed_count: number;
    created_at: string;
    completed_at: string | null;
}

interface CreditLog {
    id: string;
    campaign_id: string;
    campaign_name: string;
    user_id: string;
    wallet_id: string;
    currency: string;
    amount: number;
    status: string;
    error_message: string;
    reason_type: string;
    created_at: string;
}

interface MassPreviewResult {
    user_count: number;
    amount_per_user: number;
    total_amount: number;
    currency: string;
    hot_wallet_balance: number;
    sufficient_funds: boolean;
    users_preview: { user_id: string; email: string; balance: number }[];
}

export default function CreditsPage() {
    const [activeTab, setActiveTab] = useState<'single' | 'mass' | 'promotion' | 'campaigns' | 'logs'>('single');
    const [loading, setLoading] = useState(false);
    const [successMsg, setSuccessMsg] = useState('');
    const [errorMsg, setErrorMsg] = useState('');

    // Data
    const [reasonTypes, setReasonTypes] = useState<ReasonType[]>([]);
    const [hotWallets, setHotWallets] = useState<HotWallet[]>([]);
    const [campaigns, setCampaigns] = useState<Campaign[]>([]);
    const [logs, setLogs] = useState<CreditLog[]>([]);

    // Single Credit State
    const [scUserId, setScUserId] = useState('');
    const [scCurrency, setScCurrency] = useState('XOF');
    const [scAmount, setScAmount] = useState('');
    const [scReason, setScReason] = useState('');
    const [scReasonType, setScReasonType] = useState('compensation');

    // Mass Credit State
    const [mcCurrency, setMcCurrency] = useState('XOF');
    const [mcAmount, setMcAmount] = useState('');
    const [mcCampaignName, setMcCampaignName] = useState('');
    const [mcReason, setMcReason] = useState('');
    const [mcReasonType, setMcReasonType] = useState('compensation');
    const [mcFilters, setMcFilters] = useState({
        userTypes: [] as string[],
        countries: [] as string[],
        txDateFrom: '',
        txDateTo: '',
        minTxAmount: '',
        kycStatus: '',
    });
    const [mcPreview, setMcPreview] = useState<MassPreviewResult | null>(null);

    // Promotion State
    const [prCampaignName, setPrCampaignName] = useState('');
    const [prUserIds, setPrUserIds] = useState('');
    const [prUniformAmount, setPrUniformAmount] = useState('');
    const [prCurrency, setPrCurrency] = useState('XOF');
    const [prReason, setPrReason] = useState('');
    const [prReasonType, setPrReasonType] = useState('contest');

    const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';

    useEffect(() => {
        fetchReasonTypes();
        fetchHotWallets();
        fetchCampaigns();
        fetchLogs();
    }, []);

    const getAuthHeaders = () => {
        const token = localStorage.getItem('admin_token');
        return {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
        };
    };

    const fetchReasonTypes = async () => {
        try {
            const res = await fetch(`${API_URL}/api/v1/admin/credits/reason-types`, { headers: getAuthHeaders() });
            if (res.ok) {
                const data = await res.json();
                setReasonTypes(data.reason_types || []);
            }
        } catch (err) {
            console.error('Failed to fetch reason types:', err);
        }
    };

    const fetchHotWallets = async () => {
        try {
            const res = await fetch(`${API_URL}/api/v1/admin/credits/hot-wallets`, { headers: getAuthHeaders() });
            if (res.ok) {
                const data = await res.json();
                setHotWallets(data.hot_wallets || []);
            }
        } catch (err) {
            console.error('Failed to fetch hot wallets:', err);
        }
    };

    const fetchCampaigns = async () => {
        try {
            const res = await fetch(`${API_URL}/api/v1/admin/credits/campaigns`, { headers: getAuthHeaders() });
            if (res.ok) {
                const data = await res.json();
                setCampaigns(data.campaigns || []);
            }
        } catch (err) {
            console.error('Failed to fetch campaigns:', err);
        }
    };

    const fetchLogs = async () => {
        try {
            const res = await fetch(`${API_URL}/api/v1/admin/credits/logs`, { headers: getAuthHeaders() });
            if (res.ok) {
                const data = await res.json();
                setLogs(data.logs || []);
            }
        } catch (err) {
            console.error('Failed to fetch logs:', err);
        }
    };

    const handleSingleCredit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        resetMessages();

        try {
            const res = await fetch(`${API_URL}/api/v1/admin/credits/single`, {
                method: 'POST',
                headers: getAuthHeaders(),
                body: JSON.stringify({
                    user_id: scUserId,
                    currency: scCurrency,
                    amount: parseFloat(scAmount),
                    reason: scReason,
                    reason_type: scReasonType,
                }),
            });

            const data = await res.json();
            if (res.ok) {
                setSuccessMsg(data.message || `Crédit de ${scAmount} ${scCurrency} effectué avec succès!`);
                fetchCampaigns();
                fetchLogs();
                fetchHotWallets();
                // Reset form
                setScUserId('');
                setScAmount('');
                setScReason('');
            } else {
                setErrorMsg(data.error || 'Erreur lors du crédit');
            }
        } catch (err) {
            setErrorMsg('Erreur de connexion');
        } finally {
            setLoading(false);
        }
    };

    const handleMassPreview = async () => {
        setLoading(true);
        resetMessages();

        try {
            const filters: Record<string, unknown> = {};
            if (mcFilters.userTypes.length > 0) filters.user_types = mcFilters.userTypes;
            if (mcFilters.countries.length > 0) filters.countries = mcFilters.countries;
            if (mcFilters.txDateFrom) filters.tx_date_from = new Date(mcFilters.txDateFrom).toISOString();
            if (mcFilters.txDateTo) filters.tx_date_to = new Date(mcFilters.txDateTo).toISOString();
            if (mcFilters.minTxAmount) filters.min_tx_amount = parseFloat(mcFilters.minTxAmount);
            if (mcFilters.kycStatus) filters.kyc_status = mcFilters.kycStatus;

            const res = await fetch(`${API_URL}/api/v1/admin/credits/mass/preview`, {
                method: 'POST',
                headers: getAuthHeaders(),
                body: JSON.stringify({
                    filters,
                    amount: parseFloat(mcAmount),
                    currency: mcCurrency,
                    reason: mcReason,
                    campaign_name: mcCampaignName,
                }),
            });

            const data = await res.json();
            if (res.ok) {
                setMcPreview(data);
            } else {
                setErrorMsg(data.error || 'Erreur lors de la prévisualisation');
            }
        } catch (err) {
            setErrorMsg('Erreur de connexion');
        } finally {
            setLoading(false);
        }
    };

    const handleMassCredit = async () => {
        if (!mcPreview) return;
        setLoading(true);
        resetMessages();

        try {
            const filters: Record<string, unknown> = {};
            if (mcFilters.userTypes.length > 0) filters.user_types = mcFilters.userTypes;
            if (mcFilters.countries.length > 0) filters.countries = mcFilters.countries;
            if (mcFilters.txDateFrom) filters.tx_date_from = new Date(mcFilters.txDateFrom).toISOString();
            if (mcFilters.txDateTo) filters.tx_date_to = new Date(mcFilters.txDateTo).toISOString();
            if (mcFilters.minTxAmount) filters.min_tx_amount = parseFloat(mcFilters.minTxAmount);
            if (mcFilters.kycStatus) filters.kyc_status = mcFilters.kycStatus;

            const res = await fetch(`${API_URL}/api/v1/admin/credits/mass`, {
                method: 'POST',
                headers: getAuthHeaders(),
                body: JSON.stringify({
                    filters,
                    amount: parseFloat(mcAmount),
                    currency: mcCurrency,
                    reason: mcReason,
                    reason_type: mcReasonType,
                    campaign_name: mcCampaignName,
                }),
            });

            const data = await res.json();
            if (res.ok) {
                setSuccessMsg(data.message || `Crédit de masse lancé pour ${data.user_count} utilisateurs!`);
                setMcPreview(null);
                fetchCampaigns();
                fetchHotWallets();
                // Reset form
                setMcCampaignName('');
                setMcAmount('');
                setMcReason('');
            } else {
                setErrorMsg(data.error || 'Erreur lors du crédit de masse');
            }
        } catch (err) {
            setErrorMsg('Erreur de connexion');
        } finally {
            setLoading(false);
        }
    };

    const handlePromotionCredit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        resetMessages();

        const userIds = prUserIds.split(/[\n,]/).map(id => id.trim()).filter(id => id);

        try {
            const res = await fetch(`${API_URL}/api/v1/admin/credits/promotion`, {
                method: 'POST',
                headers: getAuthHeaders(),
                body: JSON.stringify({
                    campaign_name: prCampaignName,
                    user_ids: userIds,
                    uniform_amount: parseFloat(prUniformAmount),
                    currency: prCurrency,
                    reason: prReason,
                    reason_type: prReasonType,
                }),
            });

            const data = await res.json();
            if (res.ok) {
                setSuccessMsg(data.message || `Promotion lancée pour ${data.user_count} utilisateurs!`);
                fetchCampaigns();
                fetchHotWallets();
                // Reset form
                setPrCampaignName('');
                setPrUserIds('');
                setPrUniformAmount('');
                setPrReason('');
            } else {
                setErrorMsg(data.error || 'Erreur lors de la promotion');
            }
        } catch (err) {
            setErrorMsg('Erreur de connexion');
        } finally {
            setLoading(false);
        }
    };

    const resetMessages = () => {
        setSuccessMsg('');
        setErrorMsg('');
    };

    const getStatusBadge = (status: string) => {
        switch (status) {
            case 'completed':
                return <span className="badge badge-success">Terminé</span>;
            case 'processing':
                return <span className="badge badge-info">En cours</span>;
            case 'failed':
                return <span className="badge badge-danger">Échoué</span>;
            default:
                return <span className="badge badge-gray">En attente</span>;
        }
    };

    const getReasonTypeIcon = (type: string) => {
        const found = reasonTypes.find(r => r.id === type);
        return found?.icon || '📝';
    };

    const currencies = ['XOF', 'XAF', 'EUR', 'USD', 'NGN', 'GHS', 'KES', 'GBP'];
    const userTypes = ['individual', 'association', 'enterprise'];
    const countries = ['CI', 'SN', 'BF', 'ML', 'TG', 'BJ', 'NE', 'CM', 'GA', 'CD', 'FR'];

    // Get current wallet balance for selected currency
    const getWalletBalance = (currency: string) => {
        const wallet = hotWallets.find(w => w.currency === currency);
        return wallet?.balance || 0;
    };

    return (
        <div className="space-y-6 animate-fadeIn">
            {/* Header */}
            <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
                        <BanknotesIcon className="w-8 h-8 text-emerald-500" />
                        Gestion des Crédits
                    </h1>
                    <p className="text-gray-500 mt-1">Créditez les utilisateurs : compensation, promotion, concours & bonus.</p>
                </div>

                {/* Hot Wallets Summary */}
                <div className="flex flex-wrap gap-2">
                    {hotWallets.slice(0, 4).map(wallet => (
                        <div key={wallet.id} className="px-3 py-2 bg-gradient-to-r from-emerald-50 to-teal-50 rounded-xl border border-emerald-100">
                            <div className="text-xs text-gray-500">{wallet.currency}</div>
                            <div className="text-sm font-bold text-emerald-700">
                                {wallet.balance.toLocaleString()}
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            {/* Tabs */}
            <div className="flex flex-wrap gap-2 border-b border-gray-200 pb-1">
                <button
                    onClick={() => setActiveTab('single')}
                    className={clsx(
                        'flex items-center gap-2 px-4 py-2 rounded-t-lg text-sm font-medium transition-colors',
                        activeTab === 'single'
                            ? 'bg-emerald-50 text-emerald-700 border-b-2 border-emerald-600'
                            : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'
                    )}
                >
                    <UserIcon className="w-4 h-4" />
                    Crédit Individuel
                </button>
                <button
                    onClick={() => setActiveTab('mass')}
                    className={clsx(
                        'flex items-center gap-2 px-4 py-2 rounded-t-lg text-sm font-medium transition-colors',
                        activeTab === 'mass'
                            ? 'bg-emerald-50 text-emerald-700 border-b-2 border-emerald-600'
                            : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'
                    )}
                >
                    <UsersIcon className="w-4 h-4" />
                    Crédit de Masse
                </button>
                <button
                    onClick={() => setActiveTab('promotion')}
                    className={clsx(
                        'flex items-center gap-2 px-4 py-2 rounded-t-lg text-sm font-medium transition-colors',
                        activeTab === 'promotion'
                            ? 'bg-emerald-50 text-emerald-700 border-b-2 border-emerald-600'
                            : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'
                    )}
                >
                    <TrophyIcon className="w-4 h-4" />
                    Promotions & Concours
                </button>
                <button
                    onClick={() => setActiveTab('campaigns')}
                    className={clsx(
                        'flex items-center gap-2 px-4 py-2 rounded-t-lg text-sm font-medium transition-colors',
                        activeTab === 'campaigns'
                            ? 'bg-emerald-50 text-emerald-700 border-b-2 border-emerald-600'
                            : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'
                    )}
                >
                    <ChartBarIcon className="w-4 h-4" />
                    Campagnes
                </button>
                <button
                    onClick={() => {
                        setActiveTab('logs');
                        fetchLogs();
                    }}
                    className={clsx(
                        'flex items-center gap-2 px-4 py-2 rounded-t-lg text-sm font-medium transition-colors',
                        activeTab === 'logs'
                            ? 'bg-emerald-50 text-emerald-700 border-b-2 border-emerald-600'
                            : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'
                    )}
                >
                    <ClockIcon className="w-4 h-4" />
                    Logs d'Audit
                </button>
            </div>

            {/* Messages */}
            {successMsg && (
                <div className="p-4 bg-emerald-50 border border-emerald-200 rounded-xl flex items-center gap-3 text-emerald-700 animate-slide-up">
                    <CheckCircleIcon className="w-5 h-5 flex-shrink-0" />
                    <p>{successMsg}</p>
                </div>
            )}
            {errorMsg && (
                <div className="p-4 bg-red-50 border border-red-200 rounded-xl flex items-center gap-3 text-red-700 animate-slide-up">
                    <XCircleIcon className="w-5 h-5 flex-shrink-0" />
                    <p>{errorMsg}</p>
                </div>
            )}

            {/* Content */}
            <div className="card min-h-[500px]">

                {/* SINGLE CREDIT TAB */}
                {activeTab === 'single' && (
                    <div className="max-w-xl mx-auto space-y-6">
                        <div className="text-center mb-8">
                            <CreditCardIcon className="w-12 h-12 text-emerald-500 mx-auto mb-3" />
                            <h2 className="text-xl font-bold text-gray-900">Crédit Individuel</h2>
                            <p className="text-gray-500 text-sm">Créditez un utilisateur spécifique pour compensation, bonus, etc.</p>
                        </div>

                        <form onSubmit={handleSingleCredit} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">ID Utilisateur</label>
                                <div className="relative">
                                    <UserIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                                    <input
                                        type="text"
                                        required
                                        className="input pl-10"
                                        placeholder="UUID de l'utilisateur"
                                        value={scUserId}
                                        onChange={e => setScUserId(e.target.value)}
                                    />
                                </div>
                            </div>

                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Devise</label>
                                    <select
                                        className="select"
                                        value={scCurrency}
                                        onChange={e => setScCurrency(e.target.value)}
                                    >
                                        {currencies.map(c => (
                                            <option key={c} value={c}>{c} ({getWalletBalance(c).toLocaleString()})</option>
                                        ))}
                                    </select>
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Montant</label>
                                    <input
                                        type="number"
                                        required
                                        className="input"
                                        placeholder="0.00"
                                        min="0.01"
                                        step="any"
                                        value={scAmount}
                                        onChange={e => setScAmount(e.target.value)}
                                    />
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Type de Raison</label>
                                <div className="grid grid-cols-4 gap-2">
                                    {reasonTypes.map(rt => (
                                        <button
                                            key={rt.id}
                                            type="button"
                                            onClick={() => setScReasonType(rt.id)}
                                            className={clsx(
                                                'p-2 rounded-lg border text-center transition-all text-sm',
                                                scReasonType === rt.id
                                                    ? 'bg-emerald-50 border-emerald-500 text-emerald-700'
                                                    : 'border-gray-200 hover:border-emerald-200'
                                            )}
                                        >
                                            <span className="text-lg">{rt.icon}</span>
                                            <div className="text-xs mt-1 truncate">{rt.label.split(' ')[0]}</div>
                                        </button>
                                    ))}
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Raison / Description</label>
                                <input
                                    type="text"
                                    required
                                    className="input"
                                    placeholder="Ex: Compensation panne du 01/02/2026"
                                    value={scReason}
                                    onChange={e => setScReason(e.target.value)}
                                />
                            </div>

                            <button
                                type="submit"
                                disabled={loading}
                                className="w-full btn-primary py-3 flex items-center justify-center gap-2 mt-4 bg-gradient-to-r from-emerald-500 to-teal-600"
                            >
                                {loading && <div className="spinner w-5 h-5 border-white border-t-transparent" />}
                                Créditer l'utilisateur
                            </button>
                        </form>
                    </div>
                )}

                {/* MASS CREDIT TAB */}
                {activeTab === 'mass' && (
                    <div className="space-y-6">
                        <div className="text-center mb-6">
                            <UsersIcon className="w-12 h-12 text-indigo-500 mx-auto mb-3" />
                            <h2 className="text-xl font-bold text-gray-900">Crédit de Masse</h2>
                            <p className="text-gray-500 text-sm">Créditez plusieurs utilisateurs selon des critères spécifiques.</p>
                        </div>

                        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                            {/* Left: Filters */}
                            <div className="space-y-4">
                                <h3 className="font-bold text-gray-800 flex items-center gap-2">
                                    <FunnelIcon className="w-5 h-5" />
                                    Filtres de Sélection
                                </h3>

                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Type d'utilisateur</label>
                                    <div className="flex flex-wrap gap-2">
                                        {userTypes.map(ut => (
                                            <button
                                                key={ut}
                                                type="button"
                                                onClick={() => {
                                                    const newTypes = mcFilters.userTypes.includes(ut)
                                                        ? mcFilters.userTypes.filter(t => t !== ut)
                                                        : [...mcFilters.userTypes, ut];
                                                    setMcFilters({ ...mcFilters, userTypes: newTypes });
                                                }}
                                                className={clsx(
                                                    'px-3 py-1.5 rounded-lg border text-sm',
                                                    mcFilters.userTypes.includes(ut)
                                                        ? 'bg-indigo-100 border-indigo-500 text-indigo-700'
                                                        : 'border-gray-200 hover:border-indigo-200'
                                                )}
                                            >
                                                {ut === 'individual' ? 'Particulier' : ut === 'association' ? 'Association' : 'Entreprise'}
                                            </button>
                                        ))}
                                    </div>
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Pays</label>
                                    <div className="flex flex-wrap gap-2">
                                        {countries.map(c => (
                                            <button
                                                key={c}
                                                type="button"
                                                onClick={() => {
                                                    const newCountries = mcFilters.countries.includes(c)
                                                        ? mcFilters.countries.filter(ct => ct !== c)
                                                        : [...mcFilters.countries, c];
                                                    setMcFilters({ ...mcFilters, countries: newCountries });
                                                }}
                                                className={clsx(
                                                    'px-3 py-1.5 rounded-lg border text-sm',
                                                    mcFilters.countries.includes(c)
                                                        ? 'bg-indigo-100 border-indigo-500 text-indigo-700'
                                                        : 'border-gray-200 hover:border-indigo-200'
                                                )}
                                            >
                                                {c}
                                            </button>
                                        ))}
                                    </div>
                                </div>

                                <div className="grid grid-cols-2 gap-3">
                                    <div>
                                        <label className="block text-sm font-medium text-gray-700 mb-1">Transactions Du</label>
                                        <input
                                            type="date"
                                            className="input"
                                            value={mcFilters.txDateFrom}
                                            onChange={e => setMcFilters({ ...mcFilters, txDateFrom: e.target.value })}
                                        />
                                    </div>
                                    <div>
                                        <label className="block text-sm font-medium text-gray-700 mb-1">Au</label>
                                        <input
                                            type="date"
                                            className="input"
                                            value={mcFilters.txDateTo}
                                            onChange={e => setMcFilters({ ...mcFilters, txDateTo: e.target.value })}
                                        />
                                    </div>
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Montant Transaction Min.</label>
                                    <input
                                        type="number"
                                        className="input"
                                        placeholder="0"
                                        value={mcFilters.minTxAmount}
                                        onChange={e => setMcFilters({ ...mcFilters, minTxAmount: e.target.value })}
                                    />
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Statut KYC</label>
                                    <select
                                        className="select"
                                        value={mcFilters.kycStatus}
                                        onChange={e => setMcFilters({ ...mcFilters, kycStatus: e.target.value })}
                                    >
                                        <option value="">Tous</option>
                                        <option value="approved">Approuvé</option>
                                        <option value="pending">En attente</option>
                                        <option value="rejected">Rejeté</option>
                                    </select>
                                </div>
                            </div>

                            {/* Right: Credit Config */}
                            <div className="space-y-4">
                                <h3 className="font-bold text-gray-800 flex items-center gap-2">
                                    <BanknotesIcon className="w-5 h-5" />
                                    Configuration du Crédit
                                </h3>

                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Nom de la Campagne</label>
                                    <input
                                        type="text"
                                        className="input"
                                        placeholder="Ex: Compensation panne serveurs"
                                        value={mcCampaignName}
                                        onChange={e => setMcCampaignName(e.target.value)}
                                    />
                                </div>

                                <div className="grid grid-cols-2 gap-3">
                                    <div>
                                        <label className="block text-sm font-medium text-gray-700 mb-1">Devise</label>
                                        <select
                                            className="select"
                                            value={mcCurrency}
                                            onChange={e => setMcCurrency(e.target.value)}
                                        >
                                            {currencies.map(c => (
                                                <option key={c} value={c}>{c}</option>
                                            ))}
                                        </select>
                                    </div>
                                    <div>
                                        <label className="block text-sm font-medium text-gray-700 mb-1">Montant / Utilisateur</label>
                                        <input
                                            type="number"
                                            className="input"
                                            placeholder="500"
                                            value={mcAmount}
                                            onChange={e => setMcAmount(e.target.value)}
                                        />
                                    </div>
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Type de Raison</label>
                                    <select
                                        className="select"
                                        value={mcReasonType}
                                        onChange={e => setMcReasonType(e.target.value)}
                                    >
                                        {reasonTypes.map(rt => (
                                            <option key={rt.id} value={rt.id}>{rt.icon} {rt.label}</option>
                                        ))}
                                    </select>
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Raison / Description</label>
                                    <textarea
                                        className="input min-h-[80px]"
                                        placeholder="Description détaillée de la raison..."
                                        value={mcReason}
                                        onChange={e => setMcReason(e.target.value)}
                                    />
                                </div>

                                <button
                                    type="button"
                                    onClick={handleMassPreview}
                                    disabled={loading || !mcCampaignName || !mcAmount}
                                    className="w-full btn-secondary py-2.5 flex items-center justify-center gap-2"
                                >
                                    {loading && <div className="spinner w-5 h-5" />}
                                    <MagnifyingGlassIcon className="w-5 h-5" />
                                    Prévisualiser
                                </button>


                                {/* Preview Result */}
                                {mcPreview && (
                                    <div className={clsx(
                                        'p-4 rounded-xl border',
                                        mcPreview.sufficient_funds ? 'bg-emerald-50 border-emerald-200' : 'bg-red-50 border-red-200'
                                    )}>
                                        <h4 className="font-bold text-gray-800 mb-3">Résultat Prévisualisation</h4>
                                        <div className="grid grid-cols-2 gap-3 text-sm">
                                            <div>
                                                <span className="text-gray-500">Utilisateurs ciblés:</span>
                                                <span className="font-bold ml-2">{mcPreview.user_count}</span>
                                            </div>
                                            <div>
                                                <span className="text-gray-500">Montant/user:</span>
                                                <span className="font-bold ml-2">{mcPreview.amount_per_user} {mcPreview.currency}</span>
                                            </div>
                                            <div>
                                                <span className="text-gray-500">Total à distribuer:</span>
                                                <span className="font-bold ml-2 text-emerald-700">{mcPreview.total_amount.toLocaleString()} {mcPreview.currency}</span>
                                            </div>
                                            <div>
                                                <span className="text-gray-500">Solde Hot Wallet:</span>
                                                <span className={clsx('font-bold ml-2', mcPreview.sufficient_funds ? 'text-emerald-700' : 'text-red-600')}>
                                                    {mcPreview.hot_wallet_balance.toLocaleString()} {mcPreview.currency}
                                                </span>
                                            </div>
                                        </div>

                                        {!mcPreview.sufficient_funds && (
                                            <div className="mt-3 p-3 bg-red-100 rounded-lg flex items-center gap-2 text-red-700">
                                                <ExclamationTriangleIcon className="w-5 h-5" />
                                                <span className="text-sm font-medium">Fonds insuffisants dans le hot wallet!</span>
                                            </div>
                                        )}

                                        {mcPreview.sufficient_funds && (
                                            <button
                                                type="button"
                                                onClick={handleMassCredit}
                                                disabled={loading}
                                                className="w-full mt-4 btn-primary py-2.5 flex items-center justify-center gap-2 bg-gradient-to-r from-emerald-500 to-teal-600"
                                            >
                                                {loading && <div className="spinner w-5 h-5 border-white border-t-transparent" />}
                                                Exécuter le Crédit de Masse
                                            </button>
                                        )}
                                    </div>
                                )}
                            </div>
                        </div>
                    </div>
                )}

                {/* PROMOTION TAB */}
                {activeTab === 'promotion' && (
                    <div className="max-w-xl mx-auto space-y-6">
                        <div className="text-center mb-8">
                            <TrophyIcon className="w-12 h-12 text-amber-500 mx-auto mb-3" />
                            <h2 className="text-xl font-bold text-gray-900">Promotions & Concours</h2>
                            <p className="text-gray-500 text-sm">Distribuez des prix à des gagnants spécifiques.</p>
                        </div>

                        <form onSubmit={handlePromotionCredit} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Nom de la Campagne</label>
                                <input
                                    type="text"
                                    required
                                    className="input"
                                    placeholder="Ex: Concours Saint-Valentin 2026"
                                    value={prCampaignName}
                                    onChange={e => setPrCampaignName(e.target.value)}
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">IDs des Gagnants (un par ligne ou séparés par virgule)</label>
                                <textarea
                                    required
                                    className="input min-h-[120px]"
                                    placeholder="uuid-1-xxxx&#10;uuid-2-xxxx&#10;uuid-3-xxxx"
                                    value={prUserIds}
                                    onChange={e => setPrUserIds(e.target.value)}
                                />
                                <p className="text-xs text-gray-500 mt-1">
                                    {prUserIds.split(/[\n,]/).filter(id => id.trim()).length} utilisateur(s) sélectionné(s)
                                </p>
                            </div>

                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Devise</label>
                                    <select
                                        className="select"
                                        value={prCurrency}
                                        onChange={e => setPrCurrency(e.target.value)}
                                    >
                                        {currencies.map(c => (
                                            <option key={c} value={c}>{c}</option>
                                        ))}
                                    </select>
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Montant / Gagnant</label>
                                    <input
                                        type="number"
                                        required
                                        className="input"
                                        placeholder="10000"
                                        min="1"
                                        value={prUniformAmount}
                                        onChange={e => setPrUniformAmount(e.target.value)}
                                    />
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Type de Raison</label>
                                <select
                                    className="select"
                                    value={prReasonType}
                                    onChange={e => setPrReasonType(e.target.value)}
                                >
                                    {reasonTypes.map(rt => (
                                        <option key={rt.id} value={rt.id}>{rt.icon} {rt.label}</option>
                                    ))}
                                </select>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Message / Description</label>
                                <input
                                    type="text"
                                    required
                                    className="input"
                                    placeholder="Ex: Félicitations! Vous avez gagné au concours..."
                                    value={prReason}
                                    onChange={e => setPrReason(e.target.value)}
                                />
                            </div>

                            <button
                                type="submit"
                                disabled={loading}
                                className="w-full btn-primary py-3 flex items-center justify-center gap-2 mt-4 bg-gradient-to-r from-amber-500 to-orange-500"
                            >
                                {loading && <div className="spinner w-5 h-5 border-white border-t-transparent" />}
                                <GiftIcon className="w-5 h-5" />
                                Distribuer les Prix
                            </button>
                        </form>
                    </div>
                )}

                {/* CAMPAIGNS TAB */}
                {activeTab === 'campaigns' && (
                    <div className="overflow-x-auto">
                        <div className="flex justify-between items-center mb-4 px-2">
                            <h3 className="font-bold text-gray-700">Historique des Campagnes</h3>
                            <button onClick={fetchCampaigns} className="p-2 hover:bg-gray-100 rounded-full text-gray-500">
                                <ArrowPathIcon className="w-5 h-5" />
                            </button>
                        </div>

                        <table className="w-full text-left text-sm">
                            <thead className="bg-gray-50 text-gray-500 uppercase font-medium">
                                <tr>
                                    <th className="px-4 py-3">Campagne</th>
                                    <th className="px-4 py-3">Type</th>
                                    <th className="px-4 py-3">Statut</th>
                                    <th className="px-4 py-3 text-right">Montant</th>
                                    <th className="px-4 py-3 text-center">Utilisateurs</th>
                                    <th className="px-4 py-3">Date</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-100">
                                {campaigns.length === 0 ? (
                                    <tr>
                                        <td colSpan={6} className="px-4 py-8 text-center text-gray-500">
                                            Aucune campagne disponible
                                        </td>
                                    </tr>
                                ) : (
                                    campaigns.map(camp => (
                                        <tr key={camp.id} className="hover:bg-gray-50/50">
                                            <td className="px-4 py-3">
                                                <div className="font-medium text-gray-900">{camp.name}</div>
                                                <div className="text-xs text-gray-500">{camp.reason.substring(0, 50)}...</div>
                                            </td>
                                            <td className="px-4 py-3">
                                                {camp.type === 'single' ? (
                                                    <span className="badge badge-gray">Individuel</span>
                                                ) : camp.type === 'mass' ? (
                                                    <span className="badge badge-indigo">Masse</span>
                                                ) : (
                                                    <span className="badge badge-amber">Promotion</span>
                                                )}
                                            </td>
                                            <td className="px-4 py-3">{getStatusBadge(camp.status)}</td>
                                            <td className="px-4 py-3 text-right font-medium text-emerald-700">
                                                {camp.total_amount.toLocaleString()} {camp.currency}
                                            </td>
                                            <td className="px-4 py-3 text-center">
                                                <span className="text-emerald-600">{camp.success_count}</span>
                                                {camp.failed_count > 0 && (
                                                    <span className="text-red-500"> / {camp.failed_count}</span>
                                                )}
                                                <span className="text-gray-400"> / {camp.user_count}</span>
                                            </td>
                                            <td className="px-4 py-3 text-gray-500">
                                                {new Date(camp.created_at).toLocaleDateString()}
                                            </td>
                                        </tr>
                                    ))
                                )}
                            </tbody>
                        </table>
                    </div>
                )}

                {/* LOGS TAB */}
                {activeTab === 'logs' && (
                    <div className="overflow-x-auto">
                        <div className="flex justify-between items-center mb-4 px-2">
                            <h3 className="font-bold text-gray-700">Logs d'Audit des Opérations</h3>
                            <button onClick={fetchLogs} className="p-2 hover:bg-gray-100 rounded-full text-gray-500">
                                <ArrowPathIcon className="w-5 h-5" />
                            </button>
                        </div>

                        <table className="w-full text-left text-sm">
                            <thead className="bg-gray-50 text-gray-500 uppercase font-medium">
                                <tr>
                                    <th className="px-4 py-3">Date</th>
                                    <th className="px-4 py-3">Type</th>
                                    <th className="px-4 py-3">Utilisateur</th>
                                    <th className="px-4 py-3 text-right">Montant</th>
                                    <th className="px-4 py-3">Campagne</th>
                                    <th className="px-4 py-3">Statut</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-100">
                                {logs.length === 0 ? (
                                    <tr>
                                        <td colSpan={6} className="px-4 py-8 text-center text-gray-500">
                                            Aucun log disponible
                                        </td>
                                    </tr>
                                ) : (
                                    logs.map(log => (
                                        <tr key={log.id} className="hover:bg-gray-50/50">
                                            <td className="px-4 py-3 text-gray-500">
                                                {new Date(log.created_at).toLocaleString()}
                                            </td>
                                            <td className="px-4 py-3">
                                                <span className="text-lg mr-2">{getReasonTypeIcon(log.reason_type)}</span>
                                            </td>
                                            <td className="px-4 py-3 font-mono text-xs text-gray-600">
                                                {log.user_id.substring(0, 8)}...
                                            </td>
                                            <td className="px-4 py-3 text-right font-medium text-gray-900">
                                                {log.amount.toLocaleString()} {log.currency}
                                            </td>
                                            <td className="px-4 py-3 text-gray-600 truncate max-w-[200px]">
                                                {log.campaign_name}
                                            </td>
                                            <td className="px-4 py-3">
                                                {log.status === 'success' ? (
                                                    <span className="badge badge-success">Réussi</span>
                                                ) : (
                                                    <span className="badge badge-danger" title={log.error_message}>Échoué</span>
                                                )}
                                            </td>
                                        </tr>
                                    ))
                                )}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </div>
    );
}
