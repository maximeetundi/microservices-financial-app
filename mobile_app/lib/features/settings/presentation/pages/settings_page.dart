import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/glass_container.dart';

/// Settings Page matching web design exactly
class SettingsPage extends StatefulWidget {
  const SettingsPage({super.key});

  @override
  State<SettingsPage> createState() => _SettingsPageState();
}

class _SettingsPageState extends State<SettingsPage> {
  bool _showDeleteModal = false;
  int _securityScore = 60;
  bool _kycVerified = false;

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Scaffold(
      backgroundColor: isDark ? const Color(0xFF0F172A) : const Color(0xFFF8FAFC),
      body: Stack(
        children: [
          SafeArea(
            child: ListView(
              padding: const EdgeInsets.all(16),
              children: [
                // Header
                _buildHeader(isDark),
                const SizedBox(height: 24),
                
                // Settings Grid
                _buildSettingsGrid(isDark),
                const SizedBox(height: 32),
                
                // Quick Actions
                _buildQuickActions(isDark),
                const SizedBox(height: 32),
                
                // App Info
                _buildAppInfo(isDark),
              ],
            ),
          ),
          
          // Delete Modal
          if (_showDeleteModal) _buildDeleteModal(isDark),
        ],
      ),
    );
  }

  Widget _buildHeader(bool isDark) {
    return Row(
      children: [
        GlassContainer(
          padding: EdgeInsets.zero,
          width: 40,
          height: 40,
          borderRadius: 12,
          child: IconButton(
            icon: Icon(Icons.arrow_back_ios_new, size: 20, 
                color: isDark ? Colors.white : AppTheme.textPrimaryColor),
            onPressed: () => context.go('/dashboard'),
          ),
        ),
        const SizedBox(width: 16),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                '⚙️ Paramètres',
                style: GoogleFonts.inter(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                  color: isDark ? Colors.white : const Color(0xFF1E293B),
                ),
              ),
              Text(
                'Gérez votre compte et vos préférences',
                style: GoogleFonts.inter(
                  fontSize: 14,
                  color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildSettingsGrid(bool isDark) {
    return Column(
      children: [
        _buildSettingsCard(
          emoji: '👤',
          title: 'Profil',
          subtitle: 'Informations personnelles, coordonnées',
          color: const Color(0xFF3B82F6),
          onTap: () => context.go('/more/profile'),
          isDark: isDark,
        ),
        _buildSettingsCard(
          emoji: '🔒',
          title: 'Sécurité',
          subtitle: 'Mot de passe, 2FA, sessions',
          color: const Color(0xFF10B981),
          trailing: Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
            decoration: BoxDecoration(
              color: _securityScore >= 80 
                  ? (isDark ? const Color(0xFF22C55E).withOpacity(0.15) : const Color(0xFFF0FDF4))
                  : (isDark ? const Color(0xFFEF4444).withOpacity(0.15) : const Color(0xFFFEF2F2)),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Text(
              '$_securityScore%',
              style: GoogleFonts.inter(
                fontSize: 12,
                fontWeight: FontWeight.w600,
                color: _securityScore >= 80 ? const Color(0xFF22C55E) : const Color(0xFFEF4444),
              ),
            ),
          ),
          onTap: () => context.go('/more/security'),
          isDark: isDark,
        ),
        _buildSettingsCard(
          emoji: '📋',
          title: 'Vérification KYC',
          subtitle: 'Documents d\'identité, validation',
          color: const Color(0xFFA855F7),
          trailing: Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
            decoration: BoxDecoration(
              color: _kycVerified 
                  ? (isDark ? const Color(0xFF22C55E).withOpacity(0.15) : const Color(0xFFF0FDF4))
                  : (isDark ? const Color(0xFFF97316).withOpacity(0.15) : const Color(0xFFFFF7ED)),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Text(
              _kycVerified ? 'VÉRIFIÉ' : 'EN ATTENTE',
              style: GoogleFonts.inter(
                fontSize: 10,
                fontWeight: FontWeight.bold,
                color: _kycVerified ? const Color(0xFF22C55E) : const Color(0xFFF97316),
              ),
            ),
          ),
          onTap: () => context.go('/more/kyc'),
          isDark: isDark,
        ),
        _buildSettingsCard(
          emoji: '🎨',
          title: 'Préférences',
          subtitle: 'Thème, langue, notifications',
          color: const Color(0xFFF97316),
          onTap: () => context.go('/more/preferences'),
          isDark: isDark,
        ),
        _buildSettingsCard(
          emoji: '🔔',
          title: 'Notifications',
          subtitle: 'Alertes email, push, SMS',
          color: const Color(0xFFEC4899),
          onTap: () => context.go('/notifications'),
          isDark: isDark,
        ),
        _buildSettingsCard(
          emoji: '💳',
          title: 'Moyens de paiement',
          subtitle: 'Cartes, comptes bancaires',
          color: const Color(0xFF14B8A6),
          onTap: () => context.go('/more/payment-methods'),
          isDark: isDark,
        ),
      ],
    );
  }

  Widget _buildSettingsCard({
    required String emoji,
    required String title,
    required String subtitle,
    required Color color,
    Widget? trailing,
    required VoidCallback onTap,
    required bool isDark,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        margin: const EdgeInsets.only(bottom: 12),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: isDark ? Colors.white.withOpacity(0.03) : Colors.white,
          borderRadius: BorderRadius.circular(16),
          border: Border.all(
            color: isDark ? Colors.white.withOpacity(0.08) : const Color(0xFFE2E8F0),
          ),
        ),
        child: Row(
          children: [
            Container(
              width: 48,
              height: 48,
              decoration: BoxDecoration(
                color: color.withOpacity(isDark ? 0.15 : 0.1),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Center(
                child: Text(emoji, style: const TextStyle(fontSize: 24)),
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: GoogleFonts.inter(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color: isDark ? Colors.white : const Color(0xFF1E293B),
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    subtitle,
                    style: GoogleFonts.inter(
                      fontSize: 12,
                      color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                    ),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                ],
              ),
            ),
            if (trailing != null) ...[
              trailing,
              const SizedBox(width: 8),
            ],
            Text(
              '→',
              style: TextStyle(
                fontSize: 20,
                color: isDark ? const Color(0xFF475569) : const Color(0xFFCBD5E1),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildQuickActions(bool isDark) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'ACTIONS RAPIDES',
          style: GoogleFonts.inter(
            fontSize: 12,
            fontWeight: FontWeight.w600,
            letterSpacing: 1.2,
            color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
          ),
        ),
        const SizedBox(height: 12),
        Wrap(
          spacing: 12,
          runSpacing: 12,
          children: [
            _buildQuickButton(
              '📥 Exporter mes données',
              () {},
              isDark,
              false,
            ),
            _buildQuickButton(
              '🗑️ Supprimer mon compte',
              () => setState(() => _showDeleteModal = true),
              isDark,
              true,
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildQuickButton(String label, VoidCallback onTap, bool isDark, bool isDanger) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        decoration: BoxDecoration(
          color: isDark 
              ? Colors.transparent 
              : (isDanger ? Colors.transparent : Colors.white),
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: isDanger 
                ? (isDark ? const Color(0xFFEF4444).withOpacity(0.3) : const Color(0xFFFECACA))
                : (isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFE2E8F0)),
          ),
        ),
        child: Text(
          label,
          style: GoogleFonts.inter(
            fontSize: 14,
            fontWeight: FontWeight.w500,
            color: isDanger 
                ? const Color(0xFFEF4444)
                : (isDark ? Colors.white : const Color(0xFF1E293B)),
          ),
        ),
      ),
    );
  }

  Widget _buildAppInfo(bool isDark) {
    return Column(
      children: [
        Text(
          'Zekora v1.0.0',
          style: GoogleFonts.inter(
            fontSize: 12,
            color: isDark ? const Color(0xFF475569) : const Color(0xFF94A3B8),
          ),
        ),
        const SizedBox(height: 8),
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            GestureDetector(
              onTap: () => context.go('/more/support'),
              child: Text(
                'Aide',
                style: GoogleFonts.inter(
                  fontSize: 12,
                  color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
                ),
              ),
            ),
            Text(' • ', style: TextStyle(color: isDark ? const Color(0xFF475569) : const Color(0xFFCBD5E1))),
            Text(
              'Conditions',
              style: GoogleFonts.inter(
                fontSize: 12,
                color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
              ),
            ),
            Text(' • ', style: TextStyle(color: isDark ? const Color(0xFF475569) : const Color(0xFFCBD5E1))),
            Text(
              'Confidentialité',
              style: GoogleFonts.inter(
                fontSize: 12,
                color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
              ),
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildDeleteModal(bool isDark) {
    return GestureDetector(
      onTap: () => setState(() => _showDeleteModal = false),
      child: Container(
        color: Colors.black.withOpacity(0.5),
        child: Center(
          child: GestureDetector(
            onTap: () {},
            child: Container(
              margin: const EdgeInsets.all(24),
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(
                color: isDark ? const Color(0xFF1A1A2E) : Colors.white,
                borderRadius: BorderRadius.circular(20),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.2),
                    blurRadius: 20,
                  ),
                ],
              ),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text(
                    '⚠️ Supprimer votre compte',
                    style: GoogleFonts.inter(
                      fontSize: 20,
                      fontWeight: FontWeight.bold,
                      color: isDark ? Colors.white : const Color(0xFF1E293B),
                    ),
                  ),
                  const SizedBox(height: 12),
                  Text(
                    'Cette action est irréversible. Toutes vos données seront supprimées.',
                    style: GoogleFonts.inter(
                      fontSize: 14,
                      color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                    ),
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: 24),
                  Row(
                    children: [
                      Expanded(
                        child: GestureDetector(
                          onTap: () => setState(() => _showDeleteModal = false),
                          child: Container(
                            padding: const EdgeInsets.symmetric(vertical: 14),
                            decoration: BoxDecoration(
                              color: isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFF1F5F9),
                              borderRadius: BorderRadius.circular(12),
                            ),
                            child: Center(
                              child: Text(
                                'Annuler',
                                style: GoogleFonts.inter(
                                  fontSize: 14,
                                  fontWeight: FontWeight.w600,
                                  color: isDark ? Colors.white : const Color(0xFF1E293B),
                                ),
                              ),
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: GestureDetector(
                          onTap: () {
                            setState(() => _showDeleteModal = false);
                            ScaffoldMessenger.of(context).showSnackBar(
                              const SnackBar(content: Text('Veuillez contacter le support.')),
                            );
                          },
                          child: Container(
                            padding: const EdgeInsets.symmetric(vertical: 14),
                            decoration: BoxDecoration(
                              color: const Color(0xFFEF4444),
                              borderRadius: BorderRadius.circular(12),
                            ),
                            child: Center(
                              child: Text(
                                'Supprimer',
                                style: GoogleFonts.inter(
                                  fontSize: 14,
                                  fontWeight: FontWeight.w600,
                                  color: Colors.white,
                                ),
                              ),
                            ),
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
    );
  }
}
