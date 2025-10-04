import { apiClient, createClientAPI } from "@/lib/api-client";
import type {
  APIResponse,
  User,
  Expense,
  Approval,
  OCRResult,
  CreateUserInput,
  CreateExpenseInput,
  UpdateExpenseInput,
  ApprovalActionInput,
  UpdateUserRoleInput,
  AssignManagerInput,
} from "@/types/api";

// User API
export const userAPI = {
  list: () => apiClient<APIResponse<User[]>>("/users"),

  get: (id: string) => apiClient<APIResponse<User>>(`/users/${id}`),

  create: (data: CreateUserInput) =>
    apiClient<APIResponse<User>>("/users", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  updateRole: (id: string, data: UpdateUserRoleInput) =>
    apiClient<APIResponse<User>>(`/users/${id}/role`, {
      method: "PUT",
      body: JSON.stringify(data),
    }),

  assignManager: (id: string, data: AssignManagerInput) =>
    apiClient<APIResponse<User>>(`/users/${id}/manager`, {
      method: "PUT",
      body: JSON.stringify(data),
    }),

  delete: (id: string) =>
    apiClient<APIResponse<void>>(`/users/${id}`, {
      method: "DELETE",
    }),
};

// Expense API
export const expenseAPI = {
  list: (page = 1, limit = 10) =>
    apiClient<APIResponse<Expense[]>>(`/expenses?page=${page}&limit=${limit}`),

  get: (id: string) => apiClient<APIResponse<Expense>>(`/expenses/${id}`),

  pending: () => apiClient<APIResponse<Expense[]>>("/expenses/pending"),

  create: (data: CreateExpenseInput) =>
    apiClient<APIResponse<Expense>>("/expenses", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  update: (id: string, data: UpdateExpenseInput) =>
    apiClient<APIResponse<Expense>>(`/expenses/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    }),

  delete: (id: string) =>
    apiClient<APIResponse<void>>(`/expenses/${id}`, {
      method: "DELETE",
    }),
};

// Approval API
export const approvalAPI = {
  pending: () => apiClient<APIResponse<Approval[]>>("/approvals/pending"),

  history: (expenseId: string) =>
    apiClient<APIResponse<Approval[]>>(`/approvals/history/${expenseId}`),

  approve: (id: string, data: ApprovalActionInput) =>
    apiClient<APIResponse<Approval>>(`/approvals/${id}/approve`, {
      method: "POST",
      body: JSON.stringify(data),
    }),

  reject: (id: string, data: ApprovalActionInput) =>
    apiClient<APIResponse<Approval>>(`/approvals/${id}/reject`, {
      method: "POST",
      body: JSON.stringify(data),
    }),
};

// OCR API
export const ocrAPI = {
  upload: async (file: File, createExpense: boolean = false, accessToken: string) => {
    const formData = new FormData();
    formData.append("receipt", file);
    formData.append("create_expense", createExpense.toString());

    const api = createClientAPI(accessToken);
    return api<APIResponse<OCRResult & { expense_id?: string }>>("/ocr/upload", {
      method: "POST",
      headers: {} as any, // Remove Content-Type to let browser set it with boundary
      body: formData as any,
    });
  },
};

// Client-side versions (for use in client components)
export const createClientUserAPI = (accessToken: string) => {
  const api = createClientAPI(accessToken);
  return {
    list: () => api<APIResponse<User[]>>("/users"),
    create: (data: CreateUserInput) =>
      api<APIResponse<User>>("/users", {
        method: "POST",
        body: JSON.stringify(data),
      }),
    updateRole: (id: string, data: UpdateUserRoleInput) =>
      api<APIResponse<User>>(`/users/${id}/role`, {
        method: "PUT",
        body: JSON.stringify(data),
      }),
    assignManager: (id: string, data: AssignManagerInput) =>
      api<APIResponse<User>>(`/users/${id}/manager`, {
        method: "PUT",
        body: JSON.stringify(data),
      }),
    delete: (id: string) =>
      api<APIResponse<void>>(`/users/${id}`, {
        method: "DELETE",
      }),
  };
};

export const createClientExpenseAPI = (accessToken: string) => {
  const api = createClientAPI(accessToken);
  return {
    list: (page = 1, limit = 10) =>
      api<APIResponse<Expense[]>>(`/expenses?page=${page}&limit=${limit}`),
    create: (data: CreateExpenseInput) =>
      api<APIResponse<Expense>>("/expenses", {
        method: "POST",
        body: JSON.stringify(data),
      }),
    update: (id: string, data: UpdateExpenseInput) =>
      api<APIResponse<Expense>>(`/expenses/${id}`, {
        method: "PUT",
        body: JSON.stringify(data),
      }),
    delete: (id: string) =>
      api<APIResponse<void>>(`/expenses/${id}`, {
        method: "DELETE",
      }),
    pending: () => api<APIResponse<Expense[]>>("/expenses/pending"),
  };
};

export const createClientApprovalAPI = (accessToken: string) => {
  const api = createClientAPI(accessToken);
  return {
    pending: () => api<APIResponse<Approval[]>>("/approvals/pending"),
    approve: (id: string, data: ApprovalActionInput) =>
      api<APIResponse<Approval>>(`/approvals/${id}/approve`, {
        method: "POST",
        body: JSON.stringify(data),
      }),
    reject: (id: string, data: ApprovalActionInput) =>
      api<APIResponse<Approval>>(`/approvals/${id}/reject`, {
        method: "POST",
        body: JSON.stringify(data),
      }),
  };
};
