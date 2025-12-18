/// API Endpoints for all microservices
class ApiEndpoints {
  // Base paths for each microservice
  static const String auth = '/api/v1/auth';
  static const String users = '/api/v1/users';
  static const String wallets = '/api/v1/wallets';
  static const String transfers = '/api/v1/transfers';
  static const String cards = '/api/v1/cards';
  static const String exchange = '/api/v1/exchange';
  static const String notifications = '/api/v1/notifications';
  
  // Auth Service Endpoints
  static const String login = '$auth/login';
  static const String register = '$auth/register';
  static const String logout = '$auth/logout';
  static const String refreshToken = '$auth/refresh';
  static const String verifyEmail = '$auth/verify-email';
  static const String forgotPassword = '$auth/forgot-password';
  static const String resetPassword = '$auth/reset-password';
  static const String enable2FA = '$auth/2fa/enable';
  static const String verify2FA = '$auth/2fa/verify';
  
  // User Endpoints
  static const String profile = '$users/me';
  static const String updateProfile = '$users/me';
  static const String changePassword = '$users/me/password';
  static const String uploadAvatar = '$users/me/avatar';
  static const String kycStatus = '$users/me/kyc';
  static const String updateKYC = '$users/me/kyc';
  static const String lookup = '$users/lookup'; // Added lookup endpoint
  
  // Wallet Service Endpoints
  static String walletsList = wallets;
  static String createWallet = wallets;
  static String walletById(String id) => '$wallets/$id';
  static String walletTransactions(String id) => '$wallets/$id/transactions';
  static String walletDeposit(String id) => '$wallets/$id/deposit';
  static String walletWithdraw(String id) => '$wallets/$id/withdraw';
  static String walletBalance(String id) => '$wallets/$id/balance';
  
  // Transfer Service Endpoints
  static String transfersList = transfers;
  static String createTransfer = transfers;
  static String transferById(String id) => '$transfers/$id';
  static String cancelTransfer(String id) => '$transfers/$id/cancel';
  static const String internationalTransfer = '$transfers/international';
  static const String mobileMoneyProviders = '$transfers/mobile/providers';
  static const String sendMobileMoney = '$transfers/mobile/send';
  static const String receiveMobileMoney = '$transfers/mobile/receive';
  
  // Card Service Endpoints
  static String cardsList = cards;
  static String createCard = cards;
  static String cardById(String id) => '$cards/$id';
  static String activateCard(String id) => '$cards/$id/activate';
  static String deactivateCard(String id) => '$cards/$id/deactivate';
  static String freezeCard(String id) => '$cards/$id/freeze';
  static String unfreezeCard(String id) => '$cards/$id/unfreeze';
  static String blockCard(String id) => '$cards/$id/block';
  static String loadCard(String id) => '$cards/$id/load';
  static String setCardPIN(String id) => '$cards/$id/pin';
  static String cardLimits(String id) => '$cards/$id/limits';
  static String cardTransactions(String id) => '$cards/$id/transactions';
  static String cardBalance(String id) => '$cards/$id/balance';
  static const String virtualCard = '$cards/virtual';
  static const String orderPhysicalCard = '$cards/physical/order';
  static String shippingStatus(String id) => '$cards/$id/shipping';
  static const String giftCards = '$cards/gift';
  static const String redeemGiftCard = '$cards/gift/redeem';
  
  // Exchange Service Endpoints
  static String exchangeRates = '$exchange/rates';
  static String exchangePair(String from, String to) => '$exchange/rate/$from/$to';
  static String quote = '$exchange/quote'; // Added quote endpoint
  static String executeExchange = '$exchange/execute';
  static String exchangeHistory = '$exchange/history';
  static String exchangeById(String id) => '$exchange/$id';
  
  // Notification Endpoints
  static String notificationsList = notifications;
  static String markAsRead(String id) => '$notifications/$id/read';
  static const String markAllAsRead = '$notifications/read-all';
  static const String notificationSettings = '$notifications/settings';
}
