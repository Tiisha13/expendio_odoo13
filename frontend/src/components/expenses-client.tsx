"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { Expense, CreateExpenseInput } from "@/types/api";
import { createClientExpenseAPI, ocrAPI } from "@/lib/api";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useToast } from "@/hooks/use-toast";
import { Plus, Upload, Trash2 } from "lucide-react";
import { format } from "date-fns";

export default function ExpensesClient() {
  const { data: session } = useSession();
  const { toast } = useToast();
  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [loading, setLoading] = useState(true);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [ocrDialogOpen, setOcrDialogOpen] = useState(false);
  const [ocrLoading, setOcrLoading] = useState(false);

  const [formData, setFormData] = useState<CreateExpenseInput>({
    amount: 0,
    currency: session?.user?.company?.base_currency || "USD",
    category: "meals",
    description: "",
    expense_date: new Date().toISOString().split("T")[0],
    receipt_url: "",
  });

  const fetchExpenses = async () => {
    if (!session?.accessToken) return;

    try {
      setLoading(true);
      const api = createClientExpenseAPI(session.accessToken);
      const response = await api.list(1, 100);
      setExpenses(response.data || []); // Ensure it's always an array
    } catch (error: any) {
      setExpenses([]); // Reset to empty array on error
      toast({
        title: "Error",
        description: error.message || "Failed to fetch expenses",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchExpenses();
  }, [session]);

  const handleCreateExpense = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!session?.accessToken) return;

    try {
      const api = createClientExpenseAPI(session.accessToken);
      await api.create(formData);

      toast({
        title: "Success",
        description: "Expense created successfully",
      });

      setCreateDialogOpen(false);
      setFormData({
        amount: 0,
        currency: session?.user?.company?.base_currency || "USD",
        category: "meals",
        description: "",
        expense_date: new Date().toISOString().split("T")[0],
        receipt_url: "",
      });
      fetchExpenses();
    } catch (error: any) {
      toast({
        title: "Error",
        description: error.message || "Failed to create expense",
        variant: "destructive",
      });
    }
  };

  const handleDeleteExpense = async (expenseId: string) => {
    if (!session?.accessToken) return;
    if (!confirm("Are you sure you want to delete this expense?")) return;

    try {
      const api = createClientExpenseAPI(session.accessToken);
      await api.delete(expenseId);

      toast({
        title: "Success",
        description: "Expense deleted successfully",
      });

      fetchExpenses();
    } catch (error: any) {
      toast({
        title: "Error",
        description: error.message || "Failed to delete expense",
        variant: "destructive",
      });
    }
  };

  const handleOCRUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!e.target.files?.[0] || !session?.accessToken) return;

    const file = e.target.files[0];
    setOcrLoading(true);

    try {
      const response = await ocrAPI.upload(file, true, session.accessToken);

      toast({
        title: "Success",
        description: "Receipt processed and expense created!",
      });

      setOcrDialogOpen(false);
      fetchExpenses();
    } catch (error: any) {
      toast({
        title: "Error",
        description: error.message || "Failed to process receipt",
        variant: "destructive",
      });
    } finally {
      setOcrLoading(false);
    }
  };

  const getStatusBadge = (status: string) => {
    const variants: Record<
      string,
      "default" | "secondary" | "destructive" | "outline"
    > = {
      pending: "secondary",
      approved: "default",
      rejected: "destructive",
    };
    return <Badge variant={variants[status] || "outline"}>{status}</Badge>;
  };

  if (loading) {
    return <div className="p-4">Loading...</div>;
  }

  return (
    <div className="space-y-6 p-4">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Expenses</h1>
          <p className="text-muted-foreground">View and manage your expenses</p>
        </div>
        <div className="flex gap-2">
          <Dialog open={ocrDialogOpen} onOpenChange={setOcrDialogOpen}>
            <DialogTrigger asChild>
              <Button variant="outline">
                <Upload className="mr-2 h-4 w-4" />
                Upload Receipt
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Upload Receipt</DialogTitle>
                <DialogDescription>
                  Upload a receipt image to automatically create an expense
                </DialogDescription>
              </DialogHeader>
              <div className="grid gap-4 py-4">
                <div className="grid gap-2">
                  <Label htmlFor="receipt">Receipt Image</Label>
                  <Input
                    id="receipt"
                    type="file"
                    accept="image/*"
                    onChange={handleOCRUpload}
                    disabled={ocrLoading}
                  />
                  {ocrLoading && (
                    <p className="text-sm text-muted-foreground">
                      Processing receipt...
                    </p>
                  )}
                </div>
              </div>
            </DialogContent>
          </Dialog>

          <Dialog open={createDialogOpen} onOpenChange={setCreateDialogOpen}>
            <DialogTrigger asChild>
              <Button>
                <Plus className="mr-2 h-4 w-4" />
                Add Expense
              </Button>
            </DialogTrigger>
            <DialogContent className="max-h-[90vh] overflow-y-auto">
              <form onSubmit={handleCreateExpense}>
                <DialogHeader>
                  <DialogTitle>Create New Expense</DialogTitle>
                  <DialogDescription>
                    Add a new expense to your records
                  </DialogDescription>
                </DialogHeader>
                <div className="grid gap-4 py-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="grid gap-2">
                      <Label htmlFor="amount">Amount</Label>
                      <Input
                        id="amount"
                        type="number"
                        step="0.01"
                        value={formData.amount}
                        onChange={(e) =>
                          setFormData({
                            ...formData,
                            amount: parseFloat(e.target.value),
                          })
                        }
                        required
                      />
                    </div>
                    <div className="grid gap-2">
                      <Label htmlFor="currency">Currency</Label>
                      <Input
                        id="currency"
                        value={formData.currency}
                        onChange={(e) =>
                          setFormData({ ...formData, currency: e.target.value })
                        }
                        placeholder="USD"
                        required
                      />
                    </div>
                  </div>
                  <div className="grid gap-2">
                    <Label htmlFor="category">Category</Label>
                    <Select
                      value={formData.category}
                      onValueChange={(value) =>
                        setFormData({ ...formData, category: value as any })
                      }
                    >
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="meals">Meals</SelectItem>
                        <SelectItem value="travel">Travel</SelectItem>
                        <SelectItem value="accommodation">
                          Accommodation
                        </SelectItem>
                        <SelectItem value="entertainment">
                          Entertainment
                        </SelectItem>
                        <SelectItem value="office_supplies">
                          Office Supplies
                        </SelectItem>
                        <SelectItem value="software">Software</SelectItem>
                        <SelectItem value="other">Other</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="grid gap-2">
                    <Label htmlFor="date">Date</Label>
                    <Input
                      id="date"
                      type="date"
                      value={formData.expense_date}
                      onChange={(e) =>
                        setFormData({
                          ...formData,
                          expense_date: e.target.value,
                        })
                      }
                      required
                    />
                  </div>
                  <div className="grid gap-2">
                    <Label htmlFor="description">Description</Label>
                    <Textarea
                      id="description"
                      value={formData.description}
                      onChange={(e) =>
                        setFormData({
                          ...formData,
                          description: e.target.value,
                        })
                      }
                      rows={3}
                    />
                  </div>
                  <div className="grid gap-2">
                    <Label htmlFor="receipt_url">Receipt URL (optional)</Label>
                    <Input
                      id="receipt_url"
                      value={formData.receipt_url}
                      onChange={(e) =>
                        setFormData({
                          ...formData,
                          receipt_url: e.target.value,
                        })
                      }
                      placeholder="https://..."
                    />
                  </div>
                </div>
                <DialogFooter>
                  <Button type="submit">Create Expense</Button>
                </DialogFooter>
              </form>
            </DialogContent>
          </Dialog>
        </div>
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Date</TableHead>
              <TableHead>Category</TableHead>
              <TableHead>Description</TableHead>
              <TableHead>Amount</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {expenses.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={6}
                  className="text-center text-muted-foreground"
                >
                  No expenses found. Create your first expense!
                </TableCell>
              </TableRow>
            ) : (
              expenses.map((expense) => (
                <TableRow key={expense.id}>
                  <TableCell>
                    {format(new Date(expense.expense_date), "MMM d, yyyy")}
                  </TableCell>
                  <TableCell className="capitalize">
                    {expense.category}
                  </TableCell>
                  <TableCell className="max-w-xs truncate">
                    {expense.description}
                  </TableCell>
                  <TableCell>
                    {expense.amount.toFixed(2)} {expense.currency}
                  </TableCell>
                  <TableCell>{getStatusBadge(expense.status)}</TableCell>
                  <TableCell className="text-right">
                    {expense.status === "pending" && (
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleDeleteExpense(expense.id)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    )}
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
