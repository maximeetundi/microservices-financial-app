'use client';

import { useState, useEffect, useRef } from 'react';
import { getSupportTickets, getTicketMessages, sendTicketMessage, closeTicket, getSupportStats, getSupportAgents } from '@/lib/api';

interface Conversation {
    id: string;
    user_id: string;
    user_name: string;
    user_email: string;
    agent_type: 'ai' | 'human';
    agent_id?: string;
    subject: string;
    category: string;
    status: string;
    priority: string;
    last_message: string;
    last_message_at?: string;
    unread_count: number;
    message_count: number;
    created_at: string;
    updated_at: string;
}

interface Message {
    id: string;
    sender_id: string;
    sender_type: 'user' | 'agent' | 'system';
    sender_name: string;
    content: string;
    content_type: string;
    is_read: boolean;
    created_at: string;
}

interface Stats {
    total_conversations: number;
    open_conversations: number;
    resolved_today: number;
    pending_conversations: number;
    customer_satisfaction: number;
    active_agents: number;
    avg_response_time_minutes?: number;
}

interface Agent {
    id: string;
    name: string;
    email: string;
    type: 'ai' | 'human';
    is_available: boolean;
    active_chats: number;
    max_chats: number;
}

export default function SupportPage() {
    const [conversations, setConversations] = useState<Conversation[]>([]);
    const [selectedConv, setSelectedConv] = useState<Conversation | null>(null);
    const [messages, setMessages] = useState<Message[]>([]);
    const [newMessage, setNewMessage] = useState('');
    const [filter, setFilter] = useState('all');
    const [stats, setStats] = useState<Stats | null>(null);
    const [agents, setAgents] = useState<Agent[]>([]);
    const [sending, setSending] = useState(false);
    const [loading, setLoading] = useState(true);
    const [loadingMessages, setLoadingMessages] = useState(false);
    const messagesEndRef = useRef<HTMLDivElement>(null);

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    };

    const fetchConversations = async (statusFilter?: string) => {
        try {
            setLoading(true);
            const response = await getSupportTickets(50, 0, statusFilter !== 'all' ? statusFilter : undefined);
            const convs = response.data?.conversations || [];
            setConversations(convs);
        } catch (error) {
            console.error('Failed to fetch conversations:', error);
            setConversations([]);
        } finally {
            setLoading(false);
        }
    };

    const fetchStats = async () => {
        try {
            const response = await getSupportStats();
            setStats(response.data?.stats || response.data);
        } catch (error) {
            console.error('Failed to fetch stats:', error);
            setStats(null);
        }
    };

    const fetchAgents = async () => {
        try {
            const response = await getSupportAgents();
            setAgents(response.data?.agents || []);
        } catch (error) {
            console.error('Failed to fetch agents:', error);
        }
    };

    const fetchMessages = async (conversationId: string) => {
        try {
            setLoadingMessages(true);
            const response = await getTicketMessages(conversationId);
            const msgs = response.data?.messages || [];
            setMessages(msgs);
        } catch (error) {
            console.error('Failed to fetch messages:', error);
            setMessages([]);
        } finally {
            setLoadingMessages(false);
        }
    };

    useEffect(() => {
        fetchConversations();
        fetchStats();
        fetchAgents();

        // Refresh conversations every 30 seconds
        const interval = setInterval(() => {
            fetchConversations(filter !== 'all' ? filter : undefined);
            fetchStats();
        }, 30000);

        return () => clearInterval(interval);
    }, []);

    useEffect(() => {
        fetchConversations(filter !== 'all' ? filter : undefined);
    }, [filter]);

    useEffect(() => {
        scrollToBottom();
    }, [messages]);

    const selectConversation = async (conv: Conversation) => {
        setSelectedConv(conv);
        await fetchMessages(conv.id);
    };

    const sendMessage = async () => {
        if (!newMessage.trim() || !selectedConv) return;

        setSending(true);
        const content = newMessage;
        setNewMessage('');

        try {
            const response = await sendTicketMessage(selectedConv.id, content);

            // Add the new message to the list
            if (response.data?.message) {
                setMessages(prev => [...prev, response.data.message]);
            } else {
                // Fallback: add optimistically
                setMessages(prev => [...prev, {
                    id: 'msg-' + Date.now(),
                    sender_id: 'admin',
                    sender_type: 'agent',
                    sender_name: 'Agent Support',
                    content,
                    content_type: 'text',
                    is_read: false,
                    created_at: new Date().toISOString()
                }]);
            }

            // Refresh conversation list to update last_message
            fetchConversations(filter !== 'all' ? filter : undefined);
        } catch (error) {
            console.error('Failed to send message:', error);
            alert('Erreur lors de l\'envoi du message');
        } finally {
            setSending(false);
        }
    };

    const handleCloseConversation = async () => {
        if (!selectedConv) return;

        if (!confirm('Êtes-vous sûr de vouloir fermer cette conversation ?')) return;

        try {
            await closeTicket(selectedConv.id);
            setSelectedConv({ ...selectedConv, status: 'closed' });
            fetchConversations(filter !== 'all' ? filter : undefined);
            fetchStats();

            // Add system message
            setMessages(prev => [...prev, {
                id: 'sys-' + Date.now(),
                sender_id: 'system',
                sender_type: 'system',
                sender_name: 'Système',
                content: 'La conversation a été fermée.',
                content_type: 'text',
                is_read: true,
                created_at: new Date().toISOString()
            }]);
        } catch (error) {
            console.error('Failed to close conversation:', error);
            alert('Erreur lors de la fermeture');
        }
    };

    const getPriorityColor = (priority: string) => {
        switch (priority) {
            case 'urgent': return 'bg-red-500';
            case 'high': return 'bg-orange-500';
            case 'medium': return 'bg-yellow-500';
            default: return 'bg-green-500';
        }
    };

    const getStatusBadge = (status: string) => {
        switch (status) {
            case 'pending': return 'bg-yellow-500/20 text-yellow-400';
            case 'active': return 'bg-green-500/20 text-green-400';
            case 'escalated': return 'bg-orange-500/20 text-orange-400';
            case 'resolved': return 'bg-blue-500/20 text-blue-400';
            case 'closed': return 'bg-gray-500/20 text-gray-400';
            case 'open': return 'bg-blue-500/20 text-blue-400';
            default: return 'bg-gray-500/20 text-gray-400';
        }
    };

    const getStatusLabel = (status: string) => {
        switch (status) {
            case 'pending': return 'En attente';
            case 'active': return 'Actif';
            case 'escalated': return 'Escaladé';
            case 'resolved': return 'Résolu';
            case 'closed': return 'Fermé';
            case 'open': return 'Ouvert';
            default: return status;
        }
    };

    const formatTime = (dateString: string) => {
        const date = new Date(dateString);
        const now = new Date();
        const diffMs = now.getTime() - date.getTime();
        const diffMins = Math.floor(diffMs / 60000);
        const diffHours = Math.floor(diffMs / 3600000);
        const diffDays = Math.floor(diffMs / 86400000);

        if (diffMins < 1) return 'à l\'instant';
        if (diffMins < 60) return `il y a ${diffMins}min`;
        if (diffHours < 24) return `il y a ${diffHours}h`;
        if (diffDays < 7) return `il y a ${diffDays}j`;
        return date.toLocaleDateString('fr-FR');
    };

    const activeAgentsCount = agents.filter(a => a.is_available && a.type === 'human').length;

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between flex-wrap gap-4">
                <div>
                    <h1 className="text-2xl font-bold text-slate-900">Centre de Support</h1>
                    <p className="text-gray-500">Gérez les conversations et assistez les clients</p>
                </div>
                <div className="flex items-center gap-4">
                    <div className="flex items-center gap-2">
                        <span className={`w-2 h-2 rounded-full ${activeAgentsCount > 0 ? 'bg-green-500 animate-pulse' : 'bg-gray-400'}`}></span>
                        <span className={`text-sm ${activeAgentsCount > 0 ? 'text-green-600' : 'text-gray-500'}`}>
                            {activeAgentsCount} agent{activeAgentsCount !== 1 ? 's' : ''} en ligne
                        </span>
                    </div>
                    <button
                        onClick={() => { fetchConversations(); fetchStats(); }}
                        className="px-3 py-1.5 bg-primary-600 text-white rounded-lg text-sm hover:bg-primary-700"
                    >
                        ↻ Actualiser
                    </button>
                </div>
            </div>

            {/* Stats */}
            {stats && (
                <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
                    <div className="bg-white rounded-xl p-4 border border-gray-200 shadow-sm">
                        <p className="text-gray-500 text-sm">En attente</p>
                        <p className="text-2xl font-bold text-yellow-600">{stats.pending_conversations}</p>
                    </div>
                    <div className="bg-white rounded-xl p-4 border border-gray-200 shadow-sm">
                        <p className="text-gray-500 text-sm">Ouvertes</p>
                        <p className="text-2xl font-bold text-blue-600">{stats.open_conversations}</p>
                    </div>
                    <div className="bg-white rounded-xl p-4 border border-gray-200 shadow-sm">
                        <p className="text-gray-500 text-sm">Résolues (24h)</p>
                        <p className="text-2xl font-bold text-green-600">{stats.resolved_today}</p>
                    </div>
                    <div className="bg-white rounded-xl p-4 border border-gray-200 shadow-sm">
                        <p className="text-gray-500 text-sm">Total</p>
                        <p className="text-2xl font-bold text-slate-900">{stats.total_conversations}</p>
                    </div>
                    <div className="bg-white rounded-xl p-4 border border-gray-200 shadow-sm">
                        <p className="text-gray-500 text-sm">Satisfaction</p>
                        <p className="text-2xl font-bold text-purple-600">
                            {stats.customer_satisfaction ? `${stats.customer_satisfaction.toFixed(1)}/5` : 'N/A'} ⭐
                        </p>
                    </div>
                    <div className="bg-white rounded-xl p-4 border border-gray-200 shadow-sm">
                        <p className="text-gray-500 text-sm">Agents actifs</p>
                        <p className="text-2xl font-bold text-emerald-600">{stats.active_agents}</p>
                    </div>
                </div>
            )}

            {/* Main Content */}
            <div className="flex gap-6 h-[calc(100vh-320px)] min-h-[500px]">
                {/* Conversations List */}
                <div className="w-full lg:w-1/3 bg-white rounded-xl border border-gray-200 shadow-sm flex flex-col">
                    {/* Filters */}
                    <div className="p-4 border-b border-gray-200">
                        <div className="flex gap-2 flex-wrap">
                            {[
                                { key: 'all', label: 'Tous' },
                                { key: 'pending', label: 'En attente' },
                                { key: 'open', label: 'Ouverts' },
                                { key: 'escalated', label: '🔥 Escaladés' }
                            ].map(f => (
                                <button
                                    key={f.key}
                                    onClick={() => setFilter(f.key)}
                                    className={`px-3 py-1 rounded-full text-sm transition ${filter === f.key
                                        ? 'bg-primary-600 text-white'
                                        : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                                        }`}
                                >
                                    {f.label}
                                </button>
                            ))}
                        </div>
                    </div>

                    {/* List */}
                    <div className="flex-1 overflow-y-auto">
                        {loading ? (
                            <div className="flex items-center justify-center h-32">
                                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
                            </div>
                        ) : conversations.length === 0 ? (
                            <div className="flex items-center justify-center h-32 text-gray-500">
                                Aucune conversation
                            </div>
                        ) : (
                            conversations.map(conv => (
                                <div
                                    key={conv.id}
                                    onClick={() => selectConversation(conv)}
                                    className={`p-4 border-b border-gray-100 cursor-pointer transition hover:bg-gray-50 ${selectedConv?.id === conv.id ? 'bg-primary-50 border-l-4 border-l-primary-600' : ''
                                        }`}
                                >
                                    <div className="flex items-start gap-3">
                                        <div className={`w-2 h-2 rounded-full mt-2 flex-shrink-0 ${getPriorityColor(conv.priority)}`}></div>
                                        <div className="flex-1 min-w-0">
                                            <div className="flex items-center justify-between gap-2">
                                                <h3 className="font-medium text-slate-900 truncate">{conv.user_name}</h3>
                                                {conv.unread_count > 0 && (
                                                    <span className="bg-primary-600 text-white text-xs px-2 py-0.5 rounded-full flex-shrink-0">
                                                        {conv.unread_count}
                                                    </span>
                                                )}
                                            </div>
                                            <p className="text-sm text-gray-600 truncate">{conv.subject}</p>
                                            <p className="text-xs text-gray-400 truncate mt-1">{conv.last_message}</p>
                                            <div className="flex items-center gap-2 mt-2 flex-wrap">
                                                <span className={`px-2 py-0.5 rounded-full text-xs ${getStatusBadge(conv.status)}`}>
                                                    {getStatusLabel(conv.status)}
                                                </span>
                                                <span className="text-xs text-gray-400">
                                                    {conv.agent_type === 'ai' ? '🤖' : '👤'}
                                                </span>
                                                <span className="text-xs text-gray-400">
                                                    {formatTime(conv.updated_at || conv.created_at)}
                                                </span>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            ))
                        )}
                    </div>
                </div>

                {/* Chat Area */}
                <div className="hidden lg:flex flex-1 bg-white rounded-xl border border-gray-200 shadow-sm flex-col">
                    {selectedConv ? (
                        <>
                            {/* Chat Header */}
                            <div className="p-4 border-b border-gray-200 flex items-center justify-between">
                                <div>
                                    <h2 className="font-semibold text-slate-900">{selectedConv.user_name}</h2>
                                    <p className="text-sm text-gray-500">
                                        {selectedConv.user_email} • {selectedConv.subject}
                                    </p>
                                </div>
                                <div className="flex items-center gap-2">
                                    <span className={`px-3 py-1 rounded-full text-xs ${getStatusBadge(selectedConv.status)}`}>
                                        {getStatusLabel(selectedConv.status)}
                                    </span>
                                </div>
                            </div>

                            {/* Messages */}
                            <div className="flex-1 overflow-y-auto p-4 space-y-4 bg-gray-50">
                                {loadingMessages ? (
                                    <div className="flex items-center justify-center h-full">
                                        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
                                    </div>
                                ) : messages.length === 0 ? (
                                    <div className="flex items-center justify-center h-full text-gray-500">
                                        Aucun message
                                    </div>
                                ) : (
                                    messages.map(msg => (
                                        <div
                                            key={msg.id}
                                            className={`flex ${msg.sender_type === 'agent' ? 'justify-end' : msg.sender_type === 'system' ? 'justify-center' : 'justify-start'}`}
                                        >
                                            {msg.sender_type === 'system' ? (
                                                <div className="bg-gray-200 text-gray-600 text-sm px-4 py-2 rounded-full">
                                                    {msg.content}
                                                </div>
                                            ) : (
                                                <div className={`max-w-[70%]`}>
                                                    <div className={`rounded-xl p-3 ${msg.sender_type === 'agent'
                                                        ? 'bg-primary-600 text-white rounded-tr-none'
                                                        : 'bg-white border border-gray-200 text-slate-900 rounded-tl-none shadow-sm'
                                                        }`}>
                                                        <p className={`text-xs font-medium mb-1 ${msg.sender_type === 'agent' ? 'text-primary-100' : 'text-gray-500'}`}>
                                                            {msg.sender_name}
                                                        </p>
                                                        <p className="whitespace-pre-wrap">{msg.content}</p>
                                                    </div>
                                                    <p className={`text-xs text-gray-400 mt-1 ${msg.sender_type === 'agent' ? 'text-right' : ''}`}>
                                                        {new Date(msg.created_at).toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' })}
                                                    </p>
                                                </div>
                                            )}
                                        </div>
                                    ))
                                )}
                                <div ref={messagesEndRef} />
                            </div>

                            {/* Actions */}
                            {selectedConv.status !== 'closed' && selectedConv.status !== 'resolved' && (
                                <div className="p-3 border-t border-gray-200 flex gap-2 flex-wrap bg-gray-50">
                                    <button
                                        onClick={handleCloseConversation}
                                        className="px-3 py-1.5 bg-green-100 text-green-700 rounded-lg text-sm hover:bg-green-200"
                                    >
                                        ✓ Résoudre
                                    </button>
                                </div>
                            )}

                            {/* Input */}
                            {selectedConv.status !== 'closed' && selectedConv.status !== 'resolved' && (
                                <div className="p-4 border-t border-gray-200">
                                    <div className="flex gap-3">
                                        <input
                                            type="text"
                                            value={newMessage}
                                            onChange={(e) => setNewMessage(e.target.value)}
                                            onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
                                            placeholder="Tapez votre réponse..."
                                            className="flex-1 bg-gray-100 text-slate-900 rounded-xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-primary-500"
                                            disabled={sending}
                                        />
                                        <button
                                            onClick={sendMessage}
                                            disabled={!newMessage.trim() || sending}
                                            className="px-6 py-3 bg-primary-600 text-white rounded-xl hover:bg-primary-700 transition disabled:opacity-50"
                                        >
                                            {sending ? '...' : 'Envoyer'}
                                        </button>
                                    </div>
                                </div>
                            )}

                            {(selectedConv.status === 'closed' || selectedConv.status === 'resolved') && (
                                <div className="p-4 border-t border-gray-200 bg-gray-50 text-center text-gray-500">
                                    Cette conversation est fermée
                                </div>
                            )}
                        </>
                    ) : (
                        <div className="flex-1 flex items-center justify-center text-gray-500">
                            <div className="text-center">
                                <p className="text-4xl mb-4">💬</p>
                                <p>Sélectionnez une conversation</p>
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
