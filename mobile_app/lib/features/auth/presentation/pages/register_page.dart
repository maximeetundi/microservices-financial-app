import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/custom_text_field.dart';
import '../../../../core/widgets/custom_button.dart';
import '../../../../core/widgets/glass_container.dart';
import '../../../../core/widgets/glass_container.dart';
import '../bloc/auth_bloc.dart';
import 'package:google_fonts/google_fonts.dart';

class RegisterPage extends StatefulWidget {
  const RegisterPage({super.key});

  @override
  State<RegisterPage> createState() => _RegisterPageState();
}


class _RegisterPageState extends State<RegisterPage> {
  int _currentStep = 1;
  final _formKey = GlobalKey<FormState>();
  
  // Step 1: Identity
  final _firstNameController = TextEditingController();
  final _lastNameController = TextEditingController();
  DateTime? _dateOfBirth;
  
  // Step 2: Contact
  final _emailController = TextEditingController();
  final _phoneController = TextEditingController();
  String? _selectedCountry;
  
  // Step 3: Security
  final _passwordController = TextEditingController();
  final _confirmPasswordController = TextEditingController();
  bool _obscurePassword = true;
  bool _obscureConfirmPassword = true;
  bool _acceptTerms = false;

  final List<Map<String, String>> _countries = [
    {'code': 'CI', 'name': 'Côte d\'Ivoire', 'currency': 'XOF', 'dial_code': '+225'},
    {'code': 'SN', 'name': 'Sénégal', 'currency': 'XOF', 'dial_code': '+221'},
    {'code': 'ML', 'name': 'Mali', 'currency': 'XOF', 'dial_code': '+223'},
    {'code': 'BF', 'name': 'Burkina Faso', 'currency': 'XOF', 'dial_code': '+226'},
    {'code': 'FR', 'name': 'France', 'currency': 'EUR', 'dial_code': '+33'},
    {'code': 'US', 'name': 'États-Unis', 'currency': 'USD', 'dial_code': '+1'},
    {'code': 'GB', 'name': 'Royaume-Uni', 'currency': 'GBP', 'dial_code': '+44'},
  ];

  String? _getCurrency() {
    if (_selectedCountry == null) return null;
    final country = _countries.firstWhere(
      (c) => c['code'] == _selectedCountry,
      orElse: () => {'currency': 'USD'},
    );
    return country['currency'];
  }

  String _getDialCode(String countryCode) {
    final country = _countries.firstWhere(
      (c) => c['code'] == countryCode,
      orElse: () => {'dial_code': ''},
    );
    return country['dial_code'] ?? '';
  }

  void _onCountryChanged(String? value) {
    setState(() {
      _selectedCountry = value;
      if (value != null) {
        final dialCode = _getDialCode(value);
        if (_phoneController.text.isEmpty || _phoneController.text.startsWith('+')) {
           _phoneController.text = '$dialCode ';
        }
      }
    });
  }

  void _nextStep() {
    if (_currentStep == 1) {
      if (_firstNameController.text.isEmpty || _lastNameController.text.isEmpty || _dateOfBirth == null) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Veuillez remplir tous les champs')),
        );
        return;
      }
    } else if (_currentStep == 2) {
      if (_emailController.text.isEmpty || _selectedCountry == null || _phoneController.text.isEmpty) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Veuillez remplir tous les champs')),
        );
        return;
      }
      if (!_emailController.text.contains('@')) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Email invalide')),
        );
        return;
      }
    }
    
    setState(() => _currentStep++);
  }

  void _prevStep() {
    setState(() => _currentStep--);
  }

  @override
  void dispose() {
    _firstNameController.dispose();
    _lastNameController.dispose();
    _emailController.dispose();
    _phoneController.dispose();
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    super.dispose();
  }

  void _register() {
    if (_formKey.currentState!.validate() && _acceptTerms) {
      context.read<AuthBloc>().add(SignUpEvent(
        firstName: _firstNameController.text,
        lastName: _lastNameController.text,
        email: _emailController.text,
        phoneNumber: _phoneController.text,
        password: _passwordController.text,
        dateOfBirth: _dateOfBirth != null 
            ? '${_dateOfBirth!.toIso8601String().split('T')[0]}T00:00:00Z'
            : null,
        country: _selectedCountry,
        currency: _getCurrency(),
      ));
    } else if (!_acceptTerms) {
       ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Veuillez accepter les conditions d\'utilisation')),
        );
    }
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: Colors.transparent,
      body: Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: isDark 
                ? [const Color(0xFF110C2E), const Color(0xFF0F0C29)]
                : [const Color(0xFFFAFBFC), const Color(0xFFEFF6FF)],
          ),
        ),
        child: SafeArea(
          child: Column(
            children: [
               // Custom Back Button in App Bar area
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                child: Row(
                  children: [
                    GlassContainer(
                      padding: EdgeInsets.zero,
                      width: 40,
                      height: 40,
                      borderRadius: 12,
                      child: IconButton(
                        icon: Icon(Icons.arrow_back, color: isDark ? Colors.white : AppTheme.textPrimaryColor),
                        onPressed: () {
                           if (_currentStep > 1) {
                             _prevStep();
                           } else {
                             context.go('/auth/login');
                           }
                        },
                      ),
                    ),
                    const SizedBox(width: 16),
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Créer un compte',
                          style: GoogleFonts.inter(
                            fontSize: 20,
                            fontWeight: FontWeight.bold,
                            color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                          ),
                        ),
                        Text(
                          'Étape $_currentStep sur 3',
                          style: GoogleFonts.inter(
                            fontSize: 12,
                            color: isDark ? Colors.white70 : AppTheme.textSecondaryColor,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
              
              // Progress Bar
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                child: Row(
                  children: [
                    _buildStepIndicator(isDark, 1, 'Identité'),
                    _buildStepconnector(isDark, 1),
                    _buildStepIndicator(isDark, 2, 'Contact'),
                    _buildStepconnector(isDark, 2),
                    _buildStepIndicator(isDark, 3, 'Sécurité'),
                  ],
                ),
              ),

              Expanded(
                child: BlocListener<AuthBloc, AuthState>(
                  listener: (context, state) {
                    if (state is AuthenticatedState) {
                      context.go('/auth/biometric-setup');
                    } else if (state is AuthErrorState) {
                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(content: Text(state.message), backgroundColor: AppTheme.errorColor),
                      );
                    }
                  },
                  child: SingleChildScrollView(
                    padding: const EdgeInsets.all(24),
                    child: GlassContainer(
                      padding: const EdgeInsets.all(24),
                      borderRadius: 24,
                      child: Form(
                        key: _formKey,
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.stretch,
                          children: [
                            if (_currentStep == 1) ...[
                              Text(
                                'Informations personnelles',
                                style: GoogleFonts.inter(
                                  fontSize: 18,
                                  fontWeight: FontWeight.bold,
                                  color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                                ),
                              ),
                              const SizedBox(height: 24),
                              _buildStep1(isDark),
                            ] else if (_currentStep == 2) ...[
                              Text(
                                'Coordonnées',
                                style: GoogleFonts.inter(
                                  fontSize: 18,
                                  fontWeight: FontWeight.bold,
                                  color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                                ),
                              ),
                              const SizedBox(height: 24),
                              _buildStep2(isDark),
                            ] else ...[
                               Text(
                                'Sécurisation du compte',
                                style: GoogleFonts.inter(
                                  fontSize: 18,
                                  fontWeight: FontWeight.bold,
                                  color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                                ),
                              ),
                              const SizedBox(height: 24),
                              _buildStep3(isDark),
                            ],

                            const SizedBox(height: 32),
                            
                            // Navigation Buttons
                            Row(
                              children: [
                                if (_currentStep > 1)
                                  Expanded(
                                    child: Padding(
                                      padding: const EdgeInsets.only(right: 8.0),
                                      child: OutlinedButton(
                                        onPressed: _prevStep,
                                        style: OutlinedButton.styleFrom(
                                          padding: const EdgeInsets.symmetric(vertical: 16),
                                          side: BorderSide(color: isDark ? Colors.white30 : AppTheme.primaryColor),
                                          foregroundColor: isDark ? Colors.white : AppTheme.primaryColor,
                                          shape: RoundedRectangleBorder(
                                            borderRadius: BorderRadius.circular(12),
                                          ),
                                        ),
                                        child: const Text('Retour'),
                                      ),
                                    ),
                                  ),
                                Expanded(
                                  flex: 2,
                                  child: BlocBuilder<AuthBloc, AuthState>(
                                    builder: (context, state) {
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
                                          onPressed: state is AuthLoadingState 
                                              ? null 
                                              : (_currentStep < 3 ? _nextStep : _register),
                                          text: state is AuthLoadingState 
                                              ? 'Traitement...' 
                                              : (_currentStep < 3 ? 'Continuer' : 'Créer mon compte'),
                                          isLoading: state is AuthLoadingState,
                                          backgroundColor: AppTheme.primaryColor,
                                          textColor: Colors.white,
                                        ),
                                      );
                                    },
                                  ),
                                ),
                              ],
                            ),
                            
                            const SizedBox(height: 24),
                            
                            Row(
                              mainAxisAlignment: MainAxisAlignment.center,
                              children: [
                                Text(
                                  'Déjà un compte ?',
                                  style: GoogleFonts.inter(
                                    color: isDark ? Colors.white70 : AppTheme.textSecondaryColor,
                                  ),
                                ),
                                TextButton(
                                  onPressed: () => context.go('/auth/login'),
                                  child: Text(
                                    'Se connecter',
                                    style: GoogleFonts.inter(
                                      color: AppTheme.primaryColor,
                                      fontWeight: FontWeight.bold,
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
            ],
          ),
        ),
      ),
    );
  }

  // Step Indicators
  Widget _buildStepIndicator(bool isDark, int step, String label) {
    bool isActive = _currentStep >= step;
    return Column(
      children: [
        Container(
          width: 32,
          height: 32,
          decoration: BoxDecoration(
            shape: BoxShape.circle,
            color: isActive ? AppTheme.primaryColor : (isDark ? Colors.white10 : Colors.grey.shade200),
            gradient: isActive ? AppTheme.primaryGradient : null,
            boxShadow: isActive ? [
              BoxShadow(
                color: AppTheme.primaryColor.withOpacity(0.4),
                blurRadius: 10,
                offset: const Offset(0, 4)
              )
            ] : null,
          ),
          child: Center(
            child: isActive 
                ? const Icon(Icons.check, size: 16, color: Colors.white)
                : Text('$step', style: TextStyle(color: isDark ? Colors.white54 : Colors.grey.shade500, fontWeight: FontWeight.bold)),
          ),
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: GoogleFonts.inter(
            fontSize: 10,
            color: isActive 
                ? (isDark ? Colors.white : AppTheme.textPrimaryColor)
                : (isDark ? Colors.white24 : Colors.grey.shade400),
            fontWeight: isActive ? FontWeight.w600 : FontWeight.normal,
          ),
        ),
      ],
    );
  }

  Widget _buildStepconnector(bool isDark, int step) {
    bool isActive = _currentStep > step;
    return Expanded(
      child: Container(
        height: 2,
        margin: const EdgeInsets.symmetric(horizontal: 4, vertical: 14),
        color: isActive ? AppTheme.primaryColor : (isDark ? Colors.white10 : Colors.grey.shade200),
      ),
    );
  }

  // Step 1: Identity content
  Widget _buildStep1(bool isDark) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
         Row(
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _buildInputLabel(context, 'Prénom', Icons.person_outline),
                    const SizedBox(height: 8),
                    CustomTextField(
                      controller: _firstNameController,
                      hint: 'Prénom',
                      validator: (v) => v!.isEmpty ? 'Requis' : null,
                      fillColor: isDark ? const Color(0xFF0F0C29) : Colors.grey.shade50,
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _buildInputLabel(context, 'Nom', Icons.person_outline),
                    const SizedBox(height: 8),
                    CustomTextField(
                      controller: _lastNameController,
                      hint: 'Nom',
                      validator: (v) => v!.isEmpty ? 'Requis' : null,
                      fillColor: isDark ? const Color(0xFF0F0C29) : Colors.grey.shade50,
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          // Date of Birth
          _buildInputLabel(context, 'Date de naissance', Icons.cake_outlined),
          const SizedBox(height: 8),
          InkWell(
            onTap: () async {
              final picked = await showDatePicker(
                context: context,
                initialDate: DateTime(2000, 1, 1),
                firstDate: DateTime(1900),
                lastDate: DateTime.now(),
                builder: (context, child) {
                  return Theme(
                    data: Theme.of(context).copyWith(
                      colorScheme: isDark 
                          ? const ColorScheme.dark(primary: AppTheme.primaryColor)
                          : const ColorScheme.light(primary: AppTheme.primaryColor),
                    ),
                    child: child!,
                  );
                },
              );
              if (picked != null) {
                setState(() => _dateOfBirth = picked);
              }
            },
            child: InputDecorator(
              decoration: InputDecoration(
                hintText: 'Date de naissance',
                filled: true,
                fillColor: isDark ? const Color(0xFF0F0C29) : Colors.grey.shade50,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                  borderSide: BorderSide.none,
                ),
                enabledBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                  borderSide: BorderSide.none,
                ),
              ),
              child: Text(
                _dateOfBirth != null
                    ? '${_dateOfBirth!.day}/${_dateOfBirth!.month}/${_dateOfBirth!.year}'
                    : 'Sélectionner',
                style: TextStyle(
                  color: _dateOfBirth != null 
                      ? (isDark ? Colors.white : AppTheme.textPrimaryColor)
                      : (isDark ? Colors.white38 : Colors.grey[600]),
                ),
              ),
            ),
          ),
      ],
    );
  }

  // Step 2: Contact Content
  Widget _buildStep2(bool isDark) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        _buildInputLabel(context, 'Email', Icons.alternate_email_rounded),
          const SizedBox(height: 8),
          CustomTextField(
            controller: _emailController,
            keyboardType: TextInputType.emailAddress,
            hint: 'Entrez votre email',
            validator: (v) => v!.contains('@') ? null : 'Email invalide',
            fillColor: isDark ? const Color(0xFF0F0C29) : Colors.grey.shade50,
          ),
          const SizedBox(height: 16),

          // Country
          _buildInputLabel(context, 'Pays de résidence', Icons.public),
          const SizedBox(height: 8),
          DropdownButtonFormField<String>(
            value: _selectedCountry,
            decoration: InputDecoration(
              hintText: 'Sélectionner votre pays',
              filled: true,
              fillColor: isDark ? const Color(0xFF0F0C29) : Colors.grey.shade50,
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: BorderSide.none,
              ),
              enabledBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: BorderSide.none,
              ),
            ),
            dropdownColor: isDark ? const Color(0xFF1E293B) : Colors.white,
            items: _countries.map((country) {
              return DropdownMenuItem<String>(
                value: country['code'],
                child: Text(
                  country['name']!,
                  style: TextStyle(
                    color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                  ),
                ),
              );
            }).toList(),
            onChanged: _onCountryChanged,
            validator: (v) => v == null ? 'Requis' : null,
          ),
          const SizedBox(height: 16),
          
          // Phone
          _buildInputLabel(context, 'Téléphone', Icons.phone_outlined),
          const SizedBox(height: 8),
          CustomTextField(
            controller: _phoneController,
            keyboardType: TextInputType.phone,
            hint: 'Numéro de téléphone',
            validator: (v) => v!.isEmpty ? 'Requis' : null,
            fillColor: isDark ? const Color(0xFF0F0C29) : Colors.grey.shade50,
          ),
      ],
    );
  }

  // Step 3: Security Content
  Widget _buildStep3(bool isDark) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        _buildInputLabel(context, 'Mot de passe', Icons.lock_outline_rounded),
        const SizedBox(height: 8),
        CustomTextField(
          controller: _passwordController,
          obscureText: _obscurePassword,
          hint: '8 caractères minimum',
          suffixIcon: IconButton(
            icon: Icon(_obscurePassword ? Icons.visibility : Icons.visibility_off,
              color: isDark ? Colors.white70 : AppTheme.textSecondaryColor,
            ),
            onPressed: () => setState(() => _obscurePassword = !_obscurePassword),
          ),
          validator: (v) => v!.length >= 8 ? 'Minimum 8 caractères' : null,
          fillColor: isDark ? const Color(0xFF0F0C29) : Colors.grey.shade50,
        ),
        const SizedBox(height: 16),
        
        _buildInputLabel(context, 'Confirmer le mot de passe', Icons.lock_outline_rounded),
        const SizedBox(height: 8),
        CustomTextField(
          controller: _confirmPasswordController,
          obscureText: _obscureConfirmPassword,
          hint: 'Répétez le mot de passe',
          suffixIcon: IconButton(
            icon: Icon(_obscureConfirmPassword ? Icons.visibility : Icons.visibility_off,
              color: isDark ? Colors.white70 : AppTheme.textSecondaryColor,
            ),
            onPressed: () => setState(() => _obscureConfirmPassword = !_obscureConfirmPassword),
          ),
          validator: (v) => v == _passwordController.text ? null : 'Les mots de passe ne correspondent pas',
          fillColor: isDark ? const Color(0xFF0F0C29) : Colors.grey.shade50,
        ),
        const SizedBox(height: 24),
        
        CheckboxListTile(
          value: _acceptTerms,
          onChanged: (v) => setState(() => _acceptTerms = v!),
          title: Text(
            "J'accepte les conditions d'utilisation",
            style: GoogleFonts.inter(
              color: isDark ? Colors.white70 : AppTheme.textSecondaryColor,
              fontSize: 14,
            ),
          ),
          activeColor: AppTheme.primaryColor,
          side: BorderSide(
            color: isDark ? Colors.white60 : Colors.grey.shade400,
          ),
          controlAffinity: ListTileControlAffinity.leading,
          contentPadding: EdgeInsets.zero,
        ),
      ],
    );
  }

  Widget _buildInputLabel(BuildContext context, String label, IconData icon) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    return Row(
      children: [
        Icon(
          icon,
          size: 16,
          color: isDark ? AppTheme.primaryLightColor : AppTheme.primaryColor,
        ),
        const SizedBox(width: 8),
        Text(
          label,
          style: GoogleFonts.inter(
            fontSize: 14,
            fontWeight: FontWeight.w500,
            color: isDark ? Colors.white70 : AppTheme.textSecondaryColor,
          ),
        ),
      ],
    );
  }
}
