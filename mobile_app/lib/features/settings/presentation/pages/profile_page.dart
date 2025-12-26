import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/services/auth_api_service.dart';
import '../../../../core/services/pin_service.dart';
import '../../../../core/theme/app_theme.dart';
import '../../../auth/presentation/bloc/auth_bloc.dart';
import '../../../auth/presentation/pages/pin_setup_screen.dart';
import '../../../security/pin_code_screen.dart';

/// Profile page - READ ONLY for security (fraud prevention)
/// Requires PIN verification to access (like web version)
class ProfilePage extends StatefulWidget {
  const ProfilePage({super.key});

  @override
  State<ProfilePage> createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {
  final AuthApiService _authApi = AuthApiService();
  final PinService _pinService = PinService();
  
  bool _isUnlocked = false;
  bool _needsPinSetup = false;
  bool _loading = true;
  bool _checkingPin = true;
  
  Map<String, dynamic> _profile = {};
  
  // French month names
  static const List<String> _frenchMonths = [
    'janvier', 'février', 'mars', 'avril', 'mai', 'juin',
    'juillet', 'août', 'septembre', 'octobre', 'novembre', 'décembre'
  ];

  @override
  void initState() {
    super.initState();
    _checkPinStatus();
  }
  
  Future<void> _checkPinStatus() async {
    try {
      final hasPin = await _pinService.checkPinStatus();
      if (mounted) {
        setState(() {
          _needsPinSetup = !hasPin;
          _checkingPin = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _needsPinSetup = true;
          _checkingPin = false;
        });
      }
    }
  }
  
  Future<void> _loadProfile() async {
    setState(() => _loading = true);
    try {
      final response = await _authApi.getProfile();
      if (mounted && response.isNotEmpty) {
        setState(() {
          _profile = response;
          _loading = false;
        });
      }
    } catch (e) {
      debugPrint('Error loading profile: $e');
      if (mounted) {
        setState(() => _loading = false);
      }
    }
  }
  
  void _showPinVerification() {
    Navigator.of(context).push(
      MaterialPageRoute(
        builder: (context) => PinCodeScreen(
          isSetup: false,
          title: 'Entrez votre PIN',
          onSuccess: () {
            Navigator.of(context).pop();
            setState(() => _isUnlocked = true);
            _loadProfile();
          },
          onCancel: () {
            Navigator.of(context).pop();
            context.pop();
          },
        ),
      ),
    );
  }
  
  void _showPinSetup() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      isDismissible: false,
      enableDrag: false,
      backgroundColor: Colors.transparent,
      builder: (context) => PinSetupScreen(
        onPinSet: () {
          Navigator.of(context).pop();
          setState(() {
            _needsPinSetup = false;
          });
          _showPinVerification();
        },
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    // Still checking PIN status
    if (_checkingPin) {
      return Scaffold(
        backgroundColor: isDark ? const Color(0xFF0F172A) : const Color(0xFFF5F7FA),
        body: const Center(child: CircularProgressIndicator()),
      );
    }
    
    // Needs PIN setup
    if (_needsPinSetup) {
      return _buildPinSetupRequired(isDark);
    }
    
    // Not unlocked yet - show PIN verification prompt
    if (!_isUnlocked) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (mounted && !_isUnlocked) {
          _showPinVerification();
        }
      });
      return Scaffold(
        backgroundColor: isDark ? const Color(0xFF0F172A) : const Color(0xFFF5F7FA),
        body: const Center(child: CircularProgressIndicator()),
      );
    }
    
    // Profile content (after PIN verified)
    return Scaffold(
      backgroundColor: isDark ? const Color(0xFF0F172A) : const Color(0xFFF5F7FA),
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: Icon(Icons.arrow_back_ios, color: isDark ? Colors.white : const Color(0xFF1a1a2e)),
          onPressed: () => context.pop(),
        ),
        title: Text(
          '👤 Mon profil',
          style: TextStyle(
            color: isDark ? Colors.white : const Color(0xFF1a1a2e),
            fontWeight: FontWeight.bold,
            fontSize: 20,
          ),
        ),
        centerTitle: true,
      ),
      body: _loading 
        ? const Center(child: CircularProgressIndicator())
        : RefreshIndicator(
            onRefresh: _loadProfile,
            child: SingleChildScrollView(
              physics: const AlwaysScrollableScrollPhysics(),
              padding: const EdgeInsets.all(16),
              child: Column(
                children: [
                  // Security Notice
                  _buildSecurityNotice(isDark),
                  const SizedBox(height: 24),
                  
                  // Avatar & Basic Info
                  _buildProfileHeader(isDark),
                  const SizedBox(height: 24),
                  
                  // Personal Information
                  _buildSection(
                    title: '👤 INFORMATIONS PERSONNELLES',
                    isDark: isDark,
                    children: [
                      _buildInfoRow('Prénom', _profile['first_name'] ?? '—', isDark),
                      _buildInfoRow('Nom', _profile['last_name'] ?? '—', isDark),
                      _buildInfoRow('Email', _profile['email'] ?? '—', isDark),
                      _buildInfoRow('Téléphone', _profile['phone'] ?? '—', isDark),
                      _buildInfoRow('Date de naissance', _formatDate(_profile['date_of_birth']), isDark),
                      _buildInfoRow('Nationalité', _getCountryName(_profile['country']), isDark),
                    ],
                  ),
                  const SizedBox(height: 16),
                  
                  // Address
                  _buildSection(
                    title: '🏠 ADRESSE',
                    isDark: isDark,
                    children: [
                      _buildInfoRow('Adresse', _profile['address'] ?? 'Non renseignée', isDark),
                      _buildInfoRow('Ville', _profile['city'] ?? '—', isDark),
                      _buildInfoRow('Code postal', _profile['postal_code'] ?? '—', isDark),
                      _buildInfoRow('Pays', _getCountryName(_profile['country']), isDark),
                    ],
                  ),
                  const SizedBox(height: 16),
                  
                  // Account Info
                  _buildSection(
                    title: '📊 INFORMATIONS DU COMPTE',
                    isDark: isDark,
                    children: [
                      _buildInfoRow('Membre depuis', _formatDate(_profile['created_at']), isDark),
                      _buildInfoRow('Dernière connexion', _formatDate(_profile['last_login_at']), isDark),
                      _buildInfoRow('Niveau KYC', 'Niveau ${_profile['kyc_level'] ?? 0}', isDark, valueColor: AppTheme.primaryColor),
                      _buildInfoRow('Statut', _profile['is_active'] == true ? '✓ Actif' : '✗ Inactif', isDark, 
                        valueColor: _profile['is_active'] == true ? Colors.green : Colors.red),
                    ],
                  ),
                  const SizedBox(height: 16),
                  
                  // Verifications
                  _buildVerificationsSection(isDark),
                  const SizedBox(height: 24),
                  
                  // Contact Support
                  _buildSupportButton(isDark),
                  const SizedBox(height: 16),
                  
                  Text(
                    'Délai de traitement: 24-48h',
                    style: TextStyle(
                      fontSize: 12,
                      color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
                    ),
                  ),
                  const SizedBox(height: 32),
                ],
              ),
            ),
          ),
    );
  }
  
  Widget _buildPinSetupRequired(bool isDark) {
    return Scaffold(
      backgroundColor: const Color(0xFF0F172A),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(32),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Text('⚠️', style: TextStyle(fontSize: 64)),
              const SizedBox(height: 24),
              Text(
                'Configuration requise',
                style: GoogleFonts.inter(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                  color: Colors.white,
                ),
              ),
              const SizedBox(height: 12),
              Text(
                'Pour accéder à vos informations personnelles, vous devez d\'abord configurer votre PIN de sécurité.',
                textAlign: TextAlign.center,
                style: GoogleFonts.inter(
                  fontSize: 14,
                  color: const Color(0xFF94A3B8),
                ),
              ),
              const SizedBox(height: 32),
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: _showPinSetup,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFFF97316),
                    padding: const EdgeInsets.symmetric(vertical: 16),
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                  ),
                  child: const Text(
                    '🔐 Configurer mon PIN',
                    style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold, color: Colors.white),
                  ),
                ),
              ),
              const SizedBox(height: 16),
              TextButton(
                onPressed: () => context.pop(),
                child: const Text(
                  '← Retour aux paramètres',
                  style: TextStyle(color: Color(0xFF94A3B8)),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
  
  Widget _buildSecurityNotice(bool isDark) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xFF3B82F6).withOpacity(0.1),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: const Color(0xFF3B82F6).withOpacity(0.3)),
      ),
      child: Row(
        children: [
          const Text('🔒', style: TextStyle(fontSize: 24)),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  'Protection anti-fraude activée',
                  style: TextStyle(
                    fontWeight: FontWeight.bold,
                    color: Color(0xFF3B82F6),
                  ),
                ),
                Text(
                  'La modification des informations nécessite une vérification par le support.',
                  style: TextStyle(
                    fontSize: 12,
                    color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
  
  Widget _buildProfileHeader(bool isDark) {
    final firstName = _profile['first_name'] ?? '';
    final lastName = _profile['last_name'] ?? '';
    final email = _profile['email'] ?? '';
    final initials = '${firstName.isNotEmpty ? firstName[0] : ''}${lastName.isNotEmpty ? lastName[0] : ''}'.toUpperCase();
    
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: isDark ? const Color(0xFF1E293B) : Colors.white,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(
          color: isDark ? Colors.white.withOpacity(0.08) : const Color(0xFFE2E8F0),
        ),
      ),
      child: Row(
        children: [
          Container(
            width: 64,
            height: 64,
            decoration: BoxDecoration(
              gradient: const LinearGradient(
                colors: [Color(0xFF6366F1), Color(0xFF8B5CF6)],
              ),
              borderRadius: BorderRadius.circular(16),
            ),
            child: Center(
              child: Text(
                initials.isNotEmpty ? initials : '?',
                style: const TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                  color: Colors.white,
                ),
              ),
            ),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  '$firstName $lastName'.trim().isNotEmpty ? '$firstName $lastName' : 'Utilisateur',
                  style: GoogleFonts.inter(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  email,
                  style: GoogleFonts.inter(
                    fontSize: 14,
                    color: const Color(0xFF64748B),
                  ),
                ),
                const SizedBox(height: 8),
                Wrap(
                  spacing: 8,
                  runSpacing: 4,
                  children: [
                    if (_profile['email_verified'] == true)
                      _buildBadge('✓ Email vérifié', Colors.green, isDark),
                    if (_profile['phone_verified'] == true)
                      _buildBadge('✓ Téléphone', Colors.green, isDark),
                    if (_profile['kyc_status'] == 'verified')
                      _buildBadge('✓ KYC', Colors.green, isDark)
                    else
                      _buildBadge('⏳ KYC en attente', Colors.orange, isDark),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
  
  Widget _buildBadge(String text, Color color, bool isDark) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: color.withOpacity(0.15),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Text(
        text,
        style: TextStyle(
          fontSize: 10,
          fontWeight: FontWeight.bold,
          color: color,
        ),
      ),
    );
  }
  
  Widget _buildSection({
    required String title,
    required bool isDark,
    required List<Widget> children,
  }) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: isDark ? const Color(0xFF1E293B) : Colors.white,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(
          color: isDark ? Colors.white.withOpacity(0.08) : const Color(0xFFE2E8F0),
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            title,
            style: GoogleFonts.inter(
              fontSize: 11,
              fontWeight: FontWeight.bold,
              letterSpacing: 0.5,
              color: const Color(0xFF64748B),
            ),
          ),
          const SizedBox(height: 16),
          ...children,
        ],
      ),
    );
  }
  
  Widget _buildInfoRow(String label, String value, bool isDark, {Color? valueColor}) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: GoogleFonts.inter(
              fontSize: 13,
              color: const Color(0xFF64748B),
            ),
          ),
          Text(
            value,
            style: GoogleFonts.inter(
              fontSize: 13,
              fontWeight: FontWeight.w500,
              color: valueColor ?? (isDark ? Colors.white : const Color(0xFF1E293B)),
            ),
          ),
        ],
      ),
    );
  }
  
  Widget _buildVerificationsSection(bool isDark) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: isDark ? const Color(0xFF1E293B) : Colors.white,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(
          color: isDark ? Colors.white.withOpacity(0.08) : const Color(0xFFE2E8F0),
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '✅ VÉRIFICATIONS',
            style: GoogleFonts.inter(
              fontSize: 11,
              fontWeight: FontWeight.bold,
              letterSpacing: 0.5,
              color: const Color(0xFF64748B),
            ),
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(child: _buildVerifyItem('📧', 'Email', _profile['email_verified'] == true, isDark)),
              const SizedBox(width: 12),
              Expanded(child: _buildVerifyItem('📱', 'Téléphone', _profile['phone_verified'] == true, isDark)),
            ],
          ),
          const SizedBox(height: 12),
          Row(
            children: [
              Expanded(child: _buildVerifyItem('🛡️', '2FA', _profile['two_fa_enabled'] == true, isDark)),
              const SizedBox(width: 12),
              Expanded(child: _buildVerifyItem('🔑', 'PIN', _profile['has_pin'] == true, isDark)),
            ],
          ),
        ],
      ),
    );
  }
  
  Widget _buildVerifyItem(String icon, String label, bool verified, bool isDark) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: isDark ? Colors.white.withOpacity(0.03) : const Color(0xFFF8FAFC),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          Text(icon, style: const TextStyle(fontSize: 20)),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              label,
              style: GoogleFonts.inter(
                fontSize: 13,
                fontWeight: FontWeight.w500,
                color: isDark ? Colors.white : const Color(0xFF334155),
              ),
            ),
          ),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
            decoration: BoxDecoration(
              color: (verified ? Colors.green : const Color(0xFF64748B)).withOpacity(0.15),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Text(
              verified ? 'Vérifié' : 'En attente',
              style: TextStyle(
                fontSize: 10,
                fontWeight: FontWeight.bold,
                color: verified ? Colors.green : const Color(0xFF64748B),
              ),
            ),
          ),
        ],
      ),
    );
  }
  
  Widget _buildSupportButton(bool isDark) {
    return Container(
      width: double.infinity,
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [Color(0xFF6366F1), Color(0xFF8B5CF6)],
        ),
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: const Color(0xFF6366F1).withOpacity(0.3),
            blurRadius: 20,
            offset: const Offset(0, 8),
          ),
        ],
      ),
      child: ElevatedButton(
        onPressed: () => context.push('/support'),
        style: ElevatedButton.styleFrom(
          backgroundColor: Colors.transparent,
          shadowColor: Colors.transparent,
          padding: const EdgeInsets.symmetric(vertical: 16),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        ),
        child: const Text(
          '📞 Contacter le support pour modifier',
          style: TextStyle(
            color: Colors.white,
            fontWeight: FontWeight.bold,
            fontSize: 14,
          ),
        ),
      ),
    );
  }
  
  String _formatDate(dynamic date) {
    if (date == null) return '—';
    try {
      final DateTime dt = date is DateTime ? date : DateTime.parse(date.toString());
      return '${dt.day} ${_frenchMonths[dt.month - 1]} ${dt.year}';
    } catch (e) {
      return '—';
    }
  }
  
  String _getCountryName(String? code) {
    if (code == null || code.isEmpty) return '—';
    const countries = {
      'SEN': 'Sénégal', 'CIV': "Côte d'Ivoire", 'FRA': 'France',
      'USA': 'États-Unis', 'GBR': 'Royaume-Uni', 'MLI': 'Mali',
      'BFA': 'Burkina Faso', 'NGA': 'Nigeria', 'GHA': 'Ghana',
      'CMR': 'Cameroun', 'BEN': 'Bénin', 'TGO': 'Togo'
    };
    return countries[code] ?? code;
  }
}
