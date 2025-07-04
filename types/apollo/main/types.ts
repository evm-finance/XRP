export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };

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

// XRP Types
export type Block = {
  __typename?: 'Block';
  XRPLedger: XrpLedger;
  XRPTransactions: XrpTransactions;
  blockNumber: Scalars['Int'];
  minedAt: Scalars['Int'];
  network: Scalars['String'];
  txCount: Scalars['Int'];
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

export type XRPBalanceElem = {
  symbol: Scalars['String'];
  issuer: Scalars['String'];
  name: Scalars['String'];
  balance: Scalars['Float'];
  price: Scalars['Float'];
  value: Scalars['Float'];
};

export type XRPAccountBalances = {
  account: Scalars['String'];
  xrpBalance: Scalars['Float'];
  xrpPrice: Scalars['Float'];
  xrpTokens: Array<XRPBalanceElem>;
};

export type XRPDefiData = {
  account: Scalars['String'];
  xrpBalance: Scalars['Float'];
  xrpPrice: Scalars['Float'];
  xrpTransactions: Array<XrpTransaction>;
  xrpTokenBalances: Array<XRPBalanceElem>;
};

export type XRPAsset = {
  __typename?: 'XRPAsset';
  currency: Scalars['String'];
  issuer?: Maybe<Scalars['String']>;
};

export type XRPAssetInput = {
  currency: Scalars['String'];
  issuer?: Maybe<Scalars['String']>;
};

export type XRPAMMPool = {
  __typename?: 'XRPAMMPool';
  poolId: Scalars['String'];
  asset1: XRPAsset;
  asset2: XRPAsset;
  asset1Balance: Scalars['Float'];
  asset2Balance: Scalars['Float'];
  lpBalance: Scalars['Float'];
  fee: Scalars['Float'];
  tradingVolume24H?: Maybe<Scalars['Float']>;
  tradingVolume7D?: Maybe<Scalars['Float']>;
  createdAt: Scalars['Int'];
  lastUpdated: Scalars['Int'];
};

export type XRPAMMLiquidityCalculation = {
  __typename?: 'XRPAMMLiquidityCalculation';
  poolId: Scalars['String'];
  asset1ValueUsd: Scalars['Float'];
  asset2ValueUsd: Scalars['Float'];
  totalLiquidityUsd: Scalars['Float'];
  asset1Percentage: Scalars['Float'];
  asset2Percentage: Scalars['Float'];
  priceImpact: Scalars['Float'];
};

export type XRPAMMSwapQuote = {
  __typename?: 'XRPAMMSwapQuote';
  inputAmount: Scalars['Float'];
  outputAmount: Scalars['Float'];
  priceImpact: Scalars['Float'];
  fee: Scalars['Float'];
  minimumReceived: Scalars['Float'];
};

export type XRPToken = {
  __typename?: 'XRPToken';
  issuer: Scalars['String'];
  currency: Scalars['String'];
  name: Scalars['String'];
  icon?: Maybe<Scalars['String']>;
  description?: Maybe<Scalars['String']>;
  marketcap?: Maybe<Scalars['Float']>;
  price?: Maybe<Scalars['Float']>;
  volume24h?: Maybe<Scalars['Float']>;
  volume7d?: Maybe<Scalars['Float']>;
};

export type XRPScreenerItem = {
  __typename?: 'XRPScreenerItem';
  currency: Scalars['String'];
  issuerAddress: Scalars['String'];
  icon: Scalars['String'];
  tokenName: Scalars['String'];
  issuerName: Scalars['String'];
  marketcap: Scalars['Float'];
  price: Scalars['Float'];
  volume24H: Scalars['Float'];
};

export type XRPHeatmapItem = {
  __typename?: 'XRPHeatmapItem';
  currency: Scalars['String'];
  name: Scalars['String'];
  price_usd: Scalars['Float'];
  price_xrp: Scalars['Float'];
  price_change_1h: Scalars['Float'];
  price_change_4h: Scalars['Float'];
  price_change_24h: Scalars['Float'];
  price_change_7d: Scalars['Float'];
  price_change_30d: Scalars['Float'];
  marketcap: Scalars['Float'];
  volume_24h: Scalars['Float'];
  liquidity: Scalars['Float'];
  tvl: Scalars['Float'];
  issuer: Scalars['String'];
  hasAmmPool: Scalars['Boolean'];
  hasTrustLine: Scalars['Boolean'];
};

export type XRPAMMHeatmapItem = {
  __typename?: 'XRPAMMHeatmapItem';
  poolId: Scalars['String'];
  token1: XRPAsset;
  token2: XRPAsset;
  liquidity: Scalars['Float'];
  volume24h: Scalars['Float'];
  volume7d: Scalars['Float'];
  fee: Scalars['Float'];
  apr: Scalars['Float'];
  priceChange24h: Scalars['Float'];
  priceChange7d: Scalars['Float'];
  tvl: Scalars['Float'];
  transactions24h: Scalars['Int'];
  uniqueTraders24h: Scalars['Int'];
};

// Query Types
export type Query = {
  __typename?: 'Query';
  blocks: Array<Block>;
  block: Block;
  xrpTransaction: XrpTransaction;
  xrpAccountBalances: XRPAccountBalances;
  xrpAccountTransactions: Array<XrpTransaction>;
  xrpDefiData: XRPDefiData;
  xrpScreener: Array<XRPScreenerItem>;
  xrpAmmPools: Array<XRPAMMPool>;
  xrpAmmPoolDetails: XRPAMMPool;
  xrpAmmUserPositions: Array<XRPAMMPool>;
  xrpAmmQuote: XRPAMMSwapQuote;
  xrpAmmTransactions: Array<XrpTransaction>;
  xrpTokenPrice: XRPToken;
  xrpTokenBalances: Array<XRPBalanceElem>;
  xrpHeatmap: Array<XRPHeatmapItem>;
  xrpAmmHeatmap: Array<XRPAMMHeatmapItem>;
  xrpHeatmapUpdates: Array<XRPHeatmapItem>;
  xrpAMMLiquidityValues: Array<XRPAMMLiquidityCalculation>;
  xrpTokenMints: Array<XRPToken>;
  xrpLiquidityPools: Array<XRPAMMPool>;
};

// Query Arguments
export type QueryBlocksArgs = {
  network: Scalars['String'];
};

export type QueryBlockArgs = {
  blockNumber: Scalars['Int'];
  network: Scalars['String'];
};

export type QueryXrpTransactionArgs = {
  hash: Scalars['String'];
};

export type QueryXrpAccountBalancesArgs = {
  account: Scalars['String'];
};

export type QueryXrpAccountTransactionsArgs = {
  address: Scalars['String'];
};

export type QueryXrpDefiDataArgs = {
  address: Scalars['String'];
};

export type QueryXrpAmmPoolDetailsArgs = {
  poolId: Scalars['String'];
};

export type QueryXrpAmmUserPositionsArgs = {
  address: Scalars['String'];
};

export type QueryXrpAmmQuoteArgs = {
  poolId: Scalars['String'];
  amount: Scalars['String'];
  fromToken: Scalars['String'];
};

export type QueryXrpAmmTransactionsArgs = {
  poolId: Scalars['String'];
  limit?: Maybe<Scalars['Int']>;
};

export type QueryXrpTokenPriceArgs = {
  currency: Scalars['String'];
  issuer?: Maybe<Scalars['String']>;
};

export type QueryXrpTokenBalancesArgs = {
  address: Scalars['String'];
};

export type QueryXrpHeatmapArgs = {
  timeFrame?: Maybe<Scalars['String']>;
  blockSize?: Maybe<Scalars['String']>;
  limit?: Maybe<Scalars['Int']>;
};

export type QueryXrpAmmHeatmapArgs = {
  timeFrame?: Maybe<Scalars['String']>;
  blockSize?: Maybe<Scalars['String']>;
  limit?: Maybe<Scalars['Int']>;
};

export type QueryXrpHeatmapUpdatesArgs = {
  currencies?: Maybe<Array<Scalars['String']>>;
  timeFrame?: Maybe<Scalars['String']>;
};

export type QueryXrpTokenMintsArgs = {
  limit?: Maybe<Scalars['Int']>;
  offset?: Maybe<Scalars['Int']>;
};

export type QueryXrpLiquidityPoolsArgs = {
  currency?: Maybe<Scalars['String']>;
  issuer?: Maybe<Scalars['String']>;
  limit?: Maybe<Scalars['Int']>;
  offset?: Maybe<Scalars['Int']>;
};

// Subscription Types
export type Subscription = {
  __typename?: 'Subscription';
  block: Array<Block>;
  xrpHeatmapUpdates: Array<XRPHeatmapItem>;
  xrpAmmHeatmapUpdates: Array<XRPAMMHeatmapItem>;
};

export type SubscriptionBlockArgs = {
  network: Scalars['String'];
};

export type SubscriptionXrpHeatmapUpdatesArgs = {
  currencies?: Maybe<Array<Scalars['String']>>;
};

export type SubscriptionXrpAmmHeatmapUpdatesArgs = {
  currencies?: Maybe<Array<Scalars['String']>>;
};

// Additional types for portfolio components
export type BalanceItem = {
  symbol: Scalars['String'];
  issuer?: Scalars['String'];
  name: Scalars['String'];
  balance: Scalars['Float'];
  price: Scalars['Float'];
  quote: Scalars['Float'];
  value: Scalars['Float'];
};

export type Balance = {
  items: Array<BalanceItem>;
}; 