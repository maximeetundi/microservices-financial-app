import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/glass_container.dart';
import 'package:google_fonts/google_fonts.dart';

class SettingsPage extends StatelessWidget {
  const SettingsPage({super.key});

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
        child: SafeArea(
          child: Column(
            children: [
              // Custom App Bar
               Padding(
                 padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                 child: Row(
                   mainAxisAlignment: MainAxisAlignment.center,
                   children: [
                     Text(
                        'Paramètres',
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
                child: ListView(
                  padding: const EdgeInsets.all(16),
                  children: [
                    // Profile Section
                    GlassContainer(
                      padding: const EdgeInsets.all(20),
                      borderRadius: 24,
                      child: Row(
                        children: [
                          Container(
                            width: 70,
                            height: 70,
                            decoration: BoxDecoration(
                              shape: BoxShape.circle,
                              gradient: AppTheme.primaryGradient,
                              boxShadow: [
                                BoxShadow(
                                  color: AppTheme.primaryColor.withOpacity(0.3),
                                  blurRadius: 10,
                                  offset: const Offset(0, 5),
                                ),
                              ],
                            ),
                            child: const Center(
                              child: Text('JD', style: TextStyle(fontSize: 24, color: Colors.white, fontWeight: FontWeight.bold)),
                            ),
                          ),
                          const SizedBox(width: 16),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  'John Doe', 
                                  style: GoogleFonts.inter(
                                    fontSize: 20, 
                                    fontWeight: FontWeight.bold,
                                    color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                                  ),
                                ),
                                const SizedBox(height: 4),
                                Text(
                                  'john.doe@email.com', 
                                  style: GoogleFonts.inter(
                                    color: isDark ? Colors.white70 : AppTheme.textSecondaryColor,
                                    fontSize: 14,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          IconButton(
                            icon: Icon(Icons.edit, color: isDark ? Colors.white70 : AppTheme.primaryColor),
                            onPressed: () => context.push('/more/profile'),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: 24),
                    
                    // Settings Groups
                    _SettingsGroup(
                      title: 'Compte',
                      items: [
                        _SettingsItem(icon: Icons.person_outline, title: 'Profil', onTap: () => context.push('/more/profile')),
                        _SettingsItem(icon: Icons.shield_outlined, title: 'Sécurité', onTap: () => context.push('/more/security')),
                        _SettingsItem(icon: Icons.notifications_outlined, title: 'Notifications', onTap: () {}),
                        _SettingsItem(icon: Icons.language, title: 'Langue', subtitle: 'Français', onTap: () {}),
                      ],
                    ),
                    
                    _SettingsGroup(
                      title: 'Préférences',
                      items: [
                        _SettingsItem(
                          icon: Icons.dark_mode_outlined, 
                          title: 'Thème sombre', 
                          trailing: Switch.adaptive(
                            value: isDark, 
                            onChanged: (_) {}, 
                            activeColor: AppTheme.primaryColor,
                          ),
                        ),
                        _SettingsItem(icon: Icons.currency_exchange, title: 'Devise par défaut', subtitle: 'USD', onTap: () {}),
                        _SettingsItem(
                          icon: Icons.fingerprint, 
                          title: 'Biométrie', 
                          trailing: Switch.adaptive(
                            value: true, 
                            onChanged: (_) {},
                            activeColor: AppTheme.primaryColor,
                          ),
                        ),
                      ],
                    ),
                    
                    _SettingsGroup(
                      title: 'Services',
                      items: [
                        _SettingsItem(icon: Icons.credit_card_outlined, title: 'Mes cartes', onTap: () => context.push('/more/cards')),
                        _SettingsItem(icon: Icons.send_outlined, title: 'Transferts', onTap: () => context.push('/more/transfer')),
                        _SettingsItem(icon: Icons.storefront_outlined, title: 'Espace Marchand', onTap: () => context.push('/more/merchant')),
                      ],
                    ),
                    
                    _SettingsGroup(
                      title: 'Support',
                      items: [
                        _SettingsItem(icon: Icons.help_outline, title: 'Centre d\'aide', onTap: () => context.push('/more/support')),
                        _SettingsItem(icon: Icons.chat_bubble_outline, title: 'Contacter le support', onTap: () => context.push('/more/support')),
                        _SettingsItem(icon: Icons.privacy_tip_outlined, title: 'Politique de confidentialité', onTap: () {}),
                        _SettingsItem(icon: Icons.description_outlined, title: 'Conditions d\'utilisation', onTap: () {}),
                      ],
                    ),
                    
                     _SettingsGroup(
                      title: 'Application',
                      items: [
                        _SettingsItem(icon: Icons.info_outline, title: 'Version', subtitle: '1.0.0', onTap: () {}),
                        _SettingsItem(icon: Icons.star_outline, title: 'Noter l\'application', onTap: () {}),
                      ],
                    ),
                    
                    const SizedBox(height: 16),
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 16),
                      child: GlassContainer(
                        padding: const EdgeInsets.all(4),
                        borderRadius: 16,
                        color: Colors.red.withOpacity(0.1),
                         child: InkWell(
                          onTap: () {
                            // Logout logic
                            context.go('/auth/login');
                          },
                          borderRadius: BorderRadius.circular(16),
                           child: Padding(
                             padding: const EdgeInsets.symmetric(vertical: 16),
                             child: Row(
                               mainAxisAlignment: MainAxisAlignment.center,
                               children: [
                                 const Icon(Icons.logout, color: Colors.red),
                                 const SizedBox(width: 8),
                                 Text(
                                   'Déconnexion', 
                                   style: GoogleFonts.inter(
                                     color: Colors.red,
                                     fontWeight: FontWeight.bold,
                                     fontSize: 16,
                                   ),
                                 ),
                               ],
                             ),
                           ),
                         ),
                      ),
                    ),
                    const SizedBox(height: 32),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _SettingsGroup extends StatelessWidget {
  final String title;
  final List<Widget> items;

  const _SettingsGroup({required this.title, required this.items});

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 16, 16, 12),
          child: Text(
            title, 
            style: GoogleFonts.inter(
              fontSize: 14, 
              fontWeight: FontWeight.bold, 
              color: isDark ? Colors.white70 : AppTheme.textSecondaryColor
            ),
          ),
        ),
        Container(
          margin: const EdgeInsets.symmetric(horizontal: 16),
          decoration: BoxDecoration(
             color: isDark ? Colors.white.withOpacity(0.05) : Colors.white,
             borderRadius: BorderRadius.circular(20),
             border: Border.all(color: isDark ? Colors.white.withOpacity(0.1) : Colors.grey.shade200),
          ),
          child: Column(
            children: items.asMap().entries.map((entry) {
              final index = entry.key;
              final item = entry.value;
              final isLast = index == items.length - 1;
              
              return Column(
                children: [
                  item,
                  if (!isLast)
                    Divider(
                      height: 1, 
                      color: isDark ? Colors.white.withOpacity(0.1) : Colors.grey.shade100,
                      indent: 56,
                    ),
                ],
              );
            }).toList(),
          ),
        ),
      ],
    );
  }
}

class _SettingsItem extends StatelessWidget {
  final IconData icon;
  final String title;
  final String? subtitle;
  final Widget? trailing;
  final VoidCallback? onTap;

  const _SettingsItem({
    required this.icon,
    required this.title,
    this.subtitle,
    this.trailing,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return ListTile(
      leading: Container(
        padding: const EdgeInsets.all(8),
        decoration: BoxDecoration(
          color: isDark ? Colors.white.withOpacity(0.1) : AppTheme.primaryColor.withOpacity(0.1),
          shape: BoxShape.circle,
        ),
        child: Icon(icon, color: isDark ? Colors.white : AppTheme.primaryColor, size: 20),
      ),
      title: Text(
        title, 
        style: GoogleFonts.inter(
          color: isDark ? Colors.white : AppTheme.textPrimaryColor,
          fontWeight: FontWeight.w500,
          fontSize: 15,
        ),
      ),
      subtitle: subtitle != null ? Text(
        subtitle!,
        style: GoogleFonts.inter(
          color: isDark ? Colors.white54 : AppTheme.textSecondaryColor,
          fontSize: 13,
        ),
      ) : null,
      trailing: trailing ?? Icon(
        Icons.chevron_right, 
        color: isDark ? Colors.white30 : Colors.grey.shade400,
        size: 20,
      ),
      onTap: onTap,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
    );
  }
}
