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

  // Expanded Country List (Matches Frontend)
  // Helper for Flags
  String _getFlagEmoji(String countryCode) {
    return countryCode.toUpperCase().replaceAllMapped(RegExp(r'[A-Z]'),
        (match) => String.fromCharCode(match.group(0)!.codeUnitAt(0) + 127397));
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
            final filtered = _countries.where((c) {
               final q = searchQuery.toLowerCase();
               return c['name']!.toLowerCase().contains(q) || c['dial_code']!.contains(q);
            }).toList();

            return Container(
              height: MediaQuery.of(context).size.height * 0.7,
              padding: const EdgeInsets.all(24),
              child: Column(
                children: [
                  // Handle
                  Container(
                    width: 40, height: 4,
                    margin: const EdgeInsets.only(bottom: 24),
                    decoration: BoxDecoration(
                      color: isDark ? Colors.white24 : Colors.grey.shade300,
                      borderRadius: BorderRadius.circular(2),
                    ),
                  ),
                  // Search
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
                  // List
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

  final List<Map<String, String>> _countries = [
    {'code': 'AF', 'name': 'Afghanistan', 'currency': 'AFN', 'dial_code': '+93'},
    {'code': 'ZA', 'name': 'Afrique du Sud', 'currency': 'ZAR', 'dial_code': '+27'},
    {'code': 'AL', 'name': 'Albanie', 'currency': 'ALL', 'dial_code': '+355'},
    {'code': 'DZ', 'name': 'Algérie', 'currency': 'DZD', 'dial_code': '+213'},
    {'code': 'DE', 'name': 'Allemagne', 'currency': 'EUR', 'dial_code': '+49'},
    {'code': 'AD', 'name': 'Andorre', 'currency': 'EUR', 'dial_code': '+376'},
    {'code': 'AO', 'name': 'Angola', 'currency': 'AOA', 'dial_code': '+244'},
    {'code': 'SA', 'name': 'Arabie Saoudite', 'currency': 'SAR', 'dial_code': '+966'},
    {'code': 'AR', 'name': 'Argentine', 'currency': 'ARS', 'dial_code': '+54'},
    {'code': 'AM', 'name': 'Arménie', 'currency': 'AMD', 'dial_code': '+374'},
    {'code': 'AU', 'name': 'Australie', 'currency': 'AUD', 'dial_code': '+61'},
    {'code': 'AT', 'name': 'Autriche', 'currency': 'EUR', 'dial_code': '+43'},
    {'code': 'AZ', 'name': 'Azerbaïdjan', 'currency': 'AZN', 'dial_code': '+994'},
    {'code': 'BS', 'name': 'Bahamas', 'currency': 'BSD', 'dial_code': '+1'},
    {'code': 'BH', 'name': 'Bahreïn', 'currency': 'BHD', 'dial_code': '+973'},
    {'code': 'BD', 'name': 'Bangladesh', 'currency': 'BDT', 'dial_code': '+880'},
    {'code': 'BB', 'name': 'Barbade', 'currency': 'BBD', 'dial_code': '+1'},
    {'code': 'BE', 'name': 'Belgique', 'currency': 'EUR', 'dial_code': '+32'},
    {'code': 'BZ', 'name': 'Belize', 'currency': 'BZD', 'dial_code': '+501'},
    {'code': 'BJ', 'name': 'Bénin', 'currency': 'XOF', 'dial_code': '+229'},
    {'code': 'BT', 'name': 'Bhoutan', 'currency': 'BTN', 'dial_code': '+975'},
    {'code': 'BY', 'name': 'Biélorussie', 'currency': 'BYN', 'dial_code': '+375'},
    {'code': 'BO', 'name': 'Bolivie', 'currency': 'BOB', 'dial_code': '+591'},
    {'code': 'BA', 'name': 'Bosnie-Herzégovine', 'currency': 'BAM', 'dial_code': '+387'},
    {'code': 'BW', 'name': 'Botswana', 'currency': 'BWP', 'dial_code': '+267'},
    {'code': 'BR', 'name': 'Brésil', 'currency': 'BRL', 'dial_code': '+55'},
    {'code': 'BN', 'name': 'Brunei', 'currency': 'BND', 'dial_code': '+673'},
    {'code': 'BG', 'name': 'Bulgarie', 'currency': 'BGN', 'dial_code': '+359'},
    {'code': 'BF', 'name': 'Burkina Faso', 'currency': 'XOF', 'dial_code': '+226'},
    {'code': 'BI', 'name': 'Burundi', 'currency': 'BIF', 'dial_code': '+257'},
    {'code': 'KH', 'name': 'Cambodge', 'currency': 'KHR', 'dial_code': '+855'},
    {'code': 'CM', 'name': 'Cameroun', 'currency': 'XAF', 'dial_code': '+237'},
    {'code': 'CA', 'name': 'Canada', 'currency': 'CAD', 'dial_code': '+1'},
    {'code': 'CV', 'name': 'Cap-Vert', 'currency': 'CVE', 'dial_code': '+238'},
    {'code': 'CF', 'name': 'Centrafrique', 'currency': 'XAF', 'dial_code': '+236'},
    {'code': 'CL', 'name': 'Chili', 'currency': 'CLP', 'dial_code': '+56'},
    {'code': 'CN', 'name': 'Chine', 'currency': 'CNY', 'dial_code': '+86'},
    {'code': 'CY', 'name': 'Chypre', 'currency': 'EUR', 'dial_code': '+357'},
    {'code': 'CO', 'name': 'Colombie', 'currency': 'COP', 'dial_code': '+57'},
    {'code': 'KM', 'name': 'Comores', 'currency': 'KMF', 'dial_code': '+269'},
    {'code': 'CG', 'name': 'Congo-Brazzaville', 'currency': 'XAF', 'dial_code': '+242'},
    {'code': 'CD', 'name': 'Congo-Kinshasa', 'currency': 'CDF', 'dial_code': '+243'},
    {'code': 'KR', 'name': 'Corée du Sud', 'currency': 'KRW', 'dial_code': '+82'},
    {'code': 'CR', 'name': 'Costa Rica', 'currency': 'CRC', 'dial_code': '+506'},
    {'code': 'CI', 'name': 'Côte d\'Ivoire', 'currency': 'XOF', 'dial_code': '+225'},
    {'code': 'HR', 'name': 'Croatie', 'currency': 'EUR', 'dial_code': '+385'},
    {'code': 'CU', 'name': 'Cuba', 'currency': 'CUP', 'dial_code': '+53'},
    {'code': 'DK', 'name': 'Danemark', 'currency': 'DKK', 'dial_code': '+45'},
    {'code': 'DJ', 'name': 'Djibouti', 'currency': 'DJF', 'dial_code': '+253'},
    {'code': 'DM', 'name': 'Dominique', 'currency': 'XCD', 'dial_code': '+1'},
    {'code': 'EG', 'name': 'Égypte', 'currency': 'EGP', 'dial_code': '+20'},
    {'code': 'AE', 'name': 'Émirats Arabes Unis', 'currency': 'AED', 'dial_code': '+971'},
    {'code': 'EC', 'name': 'Équateur', 'currency': 'USD', 'dial_code': '+593'},
    {'code': 'ER', 'name': 'Érythrée', 'currency': 'ERN', 'dial_code': '+291'},
    {'code': 'ES', 'name': 'Espagne', 'currency': 'EUR', 'dial_code': '+34'},
    {'code': 'EE', 'name': 'Estonie', 'currency': 'EUR', 'dial_code': '+372'},
    {'code': 'US', 'name': 'États-Unis', 'currency': 'USD', 'dial_code': '+1'},
    {'code': 'ET', 'name': 'Éthiopie', 'currency': 'ETB', 'dial_code': '+251'},
    {'code': 'FJ', 'name': 'Fidji', 'currency': 'FJD', 'dial_code': '+679'},
    {'code': 'FI', 'name': 'Finlande', 'currency': 'EUR', 'dial_code': '+358'},
    {'code': 'FR', 'name': 'France', 'currency': 'EUR', 'dial_code': '+33'},
    {'code': 'GA', 'name': 'Gabon', 'currency': 'XAF', 'dial_code': '+241'},
    {'code': 'GM', 'name': 'Gambie', 'currency': 'GMD', 'dial_code': '+220'},
    {'code': 'GE', 'name': 'Géorgie', 'currency': 'GEL', 'dial_code': '+995'},
    {'code': 'GH', 'name': 'Ghana', 'currency': 'GHS', 'dial_code': '+233'},
    {'code': 'GR', 'name': 'Grèce', 'currency': 'EUR', 'dial_code': '+30'},
    {'code': 'GD', 'name': 'Grenade', 'currency': 'XCD', 'dial_code': '+1'},
    {'code': 'GT', 'name': 'Guatemala', 'currency': 'GTQ', 'dial_code': '+502'},
    {'code': 'GN', 'name': 'Guinée', 'currency': 'GNF', 'dial_code': '+224'},
    {'code': 'GW', 'name': 'Guinée-Bissau', 'currency': 'XOF', 'dial_code': '+245'},
    {'code': 'GQ', 'name': 'Guinée Équatoriale', 'currency': 'XAF', 'dial_code': '+240'},
    {'code': 'GY', 'name': 'Guyana', 'currency': 'GYD', 'dial_code': '+592'},
    {'code': 'HT', 'name': 'Haïti', 'currency': 'HTG', 'dial_code': '+509'},
    {'code': 'HN', 'name': 'Honduras', 'currency': 'HNL', 'dial_code': '+504'},
    {'code': 'HU', 'name': 'Hongrie', 'currency': 'HUF', 'dial_code': '+36'},
    {'code': 'IN', 'name': 'Inde', 'currency': 'INR', 'dial_code': '+91'},
    {'code': 'ID', 'name': 'Indonésie', 'currency': 'IDR', 'dial_code': '+62'},
    {'code': 'IQ', 'name': 'Irak', 'currency': 'IQD', 'dial_code': '+964'},
    {'code': 'IR', 'name': 'Iran', 'currency': 'IRR', 'dial_code': '+98'},
    {'code': 'IE', 'name': 'Irlande', 'currency': 'EUR', 'dial_code': '+353'},
    {'code': 'IS', 'name': 'Islande', 'currency': 'ISK', 'dial_code': '+354'},
    {'code': 'IL', 'name': 'Israël', 'currency': 'ILS', 'dial_code': '+972'},
    {'code': 'IT', 'name': 'Italie', 'currency': 'EUR', 'dial_code': '+39'},
    {'code': 'JM', 'name': 'Jamaïque', 'currency': 'JMD', 'dial_code': '+1'},
    {'code': 'JP', 'name': 'Japon', 'currency': 'JPY', 'dial_code': '+81'},
    {'code': 'JO', 'name': 'Jordanie', 'currency': 'JOD', 'dial_code': '+962'},
    {'code': 'KZ', 'name': 'Kazakhstan', 'currency': 'KZT', 'dial_code': '+7'},
    {'code': 'KE', 'name': 'Kenya', 'currency': 'KES', 'dial_code': '+254'},
    {'code': 'KG', 'name': 'Kirghizistan', 'currency': 'KGS', 'dial_code': '+996'},
    {'code': 'KW', 'name': 'Koweït', 'currency': 'KWD', 'dial_code': '+965'},
    {'code': 'LA', 'name': 'Laos', 'currency': 'LAK', 'dial_code': '+856'},
    {'code': 'LS', 'name': 'Lesotho', 'currency': 'LSL', 'dial_code': '+266'},
    {'code': 'LV', 'name': 'Lettonie', 'currency': 'EUR', 'dial_code': '+371'},
    {'code': 'LB', 'name': 'Liban', 'currency': 'LBP', 'dial_code': '+961'},
    {'code': 'LR', 'name': 'Liberia', 'currency': 'LRD', 'dial_code': '+231'},
    {'code': 'LY', 'name': 'Libye', 'currency': 'LYD', 'dial_code': '+218'},
    {'code': 'LI', 'name': 'Liechtenstein', 'currency': 'CHF', 'dial_code': '+423'},
    {'code': 'LT', 'name': 'Lituanie', 'currency': 'EUR', 'dial_code': '+370'},
    {'code': 'LU', 'name': 'Luxembourg', 'currency': 'EUR', 'dial_code': '+352'},
    {'code': 'MK', 'name': 'Macédoine du Nord', 'currency': 'MKD', 'dial_code': '+389'},
    {'code': 'MG', 'name': 'Madagascar', 'currency': 'MGA', 'dial_code': '+261'},
    {'code': 'MY', 'name': 'Malaisie', 'currency': 'MYR', 'dial_code': '+60'},
    {'code': 'MW', 'name': 'Malawi', 'currency': 'MWK', 'dial_code': '+265'},
    {'code': 'MV', 'name': 'Maldives', 'currency': 'MVR', 'dial_code': '+960'},
    {'code': 'ML', 'name': 'Mali', 'currency': 'XOF', 'dial_code': '+223'},
    {'code': 'MT', 'name': 'Malte', 'currency': 'EUR', 'dial_code': '+356'},
    {'code': 'MA', 'name': 'Maroc', 'currency': 'MAD', 'dial_code': '+212'},
    {'code': 'MU', 'name': 'Maurice', 'currency': 'MUR', 'dial_code': '+230'},
    {'code': 'MR', 'name': 'Mauritanie', 'currency': 'MRU', 'dial_code': '+222'},
    {'code': 'MX', 'name': 'Mexique', 'currency': 'MXN', 'dial_code': '+52'},
    {'code': 'MD', 'name': 'Moldavie', 'currency': 'MDL', 'dial_code': '+373'},
    {'code': 'MC', 'name': 'Monaco', 'currency': 'EUR', 'dial_code': '+377'},
    {'code': 'MN', 'name': 'Mongolie', 'currency': 'MNT', 'dial_code': '+976'},
    {'code': 'ME', 'name': 'Monténégro', 'currency': 'EUR', 'dial_code': '+382'},
    {'code': 'MZ', 'name': 'Mozambique', 'currency': 'MZN', 'dial_code': '+258'},
    {'code': 'NA', 'name': 'Namibie', 'currency': 'NAD', 'dial_code': '+264'},
    {'code': 'NP', 'name': 'Népal', 'currency': 'NPR', 'dial_code': '+977'},
    {'code': 'NI', 'name': 'Nicaragua', 'currency': 'NIO', 'dial_code': '+505'},
    {'code': 'NE', 'name': 'Niger', 'currency': 'XOF', 'dial_code': '+227'},
    {'code': 'NG', 'name': 'Nigéria', 'currency': 'NGN', 'dial_code': '+234'},
    {'code': 'NO', 'name': 'Norvège', 'currency': 'NOK', 'dial_code': '+47'},
    {'code': 'NZ', 'name': 'Nouvelle-Zélande', 'currency': 'NZD', 'dial_code': '+64'},
    {'code': 'OM', 'name': 'Oman', 'currency': 'OMR', 'dial_code': '+968'},
    {'code': 'UG', 'name': 'Ouganda', 'currency': 'UGX', 'dial_code': '+256'},
    {'code': 'UZ', 'name': 'Ouzbékistan', 'currency': 'UZS', 'dial_code': '+998'},
    {'code': 'PK', 'name': 'Pakistan', 'currency': 'PKR', 'dial_code': '+92'},
    {'code': 'PA', 'name': 'Panama', 'currency': 'PAB', 'dial_code': '+507'},
    {'code': 'PG', 'name': 'Papouasie-Nouvelle-Guinée', 'currency': 'PGK', 'dial_code': '+675'},
    {'code': 'PY', 'name': 'Paraguay', 'currency': 'PYG', 'dial_code': '+595'},
    {'code': 'NL', 'name': 'Pays-Bas', 'currency': 'EUR', 'dial_code': '+31'},
    {'code': 'PE', 'name': 'Pérou', 'currency': 'PEN', 'dial_code': '+51'},
    {'code': 'PH', 'name': 'Philippines', 'currency': 'PHP', 'dial_code': '+63'},
    {'code': 'PL', 'name': 'Pologne', 'currency': 'PLN', 'dial_code': '+48'},
    {'code': 'PT', 'name': 'Portugal', 'currency': 'EUR', 'dial_code': '+351'},
    {'code': 'QA', 'name': 'Qatar', 'currency': 'QAR', 'dial_code': '+974'},
    {'code': 'RO', 'name': 'Roumanie', 'currency': 'RON', 'dial_code': '+40'},
    {'code': 'GB', 'name': 'Royaume-Uni', 'currency': 'GBP', 'dial_code': '+44'},
    {'code': 'RU', 'name': 'Russie', 'currency': 'RUB', 'dial_code': '+7'},
    {'code': 'RW', 'name': 'Rwanda', 'currency': 'RWF', 'dial_code': '+250'},
    {'code': 'KN', 'name': 'Saint-Kitts-et-Nevis', 'currency': 'XCD', 'dial_code': '+1'},
    {'code': 'SM', 'name': 'Saint-Marin', 'currency': 'EUR', 'dial_code': '+378'},
    {'code': 'LC', 'name': 'Sainte-Lucie', 'currency': 'XCD', 'dial_code': '+1'},
    {'code': 'SV', 'name': 'Salvador', 'currency': 'USD', 'dial_code': '+503'},
    {'code': 'WS', 'name': 'Samoa', 'currency': 'WST', 'dial_code': '+685'},
    {'code': 'ST', 'name': 'Sao Tomé-et-Principe', 'currency': 'STN', 'dial_code': '+239'},
    {'code': 'SN', 'name': 'Sénégal', 'currency': 'XOF', 'dial_code': '+221'},
    {'code': 'RS', 'name': 'Serbie', 'currency': 'RSD', 'dial_code': '+381'},
    {'code': 'SC', 'name': 'Seychelles', 'currency': 'SCR', 'dial_code': '+248'},
    {'code': 'SL', 'name': 'Sierra Leone', 'currency': 'SLL', 'dial_code': '+232'},
    {'code': 'SG', 'name': 'Singapour', 'currency': 'SGD', 'dial_code': '+65'},
    {'code': 'SK', 'name': 'Slovaquie', 'currency': 'EUR', 'dial_code': '+421'},
    {'code': 'SI', 'name': 'Slovénie', 'currency': 'EUR', 'dial_code': '+386'},
    {'code': 'SO', 'name': 'Somalie', 'currency': 'SOS', 'dial_code': '+252'},
    {'code': 'SD', 'name': 'Soudan', 'currency': 'SDG', 'dial_code': '+249'},
    {'code': 'SS', 'name': 'Soudan du Sud', 'currency': 'SSP', 'dial_code': '+211'},
    {'code': 'LK', 'name': 'Sri Lanka', 'currency': 'LKR', 'dial_code': '+94'},
    {'code': 'SE', 'name': 'Suède', 'currency': 'SEK', 'dial_code': '+46'},
    {'code': 'CH', 'name': 'Suisse', 'currency': 'CHF', 'dial_code': '+41'},
    {'code': 'SR', 'name': 'Suriname', 'currency': 'SRD', 'dial_code': '+597'},
    {'code': 'SY', 'name': 'Syrie', 'currency': 'SYP', 'dial_code': '+963'},
    {'code': 'TJ', 'name': 'Tadjikistan', 'currency': 'TJS', 'dial_code': '+992'},
    {'code': 'TW', 'name': 'Taïwan', 'currency': 'TWD', 'dial_code': '+886'},
    {'code': 'TZ', 'name': 'Tanzanie', 'currency': 'TZS', 'dial_code': '+255'},
    {'code': 'TD', 'name': 'Tchad', 'currency': 'XAF', 'dial_code': '+235'},
    {'code': 'CZ', 'name': 'Tchéquie', 'currency': 'CZK', 'dial_code': '+420'},
    {'code': 'TH', 'name': 'Thaïlande', 'currency': 'THB', 'dial_code': '+66'},
    {'code': 'TL', 'name': 'Timor oriental', 'currency': 'USD', 'dial_code': '+670'},
    {'code': 'TG', 'name': 'Togo', 'currency': 'XOF', 'dial_code': '+228'},
    {'code': 'TO', 'name': 'Tonga', 'currency': 'TOP', 'dial_code': '+676'},
    {'code': 'TT', 'name': 'Trinité-et-Tobago', 'currency': 'TTD', 'dial_code': '+1'},
    {'code': 'TN', 'name': 'Tunisie', 'currency': 'TND', 'dial_code': '+216'},
    {'code': 'TM', 'name': 'Turkménistan', 'currency': 'TMT', 'dial_code': '+993'},
    {'code': 'TR', 'name': 'Turquie', 'currency': 'TRY', 'dial_code': '+90'},
    {'code': 'TV', 'name': 'Tuvalu', 'currency': 'AUD', 'dial_code': '+688'},
    {'code': 'UA', 'name': 'Ukraine', 'currency': 'UAH', 'dial_code': '+380'},
    {'code': 'UY', 'name': 'Uruguay', 'currency': 'UYU', 'dial_code': '+598'},
    {'code': 'VU', 'name': 'Vanuatu', 'currency': 'VUV', 'dial_code': '+678'},
    {'code': 'VE', 'name': 'Venezuela', 'currency': 'VES', 'dial_code': '+58'},
    {'code': 'VN', 'name': 'Vietnam', 'currency': 'VND', 'dial_code': '+84'},
    {'code': 'YE', 'name': 'Yémen', 'currency': 'YER', 'dial_code': '+967'},
    {'code': 'ZM', 'name': 'Zambie', 'currency': 'ZMW', 'dial_code': '+260'},
    {'code': 'ZW', 'name': 'Zimbabwe', 'currency': 'ZWL', 'dial_code': '+263'},
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

  final _dialCodeController = TextEditingController();

  void _onCountryChanged(String? value) {
    setState(() {
      _selectedCountry = value;
      if (value != null) {
        final dialCode = _getDialCode(value);
        if (_dialCodeController.text != dialCode) {
           _dialCodeController.text = dialCode;
        }
      }
    });
  }

  void _onDialCodeChanged(String value) {
    if (!value.startsWith('+')) {
       value = '+$value';
    }
    
    // Prioritize main countries for shared codes
    final Map<String, String> priorityMap = {
        '+1': 'US',
        '+7': 'RU',
        '+44': 'GB',
        '+39': 'IT',
        '+33': 'FR',
    };

    String? matchedCountryCode;
    
    // Check priority first
    if (priorityMap.containsKey(value)) {
      matchedCountryCode = priorityMap[value];
    } else {
      // Find any match
      final country = _countries.firstWhere(
        (c) => c['dial_code'] == value,
        orElse: () => {},
      );
      if (country.isNotEmpty) {
        matchedCountryCode = country['code'];
      }
    }

    if (matchedCountryCode != null && matchedCountryCode != _selectedCountry) {
       setState(() {
         _selectedCountry = matchedCountryCode;
       });
    }
  }

  // Simple As-You-Type formatter (adds space every 2 digits after first 2 or 3?)
  // Let's use a generic blocking: XXX XXX XXX XX
  void _onPhoneChanged(String value) {
     // Remove non-digits
     String clean = value.replaceAll(RegExp(r'[^\d]'), '');
     if (clean.isEmpty) return;
     
     // Simple grouping of 2 chars (common in FR) or 3 (US)?
     // Let's adapt based on length.
     // If length > 10, maybe 3-3-4?
     // If length <= 10, maybe 2-2-2-2-2?
     // Let's just do groups of 2 for now as it's visually cleaner than nothing.
     // Or mimic the frontend: it uses AsYouType which is smart.
     // We can't be perfectly smart without the lib.
     // Let's just do a visual spacer every 3 chars for now?
     // Or actually, just let user type but allow spaces.
     // Wait, I said I will implement it.
     
     // Smartish formatting:
     String formatted = '';
     for (int i = 0; i < clean.length; i++) {
       if (i > 0 && i % 2 == 0 && i < 12) { // 2 2 2 2 2 format
         formatted += ' ';
       }
       formatted += clean[i];
     }
     
     // Only update if visually different (to avoid cursor issues if not careful)
     // Handling cursor is tricky in standard onChanged. 
     // Let's just set the text and move cursor to end.
     if (_phoneController.text != formatted) {
       _phoneController.value = TextEditingValue(
         text: formatted,
         selection: TextSelection.collapsed(offset: formatted.length),
       );
     }
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
    _dialCodeController.dispose();
    super.dispose();
  }

  void _register() {
    if (_formKey.currentState!.validate() && _acceptTerms) {
      context.read<AuthBloc>().add(SignUpEvent(
        firstName: _firstNameController.text,
        lastName: _lastNameController.text,
        email: _emailController.text,
        // Concatenate DialCode + Phone (remove spaces)
        phoneNumber: '${_dialCodeController.text}${_phoneController.text.replaceAll(" ", "")}',
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
                      context.go('/auth/pin-setup');
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

          // Country Picker Custom UI
          _buildInputLabel(context, 'Pays de résidence', Icons.public), // Changed icon in buildInputLabel if needed, but param is icon
          const SizedBox(height: 8),
          InkWell(
            onTap: () => _showCountryPicker(context, isDark),
            borderRadius: BorderRadius.circular(12),
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
              decoration: BoxDecoration(
                color: isDark ? const Color(0xFF0F0C29) : Colors.grey.shade50,
                borderRadius: BorderRadius.circular(12),
                border: Border.all(
                    color: isDark ? Colors.transparent : Colors.transparent, // Or match CustomTextField border
                ),
              ),
              child: Row(
                children: [
                  if (_selectedCountry != null) ...[
                      Text(_getFlagEmoji(_selectedCountry!), style: const TextStyle(fontSize: 20)),
                      const SizedBox(width: 12),
                  ],
                  Expanded(
                    child: Text(
                      _selectedCountry != null 
                          ? _countries.firstWhere((c) => c['code'] == _selectedCountry)['name']!
                          : 'Sélectionner votre pays',
                      style: TextStyle(
                        color: _selectedCountry != null 
                            ? (isDark ? Colors.white : AppTheme.textPrimaryColor)
                            : (isDark ? Colors.white38 : Colors.grey[600]),
                        fontSize: 16,
                      ),
                    ),
                  ),
                  Icon(Icons.keyboard_arrow_down, color: isDark ? Colors.white54 : Colors.grey),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),
          
          // Phone Input Row
          _buildInputLabel(context, 'Téléphone', Icons.phone_outlined),
          const SizedBox(height: 8),
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
               // Country Dial Code & Picker
               // We remove the Dropdown overlay since we have the main country picker above.
               // But we still keep this as a visual indicator or a quick way to see code.
               // User specifically asked for "Exactly the same", on Web we have:
               // [Flag+Label] [Select Box]
               // [Phone Label] [DialCode] [Number]
               // So here we keep DialCode simple.
               SizedBox(
                 width: 100,
                 child: CustomTextField(
                    controller: _dialCodeController,
                    hint: '+33',
                    keyboardType: TextInputType.phone,
                    onChanged: _onDialCodeChanged,
                    fillColor: isDark ? const Color(0xFF0F0C29) : Colors.grey.shade50,
                    textAlign: TextAlign.center,
                  ),
               ),
               const SizedBox(width: 12),
               // National Number
               Expanded(
                 child: CustomTextField(
                    controller: _phoneController,
                    keyboardType: TextInputType.phone,
                    hint: '6 12 34 56 78',
                    onChanged: _onPhoneChanged,
                    validator: (v) => v!.isEmpty ? 'Requis' : null,
                    fillColor: isDark ? const Color(0xFF0F0C29) : Colors.grey.shade50,
                 ),
               ),
            ],
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
          validator: (v) => v != null && v.length < 8 ? 'Minimum 8 caractères' : null,
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
