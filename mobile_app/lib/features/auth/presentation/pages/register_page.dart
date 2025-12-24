import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/custom_text_field.dart';
import '../../../../core/widgets/custom_button.dart';
import '../../../../core/widgets/glass_container.dart';
import '../bloc/auth_bloc.dart';

class RegisterPage extends StatefulWidget {
  const RegisterPage({super.key});

  @override
  State<RegisterPage> createState() => _RegisterPageState();
}

class _RegisterPageState extends State<RegisterPage> {
  final _formKey = GlobalKey<FormState>();
  final _firstNameController = TextEditingController();
  final _lastNameController = TextEditingController();
  final _emailController = TextEditingController();
  final _phoneController = TextEditingController();
  final _passwordController = TextEditingController();
  final _confirmPasswordController = TextEditingController();
  DateTime? _dateOfBirth;
  String? _selectedCountry;
  bool _obscurePassword = true;
  bool _obscureConfirmPassword = true;
  bool _acceptTerms = false;

  final List<Map<String, String>> _countries = [
    {'code': 'CI', 'name': 'Côte d\'Ivoire', 'currency': 'XOF'},
    {'code': 'SN', 'name': 'Sénégal', 'currency': 'XOF'},
    {'code': 'ML', 'name': 'Mali', 'currency': 'XOF'},
    {'code': 'BF', 'name': 'Burkina Faso', 'currency': 'XOF'},
    {'code': 'FR', 'name': 'France', 'currency': 'EUR'},
    {'code': 'US', 'name': 'États-Unis', 'currency': 'USD'},
    {'code': 'GB', 'name': 'Royaume-Uni', 'currency': 'GBP'},
  ];

  String? _getCurrency() {
    if (_selectedCountry == null) return null;
    final country = _countries.firstWhere(
      (c) => c['code'] == _selectedCountry,
      orElse: () => {'currency': 'USD'},
    );
    return country['currency'];
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
    }
  }

  @override
  Widget build(BuildContext context) {
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
                ? [const Color(0xFF020617), const Color(0xFF0F172A)] 
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
                        onPressed: () => context.go('/auth/login'),
                      ),
                    ),
                    const SizedBox(width: 16),
                    Text(
                      'Créer un compte',
                      style: GoogleFonts.inter(
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                        color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                      ),
                    ),
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
                            Text(
                              'Rejoignez Zekora',
                              style: GoogleFonts.inter(
                                fontSize: 28, 
                                fontWeight: FontWeight.bold,
                                color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                              ),
                            ),
                            const SizedBox(height: 8),
                            Text(
                              'Créez votre compte en quelques étapes',
                              style: GoogleFonts.inter(
                                fontSize: 14, 
                                color: isDark ? Colors.white70 : AppTheme.textSecondaryColor,
                              ),
                            ),
                            const SizedBox(height: 32),
                            
                            Row(
                              children: [
                                Expanded(
                                  child: CustomTextField(
                                    controller: _firstNameController,
                                    label: 'Prénom',
                                    prefixIcon: Icons.person_outline,
                                    validator: (v) => v!.isEmpty ? 'Requis' : null,
                                    fillColor: isDark ? Colors.white.withOpacity(0.05) : Colors.grey.shade50,
                                  ),
                                ),
                                const SizedBox(width: 16),
                                Expanded(
                                  child: CustomTextField(
                                    controller: _lastNameController,
                                    label: 'Nom',
                                    prefixIcon: Icons.person_outline,
                                    validator: (v) => v!.isEmpty ? 'Requis' : null,
                                    fillColor: isDark ? Colors.white.withOpacity(0.05) : Colors.grey.shade50,
                                  ),
                                ),
                              ],
                            ),
                            const SizedBox(height: 16),
                            
                            CustomTextField(
                              controller: _emailController,
                              keyboardType: TextInputType.emailAddress,
                              label: 'Email',
                              prefixIcon: Icons.email_outlined,
                              validator: (v) => v!.contains('@') ? null : 'Email invalide',
                              fillColor: isDark ? Colors.white.withOpacity(0.05) : Colors.grey.shade50,
                            ),
                            const SizedBox(height: 16),
                            
                            CustomTextField(
                              controller: _phoneController,
                              keyboardType: TextInputType.phone,
                              label: 'Téléphone',
                              prefixIcon: Icons.phone_outlined,
                              validator: (v) => v!.isEmpty ? 'Requis' : null,
                              fillColor: isDark ? Colors.white.withOpacity(0.05) : Colors.grey.shade50,
                            ),
                            const SizedBox(height: 16),
                            
                            // Date of Birth
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
                                  labelText: 'Date de naissance',
                                  prefixIcon: const Icon(Icons.cake_outlined),
                                  filled: true,
                                  fillColor: isDark ? Colors.white.withOpacity(0.05) : Colors.grey.shade50,
                                  border: OutlineInputBorder(
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
                            const SizedBox(height: 16),
                            
                            // Country
                            DropdownButtonFormField<String>(
                              value: _selectedCountry,
                              decoration: InputDecoration(
                                labelText: 'Pays',
                                prefixIcon: const Icon(Icons.public),
                                filled: true,
                                fillColor: isDark ? Colors.white.withOpacity(0.05) : Colors.grey.shade50,
                                border: OutlineInputBorder(
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
                              onChanged: (value) {
                                setState(() => _selectedCountry = value);
                              },
                              validator: (v) => v == null ? 'Requis' : null,
                            ),
                            const SizedBox(height: 16),
                            
                            CustomTextField(
                              controller: _passwordController,
                              obscureText: _obscurePassword,
                              label: 'Mot de passe',
                              prefixIcon: Icons.lock_outline,
                              suffixIcon: IconButton(
                                icon: Icon(_obscurePassword ? Icons.visibility : Icons.visibility_off,
                                  color: isDark ? Colors.white70 : AppTheme.textSecondaryColor,
                                ),
                                onPressed: () => setState(() => _obscurePassword = !_obscurePassword),
                              ),
                              validator: (v) => v!.length >= 8 ? 'Minimum 8 caractères' : null,
                              fillColor: isDark ? Colors.white.withOpacity(0.05) : Colors.grey.shade50,
                            ),
                            const SizedBox(height: 16),
                            
                            CustomTextField(
                              controller: _confirmPasswordController,
                              obscureText: _obscureConfirmPassword,
                              label: 'Confirmer le mot de passe',
                              prefixIcon: Icons.lock_outline,
                              suffixIcon: IconButton(
                                icon: Icon(_obscureConfirmPassword ? Icons.visibility : Icons.visibility_off,
                                  color: isDark ? Colors.white70 : AppTheme.textSecondaryColor,
                                ),
                                onPressed: () => setState(() => _obscureConfirmPassword = !_obscureConfirmPassword),
                              ),
                              validator: (v) => v == _passwordController.text ? null : 'Les mots de passe ne correspondent pas',
                              fillColor: isDark ? Colors.white.withOpacity(0.05) : Colors.grey.shade50,
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
                            const SizedBox(height: 24),
                            
                            BlocBuilder<AuthBloc, AuthState>(
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
                                    onPressed: state is AuthLoadingState ? null : _register,
                                    text: state is AuthLoadingState ? 'Création...' : 'Créer mon compte',
                                    isLoading: state is AuthLoadingState,
                                    backgroundColor: AppTheme.primaryColor,
                                    textColor: Colors.white,
                                  ),
                                );
                              },
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
}
