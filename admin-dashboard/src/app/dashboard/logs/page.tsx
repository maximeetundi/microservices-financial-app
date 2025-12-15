'use client';

import { useEffect, useState } from 'react';
import { getAuditLogs } from '@/lib/api';
import { format } from 'date-fns';

interface AuditLog {
    id: string;
    admin_id: string;
    admin_email: string;
    action: string;
    resource: string;
    resource_id: string;
    ip_address: string;
    created_at: string;
}

export default function LogsPage() {
    const [logs, setLogs] = useState<AuditLog[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchLogs = async () => {
            try {
                const response = await getAuditLogs();
                setLogs(response.data.logs || []);
            } catch (error) {
                console.error('Failed to fetch logs:', error);
            } finally {
                setLoading(false);
            }
        };
        fetchLogs();
    }, []);

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
                <h1 className="text-2xl font-bold text-slate-900">Logs d'audit</h1>
                <p className="text-slate-500 mt-1">Historique des actions administratives</p>
            </div>

            <div className="table-container">
                <table className="w-full">
                    <thead className="bg-gray-50 border-b">
                        <tr>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Date</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Admin</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Action</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">Ressource</th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-slate-600">IP</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-100">
                        {logs.map((log) => (
                            <tr key={log.id} className="hover:bg-gray-50">
                                <td className="px-6 py-4 text-sm text-slate-500">
                                    {log.created_at ? format(new Date(log.created_at), 'dd/MM/yyyy HH:mm:ss') : '-'}
                                </td>
                                <td className="px-6 py-4 text-slate-700">{log.admin_email}</td>
                                <td className="px-6 py-4">
                                    <span className={`badge ${getActionBadge(log.action)}`}>
                                        {log.action}
                                    </span>
                                </td>
                                <td className="px-6 py-4 text-slate-700">
                                    {log.resource}
                                    {log.resource_id && (
                                        <span className="text-slate-400 text-sm ml-2">({log.resource_id.slice(0, 8)}...)</span>
                                    )}
                                </td>
                                <td className="px-6 py-4 text-slate-500 font-mono text-sm">{log.ip_address}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
                {logs.length === 0 && (
                    <div className="text-center py-12 text-slate-500">
                        Aucun log trouvé
                    </div>
                )}
            </div>
        </div>
    );
}

function getActionBadge(action: string): string {
    if (action?.includes('block') || action?.includes('delete')) return 'badge-danger';
    if (action?.includes('create') || action?.includes('approve')) return 'badge-success';
    if (action?.includes('login')) return 'badge-info';
    return 'badge-warning';
}
