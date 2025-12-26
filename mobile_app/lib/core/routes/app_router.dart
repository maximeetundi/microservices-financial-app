import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../features/auth/presentation/pages/login_page.dart';
import '../../features/auth/presentation/pages/register_page.dart';
import '../../features/auth/presentation/pages/forgot_password_page.dart';
import '../../features/auth/presentation/pages/biometric_setup_page.dart';
import '../../features/auth/presentation/pages/pin_setup_page.dart';
import '../../features/dashboard/presentation/pages/modern_home_page.dart';
import '../../features/wallet/presentation/pages/wallet_page.dart';
import '../../features/wallet/presentation/pages/wallet_detail_page.dart';
import '../../features/exchange/presentation/pages/exchange_page.dart';
import '../../features/exchange/presentation/pages/trading_page.dart';
import '../../features/transfer/presentation/pages/transfer_page.dart';
import '../../features/transfer/presentation/pages/transfer_detail_page.dart';
import '../../features/wallet/presentation/pages/transactions_page.dart';
import '../../features/cards/presentation/pages/cards_page.dart';
import '../../features/cards/presentation/pages/card_detail_page.dart';
import '../../features/portfolio/presentation/pages/portfolio_page.dart';
import '../../features/settings/presentation/pages/settings_page.dart';
import '../../features/settings/presentation/pages/security_page.dart';
import '../../features/settings/presentation/pages/profile_page.dart';
import '../../features/settings/presentation/pages/kyc_page.dart';
import '../../features/settings/presentation/pages/preferences_page.dart';
import '../../features/settings/presentation/pages/payment_methods_page.dart';
import '../../features/notifications/presentation/pages/notifications_page.dart';
import '../../features/support/support_screen.dart';
import '../../features/merchant/merchant_screen.dart';
import '../../features/merchant/scan_pay_screen.dart';
import '../../main.dart';
import '../../features/auth/presentation/bloc/auth_bloc.dart';

class AppRouter {
  static final GoRouter router = GoRouter(
    initialLocation: '/splash',
    debugLogDiagnostics: true,
    redirect: _redirect,
    routes: [
      // Splash Route
      GoRoute(
        path: '/splash',
        name: 'splash',
        builder: (context, state) => const SplashScreen(),
      ),
      
      // Auth Routes
      GoRoute(
        path: '/auth/login',
        name: 'login',
        builder: (context, state) => const LoginPage(),
      ),
      GoRoute(
        path: '/auth/register',
        name: 'register',
        builder: (context, state) => const RegisterPage(),
      ),
      GoRoute(
        path: '/auth/forgot-password',
        name: 'forgot-password',
        builder: (context, state) => const ForgotPasswordPage(),
      ),
      GoRoute(
        path: '/auth/pin-setup',
        name: 'pin-setup',
        builder: (context, state) => const PinSetupPage(),
      ),
      GoRoute(
        path: '/auth/biometric-setup',
        name: 'biometric-setup',
        builder: (context, state) => const BiometricSetupPage(),
      ),
      
      // Modern Home with Animated Drawer
      GoRoute(
        path: '/dashboard',
        name: 'dashboard',
        builder: (context, state) => const ModernHomePage(),
        routes: [
          GoRoute(
            path: 'notifications',
            name: 'notifications',
            builder: (context, state) => const NotificationsPage(),
          ),
        ],
      ),
      
      // Wallet Route
      GoRoute(
        path: '/wallet',
        name: 'wallet',
        builder: (context, state) => const WalletPage(),
        routes: [
          GoRoute(
            path: ':walletId',
            name: 'wallet-detail',
            builder: (context, state) => WalletDetailPage(
              walletId: state.pathParameters['walletId']!,
            ),
          ),
        ],
      ),
      
      // Exchange Route
      GoRoute(
        path: '/exchange',
        name: 'exchange',
        builder: (context, state) => const ExchangePage(),
        routes: [
          GoRoute(
            path: 'trading',
            name: 'trading',
            builder: (context, state) => const TradingPage(),
          ),
        ],
      ),
      
      // Portfolio Route
      GoRoute(
        path: '/portfolio',
        name: 'portfolio',
        builder: (context, state) => const PortfolioPage(),
      ),
      
      // Transactions Route
      GoRoute(
        path: '/transactions',
        name: 'transactions',
        builder: (context, state) => const TransactionsPage(),
      ),
      
      // More/Settings Route
      GoRoute(
        path: '/more',
        name: 'more',
        builder: (context, state) => const SettingsPage(),
        routes: [
          GoRoute(
            path: 'profile',
            name: 'profile',
            builder: (context, state) => const ProfilePage(),
          ),
          GoRoute(
            path: 'security',
            name: 'security',
            builder: (context, state) => const SecurityPage(),
          ),
          GoRoute(
            path: 'kyc',
            name: 'kyc',
            builder: (context, state) => const KycPage(),
          ),
          GoRoute(
            path: 'preferences',
            name: 'preferences',
            builder: (context, state) => const PreferencesPage(),
          ),
          GoRoute(
            path: 'payment-methods',
            name: 'payment-methods',
            builder: (context, state) => const PaymentMethodsPage(),
          ),
          GoRoute(
            path: 'settings',
            name: 'settings',
            builder: (context, state) => const SettingsPage(),
          ),
          GoRoute(
            path: 'cards',
            name: 'cards',
            builder: (context, state) => const CardsPage(),
            routes: [
              GoRoute(
                path: ':cardId',
                name: 'card-detail',
                builder: (context, state) => CardDetailPage(
                  cardId: state.pathParameters['cardId']!,
                ),
              ),
            ],
          ),
          GoRoute(
            path: 'transfer',
            name: 'transfer',
            builder: (context, state) => const TransferPage(),
            routes: [
              GoRoute(
                path: ':transferId',
                name: 'transfer-detail',
                builder: (context, state) => TransferDetailPage(
                  transferId: state.pathParameters['transferId']!,
                ),
              ),
            ],
          ),
          GoRoute(
            path: 'support',
            name: 'support',
            builder: (context, state) => const SupportScreen(),
            routes: [
              GoRoute(
                path: 'chat',
                name: 'support-chat',
                builder: (context, state) {
                  final ticketId = state.uri.queryParameters['id'];
                  final agentType = state.uri.queryParameters['agent'] ?? 'ai';
                  return ChatScreen(
                    agentType: agentType,
                    ticketId: ticketId,
                  );
                },
              ),
            ],
          ),
          GoRoute(
            path: 'merchant',
            name: 'merchant',
            builder: (context, state) => const MerchantScreen(),
            routes: [
              GoRoute(
                path: 'scan',
                name: 'merchant-scan',
                builder: (context, state) => const ScanPayScreen(),
              ),
            ],
          ),
        ],
      ),
    ],
    errorBuilder: (context, state) => Scaffold(
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.error_outline, size: 64, color: Colors.red),
            const SizedBox(height: 16),
            Text('Page Not Found', style: Theme.of(context).textTheme.headlineMedium),
            const SizedBox(height: 24),
            ElevatedButton(
              onPressed: () => context.go('/dashboard'),
              child: const Text('Go to Dashboard'),
            ),
          ],
        ),
      ),
    ),
  );

  // Route Guard
  static String? _redirect(BuildContext context, GoRouterState state) {
    if (state.matchedLocation == '/splash') return null;
    
    final authBloc = context.read<AuthBloc>();
    final authState = authBloc.state;
    final isOnAuthPage = state.matchedLocation.startsWith('/auth');
    
    if (authState is AuthenticatedState) {
      // Check if PIN is setup
      // Note: We check if matchedLocation is NOT pin-setup to avoid infinite loop
      if (!authState.user.hasPin && state.matchedLocation != '/auth/pin-setup') {
        return '/auth/pin-setup';
      }
      
      // If PIN is setup (or we are already on pin-setup), redirect away from other auth pages
      // But allow pin-setup if we are there (though normally we should be redirected if hasPin is true, logic below handles it)
      if (isOnAuthPage && state.matchedLocation != '/auth/pin-setup') return '/dashboard';
      
      // If PIN is setup and user tries to access pin-setup page, redirect to dashboard? 
      // Maybe user wants to change pin? That would be a different route usually.
      // For initial setup flow:
      if (authState.user.hasPin && state.matchedLocation == '/auth/pin-setup') {
         return '/dashboard';
      }
    } else if (authState is UnauthenticatedState) {
      if (!isOnAuthPage) return '/auth/login';
    }
    
    return null;
  }
}