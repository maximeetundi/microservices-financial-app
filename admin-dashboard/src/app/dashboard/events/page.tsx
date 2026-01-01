'use client';

import { useState, useEffect } from 'react';
import { getAllEvents, getEventStats, getEventTickets, verifyTicket, useTicket } from '@/lib/api';
import {
    TicketIcon,
    CalendarIcon,
    MapPinIcon,
    UserGroupIcon,
    CurrencyDollarIcon,
    QrCodeIcon,
    CheckCircleIcon,
    XCircleIcon,
    ClockIcon,
} from '@heroicons/react/24/outline';

interface Event {
    id: string;
    title: string;
    description: string;
    location: string;
    cover_image: string;
    start_date: string;
    end_date: string;
    status: string;
    event_code: string;
    currency: string;
    ticket_tiers: TicketTier[];
    total_sold?: number;
    total_revenue?: number;
}

interface TicketTier {
    id: string;
    name: string;
    icon: string;
    price: number;
    quantity: number;
    sold: number;
    color: string;
}

interface Ticket {
    id: string;
    ticket_code: string;
    tier_name: string;
    tier_icon: string;
    price: number;
    currency: string;
    status: string;
    form_data: Record<string, string>;
    used_at: string | null;
    created_at: string;
}

export default function EventsPage() {
    const [events, setEvents] = useState<Event[]>([]);
    const [selectedEvent, setSelectedEvent] = useState<Event | null>(null);
    const [eventStats, setEventStats] = useState<any>(null);
    const [eventTickets, setEventTickets] = useState<Ticket[]>([]);
    const [loading, setLoading] = useState(true);
    const [verifyCode, setVerifyCode] = useState('');
    const [verifyResult, setVerifyResult] = useState<any>(null);
    const [verifying, setVerifying] = useState(false);

    useEffect(() => {
        loadEvents();
    }, []);

    const loadEvents = async () => {
        setLoading(true);
        try {
            const res = await getAllEvents();
            setEvents(res.data?.events || []);
        } catch (e) {
            console.error('Failed to load events:', e);
        } finally {
            setLoading(false);
        }
    };

    const selectEvent = async (event: Event) => {
        setSelectedEvent(event);
        setVerifyResult(null);
        try {
            const [statsRes, ticketsRes] = await Promise.all([
                getEventStats(event.id),
                getEventTickets(event.id)
            ]);
            setEventStats(statsRes.data?.stats || statsRes.data);
            setEventTickets(ticketsRes.data?.tickets || []);
        } catch (e) {
            console.error('Failed to load event details:', e);
        }
    };

    const handleVerify = async () => {
        if (!verifyCode.trim()) return;
        setVerifying(true);
        setVerifyResult(null);
        try {
            const res = await verifyTicket(verifyCode.trim());
            setVerifyResult(res.data);
        } catch (e: any) {
            setVerifyResult({ valid: false, message: e.response?.data?.error || 'Erreur de vérification' });
        } finally {
            setVerifying(false);
        }
    };

    const handleUseTicket = async () => {
        if (!verifyResult?.ticket?.id) return;
        try {
            await useTicket(verifyResult.ticket.id);
            setVerifyResult({ ...verifyResult, can_use: false, message: 'Ticket marqué comme utilisé ✓' });
            // Refresh tickets
            if (selectedEvent) {
                const ticketsRes = await getEventTickets(selectedEvent.id);
                setEventTickets(ticketsRes.data?.tickets || []);
            }
        } catch (e: any) {
            alert(e.response?.data?.error || 'Erreur');
        }
    };

    const formatDate = (date: string) => {
        return new Date(date).toLocaleDateString('fr-FR', {
            day: 'numeric',
            month: 'short',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    const formatAmount = (amount: number) => {
        return new Intl.NumberFormat('fr-FR').format(amount);
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case 'active': return 'bg-green-500';
            case 'draft': return 'bg-gray-500';
            case 'ended': return 'bg-orange-500';
            case 'cancelled': return 'bg-red-500';
            default: return 'bg-gray-500';
        }
    };

    const getStatusLabel = (status: string) => {
        switch (status) {
            case 'active': return 'Actif';
            case 'draft': return 'Brouillon';
            case 'ended': return 'Terminé';
            case 'cancelled': return 'Annulé';
            default: return status;
        }
    };

    return (
        <div className="p-6">
            <div className="flex justify-between items-center mb-6">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900">🎫 Gestion des événements</h1>
                    <p className="text-gray-500">Visualisez et gérez les événements et tickets</p>
                </div>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                {/* Events List */}
                <div className="lg:col-span-1 bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
                    <div className="p-4 border-b border-gray-200">
                        <h2 className="font-semibold text-gray-900">Événements actifs</h2>
                    </div>
                    <div className="divide-y divide-gray-100 max-h-[600px] overflow-y-auto">
                        {loading ? (
                            <div className="p-8 text-center text-gray-500">Chargement...</div>
                        ) : events.length === 0 ? (
                            <div className="p-8 text-center text-gray-500">Aucun événement</div>
                        ) : (
                            events.map((event) => (
                                <div
                                    key={event.id}
                                    onClick={() => selectEvent(event)}
                                    className={`p-4 cursor-pointer hover:bg-gray-50 transition-colors ${selectedEvent?.id === event.id ? 'bg-indigo-50 border-l-4 border-indigo-500' : ''
                                        }`}
                                >
                                    <div className="flex items-start gap-3">
                                        <div className="w-12 h-12 bg-indigo-100 rounded-lg flex items-center justify-center text-2xl">
                                            🎪
                                        </div>
                                        <div className="flex-1 min-w-0">
                                            <h3 className="font-medium text-gray-900 truncate">{event.title}</h3>
                                            <div className="flex items-center gap-2 mt-1">
                                                <span className={`px-2 py-0.5 text-xs font-medium text-white rounded-full ${getStatusColor(event.status)}`}>
                                                    {getStatusLabel(event.status)}
                                                </span>
                                                <span className="text-xs text-gray-500">
                                                    {event.ticket_tiers?.length || 0} niveaux
                                                </span>
                                            </div>
                                            <p className="text-xs text-gray-500 mt-1 flex items-center gap-1">
                                                <CalendarIcon className="w-3 h-3" />
                                                {formatDate(event.start_date)}
                                            </p>
                                        </div>
                                    </div>
                                </div>
                            ))
                        )}
                    </div>
                </div>

                {/* Event Details */}
                <div className="lg:col-span-2 space-y-6">
                    {selectedEvent ? (
                        <>
                            {/* Event Header */}
                            <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
                                <div
                                    className="h-40 bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center"
                                    style={selectedEvent.cover_image ? { backgroundImage: `url(${selectedEvent.cover_image})`, backgroundSize: 'cover' } : {}}
                                >
                                    {!selectedEvent.cover_image && <span className="text-6xl">🎪</span>}
                                </div>
                                <div className="p-6">
                                    <div className="flex justify-between items-start">
                                        <div>
                                            <h2 className="text-xl font-bold text-gray-900">{selectedEvent.title}</h2>
                                            <p className="text-gray-500 mt-1">{selectedEvent.description}</p>
                                        </div>
                                        <span className={`px-3 py-1 text-sm font-medium text-white rounded-full ${getStatusColor(selectedEvent.status)}`}>
                                            {getStatusLabel(selectedEvent.status)}
                                        </span>
                                    </div>
                                    <div className="grid grid-cols-2 gap-4 mt-4">
                                        <div className="flex items-center gap-2 text-sm text-gray-600">
                                            <MapPinIcon className="w-4 h-4" />
                                            {selectedEvent.location || 'Non défini'}
                                        </div>
                                        <div className="flex items-center gap-2 text-sm text-gray-600">
                                            <CalendarIcon className="w-4 h-4" />
                                            {formatDate(selectedEvent.start_date)}
                                        </div>
                                        <div className="flex items-center gap-2 text-sm text-gray-600">
                                            <QrCodeIcon className="w-4 h-4" />
                                            Code: <span className="font-mono font-bold">{selectedEvent.event_code}</span>
                                        </div>
                                    </div>
                                </div>
                            </div>

                            {/* Stats */}
                            {eventStats && (
                                <div className="grid grid-cols-4 gap-4">
                                    <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-4">
                                        <div className="flex items-center gap-3">
                                            <div className="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
                                                <TicketIcon className="w-5 h-5 text-blue-600" />
                                            </div>
                                            <div>
                                                <p className="text-2xl font-bold text-gray-900">{eventStats.sold_tickets || 0}</p>
                                                <p className="text-xs text-gray-500">Vendus</p>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-4">
                                        <div className="flex items-center gap-3">
                                            <div className="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center">
                                                <CheckCircleIcon className="w-5 h-5 text-green-600" />
                                            </div>
                                            <div>
                                                <p className="text-2xl font-bold text-gray-900">{eventStats.used_tickets || 0}</p>
                                                <p className="text-xs text-gray-500">Utilisés</p>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-4">
                                        <div className="flex items-center gap-3">
                                            <div className="w-10 h-10 bg-yellow-100 rounded-lg flex items-center justify-center">
                                                <CurrencyDollarIcon className="w-5 h-5 text-yellow-600" />
                                            </div>
                                            <div>
                                                <p className="text-2xl font-bold text-gray-900">{formatAmount(eventStats.total_revenue || 0)}</p>
                                                <p className="text-xs text-gray-500">Revenus (XOF)</p>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-4">
                                        <div className="flex items-center gap-3">
                                            <div className="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center">
                                                <UserGroupIcon className="w-5 h-5 text-purple-600" />
                                            </div>
                                            <div>
                                                <p className="text-2xl font-bold text-gray-900">{selectedEvent.ticket_tiers?.length || 0}</p>
                                                <p className="text-xs text-gray-500">Niveaux</p>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            )}

                            {/* Ticket Verification */}
                            <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
                                <h3 className="font-semibold text-gray-900 mb-4">🔍 Vérifier un ticket</h3>
                                <div className="flex gap-3">
                                    <input
                                        type="text"
                                        value={verifyCode}
                                        onChange={(e) => setVerifyCode(e.target.value)}
                                        placeholder="Entrez le code du ticket (ex: TKT-XXXXXXXX)"
                                        className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                                        onKeyPress={(e) => e.key === 'Enter' && handleVerify()}
                                    />
                                    <button
                                        onClick={handleVerify}
                                        disabled={verifying}
                                        className="px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50"
                                    >
                                        {verifying ? 'Vérification...' : 'Vérifier'}
                                    </button>
                                </div>
                                {verifyResult && (
                                    <div className={`mt-4 p-4 rounded-lg ${verifyResult.valid ? 'bg-green-50 border border-green-200' : 'bg-red-50 border border-red-200'}`}>
                                        <div className="flex items-center gap-3">
                                            {verifyResult.valid ? (
                                                <CheckCircleIcon className="w-6 h-6 text-green-600" />
                                            ) : (
                                                <XCircleIcon className="w-6 h-6 text-red-600" />
                                            )}
                                            <div className="flex-1">
                                                <p className={`font-medium ${verifyResult.valid ? 'text-green-800' : 'text-red-800'}`}>
                                                    {verifyResult.message}
                                                </p>
                                                {verifyResult.ticket && (
                                                    <p className="text-sm text-gray-600 mt-1">
                                                        {verifyResult.ticket.tier_icon} {verifyResult.ticket.tier_name} - {verifyResult.ticket.ticket_code}
                                                    </p>
                                                )}
                                            </div>
                                            {verifyResult.can_use && (
                                                <button
                                                    onClick={handleUseTicket}
                                                    className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700"
                                                >
                                                    ✓ Valider l'entrée
                                                </button>
                                            )}
                                        </div>
                                    </div>
                                )}
                            </div>

                            {/* Tickets List */}
                            <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
                                <div className="p-4 border-b border-gray-200">
                                    <h3 className="font-semibold text-gray-900">Tickets vendus ({eventTickets.length})</h3>
                                </div>
                                <div className="overflow-x-auto">
                                    <table className="w-full">
                                        <thead className="bg-gray-50">
                                            <tr>
                                                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Code</th>
                                                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Niveau</th>
                                                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Prix</th>
                                                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Statut</th>
                                                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
                                            </tr>
                                        </thead>
                                        <tbody className="divide-y divide-gray-100">
                                            {eventTickets.map((ticket) => (
                                                <tr key={ticket.id} className="hover:bg-gray-50">
                                                    <td className="px-4 py-3 font-mono text-sm">{ticket.ticket_code}</td>
                                                    <td className="px-4 py-3">
                                                        <span className="text-lg">{ticket.tier_icon}</span> {ticket.tier_name}
                                                    </td>
                                                    <td className="px-4 py-3">{formatAmount(ticket.price)} {ticket.currency}</td>
                                                    <td className="px-4 py-3">
                                                        <span className={`px-2 py-1 text-xs font-medium rounded-full ${ticket.status === 'paid' ? 'bg-green-100 text-green-800' :
                                                                ticket.status === 'used' ? 'bg-gray-100 text-gray-800' :
                                                                    'bg-yellow-100 text-yellow-800'
                                                            }`}>
                                                            {ticket.status === 'paid' ? 'Valide' : ticket.status === 'used' ? 'Utilisé' : ticket.status}
                                                        </span>
                                                    </td>
                                                    <td className="px-4 py-3 text-sm text-gray-500">
                                                        {formatDate(ticket.created_at)}
                                                    </td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                    {eventTickets.length === 0 && (
                                        <div className="p-8 text-center text-gray-500">Aucun ticket vendu</div>
                                    )}
                                </div>
                            </div>
                        </>
                    ) : (
                        <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-12 text-center">
                            <div className="text-6xl mb-4">🎫</div>
                            <h3 className="text-lg font-medium text-gray-900">Sélectionnez un événement</h3>
                            <p className="text-gray-500 mt-2">Cliquez sur un événement pour voir les détails et gérer les tickets</p>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
