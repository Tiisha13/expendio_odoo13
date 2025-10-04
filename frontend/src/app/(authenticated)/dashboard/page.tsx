import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { auth } from "@/auth";
import { redirect } from "next/navigation";
import { expenseAPI } from "@/lib/api";
import { Receipt, DollarSign, CheckCircle, Clock } from "lucide-react";

export default async function Page() {
  const session = await auth();

  if (!session) {
    redirect("/login");
  }

  // Fetch dashboard stats
  let stats = {
    total: 0,
    pending: 0,
    approved: 0,
    totalAmount: 0,
  };

  try {
    const expensesRes = await expenseAPI.list(1, 100);
    const expenses = expensesRes.data || []; // Ensure it's always an array

    stats.total = expenses.length;
    stats.pending = expenses.filter((e) => e.status === "pending").length;
    stats.approved = expenses.filter((e) => e.status === "approved").length;
    stats.totalAmount = expenses.reduce((sum, e) => sum + e.amount, 0);
  } catch (error) {
    console.error("Failed to fetch expenses:", error);
    // stats already has default values
  }

  const user = session.user;
  const company = user.company;

  return (
    <div className="space-y-6 p-4">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">
          Welcome back, {user.first_name}!
        </h1>
        <p className="text-muted-foreground">
          {company.name} • {user.role}
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total Expenses
            </CardTitle>
            <Receipt className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.total}</div>
            <p className="text-xs text-muted-foreground">All time expenses</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Amount</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {stats.totalAmount.toFixed(2)} {company.base_currency}
            </div>
            <p className="text-xs text-muted-foreground">
              In {company.base_currency}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Pending</CardTitle>
            <Clock className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.pending}</div>
            <p className="text-xs text-muted-foreground">Awaiting approval</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Approved</CardTitle>
            <CheckCircle className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.approved}</div>
            <p className="text-xs text-muted-foreground">
              Successfully approved
            </p>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Quick Actions</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-4 md:grid-cols-3">
          {user.role !== "employee" && (
            <a
              href="/users"
              className="flex flex-col items-center justify-center p-6 border rounded-lg hover:bg-accent transition-colors"
            >
              <div className="text-2xl font-semibold">👥</div>
              <div className="mt-2 font-medium">Manage Users</div>
              <div className="text-sm text-muted-foreground">
                View and manage team
              </div>
            </a>
          )}

          <a
            href="/expenses"
            className="flex flex-col items-center justify-center p-6 border rounded-lg hover:bg-accent transition-colors"
          >
            <div className="text-2xl font-semibold">💰</div>
            <div className="mt-2 font-medium">My Expenses</div>
            <div className="text-sm text-muted-foreground">
              View and create expenses
            </div>
          </a>

          {(user.role === "manager" || user.role === "admin") && (
            <a
              href="/approvals"
              className="flex flex-col items-center justify-between p-6 border rounded-lg hover:bg-accent transition-colors"
            >
              <div className="text-2xl font-semibold">✅</div>
              <div className="mt-2 font-medium">Approvals</div>
              <div className="text-sm text-muted-foreground">
                Review pending expenses
              </div>
            </a>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

export const metadata = {
  title: "Dashboard | Expensio",
  description: "Your personal dashboard to manage expenses and budgets.",
};
