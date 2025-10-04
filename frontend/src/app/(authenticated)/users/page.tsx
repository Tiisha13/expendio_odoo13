"use client";

import { useEffect, useState } from "react";
import { useSession } from "next-auth/react";
import { User, CreateUserInput, UpdateUserRoleInput, AssignManagerInput } from "@/types/api";
import { createClientUserAPI } from "@/lib/api";
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
import { useToast } from "@/hooks/use-toast";
import { Plus } from "lucide-react";
import { DataTable } from "@/components/ui/data-table";
import { createUserColumns } from "@/components/users-columns";

export default function Page() {
  const { data: session } = useSession();
  const { toast } = useToast();
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [userToDelete, setUserToDelete] = useState<string | null>(null);

  const [formData, setFormData] = useState<CreateUserInput>({
    email: "",
    password: "",
    first_name: "",
    last_name: "",
    role: "employee",
  });

  const fetchUsers = async () => {
    if (!session?.accessToken) return;

    try {
      setLoading(true);
      const api = createClientUserAPI(session.accessToken, session.user?.refresh_token);
      const response = await api.list();
      setUsers(response.data || []); // Ensure it's always an array
    } catch (error: any) {
      setUsers([]); // Reset to empty array on error
      toast({
        title: "Error",
        description: error.message || "Failed to fetch users",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, [session]);

  const handleCreateUser = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!session?.accessToken) return;

    try {
      const api = createClientUserAPI(session.accessToken, session.user?.refresh_token);
      await api.create(formData);

      toast({
        title: "Success",
        description: "User created successfully",
      });

      setCreateDialogOpen(false);
      setFormData({
        email: "",
        password: "",
        first_name: "",
        last_name: "",
        role: "employee",
      });
      fetchUsers();
    } catch (error: any) {
      toast({
        title: "Error",
        description: error.message || "Failed to create user",
        variant: "destructive",
      });
    }
  };

  const handleUpdateRole = async (userId: string, role: string) => {
    if (!session?.accessToken) return;

    try {
      const api = createClientUserAPI(session.accessToken, session.user?.refresh_token);
      await api.updateRole(userId, { role } as UpdateUserRoleInput);

      toast({
        title: "Success",
        description: "User role updated successfully",
      });

      fetchUsers();
    } catch (error: any) {
      toast({
        title: "Error",
        description: error.message || "Failed to update role",
        variant: "destructive",
      });
    }
  };

  const handleAssignManager = async (userId: string, managerId: string) => {
    if (!session?.accessToken) return;

    try {
      const api = createClientUserAPI(session.accessToken, session.user?.refresh_token);
      await api.assignManager(userId, {
        manager_id: managerId,
      } as AssignManagerInput);

      toast({
        title: "Success",
        description: "Manager assigned successfully",
      });

      fetchUsers();
    } catch (error: any) {
      toast({
        title: "Error",
        description: error.message || "Failed to assign manager",
        variant: "destructive",
      });
    }
  };

  const handleDeleteUser = async (userId: string) => {
    if (!session?.accessToken) return;

    try {
      const api = createClientUserAPI(session.accessToken, session.user?.refresh_token);
      await api.delete(userId);

      toast({
        title: "Success",
        description: "User deleted successfully",
      });

      setDeleteDialogOpen(false);
      setUserToDelete(null);
      fetchUsers();
    } catch (error: any) {
      toast({
        title: "Error",
        description: error.message || "Failed to delete user",
        variant: "destructive",
      });
    }
  };

  const openDeleteDialog = (userId: string) => {
    setUserToDelete(userId);
    setDeleteDialogOpen(true);
  };

  const confirmDelete = () => {
    if (userToDelete) {
      handleDeleteUser(userToDelete);
    }
  };

  const managers = users.filter((u) => u.role === "manager" || u.role === "admin");

  if (loading) {
    return <div className="p-4">Loading...</div>;
  }

  return (
    <div className="space-y-6 p-4">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Users</h1>
          <p className="text-muted-foreground">Manage your team members and their roles</p>
        </div>
        <Dialog open={createDialogOpen} onOpenChange={setCreateDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="mr-2 h-4 w-4" />
              Add User
            </Button>
          </DialogTrigger>
          <DialogContent>
            <form onSubmit={handleCreateUser}>
              <DialogHeader>
                <DialogTitle>Create New User</DialogTitle>
                <DialogDescription>Add a new team member to your organization</DialogDescription>
              </DialogHeader>
              <div className="grid gap-4 py-4">
                <div className="grid gap-2">
                  <Label htmlFor="email">Email</Label>
                  <Input
                    id="email"
                    type="email"
                    value={formData.email}
                    onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                    required
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="password">Password</Label>
                  <Input
                    id="password"
                    type="password"
                    value={formData.password}
                    onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                    required
                  />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div className="grid gap-2">
                    <Label htmlFor="first_name">First Name</Label>
                    <Input
                      id="first_name"
                      value={formData.first_name}
                      onChange={(e) => setFormData({ ...formData, first_name: e.target.value })}
                      required
                    />
                  </div>
                  <div className="grid gap-2">
                    <Label htmlFor="last_name">Last Name</Label>
                    <Input
                      id="last_name"
                      value={formData.last_name}
                      onChange={(e) => setFormData({ ...formData, last_name: e.target.value })}
                      required
                    />
                  </div>
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="role">Role</Label>
                  <Select
                    value={formData.role}
                    onValueChange={(value) => setFormData({ ...formData, role: value as any })}
                  >
                    <SelectTrigger className="w-full">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="employee">Employee</SelectItem>
                      <SelectItem value="manager">Manager</SelectItem>
                      <SelectItem value="admin">Admin</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
              <DialogFooter>
                <Button type="submit">Create User</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <DataTable
        columns={createUserColumns(
          handleUpdateRole,
          handleAssignManager,
          openDeleteDialog,
          managers
        )}
        data={users}
        searchKey="email"
        searchPlaceholder="Search by email..."
      />

      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete User</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete this user? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => {
                setDeleteDialogOpen(false);
                setUserToDelete(null);
              }}
            >
              Cancel
            </Button>
            <Button variant="destructive" onClick={confirmDelete}>
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
