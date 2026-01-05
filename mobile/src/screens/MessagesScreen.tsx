import React, { useState, useEffect } from 'react';
import { View, Text, FlatList, TouchableOpacity, StyleSheet, RefreshControl } from 'react-native';
import api from '../api';

const MessagesScreen = ({ navigation }: any) => {
    const [conversations, setConversations] = useState([]);
    const [refreshing, setRefreshing] = useState(false);

    const loadConversations = async () => {
        try {
            const res = await api.get('/messaging-service/api/v1/conversations');
            setConversations(res.data.conversations || []);
        } catch (err) {
            console.error(err);
        }
    };

    useEffect(() => {
        loadConversations();
    }, []);

    const onRefresh = async () => {
        setRefreshing(true);
        await loadConversations();
        setRefreshing(false);
    };

    return (
        <View style={styles.container}>
            <View style={styles.header}>
                <Text style={styles.title}>Messages</Text>
            </View>

            <FlatList
                data={conversations}
                keyExtractor={(item: any) => item.id}
                refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
                ListEmptyComponent={
                    <View style={styles.empty}>
                        <Text style={styles.emptyIcon}>💬</Text>
                        <Text style={styles.emptyText}>Aucune conversation</Text>
                        <Text style={styles.emptySubtext}>Messagerie utilisateur-à-utilisateur</Text>
                    </View>
                }
                renderItem={({ item }: any) => (
                    <TouchableOpacity style={styles.conversationCard}>
                        <View style={styles.avatar}>
                            <Text>{item.participant_name?.[0]}</Text>
                        </View>
                        <View style={{ flex: 1 }}>
                            <Text style={styles.name}>{item.participant_name}</Text>
                            <Text style={styles.lastMessage} numberOfLines={1}>
                                {item.last_message}
                            </Text>
                        </View>
                        <Text style={styles.time}>
                            {item.last_message_at ? new Date(item.last_message_at).toLocaleDateString() : ''}
                        </Text>
                    </TouchableOpacity>
                )}
            />
        </View>
    );
};

const styles = StyleSheet.create({
    container: { flex: 1, backgroundColor: '#f5f5f5' },
    header: { backgroundColor: '#4f46e5', padding: 16 },
    title: { fontSize: 24, fontWeight: 'bold', color: '#fff' },
    empty: { padding: 40, alignItems: 'center' },
    emptyIcon: { fontSize: 64, marginBottom: 16 },
    emptyText: { fontSize: 18, fontWeight: '600', color: '#666', marginBottom: 8 },
    emptySubtext: { fontSize: 14, color: '#999' },
    conversationCard: {
        flexDirection: 'row',
        alignItems: 'center',
        padding: 16,
        backgroundColor: '#fff',
        borderBottomWidth: 1,
        borderBottomColor: '#e5e5e5',
    },
    avatar: {
        width: 48,
        height: 48,
        borderRadius: 24,
        backgroundColor: '#4f46e5',
        alignItems: 'center',
        justifyContent: 'center',
        marginRight: 12,
    },
    name: { fontSize: 16, fontWeight: '600', marginBottom: 4 },
    lastMessage: { fontSize: 14, color: '#666' },
    time: { fontSize: 12, color: '#999' },
});

export default MessagesScreen;
