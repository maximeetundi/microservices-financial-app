'use client';

import { useState, useEffect } from 'react';
import {
    UserGroupIcon,
    BanknotesIcon,
    MagnifyingGlassIcon,
    EyeIcon,
    NoSymbolIcon,
    CheckCircleIcon,
    ArrowPathIcon,
} from '@heroicons/react/24/outline';
import { getAllAssociations, getAssociationDetails, suspendAssociation, activateAssociation } from '@/lib/api';

interface Association {
    id: string;
    name: string;
    type: string;
    description: string;
    total_members: number;
    treasury_balance: number;
    currency: string;
    status: string;
    created_at: string;
}

export default function AssociationsPage() {
    const [associations, setAssociations] = useState<Association[]>([]);
    const [loading, setLoading] = useState(true);
    const [searchQuery, setSearchQuery] = useState('');
    const [selectedAssociation, setSelectedAssociation] = useState<Association | null>(null);
    const [showDetails, setShowDetails] = useState(false);

    useEffect(() => {
        loadAssociations();
    }, []);

    const loadAssociations = async () => {
        setLoading(true);
        try {
            const response = await getAllAssociations();
            setAssociations(response.data || []);
        } catch (error) {
            console.error('Failed to load associations:', error);
            // Mock data for demo
            setAssociations([
                {
                    id: '1',
                    name: 'Tontine Famille Toure',
                    type: 'tontine',
                    description: 'Tontine mensuelle familiale',
                    total_members: 12,
                    treasury_balance: 1200000,
                    currency: 'XOF',
                    status: 'active',
                    created_at: '2024-01-01',
                },
                {
                    id: '2',
                    name: 'Épargne Commerçants Plateau',
                    type: 'savings',
                    description: 'Caisse de solidarité commerçants',
                    total_members: 45,
                    treasury_balance: 5500000,
                    currency: 'XOF',
                    status: 'active',
                    created_at: '2024-02-15',
                },
                {
                    id: '3',
                    name: 'Crédit Mutuel Artisans',
                    type: 'credit',
                    description: 'Association de crédit pour artisans',
                    total_members: 28,
                    treasury_balance: 3200000,
                    currency: 'XOF',
                    status: 'suspended',
                    created_at: '2024-03-10',
                },
            ]);
        } finally {
            setLoading(false);
        }
    };

    const handleSuspend = async (id: string) => {
        try {
            await suspendAssociation(id);
            loadAssociations();
        } catch (error) {
            console.error('Failed to suspend:', error);
        }
    };

    const handleActivate = async (id: string) => {
        try {
            await activateAssociation(id);
            loadAssociations();
        } catch (error) {
            console.error('Failed to activate:', error);
        }
    };

    const formatCurrency = (amount: number, currency: string) => {
        return new Intl.NumberFormat('fr-FR', { style: 'currency', currency }).format(amount);
    };

    const getTypeLabel = (type: string) => {
        const labels: Record<string, string> = {
            tontine: 'Tontine',
            savings: 'Épargne',
            credit: 'Crédit Mutuel',
            general: 'Général',
        };
        return labels[type] || type;
    };

    const getStatusBadge = (status: string) => {
        const styles: Record<string, string> = {
            active: 'bg-green-100 text-green-800',
            suspended: 'bg-red-100 text-red-800',
            closed: 'bg-gray-100 text-gray-800',
        };
        const labels: Record<string, string> = {
            active: 'Actif',
            suspended: 'Suspendu',
            closed: 'Clôturé',
        };
        return (
            <span className={`px-2 py-1 rounded-full text-xs font-medium ${styles[status] || 'bg-gray-100 text-gray-800'}`}>
                {labels[status] || status}
            </span>
        );
    };

    const filteredAssociations = associations.filter(
        (a) =>
            a.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
            a.type.toLowerCase().includes(searchQuery.toLowerCase())
    );

    const stats = {
        total: associations.length,
        active: associations.filter((a) => a.status === 'active').length,
        totalTreasury: associations.reduce((sum, a) => sum + a.treasury_balance, 0),
        totalMembers: associations.reduce((sum, a) => sum + a.total_members, 0),
    };

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <h1 className="text-2xl font-bold text-gray-900">Gestion des Associations</h1>
                <button
                    onClick={loadAssociations}
                    className="flex items-center space-x-2 px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
                >
                    <ArrowPathIcon className="w-5 h-5" />
                    <span>Actualiser</span>
                </button>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                <div className="bg-white rounded-lg shadow p-4">
                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-sm text-gray-500">Total Associations</p>
                            <p className="text-2xl font-bold text-gray-900">{stats.total}</p>
                        </div>
                        <UserGroupIcon className="w-10 h-10 text-indigo-500" />
                    </div>
                </div>
                <div className="bg-white rounded-lg shadow p-4">
                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-sm text-gray-500">Actives</p>
                            <p className="text-2xl font-bold text-green-600">{stats.active}</p>
                        </div>
                        <CheckCircleIcon className="w-10 h-10 text-green-500" />
                    </div>
                </div>
                <div className="bg-white rounded-lg shadow p-4">
                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-sm text-gray-500">Trésorerie Totale</p>
                            <p className="text-2xl font-bold text-indigo-600">{formatCurrency(stats.totalTreasury, 'XOF')}</p>
                        </div>
                        <BanknotesIcon className="w-10 h-10 text-indigo-500" />
                    </div>
                </div>
                <div className="bg-white rounded-lg shadow p-4">
                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-sm text-gray-500">Membres Total</p>
                            <p className="text-2xl font-bold text-gray-900">{stats.totalMembers}</p>
                        </div>
                        <UserGroupIcon className="w-10 h-10 text-gray-500" />
                    </div>
                </div>
            </div>

            {/* Search */}
            <div className="bg-white rounded-lg shadow p-4">
                <div className="relative">
                    <MagnifyingGlassIcon className="w-5 h-5 absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
                    <input
                        type="text"
                        placeholder="Rechercher une association..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                    />
                </div>
            </div>

            {/* Table */}
            <div className="bg-white rounded-lg shadow overflow-hidden">
                <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                        <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Association
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Type
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Membres
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Trésorerie
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Statut
                            </th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Actions
                            </th>
                        </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                        {loading ? (
                            <tr>
                                <td colSpan={6} className="px-6 py-12 text-center text-gray-500">
                                    Chargement...
                                </td>
                            </tr>
                        ) : filteredAssociations.length === 0 ? (
                            <tr>
                                <td colSpan={6} className="px-6 py-12 text-center text-gray-500">
                                    Aucune association trouvée
                                </td>
                            </tr>
                        ) : (
                            filteredAssociations.map((association) => (
                                <tr key={association.id} className="hover:bg-gray-50">
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <div>
                                            <div className="text-sm font-medium text-gray-900">{association.name}</div>
                                            <div className="text-sm text-gray-500 truncate max-w-xs">
                                                {association.description}
                                            </div>
                                        </div>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <span className="px-2 py-1 bg-indigo-100 text-indigo-800 rounded-full text-xs font-medium">
                                            {getTypeLabel(association.type)}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                        {association.total_members}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                        {formatCurrency(association.treasury_balance, association.currency)}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap">{getStatusBadge(association.status)}</td>
                                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                        <button
                                            onClick={() => {
                                                setSelectedAssociation(association);
                                                setShowDetails(true);
                                            }}
                                            className="text-indigo-600 hover:text-indigo-900 mr-3"
                                            title="Voir détails"
                                        >
                                            <EyeIcon className="w-5 h-5" />
                                        </button>
                                        {association.status === 'active' ? (
                                            <button
                                                onClick={() => handleSuspend(association.id)}
                                                className="text-red-600 hover:text-red-900"
                                                title="Suspendre"
                                            >
                                                <NoSymbolIcon className="w-5 h-5" />
                                            </button>
                                        ) : (
                                            <button
                                                onClick={() => handleActivate(association.id)}
                                                className="text-green-600 hover:text-green-900"
                                                title="Réactiver"
                                            >
                                                <CheckCircleIcon className="w-5 h-5" />
                                            </button>
                                        )}
                                    </td>
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>

            {/* Details Modal */}
            {showDetails && selectedAssociation && (
                <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
                    <div className="bg-white rounded-xl shadow-xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
                        <div className="p-6 border-b border-gray-200">
                            <div className="flex justify-between items-start">
                                <div>
                                    <h2 className="text-xl font-bold text-gray-900">{selectedAssociation.name}</h2>
                                    <p className="text-gray-500">{getTypeLabel(selectedAssociation.type)}</p>
                                </div>
                                <button
                                    onClick={() => setShowDetails(false)}
                                    className="text-gray-400 hover:text-gray-600"
                                >
                                    ✕
                                </button>
                            </div>
                        </div>
                        <div className="p-6 space-y-4">
                            <div className="grid grid-cols-2 gap-4">
                                <div className="bg-gray-50 p-4 rounded-lg">
                                    <p className="text-sm text-gray-500">Membres</p>
                                    <p className="text-2xl font-bold text-gray-900">{selectedAssociation.total_members}</p>
                                </div>
                                <div className="bg-indigo-50 p-4 rounded-lg">
                                    <p className="text-sm text-indigo-600">Trésorerie</p>
                                    <p className="text-2xl font-bold text-indigo-600">
                                        {formatCurrency(selectedAssociation.treasury_balance, selectedAssociation.currency)}
                                    </p>
                                </div>
                            </div>
                            <div>
                                <p className="text-sm text-gray-500 mb-1">Description</p>
                                <p className="text-gray-900">{selectedAssociation.description || 'Aucune description'}</p>
                            </div>
                            <div>
                                <p className="text-sm text-gray-500 mb-1">Statut</p>
                                {getStatusBadge(selectedAssociation.status)}
                            </div>
                            <div>
                                <p className="text-sm text-gray-500 mb-1">Date de création</p>
                                <p className="text-gray-900">
                                    {new Date(selectedAssociation.created_at).toLocaleDateString('fr-FR')}
                                </p>
                            </div>
                        </div>
                        <div className="p-6 border-t border-gray-200 flex justify-end space-x-3">
                            <button
                                onClick={() => setShowDetails(false)}
                                className="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
                            >
                                Fermer
                            </button>
                            {selectedAssociation.status === 'active' ? (
                                <button
                                    onClick={() => {
                                        handleSuspend(selectedAssociation.id);
                                        setShowDetails(false);
                                    }}
                                    className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
                                >
                                    Suspendre
                                </button>
                            ) : (
                                <button
                                    onClick={() => {
                                        handleActivate(selectedAssociation.id);
                                        setShowDetails(false);
                                    }}
                                    className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700"
                                >
                                    Réactiver
                                </button>
                            )}
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
