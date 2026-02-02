'use client';

import { useState, useEffect } from 'react';
import {
    SparklesIcon,
    BoltIcon,
    CreditCardIcon,
    CodeBracketIcon,
    ClockIcon,
    UserIcon,
    BeakerIcon,
    CheckCircleIcon,
    XCircleIcon,
    ArrowPathIcon,
} from '@heroicons/react/24/outline';
import clsx from 'clsx';
import { useAuthStore } from '@/stores/authStore';

// Types
interface CreditPreset {
    small: number;
    medium: number;
    large: number;
}

interface CreditLog {
    id: string;
    admin_id: string;
    user_id: string;
    wallet_id: string;
    currency: string;
    amount: number;
    reason: string;
    aggregator: string;
    is_test_mode: boolean;
    status: string;
    created_at: string;
}

export default function TestModePage() {
    const [activeTab, setActiveTab] = useState<'manual' | 'quick' | 'webhook' | 'logs'>('quick');
    const [loading, setLoading] = useState(false);
    const [successMsg, setSuccessMsg] = useState('');
    const [errorMsg, setErrorMsg] = useState('');
    const [logs, setLogs] = useState<CreditLog[]>([]);

    // Quick Credit State
    const [qcUserId, setQcUserId] = useState('');
    const [qcCurrency, setQcCurrency] = useState('USD');
    const [qcPreset, setQcPreset] = useState<'small' | 'medium' | 'large'>('small');
    const [presets, setPresets] = useState<Record<string, CreditPreset>>({});

    // Manual Credit State
    const [mcUserId, setMcUserId] = useState('');
    const [mcCurrency, setMcCurrency] = useState('USD');
    const [mcAmount, setMcAmount] = useState('');
    const [mcReason, setMcReason] = useState('Test deposit');
    const [mcTestMode, setMcTestMode] = useState(true);

    // Webhook Sim State
    const [whUserId, setWhUserId] = useState('');
    const [whCurrency, setWhCurrency] = useState('USD');
    const [whAmount, setWhAmount] = useState('');
    const [whAggregator, setWhAggregator] = useState('stripe');
    const [whEvent, setWhEvent] = useState('payment.success');

    useEffect(() => {
        fetchPresets();
        fetchLogs();
    }, []);

    const fetchPresets = async () => {
        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');
            const res = await fetch(`${API_URL}/api/v1/admin/test-mode/presets`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            if (res.ok) {
                const data = await res.json();
                setPresets(data.presets);
            }
        } catch (err) {
            console.error(err);
        }
    };

    const fetchLogs = async () => {
        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');
            const res = await fetch(`${API_URL}/api/v1/admin/test-mode/logs`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            if (res.ok) {
                const data = await res.json();
                setLogs(data.logs || []);
            }
        } catch (err) {
            console.error(err);
        }
    };

    const handleQuickCredit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        resetMessages();

        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');

            const res = await fetch(`${API_URL}/api/v1/admin/test-mode/quick-credit`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    user_id: qcUserId,
                    currency: qcCurrency,
                    preset: qcPreset
                }),
            });

            const data = await res.json();
            if (res.ok) {
                setSuccessMsg(`Crédit réussi! Nouveau solde: ${data.new_balance} ${data.currency}`);
                fetchLogs();
            } else {
                setErrorMsg(data.error || 'Erreur lors du crédit');
            }
        } catch (err) {
            setErrorMsg('Erreur de connexion');
        } finally {
            setLoading(false);
        }
    };

    const handleManualCredit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        resetMessages();

        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');

            const res = await fetch(`${API_URL}/api/v1/admin/test-mode/credit`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    user_id: mcUserId,
                    currency: mcCurrency,
                    amount: parseFloat(mcAmount),
                    reason: mcReason,
                    test_mode: mcTestMode
                }),
            });

            const data = await res.json();
            if (res.ok) {
                setSuccessMsg(`Opération réussie! ${mcTestMode ? '(Mode Test: Solde inchangé)' : `Nouveau solde: ${data.new_balance} ${data.currency}`}`);
                fetchLogs();
            } else {
                setErrorMsg(data.error || 'Erreur lors de l\'opération');
            }
        } catch (err) {
            setErrorMsg('Erreur de connexion');
        } finally {
            setLoading(false);
        }
    };

    const handleWebhookSim = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        resetMessages();

        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');

            const res = await fetch(`${API_URL}/api/v1/admin/test-mode/simulate-webhook`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    user_id: whUserId,
                    currency: whCurrency,
                    amount: parseFloat(whAmount),
                    aggregator: whAggregator,
                    event_type: whEvent
                }),
            });

            const data = await res.json();
            if (res.ok) {
                setSuccessMsg('Webhook simulé avec succès!');
                fetchLogs(); // Webhook sims might not log to same table depending on backend implementation, but usually they trigger credits which log
            } else {
                setErrorMsg(data.error || 'Erreur simulation webhook');
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

    const currencies = ['USD', 'EUR', 'GBP', 'NGN', 'XOF', 'XAF', 'KES', 'GHS', 'BTC', 'ETH', 'USDT'];
    const aggregators = ['stripe', 'paypal', 'flutterwave', 'mtn_momo', 'orange_money', 'pesapal', 'chipper'];

    return (
        <div className="space-y-6 animate-fadeIn">
            {/* Header */}
            <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
                        <SparklesIcon className="w-8 h-8 text-amber-500" />
                        Mode Test & Développement
                    </h1>
                    <p className="text-gray-500 mt-1">Outils pour tester les dépôts, simuler des webhooks et gérer les soldes.</p>
                </div>
            </div>

            {/* Tabs */}
            <div className="flex flex-wrap gap-2 border-b border-gray-200 pb-1">
                <button
                    onClick={() => setActiveTab('quick')}
                    className={clsx(
                        'flex items-center gap-2 px-4 py-2 rounded-t-lg text-sm font-medium transition-colors',
                        activeTab === 'quick'
                            ? 'bg-indigo-50 text-indigo-700 border-b-2 border-indigo-600'
                            : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'
                    )}
                >
                    <BoltIcon className="w-4 h-4" />
                    Crédit Rapide
                </button>
                <button
                    onClick={() => setActiveTab('manual')}
                    className={clsx(
                        'flex items-center gap-2 px-4 py-2 rounded-t-lg text-sm font-medium transition-colors',
                        activeTab === 'manual'
                            ? 'bg-indigo-50 text-indigo-700 border-b-2 border-indigo-600'
                            : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'
                    )}
                >
                    <CreditCardIcon className="w-4 h-4" />
                    Crédit Manuel
                </button>
                <button
                    onClick={() => setActiveTab('webhook')}
                    className={clsx(
                        'flex items-center gap-2 px-4 py-2 rounded-t-lg text-sm font-medium transition-colors',
                        activeTab === 'webhook'
                            ? 'bg-indigo-50 text-indigo-700 border-b-2 border-indigo-600'
                            : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'
                    )}
                >
                    <CodeBracketIcon className="w-4 h-4" />
                    Simulateur Webhook
                </button>
                <button
                    onClick={() => setActiveTab('logs')}
                    className={clsx(
                        'flex items-center gap-2 px-4 py-2 rounded-t-lg text-sm font-medium transition-colors',
                        activeTab === 'logs'
                            ? 'bg-indigo-50 text-indigo-700 border-b-2 border-indigo-600'
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
            <div className="card min-h-[400px]">

                {/* QUICK CREDIT TAB */}
                {activeTab === 'quick' && (
                    <div className="max-w-xl mx-auto space-y-6">
                        <div className="text-center mb-8">
                            <BoltIcon className="w-12 h-12 text-indigo-500 mx-auto mb-3" />
                            <h2 className="text-xl font-bold text-gray-900">Crédit Rapide</h2>
                            <p className="text-gray-500 text-sm">Ajoutez des fonds instantanément à un utilisateur pour tester.</p>
                        </div>

                        <form onSubmit={handleQuickCredit} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">ID Utilisateur</label>
                                <div className="relative">
                                    <UserIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                                    <input
                                        type="text"
                                        required
                                        className="input pl-10"
                                        placeholder="UUID de l'utilisateur"
                                        value={qcUserId}
                                        onChange={e => setQcUserId(e.target.value)}
                                    />
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Devise</label>
                                <div className="grid grid-cols-4 gap-2">
                                    {currencies.slice(0, 8).map(curr => (
                                        <button
                                            key={curr}
                                            type="button"
                                            onClick={() => setQcCurrency(curr)}
                                            className={clsx(
                                                'px-3 py-2 text-sm rounded-lg border transition-all',
                                                qcCurrency === curr
                                                    ? 'bg-indigo-600 text-white border-indigo-600 shadow-lg shadow-indigo-200'
                                                    : 'bg-white text-gray-700 border-gray-200 hover:bg-gray-50'
                                            )}
                                        >
                                            {curr}
                                        </button>
                                    ))}
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Montant (Preset)</label>
                                <div className="grid grid-cols-3 gap-3">
                                    {(['small', 'medium', 'large'] as const).map(preset => {
                                        const amount = presets[qcCurrency]?.[preset] || 0;
                                        return (
                                            <button
                                                key={preset}
                                                type="button"
                                                onClick={() => setQcPreset(preset)}
                                                className={clsx(
                                                    'p-4 rounded-xl border-2 text-center transition-all',
                                                    qcPreset === preset
                                                        ? 'border-indigo-500 bg-indigo-50'
                                                        : 'border-gray-200 hover:border-indigo-200'
                                                )}
                                            >
                                                <div className="text-xs text-gray-500 uppercase font-bold tracking-wider mb-1">{preset}</div>
                                                <div className="text-lg font-bold text-indigo-700">
                                                    {amount} {qcCurrency}
                                                </div>
                                            </button>
                                        );
                                    })}
                                </div>
                            </div>

                            <button
                                type="submit"
                                disabled={loading}
                                className="w-full btn-primary py-3 flex items-center justify-center gap-2 mt-4"
                            >
                                {loading && <div className="spinner w-5 h-5 border-white border-t-transparent" />}
                                Créditer maintenant
                            </button>
                        </form>
                    </div>
                )}

                {/* MANUAL CREDIT TAB */}
                {activeTab === 'manual' && (
                    <div className="max-w-xl mx-auto space-y-6">
                        <div className="text-center mb-8">
                            <CreditCardIcon className="w-12 h-12 text-emerald-500 mx-auto mb-3" />
                            <h2 className="text-xl font-bold text-gray-900">Crédit Manuel Avancé</h2>
                            <p className="text-gray-500 text-sm">Crédit précis avec contrôle du mode test/réel.</p>
                        </div>

                        <form onSubmit={handleManualCredit} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">ID Utilisateur</label>
                                <input
                                    type="text"
                                    required
                                    className="input"
                                    placeholder="UUID de l'utilisateur"
                                    value={mcUserId}
                                    onChange={e => setMcUserId(e.target.value)}
                                />
                            </div>

                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Devise</label>
                                    <select
                                        className="select"
                                        value={mcCurrency}
                                        onChange={e => setMcCurrency(e.target.value)}
                                    >
                                        {currencies.map(c => <option key={c} value={c}>{c}</option>)}
                                    </select>
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Montant</label>
                                    <input
                                        type="number"
                                        required
                                        className="input"
                                        placeholder="0.00"
                                        min="0.00000001"
                                        step="any"
                                        value={mcAmount}
                                        onChange={e => setMcAmount(e.target.value)}
                                    />
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Raison</label>
                                <input
                                    type="text"
                                    required
                                    className="input"
                                    placeholder="Ex: Correction solde"
                                    value={mcReason}
                                    onChange={e => setMcReason(e.target.value)}
                                />
                            </div>

                            <div className="flex items-center gap-3 p-4 bg-gray-50 rounded-xl border border-gray-200">
                                <button
                                    type="button"
                                    onClick={() => setMcTestMode(!mcTestMode)}
                                    className={clsx(
                                        'relative inline-flex h-6 w-11 items-center rounded-full transition-colors',
                                        mcTestMode ? 'bg-amber-400' : 'bg-emerald-500'
                                    )}
                                >
                                    <span className={clsx(
                                        'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
                                        mcTestMode ? 'translate-x-6' : 'translate-x-1'
                                    )} />
                                </button>
                                <div>
                                    <span className="font-medium text-gray-900">
                                        {mcTestMode ? 'Mode Test (Simulation)' : 'Mode Réel (Transaction)'}
                                    </span>
                                    <p className="text-xs text-gray-500">
                                        {mcTestMode
                                            ? 'Le solde de l\'utilisateur ne sera PAS modifié.'
                                            : '⚠️ ATTENTION: Le solde de l\'utilisateur SERA crédité.'}
                                    </p>
                                </div>
                            </div>

                            <button
                                type="submit"
                                disabled={loading}
                                className={clsx(
                                    'w-full py-3 flex items-center justify-center gap-2 mt-4 rounded-xl font-bold text-white transition-all shadow-lg',
                                    mcTestMode
                                        ? 'bg-gradient-to-r from-amber-400 to-orange-500 shadow-amber-200 hover:shadow-amber-300'
                                        : 'bg-gradient-to-r from-emerald-500 to-teal-600 shadow-emerald-200 hover:shadow-emerald-300'
                                )}
                            >
                                {loading && <div className="spinner w-5 h-5 border-white border-t-transparent" />}
                                {mcTestMode ? 'Simuler Crédit' : 'Confirmer le Crédit'}
                            </button>
                        </form>
                    </div>
                )}

                {/* WEBHOOK SIMULATOR TAB */}
                {activeTab === 'webhook' && (
                    <div className="max-w-xl mx-auto space-y-6">
                        <div className="text-center mb-8">
                            <CodeBracketIcon className="w-12 h-12 text-purple-500 mx-auto mb-3" />
                            <h2 className="text-xl font-bold text-gray-900">Simulateur Webhook</h2>
                            <p className="text-gray-500 text-sm">Déclenchez les handlers de paiement comme si le fournisseur appelait l'API.</p>
                        </div>

                        <form onSubmit={handleWebhookSim} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">ID Utilisateur</label>
                                <input
                                    type="text"
                                    required
                                    className="input"
                                    placeholder="UUID de l'utilisateur"
                                    value={whUserId}
                                    onChange={e => setWhUserId(e.target.value)}
                                />
                            </div>

                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Agrégateur</label>
                                    <select
                                        className="select"
                                        value={whAggregator}
                                        onChange={e => setWhAggregator(e.target.value)}
                                    >
                                        {aggregators.map(a => <option key={a} value={a}>{a}</option>)}
                                    </select>
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Événement</label>
                                    <select
                                        className="select"
                                        value={whEvent}
                                        onChange={e => setWhEvent(e.target.value)}
                                    >
                                        <option value="payment.success">Paiement Réussi</option>
                                        <option value="payment.failed">Paiement Échoué</option>
                                        <option value="refund.processed">Remboursement</option>
                                    </select>
                                </div>
                            </div>

                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Devise</label>
                                    <select
                                        className="select"
                                        value={whCurrency}
                                        onChange={e => setWhCurrency(e.target.value)}
                                    >
                                        {currencies.map(c => <option key={c} value={c}>{c}</option>)}
                                    </select>
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Montant</label>
                                    <input
                                        type="number"
                                        required
                                        className="input"
                                        placeholder="0.00"
                                        min="0.00000001"
                                        step="any"
                                        value={whAmount}
                                        onChange={e => setWhAmount(e.target.value)}
                                    />
                                </div>
                            </div>

                            <button
                                type="submit"
                                disabled={loading}
                                className="w-full btn-primary py-3 flex items-center justify-center gap-2 mt-4 bg-gradient-to-r from-purple-600 to-indigo-600"
                            >
                                {loading && <div className="spinner w-5 h-5 border-white border-t-transparent" />}
                                Envoyer Webhook
                            </button>
                        </form>
                    </div>
                )}

                {/* LOGS TAB */}
                {activeTab === 'logs' && (
                    <div className="overflow-x-auto">
                        <div className="flex justify-between items-center mb-4 px-2">
                            <h3 className="font-bold text-gray-700">Historique des opérations</h3>
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
                                    <th className="px-4 py-3">Raison</th>
                                    <th className="px-4 py-3">Mode</th>
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
                                                {log.reason.includes('Webhook') ? (
                                                    <span className="badge badge-purple">Webhook</span>
                                                ) : log.reason.includes('Quick') ? (
                                                    <span className="badge badge-info">Rapide</span>
                                                ) : (
                                                    <span className="badge badge-gray">Manuel</span>
                                                )}
                                            </td>
                                            <td className="px-4 py-3 font-mono text-xs text-gray-600">
                                                {log.user_id.substring(0, 8)}...
                                            </td>
                                            <td className="px-4 py-3 text-right font-medium text-gray-900">
                                                {log.amount} {log.currency}
                                            </td>
                                            <td className="px-4 py-3 text-gray-600 truncate max-w-[200px]">
                                                {log.reason}
                                            </td>
                                            <td className="px-4 py-3">
                                                {log.is_test_mode ? (
                                                    <span className="badge badge-warning">Simulé</span>
                                                ) : (
                                                    <span className="badge badge-success">Réel</span>
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
