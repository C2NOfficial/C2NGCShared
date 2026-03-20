package payu

import "net/url"

// Reference: https://docs.payu.in/reference/check_action_status_api_with_payu_id
type RefundInitiationErrorCodes int

var RefundInitiationErrorCodeMap = map[RefundInitiationErrorCodes]string{
	100: "Refund Successful",
	101: "Refund Successful (Pending)",
	102: "Refund Queued",
	103: "Request rejected on reconfirmation",
	104: "Confirmation required",
	105: "Invalid amount",
	106: "Token already exists",
	107: "Upgraded to refund",
	108: "Refund failure",
	109: "Request is already logged",
	110: "More than one partial refund of Maestro transactions are not allowed",
	111: "Invalid transaction status",
	112: "Risk queued",
	113: "Invalid Amount - Chargeback of amount present",
	115: "Invalid status to be updated",
	116: "Transaction not found",
	117: "Amount does not match",
	119: "No such request found",
	120: "Transaction lock could not be obtained",
	121: "Incorrect or empty value passed in retry",
	122: "Approval pending",
	123: "Request set as pending - requires manual follow-up",
	124: "Input data missing",
	125: "Merchant failed the pending refund",
	126: "Refund in progress",
	127: "Refund requested",
	128: "Partial refunds not allowed",
	129: "Remark is mandatory for retry",
	130: "Refunds not allowed after this period",
	214: "Two refunds of same amount for same transaction within 5 minutes are not allowed",
	225: "Overdraft has occurred - recheck status tomorrow",
	226: "Capture has been initiated today - check refund status tomorrow",
	227: "Transactions with same amount and same token not allowed",
	230: "Purged transaction - refund request requires manual follow-up",
	231: "Refund could not be initiated due to an internal error",
	232: "Refund could not be initiated - either refunds are not supported or need manual intervention",
	233: "Refund blocked from merchant panel - contact support",
	234: "Refund blocked from merchant panel and API - contact support",
	235: "Refund blocked - contact support",
	236: "Refund not possible on this transaction",
	237: "Validation failure - special characters not allowed",
	238: "Validation failure - mandatory field missing",
	239: "API based alternate instant refunds not activated",
	240: "Store card failed",
	241: "Refund not supported by the bank - payment is too old",
	242: "Bank code not supported - raise to PayU support team",
	243: "Virtual account setup to process instant refund is incomplete",
	244: "Beneficiary code for virtual account not set",
	245: "BBPS transaction is not successful",
	246: "Invalid value for the merchant SKU",
	247: "Not allowed - no offers found for the SKU",
	248: "Balance check initiated",
	249: "Retry the transaction",
	250: "Refund failed on uploading successful chargeback",
	251: "Refund blocked for this PGMID by bank",
	252: "Refunds are not allowed from panel for this MID",
	253: "Instant refunds invalid mode",
	254: "Remarks cannot contain special characters",
	255: "Token length exceeded for refund",
	256: "Refund not supported on split transactions - initiate refund on the order transaction",
	258: "Refund initiated",
	259: "Requested retry",
	261: "Error while processing request",
	262: "Error while processing request",
	263: "Invalid requested amount",
	264: "Error while processing request",
	265: "Error while processing request",
	266: "Chargeback is pending against this transaction",
	267: "Lock acquired on transaction metadata",
	270: "Transaction not eligible for instant refund",
	299: "Blocking refund initiation for Type A merchant",
	301: "Capture already successful for this transaction",
	302: "Please try after some time",
	303: "Amount greater than maximum capturable amount",
	304: "Amount less than allowed",
	305: "Amount more than allowed",
	306: "Invalid amount tolerance configuration",
	424: "Transaction upgraded to capture or refund",
	500: "Some exception occurred",
	501: "Successfully updated",
	502: "Failed to update",
}

func (c RefundInitiationErrorCodes) Message() string {
	if msg, ok := RefundInitiationErrorCodeMap[c]; ok {
		return msg
	}
	return "Unknown refund code"
}

// Reference: https://docs.payu.in/reference/refund_transaction_api
type RefundInitiationRequest struct {
	Key     string `json:"key"`     // PayU Merchant Key
	Command string `json:"command"` //The API command name. For refund or canceling transactions, the value should be 'cancel_refund_transaction'.
	Var1    string `json:"var1"`    //mihpayid
	Var2    string `json:"var2"`    //unique token
	Var3    string `json:"var3"`    // amount. Cannot be greater than total amount paid for that transaction
	Var5    string `json:"var5"`    //webhook url
	Hash    string `json:"hash"`    // Format -> sha512(key|command|var1|salt)
}

func (rir *RefundInitiationRequest) SetHash(hash string) {
	rir.Hash = hash
}

func (rir *RefundInitiationRequest) ToURLValues() string {
	values := url.Values{}
	values.Add("key", rir.Key)
	values.Add("command", rir.Command)
	values.Add("var1", rir.Var1)
	values.Add("var2", rir.Var2)
	values.Add("var3", rir.Var3)
	values.Add("var5", rir.Var5)
	values.Add("hash", rir.Hash)
	return values.Encode()
}

// Just to avoid any typos in the code
type InitiateRefundResponseStatus uint8

const (
	InitiateRefundStatusSuccess InitiateRefundResponseStatus = 0
	InitiateRefundStatusFailure InitiateRefundResponseStatus = 1
)

// Reference: https://docs.payu.in/reference/refund_transaction_api
type InitiateRefundResponse struct {
	Status     InitiateRefundResponseStatus `json:"status"`       // 1 if API call is a success 0 if the API has failed
	Message    string                       `json:"msg"`          // This parameter contains a response message description
	RequestID  string                       `json:"request_id"`   // This parameter contains a unique refund ID generated by PayU
	BankRefNum string                       `json:"bank_ref_num"` // This parameter contains a bank reference number returned from bank
	MihpayID   string                       `json:"mihpayid"`     // This parameter contains a unique transaction ID generated by PayU during sale
	ErrorCode  RefundInitiationErrorCodes   `json:"error_code"`   // This parameter contains the code for response. For a list of error codes and their description, refer to Refund Error Codes
}

type RefundStatus string

const (
	RefundStatusSuccess RefundStatus = "success"
	RefundStatusFailure RefundStatus = "failure"
)

// Reference: https://docs.payu.in/reference/refund-status-callback
type RefundWebhookResponse struct {
	AdditionalValue1 *string      `json:"additionalValue1"`
	BankARN          int64        `json:"bank_arn"`
	RefundMode       string       `json:"refund_mode"`
	BankRefNum       string       `json:"bank_ref_num"`
	Key              string       `json:"key"`
	Amount           string       `json:"amt"`
	Remark           *string      `json:"remark"`
	Status           RefundStatus `json:"status"`
	Token            string       `json:"token"`
	MihpayID         string       `json:"mihpayid"`
	RequestID        string       `json:"request_id"`
	MerchantTxnID    string       `json:"merchantTxnId"`
	AdditionalValue2 *string      `json:"additionalValue2"`
	Action           string       `json:"action"`
}

// Only need this to verify the refund webhook response. 
// So only key, mihpayid and command are required.
func (rwr *RefundWebhookResponse) ToRefundInitateRequest() *RefundInitiationRequest {
	return &RefundInitiationRequest{
		Key: Key, 
		Command: "cancel_refund_transaction",
		Var1: rwr.MihpayID,
	}
}