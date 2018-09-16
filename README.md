# r0ix-assessment

## available endpoints

### 1. Compare
> /compare?ticker_symbol_1={TickerSymbolA}&ticker_symbol_2={TickerSymbolB}
  - This endpoint will compare 1-hour-percent-change of both coin and respond with the name of the coin which has more change.
  - Replace {TickerSymbolA} and {TickerSymbolB} with coin's symbols according to https://api.coinmarketcap.com/v2/listings/
### 2. Key Generate
> /keygen/
  - This endpoint randomly generate Stellar seed(private key) and address(public key)
### 3. Create account
> /account?key={address}
  - This endpoint create Stellar account and provide a curtain amount of balance. This is for testing purpose only.
  - Replace {address} with an account address(public key). To get a new address please see Key Generate endpoint.
  - Using duplicated key will end up receiving errors.
### 4. Account details
> /accountDetail?key={address}
  - This endpoint show Stellar account's details i.e. account balance.
  - Replace {address} with an account address(public key). To get a new address please see Key Generate endpoint.
  - Using unknown account key will end up receiving errors.
### 5. Transfer
> /transfer?source_account={sourceAccountSeed}&destination_account={destAccountAddress}&amount={amount}
  - This endpoint transfer Stellar coin(XLM) from the 'source_account' to the 'destination_account' with the amount of 'amount'. 
  - Replace {sourceAccountSeed} with the source account's seed(private key). To get a new seed please see Key Generate endpoint.
  - Replace {destAccountAddress} with the destination account's address(public key). To get a new address please see Key Generate endpoint.
  - Replace {amount} with the amount of Stellar coin that you want to transfer.
  - Enter the incorrect {sourceAccountSeed}, {destAccountAddress} or {amount} will end up receiving errors.
