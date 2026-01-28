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
  
  // Phone Login State
  String _loginMethod = 'phone'; // 'phone' or 'email'
  final _phoneController = TextEditingController();
  final _dialCodeController = TextEditingController(text: '+225');
  String? _selectedCountry = 'CI';

  bool _obscurePassword = true;
  bool _rememberMe = false;
  // _showTOTP is removed in favor of Modal
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

  String _getFlagEmoji(String countryCode) {
    return countryCode.toUpperCase().replaceAllMapped(RegExp(r'[A-Z]'),
        (match) => String.fromCharCode(match.group(0)!.codeUnitAt(0) + 127397));
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

  void _showCountryPicker(BuildContext context, bool isDark) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: isDark ? const Color(0xFF1E293B) : Colors.white,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      builder: (context) {
        String searchQuery = '';
        return StatefulBuilder(
          builder: (context, setModalState) {
            final filtered = AppConstants.countries.where((c) {
               final q = searchQuery.toLowerCase();
               return c['name']!.toLowerCase().contains(q) || c['dial_code']!.contains(q);
            }).toList();

            return Container(
              height: MediaQuery.of(context).size.height * 0.7,
              padding: const EdgeInsets.all(24),
              child: Column(
                children: [
                   Container(
                    width: 40, height: 4,
                    margin: const EdgeInsets.only(bottom: 24),
                    decoration: BoxDecoration(
                      color: isDark ? Colors.white24 : Colors.grey.shade300,
                      borderRadius: BorderRadius.circular(2),
                    ),
                  ),
                  Container(
                    decoration: BoxDecoration(
                      color: isDark ? Colors.white.withOpacity(0.05) : Colors.grey.shade100,
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: TextField(
                      style: TextStyle(color: isDark ? Colors.white : Colors.black87),
                      decoration: InputDecoration(
                        hintText: 'Rechercher un pays...',
                        hintStyle: TextStyle(color: isDark ? Colors.white38 : Colors.grey),
                        prefixIcon: Icon(Icons.search, color: isDark ? Colors.white54 : Colors.grey),
                        border: InputBorder.none,
                        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                      ),
                      onChanged: (v) => setModalState(() => searchQuery = v),
                    ),
                  ),
                  const SizedBox(height: 16),
                  Expanded(
                    child: ListView.separated(
                      itemCount: filtered.length,
                      separatorBuilder: (_, __) => Divider(color: isDark ? Colors.white10 : Colors.grey.shade100),
                      itemBuilder: (context, index) {
                        final country = filtered[index];
                        final isSelected = country['code'] == _selectedCountry;
                        return ListTile(
                          contentPadding: EdgeInsets.zero,
                          onTap: () {
                            _onCountryChanged(country['code']);
                            Navigator.pop(context);
                          },
                          leading: Text(
                            _getFlagEmoji(country['code']!),
                            style: const TextStyle(fontSize: 24),
                          ),
                          title: Text(
                            country['name']!,
                            style: TextStyle(
                              color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                              fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
                            ),
                          ),
                          trailing: isSelected 
                              ? const Icon(Icons.check, color: AppTheme.primaryColor) 
                              : null,
                        );
                      },
                    ),
                  ),
                ],
              ),
            );
          },
        );
      },
    );
  }

  void _onCountryChanged(String? value) {
    setState(() {
      _selectedCountry = value;
      if (value != null) {
        final country = AppConstants.countries.firstWhere(
            (c) => c['code'] == value, 
            orElse: () => {'dial_code': ''}
        );
        final dialCode = country['dial_code'] ?? '';
        if (_dialCodeController.text != dialCode) {
           _dialCodeController.text = dialCode;
        }
      }
    });
  }

  void _onDialCodeChanged(String value) {
    if (!value.startsWith('+')) value = '+$value';
    
    // Find matching country
    final priorityMap = {'+1': 'US', '+33': 'FR', '+225': 'CI'}; // Simplified priority
    String? matched;
    
    if (priorityMap.containsKey(value)) {
      matched = priorityMap[value];
    } else {
      final country = AppConstants.countries.firstWhere(
        (c) => c['dial_code'] == value,
        orElse: () => {},
      );
      if (country.isNotEmpty) matched = country['code'];
    }

    if (matched != null && matched != _selectedCountry) {
       setState(() => _selectedCountry = matched);
    }
  }

  void _onPhoneChanged(String value) {
      // Basic formatter if needed, or leave raw
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

  void _showTOTPModal(String? tempToken) {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) {
        final isDark = Theme.of(context).brightness == Brightness.dark;
        return Dialog(
          backgroundColor: isDark ? const Color(0xFF1E293B) : Colors.white,
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(24)),
          child: Padding(
            padding: const EdgeInsets.all(24.0),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Container(
                  width: 64, height: 64,
                  decoration: BoxDecoration(
                    color: isDark ? const Color(0xFF6366F1).withOpacity(0.2) : const Color(0xFFEEF2FF),
                    shape: BoxShape.circle,
                  ),
                  child: const Center(child: Text("🔐", style: TextStyle(fontSize: 32))),
                ),
                const SizedBox(height: 16),
                Text(
                  "Authentification 2FA",
                  style: GoogleFonts.inter(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    color: isDark ? Colors.white : Colors.black,
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  "Entrez le code généré par votre application Google Authenticator",
                  textAlign: TextAlign.center,
                  style: GoogleFonts.inter(
                    fontSize: 14,
                    color: isDark ? Colors.grey : Colors.black54,
                  ),
                ),
                const SizedBox(height: 24),
                TextField(
                    controller: _totpController,
                    keyboardType: TextInputType.number,
                    maxLength: 6,
                    textAlign: TextAlign.center,
                    autofocus: true,
                    style: GoogleFonts.robotoMono(
                        fontSize: 24,
                        letterSpacing: 8,
                        color: isDark ? Colors.white : Colors.black
                    ),
                    decoration: InputDecoration(
                        hintText: '000000',
                        counterText: '',
                        filled: true,
                        fillColor: isDark ? Colors.black26 : Colors.grey.shade100,
                        border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(12),
                            borderSide: BorderSide.none
                        )
                    ),
                ),
                const SizedBox(height: 24),
                Row(
                  children: [
                    Expanded(
                      child: TextButton(
                        onPressed: () {
                           Navigator.pop(context);
                           context.read<AuthBloc>().add(const SignOutEvent()); // Reset state potentially
                        }, 
                        child: Text("Annuler", style: TextStyle(color: isDark ? Colors.white60 : Colors.grey)),
                      ),
                    ),
                    const SizedBox(width: 16),
                    Expanded(
                      child: ElevatedButton(
                        style: ElevatedButton.styleFrom(
                          backgroundColor: const Color(0xFF6366F1),
                          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                          padding: const EdgeInsets.symmetric(vertical: 16),
                        ),
                        onPressed: () {
                           Navigator.pop(context);
                           _handleLogin(); // Retry login with code
                        },
                        child: const Text("Vérifier", style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
                      ),
                    ),
                  ],
                )
              ],
            ),
          ),
        );
      },
    );
  }

  @override
  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final size = MediaQuery.of(context).size;
    
    return Scaffold(
      backgroundColor: Colors.transparent,
      body: Container(
        decoration: BoxDecoration(
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
            _buildAnimatedBlobs(size, isDark),
            
            BlocListener<AuthBloc, AuthState>(
              listener: (context, state) {
                if (state is AuthenticatedState) {
                  context.go('/dashboard');
                } else if (state is AuthErrorState) {
                  setState(() => _errorMessage = state.message);
                } else if (state is Auth2FARequiredState) {
                  _showTOTPModal(state.tempToken);
                }
              },
              child: SafeArea(
                child: Center(
                  child: SingleChildScrollView(
                    padding: const EdgeInsets.all(24.0),
                    child: Column(
                      children: [
                        _buildLogoSection(isDark),
                        const SizedBox(height: 32),
                        
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
                                // Login Tabs
                                _buildTabs(isDark),
                                const SizedBox(height: 24),
                                
                                // Conditional Input
                                if (_loginMethod == 'phone')
                                  _buildPhoneInput(isDark)
                                else ...[
                                  _buildInputLabel('Adresse email', Icons.alternate_email_rounded, isDark),
                                  const SizedBox(height: 8),
                                  _buildTextField(
                                    controller: _emailController,
                                    hint: 'exemple@email.com',
                                    keyboardType: TextInputType.emailAddress,
                                    isDark: isDark,
                                  ),
                                ],
                                
                                const SizedBox(height: 24),
                                
                                _buildInputLabel('Mot de passe', Icons.lock_outline_rounded, isDark),
                                const SizedBox(height: 8),
                                _buildPasswordField(isDark),
                                
                                const SizedBox(height: 16),
                                _buildRememberForgotRow(isDark),
                                
                                if (_errorMessage != null) ...[
                                  const SizedBox(height: 16),
                                  _buildErrorBanner(),
                                ],
                                
                                const SizedBox(height: 32),
                                _buildLoginButton(isDark),
                                
                                if (_biometricsAvailable) ...[
                                  const SizedBox(height: 16),
                                  _buildBiometricButton(isDark),
                                ],
                                
                                const SizedBox(height: 32),
                                _buildRegisterLink(isDark),
                              ],
                            ),
                          ),
                        ),
                        
                        const SizedBox(height: 32),
                        _buildFooterLinks(isDark),
                      ],
                    ),
                  ),
                ),
              ),
            ),
            
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

  Widget _buildTabs(bool isDark) {
    return Container(
      padding: const EdgeInsets.all(4),
      decoration: BoxDecoration(
         color: isDark ? Colors.black26 : Colors.grey.shade100,
         borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
           Expanded(
             child: GestureDetector(
               onTap: () => setState(() => _loginMethod = 'phone'),
               child: Container(
                 padding: const EdgeInsets.symmetric(vertical: 10),
                 decoration: BoxDecoration(
                   color: _loginMethod == 'phone' 
                       ? (isDark ? const Color(0xFF6366F1) : Colors.white) 
                       : Colors.transparent,
                   borderRadius: BorderRadius.circular(8),
                   boxShadow: _loginMethod == 'phone' && !isDark ? [BoxShadow(color: Colors.black12, blurRadius: 4)] : [],
                 ),
                 alignment: Alignment.center,
                 child: Text("Téléphone", style: TextStyle(
                    fontWeight: FontWeight.bold,
                    color: _loginMethod == 'phone' 
                        ? (isDark ? Colors.white : Colors.black) 
                        : (isDark ? Colors.white54 : Colors.grey)
                 )),
               ),
             ),
           ),
           Expanded(
             child: GestureDetector(
               onTap: () => setState(() => _loginMethod = 'email'),
               child: Container(
                 padding: const EdgeInsets.symmetric(vertical: 10),
                 decoration: BoxDecoration(
                   color: _loginMethod == 'email' 
                       ? (isDark ? const Color(0xFF6366F1) : Colors.white) 
                       : Colors.transparent,
                   borderRadius: BorderRadius.circular(8),
                   boxShadow: _loginMethod == 'email' && !isDark ? [BoxShadow(color: Colors.black12, blurRadius: 4)] : [],
                 ),
                 alignment: Alignment.center,
                 child: Text("Email", style: TextStyle(
                    fontWeight: FontWeight.bold,
                    color: _loginMethod == 'email' 
                        ? (isDark ? Colors.white : Colors.black) 
                        : (isDark ? Colors.white54 : Colors.grey)
                 )),
               ),
             ),
           ),
        ],
      ),
    );
  }

  Widget _buildPhoneInput(bool isDark) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
         _buildInputLabel('Numéro de téléphone', Icons.phone_android_rounded, isDark),
         const SizedBox(height: 8),
         Row(
           children: [
             // Country Dial Code
             GestureDetector(
               onTap: () => _showCountryPicker(context, isDark),
               child: Container(
                 padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 16),
                 decoration: BoxDecoration(
                   color: isDark ? Colors.black.withOpacity(0.2) : const Color(0xFFF9FAFB),
                   borderRadius: BorderRadius.circular(12),
                   border: Border.all(color: isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFD1D5DB)),
                 ),
                 child: Row(
                   children: [
                      Text(_selectedCountry != null ? _getFlagEmoji(_selectedCountry!) : '🌍', style: const TextStyle(fontSize: 20)),
                      const SizedBox(width: 8),
                      SizedBox(
                        width: 50,
                        child: TextField(
                          controller: _dialCodeController,
                          enabled: false,
                          style: TextStyle(color: isDark ? Colors.white : Colors.black),
                          decoration: const InputDecoration(border: InputBorder.none, isDense: true, contentPadding: EdgeInsets.zero),
                        ),
                      ),
                      Icon(Icons.arrow_drop_down, color: isDark ? Colors.white70 : Colors.grey),
                   ],
                 ),
               ),
             ),
             const SizedBox(width: 12),
             Expanded(
               child: _buildTextField(
                  controller: _phoneController,
                  hint: 'Ex: 0102030405',
                  isDark: isDark,
                  keyboardType: TextInputType.phone
               ),
             ),
           ],
         ),
      ],
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
    
    final password = _passwordController.text;
    final totpCode = _totpController.text.trim();
    
    // Check password always
    if (password.isEmpty) {
      setState(() => _errorMessage = 'Veuillez entrer votre mot de passe.');
      return;
    }
    
    String? email;
    String? phone;
    
    if (_loginMethod == 'email') {
      email = _emailController.text.trim();
      if (email.isEmpty) {
        setState(() => _errorMessage = 'Veuillez entrer votre email.');
        return;
      }
    } else {
      final number = _phoneController.text.trim();
      if (number.isEmpty) {
         setState(() => _errorMessage = 'Veuillez entrer votre numéro de téléphone.');
         return;
      }
      // Combine dial code and number if needed, or send as is if backend handles it?
      // Backend expects 'phone'. Usually FULL E.164.
      // RegisterPage does logic. Here I'll just concat.
      // Note: DialCodeController has text like "+225".
      phone = '${_dialCodeController.text}$number';
    }
    
    context.read<AuthBloc>().add(
      SignInEvent(
        email: email,
        phone: phone,
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