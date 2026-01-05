import React, { useState, useEffect } from 'react';
import { View, Text, FlatList, TouchableOpacity, StyleSheet, RefreshControl } from 'react-native';
import api from '../api';

const NotificationsScreen = () => {
    const [notifications, setNotifications] = useState([]);
    const [refreshing, setRefreshing] = useState(false);

    const loadNotifications = async () => {
        try {
            const res = await api.get('/notification-service/api/v1/notifications');
            setNotifications(res.data || []);
        } catch (err) {
            console.error(err);
        }
    };

    useEffect(() => {
        loadNotifications();
    }, []);

    const onRefresh = async () => {
        setRefreshing(true);
        await loadNotifications();
        setRefreshing(false);
    };

    const markAsRead = async (id: string) => {
        try {
            await api.post(`/notification-service/api/v1/notifications/${id}/read`);
            loadNotifications();
        } catch (err) {
            console.error(err);
        }
    };

    return (
        <View style={styles.container}>
            <View style={styles.header}>
                <Text style={styles.title}>Notifications</Text>
            </View>

            <FlatList
                data={notifications}
                keyExtractor={(item: any) => item.id}
                refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
                ListEmptyComponent={
                    <View style={styles.empty}>
                        <Text style={styles.emptyIcon}>🔔</Text>
                        <Text style={styles.emptyText}>Aucune notification</Text>
                    </View>
                }
                renderItem={({ item }: any) => (
                    <TouchableOpacity
                        style={[styles.notificationCard, item.read && styles.readCard]}
                        onPress={() => markAsRead(item.id)}
                    >
                        <Text style={styles.notificationTitle}>{item.title}</Text>
                        <Text style={styles.notificationMessage}>{item.message}</Text>
                        <Text style={styles.time}>{new Date(item.created_at).toLocaleString()}</Text>
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
    emptyText: { fontSize: 18, fontWeight: '600', color: '#666' },
    notificationCard: {
        padding: 16,
        backgroundColor: '#fff',
        borderBottomWidth: 1,
        borderBottomColor: '#e5e5e5',
    },
    readCard: { opacity: 0.6 },
    notificationTitle: { fontSize: 16, fontWeight: '600', marginBottom: 4 },
    notificationMessage: { fontSize: 14, color: '#666', marginBottom: 8 },
    time: { fontSize: 12, color: '#999' },
});

export default NotificationsScreen;
