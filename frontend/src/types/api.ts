// API Response wrapper
export interface APIResponse<T> {
  success: boolean;
  message?: string;
  data: T;
  meta?: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}

// User types
export type UserRole = "admin" | "manager" | "employee";

export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: UserRole;
  company_id: string;
  manager_id?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

// Company types
export interface Company {
  id: string;
  name: string;
  base_currency: string;
  country: string;
  admin_user_id: string;
  approval_rule_id?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

// Expense types
export type ExpenseStatus = "pending" | "approved" | "rejected";
export type ExpenseCategory =
  | "meals"
  | "travel"
  | "accommodation"
  | "entertainment"
  | "office_supplies"
  | "software"
  | "other";

export interface Expense {
  id: string;
  user_id: string;
  company_id: string;
  amount: number;
  currency: string;
  converted_amount: number;
  exchange_rate: number;
  category: ExpenseCategory;
  description: string;
  date: string; // Added for compatibility
  expense_date: string;
  receipt_url?: string;
  merchant?: string;
  status: ExpenseStatus;
  current_approval_level: number;
  created_at: string;
  updated_at: string;
  user?: User; // Added for populated responses
}

// Approval types
export type ApprovalStatus = "pending" | "approved" | "rejected";

export interface Approval {
  id: string;
  expense_id: string;
  approver_id: string;
  level: number;
  status: ApprovalStatus;
  comments?: string;
  approved_at?: string;
  created_at: string;
  updated_at: string;
  expense?: Expense;
  approver?: User;
}

// OCR types
export interface OCRResult {
  id: string;
  user_id: string;
  receipt_url: string;
  amount?: number;
  currency?: string;
  merchant?: string;
  date?: string;
  category?: ExpenseCategory;
  raw_text: string;
  confidence: number;
  processed_at: string;
  created_at: string;
}

// Form types
export interface CreateUserInput {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  role: UserRole;
}

export interface CreateExpenseInput {
  amount: number;
  currency: string;
  category: ExpenseCategory;
  description: string;
  date?: string; // Added for form convenience
  expense_date: string;
  receipt_url?: string;
  merchant?: string;
}

export interface UpdateExpenseInput {
  amount?: number;
  description?: string;
  merchant?: string;
}

export interface ApprovalActionInput {
  comments: string;
}

export interface UpdateUserRoleInput {
  role: UserRole;
}

export interface AssignManagerInput {
  manager_id: string;
}

// Dashboard stats
export interface DashboardStats {
  total_expenses: number;
  pending_expenses: number;
  approved_expenses: number;
  rejected_expenses: number;
  total_amount: number;
  pending_approvals: number;
}
