import { Address } from "@helium/crypto";
import { PaymentV2 } from "@helium/transactions"
import { PaymentJson, PaymentV2Json } from "./transaction_types"

function paymentV2toJson(p:PaymentV2):PaymentV2Json {
    let payments:PaymentJson[] = [];
    p.payments.forEach(payment => {
        payments.push({
            amount: payment.amount,
            payee: payment.payee.b58,
            memo: payment?.memo
        })
    })
    const paymentV2Json:PaymentV2Json = {
        type: "paymentV2",
        payer: p.payer.b58,
        nonce: p.nonce,
        fee: p.fee,
        payments: payments
    }

    return paymentV2Json;
}

export {
    paymentV2toJson,
}