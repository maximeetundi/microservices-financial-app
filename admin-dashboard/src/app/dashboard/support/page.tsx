'use client';

import { useState, useEffect, useRef } from 'react';

interface Conversation {
    id: string;
    user_id: string;
    user_name: string;
    user_email: string;
    agent_type: 'ai' | 'human';
    subject: string;
    category: string;
    status: string;
    priority: string;
    last_message: string;
    unread_count: number;
    message_count: number;
    created_at: string;
    updated_at: string;
}

interface Message {
    id: string;
    sender_type: 'user' | 'agent' | 'system';
    sender_name: string;
    content: string;
    created_at: string;
}

interface Stats {
    total_conversations: number;
    open_conversations: number;
    resolved_today: number;
    pending_conversations: number;
    customer_satisfaction: number;
    active_agents: number;
}

export default function SupportPage() {
    const [conversations, setConversations] = useState<Conversation[]>([]);
    const [selectedConv, setSelectedConv] = useState<Conversation | null>(null);
    const [messages, setMessages] = useState<Message[]>([]);
    const [newMessage, setNewMessage] = useState('');
    const [filter, setFilter] = useState('all');
    const [stats, setStats] = useState<Stats | null>(null);
    const [sending, setSending] = useState(false);
    const messagesEndRef = useRef<HTMLDivElement>(null);

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    };

    useEffect(() => {
        // Demo data
        setStats({
            total_conversations: 156,
            open_conversations: 23,
            resolved_today: 45,
            pending_conversations: 8,
            customer_satisfaction: 4.8,
            active_agents: 5
        });

        setConversations([
            {
                id: '1',
                user_id: 'u1',
                user_name: 'Jean Dupont',
                user_email: 'jean@example.com',
                agent_type: 'human',
                subject: 'Problème de transfert international',
                category: 'transfer',
                status: 'pending',
                priority: 'high',
                last_message: 'Mon transfert est bloqué depuis 2 jours...',
                unread_count: 3,
                message_count: 5,
                created_at: new Date(Date.now() - 3600000).toISOString(),
                updated_at: new Date(Date.now() - 600000).toISOString()
            },
            {
                id: '2',
                user_id: 'u2',
                user_name: 'Marie Martin',
                user_email: 'marie@example.com',
                agent_type: 'ai',
                subject: 'Question sur les frais',
                category: 'fees',
                status: 'escalated',
                priority: 'medium',
                last_message: 'Je voudrais comprendre pourquoi j\'ai été facturé...',
                unread_count: 1,
                message_count: 8,
                created_at: new Date(Date.now() - 7200000).toISOString(),
                updated_at: new Date(Date.now() - 1800000).toISOString()
            },
            {
                id: '3',
                user_id: 'u3',
                user_name: 'Pierre Leroy',
                user_email: 'pierre@example.com',
                agent_type: 'human',
                subject: 'Carte volée - Demande de blocage',
                category: 'security',
                status: 'active',
                priority: 'urgent',
                last_message: 'Ma carte a été volée, pouvez-vous la bloquer?',
                unread_count: 0,
                message_count: 12,
                created_at: new Date(Date.now() - 1800000).toISOString(),
                updated_at: new Date(Date.now() - 300000).toISOString()
            }
        ]);
    }, []);

    useEffect(() => {
        scrollToBottom();
    }, [messages]);

    const selectConversation = (conv: Conversation) => {
        setSelectedConv(conv);
        // Load messages for this conversation
        setMessages([
            {
                id: 'm1',
                sender_type: 'user',
                sender_name: conv.user_name,
                content: `Bonjour, j'ai un problème avec mon compte concernant: ${conv.subject}`,
                created_at: new Date(Date.now() - 3600000).toISOString()
            },
            {
                id: 'm2',
                sender_type: 'agent',
                sender_name: 'Assistant IA',
                content: 'Bonjour ! Je comprends votre préoccupation. Pouvez-vous me donner plus de détails ?',
                created_at: new Date(Date.now() - 3500000).toISOString()
            },
            {
                id: 'm3',
                sender_type: 'user',
                sender_name: conv.user_name,
                content: conv.last_message,
                created_at: new Date(Date.now() - 600000).toISOString()
            }
        ]);
    };

    const sendMessage = async () => {
        if (!newMessage.trim() || !selectedConv) return;

        setSending(true);
        const content = newMessage;
        setNewMessage('');

        // Add message immediately
        setMessages(prev => [...prev, {
            id: 'msg-' + Date.now(),
            sender_type: 'agent',
            sender_name: 'Agent Support',
            content,
            created_at: new Date().toISOString()
        }]);

        // Simulate API call
        await new Promise(resolve => setTimeout(resolve, 500));
        setSending(false);
    };

    const handleAction = (action: string) => {
        if (!selectedConv) return;

        let systemMessage = '';
        switch (action) {
            case 'resolve':
                systemMessage = 'La conversation a été marquée comme résolue.';
                setSelectedConv({ ...selectedConv, status: 'resolved' });
                break;
            case 'escalate':
                systemMessage = 'La conversation a été escaladée au niveau supérieur.';
                break;
            case 'block-card':
                systemMessage = '🔒 La carte du client a été bloquée avec succès.';
                break;
            case 'refund':
                systemMessage = '💰 Un remboursement a été initié pour le client.';
                break;
        }

        setMessages(prev => [...prev, {
            id: 'sys-' + Date.now(),
            sender_type: 'system',
            sender_name: 'Système',
            content: systemMessage,
            created_at: new Date().toISOString()
        }]);
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
            default: return 'bg-gray-500/20 text-gray-400';
        }
    };

    const filteredConversations = conversations.filter(c => {
        if (filter === 'all') return true;
        if (filter === 'pending') return c.status === 'pending' || c.status === 'escalated';
        if (filter === 'active') return c.status === 'active';
        if (filter === 'urgent') return c.priority === 'urgent' || c.priority === 'high';
        return true;
    });

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-bold text-white">Centre de Support</h1>
                    <p className="text-gray-400">Gérez les conversations et assistez les clients</p>
                </div>
                <div className="flex items-center gap-4">
                    <div className="flex items-center gap-2">
                        <span className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
                        <span className="text-green-400 text-sm">5 agents en ligne</span>
                    </div>
                </div>
            </div>

            {/* Stats */}
            {stats && (
                <div className="grid grid-cols-1 md:grid-cols-6 gap-4">
                    <div className="bg-slate-800/50 rounded-xl p-4 border border-slate-700">
                        <p className="text-gray-400 text-sm">En attente</p>
                        <p className="text-2xl font-bold text-yellow-400">{stats.pending_conversations}</p>
                    </div>
                    <div className="bg-slate-800/50 rounded-xl p-4 border border-slate-700">
                        <p className="text-gray-400 text-sm">Ouvertes</p>
                        <p className="text-2xl font-bold text-blue-400">{stats.open_conversations}</p>
                    </div>
                    <div className="bg-slate-800/50 rounded-xl p-4 border border-slate-700">
                        <p className="text-gray-400 text-sm">Résolues (24h)</p>
                        <p className="text-2xl font-bold text-green-400">{stats.resolved_today}</p>
                    </div>
                    <div className="bg-slate-800/50 rounded-xl p-4 border border-slate-700">
                        <p className="text-gray-400 text-sm">Total</p>
                        <p className="text-2xl font-bold text-white">{stats.total_conversations}</p>
                    </div>
                    <div className="bg-slate-800/50 rounded-xl p-4 border border-slate-700">
                        <p className="text-gray-400 text-sm">Satisfaction</p>
                        <p className="text-2xl font-bold text-purple-400">{stats.customer_satisfaction}/5 ⭐</p>
                    </div>
                    <div className="bg-slate-800/50 rounded-xl p-4 border border-slate-700">
                        <p className="text-gray-400 text-sm">Agents actifs</p>
                        <p className="text-2xl font-bold text-emerald-400">{stats.active_agents}</p>
                    </div>
                </div>
            )}

            {/* Main Content */}
            <div className="flex gap-6 h-[calc(100vh-300px)] min-h-[500px]">
                {/* Conversations List */}
                <div className="w-1/3 bg-slate-800/50 rounded-xl border border-slate-700 flex flex-col">
                    {/* Filters */}
                    <div className="p-4 border-b border-slate-700">
                        <div className="flex gap-2 flex-wrap">
                            {['all', 'pending', 'active', 'urgent'].map(f => (
                                <button
                                    key={f}
                                    onClick={() => setFilter(f)}
                                    className={`px-3 py-1 rounded-full text-sm transition ${filter === f
                                            ? 'bg-blue-500 text-white'
                                            : 'bg-slate-700 text-gray-300 hover:bg-slate-600'
                                        }`}
                                >
                                    {f === 'all' && 'Tous'}
                                    {f === 'pending' && 'En attente'}
                                    {f === 'active' && 'Actifs'}
                                    {f === 'urgent' && '🔥 Urgents'}
                                </button>
                            ))}
                        </div>
                    </div>

                    {/* List */}
                    <div className="flex-1 overflow-y-auto">
                        {filteredConversations.map(conv => (
                            <div
                                key={conv.id}
                                onClick={() => selectConversation(conv)}
                                className={`p-4 border-b border-slate-700 cursor-pointer transition hover:bg-slate-700/50 ${selectedConv?.id === conv.id ? 'bg-slate-700/50 border-l-2 border-l-blue-500' : ''
                                    }`}
                            >
                                <div className="flex items-start gap-3">
                                    <div className={`w-2 h-2 rounded-full mt-2 ${getPriorityColor(conv.priority)}`}></div>
                                    <div className="flex-1 min-w-0">
                                        <div className="flex items-center justify-between">
                                            <h3 className="font-medium text-white truncate">{conv.user_name}</h3>
                                            {conv.unread_count > 0 && (
                                                <span className="bg-blue-500 text-white text-xs px-2 py-0.5 rounded-full">
                                                    {conv.unread_count}
                                                </span>
                                            )}
                                        </div>
                                        <p className="text-sm text-gray-300 truncate">{conv.subject}</p>
                                        <p className="text-xs text-gray-500 truncate mt-1">{conv.last_message}</p>
                                        <div className="flex items-center gap-2 mt-2">
                                            <span className={`px-2 py-0.5 rounded-full text-xs ${getStatusBadge(conv.status)}`}>
                                                {conv.status}
                                            </span>
                                            <span className="text-xs text-gray-500">
                                                {conv.agent_type === 'ai' ? '🤖' : '👤'}
                                            </span>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>

                {/* Chat Area */}
                <div className="flex-1 bg-slate-800/50 rounded-xl border border-slate-700 flex flex-col">
                    {selectedConv ? (
                        <>
                            {/* Chat Header */}
                            <div className="p-4 border-b border-slate-700 flex items-center justify-between">
                                <div>
                                    <h2 className="font-semibold text-white">{selectedConv.user_name}</h2>
                                    <p className="text-sm text-gray-400">{selectedConv.user_email} • {selectedConv.subject}</p>
                                </div>
                                <div className="flex items-center gap-2">
                                    <span className={`px-3 py-1 rounded-full text-xs ${getStatusBadge(selectedConv.status)}`}>
                                        {selectedConv.status}
                                    </span>
                                </div>
                            </div>

                            {/* Messages */}
                            <div className="flex-1 overflow-y-auto p-4 space-y-4">
                                {messages.map(msg => (
                                    <div
                                        key={msg.id}
                                        className={`flex ${msg.sender_type === 'agent' ? 'justify-end' : msg.sender_type === 'system' ? 'justify-center' : 'justify-start'}`}
                                    >
                                        {msg.sender_type === 'system' ? (
                                            <div className="bg-slate-700 text-gray-300 text-sm px-4 py-2 rounded-full">
                                                {msg.content}
                                            </div>
                                        ) : (
                                            <div className={`max-w-[70%] ${msg.sender_type === 'agent' ? 'order-2' : ''}`}>
                                                <div className={`rounded-xl p-3 ${msg.sender_type === 'agent'
                                                        ? 'bg-blue-500 text-white rounded-tr-none'
                                                        : 'bg-slate-700 text-white rounded-tl-none'
                                                    }`}>
                                                    <p className="text-sm font-medium mb-1">{msg.sender_name}</p>
                                                    <p className="whitespace-pre-wrap">{msg.content}</p>
                                                </div>
                                                <p className={`text-xs text-gray-500 mt-1 ${msg.sender_type === 'agent' ? 'text-right' : ''}`}>
                                                    {new Date(msg.created_at).toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' })}
                                                </p>
                                            </div>
                                        )}
                                    </div>
                                ))}
                                <div ref={messagesEndRef} />
                            </div>

                            {/* Actions */}
                            <div className="p-3 border-t border-slate-700 flex gap-2 flex-wrap">
                                <button onClick={() => handleAction('resolve')} className="px-3 py-1.5 bg-green-500/20 text-green-400 rounded-lg text-sm hover:bg-green-500/30">
                                    ✓ Résoudre
                                </button>
                                <button onClick={() => handleAction('escalate')} className="px-3 py-1.5 bg-orange-500/20 text-orange-400 rounded-lg text-sm hover:bg-orange-500/30">
                                    ⬆ Escalader
                                </button>
                                <button onClick={() => handleAction('block-card')} className="px-3 py-1.5 bg-red-500/20 text-red-400 rounded-lg text-sm hover:bg-red-500/30">
                                    🔒 Bloquer carte
                                </button>
                                <button onClick={() => handleAction('refund')} className="px-3 py-1.5 bg-purple-500/20 text-purple-400 rounded-lg text-sm hover:bg-purple-500/30">
                                    💰 Rembourser
                                </button>
                            </div>

                            {/* Input */}
                            <div className="p-4 border-t border-slate-700">
                                <div className="flex gap-3">
                                    <input
                                        type="text"
                                        value={newMessage}
                                        onChange={(e) => setNewMessage(e.target.value)}
                                        onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
                                        placeholder="Tapez votre réponse..."
                                        className="flex-1 bg-slate-700 text-white rounded-xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                        disabled={sending}
                                    />
                                    <button
                                        onClick={sendMessage}
                                        disabled={!newMessage.trim() || sending}
                                        className="px-6 py-3 bg-blue-500 text-white rounded-xl hover:bg-blue-600 transition disabled:opacity-50"
                                    >
                                        Envoyer
                                    </button>
                                </div>
                            </div>
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
