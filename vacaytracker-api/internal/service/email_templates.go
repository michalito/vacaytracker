package service

// Email template data structures
type welcomeEmailData struct {
	AppURL       string
	UserName     string
	UserEmail    string
	TempPassword string
}

type vacationEmailData struct {
	AppURL    string
	UserName  string
	StartDate string
	EndDate   string
	TotalDays int
	Reason    string // Only used for rejections
}

type adminNotificationData struct {
	AppURL        string
	RequesterName string
	StartDate     string
	EndDate       string
	TotalDays     int
	RequestReason string
}

// Welcome email templates
const welcomeEmailSubject = "Welcome to VacayTracker!"

const welcomeEmailHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to VacayTracker</title>
    <!--[if mso]>
    <noscript>
        <xml>
            <o:OfficeDocumentSettings>
                <o:PixelsPerInch>96</o:PixelsPerInch>
            </o:OfficeDocumentSettings>
        </xml>
    </noscript>
    <![endif]-->
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #e6f7f9;">
    <!-- Preheader text (shows in inbox preview) -->
    <div style="display: none; max-height: 0; overflow: hidden; mso-hide: all;">
        Your VacayTracker account is ready! Log in to start tracking your time off.
        &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847;
    </div>
    <table role="presentation" style="width: 100%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 20px;">
                <table role="presentation" style="width: 600px; max-width: 100%; border-collapse: collapse; background-color: #ffffff; border-radius: 16px; box-shadow: 0 4px 20px rgba(13, 131, 162, 0.08);">
                    <!-- Header with Logo -->
                    <tr>
                        <td style="padding: 32px 40px 24px; text-align: center;">
                            <img src="{{.AppURL}}/logo.png" width="64" height="64" alt="VacayTracker" style="height: 64px; width: 64px; display: block; margin: 0 auto 16px; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; font-size: 18px; font-weight: 600; color: #0D83A2;">
                            <h1 style="margin: 0; color: #00384F; font-size: 24px; font-weight: 600; letter-spacing: -0.5px;">Welcome Aboard!</h1>
                        </td>
                    </tr>
                    <!-- Status Bar -->
                    <tr>
                        <td style="padding: 0; height: 4px; background: linear-gradient(90deg, #0D83A2 0%, #15ABCB 100%); background-color: #0D83A2;" bgcolor="#0D83A2"></td>
                    </tr>
                    <!-- Content -->
                    <tr>
                        <td style="padding: 32px 40px;">
                            <p style="margin: 0 0 20px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Ahoy, <strong style="color: #00384F;">{{.UserName}}</strong>!
                            </p>
                            <p style="margin: 0 0 24px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Your VacayTracker account has been created. You can now start planning your well-deserved time off!
                            </p>
                            <!-- Credentials Box -->
                            <div style="background-color: #f0f9ff; border-radius: 12px; padding: 20px; margin: 0 0 24px;">
                                <p style="margin: 0 0 12px; color: #0D83A2; font-size: 14px; font-weight: 600;">Your Login Credentials</p>
                                <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                    <tr>
                                        <td style="padding: 6px 0; color: #6b7280; font-size: 14px;">Email</td>
                                        <td style="padding: 6px 0; color: #00384F; font-size: 14px; font-weight: 500; text-align: right;">
                                            <code style="background-color: #e0f2fe; padding: 3px 8px; border-radius: 4px; font-family: monospace;">{{.UserEmail}}</code>
                                        </td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 6px 0; color: #6b7280; font-size: 14px;">Temporary Password</td>
                                        <td style="padding: 6px 0; color: #00384F; font-size: 14px; font-weight: 500; text-align: right;">
                                            <code style="background-color: #e0f2fe; padding: 3px 8px; border-radius: 4px; font-family: monospace;">{{.TempPassword}}</code>
                                        </td>
                                    </tr>
                                </table>
                            </div>
                            <!-- Security Note -->
                            <p style="margin: 0 0 28px; color: #991b1b; font-size: 14px; line-height: 1.5; padding: 12px 16px; background-color: #fef2f2; border-radius: 8px;">
                                <strong>Important:</strong> Please change your password after your first login for security purposes.
                            </p>
                            <!-- CTA Button -->
                            <div style="text-align: center;">
                                <a href="{{.AppURL}}" style="display: inline-block; padding: 14px 32px; background-color: #0D83A2; color: #ffffff; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px; box-shadow: 0 2px 8px rgba(13, 131, 162, 0.25);">Login to VacayTracker</a>
                            </div>
                        </td>
                    </tr>
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 24px 40px; background-color: #e6f7f9; border-radius: 0 0 16px 16px; text-align: center; border-top: 1px solid #cceff3;">
                            <p style="margin: 0 0 4px; color: #0a6a84; font-size: 13px; font-weight: 500;">VacayTracker</p>
                            <p style="margin: 0; color: #6b7280; font-size: 12px;">Your vacation tracking companion</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`

const welcomeEmailText = `Welcome Aboard, {{.UserName}}!

Your VacayTracker account has been created. You can now start planning your well-deserved time off!

Your login credentials:
- Email: {{.UserEmail}}
- Temporary Password: {{.TempPassword}}

Please change your password after your first login for security purposes.

Login at: {{.AppURL}}

---
VacayTracker - Your vacation tracking companion`

// Request submitted email templates
const requestSubmittedSubject = "Vacation Request Submitted"

const requestSubmittedHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Vacation Request Submitted</title>
    <!--[if mso]>
    <noscript>
        <xml>
            <o:OfficeDocumentSettings>
                <o:PixelsPerInch>96</o:PixelsPerInch>
            </o:OfficeDocumentSettings>
        </xml>
    </noscript>
    <![endif]-->
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #e6f7f9;">
    <!-- Preheader text (shows in inbox preview) -->
    <div style="display: none; max-height: 0; overflow: hidden; mso-hide: all;">
        Your vacation request has been submitted and is awaiting approval.
        &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847;
    </div>
    <table role="presentation" style="width: 100%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 20px;">
                <table role="presentation" style="width: 600px; max-width: 100%; border-collapse: collapse; background-color: #ffffff; border-radius: 16px; box-shadow: 0 4px 20px rgba(13, 131, 162, 0.08);">
                    <!-- Header with Logo -->
                    <tr>
                        <td style="padding: 32px 40px 24px; text-align: center;">
                            <img src="{{.AppURL}}/logo.png" width="64" height="64" alt="VacayTracker" style="height: 64px; width: 64px; display: block; margin: 0 auto 16px; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; font-size: 18px; font-weight: 600; color: #0D83A2;">
                            <h1 style="margin: 0; color: #00384F; font-size: 24px; font-weight: 600; letter-spacing: -0.5px;">Request Submitted</h1>
                        </td>
                    </tr>
                    <!-- Status Bar (Amber for Pending) -->
                    <tr>
                        <td style="padding: 0; height: 4px; background: linear-gradient(90deg, #f59e0b 0%, #fbbf24 100%); background-color: #f59e0b;" bgcolor="#f59e0b"></td>
                    </tr>
                    <!-- Content -->
                    <tr>
                        <td style="padding: 32px 40px;">
                            <p style="margin: 0 0 20px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Hi <strong style="color: #00384F;">{{.UserName}}</strong>,
                            </p>
                            <p style="margin: 0 0 24px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Your vacation request has been submitted and is pending approval.
                            </p>
                            <!-- Details Box -->
                            <div style="background-color: #f8fafc; border: 1px solid #e2e8f0; border-radius: 12px; padding: 20px; margin: 0 0 24px;">
                                <div style="display: inline-block; padding: 4px 12px; background-color: #fffbeb; color: #92400e; font-size: 12px; font-weight: 600; border-radius: 20px; margin-bottom: 12px;">Pending Review</div>
                                <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px;">Start Date</td>
                                        <td style="padding: 8px 0; color: #00384F; font-size: 14px; font-weight: 600; text-align: right;">{{.StartDate}}</td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px;">End Date</td>
                                        <td style="padding: 8px 0; color: #00384F; font-size: 14px; font-weight: 600; text-align: right;">{{.EndDate}}</td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px;">Total Days</td>
                                        <td style="padding: 8px 0; color: #00384F; font-size: 14px; font-weight: 600; text-align: right;">{{.TotalDays}}</td>
                                    </tr>
                                </table>
                            </div>
                            <p style="margin: 0 0 28px; color: #6b7280; font-size: 14px; line-height: 1.6;">
                                You'll receive another email once your request has been reviewed.
                            </p>
                            <!-- CTA Button -->
                            <div style="text-align: center;">
                                <a href="{{.AppURL}}/employee" style="display: inline-block; padding: 14px 32px; background-color: #0D83A2; color: #ffffff; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px; box-shadow: 0 2px 8px rgba(13, 131, 162, 0.25);">View Dashboard</a>
                            </div>
                        </td>
                    </tr>
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 24px 40px; background-color: #e6f7f9; border-radius: 0 0 16px 16px; text-align: center; border-top: 1px solid #cceff3;">
                            <p style="margin: 0 0 4px; color: #0a6a84; font-size: 13px; font-weight: 500;">VacayTracker</p>
                            <p style="margin: 0; color: #6b7280; font-size: 12px;">Your vacation tracking companion</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`

const requestSubmittedText = `Hi {{.UserName}},

Your vacation request has been submitted and is pending approval.

Request Details:
- Start Date: {{.StartDate}}
- End Date: {{.EndDate}}
- Total Days: {{.TotalDays}}

You'll receive another email once your request has been reviewed.

View your dashboard at: {{.AppURL}}/employee

---
VacayTracker - Your vacation tracking companion`

// Request approved email templates
const requestApprovedSubject = "Vacation Request Approved!"

const requestApprovedHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Vacation Request Approved</title>
    <!--[if mso]>
    <noscript>
        <xml>
            <o:OfficeDocumentSettings>
                <o:PixelsPerInch>96</o:PixelsPerInch>
            </o:OfficeDocumentSettings>
        </xml>
    </noscript>
    <![endif]-->
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #e6f7f9;">
    <!-- Preheader text (shows in inbox preview) -->
    <div style="display: none; max-height: 0; overflow: hidden; mso-hide: all;">
        Great news! Your vacation request has been approved. Time to start planning!
        &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847;
    </div>
    <table role="presentation" style="width: 100%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 20px;">
                <table role="presentation" style="width: 600px; max-width: 100%; border-collapse: collapse; background-color: #ffffff; border-radius: 16px; box-shadow: 0 4px 20px rgba(13, 131, 162, 0.08);">
                    <!-- Header with Logo -->
                    <tr>
                        <td style="padding: 32px 40px 24px; text-align: center;">
                            <img src="{{.AppURL}}/logo.png" width="64" height="64" alt="VacayTracker" style="height: 64px; width: 64px; display: block; margin: 0 auto 16px; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; font-size: 18px; font-weight: 600; color: #0D83A2;">
                            <h1 style="margin: 0; color: #00384F; font-size: 24px; font-weight: 600; letter-spacing: -0.5px;">You're All Set!</h1>
                        </td>
                    </tr>
                    <!-- Status Bar (Green for Approved) -->
                    <tr>
                        <td style="padding: 0; height: 4px; background: linear-gradient(90deg, #22c55e 0%, #4ade80 100%); background-color: #22c55e;" bgcolor="#22c55e"></td>
                    </tr>
                    <!-- Content -->
                    <tr>
                        <td style="padding: 32px 40px;">
                            <p style="margin: 0 0 20px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Great news, <strong style="color: #00384F;">{{.UserName}}</strong>!
                            </p>
                            <p style="margin: 0 0 24px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Your vacation request has been approved. Time to start planning your getaway!
                            </p>
                            <!-- Details Box -->
                            <div style="background-color: #f8fafc; border: 1px solid #e2e8f0; border-radius: 12px; padding: 20px; margin: 0 0 28px;">
                                <div style="display: inline-block; padding: 4px 12px; background-color: #f0fdf4; color: #166534; font-size: 12px; font-weight: 600; border-radius: 20px; margin-bottom: 12px;">Approved</div>
                                <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px;">Start Date</td>
                                        <td style="padding: 8px 0; color: #00384F; font-size: 14px; font-weight: 600; text-align: right;">{{.StartDate}}</td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px;">End Date</td>
                                        <td style="padding: 8px 0; color: #00384F; font-size: 14px; font-weight: 600; text-align: right;">{{.EndDate}}</td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px;">Total Days</td>
                                        <td style="padding: 8px 0; color: #00384F; font-size: 14px; font-weight: 600; text-align: right;">{{.TotalDays}}</td>
                                    </tr>
                                </table>
                            </div>
                            <!-- CTA Button -->
                            <div style="text-align: center;">
                                <a href="{{.AppURL}}/employee" style="display: inline-block; padding: 14px 32px; background-color: #0D83A2; color: #ffffff; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px; box-shadow: 0 2px 8px rgba(13, 131, 162, 0.25);">View Dashboard</a>
                            </div>
                        </td>
                    </tr>
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 24px 40px; background-color: #e6f7f9; border-radius: 0 0 16px 16px; text-align: center; border-top: 1px solid #cceff3;">
                            <p style="margin: 0 0 4px; color: #0a6a84; font-size: 13px; font-weight: 500;">VacayTracker</p>
                            <p style="margin: 0; color: #6b7280; font-size: 12px;">Your vacation tracking companion</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`

const requestApprovedText = `Great news, {{.UserName}}!

Your vacation request has been approved. Time to start planning your getaway!

Approved Vacation:
- Start Date: {{.StartDate}}
- End Date: {{.EndDate}}
- Total Days: {{.TotalDays}}

View your dashboard at: {{.AppURL}}/employee

---
VacayTracker - Your vacation tracking companion`

// Request rejected email templates
const requestRejectedSubject = "Vacation Request Update"

const requestRejectedHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Vacation Request Update</title>
    <!--[if mso]>
    <noscript>
        <xml>
            <o:OfficeDocumentSettings>
                <o:PixelsPerInch>96</o:PixelsPerInch>
            </o:OfficeDocumentSettings>
        </xml>
    </noscript>
    <![endif]-->
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #e6f7f9;">
    <!-- Preheader text (shows in inbox preview) -->
    <div style="display: none; max-height: 0; overflow: hidden; mso-hide: all;">
        Your vacation request needs attention. Please review the details inside.
        &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847;
    </div>
    <table role="presentation" style="width: 100%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 20px;">
                <table role="presentation" style="width: 600px; max-width: 100%; border-collapse: collapse; background-color: #ffffff; border-radius: 16px; box-shadow: 0 4px 20px rgba(13, 131, 162, 0.08);">
                    <!-- Header with Logo -->
                    <tr>
                        <td style="padding: 32px 40px 24px; text-align: center;">
                            <img src="{{.AppURL}}/logo.png" width="64" height="64" alt="VacayTracker" style="height: 64px; width: 64px; display: block; margin: 0 auto 16px; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; font-size: 18px; font-weight: 600; color: #0D83A2;">
                            <h1 style="margin: 0; color: #00384F; font-size: 24px; font-weight: 600; letter-spacing: -0.5px;">Request Update</h1>
                        </td>
                    </tr>
                    <!-- Status Bar (Red for Rejected) -->
                    <tr>
                        <td style="padding: 0; height: 4px; background: linear-gradient(90deg, #ef4444 0%, #f87171 100%); background-color: #ef4444;" bgcolor="#ef4444"></td>
                    </tr>
                    <!-- Content -->
                    <tr>
                        <td style="padding: 32px 40px;">
                            <p style="margin: 0 0 20px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Hi <strong style="color: #00384F;">{{.UserName}}</strong>,
                            </p>
                            <p style="margin: 0 0 24px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Unfortunately, your vacation request could not be approved at this time.
                            </p>
                            <!-- Details Box -->
                            <div style="background-color: #f8fafc; border: 1px solid #e2e8f0; border-radius: 12px; padding: 20px; margin: 0 0 24px;">
                                <div style="display: inline-block; padding: 4px 12px; background-color: #fef2f2; color: #991b1b; font-size: 12px; font-weight: 600; border-radius: 20px; margin-bottom: 12px;">Not Approved</div>
                                <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px;">Start Date</td>
                                        <td style="padding: 8px 0; color: #00384F; font-size: 14px; font-weight: 600; text-align: right;">{{.StartDate}}</td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px;">End Date</td>
                                        <td style="padding: 8px 0; color: #00384F; font-size: 14px; font-weight: 600; text-align: right;">{{.EndDate}}</td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px;">Total Days</td>
                                        <td style="padding: 8px 0; color: #00384F; font-size: 14px; font-weight: 600; text-align: right;">{{.TotalDays}}</td>
                                    </tr>
                                </table>
                            </div>
                            {{if .Reason}}
                            <!-- Reason Box -->
                            <div style="background-color: #f9fafb; border-radius: 12px; padding: 16px 20px; margin: 0 0 24px;">
                                <p style="margin: 0 0 8px; color: #0D83A2; font-size: 14px; font-weight: 600;">Reason</p>
                                <p style="margin: 0; color: #374151; font-size: 14px; line-height: 1.5;">{{.Reason}}</p>
                            </div>
                            {{end}}
                            <p style="margin: 0 0 28px; color: #6b7280; font-size: 14px; line-height: 1.6;">
                                Please contact your manager if you have questions or would like to submit a new request for different dates.
                            </p>
                            <!-- CTA Button -->
                            <div style="text-align: center;">
                                <a href="{{.AppURL}}/employee" style="display: inline-block; padding: 14px 32px; background-color: #0D83A2; color: #ffffff; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px; box-shadow: 0 2px 8px rgba(13, 131, 162, 0.25);">View Dashboard</a>
                            </div>
                        </td>
                    </tr>
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 24px 40px; background-color: #e6f7f9; border-radius: 0 0 16px 16px; text-align: center; border-top: 1px solid #cceff3;" class="email-footer">
                            <p style="margin: 0 0 4px; color: #0a6a84; font-size: 13px; font-weight: 500;" class="text-heading">VacayTracker</p>
                            <p style="margin: 0; color: #6b7280; font-size: 12px;" class="text-secondary">Your vacation tracking companion</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`

const requestRejectedText = `Hi {{.UserName}},

Unfortunately, your vacation request could not be approved at this time.

Request Details:
- Start Date: {{.StartDate}}
- End Date: {{.EndDate}}
- Total Days: {{.TotalDays}}
{{if .Reason}}
Reason: {{.Reason}}
{{end}}
Please contact your manager if you have questions or would like to submit a new request for different dates.

View your dashboard at: {{.AppURL}}/employee

---
VacayTracker - Your vacation tracking companion`

// Admin notification email templates
const adminNewRequestSubject = "New Vacation Request Pending"

const adminNewRequestHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>New Vacation Request</title>
    <!--[if mso]>
    <noscript>
        <xml>
            <o:OfficeDocumentSettings>
                <o:PixelsPerInch>96</o:PixelsPerInch>
            </o:OfficeDocumentSettings>
        </xml>
    </noscript>
    <![endif]-->
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #e6f7f9;">
    <!-- Preheader text (shows in inbox preview) -->
    <div style="display: none; max-height: 0; overflow: hidden; mso-hide: all;">
        {{.RequesterName}} has submitted a vacation request for {{.TotalDays}} days. Review required.
        &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847; &#847;
    </div>
    <table role="presentation" style="width: 100%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 20px;">
                <table role="presentation" style="width: 600px; max-width: 100%; border-collapse: collapse; background-color: #ffffff; border-radius: 16px; box-shadow: 0 4px 20px rgba(13, 131, 162, 0.08);">
                    <!-- Header with Logo -->
                    <tr>
                        <td style="padding: 32px 40px 24px; text-align: center;">
                            <img src="{{.AppURL}}/logo.png" width="64" height="64" alt="VacayTracker" style="height: 64px; width: 64px; display: block; margin: 0 auto 16px; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; font-size: 18px; font-weight: 600; color: #0D83A2;">
                            <h1 style="margin: 0; color: #00384F; font-size: 24px; font-weight: 600; letter-spacing: -0.5px;">New Request</h1>
                        </td>
                    </tr>
                    <!-- Status Bar (Purple for Admin) -->
                    <tr>
                        <td style="padding: 0; height: 4px; background: linear-gradient(90deg, #8b5cf6 0%, #a78bfa 100%); background-color: #8b5cf6;" bgcolor="#8b5cf6"></td>
                    </tr>
                    <!-- Content -->
                    <tr>
                        <td style="padding: 32px 40px;">
                            <p style="margin: 0 0 24px; color: #374151; font-size: 16px; line-height: 1.6;">
                                A new vacation request requires your attention.
                            </p>
                            <!-- Details Box -->
                            <div style="background-color: #f8fafc; border: 1px solid #e2e8f0; border-radius: 12px; padding: 20px; margin: 0 0 24px;">
                                <div style="display: inline-block; padding: 4px 12px; background-color: #f3f0ff; color: #5b21b6; font-size: 12px; font-weight: 600; border-radius: 20px; margin-bottom: 12px;">Action Required</div>
                                <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px;">Employee</td>
                                        <td style="padding: 8px 0; color: #00384F; font-size: 14px; font-weight: 600; text-align: right;">{{.RequesterName}}</td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px;">Start Date</td>
                                        <td style="padding: 8px 0; color: #00384F; font-size: 14px; font-weight: 600; text-align: right;">{{.StartDate}}</td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px;">End Date</td>
                                        <td style="padding: 8px 0; color: #00384F; font-size: 14px; font-weight: 600; text-align: right;">{{.EndDate}}</td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px;">Total Days</td>
                                        <td style="padding: 8px 0; color: #00384F; font-size: 14px; font-weight: 600; text-align: right;">{{.TotalDays}}</td>
                                    </tr>
                                </table>
                                {{if .RequestReason}}
                                <div style="margin-top: 16px; padding-top: 16px; border-top: 1px solid #e2e8f0;">
                                    <p style="margin: 0 0 4px; color: #6b7280; font-size: 14px;">Reason</p>
                                    <p style="margin: 0; color: #374151; font-size: 14px;">{{.RequestReason}}</p>
                                </div>
                                {{end}}
                            </div>
                            <!-- CTA Button -->
                            <div style="text-align: center;">
                                <a href="{{.AppURL}}/admin" style="display: inline-block; padding: 14px 32px; background-color: #0D83A2; color: #ffffff; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px; box-shadow: 0 2px 8px rgba(13, 131, 162, 0.25);">Review Request</a>
                            </div>
                        </td>
                    </tr>
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 24px 40px; background-color: #e6f7f9; border-radius: 0 0 16px 16px; text-align: center; border-top: 1px solid #cceff3;">
                            <p style="margin: 0 0 4px; color: #0a6a84; font-size: 13px; font-weight: 500;">VacayTracker</p>
                            <p style="margin: 0; color: #6b7280; font-size: 12px;">Admin Notification</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`

const adminNewRequestText = `New Vacation Request Pending

A new vacation request requires your attention.

Request Details:
- Employee: {{.RequesterName}}
- Start Date: {{.StartDate}}
- End Date: {{.EndDate}}
- Total Days: {{.TotalDays}}
{{if .RequestReason}}- Reason: {{.RequestReason}}{{end}}

Review this request at: {{.AppURL}}/admin

---
VacayTracker - Admin Notification`
