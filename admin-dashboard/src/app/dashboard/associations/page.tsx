'use client';

import { useState, useEffect } from 'react';
import {
    UsersIcon,
    HeartIcon,
    MagnifyingGlassIcon,
    ArrowPathIcon,
    CheckCircleIcon,
    XCircleIcon,
    EyeIcon,
} from '@heroicons/react/24/outline';

interface Association {
    id: string;
    name: string;
    description: string;
    email: string;
    phone: string;
    status: string;
    campaigns_count: number;
    total_raised: number;
    created_at: string;
}

export default function AssociationsPage() {
    const [associations, setAssociations] = useState<Association[]>([]);
    const [loading, setLoading] = useState(true);
    const [searchQuery, setSearchQuery] = useState('');

    useEffect(() => {
        loadAssociations();
    }, []);

    const loadAssociations = async () => {
        setLoading(true);
        try {
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8088';
            const token = localStorage.getItem('admin_token');

            // Get campaigns grouped by association
            const response = await fetch(`${API_URL}/api/v1/admin/campaigns?limit=100`, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
            });

            if (response.ok) {
                const data = await response.json();
                // Group by association from campaigns
                const campaigns = data.campaigns || [];
                const grouped: { [key: string]: Association } = {};

                campaigns.forEach((campaign: any) => {
                    const userId = campaign.user_id || campaign.creator_id || 'unknown';
                    if (!grouped[userId]) {
                        grouped[userId] = {
                            id: userId,
                            name: campaign.organization_name || campaign.title || 'Association sans nom',
                            description: campaign.description?.substring(0, 100) || '',
                            email: campaign.contact_email || '',
                            phone: campaign.contact_phone || '',
                            status: 'active',
                            campaigns_count: 0,
                            total_raised: 0,
                            created_at: campaign.created_at,
                        };
                    }
                    grouped[userId].campaigns_count++;
                    grouped[userId].total_raised += campaign.current_amount || 0;
                });

                setAssociations(Object.values(grouped));
            }
        } catch (error) {
            console.error('Error loading associations:', error);
        } finally {
            setLoading(false);
        }
    };

    const filteredAssociations = associations.filter(a =>
        a.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        a.email.toLowerCase().includes(searchQuery.toLowerCase())
    );

    return (
        <div className="space-y-6 animate-fadeIn">
            {/* Header */}
            <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900">Associations & Organisations</h1>
                    <p className="text-gray-500 mt-1">Gérez les associations enregistrées sur la plateforme</p>
                </div>
                <button
                    onClick={loadAssociations}
                    className="btn-secondary flex items-center gap-2"
                >
                    <ArrowPathIcon className={`w-5 h-5 ${loading ? 'animate-spin' : ''}`} />
                    Actualiser
                </button>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="stat-card flex items-center gap-4">
                    <div className="p-3 rounded-xl bg-indigo-100">
                        <UsersIcon className="w-6 h-6 text-indigo-600" />
                    </div>
                    <div>
                        <p className="text-2xl font-bold text-gray-900">{associations.length}</p>
                        <p className="text-sm text-gray-500">Associations</p>
                    </div>
                </div>

                <div className="stat-card flex items-center gap-4">
                    <div className="p-3 rounded-xl bg-emerald-100">
                        <HeartIcon className="w-6 h-6 text-emerald-600" />
                    </div>
                    <div>
                        <p className="text-2xl font-bold text-emerald-600">
                            {associations.reduce((acc, a) => acc + a.campaigns_count, 0)}
                        </p>
                        <p className="text-sm text-gray-500">Campagnes</p>
                    </div>
                </div>

                <div className="stat-card flex items-center gap-4">
                    <div className="p-3 rounded-xl bg-amber-100">
                        <CheckCircleIcon className="w-6 h-6 text-amber-600" />
                    </div>
                    <div>
                        <p className="text-2xl font-bold text-amber-600">
                            {Math.round(associations.reduce((acc, a) => acc + a.total_raised, 0)).toLocaleString()} XOF
                        </p>
                        <p className="text-sm text-gray-500">Total Collecté</p>
                    </div>
                </div>
            </div>

            {/* Search */}
            <div className="card">
                <div className="relative">
                    <MagnifyingGlassIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                    <input
                        type="text"
                        placeholder="Rechercher une association..."
                        className="input pl-10 w-full"
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                    />
                </div>
            </div>

            {/* List */}
            {loading ? (
                <div className="flex justify-center py-12">
                    <div className="spinner w-8 h-8" />
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {filteredAssociations.map((association) => (
                        <div
                            key={association.id}
                            className="card hover:shadow-xl transition-all duration-300"
                        >
                            <div className="flex items-start justify-between mb-4">
                                <div className="flex items-center gap-3">
                                    <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-bold text-lg">
                                        {association.name[0]}
                                    </div>
                                    <div>
                                        <h3 className="font-semibold text-gray-900">{association.name}</h3>
                                        <p className="text-xs text-gray-500">{association.email || 'Pas d\'email'}</p>
                                    </div>
                                </div>
                            </div>

                            <p className="text-sm text-gray-600 line-clamp-2 mb-4">
                                {association.description || 'Aucune description'}
                            </p>

                            <div className="grid grid-cols-2 gap-3 text-sm mb-4">
                                <div className="p-2 bg-gray-50 rounded-lg">
                                    <p className="text-gray-500 text-xs">Campagnes</p>
                                    <p className="font-bold text-gray-900">{association.campaigns_count}</p>
                                </div>
                                <div className="p-2 bg-emerald-50 rounded-lg">
                                    <p className="text-gray-500 text-xs">Collecté</p>
                                    <p className="font-bold text-emerald-600">
                                        {Math.round(association.total_raised).toLocaleString()} XOF
                                    </p>
                                </div>
                            </div>

                            <div className="pt-3 border-t border-gray-100">
                                <button className="btn-secondary w-full flex items-center justify-center gap-2 text-sm">
                                    <EyeIcon className="w-4 h-4" />
                                    Voir les campagnes
                                </button>
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {/* Empty State */}
            {!loading && filteredAssociations.length === 0 && (
                <div className="text-center py-12">
                    <UsersIcon className="w-12 h-12 mx-auto text-gray-300 mb-3" />
                    <p className="text-gray-500">Aucune association trouvée</p>
                </div>
            )}
        </div>
    );
}
