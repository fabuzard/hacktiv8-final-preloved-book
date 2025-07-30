package utility

import "fmt"

func GenerateVerificationTokenHTML(token string) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Verify Your Email</title>
			<style>
				body { font-family: Arial, sans-serif; background-color: #f9f9f9; padding: 20px; }
				.container { background-color: white; padding: 30px; border-radius: 8px; box-shadow: 0 0 10px rgba(0,0,0,0.1); }
				.token { font-size: 24px; font-weight: bold; color: #2c3e50; background: #ecf0f1; padding: 12px 20px; display: inline-block; border-radius: 6px; letter-spacing: 2px; margin: 20px 0; }
				p { font-size: 16px; color: #333; }
			</style>
		</head>
		<body>
			<div class="container">
				<h2>Email Verification</h2>
				<p>Use the verification code below to complete your registration. The code will expire in 15 minutes.</p>
				<div class="token">%s</div>
				<p>If you didn’t request this, you can safely ignore this email.</p>
			</div>
		</body>
		</html>
	`, token)
}
