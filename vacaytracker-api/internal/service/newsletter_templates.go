package service

// Newsletter email templates

const newsletterSubject = "VacayTracker Monthly Summary"

const newsletterHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>VacayTracker Monthly Summary</title>
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
        Your monthly vacation summary for {{.Period}} is here. See your team's stats and upcoming time off.
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
                            <h1 style="margin: 0 0 8px; color: #00384F; font-size: 24px; font-weight: 600; letter-spacing: -0.5px;">Monthly Summary</h1>
                            <p style="margin: 0; color: #0D83A2; font-size: 16px; font-weight: 500;">{{.Period}}</p>
                        </td>
                    </tr>
                    <!-- Status Bar (Ocean brand) -->
                    <tr>
                        <td style="padding: 0; height: 4px; background: linear-gradient(90deg, #0D83A2 0%, #18C8D3 100%); background-color: #0D83A2;" bgcolor="#0D83A2"></td>
                    </tr>
                    <!-- Content -->
                    <tr>
                        <td style="padding: 32px 40px;">
                            <p style="margin: 0 0 20px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Ahoy, <strong style="color: #00384F;">{{.RecipientName}}</strong>!
                            </p>
                            <p style="margin: 0 0 28px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Here's your monthly vacation summary from VacayTracker.
                            </p>

                            <!-- Statistics Section -->
                            <div style="background-color: #f8fafc; border: 1px solid #e2e8f0; border-radius: 12px; padding: 20px; margin: 0 0 24px;">
                                <h2 style="margin: 0 0 16px; color: #0D83A2; font-size: 16px; font-weight: 600;">Monthly Statistics</h2>
                                <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px;">Requests Submitted</td>
                                        <td style="padding: 8px 0; color: #00384F; font-size: 14px; font-weight: 600; text-align: right;">{{.Stats.TotalSubmitted}}</td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px;">Approved</td>
                                        <td style="padding: 8px 0; color: #166534; font-size: 14px; font-weight: 600; text-align: right;">{{.Stats.TotalApproved}}</td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px;">Rejected</td>
                                        <td style="padding: 8px 0; color: #991b1b; font-size: 14px; font-weight: 600; text-align: right;">{{.Stats.TotalRejected}}</td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 8px 0; color: #6b7280; font-size: 14px;">Pending</td>
                                        <td style="padding: 8px 0; color: #92400e; font-size: 14px; font-weight: 600; text-align: right;">{{.Stats.TotalPending}}</td>
                                    </tr>
                                    <tr style="border-top: 1px solid #e2e8f0;">
                                        <td style="padding: 12px 0 8px; color: #6b7280; font-size: 14px;">Total Days Used</td>
                                        <td style="padding: 12px 0 8px; color: #0D83A2; font-size: 14px; font-weight: 600; text-align: right;">{{.Stats.TotalDaysUsed}}</td>
                                    </tr>
                                </table>
                            </div>

                            {{if .HasUpcoming}}
                            <!-- Upcoming Vacations Section -->
                            <div style="background-color: #f8fafc; border: 1px solid #e2e8f0; border-radius: 12px; padding: 20px; margin: 0 0 24px;">
                                <div style="display: inline-block; padding: 4px 12px; background-color: #f0fdf4; color: #166534; font-size: 12px; font-weight: 600; border-radius: 20px; margin-bottom: 12px;">Upcoming</div>
                                <h2 style="margin: 0 0 12px; color: #00384F; font-size: 16px; font-weight: 600;">Team Vacations</h2>
                                <p style="margin: 0 0 16px; color: #6b7280; font-size: 14px;">Team members on vacation next month:</p>
                                <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                    {{range .UpcomingVacations}}
                                    <tr>
                                        <td style="padding: 12px 0; color: #374151; font-size: 14px; border-bottom: 1px solid #e2e8f0;">
                                            <strong style="color: #00384F;">{{.UserName}}</strong><br>
                                            <span style="color: #6b7280; font-size: 13px;">{{.StartDate}} - {{.EndDate}} ({{.TotalDays}} days)</span>
                                        </td>
                                    </tr>
                                    {{end}}
                                </table>
                            </div>
                            {{end}}

                            {{if .HasLowBalance}}
                            <!-- Low Balance Reminder Section -->
                            <div style="background-color: #f8fafc; border: 1px solid #e2e8f0; border-radius: 12px; padding: 20px; margin: 0 0 24px;">
                                <div style="display: inline-block; padding: 4px 12px; background-color: #fffbeb; color: #92400e; font-size: 12px; font-weight: 600; border-radius: 20px; margin-bottom: 12px;">Reminder</div>
                                <h2 style="margin: 0 0 12px; color: #00384F; font-size: 16px; font-weight: 600;">Low Balances</h2>
                                <p style="margin: 0 0 16px; color: #6b7280; font-size: 14px;">The following team members have low vacation balances:</p>
                                <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                    {{range .LowBalanceUsers}}
                                    <tr>
                                        <td style="padding: 8px 0; color: #00384F; font-size: 14px; font-weight: 500;">{{.UserName}}</td>
                                        <td style="padding: 8px 0; color: #92400e; font-size: 14px; font-weight: 600; text-align: right;">{{.RemainingDays}} days</td>
                                    </tr>
                                    {{end}}
                                </table>
                            </div>
                            {{end}}

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
                            <p style="margin: 0 0 8px; color: #6b7280; font-size: 12px;">Your vacation tracking companion</p>
                            <p style="margin: 0; color: #9ca3af; font-size: 11px;">
                                You're receiving this because you opted in to digest emails.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`

const newsletterText = `VacayTracker Monthly Summary - {{.Period}}

Ahoy, {{.RecipientName}}!

Here's your monthly vacation summary from VacayTracker.

=== MONTHLY STATISTICS ===
Requests Submitted: {{.Stats.TotalSubmitted}}
Approved: {{.Stats.TotalApproved}}
Rejected: {{.Stats.TotalRejected}}
Pending: {{.Stats.TotalPending}}
Total Days Used: {{.Stats.TotalDaysUsed}}

{{if .HasUpcoming}}
=== UPCOMING VACATIONS ===
Team members on vacation next month:
{{range .UpcomingVacations}}
- {{.UserName}}: {{.StartDate}} - {{.EndDate}} ({{.TotalDays}} days)
{{end}}
{{end}}

{{if .HasLowBalance}}
=== LOW BALANCE REMINDER ===
The following team members have low vacation balances:
{{range .LowBalanceUsers}}
- {{.UserName}}: {{.RemainingDays}} days remaining
{{end}}
{{end}}

View your dashboard at: {{.AppURL}}/employee

---
VacayTracker - Your vacation tracking companion
You're receiving this because you opted in to weekly digest emails.`
