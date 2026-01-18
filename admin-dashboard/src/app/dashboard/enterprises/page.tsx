"use client";

import { useState, useEffect } from "react";
import {
    BuildingOfficeIcon,
    UsersIcon,
    BanknotesIcon,
    PlusIcon,
    XMarkIcon,
    EyeIcon,
    PencilIcon,
    TrashIcon,
    CheckCircleIcon,
    XCircleIcon
} from "@heroicons/react/24/outline";

interface Enterprise {
    id: string;
    _id?: string;
    name: string;
    type: string;
    status: string;
    employees_count?: number;
    owner_id?: string;
    created_at?: string;
    description?: string;
    logo?: string;
}

export default function EnterprisesPage() {
    const [enterprises, setEnterprises] = useState<Enterprise[]>([]);
    const [loading, setLoading] = useState(true);
    const [stats, setStats] = useState({ total: 0, active: 0, employees: 0 });

    // Modal states
    const [selectedEnterprise, setSelectedEnterprise] = useState<Enterprise | null>(null);
    const [showViewModal, setShowViewModal] = useState(false);
    const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
    const [actionLoading, setActionLoading] = useState(false);

    const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

    const fetchEnterprises = async () => {
        setLoading(true);
        try {
            const token = localStorage.getItem('admin_token');
            const res = await fetch(`${API_BASE}/enterprise-service/api/v1/admin/enterprises`, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });

            if (res.ok) {
                const data = await res.json();
                const list = Array.isArray(data) ? data : (data.enterprises || []);
                setEnterprises(list);

                // Calculate stats
                const active = list.filter((e: Enterprise) => e.status === 'ACTIVE').length;
                const totalEmployees = list.reduce((sum: number, e: Enterprise) => sum + (e.employees_count || 0), 0);
                setStats({ total: list.length, active, employees: totalEmployees });
            }
        } catch (error) {
            console.error('Failed to fetch enterprises:', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchEnterprises();
    }, []);

    const handleView = (ent: Enterprise) => {
        setSelectedEnterprise(ent);
        setShowViewModal(true);
    };

    const handleStatusChange = async (ent: Enterprise, newStatus: string) => {
        setActionLoading(true);
        try {
            const token = localStorage.getItem('admin_token');
            const res = await fetch(`${API_BASE}/enterprise-service/api/v1/admin/enterprises/${ent.id || ent._id}/status`, {
                method: 'PUT',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ status: newStatus })
            });

            if (res.ok) {
                await fetchEnterprises();
                setShowViewModal(false);
            } else {
                alert('Erreur lors du changement de statut');
            }
        } catch (error) {
            console.error('Failed to update status:', error);
            alert('Erreur réseau');
        } finally {
            setActionLoading(false);
        }
    };

    const handleDelete = async () => {
        if (!selectedEnterprise) return;
        setActionLoading(true);
        try {
            const token = localStorage.getItem('admin_token');
            const res = await fetch(`${API_BASE}/enterprise-service/api/v1/admin/enterprises/${selectedEnterprise.id || selectedEnterprise._id}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });

            if (res.ok) {
                await fetchEnterprises();
                setShowDeleteConfirm(false);
                setShowViewModal(false);
            } else {
                alert('Erreur lors de la suppression');
            }
        } catch (error) {
            console.error('Failed to delete:', error);
            alert('Erreur réseau');
        } finally {
            setActionLoading(false);
        }
    };

    const getTypeColor = (type: string) => {
        const colors: Record<string, string> = {
            'SERVICE': 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300',
            'SCHOOL': 'bg-purple-100 text-purple-700 dark:bg-purple-900 dark:text-purple-300',
            'TRANSPORT': 'bg-orange-100 text-orange-700 dark:bg-orange-900 dark:text-orange-300',
            'UTILITY': 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300',
        };
        return colors[type] || 'bg-gray-100 text-gray-700';
    };

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Enterprises Management</h1>
                <button
                    onClick={() => fetchEnterprises()}
                    className="flex items-center space-x-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
                >
                    <span>🔄</span>
                    <span>Refresh</span>
                </button>
            </div>

            {/* Metrics */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="bg-white dark:bg-gray-800 p-6 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700">
                    <div className="flex items-center space-x-4">
                        <div className="p-3 bg-blue-100 dark:bg-blue-900 text-blue-600 rounded-full">
                            <BuildingOfficeIcon className="w-6 h-6" />
                        </div>
                        <div>
                            <p className="text-sm text-gray-500">Total Enterprises</p>
                            <p className="text-2xl font-bold dark:text-white">{stats.total}</p>
                        </div>
                    </div>
                </div>
                <div className="bg-white dark:bg-gray-800 p-6 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700">
                    <div className="flex items-center space-x-4">
                        <div className="p-3 bg-green-100 dark:bg-green-900 text-green-600 rounded-full">
                            <CheckCircleIcon className="w-6 h-6" />
                        </div>
                        <div>
                            <p className="text-sm text-gray-500">Active</p>
                            <p className="text-2xl font-bold dark:text-white">{stats.active}</p>
                        </div>
                    </div>
                </div>
                <div className="bg-white dark:bg-gray-800 p-6 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700">
                    <div className="flex items-center space-x-4">
                        <div className="p-3 bg-purple-100 dark:bg-purple-900 text-purple-600 rounded-full">
                            <UsersIcon className="w-6 h-6" />
                        </div>
                        <div>
                            <p className="text-sm text-gray-500">Total Employees</p>
                            <p className="text-2xl font-bold dark:text-white">{stats.employees}</p>
                        </div>
                    </div>
                </div>
            </div>

            {/* Loading */}
            {loading && (
                <div className="flex justify-center py-12">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
                </div>
            )}

            {/* Table */}
            {!loading && (
                <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden">
                    <div className="overflow-x-auto">
                        <table className="w-full text-left">
                            <thead className="bg-gray-50 dark:bg-gray-700 text-xs uppercase text-gray-500 dark:text-gray-400">
                                <tr>
                                    <th className="px-6 py-4 font-semibold">Name</th>
                                    <th className="px-6 py-4 font-semibold">Type</th>
                                    <th className="px-6 py-4 font-semibold">Status</th>
                                    <th className="px-6 py-4 font-semibold">Employees</th>
                                    <th className="px-6 py-4 font-semibold">Created</th>
                                    <th className="px-6 py-4 font-semibold">Actions</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                                {enterprises.map((ent) => (
                                    <tr key={ent.id || ent._id} className="hover:bg-gray-50 dark:hover:bg-gray-750 transition">
                                        <td className="px-6 py-4">
                                            <div className="flex items-center space-x-3">
                                                {ent.logo ? (
                                                    <img src={ent.logo} className="w-8 h-8 rounded-full object-cover" />
                                                ) : (
                                                    <div className="w-8 h-8 rounded-full bg-blue-100 dark:bg-blue-900 flex items-center justify-center text-blue-600 font-bold">
                                                        {ent.name?.charAt(0)}
                                                    </div>
                                                )}
                                                <span className="font-medium text-gray-900 dark:text-white">{ent.name}</span>
                                            </div>
                                        </td>
                                        <td className="px-6 py-4">
                                            <span className={`px-2 py-1 rounded text-xs font-medium ${getTypeColor(ent.type)}`}>
                                                {ent.type}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4">
                                            <span className={`px-2 py-1 rounded-full text-xs font-medium ${ent.status === 'ACTIVE'
                                                    ? 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300'
                                                    : ent.status === 'SUSPENDED'
                                                        ? 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300'
                                                        : 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300'
                                                }`}>
                                                {ent.status}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 text-gray-500 dark:text-gray-400">
                                            {ent.employees_count || 0}
                                        </td>
                                        <td className="px-6 py-4 text-gray-500 dark:text-gray-400 text-sm">
                                            {ent.created_at ? new Date(ent.created_at).toLocaleDateString() : '-'}
                                        </td>
                                        <td className="px-6 py-4">
                                            <div className="flex items-center space-x-2">
                                                <button
                                                    onClick={() => handleView(ent)}
                                                    className="p-2 text-blue-600 hover:bg-blue-100 dark:hover:bg-blue-900 rounded-lg transition"
                                                    title="View Details"
                                                >
                                                    <EyeIcon className="w-5 h-5" />
                                                </button>
                                                {ent.status === 'ACTIVE' ? (
                                                    <button
                                                        onClick={() => handleStatusChange(ent, 'SUSPENDED')}
                                                        className="p-2 text-red-600 hover:bg-red-100 dark:hover:bg-red-900 rounded-lg transition"
                                                        title="Suspend"
                                                    >
                                                        <XCircleIcon className="w-5 h-5" />
                                                    </button>
                                                ) : (
                                                    <button
                                                        onClick={() => handleStatusChange(ent, 'ACTIVE')}
                                                        className="p-2 text-green-600 hover:bg-green-100 dark:hover:bg-green-900 rounded-lg transition"
                                                        title="Activate"
                                                    >
                                                        <CheckCircleIcon className="w-5 h-5" />
                                                    </button>
                                                )}
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                                {enterprises.length === 0 && (
                                    <tr>
                                        <td colSpan={6} className="px-6 py-12 text-center text-gray-500">
                                            No enterprises found
                                        </td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>
            )}

            {/* View Modal */}
            {showViewModal && selectedEnterprise && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50">
                    <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
                        <div className="p-6 border-b border-gray-100 dark:border-gray-700 flex justify-between items-center">
                            <h3 className="text-xl font-bold text-gray-900 dark:text-white">
                                {selectedEnterprise.name}
                            </h3>
                            <button
                                onClick={() => setShowViewModal(false)}
                                className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg"
                            >
                                <XMarkIcon className="w-5 h-5" />
                            </button>
                        </div>

                        <div className="p-6 space-y-4">
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <p className="text-sm text-gray-500">Type</p>
                                    <p className="font-medium text-gray-900 dark:text-white">{selectedEnterprise.type}</p>
                                </div>
                                <div>
                                    <p className="text-sm text-gray-500">Status</p>
                                    <span className={`px-2 py-1 rounded-full text-xs font-medium ${selectedEnterprise.status === 'ACTIVE'
                                            ? 'bg-green-100 text-green-700'
                                            : 'bg-yellow-100 text-yellow-700'
                                        }`}>
                                        {selectedEnterprise.status}
                                    </span>
                                </div>
                                <div>
                                    <p className="text-sm text-gray-500">Employees</p>
                                    <p className="font-medium text-gray-900 dark:text-white">{selectedEnterprise.employees_count || 0}</p>
                                </div>
                                <div>
                                    <p className="text-sm text-gray-500">Owner ID</p>
                                    <p className="font-medium text-gray-900 dark:text-white text-xs">{selectedEnterprise.owner_id || '-'}</p>
                                </div>
                            </div>

                            {selectedEnterprise.description && (
                                <div>
                                    <p className="text-sm text-gray-500">Description</p>
                                    <p className="text-gray-900 dark:text-white">{selectedEnterprise.description}</p>
                                </div>
                            )}

                            <div className="border-t border-gray-100 dark:border-gray-700 pt-4">
                                <p className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">Actions Admin</p>
                                <div className="flex flex-wrap gap-2">
                                    {selectedEnterprise.status === 'ACTIVE' ? (
                                        <button
                                            onClick={() => handleStatusChange(selectedEnterprise, 'SUSPENDED')}
                                            disabled={actionLoading}
                                            className="px-4 py-2 bg-red-100 text-red-700 rounded-lg hover:bg-red-200 disabled:opacity-50 flex items-center gap-2"
                                        >
                                            <XCircleIcon className="w-4 h-4" />
                                            Suspendre
                                        </button>
                                    ) : (
                                        <button
                                            onClick={() => handleStatusChange(selectedEnterprise, 'ACTIVE')}
                                            disabled={actionLoading}
                                            className="px-4 py-2 bg-green-100 text-green-700 rounded-lg hover:bg-green-200 disabled:opacity-50 flex items-center gap-2"
                                        >
                                            <CheckCircleIcon className="w-4 h-4" />
                                            Activer
                                        </button>
                                    )}
                                    <button
                                        onClick={() => setShowDeleteConfirm(true)}
                                        disabled={actionLoading}
                                        className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50 flex items-center gap-2"
                                    >
                                        <TrashIcon className="w-4 h-4" />
                                        Supprimer
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            )}

            {/* Delete Confirmation */}
            {showDeleteConfirm && (
                <div className="fixed inset-0 z-[60] flex items-center justify-center p-4 bg-black/50">
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-2xl w-full max-w-md p-6">
                        <div className="text-center">
                            <div className="w-16 h-16 mx-auto mb-4 bg-red-100 dark:bg-red-900 rounded-full flex items-center justify-center">
                                <TrashIcon className="w-8 h-8 text-red-600" />
                            </div>
                            <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-2">
                                Supprimer cette entreprise ?
                            </h3>
                            <p className="text-gray-500 mb-6">
                                Cette action est irréversible. Tous les employés, services et données seront perdus.
                            </p>
                            <div className="flex gap-3 justify-center">
                                <button
                                    onClick={() => setShowDeleteConfirm(false)}
                                    className="px-4 py-2 bg-gray-100 dark:bg-gray-700 rounded-lg hover:bg-gray-200"
                                >
                                    Annuler
                                </button>
                                <button
                                    onClick={handleDelete}
                                    disabled={actionLoading}
                                    className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50"
                                >
                                    {actionLoading ? 'Suppression...' : 'Confirmer'}
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
