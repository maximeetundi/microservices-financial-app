'use client';

import { useState, useEffect } from 'react';
import {
    AdjustmentsHorizontalIcon,
    ArrowPathIcon,
    PlusIcon,
    PencilSquareIcon,
    CheckIcon,
    XMarkIcon
} from '@heroicons/react/24/outline';
import { getFeeConfigs, updateFeeConfig, createFeeConfig } from '@/lib/api';

interface FeeConfig {
    id: string;
    key: string;
    name: string;
    description: string;
    type: string; // flat, percentage, hybrid
    fixed_amount: number;
    percentage_amount: number;
    currency: string;
    is_enabled: boolean;
    updated_at: string;
}

export default function SettingsPage() {
    const [fees, setFees] = useState<FeeConfig[]>([]);
    const [loading, setLoading] = useState(true);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingFee, setEditingFee] = useState<FeeConfig | null>(null);
    const [saving, setSaving] = useState(false);

    // Form state
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
            if (res.data?.fees) {
                setFees(res.data.fees);
            } else {
                setFees([]);
            }
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

    const handleCreate = () => {
        setEditingFee(null);
        setFormData({
            key: '',
            name: '',
            description: '',
            type: 'percentage',
            fixed_amount: 0,
            percentage_amount: 0,
            currency: 'EUR',
            is_enabled: true
        });
        setIsModalOpen(true);
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSaving(true);
        try {
            if (editingFee) {
                await updateFeeConfig(editingFee.key, formData);
            } else {
                await createFeeConfig(formData);
            }
            await loadFees();
            setIsModalOpen(false);
        } catch (error) {
            console.error('Failed to save fee', error);
            alert('Failed to save settings');
        } finally {
            setSaving(false);
        }
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('fr-FR', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-2xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-indigo-600">
                        Global Settings
                    </h1>
                    <p className="text-slate-500 mt-1">Configure global platform fees and limits</p>
                </div>
                <div className="flex gap-3">
                    <button
                        onClick={loadFees}
                        className="p-2 rounded-lg bg-white border border-slate-200 text-slate-500 hover:text-slate-700 transition"
                    >
                        <ArrowPathIcon className={`w-5 h-5 ${loading ? 'animate-spin' : ''}`} />
                    </button>
                    {/* <button 
                        onClick={handleCreate} 
                        className="flex items-center gap-2 px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition font-medium"
                    >
                        <PlusIcon className="w-5 h-5" />
                        New Config
                    </button> */}
                </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {fees.map((fee) => (
                    <div key={fee.id} className="group relative bg-white p-6 rounded-2xl border border-slate-100 shadow-sm hover:shadow-md hover:border-indigo-100 transition-all duration-300">
                        <div className="flex justify-between items-start mb-4">
                            <div>
                                <h3 className="font-semibold text-slate-900 group-hover:text-indigo-600 transition-colors">{fee.name}</h3>
                                <code className="text-xs text-slate-400 bg-slate-50 px-1.5 py-0.5 rounded mt-1 block w-fit">{fee.key}</code>
                            </div>
                            <span className={`px-2.5 py-1 rounded-full text-xs font-semibold ${fee.is_enabled
                                    ? 'bg-emerald-50 text-emerald-600 border border-emerald-100'
                                    : 'bg-slate-100 text-slate-500 border border-slate-200'
                                }`}>
                                {fee.is_enabled ? 'Active' : 'Disabled'}
                            </span>
                        </div>

                        <p className="text-sm text-slate-500 mb-6 h-10 line-clamp-2">{fee.description}</p>

                        <div className="bg-slate-50 rounded-xl p-4 mb-4 border border-slate-100 group-hover:bg-indigo-50/30 group-hover:border-indigo-100 transition-colors">
                            <div className="flex justify-between items-center mb-2">
                                <span className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Type</span>
                                <span className="text-sm font-medium text-slate-700 capitalize">{fee.type}</span>
                            </div>
                            <div className="flex justify-between items-center">
                                <span className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Value</span>
                                <div className="text-right">
                                    {(fee.type === 'percentage' || fee.type === 'hybrid') && (
                                        <div className="text-sm font-bold text-slate-900">{fee.percentage_amount}%</div>
                                    )}
                                    {(fee.type === 'flat' || fee.type === 'hybrid') && (
                                        <div className="text-sm font-bold text-slate-900">
                                            {fee.type === 'hybrid' && <span className="text-slate-400 font-normal mr-1">+</span>}
                                            {fee.fixed_amount} {fee.currency}
                                        </div>
                                    )}
                                </div>
                            </div>
                        </div>

                        <div className="flex items-center justify-between pt-4 border-t border-slate-100">
                            <span className="text-xs text-slate-400">Updated: {formatDate(fee.updated_at)}</span>
                            <button
                                onClick={() => handleEdit(fee)}
                                className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-50 text-slate-600 hover:bg-indigo-50 hover:text-indigo-600 transition-colors text-xs font-semibold"
                            >
                                <PencilSquareIcon className="w-3.5 h-3.5" />
                                Edit
                            </button>
                        </div>
                    </div>
                ))}
            </div>

            {/* Modal */}
            {isModalOpen && (
                <div className="fixed inset-0 z-[60] flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-sm">
                    <div className="bg-white rounded-2xl w-full max-w-lg shadow-2xl overflow-hidden animate-slide-up">
                        <div className="px-6 py-4 border-b border-slate-100 flex justify-between items-center bg-slate-50/50">
                            <h3 className="font-bold text-lg text-slate-900">{editingFee ? 'Edit Configuration' : 'New Configuration'}</h3>
                            <button onClick={() => setIsModalOpen(false)} className="text-slate-400 hover:text-slate-600 transition">
                                <XMarkIcon className="w-5 h-5" />
                            </button>
                        </div>

                        <form onSubmit={handleSubmit} className="p-6 space-y-5">
                            <div>
                                <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-2">Config Key</label>
                                <input
                                    type="text"
                                    value={formData.key}
                                    onChange={e => setFormData({ ...formData, key: e.target.value })}
                                    disabled={!!editingFee}
                                    className="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-lg text-sm text-slate-900 font-mono disabled:opacity-60 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
                                    placeholder="e.g. transfer_fee_internal"
                                    required
                                />
                            </div>

                            <div className="grid grid-cols-2 gap-4">
                                <div className="col-span-2">
                                    <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-2">Display Name</label>
                                    <input
                                        type="text"
                                        value={formData.name}
                                        onChange={e => setFormData({ ...formData, name: e.target.value })}
                                        className="w-full px-4 py-2.5 bg-white border border-slate-200 rounded-lg text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
                                        required
                                    />
                                </div>

                                <div className="col-span-2">
                                    <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-2">Type</label>
                                    <div className="grid grid-cols-3 gap-2">
                                        {['flat', 'percentage', 'hybrid'].map(type => (
                                            <button
                                                key={type}
                                                type="button"
                                                onClick={() => setFormData({ ...formData, type })}
                                                className={`px-3 py-2 text-sm font-medium rounded-lg border capitalize transition-all ${formData.type === type
                                                        ? 'bg-indigo-50 border-indigo-200 text-indigo-700'
                                                        : 'bg-white border-slate-200 text-slate-600 hover:bg-slate-50'
                                                    }`}
                                            >
                                                {type}
                                            </button>
                                        ))}
                                    </div>
                                </div>

                                {(formData.type !== 'flat') && (
                                    <div>
                                        <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-2">Percentage</label>
                                        <div className="relative">
                                            <input
                                                type="number" step="0.01" min="0" max="100"
                                                value={formData.percentage_amount}
                                                onChange={e => setFormData({ ...formData, percentage_amount: parseFloat(e.target.value) })}
                                                className="w-full px-4 py-2.5 bg-white border border-slate-200 rounded-lg text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all pl-4 pr-8"
                                            />
                                            <span className="absolute right-3 top-2.5 text-slate-400 font-medium">%</span>
                                        </div>
                                    </div>
                                )}

                                {(formData.type !== 'percentage') && (
                                    <div className="flex gap-2">
                                        <div className="flex-1">
                                            <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-2">Fixed Amount</label>
                                            <input
                                                type="number" step="0.01" min="0"
                                                value={formData.fixed_amount}
                                                onChange={e => setFormData({ ...formData, fixed_amount: parseFloat(e.target.value) })}
                                                className="w-full px-4 py-2.5 bg-white border border-slate-200 rounded-lg text-sm text-slate-900 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
                                            />
                                        </div>
                                        <div className="w-24">
                                            <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-2">Currency</label>
                                            <input
                                                type="text"
                                                value={formData.currency}
                                                onChange={e => setFormData({ ...formData, currency: e.target.value })}
                                                className="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-lg text-sm text-center font-medium uppercase focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
                                            />
                                        </div>
                                    </div>
                                )}
                            </div>

                            <div className="flex items-center justify-between p-4 bg-slate-50 rounded-xl border border-slate-100">
                                <div>
                                    <div className="text-sm font-semibold text-slate-900">Active Status</div>
                                    <div className="text-xs text-slate-500">Enable or disable this fee rule</div>
                                </div>
                                <button
                                    type="button"
                                    onClick={() => setFormData({ ...formData, is_enabled: !formData.is_enabled })}
                                    className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${formData.is_enabled ? 'bg-emerald-500' : 'bg-slate-300'}`}
                                >
                                    <span className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${formData.is_enabled ? 'translate-x-6' : 'translate-x-1'}`} />
                                </button>
                            </div>

                            <div className="pt-2 flex justify-end gap-3">
                                <button
                                    type="button"
                                    onClick={() => setIsModalOpen(false)}
                                    className="px-4 py-2 text-sm font-medium text-slate-600 hover:bg-slate-50 rounded-lg transition-colors"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    disabled={saving}
                                    className="px-6 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-medium rounded-lg shadow-lg shadow-indigo-500/30 transition-all flex items-center gap-2"
                                >
                                    {saving && <ArrowPathIcon className="w-4 h-4 animate-spin" />}
                                    Save Changes
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}
