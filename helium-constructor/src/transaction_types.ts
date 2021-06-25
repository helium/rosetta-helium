interface PaymentJson {
    amount: number
    payee: string
    memo?: string
}

interface PaymentV2Json {
    type: string
    payer: string
    nonce: number
    fee: number
    payments: any[]
    signature?: string
}

export { 
    PaymentJson,
    PaymentV2Json
}