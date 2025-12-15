import 'package:get_it/get_it.dart';

import '../api/api_client.dart';
import '../services/api_service.dart';
import '../services/secure_storage_service.dart';
import '../../features/auth/presentation/bloc/auth_bloc.dart';
import '../../features/wallet/presentation/bloc/wallet_bloc.dart';
import '../../features/exchange/presentation/bloc/exchange_bloc.dart';
import '../../features/cards/presentation/bloc/cards_bloc.dart';
import '../../features/portfolio/presentation/bloc/portfolio_bloc.dart';

final sl = GetIt.instance;

Future<void> init() async {
  // Core
  sl.registerLazySingleton<ApiClient>(() => ApiClient());
  sl.registerLazySingleton<ApiService>(() => ApiService());
  sl.registerLazySingleton<SecureStorageService>(() => SecureStorageService());
  
  // BLoCs
  sl.registerFactory<AuthBloc>(
    () => AuthBloc(
      apiService: sl<ApiService>(),
      secureStorage: sl<SecureStorageService>(),
    ),
  );
  
  sl.registerFactory<WalletBloc>(
    () => WalletBloc(
      apiService: sl<ApiService>(),
    ),
  );
  
  sl.registerFactory<ExchangeBloc>(
    () => ExchangeBloc(
      apiService: sl<ApiService>(),
    ),
  );
  
  sl.registerFactory<CardsBloc>(
    () => CardsBloc(
      apiService: sl<ApiService>(),
    ),
  );
  
  sl.registerFactory<PortfolioBloc>(
    () => PortfolioBloc(
      apiService: sl<ApiService>(),
    ),
  );
}
