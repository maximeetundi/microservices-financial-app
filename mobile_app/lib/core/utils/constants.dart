class AppConstants {
  // App Information
  static const String appName = 'Zekora';
  static const String appVersion = '1.0.0';
  static const String appDescription = 'Your Digital Banking Solution';
  
  // API Configuration
  static const String baseUrl = 'http://localhost:8080/api/v1';
  static const String websocketUrl = 'ws://localhost:8080/ws';
  
  // Storage Keys
  static const String accessTokenKey = 'access_token';
  static const String refreshTokenKey = 'refresh_token';
  static const String userDataKey = 'user_data';
  static const String biometricEnabledKey = 'biometric_enabled';
  static const String themeKey = 'theme_mode';
  static const String languageKey = 'language';
  static const String currencyKey = 'preferred_currency';
  
  // Pagination
  static const int defaultPageSize = 20;
  static const int maxPageSize = 100;
  
  // Validation
  static const int minPasswordLength = 8;
  static const int maxPasswordLength = 128;
  static const int pinLength = 4;
  static const int totpLength = 6;
  
  // Timeouts
  static const int connectionTimeout = 30; // seconds
  static const int receiveTimeout = 30; // seconds
  static const int sendTimeout = 30; // seconds
  
  // Rate Limits
  static const int maxLoginAttempts = 5;
  static const int loginCooldownMinutes = 15;
  
  // Transaction Limits
  static const double minTransferAmount = 0.01;
  static const double maxDailyTransferLimit = 100000.0;
  static const double maxMonthlyTransferLimit = 1000000.0;
  
  // Currency Codes
  static const List<String> supportedCryptoCurrencies = [
    'BTC',
    'ETH',
    'LTC',
    'XRP',
    'ADA',
    'DOT',
    'LINK',
    'BCH',
    'XLM',
    'USDT',
    'USDC',
    'BNB',
    'SOL',
    'AVAX',
    'MATIC',
    'ATOM',
  ];
  
  static const List<String> supportedFiatCurrencies = [
    'USD',
    'EUR',
    'GBP',
    'JPY',
    'CAD',
    'AUD',
    'CHF',
    'SEK',
    'NOK',
    'DKK',
    'PLN',
    'CZK',
    'HUF',
    'SGD',
    'HKD',
    'NZD',
    'MXN',
    'BRL',
  ];
  
  // Card Types
  static const List<String> cardTypes = [
    'virtual',
    'physical',
    'premium',
  ];
  
  // Transaction Types
  static const List<String> transactionTypes = [
    'send',
    'receive',
    'exchange',
    'buy',
    'sell',
    'card_payment',
    'card_topup',
    'fee',
  ];
  
  // Order Types
  static const List<String> orderTypes = [
    'market',
    'limit',
    'stop_loss',
    'stop_limit',
  ];
  
  // Order Sides
  static const List<String> orderSides = [
    'buy',
    'sell',
  ];
  
  // Status Types
  static const List<String> statusTypes = [
    'pending',
    'processing',
    'completed',
    'failed',
    'cancelled',
    'expired',
  ];
  
  // KYC Levels
  static const List<String> kycLevels = [
    'none',
    'basic',
    'intermediate',
    'advanced',
  ];
  
  // Chart Timeframes
  static const List<String> chartTimeframes = [
    '1h',
    '4h',
    '1d',
    '1w',
    '1m',
    '3m',
    '1y',
  ];
  
  // Regular Expressions
  static const String emailRegex = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$';
  static const String phoneRegex = r'^\+?[1-9]\d{1,14}$';
  static const String passwordRegex = r'^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]';
  static const String cryptoAddressRegex = r'^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$|^0x[a-fA-F0-9]{40}$';
  
  // Error Messages
  static const String networkError = 'Network connection error. Please check your internet connection.';
  static const String serverError = 'Server error occurred. Please try again later.';
  static const String unauthorizedError = 'Session expired. Please log in again.';
  static const String validationError = 'Please check your input and try again.';
  static const String biometricError = 'Biometric authentication failed. Please try again.';
  
  // Success Messages
  static const String loginSuccess = 'Login successful';
  static const String registerSuccess = 'Registration successful';
  static const String transferSuccess = 'Transfer completed successfully';
  static const String exchangeSuccess = 'Exchange completed successfully';
  static const String orderSuccess = 'Order placed successfully';
  
  // Formatting
  static const int decimalPlaces = 8;
  static const int fiatDecimalPlaces = 2;
  static const String dateFormat = 'dd/MM/yyyy';
  static const String timeFormat = 'HH:mm';
  static const String dateTimeFormat = 'dd/MM/yyyy HH:mm';
  
  // Animation Durations
  static const Duration shortAnimation = Duration(milliseconds: 200);
  static const Duration mediumAnimation = Duration(milliseconds: 400);
  static const Duration longAnimation = Duration(milliseconds: 600);
  
  // Debounce Durations
  static const Duration searchDebounce = Duration(milliseconds: 500);
  static const Duration refreshDebounce = Duration(milliseconds: 1000);
  
  // Cache Durations
  static const Duration shortCache = Duration(minutes: 5);
  static const Duration mediumCache = Duration(minutes: 30);
  static const Duration longCache = Duration(hours: 24);
  
  // WebSocket Events
  static const String priceUpdateEvent = 'price_update';
  static const String balanceUpdateEvent = 'balance_update';
  static const String orderUpdateEvent = 'order_update';
  static const String transactionUpdateEvent = 'transaction_update';
  static const String notificationEvent = 'notification';
  
  // Push Notification Categories
  static const String transactionCategory = 'transaction';
  static const String securityCategory = 'security';
  static const String marketCategory = 'market';
  static const String promotionCategory = 'promotion';
  
  // Deep Link Schemes
  static const String deepLinkScheme = 'zekora';
  static const String universalLinkDomain = 'zekora.com';
  
  // Feature Flags
  static const bool enableBiometrics = true;
  static const bool enablePushNotifications = true;
  static const bool enableWebSocket = true;
  static const bool enableCrashlytics = true;
  static const bool enableAnalytics = true;
  
  // Debug Settings
  static const bool enableDebugMode = false;
  static const bool enableNetworkLogs = false;
  static const bool enablePerformanceLogs = false;

  static const List<Map<String, String>> countries = [
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
}