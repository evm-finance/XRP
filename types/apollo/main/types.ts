export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  Any: any;
  JSONMap: any;
  Map: any;
};

export type AaveAddress = {
  __typename?: 'AaveAddress';
  aTokenAddress: Scalars['String'];
  aTokenSymbol: Scalars['String'];
  address: Scalars['String'];
  decimals: Scalars['Int'];
  stableDebtTokenAddress: Scalars['String'];
  symbol: Scalars['String'];
  variableDebtTokenAddress: Scalars['String'];
};

export type AavePool = {
  __typename?: 'AavePool';
  aEmissionPerSecond: Scalars['Float'];
  addresses: AaveAddress;
  availableLiquidity: Scalars['Float'];
  baseLTVasCollateral: Scalars['Float'];
  borrowingEnabled: Scalars['Boolean'];
  decimals: Scalars['Int'];
  id: Scalars['String'];
  liquidityRate: Scalars['Float'];
  name: Scalars['String'];
  portfolioVal: AavePortfolio;
  price: AavePoolPrice;
  reserveLiquidationBonus: Scalars['Float'];
  reserveLiquidationThreshold: Scalars['Float'];
  sEmissionPerSecond: Scalars['Float'];
  stableBorrowRate: Scalars['Float'];
  stableBorrowRateEnabled: Scalars['Boolean'];
  symbol: Scalars['String'];
  totalATokenSupply: Scalars['Float'];
  totalCurrentVariableDebt: Scalars['Float'];
  totalLiquidity: Scalars['Float'];
  totalLiquidityAsCollateral: Scalars['Float'];
  totalPrincipalStableDebt: Scalars['Float'];
  underlyingAsset: Scalars['String'];
  usageAsCollateralEnabled: Scalars['Boolean'];
  utilizationRate: Scalars['Float'];
  vEmissionPerSecond: Scalars['Float'];
  variableBorrowRate: Scalars['Float'];
};

export type AavePoolPrice = {
  __typename?: 'AavePoolPrice';
  id: Scalars['String'];
  priceInEth: Scalars['Float'];
  priceUsd: Scalars['Float'];
};

export type AavePortfolio = {
  __typename?: 'AavePortfolio';
  stableBorrow: Scalars['Float'];
  symbol: Scalars['String'];
  totalDeposits: Scalars['Float'];
  variableBorrow: Scalars['Float'];
  walletBal: Scalars['Float'];
  nativeBalance: Scalars['Float'];
  isWrapped: Scalars['Boolean'];
  networkSymbol: Scalars['String'];
};

export type Balance = {
  __typename?: 'Balance';
  aavePools: Array<AavePool>;
  address: Scalars['String'];
  chainId: Scalars['Float'];
  items: Array<BalanceItem>;
  nextUpdateAt: Scalars['String'];
  pagination?: Maybe<Pagination>;
  quoteCurrency: Scalars['String'];
  updatedAt: Scalars['String'];
};

export type BalanceItem = {
  __typename?: 'BalanceItem';
  balance: Scalars['Float'];
  balance24h: Scalars['Float'];
  contractAddress: Scalars['String'];
  contractDecimals: Scalars['Float'];
  contractName: Scalars['String'];
  contractTickerSymbol: Scalars['String'];
  disableQuoteRate: Scalars['Boolean'];
  lastTransferredAt: Scalars['String'];
  nativeToken: Scalars['Boolean'];
  nftData: Scalars['Float'];
  quote: Scalars['Float'];
  quote24h: Scalars['Float'];
  quoteRate: Scalars['Float'];
  quoteRate24h: Scalars['Float'];
  supportsErc: Array<Scalars['String']>;
  type: Scalars['String'];
};

export type Block = {
  __typename?: 'Block';
  XRPLedger: XrpLedger;
  XRPTransactions: XrpTransactions;
  blockNumber: Scalars['Int'];
  burnCount: Scalars['Int'];
  events: TransactionEvents;
  metrics: BlockMetrics;
  minedAt: Scalars['Int'];
  mintCount: Scalars['Int'];
  network: Scalars['String'];
  pairCreatedCount: Scalars['Int'];
  swapCount: Scalars['Int'];
  txCount: Scalars['Int'];
};

export type BlockMetric = {
  __typename?: 'BlockMetric';
  address: Scalars['String'];
  change1H: Scalars['Float'];
  contract: Scalars['String'];
  token0Symbol: Scalars['String'];
  token1Symbol: Scalars['String'];
  totalLiquidity: Scalars['Float'];
};

export type BlockMetrics = {
  __typename?: 'BlockMetrics';
  items: Array<BlockMetric>;
};

export type BlockTransactionContractTransfer = {
  __typename?: 'BlockTransactionContractTransfer';
  blockSignedAt: Scalars['String'];
  gasPrice: Scalars['Int'];
  gasQuote: Scalars['Float'];
  gasSpent: Scalars['Int'];
  successful: Scalars['Boolean'];
  transfers: Array<TokenTransferItem>;
  txHash: Scalars['String'];
};

// export type Chain = {
//   __typename?: 'Chain';
//   blockExplorerUrl: Scalars['String'];
//   chainIdentifier: Scalars['Int'];
//   dex: Array<Dex>;
//   id: Scalars['String'];
//   name: Scalars['String'];
//   rpcUrl: Scalars['String'];
//   symbol: Scalars['String'];
//   weth: UniswapToken;
// };

export type DailyChart = {
  __typename?: 'DailyChart';
  date: Scalars['Int'];
  priceUsd: Scalars['Float'];
};



export type EthGasStatsResult = {
  __typename?: 'EthGasStatsResult';
  fastGasPrice: Scalars['Float'];
  gasUsedRatio: Scalars['String'];
  lastBlock: Scalars['Int'];
  proposeGasPrice: Scalars['Float'];
  safeGasPrice: Scalars['Float'];
  suggestBaseFee: Scalars['Float'];
};

export type GasStats = {
  __typename?: 'GasStats';
  blockExplorer: Scalars['String'];
  gas: EthGasStatsResult;
  name: Scalars['String'];
  symbol: Scalars['String'];
};

export type GlobalStats = {
  __typename?: 'GlobalStats';
  defiDominance: Scalars['Float'];
  defiMarketCap: Scalars['Float'];
  defiToEthRatio: Scalars['Float'];
  ethMarketCap: Scalars['Float'];
  topCoinDefiDominance: Scalars['Float'];
  topCoinName: Scalars['String'];
  tradingVolume24h: Scalars['Float'];
};

export type HighAndLow = {
  __typename?: 'HighAndLow';
  high: Scalars['Float'];
  interval: Scalars['String'];
  low: Scalars['Float'];
  unixTime: Scalars['Int'];
};

export type LogEvent = {
  __typename?: 'LogEvent';
  decoded: LogEventDecoded;
  logOffset: Scalars['Int'];
  senderAddress: Scalars['String'];
  senderContractDecimals: Scalars['Int'];
  senderContractTickerSymbol: Scalars['String'];
  senderLogoUrl: Scalars['String'];
  senderName: Scalars['String'];
  txHash: Scalars['String'];
};

export type LogEventDecoded = {
  __typename?: 'LogEventDecoded';
  name: Scalars['String'];
  params: Array<Maybe<LogEventParams>>;
  signature: Scalars['String'];
};

export type LogEventParams = {
  __typename?: 'LogEventParams';
  decoded: Scalars['Boolean'];
  name: Scalars['String'];
  type: Scalars['String'];
  value: Scalars['String'];
};

export type News = {
  __typename?: 'News';
  currencies: Array<NewsCurrency>;
  id: Scalars['Int'];
  publishedAt: Scalars['String'];
  title: Scalars['String'];
  url: Scalars['String'];
};

export type NewsCurrency = {
  __typename?: 'NewsCurrency';
  code: Scalars['String'];
  slug: Scalars['String'];
  title: Scalars['String'];
  url: Scalars['String'];
};

export type Pagination = {
  __typename?: 'Pagination';
  hasMore?: Maybe<Scalars['Boolean']>;
  pageNumber?: Maybe<Scalars['Int']>;
  pageSize?: Maybe<Scalars['Int']>;
  totalCount?: Maybe<Scalars['Int']>;
};

export type Pool = {
  __typename?: 'Pool';
  address: Scalars['String'];
  change1h: Scalars['Float'];
  change5Min: Scalars['Float'];
  change24h: Scalars['Float'];
  dex: Scalars['String'];
  id: Scalars['Int'];
  network: Scalars['String'];
  quoteExactIn: Scalars['Float'];
  reserveRatio: Scalars['Float'];
  token0Address: Scalars['String'];
  token0Decimals: Scalars['Int'];
  token0Name: Scalars['String'];
  token0PriceUSD: Scalars['Float'];
  token0Symbol: Scalars['String'];
  token1Address: Scalars['String'];
  token1Decimals: Scalars['Int'];
  token1Name: Scalars['String'];
  token1PriceUSD: Scalars['Float'];
  token1Symbol: Scalars['String'];
};

export type Price = {
  __typename?: 'Price';
  dataSource: Scalars['String'];
  priceEth: Scalars['Float'];
  priceUsd: Scalars['Float'];
  qcKey: Scalars['String'];
  symbolName: Scalars['String'];
};

export type PriceStream = {
  __typename?: 'PriceStream';
  dex: Scalars['String'];
  id: Scalars['Int'];
  network: Scalars['String'];
  pairAddress: Scalars['String'];
  quoteExactIn: Scalars['Float'];
  reserveRatio: Scalars['Float'];
  token0PriceUSD: Scalars['Float'];
  token1PriceUSD: Scalars['Float'];
};

export type Query = {
  __typename?: 'Query';
  aaveAddresses: Array<AaveAddress>;
  aavePools: Array<AavePool>;
  balances: Array<Balance>;
  block: Block;
  blocks: Array<Block>;
  dailyChart: Array<DailyChart>;
  fiatPrices: Scalars['Map'];
  gas: Array<GasStats>;
  globalStats: GlobalStats;
  poolScreener: Array<Pool>;
  recentPrices: Scalars['Map'];
  token: Token;
  transactions: Transaction;
  transfers: Transfers;
  uniswapTokens: UniswapTokens;
};


export type QueryAaveAddressesArgs = {
  chainId: Scalars['Int'];
};


export type QueryAavePoolsArgs = {
  chainId: Scalars['Int'];
};


export type QueryBalancesArgs = {
  address: Scalars['String'];
  chainIds: Array<Scalars['Int']>;
};


export type QueryBlockArgs = {
  blockNumber: Scalars['Int'];
  network: Scalars['String'];
};


export type QueryBlocksArgs = {
  network: Scalars['String'];
};


export type QueryDailyChartArgs = {
  coinGeckoID?: Scalars['String'];
};


export type QueryFiatPricesArgs = {
  addresses: Array<Scalars['String']>;
  platform?: Scalars['String'];
};


export type QueryPoolScreenerArgs = {
  dex: Scalars['String'];
  network: Scalars['String'];
  pageNumber?: Scalars['Int'];
  sort: Scalars['String'];
  sortBy: Scalars['String'];
};


export type QueryTokenArgs = {
  chainId?: Scalars['Int'];
  contractAddress?: Scalars['String'];
  decimals?: Scalars['Int'];
  interval?: TimeInterval;
  qcKey: Scalars['String'];
  walletAddress?: Scalars['String'];
};


export type QueryTransactionsArgs = {
  address: Scalars['String'];
  chainId?: Scalars['Int'];
  pageNumber?: Scalars['Int'];
  pageSize?: Scalars['Int'];
};


export type QueryTransfersArgs = {
  address: Scalars['String'];
  chainId?: Scalars['Int'];
  contractAddress: Scalars['String'];
  pageNumber?: Scalars['Int'];
  pageSize?: Scalars['Int'];
};


export type QueryUniswapTokensArgs = {
  chainId?: Scalars['Int'];
  pageNumber?: Scalars['Int'];
  pageSize?: Scalars['Int'];
  searchString?: Scalars['String'];
  source: Array<UniswapTokenSource>;
  userWallet?: Scalars['String'];
};

export type ScreenerItem = {
  __typename?: 'ScreenerItem';
  PriceUSD: Scalars['Float'];
  address: Scalars['String'];
  change1h: Scalars['Float'];
  change5Min: Scalars['Float'];
  change24h: Scalars['Float'];
  decimals: Scalars['Int'];
  name: Scalars['String'];
  platform: Scalars['String'];
  symbol: Scalars['String'];
  token0Address: Scalars['String'];
  token0Decimals: Scalars['Int'];
  token0Name: Scalars['String'];
  token0Symbol: Scalars['String'];
  token1Address: Scalars['String'];
  token1Decimals: Scalars['Int'];
  token1Name: Scalars['String'];
  token1Symbol: Scalars['String'];
};

export type Subscription = {
  __typename?: 'Subscription';
  block: Array<Block>;
  currentTime: Time;
  priceStream: Array<PriceStream>;
};


export type SubscriptionBlockArgs = {
  network: Scalars['String'];
};


export type SubscriptionPriceStreamArgs = {
  address?: InputMaybe<Array<Scalars['String']>>;
  dex: Scalars['String'];
  network: Scalars['String'];
};

export type Time = {
  __typename?: 'Time';
  timeStamp: Scalars['String'];
  unixTime: Scalars['Int'];
};

export enum TimeInterval {
  Interval_1H = 'INTERVAL_1H',
  Interval_1Week = 'INTERVAL_1WEEK',
  Interval_2H = 'INTERVAL_2H',
  Interval_4H = 'INTERVAL_4H',
  Interval_5Min = 'INTERVAL_5MIN',
  Interval_15Min = 'INTERVAL_15MIN',
  Interval_24H = 'INTERVAL_24H',
  Interval_30Min = 'INTERVAL_30MIN'
}

export type Token = {
  __typename?: 'Token';
  ID?: Maybe<Scalars['ID']>;
  aavePools: Array<AavePool>;
  aaveSymbol: Scalars['String'];
  balances: Array<Balance>;
  bitbucketRepos?: Maybe<Scalars['JSONMap']>;
  chainId: Scalars['Int'];
  circulatingSupply: Scalars['Float'];
  coinDescription: Scalars['String'];
  coinGeckoID: Scalars['String'];
  contractAddress: Scalars['String'];
  decimals: Scalars['Int'];
  discordChannelId: Scalars['String'];
  explorerUrls?: Maybe<Scalars['JSONMap']>;
  facebookUrl: Scalars['String'];
  githubRepos?: Maybe<Scalars['JSONMap']>;
  interval: TimeInterval;
  isAaveToken: Scalars['Boolean'];
  isQCToken: Scalars['Boolean'];
  marketcap: Scalars['Float'];
  news: Array<News>;
  price: Price;
  price24h: Scalars['Float'];
  qcKey: Scalars['String'];
  rank: Scalars['Int'];
  resistance1h: Scalars['Float'];
  safeScore: Scalars['Float'];
  subredditUrl: Scalars['String'];
  support1h: Scalars['Float'];
  symbolName: Scalars['String'];
  telegramChannelId: Scalars['String'];
  tokenInterval: HighAndLow;
  twitterUrl: Scalars['String'];
  volume24h: Scalars['Float'];
  walletAddress: Scalars['String'];
  websiteUrl?: Maybe<Scalars['JSONMap']>;
};

export type TokenTransferItem = {
  __typename?: 'TokenTransferItem';
  contractAddress: Scalars['String'];
  contractDecimals: Scalars['Int'];
  contractName: Scalars['String'];
  contractTickerSymbol: Scalars['String'];
  delta: Scalars['String'];
  deltaQuote: Scalars['Float'];
  fromAddress: Scalars['String'];
  fromAddressIsContract: Scalars['Boolean'];
  fromAddressLabel: Scalars['String'];
  fromAddressName: Scalars['String'];
  fromAddressSymbol: Scalars['String'];
  logoUrl: Scalars['String'];
  toAddress: Scalars['String'];
  toAddressIsContract: Scalars['Boolean'];
  toAddressLabel: Scalars['String'];
  toAddressName: Scalars['String'];
  toAddressSymbol: Scalars['String'];
  transferType: Scalars['String'];
};

export type Transaction = {
  __typename?: 'Transaction';
  address: Scalars['String'];
  chainID: Scalars['Int'];
  items: Array<TransactionItem>;
  nextUpdateAt: Scalars['String'];
  pagination?: Maybe<Pagination>;
  quoteCurrency: Scalars['String'];
  updatedAt: Scalars['String'];
};

export type TransactionEvent = {
  __typename?: 'TransactionEvent';
  address: Scalars['String'];
  allFunctionParams?: Maybe<Scalars['JSONMap']>;
  contract: Scalars['String'];
  indexedParams?: Maybe<Scalars['JSONMap']>;
  name: Scalars['String'];
  network: Scalars['String'];
  outputDataMap?: Maybe<Scalars['JSONMap']>;
  outputDataMapHex: Scalars['String'];
  signature: Scalars['String'];
  topic: Scalars['String'];
};

export type TransactionEvents = {
  __typename?: 'TransactionEvents';
  items: Array<TransactionEvent>;
};

export type TransactionItem = {
  __typename?: 'TransactionItem';
  blockSignedAt: Scalars['String'];
  callFunction: Scalars['String'];
  decodedFunctionID: Scalars['String'];
  fromAddress: Scalars['String'];
  fromAddressIsContract: Scalars['Boolean'];
  fromAddressName: Scalars['String'];
  fromAddressSymbol: Scalars['String'];
  gasPrice: Scalars['Float'];
  gasQuote: Scalars['Float'];
  gasSpent: Scalars['Float'];
  logEvents?: Maybe<Array<LogEvent>>;
  successful: Scalars['Boolean'];
  toAddress: Scalars['String'];
  toAddressIsContract: Scalars['Boolean'];
  toAddressName: Scalars['String'];
  toAddressSymbol: Scalars['String'];
  txDetails?: Maybe<Array<TxDetail>>;
  txHash: Scalars['String'];
  value: Scalars['String'];
  valueQuote: Scalars['Float'];
};

export type Transfers = {
  __typename?: 'Transfers';
  address: Scalars['String'];
  chainId: Scalars['Int'];
  items: Array<BlockTransactionContractTransfer>;
  nextUpdateAt: Scalars['String'];
  pagination: Pagination;
  quoteCurrency: Scalars['String'];
  updatedAt: Scalars['String'];
};

export type TxDetail = {
  __typename?: 'TxDetail';
  fromAddress: Scalars['String'];
  method: Scalars['String'];
  priceUSD: Scalars['Float'];
  toAddress: Scalars['String'];
  tokenAddress: Scalars['String'];
  tokenContractDecimals: Scalars['Int'];
  tokenContractSymbol: Scalars['String'];
  tokenLogoUrl: Scalars['String'];
  tokenSymbolName: Scalars['String'];
  value: Scalars['String'];
};

export type UniswapToken = {
  __typename?: 'UniswapToken';
  address: Scalars['String'];
  chainId: Scalars['Int'];
  decimals: Scalars['Int'];
  name: Scalars['String'];
  symbol: Scalars['String'];
};

export enum UniswapTokenSource {
  Aave = 'aave',
  CoinGecko = 'coinGecko',
  Uniswap = 'uniswap',
  UserBalances = 'userBalances'
}

export type UniswapTokens = {
  __typename?: 'UniswapTokens';
  items: Array<UniswapToken>;
  pagination?: Maybe<Pagination>;
};

export type XrpLedger = {
  __typename?: 'XRPLedger';
  eventsCount?: Maybe<Scalars['Map']>;
  ledger: XrpLedgerData;
  ledgerHash: Scalars['String'];
  ledgerIndex: Scalars['Int'];
  validated: Scalars['Boolean'];
};

export type XrpLedgerData = {
  __typename?: 'XRPLedgerData';
  accepted: Scalars['Boolean'];
  accountHash: Scalars['String'];
  closeFlags: Scalars['Int'];
  closeTime: Scalars['Int'];
  closeTimeHuman: Scalars['String'];
  closeTimeResolution: Scalars['Int'];
  closed: Scalars['Boolean'];
  hash: Scalars['String'];
  ledgerHash: Scalars['String'];
  ledgerIndex: Scalars['String'];
  parentCloseTime: Scalars['Int'];
  parentHash: Scalars['String'];
  seqNum: Scalars['String'];
  totalCoins: Scalars['String'];
  totalCoins1: Scalars['String'];
  transactionHash: Scalars['String'];
};

export type XrpTransaction = {
  __typename?: 'XRPTransaction';
  account: Scalars['String'];
  amount?: Maybe<Scalars['Any']>;
  date: Scalars['Int'];
  destination: Scalars['String'];
  fee: Scalars['String'];
  flags: Scalars['Int'];
  hash: Scalars['String'];
  inLedger: Scalars['Int'];
  lastLedgerSequence: Scalars['Int'];
  ledgerIndex: Scalars['Int'];
  memos: Array<Maybe<Scalars['Any']>>;
  meta?: Maybe<Scalars['Any']>;
  metadata?: Maybe<Scalars['Any']>;
  offerSequence: Scalars['Int'];
  sequence: Scalars['Int'];
  signingPubKey: Scalars['String'];
  takerGets?: Maybe<Scalars['Any']>;
  takerPays?: Maybe<Scalars['Any']>;
  transactionType: Scalars['String'];
  txnSignature: Scalars['String'];
  validated?: Maybe<Scalars['Boolean']>;
  warnings: Array<Maybe<Scalars['Any']>>;
};

export type XrpTransactions = {
  __typename?: 'XRPTransactions';
  items: Array<XrpTransaction>;
};

export type SupportedChainsGqlQueryVariables = Exact<{ [key: string]: never; }>;


export type SupportedChainsGqlQuery = { __typename?: 'Query', chains: Array<{ __typename?: 'Chain', id: string, chainIdentifier: number, name: string, symbol: string, rpcUrl: string, blockExplorerUrl: string, dex: Array<{ __typename?: 'Dex', name: string, value: string, symbol: string }>, weth: { __typename?: 'UniswapToken', chainId: number, address: string, symbol: string } }> };

export type GasGqlQueryVariables = Exact<{ [key: string]: never; }>;


export type GasGqlQuery = { __typename?: 'Query', gas: Array<{ __typename?: 'GasStats', symbol: string, blockExplorer: string, name: string, gas: { __typename?: 'EthGasStatsResult', lastBlock: number, safeGasPrice: number, proposeGasPrice: number, fastGasPrice: number, suggestBaseFee: number } }> };

export type DeFiStatsQueryVariables = Exact<{ [key: string]: never; }>;


export type DeFiStatsQuery = { __typename?: 'Query', globalStats: { __typename?: 'GlobalStats', defiMarketCap: number, ethMarketCap: number, defiToEthRatio: number, tradingVolume24h: number, defiDominance: number, topCoinName: string, topCoinDefiDominance: number } };

export type UniswapTokensGqlQueryVariables = Exact<{
  chainId: Scalars['Int'];
  userWallet: Scalars['String'];
  searchString: Scalars['String'];
  source: Array<UniswapTokenSource> | UniswapTokenSource;
  pageNumber: Scalars['Int'];
  pageSize: Scalars['Int'];
}>;


export type UniswapTokensGqlQuery = { __typename?: 'Query', uniswapTokens: { __typename?: 'UniswapTokens', items: Array<{ __typename?: 'UniswapToken', chainId: number, address: string, symbol: string, name: string, decimals: number }>, pagination?: { __typename?: 'Pagination', hasMore?: boolean | null, pageNumber?: number | null, pageSize?: number | null, totalCount?: number | null } | null } };

export type RecentPricesGqlQueryVariables = Exact<{ [key: string]: never; }>;


export type RecentPricesGqlQuery = { __typename?: 'Query', recentPrices: any };

export type FiatPricesGqlQueryVariables = Exact<{
  addresses: Array<Scalars['String']> | Scalars['String'];
  platform: Scalars['String'];
}>;


export type FiatPricesGqlQuery = { __typename?: 'Query', fiatPrices: any };

export type AavePoolGqlQueryVariables = Exact<{
  chainId: Scalars['Int'];
}>;


export type AavePoolGqlQuery = { __typename?: 'Query', aavePools: Array<{ __typename?: 'AavePool', id: string, underlyingAsset: string, name: string, symbol: string, decimals: number, totalLiquidity: number, liquidityRate: number, stableBorrowRate: number, variableBorrowRate: number, aEmissionPerSecond: number, vEmissionPerSecond: number, sEmissionPerSecond: number, availableLiquidity: number, utilizationRate: number, totalATokenSupply: number, totalCurrentVariableDebt: number, totalPrincipalStableDebt: number, totalLiquidityAsCollateral: number, baseLTVasCollateral: number, reserveLiquidationThreshold: number, reserveLiquidationBonus: number, usageAsCollateralEnabled: boolean, borrowingEnabled: boolean, stableBorrowRateEnabled: boolean, price: { __typename?: 'AavePoolPrice', id: string, priceInEth: number, priceUsd: number }, addresses: { __typename?: 'AaveAddress', aTokenAddress: string, aTokenSymbol: string, stableDebtTokenAddress: string, variableDebtTokenAddress: string, decimals: number, address: string, symbol: string }, portfolioVal: { __typename?: 'AavePortfolio', totalDeposits: number, walletBal: number, stableBorrow: number, variableBorrow: number } }> };

export type BalancesGqlQueryVariables = Exact<{
  chainIds: Array<Scalars['Int']> | Scalars['Int'];
  address: Scalars['String'];
}>;


export type BalancesGqlQuery = { __typename?: 'Query', balances: Array<{ __typename?: 'Balance', address: string, updatedAt: string, nextUpdateAt: string, quoteCurrency: string, chainId: number, pagination?: { __typename?: 'Pagination', hasMore?: boolean | null } | null, items: Array<{ __typename?: 'BalanceItem', contractDecimals: number, contractName: string, contractTickerSymbol: string, contractAddress: string, supportsErc: Array<string>, lastTransferredAt: string, type: string, balance: number, balance24h: number, quoteRate: number, quoteRate24h: number, quote: number, disableQuoteRate: boolean, quote24h: number, nftData: number }>, aavePools: Array<{ __typename?: 'AavePool', symbol: string, id: string, name: string, liquidityRate: number, borrowingEnabled: boolean, stableBorrowRateEnabled: boolean, stableBorrowRate: number, variableBorrowRate: number, portfolioVal: { __typename?: 'AavePortfolio', totalDeposits: number, walletBal: number, stableBorrow: number, variableBorrow: number }, price: { __typename?: 'AavePoolPrice', id: string, priceInEth: number, priceUsd: number }, addresses: { __typename?: 'AaveAddress', address: string, decimals: number } }> }> };

export type TransactionsGqlQueryVariables = Exact<{
  chainId: Scalars['Int'];
  address: Scalars['String'];
  pageNumber: Scalars['Int'];
  pageSize: Scalars['Int'];
}>;


export type TransactionsGqlQuery = { __typename?: 'Query', transactions: { __typename?: 'Transaction', pagination?: { __typename?: 'Pagination', hasMore?: boolean | null, pageSize?: number | null, pageNumber?: number | null } | null, items: Array<{ __typename?: 'TransactionItem', blockSignedAt: string, txHash: string, successful: boolean, fromAddress: string, fromAddressName: string, fromAddressSymbol: string, fromAddressIsContract: boolean, toAddress: string, toAddressName: string, toAddressIsContract: boolean, toAddressSymbol: string, value: string, valueQuote: number, gasQuote: number, gasPrice: number, gasSpent: number, callFunction: string, txDetails?: Array<{ __typename?: 'TxDetail', fromAddress: string, toAddress: string, value: string, tokenSymbolName: string, tokenContractSymbol: string, tokenContractDecimals: number, tokenAddress: string, tokenLogoUrl: string, priceUSD: number }> | null, logEvents?: Array<{ __typename?: 'LogEvent', logOffset: number, txHash: string, senderContractDecimals: number, senderName: string, senderContractTickerSymbol: string, senderAddress: string, senderLogoUrl: string, decoded: { __typename?: 'LogEventDecoded', name: string, signature: string, params: Array<{ __typename?: 'LogEventParams', name: string, type: string, decoded: boolean, value: string } | null> } }> | null }> } };

export type TransfersGqlQueryVariables = Exact<{
  chainId: Scalars['Int'];
  address: Scalars['String'];
  contractAddress: Scalars['String'];
  pageNumber: Scalars['Int'];
  pageSize: Scalars['Int'];
}>;


export type TransfersGqlQuery = { __typename?: 'Query', transfers: { __typename?: 'Transfers', pagination: { __typename?: 'Pagination', hasMore?: boolean | null, pageSize?: number | null, pageNumber?: number | null }, items: Array<{ __typename?: 'BlockTransactionContractTransfer', blockSignedAt: string, txHash: string, successful: boolean, gasSpent: number, gasPrice: number, gasQuote: number, transfers: Array<{ __typename?: 'TokenTransferItem', fromAddress: string, fromAddressLabel: string, fromAddressName: string, fromAddressSymbol: string, fromAddressIsContract: boolean, toAddress: string, toAddressLabel: string, toAddressName: string, toAddressSymbol: string, toAddressIsContract: boolean, contractDecimals: number, contractName: string, contractTickerSymbol: string, contractAddress: string, logoUrl: string, transferType: string, delta: string, deltaQuote: number }> }> } };

export type TokenQueryGqlQueryVariables = Exact<{
  qcKey: Scalars['String'];
  walletAddress: Scalars['String'];
  interval: TimeInterval;
  contractAddress: Scalars['String'];
  chainId: Scalars['Int'];
  decimals: Scalars['Int'];
}>;


export type TokenQueryGqlQuery = { __typename?: 'Query', token: { __typename?: 'Token', ID?: string | null, walletAddress: string, interval: TimeInterval, qcKey: string, aaveSymbol: string, isAaveToken: boolean, isQCToken: boolean, coinGeckoID: string, symbolName: string, rank: number, chainId: number, contractAddress: string, decimals: number, price24h: number, marketcap: number, volume24h: number, circulatingSupply: number, support1h: number, resistance1h: number, safeScore: number, websiteUrl?: any | null, bitbucketRepos?: any | null, githubRepos?: any | null, explorerUrls?: any | null, telegramChannelId: string, twitterUrl: string, subredditUrl: string, facebookUrl: string, coinDescription: string, tokenInterval: { __typename?: 'HighAndLow', high: number, low: number, interval: string, unixTime: number }, price: { __typename?: 'Price', qcKey: string, symbolName: string, priceUsd: number, priceEth: number }, news: Array<{ __typename?: 'News', id: number, title: string, url: string, publishedAt: string, currencies: Array<{ __typename?: 'NewsCurrency', code: string, title: string, slug: string, url: string }> }>, balances: Array<{ __typename?: 'Balance', chainId: number, items: Array<{ __typename?: 'BalanceItem', balance: number }> }> } };

export type DailyChartGqlQueryVariables = Exact<{
  coinGeckoID: Scalars['String'];
}>;


export type DailyChartGqlQuery = { __typename?: 'Query', dailyChart: Array<{ __typename?: 'DailyChart', date: number, priceUsd: number }> };

export type ScreenerGqlQueryVariables = Exact<{
  network: Scalars['String'];
  dex: Scalars['String'];
  sortBy: Scalars['String'];
  sort: Scalars['String'];
  pageNumber: Scalars['Int'];
}>;


export type ScreenerGqlQuery = { __typename?: 'Query', poolScreener: Array<{ __typename?: 'Pool', network: string, dex: string, address: string, token0Symbol: string, token0Address: string, token0Name: string, token0Decimals: number, token1Symbol: string, token1Name: string, token0PriceUSD: number, token1PriceUSD: number, quoteExactIn: number, reserveRatio: number, change5Min: number, change1h: number, change24h: number }> };

export type TimeGqlSubscriptionVariables = Exact<{ [key: string]: never; }>;


export type TimeGqlSubscription = { __typename?: 'Subscription', currentTime: { __typename?: 'Time', timeStamp: string, unixTime: number } };

export type PriceStreamGqlSubscriptionVariables = Exact<{
  network: Scalars['String'];
  dex: Scalars['String'];
  address: Array<Scalars['String']> | Scalars['String'];
}>;


export type PriceStreamGqlSubscription = { __typename?: 'Subscription', priceStream: Array<{ __typename?: 'PriceStream', network: string, dex: string, pairAddress: string, token0PriceUSD: number, token1PriceUSD: number, quoteExactIn: number, reserveRatio: number }> };

export type BlocksGqlQueryVariables = Exact<{
  network: Scalars['String'];
}>;


export type BlocksGqlQuery = { __typename?: 'Query', blocks: Array<{ __typename?: 'Block', network: string, blockNumber: number, minedAt: number, txCount: number, swapCount: number, pairCreatedCount: number, mintCount: number, metrics: { __typename?: 'BlockMetrics', items: Array<{ __typename?: 'BlockMetric', totalLiquidity: number, change1H: number, token0Symbol: string, token1Symbol: string }> } }> };

export type BlocksXrpGqlQueryVariables = Exact<{
  network: Scalars['String'];
}>;


export type BlocksXrpGqlQuery = { __typename?: 'Query', blocks: Array<{ __typename?: 'Block', network: string, blockNumber: number, minedAt: number, txCount: number, XRPLedger: { __typename?: 'XRPLedger', ledgerHash: string, eventsCount?: any | null } }> };

export type BlockGqlQueryVariables = Exact<{
  network: Scalars['String'];
  blockNumber: Scalars['Int'];
}>;


export type BlockGqlQuery = { __typename?: 'Query', block: { __typename?: 'Block', network: string, blockNumber: number, minedAt: number, txCount: number, XRPLedger: { __typename?: 'XRPLedger', ledgerHash: string, eventsCount?: any | null, ledger: { __typename?: 'XRPLedgerData', ledgerHash: string, parentHash: string, transactionHash: string, closeTimeHuman: string, totalCoins: string } }, XRPTransactions: { __typename?: 'XRPTransactions', items: Array<{ __typename?: 'XRPTransaction', hash: string, account: string, amount?: any | null, fee: string, metadata?: any | null }> } } };

export type BlocksStreamGqlSubscriptionVariables = Exact<{
  network: Scalars['String'];
}>;


export type BlocksStreamGqlSubscription = { __typename?: 'Subscription', block: Array<{ __typename?: 'Block', network: string, blockNumber: number, minedAt: number, txCount: number, swapCount: number, pairCreatedCount: number, mintCount: number, metrics: { __typename?: 'BlockMetrics', items: Array<{ __typename?: 'BlockMetric', totalLiquidity: number, change1H: number, token0Symbol: string, token1Symbol: string }> }, XRPLedger: { __typename?: 'XRPLedger', ledgerHash: string, eventsCount?: any | null } }> };
