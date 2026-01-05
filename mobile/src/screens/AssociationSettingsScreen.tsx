import React, { useState, useEffect } from 'react';
import { View, Text, ScrollView, TouchableOpacity, StyleSheet, TextInput, Switch } from 'react-native';
import { associationAPI } from '../api';

const AssociationSettingsScreen = ({ route, navigation }: any) => {
    const { id } = route.params;
    const [roles, setRoles] = useState([]);
    const [approvers, setApprovers] = useState([]);
    const [members, setMembers] = useState([]);
    const [selectedApprovers, setSelectedApprovers] = useState<string[]>([]);
    const [activeTab, setActiveTab] = useState('roles');

    const [newRole, setNewRole] = useState({ name: '', permissions: [] });
    const [settings, setSettings] = useState({ late_fee_amount: 0, food_fee: 2000, drink_fee: 2000 });
    const [chatEnabled, setChatEnabled] = useState(true);

    useEffect(() => {
        loadData();
    }, [id, activeTab]);

    const loadData = async () => {
        try {
            if (activeTab === 'roles') {
                const res = await associationAPI.getRoles(id);
                setRoles(res.data);
            } else if (activeTab === 'approvers') {
                const [approversRes, membersRes] = await Promise.all([
                    associationAPI.getApprovers(id),
                    associationAPI.getMembers(id),
                ]);
                setApprovers(approversRes.data);
                setMembers(membersRes.data);
                setSelectedApprovers(approversRes.data.map((a: any) => a.member_id));
            }
        } catch (err) {
            console.error(err);
        }
    };

    const saveApprovers = async () => {
        if (selectedApprovers.length !== 5) {
            alert('Vous devez sélectionner exactement 5 approuveurs');
            return;
        }
        try {
            await associationAPI.setApprovers(id, selectedApprovers);
            alert('Approuveurs enregistrés !');
        } catch (err: any) {
            alert('Erreur: ' + err.response?.data?.error);
        }
    };

    const toggleApprover = (memberId: string) => {
        if (selectedApprovers.includes(memberId)) {
            setSelectedApprovers(selectedApprovers.filter(id => id !== memberId));
        } else if (selectedApprovers.length < 5) {
            setSelectedApprovers([...selectedApprovers, memberId]);
        }
    };

    return (
        <View style={styles.container}>
            <View style={styles.header}>
                <TouchableOpacity onPress={() => navigation.goBack()}>
                    <Text style={styles.backButton}>← Retour</Text>
                </TouchableOpacity>
                <Text style={styles.title}>Paramètres</Text>
            </View>

            {/* Tabs */}
            <View style={styles.tabsContainer}>
                {['roles', 'approvers', 'general'].map(tab => (
                    <TouchableOpacity
                        key={tab}
                        onPress={() => setActiveTab(tab)}
                        style={[styles.tab, activeTab === tab && styles.activeTab]}
                    >
                        <Text style={[styles.tabText, activeTab === tab && styles.activeTabText]}>
                            {tab === 'roles' ? 'Rôles' : tab === 'approvers' ? 'Approuveurs' : 'Général'}
                        </Text>
                    </TouchableOpacity>
                ))}
            </View>

            <ScrollView style={styles.content}>
                {activeTab === 'roles' && (
                    <View>
                        <Text style={styles.sectionTitle}>Rôles personnalisés</Text>
                        {roles.map((role: any) => (
                            <View key={role.id} style={styles.roleCard}>
                                <Text style={styles.roleName}>{role.name}</Text>
                                <Text style={styles.rolePermissions}>{role.permissions.length} permissions</Text>
                            </View>
                        ))}
                    </View>
                )}

                {activeTab === 'approvers' && (
                    <View>
                        <Text style={styles.sectionTitle}>Sélectionnez 5 approuveurs</Text>
                        <Text style={styles.description}>
                            Ces membres pourront voter sur les prêts et distributions (4/5 requis)
                        </Text>
                        {members.map((member: any) => (
                            <TouchableOpacity
                                key={member.id}
                                onPress={() => toggleApprover(member.id)}
                                style={styles.memberCard}
                            >
                                <View style={styles.checkbox}>
                                    {selectedApprovers.includes(member.id) && <Text>✓</Text>}
                                </View>
                                <View style={{ flex: 1 }}>
                                    <Text style={styles.memberName}>{member.user_id}</Text>
                                    <Text style={styles.memberRole}>{member.role}</Text>
                                </View>
                            </TouchableOpacity>
                        ))}
                        <TouchableOpacity
                            onPress={saveApprovers}
                            style={[styles.saveButton, selectedApprovers.length !== 5 && styles.disabledButton]}
                            disabled={selectedApprovers.length !== 5}
                        >
                            <Text style={styles.saveButtonText}>
                                Enregistrer ({selectedApprovers.length}/5)
                            </Text>
                        </TouchableOpacity>
                    </View>
                )}

                {activeTab === 'general' && (
                    <View>
                        <Text style={styles.sectionTitle}>Paramètres généraux</Text>

                        <View style={styles.inputGroup}>
                            <Text style={styles.label}>Frais de retard</Text>
                            <TextInput
                                style={styles.input}
                                value={String(settings.late_fee_amount)}
                                onChangeText={(val) => setSettings({ ...settings, late_fee_amount: parseFloat(val) || 0 })}
                                keyboardType="numeric"
                            />
                        </View>

                        <View style={styles.inputGroup}>
                            <Text style={styles.label}>Frais de nourriture</Text>
                            <TextInput
                                style={styles.input}
                                value={String(settings.food_fee)}
                                onChangeText={(val) => setSettings({ ...settings, food_fee: parseFloat(val) || 0 })}
                                keyboardType="numeric"
                            />
                        </View>

                        <View style={styles.inputGroup}>
                            <Text style={styles.label}>Frais de boisson</Text>
                            <TextInput
                                style={styles.input}
                                value={String(settings.drink_fee)}
                                onChangeText={(val) => setSettings({ ...settings, drink_fee: parseFloat(val) || 0 })}
                                keyboardType="numeric"
                            />
                        </View>

                        <View style={styles.switchRow}>
                            <Text>Activer le chat</Text>
                            <Switch value={chatEnabled} onValueChange={setChatEnabled} />
                        </View>
                    </View>
                )}
            </ScrollView>
        </View>
    );
};

const styles = StyleSheet.create({
    container: { flex: 1, backgroundColor: '#f5f5f5' },
    header: { flexDirection: 'row', alignItems: 'center', padding: 16, backgroundColor: '#4f46e5', gap: 16 },
    backButton: { color: '#fff', fontSize: 16 },
    title: { fontSize: 20, fontWeight: 'bold', color: '#fff' },
    tabsContainer: { flexDirection: 'row', backgroundColor: '#fff', borderBottomWidth: 1, borderBottomColor: '#e5e5e5' },
    tab: { flex: 1, paddingVertical: 12, alignItems: 'center' },
    activeTab: { borderBottomWidth: 2, borderBottomColor: '#4f46e5' },
    tabText: { color: '#666', fontSize: 14 },
    activeTabText: { color: '#4f46e5', fontWeight: '600' },
    content: { flex: 1, padding: 16 },
    sectionTitle: { fontSize: 18, fontWeight: 'bold', marginBottom: 8 },
    description: { color: '#666', marginBottom: 16 },
    roleCard: { padding: 16, backgroundColor: '#fff', borderRadius: 8, marginBottom: 8 },
    roleName: { fontSize: 16, fontWeight: '600' },
    rolePermissions: { color: '#666', fontSize: 12 },
    memberCard: { flexDirection: 'row', alignItems: 'center', padding: 16, backgroundColor: '#fff', borderRadius: 8, marginBottom: 8 },
    checkbox: { width: 24, height: 24, borderWidth: 2, borderColor: '#4f46e5', borderRadius: 4, marginRight: 12, alignItems: 'center', justifyContent: 'center' },
    memberName: { fontSize: 16, fontWeight: '600' },
    memberRole: { color: '#666', fontSize: 12 },
    saveButton: { backgroundColor: '#4f46e5', padding: 16, borderRadius: 8, alignItems: 'center', marginTop: 16 },
    disabledButton: { opacity: 0.5 },
    saveButtonText: { color: '#fff', fontWeight: '600', fontSize: 16 },
    inputGroup: { marginBottom: 16 },
    label: { fontSize: 14, fontWeight: '600', marginBottom: 8 },
    input: { backgroundColor: '#fff', padding: 12, borderRadius: 8, borderWidth: 1, borderColor: '#ddd' },
    switchRow: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', padding: 16, backgroundColor: '#fff', borderRadius: 8, marginBottom: 8 },
});

export default AssociationSettingsScreen;
