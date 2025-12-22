import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../features/auth/presentation/pages/login_page.dart';
import '../../features/auth/presentation/pages/register_page.dart';
import '../../features/auth/presentation/pages/forgot_password_page.dart';
import '../../features/auth/presentation/pages/biometric_setup_page.dart';
import '../../features/dashboard/presentation/pages/dashboard_page.dart';
import '../../features/wallet/presentation/pages/wallet_page.dart';
import '../../features/wallet/presentation/pages/wallet_detail_page.dart';
import '../../features/exchange/presentation/pages/exchange_page.dart';
import '../../features/exchange/presentation/pages/trading_page.dart';
import '../../features/transfer/presentation/pages/transfer_page.dart';
import '../../features/transfer/presentation/pages/transfer_detail_page.dart';
import '../../features/cards/presentation/pages/cards_page.dart';
import '../../features/cards/presentation/pages/card_detail_page.dart';
import '../../features/portfolio/presentation/pages/portfolio_page.dart';
import '../../features/settings/presentation/pages/settings_page.dart';
import '../../features/settings/presentation/pages/security_page.dart';
import '../../features/settings/presentation/pages/profile_page.dart';
import '../../features/notifications/presentation/pages/notifications_page.dart';
import '../../features/support/support_screen.dart';
import '../../features/merchant/merchant_screen.dart';
import '../../features/merchant/scan_pay_screen.dart';
import '../../main.dart';
import '../navigation/main_navigation.dart';
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
        path: '/auth/biometric-setup',
        name: 'biometric-setup',
        builder: (context, state) => const BiometricSetupPage(),
      ),
      
      // Main App Routes with Bottom Navigation
      StatefulShellRoute.indexedStack(
        builder: (context, state, navigationShell) {
          return MainNavigation(navigationShell: navigationShell);
        },
        branches: [
          // Dashboard Branch
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: '/dashboard',
                name: 'dashboard',
                builder: (context, state) => const DashboardPage(),
                routes: [
                  GoRoute(
                    path: 'notifications',
                    name: 'notifications',
                    builder: (context, state) => const NotificationsPage(),
                  ),
                ],
              ),
            ],
          ),
          
          // Wallet Branch
          StatefulShellBranch(
            routes: [
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
            ],
          ),
          
          // Exchange Branch
          StatefulShellBranch(
            routes: [
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
            ],
          ),
          
          // Portfolio Branch
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: '/portfolio',
                name: 'portfolio',
                builder: (context, state) => const PortfolioPage(),
              ),
            ],
          ),
          
          // More Branch
          StatefulShellBranch(
            routes: [
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
          ),
        ],
      ),
    ],
    errorBuilder: (context, state) => Scaffold(
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(
              Icons.error_outline,
              size: 64,
              color: Colors.red,
            ),
            const SizedBox(height: 16),
            Text(
              'Page Not Found',
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            const SizedBox(height: 8),
            Text(
              'The requested page could not be found.',
              style: Theme.of(context).textTheme.bodyMedium,
            ),
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
    // Skip redirect for splash screen
    if (state.matchedLocation == '/splash') {
      return null;
    }
    
    final authBloc = context.read<AuthBloc>();
    final authState = authBloc.state;
    
    // Check if user is on auth pages
    final isOnAuthPage = state.matchedLocation.startsWith('/auth');
    
    if (authState is AuthenticatedState) {
      // User is authenticated, redirect away from auth pages
      if (isOnAuthPage) {
        return '/dashboard';
      }
    } else if (authState is UnauthenticatedState) {
      // User is not authenticated, redirect to login unless already on auth page
      if (!isOnAuthPage) {
        return '/auth/login';
      }
    }
    
    // No redirect needed
    return null;
  }
}

// Custom Page Transitions
class SlidePageRoute<T> extends PageRouteBuilder<T> {
  final Widget page;
  final Offset direction;

  SlidePageRoute({
    required this.page,
    this.direction = const Offset(1.0, 0.0),
  }) : super(
          pageBuilder: (context, animation, secondaryAnimation) => page,
          transitionsBuilder: (context, animation, secondaryAnimation, child) {
            return SlideTransition(
              position: animation.drive(
                Tween(begin: direction, end: Offset.zero).chain(
                  CurveTween(curve: Curves.easeInOut),
                ),
              ),
              child: child,
            );
          },
        );
}

class FadePageRoute<T> extends PageRouteBuilder<T> {
  final Widget page;

  FadePageRoute({required this.page})
      : super(
          pageBuilder: (context, animation, secondaryAnimation) => page,
          transitionsBuilder: (context, animation, secondaryAnimation, child) {
            return FadeTransition(
              opacity: animation,
              child: child,
            );
          },
        );
}

class ScalePageRoute<T> extends PageRouteBuilder<T> {
  final Widget page;

  ScalePageRoute({required this.page})
      : super(
          pageBuilder: (context, animation, secondaryAnimation) => page,
          transitionsBuilder: (context, animation, secondaryAnimation, child) {
            return ScaleTransition(
              scale: animation.drive(
                Tween(begin: 0.8, end: 1.0).chain(
                  CurveTween(curve: Curves.elasticOut),
                ),
              ),
              child: child,
            );
          },
        );
}