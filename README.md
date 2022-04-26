## :money_with_wings: Fat Stack :money_with_wings:

The goal of this project is to have a fully backtestable trading system.

Any strategies, eg. MACD, RSI, should be configurable, and testable on different time intervals for a given set of tick data (trade data)


### Running the application

To run the application you need [go installed](https://go.dev/doc/install)

Once you have go installed, on Linux and Mac's just run

```
make build
```

This will produce an executable in the `bin` folder.

Then you just need to run it.

```
./bin/fat-stacks run
```

The output currently only writes all Binance Futures aggregated trades to the console.