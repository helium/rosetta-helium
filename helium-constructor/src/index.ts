import { Keypair, Address } from '@helium/crypto'
import proto from '@helium/proto'
import { PaymentV2, PaymentV1, Transaction } from '@helium/transactions'
import { Client } from '@helium/http'

import express = require('express');
import http = require('http');
const bodyParser = require('body-parser');
const asyncHandler = require('express-async-handler');

var app = express();
app.use(bodyParser.json());
app.set('port', process.env.PORT || 3000);
app.post('/create-tx', function(req: express.Request, res: express.Response) {
  try {
    const vars = req.body["chain_vars"];
    Transaction.config(vars); 

    switch (req.body["options"]["transaction_type"]) {
      case "payment_v2":
        const payments = []
        req.body["options"]["helium_metadata"]["payments"].forEach(payment => {
          payments.push({
            "payee": payment["payee"],
            "amount": payment["amount"]
          })
        });

        const unsignedPaymentV2Txn:PaymentV2 = new PaymentV2({
          payer: Address.fromB58(req.body["options"]["helium_metadata"]["payer"]),
          payments: payments,
          nonce: req.body["get_nonce_for"]["nonce"] + 1
        });

        const hex_bytes:string = Buffer.from(unsignedPaymentV2Txn.serialize()).toString('hex');

        res.status(200).send({"unsigned_txn": unsignedPaymentV2Txn.toString(), "type": "payment_v2", "payload": hex_bytes });
        break;
      default:
        res.status(500);
        break;
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