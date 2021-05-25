import { Keypair, Address } from '@helium/crypto'
import { PaymentV2, Transaction } from '@helium/transactions'
import express = require('express');
var bodyParser = require('body-parser');
import http = require('http');

var app = express();
app.use(bodyParser.json());
app.set('port', process.env.PORT || 3000);
app.post('/get-nonce', function(req: express.Request, res: express.Response) {
  try {
    const account:Address = Address.fromB58(req.body.requested_metadata.get_nonce_for);
    console.log(account.b58);
    res.send(account)
  } catch(e:any) {
    res.send(500, e);
  }
  
});

http.createServer(app).listen(app.get('port'), function() {
  console.log('Express server listening on port ' + app.get('port'));
});