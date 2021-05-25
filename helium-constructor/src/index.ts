import { Keypair, Address } from '@helium/crypto'
import { PaymentV2, Transaction } from '@helium/transactions'
import express = require('express');
var bodyParser = require('body-parser');
import http = require('http');

var app = express();
app.use(bodyParser.json());
app.set('port', process.env.PORT || 3000);
app.post('/', function(req: express.Request, res: express.Response) {
  console.log(req.body);
  res.send('Hello world!');
});

http.createServer(app).listen(app.get('port'), function() {
  console.log('Express server listening on port ' + app.get('port'));
});