'use client';

import { useEffect, useState } from 'react';
import { getAdmins, getRoles, createAdmin, deleteAdmin } from '@/lib/api';
import { format } from 'date-fns';

interface Admin {
    id: string;
    email: string;
    first_name: string;
    last_name: string;
    role?: { name: string };
    is_active: boolean;
    created_at: string;
}

interface Role {
    id: string;
    name: string;
}

export default function AdminsPage() {
    const [admins, setAdmins] = useState<Admin[]>([]);
    const [roles, setRoles] = useState<Role[]>([]);
    const [loading, setLoading] = useState(true);
    const [showForm, setShowForm] = useState(false);
    const [form, setForm] = useState({ email: '', password: '', first_name: '', last_name: '', role_id: '' });

    const fetchData = async () => {
        try {
            const [adminsRes, rolesRes] = await Promise.all([getAdmins(), getRoles()]);
            setAdmins(adminsRes.data.admins || []);
            setRoles(rolesRes.data.roles || []);
        } catch (error) {
            console.error('Failed to fetch data:', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, []);

    const handleCreate = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await createAdmin(form);
            setShowForm(false);
            setForm({ email: '', password: '', first_name: '', last_name: '', role_id: '' });
            await fetchData();
            alert('Admin créé');
        } catch (error) {
            alert('Erreur lors de la création');
        }
    };

    const handleDelete = async (id: string) => {
        if (!confirm('Supprimer cet admin ?')) return;
        try {
            await deleteAdmin(id);
            await fetchData();
            alert('Admin supprimé');
        } catch (error) {
            alert('Erreur');
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
            <div className="mb-8 flex justify-between items-center">
                <div>
                    <h1 className="text-2xl font-bold text-slate-900">Administrateurs</h1>
                    <p className="text-slate-500 mt-1">{admins.length} administrateurs</p>
                </div>
                <button onClick={() => setShowForm(!showForm)} className="btn-primary">
                    {showForm ? 'Annuler' : '+ Nouvel admin'}
                </button>
            </div>

            {showForm && (
                <form onSubmit={handleCreate} className="card mb-6 grid grid-cols-2 gap-4">
                    <input
                        type="email"
                        placeholder="Email"
                        value={form.email}
                        onChange={(e) => setForm({ ...form, email: e.target.value })}
                        className="input"
                        required
                    />
                    <input
                        type="password"
                        placeholder="Mot de passe"
                        value={form.password}
                        onChange={(e) => setForm({ ...form, password: e.target.value })}
                        className="input"
                        required
                    />
                    <input
                        type="text"
                        placeholder="Prénom"
                        value={form.first_name}
                        onChange={(e) => setForm({ ...form, first_name: e.target.value })}
                        className="input"
                        required
                    />
                    <input
                        type="text"
                        placeholder="Nom"
                        value={form.last_name}
                        onChange={(e) => setForm({ ...form, last_name: e.target.value })}
                        className="input"
                        required
                    />
                    <select
                        value={form.role_id}
                        onChange={(e) => setForm({ ...form, role_id: e.target.value })}
                        className="input"
                        required
                    >
                        <option value="">Sélectionner un rôle</option>
                        {roles.map((role) => (
                            <option key={role.id} value={role.id}>{role.name}</option>
                        ))}
                    </select>
                    <button type="submit" className="btn-primary">Créer</button>
                </form>
            )}

            <div className="table-container">
                <table className="w-full">
                    <thead className="bg-gray-50 border-b">
                        <tr>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Admin</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Email</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Rôle</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Status</th>
                            <th className="text-right px-6 py-4 text-sm font-medium text-slate-600">Actions</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-100">
                        {admins.map((admin) => (
                            <tr key={admin.id} className="hover:bg-gray-50">
                                <td className="px-6 py-4 font-medium text-slate-900">
                                    {admin.first_name} {admin.last_name}
                                </td>
                                <td className="px-6 py-4 text-slate-700">{admin.email}</td>
                                <td className="px-6 py-4">
                                    <span className="badge badge-info">{admin.role?.name || 'N/A'}</span>
                                </td>
                                <td className="px-6 py-4">
                                    <span className={`badge ${admin.is_active ? 'badge-success' : 'badge-danger'}`}>
                                        {admin.is_active ? 'Actif' : 'Inactif'}
                                    </span>
                                </td>
                                <td className="px-6 py-4 text-right">
                                    <button
                                        onClick={() => handleDelete(admin.id)}
                                        className="btn-danger text-sm px-3 py-1"
                                    >
                                        Supprimer
                                    </button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
}
