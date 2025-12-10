package service

// Newsletter email templates

const newsletterSubject = "VacayTracker Monthly Summary"

const newsletterHTML = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>VacayTracker Monthly Summary</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f5f4;">
    <table role="presentation" style="width: 100%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; max-width: 100%; border-collapse: collapse; background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 20px; text-align: center; background: linear-gradient(135deg, #0ea5e9 0%, #0284c7 100%); border-radius: 12px 12px 0 0;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: 600;">Monthly Summary</h1>
                            <p style="margin: 10px 0 0; color: #e0f2fe; font-size: 16px;">{{.Period}}</p>
                        </td>
                    </tr>
                    <!-- Content -->
                    <tr>
                        <td style="padding: 40px;">
                            <p style="margin: 0 0 20px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Ahoy, <strong>{{.RecipientName}}</strong>!
                            </p>
                            <p style="margin: 0 0 30px; color: #374151; font-size: 16px; line-height: 1.6;">
                                Here's your monthly vacation summary from VacayTracker.
                            </p>

                            <!-- Statistics Section -->
                            <div style="background-color: #f0f9ff; border-radius: 8px; padding: 20px; margin: 20px 0;">
                                <h2 style="margin: 0 0 15px; color: #0369a1; font-size: 18px; font-weight: 600;">Monthly Statistics</h2>
                                <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                    <tr>
                                        <td style="padding: 8px 0; color: #374151; font-size: 14px;">Requests Submitted:</td>
                                        <td style="padding: 8px 0; color: #374151; font-size: 14px; text-align: right;"><strong>{{.Stats.TotalSubmitted}}</strong></td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 8px 0; color: #374151; font-size: 14px;">Approved:</td>
                                        <td style="padding: 8px 0; color: #22c55e; font-size: 14px; text-align: right;"><strong>{{.Stats.TotalApproved}}</strong></td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 8px 0; color: #374151; font-size: 14px;">Rejected:</td>
                                        <td style="padding: 8px 0; color: #ef4444; font-size: 14px; text-align: right;"><strong>{{.Stats.TotalRejected}}</strong></td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 8px 0; color: #374151; font-size: 14px;">Pending:</td>
                                        <td style="padding: 8px 0; color: #f59e0b; font-size: 14px; text-align: right;"><strong>{{.Stats.TotalPending}}</strong></td>
                                    </tr>
                                    <tr style="border-top: 1px solid #e0f2fe;">
                                        <td style="padding: 12px 0 8px; color: #374151; font-size: 14px;">Total Days Used:</td>
                                        <td style="padding: 12px 0 8px; color: #0369a1; font-size: 14px; text-align: right;"><strong>{{.Stats.TotalDaysUsed}}</strong></td>
                                    </tr>
                                </table>
                            </div>

                            {{if .HasUpcoming}}
                            <!-- Upcoming Vacations Section -->
                            <div style="background-color: #f0fdf4; border-radius: 8px; padding: 20px; margin: 20px 0;">
                                <h2 style="margin: 0 0 15px; color: #166534; font-size: 18px; font-weight: 600;">Upcoming Vacations</h2>
                                <p style="margin: 0 0 15px; color: #374151; font-size: 14px;">Team members on vacation next month:</p>
                                <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                    {{range .UpcomingVacations}}
                                    <tr>
                                        <td style="padding: 10px 0; color: #374151; font-size: 14px; border-bottom: 1px solid #dcfce7;">
                                            <strong>{{.UserName}}</strong><br>
                                            <span style="color: #6b7280; font-size: 12px;">{{.StartDate}} - {{.EndDate}} ({{.TotalDays}} days)</span>
                                        </td>
                                    </tr>
                                    {{end}}
                                </table>
                            </div>
                            {{end}}

                            {{if .HasLowBalance}}
                            <!-- Low Balance Reminder Section -->
                            <div style="background-color: #fffbeb; border-left: 4px solid #f59e0b; border-radius: 0 8px 8px 0; padding: 20px; margin: 20px 0;">
                                <h2 style="margin: 0 0 15px; color: #b45309; font-size: 18px; font-weight: 600;">Balance Reminder</h2>
                                <p style="margin: 0 0 15px; color: #374151; font-size: 14px;">The following team members have low vacation balances:</p>
                                <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                    {{range .LowBalanceUsers}}
                                    <tr>
                                        <td style="padding: 8px 0; color: #374151; font-size: 14px;">{{.UserName}}</td>
                                        <td style="padding: 8px 0; color: #b45309; font-size: 14px; text-align: right;"><strong>{{.RemainingDays}} days remaining</strong></td>
                                    </tr>
                                    {{end}}
                                </table>
                            </div>
                            {{end}}

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
                            <p style="margin: 8px 0 0; color: #9ca3af; font-size: 11px;">
                                You're receiving this because you opted in to weekly digest emails.
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
