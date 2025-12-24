import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:local_auth/local_auth.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/utils/constants.dart';
import '../../../../core/widgets/custom_text_field.dart';
import '../../../../core/widgets/custom_button.dart';
import '../../../../core/widgets/glass_container.dart';
import '../bloc/auth_bloc.dart';
import '../widgets/auth_header.dart';
import '../widgets/social_login_section.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({Key? key}) : super(key: key);

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _formKey = GlobalKey<FormState>();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _totpController = TextEditingController();
  
  bool _obscurePassword = true;
  bool _rememberMe = false;
  bool _showTOTP = false;
  bool _biometricsAvailable = false;
  
  final LocalAuthentication _localAuth = LocalAuthentication();

  @override
  void initState() {
    super.initState();
    _checkBiometrics();
    _loadSavedCredentials();
  }

  Future<void> _checkBiometrics() async {
    try {
      final bool isAvailable = await _localAuth.canCheckBiometrics;
      final bool isDeviceSupported = await _localAuth.isDeviceSupported();
      
      setState(() {
        _biometricsAvailable = isAvailable && isDeviceSupported;
      });
    } catch (e) {
      setState(() {
        _biometricsAvailable = false;
      });
    }
  }

  Future<void> _loadSavedCredentials() async {
    // Load saved email if remember me was checked
    // Implementation would use secure storage
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Scaffold(
      backgroundColor: Colors.transparent, // Allow gradient to show
      body: Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
              colors: isDark 
                  ? [const Color(0xFF020617), const Color(0xFF0F172A)] 
                  : [const Color(0xFFFAFBFC), const Color(0xFFEFF6FF)],
            ),
        ),
        child: BlocListener<AuthBloc, AuthState>(
        listener: (context, state) {
          if (state is AuthenticatedState) {
            context.go('/dashboard');
          } else if (state is AuthErrorState) {
            _showErrorSnackBar(state.message);
          } else if (state is Auth2FARequiredState) {
            setState(() {
              _showTOTP = true;
            });
          }
        },
        child: SafeArea(
          child: Center(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(24.0),
              child: GlassContainer(
                padding: const EdgeInsets.all(32),
                borderRadius: 24,
                child: Form(
                  key: _formKey,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      // Header
                      Text(
                        'Bienvenue',
                        textAlign: TextAlign.center,
                        style: GoogleFonts.inter(
                          fontSize: 28,
                          fontWeight: FontWeight.bold,
                          color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                        ),
                      ),
                      const SizedBox(height: 8),
                      Text(
                        'Connectez-vous à votre compte Crypto Bank',
                        textAlign: TextAlign.center,
                        style: GoogleFonts.inter(
                          fontSize: 14,
                          color: isDark ? Colors.white70 : AppTheme.textSecondaryColor,
                        ),
                      ),
                      
                      const SizedBox(height: 40),
                      
                      // Email Field
                      CustomTextField(
                        controller: _emailController,
                        label: 'Adresse Email',
                        hint: 'Entrez votre email',
                        keyboardType: TextInputType.emailAddress,
                        prefixIcon: Icons.email_outlined,
                        validator: _validateEmail,
                        fillColor: isDark ? Colors.white.withOpacity(0.05) : Colors.grey.shade50,
                      ),
                      
                      const SizedBox(height: 16),
                      
                      // Password Field
                      CustomTextField(
                        controller: _passwordController,
                        label: 'Mot de passe',
                        hint: 'Entrez votre mot de passe',
                        obscureText: _obscurePassword,
                        prefixIcon: Icons.lock_outlined,
                        suffixIcon: IconButton(
                          icon: Icon(
                            _obscurePassword ? Icons.visibility_outlined : Icons.visibility_off_outlined,
                           color: isDark ? Colors.white70 : AppTheme.textSecondaryColor,
                          ),
                          onPressed: () {
                            setState(() {
                              _obscurePassword = !_obscurePassword;
                            });
                          },
                        ),
                        validator: _validatePassword,
                         fillColor: isDark ? Colors.white.withOpacity(0.05) : Colors.grey.shade50,
                      ),
                      
                      // 2FA Field (shown when required)
                      if (_showTOTP) ...[
                        const SizedBox(height: 16),
                        CustomTextField(
                          controller: _totpController,
                          label: 'Code d\'authentification',
                          hint: 'Entrez le code à 6 chiffres',
                          keyboardType: TextInputType.number,
                          prefixIcon: Icons.security_outlined,
                          maxLength: 6,
                          validator: _validateTOTP,
                          fillColor: isDark ? Colors.white.withOpacity(0.05) : Colors.grey.shade50,
                        ),
                      ],
                      
                      const SizedBox(height: 16),
                      
                      // Remember Me & Forgot Password
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Row(
                            children: [
                              Transform.scale(
                                scale: 0.9,
                                child: Checkbox(
                                  value: _rememberMe,
                                  activeColor: AppTheme.primaryColor,
                                  side: BorderSide(
                                    color: isDark ? Colors.white60 : Colors.grey.shade400,
                                  ),
                                  onChanged: (value) {
                                    setState(() {
                                      _rememberMe = value ?? false;
                                    });
                                  },
                                ),
                              ),
                              Text(
                                'Se souvenir de moi',
                                style: GoogleFonts.inter(
                                  fontSize: 13,
                                  color: isDark ? Colors.white70 : AppTheme.textSecondaryColor,
                                ),
                              ),
                            ],
                          ),
                          TextButton(
                            onPressed: () => context.push('/auth/forgot-password'),
                            child: Text(
                              'Mot de passe oublié ?',
                              style: GoogleFonts.inter(
                                fontSize: 13,
                                fontWeight: FontWeight.w600,
                                color: AppTheme.primaryColor,
                              ),
                            ),
                          ),
                        ],
                      ),
                      
                      const SizedBox(height: 32),
                      
                      // Login Button
                      BlocBuilder<AuthBloc, AuthState>(
                        builder: (context, state) {
                          final isLoading = state is AuthLoadingState;
                          
                          return Container(
                            decoration: BoxDecoration(
                              borderRadius: BorderRadius.circular(12),
                              boxShadow: [
                                BoxShadow(
                                  color: AppTheme.primaryColor.withOpacity(0.3),
                                  blurRadius: 20,
                                  offset: const Offset(0, 8),
                                ),
                              ],
                            ),
                            child: CustomButton(
                              text: isLoading ? 'Connexion...' : 'Se Connecter',
                              onPressed: isLoading ? null : _handleLogin,
                              isLoading: isLoading,
                              backgroundColor: AppTheme.primaryColor,
                              textColor: Colors.white,
                            ),
                          );
                        },
                      ),
                      
                      // Biometric Login Button
                      if (_biometricsAvailable) ...[
                        const SizedBox(height: 16),
                        OutlinedButton.icon(
                          onPressed: _handleBiometricLogin,
                          icon: const Icon(Icons.fingerprint),
                          label: const Text('Utiliser Biométrie'),
                          style: OutlinedButton.styleFrom(
                            foregroundColor: isDark ? Colors.white : AppTheme.textPrimaryColor,
                            side: BorderSide(
                              color: isDark ? Colors.white24 : Colors.grey.shade300,
                            ),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(12),
                            ),
                            padding: const EdgeInsets.symmetric(vertical: 16),
                          ),
                        ),
                      ],
                      
                      const SizedBox(height: 32),
                      
                      // Social Login Section
                      const SocialLoginSection(),
                      
                      const SizedBox(height: 32),
                      
                      // Sign Up Link
                      Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Text(
                            'Pas encore de compte ? ',
                             style: GoogleFonts.inter(
                                  color: isDark ? Colors.white70 : AppTheme.textSecondaryColor,
                                ),
                          ),
                          TextButton(
                            onPressed: () => context.push('/auth/register'),
                            child: Text(
                              'S\'inscrire',
                              style: GoogleFonts.inter(
                                fontWeight: FontWeight.bold,
                                color: AppTheme.primaryColor,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ),
        ),
      ),
      ),
    );
  }

  String? _validateEmail(String? value) {
    if (value == null || value.isEmpty) {
      return 'Please enter your email';
    }
    if (!RegExp(AppConstants.emailRegex).hasMatch(value)) {
      return 'Please enter a valid email';
    }
    return null;
  }

  String? _validatePassword(String? value) {
    if (value == null || value.isEmpty) {
      return 'Please enter your password';
    }
    if (value.length < AppConstants.minPasswordLength) {
      return 'Password must be at least ${AppConstants.minPasswordLength} characters';
    }
    return null;
  }

  String? _validateTOTP(String? value) {
    if (_showTOTP) {
      if (value == null || value.isEmpty) {
        return 'Please enter the authentication code';
      }
      if (value.length != AppConstants.totpLength) {
        return 'Code must be ${AppConstants.totpLength} digits';
      }
    }
    return null;
  }

  void _handleLogin() {
    if (_formKey.currentState!.validate()) {
      final email = _emailController.text.trim();
      final password = _passwordController.text;
      final totpCode = _totpController.text.trim();
      
      context.read<AuthBloc>().add(
        SignInEvent(
          email: email,
          password: password,
          totpCode: totpCode.isNotEmpty ? totpCode : null,
          rememberMe: _rememberMe,
        ),
      );
    }
  }

  Future<void> _handleBiometricLogin() async {
    try {
      final bool didAuthenticate = await _localAuth.authenticate(
        localizedReason: 'Please authenticate to access your account',
        options: const AuthenticationOptions(
          biometricOnly: false,
          stickyAuth: true,
        ),
      );
      
      if (didAuthenticate) {
        // Load saved credentials and login
        context.read<AuthBloc>().add(BiometricSignInEvent());
      }
    } catch (e) {
      _showErrorSnackBar('Biometric authentication failed');
    }
  }

  void _showErrorSnackBar(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: AppTheme.errorColor,
        action: SnackBarAction(
          label: 'Dismiss',
          textColor: Colors.white,
          onPressed: () => ScaffoldMessenger.of(context).hideCurrentSnackBar(),
        ),
      ),
    );
  }

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    _totpController.dispose();
    super.dispose();
  }
}