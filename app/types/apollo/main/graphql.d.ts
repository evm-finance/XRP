
declare module '*/config.query.graphql' {
  import { DocumentNode } from 'graphql';
  const defaultDocument: DocumentNode;
  export const SupportedChainsGQL: DocumentNode;
export const GasGQL: DocumentNode;
export const DeFiStats: DocumentNode;
export const UniswapTokensGQL: DocumentNode;
export const RecentPricesGQL: DocumentNode;
export const FiatPricesGQL: DocumentNode;

  export default defaultDocument;
}
    

declare module '*/pools.query.graphql' {
  import { DocumentNode } from 'graphql';
  const defaultDocument: DocumentNode;
  export const AavePoolGQL: DocumentNode;

  export default defaultDocument;
}
    

declare module '*/portfolio.query.graphql' {
  import { DocumentNode } from 'graphql';
  const defaultDocument: DocumentNode;
  export const BalancesGQL: DocumentNode;
export const TransactionsGQL: DocumentNode;
export const TransfersGQL: DocumentNode;

  export default defaultDocument;
}
    

declare module '*/token.query.graphql' {
  import { DocumentNode } from 'graphql';
  const defaultDocument: DocumentNode;
  export const TokenQueryGQL: DocumentNode;
export const DailyChartGQL: DocumentNode;
export const ScreenerGQL: DocumentNode;
export const TimeGQL: DocumentNode;
export const PriceStreamGQL: DocumentNode;
export const BlocksGQL: DocumentNode;
export const BlocksXrpGQL: DocumentNode;
export const BlockGQL: DocumentNode;
export const BlocksStreamGQL: DocumentNode;

  export default defaultDocument;
}
    