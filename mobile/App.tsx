import React from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { createStackNavigator } from '@react-navigation/stack';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import Icon from 'react-native-vector-icons/Ionicons';
import AssociationDetailScreen from './src/screens/AssociationDetailScreen';
import AssociationSettingsScreen from './src/screens/AssociationSettingsScreen';
import MessagesScreen from './src/screens/MessagesScreen';
import NotificationsScreen from './src/screens/NotificationsScreen';

const Stack = createStackNavigator();
const Tab = createBottomTabNavigator();

function TabNavigator() {
    return (
        <Tab.Navigator
            screenOptions={{
                headerShown: false,
                tabBarActiveTintColor: '#4f46e5',
                tabBarInactiveTintColor: '#9ca3af',
            }}
        >
            <Tab.Screen
                name="Associations"
                component={AssociationDetailScreen}
                options={{
                    tabBarIcon: ({ color, size }) => <Icon name="people" size={size} color={color} />,
                }}
            />
            <Tab.Screen
                name="Messages"
                component={MessagesScreen}
                options={{
                    tabBarIcon: ({ color, size }) => <Icon name="chatbubbles" size={size} color={color} />,
                }}
            />
            <Tab.Screen
                name="Notifications"
                component={NotificationsScreen}
                options={{
                    tabBarIcon: ({ color, size }) => <Icon name="notifications" size={size} color={color} />,
                }}
            />
        </Tab.Navigator>
    );
}

function App() {
    return (
        <NavigationContainer>
            <Stack.Navigator screenOptions={{ headerShown: false }}>
                <Stack.Screen name="Main" component={TabNavigator} />
                <Stack.Screen name="AssociationSettings" component={AssociationSettingsScreen} />
            </Stack.Navigator>
        </NavigationContainer>
    );
}

export default App;
