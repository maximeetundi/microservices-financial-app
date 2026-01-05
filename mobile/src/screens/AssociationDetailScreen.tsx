import React, { useEffect, useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, StyleSheet, RefreshControl } from 'react-native';
import { associationAPI } from '../api';

const AssociationDetailScreen = ({ route, navigation }: any) => {
    const { id } = route.params;
    const [association, setAssociation] = useState<any>(null);
    const [activeTab, setActiveTab] = useState('members');
    const [members, setMembers] = useState([]);
    const [messages, setMessages] = useState([]);
    const [solidarity, setSolidarity] = useState([]);
    const [rounds, setRounds] = useState([]);
    const [approvals, setApprovals] = useState([]);
    const [refreshing, setRefreshing] = useState(false);

    const loadData = async () => {
        try {
            const res = await associationAPI.get(id);
            setAssociation(res.data);

            if (activeTab === 'members') {
                const membersRes = await associationAPI.getMembers(id);
                setMembers(membersRes.data);
            } else if (activeTab === 'chat') {
                const chatRes = await associationAPI.getMessages(id);
                setMessages(chatRes.data);
            } else if (activeTab === 'solidarity') {
                const solidarityRes = await associationAPI.getSolidarityEvents(id);
                setSolidarity(solidarityRes.data);
            } else if (activeTab === 'tontine') {
                const roundsRes = await associationAPI.getCalledRounds(id);
                setRounds(roundsRes.data);
            } else if (activeTab === 'approvals') {
                const approvalsRes = await associationAPI.getPendingApprovals(id);
                setApprovals(approvalsRes.data);
            }
        } catch (err) {
            console.error(err);
        }
    };

    useEffect(() => {
        loadData();
    }, [id, activeTab]);

    const onRefresh = async () => {
        setRefreshing(true);
        await loadData();
        setRefreshing(false);
    };

    const tabs = [
        { id: 'members', label: 'Membres' },
        { id: 'chat', label: 'Chat' },
        { id: 'solidarity', label: 'Solidarité' },
        { id: 'tontine', label: 'Tontine' },
        { id: 'approvals', label: 'Votes' },
    ];

    return (
        <View style={styles.container}>
            {/* Header */}
            <View style={styles.header}>
                <TouchableOpacity onPress={() => navigation.goBack()}>
                    <Text style={styles.backButton}>← Retour</Text>
                </TouchableOpacity>
                <View>
                    <Text style={styles.title}>{association?.name}</Text>
                    <Text style={styles.subtitle}>{association?.type}</Text>
                </View>
                <TouchableOpacity onPress={() => navigation.navigate('AssociationSettings', { id })}>
                    <Text style={styles.settingsButton}>⚙️</Text>
                </TouchableOpacity>
            </View>

            {/* Stats */}
            <View style={styles.stats}>
                <View style={styles.statItem}>
                    <Text style={styles.statValue}>{association?.total_members || 0}</Text>
                    <Text style={styles.statLabel}>Membres</Text>
                </View>
                <View style={styles.statItem}>
                    <Text style={styles.statValue}>{association?.treasury_balance || 0} XOF</Text>
                    <Text style={styles.statLabel}>Caisse</Text>
                </View>
            </View>

            {/* Tabs */}
            <ScrollView horizontal showsHorizontalScrollIndicator={false} style={styles.tabsContainer}>
                {tabs.map(tab => (
                    <TouchableOpacity
                        key={tab.id}
                        onPress={() => setActiveTab(tab.id)}
                        style={[styles.tab, activeTab === tab.id && styles.activeTab]}
                    >
                        <Text style={[styles.tabText, activeTab === tab.id && styles.activeTabText]}>
                            {tab.label}
                        </Text>
                    </TouchableOpacity>
                ))}
            </ScrollView>

            {/* Content */}
            <ScrollView
                style={styles.content}
                refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
            >
                {activeTab === 'members' && (
                    <View>
                        {members.map((member: any) => (
                            <View key={member.id} style={styles.memberCard}>
                                <View>
                                    <Text style={styles.memberName}>{member.user_id}</Text>
                                    <Text style={styles.memberRole}>{member.role}</Text>
                                </View>
                                <Text style={styles.memberAmount}>{member.contributions_paid} XOF</Text>
                            </View>
                        ))}
                    </View>
                )}

                {activeTab === 'chat' && (
                    <View>
                        {messages.map((msg: any) => (
                            <View key={msg.id} style={styles.messageCard}>
                                <Text style={styles.messageSender}>{msg.sender_name}</Text>
                                <Text style={styles.messageContent}>{msg.content}</Text>
                                <Text style={styles.messageTime}>{new Date(msg.created_at).toLocaleTimeString()}</Text>
                            </View>
                        ))}
                    </View>
                )}

                {activeTab === 'solidarity' && (
                    <View>
                        {solidarity.map((event: any) => (
                            <View key={event.id} style={styles.solidarityCard}>
                                <Text style={styles.solidarityTitle}>{event.title}</Text>
                                <Text style={styles.solidarityType}>{event.event_type}</Text>
                                <View style={styles.progressBar}>
                                    <View style={[styles.progress, { width: `${(event.collected_amount / event.target_amount) * 100}%` }]} />
                                </View>
                                <Text>{event.collected_amount} / {event.target_amount} XOF</Text>
                            </View>
                        ))}
                    </View>
                )}

                {activeTab === 'tontine' && (
                    <View>
                        {rounds.map((round: any) => (
                            <View key={round.id} style={styles.roundCard}>
                                <Text style={styles.roundNumber}>Tour #{round.round_number}</Text>
                                <Text>Bénéficiaire: {round.beneficiary_name}</Text>
                                <Text style={styles.roundAmount}>{round.total_collected} XOF</Text>
                            </View>
                        ))}
                    </View>
                )}

                {activeTab === 'approvals' && (
                    <View>
                        {approvals.map((approval: any) => (
                            <View key={approval.id} style={styles.approvalCard}>
                                <Text style={styles.approvalType}>{approval.request_type}</Text>
                                <Text style={styles.approvalAmount}>{approval.amount} XOF</Text>
                                <Text>{approval.current_approvals}/{approval.required_approvals} votes</Text>
                                <View style={styles.voteButtons}>
                                    <TouchableOpacity style={styles.approveButton}>
                                        <Text style={styles.buttonText}>Approuver</Text>
                                    </TouchableOpacity>
                                    <TouchableOpacity style={styles.rejectButton}>
                                        <Text style={styles.buttonText}>Rejeter</Text>
                                    </TouchableOpacity>
                                </View>
                            </View>
                        ))}
                    </View>
                )}
            </ScrollView>
        </View>
    );
};

const styles = StyleSheet.create({
    container: { flex: 1, backgroundColor: '#f5f5f5' },
    header: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', padding: 16, backgroundColor: '#4f46e5' },
    backButton: { color: '#fff', fontSize: 16 },
    title: { fontSize: 20, fontWeight: 'bold', color: '#fff' },
    subtitle: { color: '#e0e7ff', fontSize: 14 },
    settingsButton: { fontSize: 24 },
    stats: { flexDirection: 'row', backgroundColor: '#fff', padding: 16 },
    statItem: { flex: 1, alignItems: 'center' },
    statValue: { fontSize: 24, fontWeight: 'bold', color: '#4f46e5' },
    statLabel: { fontSize: 12, color: '#666' },
    tabsContainer: { backgroundColor: '#fff', borderBottomWidth: 1, borderBottomColor: '#e5e5e5' },
    tab: { paddingVertical: 12, paddingHorizontal: 20 },
    activeTab: { borderBottomWidth: 2, borderBottomColor: '#4f46e5' },
    tabText: { color: '#666', fontSize: 14 },
    activeTabText: { color: '#4f46e5', fontWeight: '600' },
    content: { flex: 1, padding: 16 },
    memberCard: { flexDirection: 'row', justifyContent: 'space-between', padding: 16, backgroundColor: '#fff', borderRadius: 8, marginBottom: 8 },
    memberName: { fontSize: 16, fontWeight: '600' },
    memberRole: { color: '#666', fontSize: 12 },
    memberAmount: { fontSize: 16, fontWeight: 'bold', color: '#4f46e5' },
    messageCard: { padding: 12, backgroundColor: '#fff', borderRadius: 8, marginBottom: 8 },
    messageSender: { fontWeight: '600', marginBottom: 4 },
    messageContent: { fontSize: 14 },
    messageTime: { fontSize: 12, color: '#999', marginTop: 4 },
    solidarityCard: { padding: 16, backgroundColor: '#fff', borderRadius: 8, marginBottom: 8 },
    solidarityTitle: { fontSize: 16, fontWeight: '600', marginBottom: 4 },
    solidarityType: { color: '#666', fontSize: 12, marginBottom: 8 },
    progressBar: { height: 8, backgroundColor: '#e5e5e5', borderRadius: 4, marginBottom: 8 },
    progress: { height: '100%', backgroundColor: '#4f46e5', borderRadius: 4 },
    roundCard: { padding: 16, backgroundColor: '#fff', borderRadius: 8, marginBottom: 8 },
    roundNumber: { fontSize: 16, fontWeight: '600', marginBottom: 4 },
    roundAmount: { fontSize: 20, fontWeight: 'bold', color: '#4f46e5', marginTop: 8 },
    approvalCard: { padding: 16, backgroundColor: '#fff', borderRadius: 8, marginBottom: 8 },
    approvalType: { fontSize: 16, fontWeight: '600' },
    approvalAmount: { fontSize: 24, fontWeight: 'bold', color: '#4f46e5', marginVertical: 8 },
    voteButtons: { flexDirection: 'row', gap: 8, marginTop: 12 },
    approveButton: { flex: 1, backgroundColor: '#10b981', padding: 12, borderRadius: 8, alignItems: 'center' },
    rejectButton: { flex: 1, backgroundColor: '#ef4444', padding: 12, borderRadius: 8, alignItems: 'center' },
    buttonText: { color: '#fff', fontWeight: '600' },
});

export default AssociationDetailScreen;
