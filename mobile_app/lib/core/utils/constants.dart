class AppConstants {
  // App Information
  static const String appName = 'Crypto Bank';
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
  static const String deepLinkScheme = 'cryptobank';
  static const String universalLinkDomain = 'cryptobank.com';
  
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
}