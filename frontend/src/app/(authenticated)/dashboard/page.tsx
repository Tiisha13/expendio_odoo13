import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { auth } from "@/auth";
import { redirect } from "next/navigation";
import { expenseAPI } from "@/lib/api";
import { Receipt, DollarSign, CheckCircle, Clock } from "lucide-react";
import Link from "next/link";

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
  const company = user?.company;

  return (
    <div className="space-y-6 p-4">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Welcome back, {user.first_name}!</h1>
        <p className="text-muted-foreground">
          {company?.name || "Your Company"} • {user.role}
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Expenses</CardTitle>
            <Receipt className="text-muted-foreground h-4 w-4" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.total}</div>
            <p className="text-muted-foreground text-xs">All time expenses</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Amount</CardTitle>
            <DollarSign className="text-muted-foreground h-4 w-4" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {stats.totalAmount.toFixed(2)} {company.base_currency}
            </div>
            <p className="text-muted-foreground text-xs">In {company.base_currency}</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Pending</CardTitle>
            <Clock className="text-muted-foreground h-4 w-4" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.pending}</div>
            <p className="text-muted-foreground text-xs">Awaiting approval</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Approved</CardTitle>
            <CheckCircle className="text-muted-foreground h-4 w-4" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.approved}</div>
            <p className="text-muted-foreground text-xs">Successfully approved</p>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Quick Actions</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-4 md:grid-cols-3">
          {user.role !== "employee" && (
            <Link
              href="/users"
              className="hover:bg-accent flex flex-col items-center justify-center rounded-lg border p-6 transition-colors"
            >
              <div className="text-2xl font-semibold">👥</div>
              <div className="mt-2 font-medium">Manage Users</div>
              <div className="text-muted-foreground text-sm">View and manage team</div>
            </Link>
          )}

          <Link
            href="/expenses"
            className="hover:bg-accent flex flex-col items-center justify-center rounded-lg border p-6 transition-colors"
          >
            <div className="text-2xl font-semibold">💰</div>
            <div className="mt-2 font-medium">My Expenses</div>
            <div className="text-muted-foreground text-sm">View and create expenses</div>
          </Link>

          {(user.role === "manager" || user.role === "admin") && (
            <Link
              href="/approvals"
              className="hover:bg-accent flex flex-col items-center justify-between rounded-lg border p-6 transition-colors"
            >
              <div className="text-2xl font-semibold">✅</div>
              <div className="mt-2 font-medium">Approvals</div>
              <div className="text-muted-foreground text-sm">Review pending expenses</div>
            </Link>
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
