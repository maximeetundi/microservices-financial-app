import 'package:get_it/get_it.dart';

import '../api/api_client.dart';
import '../services/api_service.dart';
import '../services/secure_storage_service.dart';
import '../../features/auth/presentation/bloc/auth_bloc.dart';
import '../../features/wallet/presentation/bloc/wallet_bloc.dart';
import '../../features/exchange/presentation/bloc/exchange_bloc.dart';
import '../../features/cards/presentation/bloc/cards_bloc.dart';
import '../../features/portfolio/presentation/bloc/portfolio_bloc.dart';
import '../../features/transfer/presentation/bloc/transfer_bloc.dart';
import '../../features/transfer/domain/usecases/send_transfer_usecase.dart';
import '../../features/transfer/domain/usecases/get_transfer_history_usecase.dart';

final sl = GetIt.instance;

Future<void> init() async {
  // Core
  sl.registerLazySingleton<ApiClient>(() => ApiClient());
  sl.registerLazySingleton<ApiService>(() => ApiService());
  sl.registerLazySingleton<SecureStorageService>(() => SecureStorageService());
  
  // Use Cases (no-arg constructors)
  sl.registerLazySingleton<SendTransferUseCase>(() => SendTransferUseCase());
  sl.registerLazySingleton<GetTransferHistoryUseCase>(() => GetTransferHistoryUseCase());
  
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
      secureStorage: sl<SecureStorageService>(), // Added dependency
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
  
  sl.registerFactory<TransferBloc>(
    () => TransferBloc(
      sendTransferUseCase: sl<SendTransferUseCase>(),
      getTransferHistoryUseCase: sl<GetTransferHistoryUseCase>(),
    ),
  );
}

