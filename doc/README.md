### 說明文件


1. 接口文檔：
   https://github.com/binance/binance-spot-api-docs/blob/master/web-socket-streams_CN.md
   需將文檔中的 wss://stream.binance.com:9443網站地址：
   改為 wss://stream.yshyqxx.com/stream

2. Websocket在線測試：http://coolaf.com/tool/chattest

### 要求
1. 訂閱近期成交(BTCUSDT歸集)功能
wss://stream.yshyqxx.com/stream?streams=btcusdt@aggTrade

2. 將最新的一筆行情存入redis ，鍵名為“streams=btcusdt@aggTrade”

3. 客戶端能訪問獲取最新一筆行情

4. 客戶端獲取最新行情的推送
