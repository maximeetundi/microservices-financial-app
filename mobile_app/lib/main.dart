import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:provider/provider.dart';
import 'package:google_fonts/google_fonts.dart';

import 'core/di/injection_container.dart' as di;
import 'core/theme/app_theme.dart';
import 'core/providers/theme_provider.dart';
import 'core/routes/app_router.dart';
import 'core/utils/constants.dart';
import 'core/api/api_client.dart';
import 'core/services/push_notification_service.dart';
import 'features/auth/presentation/bloc/auth_bloc.dart';
import 'features/wallet/presentation/bloc/wallet_bloc.dart';
import 'features/exchange/presentation/bloc/exchange_bloc.dart';
import 'features/cards/presentation/bloc/cards_bloc.dart';
import 'features/portfolio/presentation/bloc/portfolio_bloc.dart';
import 'features/transfer/presentation/bloc/transfer_bloc.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Initialize dependency injection
  await di.init();
  
  // Initialize push notifications
  await PushNotificationService().initialize();
  
  // Configure auto-logout callback for token expiry
  ApiClient.onLogout = () {
    AppRouter.router.go('/auth/login');
  };
  
  // Configure notification tap handler
  PushNotificationService().onNotificationTap = (payload) {
    if (payload != null) {
      try {
        final data = jsonDecode(payload);
        if (data['route'] != null) {
          AppRouter.router.push(data['route']);
        } else {
          AppRouter.router.go('/notifications');
        }
      } catch (e) {
        debugPrint('Error parsing notification payload: $e');
        AppRouter.router.go('/notifications');
      }
    }
  };
  
  // Set preferred orientations
  await SystemChrome.setPreferredOrientations([
    DeviceOrientation.portraitUp,
    DeviceOrientation.portraitDown,
  ]);
  
  // Set system UI overlay style
  SystemChrome.setSystemUIOverlayStyle(
    const SystemUiOverlayStyle(
      statusBarColor: Colors.transparent,
      statusBarIconBrightness: Brightness.dark,
      systemNavigationBarColor: Colors.white,
      systemNavigationBarIconBrightness: Brightness.dark,
    ),
  );
  
  runApp(const ZekoraApp());
}

class ZekoraApp extends StatelessWidget {
  const ZekoraApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (_) => ThemeProvider(),
      child: MultiBlocProvider(
        providers: [
          BlocProvider(create: (_) => di.sl<AuthBloc>()),
          BlocProvider(create: (_) => di.sl<WalletBloc>()),
          BlocProvider(create: (_) => di.sl<ExchangeBloc>()),
          BlocProvider(create: (_) => di.sl<CardsBloc>()),
          BlocProvider(create: (_) => di.sl<PortfolioBloc>()),
          BlocProvider(create: (_) => di.sl<TransferBloc>()),
        ],
        child: Consumer<ThemeProvider>(
          builder: (context, themeProvider, child) {
            return MaterialApp.router(
              title: AppConstants.appName,
              debugShowCheckedModeBanner: false,
              theme: AppTheme.lightTheme,
              darkTheme: AppTheme.darkTheme,
              themeMode: themeProvider.themeMode,
              routerConfig: AppRouter.createRouter(context.read<AuthBloc>()),
              builder: (context, child) {
                return MediaQuery(
                  data: MediaQuery.of(context).copyWith(
                    textScaleFactor: 1.0, // Prevent text scaling
                  ),
                  child: child!,
                );
              },
            );
          },
        ),
      ),
    );
  }
}

class SplashScreen extends StatefulWidget {
  const SplashScreen({Key? key}) : super(key: key);

  @override
  State<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends State<SplashScreen>
    with SingleTickerProviderStateMixin {
  late AnimationController _animationController;
  late Animation<double> _fadeAnimation;
  late Animation<double> _scaleAnimation;

  @override
  void initState() {
    super.initState();
    
    _animationController = AnimationController(
      duration: const Duration(seconds: 2),
      vsync: this,
    );
    
    _fadeAnimation = Tween<double>(
      begin: 0.0,
      end: 1.0,
    ).animate(CurvedAnimation(
      parent: _animationController,
      curve: Curves.easeInOut,
    ));
    
    _scaleAnimation = Tween<double>(
      begin: 0.8,
      end: 1.0,
    ).animate(CurvedAnimation(
      parent: _animationController,
      curve: Curves.elasticOut,
    ));
    
    _animationController.forward();
    
    // Navigate to appropriate screen after animation
    Future.delayed(const Duration(seconds: 3), () {
      context.read<AuthBloc>().add(CheckAuthStatusEvent());
    });
  }

  @override
  void dispose() {
    _animationController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppTheme.primaryColor,
      body: BlocListener<AuthBloc, AuthState>(
        listener: (context, state) {
          if (state is AuthenticatedState) {
            AppRouter.router.go('/dashboard');
          } else if (state is UnauthenticatedState) {
            AppRouter.router.go('/auth/login');
          }
        },
        child: Center(
          child: AnimatedBuilder(
            animation: _animationController,
            builder: (context, child) {
              return FadeTransition(
                opacity: _fadeAnimation,
                child: ScaleTransition(
                  scale: _scaleAnimation,
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      // App Logo
                      Container(
                        width: 120,
                        height: 120,
                        decoration: BoxDecoration(
                          color: Colors.white,
                          borderRadius: BorderRadius.circular(24),
                          boxShadow: [
                            BoxShadow(
                              color: Colors.black.withOpacity(0.1),
                              blurRadius: 20,
                              offset: const Offset(0, 10),
                            ),
                          ],
                        ),
                        child: const Icon(
                          Icons.account_balance_wallet,
                          size: 60,
                          color: AppTheme.primaryColor,
                        ),
                      ),
                      const SizedBox(height: 24),
                      
                      // App Name
                      Text(
                        AppConstants.appName,
                        style: GoogleFonts.inter(
                          fontSize: 32,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                      ),
                      const SizedBox(height: 8),
                      
                      // Tagline
                      Text(
                        'Your Digital Banking Solution',
                        style: GoogleFonts.inter(
                          fontSize: 16,
                          color: Colors.white.withOpacity(0.9),
                        ),
                      ),
                      const SizedBox(height: 48),
                      
                      // Loading Indicator
                      const SizedBox(
                        width: 40,
                        height: 40,
                        child: CircularProgressIndicator(
                          strokeWidth: 3,
                          valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
                        ),
                      ),
                    ],
                  ),
                ),
              );
            },
          ),
        ),
      ),
    );
  }
}