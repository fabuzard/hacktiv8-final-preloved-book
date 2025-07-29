package utility

import "fmt"

func BuildTransactionHTMLBody(reqEmail, txnID, product string, amount float64, status, timestamp, invoiceURL string) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Transaction Receipt</title>
		</head>
		<body style="font-family: Arial, sans-serif; background-color: #f7f9fc; padding: 20px;">
			<div style="max-width: 600px; margin: auto; background-color: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);">
				<h2 style="color: #2c3e50;">🧾 Transaction Successful</h2>
				<p>Hi %s,</p>
				<p>Here’s your purchase receipt:</p>

				<table style="width: 100%%; border-collapse: collapse; margin-top: 20px;">
					<tr><td><strong>Transaction ID</strong></td><td>%s</td></tr>
					<tr><td><strong>Product</strong></td><td>%s</td></tr>
					<tr><td><strong>Amount</strong></td><td>$%.2f</td></tr>
					<tr><td><strong>Status</strong></td><td>%s</td></tr>
					<tr><td><strong>Date</strong></td><td>%s</td></tr>
					<tr><td><strong>Invoice</strong></td><td><a href="%s">View Invoice</a></td></tr>
				</table>

				<p style="margin-top: 20px;">Thank you for your transaction!</p>
			</div>
		</body>
		</html>`,
		reqEmail, txnID, product, amount, status, timestamp, invoiceURL,
	)
}
