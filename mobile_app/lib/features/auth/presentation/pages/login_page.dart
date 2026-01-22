import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:local_auth/local_auth.dart';
import 'dart:math' as math;

import '../../../../core/theme/app_theme.dart';
import '../../../../core/utils/constants.dart';
import '../../../../core/widgets/custom_text_field.dart';
import '../../../../core/widgets/glass_container.dart';
import '../../../../core/providers/theme_provider.dart';
import '../bloc/auth_bloc.dart';

/// Login page matching web frontend design exactly
/// Features: animated blob background, glass card, theme toggle, 2FA support
class LoginPage extends StatefulWidget {
  const LoginPage({Key? key}) : super(key: key);

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> with TickerProviderStateMixin {
  final _formKey = GlobalKey<FormState>();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _totpController = TextEditingController();
  
  bool _obscurePassword = true;
  bool _rememberMe = false;
  bool _showTOTP = false;
  bool _biometricsAvailable = false;
  String? _errorMessage;
  
  final LocalAuthentication _localAuth = LocalAuthentication();
  
  // Animation controllers for blob animations
  late AnimationController _blobController1;
  late AnimationController _blobController2;
  late AnimationController _blobController3;

  @override
  void initState() {
    super.initState();
    _checkBiometrics();
    _initBlobAnimations();
  }

  void _initBlobAnimations() {
    _blobController1 = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 7),
    )..repeat(reverse: true);
    
    _blobController2 = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 9),
    )..repeat(reverse: true);
    
    _blobController3 = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 11),
    )..repeat(reverse: true);
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

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final size = MediaQuery.of(context).size;
    
    return Scaffold(
      backgroundColor: Colors.transparent,
      body: Container(
        decoration: BoxDecoration(
          // Match web background gradient
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: isDark 
                ? [const Color(0xFF0F0F1A), const Color(0xFF1A1A2E)]
                : [const Color(0xFFF8FAFC), const Color(0xFFEEF2FF)],
          ),
        ),
        child: Stack(
          children: [
            // Animated Blob Backgrounds (matching web)
            _buildAnimatedBlobs(size, isDark),
            
            // Main Content
            BlocListener<AuthBloc, AuthState>(
              listener: (context, state) {
                if (state is AuthenticatedState) {
                  context.go('/dashboard');
                } else if (state is AuthErrorState) {
                  setState(() => _errorMessage = state.message);
                } else if (state is Auth2FARequiredState) {
                  setState(() {
                    _showTOTP = true;
                    _errorMessage = 'Veuillez entrer votre code d\'authentification.';
                  });
                }
              },
              child: SafeArea(
                child: Center(
                  child: SingleChildScrollView(
                    padding: const EdgeInsets.all(24.0),
                    child: Column(
                      children: [
                        // Logo Section
                        _buildLogoSection(isDark),
                        const SizedBox(height: 32),
                        
                        // Login Card (Glass Container)
                        GlassContainer(
                          padding: const EdgeInsets.all(32),
                          borderRadius: 24,
                          borderColor: isDark 
                              ? Colors.white.withOpacity(0.1)
                              : const Color(0xFFE2E8F0),
                          color: isDark 
                              ? Colors.white.withOpacity(0.05)
                              : Colors.white.withOpacity(0.9),
                          child: Form(
                            key: _formKey,
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.stretch,
                              children: [
                                // Email Field
                                _buildInputLabel('Adresse email', Icons.alternate_email_rounded, isDark),
                                const SizedBox(height: 8),
                                _buildTextField(
                                  controller: _emailController,
                                  hint: 'exemple@email.com',
                                  keyboardType: TextInputType.emailAddress,
                                  isDark: isDark,
                                ),
                                
                                const SizedBox(height: 24),
                                
                                // Password Field
                                _buildInputLabel('Mot de passe', Icons.lock_outline_rounded, isDark),
                                const SizedBox(height: 8),
                                _buildPasswordField(isDark),
                                
                                // 2FA Field
                                if (_showTOTP) ...[
                                  const SizedBox(height: 24),
                                  _buildInputLabel('Code d\'authentification', Icons.security_rounded, isDark),
                                  const SizedBox(height: 8),
                                  _buildTotpField(isDark),
                                ],
                                
                                const SizedBox(height: 16),
                                
                                // Remember Me & Forgot Password
                                _buildRememberForgotRow(isDark),
                                
                                // Error Message
                                if (_errorMessage != null) ...[
                                  const SizedBox(height: 16),
                                  _buildErrorBanner(),
                                ],
                                
                                const SizedBox(height: 32),
                                
                                // Login Button
                                _buildLoginButton(isDark),
                                
                                // Biometric Login
                                if (_biometricsAvailable) ...[
                                  const SizedBox(height: 16),
                                  _buildBiometricButton(isDark),
                                ],
                                
                                const SizedBox(height: 32),
                                
                                // Register Link
                                _buildRegisterLink(isDark),
                              ],
                            ),
                          ),
                        ),
                        
                        const SizedBox(height: 32),
                        
                        // Footer Links
                        _buildFooterLinks(isDark),
                      ],
                    ),
                  ),
                ),
              ),
            ),
            
            // Theme Toggle Button (top right)
            Positioned(
              top: MediaQuery.of(context).padding.top + 16,
              right: 16,
              child: _buildThemeToggle(isDark),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildAnimatedBlobs(Size size, bool isDark) {
    return Stack(
      children: [
        // Blob 1 - Top Right, Indigo
        AnimatedBuilder(
          animation: _blobController1,
          builder: (context, child) {
            final value = _blobController1.value;
            return Positioned(
              top: -size.height * 0.2 + (30 * math.sin(value * math.pi * 2)),
              right: -size.width * 0.1 + (50 * math.cos(value * math.pi * 2)),
              child: Container(
                width: size.width * 0.7,
                height: size.width * 0.7,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  gradient: RadialGradient(
                    colors: [
                      const Color(0xFF6366F1).withOpacity(0.2),
                      Colors.transparent,
                    ],
                  ),
                ),
              ),
            );
          },
        ),
        
        // Blob 2 - Left Side, Purple
        AnimatedBuilder(
          animation: _blobController2,
          builder: (context, child) {
            final value = _blobController2.value;
            return Positioned(
              top: size.height * 0.2 + (20 * math.cos(value * math.pi * 2)),
              left: -size.width * 0.1 + (-30 * math.sin(value * math.pi * 2)),
              child: Container(
                width: size.width * 0.5,
                height: size.width * 0.5,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  gradient: RadialGradient(
                    colors: [
                      const Color(0xFF8B5CF6).withOpacity(0.2),
                      Colors.transparent,
                    ],
                  ),
                ),
              ),
            );
          },
        ),
        
        // Blob 3 - Bottom Right, Blue
        AnimatedBuilder(
          animation: _blobController3,
          builder: (context, child) {
            final value = _blobController3.value;
            return Positioned(
              bottom: -size.height * 0.2 + (-20 * math.sin(value * math.pi * 2)),
              right: size.width * 0.2 + (20 * math.cos(value * math.pi * 2)),
              child: Container(
                width: size.width * 0.6,
                height: size.width * 0.6,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  gradient: RadialGradient(
                    colors: [
                      const Color(0xFF3B82F6).withOpacity(0.1),
                      Colors.transparent,
                    ],
                  ),
                ),
              ),
            );
          },
        ),
      ],
    );
  }

  Widget _buildLogoSection(bool isDark) {
    return Column(
      children: [
        Container(
          width: 80,
          height: 80,
          decoration: BoxDecoration(
            gradient: const LinearGradient(
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
              colors: [Color(0xFF6366F1), Color(0xFF8B5CF6)],
            ),
            borderRadius: BorderRadius.circular(20),
            boxShadow: [
              BoxShadow(
                color: const Color(0xFF6366F1).withOpacity(0.3),
                blurRadius: 20,
                offset: const Offset(0, 10),
              ),
            ],
          ),
          child: Padding(
            padding: const EdgeInsets.all(12.0),
            child: Image.asset(
              'assets/images/logo.png',
              fit: BoxFit.contain,
            ),
          ),
        ),
        const SizedBox(height: 24),
        Text(
          'Zekora',
          style: GoogleFonts.inter(
            fontSize: 28,
            fontWeight: FontWeight.bold,
            color: isDark ? Colors.white : const Color(0xFF1E293B),
          ),
        ),
        const SizedBox(height: 8),
        Text(
          'Connectez-vous à votre espace sécurisé',
          style: GoogleFonts.inter(
            fontSize: 14,
            color: isDark ? const Color(0xFF818CF8).withOpacity(0.8) : const Color(0xFF64748B),
          ),
        ),
      ],
    );
  }

  Widget _buildInputLabel(String label, IconData icon, bool isDark) {
    return Row(
      children: [
        Icon(
          icon,
          size: 16,
          color: isDark ? const Color(0xFF818CF8) : const Color(0xFF6366F1),
        ),
        const SizedBox(width: 8),
        Text(
          label,
          style: GoogleFonts.inter(
            fontSize: 14,
            fontWeight: FontWeight.w500,
            color: isDark ? const Color(0xFF818CF8) : const Color(0xFF374151),
          ),
        ),
      ],
    );
  }

  Widget _buildTextField({
    required TextEditingController controller,
    required String hint,
    required bool isDark,
    TextInputType? keyboardType,
  }) {
    return TextField(
      controller: controller,
      keyboardType: keyboardType,
      style: GoogleFonts.inter(
        color: isDark ? Colors.white : const Color(0xFF1E293B),
      ),
      decoration: InputDecoration(
        hintText: hint,
        hintStyle: GoogleFonts.inter(
          color: isDark ? const Color(0xFF64748B) : const Color(0xFF9CA3AF),
        ),
        filled: true,
        fillColor: isDark 
            ? Colors.black.withOpacity(0.2)
            : const Color(0xFFF9FAFB),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: BorderSide(
            color: isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFD1D5DB),
          ),
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: BorderSide(
            color: isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFD1D5DB),
          ),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: BorderSide(
            color: isDark ? const Color(0xFF6366F1).withOpacity(0.5) : const Color(0xFF6366F1),
            width: 2,
          ),
        ),
        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
      ),
    );
  }

  Widget _buildPasswordField(bool isDark) {
    return TextField(
      controller: _passwordController,
      obscureText: _obscurePassword,
      style: GoogleFonts.inter(
        color: isDark ? Colors.white : const Color(0xFF1E293B),
      ),
      decoration: InputDecoration(
        hintText: '••••••••',
        hintStyle: GoogleFonts.inter(
          color: isDark ? const Color(0xFF64748B) : const Color(0xFF9CA3AF),
        ),
        filled: true,
        fillColor: isDark 
            ? Colors.black.withOpacity(0.2)
            : const Color(0xFFF9FAFB),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: BorderSide(
            color: isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFD1D5DB),
          ),
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: BorderSide(
            color: isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFD1D5DB),
          ),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: BorderSide(
            color: isDark ? const Color(0xFF6366F1).withOpacity(0.5) : const Color(0xFF6366F1),
            width: 2,
          ),
        ),
        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
        suffixIcon: IconButton(
          icon: Icon(
            _obscurePassword ? Icons.visibility_outlined : Icons.visibility_off_outlined,
            color: isDark ? const Color(0xFF818CF8) : const Color(0xFF6B7280),
          ),
          onPressed: () => setState(() => _obscurePassword = !_obscurePassword),
        ),
      ),
    );
  }

  Widget _buildTotpField(bool isDark) {
    return TextField(
      controller: _totpController,
      keyboardType: TextInputType.number,
      maxLength: 6,
      textAlign: TextAlign.center,
      style: GoogleFonts.robotoMono(
        fontSize: 24,
        letterSpacing: 8,
        color: isDark ? Colors.white : const Color(0xFF1E293B),
      ),
      decoration: InputDecoration(
        hintText: '000000',
        counterText: '',
        hintStyle: GoogleFonts.robotoMono(
          fontSize: 24,
          letterSpacing: 8,
          color: isDark ? const Color(0xFF64748B) : const Color(0xFF9CA3AF),
        ),
        filled: true,
        fillColor: isDark 
            ? Colors.black.withOpacity(0.2)
            : const Color(0xFFF9FAFB),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: BorderSide(
            color: isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFD1D5DB),
          ),
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: BorderSide(
            color: isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFD1D5DB),
          ),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: BorderSide(
            color: isDark ? const Color(0xFF6366F1).withOpacity(0.5) : const Color(0xFF6366F1),
            width: 2,
          ),
        ),
      ),
    );
  }

  Widget _buildRememberForgotRow(bool isDark) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Row(
          children: [
            SizedBox(
              width: 20,
              height: 20,
              child: Checkbox(
                value: _rememberMe,
                activeColor: const Color(0xFF6366F1),
                side: BorderSide(
                  color: isDark ? Colors.white.withOpacity(0.2) : const Color(0xFFD1D5DB),
                ),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(4),
                ),
                onChanged: (value) => setState(() => _rememberMe = value ?? false),
              ),
            ),
            const SizedBox(width: 8),
            Text(
              'Se souvenir de moi',
              style: GoogleFonts.inter(
                fontSize: 13,
                color: isDark ? const Color(0xFF818CF8).withOpacity(0.8) : const Color(0xFF6B7280),
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
              color: const Color(0xFF818CF8),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildErrorBanner() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xFFEF4444).withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: const Color(0xFFEF4444).withOpacity(0.2),
        ),
      ),
      child: Row(
        children: [
          Container(
            padding: const EdgeInsets.all(8),
            decoration: BoxDecoration(
              color: const Color(0xFFEF4444).withOpacity(0.2),
              shape: BoxShape.circle,
            ),
            child: const Icon(
              Icons.error_outline,
              size: 16,
              color: Color(0xFFFCA5A5),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              _errorMessage!,
              style: GoogleFonts.inter(
                fontSize: 14,
                fontWeight: FontWeight.w500,
                color: const Color(0xFFFCA5A5),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildLoginButton(bool isDark) {
    return BlocBuilder<AuthBloc, AuthState>(
      builder: (context, state) {
        final isLoading = state is AuthLoadingState;
        
        return Container(
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(12),
            gradient: const LinearGradient(
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
              colors: [Color(0xFF6366F1), Color(0xFF8B5CF6)],
            ),
            boxShadow: [
              BoxShadow(
                color: const Color(0xFF6366F1).withOpacity(0.25),
                blurRadius: 20,
                offset: const Offset(0, 8),
              ),
            ],
          ),
          child: ElevatedButton(
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.transparent,
              shadowColor: Colors.transparent,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
              padding: const EdgeInsets.symmetric(vertical: 16),
            ),
            onPressed: isLoading ? null : _handleLogin,
            child: isLoading 
              ? const SizedBox(
                  height: 20, 
                  width: 20, 
                  child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2),
                ) 
              : Text(
                  'Se connecter',
                  style: GoogleFonts.inter(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  ),
                ),
          ),
        );
      },
    );
  }

  Widget _buildBiometricButton(bool isDark) {
    return OutlinedButton.icon(
      onPressed: _handleBiometricLogin,
      icon: Icon(
        Icons.fingerprint,
        color: isDark ? Colors.white : const Color(0xFF1E293B),
      ),
      label: Text(
        'Utiliser Biométrie',
        style: GoogleFonts.inter(
          color: isDark ? Colors.white : const Color(0xFF1E293B),
        ),
      ),
      style: OutlinedButton.styleFrom(
        side: BorderSide(
          color: isDark ? Colors.white.withOpacity(0.2) : const Color(0xFFD1D5DB),
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
        ),
        padding: const EdgeInsets.symmetric(vertical: 16),
      ),
    );
  }

  Widget _buildRegisterLink(bool isDark) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Text(
          'Pas encore de compte ? ',
          style: GoogleFonts.inter(
            fontSize: 14,
            color: isDark ? const Color(0xFF818CF8).withOpacity(0.6) : const Color(0xFF6B7280),
          ),
        ),
        TextButton(
          onPressed: () => context.push('/auth/register'),
          child: Text(
            'Créer un compte',
            style: GoogleFonts.inter(
              fontSize: 14,
              fontWeight: FontWeight.bold,
              color: const Color(0xFF818CF8),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildFooterLinks(bool isDark) {
    final linkColor = isDark 
        ? const Color(0xFF818CF8).withOpacity(0.4)
        : const Color(0xFF9CA3AF);
    
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Text('Confidentialité', style: GoogleFonts.inter(fontSize: 12, color: linkColor)),
        Text(' • ', style: GoogleFonts.inter(fontSize: 12, color: linkColor.withOpacity(0.5))),
        Text('CGU', style: GoogleFonts.inter(fontSize: 12, color: linkColor)),
        Text(' • ', style: GoogleFonts.inter(fontSize: 12, color: linkColor.withOpacity(0.5))),
        Text('Aide', style: GoogleFonts.inter(fontSize: 12, color: linkColor)),
      ],
    );
  }

  Widget _buildThemeToggle(bool isDark) {
    return GestureDetector(
      onTap: () => ThemeProvider().toggleTheme(),
      child: GlassContainer(
        padding: const EdgeInsets.all(12),
        borderRadius: 12,
        showTopHighlight: false,
        child: Icon(
          isDark ? Icons.wb_sunny_outlined : Icons.dark_mode_outlined,
          size: 20,
          color: isDark ? const Color(0xFFFBBF24) : const Color(0xFF6366F1),
        ),
      ),
    );
  }

  void _handleLogin() {
    setState(() => _errorMessage = null);
    
    final email = _emailController.text.trim();
    final password = _passwordController.text;
    final totpCode = _totpController.text.trim();
    
    if (email.isEmpty || password.isEmpty) {
      setState(() => _errorMessage = 'Veuillez remplir tous les champs.');
      return;
    }
    
    context.read<AuthBloc>().add(
      SignInEvent(
        email: email,
        password: password,
        totpCode: totpCode.isNotEmpty ? totpCode : null,
        rememberMe: _rememberMe,
      ),
    );
  }

  Future<void> _handleBiometricLogin() async {
    try {
      final bool didAuthenticate = await _localAuth.authenticate(
        localizedReason: 'Authentifiez-vous pour accéder à votre compte',
        options: const AuthenticationOptions(
          biometricOnly: false,
          stickyAuth: true,
        ),
      );
      
      if (didAuthenticate) {
        context.read<AuthBloc>().add(BiometricSignInEvent());
      }
    } catch (e) {
      setState(() => _errorMessage = 'Échec de l\'authentification biométrique');
    }
  }

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    _totpController.dispose();
    _blobController1.dispose();
    _blobController2.dispose();
    _blobController3.dispose();
    super.dispose();
  }
}