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
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to VacayTracker</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f5f4;">
    <table role="presentation" style="width: 100%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; max-width: 100%; border-collapse: collapse; background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 20px; text-align: center; background: linear-gradient(135deg, #0ea5e9 0%, #0284c7 100%); border-radius: 12px 12px 0 0;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: 600;">Welcome Aboard!</h1>
                        </td>
                    </tr>
                    <!-- Content -->
                    <tr>
                        <td style="padding: 40px;">
                            <p style="margin: 0 0 20px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Ahoy, <strong>{{.UserName}}</strong>!
                            </p>
                            <p style="margin: 0 0 20px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Your VacayTracker account has been created. You can now start planning your well-deserved time off!
                            </p>
                            <div style="background-color: #f0f9ff; border-radius: 8px; padding: 20px; margin: 20px 0;">
                                <p style="margin: 0 0 10px; color: #374151; font-size: 14px;"><strong>Your login credentials:</strong></p>
                                <p style="margin: 0 0 5px; color: #374151; font-size: 14px;">Email: <code style="background-color: #e0f2fe; padding: 2px 6px; border-radius: 4px;">{{.UserEmail}}</code></p>
                                <p style="margin: 0; color: #374151; font-size: 14px;">Temporary Password: <code style="background-color: #e0f2fe; padding: 2px 6px; border-radius: 4px;">{{.TempPassword}}</code></p>
                            </div>
                            <p style="margin: 20px 0; color: #dc2626; font-size: 14px; line-height: 1.6;">
                                Please change your password after your first login for security purposes.
                            </p>
                            <div style="text-align: center; margin: 30px 0;">
                                <a href="{{.AppURL}}" style="display: inline-block; padding: 14px 32px; background-color: #0ea5e9; color: #ffffff; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px;">Login to VacayTracker</a>
                            </div>
                        </td>
                    </tr>
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 20px 40px; background-color: #f9fafb; border-radius: 0 0 12px 12px; text-align: center;">
                            <p style="margin: 0; color: #6b7280; font-size: 12px;">
                                VacayTracker - Your vacation tracking companion
                            </p>
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
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Vacation Request Submitted</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f5f4;">
    <table role="presentation" style="width: 100%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; max-width: 100%; border-collapse: collapse; background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 20px; text-align: center; background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%); border-radius: 12px 12px 0 0;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: 600;">Request Submitted</h1>
                        </td>
                    </tr>
                    <!-- Content -->
                    <tr>
                        <td style="padding: 40px;">
                            <p style="margin: 0 0 20px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Hi <strong>{{.UserName}}</strong>,
                            </p>
                            <p style="margin: 0 0 20px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Your vacation request has been submitted and is pending approval.
                            </p>
                            <div style="background-color: #fffbeb; border-left: 4px solid #f59e0b; border-radius: 0 8px 8px 0; padding: 20px; margin: 20px 0;">
                                <p style="margin: 0 0 10px; color: #374151; font-size: 14px;"><strong>Request Details:</strong></p>
                                <p style="margin: 0 0 5px; color: #374151; font-size: 14px;">Start Date: <strong>{{.StartDate}}</strong></p>
                                <p style="margin: 0 0 5px; color: #374151; font-size: 14px;">End Date: <strong>{{.EndDate}}</strong></p>
                                <p style="margin: 0; color: #374151; font-size: 14px;">Total Days: <strong>{{.TotalDays}}</strong></p>
                            </div>
                            <p style="margin: 20px 0; color: #6b7280; font-size: 14px; line-height: 1.6;">
                                You'll receive another email once your request has been reviewed.
                            </p>
                            <div style="text-align: center; margin: 30px 0;">
                                <a href="{{.AppURL}}/employee" style="display: inline-block; padding: 14px 32px; background-color: #0ea5e9; color: #ffffff; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px;">View Dashboard</a>
                            </div>
                        </td>
                    </tr>
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 20px 40px; background-color: #f9fafb; border-radius: 0 0 12px 12px; text-align: center;">
                            <p style="margin: 0; color: #6b7280; font-size: 12px;">
                                VacayTracker - Your vacation tracking companion
                            </p>
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
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Vacation Request Approved</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f5f4;">
    <table role="presentation" style="width: 100%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; max-width: 100%; border-collapse: collapse; background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 20px; text-align: center; background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%); border-radius: 12px 12px 0 0;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: 600;">Request Approved!</h1>
                        </td>
                    </tr>
                    <!-- Content -->
                    <tr>
                        <td style="padding: 40px;">
                            <p style="margin: 0 0 20px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Great news, <strong>{{.UserName}}</strong>!
                            </p>
                            <p style="margin: 0 0 20px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Your vacation request has been approved. Time to start planning your getaway!
                            </p>
                            <div style="background-color: #f0fdf4; border-left: 4px solid #22c55e; border-radius: 0 8px 8px 0; padding: 20px; margin: 20px 0;">
                                <p style="margin: 0 0 10px; color: #374151; font-size: 14px;"><strong>Approved Vacation:</strong></p>
                                <p style="margin: 0 0 5px; color: #374151; font-size: 14px;">Start Date: <strong>{{.StartDate}}</strong></p>
                                <p style="margin: 0 0 5px; color: #374151; font-size: 14px;">End Date: <strong>{{.EndDate}}</strong></p>
                                <p style="margin: 0; color: #374151; font-size: 14px;">Total Days: <strong>{{.TotalDays}}</strong></p>
                            </div>
                            <div style="text-align: center; margin: 30px 0;">
                                <a href="{{.AppURL}}/employee" style="display: inline-block; padding: 14px 32px; background-color: #0ea5e9; color: #ffffff; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px;">View Dashboard</a>
                            </div>
                        </td>
                    </tr>
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 20px 40px; background-color: #f9fafb; border-radius: 0 0 12px 12px; text-align: center;">
                            <p style="margin: 0; color: #6b7280; font-size: 12px;">
                                VacayTracker - Your vacation tracking companion
                            </p>
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
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Vacation Request Update</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f5f4;">
    <table role="presentation" style="width: 100%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; max-width: 100%; border-collapse: collapse; background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 20px; text-align: center; background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%); border-radius: 12px 12px 0 0;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: 600;">Request Not Approved</h1>
                        </td>
                    </tr>
                    <!-- Content -->
                    <tr>
                        <td style="padding: 40px;">
                            <p style="margin: 0 0 20px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Hi <strong>{{.UserName}}</strong>,
                            </p>
                            <p style="margin: 0 0 20px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Unfortunately, your vacation request could not be approved at this time.
                            </p>
                            <div style="background-color: #fef2f2; border-left: 4px solid #ef4444; border-radius: 0 8px 8px 0; padding: 20px; margin: 20px 0;">
                                <p style="margin: 0 0 10px; color: #374151; font-size: 14px;"><strong>Request Details:</strong></p>
                                <p style="margin: 0 0 5px; color: #374151; font-size: 14px;">Start Date: <strong>{{.StartDate}}</strong></p>
                                <p style="margin: 0 0 5px; color: #374151; font-size: 14px;">End Date: <strong>{{.EndDate}}</strong></p>
                                <p style="margin: 0; color: #374151; font-size: 14px;">Total Days: <strong>{{.TotalDays}}</strong></p>
                            </div>
                            {{if .Reason}}
                            <div style="background-color: #f9fafb; border-radius: 8px; padding: 20px; margin: 20px 0;">
                                <p style="margin: 0 0 10px; color: #374151; font-size: 14px;"><strong>Reason:</strong></p>
                                <p style="margin: 0; color: #374151; font-size: 14px;">{{.Reason}}</p>
                            </div>
                            {{end}}
                            <p style="margin: 20px 0; color: #6b7280; font-size: 14px; line-height: 1.6;">
                                Please contact your manager if you have questions or would like to submit a new request for different dates.
                            </p>
                            <div style="text-align: center; margin: 30px 0;">
                                <a href="{{.AppURL}}/employee" style="display: inline-block; padding: 14px 32px; background-color: #0ea5e9; color: #ffffff; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px;">View Dashboard</a>
                            </div>
                        </td>
                    </tr>
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 20px 40px; background-color: #f9fafb; border-radius: 0 0 12px 12px; text-align: center;">
                            <p style="margin: 0; color: #6b7280; font-size: 12px;">
                                VacayTracker - Your vacation tracking companion
                            </p>
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
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>New Vacation Request</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f5f4;">
    <table role="presentation" style="width: 100%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; max-width: 100%; border-collapse: collapse; background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 20px; text-align: center; background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%); border-radius: 12px 12px 0 0;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: 600;">New Request</h1>
                        </td>
                    </tr>
                    <!-- Content -->
                    <tr>
                        <td style="padding: 40px;">
                            <p style="margin: 0 0 20px; color: #374151; font-size: 16px; line-height: 1.6;">
                                A new vacation request requires your attention.
                            </p>
                            <div style="background-color: #f5f3ff; border-left: 4px solid #8b5cf6; border-radius: 0 8px 8px 0; padding: 20px; margin: 20px 0;">
                                <p style="margin: 0 0 10px; color: #374151; font-size: 14px;"><strong>Request Details:</strong></p>
                                <p style="margin: 0 0 5px; color: #374151; font-size: 14px;">Employee: <strong>{{.RequesterName}}</strong></p>
                                <p style="margin: 0 0 5px; color: #374151; font-size: 14px;">Start Date: <strong>{{.StartDate}}</strong></p>
                                <p style="margin: 0 0 5px; color: #374151; font-size: 14px;">End Date: <strong>{{.EndDate}}</strong></p>
                                <p style="margin: 0; color: #374151; font-size: 14px;">Total Days: <strong>{{.TotalDays}}</strong></p>
                                {{if .RequestReason}}
                                <p style="margin: 10px 0 0; color: #374151; font-size: 14px;">Reason: {{.RequestReason}}</p>
                                {{end}}
                            </div>
                            <div style="text-align: center; margin: 30px 0;">
                                <a href="{{.AppURL}}/admin" style="display: inline-block; padding: 14px 32px; background-color: #0ea5e9; color: #ffffff; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px;">Review Request</a>
                            </div>
                        </td>
                    </tr>
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 20px 40px; background-color: #f9fafb; border-radius: 0 0 12px 12px; text-align: center;">
                            <p style="margin: 0; color: #6b7280; font-size: 12px;">
                                VacayTracker - Admin Notification
                            </p>
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
