'use client';

import { useEffect, useState } from 'react';
import { getUsers, approveKYC, rejectKYC } from '@/lib/api';
import { format } from 'date-fns';

interface User {
    id: string;
    email: string;
    first_name: string;
    last_name: string;
    kyc_level: number | string;
    created_at: string;
}

export default function KYCPage() {
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(true);
    const [actionLoading, setActionLoading] = useState<string | null>(null);

    const fetchUsers = async () => {
        try {
            const response = await getUsers();
            const allUsers = response.data.users || [];
            // Filter users with pending KYC
            setUsers(allUsers.filter((u: User) =>
                u.kyc_level === 'pending' || u.kyc_level === 1
            ));
        } catch (error) {
            console.error('Failed to fetch users:', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchUsers();
    }, []);

    const handleApprove = async (userId: string) => {
        setActionLoading(userId);
        try {
            await approveKYC(userId, 'verified');
            await fetchUsers();
            alert('KYC approuvé');
        } catch (error) {
            alert('Erreur');
        } finally {
            setActionLoading(null);
        }
    };

    const handleReject = async (userId: string) => {
        const reason = prompt('Raison du rejet:');
        if (!reason) return;

        setActionLoading(userId);
        try {
            await rejectKYC(userId, reason);
            await fetchUsers();
            alert('KYC rejeté');
        } catch (error) {
            alert('Erreur');
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
                <h1 className="text-2xl font-bold text-slate-900">Vérification KYC</h1>
                <p className="text-slate-500 mt-1">{users.length} demandes en attente</p>
            </div>

            {users.length === 0 ? (
                <div className="card text-center py-12">
                    <p className="text-slate-500">Aucune demande KYC en attente</p>
                </div>
            ) : (
                <div className="grid gap-4">
                    {users.map((user) => (
                        <div key={user.id} className="card flex items-center justify-between">
                            <div className="flex items-center gap-4">
                                <div className="w-12 h-12 rounded-full bg-yellow-100 flex items-center justify-center text-yellow-600 font-medium">
                                    {user.first_name?.[0]}{user.last_name?.[0]}
                                </div>
                                <div>
                                    <p className="font-medium text-slate-900">{user.first_name} {user.last_name}</p>
                                    <p className="text-sm text-slate-500">{user.email}</p>
                                    <p className="text-xs text-slate-400">
                                        Inscrit le {user.created_at ? format(new Date(user.created_at), 'dd/MM/yyyy') : '-'}
                                    </p>
                                </div>
                            </div>
                            <div className="flex gap-2">
                                <button
                                    onClick={() => handleApprove(user.id)}
                                    disabled={actionLoading === user.id}
                                    className="btn-primary text-sm"
                                >
                                    {actionLoading === user.id ? '...' : 'Approuver'}
                                </button>
                                <button
                                    onClick={() => handleReject(user.id)}
                                    disabled={actionLoading === user.id}
                                    className="btn-danger text-sm"
                                >
                                    Rejeter
                                </button>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
