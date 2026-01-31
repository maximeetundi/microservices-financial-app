'use client';

import { useEffect, useState, useMemo } from 'react';
import { getAuditLogs } from '@/lib/api';
import { format, subDays, isWithinInterval, startOfDay, endOfDay } from 'date-fns';
import { fr } from 'date-fns/locale';
import {
    DocumentTextIcon,
    MagnifyingGlassIcon,
    ArrowDownTrayIcon,
    FunnelIcon,
    ChevronLeftIcon,
    ChevronRightIcon,
    UserCircleIcon,
    ShieldCheckIcon,
    TrashIcon,
    PlusCircleIcon,
    PencilSquareIcon,
    ArrowRightOnRectangleIcon,
    ExclamationTriangleIcon,
    EyeIcon,
    LockClosedIcon,
    CheckCircleIcon,
    XCircleIcon,
    ClockIcon,
    GlobeAltIcon,
} from '@heroicons/react/24/outline';
import PageHeader from '@/components/ui/PageHeader';
import StatsCard from '@/components/ui/StatsCard';
import Modal, { ModalFooter, ModalButton } from '@/components/ui/Modal';
import { SimpleToast } from '@/components/ui/Toast';

interface AuditLog {
    id: string;
    admin_id: string;
    admin_email: string;
    action: string;
    resource: string;
    resource_id: string;
    ip_address: string;
    user_agent?: string;
    details?: string;
    created_at: string;
}

type ActionFilter = 'all' | 'login' | 'create' | 'update' | 'delete' | 'block' | 'approve' | 'reject';
type DateFilter = 'all' | 'today' | 'week' | 'month';

const ITEMS_PER_PAGE = 15;

export default function LogsPage() {
    // Data state
    const [logs, setLogs] = useState<AuditLog[]>([]);
    const [loading, setLoading] = useState(true);
    const [refreshing, setRefreshing] = useState(false);

    // Filters
    const [search, setSearch] = useState('');
    const [actionFilter, setActionFilter] = useState<ActionFilter>('all');
    const [dateFilter, setDateFilter] = useState<DateFilter>('all');

    // Pagination
    const [currentPage, setCurrentPage] = useState(1);

    // Modal
    const [selectedLog, setSelectedLog] = useState<AuditLog | null>(null);
    const [detailsOpen, setDetailsOpen] = useState(false);

    // Toast
    const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' } | null>(null);

    const fetchLogs = async (silent = false) => {
        if (!silent) setLoading(true);
        else setRefreshing(true);

        try {
            const response = await getAuditLogs();
            setLogs(response.data?.logs || []);
        } catch (error) {
            console.error('Failed to fetch logs:', error);
            setToast({ message: 'Erreur lors du chargement des logs', type: 'error' });
        } finally {
            setLoading(false);
            setRefreshing(false);
        }
    };

    useEffect(() => {
        fetchLogs();
    }, []);

    // Computed stats
    const stats = useMemo(() => {
        const today = new Date();
        const todayLogs = logs.filter(log =>
            log.created_at && isWithinInterval(new Date(log.created_at), {
                start: startOfDay(today),
                end: endOfDay(today)
            })
        );

        const logins = logs.filter(l => l.action?.toLowerCase().includes('login')).length;
        const critical = logs.filter(l =>
            l.action?.toLowerCase().includes('delete') ||
            l.action?.toLowerCase().includes('block')
        ).length;
        const uniqueAdmins = new Set(logs.map(l => l.admin_id)).size;

        return { total: logs.length, today: todayLogs.length, logins, critical, uniqueAdmins };
    }, [logs]);

    // Filtered logs
    const filteredLogs = useMemo(() => {
        const today = new Date();

        return logs.filter(log => {
            // Search filter
            if (search) {
                const s = search.toLowerCase();
                const matchesSearch =
                    log.admin_email?.toLowerCase().includes(s) ||
                    log.action?.toLowerCase().includes(s) ||
                    log.resource?.toLowerCase().includes(s) ||
                    log.ip_address?.includes(s);
                if (!matchesSearch) return false;
            }

            // Action filter
            if (actionFilter !== 'all') {
                if (!log.action?.toLowerCase().includes(actionFilter)) return false;
            }

            // Date filter
            if (dateFilter !== 'all' && log.created_at) {
                const logDate = new Date(log.created_at);
                switch (dateFilter) {
                    case 'today':
                        if (!isWithinInterval(logDate, { start: startOfDay(today), end: endOfDay(today) })) return false;
                        break;
                    case 'week':
                        if (!isWithinInterval(logDate, { start: startOfDay(subDays(today, 7)), end: endOfDay(today) })) return false;
                        break;
                    case 'month':
                        if (!isWithinInterval(logDate, { start: startOfDay(subDays(today, 30)), end: endOfDay(today) })) return false;
                        break;
                }
            }

            return true;
        });
    }, [logs, search, actionFilter, dateFilter]);

    // Pagination
    const totalPages = Math.ceil(filteredLogs.length / ITEMS_PER_PAGE);
    const paginatedLogs = filteredLogs.slice(
        (currentPage - 1) * ITEMS_PER_PAGE,
        currentPage * ITEMS_PER_PAGE
    );

    // Reset page when filters change
    useEffect(() => {
        setCurrentPage(1);
    }, [search, actionFilter, dateFilter]);

    // Export to CSV
    const handleExport = () => {
        const headers = ['Date', 'Admin', 'Action', 'Ressource', 'ID Ressource', 'IP'];
        const rows = filteredLogs.map(log => [
            log.created_at ? format(new Date(log.created_at), 'dd/MM/yyyy HH:mm:ss') : '',
            log.admin_email || '',
            log.action || '',
            log.resource || '',
            log.resource_id || '',
            log.ip_address || '',
        ]);

        const csv = [headers, ...rows].map(row => row.map(cell => `"${cell}"`).join(',')).join('\n');
        const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
        const url = URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = `audit_logs_${format(new Date(), 'yyyy-MM-dd')}.csv`;
        link.click();
        URL.revokeObjectURL(url);

        setToast({ message: 'Export CSV téléchargé', type: 'success' });
    };

    const getActionConfig = (action: string) => {
        const configs: Record<string, { icon: React.ReactNode; bg: string; text: string; label: string }> = {
            login: {
                icon: <ArrowRightOnRectangleIcon className="w-4 h-4" />,
                bg: 'bg-blue-100',
                text: 'text-blue-700',
                label: 'Connexion',
            },
            create: {
                icon: <PlusCircleIcon className="w-4 h-4" />,
                bg: 'bg-green-100',
                text: 'text-green-700',
                label: 'Création',
            },
            update: {
                icon: <PencilSquareIcon className="w-4 h-4" />,
                bg: 'bg-amber-100',
                text: 'text-amber-700',
                label: 'Modification',
            },
            delete: {
                icon: <TrashIcon className="w-4 h-4" />,
                bg: 'bg-red-100',
                text: 'text-red-700',
                label: 'Suppression',
            },
            block: {
                icon: <LockClosedIcon className="w-4 h-4" />,
                bg: 'bg-red-100',
                text: 'text-red-700',
                label: 'Blocage',
            },
            approve: {
                icon: <CheckCircleIcon className="w-4 h-4" />,
                bg: 'bg-green-100',
                text: 'text-green-700',
                label: 'Approbation',
            },
            reject: {
                icon: <XCircleIcon className="w-4 h-4" />,
                bg: 'bg-red-100',
                text: 'text-red-700',
                label: 'Rejet',
            },
            view: {
                icon: <EyeIcon className="w-4 h-4" />,
                bg: 'bg-gray-100',
                text: 'text-gray-700',
                label: 'Consultation',
            },
        };

        const actionLower = action?.toLowerCase() || '';
        for (const [key, config] of Object.entries(configs)) {
            if (actionLower.includes(key)) return config;
        }

        return {
            icon: <DocumentTextIcon className="w-4 h-4" />,
            bg: 'bg-gray-100',
            text: 'text-gray-700',
            label: action,
        };
    };

    const handleViewDetails = (log: AuditLog) => {
        setSelectedLog(log);
        setDetailsOpen(true);
    };

    if (loading) {
        return (
            <div className="flex items-center justify-center h-96">
                <div className="text-center">
                    <div className="w-16 h-16 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin mx-auto mb-4" />
                    <p className="text-gray-500">Chargement des logs d'audit...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="pb-10">
            <PageHeader
                title="Logs d'Audit"
                subtitle="Historique des actions administratives"
                icon="📋"
                onRefresh={() => fetchLogs(true)}
                loading={refreshing}
                actions={
                    <button
                        onClick={handleExport}
                        className="flex items-center gap-2 px-4 py-2.5 bg-gradient-to-r from-green-500 to-emerald-600 text-white rounded-xl font-semibold hover:from-green-600 hover:to-emerald-700 transition-all shadow-lg shadow-green-500/25"
                    >
                        <ArrowDownTrayIcon className="w-4 h-4" />
                        Exporter CSV
                    </button>
                }
            />

            {/* Stats Cards */}
            <div className="grid grid-cols-2 lg:grid-cols-5 gap-4 mb-6">
                <StatsCard
                    title="Total"
                    value={stats.total}
                    icon={DocumentTextIcon}
                    color="primary"
                />
                <StatsCard
                    title="Aujourd'hui"
                    value={stats.today}
                    icon={ClockIcon}
                    color="info"
                />
                <StatsCard
                    title="Connexions"
                    value={stats.logins}
                    icon={ArrowRightOnRectangleIcon}
                    color="success"
                />
                <StatsCard
                    title="Actions Critiques"
                    value={stats.critical}
                    icon={ExclamationTriangleIcon}
                    color="danger"
                />
                <StatsCard
                    title="Admins Actifs"
                    value={stats.uniqueAdmins}
                    icon={UserCircleIcon}
                    color="purple"
                />
            </div>

            {/* Filters */}
            <div className="bg-white rounded-2xl border border-gray-200 shadow-sm p-4 mb-6">
                <div className="flex flex-col lg:flex-row gap-4">
                    {/* Search */}
                    <div className="flex-1 relative">
                        <MagnifyingGlassIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                        <input
                            type="text"
                            value={search}
                            onChange={(e) => setSearch(e.target.value)}
                            placeholder="Rechercher par admin, action, ressource, IP..."
                            className="w-full pl-10 pr-4 py-2.5 rounded-xl border border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 outline-none transition-all"
                        />
                    </div>

                    {/* Date Filter */}
                    <select
                        value={dateFilter}
                        onChange={(e) => setDateFilter(e.target.value as DateFilter)}
                        className="px-4 py-2.5 rounded-xl border border-gray-200 bg-white text-gray-700 focus:border-indigo-500 outline-none"
                    >
                        <option value="all">Toutes les dates</option>
                        <option value="today">Aujourd'hui</option>
                        <option value="week">7 derniers jours</option>
                        <option value="month">30 derniers jours</option>
                    </select>

                    {/* Action Filter */}
                    <select
                        value={actionFilter}
                        onChange={(e) => setActionFilter(e.target.value as ActionFilter)}
                        className="px-4 py-2.5 rounded-xl border border-gray-200 bg-white text-gray-700 focus:border-indigo-500 outline-none"
                    >
                        <option value="all">Toutes les actions</option>
                        <option value="login">🔑 Connexions</option>
                        <option value="create">➕ Créations</option>
                        <option value="update">✏️ Modifications</option>
                        <option value="delete">🗑️ Suppressions</option>
                        <option value="block">🔒 Blocages</option>
                        <option value="approve">✅ Approbations</option>
                        <option value="reject">❌ Rejets</option>
                    </select>
                </div>
            </div>

            {/* Table */}
            <div className="bg-white rounded-2xl border border-gray-200 shadow-sm overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="w-full">
                        <thead>
                            <tr className="bg-gradient-to-r from-gray-50 to-gray-100 border-b border-gray-200">
                                <th className="text-left px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Date</th>
                                <th className="text-left px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Admin</th>
                                <th className="text-left px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Action</th>
                                <th className="text-left px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Ressource</th>
                                <th className="text-left px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">IP</th>
                                <th className="text-right px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100">
                            {paginatedLogs.map((log) => {
                                const actionConfig = getActionConfig(log.action);
                                return (
                                    <tr key={log.id} className="hover:bg-gray-50/50 transition-colors">
                                        <td className="px-6 py-4">
                                            <div>
                                                <p className="text-sm font-medium text-gray-900">
                                                    {log.created_at ? format(new Date(log.created_at), 'dd MMM yyyy', { locale: fr }) : '-'}
                                                </p>
                                                <p className="text-xs text-gray-400">
                                                    {log.created_at ? format(new Date(log.created_at), 'HH:mm:ss') : ''}
                                                </p>
                                            </div>
                                        </td>
                                        <td className="px-6 py-4">
                                            <div className="flex items-center gap-3">
                                                <div className="w-8 h-8 rounded-full bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-bold text-xs">
                                                    {log.admin_email?.charAt(0).toUpperCase() || 'A'}
                                                </div>
                                                <div>
                                                    <p className="text-sm font-medium text-gray-900">{log.admin_email || 'N/A'}</p>
                                                    <p className="text-xs text-gray-400 font-mono">{log.admin_id?.slice(0, 8)}...</p>
                                                </div>
                                            </div>
                                        </td>
                                        <td className="px-6 py-4">
                                            <span className={`inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold ${actionConfig.bg} ${actionConfig.text}`}>
                                                {actionConfig.icon}
                                                {log.action || 'N/A'}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4">
                                            <div>
                                                <p className="text-sm text-gray-900">{log.resource || 'N/A'}</p>
                                                {log.resource_id && (
                                                    <p className="text-xs text-gray-400 font-mono">{log.resource_id.slice(0, 12)}...</p>
                                                )}
                                            </div>
                                        </td>
                                        <td className="px-6 py-4">
                                            <div className="flex items-center gap-2">
                                                <GlobeAltIcon className="w-4 h-4 text-gray-400" />
                                                <code className="text-sm text-gray-600">{log.ip_address || 'N/A'}</code>
                                            </div>
                                        </td>
                                        <td className="px-6 py-4">
                                            <div className="flex items-center justify-end">
                                                <button
                                                    onClick={() => handleViewDetails(log)}
                                                    className="p-2 rounded-lg hover:bg-gray-100 text-gray-500 hover:text-indigo-600 transition-colors"
                                                    title="Voir détails"
                                                >
                                                    <EyeIcon className="w-5 h-5" />
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                );
                            })}
                        </tbody>
                    </table>
                </div>

                {/* Empty State */}
                {filteredLogs.length === 0 && (
                    <div className="py-16 text-center">
                        <DocumentTextIcon className="w-16 h-16 text-gray-200 mx-auto mb-4" />
                        <h3 className="text-lg font-semibold text-gray-900 mb-2">Aucun log trouvé</h3>
                        <p className="text-gray-500 text-sm">Modifiez vos filtres pour voir plus de résultats</p>
                    </div>
                )}

                {/* Pagination */}
                {totalPages > 1 && (
                    <div className="flex items-center justify-between px-6 py-4 border-t border-gray-100">
                        <p className="text-sm text-gray-500">
                            Affichage de {(currentPage - 1) * ITEMS_PER_PAGE + 1} à{' '}
                            {Math.min(currentPage * ITEMS_PER_PAGE, filteredLogs.length)} sur{' '}
                            {filteredLogs.length} résultats
                        </p>
                        <div className="flex items-center gap-2">
                            <button
                                onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                                disabled={currentPage === 1}
                                className="p-2 rounded-lg border border-gray-200 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                                <ChevronLeftIcon className="w-5 h-5" />
                            </button>
                            <span className="px-4 py-2 text-sm font-medium text-gray-700">
                                Page {currentPage} / {totalPages}
                            </span>
                            <button
                                onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
                                disabled={currentPage === totalPages}
                                className="p-2 rounded-lg border border-gray-200 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                                <ChevronRightIcon className="w-5 h-5" />
                            </button>
                        </div>
                    </div>
                )}
            </div>

            {/* Details Modal */}
            <Modal
                isOpen={detailsOpen}
                onClose={() => setDetailsOpen(false)}
                title="Détails du Log"
                subtitle={selectedLog?.action}
                size="lg"
            >
                {selectedLog && (
                    <div className="space-y-6">
                        {/* Header */}
                        <div className="flex items-center gap-4 p-4 bg-gradient-to-r from-gray-50 to-gray-100 rounded-xl">
                            <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-bold text-lg">
                                {selectedLog.admin_email?.charAt(0).toUpperCase() || 'A'}
                            </div>
                            <div className="flex-1">
                                <h3 className="text-lg font-bold text-gray-900">{selectedLog.admin_email}</h3>
                                <p className="text-gray-500 text-sm font-mono">{selectedLog.admin_id}</p>
                            </div>
                            {(() => {
                                const config = getActionConfig(selectedLog.action);
                                return (
                                    <span className={`inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-semibold ${config.bg} ${config.text}`}>
                                        {config.icon}
                                        {selectedLog.action}
                                    </span>
                                );
                            })()}
                        </div>

                        {/* Details Grid */}
                        <div className="grid grid-cols-2 gap-4">
                            <div className="p-4 bg-gray-50 rounded-xl">
                                <p className="text-xs text-gray-400 uppercase font-bold mb-1">Date & Heure</p>
                                <p className="text-sm text-gray-900">
                                    {selectedLog.created_at
                                        ? format(new Date(selectedLog.created_at), 'dd MMMM yyyy à HH:mm:ss', { locale: fr })
                                        : '-'
                                    }
                                </p>
                            </div>
                            <div className="p-4 bg-gray-50 rounded-xl">
                                <p className="text-xs text-gray-400 uppercase font-bold mb-1">Adresse IP</p>
                                <p className="text-sm text-gray-900 font-mono">{selectedLog.ip_address || 'N/A'}</p>
                            </div>
                            <div className="p-4 bg-gray-50 rounded-xl">
                                <p className="text-xs text-gray-400 uppercase font-bold mb-1">Ressource</p>
                                <p className="text-sm text-gray-900">{selectedLog.resource || 'N/A'}</p>
                            </div>
                            <div className="p-4 bg-gray-50 rounded-xl">
                                <p className="text-xs text-gray-400 uppercase font-bold mb-1">ID Ressource</p>
                                <p className="text-sm text-gray-900 font-mono break-all">{selectedLog.resource_id || 'N/A'}</p>
                            </div>
                            {selectedLog.user_agent && (
                                <div className="col-span-2 p-4 bg-gray-50 rounded-xl">
                                    <p className="text-xs text-gray-400 uppercase font-bold mb-1">User Agent</p>
                                    <p className="text-xs text-gray-600 font-mono break-all">{selectedLog.user_agent}</p>
                                </div>
                            )}
                            {selectedLog.details && (
                                <div className="col-span-2 p-4 bg-gray-50 rounded-xl">
                                    <p className="text-xs text-gray-400 uppercase font-bold mb-1">Détails</p>
                                    <p className="text-sm text-gray-900">{selectedLog.details}</p>
                                </div>
                            )}
                        </div>

                        <ModalFooter>
                            <ModalButton variant="secondary" onClick={() => setDetailsOpen(false)}>
                                Fermer
                            </ModalButton>
                        </ModalFooter>
                    </div>
                )}
            </Modal>

            {/* Toast */}
            {toast && (
                <SimpleToast
                    message={toast.message}
                    type={toast.type}
                    onClose={() => setToast(null)}
                />
            )}
        </div>
    );
}
