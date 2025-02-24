Example project written in go using gin-gonic

Endpoint rates:
```
GET /rates?currencies=EUR,GBP
```

Example response:
```json
[
  { "from": "EUR", "to": "GBP", "rate": 1.0 },
  { "from": "GBP", "to": "EUR", "rate": 1.0 }
]
```

returns rates for all currency pairs: EUR-GBP, GBP-EUR fetched from openexchangesrates.org API.

Endpoint exchange:
```
GET /exchange?from=WBTC&to=USDT&amount=1.0
```

Example response:
```json
{"from": "WBTC", "to": "USDT", "amount": 57094.314314}
```

How to run:
```bash
cp dev.env .env
```

Create APP ID on openexchangesrates.org website (registration is needed)

Paste your APP ID into .env file (variable OPEN_EXCHANGE_RATES_APP_ID).

Build docker image:
```bash
make build-image
```

Run docker image:
```bash
make run-image
```

Test:
```bash
curl "localhost:8080/rates?currencies=EUR,GBP"
```
