import { Address, NetType } from '@helium/crypto'
import proto from '@helium/proto'
import * as utils from './utils'
import { PaymentV2, PaymentV1, Transaction } from '@helium/transactions'
import { Client, Network, PendingTransaction } from '@helium/http'
import * as express from "express"
import * as http from "http"
import { PaymentV2Json } from './transaction_types'
import * as JSLong from "long"
import * as crypto from "crypto"
import base64url from "base64url"
import * as winston from "winston"

const express = require('express');
const bodyParser = require('body-parser');
const asyncHandler = require('express-async-handler');
var app = express();

const logger = winston.createLogger();
logger.add(new winston.transports.Console({
  format: winston.format.combine(
    winston.format.colorize(),
    winston.format.timestamp(),
    winston.format.align(),
    winston.format.printf(info => `${info.timestamp} ${info.level}: ${info.message}`)
  )
}));

var netType:number;
var clientType:Network;

if (process.env.NETWORK == "testnet") {
    logger.info("Starting testnet server");
    netType = 1;
    clientType = Network.testnet;
} else {
    logger.info("Starting mainnet server");
    netType = 0;
    clientType = Network.production;
}

app.use(bodyParser.json());
app.set('port', process.env.PORT || 3000);
app.post('/create-tx', function(req: express.Request, res: express.Response) {
  try {
    const vars = req.body["chain_vars"];
    const transactionType = req.body["options"]["transaction_type"];
    Transaction.config(vars); 

    switch (transactionType) {
      case "payment_v2":
        const payments = []
        req.body["options"]["helium_metadata"]["payments"].forEach(payment => {
          
          // Ensure 8-byte base64 encoded memo is available for each payment to generate proper signature
          const memo = payment["memo"] ? payment["memo"] : "MDAwMDAwMDA=";
          if (Buffer.from(memo, "base64").length != 8) {
            res.status(200).send({ error: "invalid memo" });
          }

          payments.push({
            "payee": Address.fromB58(payment["payee"]),
            "amount": payment["amount"],
            "memo": memo
          })
        });

        // Create helium-js transaction
        const unsignedPaymentV2Txn:PaymentV2 = new PaymentV2({
          payer: Address.fromB58(req.body["options"]["helium_metadata"]["payer"]),
          payments: payments,
          nonce: req.body["get_nonce_for"]["nonce"] + 1
        });

        // Construct protobuf payload from helium-js transaction
        const Payment = proto.helium.payment
        const paymentsProto = unsignedPaymentV2Txn.payments.map(({ payee, amount, memo }) => {
          const memoBuffer = memo ? Buffer.from(memo, 'base64') : undefined
          return Payment.create({
            payee: Uint8Array.from(Buffer.from(payee.bin)),
            amount,
            memo: memoBuffer ? JSLong.fromBytes(Array.from(memoBuffer), true, true) : undefined,
          })
        })

        const PaymentV2Txn = proto.helium.blockchain_txn_payment_v2
        const paymentV2Proto = PaymentV2Txn.create({
            payer: unsignedPaymentV2Txn.payer ? Uint8Array.from(Buffer.from(unsignedPaymentV2Txn.payer.bin)) : undefined,
            payments: paymentsProto,
            fee: unsignedPaymentV2Txn.fee,
            nonce: unsignedPaymentV2Txn.nonce
        })

        // Enocde protobuf payload for signing and convert to hex string
        const serialized = proto.helium.blockchain_txn_payment_v2.encode(paymentV2Proto).finish();
        const hex_bytes:string = Buffer.from(serialized).toString("hex");

        res.status(200).send({"unsigned_txn": unsignedPaymentV2Txn.toString(), "type": "payment_v2", "payload": hex_bytes });
        break;
      default:
        res.status(500).send({ error: "Unrecognized transaction type: " +  transactionType });
        break;
    }
  } catch(e:any) {
    res.status(500).send({ error: e });
  }
});

app.post('/combine-tx', function(req: express.Request, res: express.Response) {
  try {
    const rawUnsignedTxn:string = req.body["unsigned_transaction"]
    const unsignedTxnType:string = Transaction.stringType(rawUnsignedTxn);

    switch (unsignedTxnType) {
      case "paymentV2":
        const signature:string = req.body["signatures"][0]["hex_bytes"];
        const payment:PaymentV2 = PaymentV2.fromString(rawUnsignedTxn);
        payment.signature = Uint8Array.from(Buffer.from(signature, "hex"));
        res.status(200).send({ signed_transaction: payment.toString() });
        break;
      default:
        res.status(500).send({ error: "unrecognized transaction type: " + unsignedTxnType });
        break;
    }
  } catch(e:any) {
    res.status(500).send({ error: e });
  }
});

app.post('/parse-tx', function(req: express.Request, res: express.Response) {
  try {
    const rawTxn:string = req.body["raw_transaction"];
    const txnType:string = Transaction.stringType(rawTxn);

    switch (txnType) {
      case "paymentV2":
        const paymentV2:PaymentV2 = PaymentV2.fromString(rawTxn);
        const payload:PaymentV2Json = utils.paymentV2toJson(paymentV2);
        const resp:Object = {
          "payload": payload          
        }
        
        if (req.body["signed"]) {
          resp["signer"] = payload.payer;
        }

        res.status(200).send(resp);
      default:
        res.status(500);
    }

  } catch(e:any) {
    res.status(500).send({ error: e });
  }
});

app.get('/chain-vars', asyncHandler(async function(req: express.Request, res: express.Response) {
  const client:Client = new Client(clientType);
  const vars = await client.vars.get();
  res.status(200).send(vars);
}));

app.get('/current-height', asyncHandler(async function(req: express.Request, res: express.Response) {
  const client:Client = new Client(clientType);
  const currentHeight = await client.blocks.getHeight();
  res.status(200).send({ current_height: currentHeight});
}));

app.post('/hash', function(req: express.Request, res: express.Response){
  try {
    const txnString:string = req.body["txn"];
    const txnType:string = Transaction.stringType(txnString);
    
    switch (txnType) {
      case "paymentV2":
        const p = PaymentV2.fromString(req.body["txn"]);
        p.signature = undefined;

        const Payment = proto.helium.payment
        const payments = p.payments.map(({ payee, amount, memo }) => {
          const memoBuffer = memo ? Buffer.from(memo, 'base64') : undefined
          return Payment.create({
            payee: Uint8Array.from(Buffer.from(payee.bin)),
            amount,
            memo: memoBuffer ? JSLong.fromBytes(Array.from(memoBuffer), true, true) : undefined,
          })
        });

        const PaymentTxn = proto.helium.blockchain_txn_payment_v2 
        const PaymentTxnPB = PaymentTxn.create({
            payer: Uint8Array.from(Buffer.from(p.payer.bin)),
            payments,
            fee: p.fee,
            nonce: p.nonce
        })
        const serializedPaymentTxnPB = proto.helium.blockchain_txn_payment_v2.encode(PaymentTxnPB);
        
        res.status(200).send({ 
          hash: base64url.fromBase64(crypto.createHash("sha256").update(serializedPaymentTxnPB.finish()).digest("base64")) 
        });
        break;
      default:
        res.status(500).send({
          error: "Transaction not recognized"
        });
        break;
    }
  } catch(e:any) {
    res.status(500).send({ error: e });
  }
});

app.post('/derive', function(req: express.Request, res: express.Response) {
  try {
    const curveType: string = req.body["curve_type"];
    const publicKey: string = req.body["public_key"];

    if (curveType != "edwards25519") {
      throw "curve type " + curveType + " not surrported";
    }

    // Add nettype and keytype as part of hex string for first byte
    res.status(200).send({ 
      address: Address.fromBin(Buffer.from(netType+'1'+publicKey, "hex")).b58
    });
    
  } catch(e:any) {
    res.status(500).send({ error: e });
  }
});

app.post('/submit-tx', asyncHandler(async function(req: express.Request, res: express.Response) {
  const signedTransaction: string = req.body["signed_transaction"];
  const client = new Client(clientType);

  const pendingTransaction: PendingTransaction = await client.transactions.submit(signedTransaction);
  res.status(200).send({
    hash: pendingTransaction.hash
  });
}));

http.createServer(app).listen(app.get('port'), function() {
  logger.info('Express server listening on port ' + app.get('port'));
});