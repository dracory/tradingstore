package tradingstore

// Asset Class (from https://www.investopedia.com/terms/a/asset-class.asp)
const ASSET_CLASS_BOND = "BOND"
const ASSET_CLASS_COMMODITY = "COMMODITY"   // Commodities, such as gold, oil, and agricultural products
const ASSET_CLASS_CRYPTO = "CRYPTO"         // Cryptocurrency
const ASSET_CLASS_CURRENCY = "CURRENCY"     // Currencies
const ASSET_CLASS_DERIVATIVE = "DERIVATIVE" // Derivatives
const ASSET_CLASS_FOREX = "FOREX"           // Foreign Exchange
const ASSET_CLASS_ETF = "ETF"               // Exchange Traded Funds
const ASSET_CLASS_FUTURE = "FUTURE"         // Futures
const ASSET_CLASS_INDEX = "INDEX"           // Index
const ASSET_CLASS_OPTION = "OPTION"         // Options
const ASSET_CLASS_REIT = "REIT"             // Real Estate Investment Trust
const ASSET_CLASS_STOCK = "STOCK"           // Stocks
const ASSET_CLASS_UNKNOWN = "UNKNOWN"       // Unknown

// Column names
const COLUMN_ASSET_CLASS = "asset_class"
const COLUMN_CLOSE = "close"
const COLUMN_CREATED_AT = "created_at"
const COLUMN_DESCRIPTION = "description"
const COLUMN_EXCHANGE = "exchange"
const COLUMN_ID = "id"
const COLUMN_HIGH = "high"
const COLUMN_LOW = "low"
const COLUMN_MEMO = "memo"
const COLUMN_NAME = "name"
const COLUMN_METAS = "metas"
const COLUMN_OPEN = "open"
const COLUMN_SOFT_DELETED_AT = "soft_deleted_at"
const COLUMN_STATUS = "status"
const COLUMN_SYMBOL = "symbol"
const COLUMN_TIME = "time"
const COLUMN_TIMEFRAMES = "timeframes"
const COLUMN_UPDATED_AT = "updated_at"
const COLUMN_VOLUME = "volume"

// Nil float
const NIL_FLOAT = -0.0000000001

// Instrument Status
const INSTRUMENT_STATUS_DRAFT = "draft"
const INSTRUMENT_STATUS_ACTIVE = "active"
const INSTRUMENT_STATUS_INACTIVE = "inactive"
const INSTRUMENT_STATUS_DISABLED = "disabled"

// Timeframe
const TIMEFRAME_1_MINUTE = "1min"
const TIMEFRAME_5_MINUTES = "5min"
const TIMEFRAME_15_MINUTES = "15min"
const TIMEFRAME_30_MINUTES = "30min"
const TIMEFRAME_1_HOUR = "1hour"
const TIMEFRAME_4_HOURS = "4hour"
const TIMEFRAME_1_DAY = "1day"
const TIMEFRAME_1_WEEK = "1week"
const TIMEFRAME_1_MONTH = "1month"
const TIMEFRAME_1_YEAR = "1year"
