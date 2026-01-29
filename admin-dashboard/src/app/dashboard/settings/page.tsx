'use client';

import { useState, useEffect } from 'react';
import {
    AdjustmentsHorizontalIcon,
    ArrowPathIcon,
    PencilSquareIcon,
    XMarkIcon,
    CurrencyDollarIcon,
    ArrowsRightLeftIcon,
    BanknotesIcon,
    GlobeAltIcon,
    ShieldCheckIcon,
    CreditCardIcon,
    WalletIcon,
    ChartBarIcon,
    SwatchIcon,
    Cog6ToothIcon,
    UserIcon,
    BuildingOfficeIcon,
    ShoppingBagIcon
} from '@heroicons/react/24/outline';
import { getFeeConfigs, updateFeeConfig, createFeeConfig } from '@/lib/api';

interface FeeConfig {
    id: string;
    key: string;
    value?: string; // Add value field for system settings
    name: string;
    description: string;
    type: string; // flat, percentage, hybrid
    fixed_amount: number;
    percentage_amount: number;
    currency: string;
    is_enabled: boolean;
    updated_at: string;
}

// ----------------------------------------------------------------------
// CONSTANTS & HELPERS
// ----------------------------------------------------------------------

const CATEGORIES = {
    CRYPTO: 'Crypto',
    TRANSFER: 'Transferts',
    FIAT: 'Fiat',
    CARD: 'Cartes',
    SYSTEM: 'Système',
    LIMITS: 'Limites',
    SERVICES: 'Services'
};

const getCategory = (key: string): string => {
    if (key.includes('crypto')) return CATEGORIES.CRYPTO;
    if (key.includes('transfer') || key.includes('international')) return CATEGORIES.TRANSFER;
    if (key.includes('card')) return CATEGORIES.CARD;
    if (key.includes('fiat') || key.includes('deposit') || key.includes('withdrawal')) return CATEGORIES.FIAT;
    if (key.includes('donation') || key.includes('ecommerce') || key.includes('event') || key.includes('bill') || key.includes('mobile') || key.includes('association') || key.includes('crowdfunding')) return CATEGORIES.SERVICES;
    return CATEGORIES.SYSTEM;
};

// Mapping des clés techniques vers des noms conviviaux
const FRIENDLY_NAMES: Record<string, string> = {
    'transfer_international': 'Virement International',
    'crypto_withdrawal_btc': 'Retrait Bitcoin (BTC)',
    'crypto_withdrawal_eth': 'Retrait Ethereum (ETH)',
    'crypto_exchange_fee': 'Commission de Change',
    'fiat_deposit_fee': 'Dépôt Espèces/Virement',
    'card_issuance_fee': 'Création de Carte',
    'card_monthly_fee': 'Mensualité Carte',
    'transfer_p2p_fee': 'Transfert P2P',
    'sms_notification_fee': 'Alerte SMS',
    'donation_fee': 'Dons',
    'ecommerce_purchase_fee': 'Achats E-commerce',
    'event_ticket_fee': 'Billetterie Événements',
    'bill_payment_fee': 'Paiement de Factures',
    'mobile_money_cashin_fee': 'Dépôt Mobile Money',
    'mobile_money_cashout_fee': 'Retrait Mobile Money',
    'association_membership_fee': 'Adhésion Association',
    'crowdfunding_contribution_fee': 'Contribution Participative'
};

const getFriendlyName = (key: string, defaultName: string) => {
    return FRIENDLY_NAMES[key] || defaultName;
};

const getIcon = (category: string) => {
    switch (category) {
        case CATEGORIES.CRYPTO: return <CurrencyDollarIcon className="w-6 h-6 text-amber-500" />;
        case CATEGORIES.TRANSFER: return <ArrowsRightLeftIcon className="w-6 h-6 text-indigo-500" />;
        case CATEGORIES.FIAT: return <BanknotesIcon className="w-6 h-6 text-emerald-500" />;
        case CATEGORIES.CARD: return <CreditCardIcon className="w-6 h-6 text-purple-500" />;
        case CATEGORIES.SERVICES: return <ShoppingBagIcon className="w-6 h-6 text-pink-500" />;
        default: return <Cog6ToothIcon className="w-6 h-6 text-slate-500" />;
    }
};

// ----------------------------------------------------------------------
// COMPONENTS
// ----------------------------------------------------------------------

export default function FeeManagementPage() {
    // API State
    const [fees, setFees] = useState<FeeConfig[]>([]);
    const [loading, setLoading] = useState(true);

    // UI State
    const [activeTab, setActiveTab] = useState<string>('ALL');
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingFee, setEditingFee] = useState<FeeConfig | null>(null);
    const [saving, setSaving] = useState(false);

    // Form State
    const [formData, setFormData] = useState({
        key: '',
        name: '',
        description: '',
        type: 'percentage',
        fixed_amount: 0,
        percentage_amount: 0,
        currency: 'EUR',
        is_enabled: true
    });

    useEffect(() => {
        loadFees();
    }, []);

    const loadFees = async () => {
        setLoading(true);
        try {
            const res = await getFeeConfigs();
            if (res.data?.fees) setFees(res.data.fees);
            else setFees([]);
        } catch (error) {
            console.error('Failed to load fees', error);
        } finally {
            setLoading(false);
        }
    };

    const handleEdit = (fee: FeeConfig) => {
        setEditingFee(fee);
        setFormData({
            key: fee.key,
            name: fee.name,
            description: fee.description,
            type: fee.type,
            fixed_amount: fee.fixed_amount,
            percentage_amount: fee.percentage_amount,
            currency: fee.currency,
            is_enabled: fee.is_enabled
        });
        setIsModalOpen(true);
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSaving(true);
        try {
            if (editingFee) await updateFeeConfig(editingFee.key, formData);
            else await createFeeConfig(formData);

            await loadFees();
            setIsModalOpen(false);
        } catch (error) {
            console.error('Failed to save fee', error);
            alert('Erreur lors de la sauvegarde');
        } finally {
            setSaving(false);
        }
    };

    // Filter Logic
    const tabs = [
        { id: 'ALL', label: 'Vue d\'ensemble' },
        { id: CATEGORIES.CRYPTO, label: 'Crypto-monnaies' },
        { id: CATEGORIES.FIAT, label: 'Opérations Fiat' },
        { id: CATEGORIES.TRANSFER, label: 'Transferts' },
        { id: CATEGORIES.CARD, label: 'Cartes Bancaires' },
        { id: CATEGORIES.SERVICES, label: 'Services' },
        { id: CATEGORIES.SYSTEM, label: 'Système' },
    ];

    const filteredFees = activeTab === 'ALL'
        ? fees
        : fees.filter(f => getCategory(f.key) === activeTab);

    return (
        <div className="min-h-screen bg-slate-50/50 pb-20">
            {/* ... Header ... */}

            {/* ... Content ... */}
            <div className="max-w-7xl mx-auto px-6 py-8">
                {activeTab === 'ALL' && (
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-10">
                        {/* ... Summary Cards ... */}
                    </div>
                )}

                <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
                    {(activeTab === 'ALL' || activeTab === CATEGORIES.SYSTEM) && (
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 col-span-full mb-6">
                            <div className="col-span-full flex items-center gap-2 mb-2">
                                <div className="p-2 bg-indigo-50 rounded-lg text-indigo-600">
                                    <AdjustmentsHorizontalIcon className="w-5 h-5" />
                                </div>
                                <h3 className="font-bold text-slate-800 text-lg">Configurations Système</h3>
                            </div>
                            {/* Hardcoded System Toggles for UI Verification */}
                            {[
                                { key: 'system_maintenance_mode', label: 'Mode Maintenance', desc: 'Suspend l\'accès utilisateur pour maintenance.', color: 'bg-red-500' },
                                { key: 'crypto_network', label: 'Réseaux de Test (Testnet)', desc: 'Active BTC Testnet, Sepolia, etc. pour le développement.', color: 'bg-orange-500' },
                                { key: 'system_signup_enabled', label: 'Inscriptions', desc: 'Autoriser les nouveaux utilisateurs à s\'inscrire.', color: 'bg-emerald-500' },
                                { key: 'system_notifications', label: 'Notifications Globales', desc: 'Activer les emails et push notifications.', color: 'bg-blue-500' }
                            ].map(sys => {
                                const existing = fees.find(f => f.key === sys.key);

                                // Special logic for crypto_network (string value) vs boolean toggles
                                let isEnabled = false;
                                if (sys.key === 'crypto_network') {
                                    isEnabled = existing ? existing.value === 'testnet' : false;
                                } else {
                                    isEnabled = existing ? existing.is_enabled : false;
                                    // Note: if backend doesn't support is_enabled for system_settings, this fallbacks might need review for other keys
                                    // but specifically fixing crypto_network now.
                                }

                                return (
                                    <div key={sys.key} className="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition-all">
                                        <div className="flex justify-between items-start mb-4">
                                            <div className={`p-3 rounded-xl ${isEnabled ? 'bg-indigo-50 text-indigo-600' : 'bg-slate-50 text-slate-400'}`}>
                                                <AdjustmentsHorizontalIcon className="w-6 h-6" />
                                            </div>
                                            <div onClick={() => {
                                                let payload: any = {
                                                    key: sys.key,
                                                    name: sys.label,
                                                    description: sys.desc,
                                                    type: 'system_toggle',
                                                    currency: 'USD',
                                                    fixed_amount: 0,
                                                    percentage_amount: 0
                                                };

                                                if (sys.key === 'crypto_network') {
                                                    payload.value = isEnabled ? 'mainnet' : 'testnet'; // Toggle
                                                    // Ensure we keep existing other fields if needed, but 'value' is what matters
                                                } else {
                                                    payload.is_enabled = !isEnabled;
                                                }

                                                if (existing) {
                                                    handleEdit({ ...existing, ...payload } as FeeConfig);
                                                } else {
                                                    createFeeConfig(payload).then(loadFees);
                                                }
                                            }} className={`w-12 h-6 rounded-full cursor-pointer transition-colors relative ${isEnabled ? 'bg-indigo-600' : 'bg-slate-300'}`}>
                                                <div className={`absolute top-1 left-1 bg-white w-4 h-4 rounded-full transition-transform ${isEnabled ? 'translate-x-6' : ''}`} />
                                            </div>
                                        </div>
                                        <h3 className="font-bold text-slate-900 mb-1">{sys.label}</h3>
                                        <p className="text-sm text-slate-500">{sys.desc}</p>
                                    </div>
                                );
                            })}
                        </div>
                    )}

                    {activeTab === CATEGORIES.LIMITS && (
                        <div className="space-y-8">
                            {/* Limits configuration for different Tiers */}
                            {['PARTICULIER', 'ENTREPRISE'].map(type => (
                                <div key={type} className="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
                                    <div className="px-6 py-4 border-b border-slate-100 bg-slate-50 flex justify-between items-center">
                                        <h3 className="font-bold text-slate-800 flex items-center gap-2">
                                            {type === 'PARTICULIER' ? <UserIcon className="w-5 h-5 text-indigo-600" /> : <BuildingOfficeIcon className="w-5 h-5 text-emerald-600" />}
                                            Comptes {type === 'PARTICULIER' ? 'Particuliers' : 'Entreprises'}
                                        </h3>
                                        <span className="text-xs font-bold uppercase tracking-wider text-slate-400">Configuration des Plafonds</span>
                                    </div>
                                    <div className="p-6">
                                        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                                            {['Niveau Standard (KYC 1)', 'Niveau Vérifié (KYC 2)', 'Niveau Premium (KYC 3)'].map((level, idx) => (
                                                <div key={idx} className="p-5 rounded-xl border border-slate-200 hover:border-indigo-300 transition-all group">
                                                    <div className="flex items-center justify-between mb-4">
                                                        <span className="text-xs font-bold uppercase text-slate-500 bg-slate-100 px-2 py-1 rounded">{level}</span>
                                                        <PencilSquareIcon className="w-4 h-4 text-slate-300 group-hover:text-indigo-600 cursor-pointer" />
                                                    </div>
                                                    <div className="space-y-3">
                                                        <div>
                                                            <div className="text-[10px] text-slate-400 uppercase font-bold">Plafond Journalier</div>
                                                            <div className="font-mono font-bold text-slate-800 text-lg">
                                                                {type === 'ENTREPRISE' ? (25000 * (idx + 1)).toLocaleString() : (1000 * (idx + 1)).toLocaleString()} €
                                                            </div>
                                                        </div>
                                                        <div>
                                                            <div className="text-[10px] text-slate-400 uppercase font-bold">Plafond Mensuel</div>
                                                            <div className="font-mono font-bold text-slate-800">
                                                                {type === 'ENTREPRISE' ? (100000 * (idx + 1)).toLocaleString() : (5000 * (idx + 1)).toLocaleString()} €
                                                            </div>
                                                        </div>
                                                        <div className="pt-3 border-t border-slate-100 mt-3">
                                                            <div className="text-[10px] text-slate-400 uppercase font-bold">Limite Solde</div>
                                                            <div className="font-mono font-bold text-emerald-600">
                                                                {idx === 2 ? 'Illimité' : (type === 'ENTREPRISE' ? (500000 * (idx + 1)).toLocaleString() + ' €' : (10000 * (idx + 1)).toLocaleString() + ' €')}
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}

                    {filteredFees.map((fee) => {
                        const category = getCategory(fee.key);
                        if (category === CATEGORIES.SYSTEM) return null; // Handled above specifically
                        return (
                            <div key={fee.id} className="group bg-white rounded-2xl border border-slate-200 shadow-sm hover:shadow-xl hover:border-indigo-300 transition-all duration-300 flex flex-col overflow-hidden">
                                {/* Card Status Stripe */}
                                <div className={`h-1.5 w-full ${fee.is_enabled ? 'bg-gradient-to-r from-emerald-400 to-teal-500' : 'bg-slate-200'}`} />

                                <div className="p-6 flex-1 flex flex-col">
                                    <div className="flex justify-between items-start mb-6">
                                        <div className="flex gap-4">
                                            <div className="p-3 bg-slate-50 rounded-2xl group-hover:bg-indigo-50 group-hover:scale-110 transition-all duration-300">
                                                {getIcon(category)}
                                            </div>
                                            <div>
                                                <h3 className="font-bold text-slate-900 text-lg leading-tight group-hover:text-indigo-700 transition-colors">
                                                    {getFriendlyName(fee.key, fee.name)}
                                                </h3>
                                                <span className="inline-flex items-center gap-1.5 mt-1.5 px-2.5 py-0.5 rounded-md bg-slate-100 text-slate-500 text-[10px] font-bold uppercase tracking-wide">
                                                    {category}
                                                </span>
                                            </div>
                                        </div>
                                    </div>

                                    {/* Stats Grid */}
                                    <div className="grid grid-cols-2 gap-3 mb-6">
                                        <div className="bg-slate-50 rounded-xl p-3 border border-slate-100">
                                            <p className="text-[10px] font-bold text-slate-400 uppercase">Comm. Variable</p>
                                            <p className={`text-xl font-black ${fee.percentage_amount > 0 ? 'text-indigo-600' : 'text-slate-300'}`}>
                                                {fee.percentage_amount}%
                                            </p>
                                        </div>
                                        <div className="bg-slate-50 rounded-xl p-3 border border-slate-100">
                                            <p className="text-[10px] font-bold text-slate-400 uppercase">Frais Fixes</p>
                                            <p className={`text-xl font-black ${fee.fixed_amount > 0 ? 'text-slate-700' : 'text-slate-300'}`}>
                                                {fee.fixed_amount} <span className="text-xs font-bold text-slate-400">{fee.currency}</span>
                                            </p>
                                        </div>
                                    </div>

                                    <div className="mt-auto flex items-center justify-between pt-4 border-t border-slate-100">
                                        <div className="flex items-center gap-2">
                                            <div className={`w-2 h-2 rounded-full ${fee.is_enabled ? 'bg-emerald-500 animate-pulse' : 'bg-slate-300'}`}></div>
                                            <span className="text-xs font-semibold text-slate-500">{fee.is_enabled ? 'Actif' : 'Désactivé'}</span>
                                        </div>
                                        <button
                                            onClick={() => handleEdit(fee)}
                                            className="text-sm font-bold text-indigo-600 hover:text-indigo-800 hover:underline decoration-2 underline-offset-2 transition-all flex items-center gap-1"
                                        >
                                            Configurer
                                            <PencilSquareIcon className="w-4 h-4" />
                                        </button>
                                    </div>
                                </div>
                            </div>
                        );
                    })}
                </div>
            </div>

            {/* Premium Modal */}
            {isModalOpen && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm transition-all">
                    <div className="bg-white rounded-3xl w-full max-w-2xl shadow-2xl overflow-hidden animate-slide-up border border-slate-100">
                        {/* Modal Header */}
                        <div className="px-8 py-6 bg-slate-50 border-b border-slate-100 flex justify-between items-center">
                            <div>
                                <h3 className="font-extrabold text-2xl text-slate-900">
                                    {editingFee ? 'Modifier le Tarif' : 'Nouveau Tarif'}
                                </h3>
                                <p className="text-slate-500 text-sm font-medium mt-1">
                                    {editingFee ? `Configuration clé: ${editingFee.key}` : 'Créez une nouvelle règle de frais'}
                                </p>
                            </div>
                            <button onClick={() => setIsModalOpen(false)} className="p-2 rounded-full hover:bg-slate-200 text-slate-400 hover:text-slate-600 transition">
                                <XMarkIcon className="w-6 h-6" />
                            </button>
                        </div>

                        <form onSubmit={handleSubmit} className="p-8 space-y-8">
                            {/* Main Info */}
                            <div className="grid grid-cols-2 gap-6">
                                <div className="col-span-2">
                                    <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Nom Affiché (Public)</label>
                                    <input
                                        type="text"
                                        value={formData.name}
                                        onChange={e => setFormData({ ...formData, name: e.target.value })}
                                        className="w-full px-5 py-3 bg-white border-2 border-slate-100 rounded-xl text-base font-semibold text-slate-900 focus:outline-none focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 transition-all placeholder:text-slate-300"
                                        placeholder="Ex: Frais de transaction standard"
                                    />
                                </div>

                                {/* Key (ReadOnly if editing) */}
                                <div className="col-span-2">
                                    <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Identifiant Technique</label>
                                    <div className={`px-5 py-3 rounded-xl border-2 border-dashed border-slate-200 bg-slate-50 text-slate-500 font-mono text-sm ${!editingFee ? 'hidden' : ''}`}>
                                        {formData.key}
                                    </div>
                                    {!editingFee && (
                                        <input
                                            type="text"
                                            value={formData.key}
                                            onChange={e => setFormData({ ...formData, key: e.target.value })}
                                            className="w-full px-5 py-3 bg-white border-2 border-slate-100 rounded-xl text-sm font-mono text-slate-900 focus:outline-none focus:border-indigo-500 transition-all"
                                            placeholder="transfer_fee_example"
                                        />
                                    )}
                                </div>
                            </div>

                            {/* Pricing Engine */}
                            <div className="bg-slate-50 rounded-2xl p-6 border border-slate-100">
                                <h4 className="text-sm font-black text-slate-900 uppercase tracking-wide mb-6 flex items-center gap-2">
                                    <SwatchIcon className="w-5 h-5 text-indigo-500" />
                                    Moteur de Calcul
                                </h4>

                                <div className="grid grid-cols-2 gap-6">
                                    {/* Type Selector */}
                                    <div className="col-span-2 flex bg-white p-1.5 rounded-xl border border-slate-200 shadow-sm">
                                        {['percentage', 'flat', 'hybrid'].map(type => (
                                            <button
                                                key={type}
                                                type="button"
                                                onClick={() => setFormData({ ...formData, type })}
                                                className={`flex-1 py-2 rounded-lg text-sm font-bold capitalize transition-all ${formData.type === type
                                                    ? 'bg-indigo-600 text-white shadow-md'
                                                    : 'text-slate-500 hover:text-slate-700 hover:bg-slate-50'}`}
                                            >
                                                {type === 'percentage' ? 'Proportionnel (%)' : type === 'flat' ? 'Fixe' : 'Hybride'}
                                            </button>
                                        ))}
                                    </div>

                                    {/* Values */}
                                    {(formData.type !== 'flat') && (
                                        <div>
                                            <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Taux (%)</label>
                                            <div className="relative group">
                                                <input
                                                    type="number" step="0.01" min="0" max="100"
                                                    value={formData.percentage_amount}
                                                    onChange={e => setFormData({ ...formData, percentage_amount: parseFloat(e.target.value) })}
                                                    className="w-full px-5 py-3 bg-white border-2 border-slate-100 rounded-xl text-lg font-bold text-indigo-600 focus:outline-none focus:border-indigo-500 transition-all text-right pr-12"
                                                />
                                                <span className="absolute right-5 top-1/2 -translate-y-1/2 text-slate-400 font-bold">%</span>
                                            </div>
                                        </div>
                                    )}

                                    {(formData.type !== 'percentage') && (
                                        <div className="flex gap-4">
                                            <div className="flex-1">
                                                <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Montant</label>
                                                <input
                                                    type="number" step="0.01" min="0"
                                                    value={formData.fixed_amount}
                                                    onChange={e => setFormData({ ...formData, fixed_amount: parseFloat(e.target.value) })}
                                                    className="w-full px-5 py-3 bg-white border-2 border-slate-100 rounded-xl text-lg font-bold text-slate-900 focus:outline-none focus:border-indigo-500 transition-all text-right"
                                                />
                                            </div>
                                            <div className="w-24">
                                                <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Devise</label>
                                                <input
                                                    type="text"
                                                    value={formData.currency}
                                                    onChange={e => setFormData({ ...formData, currency: e.target.value })}
                                                    className="w-full px-5 py-3 bg-slate-100 rounded-xl text-lg font-bold text-slate-500 text-center border-2 border-transparent"
                                                    readOnly // Simplify: usually locked to system default or need dropdown. Keeping readonly-ish look for now
                                                />
                                            </div>
                                        </div>
                                    )}
                                </div>
                            </div>

                            {/* Activation Toggle */}
                            <div className="flex items-center justify-between p-5 bg-white rounded-2xl border-2 border-slate-100 hover:border-indigo-100 transition-colors cursor-pointer" onClick={() => setFormData({ ...formData, is_enabled: !formData.is_enabled })}>
                                <div className="flex items-center gap-4">
                                    <div className={`p-3 rounded-xl ${formData.is_enabled ? 'bg-emerald-100 text-emerald-600' : 'bg-slate-100 text-slate-400'}`}>
                                        <ShieldCheckIcon className="w-6 h-6" />
                                    </div>
                                    <div>
                                        <h4 className="font-bold text-slate-900">Activer cette règle</h4>
                                        <p className="text-sm text-slate-500">Rend ces frais applicables immédiatement</p>
                                    </div>
                                </div>
                                <div className={`w-14 h-8 rounded-full p-1 transition-colors duration-300 ${formData.is_enabled ? 'bg-emerald-500' : 'bg-slate-200'}`}>
                                    <div className={`w-6 h-6 bg-white rounded-full shadow-sm transform transition-transform duration-300 ${formData.is_enabled ? 'translate-x-6' : ''}`} />
                                </div>
                            </div>

                            {/* Footer Actions */}
                            <div className="flex justify-end gap-3 pt-4 border-t border-slate-100">
                                <button
                                    type="button"
                                    onClick={() => setIsModalOpen(false)}
                                    className="px-6 py-3 font-bold text-slate-500 hover:bg-slate-50 rounded-xl transition-all"
                                >
                                    Annuler
                                </button>
                                <button
                                    type="submit"
                                    disabled={saving}
                                    className="px-8 py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-bold rounded-xl shadow-lg shadow-indigo-500/20 transform hover:-translate-y-0.5 transition-all flex items-center gap-2"
                                >
                                    {saving && <ArrowPathIcon className="w-5 h-5 animate-spin" />}
                                    Sauvegarder
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}
