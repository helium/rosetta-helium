import { Keypair, Address } from '@helium/crypto'
import { PaymentV2, PaymentV1, Transaction } from '@helium/transactions'
import { Client } from '@helium/http'

import express = require('express');
import http = require('http');
const bodyParser = require('body-parser');
const asyncHandler = require('express-async-handler');

var app = express();
app.use(bodyParser.json());
app.set('port', process.env.PORT || 3000);
app.post('/get-fee', function(req: express.Request, res: express.Response) {
  try {
    if (!req.body.transaction_type) {
      throw { "error": "`transaction_type` required" }
    }
    const transaction_type:string = req.body.transaction_type;

    switch (transaction_type) {
      case "payment_v2":
        console.log(transaction_type)
        const payment:PaymentV1 = new PaymentV1({
          payer: Address.fromB58("13HPSdf8Ng8E2uKpLm8Ba3sQ6wdNimTcaKXYmMkHyTUUeUELPwJ"),
          payee: Address.fromB58("1aCjThQENE7h1r8qQ52H2P1hCN53uBR6sVrr4MKJPh4Bg8dVqbY"),
          amount: 100,
          nonce: 1311,
        })
        console.log(payment.fee)
        res.send(200, payment.fee)
        break;
      default:
        throw { "error": "transaction type '" + transaction_type + "' is not valid"}
    }
  } catch(e:any) {
    res.status(500).send(e);
  }
  
});

app.get('/chain-vars', asyncHandler(async function(req: express.Request, res: express.Response) {
  const client:Client = new Client();
  const vars = await client.vars.get();
  res.status(200).send(vars);
}));

http.createServer(app).listen(app.get('port'), function() {
  console.log('Express server listening on port ' + app.get('port'));
});