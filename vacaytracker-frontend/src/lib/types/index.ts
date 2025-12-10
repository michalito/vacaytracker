// User Types
export type Role = 'admin' | 'employee';

export interface EmailPreferences {
	vacationUpdates: boolean;
	weeklyDigest: boolean;
	teamNotifications: boolean;
}

export interface User {
	id: string;
	email: string;
	name: string;
	role: Role;
	vacationBalance: number;
	startDate?: string;
	emailPreferences: EmailPreferences;
	createdAt?: string;
	updatedAt?: string;
}

// Vacation Types
export type VacationStatus = 'pending' | 'approved' | 'rejected';

export interface VacationRequest {
	id: string;
	userId: string;
	userName?: string;
	userEmail?: string;
	startDate: string;
	endDate: string;
	totalDays: number;
	reason?: string;
	status: VacationStatus;
	reviewedBy?: string;
	reviewedAt?: string;
	rejectionReason?: string;
	createdAt: string;
	updatedAt: string;
}

export interface TeamVacation {
	id: string;
	userId: string;
	userName: string;
	startDate: string;
	endDate: string;
	totalDays: number;
}

// Settings Types
export interface WeekendPolicy {
	excludeWeekends: boolean;
	excludedDays: number[];
}

export interface NewsletterConfig {
	enabled: boolean;
	frequency: 'weekly' | 'monthly';
	dayOfMonth: number;
	lastSentAt?: string | null;
}

export interface Settings {
	id: string;
	weekendPolicy: WeekendPolicy;
	newsletter: NewsletterConfig;
	defaultVacationDays: number;
	vacationResetMonth: number;
	updatedAt: string;
}

// API Types
export interface ApiError {
	code: string;
	message: string;
	details?: Record<string, unknown>;
}

export interface PaginationInfo {
	page: number;
	limit: number;
	total: number;
	totalPages: number;
}

export interface LoginResponse {
	token: string;
	user: User;
}

export interface VacationListResponse {
	requests: VacationRequest[];
	total: number;
}

export interface TeamVacationResponse {
	vacations: TeamVacation[];
	month: number;
	year: number;
}

export interface UserListResponse {
	users: User[];
	pagination: PaginationInfo;
}

// Form Types
export interface CreateVacationForm {
	startDate: string;
	endDate: string;
	reason?: string;
}

export interface CreateUserForm {
	email: string;
	password: string;
	name: string;
	role: Role;
	vacationBalance?: number;
	startDate?: string;
}

export interface UpdateUserForm {
	email?: string;
	name?: string;
	role?: Role;
	vacationBalance?: number;
	startDate?: string;
}

// Newsletter Types
export interface NewsletterPreview {
	subject: string;
	htmlBody: string;
	textBody: string;
	recipients: string[];
	recipientCount: number;
}

export interface NewsletterSendResponse {
	success: boolean;
	recipientCount: number;
	message: string;
}
